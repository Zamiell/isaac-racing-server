package main

import (
	"strconv"
	"time"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/Zamiell/isaac-racing-server/src/models"
)

/*
	Data structures
*/

// Used to track the current races in memory
type Race struct {
	ID              int
	Name            string
	Status          string /* open, starting, in progress, finished */
	Ruleset         Ruleset
	Captain         string
	SoundPlayed     bool
	DatetimeCreated int64
	DatetimeStarted int64
	Racers          map[string]*Racer
}
type Ruleset struct {
	Ranked              bool   `json:"ranked"`
	Solo                bool   `json:"solo"`
	Format              string `json:"format"`
	Character           string `json:"character"`
	CharacterRandom     bool   `json:"characterRandom"`
	Goal                string `json:"goal"`
	StartingBuild       int    `json:"startingBuild"`
	StartingBuildRandom bool   `json:"startingBuildRandom"`
	StartingItems       []int  `json:"startingItems"`
	Seed                string `json:"seed"`
}
type Racer struct {
	ID                   int
	Name                 string
	DatetimeJoined       int64
	Status               string /* not ready, ready, finished, quit, disqualified */
	Seed                 string
	FloorNum             int
	StageType            int
	DatetimeArrivedFloor int64
	Items                []*Item
	StartingItem         int
	Rooms                []*Room
	Place                int
	PlaceMid             int
	DatetimeFinished     int64
	RunTime              int64 /* in milliseconds */
	Comment              string
}
type Item struct {
	ID               int   `json:"id"`
	FloorNum         int   `json:"floorNum"`
	StageType        int   `json:"stageType"`
	DatetimeAcquired int64 `json:"datetimeAcquired"`
}
type Room struct {
	ID              string /* e.g. "5.999" */
	FloorNum        int
	StageType       int
	DatetimeArrived int64
}

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
	log.Info("Race "+strconv.Itoa(race.ID)+" starting in", int(secondsToWait), "seconds.")

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

	// Return for now and do more things later on when it is time to check to see if the race has been going for too long
	go race.Start3()
}

func (race *Race) Start3() {
	if race.Ruleset.Format == "custom" {
		time.Sleep(4 * time.Hour) // We need to make the timeout longer to accomodate multi-character speedrun races
	} else {
		time.Sleep(30 * time.Minute)
	}

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

		log.Info("Forcing racer \"" + racer.Name + "\" to quit since the race time limit has been reached.")

		d := &IncomingWebsocketData{
			Command: "race.Start3",
			ID:      race.ID,
			v: &models.SessionValues{
				Username: racer.Name,
			},
		}
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
				item.DatetimeAcquired,
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
				room.DatetimeArrived,
			); err != nil {
				log.Error("Failed to write the RaceParticipantRooms row for \""+race.Name+"\" to the database:", err)
				return
			}
		}
	}

	if race.Ruleset.Solo {
		if race.Ruleset.Ranked {
			if race.Ruleset.Format == "seeded" {
				leaderboardUpdateSoloSeeded(race)
			} else if race.Ruleset.Format == "unseeded" {
				leaderboardUpdateSoloUnseeded(race)
			}
		}
	} else {
		if race.Ruleset.Format == "seeded" ||
			race.Ruleset.Format == "seeded" ||
			race.Ruleset.Format == "diversity" {

			leaderboardUpdateTrueSkill(race)
		}
	}
}
