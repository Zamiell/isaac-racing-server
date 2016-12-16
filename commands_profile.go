package main

/*
	Imports
*/

import (
	"regexp"
	"strings"
)

/*
	WebSocket profile command functions
*/

func profileGet(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "profileSetUsername"
	username := conn.Username
	profileUsername := data.Name

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the requested profile is sane
	if profileUsername == "" {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to request an empty profile.")
		connError(conn, functionName, "That is not a valid profile name.")
		return
	}

	// Validate that the requested person exists in the database
	if userExists, err := db.Users.Exists(profileUsername); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userExists == false {
		commandMutex.Unlock()
		connError(conn, functionName, "That user does not exist.")
		return
	}

	/*
		Build the profile
	*/

	// Get the number of races
	// TODO
	var profile Profile

	// Send them the profile
	conn.Connection.Emit("profile", profile)

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func profileSetUsername(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "profileSetUsername"
	userID := conn.UserID
	username := conn.Username
	newUsername := data.Name

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the submitted stylization is different than before
	if newUsername == username {
		commandMutex.Unlock()
		connError(conn, functionName, "Your username is already set to that stylization.")
		return
	}

	// Validate that the submitted stylization is not a different username (or an empty username)
	if strings.ToLower(newUsername) != strings.ToLower(username) {
		commandMutex.Unlock()
		connError(conn, functionName, "You can only change the capitalization of your username, not change it entirely.")
		return
	}

	// Validate that the user is not in any races that are currently going on
	raceList, err := db.RaceParticipants.GetCurrentRaces(username)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}
	if len(raceList) > 1 {
		commandMutex.Unlock()
		connError(conn, functionName, "You cannot change your name if you are currently in a race.")
		return
	}

	// Set the new username in the database
	if err := db.Users.SetUsername(userID, newUsername); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Look for this user in all chat rooms
	chatRoomMap.Lock()
	for room, users := range chatRoomMap.m {
		// See if the user is in this chat room
		index := -1
		for i, user := range users {
			if user.Name == username {
				index = i
				break
			}
		}
		if index != -1 {
			// Update their username
			chatRoomMap.m[room][index].Name = newUsername
		}
	}
	chatRoomMap.Unlock()

	// Send everyone a notification that the user changed their name
	connectionMap.RLock()
	for _, conn := range connectionMap.m {
		conn.Connection.Emit("profileSetName", &ProfileSetNameMessage{
			Name:    username,
			NewName: newUsername,
		})
	}
	connectionMap.RUnlock()

	// Change their username in the connection map
	// (Connections are indexed by username, so we have to delete the old entry and add a new one)
	connectionMap.Lock()
	tempConn := connectionMap.m[username]
	tempConn.Username = newUsername
	delete(connectionMap.m, username) // This will do nothing if the entry doesn't exist
	connectionMap.m[newUsername] = tempConn
	connectionMap.Unlock()

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func profileSetStream(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "profileSetStream"
	userID := conn.UserID
	username := conn.Username
	newStream := data.Name

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Get the user's current stream
	oldStream, err := db.Users.GetStream(userID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Validate that the submitted stream is different than the current one
	if newStream == oldStream {
		commandMutex.Unlock()
		connError(conn, functionName, "Your stream is already set to that.")
		return
	}

	// Validate that the submitted stream is not a nasty URL
	if newStream == "-" {
		// Do nothing
	} else if strings.HasPrefix(newStream, "https://www.twitch.tv/") {
		// Do nothing
	} else {
		commandMutex.Unlock()
		connError(conn, functionName, "Stream URLs must either be \"-\" or begin with \"https://www.twitch.tv/\".")
		return
	}
	// TODO Add Hitbox/Youtube?

	// If this is a Twitch stream, validate that the Twitch username is valid
	if strings.HasPrefix(newStream, "https://www.twitch.tv/") {
		// Parse the for username
		re, err := regexp.Compile(`https://www.twitch.tv/(.+)`)
		if err != nil {
			commandMutex.Unlock()
			log.Error("Failed to compile the Twitch stream regular expression.")
			connError(conn, functionName, "Something went wrong. Please contact an administrator.")
			return
		}
		user := re.FindStringSubmatch(newStream)[1]

		// Validate the username (from https://www.reddit.com/r/Twitch/comments/32w5b2/username_requirements/)
		re, err = regexp.Compile(`^[a-zA-Z0-9_]{4,25}$`)
		if err != nil {
			commandMutex.Unlock()
			log.Error("Failed to compile the Twitch username regular expression.")
			connError(conn, functionName, "Something went wrong. Please contact an administrator.")
			return
		}
		if re.FindString(user) == "" {
			commandMutex.Unlock()
			connError(conn, functionName, "The stream URL submitted does not have a valid Twitch username.")
			return
		}
	}

	// Validate that the user is not in any races that are currently going on
	raceList, err := db.RaceParticipants.GetCurrentRaces(username)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}
	if len(raceList) > 1 {
		commandMutex.Unlock()
		connError(conn, functionName, "You cannot change your stream URL if you are currently in a race.")
		return
	}

	// Set the new stream URL in the database
	if err := db.Users.SetStream(userID, newStream); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func profileSetTwitchBotEnabled(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "profileSetTwitchBotEnabled"
	userID := conn.UserID
	username := conn.Username
	newValue := data.Enabled

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Get the user's current Twitch bot setting
	twitchBotEnabled, err := db.Users.GetTwitchBotEnabled(username)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Validate that the setting is different than the current one
	if newValue == twitchBotEnabled {
		commandMutex.Unlock()
		connError(conn, functionName, "Your Twitch bot setting is already set to that.")
		return
	}

	// If they are turning it off, then we can do that and our work is finished
	if newValue == false {
		// Set the new Twitch bot setting in the database
		if err := db.Users.SetTwitchBotEnabled(userID, newValue); err != nil {
			commandMutex.Unlock()
			log.Error("Database error:", err)
			connError(conn, functionName, "Something went wrong. Please contact an administrator.")
			return
		} else {
			commandMutex.Unlock()
			return
		}
	}

	// Get the user's current stream
	stream, err := db.Users.GetStream(userID)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Parse their Twitch username from their stream URL
	var user string
	if strings.HasPrefix(stream, "https://www.twitch.tv/") {
		// Parse for the username
		re, err := regexp.Compile(`https://www.twitch.tv/(.+)`)
		if err != nil {
			commandMutex.Unlock()
			log.Error("Failed to compile the Twitch username regular expression:", err)
			connError(conn, functionName, "Something went wrong. Please contact an administrator.")
			return
		}
		user = re.FindStringSubmatch(stream)[1]
		user = strings.ToLower(user)
	} else {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to enable the Twitch bot without having a Twitch stream set.")
		connError(conn, functionName, "You must have a Twitch stream set in order to use the Twitch chat bot.")
		return
	}

	// Set the new Twitch bot setting in the database
	if err := db.Users.SetTwitchBotEnabled(userID, newValue); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// If it is a Twitch stream, make the Twitch IRC bot join their channel
	ircSend("JOIN #" + user)

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func profileSetTwitchBotDelay(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "profileSetTwitchBotDelay"
	userID := conn.UserID
	username := conn.Username
	newValue := data.Value

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Get the user's current Twitch bot delay
	twitchBotDelay, err := db.Users.GetTwitchBotDelay(username)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Validate that the delay is different than the current one
	if newValue == twitchBotDelay {
		commandMutex.Unlock()
		connError(conn, functionName, "Your Twitch bot delay is already set to that.")
		return
	}

	// Validate that it is a sane delay
	if newValue < 0 || newValue > 60 {
		commandMutex.Unlock()
		connError(conn, functionName, "Your Twitch bot delay must be between 0 and 60.")
		return
	}

	// Set the new Twitch bot delay in the database
	if err := db.Users.SetTwitchBotDelay(userID, newValue); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}
