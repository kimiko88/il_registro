package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"registro-backend/pkg/jwt"
	"registro-backend/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

// Service handles authentication business logic
type Service struct {
	repo         Repository
	tokenManager *jwt.TokenManager
	mfaService   *MFAService
	bcryptCost   int
}

// NewService creates a new auth service
func NewService(repo Repository, tokenManager *jwt.TokenManager, mfaService *MFAService) *Service {
	return &Service{
		repo:         repo,
		tokenManager: tokenManager,
		mfaService:   mfaService,
		bcryptCost:   12,
	}
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
	// Validate input
	if err := ValidateRegisterRequest(req); err != nil {
		return nil, err
	}

	// Check if email already exists
	_, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.bcryptCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &User{
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
	}
	if req.SchoolID != "" {
		user.SchoolID = &req.SchoolID
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *Service) Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*AuthResponse, error) {
	// Check rate limiting
	since := time.Now().Add(-15 * time.Minute)
	attempts, err := s.repo.GetRecentLoginAttempts(ctx, req.Email, ipAddress, since)
	if err != nil {
		return nil, err
	}
	if attempts >= 5 {
		return nil, ErrTooManyAttempts
	}

	// Get user
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.recordFailedAttempt(ctx, req.Email, ipAddress)
		return nil, ErrInvalidCredentials
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		s.recordFailedAttempt(ctx, req.Email, ipAddress)
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Check if password has expired (90 days for admin/secretary/superadmin)
	if user.Role == "admin" || user.Role == "superadmin" || user.Role == "secretary" {
		if user.PasswordChangedAt != nil && time.Since(*user.PasswordChangedAt) > 90*24*time.Hour {
			return nil, ErrPasswordExpired
		}
	}

	// Check MFA
	if user.MFAEnabled {
		if req.MFAToken == "" {
			return nil, ErrMFARequired
		}
		secret, err := s.repo.GetMFASecret(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		if !s.mfaService.VerifyTOTP(secret, req.MFAToken) {
			return nil, ErrInvalidMFAToken
		}
	}

	// Generate tokens
	accessToken, err := s.tokenManager.GenerateAccessToken(
		user.ID, user.Email, user.Role,
		func() string {
			if user.SchoolID != nil {
				return *user.SchoolID
			}
			return ""
		}(),
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	rt := &RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
	if err := s.repo.CreateRefreshToken(ctx, rt); err != nil {
		return nil, err
	}

	// Update last login
	_ = s.repo.UpdateLastLogin(ctx, user.ID)

	// Record successful attempt
	s.recordSuccessfulAttempt(ctx, req.Email, ipAddress)

	return &AuthResponse{
		User: &UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Role:          user.Role,
			SchoolID:      user.SchoolID,
			EmailVerified: user.EmailVerified,
			MFAEnabled:    user.MFAEnabled,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.tokenManager.GetAccessTokenTTL(),
	}, nil
}

// RefreshToken generates new access token from refresh token
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Get refresh token from DB
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	// Check if revoked
	if rt.Revoked {
		return nil, ErrTokenRevoked
	}

	// Check if expired
	if time.Now().After(rt.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	// Get user
	user, err := s.repo.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new access token
	accessToken, err := s.tokenManager.GenerateAccessToken(
		user.ID, user.Email, user.Role,
		func() string {
			if user.SchoolID != nil {
				return *user.SchoolID
			}
			return ""
		}(),
	)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.tokenManager.GetAccessTokenTTL(),
	}, nil
}

// Logout revokes ALL refresh tokens for the user (not only the one presented).
// This ensures that stolen tokens from other devices/sessions are also invalidated.
// If the refresh token is not found or already expired, logout still succeeds
// as long as we can extract the userID from the JWT access token via context.
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		// Token not found or already revoked: still a valid logout intent.
		// We cannot revoke all sessions without the userID, so we return nil
		// (no-op — the token was already invalid).
		return nil
	}
	// Revoke ALL sessions for this user, not just the one presented.
	// This is the secure behavior: if a device is compromised the user
	// can log out from any session and all others are terminated.
	return s.repo.RevokeAllUserTokens(ctx, rt.UserID)
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

