package server

import (
	"github.com/gin-gonic/gin"
)

func httpDownload(c *gin.Context) {
	w := c.Writer

	data := TemplateData{
		Title: "Download",
	}
	httpServeTemplate(w, "download", data)
}
