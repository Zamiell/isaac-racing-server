package main

/*
	WebSocket admin command functions
*/

func adminBan(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "adminBan"
	userID := conn.UserID
	username := conn.Username
	recipient := data.Name

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command (for recipient \""+recipient+"\").")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the user is staff/admin
	if conn.Admin == 0 {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to ban someone, but they are not staff/admin.")
		connError(conn, functionName, "Only staff members or administrators can do that.")
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to ban a blank person.")
		connError(conn, functionName, "That person is not valid.")
		return
	}

	// Validate that the requested person exists in the database
	if userExists, err := db.Users.Exists(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userExists == false {
		commandMutex.Unlock()
		connError(conn, functionName, "That user does not exist.")
		return
	}

	// Validate that the requested person is not a staff member or an administrator
	if userIsStaff, err := db.Users.CheckStaff(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userIsStaff == true {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to ban \"" + recipient + "\", but staff/admins cannot be banned.")
		connError(conn, functionName, "You cannot ban a staff member or an administrator.")
		return
	}

	// Validate that the requested person is not already banned
	if userIsBanned, err := db.BannedUsers.Check(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userIsBanned == true {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to ban \"" + recipient + "\", but they are already banned.")
		connError(conn, functionName, "That user is already banned.")
		return
	}

	// Add this username to the ban list in the database
	if err := db.BannedUsers.Insert(recipient, userID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Add their IP to the banned IP list
	if err := db.BannedIPs.Insert(recipient, userID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Find out if the banned user is in any races that are currently going on
	raceList, err := db.RaceParticipants.GetCurrentRaces(recipient)
	if err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Iterate over the races that they are currently in
	for _, race := range raceList {
		raceID := race.ID

		// Find out if the race is started
		if race.Status == "open" {
			// Remove this user from the participants list for that race
			if err := db.RaceParticipants.Delete(recipient, raceID); err != nil {
				commandMutex.Unlock()
				log.Error("Database error:", err)
				connError(conn, functionName, "Something went wrong. Please contact an administrator.")
				return
			}

			// Send everyone a notification that the user left the race
			for _, conn := range connectionMap.m {
				conn.Connection.Emit("raceLeft", RaceMessage{raceID, recipient})
			}
		} else {
			// Set this racer's status to disqualified
			if err := db.RaceParticipants.SetStatus(recipient, raceID, "disqualified"); err != nil {
				commandMutex.Unlock()
				log.Error("Database error:", err)
				connError(conn, functionName, "Something went wrong. Please contact an administrator.")
				return
			}

			// Set their finish time
			if err := db.RaceParticipants.SetDatetimeFinished(recipient, raceID, int(makeTimestamp())); err != nil {
				commandMutex.Unlock()
				log.Error("Database error:", err)
				connError(conn, functionName, "Something went wrong. Please contact an administrator.")
				return
			}

			// Set their (final) place to -2 (which indicates a disqualified status)
			if err := db.RaceParticipants.SetPlace(username, raceID, -2); err != nil {
				commandMutex.Unlock()
				log.Error("Database error:", err)
				connError(conn, functionName, "Something went wrong. Please contact an administrator.")
				return
			}

			// Get the list of racers for this race
			racerNames, err := db.RaceParticipants.GetRacerNames(raceID)
			if err != nil {
				commandMutex.Unlock()
				log.Error("Database error:", err)
				return
			}

			// Send a notification to all the people in this particular race that the user got disqualified
			for _, racer := range racerNames {
				conn, ok := connectionMap.m[racer]
				if ok == true { // Not all racers may be online during a race
					conn.Connection.Emit("racerSetStatus", &RacerSetStatusMessage{raceID, username, "disqualified", -2})
				}
			}
		}

		// Check to see if the race should start or finish
		raceCheckStartFinish(raceID)
	}

	// Check to see if the user is online
	bannedConnection, ok := connectionMap.m[recipient]
	if ok == true {
		// Disconnect the user
		connError(bannedConnection, functionName, "You have been banned. If you think this was a mistake, please contact the administration to appeal.")
		bannedConnection.Connection.Close()
	}

	// Log the ban
	log.Info("User \"" + username + "\" banned user \"" + recipient + "\".")

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func adminUnban(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "adminUnban"
	username := conn.Username
	recipient := data.Name

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command (for recipient \""+recipient+"\").")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the user is staff/admin
	if conn.Admin == 0 {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to unban someone, but they are not staff/admin.")
		connError(conn, functionName, "Only staff members or administrators can do that.")
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to unban a blank person.")
		connError(conn, functionName, "That person is not valid.")
		return
	}

	// Validate that the requested person exists in the database
	if userExists, err := db.Users.Exists(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userExists == false {
		commandMutex.Unlock()
		connError(conn, functionName, "That user does not exist.")
		return
	}

	// Validate that the requested person is not a staff member or an administrator
	if userIsStaff, err := db.Users.CheckStaff(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userIsStaff == true {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to unban \"" + recipient + "\", but staff/admins cannot be unbanned.")
		connError(conn, functionName, "You cannot unban a staff member or an administrator.")
		return
	}

	// Validate that the requested person is banned
	if userIsBanned, err := db.BannedUsers.Check(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userIsBanned == false {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to unban \"" + recipient + "\", but they are not banned.")
		connError(conn, functionName, "That user is not banned.")
		return
	}

	// Remove this username from the ban list in the database
	if err := db.BannedUsers.Delete(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Remove the user's last IP from the banned IP list, if present
	if err := db.BannedIPs.Delete(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Log the unban
	log.Info("User \"" + username + "\" unbanned user \"" + recipient + "\".")

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func adminBanIP(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "adminBanIP"
	userID := conn.UserID
	username := conn.Username
	ip := data.IP

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command (for IP address \""+ip+"\").")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the user is staff/admin
	if conn.Admin == 0 {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to ban an IP, but they are not staff/admin.")
		connError(conn, functionName, "Only staff members or administrators can do that.")
		return
	}

	// Validate that the requested IP is sane
	if ip == "" {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to ban a blank IP.")
		connError(conn, functionName, "That IP is not valid.")
		return
	}

	// Validate that the requested IP is not already banned
	if IPBanned, err := db.BannedIPs.Check(ip); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if IPBanned == true {
		commandMutex.Unlock()
		connError(conn, functionName, "That IP is already banned.")
		return
	}

	// Add the IP to the list in the database
	if err := db.BannedIPs.InsertIP(ip, userID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Log the ban
	log.Info("User \"" + username + "\" banned IP \"" + ip + "\".")

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func adminUnbanIP(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "adminUnbanIP"
	username := conn.Username
	ip := data.IP

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command (for IP address \""+ip+"\").")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the user is staff/admin
	if conn.Admin == 0 {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to unban an IP, but they are not staff/admin.")
		connError(conn, functionName, "Only staff members or administrators can do that.")
		return
	}

	// Validate that the requested IP is sane
	if ip == "" {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to unban a blank IP.")
		connError(conn, functionName, "That IP is not valid.")
		return
	}

	// Validate that the requested IP is not already banned
	if IPBanned, err := db.BannedIPs.Check(ip); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if IPBanned == false {
		commandMutex.Unlock()
		connError(conn, functionName, "That IP is not banned.")
		return
	}

	// Remove the IP from the list in the database
	if err := db.BannedIPs.DeleteIP(ip); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Log the unban
	log.Info("User \"" + username + "\" unbanned IP \"" + ip + "\".")

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func adminMute(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "adminMute"
	userID := conn.UserID
	username := conn.Username
	recipient := data.Name

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command (for recipient \""+recipient+"\").")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the user is staff/admin
	if conn.Admin == 0 {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to mute \"" + recipient + "\", but they are not staff/admin.")
		connError(conn, functionName, "Only staff members and administrators can do that.")
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to mute a blank person.")
		connError(conn, functionName, "That person is not valid.")
		return
	}

	// Validate that the requested person exists in the database
	if userExists, err := db.Users.Exists(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userExists == false {
		commandMutex.Unlock()
		connError(conn, functionName, "That user does not exist.")
		return
	}

	// Validate that the requested person is not a staff member or an administrator
	if userIsStaff, err := db.Users.CheckStaff(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userIsStaff == true {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to mute \"" + recipient + "\", but staff/admins cannot be muted.")
		connError(conn, functionName, "You cannot mute a staff member or an administrator.")
		return
	}

	// Validate that they are not already muted
	if userIsMuted, err := db.MutedUsers.Check(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userIsMuted == true {
		commandMutex.Unlock()
		connError(conn, functionName, "That user is already muted.")
		return
	}

	// Add this username to the muted list in the database
	if err := db.MutedUsers.Insert(recipient, userID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Check to see if this user is currently connected
	_, ok := connectionMap.m[recipient]
	if ok == true {
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
				if ok == false {
					log.Error("Failed to retrieve the user list from the chat room map for room \"" + room + "\".")
					continue
				}

				for _, user := range users {
					userConnection, ok := connectionMap.m[user.Name] // This should always succeed, but there might be a race condition
					if ok == true {
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

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func adminUnmute(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "adminUnmute"
	username := conn.Username
	recipient := data.Name

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command (for recipient \""+recipient+"\").")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the user is staff/admin
	if conn.Admin == 0 {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to mute someone, but they are not staff/admin.")
		connError(conn, functionName, "Only staff members and administrators can do that.")
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to unmute a blank person.")
		connError(conn, functionName, "That person is not valid.")
		return
	}

	// Validate that the requested person exists in the database
	if userExists, err := db.Users.Exists(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userExists == false {
		commandMutex.Unlock()
		connError(conn, functionName, "That user does not exist.")
		return
	}

	// Validate that the requested person is not a staff member or an administrator
	if userIsStaff, err := db.Users.CheckStaff(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userIsStaff == true {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to unmute \"" + recipient + "\", but staff/admins cannot be unmuted.")
		connError(conn, functionName, "You cannot unmute a staff member or an administrator.")
		return
	}

	// Validate that they are muted
	if userIsMuted, err := db.MutedUsers.Check(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userIsMuted == false {
		commandMutex.Unlock()
		connError(conn, functionName, "That user is not muted.")
		return
	}

	// Remove this username from the muted list in the database
	if err := db.MutedUsers.Delete(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Check to see if this user is currently connected
	_, ok := connectionMap.m[recipient]
	if ok == true {
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
				if ok == false {
					log.Error("Failed to retrieve the user list from the chat room map for room \"" + room + "\".")
					continue
				}

				for _, user := range users {
					userConnection, ok := connectionMap.m[user.Name] // This should always succeed, but there might be a race condition
					if ok == true {
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

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func adminPromote(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "adminPromote"
	username := conn.Username
	recipient := data.Name

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command (for recipient \""+recipient+"\").")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the user is an admin
	if conn.Admin != 2 {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to promote someone, but they are not an administrator.")
		connError(conn, functionName, "Only administrators can do that.")
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to promote a blank person.")
		connError(conn, functionName, "That person is not valid.")
		return
	}

	// Validate that the requested person exists in the database
	if userExists, err := db.Users.Exists(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userExists == false {
		commandMutex.Unlock()
		connError(conn, functionName, "That user does not exist.")
		return
	}

	// Validate that the requested person is not a staff member or an administrator
	if userIsStaff, err := db.Users.CheckStaff(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userIsStaff == true {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to promote \"" + recipient + "\", but they are already staff/admin.")
		connError(conn, functionName, "That user is already a staff member or an administrator.")
		return
	}

	// Set them to be a staff member
	if err := db.Users.SetAdmin(recipient, 1); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Check to see if this user is currently connected
	_, ok := connectionMap.m[recipient]
	if ok == true {
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
				if ok == false {
					log.Error("Failed to retrieve the user list from the chat room map for room \"" + room + "\".")
					continue
				}

				for _, user := range users {
					userConnection, ok := connectionMap.m[user.Name] // This should always succeed, but there might be a race condition
					if ok == true {
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

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func adminDemote(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "adminDemote"
	username := conn.Username
	recipient := data.Name

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command (for recipient \""+recipient+"\").")

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the user is an admin
	if conn.Admin != 2 {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to demote someone, but they are not an administrator.")
		connError(conn, functionName, "Only administrators can do that.")
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to demote a blank person.")
		connError(conn, functionName, "That person is not valid.")
		return
	}

	// Validate that the requested person exists in the database
	if userExists, err := db.Users.Exists(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userExists == false {
		commandMutex.Unlock()
		connError(conn, functionName, "That user does not exist.")
		return
	}

	// Validate that the requested person is not a staff member or an administrator
	if userIsStaff, err := db.Users.CheckStaff(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userIsStaff == false {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to demote \"" + recipient + "\", but they not staff/admin.")
		connError(conn, functionName, "That user is not a staff member or an administrator.")
		return
	}

	// Set their admin status to 0
	if err := db.Users.SetAdmin(recipient, 0); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Check to see if this user is currently connected
	_, ok := connectionMap.m[recipient]
	if ok == true {
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
				if ok == false {
					log.Error("Failed to retrieve the user list from the chat room map for room \"" + room + "\".")
					continue
				}

				for _, user := range users {
					userConnection, ok := connectionMap.m[user.Name] // This should always succeed, but there might be a race condition
					if ok == true {
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

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func adminMessage(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "adminMessage"
	username := conn.Username
	message := data.Message

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Log the received command
	log.Debug("User \""+username+"\" sent a", functionName, "command.")

	// Validate that the user is an admin
	if conn.Admin != 2 {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to send a server broadcast, but they are not an administrator.")
		connError(conn, functionName, "Only administrators can do that.")
		return
	}

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Send everyone the server broadcast notification
	for _, conn := range connectionMap.m {
		conn.Connection.Emit("adminMessage", &RoomMessageMessage{
			Message: message,
		})
	}

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}
