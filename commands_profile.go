package main

/*
 *  Imports
 */

import (
	"strings"
)

/*
 *  WebSocket room/chat command functions
 */

func profileGet(conn *ExtendedConnection, data *ProfileMessage) {
	// Local variables
	functionName := "profileSetUsername"
	username := conn.Username
	profileUsername := data.Name

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the requested person exists in the database
	if userExists, err := db.Users.Exists(profileUsername); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userExists == false {
		connError(conn, functionName, "That user does not exist.")
		return
	}

	/*
	 *   Build the profile
	 */

	// Get the number of races
	// TODO
	var profile string

	// Send them the profile
	conn.Connection.Emit("profile", profile)

	// Send success confirmation
	connSuccess(conn, functionName, data)
}

func profileSetUsername(conn *ExtendedConnection, data *ProfileMessage) {
	// Local variables
	functionName := "profileSetUsername"
	userID := conn.UserID
	username := conn.Username
	newUsername := data.Name

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the submitted stylization
	if username == newUsername {
		connError(conn, functionName, "Your username is already set to that stylization.")
		return
	}

	// Validate the submitted stylization is not a different username
	if strings.ToLower(username) != strings.ToLower(newUsername) {
		connError(conn, functionName, "You can only change the capitalization of your username, not change it entirely.")
		return
	}

	// Set the new username in the database
	if err := db.Users.SetUsername(userID, newUsername); err != nil {
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
				chatRoomMap.Unlock()
				return
			}

			connectionMap.RLock()
			for _, user := range users {
				connectionMap.m[user.Name].Connection.Emit("roomList", &RoomList{
					room,
					users,
				})
			}
			connectionMap.RUnlock()
		}
	}
	chatRoomMap.Unlock()

	// Send success confirmation
	connSuccess(conn, functionName, data)
}
