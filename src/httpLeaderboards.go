package server

import (
	"net/http"

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
	seededRacesNeeded := 5
	unseededRacesNeeded := 5
	unseededSoloRacesNeeded := 20
	diversityRacesNeeded := 10

	leaderboardSeeded, err := db.Users.GetLeaderboardSeeded(seededRacesNeeded)
	if err != nil {
		logger.Error("Failed to get the seeded leaderboard:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	leaderboardUnseeded, err := db.Users.GetLeaderboardUnseeded(unseededRacesNeeded)
	if err != nil {
		logger.Error("Failed to get the unseeded leaderboard:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	leaderboardDiversity, err := db.Users.GetLeaderboardDiversity(diversityRacesNeeded)
	if err != nil {
		logger.Error("Failed to get the diversity leaderboard:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	leaderboardRankedSolo, err := db.Users.GetLeaderboardRankedSolo(unseededSoloRacesNeeded)
	if err != nil {
		logger.Error("Failed to get the ranked solo leaderboard:", err)
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
		Title:                 "Leaderboards",
		LeaderboardSeeded:     leaderboardSeeded,
		LeaderboardUnseeded:   leaderboardUnseeded,
		LeaderboardDiversity:  leaderboardDiversity,
		LeaderboardRankedSolo: leaderboardRankedSolo,
	}

	httpServeTemplate(w, "leaderboards", data)
}
