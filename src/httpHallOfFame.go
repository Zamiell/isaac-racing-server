package main

import (
	//	"math"
	//	"net/http"
	//	"strconv"
	//	"strings"

	//	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/gin-gonic/gin"
)

func httpHallOfFame(c *gin.Context) {
	// Local variables
	w := c.Writer

	data := TemplateData{
		Title: "Hall Of Fame",
	}
	httpServeTemplate(w, "halloffame", data)
}
