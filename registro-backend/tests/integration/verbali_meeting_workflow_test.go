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

	"registro-backend/internal/verbali"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockVerbaliRepo struct {
	meetings map[string]*verbali.CouncilMeeting
	verbali  map[string]*verbali.MeetingVerbale
	sigs     map[string][]verbali.VerbaleSignature
}

func (m *mockVerbaliRepo) CreateMeeting(ctx context.Context, cm *verbali.CouncilMeeting) error {
	if cm.ID == "" {
		cm.ID = "meeting-1"
	}
	cm.CreatedAt = time.Now()
	m.meetings[cm.ID] = cm
	return nil
}

func (m *mockVerbaliRepo) ListMeetings(ctx context.Context, schoolID, classID string) ([]*verbali.CouncilMeeting, error) {
	res := make([]*verbali.CouncilMeeting, 0)
	for _, cm := range m.meetings {
		res = append(res, cm)
	}
	return res, nil
}

func (m *mockVerbaliRepo) GetMeetingByID(ctx context.Context, id string) (*verbali.CouncilMeeting, error) {
	cm, ok := m.meetings[id]
	if !ok {
		return nil, errors.New("meeting not found")
	}
	return cm, nil
}

func (m *mockVerbaliRepo) CreateVerbale(ctx context.Context, v *verbali.MeetingVerbale) error {
	if v.ID == "" {
		v.ID = "verb-1"
	}
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()
	m.verbali[v.ID] = v
	return nil
}

func (m *mockVerbaliRepo) GetVerbaleByID(ctx context.Context, id string, userID string) (*verbali.MeetingVerbale, error) {
	v, ok := m.verbali[id]
	if !ok {
		return nil, errors.New("verbale not found")
	}
	return v, nil
}

func (m *mockVerbaliRepo) ListVerbali(ctx context.Context, meetingID string, userID string) ([]*verbali.MeetingVerbale, error) {
	res := make([]*verbali.MeetingVerbale, 0)
	for _, v := range m.verbali {
		if v.MeetingID == meetingID {
			res = append(res, v)
		}
	}
	return res, nil
}

func (m *mockVerbaliRepo) SignVerbale(ctx context.Context, verbaleID, userID, ipAddress string) error {
	m.sigs[verbaleID] = append(m.sigs[verbaleID], verbali.VerbaleSignature{
		ID:        "sig-1",
		VerbaleID: verbaleID,
		UserID:    userID,
		SignedAt:  time.Now(),
		IPAddress: ipAddress,
	})
	return nil
}

func (m *mockVerbaliRepo) GetSignatures(ctx context.Context, verbaleID string) ([]verbali.VerbaleSignature, error) {
	return m.sigs[verbaleID], nil
}

func (m *mockVerbaliRepo) ClassBelongsToSchool(ctx context.Context, classID, schoolID string) (bool, error) {
	return true, nil
}

func TestIntegration_Verbali_Meeting_Workflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockVerbaliRepo{
		meetings: make(map[string]*verbali.CouncilMeeting),
		verbali:  make(map[string]*verbali.MeetingVerbale),
		sigs:     make(map[string][]verbali.VerbaleSignature),
	}
	svc := verbali.NewService(repo)
	handler := verbali.NewHandler(svc)

	r := gin.New()
	r.POST("/verbali/meetings", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.CreateMeeting(c)
	})
	r.POST("/verbali", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.CreateVerbale(c)
	})
	r.POST("/verbali/:id/sign", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.SignVerbale(c)
	})

	// 1. Create council meeting
	mReq := map[string]interface{}{
		"class_id":   "class-1",
		"title":      "Consiglio di Classe 1° Trimestre - Classe 2A",
		"date":       "2026-10-15",
		"start_time": "15:00",
		"end_time":   "17:00",
		"agenda":     "Andamento didattico-disciplinare e approvazione PDP",
	}
	body1, _ := json.Marshal(mReq)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/verbali/meetings", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// 2. Draft verbale for meeting
	secID := "teacher-1"
	vReq := map[string]interface{}{
		"meeting_id":   "meeting-1",
		"title":        "Verbale Consiglio di Classe n. 1",
		"content":      "Alle ore 15:00 si insedia il Consiglio di Classe. Approvato all'unanimità il PDP.",
		"secretary_id": secID,
		"is_published": true,
	}
	body2, _ := json.Marshal(vReq)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/verbali", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusCreated, w2.Code)

	// 3. Council member signs verbale
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("POST", "/verbali/verb-1/sign", nil)
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
}
