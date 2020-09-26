package models

import (
	"database/sql"
	"fmt"
	"os"

	//
	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

/*
	Data types
*/

// Each property represents a database table
type Models struct {
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
	Tournament
	UserAchievements
	Users
}

func (*Models) Close() {
	db.Close()
}

/*
	Initialization function
*/

func Init() (*Models, error) {
	// Read the database configuration from environment variables
	// (it was loaded from the .env file in main.go)
	dbHost := os.Getenv("DB_HOST")
	if len(dbHost) == 0 {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if len(dbPort) == 0 {
		dbPort = "3306" // This is the default port for MySQL
	}
	dbUser := os.Getenv("DB_USER")
	if len(dbUser) == 0 {
		defaultUser := "isaacuser"
		fmt.Println("DB_USER not specified; using default value of \"" + defaultUser + "\".")
		dbUser = defaultUser
	}
	dbPass := os.Getenv("DB_PASS")
	if len(dbPass) == 0 {
		defaultPass := "1234567890"
		fmt.Println("DB_PASS not specified; using default value of \"" + defaultPass + "\".")
		dbPass = defaultPass
	}
	dbName := os.Getenv("DB_NAME")
	if len(dbPass) == 0 {
		defaultName := "isaac"
		fmt.Println("DB_NAME not specified; using default value of \"" + defaultName + "\".")
		dbName = defaultName
	}

	// Initialize the database
	dsn := dbUser + ":" + dbPass + "@(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true"
	if v, err := sql.Open("mysql", dsn); err != nil {
		return nil, err
	} else {
		db = v
	}

	// Create the model
	return &Models{}, nil
}
