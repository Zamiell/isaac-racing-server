package main

/*
 *  Imports
 */

import (
	"github.com/Zamiell/isaac-racing-server/model"

	"time"
	"github.com/trevex/golem"
)

/*
 *  Golem data types
 */

// We must extend the default Golem connection so that it hold information about the user
type ExtendedConnection struct {
	Connection *golem.Connection
	UserID             int
	Username           string
	Admin              int
	Squelched          int
	RateLimitAllowance float64
	RateLimitLastCheck time.Time
}

// Used in "connSuccess" and "connError"
type SystemMessage struct {
	Type string      `json:"type"`
	Msg  interface{} `json:"msg"` 
}

/*
 *  Chat room data types
 */

// Sent on connection and whenever someone joins or leaves a chat channel
type RoomList struct {
	Room  string `json:"room"`
	Users []User `json:"users"`
}
type User struct {
	Name      string `json:"name"`
	Admin     int    `json:"admin"`
	Squelched int    `json:"squelched"`
}

// Used in "roomJoom" and "roomLeave"
type RoomMessage struct {
	Name string `json:"name"`
}

// Used in "roomMessage" and "privateMessage"
type ChatMessage struct {
	To   string `json:"to"`
	From string `json:"from"`
	Msg  string `json:"msg"`
}

/*
 *  Race data types
 */

// Used in "raceJoin", "raceLeave", "raceReady", "raceUnready", "raceDone", and "raceQuit"
type RaceMessage struct {
	ID      int    `json:"id"`
	Name    string `json:"name"` // Only used when returning information back to the client
}

// Used in "raceCreate"
type RaceCreateMessage struct {
	Name    string `json:"name"`
	Ruleset string `json:"ruleset"`
	ID      int    `json:"id"`      // Only used when returning information back to the client
}

// Used in "raceRuleset"
type RaceRulesetMessage struct {
	ID      int    `json:"id"`
	Ruleset string `json:"ruleset"`
}

// Used in "raceComment"
type RaceCommentMessage struct {
	ID      int    `json:"id"`
	Comment string `json:"comment"`
}

// Used in "raceItem"
type RaceItemMessage struct {
	ID     int `json:"id"`
	ItemID int `json:"itemID"`
}

// Used in "raceFloor"
type RaceFloorMessage struct {
	ID    int `json:"id"`
	Floor int `json:"floor"`
}

// Sent to tell the client exactly when the race is starting
type RaceStartMessage struct {
	ID   int   `json:"id"`
	Time int64 `json:"time"`
}

// Sent to tell the client that something has happened within the particular race
type RacerList struct {
	ID        int           `json:"id"`
	RacerList []model.Racer `json:"racerList"`
}

/*
 *  Admin data types
 */

// Used in "adminBan", "adminUnban", "adminSquelch", "adminSquelch", "adminUnsquelch", "adminPromote", and "adminDemote"
type AdminMessage struct {
	Name string `json:"name"`
}

// Used in "adminBanIP" and "adminUnbanIP"
type AdminIPMessage struct {
	IP string `json:"ip"`
}
