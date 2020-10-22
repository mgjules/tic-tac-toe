package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/mgjules/tic-tac-toe/build"
	"github.com/mgjules/tic-tac-toe/config"
	"github.com/mgjules/tic-tac-toe/docs"
	"github.com/mgjules/tic-tac-toe/game/repository"
	rhttp "github.com/mgjules/tic-tac-toe/http"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const appName = "tictactoe"

// @title TicTacToe
// @version 1.0
// @description Microservice for TicTacToe
// @termsOfService http://swagger.io/terms/

// @contact.name Jules Michael
// @contact.email julesmichaelgiovanni@gmail.com

// @host localhost:3001
// @BasePath /

func main() {
	if err := run(); err != nil {
		fmt.Printf("[%s] %v", appName, err)
		os.Exit(1)
	}
}

func run() error {
	// Flags
	path := flag.String("path", ".env", "path of the config file")
	version := flag.Bool("v", false, "prints build information including version")
	flag.Parse()

	if *version {
		fmt.Printf("Version: %s\n", build.Version)
		fmt.Printf("Commit: %s\n", build.Commit)
		fmt.Printf("Branch: %s\n", build.Branch)
		fmt.Printf("Date: %s\n", build.Date)
		return nil
	}

	// Config
	cfg, err := config.Load(*path)
	if err != nil {
		return errors.Wrap(err, "can't load config")
	}

	// Swagger
	docs.SwaggerInfo.Host = cfg.Host + ":" + cfg.Port

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
	server.Routes(firebaseRepository, cfg.Host, cfg.Port)

	logger.Info("Gin server started on ", zap.String("host", cfg.Host), zap.String("port", cfg.Port))

	go func() {
		if err := server.Start(cfg.Host, cfg.Port); err != nil && err != http.ErrServerClosed {
			fmt.Printf("[%s] listen: %v", appName, err)
			os.Exit(1)
		}
	}()

	return server.WatchAndStop()
}
