package main

/*
	Imports
*/

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
)

/*
	The user will login with their Steam account according to this authentication flow:
	https://partner.steamgames.com/documentation/auth#client_to_backend_webapi
	(you have to be logged in for the link to work)
*/

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Local variables
	functionName := "loginHandler"
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	// Check to see if their IP is banned
	if userIsBanned, err := db.BannedIPs.Check(ip); err != nil {
		log.Error("Database error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if userIsBanned == true {
		log.Info("IP \"" + ip + "\" tried to log in, but they are banned.")
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
		log.Info("User from IP \"" + ip + "\" tried to get a session cookie, but they are already logged in.")
		http.Error(w, "You are already logged in. Please wait 5 seconds, then try again.", http.StatusUnauthorized)
		return
	}

	// Validate that the user sent the Steam ID and the ticket
	r.ParseForm()
	steamID := r.FormValue("steamID")
	if steamID == "" {
		log.Error("User from IP \"" + ip + "\" tried to log in, but they did not provide the \"steamID\" parameter.")
		http.Error(w, "You must provide the \"steamID\" parameter to log in.", http.StatusUnauthorized)
		return
	}
	ticket := r.FormValue("ticket")
	if ticket == "" {
		log.Error("User from IP \"" + ip + "\" tried to log in, but they did not provide the \"ticket\" parameter.")
		http.Error(w, "You must provide the \"ticket\" parameter to log in.", http.StatusUnauthorized)
		return
	}

	// Validate the ticket with the Steam API
	if validateSteamTicket(steamID, ticket, ip, w) == false {
		return
	}

	// Check to see if this Steam ID exists in the database
	var muted int
	userID, username, admin, err := db.Users.Login(steamID)
	if err != nil {
		log.Error("Database error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if userID == 0 {
		// This is a new user, so return a success, but don't give them a WebSocket cookie
		// (the client is expected to now make a POST request to "/register")
		http.Error(w, http.StatusText(http.StatusAccepted), http.StatusAccepted)
		return
	}

	// Check to see if this user is banned
	if userIsBanned, err := db.BannedUsers.Check(username); err != nil {
		log.Error("Database error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if userIsBanned == true {
		log.Info("User \"" + username + "\" tried to log in, but they are banned.")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Check to see if this user is muted
	if userIsMuted, err := db.MutedUsers.Check(username); err != nil {
		log.Error("Database error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if userIsMuted == true {
		muted = 1
	} else {
		muted = 0
	}

	// Update the database with datetime_last_login and last_ip
	if err := db.Users.SetLogin(username, ip); err != nil {
		log.Error("Database error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Save the information to the session
	session.Values["userID"] = userID
	session.Values["username"] = username
	session.Values["admin"] = admin
	session.Values["muted"] = muted
	if err := session.Save(r, w); err != nil {
		log.Error("Failed to save the session cookie:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Log the login request
	log.Info("User \""+username+"\" logged in from:", ip)
}

/*
	We need to create some structures that emulate what the JSON that the Steam API returns
*/

type SteamAPIReply struct {
	Response SteamAPIResponse `json:"response"`
}
type SteamAPIResponse struct {
	Params SteamAPIParams `json:"params"`
	Error  SteamAPIError  `json:"error"`
}
type SteamAPIParams struct {
	Result          string `json:"result"`
	SteamID         string `json:"steamid"`
	OwnerSteamId    string `json:"ownersteamid"`
	VACBanned       bool   `json:"vacbanned"`
	PublisherBanned bool   `json:"publisherbanned"`
}
type SteamAPIError struct {
	Code int    `json:"errorcode"`
	Desc string `json:"errordesc"`
}

/*
	Validate that the ticket is valid using the Steam web API
	E.g. https://api.steampowered.com/ISteamUserAuth/AuthenticateUserTicket/v1?key=secret&appid=113200&ticket=longhex
*/

func validateSteamTicket(steamID string, ticket string, ip string, w http.ResponseWriter) bool {
	// Make the request
	apiKey := os.Getenv("STEAM_WEB_API_KEY")
	appID := "113200" // This is the vanilla Binding of Isaac game ID on Steam
	resp, err := myHTTPClient.Get("https://api.steampowered.com/ISteamUserAuth/AuthenticateUserTicket/v1?key=" + apiKey + "&appid=" + appID + "&ticket=" + ticket)
	if err != nil {
		log.Error("Failed to query the Steam web API for IP \""+ip+"\": ", err)
		http.Error(w, "An error occured while verifying your Steam account. Please try again later.", http.StatusUnauthorized)
		return false
	}
	defer resp.Body.Close()

	// Read the body
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read the body of the response from the Steam web API for IP \""+ip+"\": ", err)
		http.Error(w, "An error occured while verifying your Steam account. Please try again later.", http.StatusUnauthorized)
		return false
	}

	// Unmarshall the JSON of the body from the response
	var steamAPIReply SteamAPIReply
	if err := json.Unmarshal(raw, &steamAPIReply); err != nil {
		log.Error("Failed to unmarshall the body of the response:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}

	// Check to see if we got an error
	steamError := steamAPIReply.Response.Error
	if steamError.Code != 0 {
		log.Error("The Steam web API returned error code " + strconv.Itoa(steamError.Code) + ": " + steamError.Desc)
		http.Error(w, "Your Steam account appears to be invalid. Please make sure you have the latest version of Steam installed and are correctly logged in.", http.StatusUnauthorized)
		return false
	}

	// Check to see if the ticket is valid
	result := steamAPIReply.Response.Params.Result
	if result == "" {
		log.Error("The Steam web API response does not have a \"result\" property.")
		http.Error(w, "An error occured while verifying your Steam account. Please try again later.", http.StatusUnauthorized)
		return false
	} else if result != "OK" {
		log.Warning("A user from IP \"" + ip + "\" tried to log in, but their Steam ticket was invalid.")
		http.Error(w, "Your Steam account appears to be invalid. Please make sure you have the latest version of Steam installed and are correctly logged in.", http.StatusUnauthorized)
		return false
	}

	// Check to see if the steam ID matches who they claim to be
	ticketSteamID := steamAPIReply.Response.Params.SteamID
	if ticketSteamID == "" {
		log.Error("The Steam web API response does not have a \"steamID\" property.")
		http.Error(w, "An error occured while verifying your Steam account. Please try again later.", http.StatusUnauthorized)
		return false
	} else if ticketSteamID != steamID {
		log.Warning("A user from IP \"" + ip + "\" submitted a Steam ticket that does not match their submitted Steam ID.")
		http.Error(w, "Your Steam account appears to be invalid. Please make sure you have the latest version of Steam installed and are correctly logged in.", http.StatusUnauthorized)
		return false
	}

	return true
}
