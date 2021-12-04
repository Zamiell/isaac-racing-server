package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func httpDebug(c *gin.Context) {
	if !isDev {
		c.String(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	logger.Debug("httpDebug function entered.")
	httpDebugFunc()
	logger.Debug("httpDebug function finished.")
	c.String(http.StatusOK, http.StatusText(http.StatusOK))
}

func httpDebugFunc() {
	logger.Debug("Doing seeded...")
	leaderboardRecalculateTrueSkill("seeded")
	logger.Debug("Doing unseeded...")
	leaderboardRecalculateTrueSkill("unseeded")
	logger.Debug("Doing diversity...")
	leaderboardRecalculateTrueSkill("diversity")
}
