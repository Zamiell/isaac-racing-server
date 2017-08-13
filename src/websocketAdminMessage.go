package main

import (
	"strings"
	"unicode/utf8"

	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

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
	d.Message = strings.TrimSpace(d.Message)

	// Don't allow empty messages
	if d.Message == "" {
		log.Warning("User \"" + username + "\" tried to send an empty message.")
		websocketWarning(s, d.Command, "You cannot send an empty message.")
		return
	}

	// Validate that the message is not excessively long
	if utf8.RuneCountInString(d.Message) > 150 {
		websocketWarning(s, d.Command, "Messages must not be longer than 150 characters.")
		return
	}

	/*
		Send the message
	*/

	// Add the new message to the database
	if err := db.ChatLog.Insert("server", userID, d.Message); err != nil {
		log.Error("Database error:", err)
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
	if d.Room == "lobby" {
		discordSend(discordLobbyChannelID, "[Server Notice] "+d.Message)
	}

	// Log the message
	log.Info("#" + d.Room + " <" + username + "> [Server Notice] " + d.Message)
}
