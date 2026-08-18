package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"registro-backend/pkg/crypto"
	"registro-backend/pkg/jwt"
	"registro-backend/pkg/logger"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type EmailSender interface {
	SendPasswordReset(ctx context.Context, email, token string) error
}

// Service handles authentication business logic
type Service struct {
	repo            Repository
	tokenManager    *jwt.TokenManager
	mfaService      *MFAService
	emailSender     EmailSender
	bcryptCost      int
	dummyBcryptHash string
}

// NewService creates a new auth service
func NewService(repo Repository, tokenManager *jwt.TokenManager, mfaService *MFAService, emailSender EmailSender) *Service {
	cost := 12
	dummyHash := "$2a$12$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeg6Lruj3vjPGga31lW"
	if h, err := bcrypt.GenerateFromPassword([]byte("dummy_password_for_timing_protection"), cost); err == nil {
		dummyHash = string(h)
	}
	return &Service{
		repo:            repo,
		tokenManager:    tokenManager,
		mfaService:      mfaService,
		emailSender:     emailSender,
		bcryptCost:      cost,
		dummyBcryptHash: dummyHash,
	}
}

// SetEmailSender sets the EmailSender for password reset emails
func (s *Service) SetEmailSender(sender EmailSender) {
	s.emailSender = sender
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// Register creates a new user account.
// Input validation and RBAC checks are performed by the handler layer
// (ValidateRegisterRequest in validator.go) before this method is called.
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
	normalizedEmail := normalizeEmail(req.Email)

	// Check if email already exists
	_, err := s.repo.GetUserByEmail(ctx, normalizedEmail)
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
		Email:        normalizedEmail,
		PasswordHash: string(passwordHash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
	}
	if req.SchoolID != "" {
		user.SchoolID = &req.SchoolID
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, ErrEmailAlreadyExists
		}
		if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "duplicate key") {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	// Seed password history so ResetPassword can detect immediate re-use.
	if err := s.repo.AddPasswordHistory(ctx, user.ID, string(passwordHash)); err != nil {
		logger.Log.Warnf("Register: failed to seed initial password history for user %s: %v", user.ID, err)
	}

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *Service) Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*AuthResponse, error) {
	req.Email = normalizeEmail(req.Email)

	// --- Rate limiting ---
	// 1. Per (email, IP): max 5 attempts in 15 minutes — blocks single-IP bursts.
	// 2. Per email only: max 20 attempts in 1 hour — blocks distributed IP-rotation attacks.
	since15m := time.Now().Add(-15 * time.Minute)
	attemptsPerIP, err := s.repo.GetRecentLoginAttempts(ctx, req.Email, ipAddress, since15m)
	if err != nil {
		return nil, err
	}
	if attemptsPerIP >= 5 {
		return nil, ErrTooManyAttempts
	}

	since1h := time.Now().Add(-1 * time.Hour)
	attemptsPerEmail, err := s.repo.GetRecentLoginAttemptsByEmail(ctx, req.Email, since1h)
	if err != nil {
		return nil, err
	}
	if attemptsPerEmail >= 20 {
		return nil, ErrTooManyAttempts
	}

	// Get user
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// Run bcrypt against dummy hash to prevent timing attacks that enumerate valid user accounts
		_ = bcrypt.CompareHashAndPassword([]byte(s.dummyBcryptHash), []byte(req.Password))
		s.recordFailedAttempt(ctx, req.Email, ipAddress)
		return nil, ErrInvalidCredentials
	}

	// Verify password BEFORE active check to maintain constant response timing
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if bcryptErr != nil {
		s.recordFailedAttempt(ctx, req.Email, ipAddress)
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Check if password has expired (90 days for privileged roles)
	if expired, err := isPasswordExpiredReason(user); expired {
		return nil, err
	}

	// Check MFA — decrypt stored secret before verifying the TOTP code.
	if user.MFAEnabled {
		if req.MFAToken == "" {
			return nil, ErrMFARequired
		}
		encryptedSecret, err := s.repo.GetMFASecret(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		secret, err := crypto.DecryptString(encryptedSecret)
		if err != nil {
			return nil, err
		}
		if !s.mfaService.VerifyTOTPUser(user.ID, secret, req.MFAToken) {
			s.recordFailedAttempt(ctx, req.Email, ipAddress)
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

	if len(userAgent) > 512 {
		userAgent = userAgent[:512]
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
			IsStaff:       user.IsStaff,
			EmailVerified: user.EmailVerified,
			MFAEnabled:    user.MFAEnabled,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.tokenManager.GetAccessTokenTTL(),
	}, nil
}

// RefreshToken generates new access token from refresh token.
// It also verifies that the user account is still active before issuing
// a new token — prevents disabled accounts from silently regaining access.
func (s *Service) RefreshToken(ctx context.Context, refreshToken, ipAddress, userAgent string) (*TokenPair, error) {
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

	// Soft anomaly check: log warning if refresh token IP differs from original IP
	if rt.IPAddress != "" && ipAddress != "" && rt.IPAddress != ipAddress {
		logger.Log.Warnf("Refresh token IP mismatch for user %s: issued at %s, used from %s", rt.UserID, rt.IPAddress, ipAddress)
	}

	// Get user
	user, err := s.repo.GetUserByID(ctx, rt.UserID)
	if err != nil {
		return nil, err
	}

	// Verify account is still active before issuing a new access token.
	// A disabled account must not be able to obtain new tokens via refresh.
	if !user.IsActive {
		// Revoke the refresh token so it cannot be retried.
		_ = s.repo.RevokeRefreshToken(ctx, rt.ID)
		return nil, ErrUserInactive
	}

	if isPasswordExpired(user) {
		_ = s.repo.RevokeRefreshToken(ctx, rt.ID)
		return nil, ErrPasswordExpired
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

	// Generate a new refresh token
	newRefreshToken, err := s.tokenManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Update IP and UserAgent from current request if provided, falling back to previous token's values
	newIp := ipAddress
	if newIp == "" {
		newIp = rt.IPAddress
	}
	newUserAgent := userAgent
	if newUserAgent == "" {
		newUserAgent = rt.UserAgent
	}
	if len(newUserAgent) > 512 {
		newUserAgent = newUserAgent[:512]
	}

	// Store the new refresh token in DB and revoke old token in single transaction
	newRt := &RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IPAddress: newIp,
		UserAgent: newUserAgent,
	}
	if err := s.repo.RotateRefreshTokenTx(ctx, rt.ID, newRt); err != nil {
		return nil, fmt.Errorf("failed to rotate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    s.tokenManager.GetAccessTokenTTL(),
	}, nil
}

// Logout revokes ALL refresh tokens for the user.
// It verifies that the refresh token presented belongs to the authenticated
// user (from the JWT access token in context) to prevent cross-user logout.
func (s *Service) Logout(ctx context.Context, refreshToken string, callerUserID string) error {
	rt, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) || errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		logger.Log.Errorf("Logout: database error retrieving refresh token: %v", err)
		return fmt.Errorf("failed to retrieve refresh token: %w", err)
	}
	// Ownership check: refuse to revoke another user's sessions.
	if rt.UserID != callerUserID {
		return ErrUnauthorized
	}
	// Revoke ALL sessions for this user for maximum security.
	return s.repo.RevokeAllUserTokens(ctx, rt.UserID)
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

// SetupMFA initiates MFA setup for a user.
// The TOTP secret is encrypted with AES-256-GCM before being stored so that
// a database compromise does not expose raw secrets.
func (s *Service) SetupMFA(ctx context.Context, userID string) (*MFASetupResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.MFAEnabled {
		return nil, ErrMFAAlreadyEnabled
	}

	// Generate plaintext secret
	secret, err := s.mfaService.GenerateSecret()
	if err != nil {
		return nil, err
	}

	// Generate QR code URL using the plaintext secret
	qrURL := s.mfaService.GenerateQRCodeURL(user.Email, secret)

	// Encrypt and store the TOTP secret.
	encryptedSecret, err := crypto.EncryptString(secret)
	if err != nil {
		return nil, err
	}
	if err := s.repo.SaveTempMFASecret(ctx, userID, encryptedSecret); err != nil {
		return nil, err
	}

	return &MFASetupResponse{
		Secret:    secret, // plaintext shown once to the user for manual entry
		QRCodeURL: qrURL,
	}, nil
}

// VerifyMFA verifies MFA token and enables it on successful confirmation.
// Recovery codes are generated and stored ONLY after TOTP secret is confirmed.
func (s *Service) VerifyMFA(ctx context.Context, userID, token string) ([]string, error) {
	encryptedSecret, err := s.repo.GetMFASecret(ctx, userID)
	if err != nil {
		return nil, err
	}

	secret, err := crypto.DecryptString(encryptedSecret)
	if err != nil {
		return nil, err
	}

	if !s.mfaService.VerifyTOTPUser(userID, secret, token) {
		return nil, ErrInvalidMFAToken
	}

	// Generate and store recovery codes after successful verification
	plainRecoveryCodes, err := s.mfaService.GenerateRecoveryCodes(10)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recovery codes: %w", err)
	}
	var validHashed []string
	for _, code := range plainRecoveryCodes {
		h, err := bcrypt.GenerateFromPassword([]byte(code), s.bcryptCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash recovery code: %w", err)
		}
		validHashed = append(validHashed, string(h))
	}
	if err := s.repo.CreateRecoveryCodes(ctx, userID, validHashed); err != nil {
		return nil, fmt.Errorf("failed to save recovery codes: %w", err)
	}

	// Token is correct, enable MFA officially
	if err := s.repo.ConfirmMFA(ctx, userID); err != nil {
		return nil, err
	}
	return plainRecoveryCodes, nil
}

// RequestPasswordReset creates a password reset token.
// To prevent email flooding attacks, at most 3 reset requests are allowed
// per email address in a 15-minute window.
func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		// Don't reveal whether the email exists — always return success.
		return nil
	}

	// Rate limit: max 3 password reset requests per email in 15 minutes.
	since15m := time.Now().Add(-15 * time.Minute)
	recentResets, err := s.repo.GetRecentPasswordResets(ctx, user.ID, since15m)
	if err != nil {
		return err
	}
	if recentResets >= 3 {
		logger.Log.Warnf("SECURITY AUDIT: Password reset rate limit hit for email %s (user %s)", email, user.ID)
		// Silently drop the request — do not reveal the throttle to the caller.
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

	// Log the reset request (no sensitive data in log)
	logger.Log.Infof("PASSWORD RESET REQUEST for user %s", user.ID)

	if s.emailSender != nil {
		// SECURITY: SendPasswordReset receives the unhashed single-use token to generate the reset URL.
		// Implementation of EmailSender MUST NOT log the raw token parameter under any log level.
		if err := s.emailSender.SendPasswordReset(ctx, user.Email, token); err != nil {
			logger.Log.Errorf("failed to send password reset email to %s: %v", user.Email, err)
		}
	}

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

	user, err := s.repo.GetUserByID(ctx, prt.UserID)
	if err != nil {
		return err
	}
	if !user.IsActive {
		return ErrUserInactive
	}

	// Check password history (prevent reuse of last 5)
	history, err := s.repo.GetPasswordHistory(ctx, prt.UserID)
	if err != nil {
		return fmt.Errorf("failed to fetch password history: %w", err)
	}
	for _, oldHash := range history {
		if bcrypt.CompareHashAndPassword([]byte(oldHash), []byte(newPassword)) == nil {
			return ErrPasswordReused
		}
	}

	// Hash new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), s.bcryptCost)
	if err != nil {
		return err
	}

	// Atomically update password, burn reset token, and add to password history in a single DB transaction
	if err := s.repo.ResetPasswordTx(ctx, prt.UserID, string(passwordHash), prt.ID); err != nil {
		return err
	}

	// Revoke all refresh tokens for security
	if err := s.repo.RevokeAllUserTokens(ctx, prt.UserID); err != nil {
		logger.Log.Warnf("ResetPassword: failed to revoke user tokens for user %s: %v", prt.UserID, err)
	}

	return nil
}

// isPasswordExpired checks if a user's password has expired (90 days for privileged roles).
func isPasswordExpired(user *User) bool {
	expired, _ := isPasswordExpiredReason(user)
	return expired
}

func isPasswordExpiredReason(user *User) (bool, error) {
	if user == nil {
		return false, nil
	}
	if user.Role == RoleSuperAdmin || user.Role == RoleAdmin || user.Role == RoleSecretary || user.Role == RoleTeacher || user.Role == RoleCoordinator || user.Role == RoleSystemAuditor || user.Role == RolePrincipal || user.Role == RoleVicePrincipal {
		if user.PasswordChangedAt != nil {
			if time.Since(*user.PasswordChangedAt) > 90*24*time.Hour {
				return true, ErrPasswordExpired
			}
			return false, nil
		}
		// Legacy accounts created prior to PasswordChangedAt addition: do not lock out automatically
		return false, nil
	}
	return false, nil
}

// ChangePassword changes password for an authenticated user
func (s *Service) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if !user.IsActive {
		return ErrUserInactive
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return ErrInvalidCredentials
	}

	// Validate new password policy
	passwordValidator := NewPasswordValidator()
	if err := passwordValidator.Validate(newPassword); err != nil {
		return err
	}

	// Check password history (prevent reuse of last 5)
	history, err := s.repo.GetPasswordHistory(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch password history: %w", err)
	}
	for _, oldHash := range history {
		if bcrypt.CompareHashAndPassword([]byte(oldHash), []byte(newPassword)) == nil {
			return ErrPasswordReused
		}
	}

	// Hash new password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), s.bcryptCost)
	if err != nil {
		return err
	}

	if err := s.repo.ChangePasswordTx(ctx, userID, string(passwordHash)); err != nil {
		return err
	}

	if err := s.repo.RevokeAllUserTokens(ctx, userID); err != nil {
		logger.Log.Warnf("ChangePassword: failed to revoke user tokens for user %s: %v", userID, err)
	}
	return nil
}

func (s *Service) RecordFailedAttempt(ctx context.Context, email, ipAddress string) {
	s.recordFailedAttempt(ctx, email, ipAddress)
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
