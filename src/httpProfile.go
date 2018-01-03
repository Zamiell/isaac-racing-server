package main

import (
	"net/http"
	"strings"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/gin-gonic/gin"
)

func httpProfile(c *gin.Context) {
	// Local variables
	w := c.Writer
	racesRankedTotal := 50
	racesAllTotal := 5

	// Parse the player name from the URL
	player := c.Params.ByName("player")
	if player == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Check if the player exists
	exists, playerID, err := db.Users.Exists(player)
	if !exists {
		data := TemplateData{
			Title:         "Profile Missing",
			MissingPlayer: player,
		}
		httpServeTemplate(w, "noprofile", data)
		return
	}

	// Get the player data
	playerData, err := db.Users.GetProfileData(playerID)
	if err != nil {
		log.Error("Failed to get player, '" + player + "' data from the database: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	totalTime, err := db.Users.GetTotalTime(player)
	if err != nil {
		log.Error("Failed to get player time from database: ", totalTime, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Get the race data for the last x races
	raceDataRanked, err := db.Races.GetRankedRaceProfileHistory(player, racesRankedTotal)
	if err != nil {
		log.Error("Failed to get the race data: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	raceDataAll, err := db.Races.GetAllRaceProfileHistory(player, racesAllTotal)
	if err != nil {
		log.Error("Failed to get the race data: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Do special things on the race things with fancy things
	for i := range raceDataRanked {
		raceDataRanked[i].RaceFormat.String = strings.Title(raceDataRanked[i].RaceFormat.String)
		for p := range raceDataRanked[i].RaceParticipants {
			raceDataRanked[i].RaceParticipants[p].RacerStartingItemName = allItemNames[int(raceDataRanked[i].RaceParticipants[p].RacerStartingItem.Int64)]
			if raceDataRanked[i].RaceParticipants[p].RacerStartingBuild.Int64 > 0 {
				raceDataRanked[i].RaceParticipants[p].RacerStartingBuildName = seededBuilds[raceDataRanked[i].RaceParticipants[p].RacerStartingBuild.Int64-1].Name
			}
		}
	}
	for i := range raceDataAll {
		raceDataAll[i].RaceFormat.String = strings.Title(raceDataAll[i].RaceFormat.String)
		for p := range raceDataAll[i].RaceParticipants {
			raceDataAll[i].RaceParticipants[p].RacerStartingItemName = allItemNames[int(raceDataAll[i].RaceParticipants[p].RacerStartingItem.Int64)]
			if raceDataAll[i].RaceParticipants[p].RacerStartingBuild.Int64 > 0 {
				raceDataAll[i].RaceParticipants[p].RacerStartingBuildName = seededBuilds[raceDataAll[i].RaceParticipants[p].RacerStartingBuild.Int64-1]
			}
		}
	}

	// Set data to serve to the template
	data := TemplateData{
		Title:             "Profile",
		ResultsProfile:    playerData,
		TotalTime:         totalTime,
		RaceResultsRanked: raceDataRanked,
		RaceResultsAll:    raceDataAll,
	}

	httpServeTemplate(w, "profile", data)
}
