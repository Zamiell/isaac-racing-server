package models

/*
	Imports
*/

import (
	"database/sql"
	"strings"
)

/*
	Data types
*/

type Races struct{}

/*
	"races" table functions
*/

func (*Races) Exists(raceID int) (bool, error) {
	// Find out if the requested race exists
	var id int
	err := db.QueryRow(`
		SELECT id
		FROM races
		WHERE id = ?
	`, raceID).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (*Races) GetCurrentRaces() ([]Race, error) {
	// Get the current races
	rows, err := db.Query(`
		SELECT
			id,
			name,
			status,
			type,
			solo,
			format,
			character,
			goal,
			starting_build,
			seed,
			datetime_created,
			datetime_started,
			(SELECT username FROM users WHERE id = captain) as captain
		FROM races
		WHERE status != 'finished'
		ORDER BY datetime_created
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// We have to initialize this way to avoid sending a null on an empty array: https://danott.co/posts/json-marshalling-empty-slices-to-empty-arrays-in-go.html
	raceList := make([]Race, 0)
	for rows.Next() {
		var race Race
		err2 := rows.Scan(
			&race.ID,
			&race.Name,
			&race.Status,
			&race.Ruleset.Type,
			&race.Ruleset.Solo,
			&race.Ruleset.Format,
			&race.Ruleset.Character,
			&race.Ruleset.Goal,
			&race.Ruleset.StartingBuild,
			&race.Seed,
			&race.DatetimeCreated,
			&race.DatetimeStarted,
			&race.Captain,
		)
		if err2 != nil {
			return nil, err
		}

		// Add it to the list
		raceList = append(raceList, race)
	}

	// Get the names of the people in this race
	rows, err = db.Query(`
		SELECT races.id, users.username
		FROM races
			JOIN race_participants ON race_participants.race_id = races.id
			JOIN users ON users.id = race_participants.user_id
		WHERE races.status != 'finished'
		ORDER BY race_participants.datetime_joined
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// For each name that we found, append it to the appropriate place in the raceList object
	for rows.Next() {
		var raceID int
		var username string
		err := rows.Scan(&raceID, &username)
		if err != nil {
			return nil, err
		}

		// Find the race in the raceList object
		for i, race := range raceList {
			if race.ID == raceID {
				raceList[i].Racers = append(race.Racers, username)
				break
			}
		}
	}

	return raceList, nil
}

func (*Races) GetStatus(raceID int) (string, error) {
	// Get the status of the race
	var status string
	err := db.QueryRow(`
		SELECT status
		FROM races
		WHERE id = ?
	`, raceID).Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (*Races) GetDatetimeStarted(raceID int) (int, error) {
	// Get the time that the race started
	var datetimeStarted int
	err := db.QueryRow(`
		SELECT datetime_started
		FROM races
		WHERE id = ?
	`, raceID).Scan(&datetimeStarted)
	if err != nil {
		return 0, err
	}

	return datetimeStarted, nil
}

func (*Races) GetRuleset(raceID int) (Ruleset, error) {
	// Get the ruleset of the race
	var ruleset Ruleset
	err := db.QueryRow(`
		SELECT solo, format, character, goal, starting_build
		FROM races
		WHERE id = ?
	`, raceID).Scan(&ruleset.Solo, &ruleset.Format, &ruleset.Character, &ruleset.Goal, &ruleset.StartingBuild)
	if err != nil {
		return ruleset, err
	}

	return ruleset, nil
}

func (*Races) CheckName(name string) (bool, error) {
	// Check to see if there are non-finished races with the same name
	rows, err := db.Query(`
		SELECT name
		FROM races
		WHERE status != 'finished'
	`)
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

func (*Races) CheckStatus(raceID int, status string) (bool, error) {
	// Check to see if the race is set to this status
	var id int
	err := db.QueryRow(`
		SELECT id
		FROM races
		WHERE id = ? AND status = ?
	`, raceID, status).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (*Races) CheckCaptain(raceID int, captain int) (bool, error) {
	// Check to see if this user is the captain of the race
	var id int
	err := db.QueryRow(`
		SELECT id
		FROM races
		WHERE id = ? AND captain = ?
	`, raceID, captain).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (*Races) CheckSolo(raceID int) (bool, error) {
	// Check to see if this is a solo race
	// (only 1 person is able to join)
	var id int
	err := db.QueryRow(`
		SELECT solo
		FROM races
		WHERE id = ?
	`, raceID).Scan(&id)
	if err != nil {
		return false, err
	} else if id == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (*Races) CheckSeededRanked(raceID int) (bool, error) {
	// Check to see if this is a seeded and ranked race
	var raceType int
	var raceFormat string
	err := db.QueryRow(`
		SELECT type, format
		FROM races
		WHERE id = ?
	`, raceID).Scan(&raceType, &raceFormat)
	if err != nil {
		return false, err
	} else if raceType == 1 && raceFormat == "seeded" {
		return true, nil
	} else {
		return false, nil
	}
}

func (*Races) CheckUnseededRanked(raceID int) (bool, error) {
	// Check to see if this is an unseeded and ranked race
	var raceType int
	var raceFormat string
	err := db.QueryRow(`
		SELECT type, format
		FROM races
		WHERE id = ?
	`, raceID).Scan(&raceType, &raceFormat)
	if err != nil {
		return false, err
	} else if raceType == 1 && raceFormat == "unseeded" {
		return true, nil
	} else {
		return false, nil
	}
}

func (*Races) SetStatus(raceID int, status string) error {
	// Set the new status for this race
	stmt, err := db.Prepare(`
		UPDATE races
		SET status = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(status, raceID)
	if err != nil {
		return err
	}

	return nil
}

func (*Races) SetRuleset(raceID int, ruleset Ruleset) error {
	// Set the new ruleset for this race
	stmt, err := db.Prepare(`
		UPDATE races
		SET ruleset = ?, character = ?, goal = ?, starting_build = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(ruleset.Format, ruleset.Character, ruleset.Goal, ruleset.StartingBuild, raceID)
	if err != nil {
		return err
	}

	return nil
}

func (*Races) SetSeed(raceID int, seed string) error {
	// Set a seed for this race
	stmt, err := db.Prepare(`
		UPDATE races
		SET seed = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(seed, raceID)
	if err != nil {
		return err
	}

	return nil
}

func (*Races) Start(raceID int) error {
	// Change the status for this race to "in progress" and set "datetime_started" equal to now
	stmt, err := db.Prepare(`
		UPDATE races
		SET status = 'in progress', datetime_started = ?
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(makeTimestamp(), raceID)
	if err != nil {
		return err
	}

	return nil
}

func (*Races) Finish(raceID int) error {
	// Change the status for this race to "finished" and set "datetime_finished" equal to now
	stmt, err := db.Prepare(`
		UPDATE races
		SET status = 'finished', datetime_finished = (strftime('%s', 'now'))
		WHERE id = ?
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(raceID)
	if err != nil {
		return err
	}

	return nil
}

func (*Races) Insert(name string, ruleset Ruleset, userID int) (int, error) {
	// Add the race to the database
	var raceID int
	if stmt, err := db.Prepare(`
		INSERT INTO races (
			name,
			type,
			solo,
			format,
			character,
			goal,
			starting_build,
			captain,
			datetime_created
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`); err != nil {
		return 0, err
	} else {
		result, err := stmt.Exec(
			name,
			ruleset.Type,
			ruleset.Solo,
			ruleset.Format,
			ruleset.Character,
			ruleset.Goal,
			ruleset.StartingBuild,
			userID,
			makeTimestamp(),
		)
		if err != nil {
			return 0, err
		}
		raceID64, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}
		raceID = int(raceID64)
	}

	return raceID, nil
}

func (*Races) Cleanup() ([]int, error) {
	// Get the current races
	rows, err := db.Query(`
		SELECT id
		FROM races
		WHERE status = 'open' OR status = 'starting'
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the current races
	var leftoverRaces []int
	for rows.Next() {
		var raceID int
		err := rows.Scan(&raceID)
		if err != nil {
			return nil, err
		}

		leftoverRaces = append(leftoverRaces, raceID)
	}

	// Delete all of the entries from the race_participants table (we don't want to use RaceParticipants.Delete because we don't care about captains)
	for _, raceID := range leftoverRaces {
		stmt, err := db.Prepare(`
			DELETE FROM race_participants
			WHERE race_id = ?
		`)
		if err != nil {
			return nil, err
		}
		_, err = stmt.Exec(raceID)
		if err != nil {
			return nil, err
		}
	}

	// Delete the entries from the races table
	for _, raceID := range leftoverRaces {
		stmt, err := db.Prepare(`
			DELETE FROM races
			WHERE id = ?
		`)
		if err != nil {
			return nil, err
		}
		_, err = stmt.Exec(raceID)
		if err != nil {
			return nil, err
		}
	}

	return leftoverRaces, nil
}

func (*Races) GetRaceHistory(currentPage int, racesPerPage int) ([]RaceHistory, int, error) {
	raceOffset := currentPage * racesPerPage
	rows, err := db.Query(`
		SELECT 
			r.id, 
			r.datetime_started,
			r.type, 
			r.format, 
			r.character, 
			r.goal
		FROM 
			races r
		GROUP BY
			r.id
		ORDER BY
			r.id DESC
		LIMIT
			?
		OFFSET
			?
	`, racesPerPage, raceOffset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()	
	raceHistory := make([]RaceHistory, 0)
	for rows.Next() {
		var race RaceHistory
		err = rows.Scan(
			&race.RaceID,
			&race.RaceDate,
			&race.RaceType,
			&race.RaceFormat,			
			&race.RaceChar,		
			&race.RaceGoal,			
		)
		if err != nil {
			return raceHistory, 0, err
		}
		raceHistory = append(raceHistory, race)
	}
	rows, err = db.Query(`
		SELECT 
			count(id) 
		FROM 
			races
		WHERE
			status = 'finished'
	`)
	if err != nil {
		return raceHistory, 0, err
	}
	defer rows.Close()
	var allRaceCount int
	for rows.Next() {
		err = rows.Scan(&allRaceCount)
		if err != nil {
			return raceHistory, allRaceCount, err
		}
	}	
	return raceHistory, allRaceCount, nil

}