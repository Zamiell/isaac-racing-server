package models

/*
 *  Data types
 */

type Seeds struct{}

/*
 *  seeds table functions
 */

func (*Seeds) Get() (string, error) {
	// Get the first seed from the database
	var id int
	var seed string
	err := db.QueryRow("SELECT id, seed FROM seeds LIMIT 1 ORDER BY id").Scan(&id, &seed)
	if err != nil {
		return "", err
	}

	// Delete the seed from the database
	stmt, err := db.Prepare("DELETE FROM seeds WHERE id = ?")
	if err != nil {
		return "", err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return "", err
	}

	return seed, nil
}
