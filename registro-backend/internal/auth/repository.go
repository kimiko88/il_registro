package auth

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// Repository handles auth-related database operations
type Repository interface {
	// User operations
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	UpdateLastLogin(ctx context.Context, userID string) error
	SchoolExists(ctx context.Context, schoolID string) (bool, error)

	// MFA operations
	EnableMFA(ctx context.Context, userID, secret string) error
	DisableMFA(ctx context.Context, userID string) error
	GetMFASecret(ctx context.Context, userID string) (string, error)
	SaveTempMFASecret(ctx context.Context, userID, secret string) error
	ConfirmMFA(ctx context.Context, userID string) error
	ConfirmMFAAndSaveRecoveryCodesTx(ctx context.Context, userID string, hashedCodes []string) error

	// Recovery codes
	// NOTE: codes passed to CreateRecoveryCodes must already be bcrypt-hashed.
	CreateRecoveryCodes(ctx context.Context, userID string, hashedCodes []string) error
	GetRecoveryCodes(ctx context.Context, userID string) ([]*MFARecoveryCode, error)
	// UseRecoveryCode finds a recovery code for userID that matches plainCode via
	// constant-time bcrypt comparison and marks it as used. Returns ErrInvalidMFAToken
	// if no matching unused code is found.
	UseRecoveryCode(ctx context.Context, userID, plainCode string) error

	// Refresh tokens
	CreateRefreshToken(ctx context.Context, token *RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error
	RotateRefreshTokenTx(ctx context.Context, oldID string, newRt *RefreshToken) error

	// Password reset
	CreatePasswordResetToken(ctx context.Context, token *PasswordResetToken) error
	GetPasswordResetToken(ctx context.Context, token string) (*PasswordResetToken, error)
	UsePasswordResetToken(ctx context.Context, tokenID string) error
	UpdatePassword(ctx context.Context, userID, passwordHash string) error
	ResetPasswordTx(ctx context.Context, userID, passwordHash, tokenID string) error
	// GetRecentPasswordResets counts how many reset tokens have been created for
	// userID since the given time. Used to rate-limit reset emails.
	GetRecentPasswordResets(ctx context.Context, userID string, since time.Time) (int, error)

	// Password history
	GetPasswordHistory(ctx context.Context, userID string) ([]string, error)
	AddPasswordHistory(ctx context.Context, userID, passwordHash string) error
	ChangePasswordTx(ctx context.Context, userID, passwordHash string) error

	// Rate limiting
	RecordLoginAttempt(ctx context.Context, attempt *LoginAttempt) error
	// GetRecentLoginAttempts counts failed attempts for a specific (email, IP) pair since `since`.
	// Used for per-IP rate limiting (max 5 in 15 min).
	GetRecentLoginAttempts(ctx context.Context, email, ipAddress string, since time.Time) (int, error)
	// GetRecentLoginAttemptsByEmail counts failed attempts for an email across ALL IPs since `since`.
	// Used for anti-IP-rotation rate limiting (max 20 in 1 hour).
	GetRecentLoginAttemptsByEmail(ctx context.Context, email string, since time.Time) (int, error)
}

type repository struct {
	db *sql.DB
}

// NewRepository creates a new auth repository
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, email, password_hash, first_name, last_name, role, school_id, is_active, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true
	user.EmailVerified = false

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Role, user.SchoolID, user.IsActive, user.EmailVerified,
		user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, school_id, 
		       is_active, COALESCE(is_staff, false), email_verified, mfa_enabled, mfa_secret, created_at, updated_at, last_login, password_changed_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	user := &User{}
	var mfaSecret sql.NullString
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Role, &user.SchoolID, &user.IsActive, &user.IsStaff, &user.EmailVerified,
		&user.MFAEnabled, &mfaSecret, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.PasswordChangedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	user.MFASecret = mfaSecret.String
	return user, nil
}

func (r *repository) GetUserByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, school_id,
		       is_active, COALESCE(is_staff, false), email_verified, mfa_enabled, mfa_secret, created_at, updated_at, last_login, password_changed_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	user := &User{}
	var mfaSecret sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Role, &user.SchoolID, &user.IsActive, &user.IsStaff, &user.EmailVerified,
		&user.MFAEnabled, &mfaSecret, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.PasswordChangedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	user.MFASecret = mfaSecret.String
	return user, nil
}

