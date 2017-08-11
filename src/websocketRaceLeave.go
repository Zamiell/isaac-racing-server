package main

/*
	Imports
*/

import (
	"strconv"

	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceLeave(s *melody.Session, d *IncomingWebsocketData) {
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

	// Remove this user from the participants list for that race
	if err := db.RaceParticipants.Delete(username, raceID); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Disconnect the user from the channel for that race
	d.Room = "_race_" + strconv.Itoa(raceID)
	websocketRoomLeaveSub(s, d)

	// Send everyone a notification that the user left the race
	for _, s := range websocketSessions {
		websocketEmit(s, "raceLeft", &RaceLeftMessage{raceID, username})
	}

	// If the race went from 2 people to 1, automatically unready the last person so that they don't start the race by themsevles
	if racerNames, err := db.RaceParticipants.GetRacerNames(raceID); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if len(racerNames) == 1 {
		if currentStatus, err := db.RaceParticipants.GetStatus(racerNames[0], raceID); err != nil {
			log.Error("Database error:", err)
			websocketError(s, d.Command, "")
			return
		} else if currentStatus == "ready" {
			// Set them from "ready" to "not ready"
			if err := db.RaceParticipants.SetStatus(racerNames[0], raceID, "not ready"); err != nil {
				log.Error("Database error:", err)
				websocketError(s, d.Command, "")
				return
			}

			// Tell them
			// (they should definately be online, but check just in case)
			if s, ok := websocketSessions[racerNames[0]]; ok {
				websocketEmit(s, "racerSetStatus", &RacerSetStatusMessage{raceID, username, "not ready", 0})
			}
		}
	}

	// Check to see if the race is ready to start
	raceCheckStart(raceID)
}
