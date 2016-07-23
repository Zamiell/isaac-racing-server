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

 func roomJoin(conn *ExtendedConnection, data *RoomMessage) {
	// Local variables
	functionName := "roomJoin"
	username := conn.Username
	room := data.Name

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command for room \""+room+"\".")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the requested room is sane
	if room == "" {
		log.Warning("User \"" + username + "\" tried to join an empty room.")
		connError(conn, functionName, "That is not a valid room name.")
		return
	}

	// Validate that they are not trying to join a system room
	if strings.HasPrefix(room, "_") {
		log.Warning("Access denied to system room.")
		connError(conn, functionName, "You are not allowed to manually join system rooms.")
		return
	}

	// Validate that the room exists
	chatRoomMap.RLock()
	users, ok := chatRoomMap.m[room]
	chatRoomMap.RUnlock()
	if ok == true {
		// Validate that they are not already in the room
		userInRoom := false
		for _, user := range users {
			if user.Name == username {
				userInRoom = true
				break
			}
		}
		if userInRoom == true {
			log.Warning("User \"" + username + "\" tried to join a room they were already in.")
			connError(conn, functionName, "You are already in that room.")
			return
		}
	}

	// Let them join the room
	roomJoinSub(conn, room)

	// Send success confirmation
	connSuccess(conn, functionName, data)
}

func roomLeave(conn *ExtendedConnection, data *RoomMessage) {
	// Local variables
	functionName := "roomLeave"
	username := conn.Username
	room := data.Name

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command for room \""+room+"\".")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the requested room is sane
	if room == "" {
		log.Warning("User \"" + username + "\" tried to leave an empty room.")
		connError(conn, functionName, "That is not a valid room name.")
		return
	}

	// Validate that they are not trying to leave a system room
	if strings.HasPrefix(room, "_") {
		log.Warning("Access denied to system room.")
		connError(conn, functionName, "You are not allowed to manually leave system rooms.")
		return
	}

	// Validate that the room exists
	chatRoomMap.RLock()
	users, ok := chatRoomMap.m[room]
	chatRoomMap.RUnlock()
	if ok == false {
		log.Warning("User \"" + username + "\" tried to leave an invalid room.")
		connError(conn, functionName, "That is not a valid room name.")
		return
	}

	// Validate that they are actually in this room
	userInRoom := false
	for _, user := range users {
		if user.Name == username {
			userInRoom = true
			break
		}
	}
	if userInRoom == false {
		log.Warning("User \"" + username + "\" tried to leave a room they were not in.")
		connError(conn, functionName, "You are not in that room.")
		return
	}

	// Let them leave the room
	roomLeaveSub(conn, room)

	// Send success confirmation
	connSuccess(conn, functionName, data)
}

func roomMessage(conn *ExtendedConnection, data *ChatMessage) {
	// Local variables
	functionName := "roomMessage"
	username := conn.Username
	room := data.To

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the requested room is sane
	if room == "" {
		log.Warning("User \"" + username + "\" tried to message an empty room.")
		connError(conn, functionName, "That is not a valid room name.")
		return
	}

	// Validate that the room exists
	chatRoomMap.RLock()
	users, ok := chatRoomMap.m[room]
	chatRoomMap.RUnlock()
	if ok == false {
		connError(conn, functionName, "That is not a valid room name.")
		return
	}

	// Validate that they are actually in this room
	userInRoom := false
	for _, user := range users {
		if user.Name == username {
			userInRoom = true
			break
		}
	}
	if userInRoom == false {
		log.Warning("User \"" + username + "\" tried to message a room they were not in.")
		connError(conn, functionName, "You are not in that room.")
		return
	}

	// Don't allow empty messages
	if data.Msg == "" {
		return
	}

	// Validate that the user is not squelched
	if conn.Squelched == 1 {
		connError(conn, functionName, "You have been squelched by an administrator, so you cannot chat with others.")
		return
	}

	// Make sure that clients cannot masquerade as others
	data.From = conn.Username

	// Add the new message to the database
	if err := db.ChatLog.Insert(room, data.From, data.Msg); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Log the message
	log.Info("#" + room + " <" + data.From + "> " + data.Msg)

	// Send the message
	roomManager.Emit(room, functionName, &data)

	// Send success confirmation
	connSuccess(conn, functionName, data)
}

