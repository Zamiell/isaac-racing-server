package server

import (
	"github.com/Zamiell/isaac-racing-server/models"
)

const (
	numUnseededRacesForAverage = 100
)

func leaderboardUpdateSoloSeeded(race *Race) {
	// Update the stats for every person in the race
	/*
		for _, racer := range race.Racers {
		}
	*/
}

func leaderboardUpdateSoloUnseeded(race *Race) {
	// Update the stats for the solo racer
	// (we still have to iterate over the map to get the racer)
	for _, racer := range race.Racers {
		var unseededTimes []models.UnseededTime
		if v, err := db.RaceParticipants.GetNUnseededTimes(racer.ID, numUnseededRacesForAverage); err != nil {
			logger.Error("Database error while getting the unseeded times:", err)
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
		if err := db.Users.SetStatsSoloUnseeded(racer.ID, int(averageTime), numForfeits, int(forfeitPenalty)); err != nil {
			logger.Error("Database error while setting the unseeded stats for \""+racer.Name+"\":", err)
			return
		}
	}
}

/*
func leaderboardRecalculateSoloUnseeded() {
	format := "unseeded_solo"
	if err := db.Users.ResetSoloUnseeded(); err != nil {
		logger.Error("Database error while resetting the unseeded solo stats:", err)
		return
	}

	var allRaces []models.RaceHistory
	if v, err := db.Races.GetAllRaces(format); err != nil {
		logger.Error("Database error while getting all of the races:", err)
		return
	} else {
		allRaces = v
	}

	// Go through every race for this particular format in the database
	for _, modelsRace := range allRaces {
		// Convert the "RaceHistory" struct to a "Race" struct
		race := &Race{}
		race.ID = int(modelsRace.RaceID.Int64)
		race.Ruleset.Format = format
		race.Racers = make(map[string]*Racer)
		for _, modelsRacer := range modelsRace.RaceParticipants {
			racer := &Racer{
				ID:    int(modelsRacer.ID.Int64),
				Name:  modelsRacer.RacerName.String,
				Place: int(modelsRacer.RacerPlace.Int64),
			}
			race.Racers[modelsRacer.RacerName.String] = racer
		}

		// Pretend like this race just finished
		leaderboardUpdateSoloUnseeded(race)
	}

	// Fix the "Date of Last Race" column
	if err := db.Users.SetLastRace(format); err != nil {
		logger.Error("Database error while setting the last race:", err)
		return
	}

	logger.Info("Successfully reset the leaderboard for " + format + ".")
}
*/
