package server

import (
	"github.com/gin-gonic/gin"
)

func httpHallOfFame(c *gin.Context) {
	w := c.Writer

	data := TemplateData{
		Title:             "Hall of Fame",
		Season1r9:         season1r9,
		Season1r14:        season1r14,
		Season2r7:         season2r7,
		Season3r7:         season3r7,
		Season4r7:         season4r7,
		Season5r7:         season5r7,
		Season6r7:         season6r7,
		Season7r7:         season7r7,
		Season8r7:         season8r7,
		Season1RankedSolo: season1RankedSolo,
		Season2RankedSolo: season2RankedSolo,
	}
	httpServeTemplate(w, "halloffame", data)
}
