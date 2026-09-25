package timetablegen

import (
	"encoding/json"
	"time"
)

// PreferenceType constants
const (
	PrefPreferred   = "preferred"
	PrefNeutral     = "neutral"
	PrefUnavailable = "unavailable"
)

// JobStatus constants
const (
	JobStatusPending   = "pending"
	JobStatusRunning   = "running"
	JobStatusCompleted = "completed"
	JobStatusFailed    = "failed"
	JobStatusCancelled = "cancelled"
)

// TeacherPreference represents teacher's desired availability
type TeacherPreference struct {
	ID             string    `json:"id" db:"id"`
	SchoolID       string    `json:"school_id" db:"school_id"`
	TeacherID      string    `json:"teacher_id" db:"teacher_id"` // user_id or teacher_id
	AcademicYearID *string   `json:"academic_year_id,omitempty" db:"academic_year_id"`
	DayOfWeek      int       `json:"day_of_week" db:"day_of_week"`
	HourIndex      int       `json:"hour_index" db:"hour_index"`
	PreferenceType string    `json:"preference_type" db:"preference_type"`
	Reason         *string   `json:"reason,omitempty" db:"reason"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// SubjectRoomRequirement maps a subject to a required bookable room type
type SubjectRoomRequirement struct {
	ID               string    `json:"id" db:"id"`
	SchoolID         string    `json:"school_id" db:"school_id"`
	SubjectID        string    `json:"subject_id" db:"subject_id"`
	SubjectName      string    `json:"subject_name,omitempty" db:"subject_name"`
	RequiredRoomType string    `json:"required_room_type" db:"required_room_type"`
	LabHours         int       `json:"lab_hours" db:"lab_hours"`
	IsMandatory      bool      `json:"is_mandatory" db:"is_mandatory"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// TimetableConstraint represents a configurable rule for schedule generation
type TimetableConstraint struct {
	ID             string          `json:"id" db:"id"`
	SchoolID       string          `json:"school_id" db:"school_id"`
	ConstraintType string          `json:"constraint_type" db:"constraint_type"`
	TargetType     *string         `json:"target_type,omitempty" db:"target_type"`
	TargetID       *string         `json:"target_id,omitempty" db:"target_id"`
	Parameters     json.RawMessage `json:"parameters" db:"parameters"`
	IsHard         bool            `json:"is_hard" db:"is_hard"`
	Priority       int             `json:"priority" db:"priority"`
	IsActive       bool            `json:"is_active" db:"is_active"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
}

// TimetableJob tracks schedule generation runs
type TimetableJob struct {
	ID             string          `json:"id" db:"id"`
	SchoolID       string          `json:"school_id" db:"school_id"`
	AcademicYearID *string         `json:"academic_year_id,omitempty" db:"academic_year_id"`
	TriggeredBy    *string         `json:"triggered_by,omitempty" db:"triggered_by"`
	Status         string          `json:"status" db:"status"`
	Algorithm      string          `json:"algorithm" db:"algorithm"`
	Parameters     json.RawMessage `json:"parameters" db:"parameters"`
	ResultSummary  json.RawMessage `json:"result_summary" db:"result_summary"`
	StartedAt      *time.Time      `json:"started_at,omitempty" db:"started_at"`
	CompletedAt    *time.Time      `json:"completed_at,omitempty" db:"completed_at"`
	ErrorMessage   *string         `json:"error_message,omitempty" db:"error_message"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
}

// DTOs
type SavePreferencesRequest struct {
	AcademicYearID *string           `json:"academic_year_id"`
	Preferences    []PreferenceEntry `json:"preferences" binding:"required"`
}

type PreferenceEntry struct {
	DayOfWeek      int     `json:"day_of_week" binding:"required,min=1,max=7"`
	HourIndex      int     `json:"hour_index" binding:"required,min=1,max=12"`
	PreferenceType string  `json:"preference_type" binding:"required"`
	Reason         *string `json:"reason"`
}

type SaveRoomRequirementRequest struct {
	SubjectID        string `json:"subject_id" binding:"required"`
	RequiredRoomType string `json:"required_room_type" binding:"required"`
	LabHours         int    `json:"lab_hours"`
	IsMandatory      bool   `json:"is_mandatory"`
}

type AssociatedGroup struct {
	ID           string   `json:"id"`
	SchoolID     string   `json:"school_id"`
	Name         string   `json:"name"`
	SubjectID    string   `json:"subject_id"`
	SubjectName  string   `json:"subject_name,omitempty"`
	TeacherID    string   `json:"teacher_id"`
	TeacherName  string   `json:"teacher_name,omitempty"`
	ClassIDs     []string `json:"class_ids"`
	ClassNames   []string `json:"class_names,omitempty"`
	HoursPerWeek int      `json:"hours_per_week"`
}

type SaveConstraintRequest struct {
	ConstraintType string          `json:"constraint_type" binding:"required"`
	TargetType     *string         `json:"target_type"`
	TargetID       *string         `json:"target_id"`
	Parameters     json.RawMessage `json:"parameters"`
	IsHard         bool            `json:"is_hard"`
	Priority       int             `json:"priority"`
	IsActive       *bool           `json:"is_active"`
}

type GenerateTimetableRequest struct {
	AcademicYearID   *string `json:"academic_year_id"`
	MaxIterations    int     `json:"max_iterations"`
	TimeLimitSeconds int     `json:"time_limit_seconds"`
}

type AdjustTimetableRequest struct {
	Slots []GeneratedSlot `json:"slots" binding:"required"`
}

type DesiderataWindowResponse struct {
	IsOpen bool `json:"is_open"`
}

type SetDesiderataWindowRequest struct {
	IsOpen bool `json:"is_open"`
}

// Generation structures
type GeneratedSlot struct {
	ClassID       string  `json:"class_id"`
	ClassName     string  `json:"class_name"`
	BuildingID    *string `json:"building_id,omitempty"`
	SubjectID     string  `json:"subject_id"`
	SubjectName   string  `json:"subject_name"`
	TeacherID     *string `json:"teacher_id,omitempty"`
	TeacherUserID *string `json:"teacher_user_id,omitempty"`
	TeacherName   string  `json:"teacher_name"`
	DayOfWeek     int     `json:"day_of_week"`
	HourIndex     int     `json:"hour_index"`
	RoomID        *string `json:"room_id,omitempty"`
	RoomName      string  `json:"room_name,omitempty"`
}

type HardConflict struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	ClassID     string `json:"class_id,omitempty"`
	TeacherID   string `json:"teacher_id,omitempty"`
	DayOfWeek   int    `json:"day_of_week"`
	HourIndex   int    `json:"hour_index"`
}

