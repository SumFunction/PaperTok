# Auth Core Service

## Overview
Core authentication service that provides password hashing and JWT token management.

## Module Structure

### Files
- `interface.go` - Service interface definition
- `deps.go` - Configuration types and defaults
- `types.go` - Domain types (Claims, TokenInfo, TokenValidationResult)
- `errors.go` - Error definitions
- `service.go` - Implementation using bcrypt and JWT
- `service_test.go` - Unit tests

## Configuration

```go
type Config struct {
    Secret             string        // JWT signing secret
    AccessTokenExpiry  time.Duration // Token lifetime
    RefreshTokenExpiry time.Duration // Refresh token lifetime
    Issuer             string        // Token issuer
    PasswordCost       int           // bcrypt cost factor (4-31)
}
```

## API

### HashPassword
Creates a bcrypt hash of the given password.
```go
HashPassword(ctx context.Context, password string) (string, error)
```

### VerifyPassword
Verifies a password against its hash.
```go
VerifyPassword(ctx context.Context, hashedPassword, password string) error
```

### GenerateToken
Generates a JWT token for a user.
```go
GenerateToken(ctx context.Context, userID int64, username, email string) (*TokenInfo, error)
```

### ValidateToken
Validates a JWT token and returns its claims.
```go
ValidateToken(ctx context.Context, token string) (*Claims, error)
```

### RefreshToken
Generates a new token from an existing (possibly expired) token.
```go
RefreshToken(ctx context.Context, token string) (*TokenInfo, error)
```

## JWT Claims Structure

```go
type Claims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Issuer   string `json:"iss"`
    Subject  string `json:"sub"`
}
```

## Security Notes

1. **Password Cost**: bcrypt cost factor of 10 is recommended for production
2. **Secret Management**: Use environment variable for JWT secret in production
3. **Token Expiration**: Default 24 hours for access tokens, 7 days for refresh tokens
4. **Algorithm**: Uses HS256 (HMAC-SHA256) for signing

## Testing

### Test Config
For testing, use `TestConfig()` which provides:
- Lower bcrypt cost (4) for faster tests
- Test-friendly secret key
- Standard token expiration

```go
cfg := auth.TestConfig()
svc, err := auth.New(cfg)
```

## Dependencies

- `golang.org/x/crypto/bcrypt` - Password hashing
- `github.com/golang-jwt/jwt/v5` - JWT token handling
