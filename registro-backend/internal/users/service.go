package users

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

// allowedCreators maps each role to the set of roles it is allowed to create.
var allowedCreators = map[string]map[string]bool{
	"superadmin": {"superadmin": true, "admin": true, "secretary": true, "teacher": true, "student": true, "parent": true},
	"admin":      {"secretary": true, "teacher": true, "student": true, "parent": true},
	"secretary":  {"student": true, "parent": true},
}

// Service handles business logic for user management.
type Service struct {
	repo Repository
}

// NewService creates a new Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ─── Auth helpers ────────────────────────────────────────────────────────────

func isPrivileged(role string) bool {
	return role == "admin" || role == "superadmin" || role == "secretary"
}

// ─── CRUD ────────────────────────────────────────────────────────────────────

// CreateUser creates a new user, enforcing role-based permission checks.
func (s *Service) CreateUser(ctx context.Context, actorRole string, req CreateUserRequest) (*User, error) {
	allowed, ok := allowedCreators[actorRole]
	if !ok {
		return nil, ErrUnauthorized
	}
	if !allowed[req.Role] {
		return nil, ErrUnauthorized
	}

	if len(req.Password) < 8 {
		return nil, fmt.Errorf("password must be at least 8 characters long")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	var fiscalCode *string
	if req.FiscalCode != "" {
		fc := req.FiscalCode
		fiscalCode = &fc
	}
	var phoneNumber *string
	if req.PhoneNumber != "" {
		ph := req.PhoneNumber
		phoneNumber = &ph
	}
	var jobTitle *string
	if req.JobTitle != "" {
		jt := req.JobTitle
		jobTitle = &jt
	}

	now := time.Now()
	user := &User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: string(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		FiscalCode:   fiscalCode,
		Role:         req.Role,
		SchoolID:     req.SchoolID,
		ClassID:      req.ClassID,
		PhoneNumber:  phoneNumber,
		JobTitle:     jobTitle,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// ListUsers returns a paginated, filtered list of users.
func (s *Service) ListUsers(ctx context.Context, actorRole string, filter UserFilter) ([]User, int, error) {
	if !isPrivileged(actorRole) {
		return nil, 0, ErrUnauthorized
	}
	return s.repo.List(ctx, filter)
}

// GetUser returns a single user by ID.
func (s *Service) GetUser(ctx context.Context, actorRole string, id string) (*User, error) {
	if !isPrivileged(actorRole) && actorRole != "teacher" {
		return nil, ErrUnauthorized
	}
	return s.repo.GetByID(ctx, id)
}

// UpdateUser updates an existing user's fields.
func (s *Service) UpdateUser(ctx context.Context, actorRole string, id string, req UpdateUserRequest) (*User, error) {
	if !isPrivileged(actorRole) {
		return nil, ErrUnauthorized
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.PhoneNumber != nil {
		user.PhoneNumber = req.PhoneNumber
	}
	if req.JobTitle != nil {
		user.JobTitle = req.JobTitle
	}
	if req.ClassID != nil {
		user.ClassID = req.ClassID
	}
	if req.SchoolID != nil {
		user.SchoolID = req.SchoolID
	}
	if req.FiscalCode != nil {
		user.FiscalCode = req.FiscalCode
	}

	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser soft-deletes a user.
func (s *Service) DeleteUser(ctx context.Context, actorRole string, id string) error {
	if !isPrivileged(actorRole) {
		return ErrUnauthorized
	}
	return s.repo.Delete(ctx, id)
}

// BulkDeleteUsers soft-deletes multiple users at once.
func (s *Service) BulkDeleteUsers(ctx context.Context, actorRole string, ids []string) (int, error) {
	if !isPrivileged(actorRole) {
		return 0, ErrUnauthorized
	}
	return s.repo.BulkDelete(ctx, ids)
}

// RestoreUser un-deletes a soft-deleted user.
func (s *Service) RestoreUser(ctx context.Context, actorRole string, id string) error {
	if !isPrivileged(actorRole) {
		return ErrUnauthorized
	}
	return s.repo.Restore(ctx, id)
}

// ─── Password ────────────────────────────────────────────────────────────────

// ChangePassword allows a user to change their own password.
// It enforces that the new password differs from the last 5 used passwords.
func (s *Service) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Field is CurrentPassword in ChangePasswordRequest DTO
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Check password history
	history, err := s.repo.GetPasswordHistory(ctx, userID)
	if err != nil {
		return err
	}
	for _, old := range history {
		if bcrypt.CompareHashAndPassword([]byte(old), []byte(req.NewPassword)) == nil {
			return errors.New("new password must not match any of the last 5 passwords")
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	user.PasswordHash = string(hash)
	user.PasswordChangedAt = &now
	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}
	return s.repo.AddPasswordHistory(ctx, userID, string(hash))
}

// ResetPassword allows an admin to force-reset a user's password.
func (s *Service) ResetPassword(ctx context.Context, actorRole string, userID string, newPassword string) error {
	if actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	now := time.Now()
	user.PasswordHash = string(hash)
	user.PasswordChangedAt = &now
	return s.repo.Update(ctx, user)
}

// DisableMFA disables MFA for the given user (admin only).
func (s *Service) DisableMFA(ctx context.Context, actorRole string, userID string) error {
	if actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	user.MFAEnabled = false
	user.MFASecret = ""
	return s.repo.Update(ctx, user)
}

// ─── Bulk Import ─────────────────────────────────────────────────────────────

// BulkImport imports users from a CSV or XLSX file.
// Returns an ImportResult (the type defined in dto.go).
func (s *Service) BulkImport(ctx context.Context, actorRole string, file multipart.File, filename string) (*ImportResult, error) {
	if !isPrivileged(actorRole) {
		return nil, ErrUnauthorized
	}

	// NewImporter takes no arguments
	importer := NewImporter()
	users, parseErrors, err := importer.ParseFile(file, filename)
	if err != nil {
		return nil, err
	}

	result := &ImportResult{
		Total:  len(users) + len(parseErrors),
		Failed: len(parseErrors),
		Errors: parseErrors,
	}

	if len(users) > 0 {
		count, bulkErrors, err := s.repo.BulkCreate(ctx, users)
		result.Created = count
		result.Failed += len(bulkErrors)
		result.Errors = append(result.Errors, bulkErrors...)
		if err != nil && len(bulkErrors) == 0 {
			result.Errors = append(result.Errors, err.Error())
			result.Failed += len(users)
		}
	}

	return result, nil
}

// ─── GDPR ────────────────────────────────────────────────────────────────────

// GDPRDataExport exports all personal data for a user (GDPR right of access).
func (s *Service) GDPRDataExport(ctx context.Context, actorID, actorRole, targetID string) (map[string]interface{}, error) {
	if actorID != targetID && actorRole != "admin" && actorRole != "superadmin" {
		return nil, ErrUnauthorized
	}
	user, err := s.repo.GetByID(ctx, targetID)
	if err != nil {
		return nil, err
	}
	// Fetch audit logs for the export (max 1000, offset 0)
	logs, _, err := s.repo.GetAuditLogs(ctx, targetID, 1000, 0)
	if err != nil {
		return nil, err
	}
	// NewExporter takes no arguments; GenerateDataExport is on GDPRHandler
	gdpr := NewGDPRHandler(s.repo)
	return gdpr.GenerateDataExport(user, logs), nil
}

// GDPRDelete pseudonymizes a user (GDPR right to erasure).
func (s *Service) GDPRDelete(ctx context.Context, actorRole string, targetID string) error {
	if actorRole != "admin" && actorRole != "superadmin" {
		return ErrUnauthorized
	}
	user, err := s.repo.GetByID(ctx, targetID)
	if err != nil {
		return err
	}
	gdpr := NewGDPRHandler(s.repo)
	gdpr.PseudonymizeUser(user)
	// Persist the pseudonymized data
	return s.repo.Update(ctx, user)
}

// ─── Audit ───────────────────────────────────────────────────────────────────

// GetAuditLogs returns the audit log for a specific user (paginated).
// Access: the user themselves, or admin/superadmin.
func (s *Service) GetAuditLogs(ctx context.Context, actorID, actorRole, targetID string, limit, offset int) ([]AuditLog, int, error) {
	if actorID != targetID && actorRole != "admin" && actorRole != "superadmin" {
		return nil, 0, ErrUnauthorized
	}
	return s.repo.GetAuditLogs(ctx, targetID, limit, offset)
}

// ─── Relationships ───────────────────────────────────────────────────────────

// GetChildren returns the list of students linked to a parent.
func (s *Service) GetChildren(ctx context.Context, parentUserID string) ([]StudentChild, error) {
	return s.repo.GetChildren(ctx, parentUserID)
}

// GetGuardians returns all guardians linked to a student.
func (s *Service) GetGuardians(ctx context.Context, actorRole string, studentUserID string) ([]GuardianInfo, error) {
	if !isPrivileged(actorRole) {
		return nil, ErrUnauthorized
	}
	studentProfileID, err := s.repo.GetStudentProfile(ctx, studentUserID)
	if err != nil {
		return nil, fmt.Errorf("student profile not found: %w", err)
	}
	return s.repo.GetGuardians(ctx, studentProfileID)
}

// AddGuardian links a parent to a student.
func (s *Service) AddGuardian(ctx context.Context, actorRole, studentUserID, parentUserID, relationship string) error {
	if !isPrivileged(actorRole) {
		return ErrUnauthorized
	}
	studentProfileID, err := s.repo.GetStudentProfile(ctx, studentUserID)
	if err != nil {
		return fmt.Errorf("student profile not found: %w", err)
	}
	parentProfileID, err := s.repo.GetParentProfile(ctx, parentUserID)
	if err != nil {
		return fmt.Errorf("parent profile not found: %w", err)
	}
	return s.repo.AddGuardian(ctx, studentProfileID, parentProfileID, relationship)
}

// RemoveGuardian unlinks a parent from a student.
func (s *Service) RemoveGuardian(ctx context.Context, actorRole, studentUserID, parentUserID string) error {
	if !isPrivileged(actorRole) {
		return ErrUnauthorized
	}
	studentProfileID, err := s.repo.GetStudentProfile(ctx, studentUserID)
	if err != nil {
		return fmt.Errorf("student profile not found: %w", err)
	}
	parentProfileID, err := s.repo.GetParentProfile(ctx, parentUserID)
	if err != nil {
		return fmt.Errorf("parent profile not found: %w", err)
	}
	return s.repo.RemoveGuardian(ctx, studentProfileID, parentProfileID)
}

// IsGuardian checks if parentUserID is a guardian for studentUserID.
func (s *Service) IsGuardian(ctx context.Context, parentUserID, studentUserID string) (bool, error) {
	return s.repo.IsGuardian(ctx, parentUserID, studentUserID)
}

// ExportUsers exports users matching filter in specified format ("csv" or "json").
func (s *Service) ExportUsers(ctx context.Context, actorRole string, filter UserFilter, format string) ([]byte, error) {
	if !isPrivileged(actorRole) {
		return nil, ErrUnauthorized
	}
	users, _, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	if strings.ToLower(format) == "csv" {
		var buf bytes.Buffer
		buf.WriteString("ID,Email,FirstName,LastName,Role\n")
		for _, u := range users {
			buf.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s\n", u.ID, u.Email, u.FirstName, u.LastName, u.Role))
		}
		return buf.Bytes(), nil
	}
	return json.Marshal(users)
}
