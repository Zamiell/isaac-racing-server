package models

/*
	Data types
*/

type RaceParticipantItems struct{}

/*
	"race_participant_items" table functions
*/

func (*RaceParticipantItems) Insert(userID int, raceID int, itemID int, floor string) error {
	// Add the user to the participants list for that race
	stmt, err := db.Prepare(`
		INSERT INTO race_participant_items (race_participant_id, item_id, floor)
		VALUES ((SELECT id FROM race_participants WHERE user_id = ? AND race_id = ?), ?, ?)
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userID, raceID, itemID, floor)
	if err != nil {
		return err
	}

	return nil
}
