package main

/*
 *  Imports
 */

import (
	/*"github.com/tdewolff/minify" // For minification of HTTP files
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"*/
	"html/template"
	//"io/ioutil"
	"net/http"
	"os"
	"path"
	//"path/filepath"
)

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	// Prepare the layout path and the requested controller path
	lp := path.Join("views", "layout.tmpl")
	if r.URL.Path == "/" {
		r.URL.Path = "/home"
	}
	fp := path.Join("views", r.URL.Path+".tmpl")

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
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Prepare the data for the template
	var data interface{}
	data = nil

	// Execute the template and send it to the user
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
