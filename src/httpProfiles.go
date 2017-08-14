package main

import (
	"math"
	"net/http"
	"strconv"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/gin-gonic/gin"
)

func httpProfiles(c *gin.Context) {
	// Local variables
	w := c.Writer
	r := c.Request

	var currentPage int
	// Hard-coded for now, maybe will change this in the future allowing # of results per page
	usersPerPage := 20
	// Find what page we're currently on and then set it accordingly (always set to 1 otherwise)
	i, err := strconv.ParseInt(r.URL.Query().Get(":page"), 10, 32)
	if err == nil && int(i) > 1 {
		currentPage = int(i)
	} else {
		currentPage = 1
	}

	// Get profile data from the database
	userProfiles, totalProfileCount, err := db.Users.GetUserProfiles(currentPage, usersPerPage)
	if err != nil {
		log.Error("Failed to get the user profile data: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	totalPages := math.Floor(float64(totalProfileCount) / float64(usersPerPage))

	// Data to pass to the template, some of it may not be used due to changes
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
	r := c.Request

	// Parse the player name from the URL
	var player string
	player = r.URL.Query().Get(":player")
	if player == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	playerData, err := db.Users.GetProfileData(player)
	if err != nil {
		log.Error("Failed to get player data from the database: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	httpServeTemplate(w, "profile", TemplateData{
		Title:          "Profile",
		ResultsProfile: playerData,
	})
}
