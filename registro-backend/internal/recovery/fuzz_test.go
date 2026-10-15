package recovery

import (
	"context"
	"math"
	"testing"
)

type mockRecoveryRepo struct {
	course *RecoveryCourse
	tests  []RecoveryTest
}

func (m *mockRecoveryRepo) CreateCourse(ctx context.Context, c *RecoveryCourse, s []RecoveryCourseSession, st []string) error {
	c.ID = "c-123"
	m.course = c
	return nil
}
func (m *mockRecoveryRepo) GetCourseByID(ctx context.Context, id string) (*RecoveryCourse, error) {
	if m.course != nil {
		return m.course, nil
	}
	return &RecoveryCourse{ID: id, Title: "Corso Matematica"}, nil
}
func (m *mockRecoveryRepo) ListCourses(ctx context.Context, schoolID, academicYear, teacherID string) ([]RecoveryCourse, error) {
	return []RecoveryCourse{}, nil
}
func (m *mockRecoveryRepo) UpdateCourseStatus(ctx context.Context, id, status string) error {
	return nil
}
func (m *mockRecoveryRepo) UpdateStudentAttendance(ctx context.Context, courseID, studentID string, hours float64, notes string) error {
	return nil
}
func (m *mockRecoveryRepo) RecordRecoveryTest(ctx context.Context, t *RecoveryTest) error {
	m.tests = append(m.tests, *t)
	return nil
}
func (m *mockRecoveryRepo) ListRecoveryTests(ctx context.Context, schoolID, classID, studentID string) ([]RecoveryTest, error) {
	return m.tests, nil
}

func FuzzRecordTestOutcomeGradeBounds(f *testing.F) {
	seeds := []float64{
		6.0, 5.9, 1.0, 10.0, 0.9, 10.1, -5.0, 15.0, 7.5, 4.0,
	}
	for _, s := range seeds {
		f.Add(s)
	}

	repo := &mockRecoveryRepo{}
	svc := NewService(repo)
	ctx := context.Background()

	f.Fuzz(func(t *testing.T, grade float64) {
		if math.IsNaN(grade) || math.IsInf(grade, 0) {
			return
		}

		req := RecordTestOutcomeRequest{
			StudentID: "stu-1",
			SubjectID: "sub-1",
			ClassID:   "cls-1",
			TestDate:  "2026-09-02",
			TestType:  "written",
			Grade:     grade,
		}

		res, err := svc.RecordTestOutcome(ctx, "school-1", "teacher-1", req)

		if grade < 1.0 || grade > 10.0 {
			if err != ErrInvalidGrade {
				t.Fatalf("expected ErrInvalidGrade for grade %.2f, got: %v", grade, err)
			}
			if res != nil {
				t.Fatalf("expected nil result on invalid grade")
			}
			return
		}

		if err != nil {
			t.Fatalf("unexpected error for valid grade %.2f: %v", grade, err)
		}
		if res == nil {
			t.Fatalf("expected non-nil result for valid grade")
		}

		if grade >= 6.0 {
			if res.Outcome != "recuperato" || res.FinalDeliberation != "Ammesso" {
				t.Fatalf("expected recuperato / Ammesso for grade %.2f, got %s / %s", grade, res.Outcome, res.FinalDeliberation)
			}
		} else {
			if res.Outcome != "non_recuperato" || res.FinalDeliberation != "Non Ammesso" {
				t.Fatalf("expected non_recuperato / Non Ammesso for grade %.2f, got %s / %s", grade, res.Outcome, res.FinalDeliberation)
			}
		}
	})
}
