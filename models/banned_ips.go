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

type BannedIPs struct{}

/*
 *  banned_ips table functions
 */

func (*BannedIPs) Check(ip string) (bool, error) {
	// Check if this IP is banned
	var id int
	err := db.QueryRow("SELECT id FROM banned_ips WHERE ip = ?", ip).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (*BannedIPs) Insert(username string, adminResponsible int) error {
	// Add the IP address to the banned list in the database
	stmt, err := db.Prepare(`
		INSERT INTO banned_ips (ip, admin_responsible)
		VALUES ((SELECT last_ip FROM users WHERE username = ?), ?)
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

func (*BannedIPs) InsertIP(ip string, adminResponsible int) error {
	// Add the IP address to the banned list in the database
	stmt, err := db.Prepare("INSERT INTO banned_ips (ip, admin_responsible) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(ip, adminResponsible)
	if err != nil {
		return err
	}

	return nil
}

func (*BannedIPs) Delete(username string) error {
	// Remove the IP address from the banned IP list in the database
	stmt, err := db.Prepare("DELETE FROM banned_ips WHERE ip = (SELECT last_ip FROM users WHERE username = ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(username)
	if err != nil {
		return err
	}

	return nil
}

func (*BannedIPs) DeleteIP(ip string) error {
	// Remove the IP address from the banned IP list in the database
	stmt, err := db.Prepare("DELETE FROM banned_ips WHERE ip = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(ip)
	if err != nil {
		return err
	}

	return nil
}
