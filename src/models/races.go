package models

import (
	"database/sql"

	"github.com/Zamiell/isaac-racing-server/src/log"
)

type Races struct{}

// This mirrors the "races" table row
// (it contains a subset of the information in the non-models Race struct)
type Race struct {
	ID            int
	Name          string
	Ranked        bool
	Solo          bool
	Format        string
	Character     string
	Goal          string
	StartingBuild int
	Seed          string
	Captain       string
	/* This is stored in the database as a user_id reference, but we convert it during the SELECT */
	DatetimeCreated  int64
	DatetimeStarted  int64
	DatetimeFinished int64
}

// Create a new row in the races table with no data associated with it
// (see the large comment in the "websocketRaceCreate" function for an
// explanation)
func (*Races) Insert() (int, error) {
	var result sql.Result
	if v, err := db.Exec(`
		INSERT INTO races ()
		VALUES ()
	`); err != nil {
		return 0, err
	} else {
		result = v
	}

	var raceID int
	if raceID64, err := result.LastInsertId(); err != nil {
		return 0, err
	} else {
		raceID = int(raceID64)
	}

	return raceID, nil
}

func (*Races) Delete(raceID int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		DELETE FROM races
		WHERE id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(raceID); err != nil {
		return err
	}

	return nil
}

/*
// Used in the "?" function
// TODO probably needed for set comment
func (*Races) Exists(raceID int) (bool, error) {
	var id int
	if err := db.QueryRow(`
		SELECT id
		FROM races
		WHERE id = ?
	`, raceID).Scan(&id); err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
*/

// Now that the race is over, fill in the blank race in the database with all of
// the information that the server had on hand
func (*Races) Finish(race *Race) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE
			races
		SET
			finished = 1,
			name = ?,
			ranked = ?,
			solo = ?,
			format = ?,
			player_type = ?,
			goal = ?,
			starting_build = ?,
			seed = ?,
			captain = (SELECT id FROM users where username = ?),
			datetime_started = FROM_UNIXTIME(? / 1000),
			datetime_finished = NOW()
		WHERE
			id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	// Convert some bools to ints
	ranked := 0
	if race.Ranked {
		ranked = 1
	}
	solo := 0
	if race.Solo {
		solo = 1
	}

	if _, err := stmt.Exec(
		race.Name,
		ranked,
		solo,
		race.Format,
		race.Character,
		race.Goal,
		race.StartingBuild,
		race.Seed,
		race.Captain,
		race.DatetimeStarted,
		race.ID,
	); err != nil {
		return err
	}

	log.Debug("race.DatetimeStarted:", race.DatetimeStarted)

	return nil
}

// Clean up any unfinished races from the database
func (*Races) Cleanup() ([]int, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT id
		FROM races
		WHERE finished = 0
		ORDER BY id
	`); err != nil {
		return nil, err
	} else {
		rows = v
	}
	defer rows.Close()

	var leftoverRaces []int
	for rows.Next() {
		var raceID int
		if err := rows.Scan(&raceID); err != nil {
			return nil, err
		}
		leftoverRaces = append(leftoverRaces, raceID)
	}

	// Delete the entries from the races table
	for _, raceID := range leftoverRaces {
		var stmt *sql.Stmt
		if v, err := db.Prepare(`
			DELETE FROM races
			WHERE id = ?
		`); err != nil {
			return nil, err
		} else {
			stmt = v
		}
		defer stmt.Close()

		if _, err := stmt.Exec(raceID); err != nil {
			return nil, err
		}
	}

	return leftoverRaces, nil
}
