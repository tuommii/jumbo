package model

// Player ...
type Player struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// Game ...
type Game struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// Match ...
type Match struct {
	ID       int    `json:"id"`
	GameName string `json:"name"`
	Winner   string `json:"winner"`
	Loser    string `json:"loser"`
	IsTie    bool   `json:"isTie"`
	Added    string `json:"added"`
	Comment  string `json:"comment"`
	// WinnerScore int    `json:"winnerScore,omitempty"`
	// LoserScore  int    `json:"loserScore,omitempty"`
}
