package attendance

import (
	"errors"
	"time"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) IsValidStatus(status AttendanceStatus) bool {
	switch status {
	case StatusPresent, StatusAbsent, StatusLate, StatusEarlyExit, StatusOutOfClass, StatusExempt:
		return true
	default:
		return false
	}
}

func (v *Validator) ValidateEntry(a *Attendance) error {
	if !v.IsValidStatus(a.Status) {
		return errors.New("invalid attendance status")
	}

	loc, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		loc = time.Local
	}
	nowInLoc := time.Now().In(loc)
	date := a.Date

	// 1. Future Check: cannot mark attendance after today's end
	todayEnd := time.Date(nowInLoc.Year(), nowInLoc.Month(), nowInLoc.Day(), 23, 59, 59, 999999999, loc)
	if date.After(todayEnd) {
		return errors.New("cannot mark attendance in future")
	}

	// 2. Historical Limit (reject records older than 2 years / obsolete school years)
	minDate := nowInLoc.AddDate(-2, 0, 0)
	if date.Before(minDate) {
		return errors.New("cannot mark attendance for dates older than 2 years")
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
