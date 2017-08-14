package models

import "database/sql"

/*
	These are more functions for querying the "races" table,
	but these functions are only used for the website
*/

type RaceHistory struct {
	RaceID           int
	RaceDate         int
	RaceType         string
	RaceFormat       string
	RaceChar         string
	RaceGoal         string
	RaceParticipants []RaceHistoryParticipants
}
type RaceHistoryParticipants struct {
	RacerName     string
	RacerPlace    int
	RacerTime     string
	RacerComment  string
	RaceRoomCount int
	RaceItemCount int
}

func (*Races) GetRaceHistory(currentPage int, racesPerPage int) ([]RaceHistory, int, error) {
	raceOffset := currentPage * racesPerPage

	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			id,
			ranked,
			format,
			player_type,
			goal,
			datetime_finished
		FROM
			races
		WHERE
			finished = 1
		GROUP BY
			id
		ORDER BY
			id DESC
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
			&race.RaceDate,
			&race.RaceType,
			&race.RaceFormat,
			&race.RaceChar,
			&race.RaceGoal,
		); err != nil {
			return nil, 0, err
		}
		race.RaceParticipants = nil

		var rows2 *sql.Rows
		if v, err := db.Query(`
			SELECT
			    u.username,
			    rp.place,
			    rp.run_time,
			    rp.comment
			FROM
			    race_participants rp
			LEFT JOIN
			    users u
			    ON u.id = rp.user_id
			WHERE
			    rp.race_id = 33
			ORDER BY
			    CASE WHEN rp.place = -1 THEN 1 ELSE 0 END,
			    rp.place;
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
				&racer.RacerComment,
				&racer.RaceRoomCount,
				&racer.RaceItemCount,
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
		SELECT count(id)
		FROM races
		WHERE finished = 1
	`).Scan(&allRaceCount); err != nil {
		return nil, 0, err
	}

	return raceHistory, allRaceCount, nil
}
