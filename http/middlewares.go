package http

import (
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"go.uber.org/zap"
)

// Middlewares attaches server-wide middlewares
func (s *Server) Middlewares(logger *zap.Logger, corsAllowedOrigins []string) {
	// create cors config
	corsCfg := cors.DefaultConfig()
	if corsAllowedOrigins != nil && len(corsAllowedOrigins) > 0 {
		corsCfg.AllowOrigins = corsAllowedOrigins
	} else {
		corsCfg.AllowAllOrigins = true
	}

	s.router.Use(cors.New(corsCfg))
	s.router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	s.router.Use(ginzap.RecoveryWithZap(logger, true))
}
