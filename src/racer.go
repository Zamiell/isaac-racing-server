package server

import (
	melody "gopkg.in/olahol/melody.v1"
)

/*
	Data structures
*/

type Racer struct {
	ID                   int
	Name                 string
	DatetimeJoined       int64
	Status               RacerStatus
	Seed                 string
	FloorNum             int
	StageType            int
	BackwardsPath        bool
	DatetimeArrivedFloor int64 // Epoch timestamp in milliseconds
	Items                []*Item
	StartingItem         int
	Rooms                []*Room
	CharacterNum         int // Only used in multi-character races
	Place                int
	PlaceMid             int // -1 if quit or finished
	PlaceMidOld          int
	DatetimeFinished     int64
	RunTime              int64 // In milliseconds
	Comment              string
}

type Item struct {
	ID               int   `json:"id"`
	FloorNum         int   `json:"floorNum"`
	StageType        int   `json:"stageType"`
	DatetimeAcquired int64 `json:"datetimeAcquired"`
}

type Room struct {
	ID              string // e.g. "5.999"
	FloorNum        int
	StageType       int
	DatetimeArrived int64
}

// Prepare some data about all of the ongoing racers to send to a user who just joined the race
// (or just reconnected after a disconnect)
// (we only want to send the client a subset of the total information in
// order to conserve bandwidth)
// Called from "websocketHandleConnect" and "websocketRaceJoin"
func racerListMessage(s *melody.Session, race *Race) {
	type RacerMessage struct {
		Name                 string      `json:"name"`
		DatetimeJoined       int64       `json:"datetimeJoined"`
		Status               RacerStatus `json:"status"`
		FloorNum             int         `json:"floorNum"`
		StageType            int         `json:"stageType"`
		DatetimeArrivedFloor int64       `json:"datetimeArrivedFloor"`
		Items                []*Item     `json:"items"`
		StartingItem         int         `json:"startingItem"`
		CharacterNum         int         `json:"characterNum"`
		Place                int         `json:"place"`
		PlaceMid             int         `json:"placeMid"`
		DatetimeFinished     int64       `json:"datetimeFinished"`
		RunTime              int64       `json:"runTime"` // In milliseconds, reported by the mod
		Comment              string      `json:"comment"`
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
