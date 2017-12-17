package models

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

type LeaderboardRowSeeded struct {
	Name                 string
	SeededTrueSkill      float64
	SeededTrueSkillDelta float64
	SeededNumRaces       sql.NullInt64
	SeededLowestTime     sql.NullInt64
	SeededLastRace       mysql.NullTime
	SeededLastRaceID     int
	Verified             int
	StreamURL            string
}

type LeaderboardRowSeededSolo struct {
	Name                     string
	SeededSoloTrueSkill      float64
	SeededSoloTrueSkillDelta float64
	SeededSoloNumRaces       sql.NullInt64
	SeededSoloLowestTime     sql.NullInt64
	SeededSoloLastRace       mysql.NullTime
	SeededSoloLastRaceID     int
	Verified                 int
	StreamURL                string
}

type LeaderboardRowUnseeded struct {
	Name                   string
	UnseededTrueSkill      float64
	UnseededTrueSkillDelta float64
	UnseededNumRaces       sql.NullInt64
	UnseededLowestTime     sql.NullInt64
	UnseededLastRace       mysql.NullTime
	UnseededLastRaceID     int
	Verified               int
	StreamURL              string
}

type LeaderboardRowUnseededSolo struct {
	Name            string
	AdjustedAverage int
	RealAverage     int
	NumRaces        int
	NumForfeits     int
	ForfeitPenalty  int
	LowestTime      int
	LastRace        mysql.NullTime
	LastRaceID      int
	Verified        int
	StreamURL       string
}

type LeaderboardRowDiversity struct {
	Name              string
	DivTrueSkill      float64
	DivTrueSkillDelta float64
	DivNumRaces       sql.NullInt64
	DivLowestTime     sql.NullInt64
	DivLastRace       mysql.NullTime
	DivLastRaceID     int
	Verified          int
	StreamURL         string
}

func (*Users) GetLeaderboardSeeded(racesNeeded int, racesLimit int) ([]LeaderboardRowSeeded, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			u.username,
			ROUND(u.seeded_trueskill, 2),
			ROUND(u.seeded_trueskill_change, 2),
			u.seeded_num_races,
			(SELECT
					MIN(run_time)
				FROM
					race_participants
				LEFT JOIN races
					ON race_participants.race_id = races.id
				WHERE
					place > 0
					AND u.id = user_id
					AND races.format = 'seeded') as r_time,
			u.seeded_last_race,
			MAX(rp.race_id),
			u.verified,
			u.stream_url
		FROM
			users u
			LEFT JOIN
				race_participants rp ON rp.user_id = u.id
			LEFT JOIN
				races r ON r.id = rp.race_id
		WHERE
			seeded_num_races >= ?
				AND r.format = 'seeded'
				AND rp.place > 0
				AND u.id NOT IN (SELECT user_id FROM banned_users)
		GROUP BY u.username
		ORDER BY u.seeded_trueskill DESC
		LIMIT ?
	`, racesNeeded, racesLimit); err != nil {
		return nil, err
	} else {
		rows = v
	}
	defer rows.Close()

	// Iterate over the users
	leaderboard := make([]LeaderboardRowSeeded, 0)
	for rows.Next() {
		var row LeaderboardRowSeeded
		if err := rows.Scan(
			&row.Name,
			&row.SeededTrueSkill,
			&row.SeededTrueSkillDelta,
			&row.SeededNumRaces,
			&row.SeededLowestTime,
			&row.SeededLastRace,
			&row.SeededLastRaceID,
			&row.Verified,
			&row.StreamURL,
		); err != nil {
			return nil, err
		}

		// Append this row to the leaderboard
		leaderboard = append(leaderboard, row)
	}
	return leaderboard, nil
}

func (*Users) GetLeaderboardSeededSolo(racesNeeded int, racesLimit int) ([]LeaderboardRowSeededSolo, error) {
	var rows *sql.Rows
	// Readd -- ROUND(u.seeded_solo_trueskill_change, 2),
	if v, err := db.Query(`
		SELECT
			u.username,
			ROUND(u.seeded_solo_trueskill, 2),
			ROUND(u.seeded_solo_trueskill, 2),
			u.seeded_solo_num_races,
			(SELECT
					MIN(run_time)
				FROM
					race_participants
				LEFT JOIN races
					ON race_participants.race_id = races.id
				WHERE
					place > 0
					AND u.id = user_id
					AND races.format = 'seeded-solo') as r_time,
			u.seeded_solo_last_race,
			MAX(rp.race_id),
			u.verified,
			u.stream_url
		FROM
			users u
			LEFT JOIN
				race_participants rp ON rp.user_id = u.id
			LEFT JOIN
				races r ON r.id = rp.race_id
		WHERE
			seeded_num_races >= ?
				AND r.format = 'seeded-solo'
				AND rp.place > 0
				AND u.id NOT IN (SELECT user_id FROM banned_users)
		GROUP BY u.username
		ORDER BY u.seeded_solo_trueskill DESC
		LIMIT ?
	`, racesNeeded, racesLimit); err != nil {
		return nil, err
	} else {
		rows = v
	}
	defer rows.Close()

	// Iterate over the users
	leaderboard := make([]LeaderboardRowSeededSolo, 0)
	for rows.Next() {
		var row LeaderboardRowSeededSolo
		if err := rows.Scan(
			&row.Name,
			&row.SeededSoloTrueSkill,
			&row.SeededSoloTrueSkillDelta,
			&row.SeededSoloNumRaces,
			&row.SeededSoloLowestTime,
			&row.SeededSoloLastRace,
			&row.SeededSoloLastRaceID,
			&row.Verified,
			&row.StreamURL,
		); err != nil {
			return nil, err
		}

		// Append this row to the leaderboard
		leaderboard = append(leaderboard, row)
	}
	return leaderboard, nil
}

func (*Users) GetLeaderboardUnseeded(racesNeeded int, racesLimit int) ([]LeaderboardRowUnseeded, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			u.username,
			ROUND(u.unseeded_trueskill, 2),
			ROUND(u.unseeded_trueskill_change, 2),
			u.unseeded_num_races,
			(SELECT
					MIN(run_time)
				FROM
					race_participants
				LEFT JOIN races
					ON race_participants.race_id = races.id
				WHERE
					place > 0
					AND u.id = user_id
					AND races.format = 'unseeded') as r_time,
			u.unseeded_last_race,
			MAX(rp.race_id),
			u.verified,
			u.stream_url
		FROM
			users u
			LEFT JOIN
				race_participants rp ON rp.user_id = u.id
			LEFT JOIN
				races r ON r.id = rp.race_id
		WHERE
			unseeded_num_races >= ?
				AND r.format = 'unseeded'
				AND rp.place > 0
				AND u.id NOT IN (SELECT user_id FROM banned_users)
		GROUP BY u.username
		ORDER BY u.unseeded_trueskill DESC
		LIMIT ?
	`, racesNeeded, racesLimit); err != nil {
		return nil, err
	} else {
		rows = v
	}
	defer rows.Close()

	// Iterate over the users
	leaderboard := make([]LeaderboardRowUnseeded, 0)
	for rows.Next() {
		var row LeaderboardRowUnseeded
		if err := rows.Scan(
			&row.Name,
			&row.UnseededTrueSkill,
			&row.UnseededTrueSkillDelta,
			&row.UnseededNumRaces,
			&row.UnseededLowestTime,
			&row.UnseededLastRace,
			&row.UnseededLastRaceID,
			&row.Verified,
			&row.StreamURL,
		); err != nil {
			return nil, err
		}

		// Append this row to the leaderboard
		leaderboard = append(leaderboard, row)
	}
	return leaderboard, nil
}

