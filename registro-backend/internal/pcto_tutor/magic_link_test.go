package pcto_tutor

import (
	"testing"
	"time"
)

func TestGenerateMagicLinkToken(t *testing.T) {
	token, expiresAt, err := GenerateMagicLinkToken(72 * time.Hour)
	if err != nil {
		t.Fatalf("failed to generate magic link token: %v", err)
	}

	if len(token) != 64 {
		t.Errorf("expected 64 char hex token, got length %d", len(token))
	}

	if expiresAt.Before(time.Now().Add(71 * time.Hour)) {
		t.Errorf("expected expiration > 71h from now, got %v", expiresAt)
	}
}

func TestValidateEvaluationRatings(t *testing.T) {
	eval := CompanyEvaluation{
		ReliabilityLevel: 5,
		TechnicalSkills:  4,
		TeamworkSkills:   5,
		FinalFeedback:    "Ottima attitudine al lavoro di squadra e puntualità",
	}

	if err := ValidateEvaluation(&eval); err != nil {
		t.Errorf("valid evaluation failed validation: %v", err)
	}

	invalidEval := CompanyEvaluation{
		ReliabilityLevel: 6, // Exceeds max 5
		TechnicalSkills:  0, // Below min 1
		TeamworkSkills:   3,
	}
	if err := ValidateEvaluation(&invalidEval); err == nil {
		t.Errorf("invalid evaluation should fail rating range check")
	}
}

func TestVerifyTimesheetHours(t *testing.T) {
	entry := TimesheetVerification{
		HoursDeclared: 8.0,
		HoursApproved: 8.0,
	}
	if err := ValidateTimesheetEntry(&entry); err != nil {
		t.Errorf("valid timesheet entry failed: %v", err)
	}

	invalidEntry := TimesheetVerification{
		HoursDeclared: 8.0,
		HoursApproved: 10.0, // Approved cannot exceed declared
	}
	if err := ValidateTimesheetEntry(&invalidEntry); err == nil {
		t.Errorf("approved hours exceeding declared hours should fail")
	}

	negativeEntry := TimesheetVerification{
		HoursDeclared: -2.0,
		HoursApproved: 0,
	}
	if err := ValidateTimesheetEntry(&negativeEntry); err == nil {
		t.Errorf("negative hours should fail")
	}
}
