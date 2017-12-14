package models

import (
	"database/sql"
	"errors"
)

/*
	These are more functions for querying the "races" table,
	but these functions are only used in "leaderboard.go"
*/

func (*Races) GetAllRaces(format string) ([]RaceHistory, error) {
	var SQLString string
	if format == "seeded" {
		SQLString = `
			SELECT
				id
			FROM
				races
			WHERE
				format = "seeded"
				AND finished = 1
				AND solo = 0
			ORDER BY
				id
		`
	} else if format == "unseeded" {
		SQLString = `
			SELECT
				id
			FROM
				races
			WHERE
				format = "unseeded"
				AND finished = 1
				AND solo = 0
			ORDER BY
				id
		`
	} else if format == "diversity" {
		SQLString = `
			SELECT
				id
			FROM
				races
			WHERE
				format = "diversity"
				AND finished = 1
				AND solo = 0
			ORDER BY
				id
		`
	} else {
		return nil, errors.New("unknown format")
	}

	var rows *sql.Rows
	if v, err := db.Query(SQLString); err != nil {
		return nil, err
	} else {
		rows = v
	}
	defer rows.Close()

	allRaces := make([]RaceHistory, 0)
	for rows.Next() {
		var race RaceHistory
		if err := rows.Scan(
			&race.RaceID,
		); err != nil {
			return nil, err
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
			return nil, err
		} else {
			rows2 = v
		}
		defer rows2.Close()

		racers := make([]RaceHistoryParticipants, 0)
		for rows2.Next() {
			var racer RaceHistoryParticipants
			if err := rows2.Scan(
				&racer.ID,
				&racer.RacerName,
				&racer.RacerPlace,
			); err != nil {
				return nil, err
			}
			racers = append(racers, racer)
		}
		race.RaceParticipants = racers
		allRaces = append(allRaces, race)
	}

	return allRaces, nil
}
