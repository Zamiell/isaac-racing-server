package model

/*
 *  Imports
 */

import (
	"database/sql"
)

/*
 *  Data types
 */

type BannedUsers struct {
	db *Model
}

/*
 *  banned_users table functions
 */

func (self *BannedUsers) Check(username string) (bool, error) {
	// Check if this user is banned
	var id int
	err := db.QueryRow("SELECT id FROM banned_users WHERE user_id = (SELECT id FROM users WHERE username = ?)", username).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		log.Error("Database error:", err)
		return false, err
	} else {
		return true, nil
	}
}

func (self *BannedUsers) Insert(username string, adminResponsible int) error {
	// Add this user to the ban list in the database
	stmt, err := db.Prepare("INSERT INTO banned_users (user_id, admin_responsible) VALUES ((SELECT id from users WHERE username = ?), ?)")
	if err != nil {
		log.Error("Database error:", err)
		return err
	}
	_, err = stmt.Exec(username, adminResponsible)
	if err != nil {
		log.Error("Database error:", err)
		return err
	}

	return nil
}

func (self *BannedUsers) Delete(username string) error {
	// Remove the user from the banned users list in the database
	stmt, err := db.Prepare("DELETE FROM banned_users WHERE user_id = (SELECT id from users WHERE username = ?)")
	if err != nil {
		log.Error("Database error:", err)
		return err
	}
	_, err = stmt.Exec(username)
	if err != nil {
		log.Error("Database error:", err)
		return err
	}

	return nil
}
