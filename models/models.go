package models

/*
	Imports
*/

import (
	"database/sql"                  // For connecting to the database (1/2)
	_ "github.com/mattn/go-sqlite3" // For connecting to the database (2/2)
	"time"
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
	RaceParticipants
	Races
	SquelchedUsers
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
	Comment          string `json:"comment"`
	Items            []Item `json:"items"`
	Floor            string `json:"floor"`
}
type Item struct {
	ID    int    `json:"id"`
	Floor string `json:"floor"`
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
