package communications

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupCommsRouter(h *Handler, userID, schoolID, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		if schoolID != "" {
			c.Set("school_id", schoolID)
		}
		if role != "" {
			c.Set("role", role)
		}
		c.Next()
	})
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestHandler_ListAndRead(t *testing.T) {
	mockRepo := new(MockRepository)
	mockUsers := new(MockUserRepoForComms)
	svc := NewService(mockRepo, mockUsers)
	h := NewHandler(svc)
	schoolID := "school-1"

	// 1. Unauthorized
	rNoUser := setupCommsRouter(h, "", schoolID, "teacher")
	req, _ := http.NewRequest("GET", "/api/v1/communications", nil)
	w := httptest.NewRecorder()
	rNoUser.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 2. List success
	mockRepo.On("List", mock.Anything, "user-1", schoolID).Return([]*Message{
		{ID: "msg-1", SchoolID: &schoolID, Subject: "Avviso", Body: "Test body"},
	}, nil)
	rUser := setupCommsRouter(h, "user-1", schoolID, "teacher")
	req, _ = http.NewRequest("GET", "/api/v1/communications", nil)
	w = httptest.NewRecorder()
	rUser.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. ListBacheca success
	mockRepo.On("ListBacheca", mock.Anything, schoolID, "user-1").Return([]*Message{
		{ID: "msg-1", SchoolID: &schoolID, Subject: "Bacheca", Type: "notice"},
	}, nil)
	req, _ = http.NewRequest("GET", "/api/v1/communications/bacheca", nil)
	w = httptest.NewRecorder()
	rUser.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. ListCircolari success
	mockRepo.On("ListCircolari", mock.Anything, schoolID, "user-1", "2026").Return([]*Message{
		{ID: "msg-2", SchoolID: &schoolID, Subject: "Circolare N.1", Type: "circular"},
	}, nil)
	req, _ = http.NewRequest("GET", "/api/v1/communications/circolari?year=2026", nil)
	w = httptest.NewRecorder()
	rUser.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. GetUnreadCount success
	mockRepo.On("GetUnreadCount", mock.Anything, "user-1").Return(3, nil)
	req, _ = http.NewRequest("GET", "/api/v1/communications/unread-count", nil)
	w = httptest.NewRecorder()
	rUser.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. GetByID success
	mockRepo.On("Get", mock.Anything, "msg-1").Return(&Message{
		ID:          "msg-1",
		SchoolID:    &schoolID,
		SenderID:    "user-1",
		ReceiverIDs: []string{"user-1"},
		Subject:     "Dettaglio",
	}, nil)
	req, _ = http.NewRequest("GET", "/api/v1/communications/msg-1", nil)
	w = httptest.NewRecorder()
	rUser.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_SendUpdateDelete(t *testing.T) {
	mockRepo := new(MockRepository)
	mockUsers := new(MockUserRepoForComms)
	svc := NewService(mockRepo, mockUsers)
	h := NewHandler(svc)
	schoolID := "school-1"
	rUser := setupCommsRouter(h, "user-1", schoolID, "principal")

	// 1. Send invalid JSON -> 400
	req, _ := http.NewRequest("POST", "/api/v1/communications", bytes.NewBufferString("{invalid"))
	w := httptest.NewRecorder()
	rUser.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 2. Send success -> 201
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	sendBody, _ := json.Marshal(CreateMessageRequest{
		Recipients: []string{"user-2"},
		Subject:    "Nuova Circolare",
		Body:       "Testo della circolare",
		Type:       "circular",
	})
	req, _ = http.NewRequest("POST", "/api/v1/communications", bytes.NewBuffer(sendBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rUser.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 3. Update -> 200
	mockRepo.On("Get", mock.Anything, "msg-1").Return(&Message{
		ID:          "msg-1",
		SchoolID:    &schoolID,
		SenderID:    "user-1",
		Subject:     "Vecchio titolo",
		ReceiverIDs: []string{"user-2"},
	}, nil)
	mockRepo.On("GetSignatures", mock.Anything, "msg-1").Return([]string{}, nil)
	mockRepo.On("Update", mock.Anything, "msg-1", "Nuovo titolo", "Nuovo corpo").Return(nil)
	updateBody, _ := json.Marshal(map[string]string{
		"subject": "Nuovo titolo",
		"body":    "Nuovo corpo",
	})
	req, _ = http.NewRequest("PUT", "/api/v1/communications/msg-1", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rUser.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Delete -> 200
	mockRepo.On("Delete", mock.Anything, "msg-1").Return(nil)
	req, _ = http.NewRequest("DELETE", "/api/v1/communications/msg-1", nil)
	w = httptest.NewRecorder()
	rUser.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Interactions(t *testing.T) {
	mockRepo := new(MockRepository)
	mockUsers := new(MockUserRepoForComms)
	svc := NewService(mockRepo, mockUsers)
	h := NewHandler(svc)
	schoolID := "school-1"
	rParent := setupCommsRouter(h, "parent-1", schoolID, "parent")

	msgToSign := &Message{
		ID:                "msg-1",
		SchoolID:          &schoolID,
		SenderID:          "admin-1",
		ReceiverIDs:       []string{"parent-1"},
		RequiresSignature: true,
	}

	// 1. Sign
	mockRepo.On("Get", mock.Anything, "msg-1").Return(msgToSign, nil)
	mockRepo.On("SignWithIP", mock.Anything, "msg-1", "parent-1", mock.Anything).Return(nil)
	req, _ := http.NewRequest("POST", "/api/v1/communications/msg-1/sign", nil)
	w := httptest.NewRecorder()
	rParent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. MarkAsRead
	mockRepo.On("MarkAsRead", mock.Anything, "msg-1", "parent-1", mock.Anything).Return(nil)
	req, _ = http.NewRequest("POST", "/api/v1/communications/msg-1/read", nil)
	w = httptest.NewRecorder()
	rParent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. Ack
	mockRepo.On("Ack", mock.Anything, "msg-1", "parent-1").Return(nil)
	req, _ = http.NewRequest("POST", "/api/v1/communications/msg-1/ack", nil)
	w = httptest.NewRecorder()
	rParent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. GetSignatures (restricted to staff)
	rAdmin := setupCommsRouter(h, "admin-1", schoolID, "admin")
	mockRepo.On("GetSignatures", mock.Anything, "msg-1").Return([]string{"parent-1"}, nil)
	req, _ = http.NewRequest("GET", "/api/v1/communications/msg-1/signatures", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. GetSignatureReport (requires admin or sender)
	now := time.Now()
	mockRepo.On("GetSignatureReport", mock.Anything, "msg-1").Return(&SignatureReportResponse{
		CommunicationID: "msg-1",
		TotalRecipients: 10,
		SignedCount:     8,
		Signatures: []CommunicationSignature{
			{UserID: "parent-1", UserName: "Genitore 1", SignedAt: now},
		},
	}, nil)
	req, _ = http.NewRequest("GET", "/api/v1/communications/msg-1/signature-report", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. GetUnreadUsers
	mockRepo.On("GetUnreadUsers", mock.Anything, "msg-1").Return([]string{"parent-2"}, nil)
	req, _ = http.NewRequest("GET", "/api/v1/communications/msg-1/unread-users", nil)
	w = httptest.NewRecorder()
	rAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
