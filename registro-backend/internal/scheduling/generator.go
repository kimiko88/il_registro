package scheduling

import (
	"errors"
	"time"
)

// Generator handles the automated creation of schedules
type Generator struct {
	validator *Validator
}

func NewGenerator() *Generator {
	return &Generator{
		validator: NewValidator(),
	}
}

// GenerationRequest defines constraints for schedule generation
type GenerationRequest struct {
	SchoolID    string         `json:"school_id"`
	StartDate   time.Time      `json:"start_date"`
	Days        []time.Weekday `json:"days"`
	StartHour   int            `json:"start_hour"` // e.g. 8 (8:00)
	EndHour     int            `json:"end_hour"`   // e.g. 14 (14:00)
	Classes     []string       `json:"classes"`
	Teachers    []string       `json:"teachers"`
	Assignments []Assignment   `json:"assignments"` // Teacher X teaches Subject Y in Class Z (hours per week)
}

type Assignment struct {
	TeacherID string `json:"teacher_id"`
	ClassID   string `json:"class_id"`
	SubjectID string `json:"subject_id"`
	Hours     int    `json:"hours"`
}

// Generate uses a simple backtracking algorithm to assign slots
func (g *Generator) Generate(req GenerationRequest) ([]Slot, error) {
	if len(req.Assignments) == 0 {
		return nil, errors.New("no assignments provided")
	}

	// Flatten needed slots: Assignment(T1, C1, 2h) -> [Slot(T1,C1), Slot(T1,C1)]
	var pendingSlots []Slot
	for _, assign := range req.Assignments {
		for i := 0; i < assign.Hours; i++ {
			pendingSlots = append(pendingSlots, Slot{
				TeacherID: assign.TeacherID,
				ClassID:   assign.ClassID,
				SubjectID: assign.SubjectID,
			})
		}
	}

	// Available Time Slots: Day/Hour combinations
	type TimeSlot struct {
		Day  time.Weekday
		Hour int
	}
	var timeSlots []TimeSlot
	for _, day := range req.Days {
		for h := req.StartHour; h < req.EndHour; h++ {
			timeSlots = append(timeSlots, TimeSlot{Day: day, Hour: h})
		}
	}

	// Backtracking
	// result, ok := g.solve(pendingSlots, timeSlots, []Slot{})
	// Simplification: Greedy approach for MVP
	// For each pending slot, try to find the FIRST available time slot that doesn't conflict

	generatedSchedule := []Slot{}

	for _, p := range pendingSlots {
		assigned := false
		for _, ts := range timeSlots {
			// Try to place p at ts
			candidate := p // copy
			candidate.DayOfWeek = ts.Day
			candidate.Hour = ts.Hour

			// Check against already generated schedule for this specific candidate
			// We only check against the partial schedule we built so far
			// This is effectively checking constraints incrementally

			// optimization: check specific conflicts instead of full array scan in validator
			// but utilizing ValidateSchedule for correctness on the whole set

			tempSchedule := append(generatedSchedule, candidate)
			conflicts := g.validator.ValidateSchedule(tempSchedule)

			if len(conflicts) == 0 {
				generatedSchedule = tempSchedule
				assigned = true
				break
			}
		}
		if !assigned {
			// In a real generator, we would backtrack here or return "Partial Success / Error"
			// For MVP, we'll continue and just note it wasn't assigned (or return error)
			return generatedSchedule, errors.New("could not find valid slot for assignment (greedy failure)")
		}
	}

	return generatedSchedule, nil
}
