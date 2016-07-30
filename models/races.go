package model

/*
 *  Imports
 */

import (
	"database/sql"
	"strings"
)

/*
 *  Data types
 */

type Races struct {
	db *Model
}

/*
 *  races table functions
 */

func (self *Races) Exists(raceID int) (bool, error) {
	// Find out if the requested race exists
	var id int
	err := db.QueryRow("SELECT id FROM races WHERE id = ?", raceID).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (self *Races) GetCurrentRaces() ([]Race, error) {
	// Get the current races
	rows, err := db.Query(
		"SELECT id, name, status, " +
			"ruleset, character, goal, " +
			"seed, instant_start, datetime_created, " +
			"datetime_started, (SELECT username FROM users WHERE id = captain) as captain " +
			"FROM races WHERE status != 'finished'",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// We have to initialize this way to avoid sending a null on an empty array: https://danott.co/posts/json-marshalling-empty-slices-to-empty-arrays-in-go.html
	raceList := make([]Race, 0)
	for rows.Next() {
		var race Race
		err := rows.Scan(
			&race.ID, &race.Name, &race.Status,
			&race.Ruleset.Type, &race.Ruleset.Character, &race.Ruleset.Goal,
			&race.Ruleset.Seed, &race.Ruleset.InstantStart, &race.DatetimeCreated,
			&race.DatetimeStarted, &race.Captain,
		)
		if err != nil {
			return nil, err
		}

		// Now get the racers for this race and add it to the race
		racers, err := self.db.RaceParticipants.GetRacerNames(race.ID)
		if err != nil {
			return nil, err
		}
		race.Racers = racers

		// We are finished building the race, so add it to the race list
		raceList = append(raceList, race)
	}

	return raceList, nil
}

func (self *Races) GetStatus(raceID int) (string, error) {
	// Get the status of the race
	var status string
	err := db.QueryRow("SELECT status FROM races WHERE id = ?", raceID).Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (self *Races) CheckName(name string) (bool, error) {
	// Check to see if there are non-finished races with the same name
	rows, err := db.Query("SELECT name FROM races WHERE status != 'finished'")
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var raceName string
		err := rows.Scan(&raceName)
		if err != nil {
			return false, err
		}

		if strings.ToLower(name) == strings.ToLower(raceName) && name != "-" {
			return true, nil
		}
	}

	return false, nil
}

func (self *Races) CheckStatus(raceID int, status string) (bool, error) {
	// Check to see if the race is set to this status
	var id int
	err := db.QueryRow("SELECT id FROM races WHERE id = ? AND status = ?", raceID, status).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (self *Races) CheckCaptain(raceID int, captain int) (bool, error) {
	// Check to see if this user is the captain of the race
	var id int
	err := db.QueryRow("SELECT id FROM races WHERE id = ? AND captain = ?", raceID, captain).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (self *Races) CaptainCount(userID int) (int, error) {
	// Get how many races this user is a captain of
	var count int
	err := db.QueryRow("SELECT COUNT(id) as count FROM races WHERE status != 'finished' AND captain = ?", userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (self *Races) SetStatus(raceID int, status string) error {
	// Set the new status for this race
	stmt, err := db.Prepare("UPDATE races SET status = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(status, raceID)
	if err != nil {
		return err
	}

	return nil
}

func (self *Races) SetRuleset(raceID int, ruleset Ruleset) error {
	// Set the new ruleset for this race
	stmt, err := db.Prepare("UPDATE races SET ruleset = ?, character = ?, goal = ?, seed = ?, instant_start = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(ruleset.Type, ruleset.Character, ruleset.Goal, ruleset.Seed, ruleset.InstantStart, raceID)
	if err != nil {
		return err
	}

	return nil
}

func (self *Races) Start(raceID int) error {
	// Change the status for this race to "in progress" and set "datetime_started" equal to now
	stmt, err := db.Prepare("UPDATE races SET status = 'in progress', datetime_started = (strftime('%s', 'now')) WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(raceID)
	if err != nil {
		return err
	}

	// Update the status for everyone in the race to "racing"
	if err := self.db.RaceParticipants.SetAllStatus(raceID, "racing"); err != nil {
		return err
	}

	// Update the status for everyone in the race to "racing"
	if err := self.db.RaceParticipants.SetAllFloor(raceID, 1); err != nil {
		return err
	}

	return nil
}

func (self *Races) Finish(raceID int) error {
	// Change the status for this race to "finished" and set "datetime_finished" equal to now
	stmt, err := db.Prepare("UPDATE races SET status = 'finished', datetime_finished = (strftime('%s', 'now')) WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(raceID)
	if err != nil {
		return err
	}

	return nil
}

func (self *Races) Insert(name string, ruleset Ruleset, userID int) (int, error) {
	// Add the race to the database
	var raceID int
	if stmt, err := db.Prepare("INSERT INTO races (name, ruleset, character, goal, seed, instant_start, captain) VALUES (?, ?, ?, ?, ?, ?, ?)"); err != nil {
		return 0, err
	} else {
		result, err := stmt.Exec(name, ruleset.Type, ruleset.Character, ruleset.Goal, ruleset.Seed, ruleset.InstantStart, userID)
		if err != nil {
			return 0, err
		}
		raceID64, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}
		raceID = int(raceID64)
	}

	// Add the creator of the race to the participants list for that race
	if err := self.db.RaceParticipants.Insert(userID, raceID); err != nil {
		return 0, err
	}

	return raceID, nil
}

func (self *Races) Cleanup() error {
	// Get the current races
	rows, err := db.Query("SELECT id FROM races WHERE status = 'open' OR status = 'starting'")
	if err != nil {
		return err
	}
	defer rows.Close()

	// Iterate over the current races
	var leftoverRaces []int
	for rows.Next() {
		var raceID int
		err := rows.Scan(&raceID)
		if err != nil {
			return err
		}

		leftoverRaces = append(leftoverRaces, raceID)
	}

	// Delete all of the entries from the race_participants table (we don't want to use modelRaceParticipantsDelete because it might start the race)
	for _, raceID := range leftoverRaces {
		stmt, err := db.Prepare("DELETE FROM race_participants WHERE race_id = ?")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(raceID)
		if err != nil {
			return err
		}
	}

	// Delete the entries from the races table
	for _, raceID := range leftoverRaces {
		stmt, err := db.Prepare("DELETE FROM races WHERE id = ?")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(raceID)
		if err != nil {
			return err
		}
		log.Info("Deleted race", raceID, "during starting cleanup.")
	}

	return nil
}
