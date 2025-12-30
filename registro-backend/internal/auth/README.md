# Authentication Service

Complete authentication service for the Electronic School Register with enterprise-grade security.

## Features

- **Email/Password Authentication**: bcrypt hashing with cost factor 12
- **JWT Tokens**: RS256 asymmetric signing
  - Access tokens: 15 minutes expiration
  - Refresh tokens: 7 days expiration with revocation support
- **Multi-Factor Authentication (MFA)**: TOTP compatible with Google Authenticator
- **Rate Limiting**: 5 failed login attempts per 15 minutes
- **Password Reset**: Secure token-based flow
- **SPID/CIE Integration**: Ready for Italian digital identity (SAML 2.0 stubs)

## API Endpoints

### Public Endpoints

#### Register
```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "first_name": "Mario",
  "last_name": "Rossi",
  "role": "student",
  "school_id": "uuid-optional"
}
```

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "mfa_token": "123456" // optional, required if MFA is enabled
}
```

Response:
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "Mario",
    "last_name": "Rossi",
    "role": "student",
    "mfa_enabled": false
  },
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_in": 900
}
```

#### Refresh Token
```http
POST /auth/refresh-token
Content-Type: application/json

{
  "refresh_token": "eyJhbGc..."
}
```

#### Request Password Reset
```http
POST /auth/password-reset
Content-Type: application/json

{
  "email": "user@example.com"
}
```

#### Confirm Password Reset
```http
POST /auth/password-reset/confirm
Content-Type: application/json

{
  "token": "reset-token-from-email",
  "new_password": "NewSecurePass123!"
}
```

### Protected Endpoints

Require `Authorization: Bearer <access_token>` header.

#### Logout
```http
POST /auth/logout
Authorization: Bearer eyJhbGc...

{
  "refresh_token": "eyJhbGc..."
}
```

#### Setup MFA
```http
POST /auth/mfa/setup
Authorization: Bearer eyJhbGc...
```

Response:
```json
{
  "secret": "BASE32SECRET",
  "qr_code_url": "otpauth://totp/...",
  "recovery_codes": [
    "ABCD1234",
    "EFGH5678",
    ...
  ]
}
```

#### Verify MFA
```http
POST /auth/mfa/verify
Authorization: Bearer eyJhbGc...

{
  "token": "123456"
}
```

## Security

### Password Requirements
- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- At least one special character

### Rate Limiting
- Maximum 5 failed login attempts per 15 minutes per IP/email combination
- Attempts are tracked in the database
- Successful login resets the counter

### Token Management
- Access tokens use RS256 (asymmetric) signing
- Refresh tokens are stored in database with revocation support
- All user tokens can be revoked on password change
- Tokens contain: user_id, email, role, school_id

### MFA/TOTP
- Compatible with Google Authenticator, Authy, etc.
- 6-digit codes with 30-second validity window
- 10 recovery codes generated on setup
- Recovery codes are single-use

## Database Tables

### refresh_tokens
- Stores refresh tokens with metadata
- Indexed by token, user_id, expires_at
- Supports revocation

### password_reset_tokens
- One-time use tokens
- 1-hour expiration
- Marked as used after password reset

### login_attempts
- Tracks all login attempts (success/failure)
- Used for rate limiting
- Indexed by email, ip_address, attempted_at

### mfa_recovery_codes
- Backup codes for MFA
- Single-use
- 10 codes per user

## Usage Example

```go
// Initialize dependencies
db, _ := sql.Open("postgres", connString)
privateKey, publicKey := loadKeys()
tokenManager := jwt.NewTokenManager(privateKey, publicKey)
mfaService := auth.NewMFAService("My School")

// Create service
repo := auth.NewRepository(db)
service := auth.NewService(repo, tokenManager, mfaService)

// Create handler
handler := auth.NewHandler(service)
middleware := auth.NewMiddleware(tokenManager)

// Register routes
router := gin.Default()
api := router.Group("/api/v1")
handler.RegisterRoutes(api, middleware)
```

## SPID/CIE Integration

SPID and CIE handlers are provided as stubs. To complete integration:

1. Add dependency: `github.com/crewjam/saml`
2. Configure Service Provider metadata
3. Add IDP metadata endpoints
4. Implement SAML assertion validation
5. Map SPID/CIE attributes to User model

See `spid_handler.go` and `cie_handler.go` for detailed implementation notes.

## Testing

Run tests:
```bash
go test ./internal/auth/...
```

Tests use testify/mock for repository mocking and table-driven test patterns.
