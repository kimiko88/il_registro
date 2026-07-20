package schoolsettings

import "time"

type SchoolSettings struct {
	SchoolID                         string    `json:"school_id" db:"school_id"`
	RequirePrincipalApprovalForNotes bool      `json:"require_principal_approval_for_notes" db:"require_principal_approval_for_notes"`
	AllowParentsViewGrades           bool      `json:"allow_parents_view_grades" db:"allow_parents_view_grades"`
	AllowStudentsViewClassAverages   bool      `json:"allow_students_view_class_averages" db:"allow_students_view_class_averages"`
	RequireMFAForStaff               bool      `json:"require_mfa_for_staff" db:"require_mfa_for_staff"`
	LockScrutinyEditingAfterValidation bool    `json:"lock_scrutiny_editing_after_validation" db:"lock_scrutiny_editing_after_validation"`
	EnableSubstituteNotifications    bool      `json:"enable_substitute_notifications" db:"enable_substitute_notifications"`
	UpdatedAt                        time.Time `json:"updated_at" db:"updated_at"`
}

type UpdateSchoolSettingsRequest struct {
	RequirePrincipalApprovalForNotes *bool `json:"require_principal_approval_for_notes"`
	AllowParentsViewGrades           *bool `json:"allow_parents_view_grades"`
	AllowStudentsViewClassAverages   *bool `json:"allow_students_view_class_averages"`
	RequireMFAForStaff               *bool `json:"require_mfa_for_staff"`
	LockScrutinyEditingAfterValidation *bool `json:"lock_scrutiny_editing_after_validation"`
	EnableSubstituteNotifications    *bool `json:"enable_substitute_notifications"`
}
