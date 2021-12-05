package server

import (
	"github.com/gin-gonic/gin"
)

/*
	The landing page for the website
*/

func httpHome(c *gin.Context) {
	w := c.Writer

	data := TemplateData{
		Title: "Home",
	}
	httpServeTemplate(w, "home", data)
}
