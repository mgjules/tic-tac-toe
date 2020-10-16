package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mgjules/tic-tac-toe/game"
	ttt "github.com/shurcooL/tictactoe"
)

// HandleHealthCheck handles the health-check route
func (s *Server) HandleHealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"pong": time.Now()})
	}
}

// HandleMove handles the move route
func (s *Server) HandleMove(gameRepository game.Repository) gin.HandlerFunc {
	type Request struct {
		GameID   string `json:"game_id" binding:"required"`
		Mark     uint8  `json:"mark" binding:"numeric,oneof=1 2"`
		Position uint8  `json:"position" binding:"numeric,min=0,max=8"`
	}
	return func(c *gin.Context) {
		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		state := ttt.State(req.Mark)

		// load board from saved game state or create new one if not exist
		g, err := gameRepository.LoadGame(c, req.GameID)
		if err != nil && err != game.ErrGameNotFound {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if err == game.ErrGameNotFound {
			// no game found
			g = &game.Game{
				ID:         req.GameID,
				BoardCells: [9]ttt.State{},
			}
		}

		board := &ttt.Board{Cells: g.BoardCells}

		// has the game ended? (i.e 'X' won, 'O' won or tie)
		condition := board.Condition()
		if condition != ttt.NotEnd {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"result": condition.String(), "state": g})
			return
		}

		// is same player again?
		if g.LastMoveBy == state {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "please wait your turn", "state": g})
			return
		}

		// apply the move on the board
		move := ttt.Move(req.Position)
		if err := board.Apply(move, state); err != nil {
			// forbidden move
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error(), "state": g})
			return
		}

		// save game state if not from cloud function
		g.LastMoveBy = state
		g.BoardCells = board.Cells
		if c.Request.UserAgent() != "TicTacToe/1.0" {
			if err := gameRepository.SaveGame(c, *g); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// has the game ended now?
		condition = board.Condition()
		if condition != ttt.NotEnd {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"result": condition.String(), "state": g})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": "move successful", "state": g})
	}
}

// HandleNotFound handles the not found route
func (s *Server) HandleNotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "unknown route"})
	}
}
