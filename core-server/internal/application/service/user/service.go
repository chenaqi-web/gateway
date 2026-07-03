// Package user provides application-level use cases for the user subdomain.
package user

import domainuser "backend/core-server/internal/domain/user"

// Service coordinates user-related use cases.
type Service struct {
	repo domainuser.Repository
}

// New creates a user application service.
func New(repo domainuser.Repository) *Service {
	return &Service{repo: repo}
}
