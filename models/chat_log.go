package models

/*
 *  Data types
 */

type ChatLog struct{}

/*
 *  chat_log table functions
 */

func (*ChatLog) Get(room string, count int) ([]RoomHistory, error) {
	// Get the past messages sent in this room
	rows, err := db.Query(`
		SELECT users.username, chat_log.message, chat_log.datetime
		FROM chat_log
			JOIN users ON users.id = chat_log.user_id
		WHERE room = ?
		ORDER BY chat_log.datetime
		LIMIT ?
	`, room, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// We have to initialize this way to avoid sending a null on an empty array: https://danott.co/posts/json-marshalling-empty-slices-to-empty-arrays-in-go.html
	roomHistoryList := make([]RoomHistory, 0)
	for rows.Next() {
		var message RoomHistory
		err := rows.Scan(&message.Name, &message.Msg, &message.Datetime)
		if err != nil {
			return nil, err
		}
		roomHistoryList = append(roomHistoryList, message)
	}

	return roomHistoryList, nil
}

func (*ChatLog) Insert(room string, username string, msg string) error {
	// Add the message
	stmt, err := db.Prepare(`
		INSERT INTO chat_log (room, user_id, message)
		VALUES (?, (SELECT id FROM users WHERE username = ?), ?)
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(room, username, msg)
	if err != nil {
		return err
	}

	return nil
}
