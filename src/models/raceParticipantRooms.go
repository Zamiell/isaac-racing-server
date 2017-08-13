package models

import (
	"database/sql"
)

type RaceParticipantRooms struct{}

// Add the room to the list of rooms visited for this person's race
func (*RaceParticipantRooms) Insert(userID int, raceID int, roomID string, floorNum int, stageType int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO race_participant_rooms (
			race_participant_id,
			room_id,
			floor_num,
			stage_type
		)
		VALUES (
			(SELECT id FROM race_participants WHERE user_id = ? AND race_id = ?),
			?,
			?,
			?
		)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		userID,
		raceID,
		roomID,
		floorNum,
		stageType,
	); err != nil {
		return err
	}

	return nil
}
