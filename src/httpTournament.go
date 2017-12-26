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

	// Parse the player name from the URL
	tournamentName := c.Params.ByName("tourn_name")
	if len(tournamentName) == 0 {
		httpServeTemplate(w, "notournament", TemplateData{Title: "Tournaments"})
		return
	}

	// This can be removed once it's figured out.
	tournamentInfo := TournamentStats{
		tournamentName,
		"UnhingedMoose",
		"03/20/2018",
		"05/01/2018",
		14,
	}

	// Get the tournament race data
	tournamentRaces, err := db.Tournament.GetTournamentRaces(tournamentName)
	if err != nil {
		log.Error("Failed to get tournament, '" + tournamentName + "' data from the database: " + err.Error())
		httpServeTemplate(w, "notournament", TemplateData{Title: "Tournaments", TournamentInfos: tournamentInfo})
		return
	}

	// Set data to serve to the template
	data := TemplateData{
		Title:           "Tournaments",
		TournamentInfos: tournamentInfo,
		TournamentRaces: tournamentRaces,
	}

	httpServeTemplate(w, "tournament", data)
}
