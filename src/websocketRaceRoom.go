package main

/*
	Imports
*/

import (
	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceRoom(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	raceID := d.ID
	roomID := d.RoomID
	userID := d.v.UserID

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

	// Add the room to their list of visited rooms
	if err := db.RaceParticipantRooms.Insert(userID, raceID, roomID); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}
}
