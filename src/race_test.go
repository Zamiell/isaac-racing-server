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

	// If either racer 1 or racer 2 is on the backwards path, it should not affect the results
	races := make([]*server.Race, 0)
	races = append(races, getStandardRaces(goal, false, false)...)
	races = append(races, getStandardRaces(goal, true, false)...)
	races = append(races, getStandardRaces(goal, false, true)...)
	races = append(races, getStandardRaces(goal, true, true)...)

	testRaces(t, races)
}

func TestRaceTheBeast(t *testing.T) {
	t.Parallel()

	goal := server.RaceGoalBeast

	// If neither racer is on the backwards path, all the standard tests should pass
	races := getStandardRaces(goal, false, false)

	// If racer 1 is on the backwards path, all the standard tests should pass
	races = append(races, getStandardRaces(goal, true, false)...)

	// On backwards path with opponent not
	aliceOnBackwardsPathRaces := []*server.Race{
		getRace(3, 0, true, 4, 0, false, goal),
		getRace(3, 4, true, 4, 0, false, goal),
		getRace(3, 0, true, 4, 4, false, goal),
	}
	races = append(races, aliceOnBackwardsPathRaces...)

	bothOnBackwardsPathRaces := []*server.Race{
		// Ahead on normal floors (lower is better)
		// (e.g. StageTypeOriginal, StageTypeWoTL, StageTypeAfterbirth)
		getRace(2, 0, true, 3, 0, true, goal),
		getRace(2, 0, true, 3, 1, true, goal),
		getRace(2, 0, true, 3, 2, true, goal),
		getRace(2, 1, true, 3, 0, true, goal),
		getRace(2, 1, true, 3, 1, true, goal),
		getRace(2, 1, true, 3, 2, true, goal),
		getRace(2, 2, true, 3, 0, true, goal),
		getRace(2, 2, true, 3, 1, true, goal),
		getRace(2, 2, true, 3, 2, true, goal),

		// Same floor on normal floors
		getRace(3, 0, true, 3, 0, true, goal),
		getRace(3, 0, true, 3, 1, true, goal),
		getRace(3, 0, true, 3, 2, true, goal),
		getRace(3, 1, true, 3, 0, true, goal),
		getRace(3, 1, true, 3, 1, true, goal),
		getRace(3, 1, true, 3, 2, true, goal),
		getRace(3, 2, true, 3, 0, true, goal),
		getRace(3, 2, true, 3, 1, true, goal),
		getRace(3, 2, true, 3, 2, true, goal),

		// Ahead on Repentance floors
		getRace(1, 4, true, 2, 0, true, goal),
		getRace(1, 4, true, 2, 1, true, goal),
		getRace(1, 4, true, 2, 2, true, goal),
		getRace(1, 5, true, 2, 0, true, goal),
		getRace(1, 5, true, 2, 1, true, goal),
		getRace(1, 5, true, 2, 2, true, goal),

		// Ahead on Repentance floors
		// (lower real stage, same adjusted stage, Alice on Repentance)
		getRace(2, 4, true, 3, 0, true, goal),
		getRace(2, 4, true, 3, 1, true, goal),
		getRace(2, 4, true, 3, 2, true, goal),
		getRace(2, 5, true, 3, 0, true, goal),
		getRace(2, 5, true, 3, 1, true, goal),
		getRace(2, 5, true, 3, 2, true, goal),

		// Ahead on Repentance floors
		// (lower real stage, same adjusted stage, Bob on Repentance)
		getRace(3, 0, true, 2, 4, true, goal),
		getRace(3, 1, true, 2, 4, true, goal),
		getRace(3, 2, true, 2, 4, true, goal),
		getRace(3, 0, true, 2, 5, true, goal),
		getRace(3, 1, true, 2, 5, true, goal),
		getRace(3, 2, true, 2, 5, true, goal),

		// Same floor on Repentance floor
		// (same real stage)
		getRace(3, 4, true, 3, 4, true, goal),
		getRace(3, 4, true, 3, 5, true, goal),
		getRace(3, 5, true, 3, 4, true, goal),
		getRace(3, 5, true, 3, 5, true, goal),
	}

	races = append(races, getRace(2, 0, true, 3, 0, true, goal))
	races = append(races, getRace(2, 0, true, 3, 1, true, goal))
	races = append(races, getRace(2, 0, true, 3, 2, true, goal))
	races = append(races, getRace(2, 1, true, 3, 0, true, goal))
	races = append(races, getRace(2, 1, true, 3, 1, true, goal))
	races = append(races, getRace(2, 1, true, 3, 2, true, goal))
	races = append(races, getRace(2, 2, true, 3, 0, true, goal))
	races = append(races, getRace(2, 2, true, 3, 1, true, goal))
	races = append(races, getRace(2, 2, true, 3, 2, true, goal))

	// Both on backwards path (Repentance floors)
	races = append(races, getRace(2, 4, true, 3, 0, true, goal))
	races = append(races, getRace(2, 4, true, 3, 1, true, goal))
	races = append(races, getRace(2, 4, true, 3, 2, true, goal))
	races = append(races, getRace(2, 5, true, 3, 0, true, goal))
	races = append(races, getRace(2, 5, true, 3, 1, true, goal))
	races = append(races, getRace(2, 5, true, 3, 2, true, goal))

	races = append(races, bothOnBackwardsPathRaces...)

	testRaces(t, races)
}