type SoftViolation struct {
	ConstraintType string  `json:"constraint_type"`
	TeacherID      *string `json:"teacher_id,omitempty"`
	ClassID        *string `json:"class_id,omitempty"`
	Description    string  `json:"description"`
	Penalty        float64 `json:"penalty"`
}

type UnassignedSlot struct {
	ClassID     string `json:"class_id"`
	ClassName   string `json:"class_name"`
	SubjectID   string `json:"subject_id"`
	SubjectName string `json:"subject_name"`
	TeacherID   string `json:"teacher_id"`
	TeacherName string `json:"teacher_name"`
	HoursNeeded int    `json:"hours_needed"`
	Reason      string `json:"reason"`
}

type TimetableGenerationResult struct {
	JobID          string           `json:"job_id"`
	TotalSlots     int              `json:"total_slots"`
	AssignedSlots  int              `json:"assigned_slots"`
	CoveragePct    float64          `json:"coverage_pct"`
	Slots          []GeneratedSlot  `json:"slots"`
	HardConflicts  []HardConflict   `json:"hard_conflicts"`
	SoftViolations []SoftViolation  `json:"soft_violations"`
	Warnings       []string         `json:"warnings"`
	Unassigned     []UnassignedSlot `json:"unassigned"`
	DurationMs     int64            `json:"duration_ms"`
}

