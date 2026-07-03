// Package domain is the DDD Domain Layer — pure business logic, no framework dependencies.
//
// Structure per subdomain (e.g. user/):
//   - aggregate/   — aggregate roots combining entities and value objects
//   - repository.go — repository interface (persistence contract)
//   - errors.go    — domain-specific errors
//
// Shared across subdomains:
//   - entity/      — shared entities with identity
//   - valueobject/ — immutable value objects without identity
package domain
