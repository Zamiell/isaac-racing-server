package server_test

import (
	"fmt"
	"testing"

	server "github.com/Zamiell/isaac-racing-server"
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

		if race.Racers["racer1"].PlaceMid != 1 {
			t.Errorf(fmt.Sprintf("Place for racer1: %d on race number %d", race.Racers["racer1"].PlaceMid, index+1))
		}
	}
}

func getRace(racer1FloorNum int, racer1StageType int, racer1IsOnBackwardsPath bool, racer2FloorNum int, racer2StageType int, racer2IsOnBackwardsPath bool) server.Race {
	ruleset := server.Ruleset{
		Ranked:              false,
		Solo:                true,
		Format:              server.RaceFormatSeeded,
		Character:           "Judas",
		CharacterRandom:     false,
		Goal:                server.RaceGoalBeast,
		StartingBuild:       1,
		StartingBuildRandom: false,
		Seed:                "TESTTEST",
		Difficulty:          "normal",
	}

	racer1 := server.Racer{
		ID:                   0,
		Name:                 "racer1",
		Status:               "racing",
		Seed:                 "TESTTEST",
		FloorNum:             racer1FloorNum,
		StageType:            racer1StageType,
		BackwardsPath:        racer1IsOnBackwardsPath,
		DatetimeArrivedFloor: 30000,
		CharacterNum:         1,
		Place:                0,
		PlaceMid:             -1,
	}

	racer2 := server.Racer{
		ID:                   10000,
		Name:                 "racer1",
		Status:               "racing",
		Seed:                 "TESTTEST",
		FloorNum:             racer2FloorNum,
		StageType:            racer2StageType,
		BackwardsPath:        racer2IsOnBackwardsPath,
		DatetimeArrivedFloor: 30050,
		CharacterNum:         1,
		Place:                0,
		PlaceMid:             -1,
	}

	Racers := make(map[string]*server.Racer)
	Racers["racer1"] = &racer1
	Racers["racer2"] = &racer2

	race := server.Race{
		ID:              0,
		Name:            "name",
		Status:          server.RaceStatusOpen,
		Ruleset:         ruleset,
		Captain:         "username",
		Password:        "password",
		SoundPlayed:     false,
		DatetimeCreated: 0,
		DatetimeStarted: 0,
		Racers:          Racers,
	}

	return race
}
