package main

import (
	"github.com/Zamiell/isaac-racing-server/src/log"
)

/*
	Chat room subroutines
*/

func chatRoomsUpdate(username string, property string, newValue interface{}) {
	// Look for this user in all chat rooms
	for room, users := range chatRooms {
		// See if the user is in this chat room
		index := -1
		for i, user := range users {
			if user.Name == username {
				index = i
				break
			}
		}
		if index == -1 {
			// They are not in this chat room
			continue
		}

		// Update the property
		if property == "StreamURL" {
			streamURL := newValue.(string)
			chatRooms[room][index].StreamURL = streamURL
		} else {
			log.Error("The \"chatRoomsUpdate\" function was called without a valid property name.")
			return
		}

		// Send everyone in the room an update
		for _, user := range users {
			// All users in the chat room should be online, but check just in case
			if s2, ok := websocketSessions[user.Name]; ok {
				type RoomUpdateMessage struct {
					Room string `json:"room"`
					User User   `json:"user"`
				}
				websocketEmit(s2, "roomUpdate", &RoomUpdateMessage{
					room,
					chatRooms[room][index],
				})
			} else {
				log.Error("Failed to get the connection for user \"" + user.Name + "\" while setting a new chat room value for user \"" + username + "\".")
				continue
			}
		}
	}
}
