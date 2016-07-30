package main

/*
 *  Imports
 */

import (
	"encoding/json"           // For Auth0 authentication (1/3)
	"github.com/trevex/golem" // The Golem WebSocket framework
	"golang.org/x/oauth2"     // For Auth0 authentication (2/3)
	"io/ioutil"               // For Auth0 authentication (3/3)
	"net"                     // For splitting the IP address from the port
	"net/http"                // For establishing an HTTP server
	"os"                      // For logging and reading environment variables
	"strconv"                 // For converting integers to strings
	"time"                    // For rate limiting
)

/*
 *  Constants
 */

const (
	rateLimitRate = 30 // In commands sent
	rateLimitPer  = 60 // In seconds
)

/*
 *  Custom Golem connection constructor
 */

func NewExtendedConnection(conn *golem.Connection) *ExtendedConnection {
	return &ExtendedConnection{
		Connection: conn,
		UserID:     0, // These values will be set (again) during the connOpen function
		Username:   "",
		Admin:      0,
	}
}

/*
 *  Login users using Auth0 access tokens
 */

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Local variables
	functionName := "loginHandler"
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Check to see if their IP is banned
	if userIsBanned, err := db.BannedIPs.Check(ip); err != nil {
		commandMutex.Unlock()
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if userIsBanned == true {
		commandMutex.Unlock()
		log.Info("IP \"" + ip + "\" tried to log in, but they are banned.")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Get the session (this may be an empty session)
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Unable to get the session during the", functionName, "function:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Check to see if they are already logged in (which should probably never happen since the cookie lasts 5 seconds)
	if _, ok := session.Values["userID"]; ok == true {
		commandMutex.Unlock()
		log.Warning("User from IP \"" + ip + "\" tried to get a session cookie, but they are already logged in.")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Instantiate the OAuth2 package
	conf := &oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		Scopes:       []string{"openid", "name", "email", "nickname"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + auth0Domain + "/authorize",
			TokenURL: "https://" + auth0Domain + "/oauth/token",
		},
	}

	// Get the POST JSON of the access token (that the client got from https://isaacserver.auth0.com/oauth/ro)
	decoder := json.NewDecoder(r.Body)
	var token oauth2.Token
	err = decoder.Decode(&token)
	if err != nil {
		commandMutex.Unlock()
		log.Warning("Failed to receive access token from user:", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Get information about the user
	client := conf.Client(oauth2.NoContext, &token)
	resp, err := client.Get("https://" + auth0Domain + "/userinfo")
	if err != nil {
		commandMutex.Unlock()
		log.Error("Failed to login with Auth0 token \""+token.AccessToken+"\":", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Reading the body
	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		commandMutex.Unlock()
		log.Error("Failed to read the body of the profile for Auth0 token \""+token.AccessToken+"\":", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Unmarshall the JSON of the profile
	var profile map[string]interface{}
	if err := json.Unmarshal(raw, &profile); err != nil {
		commandMutex.Unlock()
		log.Error("Failed to unmarshall the profile:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Get the Auth0 user ID and username from the profile
	auth0ID := profile["user_id"].(string)
	auth0Username := profile["username"].(string)

	// Check to see if the requested person exists in the database
	var squelched int
	userID, username, admin, err := db.Users.Login(auth0ID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if userID == 0 {
		// This is a new user, so add them to the database
		if userID, err = db.Users.Insert(auth0ID, auth0Username, ip); err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// By default, new users have the same stylizaition as their Auth0 username (all lowercase)
		username = auth0Username

		// By default, new users are not administrators
		admin = 0

		// By default, new users are not squelched
		squelched = 0

		// Log the user creation
		log.Info("Added \"" + username + "\" to the database (first login).")
	} else {
		// Check to see if this user is banned
		if userIsBanned, err := db.BannedUsers.Check(username); err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		} else if userIsBanned == true {
			commandMutex.Unlock()
			log.Info("User \"" + username + "\" tried to log in, but they are banned.")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Check to see if this user is squelched
		if userIsSquelched, err := db.SquelchedUsers.Check(username); err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		} else if userIsSquelched == true {
			squelched = 1
		} else {
			squelched = 0
		}

		// Update the database with last_login and last_ip
		if err := db.Users.SetLogin(username, ip); err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	// Save the information to the session
	session.Values["userID"] = userID
	session.Values["username"] = username
	session.Values["admin"] = admin
	session.Values["squelched"] = squelched
	if err := session.Save(r, w); err != nil {
		commandMutex.Unlock()
		log.Error("Failed to save the session cookie:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Log the login request
	log.Info("User \""+username+"\" logged in from:", ip)

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

/*
 *  Validate WebSocket connection
 */

func validateSession(w http.ResponseWriter, r *http.Request) bool {
	// Local variables
	functionName := "validateSession"
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Check to see if their IP is banned
	if userIsBanned, err := db.BannedIPs.Check(ip); err != nil {
		commandMutex.Unlock()
		log.Info("IP \"" + ip + "\" tried to establish a WebSocket connection, but they are banned.")
		return false
	} else if userIsBanned == true {
		commandMutex.Unlock()
		return false
	}

	// Get the session (this may be an empty session)
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Unable to get the session during the", functionName, "function:", err)
		return false
	}

	// If they have logged in, their cookie should have a "userID", "username", "admin", and "squelched" value
	if v, ok := session.Values["userID"]; ok == true && v.(int) > 0 {
		// Do nothing
	} else {
		commandMutex.Unlock()
		log.Debug("Unauthorized WebSocket handshake detected from:", ip, "(failed userID check)")
		return false
	}
	var username string
	if v, ok := session.Values["username"]; ok == true {
		username = v.(string)
	} else {
		commandMutex.Unlock()
		log.Debug("Unauthorized WebSocket handshake detected from:", ip, "(failed username check)")
		return false
	}
	if _, ok := session.Values["admin"]; ok == true {
		// Do nothing
	} else {
		commandMutex.Unlock()
		log.Debug("Unauthorized WebSocket handshake detected from:", ip, "(failed admin check)")
		return false
	}
	if _, ok := session.Values["squelched"]; ok == true {
		// Do nothing
	} else {
		commandMutex.Unlock()
		log.Debug("Unauthorized WebSocket handshake detected from:", ip, "(failed squelched check)")
		return false
	}

	// Check for sessions that belong to orphaned accounts
	if userExists, err := db.Users.Exists(username); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return false
	} else if userExists == false {
		commandMutex.Unlock()
		log.Error("User \"" + username + "\" does not exist in the database; they are trying to establish a WebSocket connection with an orphaned account.")
		return false
	}

	// Check to see if this user is banned
	if userIsBanned, err := db.BannedUsers.Check(username); err != nil {
		commandMutex.Unlock()
		return false
	} else if userIsBanned == true {
		commandMutex.Unlock()
		log.Info("User \"" + username + "\" tried to log in, but they are banned.")
		return false
	}

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()

	// If they got this far, they are a valid user
	return true
}

/*
 * Router connection functions
 */

func connOpen(conn *ExtendedConnection, r *http.Request) {
	// Local variables
	functionName := "connOpen"

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Get the session
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		// This should not fail, since we checked the session previously in the validateSession function
		commandMutex.Unlock()
		log.Error("Unable to get the session during the", functionName, "function:", err)
		return
	}

	// Get user information from the session
	var userID int
	if v, ok := session.Values["userID"]; ok == true && v.(int) > 0 {
		userID = v.(int)
	} else {
		commandMutex.Unlock()
		log.Error("Failed to retrieve \"userID\" from the session during the", functionName, "function.")
		return
	}
	var username string
	if v, ok := session.Values["username"]; ok == true {
		username = v.(string)
	} else {
		commandMutex.Unlock()
		log.Error("Failed to retrieve \"username\" from the session during the", functionName, "function.")
		return
	}
	var admin int
	if v, ok := session.Values["admin"]; ok == true {
		admin = v.(int)
	} else {
		commandMutex.Unlock()
		log.Error("Failed to retrieve \"admin\" from the session during the", functionName, "function.")
		return
	}
	var squelched int
	if v, ok := session.Values["squelched"]; ok == true {
		squelched = v.(int)
	} else {
		commandMutex.Unlock()
		log.Error("Failed to retrieve \"squelched\" from the cookie during the", functionName, "function.")
		return
	}

	// Store user information in the Golem connection so that we can use it in the Golem WebSocket functions later on
	conn.UserID = userID
	conn.Username = username
	conn.Admin = admin
	conn.Squelched = squelched
	conn.RateLimitAllowance = rateLimitRate
	conn.RateLimitLastCheck = time.Now()

	// Disconnect any existing connections with this username
	connectionMap.RLock()
	existingConnection, ok := connectionMap.m[username]
	connectionMap.RUnlock()
	if ok == true {
		log.Info("Closing existing connection for user \"" + username + "\".")
		connError(existingConnection, "logout", "You have logged on from somewhere else, so I'll disconnect you here.")
		existingConnection.Connection.Close()

		// Wait until the existing connection is terminated
		for {
			connectionMap.RLock()
			_, ok := connectionMap.m[username]
			connectionMap.RUnlock()
			if ok == false {
				break
			}
		}
	}

	// Add the connection to a connection map so that we can keep track of all of the connections
	connectionMap.Lock()
	connectionMap.m[username] = conn
	connectionMap.Unlock()

	// Log the connection
	log.Info("User \""+username+"\" connected;", len(connectionMap.m), "user(s) now connected.")

	// Join the user to the PMManager room corresponding to their username for private messages
	pmManager.Join(username, conn.Connection)

	// Get the current list of races
	raceList, err := db.Races.GetCurrentRaces()
	if err != nil {
		commandMutex.Unlock()
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Send it to the user
	conn.Connection.Emit("raceList", raceList)

	// Find out if the user is in any races that are currently going on
	for _, race := range raceList {
		for _, racer := range race.Racers {
			if racer == username {
				// Join the user to the chat room coresponding to this race
				roomJoinSub(conn, "_race_"+strconv.Itoa(race.ID))

				// Get all the information about the racers in this race
				racerList, err := db.RaceParticipants.GetRacerList(race.ID)
				if err != nil {
					commandMutex.Unlock()
					return
				}

				// Send it to the user
				conn.Connection.Emit("racerList", &RacerList{race.ID, racerList})

				// If the race is currently in the 10 second countdown
				if race.Status == "starting" {
					// Send them a message describing exactly when it will start
					conn.Connection.Emit("raceStart", &RaceStartMessage{
						race.ID,
						time.Now().Add(10 * time.Second).UnixNano(), // 10 seconds in the future
					})
				}

				break
			}
		}
	}

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func connClose(conn *ExtendedConnection) {
	// Local variables
	userID := conn.UserID
	username := conn.Username

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Delete the connection from the connection map
	connectionMap.Lock()
	delete(connectionMap.m, username) // This will do nothing if the entry doesn't exist
	connectionMap.Unlock()

	// Make a list of all the chat rooms this person is in
	var chatRoomList []string
	chatRoomMap.RLock()
	for room, users := range chatRoomMap.m {
		for _, user := range users {
			if user.Name == username {
				chatRoomList = append(chatRoomList, room)
				break
			}
		}
	}
	chatRoomMap.RUnlock()

	// Leave all the chat rooms
	for _, room := range chatRoomList {
		roomLeaveSub(conn, room)
	}

	// Leave the chat room dedicated for private messages
	pmManager.LeaveAll(conn.Connection)

	// Check to see if this user is in any races that are not already in progress
	raceIDs, err := db.RaceParticipants.GetNotStartedRaces(userID)
	if err != nil {
		commandMutex.Unlock()
		return
	}

	// Iterate over the races that they are currently in
	for _, raceID := range raceIDs {
		// Remove this user from the participants list for that race
		if err := db.RaceParticipants.Delete(userID, raceID); err != nil {
			commandMutex.Unlock()
			return
		}

		// Send everyone a notification that the user left the race
		connectionMap.RLock()
		for _, conn := range connectionMap.m {
			conn.Connection.Emit("raceLeft", RaceMessage{raceID, username})
		}
		connectionMap.RUnlock()

		// Check to see if the race should start or finish
		raceCheckStartFinish(raceID)
	}

	// Log the disconnection
	log.Info("User \""+username+"\" disconnected;", len(connectionMap.m), "user(s) now connected.")

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

/*
 *  WebSocket logout function
 */

func logout(conn *ExtendedConnection) {
	log.Debug("User \"" + conn.Username + "\" sent a logout command.")
	conn.Connection.Close()
}

/*
 *  WebSocket miscellaneous subroutines
 */

// Sent to the client after a successful command
func connSuccess(conn *ExtendedConnection, functionName string, msg interface{}) {
	conn.Connection.Emit("success", &SuccessMessage{functionName, msg})
}

// Sent to the client if either their command was unsuccessful or something else went wrong
func connError(conn *ExtendedConnection, functionName string, msg string) {
	conn.Connection.Emit("error", &ErrorMessage{functionName, msg})
}

// Called at the beginning of every command handler
func commandRateLimit(conn *ExtendedConnection) bool {
	// Local variables
	username := conn.Username

	// Rate limit commands; algorithm from: http://stackoverflow.com/questions/667508/whats-a-good-rate-limiting-algorithm
	now := time.Now()
	timePassed := now.Sub(conn.RateLimitLastCheck).Seconds()
	conn.RateLimitLastCheck = now
	conn.RateLimitAllowance += timePassed * (rateLimitRate / rateLimitPer)
	if conn.RateLimitAllowance > rateLimitRate {
		conn.RateLimitAllowance = rateLimitRate
	}
	if conn.RateLimitAllowance < 1 {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" triggered rate-limiting; disconnecting them.")
		connError(conn, "logout", "You have been disconnected due to flooding.")
		conn.Connection.Close()
		return true
	} else {
		conn.RateLimitAllowance--
		return false
	}
}
