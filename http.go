package main

/*
 *  Imports
 */

import (
	"html/template"
	"net/http"
)

var indexTemplate = template.Must(template.ParseFiles("views/index.html"))

func httpHome(w http.ResponseWriter, r *http.Request) {
	/*data := &Index{
		Title: "Image gallery",
		Body:  "Welcome to the image gallery.",
	}
	for name, img := range images {
		data.Links = append(data.Links, Link{
			URL:   "/image/" + name,
			Title: img.Title,
		})
	}*/
	if err := indexTemplate.Execute(w, nil); err != nil {
		log.Debug(err)
	}
}
