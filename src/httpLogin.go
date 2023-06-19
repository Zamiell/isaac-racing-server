package server

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/Zamiell/isaac-racing-server/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

/*
	The user will login with their Steam account according to this authentication flow:
	https://partner.steamgames.com/documentation/auth#client_to_backend_webapi
	(You have to be logged in for the link to work.)
*/

func httpLogin(c *gin.Context) {
	r := c.Request
	w := c.Writer
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	/*
		Validation
	*/

	// Check to see if their IP is banned.
	if userIsBanned, err := db.BannedIPs.Check(ip); err != nil {
		logger.Error("Database error when checking to see if IP \""+ip+"\" was banned:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if userIsBanned {
		logger.Info("IP \"" + ip + "\" tried to log in, but they are banned.")
		http.Error(w, "Your IP address has been banned. Please contact an administrator if you think this is a mistake.", http.StatusUnauthorized)
		return
	}

	// Check to see if they are already logged in (which should probably never happen since the
	// cookie lasts 5 seconds).
	session := sessions.Default(c)
	if v := session.Get("userID"); v != nil {
		logger.Info("User from IP \"" + ip + "\" tried to get a session cookie, but they are already logged in.")
		http.Error(w, "You are already logged in. Please wait 5 seconds, then try again.", http.StatusUnauthorized)
		return
	}

	// Validate that the user sent the Steam ID, the ticket, and the version number of the client.
	steamID := c.PostForm("steamID")
	if steamID == "" {
		logger.Error("User from IP \"" + ip + "\" tried to log in, but they did not provide the \"steamID\" parameter.")
		http.Error(w, "You must provide the \"steamID\" parameter to log in.", http.StatusUnauthorized)
		return
	}
	ticket := c.PostForm("ticket")
	if ticket == "" {
		logger.Error("User from IP \"" + ip + "\" tried to log in, but they did not provide the \"ticket\" parameter.")
		http.Error(w, "You must provide the \"ticket\" parameter to log in.", http.StatusUnauthorized)
		return
	}
	version := c.PostForm("version")
	if version == "" {
		logger.Error("User from IP \"" + ip + "\" tried to log in, but they did not provide the \"version\" parameter.")
		http.Error(w, "You must provide the \"version\" parameter to log in.", http.StatusUnauthorized)
		return
	}

	// Validate that the provided Steam ID is sane.
	var steamIDint int
	if v, err := strconv.Atoi(steamID); err != nil {
		logger.Error("Failed to convert the steam ID to an integer.")
		http.Error(w, "You provided an invalid \"steamID\".", http.StatusUnauthorized)
		return
	} else {
		steamIDint = v
	}

	// Validate that the Racing+ client version is the latest version.
	if steamIDint > 0 {
		if !validateLatestVersion(version, w) {
			return
		}
	}

	// Validate the ticket with the Steam API.
	if !validateSteamTicket(steamID, ticket, ip, w) {
		return
	}

	// Check to see if this Steam ID exists in the database.
	var sessionValues *models.SessionValues
	if v, err := db.Users.Login(steamID); err != nil {
		logger.Error("Database error when checking to see if steam ID "+steamID+" exists:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if v == nil {
		// This is a new user, so return a success, but don't give them a WebSocket cookie. (The
		// client is expected to now make a POST request to "/register".)
		http.Error(w, http.StatusText(http.StatusAccepted), http.StatusAccepted)
		return
	} else {
		sessionValues = v
	}

	// Check to see if this user is banned.
	if sessionValues.Banned {
		logger.Info("User \"" + sessionValues.Username + "\" tried to log in, but they are banned.")
		http.Error(w, "Your user account has been banned. Please contact an administrator if you think this is a mistake.", http.StatusUnauthorized)
		return
	}

	/*
		Login
	*/

	// Update the database with `datetime_last_login` and `last_ip`.
	if err := db.Users.SetLogin(sessionValues.UserID, ip); err != nil {
		logger.Error("Database error when setting the login values:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Save the information to the session.
	session.Set("userID", sessionValues.UserID)
	session.Set("username", sessionValues.Username)
	session.Set("admin", sessionValues.Admin)
	session.Set("muted", sessionValues.Muted)
	session.Set("streamURL", sessionValues.StreamURL)
	session.Set("twitchBotEnabled", sessionValues.TwitchBotEnabled)
	session.Set("twitchBotDelay", sessionValues.TwitchBotDelay)
	if err := session.Save(); err != nil {
		logger.Error("Failed to save the session:", err)
	}

	// Log the login request.
	logger.Info("User \""+sessionValues.Username+"\" logged in from:", ip)

	// Now, the end user will attempt to establish a WebSocket connection.
}

/*
	We need to create some structures that emulate the JSON that the Steam API returns.
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
	OwnerSteamID    string `json:"ownersteamid"`
	VACBanned       bool   `json:"vacbanned"`
	PublisherBanned bool   `json:"publisherbanned"`
}
type SteamAPIError struct {
	Code int    `json:"errorcode"`
	Desc string `json:"errordesc"`
}

/*
	Validate that the ticket is valid using the Steam web API.
	e.g. https://api.steampowered.com/ISteamUserAuth/AuthenticateUserTicket/v1?key=secret&appid=250900&ticket=longhex
*/

func validateSteamTicket(steamID string, ticket string, ip string, w http.ResponseWriter) bool {
	// Automatically validate test accounts.
	if ticket == "debug" &&
		steamID == "-1" || // These 10 fake steam IDs allow for 10 test accounts.
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

		logger.Warning("IP \"" + ip + "\" tried to use a debug ticket, but they are not on the whitelist.")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return false
	}

	// Make the request.
	apiKey := os.Getenv("STEAM_WEB_API_KEY")
	if len(apiKey) == 0 {
		logger.Error("The \"STEAM_WEB_API_KEY\" environment variable is blank; aborting the login request.")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}
	appID := "250900" // This is the app ID on Steam for "The Binding of Isaac: Rebirth".
	url := "https://api.steampowered.com/ISteamUserAuth/AuthenticateUserTicket/v1"
	args := "?key=" + apiKey + "&appid=" + appID + "&ticket=" + ticket
	resp, err := HTTPClientWithTimeout.Get(url + args)
	if err != nil {
		logger.Error("Failed to query the Steam web API for IP \""+ip+"\": ", err)
		http.Error(w, "An error occurred while verifying your Steam account. Please try again later.", http.StatusUnauthorized)
		return false
	}
	defer resp.Body.Close()

	// Read the body.
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read the body of the response from the Steam web API for IP \""+ip+"\": ", err)
		http.Error(w, "An error occurred while verifying your Steam account. Please try again later.", http.StatusUnauthorized)
		return false
	}

	// Unmarshall the JSON of the body from the response.
	var steamAPIReply SteamAPIReply
	if err := json.Unmarshal(raw, &steamAPIReply); err != nil {
		logger.Error("Failed to unmarshall the body of the response from the Steam web API for IP \""+ip+":", err)
		logger.Error("The response was as follows:", string(raw))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}

	invalidMessage := "Your Steam account appears to be invalid. Please make sure you have the latest version of Steam installed and are correctly logged in."

	// Check to see if we got an error.
	steamError := steamAPIReply.Response.Error
	if steamError.Code != 0 {
		logger.Error("The Steam web API returned error code " + strconv.Itoa(steamError.Code) + " for IP " + ip + " and Steam ID \"" + steamID + "\" and ticket \"" + ticket + "\": " + steamError.Desc)
		http.Error(w, invalidMessage, http.StatusUnauthorized)
		return false
	}

	// Check to see if the ticket is valid.
	result := steamAPIReply.Response.Params.Result
	if result == "" {
		logger.Error("The Steam web API response does not have a \"result\" property.")
		http.Error(w, "An error occurred while verifying your Steam account. Please try again later.", http.StatusUnauthorized)
		return false
	} else if result != "OK" {
		logger.Warning("A user from IP \"" + ip + "\" tried to log in, but their Steam ticket was invalid.")
		http.Error(w, invalidMessage, http.StatusUnauthorized)
		return false
	}

	// Check to see if the Steam ID matches who they claim to be.
	ticketSteamID := steamAPIReply.Response.Params.SteamID
	if ticketSteamID == "" {
		logger.Error("The Steam web API response does not have a \"steamID\" property.")
		http.Error(w, "An error occurred while verifying your Steam account. Please try again later.", http.StatusUnauthorized)
		return false
	} else if ticketSteamID != steamID {
		logger.Warning("A user from IP \"" + ip + "\" submitted a Steam ticket that does not match their submitted Steam ID.")
		http.Error(w, invalidMessage, http.StatusUnauthorized)
		return false
	}

	return true
}

func validateLatestVersion(version string, w http.ResponseWriter) bool {
	// Make an exception for users on macOS.
	if version == "macOS" {
		return true
	}

	latestVersionRaw, err := ioutil.ReadFile(path.Join(projectPath, "latest_client_version.txt"))
	if err != nil {
		logger.Error("Failed to read the \"latest_client_version.txt\" file:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}
	latestVersion := string(latestVersionRaw)
	latestVersion = strings.TrimSpace(latestVersion)
	if len(latestVersion) == 0 {
		logger.Error("The \"latest_client_version.txt\" file is empty, so users will not be able to login to the WebSocket server.")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}
	if version != latestVersion {
		errorMsg := "Your client version is <strong>" + version + "</strong> and the latest version is <strong>" + latestVersion + "</strong>.<br /><br />Please restart the Racing+ program and it should automatically update to the latest version. If that does not work, you can try manually downloading the latest version from the Racing+ website."
		http.Error(w, errorMsg, http.StatusUnauthorized)
		return false
	}

	return true
}
