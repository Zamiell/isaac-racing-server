package models

import (
	"database/sql"
)

type RaceParticipantItems struct{}

// Add the item to the list of items that they have collected so far in this race
func (*RaceParticipantItems) Insert(userID int, raceID int, itemID int, floorNum int, stageType int, datetimeAcquired int64) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO race_participant_items (
			race_participant_id,
			item_id,
			floor_num,
			stage_type,
			datetime_acquired
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
		itemID,
		floorNum,
		stageType,
		datetimeAcquired,
	); err != nil {
		return err
	}

	return nil
}
