package scrutiny

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDeficiencyRepo struct {
	mock.Mock
}

func (m *mockDeficiencyRepo) SaveRecord(ctx context.Context, record *ScrutinyRecord) error {
	return nil
}
func (m *mockDeficiencyRepo) GetRecord(ctx context.Context, studentID, classID string, semester int) (*ScrutinyRecord, error) {
	return nil, nil
}
func (m *mockDeficiencyRepo) ListRecordsByClass(ctx context.Context, classID string, semester int) ([]ScrutinyRecord, error) {
	return nil, nil
}
func (m *mockDeficiencyRepo) ValidateClassScrutiny(ctx context.Context, classID string, semester int, validatorID string) error {
	return nil
}
func (m *mockDeficiencyRepo) UpdateClassScrutinyStatus(ctx context.Context, classID string, semester int, status string) error {
	return nil
}

func (m *mockDeficiencyRepo) SaveDeficiency(ctx context.Context, def *StudentDeficiency) error {
	args := m.Called(ctx, def)
	return args.Error(0)
}
func (m *mockDeficiencyRepo) GetDeficienciesByStudent(ctx context.Context, studentID string) ([]StudentDeficiency, error) {
	args := m.Called(ctx, studentID)
	return args.Get(0).([]StudentDeficiency), args.Error(1)
}
func (m *mockDeficiencyRepo) GetDeficienciesByClass(ctx context.Context, classID string, semester int) ([]StudentDeficiency, error) {
	args := m.Called(ctx, classID, semester)
	return args.Get(0).([]StudentDeficiency), args.Error(1)
}
func (m *mockDeficiencyRepo) SaveDeferredScrutiny(ctx context.Context, req *SaveDeferredScrutinyRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func TestStudentDeficiency_StructValidation(t *testing.T) {
	gradeVal := 7.5
	def := StudentDeficiency{
		ID:            "def-101",
		SchoolID:      "school-01",
		StudentID:     "stud-01",
		ClassID:       "class-1a",
		SubjectID:     "sub-math",
		SubjectName:   "Matematica",
		Semester:      1,
		PeriodType:    "semester_1",
		Topics:        "Trigonometria, Funzioni esponenziali",
		RecoveryMode:  "corso_recupero",
		Status:        "da_recuperare",
		RecoveryGrade: &gradeVal,
	}

	assert.Equal(t, "def-101", def.ID)
	assert.Equal(t, "school-01", def.SchoolID)
	assert.Equal(t, "stud-01", def.StudentID)
	assert.Equal(t, "class-1a", def.ClassID)
	assert.Equal(t, "sub-math", def.SubjectID)
	assert.Equal(t, "Matematica", def.SubjectName)
	assert.Equal(t, 1, def.Semester)
	assert.Equal(t, "semester_1", def.PeriodType)
	assert.Equal(t, "Trigonometria, Funzioni esponenziali", def.Topics)
	assert.Equal(t, "corso_recupero", def.RecoveryMode)
	assert.Equal(t, "da_recuperare", def.Status)
	assert.Equal(t, 7.5, *def.RecoveryGrade)
}

func TestSaveDeferredScrutinyRequest_Validation(t *testing.T) {
	recGrade := 6.0
	req := SaveDeferredScrutinyRequest{
		StudentID:     "stud-01",
		ClassID:       "class-1a",
		FinalDecision: "promosso_con_debiti_saldati",
		Notes:         "Debito di Matematica saldato con successo",
		Deficiencies: []SaveDeferredDeficiencyItemReq{
			{
				DeficiencyID:  "def-101",
				Status:        "recuperato",
				RecoveryGrade: &recGrade,
			},
		},
	}

	assert.Equal(t, "stud-01", req.StudentID)
	assert.Equal(t, "class-1a", req.ClassID)
	assert.Equal(t, "promosso_con_debiti_saldati", req.FinalDecision)
	assert.Equal(t, "Debito di Matematica saldato con successo", req.Notes)
	assert.Len(t, req.Deficiencies, 1)
	assert.Equal(t, "recuperato", req.Deficiencies[0].Status)
	assert.Equal(t, 6.0, *req.Deficiencies[0].RecoveryGrade)
}

func TestHandler_GetStudentDeficiencies_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := NewService(nil, nil, nil, nil, nil)
	h := NewHandler(svc)

	r.GET("/api/v1/scrutiny/deficiencies/student/:studentId", h.GetStudentDeficiencies)

	req := httptest.NewRequest("GET", "/api/v1/scrutiny/deficiencies/student/stud-01", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandler_SaveDeficiency_MissingParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := NewService(nil, nil, nil, nil, nil)
	h := NewHandler(svc)

	r.POST("/api/v1/scrutiny/deficiencies", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		h.SaveDeficiency(c)
	})

	// Missing required fields (student_id, class_id, subject_id)
	body, _ := json.Marshal(SaveDeficiencyRequest{Topics: "Tutti gli argomenti"})
	req := httptest.NewRequest("POST", "/api/v1/scrutiny/deficiencies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_SaveDeficiency_ForbiddenForStudent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := NewService(nil, nil, nil, nil, nil)
	h := NewHandler(svc)

	r.POST("/api/v1/scrutiny/deficiencies", func(c *gin.Context) {
		c.Set("user_id", "stud-1")
		c.Set("role", "student")
		h.SaveDeficiency(c)
	})

	body, _ := json.Marshal(SaveDeficiencyRequest{StudentID: "stud-1", ClassID: "c1", SubjectID: "s1"})
	req := httptest.NewRequest("POST", "/api/v1/scrutiny/deficiencies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestHandler_GetStudentDeficiencies_ForbiddenForOtherStudent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := NewService(nil, nil, nil, nil, nil)
	h := NewHandler(svc)

	r.GET("/api/v1/scrutiny/deficiencies/student/:studentId", func(c *gin.Context) {
		c.Set("user_id", "stud-1")
		c.Set("role", "student")
		h.GetStudentDeficiencies(c)
	})

	req := httptest.NewRequest("GET", "/api/v1/scrutiny/deficiencies/student/stud-other", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
