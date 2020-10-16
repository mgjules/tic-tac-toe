package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// Server represents a simple GIN HTTP server
type Server struct {
	Router *gin.Engine

	*http.Server
}

// NewServer setups a new Server
func NewServer(isProd bool) *Server {
	if isProd {
		gin.SetMode(gin.ReleaseMode)
	}

	s := &Server{
		Router: gin.New(),
	}

	return s
}

// Start starts the HTTP Server
func (s *Server) Start(host, port string) error {
	s.Server = &http.Server{
		Addr:           host + ":" + port,
		Handler:        s.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 5,
	}

	return s.ListenAndServe()
}

// WatchAndStop gracefully shutdowns the server
func (s *Server) WatchAndStop() error {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.Server.Shutdown(ctx)
}
