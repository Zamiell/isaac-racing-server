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

type BannedIPs struct {
	db *Model
}

/*
 *  banned_ips table functions
 */

func (self *BannedIPs) Check(ip string) (bool, error) {
	// Local variables
	functionName := "modelBannedIPsCheck"

	// Check if this IP is banned
	var id int
	err := db.QueryRow("SELECT id FROM banned_ips WHERE ip = ?", ip).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return false, err
	} else {
		log.Info("IP \"" + ip + "\" tried to log in, but they are banned.")
		return true, nil
	}
}

func (self *BannedIPs) Insert(username string, adminResponsible int) error {
	// Local variables
	functionName := "modelBannedIPsInsert"

	// Add the IP address to the banned list in the database
	stmt, err := db.Prepare("INSERT INTO banned_ips (ip, admin_responsible) VALUES ((SELECT last_ip FROM users WHERE username = ?), ?)")
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}
	_, err = stmt.Exec(username, adminResponsible)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}

	return nil
}

func (self *BannedIPs) InsertIP(ip string, adminResponsible int) error {
	// Local variables
	functionName := "modelBannedIPsInsertIP"

	// Add the IP address to the banned list in the database
	stmt, err := db.Prepare("INSERT INTO banned_ips (ip, admin_responsible) VALUES (?, ?)")
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}
	_, err = stmt.Exec(ip, adminResponsible)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}

	return nil
}

func (self *BannedIPs) Delete(username string) error {
	// Local variables
	functionName := "modelBannedIPsDelete"

	// Remove the IP address from the banned IP list in the database
	stmt, err := db.Prepare("DELETE FROM banned_ips WHERE ip = (SELECT last_ip FROM users WHERE username = ?)")
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}
	_, err = stmt.Exec(username)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}

	return nil
}

func (self *BannedIPs) DeleteIP(ip string) error {
	// Local variables
	functionName := "modelBannedIPsDelete"

	// Remove the IP address from the banned IP list in the database
	stmt, err := db.Prepare("DELETE FROM banned_ips WHERE ip = ?")
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}
	_, err = stmt.Exec(ip)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}

	return nil
}
