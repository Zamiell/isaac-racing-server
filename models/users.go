package models

/*
	Imports
*/

import (
	"database/sql"
	"errors"
	"strconv"
	
)

/*
	Data types
*/

type Users struct{}

/*
	"users" table functions
*/

func (*Users) Login(steamID string) (int, string, int, error) {
	// Check to see if they are in the user database already
	var userID int
	var username string
	var admin int
	err := db.QueryRow(`
		SELECT id, username, admin
		FROM users
		WHERE steam_id = ?
	`, steamID).Scan(&userID, &username, &admin)
	if err == sql.ErrNoRows {
		return 0, "", 0, nil
	} else if err != nil {
		return 0, "", 0, err
	} else {
		return userID, username, admin, nil
	}
}

func (*Users) Exists(username string) (bool, error) {
	// Check if the user exists in the database (and do a case-insensitive search)
	var id int
	err := db.QueryRow(`
		SELECT id
		FROM users
		WHERE username = ?
		COLLATE NOCASE
	`, username).Scan(&id)
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
	err := db.QueryRow(`
		SELECT admin
		FROM users
		WHERE username = ?
	`, username).Scan(&admin)
	if err != nil {
		return false, err
	} else if admin == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (*Users) CheckStreamURLUsed(streamURL string) (bool, error) {
	// Get a list of all the stream URLs
	rows, err := db.Query(`
		SELECT stream_url
		FROM users
		WHERE stream_url != '-'
	`)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// Iterate over the stream URLs
	var stream_urls []string
	for rows.Next() {
		var stream_url string
		err := rows.Scan(&stream_url)
		if err != nil {
			return false, err
		}

		// Append this stream to the slice
		stream_urls = append(stream_urls, stream_url)
	}

	// Find out if it exists
	// TODO

	return false, nil
}

func (*Users) GetStreamURL(userID int) (string, error) {
	// Get the user's stream URL
	var streamURL string
	err := db.QueryRow(`
		SELECT stream_url
		FROM users
		WHERE id = ?
	`, userID).Scan(&streamURL)
	if err != nil {
		return "", err
	} else {
		return streamURL, nil
	}
}

func (*Users) GetAllStreamURLs() ([]string, error) {
	// Get a list of all the stream URLs
	rows, err := db.Query(`
		SELECT stream_url
		FROM users
		WHERE stream_url != '-'
		AND twitch_bot_enabled = 1
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the stream URLs
	var stream_urls []string
	for rows.Next() {
		var stream_url string
		err := rows.Scan(&stream_url)
		if err != nil {
			return nil, err
		}

		// Append this stream to the slice
		stream_urls = append(stream_urls, stream_url)
	}

	return stream_urls, nil
}

func (*Users) GetTwitchBotEnabled(userID int) (bool, error) {
	// Get the user's Twitch bot setting
	var enabled int
	err := db.QueryRow(`
		SELECT twitch_bot_enabled
		FROM users
		WHERE id = ?
	`, userID).Scan(&enabled)
	if err != nil {
		return false, err
	} else if enabled == 0 {
		return false, nil
	} else if enabled == 1 {
		return true, nil
	} else {
		return false, errors.New("The \"twitch_bot_enabled\" field for user ID \"" + strconv.Itoa(userID) + "\" was not set to 0 or 1.")
	}
}

func (*Users) GetTwitchBotDelay(userID int) (int, error) {
	// Get the user's Twitch bot setting
	var delay int
	err := db.QueryRow(`
		SELECT twitch_bot_delay
		FROM users
		WHERE id = ?
	`, userID).Scan(&delay)
	if err != nil {
		return delay, err
	} else {
		return delay, nil
	}
}

func (*Users) GetLeaderboardSeeded() ([]LeaderboardRowSeeded, error) {
	// Make a leaderboard for the seeded format based on all of the users
	rows, err := db.Query(`
		SELECT
			username,
			elo,
			last_elo_change,
			num_seeded_races,
			last_seeded_race
		FROM users
	`)
	// WHERE num_seeded_races > 1
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the users
	leaderboard := make([]LeaderboardRowSeeded, 0)
	for rows.Next() {
		var row LeaderboardRowSeeded
		err := rows.Scan(
			&row.Name,
			&row.ELO,
			&row.LastELOChange,
			&row.NumSeededRaces,
			&row.LastSeededRace,
		)
		if err != nil {
			return leaderboard, err
		}

		// Append this row to the leaderboard
		leaderboard = append(leaderboard, row)
	}

	return leaderboard, nil
}

func (*Users) GetLeaderboardUnseeded() ([]LeaderboardRowUnseeded, error) {
	// Make a leaderboard for the unseeded format based on all of the users
	rows, err := db.Query(`
		SELECT
			username,
			unseeded_adjusted_average,
			unseeded_real_average,
			num_unseeded_races,
			num_forfeits,
			forfeit_penalty,
			lowest_unseeded_time,
			last_unseeded_race
		FROM users
	`)
	// WHERE num_unseeded_races > 15
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the users
	leaderboard := make([]LeaderboardRowUnseeded, 0)
	for rows.Next() {
		var row LeaderboardRowUnseeded
		err := rows.Scan(
			&row.Name,
			&row.UnseededAdjustedAverage,
			&row.UnseededRealAverage,
			&row.NumUnseededRaces,
			&row.NumForfeits,
			&row.ForfeitPenalty,
			&row.LowestUnseededTime,
			&row.LastUnseededRace,
		)
		if err != nil {
			return leaderboard, err
		}

		// Append this row to the leaderboard
		leaderboard = append(leaderboard, row)
	}

	return leaderboard, nil
}

func (*Users) GetProfileData(player string) (UserProfileData, error) {
	// Get player data to populate player profile page.
	var profileData UserProfileData
	err := db.QueryRow(`
		SELECT
			username,
			datetime_created,
			verified,
			elo,
			last_elo_change,
			num_seeded_races,
			num_unseeded_races,
			stream_url
		FROM
			users
		WHERE
			steam_id > 0 and
			username = ?
	`, player).Scan(
			&profileData.Username,
			&profileData.DateCreated,
			&profileData.Verified,
			&profileData.ELO,
			&profileData.LastELOChange,
			&profileData.SeededRaces,
			&profileData.UnseededRaces,
			&profileData.StreamUrl,
		)
	if err != nil {
		return profileData, err
	} else {
		return profileData, nil
	}
}
func (*Users) GetUserProfiles(currentPage int, usersPerPage int) ([]UserProfilesRow, int, error) {
	// Get players data to populate the profiles page.
	usersOffset := (currentPage - 1) * usersPerPage
	rows, err := db.Query(`
		SELECT 
			u.username, 
			u.datetime_created, 
			u.stream_url,
			count(ua.achievement_id) 
		FROM 
			users u 
		LEFT JOIN 
			user_achievements ua 
			ON
				u.id = ua.user_id 
		WHERE 
			u.steam_id > 0
		GROUP BY
			u.username
		ORDER BY
			u.username ASC
		LIMIT
			?
		OFFSET
			? 
	`, strconv.Itoa(usersPerPage), strconv.Itoa(usersOffset))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	// Iterate over the user profile results
	profiles := make([]UserProfilesRow, 0)
	for rows.Next() {
    	var row UserProfilesRow
		err := rows.Scan(
			&row.Username,
			&row.DateCreated,
			&row.StreamUrl,
			&row.Achievements,
		)
		if err != nil {
			return profiles, 0, err
		}
		// Append this row to the leaderboard
		profiles = append(profiles, row)
	}
	// Find total amount of users
	rows, err = db.Query(`
		SELECT 
			count(id) 
		FROM 
			users
		WHERE
			steam_id > 0
	`)
	if err != nil {
		return profiles, 0, err
	}
	defer rows.Close()
	var allProfilesCount int
	for rows.Next() {
		err = rows.Scan(&allProfilesCount)
		if err != nil {
			return profiles, allProfilesCount, err
		}
	}
	return profiles, allProfilesCount, nil
}
func (*Users) SetStreamURL(userID int, streamURL string) error {
	// Set the new stream URL
	stmt, err := db.Prepare(`
		UPDATE users
		SET stream_url = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(streamURL, userID)
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
	stmt, err := db.Prepare(`
		UPDATE users
		SET twitch_bot_enabled = ?
		WHERE id = ?
	`)
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
	stmt, err := db.Prepare(`
		UPDATE users
		SET twitch_bot_delay = ?
		WHERE id = ?
	`)
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
	// Update the database with datetime_last_login and last_ip
	stmt, err := db.Prepare(`
		UPDATE users
		SET datetime_last_login = ?, last_ip = ?
		WHERE username = ?
	`)
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
	stmt, err := db.Prepare(`
		UPDATE users
		SET username = ?
		WHERE id = ?
	`)
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
	stmt, err := db.Prepare(`
		UPDATE users
		SET admin = ?
		WHERE username = ?
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(admin, username)
	if err != nil {
		return err
	}

	return nil
}

func (*Users) Insert(steamID string, username string, ip string) (int, error) {
	// Get the current time
	currentTime := makeTimestamp()

	// Add them to the database
	stmt, err := db.Prepare(`
		INSERT INTO users (steam_id, username, last_ip, datetime_created, datetime_last_login)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(steamID, username, ip, currentTime, currentTime)
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
