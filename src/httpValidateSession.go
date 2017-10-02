package main

import (
	"net"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/Zamiell/isaac-racing-server/src/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

/*
	Validate that they have logged in before opening a WebSocket connection

	Essentially, all we need to do is check to see if they have any cookie
	values stored, because that implies that they got through the "httpLogin"
	less than 5 seconds ago. But we also do a few other checks to be thorough.
*/

// Called from the "httpWS()" function
func httpValidateSession(c *gin.Context) *models.SessionValues {
	// Local variables
	r := c.Request
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	// Check to see if their IP is banned
	if userIsBanned, err := db.BannedIPs.Check(ip); err != nil {
		log.Error("Database error:", err)
		return nil
	} else if userIsBanned {
		log.Info("IP \"" + ip + "\" tried to establish a WebSocket connection, but they are banned.")
		return nil
	}

	// If they have logged in, their cookie should have values matching the
	// SessionValues struct
	session := sessions.Default(c)
	var userID int
	if v := session.Get("userID"); v == nil {
		log.Info("Unauthorized WebSocket handshake detected from \"" + ip + "\" (failed userID check).")
		return nil
	} else {
		userID = v.(int)
	}
	var username string
	if v := session.Get("username"); v == nil {
		log.Info("Unauthorized WebSocket handshake detected from \"" + ip + "\" (failed username check).")
		return nil
	} else {
		username = v.(string)
	}
	var admin int
	if v := session.Get("admin"); v == nil {
		log.Info("Unauthorized WebSocket handshake detected from \"" + ip + "\" (failed admin check).")
		return nil
	} else {
		admin = v.(int)
	}
	var muted bool
	if v := session.Get("muted"); v == nil {
		log.Info("Unauthorized WebSocket handshake detected from \"" + ip + "\" (failed muted check).")
		return nil
	} else {
		muted = v.(bool)
	}
	var streamURL string
	if v := session.Get("streamURL"); v == nil {
		log.Info("Unauthorized WebSocket handshake detected from \"" + ip + "\" (failed streamURL check).")
		return nil
	} else {
		streamURL = v.(string)
	}
	var twitchBotEnabled bool
	if v := session.Get("twitchBotEnabled"); v == nil {
		log.Info("Unauthorized WebSocket handshake detected from \"" + ip + "\" (failed twitchBotEnabled check).")
		return nil
	} else {
		twitchBotEnabled = v.(bool)
	}
	var twitchBotDelay int
	if v := session.Get("twitchBotDelay"); v == nil {
		log.Info("Unauthorized WebSocket handshake detected from \"" + ip + "\" (failed twitchBotDelay check).")
		return nil
	} else {
		twitchBotDelay = v.(int)
	}

	// Check for sessions that belong to orphaned accounts
	if userExists, databaseID, err := db.Users.Exists(username); err != nil {
		log.Error("Database error:", err)
		return nil
	} else if !userExists {
		log.Error("User \"" + username + "\" does not exist in the database; they are trying to establish a WebSocket connection with an orphaned account.")
		return nil
	} else if userID != databaseID {
		log.Error("User \"" + username + "\" exists in the database, but they are trying to establish a WebSocket connection with an account ID that does not match the ID in the database.")
		return nil
	}

	// Check to see if this user is banned
	if userIsBanned, err := db.BannedUsers.Check(userID); err != nil {
		log.Error("Database error:", err)
		return nil
	} else if userIsBanned {
		log.Info("User \"" + username + "\" tried to establish a WebSocket connection, but they are banned.")
		return nil
	}

	// If they got this far, they are a valid user
	return &models.SessionValues{
		UserID:           userID,
		Username:         username,
		Admin:            admin,
		Muted:            muted,
		StreamURL:        streamURL,
		TwitchBotEnabled: twitchBotEnabled,
		TwitchBotDelay:   twitchBotDelay,
		Banned:           false,
	}
}
