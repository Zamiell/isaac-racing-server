package models

import (
	"database/sql"
)

/*
	These are more functions for querying the "races" table,
	but these functions are only used in "leaderboard.go"
*/

func (*Races) GetAllRacesForLeaderboard(format string) ([]RaceHistory, error) {
	allRaces := make([]RaceHistory, 0)

	var SQLString string
	if format == "ranked_solo" {
		SQLString = `
			SELECT
				id
			FROM
				races
			WHERE
				finished = 1
				AND ranked = 1
				AND solo = 1
				AND datetime_finished > "` + SoloSeasonStartDatetime + `"
				AND datetime_finished < "` + SoloSeasonEndDatetime + `"
			ORDER BY
				id
		`
	} else {
		SQLString = `
			SELECT
				id
			FROM
				races
			WHERE
				format = "` + format + `"
				AND finished = 1
				AND solo = 0
				AND datetime_finished > "` + RepentanceReleasedDatetime + `"
			ORDER BY
				id
		`
	}

	var rows *sql.Rows
	if v, err := db.Query(SQLString); err != nil {
		return allRaces, err
	} else {
		rows = v
	}
	defer rows.Close()

	for rows.Next() {
		racers := make([]RaceHistoryParticipants, 0)

		var race RaceHistory
		if err := rows.Scan(&race.RaceID); err != nil {
			return allRaces, err
		}
		race.RaceParticipants = nil

		var rows2 *sql.Rows
		if v, err := db.Query(`
			SELECT
				users.id,
				users.username,
				race_participants.place
			FROM
				race_participants
			JOIN
				users ON users.id = race_participants.user_id
			WHERE
				race_participants.race_id = ?
		`, race.RaceID); err != nil {
			return allRaces, err
		} else {
			rows2 = v
		}
		defer rows2.Close()

		for rows2.Next() {
			var racer RaceHistoryParticipants
			if err := rows2.Scan(
				&racer.ID,
				&racer.RacerName,
				&racer.RacerPlace,
			); err != nil {
				return allRaces, err
			}
			racers = append(racers, racer)
		}

		if err := rows2.Err(); err != nil {
			return allRaces, err
		}

		race.RaceParticipants = racers
		allRaces = append(allRaces, race)
	}

	if err := rows.Err(); err != nil {
		return allRaces, err
	}

	return allRaces, nil
}

func (*Races) GetRankedSoloRacesForUser(format string, userID int) ([]RaceHistory, error) {
	allRaces := make([]RaceHistory, 0)

	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			races.id
		FROM
			races
		JOIN
			race_participants ON race_participants.race_id = races.id
		JOIN
			users ON users.id = race_participants.user_id
		WHERE
			races.finished = 1
			AND races.ranked = 1
			AND races.solo = 1
			AND races.datetime_finished > "`+SoloSeasonStartDatetime+`"
			AND races.datetime_finished < "`+SoloSeasonEndDatetime+`"
			AND users.id = ?
		ORDER BY
			id
	`, userID); err != nil {
		return allRaces, err
	} else {
		rows = v
	}
	defer rows.Close()

	for rows.Next() {
		racers := make([]RaceHistoryParticipants, 0)

		var race RaceHistory
		if err := rows.Scan(&race.RaceID); err != nil {
			return allRaces, err
		}
		race.RaceParticipants = nil

		var rows2 *sql.Rows
		if v, err := db.Query(`
			SELECT
				users.id,
				users.username,
				race_participants.place
			FROM
				race_participants
			JOIN
				users ON users.id = race_participants.user_id
			WHERE
				race_participants.race_id = ?
		`, race.RaceID); err != nil {
			return allRaces, err
		} else {
			rows2 = v
		}
		defer rows2.Close()

		for rows2.Next() {
			var racer RaceHistoryParticipants
			if err := rows2.Scan(
				&racer.ID,
				&racer.RacerName,
				&racer.RacerPlace,
			); err != nil {
				return allRaces, err
			}
			racers = append(racers, racer)
		}

		if err := rows2.Err(); err != nil {
			return allRaces, err
		}

		race.RaceParticipants = racers
		allRaces = append(allRaces, race)
	}

	if err := rows.Err(); err != nil {
		return allRaces, err
	}

	return allRaces, nil
}

func (*Races) DeleteOldRankedSoloRaces(userID int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		DELETE races
		FROM races
		JOIN race_participants ON race_participants.race_id = races.id
		JOIN users ON users.id = race_participants.user_id
		WHERE
			races.ranked = 1
			AND races.solo = 1
			AND users.id = ?
			AND races.datetime_finished > "` + RepentanceReleasedDatetime + `"
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(userID); err != nil {
		return err
	}

	return nil
}
