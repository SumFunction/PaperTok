package auth

import "time"

// Claims represents the JWT claims structure.
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	// Issuer is the issuer of the token.
	Issuer string `json:"iss,omitempty"`
	// Subject is the subject of the token.
	Subject string `json:"sub,omitempty"`
	// Audience is the audience of the token.
	Audience string `json:"aud,omitempty"`
}

// TokenInfo contains information about a generated token.
type TokenInfo struct {
	Token     string
	ExpiresAt time.Time
}

// TokenValidationResult contains the result of token validation.
type TokenValidationResult struct {
	Valid  bool
	Claims *Claims
	Error  error
}
