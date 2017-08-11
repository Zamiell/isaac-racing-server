package main

/*
	Imports
*/

import (
	"strconv"

	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceJoin(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	raceID := d.ID
	userID := d.v.UserID
	username := d.v.Username

	// Validate basic things about the race ID
	if !raceValidate(s, d) {
		return
	}

	// Validate that the race is open
	if !raceValidateStatus(s, d, "open") {
		return
	}

	// Validate that they are not in the race
	if !raceValidateOut2(s, d) {
		return
	}

	// Validate that this is not a solo race
	if solo, err := db.Races.CheckSolo(raceID); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if solo {
		log.Warning("User \"" + username + "\" attempted to call " + d.Command + " on race ID " + strconv.Itoa(raceID) + ", but it is a solo race.")
		websocketError(s, d.Command, "Race ID "+strconv.Itoa(raceID)+" is a solo race, so you cannot join it.")
		return
	}

	// Add this user to the participants list for that race
	if err := db.RaceParticipants.Insert(userID, raceID); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Send everyone a notification that the user joined
	for _, s := range websocketSessions {
		type RaceJoinedMessage struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		websocketEmit(s, "raceJoined", &RaceJoinedMessage{raceID, username})
	}

	// Get all the information about the racers in this race
	racerList, err := db.RaceParticipants.GetRacerList(raceID)
	if err != nil {
		log.Error("Database error:", err)
		return
	}

	// Send it to the user
	websocketEmit(s, "racerList", &RacerListMessage{raceID, racerList})

	// Join the user to the channel for that race
	d.Room = "_race_" + strconv.Itoa(raceID)
	websocketRoomJoinSub(s, d)
}
