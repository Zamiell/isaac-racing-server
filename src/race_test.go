package server_test

import (
	"testing"

	server "github.com/Zamiell/isaac-racing-server"
)

const (
	PLAYER_1_NAME         = "Alice"
	PLAYER_2_NAME         = "Bob"
	RACER_1_ARRIVED_FLOOR = 30000
)

func TestRaceTheBeast(t *testing.T) {
	t.Parallel()

	races := make([]server.Race, 7)
	races[0] = getRace(3, 0, false, 2, 0, false)
	races[1] = getRace(3, 0, false, 2, 4, false)
	races[2] = getRace(3, 0, false, 3, 0, false)
	races[3] = getRace(3, 5, false, 3, 0, false)
	races[4] = getRace(6, 4, true, 6, 4, false)
	races[5] = getRace(2, 4, true, 3, 0, true)
	races[6] = getRace(2, 4, false, 2, 4, false)

	for index, race := range races {
		race.SetAllPlaceMid()

		if race.Racers[PLAYER_1_NAME].PlaceMid != 1 {
			t.Errorf(
				"Place for %s: %d on race number %d",
				PLAYER_1_NAME,
				race.Racers[PLAYER_1_NAME].PlaceMid,
				index+1,
			)
		}
	}
}

func getRace(
	racer1FloorNum int,
	racer1StageType int,
	racer1IsOnBackwardsPath bool,
	racer2FloorNum int,
	racer2StageType int,
	racer2IsOnBackwardsPath bool,
) server.Race {
	racer1 := server.Racer{
		Name:                 PLAYER_1_NAME,
		Status:               "racing",
		FloorNum:             racer1FloorNum,
		StageType:            racer1StageType,
		BackwardsPath:        racer1IsOnBackwardsPath,
		DatetimeArrivedFloor: RACER_1_ARRIVED_FLOOR,
		PlaceMid:             -1,
	}

	racer2 := server.Racer{
		Name:                 PLAYER_2_NAME,
		Status:               "racing",
		FloorNum:             racer2FloorNum,
		StageType:            racer2StageType,
		BackwardsPath:        racer2IsOnBackwardsPath,
		DatetimeArrivedFloor: RACER_1_ARRIVED_FLOOR + 50,
		PlaceMid:             -1,
	}

	Racers := make(map[string]*server.Racer)
	Racers["Alice"] = &racer1
	Racers["Bob"] = &racer2

	race := server.Race{
		Ruleset: server.Ruleset{
			Goal: server.RaceGoalBeast,
		},
		Racers: Racers,
	}

	return race
}
