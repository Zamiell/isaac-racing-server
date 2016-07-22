package model

/*
 *  Imports
 */

// n/a

/*
 *  Data types
 */

type RaceParticipantItems struct {
	db *Model
}

type Item struct {
	ID    int `json:"id"`
	Floor int `json:"floor"`
}

/*
 *  race_participant_items table functions
 */

// TODO
func (self *RaceParticipantItems) GetItemList(userID int, raceID int) ([]Item, error) {
	// Local variables
	//functionName := "modelRaceParticipantItemsGetItemList"

	// Get the people in this race
	/*rows, err := db.Query("SELECT users.id, users.username, race_participants.status, race_participants.datetime_joined, race_participants.datetime_finished, race_participants.place, race_participants.comment, race_participants.floor FROM race_participants JOIN users ON users.id = race_participants.user_id WHERE race_participants.race_id = ?", raceID)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return nil, err
	}
	defer rows.Close()

	// We have to initialize this way to avoid sending a null on an empty array: https://danott.co/posts/json-marshalling-empty-slices-to-empty-arrays-in-go.html
	racerList := make([]Racer, 0)
	for rows.Next() {
		var racer Racer
		err := rows.Scan(&racer.ID, &racer.Name, &racer.Status, &racer.DatetimeJoined, &racer.DatetimeFinished, &racer.Place, &racer.Comment, &racer.Floor)
		if err != nil {
			log.Error("Database error in the", functionName, "function:", err)
			return nil, err
		}
		racerList = append(racerList, racer)
	}

	return racerList, nil*/
	return nil, nil
}

func (self *RaceParticipantItems) Insert(userID int, raceID int, itemID int) error {
	// Local variables
	functionName := "modelRaceParticipantItemsInsert"

	// Add the user to the participants list for that race
	stmt, err := db.Prepare("INSERT INTO race_participant_items (race_participant_id, item_id, floor) VALUES ((SELECT id FROM race_participants WHERE user_id = ? AND race_id = ?), ?, (SELECT floor FROM race_participants WHERE user_id = ? AND race_id = ?))")
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}
	_, err = stmt.Exec(userID, raceID, itemID, userID, raceID)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}

	return nil
}
