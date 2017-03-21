package main

/*
	Imports
*/

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"math"
	"strconv"
	"github.com/Zamiell/isaac-racing-server/models"
)

/*
	Data types
*/

type TemplateData struct {
	Title	string
}

type TemplateDataProfiles struct {
	Title				string
	Results				[]models.UserProfilesRow
	TotalProfileCount	int
	TotalPages			int
	PreviousPage		int
	NextPage			int
	UsersPerPage		int
}

type TemplateDataProfile struct {
	Title		string
	Results		models.UserProfileData
}

/*
	Main page handlers
*/

func httpHome(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title: "Home",
	}
	serveTemplate(w, "home", data)
}

func httpNews(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title: "News",
	}
	serveTemplate(w, "news", data)
}

func httpRaces(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title: "Races",
	}
	serveTemplate(w, "races", data)
}

func httpProfile (w http.ResponseWriter, r *http.Request) {
	// Get player from url
	var player string
	player = r.URL.Query().Get(":player")
	if player == "" {
		player = "Zamiell"
		log.Error("Failed to a parse the player data: ", player)
	
	}
	// Get the data from the database
	playerData, err := db.Users.GetProfileData(player)
	if err != nil {
		log.Error("Failed to get player data from the database: ", err)
	}
	// Create the title with player's name
	data := TemplateDataProfile{
		Title: "Profile",
		Results: playerData,
	}
	serveTemplate(w, "profile", data)
}
func httpProfiles(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, http.StatusText(http.StatusInternalServerError),  http.StatusInternalServerError)
		return
	}
	totalPages := math.Ceil(float64(totalProfileCount) / float64(usersPerPage))
	// Data to pass to the template, some of it may not be used due to changes
	data := TemplateDataProfiles{
		Title: "Profiles",
		Results: userProfiles,
		TotalProfileCount: totalProfileCount,
		TotalPages: int(totalPages),
		PreviousPage: currentPage - 1,
		NextPage: currentPage + 1,
		UsersPerPage: usersPerPage,
	}
	serveTemplate(w, "profiles", data)
}

func httpLeaderboards(w http.ResponseWriter, r *http.Request) {
	/*leaderboardSeeded, err := db.Users.GetLeaderboardSeeded()
	if err != nil {
		log.Error("Failed to get the seeded leaderboard:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	leaderboardUnseeded, err := db.Users.GetLeaderboardUnseeded()
	if err != nil {
		log.Error("Failed to get the unseeded leaderboard:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}*/

	// Construct the "Top 10 Unseeded Times" leaderboard
	/*var leaderboardTop10Times string
	for _, row := range leaderboardUnseeded {

	}*/

	// Construct the "Most Races Played" leaderboard

	data := TemplateData{
		Title: "Leaderboards",
	}
	serveTemplate(w, "leaderboards", data)
}

func httpInfo(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title: "Info",
	}
	serveTemplate(w, "info", data)
}

func httpDownload(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title: "Download",
	}
	serveTemplate(w, "download", data)
}

/*
	HTTP miscellaneous subroutines
*/

func serveTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	lp := path.Join("views", "layout.tmpl")
	fp := path.Join("views", templateName+".tmpl")

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Create the template
	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		log.Error("Failed to create the template:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Execute the template and send it to the user
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Error("Failed to execute the template:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}