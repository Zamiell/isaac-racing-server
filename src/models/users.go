package models

import (
	"database/sql"
	"errors"
	"strconv"
	"time"
)

type Users struct{}

// Used during login to get all the information about a user
// Used when passing the cookie values from "httpValidateSession" to "httpWS"
// Used in "websocketHandleMessage" to stuff WebSocket session values into
// the IncomingWebsocketData object as a convenience for command handler
// functions
type SessionValues struct {
	UserID             int
	Username           string
	Admin              int
	Muted              bool
	StreamURL          string
	TwitchBotEnabled   bool
	TwitchBotDelay     int
	Banned             bool
	RateLimitAllowance float64   // Set in the "websocketGetSessionValues()" function
	RateLimitLastCheck time.Time // Set in the "websocketGetSessionValues()" function
}

// Used in the "httpRegister" function
func (*Users) Insert(steamID string, username string, ip string) (int, error) {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO users (steam_id, username, last_ip)
		VALUES (?, ?, ?)
	`); err != nil {
		return 0, err
	} else {
		stmt = v
	}
	defer stmt.Close()

	var result sql.Result
	if v, err := stmt.Exec(steamID, username, ip); err != nil {
		return 0, err
	} else {
		result = v
	}

	var userID int
	if userID64, err := result.LastInsertId(); err != nil {
		return 0, err
	} else {
		userID = int(userID64)
	}

	return userID, nil
}

// Essentially, this checks to see if they are in the user database already
// If they are, it returns a bunch of information about them
// Used in the "httpLogin" and "httpRegister" functions
func (*Users) Login(steamID string) (*SessionValues, error) {
	sessionValues := &SessionValues{}
	var rawMuted int
	var rawTwitchBotEnabled int
	var rawBanned int
	if err := db.QueryRow(`
		SELECT
			id AS matched_id,
			username,
			admin,
			(
				SELECT COUNT(id)
				FROM muted_users
				WHERE user_id = matched_id
			) AS muted,
			stream_url,
			twitch_bot_enabled,
			twitch_bot_delay,
			(
				SELECT COUNT(id)
				FROM banned_users
				WHERE user_id = matched_id
			) AS banned
		FROM
			users
		WHERE
			steam_id = ?
	`, steamID).Scan(
		&sessionValues.UserID,
		&sessionValues.Username,
		&sessionValues.Admin,
		&rawMuted,
		&sessionValues.StreamURL,
		&rawTwitchBotEnabled,
		&sessionValues.TwitchBotDelay,
		&rawBanned,
	); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Convert the ints to bools
	// (MariaDB stores booleans as a 0 or 1)
	sessionValues.Muted = false
	if rawMuted == 1 {
		sessionValues.Muted = true
	}
	sessionValues.TwitchBotEnabled = false
	if rawTwitchBotEnabled == 1 {
		sessionValues.TwitchBotEnabled = true
	}
	sessionValues.Banned = false
	if rawBanned == 1 {
		sessionValues.Banned = true
	}

	return sessionValues, nil
}

// Also returns the user ID so that it can be used in the "websocketAdminBan()" function
func (*Users) Exists(username string) (bool, int, error) {
	var id int
	if err := db.QueryRow(`
		SELECT id
		FROM users
		WHERE username = ?
	`, username).Scan(&id); err == sql.ErrNoRows {
		return false, -1, nil
	} else if err != nil {
		return false, -1, err
	}

	return true, id, nil
}

// Check to see if a user is either staff or an admin
// Used in the "websocketAdminBan()" function
func (*Users) GetAdmin(userID int) (int, error) {
	var admin int
	if err := db.QueryRow(`
		SELECT admin
		FROM users
		WHERE id = ?
	`, userID).Scan(&admin); err == sql.ErrNoRows {
		return -1, errors.New("A user with an ID of \"" + strconv.Itoa(userID) + "\" does not exist.")
	} else if err != nil {
		return -1, err
	}

	return admin, nil
}

// Used in the "twitchConnect" and "websocketProfileSetStream" functions
func (*Users) GetAllStreamURLs() ([]string, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT stream_url
		FROM users
		WHERE stream_url != '-'
		AND twitch_bot_enabled = 1
	`); err != nil {
		return nil, err
	} else {
		rows = v
	}
	defer rows.Close()

	var streamURLs []string
	for rows.Next() {
		var streamURL string
		if err := rows.Scan(&streamURL); err != nil {
			return nil, err
		}

		streamURLs = append(streamURLs, streamURL)
	}

	return streamURLs, nil
}

// Get the user ID and username that matches a stream URL
// (MariaDB will perform a case insensitive comparison by default,
// which is what we want)
// Used in the "twitchNotMod" function
func (*Users) GetUserFromStreamURL(streamURL string) (int, string, error) {
	var userID int
	var username string
	if err := db.QueryRow(`
		SELECT id, username
		FROM users
		WHERE stream_url = ?
	`, streamURL).Scan(&userID, &username); err == sql.ErrNoRows {
		return 0, "", nil
	} else if err != nil {
		return 0, "", err
	}

	return userID, username, nil
}

// Update the database with datetime_last_login and last_ip
// Used in the "httpLogin" function
func (*Users) SetLogin(userID int, lastIP string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE users
		SET datetime_last_login = NOW(), last_ip = ?
		WHERE id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(lastIP, userID); err != nil {
		return err
	}

	return nil
}

// Used in the "websocketProfileSetStream" function
func (*Users) SetStreamURL(userID int, streamURL string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE users
		SET stream_url = ?
		WHERE id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(streamURL, userID); err != nil {
		return err
	}

	return nil
}

// Used in the "websocketProfileSetStream" function
func (*Users) SetTwitchBotEnabled(userID int, enabled bool) error {
	var twitchBotEnabled int
	if enabled {
		twitchBotEnabled = 1
	} else {
		twitchBotEnabled = 0
	}

	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE users
		SET twitch_bot_enabled = ?
		WHERE id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(twitchBotEnabled, userID); err != nil {
		return err
	}

	return nil
}

// Used in the "websocketProfileSetStream" function
func (*Users) SetTwitchBotDelay(userID int, delay int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE users
		SET twitch_bot_delay = ?
		WHERE id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(delay, userID); err != nil {
		return err
	}

	return nil
}

/*
// Used in ?
func (*Users) SetAdmin(username string, admin int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE users
		SET admin = ?
		WHERE username = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(admin, username); err != nil {
		return err
	}

	return nil
}
*/
