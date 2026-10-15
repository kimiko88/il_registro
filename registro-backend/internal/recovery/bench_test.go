package recovery

import (
	"context"
	"testing"
)

func BenchmarkRecordTestOutcome(b *testing.B) {
	repo := &mockRecoveryRepo{}
	svc := NewService(repo)
	ctx := context.Background()

	req := RecordTestOutcomeRequest{
		StudentID: "stu-1",
		SubjectID: "sub-1",
		ClassID:   "cls-1",
		TestDate:  "2026-09-02",
		TestType:  "written",
		Grade:     7.5,
	}

	for b.Loop() {
		_, _ = svc.RecordTestOutcome(ctx, "school-1", "teacher-1", req)
	}
}

func BenchmarkGetCourse(b *testing.B) {
	repo := &mockRecoveryRepo{
		course: &RecoveryCourse{
			ID:         "c-bench",
			Title:      "Recupero Estivo Matematica",
			TotalHours: 15,
		},
	}
	svc := NewService(repo)
	ctx := context.Background()

	for b.Loop() {
		_, _ = svc.GetCourse(ctx, "c-bench")
	}
}
