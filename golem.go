package main

/*
	Imports
*/

import (
	"github.com/trevex/golem"
	"net"
	"net/http"
	"strconv"
	"time"
)

/*
	Custom Golem connection constructor
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
	Validate WebSocket connection
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
		log.Error("Database error:", err)
		return false
	} else if userIsBanned == true {
		commandMutex.Unlock()
		log.Info("IP \"" + ip + "\" tried to establish a WebSocket connection, but they are banned.")
		return false
	}

	// Get the session (this may be an empty session)
	session, err := sessionStore.Get(r, sessionName)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Unable to get the session during the", functionName, "function:", err)
		return false
	}

	// If they have logged in, their cookie should have a "userID", "username", "admin", and "muted" value
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
	if _, ok := session.Values["muted"]; ok == true {
		// Do nothing
	} else {
		commandMutex.Unlock()
		log.Debug("Unauthorized WebSocket handshake detected from:", ip, "(failed muted check)")
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
		log.Error("Database error:", err)
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
	Router connection functions
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
	var muted int
	if v, ok := session.Values["muted"]; ok == true {
		muted = v.(int)
	} else {
		commandMutex.Unlock()
		log.Error("Failed to retrieve \"muted\" from the cookie during the", functionName, "function.")
		return
	}

	// Store user information in the Golem connection so that we can use it in the Golem WebSocket functions later on
	conn.UserID = userID
	conn.Username = username
	conn.Admin = admin
	conn.Muted = muted
	conn.RateLimitAllowance = rateLimitRate
	conn.RateLimitLastCheck = time.Now()

	// Disconnect any existing connections with this username
	connectionMap.RLock()
	existingConnection, ok := connectionMap.m[username]
	connectionMap.RUnlock()
	if ok == true {
		log.Info("Closing existing connection for user \"" + username + "\".")
		connError(existingConnection, "logout", "You have logged on from somewhere else, so you have been disconnected here.")
		existingConnection.Connection.Close()

		// Wait until the existing connection is terminated
		commandMutex.Unlock()
		for {
			connectionMap.RLock()
			_, ok := connectionMap.m[username]
			connectionMap.RUnlock()
			if ok == false {
				break
			}
		}
		commandMutex.Lock()
	}

	// Add the connection to a connection map so that we can keep track of all of the connections
	connectionMap.Lock()
	connectionMap.m[username] = conn
	log.Info("User \""+username+"\" connected;", len(connectionMap.m), "user(s) now connected.") // Log the connection
	connectionMap.Unlock()

	// Get their stream URL
	streamURL, err := db.Users.GetStreamURL(userID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	}

	// Get their Twitch bot settings
	twitchBotEnabled, err := db.Users.GetTwitchBotEnabled(username)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	}
	twitchBotDelay, err := db.Users.GetTwitchBotDelay(username)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		return
	}

	// Send them various settings tied to their account
	conn.Connection.Emit("settings", &SettingsMessage{
		Username:         username, // The client already knows their username, but it may not be the same as the server-side stylization
		StreamURL:        streamURL,
		TwitchBotEnabled: twitchBotEnabled,
		TwitchBotDelay:   twitchBotDelay,
		Time:             makeTimestamp(), // Send them the current time so that they can calculate the local offset
	})

	// Join the user to the PMManager room corresponding to their username for private messages
	pmManager.Join(username, conn.Connection)

	// Get the current list of races
	raceList, err := db.Races.GetCurrentRaces()
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
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
					log.Error("Database error:", err)
					return
				}

				// Send it to the user
				conn.Connection.Emit("racerList", &RacerList{race.ID, racerList})

				// If the race is currently in the 10 second countdown
				if race.Status == "starting" {
					// Get the time 10 seconds in the future
					startTime := time.Now().Add(10*time.Second).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
					// This will technically put them behind the other racers by some amount of seconds, but it gives them 10 seconds to get ready after a disconnect

					// Send them a message describing exactly when it will start
					conn.Connection.Emit("raceStart", &RaceStartMessage{
						race.ID,
						startTime,
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
		log.Error("Database error:", err)
		return
	}

	// Iterate over the races that they are currently in
	for _, raceID := range raceIDs {
		// Remove this user from the participants list for that race
		if err := db.RaceParticipants.Delete(username, raceID); err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			return
		}
		// Send everyone a notification that the user left the race
		connectionMap.RLock()
		for _, conn := range connectionMap.m {
			conn.Connection.Emit("raceLeft", RaceMessage{raceID, username})
		}
		connectionMap.RUnlock()

		// Check to see if the race should start
		raceCheckStart(raceID)
	}

	// Log the disconnection
	connectionMap.RLock()
	log.Info("User \""+username+"\" disconnected;", len(connectionMap.m), "user(s) now connected.")
	connectionMap.RUnlock()

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

/*
	WebSocket logout function
*/

func logout(conn *ExtendedConnection) {
	log.Debug("User \"" + conn.Username + "\" sent a logout command.")
	conn.Connection.Close()
}

/*
	WebSocket miscellaneous subroutines
*/

// Sent to the client if either their command was unsuccessful or something else went wrong
// (this will cause a WebSocket disconnect and their client to completely restart)
func connError(conn *ExtendedConnection, functionName string, message string) {
	conn.Connection.Emit("error", &ErrorMessage{functionName, message})
}

// Sent to the client if something unexpected happened
// (this will cause a popup on the client but still allow them to continue what they were doing)
func connWarning(conn *ExtendedConnection, functionName string, message string) {
	conn.Connection.Emit("warning", &ErrorMessage{functionName, message})
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
