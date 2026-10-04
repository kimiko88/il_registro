package middleschoolexam

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

type mockExamRepo struct {
	exams      map[string]*MiddleSchoolExam
	candidates map[string]*MiddleSchoolExamCandidate
}

func newMockExamRepo() *mockExamRepo {
	return &mockExamRepo{
		exams:      make(map[string]*MiddleSchoolExam),
		candidates: make(map[string]*MiddleSchoolExamCandidate),
	}
}

func (m *mockExamRepo) GetOrCreateExamForClass(ctx context.Context, schoolID, classID, academicYear, president string) (*MiddleSchoolExam, error) {
	e := &MiddleSchoolExam{
		ID:                  "exam-1",
		SchoolID:            schoolID,
		ClassID:             classID,
		AcademicYear:        academicYear,
		SubcommissionNumber: 1,
		PresidentName:       president,
		Status:              "admission",
		CreatedAt:           time.Now(),
	}
	m.exams[e.ID] = e
	return e, nil
}

func (m *mockExamRepo) GetExamByID(ctx context.Context, id string) (*MiddleSchoolExam, error) {
	return m.exams[id], nil
}

func (m *mockExamRepo) UpdateExamStatus(ctx context.Context, id, status string) error {
	if e, ok := m.exams[id]; ok {
		e.Status = status
	}
	return nil
}

func (m *mockExamRepo) ListCandidates(ctx context.Context, examID string) ([]MiddleSchoolExamCandidate, error) {
	var list []MiddleSchoolExamCandidate
	for _, c := range m.candidates {
		if c.ExamID == examID {
			list = append(list, *c)
		}
	}
	return list, nil
}

func (m *mockExamRepo) SaveCandidateAdmission(ctx context.Context, examID, studentID string, grade int, judgment string, admitted bool) error {
	key := examID + ":" + studentID
	m.candidates[key] = &MiddleSchoolExamCandidate{
		ID:                "cand-1",
		ExamID:            examID,
		StudentID:         studentID,
		StudentName:       "Giacomo Leopardi",
		AdmissionGrade:    grade,
		AdmissionJudgment: judgment,
		IsAdmitted:        admitted,
		CreatedAt:         time.Now(),
	}
	return nil
}

func (m *mockExamRepo) SaveCandidateGrades(ctx context.Context, c *MiddleSchoolExamCandidate) error {
	key := c.ExamID + ":" + c.StudentID
	m.candidates[key] = c
	return nil
}

func (m *mockExamRepo) GetCandidate(ctx context.Context, examID, studentID string) (*MiddleSchoolExamCandidate, error) {
	key := examID + ":" + studentID
	if c, ok := m.candidates[key]; ok {
		return c, nil
	}
	return &MiddleSchoolExamCandidate{
		ID:             "cand-1",
		ExamID:         examID,
		StudentID:      studentID,
		StudentName:    "Giacomo Leopardi",
		AdmissionGrade: 10,
		IsAdmitted:     true,
	}, nil
}

func (m *mockExamRepo) GetCandidateDiplomaData(ctx context.Context, candidateID string) (*DiplomaData, error) {
	return &DiplomaData{
		SchoolName:      "Scuola Secondaria Statale Dante Alighieri",
		SchoolCode:      "RMIC801004",
		AcademicYear:    "2026/2027",
		StudentFullName: "Giacomo Leopardi",
		BirthDate:       "29/06/2012",
		BirthPlace:      "Recanati",
		FinalGrade:      10,
		HasHonors:       true,
		PresidentName:   "Prof.ssa Maria Rossi",
		IssueDate:       time.Now(),
	}, nil
}

func setupExamTestRouter(svc *Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	rg := r.Group("/api/v1")
	rg.Use(func(c *gin.Context) {
		c.Set("user_id", "teacher-user-1")
		c.Set("role", "coordinator")
		c.Set("school_id", "school-1")
		c.Next()
	})
	h := NewHandler(svc)
	h.RegisterRoutes(rg)
	return r
}

func TestMiddleSchoolExam_Integration_FullWorkflow(t *testing.T) {
	repo := newMockExamRepo()
	svc := NewService(repo)
	router := setupExamTestRouter(svc)

	// 1. Get or create exam commission
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/middle-school-exam/classes/class-3a", bytes.NewReader([]byte(`{"president_name":"Prof.ssa Rossi"}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on exam setup, got %d: %s", w.Code, w.Body.String())
	}
	var exam MiddleSchoolExam
	_ = json.Unmarshal(w.Body.Bytes(), &exam)
	if exam.ID != "exam-1" {
		t.Errorf("expected exam-1, got %s", exam.ID)
	}

	// 2. Save Admission
	admissionPayload := `{"grade":10,"judgment":"Profilo eccellente","admitted":true}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/middle-school-exam/exam-1/candidates/student-1/admission", bytes.NewReader([]byte(admissionPayload)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on admission, got %d: %s", w.Code, w.Body.String())
	}

	// 3. Evaluate candidate with 5 test grades + proposed honors
	evalPayload := `{
		"grades": {
			"grade_italian": 10,
			"grade_math": 10,
			"grade_english": 10,
			"grade_second_lang": 10,
			"grade_interview": 10,
			"proposed_honors": true
		},
		"notes": "Colloquio interdisciplinare brillante"
	}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/middle-school-exam/exam-1/candidates/student-1/evaluate", bytes.NewReader([]byte(evalPayload)))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on evaluate, got %d: %s", w.Code, w.Body.String())
	}
	var evaluated MiddleSchoolExamCandidate
	_ = json.Unmarshal(w.Body.Bytes(), &evaluated)
	if evaluated.FinalGrade != 10 || !evaluated.HasHonors || evaluated.Outcome != "licenziato" {
		t.Errorf("unexpected evaluation: %+v", evaluated)
	}

	// 4. Download diploma certificate text
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/middle-school-exam/candidates/cand-1/diploma", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on diploma, got %d", w.Code)
	}
	diplomaText := w.Body.String()
	if !bytes.Contains([]byte(diplomaText), []byte("DIPLOMA CONCLUSIVO DEL PRIMO CICLO")) {
		t.Errorf("diploma text missing title header")
	}
	if !bytes.Contains([]byte(diplomaText), []byte("10/10 E LODE")) {
		t.Errorf("diploma text missing 10/10 E LODE")
	}
}
