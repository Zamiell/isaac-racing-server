package models

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

/*
	These are functions used for the listing the schedules for tournaments
	on the website
*/

// TournamentRace holds data for a single race
type TournamentRace struct {
	TournamentName  sql.NullString
	TournamentID    sql.NullString
	TournamentType  sql.NullString
	RaceID          sql.NullInt64
	Racer1          sql.NullString
	Racer2          sql.NullString
	RaceState       sql.NullString
	TournamentRound sql.NullInt64
	RaceDateTime    mysql.NullTime
	ChallongeID     sql.NullInt64
	RaceCasterName  sql.NullString
	RaceCasterURL   sql.NullString
}

// GetTournamentRaces gets all data for all races
func (*Tournament) GetTournamentRaces() ([]TournamentRace, error) {
	tournamentRaces := make([]TournamentRace, 0)

	var rows *sql.Rows
	if v, err := db.Query(`
		SELECT
			tr.tournament_name,
			tr.challonge_url,
			tr.id,
			trs1.username,
			trs2.username,
			tr.state,
			tr.bracket_round,
			tr.datetime_scheduled,
			tr.challonge_match_id,
			trs3.username,
			trs3.stream_url
		FROM
			isaac.tournament_races tr
		LEFT JOIN
			isaac.tournament_users trs1 ON trs1.id = tr.racer1
		LEFT JOIN
			isaac.tournament_users trs2 ON trs2.id = tr.racer2
		LEFT JOIN
			isaac.tournament_casts c ON c.race_id = tr.id
		LEFT JOIN
			isaac.tournament_users trs3 ON trs3.id = c.caster
		WHERE
			tr.state not in ('initial', 'completed')
		ORDER BY datetime_scheduled ASC
	`); err != nil {
		return tournamentRaces, err
	} else {
		rows = v
	}
	defer rows.Close()

	for rows.Next() {
		var race TournamentRace
		if err := rows.Scan(
			&race.TournamentName,
			&race.TournamentID,
			&race.RaceID,
			&race.Racer1,
			&race.Racer2,
			&race.RaceState,
			&race.TournamentRound,
			&race.RaceDateTime,
			&race.ChallongeID,
			&race.RaceCasterName,
			&race.RaceCasterURL,
		); errors.Is(err, sql.ErrNoRows) {
			return tournamentRaces, sql.ErrNoRows
		} else if err != nil {
			return tournamentRaces, err
		}
		tournamentRaces = append(tournamentRaces, race)
	}

	if err := rows.Err(); err != nil {
		return tournamentRaces, err
	}

	if len(tournamentRaces) == 0 {
		return tournamentRaces, sql.ErrNoRows
	}

	return tournamentRaces, nil
}
