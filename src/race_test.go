package main

import (
	"fmt"
	"testing"
)

func TestRacersPlaces(t *testing.T) {

	races := make([]Race, 7)
	races[0] = getRacersData(3, 0, 2, 0, false, false)
	races[1] = getRacersData(3, 0, 2, 4, false, false)
	races[2] = getRacersData(3, 0, 3, 0, false, false)
	races[3] = getRacersData(3, 5, 3, 0, false, false)
	races[4] = getRacersData(6, 4, 6, 4, true, false)
	races[5] = getRacersData(2, 4, 3, 0, true, true)
	races[6] = getRacersData(2, 4, 2, 4, false, false)

	for index, race := range races {
		race.SetAllPlaceMid()

		if race.Racers["racer1"].PlaceMid != 1 {
			t.Errorf(fmt.Sprintf("Place for racer1: %d on race number %d", race.Racers["racer1"].PlaceMid, index+1))
		}
	}
}

func getRacersData(racer1FloorNum, racer1StageType, racer2FloorNum, racer2StageType int, racer1IsOnBackwardsPath, racer2IsOnBackwardsPath bool) Race {
	ruleset := Ruleset{
		Ranked:              false,
		Solo:                true,
		Format:              RaceFormatSeeded,
		Character:           "Judas",
		CharacterRandom:     false,
		Goal:                RaceGoalBeast,
		StartingBuild:       1,
		StartingBuildRandom: false,
		Seed:                "TESTTEST",
		Difficulty:          "normal",
	}

	racer1 := Racer{
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

	racer2 := Racer{
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

	Racers := make(map[string]*Racer)
	Racers["racer1"] = &racer1
	Racers["racer2"] = &racer2

	race := Race{
		ID:              0,
		Name:            "name",
		Status:          RaceStatusOpen,
		Ruleset:         ruleset,
		Captain:         "username",
		Password:        "password",
		SoundPlayed:     false,
		DatetimeCreated: getTimestamp(),
		DatetimeStarted: 0,
		Racers:          Racers,
	}

	return race
}
