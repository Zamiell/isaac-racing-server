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
	//Achievements
	//UserAchievements
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

func GetModel(logger *logging.Logger) *Model {
	// Initialize the database
	var err error
	db, err = sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		logger.Fatal("Failed to open database:", err)
	}

	// Initialize the logger
	log = logger

	// Create the model and fill it with helpful self-referneces
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
	//model.Achievements.db = model
	//model.UserAchievements.db = model

	return model
}
