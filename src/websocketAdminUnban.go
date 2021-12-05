package server

import (
	"strconv"

	melody "gopkg.in/olahol/melody.v1"
)

func websocketAdminUnban(s *melody.Session, d *IncomingWebsocketData) {
	username := d.v.Username
	admin := d.v.Admin
	recipient := d.Name

	// Validate that the user is an admin
	if admin == 0 {
		logger.Warning("User \"" + username + "\" tried to unban \"" + recipient + "\", but they are not an administrator.")
		websocketError(s, d.Command, "Only administrators can do that.")
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		logger.Warning("User \"" + username + "\" tried to unban a blank person.")
		websocketWarning(s, d.Command, "That person is not valid.")
		return
	}

	// Validate that the requested person exists in the database
	var recipientID int
	if userExists, v, err := db.Users.Exists(recipient); err != nil {
		logger.Error("Database error when checking to see if user \""+recipient+"\" exists:", err)
		websocketError(s, d.Command, "")
		return
	} else if !userExists {
		websocketWarning(s, d.Command, "That user does not exist.")
		return
	} else {
		recipientID = v
	}

	// Validate that the requested person is banned
	if userIsBanned, err := db.BannedUsers.Check(recipientID); err != nil {
		logger.Error("Database error when checking to see if user "+strconv.Itoa(recipientID)+" is banned:", err)
		websocketError(s, d.Command, "")
		return
	} else if !userIsBanned {
		logger.Warning("User \"" + username + "\" tried to unban \"" + recipient + "\", but they are not banned.")
		websocketError(s, d.Command, "That user is not banned.")
		return
	}

	// Remove this username from the ban list in the database
	if err := db.BannedUsers.Delete(recipientID); err != nil {
		logger.Error("Database error when deleting user "+strconv.Itoa(recipientID)+" from the banned list:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Remove the user's last IP from the banned IP list, if present
	if err := db.BannedIPs.DeleteUserIP(recipientID); err != nil {
		logger.Error("Database error when deleting the IP for user "+strconv.Itoa(recipientID)+" from the banned IPs list:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Send the admin a message to let them know that the unban was successful
	websocketEmit(s, "roomMessage", &RoomMessageMessage{
		"lobby",
		"!server",
		"User \"" + recipient + "\" successfully unbanned.",
	})

	// Log the unban
	logger.Info("User \"" + username + "\" successfully unbanned user \"" + recipient + "\".")
}
