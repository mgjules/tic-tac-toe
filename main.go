package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mgjules/tic-tac-toe/config"
	rhttp "github.com/mgjules/tic-tac-toe/http"
	"github.com/mgjules/tic-tac-toe/repository"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const appName = "tictactoe"

func main() {
	if err := run(); err != nil {
		fmt.Printf("[%s] %v", appName, err)
		os.Exit(1)
	}
}

func run() error {
	// Config
	cfg, err := config.Load()
	if err != nil {
		return errors.Wrap(err, "can't load config")
	}

	// Logger
	var logger *zap.Logger
	if cfg.Prod {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return errors.Wrap(err, "can't create logger")
	}
	defer func() {
		if err = logger.Sync(); err != nil {
			fmt.Printf("[%s] logger: %v", appName, err)
			os.Exit(1)
		}
	}()

	// Firebase repository
	firebaseRepository, err := repository.NewFirebase(cfg.FirebaseDBURL, cfg.FirebaseServiceAccountKeyPath)
	if err != nil {
		return errors.Wrap(err, "can't create firebase repository")
	}

	// Server
	server := rhttp.NewServer(cfg.Prod)
	server.Middlewares(logger, cfg.CorsAllowedOrigins)
	server.Routes(firebaseRepository)

	logger.Info("Gin server started on ", zap.String("host", cfg.Host), zap.String("port", cfg.Port))

	go func() {
		if err := server.Start(cfg.Host, cfg.Port); err != nil && err != http.ErrServerClosed {
			fmt.Printf("[%s] listen: %v", appName, err)
			os.Exit(1)
		}
	}()

	return server.WatchAndStop()
}
