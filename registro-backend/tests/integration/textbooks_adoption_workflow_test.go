package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/textbooks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockTextbooksRepo struct {
	textbooks   map[string]*textbooks.Textbook
	assignments map[string][]textbooks.ClassTextbook
}

func (m *mockTextbooksRepo) Create(ctx context.Context, t *textbooks.Textbook) error {
	if t.ID == "" {
		t.ID = "tb-1"
	}
	t.CreatedAt = time.Now()
	m.textbooks[t.ID] = t
	return nil
}

func (m *mockTextbooksRepo) Update(ctx context.Context, t *textbooks.Textbook) error {
	m.textbooks[t.ID] = t
	return nil
}

func (m *mockTextbooksRepo) List(ctx context.Context, schoolID string) ([]textbooks.Textbook, error) {
	res := make([]textbooks.Textbook, 0)
	for _, t := range m.textbooks {
		res = append(res, *t)
	}
	return res, nil
}

func (m *mockTextbooksRepo) Delete(ctx context.Context, id string) error {
	delete(m.textbooks, id)
	return nil
}

func (m *mockTextbooksRepo) AssignToClass(ctx context.Context, classID, subjectID, textbookID string, optional bool) error {
	m.assignments[classID] = append(m.assignments[classID], textbooks.ClassTextbook{
		ID:         "asgn-1",
		ClassID:    classID,
		SubjectID:  subjectID,
		TextbookID: textbookID,
		IsOptional: optional,
	})
	return nil
}

func (m *mockTextbooksRepo) RemoveFromClass(ctx context.Context, assignmentID string) error {
	return nil
}

func (m *mockTextbooksRepo) ListByClass(ctx context.Context, classID string) ([]textbooks.ClassTextbook, error) {
	return m.assignments[classID], nil
}

func TestIntegration_Textbooks_Adoption_Workflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockTextbooksRepo{
		textbooks:   make(map[string]*textbooks.Textbook),
		assignments: make(map[string][]textbooks.ClassTextbook),
	}
	svc := textbooks.NewService(repo)
	handler := textbooks.NewHandler(svc)

	r := gin.New()
	r.POST("/textbooks", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.Create(c)
	})
	r.POST("/textbooks/class/:classId", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.AssignToClass(c)
	})

	// 1. Register textbook in school library
	tbReq := textbooks.CreateTextbookRequest{
		Title:     "Matematica.blu 2.0",
		Author:    "Massimo Bergamini",
		Subject:   "Matematica",
		ISBN:      "9788808930438",
		Publisher: "Zanichelli",
		Price:     32.50,
	}
	body1, _ := json.Marshal(tbReq)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/textbooks", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// 2. Adopt textbook for Class 2A
	asgnReq := textbooks.AssignTextbookRequest{
		SubjectID:  "subj-1",
		TextbookID: "tb-1",
		IsOptional: false,
	}
	body2, _ := json.Marshal(asgnReq)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/textbooks/class/class-1", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusCreated, w2.Code)
}
