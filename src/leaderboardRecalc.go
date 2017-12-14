package main

import (
	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/Zamiell/isaac-racing-server/src/models"
)

func leaderboardSeededRecalculate() {

}

func leaderboardUnseededRecalculate() {

}

func leaderboardDiversityRecalculate() {
	if err := db.Users.ResetStatsDiversity(); err != nil {
		log.Error("Database error while resetting the diversity stats:", err)
	}

	var allDivRaces []models.RaceHistory
	if v, err := db.Races.GetAllDiversityRaces(); err != nil {
		log.Error("Database error while getting all of the diversity races:", err)
	} else {
		allDivRaces = v
	}

	// Go through every diversity race in the database
	for _, modelsRace := range allDivRaces {
		// Convert the "RaceHistory" struct to a "Race" struct
		race := &Race{}
		race.ID = modelsRace.RaceID
		race.Racers = make(map[string]*Racer)
		for _, modelsRacer := range modelsRace.RaceParticipants {
			racer := &Racer{
				ID:    modelsRacer.ID,
				Name:  modelsRacer.RacerName,
				Place: modelsRacer.RacerPlace,
			}
			race.Racers[modelsRacer.RacerName] = racer
		}

		// Pretend like this race just finished
		leaderboardUpdateTrueSkill(race)
	}

	// Fix the "Date of Last Race" column
	if err := db.Users.SetAllDiversityLastRace(); err != nil {
		log.Error("Database error while resetting the diversity stats:", err)
	}
}
