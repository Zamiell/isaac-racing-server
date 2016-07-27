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

func (self *ChatLog) Get(room string, count int) ([]ChatHistoryMessage, error) {
	// Validate function arguments
	if count == -1 {
		// TODO ?
	}

	// Get the past messages sent in this room
	rows, err := db.Query(
		"SELECT users.username, chat_log.message, chat_log.datetime FROM chat_log JOIN users ON users.id = chat_log.user_id WHERE room = ? LIMIT ?",
		room,
		count,
	)
	if err != nil {
		log.Error("Database error:", err)
		return nil, err
	}
	defer rows.Close()

	// We have to initialize this way to avoid sending a null on an empty array: https://danott.co/posts/json-marshalling-empty-slices-to-empty-arrays-in-go.html
	chatHistoryList := make([]ChatHistoryMessage, 0)
	for rows.Next() {
		var message ChatHistoryMessage
		err := rows.Scan(&message.Name, &message.Msg, &message.Datetime)
		if err != nil {
			log.Error("Database error:", err)
			return nil, err
		}
		chatHistoryList = append(chatHistoryList, message)
	}

	return chatHistoryList, nil
}

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
