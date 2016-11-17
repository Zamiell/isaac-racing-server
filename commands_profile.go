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

	// Validate that the submitted stylization
	if newUsername == username {
		commandMutex.Unlock()
		connError(conn, functionName, "Your username is already set to that stylization.")
		return
	}

	// Validate the submitted stylization is not a different username (or an empty username)
	if strings.ToLower(newUsername) != strings.ToLower(username) {
		commandMutex.Unlock()
		connError(conn, functionName, "You can only change the capitalization of your username, not change it entirely.")
		return
	}

	// Set the new username in the database
	if err := db.Users.SetUsername(userID, newUsername); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Set the new username in the connection
	conn.Username = newUsername

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

			// Send everyone an room update
			users, ok := chatRoomMap.m[room]
			if ok == false {
				log.Error("Failed to retrieve the user list from the chat room map for room \"" + room + "\".")
				continue
			}

			connectionMap.RLock()
			for _, user := range users {
				userConnection, ok := connectionMap.m[user.Name] // This should always succeed, but there might be a race condition
				if ok == true {
					userConnection.Connection.Emit("roomSetName", &RoomSetNameMessage{room, username, newUsername})
				} else {
					log.Error("Failed to get the connection for user \"" + user.Name + "\" while setting a new username for user \"" + username + "\".")
					continue
				}
			}
			connectionMap.RUnlock()
		}
	}
	chatRoomMap.Unlock()

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}
