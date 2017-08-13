package main

import (
	"github.com/gin-gonic/gin"
)

func httpDownload(c *gin.Context) {
	// Local variables
	w := c.Writer

	data := TemplateData{
		Title: "Download",
	}
	httpServeTemplate(w, "download", data)
}
