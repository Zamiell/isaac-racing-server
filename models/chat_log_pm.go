package models

/*
 *  Data types
 */

type ChatLogPM struct{}

/*
 *  chat_log_pm table functions
 */

func (*ChatLogPM) Insert(recipient string, username string, msg string) error {
	// Add the message
	stmt, err := db.Prepare(`
		INSERT INTO chat_log_pm (recipient_id, user_id, message)
		VALUES ((SELECT id FROM users WHERE username = ?), (SELECT id FROM users WHERE username = ?), ?)
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(recipient, username, msg)
	if err != nil {
		return err
	}

	return nil
}
