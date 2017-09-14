package main

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/gin-gonic/gin"
)

func httpProfiles(c *gin.Context) {
	// Local variables
	w := c.Writer
	currentPage := 1
	usersPerPage := 20

	// Find what page we're currently on and then set it accordingly (always set to 1 otherwise)
	i, err := strconv.ParseInt(c.Params.ByName("page"), 10, 32)
	if err == nil && int(i) > 1 {
		currentPage = int(i)
	}

	// Get profile data from the database
	userProfiles, totalProfileCount, err := db.Users.GetUserProfiles(currentPage, usersPerPage)

	if err != nil {
		log.Error("Failed to get the user profile data: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Get total number of pages needed for navigation
	totalPages := math.Floor(float64(totalProfileCount) / float64(usersPerPage))

	// Data to pass to the template
	data := TemplateData{
		Title:             "Profiles",
		ResultsProfiles:   userProfiles,
		TotalProfileCount: totalProfileCount,
		TotalPages:        int(totalPages),
		PreviousPage:      currentPage - 1,
		NextPage:          currentPage + 1,
		UsersPerPage:      usersPerPage,
	}
	httpServeTemplate(w, "profiles", data)
}

func httpProfile(c *gin.Context) {
	// Local variables
	w := c.Writer
	racesPerPage := 5

	// Parse the player name from the URL
	player := c.Params.ByName("player")
	if player == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Get the player data
	playerData, err := db.Users.GetProfileData(player)
	if err != nil {
		log.Error("Failed to get player data from the database: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Get the race data for the last x races
	raceData, err := db.Races.GetRaceProfileHistory(player, racesPerPage)
	if err != nil {
		log.Error("Failed to get the race data: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Capitalize the RaceFormat data
	for i := range raceData {
		raceData[i].RaceFormat = strings.Title(raceData[i].RaceFormat)
	}

	// Set data to serve to the template
	data := TemplateData{
		Title:          "Profile",
		ResultsProfile: playerData,
		RaceResults:    raceData,
	}

	httpServeTemplate(w, "profile", data)
}
