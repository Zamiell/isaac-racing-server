package models

import (
	"database/sql"
)

type ChatLog struct{}

// Used in the "websocketRoomMessage" function
func (*ChatLog) Insert(room string, userID int, message string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO chat_log (room, user_id, message)
		VALUES (?, ?, ?)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(room, userID, message); err != nil {
		return err
	}

	return nil
}

// Sent in the "roomHistory" command (in the "websocketRoomJoinSub" function)
type RoomHistory struct {
	Name     string `json:"name"`
	Message  string `json:"message"`
	Datetime int64  `json:"datetime"`
}

// Get the past messages sent in this room
func (*ChatLog) Get(room string, count int) ([]RoomHistory, error) {
	roomHistoryList := make([]RoomHistory, 0)

	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			users.username,
			chat_log.message,
			UNIX_TIMESTAMP(chat_log.datetime_sent)
		FROM
			chat_log
		JOIN
			users ON users.id = chat_log.user_id
		WHERE
			room = ?
		ORDER BY
			chat_log.datetime_sent DESC
		LIMIT
			?
	`, room, count); err != nil {
		return roomHistoryList, err
	} else {
		rows = v
	}
	defer rows.Close()

	for rows.Next() {
		var message RoomHistory
		if err := rows.Scan(
			&message.Name,
			&message.Message,
			&message.Datetime,
		); err != nil {
			return roomHistoryList, err
		}
		roomHistoryList = append(roomHistoryList, message)
	}

	if err := rows.Err(); err != nil {
		return roomHistoryList, err
	}

	return roomHistoryList, nil
}
