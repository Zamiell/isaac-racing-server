package main

import (
	"strconv"
	"time"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/Zamiell/isaac-racing-server/src/models"
	melody "gopkg.in/olahol/melody.v1"
)

/*
	Race subroutines
*/

func raceValidate(s *melody.Session, d *IncomingWebsocketData) bool {
	// Local variables
	raceID := d.ID
	username := d.v.Username

	// Validate that the requested race is sane
	if raceID <= 0 {
		log.Warning("User \"" + username + "\" attempted to call " + d.Command + " with a bogus ID of " + strconv.Itoa(raceID) + ".")
		websocketError(s, d.Command, "You must provide a valid race number.")
		return false
	}

	// Validate that the requested race exists
	if exists, err := db.Races.Exists(raceID); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return false
	} else if !exists {
		log.Info("User \"" + username + "\" attempted to call " + d.Command + " on race ID " + strconv.Itoa(raceID) + ", but it doesn't exist.")
		// Don't send an error to the user; this kind of thing can happen if their internet is lagging
		return false
	}

	// The user's request seems to be valid
	return true
}

func raceValidateStatus(s *melody.Session, d *IncomingWebsocketData, status string) bool {
	// Local variables
	raceID := d.ID
	username := d.v.Username

	// Validate that the race is set to the correct status
	if correctStatus, err := db.Races.CheckStatus(raceID, status); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return false
	} else if !correctStatus {
		log.Info("User \"" + username + "\" attempted to call " + d.Command + " on race ID " + strconv.Itoa(raceID) + ", but race is not set to status \"" + status + "\".")
		// Don't send an error to the user; this kind of thing can happen if their internet is lagging
		return false
	}

	// The race is the correct status
	return true
}

func raceValidateRuleset(s *melody.Session, d *IncomingWebsocketData) bool {
	// Local variables
	ruleset := d.Ruleset

	// Validate the ruleset type
	if ruleset.Type != "ranked" &&
		ruleset.Type != "unranked" {

		websocketError(s, d.Command, "That is not a valid type.")
		return false
	}

	// Validate the ruleset format
	if ruleset.Format != "unseeded" &&
		ruleset.Format != "seeded" &&
		ruleset.Format != "diversity" &&
		ruleset.Format != "unseeded-lite" &&
		ruleset.Format != "custom" {

		websocketError(s, d.Command, "That is not a valid ruleset.")
		return false
	}

	// Validate the character
	if ruleset.Character != "Isaac" &&
		ruleset.Character != "Magdalene" &&
		ruleset.Character != "Cain" &&
		ruleset.Character != "Judas" &&
		ruleset.Character != "Blue Baby" &&
		ruleset.Character != "Eve" &&
		ruleset.Character != "Samson" &&
		ruleset.Character != "Azazel" &&
		ruleset.Character != "Lazarus" &&
		ruleset.Character != "Eden" &&
		ruleset.Character != "The Lost" &&
		ruleset.Character != "Lilith" &&
		ruleset.Character != "Keeper" &&
		ruleset.Character != "Apollyon" &&
		ruleset.Character != "custom" {

		websocketError(s, d.Command, "That is not a valid character.")
		return false
	}

	// Validate the goal
	if ruleset.Goal != "Blue Baby" &&
		ruleset.Goal != "The Lamb" &&
		ruleset.Goal != "Mega Satan" &&
		ruleset.Goal != "Everything" &&
		ruleset.Goal != "custom" {

		websocketError(s, d.Command, "That is not a valid goal.")
		return false
	}

	// Validate the starting build
	if ruleset.Format != "seeded" && ruleset.StartingBuild != -1 {
		websocketError(s, d.Command, "You cannot set a starting build for a non-seeded race.")
		return false
	} else if ruleset.Format == "seeded" && (ruleset.StartingBuild < 1 || ruleset.StartingBuild > 33) { // There are 33 builds
		websocketError(s, d.Command, "That is not a valid starting build.")
		return false
	}

	// Validate unseeded ranked games
	if ruleset.Type == "ranked" && ruleset.Format == "unseeded" && ruleset.Character != "Judas" {
		websocketError(s, d.Command, "Ranked unseeded races must have a character of Judas.")
		return false
	}
	if ruleset.Type == "ranked" && ruleset.Format == "unseeded" && ruleset.Goal != "Blue Baby" {
		websocketError(s, d.Command, "Ranked unseeded races must have a goal of Blue Baby.")
		return false
	}

	return true
}

