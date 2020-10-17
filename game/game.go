package game

import (
	"errors"

	ttt "github.com/shurcooL/tictactoe"
)

// Game represents a game
type Game struct {
	ID         string   `json:"id" example:"123123"`
	LastMoveBy uint8    `json:"last_move_by" enums:"1,2" example:"1"`
	BoardCells [9]uint8 `json:"board_cells" example:"0,1,0,2,2,0,1,1,0"`
}

// Validate validates the game
func (g *Game) Validate() error {
	if g.ID == "" {
		return errors.New("invalid game")
	}

	if ttt.State(g.LastMoveBy) != ttt.O && ttt.State(g.LastMoveBy) != ttt.X {
		return errors.New("invalid state")
	}

	// check each cell in board
	for _, cell := range g.BoardCells {
		switch ttt.State(cell) {
		case ttt.O, ttt.X, ttt.F:
			// we are good
		default:
			return errors.New("invalid board")
		}
	}

	return nil
}
