package main

import (
	"net/http"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/gin-gonic/gin"
)

/*
	Data types
*/

/*
type LeaderboardSeeded []models.LeaderboardRowSeeded

func (l LeaderboardSeeded) Len() int {
	return len(s)
}
func (l LeaderboardSeeded) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (l LeaderboardSeeded) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

type LeaderboardUnseeded []models.LeaderboardRowUnseeded

*/

/*
	HTTP handlers
*/

func httpLeaderboards(c *gin.Context) {
	// Local variables
	w := c.Writer

	unseededRacesNeeded := 5
	unseededRacesLimit := 1000
	unseededSoloRacesNeeded := 20
	unseededSoloRacesLimit := 1000
	diversityRacesNeeded := 10
	diversityRacesLimit := 1000
	/*
		leaderboardSeeded, err := db.Users.GetLeaderboardSeeded()
		if err != nil {
			log.Error("Failed to get the seeded leaderboard:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	*/

	leaderboardUnseeded, err := db.Users.GetLeaderboardUnseeded(unseededRacesNeeded, unseededRacesLimit)
	if err != nil {
		log.Error("Failed to get the unseeded leaderboard:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	leaderboardUnseededSolo, err := db.Users.GetLeaderboardUnseededSolo(unseededSoloRacesNeeded, unseededSoloRacesLimit)
	if err != nil {
		log.Error("Failed to get the unseeded solo leaderboard:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	leaderboardDiversity, err := db.Users.GetLeaderboardDiversity(diversityRacesNeeded, diversityRacesLimit)
	if err != nil {
		log.Error("Failed to get the diversity leaderboard:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Construct the "Top 10 Unseeded Times" leaderboard
	/*var leaderboardTop10Times string
	for _, row := range leaderboardUnseeded {

	}*/

	// Construct the "Most Races Played" leaderboard
	// TODO

	data := TemplateData{
		Title:                   "Leaderboards",
		LeaderboardUnseeded:     leaderboardUnseeded,
		LeaderboardUnseededSolo: leaderboardUnseededSolo,
		LeaderboardDiversity:    leaderboardDiversity,
	}

	httpServeTemplate(w, "leaderboards", data)
}
