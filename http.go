package main

/*
 *  Imports
 */

import (
	"github.com/tdewolff/minify"      // For minification of HTTP files (1/2)
	"github.com/tdewolff/minify/html" // For minification of HTTP files (2/2)
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func serveTemplate(w http.ResponseWriter, r *http.Request) {
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

	// Create and minify the template
	tmpl, err := template.ParseFiles(lp, fp)
	/*tmpl, err := compileTemplates(lp, fp)
	if err != nil {
		log.Error("Failed to create and minify the template:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}*/
	//tmpl := template.Must(compileTemplates("views/layout.tmpl"))

	// Execute the template and send it to the user
	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		log.Error("Failed to execute the template:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// From: https://github.com/tdewolff/minify
func compileTemplates(filenames ...string) (*template.Template, error) {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)

	var tmpl *template.Template
	for _, filename := range filenames {
		name := filepath.Base(filename)
		if tmpl == nil {
			tmpl = template.New(name)
		} else {
			tmpl = tmpl.New(name)
		}

		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		mb, err := m.Bytes("text/html", b)
		if err != nil {
			return nil, err
		}
		tmpl.Parse(string(mb))
	}
	return tmpl, nil
}
