package timetablegen

import (
	"context"
	"fmt"
	"testing"
)

func FuzzTimetableGenerationHours(f *testing.F) {
	f.Add(1, 1)
	f.Add(4, 2)
	f.Add(6, 0)
	f.Add(10, 5)

	f.Fuzz(func(t *testing.T, hoursPerWeek, labHours int) {
		if hoursPerWeek < 1 || hoursPerWeek > 25 {
			return
		}
		if labHours < 0 || labHours > hoursPerWeek {
			return
		}

		gen := NewGenerator(DefaultConfig())
		ctx := context.Background()

		building := "b-1"
		assignments := []AssignmentData{
			{
				ClassID:       "c-1",
				ClassName:     "1A",
				BuildingID:    &building,
				SubjectID:     "s-1",
				SubjectName:   "Materia 1",
				TeacherID:     "t-1",
				TeacherUserID: "u-1",
				HoursPerWeek:  hoursPerWeek,
			},
		}

		rooms := []RoomData{
			{
				ID:         "r-1",
				BuildingID: &building,
				Name:       "Lab 1",
				RoomType:   "lab_type",
				Capacity:   30,
				IsActive:   true,
			},
		}

		roomReqs := map[string]SubjectRoomRequirement{
			"s-1": {
				SubjectID:        "s-1",
				RequiredRoomType: "lab_type",
				LabHours:         labHours,
				IsMandatory:      true,
			},
		}

		res, err := gen.Generate(ctx, assignments, rooms, roomReqs, nil, nil)
		if err == nil {
			if res.AssignedSlots > hoursPerWeek {
				t.Fatalf("assigned slots %d exceeded hours per week %d", res.AssignedSlots, hoursPerWeek)
			}
		}
	})
}

func FuzzDefaultConfig(f *testing.F) {
	f.Add(5, 6, 2)
	f.Add(6, 8, 4)
	f.Add(0, 0, 0)

	f.Fuzz(func(t *testing.T, days, hoursPerDay, maxConsecutive int) {
		cfg := GeneratorConfig{
			MaxDaysPerWeek:   days,
			MaxHoursPerDay:   hoursPerDay,
			MaxIterations:    maxConsecutive,
			TimeLimitSeconds: 10,
		}

		gen := NewGenerator(cfg)
		if gen == nil {
			t.Fatalf("failed to create generator")
		}
		_ = fmt.Sprintf("%v", gen)
	})
}
