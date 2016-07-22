package model

/*
 *  Imports
 */

import (
	"database/sql"
)

/*
 *  Data types
 */

type RaceParticipants struct {
	db *Model
}

type Racer struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Status           string `json:"status"`
	DatetimeJoined   int    `json:"datetime_joined"`
	DatetimeFinished int    `json:"datetime_finished"`
	Place            int    `json:"place"`
	Comment          string `json:"comment"`
	Items            []Item `json:"items"`
	Floor            int    `json:"floor"`
}

/*
 *  race_participants table functions
 */

func (self *RaceParticipants) GetCurrentRaces(userID int) ([]int, error) {
	// Local variables
	functionName := "modelRaceParticipantsGetCurrentRaces"

	// Get a list of the non-finished races that the user is currently in
	rows, err := db.Query("SELECT races.id FROM race_participants JOIN races ON race_participants.race_id = races.id WHERE race_participants.user_id = ? AND races.status != 'finished'", userID);
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over the races
	var raceIDs []int
	for rows.Next() {
		var raceID int
		err := rows.Scan(&raceID)
		if err != nil {
			log.Error("Database error in the", functionName, "function:", err)
			return nil, err
		}

		// Append this race to the slice
		raceIDs = append(raceIDs, raceID)
	}

	return raceIDs, nil
}

func (self *RaceParticipants) GetNotStartedRaces(userID int) ([]int, error) {
	// Local variables
	functionName := "modelRaceParticipantsGetNotStartedRaces"

	// Get a list of the non-finished and non-started races that the user is currently in
	rows, err := db.Query("SELECT races.id FROM race_participants JOIN races ON race_participants.race_id = races.id WHERE race_participants.user_id = ? AND races.status != 'finished' AND races.status != 'in progress'", userID);
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return nil, err
	}
	defer rows.Close()

	// Iterate over the races
	var raceIDs []int
	for rows.Next() {
		var raceID int
		err := rows.Scan(&raceID)
		if err != nil {
			log.Error("Database error in the", functionName, "function:", err)
			return nil, err
		}

		// Append this race to the slice
		raceIDs = append(raceIDs, raceID)
	}

	return raceIDs, nil
}

func (self *RaceParticipants) GetRacerList(raceID int) ([]Racer, error) {
	// Local variables
	functionName := "modelRaceParticipantsGetRacerList"

	// Get the people in this race
	rows, err := db.Query("SELECT users.id, users.username, race_participants.status, race_participants.datetime_joined, race_participants.datetime_finished, race_participants.place, race_participants.comment, race_participants.floor FROM race_participants JOIN users ON users.id = race_participants.user_id WHERE race_participants.race_id = ?", raceID)
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

	return racerList, nil
}

func (self *RaceParticipants) CheckInRace(userID int, raceID int) (bool, error) {
	// Local variables
	functionName := "modelRaceParticipantsCheckInRace"

	// Check to see if the user is in this race
	var id int
	err := db.QueryRow("SELECT id FROM race_participants WHERE user_id = ? AND race_id = ?", userID, raceID).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return false, err
	} else {
		return true, nil
	}
}

func (self *RaceParticipants) CheckStatus(userID int, raceID int, correctStatus string) (bool, error) {
	// Local variables
	functionName := "modelRaceParticipantsCheckStatus"

	// Check to see if the user has this status
	var status string
	err := db.QueryRow("SELECT status FROM race_participants WHERE user_id = ? AND race_id = ?", userID, raceID).Scan(&status)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return false, err
	} else if status != correctStatus {
		return false, nil
	} else {
		return true, nil
	}
}

func (self *RaceParticipants) CheckAllStatus(raceID int, correctStatus string) (bool, error) {
	// Local variables
	functionName := "modelRaceParticipantsCheckAllStatus"

	// Check to see if everyone in the race has this status
	rows, err := db.Query("SELECT status FROM race_participants WHERE race_id = ?", raceID)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return false, err
	}
	defer rows.Close()

	// Iterate over the racers
	sameStatus := true
	for rows.Next() {
		var status string
		err := rows.Scan(&status)
		if err != nil {
			log.Error("Database error in the", functionName, "function:", err)
			return false, err
		} else if status != correctStatus {
			sameStatus = false
			break
		}
	}

	return sameStatus, nil
}

