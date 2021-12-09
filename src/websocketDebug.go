package server

import (
	melody "gopkg.in/olahol/melody.v1"
)

func websocketDebug(s *melody.Session, d *IncomingWebsocketData) {
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

	debugRecalculateRankedSoloForSpecificUser(s, d)
}

func debugRecalculateRankedSoloForSpecificUser(s *melody.Session, d *IncomingWebsocketData) {
	username := d.Name
	if username == "" {
		return
	}

	// Get the user ID
	var userID int
	if exists, v, err := db.Users.Exists(username); err != nil {
		logger.Error("Failed to check to see if \""+username+"\" exists:", err)
		websocketError(s, d.Command, "")
		return
	} else if !exists {
		websocketError(s, d.Command, "That user does not exist.")
		return
	} else {
		userID = v
	}

	leaderboardRecalculateRankedSoloSpecificUser(userID)

	type PrivateMessageMessage struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}
	websocketEmit(s, "privateMessage", &PrivateMessageMessage{
		"SERVER",
		"Successfully reset ranked solo data.",
	})
}
