package models

import (
	"database/sql"
)

type BannedUsers struct{}

/*
func (*BannedUsers) Insert(username string, adminResponsible int, reason string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO banned_users (user_id, admin_responsible, reason)
		VALUES ((SELECT id from users WHERE username = ?), ?, ?)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		username,
		adminResponsible,
		reason,
	); err != nil {
		return err
	}

	return nil
}

func (*BannedUsers) Delete(username string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		DELETE FROM banned_users
		WHERE user_id = (SELECT id from users WHERE username = ?)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(username); err != nil {
		return err
	}

	return nil
}
*/

// Called from the "httpValidateSession" function
func (*BannedUsers) Check(userID int) (bool, error) {
	var id int
	if err := db.QueryRow(`
		SELECT id
		FROM banned_users
		WHERE user_id = ?
	`, userID).Scan(&id); err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
