package game

// Repository represents a game repository interface
type Repository interface {
	// LoadGame returns a game from a game id
	LoadGame(id string) (Game, error)

	// SaveGame saves a game
	SaveGame(state Game) error
}
