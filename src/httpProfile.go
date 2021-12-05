package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func httpProfile(c *gin.Context) {
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
	raceDataSoloRanked, err := db.Races.GetSoloRankedRaceProfileHistory(player, NumUnseededRacesForAverage)
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

	for i := range raceDataSoloRanked {
		raceDataSoloRanked[i].RaceFormat.String = strings.Title(raceDataSoloRanked[i].RaceFormat.String)
		for p := range raceDataSoloRanked[i].RaceParticipants {
			raceDataSoloRanked[i].RaceParticipants[p].RacerStartingItemName = allItemNames[int(raceDataSoloRanked[i].RaceParticipants[p].RacerStartingItem.Int64)]
			if raceDataSoloRanked[i].RaceParticipants[p].RacerStartingBuild.Int64 > 0 {
				startingBuildIndex := int(raceDataSoloRanked[i].RaceParticipants[p].RacerStartingBuild.Int64)
				raceDataSoloRanked[i].RaceParticipants[p].RacerStartingBuildName = getBuildName(startingBuildIndex)
				raceDataSoloRanked[i].RaceParticipants[p].RacerStartingBuildID = getBuildID(startingBuildIndex)
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
		RaceResultsRanked: raceDataSoloRanked,
		RaceResultsAll:    raceDataAll,
	}

	httpServeTemplate(w, "profile", data)
}
