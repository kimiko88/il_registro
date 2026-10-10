package timetablegen

import (
	"context"
	"fmt"
	"testing"
)

func BenchmarkTimetableGenerationSmall(b *testing.B) {
	ctx := context.Background()
	building := "b-central"

	assignments := []AssignmentData{
		{
			ClassID:       "c-1",
			ClassName:     "1A",
			BuildingID:    &building,
			SubjectID:     "s-ita",
			SubjectName:   "Italiano",
			TeacherID:     "t-1",
			TeacherUserID: "u-1",
			HoursPerWeek:  4,
		},
		{
			ClassID:       "c-1",
			ClassName:     "1A",
			BuildingID:    &building,
			SubjectID:     "s-mat",
			SubjectName:   "Matematica",
			TeacherID:     "t-2",
			TeacherUserID: "u-2",
			HoursPerWeek:  4,
		},
		{
			ClassID:       "c-1",
			ClassName:     "1A",
			BuildingID:    &building,
			SubjectID:     "s-ing",
			SubjectName:   "Inglese",
			TeacherID:     "t-3",
			TeacherUserID: "u-3",
			HoursPerWeek:  3,
		},
	}

	rooms := []RoomData{
		{
			ID:         "r-aula1",
			BuildingID: &building,
			Name:       "Aula 1",
			RoomType:   "aula",
			Capacity:   25,
			IsActive:   true,
		},
	}

	cfg := DefaultConfig()
	cfg.MaxIterations = 50 // Keep bounded for fast benchmark
	gen := NewGenerator(cfg)

	for b.Loop() {
		_, err := gen.Generate(ctx, assignments, rooms, nil, nil, nil)
		if err != nil {
			b.Fatalf("generation failed: %v", err)
		}
	}
}

func BenchmarkTimetableConstraintEvaluation(b *testing.B) {
	ctx := context.Background()
	building := "b-central"

	// 5 teachers across 2 classes
	assignments := make([]AssignmentData, 10)
	for i := 0; i < 10; i++ {
		classID := "c-1"
		if i >= 5 {
			classID = "c-2"
		}
		assignments[i] = AssignmentData{
			ClassID:       classID,
			ClassName:     "Class",
			BuildingID:    &building,
			SubjectID:     fmt.Sprintf("sub-%d", i),
			SubjectName:   fmt.Sprintf("Subject %d", i),
			TeacherID:     fmt.Sprintf("t-%d", i%5),
			TeacherUserID: fmt.Sprintf("u-%d", i%5),
			HoursPerWeek:  2,
		}
	}

	cfg := DefaultConfig()
	cfg.MaxIterations = 20
	gen := NewGenerator(cfg)

	for b.Loop() {
		_, _ = gen.Generate(ctx, assignments, nil, nil, nil, nil)
	}
}
