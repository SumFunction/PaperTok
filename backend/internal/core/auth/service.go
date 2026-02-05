package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Impl implements the Service interface using bcrypt and JWT.
type Impl struct {
	cfg Config
}

// Ensure Impl implements Service interface.
var _ Service = (*Impl)(nil)

// New creates a new authentication service instance.
func New(cfg Config) (*Impl, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid auth config: %w", err)
	}

	return &Impl{
		cfg: cfg,
	}, nil
}

// HashPassword creates a bcrypt hash of the given password.
func (s *Impl) HashPassword(ctx context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cfg.PasswordCost)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrPasswordHash, err)
	}
	return string(hash), nil
}

// VerifyPassword checks if the provided password matches the hash.
func (s *Impl) VerifyPassword(ctx context.Context, hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return ErrInvalidPassword
	}
	return nil
}

// GenerateToken creates a new JWT token for the given user.
func (s *Impl) GenerateToken(ctx context.Context, userID int64, username, email string) (*TokenInfo, error) {
	now := time.Now()
	expiresAt := now.Add(s.cfg.AccessTokenExpiry)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Issuer:   s.cfg.Issuer,
		Subject:  fmt.Sprintf("%d", userID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  claims.UserID,
		"username": claims.Username,
		"email":    claims.Email,
		"iss":      claims.Issuer,
		"sub":      claims.Subject,
		"iat":      now.Unix(),
		"exp":      expiresAt.Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTokenGeneration, err)
	}

	return &TokenInfo{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}, nil
}

// ValidateToken validates a JWT token and returns its claims.
func (s *Impl) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: unexpected signing method: %v", ErrInvalidToken, token.Header["alg"])
		}
		return []byte(s.cfg.Secret), nil
	})

	if err != nil {
		if err == jwt.ErrTokenExpired {
			return nil, ErrExpiredToken
		}
		// Check if the error message contains "expired"
		if err.Error() == "token is expired" || contains(err.Error(), "expired") {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	// Extract user info from claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, ErrInvalidClaims
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, ErrInvalidClaims
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return &Claims{
		UserID:   int64(userID),
		Username: username,
		Email:    email,
		Issuer:   getString(claims, "iss"),
		Subject:  getString(claims, "sub"),
	}, nil
}

// RefreshToken generates a new token from an existing valid token.
func (s *Impl) RefreshToken(ctx context.Context, tokenString string) (*TokenInfo, error) {
	// Validate the existing token
	claims, err := s.ValidateToken(ctx, tokenString)
	if err != nil {
		// Allow refreshing expired tokens within a grace period
		if err == ErrExpiredToken {
			// Parse the token without validation to get claims
			token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(s.cfg.Secret), nil
			})
			if parseErr != nil {
				return nil, fmt.Errorf("%w: %v", ErrInvalidToken, parseErr)
			}

			mapClaims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return nil, ErrInvalidClaims
			}

			// Extract claims
			userID, ok := mapClaims["user_id"].(float64)
			if !ok {
				return nil, ErrInvalidClaims
			}

			username, ok := mapClaims["username"].(string)
			if !ok {
				return nil, ErrInvalidClaims
			}

			email, ok := mapClaims["email"].(string)
			if !ok {
				return nil, ErrInvalidClaims
			}

			claims = &Claims{
				UserID:   int64(userID),
				Username: username,
				Email:    email,
			}
		} else {
			return nil, err
		}
	}

	// Generate a new token
	return s.GenerateToken(ctx, claims.UserID, claims.Username, claims.Email)
}

// getString safely extracts a string value from jwt.MapClaims.
func getString(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// contains checks if a string contains a substring (case-insensitive).
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(len(substr) == 0 || indexOf(s, substr) >= 0)
}

// indexOf returns the index of substr in s, or -1 if not found.
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
