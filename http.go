package main

/*
	Imports
*/

import (
	"html/template"
	"net/http"
	"os"
	"path"
)

/*
	Data types
*/

type TemplateData struct {
	Title string
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

func httpProfiles(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title: "Profiles",
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

func getPlayers() {

}
