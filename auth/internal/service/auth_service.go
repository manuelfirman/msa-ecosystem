package service

import (
	"auth/internal/domain"
	"errors"
)

var (
	// ErrInvalidCredentials is returned when the credentials are invalid
	ErrServiceCantHashPassword = errors.New("service: can't hash password")
	// ErrServiceDuplicateEntry is returned when a duplicate entry is found
	ErrServiceDuplicateEntry = errors.New("service: duplicate entry")
	// ErrServiceFailedToCreateUser is returned when the user can't be created
	ErrServiceFailedToCreateUser = errors.New("service: failed to create user")
	// ErrServiceInvalidCredentials is returned when the credentials are invalid
	ErrServiceInvalidCredentials = errors.New("service: invalid credentials")
)

// AuthService represents the Auth service
type AuthService interface {
	// Login logs in a user
	Login(*domain.Auth) error
	// Register registers a new user
	Register(*domain.User) error
}
