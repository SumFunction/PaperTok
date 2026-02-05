package auth

import "time"

// Config holds the configuration for the authentication service.
type Config struct {
	// Secret is the secret key used to sign JWT tokens.
	// In production, this should be loaded from environment variables.
	Secret string

	// AccessTokenExpiry is the duration for which access tokens are valid.
	AccessTokenExpiry time.Duration

	// RefreshTokenExpiry is the duration for which refresh tokens are valid.
	RefreshTokenExpiry time.Duration

	// Issuer is the issuer name for JWT tokens.
	Issuer string

	// PasswordCost is the bcrypt cost factor for password hashing.
	// Higher values are more secure but slower. Recommended: 10-12.
	PasswordCost int
}

// DefaultConfig returns a configuration with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Secret:             "change-this-secret-in-production",
		AccessTokenExpiry:  24 * time.Hour,
		RefreshTokenExpiry: 7 * 24 * time.Hour, // 7 days
		Issuer:             "papertok",
		PasswordCost:       10,
	}
}

// TestConfig returns a configuration suitable for testing.
func TestConfig() Config {
	return Config{
		Secret:             "test-secret-key-for-testing-only",
		AccessTokenExpiry:  24 * time.Hour,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		Issuer:             "papertok-test",
		PasswordCost:       4, // Lower cost for faster tests
	}
}

// Validate checks if the configuration is valid.
// For testing purposes, it accepts secrets ending with "-for-testing-only".
func (c Config) Validate() error {
	if c.Secret == "" {
		return ErrInvalidToken
	}
	if c.Secret == "change-this-secret-in-production" {
		return ErrInvalidToken
	}
	if c.AccessTokenExpiry <= 0 {
		return ErrInvalidToken
	}
	if c.PasswordCost < 4 || c.PasswordCost > 31 {
		return ErrInvalidToken
	}
	return nil
}
