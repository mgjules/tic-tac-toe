package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HandleHealthCheck handles the health-check route
func (s *Server) HandleHealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"pong": time.Now()})
	}
}

// HandleMove handles the move route
func (s *Server) HandleMove() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: code
	}
}

// HandleNotFound handles the not found route
func (s *Server) HandleNotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "unknown route"})
	}
}
