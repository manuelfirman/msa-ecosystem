package repository

import (
	"auth/internal/domain"
	"errors"
)

var (
	// ErrRepositoryDuplicateEntry is returned when a duplicate entry is found
	ErrRepositoryDuplicateEntry = errors.New("repository: duplicate entry")
	// ErrRepositoryTableNotFound is returned when the table is not found
	ErrRepositoryTableNotFound = errors.New("repository: table not found")
	// ErrRepositoryEntryNotFound is returned when the entry is not found
	ErrRepositoryEntryNotFound = errors.New("repository: entry not found")
)

// AuthRepository represents the Auth repository
type AuthRepository interface {
	// Login logs in a user
	Login(*domain.Auth) error
	// Register registers a new user
	Register(*domain.Auth) error
}
