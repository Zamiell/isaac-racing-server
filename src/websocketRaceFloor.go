package main

import (
	"strconv"

	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceFloor(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username
	floorNum := d.FloorNum
	stageType := d.StageType

	/*
		Validation
	*/

	// Validate that the race exists
	var race *Race
	if v, ok := races[d.ID]; !ok {
		return
	} else {
		race = v
	}

	// Validate that the race has started
	if race.Status != RaceStatusInProgress {
		return
	}

	// Validate that they are in the race
	var racer *Racer
	if v, ok := race.Racers[username]; !ok {
		return
	} else {
		racer = v
	}

	// Validate that they are still racing
	if racer.Status != RacerStatusRacing {
		return
	}

	// Validate that the floor is sane
	// (floor 13 is Home, which is the final floor)
	if floorNum < 1 || floorNum > 13 {
		logger.Warning("User \"" + username + "\" attempted to update their floor, but \"" + strconv.Itoa(floorNum) + "\" is a bogus floor number.")
		websocketError(s, d.Command, "That is not a valid floor number.")
		return
	} else if (stageType < 0 || stageType > 5) && stageType != 3 {
		logger.Warning("User \"" + username + "\" attempted to update their floor, but \"" + strconv.Itoa(stageType) + "\" is a bogus stage type.")
		websocketError(s, d.Command, "That is not a valid stage type.")
		return
	}

	/*
		Set the floor
	*/

	oldFloor := racer.FloorNum

	racer.FloorNum = floorNum
	racer.StageType = stageType
	racer.DatetimeArrivedFloor = getTimestamp()

	// If they reset from floor 1 to floor 1,
	// don't send the new floor to everyone as an optimization
	// We also do not have to recalculate the placeMids,
	// because placeMid is not assigned until they get to the second floor
	if floorNum == 1 && oldFloor == 1 && stageType != 4 && stageType != 5 {
		return
	}

	for racerName := range race.Racers {
		// Not all racers may be online during a race
		if s, ok := websocketSessions[racerName]; ok {
			type RacerSetFloorMessage struct {
				ID                   int    `json:"id"`
				Name                 string `json:"name"`
				FloorNum             int    `json:"floorNum"`
				StageType            int    `json:"stageType"`
				DatetimeArrivedFloor int64  `json:"datetimeArrivedFloor"`
			}
			websocketEmit(s, "racerSetFloor", &RacerSetFloorMessage{
				ID:                   race.ID,
				Name:                 racer.Name,
				FloorNum:             racer.FloorNum,
				StageType:            racer.StageType,
				DatetimeArrivedFloor: racer.DatetimeArrivedFloor,
			})
		}
	}

	race.SetAllPlaceMid()
}
