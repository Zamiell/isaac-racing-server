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

type Users struct {
	db *Model
}

type LoginInformation struct {
	UserID    int
	Admin     int
	Squelched int
}

/*
 *  users table functions
 */

func (self *Users) Login(auth0ID string, username string, ip string) (bool, *LoginInformation, error) {
	// Local variables
	functionName := "modelUsersLogin"

	// Check to see if they are in the user database already
	var userID int
	var admin int
	var squelched int
	err := db.QueryRow("SELECT id, admin FROM users WHERE auth0_id = ?", auth0ID).Scan(&userID, &admin)
	if err == sql.ErrNoRows {
		// Add them to the database
		stmt, err := db.Prepare("INSERT INTO users (auth0_id, username, last_ip) VALUES (?, ?, ?)")
		if err != nil {
			log.Error("Database error in the", functionName, "function:", err)
			return false, nil, err
		}
		result, err := stmt.Exec(auth0ID, username, ip)
		if err != nil {
			log.Error("Failed to insert a new row for \"" + username + "\" in the users table:", err)
			return false, nil, err
		}
		userID64, err := result.LastInsertId()
		if err != nil {
			log.Error("Database error in the", functionName, "function:", err)
			return false, nil, err
		}
		userID = int(userID64)

		// By default, new users are not administrators
		admin = 0

		// By default, new users are not squelched
		squelched = 0

		// Log the user creation
		log.Info("Added \"" + username + "\" to the database (first login).")
	} else if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return false, nil, err
	} else {
		// Check to see if this user is banned
		if userIsBanned, err := self.db.BannedUsers.Check(username); err != nil {
			return false, nil, err
		} else if userIsBanned == true {
			return false, nil, nil
		}

		// Check to see if this user is squelched
		if userIsSquelched, err := self.db.SquelchedUsers.Check(username); err != nil {
			return false, nil, err
		} else if userIsSquelched == true {
			squelched = 1
		} else {
			squelched = 0
		}

		// Update the database with last_login and last_ip
		stmt, err := db.Prepare("UPDATE users SET last_login = (strftime('%s', 'now')), last_ip = ? WHERE auth0_id = ?")
		if err != nil {
			log.Error("Database error in the", functionName, "function:", err)
			return false, nil, err
		}
		_, err = stmt.Exec(ip, auth0ID)
		if err != nil {
			log.Error("Database error in the", functionName, "function:", err)
			return false, nil, err
		}
	}

	return true, &LoginInformation{userID, admin, squelched}, nil
}

func (self *Users) Exists(username string) (bool, error) {
	// Local variables
	functionName := "modelUsersExists"

	// Check if the user exists in the database
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
	if err == sql.ErrNoRows {
		log.Error("User \"" + username + "\" does not exist in the database.")
		return false, nil
	} else if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return false, err
	} else {
		return true, nil
	}
}

func (self *Users) CheckStaff(username string) (bool, error) {
	// Local variables
	functionName := "modelUsersCheckStaff"

	// Check if the user is a staff member or an administrator
	var admin int
	err := db.QueryRow("SELECT admin FROM users WHERE username = ?", username).Scan(&admin)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return false, err
	} else if admin == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (self *Users) CheckSuperAdmin(username string) (bool, error) {
	// Local variables
	functionName := "modelUsersCheckAdmin"

	// Check if the user is an super administrator
	var admin int
	err := db.QueryRow("SELECT admin FROM users WHERE username = ?", username).Scan(&admin)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return false, err
	} else if admin == 2 {
		return true, nil
	} else {
		return false, nil
	}
}

func (self *Users) SetAdmin(username string, admin int) error {
	// Local variables
	functionName := "modelUsersSetStaff"

	// Set the new ruleset for this race
	stmt, err := db.Prepare("UPDATE users SET admin = ? WHERE id = (SELECT id FROM users WHERE username = ?)")
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}
	_, err = stmt.Exec(admin, username)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}

	return nil
}
