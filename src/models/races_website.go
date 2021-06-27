package models

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

/*
	These are more functions for querying the "races" table,
	but these functions are only used for the website
*/

// RaceHistory gets the history for each race in the database
type RaceHistory struct {
	RaceID           sql.NullInt64
	RaceSize         sql.NullInt64
	RaceType         sql.NullString
	RaceFormat       sql.NullString
	RaceChar         sql.NullString
	RaceGoal         sql.NullString
	RaceDateStart    mysql.NullTime
	RaceDateFinished mysql.NullTime
	RaceParticipants []RaceHistoryParticipants
}

// RaceHistoryParticipants gets the user stats for each racer in each race
type RaceHistoryParticipants struct {
	ID                     sql.NullInt64
	RacerName              sql.NullString
	RacerPlace             sql.NullInt64
	RacerTime              sql.NullString
	RacerSeed              sql.NullString
	RacerStartingItem      sql.NullInt64
	RacerStartingItemName  string
	RacerStartingBuild     sql.NullInt64
	RacerStartingBuildID   int
	RacerStartingBuildName string
	RacerComment           sql.NullString
}

// GetRacesHistory gets all data for all races
func (*Races) GetRacesHistory(currentPage int, racesPerPage int, raceOffset int) ([]RaceHistory, int, error) {
	raceHistory := make([]RaceHistory, 0)

	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			r.id,
			r.solo,
			r.ranked,
			r.format,
			r.player_type,
			r.goal,
			r.datetime_created,
			r.datetime_finished
		FROM
			races r
		WHERE
			r.finished = 1
			AND r.solo = 0
		GROUP BY
			r.id
		ORDER BY
			r.datetime_created DESC
		LIMIT
			?
		OFFSET
			?
	`, racesPerPage, raceOffset); err != nil {
		return raceHistory, 0, err
	} else {
		rows = v
	}
	defer rows.Close()

	for rows.Next() {
		raceRacers := make([]RaceHistoryParticipants, 0)

		var race RaceHistory
		if err := rows.Scan(
			&race.RaceID,
			&race.RaceSize,
			&race.RaceType,
			&race.RaceFormat,
			&race.RaceChar,
			&race.RaceGoal,
			&race.RaceDateStart,
			&race.RaceDateFinished,
		); err != nil {
			return raceHistory, 0, err
		}
		race.RaceParticipants = nil

		var rows2 *sql.Rows
		if v, err := db.Query(`
			SELECT
				u.username,
				rp.seed,
				rp.place,
				rp.run_time,
				rp.starting_item,
				r.starting_build,
				rp.comment
			FROM
				race_participants rp
			LEFT JOIN
				users u
					ON u.id = rp.user_id
			LEFT JOIN
				races r
					ON r.id = rp.race_id
			WHERE
				rp.race_id = ?
			ORDER BY
				CASE WHEN rp.place = -1 THEN 1 ELSE 0 END,
				rp.place,
				rp.run_time;
		`, race.RaceID); err != nil {
			return raceHistory, 0, err
		} else {
			rows2 = v
		}
		defer rows2.Close()

		for rows2.Next() {
			var racer RaceHistoryParticipants
			if err := rows2.Scan(
				&racer.RacerName,
				&racer.RacerSeed,
				&racer.RacerPlace,
				&racer.RacerTime,
				&racer.RacerStartingItem,
				&racer.RacerStartingBuild,
				&racer.RacerComment,
			); err != nil {
				return raceHistory, 0, err
			}
			raceRacers = append(raceRacers, racer)
		}

		if err := rows2.Err(); err != nil {
			return raceHistory, 0, err
		}

		race.RaceParticipants = raceRacers
		raceHistory = append(raceHistory, race)
	}

	if err := rows.Err(); err != nil {
		return raceHistory, 0, err
	}

	var allRaceCount int
	if err := db.QueryRow(`
		SELECT
			count(id)
		FROM
			races
		WHERE
			finished = 1
			AND solo = 0

	`).Scan(&allRaceCount); err != nil {
		return raceHistory, 0, err
	}

	return raceHistory, allRaceCount, nil
}

// GetRaceHistory gets race history for a single race
func (*Races) GetRaceHistory(raceID int) (RaceHistory, error) {
	var race RaceHistory
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			r.id,
			r.ranked,
			r.format,
			r.player_type,
			r.goal,
			r.datetime_created,
			r.datetime_finished
		FROM
			races r
		WHERE
			r.id = ?
		GROUP BY
			r.id
		ORDER BY
			r.datetime_created DESC
		LIMIT
			1
	`, raceID); err != nil {
		return race, err
	} else {
		rows = v
	}
	defer rows.Close()

	for rows.Next() {
		raceRacers := make([]RaceHistoryParticipants, 0)

		if err := rows.Scan(
			&race.RaceID,
			&race.RaceType,
			&race.RaceFormat,
			&race.RaceChar,
			&race.RaceGoal,
			&race.RaceDateStart,
			&race.RaceDateFinished,
		); err != nil {
			return race, err
		}
		race.RaceParticipants = nil

		var rows2 *sql.Rows
		if v, err := db.Query(`
			SELECT
				u.username,
				rp.seed,
				rp.place,
				rp.run_time,
				rp.starting_item,
				r.starting_build,
				rp.comment
			FROM
				race_participants rp
			LEFT JOIN
				users u
					ON u.id = rp.user_id
			LEFT JOIN
				races r
					ON r.id = rp.race_id
			WHERE
				rp.race_id = ?
			ORDER BY
				CASE WHEN rp.place = -1 THEN 1 ELSE 0 END,
				rp.place,
				rp.run_time;
		`, race.RaceID); err != nil {
			return race, err
		} else {
			rows2 = v
		}
		defer rows2.Close()

		for rows2.Next() {
			var racer RaceHistoryParticipants
			if err := rows2.Scan(
				&racer.RacerName,
				&racer.RacerSeed,
				&racer.RacerPlace,
				&racer.RacerTime,
				&racer.RacerStartingItem,
				&racer.RacerStartingBuild,
				&racer.RacerComment,
			); err != nil {
				return race, err
			}
			raceRacers = append(raceRacers, racer)
		}

		if err := rows2.Err(); err != nil {
			return race, err
		}

		race.RaceParticipants = raceRacers
	}

	if err := rows.Err(); err != nil {
		return race, err
	}

	return race, nil
}

