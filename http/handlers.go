package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) handleHealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"pong": time.Now()})
	}
}

func (s *Server) handleNotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "unknown route"})
	}
}
