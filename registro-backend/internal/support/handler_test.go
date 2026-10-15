package support

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"registro-backend/internal/users"
)

type mockSupportService struct {
	diaries []SupportDiaryEntry
	goals   []SupportPeiGoal
}

func (m *mockSupportService) CreateDiaryEntry(ctx context.Context, schoolID, teacherID string, req CreateDiaryEntryRequest) (*SupportDiaryEntry, error) {
	entry := &SupportDiaryEntry{
		ID:                 "entry-1",
		SchoolID:           schoolID,
		TeacherID:          teacherID,
		StudentID:          req.StudentID,
		ClassID:            req.ClassID,
		EntryDate:          req.EntryDate,
		TimeSlot:           req.TimeSlot,
		ActivityType:       req.ActivityType,
		TopicAndActivities: req.TopicAndActivities,
		IsSharedWithFamily: req.IsSharedWithFamily,
		CreatedAt:          time.Now(),
	}
	m.diaries = append(m.diaries, *entry)
	return entry, nil
}

func (m *mockSupportService) ListDiaryEntries(ctx context.Context, schoolID, studentID, teacherID, classID string, isFamily bool) ([]SupportDiaryEntry, error) {
	return m.diaries, nil
}

func (m *mockSupportService) DeleteDiaryEntry(ctx context.Context, id, teacherID string) error {
	return nil
}

func (m *mockSupportService) CreatePeiGoal(ctx context.Context, schoolID string, req CreatePeiGoalRequest) (*SupportPeiGoal, error) {
	goal := &SupportPeiGoal{
		ID:             "goal-1",
		SchoolID:       schoolID,
		StudentID:      req.StudentID,
		PeiType:        req.PeiType,
		Axis:           req.Axis,
		Title:          req.Title,
		Description:    req.Description,
		ProgressStatus: "in_corso",
	}
	m.goals = append(m.goals, *goal)
	return goal, nil
}

func (m *mockSupportService) UpdateGoalProgress(ctx context.Context, id, progressStatus string) error {
	return nil
}

func (m *mockSupportService) ListPeiGoals(ctx context.Context, schoolID, studentID string) ([]SupportPeiGoal, error) {
	return m.goals, nil
}

func (m *mockSupportService) DeletePeiGoal(ctx context.Context, id string) error {
	return nil
}

func (m *mockSupportService) GetUserRepo() users.Repository {
	return nil
}

func setupSupportRouter(svc Service, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("school_id", "school-1")
		c.Set("user_id", "teacher-1")
		c.Set("role", role)
		c.Next()
	})
	h := NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestSupportHandler_RoleAccess(t *testing.T) {
	svc := &mockSupportService{}
	rStudent := setupSupportRouter(svc, "student")

	// Student cannot create diary entry
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/support/diaries", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	rStudent.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden status 403, got %d", w.Code)
	}

	// Teacher can create diary entry
	rTeacher := setupSupportRouter(svc, "teacher")
	body, _ := json.Marshal(CreateDiaryEntryRequest{
		StudentID:          "stu-1",
		ClassID:            "cls-1",
		EntryDate:          "2026-03-10",
		TimeSlot:           "2ª Ora (09:00-10:00)",
		ActivityType:       "individuale",
		TopicAndActivities: "Supporto logico-matematico",
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/support/diaries", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rTeacher.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSupportHandler_PeiGoals(t *testing.T) {
	svc := &mockSupportService{}
	r := setupSupportRouter(svc, "teacher")

	// Create PEI goal
	body, _ := json.Marshal(CreatePeiGoalRequest{
		StudentID:   "stu-1",
		PeiType:     "equipollente",
		Axis:        "autonomia",
		Title:       "Gestione materiale",
		Description: "Gestione autonoma del materiale scolastico",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/support/pei-goals", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 created, got %d: %s", w.Code, w.Body.String())
	}

	// List PEI goals
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/support/pei-goals?student_id=stu-1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 ok, got %d", w.Code)
	}
}
