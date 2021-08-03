package server

import (
	melody "gopkg.in/olahol/melody.v1"
)

// This is also called manually by the "race.Start3" function
func websocketRaceQuit(s *melody.Session, d *IncomingWebsocketData) {
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

	/*
		Quit
	*/

	racer.Place = -1
	racer.PlaceMid = -1
	race.SetRacerStatus(username, "quit")
	racer.DatetimeFinished = getTimestamp()
	racer.RunTime = racer.DatetimeFinished - race.DatetimeStarted
	race.SetAllPlaceMid()
	twitchRacerQuit(race, racer)
	race.CheckFinish()
}
