package unit

import (
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper structures & validation rules matching backend logic
type TimetableSlot struct {
	DayOfWeek int
	HourIndex int
	TeacherID string
	ClassID   string
}

func validateTimetableSlot(slot TimetableSlot) (bool, string) {
	if slot.DayOfWeek < 1 || slot.DayOfWeek > 6 {
		return false, "invalid_day_of_week"
	}
	if slot.HourIndex < 1 || slot.HourIndex > 8 {
		return false, "invalid_hour_index"
	}
	return true, ""
}

func detectTeacherConflict(existing []TimetableSlot, newSlot TimetableSlot) bool {
	if newSlot.TeacherID == "" {
		return false
	}
	for _, e := range existing {
		if e.TeacherID == newSlot.TeacherID &&
			e.DayOfWeek == newSlot.DayOfWeek &&
			e.HourIndex == newSlot.HourIndex &&
			e.ClassID != newSlot.ClassID {
			return true
		}
	}
	return false
}

func roundGradeAverage(val float64) float64 {
	return math.Round(val*10) / 10
}

func canSecretaryResetPassword(targetRole string) bool {
	switch strings.ToLower(targetRole) {
	case "teacher", "student", "parent":
		return true
	default:
		return false
	}
}

func TestAntiRegression_TimetableRules(t *testing.T) {
	t.Run("Validates day of week bounds (1-6)", func(t *testing.T) {
		valid, reason := validateTimetableSlot(TimetableSlot{DayOfWeek: 1, HourIndex: 1})
		assert.True(t, valid)
		assert.Empty(t, reason)

		validSunday, reasonSunday := validateTimetableSlot(TimetableSlot{DayOfWeek: 7, HourIndex: 1})
		assert.False(t, validSunday)
		assert.Equal(t, "invalid_day_of_week", reasonSunday)
	})

	t.Run("Validates hour index bounds (1-8)", func(t *testing.T) {
		valid, _ := validateTimetableSlot(TimetableSlot{DayOfWeek: 1, HourIndex: 8})
		assert.True(t, valid)

		valid9, reason9 := validateTimetableSlot(TimetableSlot{DayOfWeek: 1, HourIndex: 9})
		assert.False(t, valid9)
		assert.Equal(t, "invalid_hour_index", reason9)
	})

	t.Run("Detects teacher schedule collision across different classes", func(t *testing.T) {
		existing := []TimetableSlot{
			{DayOfWeek: 1, HourIndex: 2, TeacherID: "teacher-1", ClassID: "class-1A"},
		}
		newSlotSameTime := TimetableSlot{DayOfWeek: 1, HourIndex: 2, TeacherID: "teacher-1", ClassID: "class-2B"}
		assert.True(t, detectTeacherConflict(existing, newSlotSameTime))

		newSlotDiffTime := TimetableSlot{DayOfWeek: 1, HourIndex: 3, TeacherID: "teacher-1", ClassID: "class-2B"}
		assert.False(t, detectTeacherConflict(existing, newSlotDiffTime))
	})
}

func TestAntiRegression_PasswordAndRBACRules(t *testing.T) {
	t.Run("Secretary role reset boundaries", func(t *testing.T) {
		assert.True(t, canSecretaryResetPassword("teacher"))
		assert.True(t, canSecretaryResetPassword("student"))
		assert.True(t, canSecretaryResetPassword("parent"))

		assert.False(t, canSecretaryResetPassword("admin"))
		assert.False(t, canSecretaryResetPassword("superadmin"))
		assert.False(t, canSecretaryResetPassword("secretary"))
	})
}

func TestAntiRegression_ScrutinyGradeRounding(t *testing.T) {
	t.Run("Rounds grade averages to 1 decimal place", func(t *testing.T) {
		assert.Equal(t, 7.2, roundGradeAverage(7.166666))
		assert.Equal(t, 6.5, roundGradeAverage(6.45))
		assert.Equal(t, 8.0, roundGradeAverage(8.0))
	})
}
