package main

/*
	Imports
*/

import (
	"strconv"

	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceQuit(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	raceID := d.ID
	username := d.v.Username

	// Validate basic things about the race ID
	if !raceValidate(s, d) {
		return
	}

	// Validate that the race has started
	if !raceValidateStatus(s, d, "in progress") {
		return
	}

	// Validate that they are in the race
	if !raceValidateIn2(s, d) {
		return
	}

	// Validate that their status is set to "racing" status
	if !racerValidateStatus(s, d, "racing") {
		return
	}

	// Change their status to "quit"
	if !racerSetStatus(s, d, "quit") {
		return
	}

	// Set their finish time
	if err := db.RaceParticipants.SetDatetimeFinished(username, raceID, int(makeTimestamp())); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Set their (final) place to -1 (which indicates a quit status)
	if err := db.RaceParticipants.SetPlace(username, raceID, -1); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Get the list of racers for this race
	racerList, err := db.RaceParticipants.GetRacerList(raceID)
	if err != nil {
		log.Error("Database error:", err)
		return
	}

	// Recalculate everyones mid-race places
	if !racerSetAllPlaceMid(s, d, racerList) {
		return
	}

	// Send a notification to all the people in this particular race that the user quit
	for _, racer := range racerList {
		// Not all racers may be online during a race
		if s2, ok := websocketSessions[racer.Name]; ok {
			websocketEmit(s2, "racerSetStatus", &RacerSetStatusMessage{raceID, username, "quit", -1})
		}
	}

	// Get the number of people left in the race
	peopleLeft, err := db.RaceParticipants.GetPeopleLeft(raceID)
	if err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Make the Twitch bot announce that the person quit
	twitchString := "/me - " + username + " quit - "
	if peopleLeft == 0 {
		twitchString += "Race completed."
	} else {
		twitchString += strconv.Itoa(peopleLeft) + " left"
	}
	for _, racer := range racerList {
		twitchRacerSend(racer, twitchString)
	}

	// Check to see if the race is ready to finish
	raceCheckFinish(raceID)

	// Update fields in the users table (e.g. average, ELO)
	// (we calculate stats for seeded races only when the race is completed)
	raceUpdateUnseededStats(raceID, username)
}
