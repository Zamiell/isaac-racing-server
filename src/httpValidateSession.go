package main

/*
	Imports
*/

import (
	"net"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

/*
	Validate that they have logged in before opening a WebSocket connection
*/

func httpValidateSession(c *gin.Context) *SessionValues {
	// Local variables
	r := c.Request
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	// Lock the command mutex for the duration of the function to prevent database locks
	commandMutex.Lock()
	defer commandMutex.Unlock()

	// Check to see if their IP is banned
	if userIsBanned, err := db.BannedIPs.Check(ip); err != nil {
		log.Error("Database error:", err)
		return nil
	} else if userIsBanned {
		log.Info("IP \"" + ip + "\" tried to establish a WebSocket connection, but they are banned.")
		return nil
	}

	// If they have logged in, their cookie should have values of
	// "userID", "username", "admin", "muted", and "streamURL"
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

	// Check for sessions that belong to orphaned accounts
	if userExists, err := db.Users.Exists(username); err != nil {
		log.Error("Database error:", err)
		return nil
	} else if !userExists {
		log.Error("User \"" + username + "\" does not exist in the database; they are trying to establish a WebSocket connection with an orphaned account.")
		return nil
	}

	// Check to see if this user is banned
	if userIsBanned, err := db.BannedUsers.Check(username); err != nil {
		log.Error("Database error:", err)
		return nil
	} else if userIsBanned {
		log.Info("User \"" + username + "\" tried to log in, but they are banned.")
		return nil
	}

	// If they got this far, they are a valid user
	return &SessionValues{
		userID,
		username,
		admin,
		muted,
		streamURL,
	}
}
