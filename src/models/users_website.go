package models

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

/*
	These are more functions for querying the "users" table,
	but these functions are only used for the website
*/

/*
	Data structures
*/

type StatsSeeded struct {
	TrueSkill float32
	Sigma     float32
	NumRaces  int
	LastRace  mysql.NullTime
}

type StatsUnseeded struct {
	AdjustedAverage int
	RealAverage     int
	NumRaces        int
	NumForfeits     int
	ForfeitPenalty  int
	LowestTime      int
	LastRace        mysql.NullTime
}

type StatsDiversity struct {
	TrueSkill    float64
	Sigma        float64
	Change       float64
	NumRaces     int
	LastRace     mysql.NullTime
	NewTrueSkill float64 // Only used when doing new TrueSkill calculation
}

// ProfilesRow gets each row for all profiles
type ProfilesRow struct {
	Username            string
	DatetimeCreated     time.Time
	StreamURL           string
	NumAchievements     int
	TotalRaces          int
	ProfileLastRaceId   sql.NullInt64
	ProfileLastRaceDate mysql.NullTime
}

// ProfileData has all data for each racer
type ProfileData struct {
	Username          string
	DatetimeCreated   time.Time
	DatetimeLastLogin time.Time
	Admin             int
	Verified          bool
	StatsSeeded       StatsSeeded
	StatsUnseeded     StatsUnseeded
	StatsDiversity    StatsDiversity
	TotalRaces        sql.NullInt64
	StreamURL         string
	Banned            bool
}

/*
type LeaderboardRowMostPlayed struct {
	Name     string
	Total    int
	Verified int
}
*/

/*
	Functions
*/

/*
func (*Users) GetStatsSeeded(username string) (StatsSeeded, error) {
	var stats StatsSeeded
	if err := db.QueryRow(`
		SELECT
			seeded_trueskill,
			seeded_trueskill_sigma,
			seeded_num_races,
			seeded_last_race
		FROM
			users
		WHERE
			username = ?
	`, username).Scan(
		&stats.ELO,
		&stats.NumSeededRaces,
		&stats.LastSeededRace,
	); err != nil {
		return stats, err
	} else {
		return stats, nil
	}
}

func (*Users) GetStatsUnseeded(username string) (StatsUnseeded, error) {
	var stats StatsUnseeded
	if err := db.QueryRow(`
		SELECT
			unseeded_adjusted_average,
			unseeded_real_average,
			num_unseeded_races,
			num_forfeits,
			forfeit_penalty,
			lowest_unseeded_time,
			last_unseeded_race
		FROM
			users
		WHERE
			username = ?
	`, username).Scan(
		&stats.UnseededAdjustedAverage,
		&stats.UnseededRealAverage,
		&stats.NumUnseededRaces,
		&stats.NumForfeits,
		&stats.ForfeitPenalty,
		&stats.LowestUnseededTime,
		&stats.LastUnseededRace,
	); err != nil {
		return stats, err
	} else {
		return stats, nil
	}
}
*/

// GetProfileData gets player data to populate the player's profile page
func (*Users) GetProfileData(userID int) (ProfileData, error) {
	var profileData ProfileData
	if err := db.QueryRow(`
		SELECT
			u.username,
			u.datetime_created,
			u.datetime_last_login,
			u.admin,
			u.verified,
			ROUND(u.seeded_trueskill, 2),
			ROUND(u.seeded_trueskill_sigma, 2),
			u.seeded_num_races,
			u.seeded_last_race,
			u.unseeded_adjusted_average,
			u.unseeded_real_average,
			u.unseeded_num_races,
			u.unseeded_num_forfeits,
			u.unseeded_forfeit_penalty,
			u.unseeded_lowest_time,
			u.unseeded_last_race,
			ROUND(u.diversity_trueskill, 2),
			ROUND(u.diversity_trueskill_sigma, 2),
			ROUND(u.diversity_trueskill_change, 2),
			u.diversity_num_races,
			u.diversity_last_race,
			COUNT(rp.race_id),
			u.stream_url,
			CASE WHEN u.id IN (SELECT user_id FROM banned_users) THEN 1 ELSE 0 END

			FROM
				users u
			LEFT JOIN
				race_participants rp
				ON rp.user_id = u.id
			LEFT JOIN
				races r
				ON r.id = rp.race_id
		WHERE
			u.id = ?
	`, userID).Scan(
		&profileData.Username,
		&profileData.DatetimeCreated,
		&profileData.DatetimeLastLogin,
		&profileData.Admin,
		&profileData.Verified,
		&profileData.StatsSeeded.TrueSkill,
		&profileData.StatsSeeded.Sigma,
		&profileData.StatsSeeded.NumRaces,
		&profileData.StatsSeeded.LastRace,
		&profileData.StatsUnseeded.AdjustedAverage,
		&profileData.StatsUnseeded.RealAverage,
		&profileData.StatsUnseeded.NumRaces,
		&profileData.StatsUnseeded.NumForfeits,
		&profileData.StatsUnseeded.ForfeitPenalty,
		&profileData.StatsUnseeded.LowestTime,
		&profileData.StatsUnseeded.LastRace,
		&profileData.StatsDiversity.TrueSkill,
		&profileData.StatsDiversity.Sigma,
		&profileData.StatsDiversity.Change,
		&profileData.StatsDiversity.NumRaces,
		&profileData.StatsDiversity.LastRace,
		&profileData.TotalRaces,
		&profileData.StreamURL,
		&profileData.Banned,
	); err != nil {
		return profileData, err
	}

	return profileData, nil
}

