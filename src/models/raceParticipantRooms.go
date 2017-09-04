package models

import (
	"database/sql"
)

type RaceParticipantRooms struct{}

// Add the room to the list of rooms visited for this person's race
func (*RaceParticipantRooms) Insert(userID int, raceID int, roomID string, floorNum int, stageType int, datetimeArrived int64) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO race_participant_rooms (
			race_participant_id,
			room_id,
			floor_num,
			stage_type,
			datetime_arrived
		)
		VALUES (
			(SELECT id FROM race_participants WHERE user_id = ? AND race_id = ?),
			?,
			?,
			?,
			FROM_UNIXTIME(? / 1000)
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
		datetimeArrived,
	); err != nil {
		return err
	}

	return nil
}
