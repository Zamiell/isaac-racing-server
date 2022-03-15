package server

import (
	"github.com/gin-gonic/gin"
)

func httpHallOfFame(c *gin.Context) {
	w := c.Writer

	data := TemplateData{
		Title:             "Hall of Fame",
		Season1R9AB:       season1R9AB,
		Season1R14AB:      season1R14AB,
		Season2R7AB:       season2R7AB,
		Season3R7AB:       season3R7AB,
		Season4R7AB:       season4R7AB,
		Season5R7AB:       season5R7AB,
		Season6R7AB:       season6R7AB,
		Season7R7AB:       season7R7AB,
		Season8R7AB:       season8R7AB,
		Season1R7Rep:      season1R7Rep,
		Season1RankedSolo: season1RankedSolo,
		Season2RankedSolo: season2RankedSolo,
	}
	httpServeTemplate(w, "halloffame", data)
}
