# UserAuth Feature Module

## Overview
This module implements user authentication functionality for PaperTok, including registration, login, and profile management.

## Architecture
The UserAuth feature follows the Vertical Slice Architecture (VSA) pattern with clear separation of concerns:

```
API Layer (handlers) -> Facade -> Feature (userauth) -> Core Services (auth) -> Repository (user)
```

## Module Structure

### Files
- `interface.go` - Service interface definition
- `deps.go` - Dependency interface definitions (authService, userRepository)
- `types.go` - Domain types (User, RegisterRequest, LoginRequest, AuthResponse, ProfileResponse)
- `errors.go` - Error definitions with error codes
- `service.go` - Business logic implementation
- `service_test.go` - Unit tests

## Dependencies

### Core Services
- `auth.Service` - JWT token generation/validation and password hashing

### Repositories
- `user.Repository` - User data access (Create, FindByEmail, FindByUsername, FindByID, Exists*)

## API Endpoints

### Public Routes
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token

### Protected Routes (require JWT)
- `GET /api/v1/auth/profile` - Get current user profile

## Request/Response Formats

### Register Request
```json
{
  "username": "string (3-50 chars, alphanumeric + underscore)",
  "email": "string (valid email)",
  "password": "string (min 8 chars, recommended 2+ character types)"
}
```

### Login Request
```json
{
  "identifier": "string (email or username)",
  "password": "string"
}
```

### Auth Response
```json
{
  "user": {
    "id": 123,
    "username": "string",
    "email": "string",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "token": "jwt_token_string"
}
```

### Profile Response
```json
{
  "id": 123,
  "username": "string",
  "email": "string",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

## Error Codes

| Code | Description | HTTP Status |
|------|-------------|-------------|
| USER_EXISTS | User with email/username already exists | 409 |
| INVALID_CREDENTIALS | Wrong email/username or password | 401 |
| USER_NOT_FOUND | User not found | 404 |
| INVALID_EMAIL | Invalid email format | 400 |
| INVALID_USERNAME | Invalid username format | 400 |
| WEAK_PASSWORD | Password too weak | 400 |
| VALIDATION_ERROR | Input validation failed | 400 |
| UNAUTHORIZED | Not authenticated | 401 |
| INTERNAL_ERROR | Server error | 500 |

## Security Features

1. **Password Hashing**: Uses bcrypt with cost factor 10
2. **JWT Tokens**: Signed with HS256, configurable expiration
3. **Token Validation**: Middleware validates tokens on protected routes
4. **Input Validation**: Username, email, and password validation

## Usage Example

### Registration
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"SecurePass123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"test@example.com","password":"SecurePass123"}'
```

### Get Profile (Protected)
```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer <jwt_token>"
```

## Testing

Run tests:
```bash
go test -v ./internal/features/userauth/...
```

## Notes

- Current implementation uses in-memory repository for testing
- To use MySQL repository, update Facade initialization with SQL repository
- JWT secret should be loaded from environment variable (JWT_SECRET) in production
