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
	currentPage := 1
	usersPerPage := 50

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
	totalPages := 0
	// Get total number of pages needed for navigation and remove a page if total is divisible by perPage
	if totalProfileCount%usersPerPage == 0 {
		totalPages = int(math.Floor(float64(totalProfileCount)/float64(usersPerPage)) - 1)
	} else {
		totalPages = int(math.Floor(float64(totalProfileCount) / float64(usersPerPage)))
	}

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
