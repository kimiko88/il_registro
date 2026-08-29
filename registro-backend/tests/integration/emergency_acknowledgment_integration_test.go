package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/communications"
	"registro-backend/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCommsRepoForAck struct {
	messages   map[string]*communications.Message
	signatures map[string][]string
	readUsers  map[string][]string
}

func newMockCommsRepoForAck() *mockCommsRepoForAck {
	return &mockCommsRepoForAck{
		messages:   make(map[string]*communications.Message),
		signatures: make(map[string][]string),
		readUsers:  make(map[string][]string),
	}
}

func (m *mockCommsRepoForAck) Create(ctx context.Context, msg *communications.Message) error {
	if msg.ID == "" {
		msg.ID = "comm-" + time.Now().Format("150405.000000")
	}
	msg.CreatedAt = time.Now()
	m.messages[msg.ID] = msg
	return nil
}

func (m *mockCommsRepoForAck) Get(ctx context.Context, id string) (*communications.Message, error) {
	msg, ok := m.messages[id]
	if !ok {
		return nil, assert.AnError
	}
	return msg, nil
}

func (m *mockCommsRepoForAck) Update(ctx context.Context, id string, subject, body string) error {
	if msg, ok := m.messages[id]; ok {
		msg.Subject = subject
		msg.Body = body
	}
	return nil
}

func (m *mockCommsRepoForAck) Delete(ctx context.Context, id string) error {
	delete(m.messages, id)
	return nil
}

func (m *mockCommsRepoForAck) Sign(ctx context.Context, communicationID, userID string) error {
	m.signatures[communicationID] = append(m.signatures[communicationID], userID)
	return nil
}

func (m *mockCommsRepoForAck) SignWithIP(ctx context.Context, communicationID, userID, ipAddress string) error {
	m.signatures[communicationID] = append(m.signatures[communicationID], userID)
	return nil
}

func (m *mockCommsRepoForAck) Ack(ctx context.Context, communicationID, userID string) error {
	m.signatures[communicationID] = append(m.signatures[communicationID], userID)
	return nil
}

func (m *mockCommsRepoForAck) GetSignatures(ctx context.Context, communicationID string) ([]string, error) {
	return m.signatures[communicationID], nil
}

func (m *mockCommsRepoForAck) MarkAsRead(ctx context.Context, communicationID, userID, ipAddress string) error {
	m.readUsers[communicationID] = append(m.readUsers[communicationID], userID)
	return nil
}

func (m *mockCommsRepoForAck) GetUnreadUsers(ctx context.Context, communicationID string) ([]string, error) {
	return []string{}, nil
}

func (m *mockCommsRepoForAck) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	return 0, nil
}

func (m *mockCommsRepoForAck) GetSignatureReport(ctx context.Context, communicationID string) (*communications.SignatureReportResponse, error) {
	msg := m.messages[communicationID]
	sigs := m.signatures[communicationID]
	return &communications.SignatureReportResponse{
		CommunicationID: communicationID,
		TotalRecipients: len(msg.ReceiverIDs),
		SignedCount:     len(sigs),
		PendingCount:    len(msg.ReceiverIDs) - len(sigs),
	}, nil
}

func (m *mockCommsRepoForAck) List(ctx context.Context, userID, schoolID string) ([]*communications.Message, error) {
	var res []*communications.Message
	for _, msg := range m.messages {
		res = append(res, msg)
	}
	return res, nil
}

func (m *mockCommsRepoForAck) ListBacheca(ctx context.Context, schoolID, userID string) ([]*communications.Message, error) {
	return m.List(ctx, userID, schoolID)
}

func (m *mockCommsRepoForAck) ListCircolari(ctx context.Context, schoolID, userID, year string) ([]*communications.Message, error) {
	return m.List(ctx, userID, schoolID)
}

type mockUserRepoForAck struct {
	users.Repository
}

