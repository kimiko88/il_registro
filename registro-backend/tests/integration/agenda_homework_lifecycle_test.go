package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/agenda"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockAgendaRepo struct {
	items map[string]*agenda.AgendaItem
}

func (m *mockAgendaRepo) Create(ctx context.Context, item *agenda.AgendaItem) error {
	if item.ID == "" {
		item.ID = "agenda-1"
	}
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	m.items[item.ID] = item
	return nil
}

func (m *mockAgendaRepo) GetByID(ctx context.Context, id string) (*agenda.AgendaItem, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, agenda.ErrNotFound
	}
	return item, nil
}

func (m *mockAgendaRepo) Update(ctx context.Context, item *agenda.AgendaItem) error {
	m.items[item.ID] = item
	return nil
}

func (m *mockAgendaRepo) Delete(ctx context.Context, id string) error {
	delete(m.items, id)
	return nil
}

func (m *mockAgendaRepo) ListCalendar(ctx context.Context, schoolID string, filter agenda.CalendarFilter) ([]*agenda.AgendaItem, error) {
	result := make([]*agenda.AgendaItem, 0)
	for _, item := range m.items {
		result = append(result, item)
	}
	return result, nil
}

func (m *mockAgendaRepo) SetCompletion(ctx context.Context, itemID, studentID string, completed bool) error {
	return nil
}

func (m *mockAgendaRepo) IsStudentInClass(ctx context.Context, studentID, classID string) (bool, error) {
	return true, nil
}

func TestIntegration_Agenda_Homework_Lifecycle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockAgendaRepo{items: make(map[string]*agenda.AgendaItem)}
	svc := agenda.NewService(repo)
	handler := agenda.NewHandler(svc)

	r := gin.New()
	r.POST("/agenda", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.Create(c)
	})
	r.POST("/agenda/:id/completion", func(c *gin.Context) {
		c.Set("user_id", "student-1")
		c.Set("role", "student")
		c.Set("school_id", "school-1")
		handler.MarkComplete(c)
	})

	// 1. Create agenda homework item with valid date
	itemReq := map[string]interface{}{
		"class_id":    "class-1",
		"title":       "Compiti di Matematica: Equazioni di 2° grado",
		"description": "Risolvere gli esercizi da pag 45 a pag 48",
		"type":        "homework",
		"date":        "2026-09-10",
		"all_day":     true,
	}
	body, _ := json.Marshal(itemReq)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/agenda", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// 2. Student toggles completion for the assigned homework
	w2 := httptest.NewRecorder()
	complReq := map[string]interface{}{
		"completed": true,
	}
	body2, _ := json.Marshal(complReq)
	req2, _ := http.NewRequest("POST", "/agenda/agenda-1/completion", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}
