package main

import (
	"net"
	"net/http"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Local variables
	functionName := "registerHandler"
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	// Check to see if their IP is banned
	if userIsBanned, err := db.BannedIPs.Check(ip); err != nil {
		log.Error("Database error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if userIsBanned == true {
		log.Info("IP \"" + ip + "\" tried to register, but they are banned.")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Get the session (this may be an empty session)
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		log.Error("Unable to get the session during the", functionName, "function:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Check to see if they are already logged in (which should probably never happen since the cookie lasts 5 seconds)
	if _, ok := session.Values["userID"]; ok == true {
		log.Info("User from IP \"" + ip + "\" tried to register, but they are already logged in.")
		http.Error(w, "You cannot register because you are already logged in.", http.StatusUnauthorized)
		return
	}

	// Validate that the user sent the Steam ID and the ticket
	r.ParseForm()
	steamID := r.FormValue("steamID")
	if steamID == "" {
		log.Error("User from IP \"" + ip + "\" tried to register, but they did not provide the \"steamID\" parameter.")
		http.Error(w, "You must provide the \"steamID\" parameter to register.", http.StatusUnauthorized)
		return
	}
	ticket := r.FormValue("ticket")
	if ticket == "" {
		log.Error("User from IP \"" + ip + "\" tried to register, but they did not provide the \"ticket\" parameter.")
		http.Error(w, "You must provide the \"ticket\" parameter to register.", http.StatusUnauthorized)
		return
	}
	username := r.FormValue("username")
	if username == "" {
		log.Error("User from IP \"" + ip + "\" tried to register, but they did not provide the \"username\" parameter.")
		http.Error(w, "You must provide the \"username\" parameter to register.", http.StatusUnauthorized)
		return
	}

	// Validate the ticket with the Steam API
	// (this is in the "login.go" file)
	if validateSteamTicket(steamID, ticket, ip, w) == false {
		return
	}

	// Check to see if this Steam ID exists in the database
	userID, _, _, err := db.Users.Login(steamID)
	if err != nil {
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
	if len(username) < 2 || len(username) > 16 {
		http.Error(w, "The username must be between 2 and 16 characters.", http.StatusUnauthorized)
		return
	}

	// Validate that the username only contain alphanumeric characters and underscores
	if isAlphaNumericUnderscore(username) == false {
		http.Error(w, "The username must only contain alphanumeric characters and underscores.", http.StatusUnauthorized)
		return
	}

	// Validate that the username is not already taken
	if userExists, err := db.Users.Exists(username); err != nil {
		log.Error("Database error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if userExists == true {
		http.Error(w, "Someone has already claimed that username.", http.StatusUnauthorized)
		return
	}

	// Add them to the database
	if userID, err = db.Users.Insert(steamID, username, ip); err != nil {
		log.Error("Database error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Save the information to the session
	session.Values["userID"] = userID
	session.Values["username"] = username
	session.Values["admin"] = 0 // By default, new users are not administrators
	session.Values["muted"] = 0 // By default, new users are not muted
	if err := session.Save(r, w); err != nil {
		log.Error("Failed to save the session cookie:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Log the user creation
	log.Info("Added \"" + username + "\" to the database (first login).")
}
