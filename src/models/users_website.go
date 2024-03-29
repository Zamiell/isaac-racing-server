package models

import (
	"database/sql"
	"errors"
)

/*
	These are more functions for querying the "users" table,
	but these functions are only used for the website
*/

/*
	Data structures
*/

// ProfilesRow gets each row for all profiles
type ProfilesRow struct {
	Username            sql.NullString
	DatetimeCreated     sql.NullTime
	StreamURL           sql.NullString
	NumAchievements     sql.NullInt64
	TotalRaces          sql.NullInt64
	ProfileLastRaceID   sql.NullInt64
	ProfileLastRaceDate sql.NullTime
}

// ProfileData has all data for each racer
type ProfileData struct {
	Username          sql.NullString
	DatetimeCreated   sql.NullTime
	DatetimeLastLogin sql.NullTime
	Admin             sql.NullInt64
	Verified          bool
	StatsSeeded       StatsTrueSkill
	StatsUnseeded     StatsTrueSkill
	StatsSoloUnseeded StatsUnseeded
	StatsDiversity    StatsTrueSkill
	TotalRaces        sql.NullInt64
	StreamURL         sql.NullString
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
			ROUND(u.seeded_trueskill_change, 2),
			ROUND(u.seeded_trueskill_sigma, 2),
			u.seeded_num_races,
			u.seeded_last_race,
			ROUND(u.unseeded_trueskill,2),
			ROUND(u.unseeded_trueskill_change, 2),
			ROUND(u.unseeded_trueskill_sigma, 2),
			u.unseeded_num_races,
			u.unseeded_last_race,
			u.ranked_solo_adjusted_average,
			u.ranked_solo_real_average,
			u.ranked_solo_num_races,
			u.ranked_solo_num_forfeits,
			u.ranked_solo_forfeit_penalty,
			u.ranked_solo_lowest_time,
			u.ranked_solo_last_race,
			ROUND(u.diversity_trueskill, 2),
			ROUND(u.diversity_trueskill_change, 2),
			ROUND(u.diversity_trueskill_sigma, 2),
			u.diversity_num_races,
			u.diversity_last_race,
			COUNT(rp.race_id),
			u.stream_url,
			CASE WHEN u.id IN (SELECT user_id FROM banned_users) THEN 1 ELSE 0 END AS BIT
			FROM
				users u
			LEFT JOIN
				race_participants rp
				ON rp.user_id = u.id
			LEFT JOIN
				races r
				ON r.id = rp.race_id
		WHERE
			u.id =?
	`, userID).Scan(
		&profileData.Username,
		&profileData.DatetimeCreated,
		&profileData.DatetimeLastLogin,
		&profileData.Admin,
		&profileData.Verified,
		&profileData.StatsSeeded.TrueSkill,
		&profileData.StatsSeeded.Change,
		&profileData.StatsSeeded.Sigma,
		&profileData.StatsSeeded.NumRaces,
		&profileData.StatsSeeded.LastRace,
		&profileData.StatsUnseeded.TrueSkill,
		&profileData.StatsUnseeded.Change,
		&profileData.StatsUnseeded.Sigma,
		&profileData.StatsUnseeded.NumRaces,
		&profileData.StatsUnseeded.LastRace,
		&profileData.StatsSoloUnseeded.AdjustedAverage,
		&profileData.StatsSoloUnseeded.RealAverage,
		&profileData.StatsSoloUnseeded.NumRaces,
		&profileData.StatsSoloUnseeded.NumForfeits,
		&profileData.StatsSoloUnseeded.ForfeitPenalty,
		&profileData.StatsSoloUnseeded.LowestTime,
		&profileData.StatsSoloUnseeded.LastRace,
		&profileData.StatsDiversity.TrueSkill,
		&profileData.StatsDiversity.Change,
		&profileData.StatsDiversity.Sigma,
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
	returnTime := int(totalTime.Int64)
	return returnTime, nil
}

// GetUserProfiles gets players data to populate the profiles page
func (*Users) GetUserProfiles(currentPage int, usersPerPage int) ([]ProfilesRow, int, error) {
	profiles := make([]ProfilesRow, 0)

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
	`, usersPerPage, usersOffset); errors.Is(err, sql.ErrNoRows) {
		return profiles, 0, nil
	} else if err != nil {
		return profiles, 0, err
	} else {
		rows = v
	}
	defer rows.Close()

	// Iterate over the user profile results
	for rows.Next() {
		var row ProfilesRow
		if err := rows.Scan(
			&row.Username,
			&row.DatetimeCreated,
			&row.StreamURL,
			&row.NumAchievements,
			&row.TotalRaces,
			&row.ProfileLastRaceID,
			&row.ProfileLastRaceDate,
		); err != nil {
			return profiles, 0, err
		}

		profiles = append(profiles, row)
	}

	if err := rows.Err(); err != nil {
		return profiles, 0, err
	}

	// Find total amount of users
	var allProfilesCount int
	if err := db.QueryRow(`
		SELECT count(id)
		FROM users
		WHERE steam_id > 0
	`).Scan(&allProfilesCount); err != nil {
		return profiles, 0, err
	}

	return profiles, allProfilesCount, nil
}
