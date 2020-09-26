package main

import (
	"io/ioutil"
	"path"
	"sort"
	"strconv"

	melody "gopkg.in/olahol/melody.v1"
)

func websocketHandleConnect(s *melody.Session) {
	// Local variables
	d := &IncomingWebsocketData{}
	d.Command = "websocketHandleConnect"
	if !websocketGetSessionValues(s, d) {
		logger.Error("Did not complete the \"" + d.Command + "\" function.")
		websocketClose(s)
		return
	}
	userID := d.v.UserID
	username := d.v.Username
	streamURL := d.v.StreamURL
	twitchBotEnabled := d.v.TwitchBotEnabled
	twitchBotDelay := d.v.TwitchBotDelay

	/*
		Establish the WebSocket session
	*/

	// Lock the command mutex for the duration of the function to ensure synchronous execution
	commandMutex.Lock()
	defer commandMutex.Unlock()

	// Disconnect any existing connections with this username
	if s2, ok := websocketSessions[username]; ok {
		logger.Info("Closing existing connection for user \"" + username + "\".")
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
	logger.Info("User \""+username+"\" connected;", len(websocketSessions), "user(s) now connected.")

	// Send them various settings tied to their account
	type SettingsMessage struct {
		UserID           int    `json:"userID"`
		Username         string `json:"username"`
		StreamURL        string `json:"streamURL"`
		TwitchBotEnabled bool   `json:"twitchBotEnabled"`
		TwitchBotDelay   int    `json:"twitchBotDelay"`
		Time             int64  `json:"time"`
	}
	websocketEmit(s, "settings", &SettingsMessage{
		UserID:           userID,
		Username:         username,
		StreamURL:        streamURL,
		TwitchBotEnabled: twitchBotEnabled,
		TwitchBotDelay:   twitchBotDelay,
	})

	// Prepare some data about all of the ongoing races to send to the newly
	// connected user
	// (we only want to send the client a subset of the race information in
	// order to conserve bandwidth and hide some things that they don't need to
	// see)
	// https://stackoverflow.com/questions/18342784/how-to-iterate-through-a-map-in-golang-in-order/18342865
	raceIDs := make([]int, 0)
	for id := range races {
		raceIDs = append(raceIDs, id)
	}
	sort.Ints(raceIDs)
	raceListMessage := make([]RaceCreatedMessage, 0)
	for _, id := range raceIDs {
		race := races[id]
		msg := RaceCreatedMessage{
			ID:                  race.ID,
			Name:                race.Name,
			Status:              race.Status,
			Ruleset:             race.Ruleset,
			Captain:             race.Captain,
			IsPasswordProtected: len(race.Password) > 0,
			DatetimeCreated:     race.DatetimeCreated,
			DatetimeStarted:     race.DatetimeStarted,
		}
		racers := make([]string, 0)
		for racerName := range race.Racers {
			racers = append(racers, racerName)
		}
		msg.Racers = racers

		raceListMessage = append(raceListMessage, msg)
	}

	// Send it to the user
	websocketEmit(s, "raceList", raceListMessage)

	// Check to see if this user is in any ongoing races
	for _, id := range raceIDs {
		race := races[id]
		if _, ok := race.Racers[username]; !ok {
			// They are not in this race
			continue
		}

		// Join the user to the chat room coresponding to this race
		d.Room = "_race_" + strconv.Itoa(race.ID)
		websocketRoomJoinSub(s, d)

		// Send them all the information about the racers in this race
		racerListMessage(s, race)

		// If the race is currently in the 10 second countdown
		if race.Status == "starting" {
			// Send them a message describing when it will start
			websocketEmit(s, "raceStart", &RaceStartMessage{
				ID:            race.ID,
				SecondsToWait: 10,
				// This will make them start behind the other racers,
				// but it gives them 10 seconds to get ready after a disconnect;
				// times are reported via client side start and finish anyway
			})
		}
	}

	// Send them the message(s) of the day
	websocketEmit(s, "adminMessage", &AdminMessageMessage{
		"[Server Notice] Most racers hang out in the Isaac Discord chat: https://discord.gg/JzbhWQb",
	})
	messageRaw, err := ioutil.ReadFile(path.Join(projectPath, "message_of_the_day.txt"))
	if err != nil {
		logger.Error("Failed to read the \"message_of_the_day.txt\" file:", err)
		return
	}
	message := string(messageRaw)
	if len(message) > 0 {
		websocketEmit(s, "adminMessage", &AdminMessageMessage{
			string(messageRaw),
		})
	}
}
