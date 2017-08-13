package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func httpWS(c *gin.Context) {
	// Local variables
	w := c.Writer
	r := c.Request

	// The below function will return nil if there is an error or if the user is
	// not authorized
	sessionValues := httpValidateSession(c)
	if sessionValues == nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Transfer the values from the login cookie into WebSocket session
	// variables
	keys := make(map[string]interface{})
	keys["userID"] = sessionValues.UserID
	keys["username"] = sessionValues.Username
	keys["admin"] = sessionValues.Admin
	keys["muted"] = sessionValues.Muted
	keys["streamURL"] = sessionValues.StreamURL
	keys["twitchBotEnabled"] = sessionValues.TwitchBotEnabled
	keys["twitchBotDelay"] = sessionValues.TwitchBotDelay

	// Validation succeeded, so establish the WebSocket connection
	m.HandleRequestWithKeys(w, r, keys)
}
