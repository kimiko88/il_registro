package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/orientamento"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockOrientamentoRepo struct {
	events map[string]*orientamento.Event
	parts  map[string]*orientamento.Participation
	prefs  map[string]*orientamento.StudentPreference
}

func (m *mockOrientamentoRepo) CreateEvent(ctx context.Context, e *orientamento.Event) error {
	if e.ID == "" {
		e.ID = "event-1"
	}
	m.events[e.ID] = e
	return nil
}

func (m *mockOrientamentoRepo) GetEvents(ctx context.Context, schoolID string) ([]orientamento.Event, error) {
	res := make([]orientamento.Event, 0)
	for _, e := range m.events {
		res = append(res, *e)
	}
	return res, nil
}

func (m *mockOrientamentoRepo) RegisterStudent(ctx context.Context, p *orientamento.Participation) error {
	if p.ID == "" {
		p.ID = "part-1"
	}
	m.parts[p.ID] = p
	return nil
}

func (m *mockOrientamentoRepo) GetParticipations(ctx context.Context, studentID string) ([]orientamento.Participation, error) {
	res := make([]orientamento.Participation, 0)
	for _, p := range m.parts {
		res = append(res, *p)
	}
	return res, nil
}

func (m *mockOrientamentoRepo) MarkAttendance(ctx context.Context, eventID, studentID string, attended bool) error {
	for _, p := range m.parts {
		if p.EventID == eventID && p.StudentID == studentID {
			p.Attended = attended
		}
	}
	return nil
}

func (m *mockOrientamentoRepo) SavePreference(ctx context.Context, p *orientamento.StudentPreference) error {
	m.prefs[p.StudentID] = p
	return nil
}

func (m *mockOrientamentoRepo) GetPreference(ctx context.Context, studentID string) (*orientamento.StudentPreference, error) {
	pref, ok := m.prefs[studentID]
	if !ok {
		return nil, errors.New("preference not found")
	}
	return pref, nil
}

func (m *mockOrientamentoRepo) SaveCapolavoro(ctx context.Context, c *orientamento.Capolavoro) error {
	return nil
}
func (m *mockOrientamentoRepo) GetCapolavori(ctx context.Context, studentID string) ([]orientamento.Capolavoro, error) {
	return nil, nil
}
func (m *mockOrientamentoRepo) GetCurriculumStudente(ctx context.Context, studentID string) (*orientamento.CurriculumStudenteSummary, error) {
	return nil, nil
}

func TestIntegration_Orientamento_Guidance_Workflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockOrientamentoRepo{
		events: make(map[string]*orientamento.Event),
		parts:  make(map[string]*orientamento.Participation),
		prefs:  make(map[string]*orientamento.StudentPreference),
	}
	svc := orientamento.NewService(repo)
	handler := orientamento.NewHandler(svc)

	r := gin.New()
	r.POST("/orientamento/events", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.CreateEvent(c)
	})
	r.POST("/orientamento/register", func(c *gin.Context) {
		c.Set("user_id", "student-1")
		c.Set("role", "student")
		c.Set("school_id", "school-1")
		handler.RegisterStudent(c)
	})

	// 1. Teacher registers orientation activity (30h)
	eventReq := orientamento.CreateEventRequest{
		Title:        "Fiera dell'Orientamento Universitario",
		Description:  "Presentazione corsi di laurea in Ingegneria e Informatica",
		Category:     "University",
		Date:         "2026-11-05T09:00:00Z",
		EndDate:      "2026-11-05T15:00:00Z",
		Location:     "Aula Magna",
		Hours:        6,
		MaxAttendees: 100,
	}
	body1, _ := json.Marshal(eventReq)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/orientamento/events", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// 2. Student registers for orientation event
	regReq := map[string]interface{}{
		"event_id": "event-1",
	}
	body2, _ := json.Marshal(regReq)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/orientamento/register", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}
