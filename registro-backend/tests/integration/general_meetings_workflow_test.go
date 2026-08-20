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

	"registro-backend/internal/general_meetings"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockGeneralMeetingsRepo struct {
	meetings map[string]*general_meetings.GeneralMeeting
	regs     map[string][]*general_meetings.GeneralMeetingRegistration
}

func (m *mockGeneralMeetingsRepo) Create(ctx context.Context, gm *general_meetings.GeneralMeeting) error {
	if gm.ID == "" {
		gm.ID = "gm-1"
	}
	gm.CreatedAt = time.Now()
	gm.UpdatedAt = time.Now()
	m.meetings[gm.ID] = gm
	return nil
}

func (m *mockGeneralMeetingsRepo) GetByID(ctx context.Context, id string, userID string) (*general_meetings.GeneralMeeting, error) {
	gm, ok := m.meetings[id]
	if !ok {
		return nil, errors.New("meeting not found")
	}
	return gm, nil
}

func (m *mockGeneralMeetingsRepo) List(ctx context.Context, schoolID, userID, role string) ([]*general_meetings.GeneralMeeting, error) {
	res := make([]*general_meetings.GeneralMeeting, 0)
	for _, gm := range m.meetings {
		res = append(res, gm)
	}
	return res, nil
}

func (m *mockGeneralMeetingsRepo) Delete(ctx context.Context, id string) error {
	delete(m.meetings, id)
	return nil
}

func (m *mockGeneralMeetingsRepo) RegisterUser(ctx context.Context, meetingID, userID string) error {
	m.regs[meetingID] = append(m.regs[meetingID], &general_meetings.GeneralMeetingRegistration{
		ID:           "reg-1",
		MeetingID:    meetingID,
		UserID:       userID,
		RegisteredAt: time.Now(),
	})
	return nil
}

func (m *mockGeneralMeetingsRepo) UnregisterUser(ctx context.Context, meetingID, userID string) error {
	return nil
}

func (m *mockGeneralMeetingsRepo) ListRegistrations(ctx context.Context, meetingID string) ([]*general_meetings.GeneralMeetingRegistration, error) {
	return m.regs[meetingID], nil
}

func TestIntegration_General_Meetings_Workflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockGeneralMeetingsRepo{
		meetings: make(map[string]*general_meetings.GeneralMeeting),
		regs:     make(map[string][]*general_meetings.GeneralMeetingRegistration),
	}
	svc := general_meetings.NewService(repo)
	handler := general_meetings.NewHandler(svc)

	r := gin.New()
	r.POST("/general-meetings", func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
		handler.Create(c)
	})
	r.POST("/general-meetings/:id/register", func(c *gin.Context) {
		c.Set("user_id", "parent-1")
		c.Set("role", "parent")
		c.Set("school_id", "school-1")
		handler.Register(c)
	})

	// 1. Administration schedules General Parent Assembly
	gmReq := general_meetings.CreateGeneralMeetingRequest{
		Title:       "Assemblea Generale Genitori d'Istituto",
		Description: "Presentazione del Piano Triennale dell'Offerta Formativa (PTOF)",
		Location:    "Auditorium Scolastico",
		MeetingDate: "2026-11-12T17:30:00Z",
		IsMandatory: false,
		TargetRoles: []string{"parent", "teacher"},
	}
	body1, _ := json.Marshal(gmReq)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/general-meetings", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// 2. Parent registers for the general assembly
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/general-meetings/gm-1/register", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}
