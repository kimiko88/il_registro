package attendance

import (
	"testing"
	"time"
)

func BenchmarkIsValidStatus(b *testing.B) {
	validator := NewValidator()
	statuses := []AttendanceStatus{
		StatusPresent,
		StatusAbsent,
		StatusLate,
		StatusEarlyExit,
		StatusOutOfClass,
		StatusExempt,
		"invalid",
	}

	idx := 0
	for b.Loop() {
		_ = validator.IsValidStatus(statuses[idx%len(statuses)])
		idx++
	}
}

func BenchmarkValidateEntry(b *testing.B) {
	validator := NewValidator()
	hour := 2
	entry := &Attendance{
		Status: StatusLate,
		Date:   time.Now(),
		Hour:   &hour,
	}

	for b.Loop() {
		_ = validator.ValidateEntry(entry)
	}
}

func BenchmarkValidateJustification(b *testing.B) {
	validator := NewValidator()
	start := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC)
	justification := &Justification{
		StartDate: start,
		EndDate:   end,
	}

	for b.Loop() {
		_ = validator.ValidateJustification(justification)
	}
}
