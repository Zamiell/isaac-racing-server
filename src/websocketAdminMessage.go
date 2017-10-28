package main

import (
	"strings"
	"unicode/utf8"

	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

/*
	This is sent on the client with the "/notice" command.

	Command example:
	adminMessage {
		message: "the tournament is starting soon",
	}
*/

// Also called from the "websocketAdminShutdown" function
func websocketAdminMessage(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	userID := d.v.UserID
	username := d.v.Username
	admin := d.v.Admin
	message := d.Message

	/*
		Validation
	*/

	// Validate that the user is an admin
	if admin != 2 {
		log.Warning("User \"" + username + "\" tried to send a server broadcast, but they are not an administrator.")
		websocketError(s, d.Command, "Only administrators can do that.")
		return
	}

	// Strip leading and trailing whitespace from the message
	message = strings.TrimSpace(message)

	// Don't allow empty messages
	if message == "" {
		log.Warning("User \"" + username + "\" tried to send an empty message.")
		websocketWarning(s, d.Command, "You cannot send an empty message.")
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

	// Prefix the message to designate that it is a special message
	message = "[Server Notice] " + message

	// Add the new message to the database
	if err := db.ChatLog.Insert("server", userID, message); err != nil {
		log.Error("Database error while inserting the chat message:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Send everyone the server broadcast notification
	for _, s := range websocketSessions {
		websocketEmit(s, "adminMessage", &AdminMessageMessage{
			Message: message,
		})
	}

	// Also send lobby messages to Discord
	discordSend(discordLobbyChannelID, message)

	// Log the message
	log.Info("#SERVER <" + username + "> " + message)
}