func privateMessage(conn *ExtendedConnection, data *ChatMessage) {
	// Local variables
	functionName := "privateMessage"
	username := conn.Username
	recipient := data.To

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		log.Warning("User \"" + username + "\" tried to private message an empty string.")
		connError(conn, functionName, "That is not a valid person.")
		return
	}

	// Validate that the person is online
	connectionMap.RLock()
	_, ok := connectionMap.m[data.To]
	connectionMap.RUnlock()
	if ok == false {
		log.Warning("User \"" + username + "\" tried to private message \"" + recipient + "\", who is offline.")
		connError(conn, functionName, "That user is not online.")
		return
	}

	// Don't allow empty messages
	if data.Msg == "" {
		return
	}

	// Validate that the user is not squelched
	if conn.Squelched == 1 {
		connError(conn, functionName, "You have been squelched by an administrator, so you cannot chat with others.")
		return
	}

	// Make sure that clients cannot masquerade as others
	data.From = conn.Username

	// Don't allow people to send PMs to themselves
	if data.From == data.To {
		connError(conn, functionName, "You cannot send a private message to yourself.")
		return
	}

	// Add the new message to the database
	if err := db.ChatLogPM.Insert(recipient, data.From, data.Msg); err != nil {
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Log the message
	log.Info("PM <" + data.From + "> <" + data.To + "> " + data.Msg)

	// Send the message
	pmManager.Emit(data.To, functionName, &data)

	// Send success confirmation
	connSuccess(conn, functionName, data)
}

func roomGetAll(conn *ExtendedConnection) {
	// Local variables
	functionName := "roomGetAll"
	username := conn.Username

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// We have to initialize this way to avoid sending a null on an empty array: https://danott.co/posts/json-marshalling-empty-slices-to-empty-arrays-in-go.html
	roomList := make([]Room, 0)
	chatRoomMap.RLock()
	for roomName, users := range chatRoomMap.m {
		// Add the room to the list
		roomList = append(roomList, Room{
			roomName,
			len(users),
		})
	}
	chatRoomMap.RUnlock()

	// Send it to the user
	conn.Connection.Emit("roomListAll", roomList)

	// Send success confirmation
	connSuccess(conn, functionName, "")
}


/*
 *  WebSocket room/chat subroutines
 */

func roomJoinSub(conn *ExtendedConnection, room string) {
	// Local variables
	username := conn.Username
	admin := conn.Admin
	squelched := conn.Squelched

	// Join the room
	roomManager.Join(room, conn.Connection)

	// Add the user to the chat room mapping
	chatRoomMap.Lock()
	chatRoomMap.m[room] = append(chatRoomMap.m[room], User{username, admin, squelched})
	users := chatRoomMap.m[room]
	chatRoomMap.Unlock()

	// Since the amount of people in the chat room changed, send everyone an update
	connectionMap.RLock()
	for _, user := range users {
		connectionMap.m[user.Name].Connection.Emit("roomList", &RoomList{
			room,
			users,
		})
	}
	connectionMap.RUnlock()

	// Log the join
	log.Debug("User \"" + conn.Username + "\" joined room: #" + room)
}

func roomLeaveSub(conn *ExtendedConnection, room string) {
	// Local variables
	username := conn.Username

	// Leave the room
	roomManager.Leave(room, conn.Connection)

	// Get the index of the user in the chat room mapping for this room
	chatRoomMap.RLock()
	users, ok := chatRoomMap.m[room]
	chatRoomMap.RUnlock()
	if ok == false {
		log.Error("Failed to get the chat room map for room \"" + room + "\".")
		return
	}
	index := -1
	for i, user := range users {
		if user.Name == username {
			index = i
			break
		}
	}
	if index == -1 {
		log.Error("Failed to get the index for the current user in the chat room map for room \"" + room + "\".")
		return
	}

	// Remove the user from the chat room mapping
	chatRoomMap.Lock()
	chatRoomMap.m[room] = append(users[:index], users[index+1:]...)
	users = chatRoomMap.m[room]
	chatRoomMap.Unlock()

	// Since the amount of people in the chat room changed, send everyone an update
	connectionMap.RLock()
	for _, user := range users {
		connectionMap.m[user.Name].Connection.Emit("roomList", &RoomList{
			room,
			users,
		})
	}
	connectionMap.RUnlock()

	// Log the leave
	log.Debug("User \"" + conn.Username + "\" left room: #" + room)
}
