package auth

import "errors"

// Common errors for authentication operations.
var (
	// ErrInvalidToken is returned when a JWT token is invalid or malformed.
	ErrInvalidToken = errors.New("invalid token")

	// ErrExpiredToken is returned when a JWT token has expired.
	ErrExpiredToken = errors.New("token expired")

	// ErrInvalidPassword is returned when the password doesn't match the hash.
	ErrInvalidPassword = errors.New("invalid password")

	// ErrInvalidClaims is returned when the token claims are invalid.
	ErrInvalidClaims = errors.New("invalid token claims")

	// ErrMissingToken is returned when no token is provided in the request.
	ErrMissingToken = errors.New("missing authentication token")

	// ErrTokenGeneration is returned when token generation fails.
	ErrTokenGeneration = errors.New("failed to generate token")

	// ErrPasswordHash is returned when password hashing fails.
	ErrPasswordHash = errors.New("failed to hash password")
)
