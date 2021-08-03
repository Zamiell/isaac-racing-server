package server

import (
	melody "gopkg.in/olahol/melody.v1"
)

func websocketAdminUnshutdown(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username
	admin := d.v.Admin

	// Validate that the user is an admin
	if admin != 2 {
		logger.Warning("User \"" + username + "\" tried to send turn off the shutdown mode, but they are not an administrator.")
		websocketError(s, d.Command, "Only administrators can do that.")
		return
	}

	shutdownMode = 0
}
