package model

/*
 *  Data types
 */

type Seeds struct {
	db *Model
}

/*
 *  seeds table functions
 */

func (self *Seeds) Get() (string, error) {
	// Get the first seed from the database
	var id int
	var seed string
	err := db.QueryRow("SELECT id, seed FROM seeds LIMIT 1").Scan(&id, &seed)
	if err != nil {
		log.Error("Database error:", err)
		return "", err
	}

	// Delete the seed from the database
	stmt, err := db.Prepare("DELETE FROM seeds WHERE id = ?")
	if err != nil {
		log.Error("Database error:", err)
		return "", err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		log.Error("Database error:", err)
		return "", err
	}

	return seed, nil
}
