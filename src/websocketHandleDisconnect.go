package main

import (
	melody "gopkg.in/olahol/melody.v1"
)

func websocketHandleDisconnect(s *melody.Session) {
	// Local variables
	d := &IncomingWebsocketData{}
	d.Command = "websocketHandleDisconnect"
	if !websocketGetSessionValues(s, d) {
		logger.Error("Did not complete the \"" + d.Command + "\" function. There is now likely orphaned entries in various data structures.")
		return
	}
	username := d.v.Username

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()
	defer commandMutex.Unlock()

	// Eject this player from any races that have not started yet
	for _, race := range races {
		if race.Status != "open" {
			continue
		}

		if _, ok := race.Racers[username]; ok {
			d.ID = race.ID
			websocketRaceLeave(s, d)
		}
	}

	// Leave all the chat rooms that this person is in
	// (we want this part after the race ejection because that step involves leaving rooms)
	// (at this point the user should only be in the lobby, but iterate through all of the chat rooms to make sure)
	for room, users := range chatRooms {
		for _, user := range users {
			if user.Name == username {
				d.Room = room
				websocketRoomLeaveSub(s, d)
				break
			}
		}
	}

	// Delete the connection from the session map
	delete(websocketSessions, username)

	// Log the disconnection
	logger.Info("User \""+username+"\" disconnected;", len(websocketSessions), "user(s) now connected.")
}