func (m *mockUserRepoForAck) GetByID(ctx context.Context, id string) (*users.User, error) {
	schoolID := "school-1"
	return &users.User{
		ID:        id,
		SchoolID:  &schoolID,
		FirstName: "Genitore",
		LastName:  "Bianchi",
		Role:      "parent",
	}, nil
}

func (m *mockUserRepoForAck) ListByIDs(ctx context.Context, ids []string) ([]users.User, error) {
	schoolID := "school-1"
	var res []users.User
	for _, id := range ids {
		res = append(res, users.User{
			ID:        id,
			SchoolID:  &schoolID,
			FirstName: "Genitore",
			LastName:  "Bianchi",
			Role:      "parent",
		})
	}
	return res, nil
}

func (m *mockUserRepoForAck) FilterValidRecipients(ctx context.Context, schoolID string, recipientIDs []string) ([]string, error) {
	return recipientIDs, nil
}

func setupCommsAckRouter(repo communications.Repository, uRepo users.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := communications.NewService(repo, uRepo)
	handler := communications.NewHandler(svc)

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

func TestIntegration_Emergency_Acknowledgment_Workflow(t *testing.T) {
	repo := newMockCommsRepoForAck()
	uRepo := &mockUserRepoForAck{}
	r := setupCommsAckRouter(repo, uRepo)

	// 1. Admin/Presidenza publishes an Urgent Circular with mandatory signature/ack
	sendReq := communications.CreateMessageRequest{
		Type:              "circolare",
		Subject:           "CIRCOLARE URGENTE: Chiusura Straordinaria per Allerta Meteo Rossa",
		Body:              "Si comunica che per disposizione prefettizia le lezioni sono sospese per la giornata di domani. È richiesta presa d'atto immediata da parte di tutte le famiglie.",
		Recipients:        []string{"parent-1", "parent-2"},
		RequiresSignature: true,
	}
	body, _ := json.Marshal(sendReq)
	reqSend, _ := http.NewRequest("POST", "/api/v1/communications", bytes.NewReader(body))
	reqSend.Header.Set("Content-Type", "application/json")
	reqSend.Header.Set("X-Role", "admin")
	reqSend.Header.Set("X-User-ID", "admin-1")
	wSend := httptest.NewRecorder()
	r.ServeHTTP(wSend, reqSend)
	require.Equal(t, http.StatusCreated, wSend.Code)

	var createdMsg communications.Message
	err := json.Unmarshal(wSend.Body.Bytes(), &createdMsg)
	require.NoError(t, err)
	assert.NotEmpty(t, createdMsg.ID)
	assert.True(t, createdMsg.RequiresSignature)

	// 2. Parent-1 acknowledges / signs the urgent circular
	reqAck, _ := http.NewRequest("POST", "/api/v1/communications/"+createdMsg.ID+"/ack", nil)
	reqAck.Header.Set("X-Role", "parent")
	reqAck.Header.Set("X-User-ID", "parent-1")
	wAck := httptest.NewRecorder()
	r.ServeHTTP(wAck, reqAck)
	require.Equal(t, http.StatusOK, wAck.Code)

	// 3. Admin queries signature report to check compliance rate
	reqReport, _ := http.NewRequest("GET", "/api/v1/communications/"+createdMsg.ID+"/signature-report", nil)
	reqReport.Header.Set("X-Role", "admin")
	reqReport.Header.Set("X-User-ID", "admin-1")
	wReport := httptest.NewRecorder()
	r.ServeHTTP(wReport, reqReport)
	require.Equal(t, http.StatusOK, wReport.Code)

	var report communications.SignatureReportResponse
	err = json.Unmarshal(wReport.Body.Bytes(), &report)
	require.NoError(t, err)
	assert.Equal(t, 2, report.TotalRecipients)
	assert.Equal(t, 1, report.SignedCount)
	assert.Equal(t, 1, report.PendingCount)
}
