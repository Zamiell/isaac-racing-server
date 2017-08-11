package main

/*
	Imports
*/

import (
	"io/ioutil"
	"path"
	"strconv"
	"time"

	"github.com/Zamiell/isaac-racing-server/src/log"
	"github.com/Zamiell/isaac-racing-server/src/models"
	melody "gopkg.in/olahol/melody.v1"
)

func websocketHandleConnect(s *melody.Session) {
	// Local variables
	d := &IncomingWebsocketData{}
	d.Command = "websocketHandleConnect"
	if !websocketGetSessionValues(s, d) {
		log.Error("Did not complete the \"" + d.Command + "\" function.")
		websocketClose(s)
		return
	}
	userID := d.v.UserID
	username := d.v.Username
	streamURL := d.v.StreamURL

	/*
		Establish the WebSocket session
	*/

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()
	defer commandMutex.Unlock()

	// Disconnect any existing connections with this username
	if s2, ok := websocketSessions[username]; ok {
		log.Info("Closing existing connection for user \"" + username + "\".")
		websocketError(s2, "logout", "You have logged on from somewhere else, so you have been disconnected here.")
		websocketClose(s2)

		// Wait until the existing connection is terminated
		commandMutex.Unlock()
		for {
			commandMutex.Lock()
			_, ok := websocketSessions[username]
			commandMutex.Unlock()
			if !ok {
				break
			}
		}
		commandMutex.Lock()
	}

	// Add the connection to a session map so that we can keep track of all of the connections
	websocketSessions[username] = s
	log.Info("User \""+username+"\" connected;", len(websocketSessions), "user(s) now connected.")

	// Get their Twitch bot settings
	var twitchBotEnabled bool
	if v, err := db.Users.GetTwitchBotEnabled(userID); err != nil {
		log.Error("Database error:", err)
		return
	} else {
		twitchBotEnabled = v
	}
	var twitchBotDelay int
	if v, err := db.Users.GetTwitchBotDelay(userID); err != nil {
		log.Error("Database error:", err)
		return
	} else {
		twitchBotDelay = v
	}

	// Send them various settings tied to their account
	type SettingsMessage struct {
		Username         string `json:"username"`
		StreamURL        string `json:"streamURL"`
		TwitchBotEnabled bool   `json:"twitchBotEnabled"`
		TwitchBotDelay   int    `json:"twitchBotDelay"`
		Time             int64  `json:"time"`
	}
	websocketEmit(s, "settings", &SettingsMessage{
		Username:         username,
		StreamURL:        streamURL,
		TwitchBotEnabled: twitchBotEnabled,
		TwitchBotDelay:   twitchBotDelay,
		// Send them the current time so that they can calculate the local offset
		Time: makeTimestamp(),
	})

	// Get the current list of races
	var raceList []models.Race
	if v, err := db.Races.GetCurrentRaces(); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else {
		raceList = v
	}

	// Send it to the user
	websocketEmit(s, "raceList", raceList)

	// Find out if the user is in any races that are currently going on
	for _, race := range raceList {
		for _, racer := range race.Racers {
			if racer == username {
				// Join the user to the chat room coresponding to this race
				d.Room = "_race_" + strconv.Itoa(race.ID)
				websocketRoomJoinSub(s, d)

				// Send them all the information about the racers in this race
				if racerList, err := db.RaceParticipants.GetRacerList(race.ID); err != nil {
					log.Error("Database error:", err)
					return
				} else {
					websocketEmit(s, "racerList", &RacerListMessage{race.ID, racerList})
				}

				// If the race is currently in the 10 second countdown
				if race.Status == "starting" {
					// Get the time 10 seconds in the future
					startTime := time.Now().Add(10*time.Second).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
					// This will technically put them behind the other racers by some amount of seconds, but it gives them 10 seconds to get ready after a disconnect

					// Send them a message describing exactly when it will start
					websocketEmit(s, "raceStart", &RaceStartMessage{race.ID, startTime})
				}

				break
			}
		}
	}

	// Send them the message(s) of the day
	type AdminMessageMessage struct {
		Message string `json:"message"`
	}
	websocketEmit(s, "adminMessage", &AdminMessageMessage{
		Message: "[Server Notice] Racing+ is in alpha and is NOT finished - lots of features are still missing or bugged.",
	})
	websocketEmit(s, "adminMessage", &AdminMessageMessage{
		Message: "[Server Notice] Most racers hang out in the Isaac Discord chat: https://discord.gg/JzbhWQb",
	})
	messageRaw, err := ioutil.ReadFile(path.Join(projectPath, "message_of_the_day.txt"))
	if err != nil {
		log.Error("Failed to read the \"message_of_the_day.txt\" file:", err)
		return
	}
	message := string(messageRaw)
	if len(message) > 0 {
		websocketEmit(s, "adminMessage", &AdminMessageMessage{
			Message: string(messageRaw),
		})
	}
}
