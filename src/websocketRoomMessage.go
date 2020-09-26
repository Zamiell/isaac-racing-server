package main

import (
	"strings"
	"unicode/utf8"

	melody "gopkg.in/olahol/melody.v1"
)

/*
	Command example:
	roomMessage {
		room: "lobby",
		message: "hey guys",
	}
*/

func websocketRoomMessage(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	userID := d.v.UserID
	username := d.v.Username
	muted := d.v.Muted
	message := d.Message

	/*
		Perform validation
	*/

	// Validate that the requested room is sane
	if d.Room == "" {
		logger.Warning("User \"" + username + "\" tried to send a message, but did not provide a room.")
		websocketError(s, d.Command, "That is not a valid room name.")
		return
	}

	// Strip leading and trailing whitespace from the message
	message = strings.TrimSpace(message)

	// Don't allow empty messages
	if message == "" {
		logger.Warning("User \"" + username + "\" tried to send an empty message.")
		websocketWarning(s, d.Command, "You cannot send an empty message.")
		return
	}

	// Validate that the user is not muted
	if muted {
		websocketWarning(s, d.Command, "You have been muted by an administrator, so you cannot chat with others.")
		return
	}

	// Validate that the room exists
	users, ok := chatRooms[d.Room]
	if !ok {
		websocketError(s, d.Command, "That is not a valid room name.")
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
	if !userInRoom {
		logger.Warning("User \"" + username + "\" tried to message a room they were not in.")
		websocketError(s, d.Command, "You are not in that room.")
		return
	}

	// Validate that the message is not excessively long
	if utf8.RuneCountInString(message) > 150 {
		websocketWarning(s, d.Command, "Messages must not be longer than 150 characters.")
		return
	}

	/*
		Send the message
	*/

	// Add the new message to the database
	if err := db.ChatLog.Insert(d.Room, userID, message); err != nil {
		logger.Error("Database error when inserting a message:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Send the message to everyone in the room
	for _, user := range users {
		// All users in the chat room should be online, but check just in case
		if s2, ok := websocketSessions[user.Name]; ok {
			websocketEmit(s2, "roomMessage", &RoomMessageMessage{
				d.Room,
				username,
				message,
			})
		}
	}

	// Also send lobby messages to Discord
	if d.Room == "lobby" {
		discordSend(discordLobbyChannelID, "<"+username+"> "+message)
	}

	// Log the message
	logger.Info("#" + d.Room + " <" + username + "> " + message)
}
