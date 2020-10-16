package http

// Routes setups the http routes
func (s *Server) Routes() {
	// Health-check
	s.Router.GET("/", s.HandleHealthCheck())

	// Move
	s.Router.POST("/move", s.HandleMove())

	// NotFound
	s.Router.NoRoute(s.HandleNotFound())
}
