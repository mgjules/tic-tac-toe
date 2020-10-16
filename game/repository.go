package game

import (
	"context"
	"errors"
)

// Repository errors
var (
	ErrGameNotFound = errors.New("game not found")
)

// Repository represents a game repository interface
type Repository interface {
	// LoadGame returns a game from a game id
	LoadGame(ctx context.Context, id string) (*Game, error)

	// SaveGame saves a game
	SaveGame(ctx context.Context, g Game) error
}
