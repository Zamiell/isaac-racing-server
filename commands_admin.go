package main

/*
 *  WebSocket admin command functions
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
			connectionMap.RLock()
			for _, conn := range connectionMap.m {
				conn.Connection.Emit("raceLeft", RaceMessage{raceID, recipient})
			}
			connectionMap.RUnlock()
		} else {
			// Set this racer's status to disqualified
			if err := db.RaceParticipants.SetStatus(recipient, raceID, "disqualified"); err != nil {
				commandMutex.Unlock()
				log.Error("Database error:", err)
				connError(conn, functionName, "Something went wrong. Please contact an administrator.")
				return
			}

			// Get the list of racers for this race
			racerList, err := db.RaceParticipants.GetRacerList(raceID)
			if err != nil {
				commandMutex.Unlock()
				log.Error("Database error:", err)
				return
			}

			// Send a notification to all the people in this particular race that the user got disqualified
			connectionMap.RLock()
			for _, racer := range racerList {
				conn, ok := connectionMap.m[racer.Name]
				if ok == true { // Not all racers may be online during a race
					conn.Connection.Emit("racerSetStatus", &RacerSetStatusMessage{raceID, username, "disqualified"})
				}
			}
			connectionMap.RUnlock()
		}

		// Check to see if the race should start or finish
		raceCheckStartFinish(raceID)
	}

	// Check to see if the user is online
	connectionMap.RLock()
	bannedConnection, ok := connectionMap.m[recipient]
	connectionMap.RUnlock()
	if ok == true {
		// Disconnect the user
		connError(bannedConnection, functionName, "You have been banned. If you think this was a mistake, please contact the administration to appeal.")
		bannedConnection.Connection.Close()
	}

	// Log the ban
	log.Info("User \"" + username + "\" banned user \"" + recipient + "\".")

	// Send success confirmation
	connSuccess(conn, functionName, data)

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

	// Send success confirmation
	connSuccess(conn, functionName, data)

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

	// Send success confirmation
	connSuccess(conn, functionName, data)

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

	// Send success confirmation
	connSuccess(conn, functionName, data)

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func adminSquelch(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "adminSquelch"
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
		log.Warning("User \"" + username + "\" tried to squelch \"" + recipient + "\", but they are not staff/admin.")
		connError(conn, functionName, "Only staff members and administrators can do that.")
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to squelch a blank person.")
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
		log.Warning("User \"" + username + "\" tried to squelch \"" + recipient + "\", but staff/admins cannot be squelched.")
		connError(conn, functionName, "You cannot squelch a staff member or an administrator.")
		return
	}

	// Validate that they are not already squelched
	if userIsSquelched, err := db.SquelchedUsers.Check(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userIsSquelched == true {
		commandMutex.Unlock()
		connError(conn, functionName, "That user is already squelched.")
		return
	}

	// Add this username to the squelched list in the database
	if err := db.SquelchedUsers.Insert(recipient, userID); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Check to see if this user is currently connected
	connectionMap.RLock()
	_, ok := connectionMap.m[recipient]
	connectionMap.RUnlock()
	if ok == true {
		// Update their connection map setting to be squelched
		connectionMap.Lock()
		connectionMap.m[recipient].Squelched = 1
		connectionMap.Unlock()

		// Look for this user in all chat rooms
		chatRoomMap.Lock()
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
				// Update them to be squelched
				chatRoomMap.m[room][index].Squelched = 1

				// Send everyone an room update
				users, ok := chatRoomMap.m[room]
				if ok == false {
					log.Error("Failed to retrieve the user list from the chat room map for room \"" + room + "\".")
					continue
				}

				connectionMap.RLock()
				for _, user := range users {
					userConnection, ok := connectionMap.m[user.Name] // This should always succeed, but there might be a race condition
					if ok == true {
						userConnection.Connection.Emit("roomSetSquelched", &RoomSetSquelchedMessage{room, recipient, 1})
					} else {
						log.Error("Failed to get the connection for user \"" + user.Name + "\" while squelching user \"" + recipient + "\".")
						continue
					}
				}
				connectionMap.RUnlock()
			}
		}
		chatRoomMap.Unlock()
	}

	// Log the squelch
	log.Info("User \"" + username + "\" squelched user \"" + recipient + "\".")

	// Send success confirmation
	connSuccess(conn, functionName, data)

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}

func adminUnsquelch(conn *ExtendedConnection, data *IncomingCommandMessage) {
	// Local variables
	functionName := "adminUnsquelch"
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
		log.Warning("User \"" + username + "\" tried to squelch someone, but they are not staff/admin.")
		connError(conn, functionName, "Only staff members and administrators can do that.")
		return
	}

	// Validate that the requested person is sane
	if recipient == "" {
		commandMutex.Unlock()
		log.Warning("User \"" + username + "\" tried to squelch a blank person.")
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
		log.Warning("User \"" + username + "\" tried to unsquelch \"" + recipient + "\", but staff/admins cannot be unsquelched.")
		connError(conn, functionName, "You cannot unsquelch a staff member or an administrator.")
		return
	}

	// Validate that they are not already unsquelched
	if userIsSquelched, err := db.SquelchedUsers.Check(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	} else if userIsSquelched == false {
		commandMutex.Unlock()
		connError(conn, functionName, "That user is not squelched.")
		return
	}

	// Remove this username from the squelched list in the database
	if err := db.SquelchedUsers.Delete(recipient); err != nil {
		commandMutex.Unlock()
		log.Error("Database error:", err)
		connError(conn, functionName, "Something went wrong. Please contact an administrator.")
		return
	}

	// Check to see if this user is currently connected
	connectionMap.RLock()
	_, ok := connectionMap.m[recipient]
	connectionMap.RUnlock()
	if ok == true {
		// Update their connection map setting to be unsquelched
		connectionMap.Lock()
		connectionMap.m[recipient].Squelched = 0
		connectionMap.Unlock()

		// Look for this user in all chat rooms
		chatRoomMap.Lock()
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
				// Update them to be squelched
				chatRoomMap.m[room][index].Squelched = 0

				// Send everyone an room update
				users, ok := chatRoomMap.m[room]
				if ok == false {
					log.Error("Failed to retrieve the user list from the chat room map for room \"" + room + "\".")
					continue
				}

				connectionMap.RLock()
				for _, user := range users {
					userConnection, ok := connectionMap.m[user.Name] // This should always succeed, but there might be a race condition
					if ok == true {
						userConnection.Connection.Emit("roomSetSquelched", &RoomSetSquelchedMessage{room, recipient, 0})
					} else {
						log.Error("Failed to get the connection for user \"" + user.Name + "\" while unsquelching user \"" + recipient + "\".")
						continue
					}
				}
				connectionMap.RUnlock()
			}
		}
		chatRoomMap.Unlock()
	}

	// Log the unsquelch
	log.Info("User \"" + username + "\" unsquelched user \"" + recipient + "\".")

	// Send success confirmation
	connSuccess(conn, functionName, data)

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
	connectionMap.RLock()
	_, ok := connectionMap.m[recipient]
	connectionMap.RUnlock()
	if ok == true {
		// Update their connection map setting to be an admin
		connectionMap.Lock()
		connectionMap.m[recipient].Admin = 1
		connectionMap.Unlock()

		// Look for this user in all chat rooms
		chatRoomMap.Lock()
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

				connectionMap.RLock()
				for _, user := range users {
					userConnection, ok := connectionMap.m[user.Name] // This should always succeed, but there might be a race condition
					if ok == true {
						userConnection.Connection.Emit("roomSetAdmin", &RoomSetAdminMessage{room, recipient, 1})
					} else {
						log.Error("Failed to get the connection for user \"" + user.Name + "\" while promoting user \"" + recipient + "\".")
						continue
					}
				}
				connectionMap.RUnlock()
			}
		}
		chatRoomMap.Unlock()
	}

	// Log the promotion
	log.Info("User \"" + username + "\" promoted \"" + recipient + "\" to be a staff member.")

	// Send success confirmation
	connSuccess(conn, functionName, data)

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
	connectionMap.RLock()
	_, ok := connectionMap.m[recipient]
	connectionMap.RUnlock()
	if ok == true {
		// Update their connection map setting to be a normal user
		connectionMap.Lock()
		connectionMap.m[recipient].Admin = 0
		connectionMap.Unlock()

		// Look for this user in all chat rooms
		chatRoomMap.Lock()
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

				connectionMap.RLock()
				for _, user := range users {
					userConnection, ok := connectionMap.m[user.Name] // This should always succeed, but there might be a race condition
					if ok == true {
						userConnection.Connection.Emit("roomSetAdmin", &RoomSetAdminMessage{room, recipient, 0})
					} else {
						log.Error("Failed to get the connection for user \"" + user.Name + "\" while demoting user \"" + recipient + "\".")
						continue
					}
				}
				connectionMap.RUnlock()
			}
		}
		chatRoomMap.Unlock()
	}

	// Log the demotion
	log.Info("User \"" + username + "\" demoted \"" + recipient + "\" to a normal user.")

	// Send success confirmation
	connSuccess(conn, functionName, data)

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}
