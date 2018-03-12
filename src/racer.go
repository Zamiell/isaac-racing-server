package main

import (
	melody "gopkg.in/olahol/melody.v1"
)

// Prepare some data about all of the ongoing racers to send to user who just
// join the race (or just reconnected after a disconnect)
// (we only want to send the client a subset of the total information in
// order to conserve bandwidth)
// Called from "websocketHandleConnect" and "websocketRaceJoin"
func racerListMessage(s *melody.Session, race *Race) {
	type RacerMessage struct {
		Name                 string  `json:"name"`
		DatetimeJoined       int64   `json:"datetimeJoined"`
		Status               string  `json:"status"`
		FloorNum             int     `json:"floorNum"`
		StageType            int     `json:"stageType"`
		DatetimeArrivedFloor int64   `json:"datetimeArrivedFloor"`
		Items                []*Item `json:"items"`
		StartingItem         int     `json:"startingItem"`
		Place                int     `json:"place"`
		PlaceMid             int     `json:"placeMid"`
		DatetimeFinished     int64   `json:"datetimeFinished"`
		RunTime              int64   `json:"runTime"` // In milliseconds, reported by the mod
		Comment              string  `json:"comment"`
	}
	racers := make([]RacerMessage, 0)
	for _, racer := range race.Racers {
		racers = append(racers, RacerMessage{
			racer.Name,
			racer.DatetimeJoined,
			racer.Status,
			racer.FloorNum,
			racer.StageType,
			racer.DatetimeArrivedFloor,
			racer.Items,
			racer.StartingItem,
			racer.Place,
			racer.PlaceMid,
			racer.DatetimeFinished,
			racer.RunTime,
			racer.Comment,
		})
	}

	type RacerListMessage struct {
		ID     int            `json:"id"`
		Racers []RacerMessage `json:"racers"`
	}
	websocketEmit(s, "racerList", &RacerListMessage{
		race.ID,
		racers,
	})
}
