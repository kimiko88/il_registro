package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"registro-backend/internal/permissions"
)

var (
	ErrUnauthorized    = errors.New("unauthorized")
	ErrInvalidPassword = errors.New("invalid password")
)

type Service struct {
	repo        Repository
	permManager *permissions.Manager
	validator   *Validator
	importer    *Importer
	exporter    *Exporter
	gdpr        *GDPRHandler
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:        repo,
		permManager: permissions.NewManager(),
		validator:   NewValidator(),
		importer:    NewImporter(),
		exporter:    NewExporter(),
		gdpr:        NewGDPRHandler(repo),
	}
}

// CreateUser creates a new user, hashing their password
func (s *Service) CreateUser(ctx context.Context, actorRole string, req CreateUserRequest) (*User, error) {
	if !s.permManager.HasPermission(actorRole, permissions.UserCreate) {
		return nil, ErrUnauthorized
	}

	if req.FiscalCode != "" && !s.validator.ValidateFiscalCode(req.FiscalCode) {
		return nil, fmt.Errorf("invalid fiscal code")
	}
	if !s.validator.ValidatePassword(req.Password) {
		return nil, fmt.Errorf("password does not meet complexity requirements")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           uuid.New().String(),
		Email:        SanitizeEmail(req.Email),
		PasswordHash: string(hashed),
		FirstName:    SanitizeText(req.FirstName),
		LastName:     SanitizeText(req.LastName),
		FiscalCode:   SanitizeTextPtr(req.FiscalCode),
		Role:         req.Role,
		SchoolID:     req.SchoolID,
		ClassID:      req.ClassID,
		PhoneNumber:  SanitizeTextPtr(req.PhoneNumber),
		JobTitle:     SanitizeTextPtr(req.JobTitle),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	_ = s.repo.AddPasswordHistory(ctx, user.ID, user.PasswordHash)

	// Audit
	_ = s.repo.LogAudit(ctx, &AuditLog{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		ActorID:   "system", // In real app, pass actor ID
		Action:    "create_user",
		Details:   fmt.Sprintf("Created user %s (%s)", user.Email, user.Role),
		CreatedAt: time.Now(),
	})

	return user, nil
}

func (s *Service) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	return s.repo.IsGuardian(ctx, parentID, studentID)
}

func (s *Service) GetChildren(ctx context.Context, parentUserID string) ([]StudentChild, error) {
	return s.repo.GetChildren(ctx, parentUserID)
}



