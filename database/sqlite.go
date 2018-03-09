package database

import (
	"database/sql"

	// SQLite driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/tuommii/jumbo/model"
)

// SQLiteDB implements Database interface
type SQLiteDB struct {
	Connection *sql.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS player(
		id INTEGER NOT NULL,
		name TEXT NOT NULL,
		CONSTRAINT player_PK PRIMARY KEY(id),
		CONSTRAINT player_name_MIN_LENGTH CHECK(length(name) >= 2),
		CONSTRAINT player_name_MAX_LENGTH CHECK(length(name) <= 16),
		CONSTRAINT player_name_UNIQUE UNIQUE(name));

CREATE TABLE IF NOT EXISTS game(
		id INTEGER NOT NULL,
		name TEXT NOT NULL,
		CONSTRAINT game_PK PRIMARY KEY(id),
		CONSTRAINT game_name_MIN_LENGTH CHECK(length(name) >= 2),
		CONSTRAINT game_name_MAX_LENGTH CHECK(length(name) <= 64),
		CONSTRAINT game_name_UNIQUE UNIQUE(name));

CREATE TABLE IF NOT EXISTS match(
		id INTEGER NOT NULL,
		game_name TEXT NOT NULL,
		winner TEXT NOT NULL,
		loser TEXT NOT NULL,
		comment TEXT NOT NULL,
		is_tie BOOLEAN NOT NULL,
		added TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT game_PK PRIMARY KEY(id));
`

// NewSQLiteDB returns connection to SQLite database
func NewSQLiteDB(name string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return &SQLiteDB{Connection: db}, nil
}

/*
**
** #PLAYER
**
 */

// GetPlayers returns all players
func (db *SQLiteDB) GetPlayers() ([]model.Player, error) {
	rows, err := db.Connection.Query("SELECT id, name FROM player")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]model.Player, 0)

	for rows.Next() {
		player := model.Player{}
		err := rows.Scan(&player.ID, &player.Name)
		if err != nil {
			return nil, err
		}

		players = append(players, player)
	}

	return players, nil
}

// CreatePlayer creates new player
func (db *SQLiteDB) CreatePlayer(name string) (int64, error) {
	stmt, err := db.Connection.Prepare("INSERT INTO player(name) VALUES(?)")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name)
	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

// DeletePlayer deletes player
func (db *SQLiteDB) DeletePlayer(name string) (int64, error) {
	stmt, err := db.Connection.Prepare("DELETE FROM player WHERE name = ?")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name)
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}

/*
**
** #GAME
**
 */

// GetGames returns all games
func (db *SQLiteDB) GetGames() ([]model.Game, error) {
	rows, err := db.Connection.Query("SELECT id, name FROM game")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := make([]model.Game, 0)

	for rows.Next() {
		game := model.Game{}
		err := rows.Scan(&game.ID, &game.Name)
		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	return games, nil
}

// CreateGame creates new game
func (db *SQLiteDB) CreateGame(name string) (int64, error) {
	stmt, err := db.Connection.Prepare("INSERT INTO game(name) VALUES(?)")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name)
	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

// DeleteGame deletes player
func (db *SQLiteDB) DeleteGame(name string) (int64, error) {
	stmt, err := db.Connection.Prepare("DELETE FROM game WHERE name = ?")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name)
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}

/*
**
** #MATCH
**
 */

// CreateMatch creates new match
func (db *SQLiteDB) CreateMatch(match model.Match) (int64, error) {
	stmt, err := db.Connection.Prepare("INSERT INTO match(game_name, winner, loser, comment, is_tie) VALUES(?,?,?,?,?)")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(match.GameName, match.Winner, match.Loser, match.Comment, match.IsTie)
	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

// GetMatches returns all matches
func (db *SQLiteDB) GetMatches(f model.Filter) ([]model.Match, error) {
	rows, err := db.Connection.Query(f.GetQuery())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]model.Match, 0)

	for rows.Next() {
		match := model.Match{}
		err := rows.Scan(
			&match.ID,
			&match.GameName,
			&match.IsTie,
			&match.Winner,
			&match.Loser,
			&match.Comment,
			&match.Added,
		)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}
	return matches, nil
}

// DeleteMatch deletes match
func (db *SQLiteDB) DeleteMatch(id int) (int64, error) {
	stmt, err := db.Connection.Prepare("DELETE FROM game WHERE id = ?")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}

// StatsFromMatches return games and those stats for each player
func StatsFromMatches(matchRows *sql.Rows) (model.SortedStats, []*model.Match, error) {
	matches := make([]*model.Match, 0)
	players := make(model.StatsMap)

	for matchRows.Next() {
		match := new(model.Match)
		err := matchRows.Scan(
			&match.ID, &match.GameName, &match.IsTie, &match.Winner, &match.Loser,
			&match.Comment, &match.Added,
		)
		if err != nil {
			return nil, nil, err
		}

		// Set stats to 0 for new players
		players.Init(match.Winner, match.Loser)
		// Increase stats
		players.Increase(match.Winner, match.Loser, match.IsTie)

		matches = append(matches, match)
	}

	players.CalculateComputed()

	// return SortedStats
	ss := make(model.SortedStats, 0, len(players))

	for key, p := range players {
		p.Name = key
		ss = append(ss, p)
	}

	return ss, matches, nil
}

// Search returns stats and games based on filters
// func Search(db *sql.DB, f *models.Filter) (models.SortedStats, []*models.Match, error) {
// 	var rows *sql.Rows
// 	var err error

// 	rows, err = db.Query(f.GetQuery())
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	defer rows.Close()

// 	stats, matches, err := models.StatsFromMatches(rows)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return stats, matches, nil
// }
