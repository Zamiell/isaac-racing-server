package main

import (
	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceUnready(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username

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

	// Validate that they are in the race
	var racer *Racer
	if v, ok := race.Racers[username]; !ok {
		return
	} else {
		racer = v
	}

	// Validate that their status is set to "ready"
	if racer.Status != RacerStatusReady {
		return
	}

	/*
		Unready
	*/

	race.SetRacerStatus(username, "not ready")
}
