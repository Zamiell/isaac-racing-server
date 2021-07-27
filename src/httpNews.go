package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func httpNews(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "https://github.com/Zamiell/isaac-racing-client/blob/master/docs/HISTORY.md")
}
