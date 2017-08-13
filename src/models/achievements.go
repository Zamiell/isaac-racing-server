package models

import (
	"database/sql"
)

type Achievements struct{}

func (*Achievements) Insert(id int, name string, description string) error {
	var stmt *sql.Stmt
	if v, err := db.Prepare(`
		INSERT INTO achievements (id, name, description)
		VALUES (?, ?, ?)
	`); err != nil {
		return err
	} else {
		stmt = v
	}
	defer stmt.Close()

	if _, err := stmt.Exec(id, name, description); err != nil {
		return err
	}

	return nil
}

// Delete every row in this table
func (*Achievements) DeleteAll() error {
	if _, err := db.Exec("DELETE FROM achievements"); err != nil {
		return err
	}

	return nil
}
