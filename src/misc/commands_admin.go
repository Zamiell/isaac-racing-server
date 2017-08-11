package main

/*
	WebSocket admin command functions
*/

/*
func websocketAdminBan(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	functionName := "adminBan"
	userID := conn.UserID
	username := conn.Username
	recipient := data.Name

	// Validate that the user is staff/admin
	if conn.Admin == 0 {
		log.Warning("User \"" + username + "\" tried to ban someone, but they are not staff/admin.")
		websocketError(s, d.Command, "Only staff members or administrators can do that.")
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		log.Warning("User \"" + username + "\" tried to ban a blank person.")
		websocketError(s, d.Command, "That person is not valid.")
		return
	}

	// Validate that the requested person exists in the database
	if userExists, err := db.Users.Exists(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if !userExists {
		websocketError(s, d.Command, "That user does not exist.")
		return
	}

	// Validate that the requested person is not a staff member or an administrator
	if userIsStaff, err := db.Users.CheckStaff(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if userIsStaff {
		log.Warning("User \"" + username + "\" tried to ban \"" + recipient + "\", but staff/admins cannot be banned.")
		websocketError(s, d.Command, "You cannot ban a staff member or an administrator.")
		return
	}

	// Validate that the requested person is not already banned
	if userIsBanned, err := db.BannedUsers.Check(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if userIsBanned {
		log.Warning("User \"" + username + "\" tried to ban \"" + recipient + "\", but they are already banned.")
		websocketError(s, d.Command, "That user is already banned.")
		return
	}

	// Add this username to the ban list in the database
	if err := db.BannedUsers.Insert(recipient, userID); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Add their IP to the banned IP list
	if err := db.BannedIPs.Insert(recipient, userID); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Find out if the banned user is in any races that are currently going on
	raceList, err := db.RaceParticipants.GetCurrentRaces(recipient)
	if err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Iterate over the races that they are currently in
	for _, race := range raceList {
		raceID := race.ID

		// Find out if the race is started
		if race.Status == "open" {
			// Remove this user from the participants list for that race
			if err := db.RaceParticipants.Delete(recipient, raceID); err != nil {
				log.Error("Database error:", err)
				websocketError(s, d.Command, "")
				return
			}

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

			// Get the list of racers for this race
			racerNames, err := db.RaceParticipants.GetRacerNames(raceID)
			if err != nil {
				log.Error("Database error:", err)
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

	// Check to see if the user is online
	bannedConnection, ok := connectionMap.m[recipient]
	if ok {
		// Disconnect the user
		connError(bannedConnection, functionName, "You have been banned. If you think this was a mistake, please contact the administration to appeal.")
		bannedConnection.Connection.Close()
	}

	// Log the ban
	log.Info("User \"" + username + "\" banned user \"" + recipient + "\".")
}

func websocketAdminUnban(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	functionName := "adminUnban"
	username := conn.Username
	recipient := data.Name

	// Validate that the user is staff/admin
	if conn.Admin == 0 {
		log.Warning("User \"" + username + "\" tried to unban someone, but they are not staff/admin.")
		websocketError(s, d.Command, "Only staff members or administrators can do that.")
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		log.Warning("User \"" + username + "\" tried to unban a blank person.")
		websocketError(s, d.Command, "That person is not valid.")
		return
	}

	// Validate that the requested person exists in the database
	if userExists, err := db.Users.Exists(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if !userExists {
		websocketError(s, d.Command, "That user does not exist.")
		return
	}

	// Validate that the requested person is not a staff member or an administrator
	if userIsStaff, err := db.Users.CheckStaff(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if userIsStaff {
		log.Warning("User \"" + username + "\" tried to unban \"" + recipient + "\", but staff/admins cannot be unbanned.")
		websocketError(s, d.Command, "You cannot unban a staff member or an administrator.")
		return
	}

	// Validate that the requested person is banned
	if userIsBanned, err := db.BannedUsers.Check(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if !userIsBanned {
		log.Warning("User \"" + username + "\" tried to unban \"" + recipient + "\", but they are not banned.")
		websocketError(s, d.Command, "That user is not banned.")
		return
	}

	// Remove this username from the ban list in the database
	if err := db.BannedUsers.Delete(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Remove the user's last IP from the banned IP list, if present
	if err := db.BannedIPs.Delete(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Log the unban
	log.Info("User \"" + username + "\" unbanned user \"" + recipient + "\".")
}

func websocketAdminBanIP(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	functionName := "adminBanIP"
	userID := conn.UserID
	username := conn.Username
	ip := data.IP

	// Validate that the user is staff/admin
	if conn.Admin == 0 {
		log.Warning("User \"" + username + "\" tried to ban an IP, but they are not staff/admin.")
		websocketError(s, d.Command, "Only staff members or administrators can do that.")
		return
	}

	// Validate that the requested IP is sane
	if ip == "" {
		log.Warning("User \"" + username + "\" tried to ban a blank IP.")
		websocketError(s, d.Command, "That IP is not valid.")
		return
	}

	// Validate that the requested IP is not already banned
	if IPBanned, err := db.BannedIPs.Check(ip); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if IPBanned {
		websocketError(s, d.Command, "That IP is already banned.")
		return
	}

	// Add the IP to the list in the database
	if err := db.BannedIPs.InsertIP(ip, userID); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Log the ban
	log.Info("User \"" + username + "\" banned IP \"" + ip + "\".")
}

func websocketAdminUnbanIP(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	functionName := "adminUnbanIP"
	username := conn.Username
	ip := data.IP

	// Validate that the user is staff/admin
	if conn.Admin == 0 {
		log.Warning("User \"" + username + "\" tried to unban an IP, but they are not staff/admin.")
		websocketError(s, d.Command, "Only staff members or administrators can do that.")
		return
	}

	// Validate that the requested IP is sane
	if ip == "" {
		log.Warning("User \"" + username + "\" tried to unban a blank IP.")
		websocketError(s, d.Command, "That IP is not valid.")
		return
	}

	// Validate that the requested IP is not already banned
	if IPBanned, err := db.BannedIPs.Check(ip); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if !IPBanned {
		websocketError(s, d.Command, "That IP is not banned.")
		return
	}

	// Remove the IP from the list in the database
	if err := db.BannedIPs.DeleteIP(ip); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Log the unban
	log.Info("User \"" + username + "\" unbanned IP \"" + ip + "\".")
}

func websocketAdminMute(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	functionName := "adminMute"
	userID := conn.UserID
	username := conn.Username
	recipient := data.Name

	// Validate that the user is staff/admin
	if conn.Admin == 0 {
		log.Warning("User \"" + username + "\" tried to mute \"" + recipient + "\", but they are not staff/admin.")
		websocketError(s, d.Command, "Only staff members and administrators can do that.")
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		log.Warning("User \"" + username + "\" tried to mute a blank person.")
		websocketError(s, d.Command, "That person is not valid.")
		return
	}

	// Validate that the requested person exists in the database
	if userExists, err := db.Users.Exists(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if !userExists {
		websocketError(s, d.Command, "That user does not exist.")
		return
	}

	// Validate that the requested person is not a staff member or an administrator
	if userIsStaff, err := db.Users.CheckStaff(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if userIsStaff {
		log.Warning("User \"" + username + "\" tried to mute \"" + recipient + "\", but staff/admins cannot be muted.")
		websocketError(s, d.Command, "You cannot mute a staff member or an administrator.")
		return
	}

	// Validate that they are not already muted
	if userIsMuted, err := db.MutedUsers.Check(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if userIsMuted {
		websocketError(s, d.Command, "That user is already muted.")
		return
	}

	// Add this username to the muted list in the database
	if err := db.MutedUsers.Insert(recipient, userID); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Check to see if this user is currently connected
	_, ok := connectionMap.m[recipient]
	if ok {
		// Update their connection map setting to be muted
		connectionMap.m[recipient].Muted = 1

		// Look for this user in all chat rooms
		for room, users := range chatRoomMap.m {
			// See if the user is in this chat room
			index := -1
			for i, user := range users {
				if user.Name == username {
					index = i
					break
				}
			}
			if index != -1 {
				// Update them to be muted
				chatRoomMap.m[room][index].Muted = 1

				// Send everyone an room update
				users, ok := chatRoomMap.m[room]
				if !ok {
					log.Error("Failed to retrieve the user list from the chat room map for room \"" + room + "\".")
					continue
				}

				for _, user := range users {
					userConnection, ok := connectionMap.m[user.Name] // This should always succeed, but there might be a race condition
					if ok {
						userConnection.Connection.Emit("roomSetMuted", &RoomSetMutedMessage{room, recipient, 1})
					} else {
						log.Error("Failed to get the connection for user \"" + user.Name + "\" while muting user \"" + recipient + "\".")
						continue
					}
				}
			}
		}
	}

	// Log the mute
	log.Info("User \"" + username + "\" muted user \"" + recipient + "\".")
}

func websocketAdminUnmute(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	functionName := "adminUnmute"
	username := conn.Username
	recipient := data.Name

	// Validate that the user is staff/admin
	if conn.Admin == 0 {
		log.Warning("User \"" + username + "\" tried to mute someone, but they are not staff/admin.")
		websocketError(s, d.Command, "Only staff members and administrators can do that.")
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		log.Warning("User \"" + username + "\" tried to unmute a blank person.")
		websocketError(s, d.Command, "That person is not valid.")
		return
	}

	// Validate that the requested person exists in the database
	if userExists, err := db.Users.Exists(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if !userExists {
		websocketError(s, d.Command, "That user does not exist.")
		return
	}

	// Validate that the requested person is not a staff member or an administrator
	if userIsStaff, err := db.Users.CheckStaff(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if userIsStaff {
		log.Warning("User \"" + username + "\" tried to unmute \"" + recipient + "\", but staff/admins cannot be unmuted.")
		websocketError(s, d.Command, "You cannot unmute a staff member or an administrator.")
		return
	}

	// Validate that they are muted
	if userIsMuted, err := db.MutedUsers.Check(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if !userIsMuted {
		websocketError(s, d.Command, "That user is not muted.")
		return
	}

	// Remove this username from the muted list in the database
	if err := db.MutedUsers.Delete(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Check to see if this user is currently connected
	_, ok := connectionMap.m[recipient]
	if ok {
		// Update their connection map setting to be unmuted
		connectionMap.m[recipient].Muted = 0

		// Look for this user in all chat rooms
		for room, users := range chatRoomMap.m {
			// See if the user is in this chat room
			index := -1
			for i, user := range users {
				if user.Name == username {
					index = i
					break
				}
			}
			if index != -1 {
				// Update them to be unmuted
				chatRoomMap.m[room][index].Muted = 0

				// Send everyone an room update
				users, ok := chatRoomMap.m[room]
				if !ok {
					log.Error("Failed to retrieve the user list from the chat room map for room \"" + room + "\".")
					continue
				}

				for _, user := range users {
					userConnection, ok := connectionMap.m[user.Name] // This should always succeed, but there might be a race condition
					if ok {
						userConnection.Connection.Emit("roomSetMuted", &RoomSetMutedMessage{room, recipient, 0})
					} else {
						log.Error("Failed to get the connection for user \"" + user.Name + "\" while unmuting user \"" + recipient + "\".")
						continue
					}
				}
			}
		}
	}

	// Log the unmute
	log.Info("User \"" + username + "\" unmuted user \"" + recipient + "\".")
}

func websocketAdminPromote(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	functionName := "adminPromote"
	username := conn.Username
	recipient := data.Name

	// Validate that the user is an admin
	if conn.Admin != 2 {
		log.Warning("User \"" + username + "\" tried to promote someone, but they are not an administrator.")
		websocketError(s, d.Command, "Only administrators can do that.")
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		log.Warning("User \"" + username + "\" tried to promote a blank person.")
		websocketError(s, d.Command, "That person is not valid.")
		return
	}

	// Validate that the requested person exists in the database
	if userExists, err := db.Users.Exists(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if !userExists {
		websocketError(s, d.Command, "That user does not exist.")
		return
	}

	// Validate that the requested person is not a staff member or an administrator
	if userIsStaff, err := db.Users.CheckStaff(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if userIsStaff {
		log.Warning("User \"" + username + "\" tried to promote \"" + recipient + "\", but they are already staff/admin.")
		websocketError(s, d.Command, "That user is already a staff member or an administrator.")
		return
	}

	// Set them to be a staff member
	if err := db.Users.SetAdmin(recipient, 1); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Check to see if this user is currently connected
	_, ok := connectionMap.m[recipient]
	if ok {
		// Update their connection map setting to be an admin
		connectionMap.m[recipient].Admin = 1

		// Look for this user in all chat rooms
		for room, users := range chatRoomMap.m {
			// See if the user is in this chat room
			index := -1
			for i, user := range users {
				if user.Name == username {
					index = i
					break
				}
			}
			if index != -1 {
				// Update them to be an admin
				chatRoomMap.m[room][index].Admin = 1

				// Send everyone an room update
				users, ok := chatRoomMap.m[room]
				if !ok {
					log.Error("Failed to retrieve the user list from the chat room map for room \"" + room + "\".")
					continue
				}

				for _, user := range users {
					userConnection, ok := connectionMap.m[user.Name] // This should always succeed, but there might be a race condition
					if ok {
						userConnection.Connection.Emit("roomSetAdmin", &RoomSetAdminMessage{room, recipient, 1})
					} else {
						log.Error("Failed to get the connection for user \"" + user.Name + "\" while promoting user \"" + recipient + "\".")
						continue
					}
				}
			}
		}
	}

	// Log the promotion
	log.Info("User \"" + username + "\" promoted \"" + recipient + "\" to be a staff member.")
}

func websocketAdminDemote(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	functionName := "adminDemote"
	username := conn.Username
	recipient := data.Name

	// Validate that the user is an admin
	if conn.Admin != 2 {
		log.Warning("User \"" + username + "\" tried to demote someone, but they are not an administrator.")
		websocketError(s, d.Command, "Only administrators can do that.")
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		log.Warning("User \"" + username + "\" tried to demote a blank person.")
		websocketError(s, d.Command, "That person is not valid.")
		return
	}

	// Validate that the requested person exists in the database
	if userExists, err := db.Users.Exists(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if !userExists {
		websocketError(s, d.Command, "That user does not exist.")
		return
	}

	// Validate that the requested person is not a staff member or an administrator
	if userIsStaff, err := db.Users.CheckStaff(recipient); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else if !userIsStaff {
		log.Warning("User \"" + username + "\" tried to demote \"" + recipient + "\", but they not staff/admin.")
		websocketError(s, d.Command, "That user is not a staff member or an administrator.")
		return
	}

	// Set their admin status to 0
	if err := db.Users.SetAdmin(recipient, 0); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	}

	// Check to see if this user is currently connected
	_, ok := connectionMap.m[recipient]
	if ok {
		// Update their connection map setting to be a normal user
		connectionMap.m[recipient].Admin = 0

		// Look for this user in all chat rooms
		for room, users := range chatRoomMap.m {
			// See if the user is in this chat room
			index := -1
			for i, user := range users {
				if user.Name == username {
					index = i
					break
				}
			}
			if index != -1 {
				// Update them to be a normal user
				chatRoomMap.m[room][index].Admin = 0

				// Send everyone an room update
				users, ok := chatRoomMap.m[room]
				if !ok {
					log.Error("Failed to retrieve the user list from the chat room map for room \"" + room + "\".")
					continue
				}

				for _, user := range users {
					userConnection, ok := connectionMap.m[user.Name] // This should always succeed, but there might be a race condition
					if ok {
						userConnection.Connection.Emit("roomSetAdmin", &RoomSetAdminMessage{room, recipient, 0})
					} else {
						log.Error("Failed to get the connection for user \"" + user.Name + "\" while demoting user \"" + recipient + "\".")
						continue
					}
				}
			}
		}
	}

	// Log the demotion
	log.Info("User \"" + username + "\" demoted \"" + recipient + "\" to a normal user.")
}

func websocketAdminMessage(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	functionName := "adminMessage"
	username := conn.Username
	message := data.Message

	// Validate that the user is an admin
	if conn.Admin != 2 {
		log.Warning("User \"" + username + "\" tried to send a server broadcast, but they are not an administrator.")
		websocketError(s, d.Command, "Only administrators can do that.")
		return
	}

	// Send everyone the server broadcast notification
	for _, conn := range connectionMap.m {
		conn.Connection.Emit("adminMessage", &RoomMessageMessage{
			Message: message,
		})
	}
}
*/
