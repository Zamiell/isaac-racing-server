package models

/*
 *  Data types
 */

type UserAchievements struct{}

/*
 *  user_achievements table functions
 */

func (*UserAchievements) GetAll(username string) ([]int, error) {
	// Get all the achivements for this user
	rows, err := db.Query(`
 		SELECT achievement_id
 		FROM user_achievements
 		WHERE user_id = (SELECT id FROM users WHERE username = ?)
 		ORDER BY datetime
 	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Make the list
	var achievementList []int
	for rows.Next() {
		var achievement int
		err := rows.Scan(&achievement)
		if err != nil {
			return nil, err
		}

		// Add it to the list
		achievementList = append(achievementList, achievement)
	}

	return achievementList, nil
}

func (*UserAchievements) Insert(username string, achievementID int) error {
	// Give the achievement to that user
	stmt, err := db.Prepare(`
		INSERT INTO user_achievements (user_id, achievement_id)
		VALUES ((SELECT id FROM users WHERE username = ?), ?)
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(username, achievementID)
	if err != nil {
		return err
	}

	return nil
}
