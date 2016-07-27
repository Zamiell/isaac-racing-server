package main

/*
 *  Imports
 */

import (
	"html/template"
	"net/http"
	"os"
)

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := projectPath + "/views/layout.tmpl"
	if r.URL.Path == "/" {
		r.URL.Path = "/home"
	}
	fp := projectPath + "/views/" + r.URL.Path + ".tmpl"

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
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
	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		log.Error("Failed to execute the template:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
