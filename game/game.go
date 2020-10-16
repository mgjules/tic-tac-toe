package game

import (
	"errors"

	ttt "github.com/shurcooL/tictactoe"
)

// Game represents a game
type Game struct {
	ID         string
	LastMoveBy uint8
	Board      [9]uint8
}

// Validate validates the game
func (g *Game) Validate() error {
	if len(g.ID) == 0 {
		return errors.New("invalid game")
	}

	if ttt.State(g.LastMoveBy) != ttt.O && ttt.State(g.LastMoveBy) != ttt.X {
		return errors.New("invalid player")
	}

	// check each cell in board
	for _, cell := range g.Board {
		switch ttt.State(cell) {
		case ttt.O, ttt.X, ttt.F:
			// we are good
		default:
			return errors.New("invalid board")
		}
	}

	return nil
}
