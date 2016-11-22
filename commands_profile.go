package main

/*
	Imports
*/

import (
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
