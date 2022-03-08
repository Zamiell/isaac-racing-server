package server

import (
	"strconv"

	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceFloor(s *melody.Session, d *IncomingWebsocketData) {
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
	// (floor 14 is a fake floor that we use to represent Mega Satan)
	if floorNum < 1 || floorNum > 14 {
		logger.Warning("User \"" + username + "\" attempted to update their floor, but \"" + strconv.Itoa(floorNum) + "\" is a bogus floor number.")
		websocketError(s, d.Command, "That is not a valid floor number.")
		return
	} else if stageType < 0 || stageType > 5 {
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
	racer.BackwardsPath = d.BackwardsPath
	racer.DatetimeArrivedFloor = getTimestamp()

	// If they reset from floor 1 to floor 1,
	// don't send the new floor to everyone as an optimization
	// We also do not have to recalculate the placeMids,
	// because placeMid should not be updated until they get to the second floor
	if floorNum == 1 && oldFloor == 1 && !isRepentanceStageType(stageType) && !d.BackwardsPath {
		return
	}

	race.SetAllPlaceMid()
	race.SendAllFloor(racer)
}

func (race *Race) SendAllFloor(racer *Racer) {
	leader := race.GetLeader()

	for racerName := range race.Racers {
		// Not all racers may be online during a race
		if s, ok := websocketSessions[racerName]; ok {
			millisecondsBehindLeader := int64(0)
			if leader != nil && racer.PlaceMid > 1 {
				millisecondsBehindLeader = racer.DatetimeArrivedFloor - leader.DatetimeArrivedFloor
			}

			type RacerSetFloorMessage struct {
				ID                       int    `json:"id"`
				Name                     string `json:"name"`
				FloorNum                 int    `json:"floorNum"`
				StageType                int    `json:"stageType"`
				DatetimeArrivedFloor     int64  `json:"datetimeArrivedFloor"`
				MillisecondsBehindLeader int64  `json:"millisecondsBehindLeader"`
			}
			websocketEmit(s, "racerSetFloor", &RacerSetFloorMessage{
				ID:                       race.ID,
				Name:                     racer.Name,
				FloorNum:                 racer.FloorNum,
				StageType:                racer.StageType,
				DatetimeArrivedFloor:     racer.DatetimeArrivedFloor,
				MillisecondsBehindLeader: millisecondsBehindLeader,
			})
		}
	}
}

func (race *Race) GetLeader() *Racer {
	var leader *Racer
	for _, racer := range race.Racers {
		// Skip racers who have finished or quit
		if racer.PlaceMid == -1 {
			continue
		}

		if leader == nil {
			leader = racer
			continue
		}

		if racer.PlaceMid < leader.PlaceMid {
			leader = racer
		}
	}

	return leader
}