// Make a leaderboard for the unseeded solo format based on all of the users
func (*Users) GetLeaderboardUnseededSolo(racesNeeded int, racesLimit int) ([]LeaderboardRowUnseededSolo, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			u.username,
			u.unseeded_solo_adjusted_average,
			u.unseeded_solo_real_average,
			u.unseeded_solo_num_races,
			u.unseeded_solo_num_forfeits,
			u.unseeded_solo_forfeit_penalty,
			u.unseeded_solo_lowest_time,
			u.unseeded_solo_last_race,
			MAX(rp.race_id),
			u.verified,
			u.stream_url
		FROM
			users u
			LEFT JOIN race_participants rp
				ON rp.user_id = u.id
			LEFT JOIN races r
				ON r.id = rp.race_id
		WHERE
			u.unseeded_solo_num_races >= ?
			AND u.id NOT IN (SELECT user_id FROM banned_users)
		GROUP BY
			u.username
		ORDER BY
			unseeded_solo_adjusted_average ASC
		LIMIT ?
	`, racesNeeded, racesLimit); err != nil {
		return nil, err
	} else {
		rows = v
	}
	defer rows.Close()

	// Iterate over the users
	leaderboard := make([]LeaderboardRowUnseededSolo, 0)
	for rows.Next() {
		var row LeaderboardRowUnseededSolo
		if err := rows.Scan(
			&row.Name,
			&row.AdjustedAverage,
			&row.RealAverage,
			&row.NumRaces,
			&row.NumForfeits,
			&row.ForfeitPenalty,
			&row.LowestTime,
			&row.LastRace,
			&row.LastRaceID,
			&row.Verified,
			&row.StreamURL,
		); err != nil {
			return nil, err
		}

		// Append this row to the leaderboard
		leaderboard = append(leaderboard, row)
	}
	return leaderboard, nil
}

func (*Users) GetLeaderboardDiversity(racesNeeded int, racesLimit int) ([]LeaderboardRowDiversity, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			u.username,
			ROUND(u.diversity_trueskill, 2),
			ROUND(u.diversity_trueskill_change, 2),
			u.diversity_num_races,
			(SELECT
					MIN(run_time)
				FROM
					race_participants
				LEFT JOIN races
					ON race_participants.race_id = races.id
				WHERE
					place > 0
					AND u.id = user_id
					AND races.format = 'diversity') as r_time,
			u.diversity_last_race,
			MAX(rp.race_id),
			u.verified,
			u.stream_url
		FROM
			users u
			LEFT JOIN
				race_participants rp ON rp.user_id = u.id
			LEFT JOIN
				races r ON r.id = rp.race_id
		WHERE
			diversity_num_races >= ?
				AND r.format = 'diversity'
				AND rp.place > 0
				AND u.id NOT IN (SELECT user_id FROM banned_users)
		GROUP BY u.username
		ORDER BY u.diversity_trueskill DESC
		LIMIT ?
	`, racesNeeded, racesLimit); err != nil {
		return nil, err
	} else {
		rows = v
	}
	defer rows.Close()

	// Iterate over the users
	leaderboard := make([]LeaderboardRowDiversity, 0)
	for rows.Next() {
		var row LeaderboardRowDiversity
		if err := rows.Scan(
			&row.Name,
			&row.DivTrueSkill,
			&row.DivTrueSkillDelta,
			&row.DivNumRaces,
			&row.DivLowestTime,
			&row.DivLastRace,
			&row.DivLastRaceID,
			&row.Verified,
			&row.StreamURL,
		); err != nil {
			return nil, err
		}

		// Append this row to the leaderboard
		leaderboard = append(leaderboard, row)
	}
	return leaderboard, nil
}
