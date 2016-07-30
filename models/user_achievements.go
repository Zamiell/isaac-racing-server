package model

/*
 *  Data types
 */

type UserAchievements struct {
	db *Model
}

/*
 *  user_achievements table functions
 */

func (self *UserAchievements) Insert(username string, achievementID int) error {
	// Give the achievement to that user
	stmt, err := db.Prepare("INSERT INTO user_achievements (user_id, achievement_id) VALUES ((SELECT id FROM users WHERE username = ?), ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(username, achievementID)
	if err != nil {
		return err
	}

	return nil
}
