package models

/*
	Data types
*/

type RaceParticipantRooms struct{}

/*
	"race_participant_rooms" table functions
*/

func (*RaceParticipantRooms) Insert(userID int, raceID int, roomID string) error {
	// Add the room to the list of rooms visited for this person's race
	stmt, err := db.Prepare(`
		INSERT INTO race_participant_rooms (race_participant_id, room_id)
		VALUES ((SELECT id FROM race_participants WHERE user_id = ? AND race_id = ?), ?)
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userID, raceID, roomID)
	if err != nil {
		return err
	}

	return nil
}

func (*RaceParticipantRooms) Reset(userID int, raceID int) error {
	// The racer reset, so remove all of their rooms visited for this race
	stmt, err := db.Prepare(`
		DELETE FROM race_participant_rooms
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
