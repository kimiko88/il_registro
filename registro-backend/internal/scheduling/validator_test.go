package scheduling

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidator_ValidateSchedule(t *testing.T) {
	v := NewValidator()

	t.Run("No Conflicts", func(t *testing.T) {
		slots := []Slot{
			{ID: "1", TeacherID: "T1", ClassID: "C1", DayOfWeek: time.Monday, Hour: 1},
			{ID: "2", TeacherID: "T2", ClassID: "C2", DayOfWeek: time.Monday, Hour: 1}, // Different everything
			{ID: "3", TeacherID: "T1", ClassID: "C1", DayOfWeek: time.Monday, Hour: 2}, // Different hour
		}
		conflicts := v.ValidateSchedule(slots)
		assert.Empty(t, conflicts)
	})

	t.Run("Teacher Overlap", func(t *testing.T) {
		slots := []Slot{
			{ID: "1", TeacherID: "T1", ClassID: "C1", DayOfWeek: time.Monday, Hour: 1},
			{ID: "2", TeacherID: "T1", ClassID: "C2", DayOfWeek: time.Monday, Hour: 1}, // Same Teacher, Same Time
		}
		conflicts := v.ValidateSchedule(slots)
		assert.Len(t, conflicts, 1)
		assert.Equal(t, ConflictTeacherOverlap, conflicts[0].Type)
	})

	t.Run("Class Overlap", func(t *testing.T) {
		slots := []Slot{
			{ID: "1", TeacherID: "T1", ClassID: "C1", DayOfWeek: time.Monday, Hour: 1},
			{ID: "2", TeacherID: "T2", ClassID: "C1", DayOfWeek: time.Monday, Hour: 1}, // Same Class, Same Time
		}
		conflicts := v.ValidateSchedule(slots)
		assert.Len(t, conflicts, 1)
		assert.Equal(t, ConflictClassOverlap, conflicts[0].Type)
	})

	t.Run("Room Overlap", func(t *testing.T) {
		slots := []Slot{
			{ID: "1", TeacherID: "T1", ClassID: "C1", RoomID: "R101", DayOfWeek: time.Monday, Hour: 1},
			{ID: "2", TeacherID: "T2", ClassID: "C2", RoomID: "R101", DayOfWeek: time.Monday, Hour: 1}, // Same Room
		}
		conflicts := v.ValidateSchedule(slots)
		assert.Len(t, conflicts, 1)
		assert.Equal(t, ConflictRoomOverlap, conflicts[0].Type)
	})
}
