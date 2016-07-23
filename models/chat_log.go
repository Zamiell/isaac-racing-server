package model

/*
 *  Data types
 */

type ChatLog struct {
	db *Model
}

/*
 *  chat_log table functions
 */

func (self *ChatLog) Insert(room string, username string, msg string) error {
	// Add the message
	stmt, err := db.Prepare("INSERT INTO chat_log (room, user_id, message) VALUES (?, (SELECT id FROM users WHERE username = ?), ?)")
	if err != nil {
		log.Error("Database error:", err)
		return err
	}
	_, err = stmt.Exec(room, username, msg)
	if err != nil {
		log.Error("Database error:", err)
		return err
	}

	return nil
}
