package privacy

import (
	"time"
)

type TrafficLightBadge string

const (
	TrafficLightGreen  TrafficLightBadge = "VERDE"  // Fully consented (photo/video ok)
	TrafficLightYellow TrafficLightBadge = "GIALLO" // Educational ok, but strictly NO photos on social
	TrafficLightRed    TrafficLightBadge = "ROSSO"  // Absolute ban on photos/recordings
)

type PrivacyTreatment struct {
	ID                   string    `json:"id" db:"id"`
	SchoolID             string    `json:"school_id" db:"school_id"`
	ActivityName         string    `json:"activity_name" db:"activity_name"`
	LegalBasis           string    `json:"legal_basis" db:"legal_basis"` // "obbligo_legale", "interesse_pubblico", "consenso"
	InterestedCategories string    `json:"interested_categories" db:"interested_categories"`
	RetentionPeriod      string    `json:"retention_period" db:"retention_period"`
	DPOContact           string    `json:"dpo_contact" db:"dpo_contact"`
	SecurityMeasures     string    `json:"security_measures" db:"security_measures"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
}

type StudentConsent struct {
	ID                      string            `json:"id" db:"id"`
	SchoolID                string            `json:"school_id" db:"school_id"`
	StudentID               string            `json:"student_id" db:"student_id"`
	SchoolYear              string            `json:"school_year" db:"school_year"`
	PhotoVideoSocialConsent bool              `json:"photo_video_social_consent" db:"photo_video_social_consent"`
	CloudWorkspaceConsent   bool              `json:"cloud_workspace_consent" db:"cloud_workspace_consent"`
	WalkingTripsConsent     bool              `json:"walking_trips_consent" db:"walking_trips_consent"`
	TrafficLightBadge       TrafficLightBadge `json:"traffic_light_badge" db:"traffic_light_badge"`
	SignedBy                string            `json:"signed_by,omitempty" db:"signed_by"`
	SignedAt                *time.Time        `json:"signed_at,omitempty" db:"signed_at"`
	Notes                   string            `json:"notes,omitempty" db:"notes"`
	UpdatedAt               time.Time         `json:"updated_at" db:"updated_at"`
}

type SaveConsentRequest struct {
	StudentID               string `json:"student_id" binding:"required"`
	SchoolYear              string `json:"school_year" binding:"required"`
	PhotoVideoSocialConsent bool   `json:"photo_video_social_consent"`
	CloudWorkspaceConsent   bool   `json:"cloud_workspace_consent"`
	WalkingTripsConsent     bool   `json:"walking_trips_consent"`
	Notes                   string `json:"notes"`
}

type TreatmentRequest struct {
	ActivityName         string `json:"activity_name" binding:"required"`
	LegalBasis           string `json:"legal_basis" binding:"required"`
	InterestedCategories string `json:"interested_categories" binding:"required"`
	RetentionPeriod      string `json:"retention_period" binding:"required"`
	DPOContact           string `json:"dpo_contact"`
	SecurityMeasures     string `json:"security_measures"`
}