func getStandardRaces(
	goal string,
	racer1OnBackwardsPath bool,
	racer2OnBackwardsPath bool,
) []*server.Race {
	return []*server.Race{
		// Ahead on normal floors
		// (e.g. StageTypeOriginal, StageTypeWoTL, StageTypeAfterbirth)
		getRace(3, 0, racer1OnBackwardsPath, 2, 0, racer2OnBackwardsPath, goal),
		getRace(3, 0, racer1OnBackwardsPath, 2, 1, racer2OnBackwardsPath, goal),
		getRace(3, 0, racer1OnBackwardsPath, 2, 2, racer2OnBackwardsPath, goal),
		getRace(3, 1, racer1OnBackwardsPath, 2, 0, racer2OnBackwardsPath, goal),
		getRace(3, 1, racer1OnBackwardsPath, 2, 1, racer2OnBackwardsPath, goal),
		getRace(3, 1, racer1OnBackwardsPath, 2, 2, racer2OnBackwardsPath, goal),
		getRace(3, 2, racer1OnBackwardsPath, 2, 0, racer2OnBackwardsPath, goal),
		getRace(3, 2, racer1OnBackwardsPath, 2, 1, racer2OnBackwardsPath, goal),
		getRace(3, 2, racer1OnBackwardsPath, 2, 2, racer2OnBackwardsPath, goal),

		// Same floor on normal floors
		getRace(3, 0, racer1OnBackwardsPath, 3, 0, racer2OnBackwardsPath, goal),
		getRace(3, 0, racer1OnBackwardsPath, 3, 1, racer2OnBackwardsPath, goal),
		getRace(3, 0, racer1OnBackwardsPath, 3, 2, racer2OnBackwardsPath, goal),
		getRace(3, 1, racer1OnBackwardsPath, 3, 0, racer2OnBackwardsPath, goal),
		getRace(3, 1, racer1OnBackwardsPath, 3, 1, racer2OnBackwardsPath, goal),
		getRace(3, 1, racer1OnBackwardsPath, 3, 2, racer2OnBackwardsPath, goal),
		getRace(3, 2, racer1OnBackwardsPath, 3, 0, racer2OnBackwardsPath, goal),
		getRace(3, 2, racer1OnBackwardsPath, 3, 1, racer2OnBackwardsPath, goal),
		getRace(3, 2, racer1OnBackwardsPath, 3, 2, racer2OnBackwardsPath, goal),
		getRace(3, 3, racer1OnBackwardsPath, 3, 0, racer2OnBackwardsPath, goal),
		getRace(3, 3, racer1OnBackwardsPath, 3, 1, racer2OnBackwardsPath, goal),
		getRace(3, 3, racer1OnBackwardsPath, 3, 2, racer2OnBackwardsPath, goal),

		// Ahead on Repentance floors
		getRace(3, 4, racer1OnBackwardsPath, 2, 0, racer2OnBackwardsPath, goal),
		getRace(3, 4, racer1OnBackwardsPath, 2, 1, racer2OnBackwardsPath, goal),
		getRace(3, 4, racer1OnBackwardsPath, 2, 2, racer2OnBackwardsPath, goal),
		getRace(3, 5, racer1OnBackwardsPath, 2, 0, racer2OnBackwardsPath, goal),
		getRace(3, 5, racer1OnBackwardsPath, 2, 1, racer2OnBackwardsPath, goal),
		getRace(3, 5, racer1OnBackwardsPath, 2, 2, racer2OnBackwardsPath, goal),

		// Ahead on Repentance floors
		// (same real stage, higher adjusted stage)
		getRace(3, 4, racer1OnBackwardsPath, 3, 0, racer2OnBackwardsPath, goal),
		getRace(3, 4, racer1OnBackwardsPath, 3, 1, racer2OnBackwardsPath, goal),
		getRace(3, 4, racer1OnBackwardsPath, 3, 2, racer2OnBackwardsPath, goal),
		getRace(3, 5, racer1OnBackwardsPath, 3, 0, racer2OnBackwardsPath, goal),
		getRace(3, 5, racer1OnBackwardsPath, 3, 1, racer2OnBackwardsPath, goal),
		getRace(3, 5, racer1OnBackwardsPath, 3, 2, racer2OnBackwardsPath, goal),

		// Same floor on Repentance floor
		// (same real stage)
		getRace(3, 4, racer1OnBackwardsPath, 3, 4, racer2OnBackwardsPath, goal),
		getRace(3, 4, racer1OnBackwardsPath, 3, 5, racer2OnBackwardsPath, goal),
		getRace(3, 5, racer1OnBackwardsPath, 3, 4, racer2OnBackwardsPath, goal),
		getRace(3, 5, racer1OnBackwardsPath, 3, 5, racer2OnBackwardsPath, goal),

		// Same floor on Repentance floors
		// (lower real stage, same adjusted stage, Alice on Repentance)
		getRace(3, 4, racer1OnBackwardsPath, 4, 0, racer2OnBackwardsPath, goal),
		getRace(3, 4, racer1OnBackwardsPath, 4, 1, racer2OnBackwardsPath, goal),
		getRace(3, 4, racer1OnBackwardsPath, 4, 2, racer2OnBackwardsPath, goal),
		getRace(3, 5, racer1OnBackwardsPath, 4, 0, racer2OnBackwardsPath, goal),
		getRace(3, 5, racer1OnBackwardsPath, 4, 1, racer2OnBackwardsPath, goal),
		getRace(3, 5, racer1OnBackwardsPath, 4, 2, racer2OnBackwardsPath, goal),

		// Repentance floors that should be equal
		// (lower real stage, same adjusted stage, Bob on Repentance)
		getRace(3, 0, racer1OnBackwardsPath, 2, 4, racer2OnBackwardsPath, goal),
		getRace(3, 1, racer1OnBackwardsPath, 2, 4, racer2OnBackwardsPath, goal),
		getRace(3, 2, racer1OnBackwardsPath, 2, 4, racer2OnBackwardsPath, goal),
		getRace(3, 0, racer1OnBackwardsPath, 2, 5, racer2OnBackwardsPath, goal),
		getRace(3, 1, racer1OnBackwardsPath, 2, 5, racer2OnBackwardsPath, goal),
		getRace(3, 2, racer1OnBackwardsPath, 2, 5, racer2OnBackwardsPath, goal),
	}
}

func getRace(
	racer1FloorNum int,
	racer1StageType int,
	racer1OnBackwardsPath bool,
	racer2FloorNum int,
	racer2StageType int,
	racer2OnBackwardsPath bool,
	goal string,
) *server.Race {
	racer1 := &server.Racer{
		ID:                   1,
		Name:                 Racer1Name,
		Status:               "racing",
		FloorNum:             racer1FloorNum,
		StageType:            racer1StageType,
		BackwardsPath:        racer1OnBackwardsPath,
		DatetimeArrivedFloor: Racer1ArrivedFloor,
		PlaceMid:             -1,
	}

	racer2 := &server.Racer{
		ID:                   2,
		Name:                 Racer2Name,
		Status:               "racing",
		FloorNum:             racer2FloorNum,
		StageType:            racer2StageType,
		BackwardsPath:        racer2OnBackwardsPath,
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
	t.Helper()

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

		racer2PlaceMid := race.Racers[Racer2Name].PlaceMid
		if racer2PlaceMid != 2 {
			t.Errorf(
				"Race #%d failed: %s should be in 2nd place, but was place: %d",
				index+1,
				Racer2Name,
				racer2PlaceMid,
			)
		}
	}
}
