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

	"registro-backend/internal/groups"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockGroupsRepo struct {
	groupMap map[string]*groups.Group
	members  map[string][]string
}

func (m *mockGroupsRepo) Create(ctx context.Context, g *groups.Group) error {
	if g.ID == "" {
		g.ID = "group-1"
	}
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
	m.groupMap[g.ID] = g
	return nil
}

func (m *mockGroupsRepo) Update(ctx context.Context, g *groups.Group) error {
	m.groupMap[g.ID] = g
	return nil
}

func (m *mockGroupsRepo) Delete(ctx context.Context, id string) error {
	delete(m.groupMap, id)
	return nil
}

func (m *mockGroupsRepo) GetByID(ctx context.Context, id string) (*groups.Group, error) {
	g, ok := m.groupMap[id]
	if !ok {
		return nil, errors.New("group not found")
	}
	return g, nil
}

func (m *mockGroupsRepo) ListBySchool(ctx context.Context, schoolID string) ([]groups.Group, error) {
	res := make([]groups.Group, 0)
	for _, g := range m.groupMap {
		res = append(res, *g)
	}
	return res, nil
}

func (m *mockGroupsRepo) ListByTeacher(ctx context.Context, teacherID string) ([]groups.Group, error) {
	res := make([]groups.Group, 0)
	for _, g := range m.groupMap {
		res = append(res, *g)
	}
	return res, nil
}

func (m *mockGroupsRepo) ListByStudent(ctx context.Context, studentID string) ([]groups.Group, error) {
	res := make([]groups.Group, 0)
	for _, g := range m.groupMap {
		res = append(res, *g)
	}
	return res, nil
}

func (m *mockGroupsRepo) AddStudents(ctx context.Context, groupID string, studentIDs []string) error {
	m.members[groupID] = append(m.members[groupID], studentIDs...)
	return nil
}

func (m *mockGroupsRepo) RemoveStudent(ctx context.Context, groupID, studentID string) error {
	return nil
}

func (m *mockGroupsRepo) GetStudentsInGroup(ctx context.Context, groupID string) ([]groups.GroupStudentInfo, error) {
	return []groups.GroupStudentInfo{}, nil
}

func TestIntegration_Groups_Assignment_Workflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockGroupsRepo{
		groupMap: make(map[string]*groups.Group),
		members:  make(map[string][]string),
	}
	svc := groups.NewService(repo)
	handler := groups.NewHandler(svc)

	r := gin.New()
	r.POST("/groups", func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
		handler.CreateGroup(c)
	})
	r.POST("/groups/:id/students", func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
		handler.AddStudents(c)
	})

	// 1. Create study group
	createReq := map[string]interface{}{
		"school_id":     "school-1",
		"name":          "Gruppo Recupero Matematica 2A",
		"subject_id":    "subj-1",
		"teacher_id":    "teacher-1",
		"academic_year": "2026/2027",
		"description":   "Gruppo per il recupero dei debiti formativi in matematica",
	}
	body1, _ := json.Marshal(createReq)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/groups", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// 2. Add students to group
	addReq := map[string]interface{}{
		"student_ids": []string{"st-1", "st-2"},
	}
	body2, _ := json.Marshal(addReq)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/groups/group-1/students", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}
