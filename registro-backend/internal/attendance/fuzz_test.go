package attendance

import (
	"testing"
	"time"
)

func FuzzAttendanceStatusValidation(f *testing.F) {
	seeds := []string{
		string(StatusPresent),
		string(StatusAbsent),
		string(StatusLate),
		string(StatusEarlyExit),
		string(StatusOutOfClass),
		string(StatusExempt),
		"presente",
		"assente",
		"",
		"STATUS_UNKNOWN",
		"\x00\xff",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	validator := NewValidator()

	f.Fuzz(func(t *testing.T, statusStr string) {
		status := AttendanceStatus(statusStr)
		isValid := validator.IsValidStatus(status)

		switch status {
		case StatusPresent, StatusAbsent, StatusLate, StatusEarlyExit, StatusOutOfClass, StatusExempt:
			if !isValid {
				t.Fatalf("expected true for valid status %s", status)
			}
		default:
			if isValid {
				t.Fatalf("expected false for invalid status %s", status)
			}
		}
	})
}

func FuzzValidateEntry(f *testing.F) {
	f.Add(string(StatusPresent), 0, false)
	f.Add(string(StatusLate), 2, true)
	f.Add(string(StatusLate), 0, false)
	f.Add(string(StatusEarlyExit), 4, true)
	f.Add(string(StatusEarlyExit), 0, false)
	f.Add("invalid", 1, true)

	validator := NewValidator()

	f.Fuzz(func(t *testing.T, statusStr string, hourVal int, hasHour bool) {
		status := AttendanceStatus(statusStr)
		var hour *int
		if hasHour {
			h := hourVal
			hour = &h
		}

		entry := &Attendance{
			Status: status,
			Date:   time.Now(),
			Hour:   hour,
		}

		err := validator.ValidateEntry(entry)
		if !validator.IsValidStatus(status) && err == nil {
			t.Fatalf("expected error for invalid status: %s", statusStr)
		}
		if (status == StatusLate || status == StatusEarlyExit) && hour == nil && err == nil {
			t.Fatalf("expected error for late/early exit without hour")
		}
	})
}