// GetTotalTime gets the total run_time of all races completed regardless of outcome
func (*Users) GetTotalTime(username string) (int, error) {
	var totalTime sql.NullInt64
	if err := db.QueryRow(`
	SELECT
		SUM(rp.run_time)
	FROM
		users u
	LEFT JOIN
		race_participants rp
		ON rp.user_id = u.id
	LEFT JOIN
		races r
		ON r.id = rp.race_id
	WHERE
		r.finished = 1
		AND u.username = ?
`, username).Scan(&totalTime); err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	var returnTime int
	returnTime = int(totalTime.Int64)
	return returnTime, nil

}

// GetUserProfiles gets players data to populate the profiles page
func (*Users) GetUserProfiles(currentPage int, usersPerPage int) ([]ProfilesRow, int, error) {
	usersOffset := (currentPage - 1) * usersPerPage
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			u.username,
			u.datetime_created,
			u.stream_url,
			count(ua.achievement_id),
			(
				SELECT COUNT(id)
				FROM race_participants
				WHERE user_id = u.id
			) AS num_total_race,
			MAX(rp.race_id),
			r.datetime_started
		FROM
			users u
		LEFT JOIN
			user_achievements ua
			ON u.id = ua.user_id
		LEFT JOIN race_participants rp
			ON rp.user_id = u.id
		LEFT JOIN races r
			ON r.id = rp.race_id
		WHERE
			u.steam_id > 0
		GROUP BY
			u.username
		ORDER BY
			u.username ASC
		LIMIT
			?
		OFFSET
			?
	`, usersPerPage, usersOffset); err == sql.ErrNoRows {
		return nil, 0, nil
	} else if err != nil {
		return nil, 0, err
	} else {
		rows = v
	}
	defer rows.Close()

	// Iterate over the user profile results
	profiles := make([]ProfilesRow, 0)
	for rows.Next() {
		var row ProfilesRow
		if err := rows.Scan(
			&row.Username,
			&row.DatetimeCreated,
			&row.StreamURL,
			&row.NumAchievements,
			&row.TotalRaces,
			&row.ProfileLastRaceId,
			&row.ProfileLastRaceDate,
		); err != nil {
			return nil, 0, err
		}

		profiles = append(profiles, row)
	}

	// Find total amount of users
	var allProfilesCount int
	if err := db.QueryRow(`
		SELECT count(id)
		FROM users
		WHERE steam_id > 0
	`).Scan(&allProfilesCount); err != nil {
		return nil, 0, err
	}

	return profiles, allProfilesCount, nil
}

// Make a leaderboard for the unseeded format based on all of the users
type LeaderboardRowUnseeded struct {
	Name            string
	AdjustedAverage int
	RealAverage     int
	NumRaces        int
	NumForfeits     int
	ForfeitPenalty  int
	LowestTime      int
	LastRace        time.Time
	LastRaceId      int
	Verified        int
	StreamURL       string
}

func (*Users) GetLeaderboardUnseeded(racesNeeded int, racesLimit int) ([]LeaderboardRowUnseeded, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			u.username,
			u.unseeded_adjusted_average,
			u.unseeded_real_average,
			u.unseeded_num_races,
			u.unseeded_num_forfeits,
			u.unseeded_forfeit_penalty,
			u.unseeded_lowest_time,
			u.unseeded_last_race,
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
			u.unseeded_num_races >= ?
			AND u.id NOT IN (SELECT user_id FROM banned_users)
		GROUP BY
			u.username
		ORDER BY
			unseeded_adjusted_average ASC
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
			&row.AdjustedAverage,
			&row.RealAverage,
			&row.NumRaces,
			&row.NumForfeits,
			&row.ForfeitPenalty,
			&row.LowestTime,
			&row.LastRace,
			&row.LastRaceId,
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

// Make a leaderboard for the seeded format based on all of the users
type LeaderboardRowSeeded struct {
	Name      string
	TrueSkill float64
	NumRaces  int
	LastRace  time.Time
	Verified  int
}

func (*Users) GetLeaderboardSeeded(racesNeeded int, racesLimit int) ([]LeaderboardRowSeeded, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			u.username,
			u.seeded_trueskill,
			u.seeded_trueskill_sigma,
			u.seeded_num_races,
			u.seeded_last_race,
			u.verified
		FROM
			users u
		WHERE
			u.seeded_num_races > ?
			AND u.id NOT IN (SELECT user_id FROM banned_users)
		GROUP BY u.username
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
			&row.TrueSkill,
			&row.NumRaces,
			&row.LastRace,
			&row.Verified,
		); err != nil {
			return nil, err
		}

		// Append this row to the leaderboard
		leaderboard = append(leaderboard, row)
	}

	return leaderboard, nil
}

type LeaderboardRowDiversity struct {
	Name              string
	DivTrueSkill      float64
	DivTrueSkillDelta float64
	DivNumRaces       sql.NullInt64
	DivLowestTime     sql.NullInt64
	DivLastRace       time.Time
	DivLastRaceId     int
	Verified          int
	StreamURL         string
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
			&row.DivLastRaceId,
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
