package models

/*
	Imports
*/

import (
	"database/sql"
	"errors"
)

/*
	Data types
*/

type Users struct{}

/*
	"users" table functions
*/

func (*Users) Login(auth0ID string) (int, string, int, error) {
	// Check to see if they are in the user database already
	var userID int
	var username string
	var admin int
	err := db.QueryRow("SELECT id, username, admin FROM users WHERE auth0_id = ?", auth0ID).Scan(&userID, &username, &admin)
	if err == sql.ErrNoRows {
		return 0, "", 0, nil
	} else if err != nil {
		return 0, "", 0, err
	} else {
		return userID, username, admin, nil
	}
}

func (*Users) Exists(username string) (bool, error) {
	// Check if the user exists in the database
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (*Users) CheckStaff(username string) (bool, error) {
	// Check if the user is a staff member or an administrator
	var admin int
	err := db.QueryRow("SELECT admin FROM users WHERE username = ?", username).Scan(&admin)
	if err != nil {
		return false, err
	} else if admin == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (*Users) GetStream(userID int) (string, error) {
	// Get the user's stream URL
	var stream string
	err := db.QueryRow("SELECT stream FROM users WHERE id = ?", userID).Scan(&stream)
	if err != nil {
		return "", err
	} else {
		return stream, nil
	}
}

func (*Users) GetAllStreams() ([]string, error) {
	// Get a list of all the stream URLs
	rows, err := db.Query(`
		SELECT stream
		FROM users
		WHERE stream != '-'
		AND twitch_bot_enabled = 1
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the stream URLs
	var streams []string
	for rows.Next() {
		var stream string
		err := rows.Scan(&stream)
		if err != nil {
			return nil, err
		}

		// Append this stream to the slice
		streams = append(streams, stream)
	}

	return streams, nil
}

func (*Users) GetTwitchBotEnabled(username string) (bool, error) {
	// Get the user's Twitch bot setting
	var enabled int
	err := db.QueryRow(`
		SELECT twitch_bot_enabled
		FROM users
		WHERE username = ?
	`, username).Scan(&enabled)
	if err != nil {
		return false, err
	} else if enabled == 0 {
		return false, nil
	} else if enabled == 1 {
		return true, nil
	} else {
		return false, errors.New("The \"twitch_bot_enabled\" field for user \"" + username + "\" was not set to 0 or 1.")
	}
}

func (*Users) GetTwitchBotDelay(username string) (int, error) {
	// Get the user's Twitch bot setting
	var delay int
	err := db.QueryRow(`
		SELECT twitch_bot_delay
		FROM users
		WHERE username = ?
	`, username).Scan(&delay)
	if err != nil {
		return delay, err
	} else {
		return delay, nil
	}
}

func (*Users) SetStream(userID int, stream string) error {
	// Set the new stream URL
	stmt, err := db.Prepare("UPDATE users SET stream = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(stream, userID)
	if err != nil {
		return err
	}

	return nil
}

func (*Users) SetTwitchBotEnabled(userID int, enabled bool) error {
	// Set the new Twitch bot setting
	var twitchBotEnabled int
	if enabled == true {
		twitchBotEnabled = 1
	} else {
		twitchBotEnabled = 0
	}
	stmt, err := db.Prepare("UPDATE users SET twitch_bot_enabled = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(twitchBotEnabled, userID)
	if err != nil {
		return err
	}

	return nil
}

func (*Users) SetTwitchBotDelay(userID int, delay int) error {
	// Set the new Twitch bot setting
	stmt, err := db.Prepare("UPDATE users SET twitch_bot_delay = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(delay, userID)
	if err != nil {
		return err
	}

	return nil
}

func (*Users) SetLogin(username string, lastIP string) error {
	// Update the database with last_login and last_ip
	stmt, err := db.Prepare("UPDATE users SET last_login = ?, last_ip = ? WHERE username = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(makeTimestamp(), lastIP, username)
	if err != nil {
		return err
	}

	return nil
}

func (*Users) SetUsername(userID int, username string) error {
	// Set the new username
	stmt, err := db.Prepare("UPDATE users SET username = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(username, userID)
	if err != nil {
		return err
	}

	return nil
}

func (*Users) SetAdmin(username string, admin int) error {
	// Set the admin field for this user
	stmt, err := db.Prepare("UPDATE users SET admin = ? WHERE username = ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(admin, username)
	if err != nil {
		return err
	}

	return nil
}

func (*Users) Insert(auth0ID string, auth0Username string, ip string) (int, error) {
	// Get the current time
	currentTime := makeTimestamp()

	// Add them to the database
	stmt, err := db.Prepare("INSERT INTO users (auth0_id, username, last_ip, datetime_created, last_login) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(auth0ID, auth0Username, ip, currentTime, currentTime)
	if err != nil {
		return 0, err
	}
	userID64, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	userID := int(userID64)

	return userID, nil
}
