package main

/*
	WebSocket debug command functions
*/

func debug(conn *ExtendedConnection) {
	// Local variables
	functionName := "debug"
	username := conn.Username

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()

	// Rate limit all commands
	if commandRateLimit(conn) == true {
		return
	}

	// Validate that the user is an admin
	if conn.Admin == 0 {
		commandMutex.Unlock()
		log.Info("User \"" + username + "\" tried to do a debug command, but they are not staff/admin.")
		connError(conn, functionName, "Only staff members or administrators can do that.")
		return
	}

	// Print out the connection map
	/*
		connectionMap.RLock()
		fmt.Println(connectionMap.m)
		for _, conn := range connectionMap.m {
			fmt.Println("on connection:", conn.Username)
		}
		connectionMap.RUnlock()
	*/

	// Test IRC stuff
	//ircSend("JOIN #zamiell")

	// The command is over, so unlock the command mutex
	commandMutex.Unlock()
}
