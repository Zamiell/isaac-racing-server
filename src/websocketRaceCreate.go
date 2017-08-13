package main

import (
	"unicode/utf8"

	"github.com/Zamiell/isaac-racing-server/src/log"
	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceCreate(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username
	name := d.Name
	ruleset := d.Ruleset

	/*
		Validation
	*/

	// Validate that the race name cannot be empty
	if name == "" {
		name = "-"
	}

	// Validate that the race name is not longer than 100 characters
	if utf8.RuneCountInString(name) > 100 {
		log.Warning("User \"" + username + "\" sent a race name longer than 100 characters.")
		websocketError(s, d.Command, "Race names must not be longer than 100 characters.")
		return
	}

	// Validate that the ruleset options cannot be empty
	if ruleset.Format == "" {
		ruleset.Format = "unseeded"
	}
	if ruleset.Character == "" {
		ruleset.Character = "Judas"
	}
	if ruleset.Goal == "" {
		ruleset.Goal = "The Chest"
	}

	// Validate the submitted ruleset
	if !raceValidateRuleset(s, d) {
		return
	}

	// Check if there are any ongoing races with this name
	for _, race := range races {
		if race.Name == name {
			websocketError(s, d.Command, "There is already a non-finished race with that name.")
			return
		}
	}

	/*
		Create
	*/

	// Create and set a seed if necessary
	ruleset.Seed = "-"
	if ruleset.Format == "seeded" {
		// Create a random Isaac seed
		// (using the current Epoch timestamp as a seed)
		ruleset.Seed = isaacGetRandomSeed()
	} else if ruleset.Format == "diversity" {
		ruleset.Seed = diversityGetSeed()
	}

	/*
		Create the race in the database
		(it will have no data associated with it other than the automatically
		generated row ID; we want to use this ID as a unique map key)

		The benefit of doing this is that we won't reuse any race IDs after a
		server restart or crash.
		Furthermore, we want the ability for racers to be able to submit a race
		comment after the race has already ended. (Races are deleted from the
		internal map upon finishing.)
	*/
	var raceID int
	if v, err := db.Races.Insert(); err != nil {
		log.Error("Database error:", err)
		websocketError(s, d.Command, "")
		return
	} else {
		raceID = v
	}

	// Create the race and keep track of it in the races map
	race := &Race{
		ID:              raceID,
		Name:            name,
		Status:          "open",
		Ruleset:         ruleset,
		Captain:         username,
		SoundPlayed:     false,
		DatetimeCreated: getTimestamp(),
		DatetimeStarted: 0,
		Racers:          make(map[string]*Racer, 0),
	}
	races[raceID] = race

	// Send everyone a notification that a new race has been started
	for _, s := range websocketSessions {
		websocketEmit(s, "raceCreated", &RaceCreatedMessage{
			ID:              race.ID,
			Name:            race.Name,
			Status:          race.Status,
			Ruleset:         race.Ruleset,
			Captain:         race.Captain,
			DatetimeCreated: race.DatetimeCreated,
			DatetimeStarted: race.DatetimeStarted,
			Racers:          make([]string, 0),
		})
	}

	d.ID = race.ID
	websocketRaceJoin(s, d)
}
