package models

/*
	Data types
*/

type RaceParticipantItems struct{}

/*
	"race_participant_items" table functions
*/

func (*RaceParticipantItems) Insert(userID int, raceID int, itemID int, floorNum int, stageType int) error {
	// Add the item to the list of items that they have collected so far in this race
	stmt, err := db.Prepare(`
		INSERT INTO race_participant_items (race_participant_id, item_id, floor_num, stage_type)
		VALUES ((SELECT id FROM race_participants WHERE user_id = ? AND race_id = ?), ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userID, raceID, itemID, floorNum, stageType)
	if err != nil {
		return err
	}

	return nil
}

func (*RaceParticipantItems) Reset(userID int, raceID int) error {
	// The racer reset, so remove all of their items for this race
	stmt, err := db.Prepare(`
		DELETE FROM race_participant_items
		WHERE race_participant_id = (SELECT id FROM race_participants WHERE user_id = ? AND race_id = ?)
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userID, raceID)
	if err != nil {
		return err
	}

	return nil
}
