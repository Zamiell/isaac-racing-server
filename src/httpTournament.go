package main

import (
	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/gin-gonic/gin"
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

	// Get the tournament race data or serve a blank page if no tournaments found
	tournamentRaces, err := db.Tournament.GetTournamentRaces()
	if err != nil {
		log.Error("Failed to get tournament data from the database: " + err.Error())
		httpServeTemplate(w, "notournament", TemplateData{Title: "Tournaments"})
		return
	}
	log.Info(tournamentRaces)
	// Set data to serve to the template
	data := TemplateData{
		Title:           "Tournaments",
		TournamentRaces: tournamentRaces,
	}

	httpServeTemplate(w, "tournament", data)
}
