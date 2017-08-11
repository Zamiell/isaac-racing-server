package main

/*
	Imports
*/

import (
	"strings"

	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

func websocketPrivateMessage(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	message := d.Message
	recipient := d.Name
	username := d.v.Username
	muted := d.v.Muted

	/*
		Validation
	*/

	// Validate that the requested person is sane
	if recipient == "" {
		log.Warning("User \"" + username + "\" tried to private message an empty string.")
		websocketError(s, d.Command, "That is not a valid person.")
		return
	}

	// Don't allow people to send PMs to themselves
	if recipient == username {
		websocketError(s, d.Command, "You cannot send a private message to yourself.")
		return
	}

	// Strip leading and trailing whitespace from the message
	message = strings.TrimSpace(message)

	// Don't allow empty messages
	if message == "" {
		log.Warning("User \"" + username + "\" tried to send an empty message.")
		websocketError(s, d.Command, "You cannot send an empty message.")
		return
	}

	// Validate that the user is not muted
	if muted {
		websocketError(s, d.Command, "You have been muted by an administrator, so you cannot chat with others.")
		return
	}

	// Validate that the person is online
	s2, ok := websocketSessions[recipient]
	if !ok {
		log.Info("User \"" + username + "\" tried to private message \"" + recipient + "\", but they are offline.")
		websocketError(s, d.Command, "That user is not online.")
		return
	}

	/*
		Private message
	*/

	// Add the new message to the database
	if err := db.ChatLogPM.Insert(recipient, username, message); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Send the message
	type PrivateMessageMessage struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}
	websocketEmit(s2, "privateMessage", &PrivateMessageMessage{username, message})

	// Log the message
	log.Info("PM <" + username + "> <" + recipient + "> " + message)
}
