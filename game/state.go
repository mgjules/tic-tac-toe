package game

import (
	"errors"

	ttt "github.com/shurcooL/tictactoe"
)

// State represents a game state
type State struct {
	ID         string
	LastMoveBy uint8
	BoardCells [9]uint8
}

// Validate validates the game state
func (s *State) Validate() error {
	if len(s.ID) == 0 {
		return errors.New("invalid game")
	}

	if ttt.State(s.LastMoveBy) != ttt.O && ttt.State(s.LastMoveBy) != ttt.X {
		return errors.New("invalid player")
	}

	// check each cell in board
	for _, cell := range s.BoardCells {
		switch ttt.State(cell) {
		case ttt.O, ttt.X, ttt.F:
			// we are good
		default:
			return errors.New("invalid board")
		}
	}

	return nil
}
