package main

/*
 *  Imports
 */

import (
	"github.com/Zamiell/isaac-racing-server/models"
	"github.com/trevex/golem"
	"time"
)

/*
 *  Golem data types
 */

// We must extend the default Golem connection so that it hold information about the user
type ExtendedConnection struct {
	Connection         *golem.Connection
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

// Received in "roomJoom" and "roomLeave"
type RoomMessage struct {
	Name string `json:"name"`
}

// Received in "roomMessage" and "privateMessage"
type ChatMessage struct {
	To   string `json:"to"`
	From string `json:"from"`
	Msg  string `json:"msg"`
}

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

// Sent in "roomGetAll"
type Room struct {
	Room     string `json:"room"`
	NumUsers int    `json:"numUsers"`
}

/*
 *  Race data types
 */

// Received in "raceJoin", "raceLeave", "raceReady", "raceUnready", "raceDone", and "raceQuit"
type RaceMessage struct {
	ID   int    `json:"id"`
	Name string `json:"name"` // Only used when returning information back to the client
}

// Received in "raceCreate"
type RaceCreateMessage struct {
	Name    string `json:"name"`
	Ruleset string `json:"ruleset"`
	ID      int    `json:"id"` // Only used when returning information back to the client
}

// Received in "raceRuleset"
type RaceRulesetMessage struct {
	ID      int    `json:"id"`
	Ruleset string `json:"ruleset"`
}

// Received in "raceComment"
type RaceCommentMessage struct {
	ID      int    `json:"id"`
	Comment string `json:"comment"`
}

// Received in "raceItem"
type RaceItemMessage struct {
	ID     int `json:"id"`
	ItemID int `json:"itemID"`
}

// Received in "raceFloor"
type RaceFloorMessage struct {
	ID    int `json:"id"`
	Floor int `json:"floor"`
}

// Sent to tell the client that something has happened within the particular race
type RacerList struct {
	ID     int           `json:"id"`
	Racers []model.Racer `json:"racers"`
}

// Sent to tell the client exactly when the race is starting
type RaceStartMessage struct {
	ID   int   `json:"id"`
	Time int64 `json:"time"`
}

/*
 *  Profile data types
 */

// Received in "profileGet" and "profileSetUsername"
type ProfileMessage struct {
	Name string `json:"name"`
}

// Sent after a "profileGet"
type Profile struct {
	// TODO
}

/*
 *  Admin data types
 */

// Received in "adminBan", "adminUnban", "adminSquelch", "adminSquelch", "adminUnsquelch", "adminPromote", and "adminDemote"
type AdminMessage struct {
	Name string `json:"name"`
}

// Received in "adminBanIP" and "adminUnbanIP"
type AdminIPMessage struct {
	IP string `json:"ip"`
}
