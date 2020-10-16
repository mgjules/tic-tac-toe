package http

import "github.com/mgjules/tic-tac-toe/game"

// Routes setups the http routes
func (s *Server) Routes(gameRepository game.Repository) {
	// Health-check
	s.Router.GET("/", s.HandleHealthCheck())

	// Move
	s.Router.POST("/move", s.HandleMove(gameRepository))

	// NotFound
	s.Router.NoRoute(s.HandleNotFound())
}
