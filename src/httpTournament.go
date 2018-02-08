package main

import (
	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/gin-gonic/gin"
	"strings"
)

type TournamentStats struct {
	TournamentName      string
	TournamentHost      string
	TournamentStartDate string
	TournamentEndDate   string
	TournamentPlayers   int
}

func httpTournament(c *gin.Context) {
	// Local variables
	w := c.Writer

	// Get the tournament race db data or serve a blank page if no tournaments found
	tournamentRaces, err := db.Tournament.GetTournamentRaces()
	if err != nil {
		log.Error("Failed to get tournament data from the database: " + err.Error())
		httpServeTemplate(w, "notournament", TemplateData{Title: "No Tournament"})
		return
	}

	// Sets all the TournamentName strings to lowercase because people sometimes can't format
	for i := range tournamentRaces {
		tournamentRaces[i].TournamentName.String = strings.ToLower(tournamentRaces[i].TournamentName.String)
	}

	// Set data to serve to the template
	data := TemplateData{
		Title:           "Tournaments",
		TournamentRaces: tournamentRaces,
	}

	httpServeTemplate(w, "tournament", data)
}