// Playing or observing
func raceValidateIn1(s *melody.Session, d *IncomingWebsocketData) bool {
	// Local variables
	raceID := d.ID
	userID := d.v.UserID
	username := d.v.Username

	// Validate that they are in the race
	if userInRace, err := db.RaceParticipants.CheckInRace1(userID, raceID); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return false
	} else if !userInRace {
		log.Warning("User \"" + username + "\" attempted to call " + d.Command + " on race ID " + strconv.Itoa(raceID) + ", but they are not in that race.")
		websocketError(s, d.Command, "You are not playing in or observing race ID "+strconv.Itoa(raceID)+".")
		return false
	}

	// The user is in the race
	return true
}

// ONLY playing (not observing)
func raceValidateIn2(s *melody.Session, d *IncomingWebsocketData) bool {
	// Local variables
	raceID := d.ID
	userID := d.v.UserID
	username := d.v.Username

	// Validate that they are in the race
	if userInRace, err := db.RaceParticipants.CheckInRace2(userID, raceID); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return false
	} else if !userInRace {
		log.Warning("User \"" + username + "\" attempted to call " + d.Command + " on race ID " + strconv.Itoa(raceID) + ", but they are not in that race.")
		websocketError(s, d.Command, "You are not playing in race ID "+strconv.Itoa(raceID)+".")
		return false
	}

	// The user is in the race
	return true
}

func raceValidateOut1(s *melody.Session, d *IncomingWebsocketData) bool {
	// Local variables
	raceID := d.ID
	userID := d.v.UserID
	username := d.v.Username

	// Validate that they are not already in the race
	if userInRace, err := db.RaceParticipants.CheckInRace1(userID, raceID); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return false
	} else if userInRace {
		log.Warning("User \"" + username + "\" attempted to call " + d.Command + " on race ID " + strconv.Itoa(raceID) + ", but they are already in that race.")
		websocketError(s, d.Command, "You are already playing or observing race ID "+strconv.Itoa(raceID)+".")
		return false
	}

	// The user is not in the race
	return true
}

func raceValidateOut2(s *melody.Session, d *IncomingWebsocketData) bool {
	// Local variables
	raceID := d.ID
	userID := d.v.UserID
	username := d.v.Username

	// Validate that they are not already in the race
	if userInRace, err := db.RaceParticipants.CheckInRace2(userID, raceID); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return false
	} else if userInRace {
		log.Warning("User \"" + username + "\" attempted to call " + d.Command + " on race ID " + strconv.Itoa(raceID) + ", but they are already in that race.")
		websocketError(s, d.Command, "You are already playing in race ID "+strconv.Itoa(raceID)+".")
		return false
	}

	// The user is not in the race
	return true
}

func racerValidateStatus(s *melody.Session, d *IncomingWebsocketData, status string) bool {
	// Local variables
	raceID := d.ID
	userID := d.v.UserID
	username := d.v.Username

	// Validate that the user is set to the correct status
	if correctStatus, err := db.RaceParticipants.CheckStatus(userID, raceID, status); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return false
	} else if !correctStatus {
		log.Info("User \"" + username + "\" attempted to call " + d.Command + " on race ID " + strconv.Itoa(raceID) + ", but they are not set to status \"" + status + "\".")
		// Don't send an error to the user; just silently fail
		// This type of thing can occur if, for example, they try to unready immediately before the race begins
		return false
	}

	// The user has the correct status
	return true
}

func racerSetStatus(s *melody.Session, d *IncomingWebsocketData, status string) bool {
	// Local variables
	raceID := d.ID
	username := d.v.Username

	// Change the status in the database
	if err := db.RaceParticipants.SetStatus(username, raceID, status); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return false
	}

	// The change was successful
	return true
}

