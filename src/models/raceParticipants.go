package models

import (
	"database/sql"
)

type RaceParticipants struct{}

// This mirrors the "race_participants" table row
// (it contains a subset of the information in the non-models Racer struct)
type Racer struct {
	ID               int
	DatetimeJoined   int64
	Seed             string
	StartingItem     int /* Determined by seeing if room count is > 0 */
	Place            int
	DatetimeFinished int64
	RunTime          int64 /* in milliseconds */
	Comment          string
}

func (*RaceParticipants) Insert(raceID int, racer *Racer) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO race_participants (
			user_id,
			race_id,
			datetime_joined,
			seed,
			starting_item,
			place,
			datetime_finished,
			run_time,
			comment
		)
		VALUES (
			?,
			?,
			FROM_UNIXTIME(?),
			?,
			?,
			?,
			FROM_UNIXTIME(?),
			?,
			?
		)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		racer.ID,
		raceID,
		racer.DatetimeJoined,
		racer.Seed,
		racer.StartingItem,
		racer.Place,
		racer.DatetimeFinished,
		racer.RunTime,
		racer.Comment,
	); err != nil {
		return err
	}

	return nil
}

// Get a list of the finished races for this user (quit races don't count)
// Used in the "achievements1_8" function
func (*RaceParticipants) GetFinishedRaces(userID int) ([]Race, error) {
	var raceList []Race

	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT races.id, races.format
		FROM race_participants
			JOIN races ON race_participants.race_id = races.id
		WHERE race_participants.user_id = ? AND race_participants.place > 0
		ORDER BY race_participants.datetime_finished
	`, userID); err != nil {
		return raceList, err
	} else {
		rows = v
	}
	defer rows.Close()

	// Iterate over the races
	for rows.Next() {
		var race Race
		if err := rows.Scan(&race.ID, &race.Format); err != nil {
			return raceList, err
		}

		// Append this race to the slice
		raceList = append(raceList, race)
	}

	if err := rows.Err(); err != nil {
		return raceList, err
	}

	return raceList, nil
}

// Used in the "achievements1_8()" function
type UnseededTime struct {
	Place   int // -1 is quit, -2 is disqualified
	RunTime int64
}

// Get a list of the a player's times for ranked solo races
// Used in the "leaderboardUpdateSoloUnseeded()" function
func (*RaceParticipants) GetNRankedSoloTimes(userID int, n int) ([]UnseededTime, error) {
	var timeList []UnseededTime

	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT race_participants.place, race_participants.run_time
		FROM race_participants
			JOIN races ON race_participants.race_id = races.id
		WHERE
			race_participants.user_id = ?
			AND races.finished = 1
			AND races.ranked = 1
			AND races.solo = 1
			AND races.datetime_finished > "`+SoloSeasonStartDatetime+`"
			AND races.datetime_finished < "`+SoloSeasonEndDatetime+`"
		ORDER BY races.datetime_finished DESC
		LIMIT ?
	`, userID, n); err != nil {
		return timeList, err
	} else {
		rows = v
	}
	defer rows.Close()

	// Iterate over the races
	for rows.Next() {
		var race UnseededTime
		if err := rows.Scan(&race.Place, &race.RunTime); err != nil {
			return timeList, err
		}

		// Append this race to the slice
		timeList = append(timeList, race)
	}

	if err := rows.Err(); err != nil {
		return timeList, err
	}

	return timeList, nil
}

/*
// Used in ?
func (*RaceParticipants) SetComment(userID int, raceID int, comment string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE race_participants
		SET comment = ?
		WHERE user_id = ?
			AND race_id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err = stmt.Exec(comment, userID, raceID); err != nil {
		return err
	}

	return nil
}
*/
