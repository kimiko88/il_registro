package pcto_tutor

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type mockTutorRepo struct {
	tutors       map[string]*CompanyTutor
	byToken      map[string]*CompanyTutor
	assignments  map[string]bool
	verifications []TimesheetVerification
	evaluations  []CompanyEvaluation
}

func newMockTutorRepo() *mockTutorRepo {
	return &mockTutorRepo{
		tutors:      make(map[string]*CompanyTutor),
		byToken:     make(map[string]*CompanyTutor),
		assignments: make(map[string]bool),
	}
}

func (m *mockTutorRepo) GetTutorByToken(ctx context.Context, token string) (*CompanyTutor, error) {
	if t, ok := m.byToken[token]; ok {
		return t, nil
	}
	return nil, nil
}

func (m *mockTutorRepo) GetTutorByID(ctx context.Context, id string) (*CompanyTutor, error) {
	if t, ok := m.tutors[id]; ok {
		return t, nil
	}
	return nil, nil
}

func (m *mockTutorRepo) CreateTutor(ctx context.Context, tutor *CompanyTutor) error {
	tutor.ID = "tutor-mock-1"
	tutor.CreatedAt = time.Now()
	tutor.UpdatedAt = time.Now()
	m.tutors[tutor.ID] = tutor
	m.byToken[tutor.AccessToken] = tutor
	return nil
}

func (m *mockTutorRepo) AssignStudent(ctx context.Context, tutorID, projectID, studentID string) error {
	m.assignments[tutorID+"_"+projectID+"_"+studentID] = true
	return nil
}

func (m *mockTutorRepo) ListAssignedStudents(ctx context.Context, tutorID string) ([]AssignedStudentInfo, error) {
	return []AssignedStudentInfo{
		{
			StudentID:      "student-001",
			StudentName:    "Marco Rossi",
			ProjectID:      "proj-001",
			ProjectTitle:   "Stage Sviluppo Software",
			TotalHours:     80,
			CompletedHours: 24,
			IsEvaluated:    false,
		},
	}, nil
}

func (m *mockTutorRepo) SaveTimesheetVerification(ctx context.Context, entry *TimesheetVerification) error {
	entry.ID = "verify-1"
	entry.SignedAt = time.Now()
	m.verifications = append(m.verifications, *entry)
	return nil
}

func (m *mockTutorRepo) SaveCompanyEvaluation(ctx context.Context, eval *CompanyEvaluation) error {
	eval.ID = "eval-1"
	eval.SubmittedAt = time.Now()
	m.evaluations = append(m.evaluations, *eval)
	return nil
}

func setupTutorRouter(repo Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("role", "secretary")
		c.Set("school_id", "school-123")
		c.Set("user_id", "staff-1")
		c.Next()
	})

	svc := NewService(repo)
	handler := NewHandler(svc)
	api := r.Group("/api/v1")
	handler.RegisterRoutes(api)

	return r
}

func TestPCTOTutorWorkflow(t *testing.T) {
	repo := newMockTutorRepo()
	router := setupTutorRouter(repo)

	// Step 1: Register Tutor & Generate Magic Link
	tutorPayload := map[string]interface{}{
		"company_name":     "Tech Solutions S.r.l.",
		"tutor_first_name": "Paolo",
		"tutor_last_name":  "Bianchi",
		"email":            "paolo.bianchi@techsolutions.it",
		"phone":            "+3902123456",
	}
	body, _ := json.Marshal(tutorPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/pcto-tutor/tutors", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on register tutor, got %d: %s", w.Code, w.Body.String())
	}

	var regRes struct {
		Tutor CompanyTutor `json:"tutor"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &regRes)
	token := regRes.Tutor.AccessToken
	if token == "" {
		t.Fatalf("expected magic link access token")
	}

	// Step 2: Tutor accesses session with token
	req = httptest.NewRequest(http.MethodGet, "/api/v1/pcto-tutor/session?token="+token, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on session fetch, got %d", w.Code)
	}

	// Step 3: Tutor lists assigned students
	req = httptest.NewRequest(http.MethodGet, "/api/v1/pcto-tutor/students", nil)
	req.Header.Set("X-Tutor-Token", token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on list students, got %d", w.Code)
	}

	// Step 4: Tutor verifies timesheet hours
	verifyPayload := map[string]interface{}{
		"project_id":     "proj-001",
		"student_id":     "student-001",
		"activity_date":  "2026-10-04",
		"hours_declared": 8.0,
		"hours_approved": 8.0,
		"tutor_notes":    "Attività svolta regolarmente con ottimo impegno",
	}
	body, _ = json.Marshal(verifyPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/pcto-tutor/timesheets/verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tutor-Token", token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on timesheet verification, got %d: %s", w.Code, w.Body.String())
	}

	// Step 5: Tutor submits final evaluation
	evalPayload := map[string]interface{}{
		"project_id":        "proj-001",
		"student_id":        "student-001",
		"reliability_level": 5,
		"technical_skills":  4,
		"teamwork_skills":   5,
		"final_feedback":    "Ottimo stagista, consigliata eventuale assunzione",
	}
	body, _ = json.Marshal(evalPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/pcto-tutor/evaluations", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tutor-Token", token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on evaluation submit, got %d: %s", w.Code, w.Body.String())
	}
}
