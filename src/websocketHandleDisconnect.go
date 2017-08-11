package main

import (
	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

func websocketHandleDisconnect(s *melody.Session) {
	// Local variables
	d := &IncomingWebsocketData{}
	d.Command = "websocketHandleDisconnect"
	if !websocketGetSessionValues(s, d) {
		log.Error("Did not complete the \"" + d.Command + "\" function. There is now likely orphaned entries in various data structures.")
		return
	}
	userID := d.v.UserID
	username := d.v.Username

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()
	defer commandMutex.Unlock()

	// Delete the connection from the session map
	delete(websocketSessions, username) // This will do nothing if the entry doesn't exist

	// Leave all the chat rooms that this person is in
	for room, users := range chatRooms {
		for _, user := range users {
			if user.Name == username {
				d.Room = room
				websocketRoomLeaveSub(s, d)
				break
			}
		}
	}

	// Check to see if this user is in any races that are not already in progress
	raceIDs, err := db.RaceParticipants.GetNotStartedRaces(userID)
	if err != nil {
		log.Error("Database error:", err)
		return
	}

	// Iterate over the races that they are currently in
	for _, raceID := range raceIDs {
		// Remove this user from the participants list for that race
		if err := db.RaceParticipants.Delete(username, raceID); err != nil {
			log.Error("Database error:", err)
			return
		}

		// Send everyone a notification that the user left the race
		for _, s := range websocketSessions {
			websocketEmit(s, "raceLeft", &RaceLeftMessage{raceID, username})
		}

		// Check to see if the race should start
		raceCheckStart(raceID)
	}

	// Log the disconnection
	log.Info("User \""+username+"\" disconnected;", len(websocketSessions), "user(s) now connected.")
}
