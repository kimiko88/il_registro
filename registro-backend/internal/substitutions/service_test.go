package substitutions

import (
	"context"
	"errors"
	"testing"
	"time"
)

// mockSubstitutionRepo implements Repository for unit tests.
type mockSubstitutionRepo struct {
	subs   []*Substitution
	getErr error
	saveErr error
}

func (m *mockSubstitutionRepo) Create(ctx context.Context, sub *Substitution) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	sub.ID = "generated-id"
	m.subs = append(m.subs, sub)
	return nil
}

func (m *mockSubstitutionRepo) GetByID(ctx context.Context, id string) (*Substitution, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	for _, s := range m.subs {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockSubstitutionRepo) ListBySchool(ctx context.Context, schoolID, date string) ([]*Substitution, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	var result []*Substitution
	for _, s := range m.subs {
		if s.SchoolID == schoolID {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *mockSubstitutionRepo) ListByTeacher(ctx context.Context, teacherID string) ([]*Substitution, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	var result []*Substitution
	for _, s := range m.subs {
		if s.AbsentTeacherID == teacherID {
			result = append(result, s)
		}
		if s.SubstituteTeacherID != nil && *s.SubstituteTeacherID == teacherID {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *mockSubstitutionRepo) AssignSubstitute(ctx context.Context, id string, substituteTeacherID string, notes string) error {
	return m.saveErr
}

func (m *mockSubstitutionRepo) ConfirmSubstitution(ctx context.Context, id string, substituteTeacherID string) error {
	return m.saveErr
}

var _ Repository = (*mockSubstitutionRepo)(nil)

// --- CreateSubstitution tests ---

func TestCreateSubstitution_ValidRequest(t *testing.T) {
	repo := &mockSubstitutionRepo{}
	svc := NewService(repo)

	sub, err := svc.CreateSubstitution(context.Background(), "admin", "school-1", CreateSubstitutionRequest{
		ClassID:         "class-1",
		AbsentTeacherID: "teacher-1",
		Date:            "2026-01-15",
		Slot:            3,
		Subject:         "Matematica",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sub.Status != StatusPending {
		t.Errorf("expected status 'pending', got '%s'", sub.Status)
	}
	if sub.Slot != 3 {
		t.Errorf("expected Slot 3, got %d", sub.Slot)
	}
}

func TestCreateSubstitution_WithSubstituteSetStatusAssigned(t *testing.T) {
	repo := &mockSubstitutionRepo{}
	svc := NewService(repo)

	subTeacher := "teacher-2"
	sub, err := svc.CreateSubstitution(context.Background(), "admin", "school-1", CreateSubstitutionRequest{
		ClassID:             "class-1",
		AbsentTeacherID:     "teacher-1",
		SubstituteTeacherID: &subTeacher,
		Date:                "2026-01-15",
		Slot:                2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sub.Status != StatusAssigned {
		t.Errorf("expected status 'assigned', got '%s'", sub.Status)
	}
}

func TestCreateSubstitution_InvalidDate(t *testing.T) {
	repo := &mockSubstitutionRepo{}
	svc := NewService(repo)

	_, err := svc.CreateSubstitution(context.Background(), "admin", "school-1", CreateSubstitutionRequest{
		ClassID:         "class-1",
		AbsentTeacherID: "teacher-1",
		Date:            "not-a-date",
	})
	if err == nil {
		t.Error("expected error for invalid date format, got nil")
	}
}

func TestCreateSubstitution_SlotFallbackToHour(t *testing.T) {
	repo := &mockSubstitutionRepo{}
	svc := NewService(repo)

	// When Slot=0 but Hour=4, Hour should be used
	sub, err := svc.CreateSubstitution(context.Background(), "admin", "school-1", CreateSubstitutionRequest{
		ClassID:         "class-1",
		AbsentTeacherID: "teacher-1",
		Date:            "2026-01-15",
		Slot:            0,
		Hour:            4,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sub.Slot != 4 {
		t.Errorf("expected Slot 4 (from Hour fallback), got %d", sub.Slot)
	}
}

func TestCreateSubstitution_PropagatesRepositoryError(t *testing.T) {
	repo := &mockSubstitutionRepo{saveErr: errors.New("db failure")}
	svc := NewService(repo)

	_, err := svc.CreateSubstitution(context.Background(), "admin", "school-1", CreateSubstitutionRequest{
		ClassID:         "class-1",
		AbsentTeacherID: "teacher-1",
		Date:            "2026-01-15",
	})
	if err == nil {
		t.Error("expected error from repository, got nil")
	}
}

// --- ListMyToday tests ---

func TestListMyToday_FiltersByDate(t *testing.T) {
	todayStr := time.Now().Format("2006-01-02")
	today, _ := time.Parse("2006-01-02", todayStr)
	yesterday := today.Add(-24 * time.Hour)

	repo := &mockSubstitutionRepo{
		subs: []*Substitution{
			{ID: "s1", AbsentTeacherID: "teacher-1", Date: today},
			{ID: "s2", AbsentTeacherID: "teacher-1", Date: yesterday},
		},
	}
	svc := NewService(repo)

	result, err := svc.ListMyToday(context.Background(), "teacher-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 substitution for today, got %d", len(result))
	}
	if result[0].ID != "s1" {
		t.Errorf("expected substitution ID 's1', got '%s'", result[0].ID)
	}
}

// --- SubstitutionStatus constants ---

func TestSubstitutionStatusConstants(t *testing.T) {
	if StatusPending != "pending" {
		t.Errorf("expected 'pending', got '%s'", StatusPending)
	}
	if StatusConfirmed != "confirmed" {
		t.Errorf("expected 'confirmed', got '%s'", StatusConfirmed)
	}
	if StatusAssigned != "assigned" {
		t.Errorf("expected 'assigned', got '%s'", StatusAssigned)
	}
	if StatusCancelled != "cancelled" {
		t.Errorf("expected 'cancelled', got '%s'", StatusCancelled)
	}
}
