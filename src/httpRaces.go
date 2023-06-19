package server

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func httpRaces(c *gin.Context) {
	w := c.Writer
	currentPage := 1
	racesPerPage := 20
	raceOffset := 0

	// Grab the current page from the URI and set currentPage if found.
	i, err := strconv.ParseInt(c.Params.ByName("page"), 10, 32)
	if err == nil && int(i) > 1 {
		currentPage = int(i)
		raceOffset = (racesPerPage * currentPage) + 1
	}

	raceData, totalRaces, err := db.Races.GetRacesHistory(currentPage, racesPerPage, raceOffset)
	if err != nil {
		logger.Error("Failed to get the race data: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Find total number of pages needed for navigation, if divisible by perPage, remove a page.
	var totalPages int
	if totalRaces%racesPerPage == 0 {
		totalPages = int(math.Floor(float64(totalRaces)/float64(racesPerPage)) - 1)
	} else {
		totalPages = int(math.Floor(float64(totalRaces) / float64(racesPerPage)))
	}

	for _, raceHistory := range raceData {
		raceHistory.RaceFormat.String = strings.Title(raceHistory.RaceFormat.String)
		for _, raceParticipant := range raceHistory.RaceParticipants {
			raceParticipant.RacerStartingItemName = allItemNames[int(raceParticipant.RacerStartingItem.Int64)]
			if raceParticipant.RacerStartingBuild.Int64 > 0 {
				startingBuildIndex := int(raceParticipant.RacerStartingBuild.Int64)
				raceParticipant.RacerStartingBuildName = getBuildNameFromBuildIndex(startingBuildIndex)
				raceParticipant.RacerStartingCollectibleID = getBuildFirstCollectibleID(startingBuildIndex)
			}
		}
	}

	// Build template data for serving to the template
	data := TemplateData{
		Title:          "Races",
		ResultsRaces:   raceData,
		TotalRaceCount: totalRaces,
		TotalPages:     totalPages,
		PreviousPage:   currentPage - 1,
		NextPage:       currentPage + 1,
	}

	httpServeTemplate(w, "races", data)
}

func httpRace(c *gin.Context) {
	w := c.Writer

	raceID, err := strconv.ParseInt(c.Params.ByName("raceid"), 10, 32)
	if err != nil {
		logger.Error("Failed to parse the url for raceId: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	raceData, err := db.Races.GetRaceHistory(int(raceID))
	if err != nil {
		logger.Error("Failed to get the race data: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	raceData.RaceFormat.String = strings.Title(raceData.RaceFormat.String)
	raceFormat := raceData.RaceFormat.String
	for _, raceParticipant := range raceData.RaceParticipants {
		raceParticipant.RacerStartingItemName = allItemNames[int(raceParticipant.RacerStartingItem.Int64)]
		if raceParticipant.RacerStartingBuild.Int64 > 0 {
			startingBuildIndex := int(raceParticipant.RacerStartingBuild.Int64)
			raceParticipant.RacerStartingBuildName = getBuildNameFromBuildIndex(startingBuildIndex)
			raceParticipant.RacerStartingCollectibleID = getBuildFirstCollectibleID(startingBuildIndex)
		}
	}

	data := TemplateData{
		Title:             "Race",
		SingleRaceFormat:  raceFormat,
		SingleRaceResults: raceData,
	}

	httpServeTemplate(w, "race", data)
}
