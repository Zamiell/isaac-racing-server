package models

import (
	"database/sql"
)

/*
	These are more functions for querying the "users" table,
	but these functions are only used in "leaderboard.go"
*/

func (*Users) GetStatsDiversity(userID int) (StatsDiversity, error) {
	var stats StatsDiversity
	if err := db.QueryRow(`
		SELECT
			diversity_trueskill,
			diversity_trueskill_sigma,
			diversity_trueskill_change,
			diversity_num_races,
			diversity_last_race
		FROM
			users
		WHERE
			id = ?
	`, userID).Scan(
		&stats.TrueSkill,
		&stats.Sigma,
		&stats.Change,
		&stats.NumRaces,
		&stats.LastRace,
	); err != nil {
		return stats, err
	}

	return stats, nil
}

func (*Users) SetStatsUnseeded(userID int, realAverage int, numForfeits int, forfeitPenalty int) error {
	adjustedAverage := realAverage + forfeitPenalty

	// 1800000 is 30 minutes (1000 * 60 * 30)
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE users
		SET
			unseeded_adjusted_average = ?,
			unseeded_real_average = ?,
			unseeded_num_races = (
				SELECT COUNT(race_participants.id)
				FROM race_participants
					JOIN races ON race_participants.race_id = races.id
				WHERE race_participants.user_id = ?
					AND races.ranked = 1
					AND races.format = "unseeded"
			),
			unseeded_num_forfeits = ?,
			unseeded_forfeit_penalty = ?,
			unseeded_lowest_time = (
				SELECT IFNULL(MIN(race_participants.run_time), 1800000)
				FROM race_participants
					JOIN races ON race_participants.race_id = races.id
				WHERE race_participants.user_id = ?
					AND race_participants.place > 0
					AND races.ranked = 1
					AND races.format = "unseeded"
			),
			unseeded_last_race = NOW()
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

func (*Users) SetStatsDiversity(userID int, stats StatsDiversity) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		UPDATE users
		SET
			diversity_trueskill = ?,
			diversity_trueskill_sigma = ?,
			diversity_trueskill_change = ?,
			diversity_num_races = ?,
			diversity_last_race = NOW()
		WHERE id = ?
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		stats.NewTrueSkill,
		stats.Sigma,
		stats.Change,
		stats.NumRaces,
		userID,
	); err != nil {
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
			diversity_trueskill_sigma = 8.333,
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
