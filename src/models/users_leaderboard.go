/*
	These are more functions for querying the "users" table,
	but these functions are only used in "leaderboard.go"
*/

package models

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

type StatsUnseeded struct {
	AdjustedAverage int
	RealAverage     int
	NumRaces        int
	NumForfeits     int
	ForfeitPenalty  int
	LowestTime      int
	LastRace        mysql.NullTime
}

type StatsTrueSkill struct {
	TrueSkill float64
	Mu        float64
	Sigma     float64
	Change    float64
	NumRaces  int
	LastRace  mysql.NullTime
}

func (*Users) GetTrueSkill(userID int, format string) (StatsTrueSkill, error) {
	var stats StatsTrueSkill
	var SQLString string
	if format == "seeded" {
		SQLString = `
			SELECT
				seeded_trueskill,
				seeded_trueskill_mu,
				seeded_trueskill_sigma,
				seeded_trueskill_change,
				seeded_num_races,
				seeded_last_race
			FROM
				users
			WHERE
				id = ?
		`
	} else if format == "unseeded" {
		SQLString = `
			SELECT
				unseeded_trueskill,
				unseeded_trueskill_mu,
				unseeded_trueskill_sigma,
				unseeded_trueskill_change,
				unseeded_num_races,
				unseeded_last_race
			FROM
				users
			WHERE
				id = ?
		`
	} else if format == "diversity" {
		SQLString = `
			SELECT
				diversity_trueskill,
				diversity_trueskill_mu,
				diversity_trueskill_sigma,
				diversity_trueskill_change,
				diversity_num_races,
				diversity_last_race
			FROM
				users
			WHERE
				id = ?
		`
	} else {
		return stats, errors.New("unknown format")
	}

	if err := db.QueryRow(SQLString, userID).Scan(
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

func (*Users) SetStatsSoloUnseeded(userID int, realAverage int, numForfeits int, forfeitPenalty int) error {
	adjustedAverage := realAverage + forfeitPenalty

	// 1800000 is 30 minutes (1000 * 60 * 30)
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE users
		SET
			unseeded_solo_adjusted_average = ?,
			unseeded_solo_real_average = ?,
			unseeded_solo_num_races = (
				SELECT COUNT(race_participants.id)
				FROM race_participants
					JOIN races ON race_participants.race_id = races.id
				WHERE race_participants.user_id = ?
					AND races.ranked = 1
					AND races.format = "unseeded"
			),
			unseeded_solo_num_forfeits = ?,
			unseeded_solo_forfeit_penalty = ?,
			unseeded_solo_lowest_time = (
				SELECT IFNULL(MIN(race_participants.run_time), 1800000)
				FROM race_participants
					JOIN races ON race_participants.race_id = races.id
				WHERE race_participants.user_id = ?
					AND race_participants.place > 0
					AND races.ranked = 1
					AND races.format = "unseeded"
			),
			unseeded_solo_last_race = NOW()
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
		userID,
		userID,
	); err != nil {
		return err
	}

	return nil
}

func (*Users) SetTrueSkill(userID int, stats StatsTrueSkill, format string) error {
	var SQLString string
	if format == "seeded" {
		SQLString = `
			UPDATE users
			SET
				seeded_trueskill = ?,
				seeded_trueskill_mu = ?,
				seeded_trueskill_sigma = ?,
				seeded_trueskill_change = ?,
				seeded_num_races = ?,
				seeded_last_race = NOW()
			WHERE id = ?
		`
	} else if format == "unseeded" {
		SQLString = `
			UPDATE users
			SET
				unseeded_trueskill = ?,
				unseeded_trueskill_mu = ?,
				unseeded_trueskill_sigma = ?,
				unseeded_trueskill_change = ?,
				unseeded_num_races = ?,
				unseeded_last_race = NOW()
			WHERE id = ?
		`
	} else if format == "diversity" {
		SQLString = `
			UPDATE users
			SET
				diversity_trueskill = ?,
				diversity_trueskill_mu = ?,
				diversity_trueskill_sigma = ?,
				diversity_trueskill_change = ?,
				diversity_num_races = ?,
				diversity_last_race = NOW()
			WHERE id = ?
		`
	} else {
		return errors.New("unknown format")
	}

	var stmt *sql.Stmt
	if v, err := db.Prepare(SQLString); err != nil {
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

func (*Users) SetAllDiversityLastRace() error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE users
		SET diversity_last_race = (
			SELECT races.datetime_finished
			FROM race_participants
				JOIN races ON race_participants.race_id = races.id
			WHERE
				user_id = users.id
				AND races.format = "diversity"
			ORDER BY races.datetime_finished DESC
			LIMIT 1
		)
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

func (*Users) ResetStatsDiversity() error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE users
		SET
			diversity_trueskill = 25,
			diversity_trueskill_sigma = 8.333
			diversity_trueskill_change = 0,
			diversity_num_races = 0,
			diversity_last_race = NULL
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
