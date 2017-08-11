package main

/*
	Imports
*/

import (
	"strings"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/Zamiell/isaac-racing-server/src/models"
	melody "gopkg.in/olahol/melody.v1"
)

/*
	Command example:
	roomJoin {
		"room": "lobby",
	}
*/

func websocketRoomJoin(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username

	// Validate that the requested room is sane
	if d.Room == "" {
		log.Warning("User \"" + username + "\" tried to join a room without providing a room name.")
		websocketError(s, d.Command, "That is not a valid room name.")
		return
	}

	// Validate that they are not trying to join a system room
	if strings.HasPrefix(d.Room, "_") {
		log.Warning("User \"" + username + "\" tried to join a system room.")
		websocketError(s, d.Command, "You are not allowed to manually join system rooms.")
		return
	}

	// Validate that they are not already in the room
	users, ok := chatRooms[d.Room]
	if ok {
		// The room exists (at least 1 person is in it)
		userInRoom := false
		for _, user := range users {
			if user.Name == username {
				userInRoom = true
				break
			}
		}
		if userInRoom {
			log.Warning("User \"" + username + "\" tried to join a room they were already in.")
			websocketError(s, d.Command, "You are already in that room.")
			return
		}
	}

	// Let them join the room
	websocketRoomJoinSub(s, d)
}

/*
	Subroutines
*/

func websocketRoomJoinSub(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	room := d.Room
	username := d.v.Username
	admin := d.v.Admin
	muted := d.v.Muted
	streamURL := d.v.StreamURL

	// Add the user to the chat room map
	userObject := User{username, admin, muted, streamURL}
	chatRooms[room] = append(chatRooms[room], userObject)
	users := chatRooms[room] // Save the list of users in the room for later

	// Give the user the list of everyone in the chat room
	type RoomListMessage struct {
		Room  string `json:"room"`
		Users []User `json:"users"`
	}
	websocketEmit(s, "roomList", &RoomListMessage{room, users})

	// Tell everyone else that someone joined
	for _, user := range users {
		// All users in the chat room should be online, but check just in case
		if s2, ok := websocketSessions[user.Name]; ok {
			if user.Name == username {
				// We don't need to tell the person who just joined anything
				continue
			} else {
				// Send a notification that someone joined
				type RoomJoinedMessage struct {
					Room string `json:"room"`
					User User   `json:"user"`
				}
				websocketEmit(s2, "roomJoined", &RoomJoinedMessage{room, userObject})
			}
		} else {
			log.Error("Failed to get the connection for user \"" + user.Name + "\" while connecting user \"" + username + "\" to room \"" + room + "\".")
			continue
		}
	}

	// Get the chat history for this channel
	var roomHistoryList []models.RoomHistory
	if strings.HasPrefix(room, "_race_") {
		// Get all of the history
		// (in SQLite, LIMIT -1 returns all results)
		if list, err := db.ChatLog.Get(room, -1); err != nil {
			log.Error("Database error:", err)
			return
		} else {
			roomHistoryList = list
		}
	} else {
		// Get only the last 50 entries
		if list, err := db.ChatLog.Get(room, 50); err != nil {
			log.Error("Database error:", err)
			return
		} else {
			roomHistoryList = list
		}
	}

	// Send the chat history
	type RoomHistoryMessage struct {
		Room    string               `json:"room"`
		History []models.RoomHistory `json:"history"`
	}
	websocketEmit(s, "roomHistory", &RoomHistoryMessage{room, roomHistoryList})

	// Log the join
	log.Info("User \"" + username + "\" joined room: #" + room)
}
