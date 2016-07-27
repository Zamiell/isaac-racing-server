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

// Recieved in all commands
type IncomingCommandMessage struct {
	Room    string        `json:"room"`
	Msg     string        `json:"msg"`
	Name    string        `json:"name"`
	Ruleset model.Ruleset `json:"ruleset"`
	ID      int           `json:"id"`
	Comment string        `json:"comment"`
	ItemID  int           `json:"itemID"`
	Floor   int           `json:"floor"`
	IP      string        `json:"ip"`
}

// Sent in an "success" command (in the "connSuccess" function)
type SuccessMessage struct {
	Type  string      `json:"type"`
	Input interface{} `json:"input"`
}

// Sent in an "error" command (in the "connError" function)
type ErrorMessage struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

/*
 *  Chat room data types
 */

// Sent in the "roomList" command to the person that is joining the room (in the "roomJoinSub" function)
type RoomListMessage struct {
	Room  string `json:"room"`
	Users []User `json:"users"`
}
type User struct {
	Name      string `json:"name"`
	Admin     int    `json:"admin"`
	Squelched int    `json:"squelched"`
}

// Sent in the "roomHistory" command to the person that is joining the room (in the "roomJoinSub" function)
type RoomHistoryMessage struct {
	Room    string              `json:"room"`
	History []model.RoomHistory `json:"history"`
}

// Sent in the "roomJoined" command to everyone who is already in the room (in the "roomJoinSub" function)
type RoomJoinedMessage struct {
	Room string `json:"room"`
	User User   `json:"user"`
}

// Sent in the "roomLeft" command (in the "roomLeaveSub" function)
type RoomLeftMessage struct {
	Room string `json:"room"`
	Name string `json:"name"`
}

// Sent in the "roomMessage" command (in the "roomMessage" function)
type RoomMessageMessage struct {
	Room string `json:"room"`
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

// Sent in the "privateMessage" command (in the "privateMessage" function)
type PrivateMessageMessage struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

// Sent in the "roomSetName" command (in the "profileSetUsername" function)
type RoomSetNameMessage struct {
	Room    string `json:"room"`
	Name    string `json:"name"`
	NewName string `json:"newName"`
}

// Sent in the "roomSetSquelched" command (in the "adminSquelch" and "adminUnsquelch" functions)
type RoomSetSquelchedMessage struct {
	Room      string `json:"room"`
	Name      string `json:"name"`
	Squelched int    `json:"squelched"`
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
 *  Race data types
 */

// Sent in the "racerList" command (in the "connOpen" function)
type RacerList struct {
	ID     int           `json:"id"`
	Racers []model.Racer `json:"racers"`
}

// Sent in the "raceJoined" command (in the "raceJoin" function)
// Sent in the "raceLeft" command (in the "raceLeave" and "adminBan" functions)
type RaceMessage struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Sent in the "raceSetRuleset" command (in the "raceRuleset" function)
type RaceSetRulesetMessage struct {
	ID      int           `json:"id"`
	Ruleset model.Ruleset `json:"ruleset"`
}

// Sent in the "raceSetStatus" command (in the "raceCheckStart" functions)
type RaceSetStatusMessage struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// Sent in the "racerSetStatus" command (in the "raceReady", "raceUnready", "raceFinish", and "raceQuit" functions)
type RacerSetStatusMessage struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// Sent in the "racerSetComment" command (in the "raceComment" functions)
type RacerSetCommentMessage struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

// Sent in the "racerAddItem" command (in the "raceItem" function)
type RacerAddItemMessage struct {
	ID   int        `json:"id"`
	Name string     `json:"name"`
	Item model.Item `json:"item"`
}

// Sent in the "racerSetFloor" command (in the "raceFloor" functions)
type RacerSetFloorMessage struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Floor int    `json:"floor"`
}

// Sent to tell the client exactly when the race is starting
type RaceStartMessage struct {
	ID   int   `json:"id"`
	Time int64 `json:"time"`
}

/*
 *  Profile data types
 */

// Sent in the "profile" command (in the "getProfile" function)
type Profile struct {
	// TODO
}
