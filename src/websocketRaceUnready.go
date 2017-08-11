package main

/*
	Imports
*/

import (
	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceUnready(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	raceID := d.ID
	username := d.v.Username

	// Validate basic things about the race ID
	if !raceValidate(s, d) {
		return
	}

	// Validate that the race is open
	if !raceValidateStatus(s, d, "open") {
		return
	}

	// Validate that they are in the race
	if !raceValidateIn2(s, d) {
		return
	}

	// Validate that their status is set to "ready"
	if !racerValidateStatus(s, d, "ready") {
		return
	}

	// Change their status to "not ready"
	if !racerSetStatus(s, d, "not ready") {
		return
	}

	// Get the list of racers for this race
	racerNames, err := db.RaceParticipants.GetRacerNames(raceID)
	if err != nil {
		log.Error("Database error:", err)
		return
	}

	// Send a notification to all the people in this particular race that the user is not ready
	for _, racer := range racerNames {
		// Not all racers may be online during a race
		if s, ok := websocketSessions[racer]; ok {
			websocketEmit(s, "racerSetStatus", &RacerSetStatusMessage{raceID, username, "not ready", 0})
		}
	}
}
