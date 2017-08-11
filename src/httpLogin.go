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
	"strings"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

/*
	The user will login with their Steam account according to this authentication flow:
	https://partner.steamgames.com/documentation/auth#client_to_backend_webapi
	(you have to be logged in for the link to work)
*/

func httpLogin(c *gin.Context) {
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
		log.Info("IP \"" + ip + "\" tried to log in, but they are banned.")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Check to see if they are already logged in
	// (which should probably never happen since the cookie lasts 5 seconds)
	session := sessions.Default(c)
	if v := session.Get("userID"); v != nil {
		log.Info("User from IP \"" + ip + "\" tried to get a session cookie, but they are already logged in.")
		http.Error(w, "You are already logged in. Please wait 5 seconds, then try again.", http.StatusUnauthorized)
		return
	}

	// Validate that the user sent the Steam ID and the ticket
	steamID := c.PostForm("steamID")
	if steamID == "" {
		log.Error("User from IP \"" + ip + "\" tried to log in, but they did not provide the \"steamID\" parameter.")
		http.Error(w, "You must provide the \"steamID\" parameter to log in.", http.StatusUnauthorized)
		return
	}
	ticket := c.PostForm("ticket")
	if ticket == "" {
		log.Error("User from IP \"" + ip + "\" tried to log in, but they did not provide the \"ticket\" parameter.")
		http.Error(w, "You must provide the \"ticket\" parameter to log in.", http.StatusUnauthorized)
		return
	}

	// Validate the ticket with the Steam API
	if !validateSteamTicket(steamID, ticket, ip, w) {
		return
	}

	// Check to see if this Steam ID exists in the database
	// (being banned and muted are in separate tables for reason and timestamp
	// logging purposes, so we must get them at a later step)
	userID, username, admin, streamURL, err := db.Users.Login(steamID)
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
	} else if userIsBanned {
		log.Info("User \"" + username + "\" tried to log in, but they are banned.")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Check to see if this user is muted
	var muted bool
	if userIsMuted, err := db.MutedUsers.Check(username); err != nil {
		log.Error("Database error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if userIsMuted {
		muted = true
	} else {
		muted = false
	}

	/*
		Login
	*/

	// Update the database with datetime_last_login and last_ip
	if err := db.Users.SetLogin(username, ip); err != nil {
		log.Error("Database error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Save the information to the session
	session.Set("userID", userID)
	session.Set("username", username)
	session.Set("admin", admin)
	session.Set("muted", muted)
	session.Set("streamURL", streamURL)
	session.Save()

	// Log the login request
	log.Info("User \""+username+"\" logged in from:", ip)
}

/*
	We need to create some structures that emulate the JSON that the Steam API returns
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
	E.g. https://api.steampowered.com/ISteamUserAuth/AuthenticateUserTicket/v1?key=secret&appid=250900&ticket=longhex
*/

func validateSteamTicket(steamID string, ticket string, ip string, w http.ResponseWriter) bool {
	// Automatically validate test accounts
	if ticket == "debug" &&
		steamID == "-1" || // These 10 fake steam IDs allow for 10 test accounts
		steamID == "-2" ||
		steamID == "-3" ||
		steamID == "-4" ||
		steamID == "-5" ||
		steamID == "-6" ||
		steamID == "-7" ||
		steamID == "-8" ||
		steamID == "-9" ||
		steamID == "-10" {

		IPWhitelist := os.Getenv("DEV_IP_WHITELIST")
		IPs := strings.Split(IPWhitelist, ",")
		for _, validIP := range IPs {
			if ip == validIP {
				return true
			}
		}

		log.Warning("IP \"" + ip + "\" tried to use a debug ticket, but they are not on the whitelist.")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return false
	}

	// Make the request
	apiKey := os.Getenv("STEAM_WEB_API_KEY")
	appID := "250900" // This is the app ID on Steam for The Binding of Isaac: Rebirth
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

	invalidMessage := "Your Steam account appears to be invalid. Please make sure you have the latest version of Steam installed and are correctly logged in."

	// Check to see if we got an error
	steamError := steamAPIReply.Response.Error
	if steamError.Code != 0 {
		log.Error("The Steam web API returned error code " + strconv.Itoa(steamError.Code) + " for IP " + ip + " and Steam ID \"" + steamID + "\" and ticket \"" + ticket + "\": " + steamError.Desc)
		http.Error(w, invalidMessage, http.StatusUnauthorized)
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
		http.Error(w, invalidMessage, http.StatusUnauthorized)
		return false
	}

	// Check to see if the Steam ID matches who they claim to be
	ticketSteamID := steamAPIReply.Response.Params.SteamID
	if ticketSteamID == "" {
		log.Error("The Steam web API response does not have a \"steamID\" property.")
		http.Error(w, "An error occured while verifying your Steam account. Please try again later.", http.StatusUnauthorized)
		return false
	} else if ticketSteamID != steamID {
		log.Warning("A user from IP \"" + ip + "\" submitted a Steam ticket that does not match their submitted Steam ID.")
		http.Error(w, invalidMessage, http.StatusUnauthorized)
		return false
	}

	return true
}
