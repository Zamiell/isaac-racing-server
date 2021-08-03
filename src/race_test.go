package server_test

import (
	"testing"

	server "github.com/Zamiell/isaac-racing-server"
)

const (
	Racer1Name         = "Alice"
	Racer2Name         = "Bob"
	Racer1ArrivedFloor = 30000
)

func TestRaceBlueBaby(t *testing.T) {
	t.Parallel()

	goal := server.RaceGoalBlueBaby
	races := []*server.Race{
		getRace(3, 0, false, 1, 1, false, goal), // Normal floors
		getRace(2, 2, false, 1, 4, false, goal), // Repentance floors
	}

	testRaces(t, races)
}

func TestRaceTheBeast(t *testing.T) {
	t.Parallel()

	goal := server.RaceGoalBeast
	races := []*server.Race{
		getRace(3, 0, false, 2, 0, false, goal),
		getRace(3, 0, false, 2, 4, false, goal),
		getRace(3, 0, false, 3, 0, false, goal),
		getRace(3, 5, false, 3, 0, false, goal),
		getRace(6, 4, true, 6, 4, false, goal),
		getRace(2, 4, true, 3, 0, true, goal),
		getRace(2, 4, false, 2, 4, false, goal),
	}

	testRaces(t, races)
}

func getRace(
	racer1FloorNum int,
	racer1StageType int,
	racer1IsOnBackwardsPath bool,
	racer2FloorNum int,
	racer2StageType int,
	racer2IsOnBackwardsPath bool,
	goal string,
) *server.Race {
	racer1 := &server.Racer{
		Name:                 Racer1Name,
		Status:               "racing",
		FloorNum:             racer1FloorNum,
		StageType:            racer1StageType,
		BackwardsPath:        racer1IsOnBackwardsPath,
		DatetimeArrivedFloor: Racer1ArrivedFloor,
		PlaceMid:             -1,
	}

	racer2 := &server.Racer{
		Name:                 Racer2Name,
		Status:               "racing",
		FloorNum:             racer2FloorNum,
		StageType:            racer2StageType,
		BackwardsPath:        racer2IsOnBackwardsPath,
		DatetimeArrivedFloor: Racer1ArrivedFloor + 50,
		PlaceMid:             -1,
	}

	racers := make(map[string]*server.Racer)
	racers[Racer1Name] = racer1
	racers[Racer2Name] = racer2

	return &server.Race{
		Ruleset: server.Ruleset{
			Goal: goal,
		},
		Racers: racers,
	}
}

func testRaces(t *testing.T, races []*server.Race) {
	for index, race := range races {
		race.SetAllPlaceMid()

		racer1PlaceMid := race.Racers[Racer1Name].PlaceMid
		if racer1PlaceMid != 1 {
			t.Errorf(
				"Race #%d failed: %s should be in 1st place, but was place: %d",
				index+1,
				Racer1Name,
				racer1PlaceMid,
			)
		}
	}
}
