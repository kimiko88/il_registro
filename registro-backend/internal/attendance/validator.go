package attendance

import (
	"errors"
	"time"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateEntry(a *Attendance) error {
	now := time.Now()
	date := a.Date

	// 1. Future Check
	if date.After(now.Add(24 * time.Hour)) {
		return errors.New("cannot mark attendance in future")
	}

	// 2. Retroactive Limit (30 days limit for regular modifications)
	limit := now.AddDate(0, 0, -30)
	if date.Before(limit) {
		return errors.New("cannot edit attendance older than 30 days")
	}

	// 3. Status logic
	if a.Status == StatusLate && a.Hour == nil {
		return errors.New("hour required for Late status")
	}
	if a.Status == StatusEarlyExit && a.Hour == nil {
		return errors.New("hour required for Early Exit")
	}

	return nil
}

func (v *Validator) ValidateJustification(j *Justification) error {
	if j.EndDate.Before(j.StartDate) {
		return errors.New("end date before start date")
	}
	return nil
}
