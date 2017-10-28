package main

import (
	"sort"
	"strconv"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/Zamiell/isaac-racing-server/src/models"
	trueskill "github.com/mafredri/go-trueskill"
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
	// The racer map is not in order, so first make a sorted list of usernames
	racerNames := make([]string, 0)
	for username, _ := range race.Racers {
		racerNames = append(racerNames, username)
	}
	sort.Strings(racerNames)

	// Get the stats for every person in the race
	var statsSlice []*models.StatsDiversity
	for _, racerName := range racerNames {
		racer := race.Racers[racerName]
		if v, err := db.Users.GetStatsDiversity(racer.ID); err != nil {
			log.Error("Database error while getting the diversity stats for \""+racer.Name+"\":", err)
			return
		} else {
			// Increment the number of races
			v.NumRaces++

			// "NewTrueSkill" defaults to 0
			// We need to keep track of the old and new TrueSkill so that we can calculate the change from the entire race
			v.NewTrueSkill = v.TrueSkill
			statsSlice = append(statsSlice, &v)
		}
	}

	// Do a 1v1 TrueSkill calculation for everyone in the race
	ts := trueskill.New()
	for i, racer1 := range racerNames {
		p1stats := statsSlice[i]

		for j, racer2 := range racerNames {
			// Skip this player if we are facing ourself or we have already done the matchup
			if j <= i {
				continue
			}

			p2stats := statsSlice[j]

			// Based on code from:
			// https://godoc.org/github.com/mafredri/go-trueskill
			tsPlayers := make([]trueskill.Player, 0)
			p1 := trueskill.NewPlayer(p1stats.NewTrueSkill, p1stats.Sigma)
			p2 := trueskill.NewPlayer(p2stats.NewTrueSkill, p2stats.Sigma)
			if race.Racers[racer1].Place < race.Racers[racer2].Place {
				// Player 1 wins
				tsPlayers = append(tsPlayers, p1)
				tsPlayers = append(tsPlayers, p2)

				// The second argument is "draw"
				newTsPlayers, _ := ts.AdjustSkills(tsPlayers, false)

				p1stats.NewTrueSkill = newTsPlayers[0].Mu()
				p1stats.Sigma = newTsPlayers[0].Sigma()
				p2stats.NewTrueSkill = newTsPlayers[1].Mu()
				p2stats.Sigma = newTsPlayers[1].Sigma()
			} else if race.Racers[racer1].Place > race.Racers[racer2].Place {
				// Player 2 wins
				tsPlayers = append(tsPlayers, p2)
				tsPlayers = append(tsPlayers, p1)

				// The second argument is "draw"
				newTsPlayers, _ := ts.AdjustSkills(tsPlayers, false)

				p2stats.NewTrueSkill = newTsPlayers[0].Mu()
				p2stats.Sigma = newTsPlayers[0].Sigma()
				p1stats.NewTrueSkill = newTsPlayers[1].Mu()
				p1stats.Sigma = newTsPlayers[1].Sigma()
			} else {
				// Player 1 and 2 tied; this can only happen if both players quit
				tsPlayers = append(tsPlayers, p1)
				tsPlayers = append(tsPlayers, p2)

				// The second argument is "draw"
				newTsPlayers, _ := ts.AdjustSkills(tsPlayers, true)

				p1stats.NewTrueSkill = newTsPlayers[0].Mu()
				p1stats.Sigma = newTsPlayers[0].Sigma()
				p2stats.NewTrueSkill = newTsPlayers[1].Mu()
				p2stats.Sigma = newTsPlayers[1].Sigma()
			}
		}
	}

	for i, racerName := range racerNames {
		// Update the "TrueSkill Change" values for everyone
		stats := statsSlice[i]
		stats.Change = stats.NewTrueSkill - stats.TrueSkill

		// Write the values back to the database
		racer := race.Racers[racerName]
		if err := db.Users.SetStatsDiversity(racer.ID, *stats); err != nil {
			log.Error("Database error while setting the diversity stats for user "+strconv.Itoa(racer.ID)+":", err)
		}
	}
}

func leaderboardDiversityRecalculate() {
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
		leaderboardUpdateDiversity(race)
	}
}
