package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/credits"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSchoolCreditsWorkflowRepo struct {
	credits map[string]*credits.StudentSchoolCredit
}

func newMockSchoolCreditsWorkflowRepo() *mockSchoolCreditsWorkflowRepo {
	return &mockSchoolCreditsWorkflowRepo{
		credits: make(map[string]*credits.StudentSchoolCredit),
	}
}

func (m *mockSchoolCreditsWorkflowRepo) SaveCredit(ctx context.Context, c *credits.StudentSchoolCredit) error {
	if c.ID == "" {
		c.ID = "credit-" + time.Now().Format("150405.000000")
	}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	m.credits[c.ID] = c
	return nil
}

func (m *mockSchoolCreditsWorkflowRepo) GetCreditByStudentAndYear(ctx context.Context, studentID, academicYear string, gradeLevel int) (*credits.StudentSchoolCredit, error) {
	for _, c := range m.credits {
		if c.StudentID == studentID && c.AcademicYear == academicYear && c.GradeLevel == gradeLevel {
			return c, nil
		}
	}
	return nil, nil
}

func (m *mockSchoolCreditsWorkflowRepo) ListCreditsByClass(ctx context.Context, classID, academicYear string) ([]credits.StudentSchoolCredit, error) {
	var res []credits.StudentSchoolCredit
	for _, c := range m.credits {
		if c.ClassID == classID {
			res = append(res, *c)
		}
	}
	return res, nil
}

func (m *mockSchoolCreditsWorkflowRepo) GetStudentCreditSummary(ctx context.Context, studentID string) (*credits.StudentCreditSummary, error) {
	var list []credits.StudentSchoolCredit
	total := 0
	for _, c := range m.credits {
		if c.StudentID == studentID {
			list = append(list, *c)
			total += c.AssignedCredit
		}
	}
	return &credits.StudentCreditSummary{
		StudentID:            studentID,
		CreditsByYear:        list,
		TotalTrienniumCredit: total,
		MaxPossibleCredit:    40,
	}, nil
}

func setupCreditsTestRouter(repo credits.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := credits.NewService(repo)
	handler := credits.NewHandler(svc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		c.Set("role", "teacher")
		c.Set("user_id", "teacher-coord-1")
		c.Set("school_id", "school-1")
		c.Next()
	})
	handler.RegisterRoutes(api)
	return r
}

func TestIntegration_SchoolCredits_FullWorkflow(t *testing.T) {
	repo := newMockSchoolCreditsWorkflowRepo()
	r := setupCreditsTestRouter(repo)

	// 1. Calculate suggested credit range for a 5th grade student with 8.4 average and conduct 9
	calcReq, _ := http.NewRequest("GET", "/api/v1/credits/calculate?grade_level=5&average=8.4&conduct=9&pcto_hours=90&has_extracurricular=true", nil)
	wCalc := httptest.NewRecorder()
	r.ServeHTTP(wCalc, calcReq)
	require.Equal(t, http.StatusOK, wCalc.Code)

	var calcResult credits.CreditCalculationResult
	err := json.Unmarshal(wCalc.Body.Bytes(), &calcResult)
	require.NoError(t, err)
	assert.Equal(t, 5, calcResult.GradeLevel)
	assert.GreaterOrEqual(t, calcResult.BaseCreditRangeMax, calcResult.BaseCreditRangeMin)

	// 2. Assign deliberated credit to student
	assignPayload := credits.AssignCreditRequest{
		StudentID:          "student-maturita-1",
		ClassID:            "class-5A",
		AcademicYear:       "2025/2026",
		GradeLevel:         5,
		GradeAverage:       8.4,
		ConductGrade:       9,
		AssignedCredit:     calcResult.SuggestedCredit,
		PCTOHours:          90,
		HasExtracurricular: true,
		DeliberationNotes:  "Assegnato il massimo della fascia per eccellente percorso PCTO e attività extracurricolari",
	}
	body, _ := json.Marshal(assignPayload)
	reqAssign, _ := http.NewRequest("POST", "/api/v1/credits/assign", bytes.NewReader(body))
	reqAssign.Header.Set("Content-Type", "application/json")
	wAssign := httptest.NewRecorder()
	r.ServeHTTP(wAssign, reqAssign)
	require.Equal(t, http.StatusOK, wAssign.Code)

	var savedCredit credits.StudentSchoolCredit
	err = json.Unmarshal(wAssign.Body.Bytes(), &savedCredit)
	require.NoError(t, err)
	assert.NotEmpty(t, savedCredit.ID)
	assert.Equal(t, "student-maturita-1", savedCredit.StudentID)
	assert.Equal(t, calcResult.SuggestedCredit, savedCredit.AssignedCredit)

	// 3. Get Student credit summary (triennio total)
	reqSummary, _ := http.NewRequest("GET", "/api/v1/credits/student/student-maturita-1/summary", nil)
	wSummary := httptest.NewRecorder()
	r.ServeHTTP(wSummary, reqSummary)
	require.Equal(t, http.StatusOK, wSummary.Code)

	var summary credits.StudentCreditSummary
	err = json.Unmarshal(wSummary.Body.Bytes(), &summary)
	require.NoError(t, err)
	assert.Equal(t, "student-maturita-1", summary.StudentID)
	assert.Equal(t, savedCredit.AssignedCredit, summary.TotalTrienniumCredit)
}