// Input data loaded for generation
type AssignmentData struct {
	ClassID           string
	ClassName         string
	BuildingID        *string
	SubjectID         string
	SubjectName       string
	TeacherID         string
	TeacherUserID     string
	TeacherName       string
	HiringDate        *time.Time
	HoursPerWeek      int
	AssociatedGroupID *string
	IsAssociatedGroup bool
}

type RoomData struct {
	ID         string
	BuildingID *string
	Name       string
	RoomType   string
	Capacity   int
	IsActive   bool
}

// Class Daily Limits (Min and Max hours per day for class)
type ClassDailyLimit struct {
	ClassID        string `json:"class_id"`
	MinHoursPerDay int    `json:"min_hours_per_day"`
	MaxHoursPerDay int    `json:"max_hours_per_day"`
}

// Class Curriculum Plan & Subject Hours
type ClassSubjectPlanItem struct {
	ID           string  `json:"id,omitempty"`
	SubjectID    string  `json:"subject_id"`
	SubjectName  string  `json:"subject_name,omitempty"`
	HoursPerWeek float64 `json:"hours_per_week"`
	TeacherID    *string `json:"teacher_id,omitempty"`
	TeacherName  string  `json:"teacher_name,omitempty"`
}

type ClassCurriculumPlan struct {
	ClassID        string                 `json:"class_id"`
	ClassName      string                 `json:"class_name"`
	Section        string                 `json:"section"`
	AcademicYear   string                 `json:"academic_year"`
	MinHoursPerDay int                    `json:"min_hours_per_day"`
	MaxHoursPerDay int                    `json:"max_hours_per_day"`
	TotalHoursWeek float64                `json:"total_hours_week"`
	Subjects       []ClassSubjectPlanItem `json:"subjects"`
}

type SaveClassCurriculumPlanRequest struct {
	MinHoursPerDay int                    `json:"min_hours_per_day"`
	MaxHoursPerDay int                    `json:"max_hours_per_day"`
	Subjects       []ClassSubjectPlanItem `json:"subjects"`
}

type InheritPlanRequest struct {
	SourceAcademicYear string `json:"source_academic_year,omitempty"`
}

type InheritAllResult struct {
	ClassesUpdated int    `json:"classes_updated"`
	SubjectsCopied int    `json:"subjects_copied"`
	Message        string `json:"message"`
}

// Teacher Quick Preferences (Tabular representation for all teachers)
type TeacherQuickPreferenceItem struct {
	TeacherID             string `json:"teacher_id" binding:"required"`
	TeacherUserID         string `json:"teacher_user_id,omitempty"`
	TeacherName           string `json:"teacher_name,omitempty"`
	SubjectName           string `json:"subject_name,omitempty"`
	DayOff                int    `json:"day_off"`                     // 0=none, 1=Mon, 2=Tue, 3=Wed, 4=Thu, 5=Fri, 6=Sat
	TimeSlotPref          string `json:"time_slot_pref"`              // "none", "early_hours", "late_hours"
	MaxHoursPerDay        int    `json:"max_hours_per_day,omitempty"` // optional (e.g. 4, 5, 6)
	Notes                 string `json:"notes,omitempty"`
	PreferredHoursCount   int    `json:"preferred_hours_count,omitempty"`
	UnavailableHoursCount int    `json:"unavailable_hours_count,omitempty"`
}

type SaveTeacherQuickPreferencesRequest struct {
	AcademicYearID *string                      `json:"academic_year_id,omitempty"`
	Preferences    []TeacherQuickPreferenceItem `json:"preferences" binding:"required"`
}

type TeacherQuickPreferencesOverviewResponse struct {
	AcademicYearID *string                      `json:"academic_year_id,omitempty"`
	Teachers       []TeacherQuickPreferenceItem `json:"teachers"`
	DayOffCounts   map[int]int                  `json:"day_off_counts"`
}
