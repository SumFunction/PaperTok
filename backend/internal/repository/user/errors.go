package user

import "errors"

// Common errors for user repository operations.
var (
	// ErrUserNotFound is returned when a user is not found in the repository.
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists is returned when attempting to create a user
	// that already exists (duplicate email or username).
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrInvalidID is returned when an invalid user ID is provided.
	ErrInvalidID = errors.New("invalid user ID")

	// ErrInvalidEmail is returned when an invalid email is provided.
	ErrInvalidEmail = errors.New("invalid email format")

	// ErrInvalidUsername is returned when an invalid username is provided.
	ErrInvalidUsername = errors.New("invalid username format")
)