func (r *repository) UpdateUser(ctx context.Context, user *User) error {
	query := `
		UPDATE users 
		SET first_name = $1, last_name = $2, updated_at = $3
		WHERE id = $4
	`
	user.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query, user.FirstName, user.LastName, user.UpdatedAt, user.ID)
	return err
}

func (r *repository) UpdateLastLogin(ctx context.Context, userID string) error {
	query := `UPDATE users SET last_login = $1 WHERE id = $2`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, userID)
	return err
}

func (r *repository) EnableMFA(ctx context.Context, userID, secret string) error {
	query := `UPDATE users SET mfa_enabled = true, mfa_secret = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, secret, userID)
	return err
}

func (r *repository) DisableMFA(ctx context.Context, userID string) error {
	query := `UPDATE users SET mfa_enabled = false, mfa_secret = NULL WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *repository) GetMFASecret(ctx context.Context, userID string) (string, error) {
	query := `SELECT mfa_secret FROM users WHERE id = $1 AND deleted_at IS NULL`
	var secret sql.NullString
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&secret)
	if err == sql.ErrNoRows || !secret.Valid || secret.String == "" {
		return "", ErrMFANotEnabled
	}
	return secret.String, err
}

func (r *repository) SaveTempMFASecret(ctx context.Context, userID, secret string) error {
	query := `UPDATE users SET mfa_enabled = false, mfa_secret = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, secret, userID)
	return err
}

func (r *repository) ConfirmMFA(ctx context.Context, userID string) error {
	query := `UPDATE users SET mfa_enabled = true WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// ConfirmMFAAndSaveRecoveryCodesTx enables MFA and stores pre-hashed recovery codes atomically in a single transaction.
func (r *repository) ConfirmMFAAndSaveRecoveryCodesTx(ctx context.Context, userID string, hashedCodes []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	queryMFA := `UPDATE users SET mfa_enabled = true WHERE id = $1 AND deleted_at IS NULL`
	if _, err := tx.ExecContext(ctx, queryMFA, userID); err != nil {
		return err
	}

	queryCode := `INSERT INTO mfa_recovery_codes (id, user_id, code, used, created_at) VALUES ($1, $2, $3, false, $4)`
	for _, code := range hashedCodes {
		if _, err := tx.ExecContext(ctx, queryCode, uuid.New().String(), userID, code, time.Now()); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *repository) CreateRecoveryCodes(ctx context.Context, userID string, hashedCodes []string) error {
	query := `INSERT INTO mfa_recovery_codes (id, user_id, code, used, created_at) VALUES ($1, $2, $3, false, $4)`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	for _, code := range hashedCodes {
		if _, err := tx.ExecContext(ctx, query, uuid.New().String(), userID, code, time.Now()); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *repository) GetRecoveryCodes(ctx context.Context, userID string) ([]*MFARecoveryCode, error) {
	query := `SELECT id, user_id, code, used, used_at, created_at FROM mfa_recovery_codes WHERE user_id = $1 AND used = false`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []*MFARecoveryCode
	for rows.Next() {
		code := &MFARecoveryCode{}
		if err := rows.Scan(&code.ID, &code.UserID, &code.Code, &code.Used, &code.UsedAt, &code.CreatedAt); err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, rows.Err()
}

// UseRecoveryCode compares plainCode against every unused bcrypt-hashed code
// for the user using constant-time comparison, then marks the matching code
// as used. This prevents timing-based enumeration of valid codes.
func (r *repository) UseRecoveryCode(ctx context.Context, userID, plainCode string) error {
	codes, err := r.GetRecoveryCodes(ctx, userID)
	if err != nil {
		return err
	}

	for _, rc := range codes {
		if rc.Used {
			continue
		}
		// bcrypt.CompareHashAndPassword is constant-time against timing attacks.
		if bcrypt.CompareHashAndPassword([]byte(rc.Code), []byte(plainCode)) == nil {
			// Found the matching code — mark it as used.
			query := `UPDATE mfa_recovery_codes SET used = true, used_at = $1 WHERE id = $2`
			_, err := r.db.ExecContext(ctx, query, time.Now(), rc.ID)
			return err
		}
	}

	return ErrInvalidMFAToken
}

func (r *repository) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at, revoked, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	token.ID = uuid.New().String()
	token.CreatedAt = time.Now()
	token.Revoked = false

	hashed := hashToken(token.Token)
	_, err := r.db.ExecContext(ctx, query,
		token.ID, token.UserID, hashed, token.ExpiresAt,
		token.CreatedAt, token.Revoked, token.IPAddress, token.UserAgent,
	)
	return err
}

func (r *repository) GetRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, revoked, ip_address, user_agent
		FROM refresh_tokens
		WHERE token = $1
	`
	rt := &RefreshToken{}
	hashed := hashToken(token)
	err := r.db.QueryRowContext(ctx, query, hashed).Scan(
		&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt,
		&rt.CreatedAt, &rt.Revoked, &rt.IPAddress, &rt.UserAgent,
	)
	if err == sql.ErrNoRows {
		return nil, ErrInvalidToken
	}
	return rt, err
}

func (r *repository) RevokeRefreshToken(ctx context.Context, tokenID string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, tokenID)
	return err
}

func (r *repository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *repository) RotateRefreshTokenTx(ctx context.Context, oldID string, newRt *RefreshToken) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	newRt.ID = uuid.New().String()
	newRt.CreatedAt = time.Now()
	hashed := hashToken(newRt.Token)

	if _, err := tx.ExecContext(ctx, query,
		newRt.ID, newRt.UserID, hashed, newRt.ExpiresAt, newRt.IPAddress, newRt.UserAgent, newRt.CreatedAt,
	); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `UPDATE refresh_tokens SET revoked = true WHERE id = $1`, oldID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) CreatePasswordResetToken(ctx context.Context, token *PasswordResetToken) error {
	// Invalidate any existing active reset tokens for this user before creating a new one
	_, _ = r.db.ExecContext(ctx, `UPDATE password_reset_tokens SET used = true WHERE user_id = $1 AND used = false`, token.UserID)

	query := `
		INSERT INTO password_reset_tokens (id, user_id, token, expires_at, used, created_at)
		VALUES ($1, $2, $3, $4, false, $5)
	`
	token.ID = uuid.New().String()
	token.CreatedAt = time.Now()
	token.Used = false

	hashed := hashToken(token.Token)
	_, err := r.db.ExecContext(ctx, query,
		token.ID, token.UserID, hashed, token.ExpiresAt, token.CreatedAt,
	)
	return err
}

func (r *repository) GetPasswordResetToken(ctx context.Context, token string) (*PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, used, created_at
		FROM password_reset_tokens
		WHERE token = $1 AND used = false AND expires_at > NOW()
	`
	prt := &PasswordResetToken{}
	hashed := hashToken(token)
	err := r.db.QueryRowContext(ctx, query, hashed).Scan(
		&prt.ID, &prt.UserID, &prt.Token, &prt.ExpiresAt, &prt.Used, &prt.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrInvalidToken
	}
	return prt, err
}

func (r *repository) UsePasswordResetToken(ctx context.Context, tokenID string) error {
	query := `UPDATE password_reset_tokens SET used = true WHERE id = $1 AND used = false`
	res, err := r.db.ExecContext(ctx, query, tokenID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrInvalidToken
	}
	return nil
}

func (r *repository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = $2, password_changed_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, passwordHash, time.Now(), userID)
	return err
}

func (r *repository) ResetPasswordTx(ctx context.Context, userID, passwordHash, tokenID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	now := time.Now()

	updatePassQuery := `UPDATE users SET password_hash = $1, updated_at = $2, password_changed_at = $2 WHERE id = $3`
	if _, err := tx.ExecContext(ctx, updatePassQuery, passwordHash, now, userID); err != nil {
		return err
	}

	useTokenQuery := `UPDATE password_reset_tokens SET used = true, used_at = $1 WHERE id = $2 AND used = false`
	res, err := tx.ExecContext(ctx, useTokenQuery, now, tokenID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrInvalidToken
	}

	addHistoryQuery := `INSERT INTO user_password_history (user_id, password_hash) VALUES ($1, $2)`
	if _, err := tx.ExecContext(ctx, addHistoryQuery, userID, passwordHash); err != nil {
		return fmt.Errorf("failed to record password history: %w", err)
	}

	return tx.Commit()
}

// GetRecentPasswordResets counts password_reset_tokens created for userID after `since`.
// Used to rate-limit reset email requests (max 3 in 15 minutes).
func (r *repository) GetRecentPasswordResets(ctx context.Context, userID string, since time.Time) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM password_reset_tokens
		WHERE user_id = $1 AND created_at > $2
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, userID, since).Scan(&count)
	return count, err
}

