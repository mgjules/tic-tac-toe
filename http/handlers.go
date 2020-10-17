package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mgjules/tic-tac-toe/game"
	ttt "github.com/shurcooL/tictactoe"
)

// HealthCheckResponse represents a response for HandleHealthCheck
type HealthCheckResponse struct {
	Pong time.Time `json:"pong" example:"2020-10-17T11:34:10.089762941+04:00"`
}

// HandleHealthCheck handles the health-check route
//
// @Summary Health checks the service
// @Description get a "pong" with current time
// @Produce  json
// @Success 200 {object} HealthCheckResponse "pong"
// @Router / [get]
func (s *Server) HandleHealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthCheckResponse{Pong: time.Now()})
	}
}

// MoveRequest represents a request for HandleMove
type MoveRequest struct {
	GameID   string `json:"game_id" binding:"required" example:"123123"`
	Mark     uint8  `json:"mark" binding:"numeric,oneof=1 2" enums:"1,2" example:"1"`                   // mark 'X' = 1. mark 'O' = 2.
	Position uint8  `json:"position" binding:"numeric,min=0,max=8" minimum:"0" maximum:"8" example:"5"` //
}

// MoveResponse represents a response for HandleMove
type MoveResponse struct {
	Error  string     `json:"error,omitempty"`
	Result string     `json:"result,omitempty"`
	State  *game.Game `json:"state,omitempty"`
}

// HandleMove handles the move route
//
// @Summary Makes a move
// @Description Moves a mark('x', 'o') on the tictactoe board
// @Accept  json
// @Produce  json
// @Param request body MoveRequest true "represents a move request"
// @Success 200 {object} MoveResponse "move successful / X won / O won / tie"
// @Success 400 {object} MoveResponse "validation failed on request"
// @Success 403 {object} MoveResponse "please wait your turn / forbidden move"
// @Router /move [post]
func (s *Server) HandleMove(gameRepository game.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req MoveRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, MoveResponse{Error: err.Error()})
			return
		}

		// load board from saved game state or create new one if not exist
		g, err := gameRepository.LoadGame(c, req.GameID)
		if err != nil && err != game.ErrGameNotFound {
			c.AbortWithStatusJSON(http.StatusInternalServerError, MoveResponse{Error: err.Error()})
			return
		} else if err == game.ErrGameNotFound {
			// no game found
			g = &game.Game{
				ID:         req.GameID,
				BoardCells: [9]uint8{},
			}
		}

		// convert uint8 board cells to ttt.State
		board := &ttt.Board{}
		for i, cell := range g.BoardCells {
			board.Cells[i] = ttt.State(cell)
		}

		// has the game ended? (i.e 'X' won, 'O' won or tie)
		condition := board.Condition()
		if condition != ttt.NotEnd {
			c.AbortWithStatusJSON(http.StatusOK, MoveResponse{Result: condition.String(), State: g})
			return
		}

		// is same player again?
		if g.LastMoveBy == req.Mark {
			c.AbortWithStatusJSON(http.StatusForbidden, MoveResponse{Error: "please wait your turn", State: g})
			return
		}

		// apply the move on the board
		move := ttt.Move(req.Position)
		state := ttt.State(req.Mark)
		if err := board.Apply(move, state); err != nil {
			// forbidden move
			c.AbortWithStatusJSON(http.StatusForbidden, MoveResponse{Error: err.Error(), State: g})
			return
		}

		g.LastMoveBy = req.Mark

		// convert ttt.State cells to uint8 board cells
		for i, cell := range board.Cells {
			g.BoardCells[i] = uint8(cell)
		}

		// save game state if not from cloud function
		if c.Request.UserAgent() != "TicTacToe/1.0" {
			if err := gameRepository.SaveGame(c, *g); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, MoveResponse{Error: err.Error()})
				return
			}
		}

		// has the game ended now?
		condition = board.Condition()
		if condition != ttt.NotEnd {
			c.AbortWithStatusJSON(http.StatusOK, MoveResponse{Result: condition.String(), State: g})
			return
		}

		c.JSON(http.StatusOK, MoveResponse{Result: "move successful", State: g})
	}
}

// HandleNotFound handles the not found route
func (s *Server) HandleNotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "unknown route"})
	}
}
