package server

import (
	melody "gopkg.in/olahol/melody.v1"
)

func websocketDebug(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username
	admin := d.v.Admin

	// Validate that the user is an admin
	if admin == 0 {
		logger.Info("User \"" + username + "\" tried to do a debug command, but they are not staff/admin.")
		websocketError(s, d.Command, "Only staff members or administrators can do that.")
		return
	}

	debugPrintGlobals()

	logger.Debug("debugFunc entered.")
	debugFunc()
	logger.Debug("debugFunc finished.")
}
