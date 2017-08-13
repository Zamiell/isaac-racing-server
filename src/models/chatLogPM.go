package models

import (
	"database/sql"
)

type ChatLogPM struct{}

// Used in the "websocketPrivateMessage" function
func (*ChatLogPM) Insert(recipientID int, userID int, message string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO chat_log_pm (recipient_id, user_id, message)
		VALUES (?, ?, ?)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(recipientID, userID, message); err != nil {
		return err
	}

	return nil
}
