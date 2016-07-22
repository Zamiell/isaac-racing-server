package model

/*
 *  Imports
 */

// n/a

/*
 *  Data types
 */

type ChatLogPM struct {
	db *Model
}

/*
 *  chat_log_pm table functions
 */

func (self *ChatLogPM) Insert(recipient string, username string, msg string) error {
	// Local variables
	functionName := "modelChatLogPMInsert"

	// Add the message
	stmt, err := db.Prepare("INSERT INTO chat_log_pm (recipient_id, user_id, message) VALUES ((SELECT id FROM users WHERE username = ?), (SELECT id FROM users WHERE username = ?), ?)")
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}
	_, err = stmt.Exec(recipient, username, msg)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}

	return nil
}
