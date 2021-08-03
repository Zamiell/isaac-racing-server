package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func httpProfile(c *gin.Context) {
	// Local variables
	w := c.Writer
	racesAllTotal := 5

	// Parse the player name from the URL
	player := c.Params.ByName("player")
	if player == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Check if the player exists
	var playerID int
	if exists, v, err := db.Users.Exists(player); err != nil {
		logger.Error("Failed to check if player \"" + player + "\" exists: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if !exists {
		data := TemplateData{
			Title:         "Profile Missing",
			MissingPlayer: player,
		}
		httpServeTemplate(w, "noprofile", data)
		return
	} else {
		playerID = v
	}

	// Get the player data
	playerData, err := db.Users.GetProfileData(playerID)
	if err != nil {
		logger.Error("Failed to get player, '" + player + "' data from the database: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	totalTime, err := db.Users.GetTotalTime(player)
	if err != nil {
		logger.Error("Failed to get player time from database: ", totalTime, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Get the race data for the last x races
	raceDataRanked, err := db.Races.GetRankedRaceProfileHistory(player, numUnseededRacesForAverage)
	if err != nil {
		logger.Error("Failed to get the race data: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	raceDataAll, err := db.Races.GetAllRaceProfileHistory(player, racesAllTotal)
	if err != nil {
		logger.Error("Failed to get the race data: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for i := range raceDataRanked {
		raceDataRanked[i].RaceFormat.String = strings.Title(raceDataRanked[i].RaceFormat.String)
		for p := range raceDataRanked[i].RaceParticipants {
			raceDataRanked[i].RaceParticipants[p].RacerStartingItemName = allItemNames[int(raceDataRanked[i].RaceParticipants[p].RacerStartingItem.Int64)]
			if raceDataRanked[i].RaceParticipants[p].RacerStartingBuild.Int64 > 0 {
				startingBuildIndex := int(raceDataRanked[i].RaceParticipants[p].RacerStartingBuild.Int64)
				raceDataRanked[i].RaceParticipants[p].RacerStartingBuildName = getBuildName(startingBuildIndex)
				raceDataRanked[i].RaceParticipants[p].RacerStartingBuildID = getBuildID(startingBuildIndex)
			}
		}
	}
	for i := range raceDataAll {
		raceDataAll[i].RaceFormat.String = strings.Title(raceDataAll[i].RaceFormat.String)
		for p := range raceDataAll[i].RaceParticipants {
			raceDataAll[i].RaceParticipants[p].RacerStartingItemName = allItemNames[int(raceDataAll[i].RaceParticipants[p].RacerStartingItem.Int64)]
			if raceDataAll[i].RaceParticipants[p].RacerStartingBuild.Int64 > 0 {
				startingBuildIndex := int(raceDataAll[i].RaceParticipants[p].RacerStartingBuild.Int64)
				raceDataAll[i].RaceParticipants[p].RacerStartingBuildName = getBuildName(startingBuildIndex)
				raceDataAll[i].RaceParticipants[p].RacerStartingBuildID = getBuildID(startingBuildIndex)
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
