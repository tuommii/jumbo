package database

import (
	"github.com/tuommii/jumbo/model"
)

// Database ...
type Database interface {
	GetPlayers() ([]model.Player, error)
	CreatePlayer(name string) (int64, error)
	DeletePlayer(name string) (int64, error)

	GetGames() ([]model.Game, error)
	CreateGame(name string) (int64, error)
	DeleteGame(name string) (int64, error)

	GetMatches(f model.Filter) ([]model.Match, error)
	CreateMatch(match model.Match) (int64, error)
	DeleteMatch(id int) (int64, error)
}
