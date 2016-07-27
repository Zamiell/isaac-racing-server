package model

/*
 *  Imports
 */

import (
	"database/sql"                  // For connecting to the database (1/2)
	_ "github.com/mattn/go-sqlite3" // For connecting to the database (2/2)
	"github.com/op/go-logging"
)

/*
 *  Data types
 */

type Model struct {
	// Database tables
	Users
	Races
	RaceParticipants
	RaceParticipantItems
	BannedUsers
	BannedIPs
	SquelchedUsers
	ChatLog
	ChatLogPM
	Achievements
	UserAchievements
}

// Used internally in the "Users.Login" function
type LoginInformation struct {
	UserID    int
	Admin     int
	Squelched int
}

// Sent in the "roomHistoryList" command (in the "roomJoinSub" function)
type ChatHistoryMessage struct {
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
	Ruleset         string   `json:"ruleset"`
	DatetimeCreated int      `json:"datetime_created"`
	DatetimeStarted int      `json:"datetime_started"`
	Captain         string   `json:"captain"` // This is an integer in the database but we convert it to their name during the SELECT
	Racers          []string `json:"racers"`
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
	db  *sql.DB
	log *logging.Logger
)

/*
 *  Initialization function
 */

func GetModel(dbFile string, logger *logging.Logger) *Model {
	// Initialize the database
	var err error
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		logger.Fatal("Failed to open database:", err)
	}

	// Initialize the logger
	log = logger

	// Create the model and fill it with helpful self-references
	model := &Model{}
	model.Users.db = model
	model.Races.db = model
	model.RaceParticipants.db = model
	model.RaceParticipantItems.db = model
	model.BannedUsers.db = model
	model.BannedIPs.db = model
	model.SquelchedUsers.db = model
	model.ChatLog.db = model
	model.ChatLogPM.db = model
	model.Achievements.db = model
	model.UserAchievements.db = model

	return model
}