// GetRaceProfileHistory gets the race data for the profile page
func (*Races) GetRankedRaceProfileHistory(user string, racesPerPage int) ([]RaceHistory, error) { // nolint: dupl
	raceHistory := make([]RaceHistory, 0)

	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			r.id,
			r.ranked,
			r.format,
			r.player_type,
			r.goal,
			r.datetime_created,
			r.datetime_finished
		FROM
			races r
		LEFT JOIN
			race_participants rp
				ON rp.race_id = r.id
		LEFT JOIN
			users u
				ON u.id = rp.user_id
		WHERE
			r.finished = 1
			AND r.ranked = 1
			AND r.solo = 1
			AND u.username = ?
		GROUP BY
			id
		ORDER BY
			datetime_created DESC
		LIMIT
			?
	`, user, racesPerPage); err != nil {
		return raceHistory, err
	} else {
		rows = v
	}
	defer rows.Close()

	for rows.Next() {
		raceRacers := make([]RaceHistoryParticipants, 0)

		var race RaceHistory
		if err := rows.Scan(
			&race.RaceID,
			&race.RaceType,
			&race.RaceFormat,
			&race.RaceChar,
			&race.RaceGoal,
			&race.RaceDateStart,
			&race.RaceDateFinished,
		); err != nil {
			return raceHistory, err
		}
		race.RaceParticipants = nil

		var rows2 *sql.Rows
		if v, err := db.Query(`
			SELECT
				u.username,
				rp.place,
				rp.run_time,
				rp.starting_item,
				r.starting_build,
				rp.comment
			FROM
				race_participants rp
			LEFT JOIN
				users u
					ON u.id = rp.user_id
			LEFT JOIN
				races r
					ON r.id = rp.race_id
			WHERE
				rp.race_id = ?
			ORDER BY
				CASE WHEN rp.place = -1 THEN 1 ELSE 0 END,
				rp.place,
				rp.run_time;
		`, race.RaceID); err != nil {
			return raceHistory, err
		} else {
			rows2 = v
		}
		defer rows2.Close()

		for rows2.Next() {
			var racer RaceHistoryParticipants
			if err := rows2.Scan(
				&racer.RacerName,
				&racer.RacerPlace,
				&racer.RacerTime,
				&racer.RacerStartingItem,
				&racer.RacerStartingBuild,
				&racer.RacerComment,
			); err != nil {
				return raceHistory, err
			}
			raceRacers = append(raceRacers, racer)
		}

		if err := rows2.Err(); err != nil {
			return raceHistory, err
		}

		race.RaceParticipants = raceRacers
		raceHistory = append(raceHistory, race)
	}

	if err := rows.Err(); err != nil {
		return raceHistory, err
	}

	var allRaceCount int
	if err := db.QueryRow(`
		SELECT count(id)
		FROM races
		WHERE finished = 1
	`).Scan(&allRaceCount); err != nil {
		return raceHistory, err
	}

	return raceHistory, nil
}

func (*Races) GetAllRaceProfileHistory(user string, racesPerPage int) ([]RaceHistory, error) { // nolint: dupl
	raceHistory := make([]RaceHistory, 0)

	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			r.id,
			r.ranked,
			r.format,
			r.player_type,
			r.goal,
			r.datetime_created,
			r.datetime_finished
		FROM
			races r
		LEFT JOIN
			race_participants rp
				ON rp.race_id = r.id
		LEFT JOIN
			users u
				ON u.id = rp.user_id
		WHERE
			r.finished = 1
			AND u.username = ?
		GROUP BY
			id
		ORDER BY
			datetime_created DESC
		LIMIT
			?
	`, user, racesPerPage); err != nil {
		return raceHistory, err
	} else {
		rows = v
	}
	defer rows.Close()

	for rows.Next() {
		raceRacers := make([]RaceHistoryParticipants, 0)

		var race RaceHistory
		if err := rows.Scan(
			&race.RaceID,
			&race.RaceType,
			&race.RaceFormat,
			&race.RaceChar,
			&race.RaceGoal,
			&race.RaceDateStart,
			&race.RaceDateFinished,
		); err != nil {
			return raceHistory, err
		}
		race.RaceParticipants = nil

		var rows2 *sql.Rows
		if v, err := db.Query(`
			SELECT
				u.username,
				rp.place,
				rp.run_time,
				rp.starting_item,
				r.starting_build,
				rp.comment
			FROM
				race_participants rp
			LEFT JOIN
				users u
					ON u.id = rp.user_id
			LEFT JOIN
				races r
					ON r.id = rp.race_id
			WHERE
				rp.race_id = ?
			ORDER BY
				CASE WHEN rp.place = -1 THEN 1 ELSE 0 END,
				rp.place,
				rp.run_time;
		`, race.RaceID); err != nil {
			return raceHistory, err
		} else {
			rows2 = v
		}
		defer rows2.Close()

		for rows2.Next() {
			var racer RaceHistoryParticipants
			if err := rows2.Scan(
				&racer.RacerName,
				&racer.RacerPlace,
				&racer.RacerTime,
				&racer.RacerStartingItem,
				&racer.RacerStartingBuild,
				&racer.RacerComment,
			); err != nil {
				return raceHistory, err
			}
			raceRacers = append(raceRacers, racer)
		}

		if err := rows2.Err(); err != nil {
			return raceHistory, err
		}

		race.RaceParticipants = raceRacers
		raceHistory = append(raceHistory, race)
	}

	if err := rows.Err(); err != nil {
		return raceHistory, err
	}

	var allRaceCount int
	if err := db.QueryRow(`
		SELECT count(id)
		FROM races
		WHERE finished = 1
	`).Scan(&allRaceCount); err != nil {
		return raceHistory, err
	}

	return raceHistory, nil
}
