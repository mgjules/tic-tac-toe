package repository

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/mgjules/tic-tac-toe/game"
	"google.golang.org/api/option"
)

// Firebase represents a firebase repository
type Firebase struct {
	client *db.Client
}

// NewFirebase returns a new Firebase repository
func NewFirebase(url, serviceAccountKeyPath string) (*Firebase, error) {
	ctx := context.Background()

	conf := &firebase.Config{
		DatabaseURL: url,
	}
	opt := option.WithCredentialsFile(serviceAccountKeyPath)

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Database(ctx)
	if err != nil {
		return nil, err
	}

	return &Firebase{
		client: client,
	}, nil
}

// LoadGame returns a game from a game id
func (f *Firebase) LoadGame(ctx context.Context, id string) (*game.Game, error) {
	ref := f.client.NewRef("tictactoe/games")
	results, err := ref.OrderByKey().EqualTo(id).LimitToFirst(1).GetOrdered(ctx)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, game.ErrGameNotFound
	}

	var g game.Game
	if err := results[0].Unmarshal(&g); err != nil {
		return nil, err
	}

	if err := g.Validate(); err != nil {
		return nil, err
	}

	return &g, nil
}

// SaveGame saves a game
func (f *Firebase) SaveGame(ctx context.Context, g game.Game) error {
	if err := g.Validate(); err != nil {
		return err
	}

	ref := f.client.NewRef("tictactoe/games")
	if err := ref.Child(g.ID).Set(ctx, g); err != nil {
		return err
	}

	return nil
}
