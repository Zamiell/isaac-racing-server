package models

import (
	"database/sql"
)

type BannedIPs struct{}

/*
func (*BannedIPs) Insert(ip string, adminResponsible int, reason string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO banned_ips (ip, admin_responsible, reason)
		VALUES (?, ?, ?)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(ip, adminResponsible); err != nil {
		return err
	}

	return nil
}
*/

// Used in the "websocketAdminBan()" function
func (*BannedIPs) InsertUserIP(userID int, adminResponsible int, reason string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO banned_ips (ip, user_id, admin_responsible, reason)
		VALUES ((SELECT last_ip FROM users WHERE id = ?), ?, ?, ?)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		userID,
		userID,
		adminResponsible,
		reason,
	); err != nil {
		return err
	}

	return nil
}

/*
func (*BannedIPs) Delete(ip string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		DELETE FROM banned_ips
		WHERE ip = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(ip); err != nil {
		return err
	}

	return nil
}
*/

func (*BannedIPs) DeleteUserIP(userID int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		DELETE FROM banned_ips
		WHERE ip = (SELECT last_ip FROM users WHERE id = ?)
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

func (*BannedIPs) Check(ip string) (bool, error) {
	var id int
	if err := db.QueryRow(`
		SELECT id
		FROM banned_ips
		WHERE ip = ?
	`, ip).Scan(&id); err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
