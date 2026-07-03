package user

// Repository defines persistence operations for the user aggregate.
// Implementations live in infrastructure/persistence.
type Repository interface {
	// GetByID(id string) (aggregate.User, error)
	// Save(user aggregate.User) error
}
