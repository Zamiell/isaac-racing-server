package main

import (
	"net"
	"net/http"
	"unicode/utf8"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func httpRegister(c *gin.Context) {
	// Local variables
	r := c.Request
	w := c.Writer
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	// Lock the command mutex for the duration of the function to prevent database locks
	commandMutex.Lock()
	defer commandMutex.Unlock()

	/*
		Validation
	*/

	// Check to see if their IP is banned
	if userIsBanned, err := db.BannedIPs.Check(ip); err != nil {
		log.Error("Database error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if userIsBanned {
		log.Info("IP \"" + ip + "\" tried to register, but they are banned.")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Check to see if they are already logged in
	// (which should probably never happen since the cookie lasts 5 seconds)
	session := sessions.Default(c)
	if v := session.Get("userID"); v != nil {
		log.Info("User from IP \"" + ip + "\" tried to register, but they are already logged in.")
		http.Error(w, "You cannot register because you are already logged in.", http.StatusUnauthorized)
		return
	}

	// Validate that the user sent the Steam ID, the ticket, and the username
	steamID := c.PostForm("steamID")
	if steamID == "" {
		log.Error("User from IP \"" + ip + "\" tried to register, but they did not provide the \"steamID\" parameter.")
		http.Error(w, "You must provide the \"steamID\" parameter to register.", http.StatusUnauthorized)
		return
	}
	ticket := c.PostForm("ticket")
	if ticket == "" {
		log.Error("User from IP \"" + ip + "\" tried to register, but they did not provide the \"ticket\" parameter.")
		http.Error(w, "You must provide the \"ticket\" parameter to register.", http.StatusUnauthorized)
		return
	}
	username := c.PostForm("username")
	if username == "" {
		log.Error("User from IP \"" + ip + "\" tried to register, but they did not provide the \"username\" parameter.")
		http.Error(w, "You must provide the \"username\" parameter to register.", http.StatusUnauthorized)
		return
	}

	// Validate the ticket with the Steam API
	// (this is in the "httpLogin.go" file)
	if !validateSteamTicket(steamID, ticket, ip, w) {
		return
	}

	// Check to see if this Steam ID exists in the database
	if userID, _, _, _, err := db.Users.Login(steamID); err != nil {
		log.Error("Database error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if userID != 0 {
		// They are trying to register a new account, but this Steam ID already exists in the database
		log.Error("User from IP \"" + ip + "\" tried to register, but they already exist in the database.")
		http.Error(w, "There is already a Racing+ account tied to your Steam ID, so you cannot register a new one.", http.StatusUnauthorized)
		return
	}

	// Validate that the username has between 2 and 16 characters
	if utf8.RuneCountInString(username) < 2 || utf8.RuneCountInString(username) > 16 {
		http.Error(w, "The username must be between 2 and 16 characters.", http.StatusUnauthorized)
		return
	}

	// Validate that the username only contain alphanumeric characters and underscores
	if !isAlphaNumericUnderscore(username) {
		http.Error(w, "The username must only contain alphanumeric characters and underscores.", http.StatusUnauthorized)
		return
	}

	// Validate that the username is not already taken
	if userExists, err := db.Users.Exists(username); err != nil {
		log.Error("Database error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if userExists {
		http.Error(w, "Someone has already claimed that username.", http.StatusUnauthorized)
		return
	}

	/*
		Register
	*/

	// Add them to the database
	var userID int
	if id, err := db.Users.Insert(steamID, username, ip); err != nil {
		log.Error("Database error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else {
		userID = id
	}

	// Save the information to the session
	session.Set("userID", userID)
	session.Set("username", username)
	session.Set("admin", 0)       // By default, new users are not administrators
	session.Set("muted", 0)       // By default, new users are not muted
	session.Set("streamURL", "-") // By default, new users do not have a stream URL set
	session.Save()

	// Log the user creation
	log.Info("Added \"" + username + "\" to the database (first login).")
}
