package server

import (
	"github.com/Zamiell/isaac-racing-server/models"
)

/*
	The structs here are used in more than one WebSocket file
*/

// Received in all commands
type IncomingWebsocketData struct {
	Room          string                `json:"room"`
	Message       string                `json:"message"`
	Name          string                `json:"name"`
	Ruleset       Ruleset               `json:"ruleset"`
	Password      string                `json:"password"`
	ID            int                   `json:"id"`
	Comment       string                `json:"comment"`
	Seed          string                `json:"seed"`
	ItemID        int                   `json:"itemID"` // nolint:tagliatelle
	FloorNum      int                   `json:"floorNum"`
	StageType     int                   `json:"stageType"`
	BackwardsPath bool                  `json:"backwardsPath"`
	RoomID        string                `json:"roomID"` // nolint:tagliatelle
	IP            string                `json:"ip"`
	Enabled       bool                  `json:"enabled"`
	Value         int                   `json:"value"`
	Time          int64                 `json:"time"`
	Command       string                // Added by the server after demarshaling
	v             *models.SessionValues // Added by the server after demarshaling
}

/*
	Chat room data types
*/

// Used in defining the "chatRooms" map
// Sent in the "roomUpdate" command (in the "chatRoomsUpdate()" function)
type User struct {
	Name      string `json:"name"`
	Admin     int    `json:"admin"`
	Muted     bool   `json:"muted"`
	StreamURL string `json:"streamURL"` // nolint:tagliatelle
}

// Sent in the "roomMessage" command (in the "websocketroomMessage()" function)
// Sent in the "adminBan" command (in the "websocketAdminBan()" function)
// Sent in the "adminUnban" command (in the "websocketAdminUnban()" function)
type RoomMessageMessage struct {
	Room    string `json:"room"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

/*
	Race data types
*/

// Sent in the "raceCreate" command (in the "websocketRaceCreate()" function)
// Sent in the "raceList" command (in the "websocketHandleConnect()" function)
type RaceCreatedMessage struct {
	ID                  int        `json:"id"`
	Name                string     `json:"name"`
	Status              RaceStatus `json:"status"`
	Ruleset             Ruleset    `json:"ruleset"`
	Captain             string     `json:"captain"`
	IsPasswordProtected bool       `json:"isPasswordProtected"`
	DatetimeCreated     int64      `json:"datetimeCreated"`
	DatetimeStarted     int64      `json:"datetimeStarted"`
	Racers              []string   `json:"racers"`
}

// Sent in the "raceStart" command (in the "raceCheckStart()" and "websocketHandleConnect()" functions)
type RaceStartMessage struct {
	ID            int `json:"id"`
	SecondsToWait int `json:"secondsToWait"`
}

/*
	Admin data types
*/

// Sent in the "adminMessage" command (in the "websocketHandleConnect()" and "websocketAdminMessage()" functions)
type AdminMessageMessage struct {
	Message string `json:"message"`
}
