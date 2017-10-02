package main

import (
	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/Zamiell/isaac-racing-server/src/models"
)

const (
	numUnseededRacesForAverage = 50
)

func leaderboardUpdateSeeded(race *Race) {
}

func leaderboardUpdateUnseeded(race *Race) {
	// Update the stats for every person in the race
	for _, racer := range race.Racers {
		var unseededTimes []models.UnseededTime
		if v, err := db.RaceParticipants.GetNUnseededTimes(racer.ID, numUnseededRacesForAverage); err != nil {
			log.Error("Database error:", err)
			return
		} else {
			unseededTimes = v
		}

		var numForfeits int
		var sumTimes int64
		for _, race := range unseededTimes {
			if race.Place > 0 {
				// They finished
				sumTimes += race.RunTime
			} else {
				// They quit
				numForfeits++
			}
		}

		if len(unseededTimes) == numForfeits {
			// If they forfeited every race, then we will have a divide by 0 later on,
			// so arbitrarily set it to a million seconds
		}
	}
}
