package religion_alternative

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

type mockRepo struct {
	options     map[string]*StudentReligionOption
	evaluations map[string]*AlternativeEvaluation
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		options:     make(map[string]*StudentReligionOption),
		evaluations: make(map[string]*AlternativeEvaluation),
	}
}

func (m *mockRepo) SaveOption(ctx context.Context, opt *StudentReligionOption) error {
	opt.ID = "opt-mock-1"
	opt.ChosenAt = time.Now()
	m.options[opt.StudentID+"_"+opt.AcademicYear] = opt
	return nil
}

func (m *mockRepo) GetOption(ctx context.Context, studentID, academicYear string) (*StudentReligionOption, error) {
	if o, ok := m.options[studentID+"_"+academicYear]; ok {
		return o, nil
	}
	return nil, nil
}

func (m *mockRepo) ListOptionsBySchool(ctx context.Context, schoolID, academicYear string) ([]StudentReligionOption, error) {
	var list []StudentReligionOption
	for _, o := range m.options {
		if o.SchoolID == schoolID && o.AcademicYear == academicYear {
			list = append(list, *o)
		}
	}
	return list, nil
}

func (m *mockRepo) SaveEvaluation(ctx context.Context, eval *AlternativeEvaluation) error {
	eval.ID = "eval-mock-1"
	eval.CreatedAt = time.Now()
	eval.UpdatedAt = time.Now()
	m.evaluations[eval.StudentID+"_"+eval.Period+"_"+eval.SubjectKind] = eval
	return nil
}

func (m *mockRepo) ListEvaluations(ctx context.Context, schoolID, period, subjectKind string) ([]AlternativeEvaluation, error) {
	var list []AlternativeEvaluation
	for _, e := range m.evaluations {
		if e.SchoolID == schoolID {
			if (period == "" || e.Period == period) && (subjectKind == "" || e.SubjectKind == subjectKind) {
				list = append(list, *e)
			}
		}
	}
	return list, nil
}

func setupTestRouter(repo Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("role", "teacher")
		c.Set("school_id", "school-test-123")
		c.Set("user_id", "user-teacher-1")
		c.Next()
	})

	svc := NewService(repo)
	handler := NewHandler(svc)
	api := r.Group("/api/v1")
	handler.RegisterRoutes(api)

	return r
}

func TestReligionAlternativeIntegrationFlow(t *testing.T) {
	repo := newMockRepo()
	router := setupTestRouter(repo)

	// Step 1: Save Option for Student
	optionPayload := map[string]interface{}{
		"student_id":    "student-001",
		"academic_year": "2026/2027",
		"option_type":   OptionMateriaAlternativa,
		"notes":         "Scelta concordata con la famiglia",
	}
	body, _ := json.Marshal(optionPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/religion-alternative/options", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	// Step 2: List Options
	req = httptest.NewRequest(http.MethodGet, "/api/v1/religion-alternative/options?academic_year=2026/2027", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 on list options, got %d", w.Code)
	}

	// Step 3: Record Evaluation with synthetic judgment
	evalPayload := map[string]interface{}{
		"student_id":        "student-001",
		"period":            "q1",
		"subject_kind":      "materia_alternativa",
		"judgment_level":    JudgmentOttimo,
		"descriptive_notes": "Partecipazione attiva e approfondimento eccellente",
	}
	body, _ = json.Marshal(evalPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/religion-alternative/evaluations", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 on save evaluation, got %d: %s", w.Code, w.Body.String())
	}

	// Step 4: List evaluations
	req = httptest.NewRequest(http.MethodGet, "/api/v1/religion-alternative/evaluations?period=q1&subject_kind=materia_alternativa", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 on list evaluations, got %d", w.Code)
	}
}
