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

// HandleNotFound handles the not found route
func (s *Server) HandleNotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "unknown route"})
	}
}
