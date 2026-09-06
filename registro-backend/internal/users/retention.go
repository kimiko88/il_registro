package users

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidRetentionYears = errors.New("retention years must be at least 1")
)

// RetentionRequest represents parameters for data retention policy execution
type RetentionRequest struct {
	SchoolID       *string `json:"school_id,omitempty"`
	RetentionYears int     `json:"retention_years"`
}

// RetentionResult summarizes the results of the retention policy execution
type RetentionResult struct {
	ProcessedCount int       `json:"processed_count"`
	RetentionYears int       `json:"retention_years"`
	CutoffDate     time.Time `json:"cutoff_date"`
	Message        string    `json:"message"`
}

// ApplyDataRetention executes the GDPR data retention policy by pseudonymizing
// inactive or deleted student records older than the retention threshold.
// Accessible only to school admin or superadmin.
func (s *Service) ApplyDataRetention(ctx context.Context, actorRole string, schoolID *string, retentionYears int) (*RetentionResult, error) {
	if actorRole != "admin" && actorRole != "superadmin" {
		return nil, ErrUnauthorized
	}

	if retentionYears <= 0 {
		retentionYears = 5
	}
	if retentionYears < 1 {
		return nil, ErrInvalidRetentionYears
	}

	cutoffDate := time.Now().AddDate(-retentionYears, 0, 0)

	count, err := s.repo.ApplyDataRetention(ctx, schoolID, cutoffDate)
	if err != nil {
		return nil, fmt.Errorf("failed to apply data retention: %w", err)
	}

	return &RetentionResult{
		ProcessedCount: count,
		RetentionYears: retentionYears,
		CutoffDate:     cutoffDate,
		Message:        fmt.Sprintf("Retention policy applied successfully: %d student records pseudonymized", count),
	}, nil
}
