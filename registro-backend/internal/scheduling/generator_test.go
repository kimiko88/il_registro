package scheduling

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerator_Generate_Simple(t *testing.T) {
	gen := NewGenerator()

	req := GenerationRequest{
		SchoolID:  "school-1",
		StartDate: time.Now(),
		Days:      []time.Weekday{time.Monday, time.Tuesday},
		StartHour: 8,
		EndHour:   13,
		Classes:   []string{"1A", "1B"},
		Teachers:  []string{"T1", "T2"},
		Assignments: []Assignment{
			{TeacherID: "T1", ClassID: "1A", SubjectID: "Math", Hours: 2},
			{TeacherID: "T2", ClassID: "1B", SubjectID: "Hist", Hours: 2},
		},
	}

	slots, err := gen.Generate(req)
	assert.NoError(t, err)
	assert.Len(t, slots, 4) // 2 for T1, 2 for T2

	// Check conflicts
	validator := NewValidator()
	conflicts := validator.ValidateSchedule(slots)
	assert.Empty(t, conflicts)
}

func TestGenerator_Generate_ConflictResolution(t *testing.T) {
	gen := NewGenerator()

	// Heavily constrained: T1 needs to teach 1A and 1B, but only 1 hour available
	req := GenerationRequest{
		Days:      []time.Weekday{time.Monday},
		StartHour: 8,
		EndHour:   9, // Only 1 slot available: Mon 8-9
		Assignments: []Assignment{
			{TeacherID: "T1", ClassID: "1A", Hours: 1},
			{TeacherID: "T1", ClassID: "1B", Hours: 1}, // Conflict! Same teacher, same time slot needed
		},
	}

	_, err := gen.Generate(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not find valid slot")
}
