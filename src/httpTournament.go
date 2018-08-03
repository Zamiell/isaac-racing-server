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
	currentTournament := true
	// Get the tournament race db data or serve a blank page if no tournaments found
	tournamentRaces, err := db.Tournament.GetTournamentRaces()
	if err != nil {
		log.Error("No current tournaments being tracked: " + err.Error())
		//httpServeTemplate(w, "notournament", TemplateData{Title: "No Tournament", AllTournaments: allTournaments})
		//return
		currentTournament = false
	}
	log.Info(currentTournament)
	// Set data to serve to the template
	data := TemplateData{
		Title:           "Tournaments",
		CurrentTourn:    currentTournament,
		TournamentRaces: tournamentRaces,
		AllTournaments:  allTournaments,
	}

	httpServeTemplate(w, "tournament", data)
}
