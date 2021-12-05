package server

import (
	"github.com/gin-gonic/gin"
)

func httpInfo(c *gin.Context) {
	w := c.Writer

	data := TemplateData{
		Title: "Info",
	}
	httpServeTemplate(w, "info", data)
}
