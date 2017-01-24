package main

/*
	Imports
*/

import (
	"errors"
	"fmt"
	"github.com/Zamiell/isaac-racing-server/models"
	"github.com/getsentry/raven-go"
	"github.com/op/go-logging"
	"github.com/trevex/golem"
	"time"
)

/*
	Log struct extension
*/

type CustomLogger struct {
	Logger *logging.Logger
}

func (l *CustomLogger) Fatal(message string, err error) {
	raven.CaptureError(err, map[string]string{
		"message": message,
	})
	l.Logger.Fatal(message, err)
}

func (l *CustomLogger) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}

func (l *CustomLogger) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *CustomLogger) Warning(args ...interface{}) {
	l.Logger.Warning(args...)
	err := errors.New(fmt.Sprint(args...))
	raven.CaptureError(err, nil)
}

func (l *CustomLogger) Error(args ...interface{}) {
	l.Logger.Error(args...)
	err := errors.New(fmt.Sprint(args...))
	raven.CaptureError(err, nil)
}

/*
	Golem data types
*/

// We must extend the default Golem connection so that it hold information about the user
type ExtendedConnection struct {
	Connection         *golem.Connection
	UserID             int
	Username           string
	Admin              int
	Muted              int
	RateLimitAllowance float64
	RateLimitLastCheck time.Time
}

// Recieved in all commands
type IncomingCommandMessage struct {
	Room      string         `json:"room"`
	Message   string         `json:"message"`
	Name      string         `json:"name"`
	Ruleset   models.Ruleset `json:"ruleset"`
	ID        int            `json:"id"`
	Comment   string         `json:"comment"`
	ItemID    int            `json:"itemID"`
	FloorNum  int            `json:"floorNum"`
	StageType int            `json:"stageType"`
	IP        string         `json:"ip"`
	Enabled   bool           `json:"enabled"`
	Value     int            `json:"value"`
}

// Sent upon a successful WebSocket connection
type SettingsMessage struct {
	Username         string `json:"username"`
	StreamURL        string `json:"streamURL"`
	TwitchBotEnabled bool   `json:"twitchBotEnabled"`
	TwitchBotDelay   int    `json:"twitchBotDelay"`
	Time             int64  `json:"time"`
}

// Sent in the "error" and "warning" commands (in the "connError" and "connWarning" functions)
type ErrorMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

/*
	Chat room data types
*/

// Sent in the "roomList" command to the person that is joining the room (in the "roomJoinSub" function)
type RoomListMessage struct {
	Room  string `json:"room"`
	Users []User `json:"users"`
}
type User struct {
	Name  string `json:"name"`
	Admin int    `json:"admin"`
	Muted int    `json:"muted"`
}

// Sent in the "roomHistory" command to the person that is joining the room (in the "roomJoinSub" function)
type RoomHistoryMessage struct {
	Room    string               `json:"room"`
	History []models.RoomHistory `json:"history"`
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
	Room    string `json:"room"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

// Sent in the "privateMessage" command (in the "privateMessage" function)
type PrivateMessageMessage struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

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

// Sent in the "racerList" command (in the "connOpen" function)
type RacerList struct {
	ID     int            `json:"id"`
	Racers []models.Racer `json:"racers"`
}

// Sent in the "raceJoined" command (in the "raceJoin" function)
// Sent in the "raceLeft" command (in the "raceLeave" and "adminBan" functions)
type RaceMessage struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

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

// Sent to tell the client exactly when the race is starting
type RaceStartMessage struct {
	ID   int   `json:"id"`
	Time int64 `json:"time"`
}

// Sent to tell the client that they got a new achievement
type AchievementMessage struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

/*
	Profile data types
*/

// Sent in the "profile" command (in the "getProfile" function)
type Profile struct {
	// TODO
}

// Sent in the "profileSetUsername" command (in the "profileSetUsername" function)
type ProfileSetNameMessage struct {
	Name    string `json:"name"`
	NewName string `json:"newName"`
}

/*
	HTTP data types
*/

type LeaderboardSeeded []models.LeaderboardRowSeeded

/*func (l LeaderboardSeeded) Len() int {
	return len(s)
}
func (l LeaderboardSeeded) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (l LeaderboardSeeded) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}*/

type LeaderboardUnseeded []models.LeaderboardRowUnseeded
