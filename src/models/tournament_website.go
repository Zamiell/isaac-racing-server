package models

import (
	"database/sql"
	//"github.com/Zamiell/isaac-racing-server/src/log"

	"github.com/go-sql-driver/mysql"
)

/*
	These are functions used for the listing the schedules for tournaments
	on the website
*/

// TournamentRace holds data for a single race
type TournamentRace struct {
	TournamentName sql.NullString
	TournamentType string
	RaceID         sql.NullInt64
	Racer1         sql.NullString
	Racer2         sql.NullString
	RaceState      sql.NullString
	RaceDateTime   mysql.NullTime
	ChallongeID    sql.NullInt64
	RaceCaster     sql.NullString
}

// GetTournamentRaces gets all data for all races
func (*Tournament) GetTournamentRaces() ([]TournamentRace, error) {
	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
				tr.tournament_name,
				tr.id,
				trs1.username,
				trs2.username,
				tr.state,
				tr.datetime_scheduled,
				tr.challonge_id,
				c.username
		FROM
			isaac.tournament_races tr
					LEFT JOIN
							isaac.tournament_racers trs1 ON trs1.id = tr.racer1
					LEFT JOIN
							isaac.tournament_racers trs2 ON trs2.id = tr.racer2
					LEFT JOIN
							isaac.tournament_racers c ON c.id = tr.caster
		WHERE
				tr.state = 'scheduled'
		ORDER BY datetime_scheduled ASC
	`); err != nil {
		return nil, err
	} else {
		rows = v
	}
	defer rows.Close()

	tournamentRaces := make([]TournamentRace, 0)
	for rows.Next() {
		var race TournamentRace
		if err := rows.Scan(
			&race.TournamentName,
			&race.RaceID,
			&race.Racer1,
			&race.Racer2,
			&race.RaceState,
			&race.RaceDateTime,
			&race.ChallongeID,
			&race.RaceCaster,
		); err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		} else if err != nil {
			return nil, err
		}
		tournamentRaces = append(tournamentRaces, race)
	}
	if len(tournamentRaces) == 0 {
		return nil, sql.ErrNoRows
	}
	return tournamentRaces, nil
}
