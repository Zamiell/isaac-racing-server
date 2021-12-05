package server

import (
	"net/http"
	"time"

	"github.com/Zamiell/isaac-racing-server/models"
	"github.com/gin-gonic/gin"
)

func httpWS(c *gin.Context) {
	w := c.Writer
	r := c.Request

	// The below function will return nil if there is an error or if the user is
	// not authorized
	var sessionValues *models.SessionValues
	if v, err := httpValidateSession(c); err != nil {
		if err.Error() == "" {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		return
	} else {
		sessionValues = v
	}

	// Transfer the values from the login cookie into WebSocket session variables
	keys := make(map[string]interface{})
	keys["userID"] = sessionValues.UserID
	keys["username"] = sessionValues.Username
	keys["admin"] = sessionValues.Admin
	keys["muted"] = sessionValues.Muted
	keys["streamURL"] = sessionValues.StreamURL
	keys["twitchBotEnabled"] = sessionValues.TwitchBotEnabled
	keys["twitchBotDelay"] = sessionValues.TwitchBotDelay
	keys["rateLimitAllowance"] = RateLimitRate
	keys["rateLimitLastCheck"] = time.Now()

	// Validation succeeded, so establish the WebSocket connection
	if err := m.HandleRequestWithKeys(w, r, keys); err != nil {
		logger.Error("Failed to add the keys to the websocket session:", err)
	}
}
