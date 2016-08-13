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

type SquelchedUsers struct{}

/*
 *  squelched_users table functions
 */

func (*SquelchedUsers) Check(username string) (bool, error) {
	// Check if this user is squelched
	var id int
	err := db.QueryRow(`
		SELECT id
		FROM squelched_users
		WHERE user_id = (SELECT id FROM users WHERE username = ?)
	`, username).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (*SquelchedUsers) Insert(username string, adminResponsible int) error {
	// Add the user to the squelched list in the database
	stmt, err := db.Prepare(`
		INSERT INTO squelched_users (user_id, admin_responsible)
		VALUES ((SELECT id from users WHERE username = ?), ?)
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(username, adminResponsible)
	if err != nil {
		return err
	}

	return nil
}

func (*SquelchedUsers) Delete(username string) error {
	// Remove the user from the squelched list in the database
	stmt, err := db.Prepare(`
		DELETE FROM squelched_users
		WHERE user_id = (SELECT id from users WHERE username = ?)
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(username)
	if err != nil {
		return err
	}

	return nil
}
