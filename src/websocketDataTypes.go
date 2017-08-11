package main

/*
	Imports
*/

import (
	"github.com/Zamiell/isaac-racing-server/src/models"
)

/*
	The structs here are used in more than one WebSocket file
*/

// Recieved in all commands
type IncomingWebsocketData struct {
	Room      string         `json:"room"`
	Message   string         `json:"message"`
	Name      string         `json:"name"`
	Ruleset   models.Ruleset `json:"ruleset"`
	ID        int            `json:"id"`
	Comment   string         `json:"comment"`
	ItemID    int            `json:"itemID"`
	FloorNum  int            `json:"floorNum"`
	StageType int            `json:"stageType"`
	RoomID    string         `json:"roomID"`
	IP        string         `json:"ip"`
	Enabled   bool           `json:"enabled"`
	Value     int            `json:"value"`
	Command   string         // Added by the server after demarshaling
	v         *SessionValues // Added by the server after demarshaling
}

/*
	Chat room data types
*/

type User struct {
	Name      string `json:"name"`
	Admin     int    `json:"admin"`
	Muted     bool   `json:"muted"`
	StreamURL string `json:"streamURL"`
}

/*
	Race data types
*/

// Sent in the "racerList" command (in the "raceCreate", "raceJoin", "raceJoinSpectate", and "handleConnect" functions)
type RacerListMessage struct {
	ID     int            `json:"id"`
	Racers []models.Racer `json:"racers"`
}

// Sent in the "raceLeft" command (in the "raceLeave", "handleDisconnect", and "adminBan" functions)
type RaceLeftMessage struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Sent in the "raceStart" command (in the "raceCheckStart" and "websocketHandleConnect" functions)
type RaceStartMessage struct {
	ID   int   `json:"id"`
	Time int64 `json:"time"`
}

/*
	TODO MOVE BELOW INTO THE COMMAND FUNCTIONS
*/

// Sent in the "roomSetMuted" command (in the "adminMute" and "adminUnmute" functions)
type RoomSetMutedMessage struct {
	Room  string `json:"room"`
	Name  string `json:"name"`
	Muted int    `json:"muted"`
}

// Sent in the "roomSetAdmin" command (in the "adminPromote" and "adminDemote" functions)
type RoomSetAdminMessage struct {
	Room  string `json:"room"`
	Name  string `json:"name"`
	Admin int    `json:"admin"`
}

// Sent as part of the "roomListAll" command (in the "roomListAll" function)
type Room struct {
	Room     string `json:"room"`
	NumUsers int    `json:"numUsers"`
}

/*
	Race data types
*/

// Sent in the "raceSetRuleset" command (in the "raceRuleset" function)
type RaceSetRulesetMessage struct {
	ID      int            `json:"id"`
	Ruleset models.Ruleset `json:"ruleset"`
}

// Sent in the "raceSetStatus" command (in the "raceCheckStart" functions)
type RaceSetStatusMessage struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// Sent in the "racerSetStatus" command (in the "raceReady", "raceUnready", "raceFinish", "raceQuit", and "adminBan" functions)
type RacerSetStatusMessage struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Place  int    `json:"place"`
}

// Sent in the "racerSetComment" command (in the "raceComment" functions)
type RacerSetCommentMessage struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

// Sent in the "racerAddItem" command (in the "raceItem" function)
type RacerAddItemMessage struct {
	ID   int         `json:"id"`
	Name string      `json:"name"`
	Item models.Item `json:"item"`
}

// Sent in the "racerSetFloor" command (in the "raceFloor" functions)
type RacerSetFloorMessage struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	FloorNum     int    `json:"floorNum"`
	StageType    int    `json:"stageType"`
	FloorArrived int    `json:"floorArrived"`
}
