package primaryeval

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

type mockService struct {
	objectives []LearningObjective
	matrix     *PrimaryMatrixResponse
	evals      []PrimaryEvaluation
	errToThrow error
}

func (m *mockService) CreateObjective(ctx context.Context, schoolID string, req *CreateObjectiveRequest) (*LearningObjective, error) {
	if m.errToThrow != nil {
		return nil, m.errToThrow
	}
	return &LearningObjective{
		ID:           "obj-created-1",
		SchoolID:     schoolID,
		Title:        req.Title,
		YearGrade:    req.YearGrade,
		SubjectID:    req.SubjectID,
		AcademicYear: req.AcademicYear,
	}, nil
}

func (m *mockService) ListObjectives(ctx context.Context, schoolID, subjectID string, classID *string, yearGrade int, academicYear string) ([]LearningObjective, error) {
	if m.errToThrow != nil {
		return nil, m.errToThrow
	}
	return m.objectives, nil
}

func (m *mockService) DeleteObjective(ctx context.Context, id string) error {
	return m.errToThrow
}

func (m *mockService) SaveEvaluationsBatch(ctx context.Context, schoolID, teacherID string, req *SaveEvaluationsBatchRequest) error {
	return m.errToThrow
}

func (m *mockService) ListEvaluations(ctx context.Context, classID, subjectID string, semester int) ([]PrimaryEvaluation, error) {
	if m.errToThrow != nil {
		return nil, m.errToThrow
	}
	return m.evals, nil
}

func (m *mockService) GetStudentEvaluations(ctx context.Context, studentID string, semester int) ([]PrimaryEvaluation, error) {
	if m.errToThrow != nil {
		return nil, m.errToThrow
	}
	return m.evals, nil
}

func (m *mockService) GetPrimaryMatrix(ctx context.Context, schoolID, classID, subjectID string, semester int) (*PrimaryMatrixResponse, error) {
	if m.errToThrow != nil {
		return nil, m.errToThrow
	}
	if m.matrix != nil {
		return m.matrix, nil
	}
	return &PrimaryMatrixResponse{
		ClassID:   classID,
		SubjectID: subjectID,
		Semester:  semester,
	}, nil
}

func setupTestRouter(svc Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("school_id", "test-school-id")
		c.Set("user_id", "test-teacher-id")
		c.Next()
	})
	h := NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestHandler_ListObjectives(t *testing.T) {
	svc := &mockService{
		objectives: []LearningObjective{
			{ID: "o1", Title: "Obiettivo 1", YearGrade: 1},
			{ID: "o2", Title: "Obiettivo 2", YearGrade: 1},
		},
	}
	r := setupTestRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/primary/objectives?subject_id=sub-1&year_grade=1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string][]LearningObjective
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp["data"]) != 2 {
		t.Fatalf("expected 2 objectives, got %d", len(resp["data"]))
	}
}

func TestHandler_CreateObjective(t *testing.T) {
	svc := &mockService{}
	r := setupTestRouter(svc)

	body := CreateObjectiveRequest{
		SubjectID:    "sub-ita",
		YearGrade:    3,
		Title:        "Capacità di sintesi",
		AcademicYear: "2025/2026",
	}
	jsonBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/primary/objectives", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_SaveEvaluationsBatch(t *testing.T) {
	svc := &mockService{}
	r := setupTestRouter(svc)

	batch := SaveEvaluationsBatchRequest{
		ClassID:     "cls-3b",
		SubjectID:   "sub-mat",
		ObjectiveID: "obj-calcolo",
		Date:        time.Now().Format("2006-01-02"),
		Semester:    1,
		Evaluations: []StudentObjectiveLevelItem{
			{
				StudentID: "s1",
				Level:     LevelAvanzato,
			},
			{
				StudentID: "s2",
				Level:     LevelBase,
			},
		},
	}
	jsonBytes, _ := json.Marshal(batch)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/primary/evaluations/batch", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetMatrix(t *testing.T) {
	svc := &mockService{
		matrix: &PrimaryMatrixResponse{
			ClassID:   "cls-1",
			SubjectID: "sub-1",
			Semester:  2,
			Students: []PrimaryStudentMatrixRow{
				{
					StudentID:   "st-1",
					StudentName: "Test Student",
					Evaluations: map[string]PrimaryEvaluationCell{
						"obj-1": {Level: LevelIntermedio},
					},
				},
			},
		},
	}
	r := setupTestRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/primary/matrix?class_id=cls-1&subject_id=sub-1&semester=2", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp struct {
		Data PrimaryMatrixResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode matrix json: %v", err)
	}
	if resp.Data.Semester != 2 {
		t.Fatalf("expected semester 2, got %d", resp.Data.Semester)
	}
}

func TestHandler_DeleteObjective(t *testing.T) {
	svc := &mockService{}
	r := setupTestRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/primary/objectives/obj-123", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
