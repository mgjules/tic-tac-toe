package game

import (
	"errors"

	ttt "github.com/shurcooL/tictactoe"
)

// Game represents a game
type Game struct {
	ID         string
	LastMoveBy ttt.State
	BoardCells [9]ttt.State
}

// Validate validates the game
func (g *Game) Validate() error {
	if len(g.ID) == 0 {
		return errors.New("invalid game")
	}

	if g.LastMoveBy != ttt.O && g.LastMoveBy != ttt.X {
		return errors.New("invalid state")
	}

	// check each cell in board
	for _, cell := range g.BoardCells {
		switch cell {
		case ttt.O, ttt.X, ttt.F:
			// we are good
		default:
			return errors.New("invalid board")
		}
	}

	return nil
}
