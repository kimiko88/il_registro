package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/general_meetings"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockGMQueueRepo struct {
	meetings map[string]*general_meetings.GeneralMeeting
	regs     map[string][]*general_meetings.GeneralMeetingRegistration
}

func newMockGMQueueRepo() *mockGMQueueRepo {
	return &mockGMQueueRepo{
		meetings: make(map[string]*general_meetings.GeneralMeeting),
		regs:     make(map[string][]*general_meetings.GeneralMeetingRegistration),
	}
}

func (m *mockGMQueueRepo) Create(ctx context.Context, gm *general_meetings.GeneralMeeting) error {
	if gm.ID == "" {
		gm.ID = "gm-" + time.Now().Format("150405.000000")
	}
	gm.CreatedAt = time.Now()
	gm.UpdatedAt = time.Now()
	m.meetings[gm.ID] = gm
	return nil
}

func (m *mockGMQueueRepo) GetByID(ctx context.Context, id string, userID string) (*general_meetings.GeneralMeeting, error) {
	gm, ok := m.meetings[id]
	if !ok {
		return nil, assert.AnError
	}
	return gm, nil
}

func (m *mockGMQueueRepo) List(ctx context.Context, schoolID, userID, role string) ([]*general_meetings.GeneralMeeting, error) {
	var res []*general_meetings.GeneralMeeting
	for _, gm := range m.meetings {
		res = append(res, gm)
	}
	return res, nil
}

func (m *mockGMQueueRepo) Delete(ctx context.Context, id string) error {
	delete(m.meetings, id)
	return nil
}

func (m *mockGMQueueRepo) RegisterUser(ctx context.Context, meetingID, userID string) error {
	reg := &general_meetings.GeneralMeetingRegistration{
		MeetingID:    meetingID,
		UserID:       userID,
		RegisteredAt: time.Now(),
	}
	m.regs[meetingID] = append(m.regs[meetingID], reg)
	return nil
}

func (m *mockGMQueueRepo) UnregisterUser(ctx context.Context, meetingID, userID string) error {
	list := m.regs[meetingID]
	var updated []*general_meetings.GeneralMeetingRegistration
	for _, r := range list {
		if r.UserID != userID {
			updated = append(updated, r)
		}
	}
	m.regs[meetingID] = updated
	return nil
}

func (m *mockGMQueueRepo) ListRegistrations(ctx context.Context, meetingID string) ([]*general_meetings.GeneralMeetingRegistration, error) {
	return m.regs[meetingID], nil
}

func setupGeneralMeetingsQueueRouter(repo general_meetings.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := general_meetings.NewService(repo)
	handler := general_meetings.NewHandler(svc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		role := c.GetHeader("X-Role")
		if role == "" {
			role = "admin"
		}
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = "admin-1"
		}
		c.Set("role", role)
		c.Set("user_id", userID)
		c.Set("school_id", "school-1")
		c.Next()
	})
	handler.RegisterRoutes(api)
	return r
}

func TestIntegration_GeneralMeeting_Queue_Workflow(t *testing.T) {
	repo := newMockGMQueueRepo()
	r := setupGeneralMeetingsQueueRouter(repo)

	// 1. Admin / Secretary creates a General Meeting (Assemblea Generale / Colloqui Generali)
	maxPart := 50
	createReq := general_meetings.CreateGeneralMeetingRequest{
		Title:           "Colloqui Generali Scuola-Famiglia 2° Quadrimestre",
		Description:     "Sessione plenaria e ricevimento docenti su appuntamento",
		Location:        "Plesso Centrale & Google Meet",
		MeetingDate:     "2026-05-15",
		MaxParticipants: &maxPart,
		IsMandatory:     false,
		TargetRoles:     []string{"parent", "teacher"},
	}
	body, _ := json.Marshal(createReq)
	reqCreate, _ := http.NewRequest("POST", "/api/v1/general-meetings", bytes.NewReader(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("X-Role", "admin")
	reqCreate.Header.Set("X-User-ID", "admin-1")
	wCreate := httptest.NewRecorder()
	r.ServeHTTP(wCreate, reqCreate)
	require.Equal(t, http.StatusCreated, wCreate.Code)

	var createdMeeting general_meetings.GeneralMeeting
	err := json.Unmarshal(wCreate.Body.Bytes(), &createdMeeting)
	require.NoError(t, err)
	assert.NotEmpty(t, createdMeeting.ID)

	// 2. Parent-1 registers to the general meeting slot
	reqReg, _ := http.NewRequest("POST", "/api/v1/general-meetings/"+createdMeeting.ID+"/register", nil)
	reqReg.Header.Set("X-Role", "parent")
	reqReg.Header.Set("X-User-ID", "parent-1")
	wReg := httptest.NewRecorder()
	r.ServeHTTP(wReg, reqReg)
	require.Equal(t, http.StatusOK, wReg.Code)

	// 3. Teacher/Admin lists active registrations / queue list
	reqQueue, _ := http.NewRequest("GET", "/api/v1/general-meetings/"+createdMeeting.ID+"/registrations", nil)
	reqQueue.Header.Set("X-Role", "admin")
	reqQueue.Header.Set("X-User-ID", "admin-1")
	wQueue := httptest.NewRecorder()
	r.ServeHTTP(wQueue, reqQueue)
	require.Equal(t, http.StatusOK, wQueue.Code)

	var queueList []*general_meetings.GeneralMeetingRegistration
	err = json.Unmarshal(wQueue.Body.Bytes(), &queueList)
	require.NoError(t, err)
	assert.Len(t, queueList, 1)
	assert.Equal(t, "parent-1", queueList[0].UserID)
}