func (s *Service) GetUser(ctx context.Context, actorRole string, id string) (*User, error) {
	if !s.permManager.HasPermission(actorRole, permissions.UserRead) {
		return nil, ErrUnauthorized
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListUsers(ctx context.Context, actorRole string, filter UserFilter) ([]User, int, error) {
	if !s.permManager.HasPermission(actorRole, permissions.UserRead) {
		return nil, 0, ErrUnauthorized
	}

	// Filter out higher privileged roles for non-superadmin users
	if actorRole != "superadmin" {
		filter.ExcludeRoles = append(filter.ExcludeRoles, "superadmin")
	}
	if actorRole == "secretary" {
		filter.ExcludeRoles = append(filter.ExcludeRoles, "admin")
	}

	// Improve: Restrict filter based on role (e.g. principal can only see their school)
	// For now, allow full list based on broad permission
	return s.repo.List(ctx, filter)
}

func (s *Service) UpdateUser(ctx context.Context, actorRole string, id string, req UpdateUserRequest) (*User, error) {
	if !s.permManager.HasPermission(actorRole, permissions.UserUpdate) {
		return nil, ErrUnauthorized
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = SanitizeText(*req.FirstName)
	}
	if req.LastName != nil {
		user.LastName = SanitizeText(*req.LastName)
	}
	if req.PhoneNumber != nil {
		user.PhoneNumber = SanitizeTextPtr(*req.PhoneNumber)
	}
	if req.JobTitle != nil {
		user.JobTitle = SanitizeTextPtr(*req.JobTitle)
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.SchoolID != nil {
		user.SchoolID = req.SchoolID
	}
	if req.ClassID != nil {
		user.ClassID = req.ClassID // Will be handled by repo Update
	}
	if req.FiscalCode != nil {
		if *req.FiscalCode != "" && !s.validator.ValidateFiscalCode(*req.FiscalCode) {
			return nil, fmt.Errorf("invalid fiscal code")
		}
		user.FiscalCode = SanitizeTextPtr(*req.FiscalCode)
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) DeleteUser(ctx context.Context, actorRole, id string) error {
	if !s.permManager.HasPermission(actorRole, permissions.UserDelete) {
		return ErrUnauthorized
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) BulkDeleteUsers(ctx context.Context, actorRole string, ids []string) (int, error) {
	if !s.permManager.HasPermission(actorRole, permissions.UserDelete) {
		return 0, ErrUnauthorized
	}
	return s.repo.BulkDelete(ctx, ids)
}

func (s *Service) RestoreUser(ctx context.Context, actorRole, id string) error {
	if !s.permManager.HasPermission(actorRole, permissions.UserDelete) { // Usually strictly admin
		return ErrUnauthorized
	}
	return s.repo.Restore(ctx, id)
}

func (s *Service) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return ErrInvalidPassword
	}

	if !s.validator.ValidatePassword(req.NewPassword) {
		return fmt.Errorf("new password too weak")
	}

	// Check password history (prevent reuse of last 5)
	history, err := s.repo.GetPasswordHistory(ctx, userID)
	if err == nil {
		for _, oldHash := range history {
			if bcrypt.CompareHashAndPassword([]byte(oldHash), []byte(req.NewPassword)) == nil {
				return errors.New("la nuova password non può essere uguale ad una delle ultime 5 utilizzate")
			}
		}
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	user.PasswordHash = string(hashed)
	now := time.Now()
	user.PasswordChangedAt = &now

	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	_ = s.repo.AddPasswordHistory(ctx, userID, user.PasswordHash)
	return nil
}

func (s *Service) ResetPassword(ctx context.Context, actorRole, userID, newPassword string) error {
	if !s.permManager.HasPermission(actorRole, permissions.UserUpdate) {
		return ErrUnauthorized
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if !s.validator.ValidatePassword(newPassword) {
		return fmt.Errorf("new password too weak")
	}

	// Check password history (prevent reuse of last 5)
	history, err := s.repo.GetPasswordHistory(ctx, userID)
	if err == nil {
		for _, oldHash := range history {
			if bcrypt.CompareHashAndPassword([]byte(oldHash), []byte(newPassword)) == nil {
				return errors.New("la nuova password non può essere uguale ad una delle ultime 5 utilizzate")
			}
		}
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	user.PasswordHash = string(hashed)
	now := time.Now()
	user.PasswordChangedAt = &now

	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	_ = s.repo.AddPasswordHistory(ctx, userID, user.PasswordHash)
	return nil
}

// BulkImport handles CSV/XLSX
func (s *Service) BulkImport(ctx context.Context, actorRole string, file multipart.File, filename string) (*ImportResult, error) {
	if !s.permManager.HasPermission(actorRole, permissions.UserImport) {
		return nil, ErrUnauthorized
	}

	users, parseErrs, err := s.importer.ParseFile(file, filename)
	if err != nil {
		return nil, err
	}

	count, dbErrs, err := s.repo.BulkCreate(ctx, users)
	if err != nil {
		return nil, err
	}

	result := &ImportResult{
		Total:   len(users),
		Created: count,
		Errors:  append(parseErrs, dbErrs...),
		Failed:  len(users) - count,
	}
	return result, nil
}

// ExportUsers generates a file
func (s *Service) ExportUsers(ctx context.Context, actorRole string, filter UserFilter, format string) ([]byte, error) {
	if !s.permManager.HasPermission(actorRole, permissions.UserExport) {
		return nil, ErrUnauthorized
	}

	users, _, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Transform to map for generic exporter
	var data []map[string]interface{}
	for _, u := range users {
		data = append(data, map[string]interface{}{
			"id":         u.ID,
			"first_name": u.FirstName,
			"last_name":  u.LastName,
			"email":      u.Email,
			"role":       u.Role,
			"active":     u.IsActive,
			"created_at": u.CreatedAt,
		})
	}

	switch format {
	case "csv":
		return s.exporter.ToCSV(data)
	case "xlsx":
		return s.exporter.ToXLSX(data)
	case "json":
		return s.exporter.ToJSON(data)
	default:
		return nil, fmt.Errorf("unsupported format")
	}
}

// DisableMFA - Admin override
func (s *Service) DisableMFA(ctx context.Context, actorRole, userID string) error {
	if !s.permManager.HasPermission(actorRole, permissions.UserUpdate) {
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

// GDPR Exports
func (s *Service) GDPRDataExport(ctx context.Context, actorID, actorRole, userID string) (map[string]interface{}, error) {
	// User can export their own data, or admin/DPO with permission
	if actorID != userID && !s.permManager.HasPermission(actorRole, permissions.UserExport) {
		return nil, ErrUnauthorized
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	logs, _, _ := s.repo.GetAuditLogs(ctx, userID, 1000, 0)
	return s.gdpr.GenerateDataExport(user, logs), nil
}

func (s *Service) GDPRDelete(ctx context.Context, actorRole, userID string) error {
	if !s.permManager.HasPermission(actorRole, permissions.UserDelete) {
		return ErrUnauthorized
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	s.gdpr.PseudonymizeUser(user)

	return s.repo.Update(ctx, user)
}
func (s *Service) GetGuardians(ctx context.Context, actorRole, studentUserID string) ([]GuardianInfo, error) {
	if !s.permManager.HasPermission(actorRole, permissions.UserRead) {
		return nil, ErrUnauthorized
	}
	studentProfileID, err := s.repo.GetStudentProfile(ctx, studentUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []GuardianInfo{}, nil
		}
		return nil, err
	}
	return s.repo.GetGuardians(ctx, studentProfileID)
}

func (s *Service) AddGuardian(ctx context.Context, actorRole, studentUserID, parentUserID, relationship string) error {
	if !s.permManager.HasPermission(actorRole, permissions.UserUpdate) {
		return ErrUnauthorized
	}
	studentProfileID, err := s.repo.GetStudentProfile(ctx, studentUserID)
	if err != nil {
		return err
	}
	parentProfileID, err := s.repo.GetParentProfile(ctx, parentUserID)
	if err != nil {
		return err
	}
	return s.repo.AddGuardian(ctx, studentProfileID, parentProfileID, relationship)
}

func (s *Service) RemoveGuardian(ctx context.Context, actorRole, studentUserID, parentUserID string) error {
	if !s.permManager.HasPermission(actorRole, permissions.UserUpdate) {
		return ErrUnauthorized
	}
	studentProfileID, err := s.repo.GetStudentProfile(ctx, studentUserID)
	if err != nil {
		return err
	}
	parentProfileID, err := s.repo.GetParentProfile(ctx, parentUserID)
	if err != nil {
		return err
	}
	return s.repo.RemoveGuardian(ctx, studentProfileID, parentProfileID)
}
