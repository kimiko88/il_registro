package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/pcto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockPCTORepo struct {
	projects       map[string]*pcto.Project
	participations map[string]*pcto.Participation
	hours          map[string]*pcto.HourLog
}

func (m *mockPCTORepo) CreateProject(ctx context.Context, p *pcto.Project) error {
	if p.ID == "" {
		p.ID = "proj-1"
	}
	m.projects[p.ID] = p
	return nil
}

func (m *mockPCTORepo) GetProjects(ctx context.Context, schoolID string) ([]pcto.Project, error) {
	res := make([]pcto.Project, 0)
	for _, p := range m.projects {
		res = append(res, *p)
	}
	return res, nil
}

func (m *mockPCTORepo) GetProjectByID(ctx context.Context, id string) (*pcto.Project, error) {
	p, ok := m.projects[id]
	if !ok {
		return nil, errors.New("project not found")
	}
	return p, nil
}

func (m *mockPCTORepo) UpdateProject(ctx context.Context, p *pcto.Project) error {
	m.projects[p.ID] = p
	return nil
}

func (m *mockPCTORepo) DeleteProject(ctx context.Context, id string) error {
	delete(m.projects, id)
	return nil
}

func (m *mockPCTORepo) AssignStudent(ctx context.Context, p *pcto.Participation) error {
	if p.ID == "" {
		p.ID = "part-1"
	}
	m.participations[p.ID] = p
	return nil
}

func (m *mockPCTORepo) GetParticipationsByProject(ctx context.Context, projectID string) ([]pcto.Participation, error) {
	res := make([]pcto.Participation, 0)
	for _, p := range m.participations {
		if p.ProjectID == projectID {
			res = append(res, *p)
		}
	}
	return res, nil
}

func (m *mockPCTORepo) GetParticipationsByStudent(ctx context.Context, studentID string) ([]pcto.Participation, error) {
	res := make([]pcto.Participation, 0)
	for _, p := range m.participations {
		if p.StudentID == studentID {
			res = append(res, *p)
		}
	}
	return res, nil
}

func (m *mockPCTORepo) GetParticipation(ctx context.Context, projectID, studentID string) (*pcto.Participation, error) {
	return &pcto.Participation{ID: "part-1", ProjectID: projectID, StudentID: studentID}, nil
}

func (m *mockPCTORepo) LogHours(ctx context.Context, h *pcto.HourLog) error {
	if h.ID == "" {
		h.ID = "hour-1"
	}
	h.Verified = false
	m.hours[h.ID] = h
	return nil
}

func (m *mockPCTORepo) GetHours(ctx context.Context, participationID string) ([]pcto.HourLog, error) {
	res := make([]pcto.HourLog, 0)
	for _, h := range m.hours {
		res = append(res, *h)
	}
	return res, nil
}

func (m *mockPCTORepo) VerifyHours(ctx context.Context, hourID, teacherID string) error {
	if h, ok := m.hours[hourID]; ok {
		h.Verified = true
		h.VerifiedBy = &teacherID
	}
	return nil
}

func (m *mockPCTORepo) UpdateHourLogStatus(ctx context.Context, logID, status string) error {
	return nil
}

func (m *mockPCTORepo) CreateCompany(ctx context.Context, c *pcto.Company) error {
	return nil
}

func (m *mockPCTORepo) GetCompanies(ctx context.Context, schoolID string) ([]pcto.Company, error) {
	return []pcto.Company{}, nil
}

func (m *mockPCTORepo) GetStats(ctx context.Context, schoolID string) (*pcto.PCTOStats, error) {
	return &pcto.PCTOStats{}, nil
}

func TestIntegration_PCTO_Internship_Workflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockPCTORepo{
		projects:       make(map[string]*pcto.Project),
		participations: make(map[string]*pcto.Participation),
		hours:          make(map[string]*pcto.HourLog),
	}
	svc := pcto.NewService(repo)
	handler := pcto.NewHandler(svc)

	r := gin.New()
	r.POST("/pcto/projects", func(c *gin.Context) {
		c.Set("user_id", "superadmin-1")
		c.Set("role", "superadmin")
		c.Set("school_id", "school-1")
		handler.CreateProject(c)
	})
	r.POST("/pcto/hours", func(c *gin.Context) {
		c.Set("user_id", "student-1")
		c.Set("role", "student")
		c.Set("school_id", "school-1")
		handler.LogHours(c)
	})

	// 1. Create PCTO project
	projReq := map[string]interface{}{
		"title":        "Tirocinio Sviluppo Software presso TechCorp",
		"description":  "Percorso PCTO 40 ore sviluppo web",
		"type":         "External",
		"total_hours":  40,
		"company_name": "TechCorp SRL",
		"start_date":   "2026-09-01",
		"end_date":     "2026-09-30",
	}
	body1, _ := json.Marshal(projReq)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/pcto/projects", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// 2. Student logs 8 hours for PCTO activity
	hourReq := map[string]interface{}{
		"project_id": "proj-1",
		"date":       "2026-09-02",
		"hours":      8,
		"activity":   "Attività di formazione front-end Vue.js",
	}
	body2, _ := json.Marshal(hourReq)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/pcto/hours", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusCreated, w2.Code)
}
