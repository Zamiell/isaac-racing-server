package models

import "database/sql"

type UserAchievements struct{}

func (*UserAchievements) Insert(userID int, achievementID int) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO user_achievements (user_id, achievement_id)
		VALUES (?, ?)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(userID, achievementID); err != nil {
		return err
	}

	return nil
}

func (*UserAchievements) GetAll(userID int) ([]int, error) {
	achievementList := make([]int, 0)

	var rows *sql.Rows
	if v, err := db.Query(`
 		SELECT achievement_id
 		FROM user_achievements
 		WHERE user_id = ?
 		ORDER BY datetime_achieved
 	`, userID); err != nil {
		return achievementList, err
	} else {
		rows = v
	}
	defer rows.Close()

	for rows.Next() {
		var achievementID int
		if err := rows.Scan(&achievementID); err != nil {
			return achievementList, err
		}

		achievementList = append(achievementList, achievementID)
	}

	if err := rows.Err(); err != nil {
		return achievementList, err
	}

	return achievementList, nil
}
