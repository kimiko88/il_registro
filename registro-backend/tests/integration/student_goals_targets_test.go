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

	"registro-backend/internal/student_goals"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockGoalsRepo struct {
	goals map[string]*student_goals.StudentGoal
}

func (m *mockGoalsRepo) Create(ctx context.Context, g *student_goals.StudentGoal) error {
	if g.ID == "" {
		g.ID = "goal-1"
	}
	g.CreatedAt = time.Now()
	if g.Status == "" {
		g.Status = student_goals.StatusPending
	}
	m.goals[g.ID] = g
	return nil
}

func (m *mockGoalsRepo) GetByID(ctx context.Context, id string) (*student_goals.StudentGoal, error) {
	g, ok := m.goals[id]
	if !ok {
		return nil, errors.New("goal not found")
	}
	return g, nil
}

func (m *mockGoalsRepo) ListByStudent(ctx context.Context, studentID string) ([]*student_goals.StudentGoal, error) {
	res := make([]*student_goals.StudentGoal, 0)
	for _, g := range m.goals {
		if g.StudentID == studentID {
			res = append(res, g)
		}
	}
	return res, nil
}

func (m *mockGoalsRepo) UpdateStatus(ctx context.Context, id string, status student_goals.GoalStatus) error {
	if g, ok := m.goals[id]; ok {
		g.Status = status
		if status == student_goals.StatusCompleted {
			now := time.Now()
			g.CompletedAt = &now
		}
	}
	return nil
}

func TestIntegration_Student_Goals_Targets(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockGoalsRepo{goals: make(map[string]*student_goals.StudentGoal)}
	svc := student_goals.NewService(repo)
	handler := student_goals.NewHandler(svc)

	r := gin.New()
	r.POST("/student-goals", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.Create(c)
	})
	r.PATCH("/student-goals/:id/status", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.UpdateStatus(c)
	})

	// 1. Teacher sets learning goal for student
	goalReq := student_goals.CreateGoalRequest{
		StudentID:   "st-1",
		Title:       "Miglioramento Comprensione del Testo",
		Description: "Leggere 3 testi argomentativi al mese con scheda di sintesi",
		Category:    "academic",
		Points:      50,
		DueDate:     "2026-11-30",
	}
	body1, _ := json.Marshal(goalReq)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/student-goals", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// 2. Teacher updates goal status to completed
	statusReq := student_goals.UpdateGoalStatusRequest{
		Status: student_goals.StatusCompleted,
	}
	body2, _ := json.Marshal(statusReq)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("PATCH", "/student-goals/goal-1/status", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}
