package models

import (
	"database/sql"
)

type RaceParticipantItems struct{}

// Add the item to the list of items that they have collected so far in this race
func (*RaceParticipantItems) Insert(userID int, raceID int, itemID int, floorNum int, stageType int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO race_participant_items (race_participant_id, item_id, floor_num, stage_type)
		VALUES ((SELECT id FROM race_participants WHERE user_id = ? AND race_id = ?), ?, ?, ?)
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
	); err != nil {
		return err
	}

	return nil
}
