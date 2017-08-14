package main

import (
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

	/*
		leaderboardSeeded, err := db.Users.GetLeaderboardSeeded()
		if err != nil {
			log.Error("Failed to get the seeded leaderboard:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		leaderboardUnseeded, err := db.Users.GetLeaderboardUnseeded()
		if err != nil {
			log.Error("Failed to get the unseeded leaderboard:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	*/

	// Construct the "Top 10 Unseeded Times" leaderboard
	/*var leaderboardTop10Times string
	for _, row := range leaderboardUnseeded {

	}*/

	// Construct the "Most Races Played" leaderboard
	data := TemplateData{
		Title: "Leaderboards",
	}
	httpServeTemplate(w, "leaderboards", data)
}
