package main

/*
	Imports
*/

import (
	"math"
	"net/http"
	"strconv"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/gin-gonic/gin"
)

/*
	HTTP handlers
*/

func httpRaces(c *gin.Context) {
	// Local variables
	w := c.Writer
	r := c.Request
	var currentPage int
	racesPerPage := 20

	i, err := strconv.ParseInt(r.URL.Query().Get(":page"), 10, 32)
	if err == nil && int(i) > 1 {
		currentPage = int(i)
	} else {
		currentPage = 1
	}
	raceData, totalRaces, err := db.Races.GetRaceHistory(currentPage, racesPerPage)
	if err != nil {
		log.Error("Failed to get the race data: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	totalPages := math.Floor(float64(totalRaces) / float64(racesPerPage))
	data := TemplateData{
		Title:          "Races",
		ResultsRaces:   raceData,
		TotalRaceCount: totalRaces,
		TotalPages:     int(totalPages),
		PreviousPage:   currentPage - 1,
		NextPage:       currentPage + 1,
	}

	httpServeTemplate(w, "races", data)
}
