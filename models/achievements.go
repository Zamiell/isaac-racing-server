package model

/*
 *  Data types
 */

type Achievements struct{}

/*
 *  achievements table functions
 */

func (*Achievements) Insert(id int, name string, description string) error {
	// Add the achievement
	stmt, err := db.Prepare("INSERT INTO achievements (id, name, description) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id, name, description)
	if err != nil {
		return err
	}

	return nil
}

func (*Achievements) DeleteAll() error {
	// Delete every row in the database
	_, err := db.Exec("DELETE FROM achievements")
	if err != nil {
		return err
	}

	// Vacuum, as recommended by: http://www.tutorialspoint.com/sqlite/sqlite_truncate_table.htm
	_, err = db.Exec("VACUUM")
	if err != nil {
		return err
	}

	return nil
}
