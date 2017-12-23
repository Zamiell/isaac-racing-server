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
	ID                    sql.NullInt64
	RacerName             sql.NullString
	RacerPlace            sql.NullInt64
	RacerTime             sql.NullString
	RacerStartingItem     sql.NullInt64
	RacerStartingItemName string
	RacerStartingBuild    sql.NullInt64
	RacerComment          sql.NullString
}

// GetRacesHistory gets all data for all races
func (*Races) GetRacesHistory(currentPage int, racesPerPage int, raceOffset int) ([]RaceHistory, int, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			id,
			solo,
			ranked,
			format,
			player_type,
			goal,
			datetime_created,
			datetime_finished
		FROM
			races
		WHERE
			finished = 1
			AND solo = 0
		GROUP BY
			id
		ORDER BY
			datetime_created DESC
		LIMIT
			?
		OFFSET
			?
	`, racesPerPage, raceOffset); err != nil {
		return nil, 0, err
	} else {
		rows = v
	}
	defer rows.Close()

	raceHistory := make([]RaceHistory, 0)
	for rows.Next() {
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
			return nil, 0, err
		}
		race.RaceParticipants = nil

		var rows2 *sql.Rows
		if v, err := db.Query(`
			SELECT
				u.username,
				rp.place,
				CONCAT(LPAD(FLOOR(rp.run_time/1000/60),2,0), ":", LPAD(FLOOR(rp.run_time/1000%60),2,0)),
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
			return nil, 0, err
		} else {
			rows2 = v
		}
		defer rows2.Close()

		raceRacers := make([]RaceHistoryParticipants, 0)
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
				return nil, 0, err
			}
			raceRacers = append(raceRacers, racer)
		}
		race.RaceParticipants = raceRacers
		raceHistory = append(raceHistory, race)
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
		return nil, 0, err
	}

	return raceHistory, allRaceCount, nil
}

// GetRaceHistory gets race history for a single race
func (*Races) GetRaceHistory(raceID int) ([]RaceHistory, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			id,
			ranked,
			format,
			player_type,
			goal,
			datetime_created,
			datetime_finished
		FROM
			races
		WHERE
			id = ?
		GROUP BY
			id
		ORDER BY
			datetime_created DESC
	`, raceID); err != nil {
		return nil, err
	} else {
		rows = v
	}
	defer rows.Close()

	raceHistory := make([]RaceHistory, 0)
	for rows.Next() {
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
			return nil, err
		}
		race.RaceParticipants = nil

		var rows2 *sql.Rows
		if v, err := db.Query(`
			SELECT
				u.username,
				rp.place,
				CONCAT(LPAD(FLOOR(rp.run_time/1000/60),2,0), ":", LPAD(FLOOR(rp.run_time/1000%60),2,0)),
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
			return nil, err
		} else {
			rows2 = v
		}
		defer rows2.Close()

		raceRacers := make([]RaceHistoryParticipants, 0)
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
				return nil, err
			}
			raceRacers = append(raceRacers, racer)
		}
		race.RaceParticipants = raceRacers
		raceHistory = append(raceHistory, race)
	}
	return raceHistory, nil
}

// GetRaceProfileHistory gets the race data for the profile page
func (*Races) GetRankedRaceProfileHistory(user string, racesPerPage int) ([]RaceHistory, error) {
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
			AND u.username = ?
		GROUP BY
			id
		ORDER BY
			datetime_created DESC
		LIMIT
			?
	`, user, racesPerPage); err != nil {
		return nil, err
	} else {
		rows = v
	}
	defer rows.Close()

	raceHistory := make([]RaceHistory, 0)
	for rows.Next() {
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
			return nil, err
		}
		race.RaceParticipants = nil

		var rows2 *sql.Rows
		if v, err := db.Query(`
			SELECT
				u.username,
				rp.place,
				CONCAT(LPAD(FLOOR(rp.run_time/1000/60),2,0), ":", LPAD(FLOOR(rp.run_time/1000%60),2,0)),
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
			return nil, err
		} else {
			rows2 = v
		}
		defer rows2.Close()

		raceRacers := make([]RaceHistoryParticipants, 0)
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
				return nil, err
			}
			raceRacers = append(raceRacers, racer)
		}
		race.RaceParticipants = raceRacers
		raceHistory = append(raceHistory, race)
	}

	var allRaceCount int
	if err := db.QueryRow(`
		SELECT count(id)
		FROM races
		WHERE finished = 1
	`).Scan(&allRaceCount); err != nil {
		return nil, err
	}

	return raceHistory, nil
}

func (*Races) GetAllRaceProfileHistory(user string, racesPerPage int) ([]RaceHistory, error) {
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
		return nil, err
	} else {
		rows = v
	}
	defer rows.Close()

	raceHistory := make([]RaceHistory, 0)
	for rows.Next() {
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
			return nil, err
		}
		race.RaceParticipants = nil

		var rows2 *sql.Rows
		if v, err := db.Query(`
			SELECT
				u.username,
				rp.place,
				CONCAT(LPAD(FLOOR(rp.run_time/1000/60),2,0), ":", LPAD(FLOOR(rp.run_time/1000%60),2,0)),
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
			return nil, err
		} else {
			rows2 = v
		}
		defer rows2.Close()

		raceRacers := make([]RaceHistoryParticipants, 0)
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
				return nil, err
			}
			raceRacers = append(raceRacers, racer)
		}
		race.RaceParticipants = raceRacers
		raceHistory = append(raceHistory, race)
	}

	var allRaceCount int
	if err := db.QueryRow(`
		SELECT count(id)
		FROM races
		WHERE finished = 1
	`).Scan(&allRaceCount); err != nil {
		return nil, err
	}

	return raceHistory, nil
}
