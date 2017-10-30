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
	Username        string
	DatetimeCreated time.Time
	StreamURL       string
	NumAchievements int
	TotalRaces      int
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
	StreamURL         string
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
func (*Users) GetProfileData(username string) (ProfileData, error) {
	var profileData ProfileData
	var rawVerified int
	if err := db.QueryRow(`
		SELECT
			username,
			datetime_created,
			datetime_last_login,
			admin,
			verified,
			seeded_trueskill,
			seeded_trueskill_sigma,
			seeded_num_races,
			seeded_last_race,
			unseeded_adjusted_average,
			unseeded_real_average,
			unseeded_num_races,
			unseeded_num_forfeits,
			unseeded_forfeit_penalty,
			unseeded_lowest_time,
			unseeded_last_race,
			stream_url
		FROM
			users
		WHERE
			steam_id > 0 and
			username = ?
	`, username).Scan(
		&profileData.Username,
		&profileData.DatetimeCreated,
		&profileData.DatetimeLastLogin,
		&profileData.Admin,
		&rawVerified,
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
		&profileData.StreamURL,
	); err == sql.ErrNoRows {
		return profileData, nil
	} else if err != nil {
		return profileData, err
	} else {
		// Convert the int to a bool
		if rawVerified == 1 {
			profileData.Verified = true
		}
		return profileData, nil
	}
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
            (SELECT COUNT(id)
                FROM race_participants
                WHERE user_id = u.id
            ) AS num_total_race			
		FROM
			users u
		LEFT JOIN
			user_achievements ua
			ON
				u.id = ua.user_id
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
	Verified        int
	StreamURL       string
}

func (*Users) GetLeaderboardUnseeded() ([]LeaderboardRowUnseeded, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			username,
			verified,
			unseeded_adjusted_average,
			unseeded_real_average,
			unseeded_num_races,
			unseeded_num_forfeits,
			unseeded_forfeit_penalty,
			unseeded_lowest_time,
			unseeded_last_race,
			stream_url
		FROM
			users
		WHERE
			unseeded_num_races >= 20
		ORDER BY
			unseeded_adjusted_average ASC
	`); err != nil {
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
			&row.Verified,
			&row.AdjustedAverage,
			&row.RealAverage,
			&row.NumRaces,
			&row.NumForfeits,
			&row.ForfeitPenalty,
			&row.LowestTime,
			&row.LastRace,
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

func (*Users) GetLeaderboardSeeded() ([]LeaderboardRowSeeded, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			username,
			seeded_trueskill,
			seeded_trueskill_sigma,
			seeded_num_races,
			seeded_last_race,
			verified,
		FROM
			users
		WHERE
			seeded_num_races > 1
	`); err != nil {
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
	Name                string
	DivTrueSkill        float64
	DivTrueSkillDelta   float64
	DivNumRaces         sql.NullInt64
	DivLowestTime       sql.NullInt64
	DivLastRace         time.Time
	Verified            int
	StreamURL           string
}

func (*Users) GetLeaderboardDiversity() ([]LeaderboardRowDiversity, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT 
		    u.username,
		    u.verified,
		    u.diversity_trueskill,
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
		    u.stream_url
		FROM
		    users u
		        LEFT JOIN
		    race_participants rp ON rp.user_id = u.id
		        LEFT JOIN
		    races r ON r.id = rp.race_id
		WHERE
		    diversity_num_races >= 5
		        AND r.format = 'diversity'
		        AND rp.place > 0
		GROUP BY u.username
		ORDER BY u.diversity_trueskill DESC
            
	`); err != nil {
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
			&row.Verified,
			&row.DivTrueSkill,
			&row.DivTrueSkillDelta,
			&row.DivNumRaces,
			&row.DivLowestTime,
			&row.DivLastRace,
			&row.StreamURL,
		); err != nil {
			return nil, err
		}

		// Append this row to the leaderboard
		leaderboard = append(leaderboard, row)
	}
	return leaderboard, nil
}
