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
			log.Error("Database error while getting the unseeded times:", err)
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

		var averageTime float64
		var forfeitPenalty float64
		if len(unseededTimes) == numForfeits {
			// If they forfeited every race, then we will have a divide by 0 later on,
			// so arbitrarily set it to 30 minutes (1000 * 60 * 30)
			averageTime = 1800000
			forfeitPenalty = 1800000
		} else {
			averageTime = float64(sumTimes) / float64(len(unseededTimes)-numForfeits)
			forfeitPenalty = averageTime * float64(numForfeits) / float64(len(unseededTimes))
		}

		// Update their stats in the database
		if err := db.Users.SetStatsUnseeded(racer.ID, int(averageTime), numForfeits, int(forfeitPenalty)); err != nil {
			log.Error("Database error while setting the unseeded stats for \""+racer.Name+"\":", err)
			return
		}
	}
}

func leaderboardUpdateDiversity(race *Race) {
	// Update the stats for every person in the race
	for _, racer := range race.Racers {
		var stats models.StatsDiversity
		if v, err := db.Users.GetStatsDiversity(racer.ID); err != nil {
			log.Error("Database error while getting the diversity stats for \""+racer.Name+"\":", err)
			return
		} else {
			stats = v
		}

		log.Info(stats)
	}
}