// SetupMFA initiates MFA setup for a user
func (s *Service) SetupMFA(ctx context.Context, userID string) (*MFASetupResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.MFAEnabled {
		return nil, ErrMFAAlreadyEnabled
	}

	// Generate secret
	secret, err := s.mfaService.GenerateSecret()
	if err != nil {
		return nil, err
	}

	// Generate QR code URL
	qrURL := s.mfaService.GenerateQRCodeURL(user.Email, secret)

	// Generate recovery codes
	recoveryCodes, err := s.mfaService.GenerateRecoveryCodes(10)
	if err != nil {
		return nil, err
	}

	// Store recovery codes
	if err := s.repo.CreateRecoveryCodes(ctx, userID, recoveryCodes); err != nil {
		return nil, err
	}

	// Temporarily store secret (will be confirmed on verification)
	// For now, we'll enable it immediately - in production, require verification first
	if err := s.repo.EnableMFA(ctx, userID, secret); err != nil {
		return nil, err
	}

	return &MFASetupResponse{
		Secret:        secret,
		QRCodeURL:     qrURL,
		RecoveryCodes: recoveryCodes,
	}, nil
}

// VerifyMFA verifies MFA token
func (s *Service) VerifyMFA(ctx context.Context, userID, token string) error {
	secret, err := s.repo.GetMFASecret(ctx, userID)
	if err != nil {
		return err
	}

	if !s.mfaService.VerifyTOTP(secret, token) {
		return ErrInvalidMFAToken
	}

	return nil
}

// RequestPasswordReset creates a password reset token
func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists
		return nil
	}

	// Generate reset token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return err
	}
	token := hex.EncodeToString(tokenBytes)

	// Store token
	prt := &PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	if err := s.repo.CreatePasswordResetToken(ctx, prt); err != nil {
		return err
	}

	// Log the token for local development / testing
	// In production, this would send an email.
	logger.Log.Infof("PASSWORD RESET REQUEST for %s. Reset Token: %s", email, token)

	return nil
}

// ResetPassword resets user password with token
func (s *Service) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validate new password
	passwordValidator := NewPasswordValidator()
	if err := passwordValidator.Validate(newPassword); err != nil {
		return err
	}

	// Get reset token
	prt, err := s.repo.GetPasswordResetToken(ctx, token)
	if err != nil {
		return err
	}

	// Check if already used
	if prt.Used {
		return ErrInvalidToken
	}

	// Check if expired
	if time.Now().After(prt.ExpiresAt) {
		return ErrInvalidToken
	}

	// Validate password complexity
	validator := NewPasswordValidator()
	if err := validator.Validate(newPassword); err != nil {
		return err
	}

	// Check password history (prevent reuse of last 5)
	history, err := s.repo.GetPasswordHistory(ctx, prt.UserID)
	if err == nil {
		for _, oldHash := range history {
			if bcrypt.CompareHashAndPassword([]byte(oldHash), []byte(newPassword)) == nil {
				return ErrPasswordReused
			}
		}
	}

	// Hash new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), s.bcryptCost)
	if err != nil {
		return err
	}

	// Update password
	if err := s.repo.UpdatePassword(ctx, prt.UserID, string(passwordHash)); err != nil {
		return err
	}

	// Add to password history
	_ = s.repo.AddPasswordHistory(ctx, prt.UserID, string(passwordHash))

	// Mark token as used
	if err := s.repo.UsePasswordResetToken(ctx, prt.ID); err != nil {
		return err
	}

	// Revoke all refresh tokens for security
	_ = s.repo.RevokeAllUserTokens(ctx, prt.UserID)

	return nil
}

func (s *Service) recordFailedAttempt(ctx context.Context, email, ipAddress string) {
	attempt := &LoginAttempt{
		Email:       email,
		IPAddress:   ipAddress,
		Success:     false,
		AttemptedAt: time.Now(),
	}
	_ = s.repo.RecordLoginAttempt(ctx, attempt)
}

func (s *Service) recordSuccessfulAttempt(ctx context.Context, email, ipAddress string) {
	attempt := &LoginAttempt{
		Email:       email,
		IPAddress:   ipAddress,
		Success:     true,
		AttemptedAt: time.Now(),
	}
	_ = s.repo.RecordLoginAttempt(ctx, attempt)
}
