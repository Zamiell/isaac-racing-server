package models

/*
 *  Imports
 */

import (
	"database/sql"                  // For connecting to the database (1/2)
	_ "github.com/mattn/go-sqlite3" // For connecting to the database (2/2)
)

/*
 *  Data types
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
	Seeds
	SquelchedUsers
	UserAchievements
	Users
}

// Sent in the "roomHistory" command (in the "roomJoinSub" function)
type RoomHistory struct {
	Name     string `json:"name"`
	Msg      string `json:"msg"`
	Datetime int    `json:"datetime"`
}

// Sent in the "raceList" command (in the "connOpen" function)
// Sent in the "raceCreated" command (in the "raceCreate" function)
type Race struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Status          string   `json:"status"`
	Ruleset         Ruleset  `json:"ruleset"`
	DatetimeCreated int      `json:"datetime_created"`
	DatetimeStarted int      `json:"datetime_started"`
	Captain         string   `json:"captain"` // This is an integer in the database but we convert it to their name during the SELECT
	Racers          []string `json:"racers"`
}
type Ruleset struct {
	Type         string `json:"type"`
	Character    string `json:"character"`
	Goal         string `json:"goal"`
	Seed         string `json:"seed"`
	InstantStart int    `json:"instantStart"`
}

// Sent in the "racerList" command (in the "connOpen" function)
// Used internally in the "raceStart" function
type Racer struct {
	Name             string `json:"name"`
	Status           string `json:"status"`
	DatetimeJoined   int    `json:"datetime_joined"`
	DatetimeFinished int    `json:"datetime_finished"`
	Place            int    `json:"place"`
	Comment          string `json:"comment"`
	Items            []Item `json:"items"`
	Floor            int    `json:"floor"`
}
type Item struct {
	ID    int `json:"id"`
	Floor int `json:"floor"`
}

/*
 *  Global variables
 */

var (
	db *sql.DB
)

/*
 *  Initialization function
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
