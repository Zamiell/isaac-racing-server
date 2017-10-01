package main

import (
	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceFinish(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	userID := d.v.UserID
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
	if race.Status != "in progress" {
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
	if racer.Status != "racing" {
		return
	}

	/*
		Finish
	*/

	racer.Place = race.GetCurrentPlace()
	race.SetRacerStatus(username, "finished")
	racer.DatetimeFinished = getTimestamp()
	racer.RunTime = racer.DatetimeFinished - race.DatetimeStarted
	race.SetAllPlaceMid()
	twitchRacerFinish(race, racer)
	race.CheckFinish()

	// Check to see if the user got any achievements
	// (which can only happen if they actually finished the race)
	achievementsCheck(userID, username)
}