func (r *repository) RecordLoginAttempt(ctx context.Context, attempt *LoginAttempt) error {
	query := `
		INSERT INTO login_attempts (email, ip_address, success, attempted_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, attempt.Email, attempt.IPAddress, attempt.Success, attempt.AttemptedAt)
	return err
}

// GetRecentLoginAttempts counts failed attempts for a specific (email, IP) pair.
// Used for per-IP burst protection (max 5 in 15 minutes).
func (r *repository) GetRecentLoginAttempts(ctx context.Context, email, ipAddress string, since time.Time) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM login_attempts
		WHERE LOWER(email) = LOWER($1) AND ip_address = $2 AND success = false AND attempted_at > $3
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, email, ipAddress, since).Scan(&count)
	return count, err
}

// GetRecentLoginAttemptsByEmail counts failed attempts for an email across ALL IPs.
// Used for anti-IP-rotation protection (max 20 in 1 hour).
func (r *repository) GetRecentLoginAttemptsByEmail(ctx context.Context, email string, since time.Time) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM login_attempts
		WHERE LOWER(email) = LOWER($1) AND success = false AND attempted_at > $2
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, email, since).Scan(&count)
	return count, err
}

func (r *repository) GetPasswordHistory(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT password_hash 
		FROM user_password_history 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT 5
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []string
	for rows.Next() {
		var h string
		if err := rows.Scan(&h); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, rows.Err()
}

func (r *repository) AddPasswordHistory(ctx context.Context, userID, passwordHash string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	query := `INSERT INTO user_password_history (user_id, password_hash) VALUES ($1, $2)`
	if _, err := tx.ExecContext(ctx, query, userID, passwordHash); err != nil {
		return err
	}

	pruneQuery := `
		DELETE FROM user_password_history
		WHERE user_id = $1 AND id NOT IN (
			SELECT id FROM user_password_history
			WHERE user_id = $1 ORDER BY created_at DESC LIMIT 5
		)
	`
	_, _ = tx.ExecContext(ctx, pruneQuery, userID)
	return tx.Commit()
}

func (r *repository) ChangePasswordTx(ctx context.Context, userID, passwordHash string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	now := time.Now()
	updateQuery := `UPDATE users SET password_hash = $1, password_changed_at = $2 WHERE id = $3 AND deleted_at IS NULL`
	if _, err := tx.ExecContext(ctx, updateQuery, passwordHash, now, userID); err != nil {
		return err
	}

	histQuery := `INSERT INTO user_password_history (user_id, password_hash) VALUES ($1, $2)`
	if _, err := tx.ExecContext(ctx, histQuery, userID, passwordHash); err != nil {
		return err
	}

	pruneQuery := `
		DELETE FROM user_password_history
		WHERE user_id = $1 AND id NOT IN (
			SELECT id FROM user_password_history
			WHERE user_id = $1 ORDER BY created_at DESC LIMIT 5
		)
	`
	_, _ = tx.ExecContext(ctx, pruneQuery, userID)

	return tx.Commit()
}

func (r *repository) SchoolExists(ctx context.Context, schoolID string) (bool, error) {
	if schoolID == "" {
		return false, nil
	}
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM schools WHERE id = $1::uuid)`
	err := r.db.QueryRowContext(ctx, query, schoolID).Scan(&exists)
	return exists, err
}
