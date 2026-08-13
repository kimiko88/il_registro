package schoolsettings

import (
	"context"
	"database/sql"
	"time"
)

type Repository interface {
	GetBySchoolID(ctx context.Context, schoolID string) (*SchoolSettings, error)
	Upsert(ctx context.Context, settings *SchoolSettings) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetBySchoolID(ctx context.Context, schoolID string) (*SchoolSettings, error) {
	query := `
		SELECT school_id, require_principal_approval_for_notes, allow_parents_view_grades,
		       allow_students_view_class_averages, require_mfa_for_staff,
		       lock_scrutiny_editing_after_validation, enable_substitute_notifications, updated_at
		FROM school_settings
		WHERE school_id = $1
	`
	var s SchoolSettings
	err := r.db.QueryRowContext(ctx, query, schoolID).Scan(
		&s.SchoolID, &s.RequirePrincipalApprovalForNotes, &s.AllowParentsViewGrades,
		&s.AllowStudentsViewClassAverages, &s.RequireMFAForStaff,
		&s.LockScrutinyEditingAfterValidation, &s.EnableSubstituteNotifications, &s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return default settings if not explicitly saved yet
			return &SchoolSettings{
				SchoolID:                           schoolID,
				RequirePrincipalApprovalForNotes:   false,
				AllowParentsViewGrades:             true,
				AllowStudentsViewClassAverages:     true,
				RequireMFAForStaff:                 false,
				LockScrutinyEditingAfterValidation: true,
				EnableSubstituteNotifications:      true,
				UpdatedAt:                          time.Now(),
			}, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *repository) Upsert(ctx context.Context, s *SchoolSettings) error {
	query := `
		INSERT INTO school_settings (
			school_id, require_principal_approval_for_notes, allow_parents_view_grades,
			allow_students_view_class_averages, require_mfa_for_staff,
			lock_scrutiny_editing_after_validation, enable_substitute_notifications, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		ON CONFLICT (school_id) DO UPDATE SET
			require_principal_approval_for_notes = EXCLUDED.require_principal_approval_for_notes,
			allow_parents_view_grades = EXCLUDED.allow_parents_view_grades,
			allow_students_view_class_averages = EXCLUDED.allow_students_view_class_averages,
			require_mfa_for_staff = EXCLUDED.require_mfa_for_staff,
			lock_scrutiny_editing_after_validation = EXCLUDED.lock_scrutiny_editing_after_validation,
			enable_substitute_notifications = EXCLUDED.enable_substitute_notifications,
			updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query,
		s.SchoolID, s.RequirePrincipalApprovalForNotes, s.AllowParentsViewGrades,
		s.AllowStudentsViewClassAverages, s.RequireMFAForStaff,
		s.LockScrutinyEditingAfterValidation, s.EnableSubstituteNotifications,
	)
	return err
}
