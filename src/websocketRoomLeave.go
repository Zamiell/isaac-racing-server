package main

import (
	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

/*
	Command example:
	roomLeave {
		room: "lobby",
	}
*/

func websocketRoomLeave(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username

	// Validate that the requested room is sane
	if d.Room == "" {
		log.Warning("User \"" + username + "\" tried to leave a room without providing a room name.")
		websocketError(s, d.Command, "That is not a valid room name.")
		return
	}

	// Validate that the room exists
	users, ok := chatRooms[d.Room]
	if !ok {
		log.Warning("User \"" + username + "\" tried to leave an invalid room.")
		websocketError(s, d.Command, "That is not a valid room name.")
		return
	}

	// Validate that they are actually in this room
	userInRoom := false
	for _, user := range users {
		if user.Name == username {
			userInRoom = true
			break
		}
	}
	if !userInRoom {
		log.Warning("User \"" + username + "\" tried to leave a room they were not in.")
		websocketError(s, d.Command, "You are not in that room.")
		return
	}

	// Let them leave the room
	websocketRoomLeaveSub(s, d)
}

/*
	Subroutines
*/

func websocketRoomLeaveSub(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username
	room := d.Room

	// Get the index of the user in the chat room mapping for this room
	users, ok := chatRooms[room]
	if !ok {
		log.Error("Failed to get the list of users for room \"" + room + "\".")
		return
	}
	index := -1
	for i, user := range users {
		if user.Name == username {
			index = i
			break
		}
	}
	if index == -1 {
		log.Error("Failed to get the index for the current user for room \"" + room + "\".")
		return
	}

	// Remove the user from the chat room map
	chatRooms[room] = append(users[:index], users[index+1:]...)
	users = chatRooms[room] // Save the list of users in the room for later
	if len(users) == 0 {
		// If there is no-one left in the room, remove the entry from the map entirely
		// (this prevents a memory leak)
		delete(chatRooms, room)
	}

	// Tell everyone else that someone left
	for _, user := range users {
		// All users in the chat room should be online, but check just in case
		if s2, ok := websocketSessions[user.Name]; ok {
			type RoomLeftMessage struct {
				Room string `json:"room"`
				Name string `json:"name"`
			}
			websocketEmit(s2, "roomLeft", &RoomLeftMessage{
				room,
				username,
			})
		} else {
			log.Error("Failed to get the connection for user \"" + user.Name + "\" while disconnecting user \"" + username + "\" from room \"" + room + "\".")
			continue
		}
	}

	// Log the leave
	log.Info("User \"" + username + "\" left room: #" + room)
}