// Recalculate everyones mid-race places
func racerSetAllPlaceMid(s *melody.Session, d *IncomingWebsocketData, racerList []models.Racer) bool {
	// Local variables
	raceID := d.ID

	// Get the current (final) place
	var currentPlace int
	for _, racer := range racerList {
		if racer.Place > currentPlace {
			currentPlace = racer.Place
		}
	}

	// Recalculate everyones mid-race places
	for _, racer := range racerList {
		if racer.Status != "racing" {
			continue // We don't need to calculate the mid-race place of someone who already finished or quit
		}
		racer.PlaceMid = currentPlace + 1
		for _, racer2 := range racerList {
			if racer2.Status != "racing" {
				continue // We don't count people who finished or quit since our starting point was on the currentPlace
			}
			if racer2.FloorNum > racer.FloorNum {
				racer.PlaceMid++
			} else if racer2.FloorNum == racer.FloorNum && racer2.StageType < racer.StageType {
				// This is custom logic for the "Everything" race goal
				// Sheol is StageType 0 and the Dark Room is StageType 0
				// Those are considered ahead of Cathedral and The Chest
				racer.PlaceMid++
			} else if racer2.FloorNum == racer.FloorNum && racer2.FloorArrived < racer.FloorArrived {
				racer.PlaceMid++
			}
		}
		if err := db.RaceParticipants.SetPlaceMid(racer.Name, raceID, racer.PlaceMid); err != nil {
			log.Error("Database error:", err)
			websocketError(s, d.Command, "")
			return false
		}
	}

	// Everything was set successfully
	return true
}

// Called after someone disconnects or someone is banned
// (the commandMutex should be locked when getting here)
func raceCheckStartFinish(raceID int) {
	// Get the status of the race
	if status, err := db.Races.GetStatus(raceID); err != nil {
		log.Error("Database error:", err)
		return
	} else if status == "open" {
		raceCheckStart(raceID)
	} else if status == "in progress" {
		raceCheckFinish(raceID)
	}
}

// Check to see if a race is ready to start, and if so, start it
// (the commandMutex should be locked when getting here)
func raceCheckStart(raceID int) {
	/*
		Validation
	*/

	// Check to see if the race was deleted
	if exists, err := db.Races.Exists(raceID); err != nil {
		log.Error("Database error:", err)
		return
	} else if !exists {
		return
	}

	// Check to see if this is a solo race
	solo, err2 := db.Races.CheckSolo(raceID)
	if err2 != nil {
		log.Error("Database error:", err2)
		return
	}

	// Check to see if there is only 1 person in the race
	if racerNames, err := db.RaceParticipants.GetRacerNames(raceID); err != nil {
		log.Error("Database error:", err)
		return
	} else if len(racerNames) == 1 && !solo {
		return
	}

	// Check if everyone is ready
	if sameStatus, err := db.RaceParticipants.CheckAllStatus(raceID, "ready"); err != nil {
		log.Error("Database error:", err)
		return
	} else if !sameStatus {
		return
	}

	/*
		Start the race
	*/

	// Log the race starting
	log.Info("Race " + strconv.Itoa(raceID) + " starting in 10 seconds.")

	// Change the status for this race to "starting"
	if err := db.Races.SetStatus(raceID, "starting"); err != nil {
		log.Error("Database error:", err)
		return
	}

	// Send everyone a notification that the race is starting soon
	for _, s := range websocketSessions {
		websocketEmit(s, "raceSetStatus", &RaceSetStatusMessage{raceID, "starting"})
	}

	// Get the list of people in this race
	racerList, err := db.RaceParticipants.GetRacerList(raceID)
	if err != nil {
		log.Error("Database error:", err)
		return
	}

	// Get the time X seconds in the future
	var secondsToWait time.Duration
	if solo {
		secondsToWait = 3
	} else {
		secondsToWait = 10
	}
	startTime := time.Now().Add(secondsToWait*time.Second).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))

	// Send everyone in the race a message specifying exactly when it will start
	for _, racer := range racerList {
		// A racer might go offline the moment before it starts, so check to see if the connection exists
		if s, ok := websocketSessions[racer.Name]; ok {
			websocketEmit(s, "raceStart", &RaceStartMessage{raceID, startTime})
		}
	}

	// Make the Twitch bot announce that the race is starting in 10 seconds
	if !solo {
		for _, racer := range racerList {
			twitchRacerSend(racer, "/me - The race is starting in 10 seconds!")
		}
	}

	// Return for now and do more things in 10 seconds
	go raceCheckStart2(raceID)
}

