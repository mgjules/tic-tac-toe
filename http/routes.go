package http

// Routes setups the http routes
func (s *Server) Routes() {
	// Health-check
	s.router.GET("/", s.handleHealthCheck())

	// NotFound
	s.router.NoRoute(s.handleNotFound())
}
