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

// ─── Mock Credits Repository ───────────────────────────────────────────────

type mockTrienniumCreditsRepo struct {
	credits map[string]credits.StudentSchoolCredit
}

func newMockTrienniumCreditsRepo() *mockTrienniumCreditsRepo {
	return &mockTrienniumCreditsRepo{
		credits: make(map[string]credits.StudentSchoolCredit),
	}
}

func (m *mockTrienniumCreditsRepo) SaveCredit(_ context.Context, c *credits.StudentSchoolCredit) error {
	if c.ID == "" {
		c.ID = "cred-" + time.Now().Format("150405.000")
	}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	m.credits[c.StudentID+"-"+c.AcademicYear] = *c
	return nil
}

func (m *mockTrienniumCreditsRepo) GetCreditByStudentAndYear(_ context.Context, studentID, academicYear string, _ int) (*credits.StudentSchoolCredit, error) {
	c, ok := m.credits[studentID+"-"+academicYear]
	if !ok {
		return nil, nil
	}
	return &c, nil
}

func (m *mockTrienniumCreditsRepo) ListCreditsByClass(_ context.Context, classID, academicYear string) ([]credits.StudentSchoolCredit, error) {
	var list []credits.StudentSchoolCredit
	for _, c := range m.credits {
		if classID != "" && c.ClassID != classID {
			continue
		}
		if academicYear != "" && c.AcademicYear != academicYear {
			continue
		}
		list = append(list, c)
	}
	return list, nil
}

func (m *mockTrienniumCreditsRepo) GetStudentCreditSummary(_ context.Context, studentID string) (*credits.StudentCreditSummary, error) {
	var byYear []credits.StudentSchoolCredit
	total := 0
	for _, c := range m.credits {
		if c.StudentID == studentID {
			byYear = append(byYear, c)
			total += c.AssignedCredit
		}
	}
	return &credits.StudentCreditSummary{
		StudentID:            studentID,
		CreditsByYear:        byYear,
		TotalTrienniumCredit: total,
		MaxPossibleCredit:    40,
	}, nil
}

// ─── Setup Router Helper ───────────────────────────────────────────────────

func setupCreditsRouter(repo *mockTrienniumCreditsRepo, role, userID, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := credits.NewService(repo)
	h := credits.NewHandler(svc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
			c.Set("role", role)
			c.Set("school_id", schoolID)
		}
		c.Next()
	})
	h.RegisterRoutes(api)
	return r
}

// ─── Tests ─────────────────────────────────────────────────────────────────

// CRD01 — Calculate suggested ministerial credit range for 3rd year with average 8.2
func TestCredits_Calculate_Grade3_Average8(t *testing.T) {
	repo := newMockTrienniumCreditsRepo()
	r := setupCreditsRouter(repo, "teacher", "teacher-1", "school-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/credits/calculate?grade_level=3&average=8.2&conduct=9&pcto_hours=30&has_extracurricular=true", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var calc credits.CreditCalculationResult
	err := json.Unmarshal(w.Body.Bytes(), &calc)
	require.NoError(t, err)
	assert.Equal(t, 3, calc.GradeLevel)
	assert.True(t, calc.BaseCreditRangeMin >= 10)
	assert.True(t, calc.SuggestedCredit >= calc.BaseCreditRangeMin)
}

// CRD02 — Assign credit within allowed range during scrutiny
func TestCredits_Assign_Success(t *testing.T) {
	repo := newMockTrienniumCreditsRepo()
	r := setupCreditsRouter(repo, "teacher", "coordinator-1", "school-1")

	body, _ := json.Marshal(credits.AssignCreditRequest{
		StudentID:          "student-triennio-1",
		ClassID:            "class-5A",
		AcademicYear:       "2025/2026",
		GradeLevel:         5,
		GradeAverage:       9.5,
		ConductGrade:       9,
		AssignedCredit:     14,
		PCTOHours:          90,
		HasExtracurricular: true,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/credits/assign", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var credit credits.StudentSchoolCredit
	err := json.Unmarshal(w.Body.Bytes(), &credit)
	require.NoError(t, err)
	assert.Equal(t, 14, credit.AssignedCredit)
	assert.Equal(t, 5, credit.GradeLevel)
}

// CRD03 — Assign credit out-of-range without deliberation notes is rejected
func TestCredits_Assign_OutOfRangeRequiresDeliberation(t *testing.T) {
	repo := newMockTrienniumCreditsRepo()
	r := setupCreditsRouter(repo, "teacher", "coordinator-1", "school-1")

	// Grade level 3 with average 6.0 allows max 8 credits, attempting 12 without notes
	body, _ := json.Marshal(credits.AssignCreditRequest{
		StudentID:          "student-triennio-2",
		ClassID:            "class-3A",
		AcademicYear:       "2025/2026",
		GradeLevel:         3,
		GradeAverage:       6.0,
		ConductGrade:       7,
		AssignedCredit:     12,
		DeliberationNotes:  "", // Empty notes -> must fail
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/credits/assign", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// CRD04 — Get student summary aggregating triennium credits (Years 3, 4, 5)
func TestCredits_GetStudentSummary(t *testing.T) {
	repo := newMockTrienniumCreditsRepo()
	_ = repo.SaveCredit(context.Background(), &credits.StudentSchoolCredit{
		StudentID:      "student-5A-maturita",
		AcademicYear:   "2023/2024",
		GradeLevel:     3,
		AssignedCredit: 11,
	})
	_ = repo.SaveCredit(context.Background(), &credits.StudentSchoolCredit{
		StudentID:      "student-5A-maturita",
		AcademicYear:   "2024/2025",
		GradeLevel:     4,
		AssignedCredit: 12,
	})
	_ = repo.SaveCredit(context.Background(), &credits.StudentSchoolCredit{
		StudentID:      "student-5A-maturita",
		AcademicYear:   "2025/2026",
		GradeLevel:     5,
		AssignedCredit: 14,
	})

	r := setupCreditsRouter(repo, "teacher", "coordinator-1", "school-1")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/credits/student/student-5A-maturita/summary", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var summary credits.StudentCreditSummary
	err := json.Unmarshal(w.Body.Bytes(), &summary)
	require.NoError(t, err)
	assert.Equal(t, 37, summary.TotalTrienniumCredit) // 11 + 12 + 14 = 37 / 40
	assert.Equal(t, 40, summary.MaxPossibleCredit)
	assert.Len(t, summary.CreditsByYear, 3)
}
