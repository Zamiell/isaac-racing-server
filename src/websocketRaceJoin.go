package server

import (
	"strconv"

	melody "gopkg.in/olahol/melody.v1"
)

// This is also called manually by the "websocketRaceCreate" function
func websocketRaceJoin(s *melody.Session, d *IncomingWebsocketData) {
	userID := d.v.UserID
	username := d.v.Username
	raceID := d.ID

	/*
		Validation
	*/

	// Validate that the race exists
	var race *Race
	if v, ok := races[d.ID]; !ok {
		return
	} else {
		race = v
	}

	// Validate that the race is open
	if race.Status != RaceStatusOpen {
		return
	}

	// Validate that they are not in the race
	if _, ok := race.Racers[username]; ok {
		return
	}

	// Validate that we are not trying to join a solo race
	if race.Ruleset.Solo && len(race.Racers) > 0 {
		logger.Warning("User \"" + username + "\" attempted to call " + d.Command + " on race ID " + strconv.Itoa(raceID) + ", but it is a solo race.")
		websocketError(s, d.Command, "Race ID "+strconv.Itoa(raceID)+" is a solo race, so you cannot join it.")
		return
	}

	// Validate the password if the race is password protected
	if len(race.Password) > 0 && race.Password != d.Password {
		websocketWarning(s, d.Command, "That is not the correct password.")
		return
	}

	/*
		Join
	*/

	// Add this user to the race
	racer := &Racer{
		ID:             userID,
		Name:           username,
		DatetimeJoined: getTimestamp(),
		Status:         "not ready",
		FloorNum:       1,
		Items:          make([]*Item, 0),
		Rooms:          make([]*Room, 0),
		CharacterNum:   1,
		PlaceMid:       -1, // Will be set to the number of racers once the race starts
	}
	race.Racers[username] = racer

	// Send everyone a notification that the user joined
	for _, s := range websocketSessions {
		type RaceJoinedMessage struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		websocketEmit(s, "raceJoined", &RaceJoinedMessage{
			raceID,
			username,
		})
	}

	// Send them all the information about the racers in this race
	racerListMessage(s, race)

	// Join the user to the channel for that race
	d.Room = "_race_" + strconv.Itoa(raceID)
	websocketRoomJoinSub(s, d)

	// Send a reminder message to people playing ranked solo races
	if race.Ruleset.Ranked && race.Ruleset.Solo {
		type PrivateMessageMessage struct {
			Name    string `json:"name"`
			Message string `json:"message"`
		}
		websocketEmit(s, "privateMessage", &PrivateMessageMessage{
			"SERVER",
			"In order to prevent cheating, you must stream your races on Twitch or YouTube to be eligible for the Season 3 Ranked Solo leaderboards.",
		})
	}
}
