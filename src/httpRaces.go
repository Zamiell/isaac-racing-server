package main

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/gin-gonic/gin"
)

func httpRaces(c *gin.Context) {
	// Local variables
	w := c.Writer
	currentPage := 1
	racesPerPage := 20
	raceOffset := 0

	// Grab the current page from the URI and set currentPage if found
	i, err := strconv.ParseInt(c.Params.ByName("page"), 10, 32)
	if err == nil && int(i) > 1 {
		currentPage = int(i)
		raceOffset = (racesPerPage * currentPage) + 1
	}

	raceData, totalRaces, err := db.Races.GetRacesHistory(currentPage, racesPerPage, raceOffset)
	if err != nil {
		log.Error("Failed to get the race data: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// Find total number of pages needed for navigation, if divisible by perPage, remove a page.
	totalPages := 0
	if totalRaces%racesPerPage == 0 {
		totalPages = int(math.Floor(float64(totalRaces)/float64(racesPerPage)) - 1)
	} else {
		totalPages = int(math.Floor(float64(totalRaces) / float64(racesPerPage)))
	}

	// Capitalize the RaceFormat data
	for i := range raceData {
		raceData[i].RaceFormat.String = strings.Title(raceData[i].RaceFormat.String)
		for p := range raceData[i].RaceParticipants {
			raceData[i].RaceParticipants[p].RacerStartingItemName = allItemNames[int(raceData[i].RaceParticipants[p].RacerStartingItem.Int64)].Name
		}
	}

	// Build template data for serving to the template
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

func httpRace(c *gin.Context) {
	// Local variables
	w := c.Writer

	raceId, err := strconv.ParseInt(c.Params.ByName("raceid"), 10, 32)
	if err != nil {
		log.Error("Failed to parse the url for raceId: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	raceData, err := db.Races.GetRaceHistory(int(raceId))
	if err != nil {
		log.Error("Failed to get the race data: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for i := range raceData {
		raceData[i].RaceFormat.String = strings.Title(raceData[i].RaceFormat.String)
		for p := range raceData[i].RaceParticipants {
			raceData[i].RaceParticipants[p].RacerStartingItemName = allItemNames[int(raceData[i].RaceParticipants[p].RacerStartingItem.Int64)].Name
		}
	}

	data := TemplateData{
		Title:       "Race",
		RaceResults: raceData,
	}

	httpServeTemplate(w, "race", data)
}
