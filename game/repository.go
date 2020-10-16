package game

// Repository represents a game repository interface
type Repository interface {
	// LoadState returns a game state from an id
	LoadState(id string) (State, error)

	// SaveState saves a game state
	SaveState(state State) error
}
