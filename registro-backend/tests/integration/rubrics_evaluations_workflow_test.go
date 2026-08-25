package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/rubrics"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockRubricsRepo struct {
	rubrics map[string]*rubrics.Rubric
	assess  map[string]*rubrics.RubricAssessment
}

func (m *mockRubricsRepo) CreateRubric(ctx context.Context, r *rubrics.Rubric) error {
	if r.ID == "" {
		r.ID = "rubric-1"
	}
	r.CreatedAt = time.Now()
	m.rubrics[r.ID] = r
	return nil
}

func (m *mockRubricsRepo) GetRubricByID(ctx context.Context, id string) (*rubrics.Rubric, error) {
	r, ok := m.rubrics[id]
	if !ok {
		return nil, errors.New("rubric not found")
	}
	return r, nil
}

func (m *mockRubricsRepo) ListRubrics(ctx context.Context, schoolID, teacherID string) ([]*rubrics.Rubric, error) {
	res := make([]*rubrics.Rubric, 0)
	for _, r := range m.rubrics {
		res = append(res, r)
	}
	return res, nil
}

func (m *mockRubricsRepo) UpdateRubric(ctx context.Context, r *rubrics.Rubric) error {
	m.rubrics[r.ID] = r
	return nil
}

func (m *mockRubricsRepo) DeleteRubric(ctx context.Context, id string) error {
	delete(m.rubrics, id)
	return nil
}

func (m *mockRubricsRepo) CreateAssessment(ctx context.Context, a *rubrics.RubricAssessment) error {
	if a.ID == "" {
		a.ID = "assess-1"
	}
	a.CreatedAt = time.Now()
	m.assess[a.ID] = a
	return nil
}

func (m *mockRubricsRepo) ListAssessmentsByStudent(ctx context.Context, studentID string) ([]*rubrics.RubricAssessment, error) {
	res := make([]*rubrics.RubricAssessment, 0)
	for _, a := range m.assess {
		if a.StudentID == studentID {
			res = append(res, a)
		}
	}
	return res, nil
}

func (m *mockRubricsRepo) ListAssessmentsByClass(ctx context.Context, classID string) ([]*rubrics.RubricAssessment, error) {
	res := make([]*rubrics.RubricAssessment, 0)
	for _, a := range m.assess {
		if a.ClassID == classID {
			res = append(res, a)
		}
	}
	return res, nil
}

func TestIntegration_Rubrics_Evaluations_Workflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockRubricsRepo{
		rubrics: make(map[string]*rubrics.Rubric),
		assess:  make(map[string]*rubrics.RubricAssessment),
	}
	svc := rubrics.NewService(repo)
	handler := rubrics.NewHandler(svc)

	r := gin.New()
	r.POST("/rubrics", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.Create(c)
	})
	r.POST("/rubrics/:id/assessments", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.AssessStudent(c)
	})

	// 1. Teacher creates assessment rubric template
	rubricReq := rubrics.CreateRubricRequest{
		SubjectID:   "subj-1",
		Title:       "Rubrica Valutativa Presentazione Orale",
		Description: "Criteri di valutazione per espozione orale e chiarezza espositiva",
		Criteria: []rubrics.CreateCriterionInput{
			{
				Name:        "Chiarezza Espositiva",
				Description: "Proprietà di linguaggio e fluidità del discorso",
				MaxScore:    10.0,
				Levels: []rubrics.CreateLevelInput{
					{Score: 10.0, Label: "Eccellente", Description: "Esposizione fluida e ricca"},
				},
			},
		},
	}
	body1, _ := json.Marshal(rubricReq)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/rubrics", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// 2. Teacher assesses student using rubric
	assessReq := rubrics.CreateAssessmentRequest{
		StudentID: "st-1",
		ClassID:   "class-1",
		Date:      "2026-11-10",
		Scores: []rubrics.CriterionScore{
			{
				CriterionID: "crit-1",
				LevelID:     "lvl-1",
				Score:       9.0,
			},
		},
		Notes: "Eccellente esposizione",
	}
	body2, _ := json.Marshal(assessReq)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/rubrics/rubric-1/assessments", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusCreated, w2.Code)
}
