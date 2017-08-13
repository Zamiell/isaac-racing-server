package models

type MutedUsers struct{}

/*
func (*MutedUsers) Insert(username string, adminResponsible int, reason string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO muted_users (user_id, admin_responsible, reason)
		VALUES ((SELECT id from users WHERE username = ?), ?, ?)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(username, adminResponsible); err != nil {
		return err
	}

	return nil
}

func (*MutedUsers) Delete(username string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		DELETE FROM muted_users
		WHERE user_id = (SELECT id from users WHERE username = ?)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(username); err != nil {
		return err
	}

	return nil
}

func (*MutedUsers) Check(username string) (bool, error) {
	var id int
	if err := db.QueryRow(`
		SELECT id
		FROM muted_users
		WHERE user_id = (SELECT id FROM users WHERE username = ?)
	`, username).Scan(&id); err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
*/
