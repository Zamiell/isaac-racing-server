package main

import (
	melody "gopkg.in/olahol/melody.v1"
)

func websocketRaceRoom(s *melody.Session, d *IncomingWebsocketData) {
	// Local variables
	username := d.v.Username
	roomID := d.RoomID

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
		Add the room
	*/

	room := &Room{
		roomID,
		racer.FloorNum,
		racer.StageType,
	}
	racer.Rooms = append(racer.Rooms, room)
}
