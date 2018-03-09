package model

import "strconv"

// Filter query
type Filter struct {
	GameName   string
	Player1    string
	Player2    string
	LimitDays  int
	LimitGames int
}

// PlayerCount returns number of players
func (f *Filter) playerCount() int {
	count := 0

	if f.Player1 != "" {
		count++
	}

	if f.Player2 != "" {
		count++
	}

	return count
}

// GetQuery returns SQL-query based on filters
// FIXME: Refactor
func (f *Filter) GetQuery() string {
	query := "SELECT id, game_name, is_tie, winner, loser, comment, added FROM match"

	// Needed so "AND" works in query
	noPlayers := " WHERE 1=1"
	onePlayer := ` WHERE (winner LIKE "` + f.Player1 + `" OR loser LIKE "` + f.Player1 + `")`
	twoPlayer := ` WHERE (winner LIKE "` + f.Player1 + `" OR loser LIKE "` + f.Player1 + `") AND (winner LIKE "` + f.Player2 + `" OR loser LIKE "` + f.Player2 + `")`
	gameName := ` AND (game_name LIKE "` + f.GameName + `")`
	order := " ORDER BY added DESC"
	limitDays := " AND (added > (SELECT DATETIME('now', '-" + strconv.Itoa(f.LimitDays) + " day')))"

	limitGames := " LIMIT " + strconv.Itoa(f.LimitGames)

	switch f.playerCount() {
	case 0:
		query = query + noPlayers
	case 1:
		query = query + onePlayer
	case 2:
		query = query + twoPlayer
	default:
		return ""
	}

	if f.GameName != "" {
		query = query + gameName
	}

	if f.LimitDays > 0 {
		query = query + limitDays
	}

	query = query + order

	if f.LimitGames > 0 {
		query = query + limitGames
	}

	return query
}
