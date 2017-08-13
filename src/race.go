package main

import (
	"strconv"
	"time"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/Zamiell/isaac-racing-server/src/models"
)

/*
	Race object methods
*/

// Get the place that someone would be if they finished the race right now
func (race *Race) GetCurrentPlace() int {
	currentPlace := 0
	for _, racer := range race.Racers {
		if racer.Place > currentPlace {
			currentPlace = racer.Place
		}
	}

	return currentPlace + 1
}

// Check to see if a race is ready to start, and if so, start it
// (called from the "websocketRaceReady" and "websocketRaceLeave" functions)
func (race *Race) CheckStart() {
	// Check to see if there is only 1 person in the race
	if len(race.Racers) == 1 && !race.Ruleset.Solo {
		return
	}

	// Check if everyone is ready
	for _, racer := range race.Racers {
		if racer.Status != "ready" {
			return
		}
	}

	race.Start()
}

func (race *Race) SetStatus(status string) {
	race.Status = status

	for _, s := range websocketSessions {
		type RaceSetStatusMessage struct {
			ID     int    `json:"id"`
			Status string `json:"status"`
		}
		websocketEmit(s, "raceSetStatus", &RaceSetStatusMessage{
			race.ID,
			race.Status,
		})
	}
}

func (race *Race) SetRacerStatus(username string, status string) {
	racer := race.Racers[username]
	racer.Status = status

	for racerName := range race.Racers {
		// Not all racers may be online during a race
		if s, ok := websocketSessions[racerName]; ok {
			type RacerSetStatusMessage struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				Status string `json:"status"`
				Place  int    `json:"place"`
			}
			websocketEmit(s, "racerSetStatus", &RacerSetStatusMessage{
				race.ID,
				username,
				status,
				racer.Place,
			})
		}
	}
}

// Recalculate everyone's mid-race places
func (race *Race) SetAllPlaceMid() {
	// Get the place that someone would be if they finished the race right now
	currentPlace := race.GetCurrentPlace()

	for _, racer := range race.Racers {
		if racer.Status != "racing" {
			// We don't need to calculate the mid-race place of someone who already finished or quit
			continue
		}

		racer.PlaceMid = currentPlace
		for _, racer2 := range race.Racers {
			if racer2.Status != "racing" {
				// We don't count people who finished or quit since our starting point was on "currentPlace"
				continue
			}
			if racer2.FloorNum > racer.FloorNum {
				racer.PlaceMid++
			} else if racer2.FloorNum == racer.FloorNum &&
				racer2.FloorNum > 8 &&
				racer2.StageType < racer.StageType {

				// This is custom logic for the "Everything" race goal
				// Sheol is StageType 0 and the Dark Room is StageType 0
				// Those are considered ahead of Cathedral and The Chest
				racer.PlaceMid++
			} else if racer2.FloorNum == racer.FloorNum &&
				racer2.StageType == racer.StageType &&
				racer2.DatetimeArrivedFloor < racer.DatetimeArrivedFloor {

				racer.PlaceMid++
			}
		}
	}
}

// Called from the "CheckStart" function
func (race *Race) Start() {
	var secondsToWait time.Duration
	if race.Ruleset.Solo {
		secondsToWait = 3
	} else {
		secondsToWait = 10
	}

	// Log the race starting
	log.Info("Race #"+strconv.Itoa(race.ID)+" starting in", secondsToWait, "seconds.")

	// Change the status for this race to "starting"
	race.SetStatus("starting")

	// Get the time X seconds in the future
	startTime := time.Now().Add(secondsToWait*time.Second).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))

	// Send everyone in the race a message specifying exactly when it will start
	for racerName := range race.Racers {
		// A racer might go offline the moment before it starts, so check just in case
		if s, ok := websocketSessions[racerName]; ok {
			websocketEmit(s, "raceStart", &RaceStartMessage{
				race.ID,
				startTime,
			})
		}
	}

	// Make the Twitch bot announce that the race is starting in 10 seconds
	twitchRaceStart(race)

	// Return for now and do more things in 10 seconds
	go race.Start2()
}

