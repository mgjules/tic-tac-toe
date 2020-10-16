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
	if len(corsAllowedOrigins) > 0 {
		corsCfg.AllowOrigins = corsAllowedOrigins
	} else {
		corsCfg.AllowAllOrigins = true
	}

	s.Router.Use(cors.New(corsCfg))
	s.Router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	s.Router.Use(ginzap.RecoveryWithZap(logger, true))
}
