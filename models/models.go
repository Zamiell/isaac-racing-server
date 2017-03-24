package models

/*
	Imports
*/

import (
	"database/sql" // For connecting to the database (1/2)
	"time"

	_ "github.com/mattn/go-sqlite3" // For connecting to the database (2/2)
)

/*
	Data types
*/

type Models struct {
	// Database tables
	Achievements
	BannedIPs
	BannedUsers
	ChatLogPM
	ChatLog
	RaceParticipantItems
	RaceParticipantRooms
	RaceParticipants
	Races
	MutedUsers
	UserAchievements
	Users
}

// Sent in the "roomHistory" command (in the "roomJoinSub" function)
type RoomHistory struct {
	Name     string `json:"name"`
	Message  string `json:"message"`
	Datetime int64  `json:"datetime"`
}

// Sent in the "raceList" command (in the "connOpen" function)
// Sent in the "raceCreated" command (in the "raceCreate" function)
type Race struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Status          string   `json:"status"`
	Ruleset         Ruleset  `json:"ruleset"`
	Seed            string   `json:"seed"`
	DatetimeCreated int64    `json:"datetimeCreated"`
	DatetimeStarted int64    `json:"datetimeStarted"`
	Captain         string   `json:"captain"` // This is an integer in the database but we convert it to their name during the SELECT
	Racers          []string `json:"racers"`
}
type Ruleset struct {
	Type          string `json:"type"`
	Solo          bool   `json:"solo"`
	Format        string `json:"format"`
	Character     string `json:"character"`
	Goal          string `json:"goal"`
	StartingBuild int    `json:"startingBuild"`
}

// Sent in the "racerList" command (in the "connOpen" function)
// Used internally in the "raceStart" function
type Racer struct {
	Name             string `json:"name"`
	Status           string `json:"status"`
	DatetimeJoined   int    `json:"datetimeJoined"`
	DatetimeFinished int    `json:"datetimeFinished"`
	Place            int    `json:"place"`
	PlaceMid         int    `json:"placeMid"`
	Comment          string `json:"comment"`
	Items            []Item `json:"items"`
	FloorNum         int    `json:"floorNum"`
	StageType        int    `json:"stageType"`
	FloorArrived     int    `json:"floorArrived"`
	StreamURL        string `json:"streamURL"`        // Only used internally, but still sent to clients, TODO REMOVE
	TwitchBotEnabled int    `json:"twitchBotEnabled"` // Only used internally, but still sent to clients, TODO REMOVE
	TwitchBotDelay   int    `json:"twitchBotDelay"`   // Only used internally, but still sent to clients, TODO REMOVE
}
type Item struct {
	ID        int `json:"id"`
	FloorNum  int `json:"floorNum"`
	StageType int `json:"stageType"`
}

type StatsSeeded struct {
	ELO            int
	LastELOChange  int
	NumSeededRaces int
	LastSeededRace int
}

type StatsUnseeded struct {
	UnseededAdjustedAverage int
	UnseededRealAverage     int
	NumUnseededRaces        int
	NumForfeits             int
	ForfeitPenalty          int
	LowestUnseededTime      int
	LastUnseededRace        int
}

// Used in the leaderboards (HTTP)
type LeaderboardRowSeeded struct {
	Name           string
	ELO            int
	LastELOChange  int
	NumSeededRaces int
	LastSeededRace int
	Verified       int
}
type LeaderboardRowUnseeded struct {
	Name                    string
	UnseededAdjustedAverage int
	UnseededRealAverage     int
	NumUnseededRaces        int
	NumForfeits             int
	ForfeitPenalty          int
	LowestUnseededTime      int
	LastUnseededRace        int
	Verified                int
}
type UserProfilesRow struct {
	Username		string
	DateCreated		int
	StreamUrl		string
	Achievements	int
}
type UserProfileData struct {
	Username		string
	DateCreated		int
	Verified		int
	ELO				int
	LastELOChange	int
	SeededRaces		int
	UnseededRaces	int
	StreamUrl		string
}
type LeaderboardRowTopTimes struct {
	Name     string
	Time     int
	Date     int
	Verified int
}
type LeaderboardRowMostPlayed struct {
	Name     string
	Total    int
	Verified int
}

/*
	Global variables
*/

var (
	db *sql.DB
)

/*
	Initialization function
*/

func GetModels(dbFile string) (*Models, error) {
	// Initialize the database
	var err error
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	// Enable foreign key constraints (which are disabled by default in SQLite3)
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}

	// Create the model
	return &Models{}, nil
}

/*
	Miscellaneous functions
*/

// From: https://stackoverflow.com/questions/24122821/go-golang-time-now-unixnano-convert-to-milliseconds
func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
