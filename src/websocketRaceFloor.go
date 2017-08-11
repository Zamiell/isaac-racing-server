package main

/*
	Imports
*/

import (
	"strconv"

	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceFloor(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	raceID := d.ID
	floorNum := d.FloorNum
	stageType := d.StageType
	userID := d.v.UserID
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

	// Validate that the floor is sane
	if floorNum < 1 || floorNum > 13 {
		// The Void is floor 12, and we use floor 13 to signify Mega Satan
		log.Warning("User \"" + username + "\" attempted to update their floor, but \"" + strconv.Itoa(floorNum) + "\" is a bogus floor number.")
		websocketError(s, d.Command, "That is not a valid floor number.")
		return
	} else if stageType < 0 || stageType > 3 {
		log.Warning("User \"" + username + "\" attempted to update their floor, but \"" + strconv.Itoa(stageType) + "\" is a bogus stage type.")
		websocketError(s, d.Command, "That is not a valid stage type.")
		return
	}

	// Set their floor in the database
	floorArrived, err := db.RaceParticipants.SetFloor(userID, raceID, floorNum, stageType)
	if err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// The floor gets sent as 1 when a reset occurs
	if floorNum == 1 {
		// Reset all of their accumulated items
		if err = db.RaceParticipantItems.Reset(userID, raceID); err != nil {
			log.Error("Database error:", err)
			websocketError(s, d.Command, "")
			return
		}

		// Reset all of their visited rooms
		if err = db.RaceParticipantItems.Reset(userID, raceID); err != nil {
			log.Error("Database error:", err)
			websocketError(s, d.Command, "")
			return
		}
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

	// Send a notification to all the people in this particular race that the user got to a new floor
	for _, racer := range racerList {
		// Not all racers may be online during a race
		if s, ok := websocketSessions[racer.Name]; ok {
			websocketEmit(s, "racerSetFloor", &RacerSetFloorMessage{raceID, username, floorNum, stageType, floorArrived})
		}
	}
}