func (self *RaceParticipants) CheckStillRacing(raceID int) (bool, error) {
	// Local variables
	functionName := "modelRaceParticipantsCheckAllRacing"

	// Check if anyone in the race is still racing
	var count int
	err := db.QueryRow("SELECT COUNT(id) as count FROM race_participants WHERE race_id = ? AND status == 'racing'", raceID).Scan(&count)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return false, err
	} else if count == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (self *RaceParticipants) SetStatus(userID int, raceID int, status string) error {
	// Local variables
	functionName := "modelRaceParticipantsSetStatus"

	// Set the new status for the user
	stmt, err := db.Prepare("UPDATE race_participants SET status = ? WHERE user_id = ? AND race_id = ?")
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}
	_, err = stmt.Exec(status, userID, raceID)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}

	return nil
}

func (self *RaceParticipants) SetAllStatus(raceID int, status string) error {
	// Local variables
	functionName := "modelRaceParticipantsSetAllStatus"

	// Update the status for everyone in this race
	stmt, err := db.Prepare("UPDATE race_participants SET status = ? WHERE race_id = ?")
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}
	_, err = stmt.Exec(status, raceID)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}

	return nil
}

func (self *RaceParticipants) SetComment(userID int, raceID int, comment string) error {
	// Local variables
	functionName := "modelRaceParticipantsSetComment"

	// Set the comment for the user
	stmt, err := db.Prepare("UPDATE race_participants SET comment = ? WHERE user_id = ? AND race_id = ?")
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}
	_, err = stmt.Exec(comment, userID, raceID)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}

	return nil
}

func (self *RaceParticipants) SetFloor(userID int, raceID int, floor int) error {
	// Local variables
	functionName := "modelRaceParticipantsSetFloor"

	// Set the floor for the user
	stmt, err := db.Prepare("UPDATE race_participants SET floor = ? WHERE user_id = ? AND race_id = ?")
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}
	_, err = stmt.Exec(floor, userID, raceID)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}

	return nil
}

func (self *RaceParticipants) Insert(userID int, raceID int) error {
	// Local variables
	functionName := "modelRaceParticipantsInsert"

	// Add the user to the participants list for that race
	stmt, err := db.Prepare("INSERT INTO race_participants (user_id, race_id) VALUES (?, ?)")
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}
	_, err = stmt.Exec(userID, raceID)
	if err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	}

	return nil
}

func (self *RaceParticipants) Delete(userID int, raceID int) error {
	// Local variables
	functionName := "modelRaceParticipantsDelete"

	// Remove the user from the participants list for the respective race
	if stmt, err := db.Prepare("DELETE FROM race_participants WHERE user_id = ? AND race_id = ?"); err != nil {
		log.Error("Database error in the", functionName, "function:", err)
		return err
	} else {
		_, err := stmt.Exec(userID, raceID)
		if err != nil {
			log.Error("Database error in the", functionName, "function:", err)
			return err
		}
	}

	// Check to see if anyone is still in this race
	racerList, err := self.db.RaceParticipants.GetRacerList(raceID)
	if err != nil {
		return err
	} else if len(racerList) == 0 {
		// Automatically close the race
		if stmt, err := db.Prepare("DELETE FROM races WHERE id = ?"); err != nil {
			log.Error("Database error in the", functionName, "function:", err)
			return err
		} else {
			_, err := stmt.Exec(raceID)
			if err != nil {
				log.Error("Database error in the", functionName, "function:", err)
				return err
			}
		}
	} else {
		// Check to see if this user was the captain
		var captain int
		if err := db.QueryRow("SELECT captain FROM races WHERE id = ?", raceID).Scan(&captain); err != nil {
			log.Error("Database error in the", functionName, "function:", err)
			return err
		}
		if userID == captain {
			// Change the captain to someone else
			stmt, err := db.Prepare("UPDATE races SET captain = (SELECT user_id from race_participants WHERE race_id = ? LIMIT 1) WHERE id = ?")
			if err != nil {
				log.Error("Database error in the", functionName, "function:", err)
				return err
			}
			_, err = stmt.Exec(raceID, raceID)
			if err != nil {
				log.Error("Database error in the", functionName, "function:", err)
				return err
			}
		}
	}

	return nil
}
