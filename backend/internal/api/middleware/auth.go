package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rrlian/papertok/backend/internal/core/auth"
)

const (
	// UserIDKey is the key used to store the user ID in the Gin context.
	UserIDKey = "user_id"
	// UsernameKey is the key used to store the username in the Gin context.
	UsernameKey = "username"
	// EmailKey is the key used to store the email in the Gin context.
	EmailKey = "email"
)

// AuthMiddleware creates a middleware that validates JWT tokens.
// It extracts the token from the Authorization header,
// validates it, and stores user information in the context.
func AuthMiddleware(authSvc auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "MISSING_TOKEN",
					"message": "Authorization header is required",
				},
				"timestamp": c.Request.Context().Value("timestamp"),
			})
			c.Abort()
			return
		}

		// Extract token from Bearer format
		token := extractToken(authHeader)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN_FORMAT",
					"message": "Authorization header must be in format: Bearer <token>",
				},
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := authSvc.ValidateToken(c.Request.Context(), token)
		if err != nil {
			code := "INVALID_TOKEN"
			message := "Invalid or expired token"

			if err == auth.ErrExpiredToken {
				code = "TOKEN_EXPIRED"
				message = "Token has expired"
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    code,
					"message": message,
				},
			})
			c.Abort()
			return
		}

		// Store user info in context for downstream handlers
		c.Set(UserIDKey, claims.UserID)
		c.Set(UsernameKey, claims.Username)
		c.Set(EmailKey, claims.Email)

		// Also store as string for easier access
		c.Set("user_id", formatInt64(claims.UserID))

		c.Next()
	}
}

// extractToken extracts the JWT token from the Authorization header.
// It expects the header to be in the format: "Bearer <token>".
func extractToken(authHeader string) string {
	// Check if header starts with "Bearer "
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	// Return the header as-is if it doesn't have the Bearer prefix
	// This allows for custom token formats
	return authHeader
}

// formatInt64 converts an int64 to a string.
func formatInt64(n int64) string {
	const maxInt64Digits = 19
	var b [maxInt64Digits]byte
	i := len(b)
	neg := n < 0
	if neg {
		n = -n
	}

	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}

	if i == len(b) {
		return "0"
	}

	if neg {
		i--
		b[i] = '-'
	}

	return string(b[i:])
}

// OptionalAuthMiddleware creates a middleware that validates JWT tokens if present,
// but doesn't require authentication. This is useful for endpoints that have
// different behavior for authenticated vs anonymous users.
func OptionalAuthMiddleware(authSvc auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No token provided, continue without setting user context
			c.Next()
			return
		}

		// Extract token
		token := extractToken(authHeader)
		if token == "" {
			// Invalid format, continue without setting user context
			c.Next()
			return
		}

		// Validate token
		claims, err := authSvc.ValidateToken(c.Request.Context(), token)
		if err != nil {
			// Invalid token, continue without setting user context
			c.Next()
			return
		}

		// Store user info in context
		c.Set(UserIDKey, claims.UserID)
		c.Set(UsernameKey, claims.Username)
		c.Set(EmailKey, claims.Email)
		c.Set("user_id", formatInt64(claims.UserID))

		c.Next()
	}
}

// GetUserID retrieves the user ID from the Gin context.
// Returns false if the user is not authenticated.
func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}

	id, ok := userID.(int64)
	return id, ok
}

// GetUsername retrieves the username from the Gin context.
// Returns false if the user is not authenticated.
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get(UsernameKey)
	if !exists {
		return "", false
	}

	name, ok := username.(string)
	return name, ok
}

// GetEmail retrieves the email from the Gin context.
// Returns false if the user is not authenticated.
func GetEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(EmailKey)
	if !exists {
		return "", false
	}

	e, ok := email.(string)
	return e, ok
}

// IsAuthenticated checks if the current request is authenticated.
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get(UserIDKey)
	return exists
}
