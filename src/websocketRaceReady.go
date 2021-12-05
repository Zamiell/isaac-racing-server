package server

import (
	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceReady(s *melody.Session, d *IncomingWebsocketData) {
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

	// Validate that their status is set to "not ready"
	if racer.Status != "not ready" {
		return
	}

	/*
		Ready
	*/

	race.SetRacerStatus(username, RacerStatusReady)
	race.CheckStart()
}
