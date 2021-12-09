package server

import (
	melody "gopkg.in/olahol/melody.v1"
)

func websocketRankedSoloReset(s *melody.Session, d *IncomingWebsocketData) {
	userID := d.v.UserID

	if err := db.Users.ResetRankedSolo(userID); err != nil {
		logger.Error("Failed to reset the ranked solo fields:", err)
		websocketError(s, d.Command, "")
		return
	}

	if err := db.Races.DeleteOldRankedSoloRaces(userID); err != nil {
		logger.Error("Failed to delete the old ranked solo races:", err)
		websocketError(s, d.Command, "")
		return
	}

	type PrivateMessageMessage struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}
	websocketEmit(s, "privateMessage", &PrivateMessageMessage{
		"SERVER",
		"Successfully reset ranked solo data.",
	})
}
