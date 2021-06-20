package main

import (
	"strconv"

	melody "gopkg.in/olahol/melody.v1"
)

func websocketAdminBan(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	userID := d.v.UserID
	username := d.v.Username
	admin := d.v.Admin
	recipient := d.Name
	reason := d.Comment

	// Validate that the user is an admin
	if admin == 0 {
		logger.Warning("User \"" + username + "\" tried to ban \"" + recipient + "\", but they are not an administrator.")
		websocketError(s, d.Command, "Only administrators can do that.")
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		logger.Warning("User \"" + username + "\" tried to ban a blank person.")
		websocketWarning(s, d.Command, "That person is not valid.")
		return
	}

	// Validate that the requested person exists in the database
	var recipientID int
	if userExists, v, err := db.Users.Exists(recipient); err != nil {
		logger.Error("Database error while checking to see if user "+strconv.Itoa(userID)+" exists:", err)
		websocketError(s, d.Command, "")
		return
	} else if !userExists {
		websocketWarning(s, d.Command, "That user does not exist.")
		return
	} else {
		recipientID = v
	}

	// Validate that the requested person is not already banned
	if userIsBanned, err := db.BannedUsers.Check(recipientID); err != nil {
		logger.Error("Database error while checking to see if user "+strconv.Itoa(userID)+" is banned:", err)
		websocketError(s, d.Command, "")
		return
	} else if userIsBanned {
		logger.Warning("User \"" + username + "\" tried to ban \"" + recipient + "\", but they are already banned.")
		websocketError(s, d.Command, "That user is already banned.")
		return
	}

	// Validate that the requested person is not a staff member or an administrator
	if recipientAdmin, err := db.Users.GetAdmin(recipientID); err != nil {
		logger.Error("Database error while checking to see if user "+strconv.Itoa(userID)+" is an administrator:", err)
		websocketError(s, d.Command, "")
		return
	} else if recipientAdmin > 0 {
		logger.Warning("User \"" + username + "\" tried to ban \"" + recipient + "\", but staff/admins cannot be banned.")
		websocketWarning(s, d.Command, "You cannot ban a staff member or an administrator.")
		return
	}

	// Add the player to the ban list in the database
	if err := db.BannedUsers.Insert(recipientID, userID, reason); err != nil {
		logger.Error("Database error while adding user "+strconv.Itoa(recipientID)+" to the ban list:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Add their IP to the banned IP list
	if err := db.BannedIPs.InsertUserIP(recipientID, userID, reason); err != nil {
		logger.Error("Database error while adding the IP for user "+strconv.Itoa(recipientID)+" to the banned IPs list:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Find out if the banned user is in any races that are currently going on
	for _, race := range races {
		for _, racer := range race.Racers {
			if racer.ID == recipientID {
				/*
					Disqualify
				*/

				racer.Place = -2
				racer.PlaceMid = -1
				race.SetRacerStatus(username, "disqualified")
				racer.DatetimeFinished = getTimestamp()
				racer.RunTime = racer.DatetimeFinished - race.DatetimeStarted
				race.SetAllPlaceMid()
				twitchRacerQuit(race, racer)
				race.CheckFinish()
			}
		}
	}

	// Boot them offline if they are currently connected
	if s2, ok := websocketSessions[recipient]; ok {
		websocketError(
			s2,
			"Banned",
			"You have been banned. If you think this was a mistake, please contact the administration to appeal.",
		)
		websocketClose(s2)
	}

	// Send the admin a message to let them know that the ban was successful
	websocketEmit(s, "roomMessage", &RoomMessageMessage{
		"lobby",
		"!server",
		"User \"" + recipient + "\" successfully banned.",
	})

	// Log the ban
	logger.Info("User \"" + username + "\" successfully banned user \"" + recipient + "\".")
}
