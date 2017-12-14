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
	// Update the stats for every person in the race
	/*
		for _, racer := range race.Racers {
		}
	*/
}

func leaderboardUpdateSoloSeeded(race *Race) {
	// Update the stats for every person in the race
	/*
		for _, racer := range race.Racers {
		}
	*/
}

func leaderboardUpdateUnseeded(race *Race) {
	// Update the stats for every person in the race
	/*
		for _, racer := range race.Racers {
		}
	*/
}

func leaderboardUpdateSoloUnseeded(race *Race) {
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
		if err := db.Users.SetStatsSoloUnseeded(racer.ID, int(averageTime), numForfeits, int(forfeitPenalty)); err != nil {
			log.Error("Database error while setting the unseeded stats for \""+racer.Name+"\":", err)
			return
		}
	}
}

func leaderboardUpdateTrueSkill(race *Race) {
	// The racer map is not in order, so first make a sorted list of usernames
	racerNames := make([]string, 0)
	for username, _ := range race.Racers {
		racerNames = append(racerNames, username)
	}
	sort.Strings(racerNames)

	// Get the stats for every person in the race
	var statsSlice []*models.StatsTrueSkill
	for _, racerName := range racerNames {
		racer := race.Racers[racerName]
		if v, err := db.Users.GetTrueSkill(racer.ID, race.Ruleset.Format); err != nil {
			log.Error("Database error while getting the TrueSkill stats for \""+racer.Name+"\":", err)
			return
		} else {
			// Increment the number of races
			v.NumRaces++

			statsSlice = append(statsSlice, &v)
		}
	}

	// Do a 1v1 TrueSkill calculation for everyone in the race
	for i, racer1 := range racerNames {
		p1stats := statsSlice[i]
		p1Place := race.Racers[racer1].Place
		if p1Place == -1 {
			// Change forfeits to place 999 to simplify the below calculation
			p1Place = 999
		} else if p1Place == -2 {
			// Change disqualifications to place 1000 to simplify the below calculation
			p1Place = 1000
		}

		for j, racer2 := range racerNames {
			// Skip this player if they are facing themself or we have already done the matchup
			if j <= i {
				continue
			}

			p2stats := statsSlice[j]
			p2Place := race.Racers[racer2].Place
			if p2Place == -1 {
				// Change forfeits to place 999 to simplify the below calculation
				p2Place = 999
			} else if p2Place == -2 {
				// Change disqualifications to place 1000 to simplify the below calculation
				p2Place = 1000
			}

			var p1Mu, p1Sigma, p2Mu, p2Sigma float64
			if p1Place < p2Place {
				// Player 1 wins
				p1Mu, p1Sigma, p2Mu, p2Sigma = leaderboardAdjustTrueSkill(
					p1stats.Mu,
					p1stats.Sigma,
					p2stats.Mu,
					p2stats.Sigma,
					false,
				)

				log.Info("Race #" + strconv.Itoa(race.ID) + ": player \"" + racer1 + "\" (place " + strconv.Itoa(race.Racers[racer1].Place) + ") WINS OVER player \"" + racer2 + "\" (place " + strconv.Itoa(race.Racers[racer2].Place) + ")")
			} else if p1Place > p2Place {
				// Player 2 wins
				p2Mu, p2Sigma, p1Mu, p1Sigma = leaderboardAdjustTrueSkill(
					p2stats.Mu,
					p2stats.Sigma,
					p1stats.Mu,
					p1stats.Sigma,
					false,
				)
				log.Info("Race #" + strconv.Itoa(race.ID) + ": player \"" + racer1 + "\" (place " + strconv.Itoa(race.Racers[racer1].Place) + ") LOSES TO player \"" + racer2 + "\" (place " + strconv.Itoa(race.Racers[racer2].Place) + ")")
			} else {
				// Player 1 and 2 tied; this can only happen if both players quit
				// (or they were both disqualified)
				p1Mu, p1Sigma, p2Mu, p2Sigma = leaderboardAdjustTrueSkill(
					p1stats.Mu,
					p1stats.Sigma,
					p2stats.Mu,
					p2stats.Sigma,
					true,
				)
				log.Info("Race #" + strconv.Itoa(race.ID) + ": player \"" + racer1 + "\" (place " + strconv.Itoa(race.Racers[racer1].Place) + ") TIES player \"" + racer2 + "\" (place " + strconv.Itoa(race.Racers[racer2].Place) + ")")
			}

			p1stats.Mu = p1Mu
			p1stats.Sigma = p1Sigma
			p2stats.Mu = p2Mu
			p2stats.Sigma = p2Sigma
		}
	}

	for i, racerName := range racerNames {
		// Get the player's new "TrueSkill" and the change
		stats := statsSlice[i]
		trueSkill := leaderboardGetTrueSkill(stats.Mu, stats.Sigma)
		stats.Change = stats.TrueSkill - trueSkill
		stats.TrueSkill = trueSkill

		// Write the values back to the database
		racer := race.Racers[racerName]
		if err := db.Users.SetTrueSkill(racer.ID, *stats, race.Ruleset.Format); err != nil {
			log.Error("Database error while setting the TrueSkill stats for user "+strconv.Itoa(racer.ID)+":", err)
		}
	}
}

/*
	Subroutines
*/

func leaderboardAdjustTrueSkill(p1Mu float64, p1Sigma float64, p2Mu float64, p2Sigma float64, draw bool) (float64, float64, float64, float64) {
	// Based on code from:
	// https://godoc.org/github.com/mafredri/go-trueskill
	ts := trueskill.New()
	p1 := trueskill.NewPlayer(p1Mu, p1Sigma)
	p2 := trueskill.NewPlayer(p2Mu, p2Sigma)
	tsPlayers := []trueskill.Player{p1, p2} // The first player that is put into the "tsPlayers" slice is the one who wins
	newTsPlayers, _ := ts.AdjustSkills(tsPlayers, draw)

	return newTsPlayers[0].Mu(), newTsPlayers[0].Sigma(), newTsPlayers[1].Mu(), newTsPlayers[1].Sigma()
}

func leaderboardGetTrueSkill(mu float64, sigma float64) float64 {
	// Based on code from:
	// https://godoc.org/github.com/mafredri/go-trueskill

	ts := trueskill.New()
	player := trueskill.NewPlayer(mu, sigma)
	return ts.TrueSkill(player)
}
