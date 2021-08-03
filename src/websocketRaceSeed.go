package server

import (
	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceSeed(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username
	seed := d.Seed

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

	// Validate that the race has started
	if race.Status != RaceStatusInProgress {
		return
	}

	// Validate that they are in the race
	var racer *Racer
	if v, ok := race.Racers[username]; !ok {
		return
	} else {
		racer = v
	}

	// Validate that they are still racing
	if racer.Status != RacerStatusRacing {
		return
	}

	// Validate that they did not save and continue
	// (by checking to see if they are on the same seed that they were on before)
	if racer.Seed == seed {
		return
	}

	/*
		Add the seed
	*/

	racer.Seed = seed
	racer.Items = make([]*Item, 0) // Reset all of their accumulated items
	racer.StartingItem = 0
	racer.Rooms = make([]*Room, 0) // Reset all of their visited rooms
}
