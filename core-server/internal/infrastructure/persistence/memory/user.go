// Package memory is an in-memory implementation of domain repositories.
// Swap for postgres/mongo in production without changing domain or application layers.
package memory

import domainuser "backend/core-server/internal/domain/user"

// UserRepository is an in-memory user.Repository.
type UserRepository struct{}

// NewUserRepository creates an in-memory user repository.
func NewUserRepository() domainuser.Repository {
	return &UserRepository{}
}
