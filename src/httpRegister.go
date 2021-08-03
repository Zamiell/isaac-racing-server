package server

import (
	"net"
	"net/http"
	"unicode/utf8"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func httpRegister(c *gin.Context) {
	// Local variables
	r := c.Request
	w := c.Writer
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	/*
		Validation
	*/

	// Check to see if their IP is banned
	if userIsBanned, err := db.BannedIPs.Check(ip); err != nil {
		logger.Error("Database error when checking to see if the IP \""+ip+"\" was banned:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if userIsBanned {
		logger.Info("IP \"" + ip + "\" tried to register, but they are banned.")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Check to see if they are already logged in
	// (which should probably never happen since the cookie lasts 5 seconds)
	session := sessions.Default(c)
	if v := session.Get("userID"); v != nil {
		logger.Info("User from IP \"" + ip + "\" tried to register, but they are already logged in.")
		http.Error(w, "You cannot register because you are already logged in.", http.StatusUnauthorized)
		return
	}

	// Validate that the user sent the Steam ID, the ticket, and the username
	steamID := c.PostForm("steamID")
	if steamID == "" {
		logger.Error("User from IP \"" + ip + "\" tried to register, but they did not provide the \"steamID\" parameter.")
		http.Error(w, "You must provide the \"steamID\" parameter to register.", http.StatusUnauthorized)
		return
	}
	ticket := c.PostForm("ticket")
	if ticket == "" {
		logger.Error("User from IP \"" + ip + "\" tried to register, but they did not provide the \"ticket\" parameter.")
		http.Error(w, "You must provide the \"ticket\" parameter to register.", http.StatusUnauthorized)
		return
	}
	username := c.PostForm("username")
	if username == "" {
		logger.Error("User from IP \"" + ip + "\" tried to register, but they did not provide the \"username\" parameter.")
		http.Error(w, "You must provide the \"username\" parameter to register.", http.StatusUnauthorized)
		return
	}

	// Validate the ticket with the Steam API
	// (this is in the "httpLogin.go" file)
	if !validateSteamTicket(steamID, ticket, ip, w) {
		return
	}

	// Check to see if this Steam ID exists in the database
	if sessionValues, err := db.Users.Login(steamID); err != nil {
		logger.Error("Database error when checking to see if the steam ID of \""+steamID+"\" exists:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if sessionValues != nil {
		// They are trying to register a new account, but this Steam ID already exists in the database
		logger.Error("User from IP \"" + ip + "\" tried to register, but they already exist in the database.")
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
	if userExists, _, err := db.Users.Exists(username); err != nil {
		logger.Error("Database error when checking to see if the username of \""+username+"\" exists:", err)
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
		logger.Error("Database error when inserting the username of \""+username+"\":", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else {
		userID = id
	}

	// Save the information to the session
	session.Set("userID", userID)
	session.Set("username", username)
	session.Set("admin", 0)                // By default, new users are not administrators
	session.Set("muted", false)            // By default, new users are not muted
	session.Set("streamURL", "-")          // By default, new users do not have a stream URL set
	session.Set("twitchBotEnabled", false) // By default, new users do not have the Twitch bot enabled
	session.Set("twitchBotDelay", 15)      // By default, new users have a Twitch bot delay of 15
	if err := session.Save(); err != nil {
		logger.Error("Failed to save the session:", err)
	}

	// Log the user creation
	logger.Info("Added \"" + username + "\" to the database (first login).")
}
