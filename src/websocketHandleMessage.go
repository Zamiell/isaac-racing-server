package main

import (
	"encoding/json"
	"strings"

	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1" // A WebSocket framework
)

/*
	On top of the WebSocket protocol, the client and the server communicate
	using a specific format based on the Golem WebSocket framework protocol.
	First, the name of the command is sent, then a space, then the JSON of the
	data.

	Example:
		roomJoin {"room":"lobby"}
		roomMessage {"room":"lobby","message":"hey guys"}
*/

func websocketHandleMessage(s *melody.Session, msg []byte) {
	// Local variables
	functionName := "websocketHandleMessage"

	// Get the username from the session
	// (we will get all of the values from the session later on, but for now we
	// only need the username)
	var username string
	if v, exists := s.Get("username"); !exists {
		log.Error("Failed to get \"username\" from the session (in the \"" + functionName + "\" function).")
		websocketClose(s)
		return
	} else {
		username = v.(string)
	}

	// Unpack the message to see what kind of command it is
	// (this code is taken from Golem)
	result := strings.SplitN(string(msg), " ", 2)
	// We use SplitN() with a value of 2 instead of Split() so that if there
	// is a space in the JSON, the data part of the splice doesn't get
	// messed up
	if len(result) != 2 {
		log.Warning("User \"" + username + "\" sent an invalid WebSocket message.")
		return
	}
	command := result[0]
	jsonData := []byte(result[1])

	// Check to see if there is a command handler for this command
	if _, ok := commandHandlerMap[command]; !ok {
		log.Warning("User \"" + username + "\" sent an invalid command of \"" + command + "\".")
		return
	}

	// Unmarshal the JSON (this code is taken from Golem)
	var d *IncomingWebsocketData
	if err := json.Unmarshal(jsonData, &d); err != nil {
		log.Error("User \"" + username + "\" sent an command of \"" + command + "\" with invalid data: " + string(jsonData))
		return
	}

	// Attach the command name and session values to the data so that the
	// command handlers can conveniently use this information later on
	d.Command = command
	if !websocketGetSessionValues(s, d) {
		log.Error("Aborting before entering the command handler for \"" + d.Command + "\".")
	}

	// Call the command handler for this command
	log.Info("User \"" + username + "\" sent a command of \"" + command + "\".")
	commandMutex.Lock()
	commandHandlerMap[command](s, d)
	commandMutex.Unlock()
}
