package main

import (
	melody "gopkg.in/olahol/melody.v1"
)

// Prepare some data about all of the ongoing racers to send to a user who just joined the race
// (or just reconnected after a disconnect)
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
		CharacterNum         int     `json:"characterNum"`
		Place                int     `json:"place"`
		PlaceMid             int     `json:"placeMid"`
		DatetimeFinished     int64   `json:"datetimeFinished"`
		RunTime              int64   `json:"runTime"` // In milliseconds, reported by the mod
		Comment              string  `json:"comment"`
	}
	racers := make([]RacerMessage, 0)
	for _, racer := range race.Racers {
		racers = append(racers, RacerMessage{
			Name:                 racer.Name,
			DatetimeJoined:       racer.DatetimeJoined,
			Status:               racer.Status,
			FloorNum:             racer.FloorNum,
			StageType:            racer.StageType,
			DatetimeArrivedFloor: racer.DatetimeArrivedFloor,
			Items:                racer.Items,
			StartingItem:         racer.StartingItem,
			CharacterNum:         racer.CharacterNum,
			Place:                racer.Place,
			PlaceMid:             racer.PlaceMid,
			DatetimeFinished:     racer.DatetimeFinished,
			RunTime:              racer.RunTime,
			Comment:              racer.Comment,
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
