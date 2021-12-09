/*
	These are more functions for querying the "users" table,
	but these functions are only used in "leaderboard.go"
*/

package models

import (
	"database/sql"
	"errors"
	"strconv"
)

const (
	SoloSeason1StartDatetime = "2017-10-17 23:00:00"
	SoloSeason1EndDatetime   = "2018-03-17 00:00:00"
	SoloSeason2StartDatetime = "2018-03-18 23:00:00"
	SoloSeason2EndDatetime   = "2018-10-26 00:00:00"

	// This is not actually when the Repentance DLC was released,
	// but rather when the Repentance version of Racing+ was released
	RepentanceReleasedDatetime = "2021-05-21 00:00:00"

	SoloSeason3StartDatetime = "2021-12-03 00:00:00"
	SoloSeason3EndDatetime   = "2030-00-00 00:00:00"

	SoloSeasonStartDatetime = SoloSeason3StartDatetime
	SoloSeasonEndDatetime   = SoloSeason3EndDatetime
)

type StatsUnseeded struct {
	AdjustedAverage int
	RealAverage     int
	NumRaces        int
	NumForfeits     int
	ForfeitPenalty  int
	LowestTime      int
	LastRace        sql.NullTime
}

type StatsTrueSkill struct {
	TrueSkill float64
	Mu        float64
	Sigma     float64
	Change    float64
	NumRaces  int
	LastRace  sql.NullTime
}

func (*Users) GetTrueSkill(userID int, format string) (StatsTrueSkill, error) {
	var stats StatsTrueSkill
	if err := db.QueryRow(`
		SELECT
			`+format+`_trueskill,
			`+format+`_trueskill_mu,
			`+format+`_trueskill_sigma,
			`+format+`_trueskill_change,
			`+format+`_num_races,
			`+format+`_last_race
		FROM
			users
		WHERE
			id = ?
	`, userID).Scan(
		&stats.TrueSkill,
		&stats.Mu,
		&stats.Sigma,
		&stats.Change,
		&stats.NumRaces,
		&stats.LastRace,
	); err != nil {
		return stats, err
	}

	return stats, nil
}

func (*Users) SetTrueSkill(userID int, stats StatsTrueSkill, format string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE users
		SET
			` + format + `_trueskill = ?,
			` + format + `_trueskill_mu = ?,
			` + format + `_trueskill_sigma = ?,
			` + format + `_trueskill_change = ?,
			` + format + `_num_races = ?,
			` + format + `_last_race = NOW()
		WHERE id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		stats.TrueSkill,
		stats.Mu,
		stats.Sigma,
		stats.Change,
		stats.NumRaces,
		userID,
	); err != nil {
		return err
	}

	return nil
}

// Only used in the "leaderboardRecalculate" functions
func (*Users) SetLastRace(format string) error {
	var SQLString string
	if format == "ranked_solo" {
		SQLString = `
			UPDATE users
			SET ` + format + `_last_race = (
				SELECT races.datetime_finished
				FROM race_participants
					JOIN races ON race_participants.race_id = races.id
				WHERE
					user_id = users.id
					AND races.finished = 1
					AND races.ranked = 1
					AND races.solo = 1
					AND races.datetime_finished > "` + SoloSeasonStartDatetime + `"
					AND races.datetime_finished < "` + SoloSeasonEndDatetime + `"
				ORDER BY races.datetime_finished DESC
				LIMIT 1
			)
		`
	} else {
		SQLString = `
			UPDATE users
			SET ` + format + `_last_race = (
				SELECT races.datetime_finished
				FROM race_participants
					JOIN races ON race_participants.race_id = races.id
				WHERE
					user_id = users.id
					AND races.finished = 1
					AND races.format = "` + format + `"
					AND races.solo = 0
					AND races.datetime_finished > "` + RepentanceReleasedDatetime + `"
				ORDER BY races.datetime_finished DESC
				LIMIT 1
			)
		`
	}

	var stmt *sql.Stmt
	if v, err := db.Prepare(SQLString); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(); err != nil {
		return err
	}

	return nil
}

func (*Users) ResetTrueSkill(format string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE users
		SET
			` + format + `_trueskill = 25,
			` + format + `_trueskill_sigma = 8.333,
			` + format + `_trueskill_change = 0,
			` + format + `_num_races = 0,
			` + format + `_last_race = NULL
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(); err != nil {
		return err
	}

	return nil
}

func (*Users) SetStatsRankedSolo(
	userID int,
	realAverage int,
	numForfeits int,
	forfeitPenalty int,
	lowestTime int64,
	startingBuild int,
) error {
	adjustedAverage := realAverage + forfeitPenalty

	// 1800000 is 30 minutes (1000 * 60 * 30)
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE users
		SET
			ranked_solo_adjusted_average = ?,
			ranked_solo_real_average = ?,
			ranked_solo_num_races = (
				SELECT COUNT(race_participants.id)
				FROM race_participants
					JOIN races ON race_participants.race_id = races.id
				WHERE
					race_participants.user_id = ?
					AND races.finished = 1
					AND races.ranked = 1
					AND races.solo = 1
					AND races.datetime_finished > "` + SoloSeasonStartDatetime + `"
					AND races.datetime_finished < "` + SoloSeasonEndDatetime + `"
			),
			ranked_solo_num_forfeits = ?,
			ranked_solo_forfeit_penalty = ?,
			ranked_solo_lowest_time = ?,
			ranked_solo_last_race = NOW(),
			ranked_solo_metadata = ?
		WHERE id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		adjustedAverage,
		realAverage,
		userID,
		numForfeits,
		forfeitPenalty,
		lowestTime,
		startingBuild,
		userID,
	); err != nil {
		return err
	}

	return nil
}

func (*Users) ResetRankedSolo(userID int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE users
		SET
			ranked_solo_adjusted_average = 0,
			ranked_solo_real_average = 0,
			ranked_solo_num_forfeits = 0,
			ranked_solo_forfeit_penalty = 0,
			ranked_solo_lowest_time = 0,
			ranked_solo_num_races = 0,
			ranked_solo_last_race = NULL,
			ranked_solo_metadata = NULL
		WHERE
			id = ?
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

func (*Users) ResetRankedSoloAll() error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE users
		SET
			ranked_solo_adjusted_average = 0,
			ranked_solo_real_average = 0,
			ranked_solo_num_forfeits = 0,
			ranked_solo_forfeit_penalty = 0,
			ranked_solo_lowest_time = 0,
			ranked_solo_num_races = 0,
			ranked_solo_last_race = NULL,
			ranked_solo_metadata = NULL
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(); err != nil {
		return err
	}

	return nil
}

func (*Users) GetRankedSoloMetadata(userID int) (sql.NullInt64, error) {
	var metadata sql.NullInt64
	if err := db.QueryRow(`
		SELECT ranked_solo_metadata
		FROM users
		WHERE id = ?
	`, userID).Scan(&metadata); err == sql.ErrNoRows {
		return sql.NullInt64{}, errors.New("A user with an ID of \"" + strconv.Itoa(userID) + "\" does not exist.")
	} else if err != nil {
		return sql.NullInt64{}, err
	}

	return metadata, nil
}
