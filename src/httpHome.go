package main

import (
	"github.com/gin-gonic/gin"
)

/*
	The landing page for the website
*/

func httpHome(c *gin.Context) {
	// Local variables
	w := c.Writer

	data := TemplateData{
		Title: "Home",
	}
	httpServeTemplate(w, "home", data)
}
