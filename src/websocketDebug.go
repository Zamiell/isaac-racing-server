package main

/*
	Imports
*/

import (
	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

/*
	WebSocket debug command functions
*/

func websocketDebug(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username
	admin := d.v.Admin

	// Validate that the user is an admin
	if admin == 0 {
		log.Info("User \"" + username + "\" tried to do a debug command, but they are not staff/admin.")
		websocketError(s, d.Command, "Only staff members or administrators can do that.")
		return
	}

	// Print out the connection map
	/*
		fmt.Println(connectionMap.m)
		for _, conn := range connectionMap.m {
			fmt.Println("on connection:", conn.Username)
		}
	*/

	// Test IRC stuff
	//ircSend("JOIN #zamiell")
}