func raceCheckStart2(raceID int) {
	// Check to see if this is a solo race
	solo, err2 := db.Races.CheckSolo(raceID)
	if err2 != nil {
		log.Error("Database error:", err2)
		return
	}

	// Sleep 3 or 10 seconds
	var sleepTime time.Duration
	if solo {
		sleepTime = 3
	} else {
		sleepTime = 10
	}
	time.Sleep(sleepTime * time.Second)

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()
	defer commandMutex.Unlock()

	// Check to see if the race was deleted
	if exists, err := db.Races.Exists(raceID); err != nil {
		log.Error("Database error:", err)
		return
	} else if !exists {
		return
	}

	// Get the amount of people in this race
	racerNames, err := db.RaceParticipants.GetRacerNames(raceID)
	if err != nil {
		log.Error("Database error:", err)
		return
	}

	// Log the race starting
	log.Info("Race", raceID, "started with", len(racerNames), "participants:", racerNames)

	// Change the status for this race to "in progress" and set "datetime_started" equal to now
	if err := db.Races.Start(raceID); err != nil {
		log.Error("Database error:", err)
		return
	}

	// Update the status for everyone in the race to "racing"
	if err := db.RaceParticipants.SetAllStatus(raceID, "racing"); err != nil {
		log.Error("Database error:", err)
		return
	}

	// Send everyone a notification that the race is now in progress
	for _, s := range websocketSessions {
		websocketEmit(s, "raceSetStatus", &RaceSetStatusMessage{raceID, "in progress"})
	}

	// Return for now and do more things in 30 minutes
	go raceCheckStart3(raceID)
}

func raceCheckStart3(raceID int) {
	// Sleep 30 minutes
	time.Sleep(30 * time.Minute)

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()
	defer commandMutex.Unlock()

	// Find out if the race is finished
	if status, err := db.Races.GetStatus(raceID); err != nil {
		log.Error("Database error:", err)
		return
	} else if status == "finished" {
		return
	}

	// The race is still going, so get the list of people still in this race
	racerList, err := db.RaceParticipants.GetRacerList(raceID)
	if err != nil {
		log.Error("Database error:", err)
		return
	}

	// If any are still racing, force them to quit
	for _, racer := range racerList {
		if racer.Status == "racing" {
			if err := db.RaceParticipants.SetStatus(racer.Name, raceID, "quit"); err != nil {
				log.Error("Database error:", err)
				return
			}

			// Send a notification to all the people in this particular race that the user quit
			for _, racer2 := range racerList {
				// Not all racers may be online during a race
				if s, ok := websocketSessions[racer2.Name]; ok {
					websocketEmit(s, "racerSetStatus", &RacerSetStatusMessage{raceID, racer.Name, "quit", -1})
				}
			}
		}
	}

	// Close down the race
	raceCheckFinish(raceID)
}

// Check to see if a rate is ready to finish, and if so, finish it
// (the commandMutex should be locked when getting here)
func raceCheckFinish(raceID int) {
	// Check to see if the race was deleted
	if exists, err := db.Races.Exists(raceID); err != nil {
		log.Error("Database error:", err)
		return
	} else if !exists {
		return
	}

	// Check if anyone is still racing
	if stillRacing, err := db.RaceParticipants.CheckStillRacing(raceID); err != nil {
		log.Error("Database error:", err)
		return
	} else if stillRacing {
		return
	}

	// Log the race finishing
	log.Info("Race " + strconv.Itoa(raceID) + " finished.")

	// Change the status for this race to "finished" and set "datetime_finished" equal to now
	if err := db.Races.Finish(raceID); err != nil {
		log.Error("Database error:", err)
		return
	}

	// Send everyone a notification that the race is now finished
	for _, s := range websocketSessions {
		websocketEmit(s, "raceSetStatus", &RaceSetStatusMessage{raceID, "finished"})
	}
}

// Now that a user has finished, quit, or been disqualified from a race, update fields in the users table for unseeded races
func raceUpdateUnseededStats(raceID int, username string) {
	// Don't do anything if this is not an unseeded race (or an unranked race)
	/*
		if unseededAndRanked, err := db.Races.CheckUnseededRanked(raceID); err != nil {
			log.Error("Database error:", err)
			return
		} else if !unseededAndRanked {
			return
		}
	*/

	// Get their unseeded stats
	/*
		if statsUnseeded, err := db.Users.GetStatsUnseeded(username); err != nil {
			log.Error("Database error:", err)
			return
		}
	*/

	// Update all the stats

}

// Now that the race has finished, update fields in the users table for seeded races
func raceUpdateSeededStats(raceID int, username string) {

}
