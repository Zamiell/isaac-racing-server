package main

import (
	"github.com/Zamiell/isaac-racing-server/src/log"
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
		log.Warning("User \"" + username + "\" tried to ban \"" + recipient + "\", but they are not an administrator.")
		websocketError(s, d.Command, "Only administrators can do that.")
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		log.Warning("User \"" + username + "\" tried to ban a blank person.")
		websocketWarning(s, d.Command, "That person is not valid.")
		return
	}

	// Validate that the requested person exists in the database
	var recipientID int
	if userExists, v, err := db.Users.Exists(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if !userExists {
		websocketWarning(s, d.Command, "That user does not exist.")
		return
	} else {
		recipientID = v
	}

	// Validate that the requested person is not a staff member or an administrator
	if recipientAdmin, err := db.Users.GetAdmin(recipientID); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if recipientAdmin > 0 {
		log.Warning("User \"" + username + "\" tried to ban \"" + recipient + "\", but staff/admins cannot be banned.")
		websocketWarning(s, d.Command, "You cannot ban a staff member or an administrator.")
		return
	}

	// Validate that the requested person is not already banned
	if userIsBanned, err := db.BannedUsers.Check(recipientID); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if userIsBanned {
		log.Warning("User \"" + username + "\" tried to ban \"" + recipient + "\", but they are already banned.")
		websocketError(s, d.Command, "That user is already banned.")
		return
	}

	// Add this username to the ban list in the database
	if err := db.BannedUsers.Insert(recipientID, userID, reason); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Add their IP to the banned IP list
	if err := db.BannedIPs.InsertUserIP(recipientID, userID, reason); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	/*
		// Find out if the banned user is in any races that are currently going on
		// TODO

		// Iterate over the races that they are currently in
		for _, race := range raceList {
			raceID := race.ID

			// Find out if the race is started
			if race.Status == "open" {
				// Remove this user from the participants list for that race
				// TODO

				// Send everyone a notification that the user left the race
				for _, conn := range connectionMap.m {
					conn.Connection.Emit("raceLeft", RaceMessage{raceID, recipient})
				}
			} else {
				// Set this racer's status to disqualified
				if err := db.RaceParticipants.SetStatus(recipient, raceID, "disqualified"); err != nil {
					log.Error("Database error:", err)
					websocketError(s, d.Command, "")
					return
				}

				// Set their finish time
				if err := db.RaceParticipants.SetDatetimeFinished(recipient, raceID, int(makeTimestamp())); err != nil {
					log.Error("Database error:", err)
					websocketError(s, d.Command, "")
					return
				}

				// Set their (final) place to -2 (which indicates a disqualified status)
				if err := db.RaceParticipants.SetPlace(username, raceID, -2); err != nil {
					log.Error("Database error:", err)
					websocketError(s, d.Command, "")
					return
				}

				// Send a notification to all the people in this particular race that the user got disqualified
				for _, racer := range racerNames {
					conn, ok := connectionMap.m[racer]
					if ok { // Not all racers may be online during a race
						conn.Connection.Emit("racerSetStatus", &RacerSetStatusMessage{raceID, username, "disqualified", -2})
					}
				}
			}

			// Check to see if the race should start or finish
			raceCheckStartFinish(raceID)
		}
	*/

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
	log.Info("User \"" + username + "\" successfully banned user \"" + recipient + "\".")
}
