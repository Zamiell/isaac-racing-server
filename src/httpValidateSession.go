package main

import (
	"errors"
	"net"
	"strconv"

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
func httpValidateSession(c *gin.Context) (*models.SessionValues, error) {
	// Local variables
	r := c.Request
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	// Check to see if their IP is banned
	if userIsBanned, err := db.BannedIPs.Check(ip); err != nil {
		logger.Error("Database error when checking to see if IP \""+ip+"\" is banned:", err)
		return nil, errors.New("")
	} else if userIsBanned {
		logger.Info("IP \"" + ip + "\" tried to establish a WebSocket connection, but they are banned.")
		msg := "Your IP address has been banned. Please contact an administrator if you think this is a mistake."
		return nil, errors.New(msg) // nolint: stylecheck
	}

	// If they have logged in, their cookie should have values matching the
	// SessionValues struct
	session := sessions.Default(c)
	var userID int
	if v := session.Get("userID"); v == nil {
		logger.Info("Unauthorized WebSocket handshake detected from \"" + ip + "\" (failed userID check).")
		return nil, errors.New("")
	} else {
		userID = v.(int)
	}
	var username string
	if v := session.Get("username"); v == nil {
		logger.Info("Unauthorized WebSocket handshake detected from \"" + ip + "\" (failed username check).")
		return nil, errors.New("")
	} else {
		username = v.(string)
	}
	var admin int
	if v := session.Get("admin"); v == nil {
		logger.Info("Unauthorized WebSocket handshake detected from \"" + ip + "\" (failed admin check).")
		return nil, errors.New("")
	} else {
		admin = v.(int)
	}
	var muted bool
	if v := session.Get("muted"); v == nil {
		logger.Info("Unauthorized WebSocket handshake detected from \"" + ip + "\" (failed muted check).")
		return nil, errors.New("")
	} else {
		muted = v.(bool)
	}
	var streamURL string
	if v := session.Get("streamURL"); v == nil {
		logger.Info("Unauthorized WebSocket handshake detected from \"" + ip + "\" (failed streamURL check).")
		return nil, errors.New("")
	} else {
		streamURL = v.(string)
	}
	var twitchBotEnabled bool
	if v := session.Get("twitchBotEnabled"); v == nil {
		logger.Info("Unauthorized WebSocket handshake detected from \"" + ip + "\" (failed twitchBotEnabled check).")
		return nil, errors.New("")
	} else {
		twitchBotEnabled = v.(bool)
	}
	var twitchBotDelay int
	if v := session.Get("twitchBotDelay"); v == nil {
		logger.Info("Unauthorized WebSocket handshake detected from \"" + ip + "\" (failed twitchBotDelay check).")
		return nil, errors.New("")
	} else {
		twitchBotDelay = v.(int)
	}

	// Check for sessions that belong to orphaned accounts
	if userExists, databaseID, err := db.Users.Exists(username); err != nil {
		logger.Error("Database error when checking to see if user \""+username+"\" exists:", err)
		return nil, errors.New("")
	} else if !userExists {
		logger.Error("User \"" + username + "\" does not exist in the database; they are trying to establish a WebSocket connection with an orphaned account.")
		return nil, errors.New("")
	} else if userID != databaseID {
		logger.Error("User \"" + username + "\" exists in the database, but they are trying to establish a WebSocket connection with an account ID that does not match the ID in the database.")
		return nil, errors.New("")
	}

	// Check to see if this user is banned
	if userIsBanned, err := db.BannedUsers.Check(userID); err != nil {
		logger.Error("Database error when checking to see if user "+strconv.Itoa(userID)+" is banned:", err)
		return nil, errors.New("")
	} else if userIsBanned {
		logger.Info("User \"" + username + "\" tried to establish a WebSocket connection, but they are banned.")
		msg := "Your user account has been banned. Please contact an administrator if you think this is a mistake."
		return nil, errors.New(msg) // nolint: stylecheck
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
	}, nil
}