func (race *Race) Start2() {
	// Sleep 3 or 10 seconds
	var sleepTime time.Duration
	if race.Ruleset.Solo {
		sleepTime = 3
	} else {
		sleepTime = 10
	}
	time.Sleep(sleepTime * time.Second)

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()
	defer commandMutex.Unlock()

	// Log the race starting
	log.Info("Race", race.ID, "started with", len(race.Racers), "participants.")

	race.SetStatus("in progress")
	race.DatetimeStarted = getTimestamp()
	for _, racer := range race.Racers {
		racer.Status = "racing"
	}

	// Return for now and do more things in 30 minutes
	go race.Start3()
}

func (race *Race) Start3() {
	// Sleep 30 minutes
	time.Sleep(30 * time.Minute)

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()
	defer commandMutex.Unlock()

	// Find out if the race is finished
	// (we remove finished races from the "races" map)
	if _, ok := races[race.ID]; !ok {
		return
	}

	// Force the remaining racers to quit
	for _, racer := range race.Racers {
		if racer.Status != "racing" {
			continue
		}

		d := &IncomingWebsocketData{}
		d.Command = "race.Start3"
		d.ID = race.ID
		d.v.Username = racer.Name
		websocketRaceQuit(nil, d)
	}
}

func (race *Race) CheckFinish() {
	for _, racer := range race.Racers {
		if racer.Status == "racing" {
			return
		}
	}

	race.Finish()
}

func (race *Race) Finish() {
	// Log the race finishing
	log.Info("Race " + strconv.Itoa(race.ID) + " finished.")

	// Let everyone know it ended
	race.SetStatus("finished")

	// Remove it from the map
	delete(races, race.ID)

	// Write it to the database
	databaseRace := &models.Race{
		ID:              race.ID,
		Name:            race.Name,
		Ranked:          race.Ruleset.Ranked,
		Solo:            race.Ruleset.Solo,
		Format:          race.Ruleset.Format,
		Character:       race.Ruleset.Character,
		Goal:            race.Ruleset.Goal,
		StartingBuild:   race.Ruleset.StartingBuild,
		Seed:            race.Ruleset.Seed,
		Captain:         race.Captain,
		DatetimeStarted: race.DatetimeStarted,
	}
	if err := db.Races.Finish(databaseRace); err != nil {
		log.Error("Failed to write race #"+strconv.Itoa(race.ID)+" to the database:", err)
		return
	}

	for _, racer := range race.Racers {
		databaseRacer := &models.Racer{
			ID:               racer.ID,
			DatetimeJoined:   racer.DatetimeJoined,
			Seed:             racer.Seed,
			StartingItem:     racer.StartingItem,
			Place:            racer.Place,
			DatetimeFinished: racer.DatetimeFinished,
			RunTime:          racer.RunTime,
			Comment:          racer.Comment,
		}
		if err := db.RaceParticipants.Insert(race.ID, databaseRacer); err != nil {
			log.Error("Failed to write the RaceParticipants row for \""+race.Name+"\" to the database:", err)
			return
		}

		for _, item := range racer.Items {
			if err := db.RaceParticipantItems.Insert(
				racer.ID,
				race.ID,
				item.ID,
				item.FloorNum,
				item.StageType,
			); err != nil {
				log.Error("Failed to write the RaceParticipantItems row for \""+race.Name+"\" to the database:", err)
				return
			}
		}

		for _, room := range racer.Rooms {
			if err := db.RaceParticipantRooms.Insert(
				racer.ID,
				race.ID,
				room.ID,
				room.FloorNum,
				room.StageType,
			); err != nil {
				log.Error("Failed to write the RaceParticipantRooms row for \""+race.Name+"\" to the database:", err)
				return
			}
		}
	}
}

/*
	Race subroutines
*/

// Now that a user has finished, quit, or been disqualified from a race, update fields in the users table for unseeded races
func raceUpdateUnseededStats(raceID int, username string) {
	// Don't do anything if this is not an unseeded race (or an unranked race)
	// TODO

	// Get their unseeded stats
	/*
		if statsUnseeded, err := db.Users.GetStatsUnseeded(username); err != nil {
			log.Error("Database error:", err)
			return
		}
	*/

	// Update all the stats
	// TODO
}

// Now that the race has finished, update fields in the users table for seeded races
func raceUpdateSeededStats(raceID int, username string) {

}
