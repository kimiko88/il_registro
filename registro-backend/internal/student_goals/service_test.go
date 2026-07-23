package student_goals

import (
	"context"
	"errors"
	"testing"
	"time"
)

// mockGoalRepo implements Repository for unit tests.
type mockGoalRepo struct {
	goals  []*StudentGoal
	getErr error
	saveErr error
}

func (m *mockGoalRepo) Create(ctx context.Context, goal *StudentGoal) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	goal.ID = "generated-id"
	m.goals = append(m.goals, goal)
	return nil
}

func (m *mockGoalRepo) GetByID(ctx context.Context, id string) (*StudentGoal, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	for _, g := range m.goals {
		if g.ID == id {
			return g, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockGoalRepo) ListByStudent(ctx context.Context, studentID string) ([]*StudentGoal, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	var result []*StudentGoal
	for _, g := range m.goals {
		if g.StudentID == studentID {
			result = append(result, g)
		}
	}
	return result, nil
}

func (m *mockGoalRepo) UpdateStatus(ctx context.Context, id string, status GoalStatus) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	for _, g := range m.goals {
		if g.ID == id {
			g.Status = status
			return nil
		}
	}
	return errors.New("not found")
}

var _ Repository = (*mockGoalRepo)(nil)

// --- CreateGoal tests ---

func TestCreateGoal_DefaultsToStatusPending(t *testing.T) {
	repo := &mockGoalRepo{}
	svc := NewService(repo)

	goal, err := svc.CreateGoal(context.Background(), "teacher-1", CreateGoalRequest{
		StudentID: "stu-1",
		Title:     "Migliorare Matematica",
		Category:  "academic",
		Points:    100,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if goal.Status != StatusPending {
		t.Errorf("expected status 'pending', got '%s'", goal.Status)
	}
	if goal.TeacherID != "teacher-1" {
		t.Errorf("expected TeacherID 'teacher-1', got '%s'", goal.TeacherID)
	}
}

func TestCreateGoal_WithValidDueDate(t *testing.T) {
	repo := &mockGoalRepo{}
	svc := NewService(repo)

	goal, err := svc.CreateGoal(context.Background(), "teacher-1", CreateGoalRequest{
		StudentID: "stu-1",
		Title:     "Completare progetto",
		DueDate:   "2026-06-30",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if goal.DueDate == nil {
		t.Error("expected non-nil DueDate")
	}
	expected := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)
	if !goal.DueDate.Equal(expected) {
		t.Errorf("expected DueDate %v, got %v", expected, *goal.DueDate)
	}
}

func TestCreateGoal_InvalidDueDateIsIgnored(t *testing.T) {
	repo := &mockGoalRepo{}
	svc := NewService(repo)

	// Invalid date format — should not fail, DueDate should be nil
	goal, err := svc.CreateGoal(context.Background(), "teacher-1", CreateGoalRequest{
		StudentID: "stu-1",
		Title:     "Meta",
		DueDate:   "invalid-date",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if goal.DueDate != nil {
		t.Errorf("expected nil DueDate for invalid date, got %v", goal.DueDate)
	}
}

func TestCreateGoal_PropagatesRepositoryError(t *testing.T) {
	repo := &mockGoalRepo{saveErr: errors.New("db failure")}
	svc := NewService(repo)

	_, err := svc.CreateGoal(context.Background(), "teacher-1", CreateGoalRequest{
		StudentID: "stu-1",
		Title:     "Meta",
	})
	if err == nil {
		t.Error("expected error from repository, got nil")
	}
}

// --- ListByStudent tests ---

func TestListByStudent_ReturnsMatchingGoals(t *testing.T) {
	repo := &mockGoalRepo{
		goals: []*StudentGoal{
			{ID: "g1", StudentID: "stu-1", Title: "Meta 1", Status: StatusPending},
			{ID: "g2", StudentID: "stu-2", Title: "Meta 2", Status: StatusPending},
			{ID: "g3", StudentID: "stu-1", Title: "Meta 3", Status: StatusCompleted},
		},
	}
	svc := NewService(repo)

	goals, err := svc.ListByStudent(context.Background(), "stu-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(goals) != 2 {
		t.Errorf("expected 2 goals for stu-1, got %d", len(goals))
	}
}

// --- UpdateGoalStatus tests ---

func TestUpdateGoalStatus_OwnerCanUpdate(t *testing.T) {
	repo := &mockGoalRepo{
		goals: []*StudentGoal{
			{ID: "g1", TeacherID: "teacher-1", StudentID: "stu-1", Status: StatusPending},
		},
	}
	svc := NewService(repo)

	err := svc.UpdateGoalStatus(context.Background(), "teacher-1", "teacher", "g1", StatusCompleted)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.goals[0].Status != StatusCompleted {
		t.Errorf("expected status 'completed', got '%s'", repo.goals[0].Status)
	}
}

func TestUpdateGoalStatus_StudentCanUpdate(t *testing.T) {
	repo := &mockGoalRepo{
		goals: []*StudentGoal{
			{ID: "g1", TeacherID: "teacher-1", StudentID: "stu-1", Status: StatusPending},
		},
	}
	svc := NewService(repo)

	err := svc.UpdateGoalStatus(context.Background(), "stu-1", "student", "g1", StatusInProgress)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateGoalStatus_UnauthorizedActor(t *testing.T) {
	repo := &mockGoalRepo{
		goals: []*StudentGoal{
			{ID: "g1", TeacherID: "teacher-1", StudentID: "stu-1", Status: StatusPending},
		},
	}
	svc := NewService(repo)

	err := svc.UpdateGoalStatus(context.Background(), "other-actor", "parent", "g1", StatusCompleted)
	if err == nil {
		t.Error("expected unauthorized error, got nil")
	}
}

func TestUpdateGoalStatus_AdminBypassesOwnerCheck(t *testing.T) {
	repo := &mockGoalRepo{
		goals: []*StudentGoal{
			{ID: "g1", TeacherID: "teacher-1", StudentID: "stu-1", Status: StatusPending},
		},
	}
	svc := NewService(repo)

	err := svc.UpdateGoalStatus(context.Background(), "admin-user", "admin", "g1", StatusCompleted)
	if err != nil {
		t.Fatalf("unexpected error for admin: %v", err)
	}
}

// --- GoalStatus constants ---

func TestGoalStatusConstants(t *testing.T) {
	if StatusPending != "pending" {
		t.Errorf("expected 'pending', got '%s'", StatusPending)
	}
	if StatusInProgress != "in_progress" {
		t.Errorf("expected 'in_progress', got '%s'", StatusInProgress)
	}
	if StatusCompleted != "completed" {
		t.Errorf("expected 'completed', got '%s'", StatusCompleted)
	}
}
