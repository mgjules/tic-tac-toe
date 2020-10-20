package repository

import (
	"context"
	"sync"

	"github.com/mgjules/tic-tac-toe/game"
)

// Mock represents a Mock repository
type Mock struct {
	mu sync.Mutex
	db map[string]*game.Game
}

// NewMock returns a new Mock repository
func NewMock() *Mock {
	return &Mock{
		db: map[string]*game.Game{},
	}
}

// LoadGame returns a game from a game id
func (m *Mock) LoadGame(ctx context.Context, id string) (*game.Game, error) {
	g, found := m.db[id]
	if !found {
		return nil, game.ErrGameNotFound
	}

	if err := g.Validate(); err != nil {
		return nil, err
	}

	return g, nil
}

// SaveGame saves a game
func (m *Mock) SaveGame(ctx context.Context, g game.Game) error {
	if err := g.Validate(); err != nil {
		return err
	}

	m.mu.Lock()
	m.db[g.ID] = &g
	m.mu.Unlock()

	return nil
}
