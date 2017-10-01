package models

import (
	"database/sql"
)

type BannedUsers struct{}

// Called from the "websocketAdminBan()" function
func (*BannedUsers) Insert(userID int, adminResponsible int, reason string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO banned_users (user_id, admin_responsible, reason)
		VALUES (?, ?, ?)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		userID,
		adminResponsible,
		reason,
	); err != nil {
		return err
	}

	return nil
}

// Called from the "websocketAdminUnban()" function
func (*BannedUsers) Delete(userID int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		DELETE FROM banned_users
		WHERE user_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(userID); err != nil {
		return err
	}

	return nil
}

// Called from the "httpValidateSession()" function
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
