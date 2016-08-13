package models

/*
 *  Imports
 */

import (
	"database/sql"
)

/*
 *  Data types
 */

type Users struct{}

/*
 *  users table functions
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

func (*Users) SetLogin(username string, lastIP string) error {
	// Update the database with last_login and last_ip
	stmt, err := db.Prepare("UPDATE users SET last_login = (strftime('%s', 'now')), last_ip = ? WHERE username = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(lastIP, username)
	if err != nil {
		return err
	}

	return nil
}

func (*Users) SetUsername(userID int, username string) error {
	// Set the new username
	stmt, err := db.Prepare("UPDATE users SET username = ? WHERE user_id = ?")
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
	// Add them to the database
	stmt, err := db.Prepare("INSERT INTO users (auth0_id, username, last_ip) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(auth0ID, auth0Username, ip)
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
