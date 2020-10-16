package http

// Routes setups the http routes
func (s *Server) Routes() {
	// Health-check
	s.Router.GET("/", s.HandleHealthCheck())

	// NotFound
	s.Router.NoRoute(s.HandleNotFound())
}
