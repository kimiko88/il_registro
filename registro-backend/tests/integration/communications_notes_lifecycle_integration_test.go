package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"registro-backend/internal/communications"
	"registro-backend/internal/notes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// In-memory mock for communications
type mockCommsLifecycleRepo struct {
	mu         sync.Mutex
	messages   map[string]*communications.Message
	readBy     map[string]map[string]time.Time // msgID -> userID -> readAt
	signatures map[string]map[string]communications.CommunicationSignature
}

func newMockCommsLifecycleRepo() *mockCommsLifecycleRepo {
	return &mockCommsLifecycleRepo{
		messages:   make(map[string]*communications.Message),
		readBy:     make(map[string]map[string]time.Time),
		signatures: make(map[string]map[string]communications.CommunicationSignature),
	}
}

func (m *mockCommsLifecycleRepo) Create(ctx context.Context, msg *communications.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	msg.ID = fmt.Sprintf("comm-%d", len(m.messages)+1)
	msg.CreatedAt = time.Now()
	m.messages[msg.ID] = msg
	m.readBy[msg.ID] = make(map[string]time.Time)
	m.signatures[msg.ID] = make(map[string]communications.CommunicationSignature)
	return nil
}

func (m *mockCommsLifecycleRepo) List(ctx context.Context, userID, schoolID string) ([]*communications.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []*communications.Message
	for _, msg := range m.messages {
		list = append(list, msg)
	}
	return list, nil
}

func (m *mockCommsLifecycleRepo) ListBacheca(ctx context.Context, schoolID, userID string) ([]*communications.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []*communications.Message
	for _, msg := range m.messages {
		if msg.Type == "bacheca" {
			list = append(list, msg)
		}
	}
	return list, nil
}

func (m *mockCommsLifecycleRepo) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.messages, id)
	return nil
}

func (m *mockCommsLifecycleRepo) Sign(ctx context.Context, communicationID string, userID string) error {
	return m.SignWithIP(ctx, communicationID, userID, "127.0.0.1")
}

func (m *mockCommsLifecycleRepo) SignWithIP(ctx context.Context, communicationID string, userID string, ipAddress string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.messages[communicationID]; !ok {
		return errors.New("message not found")
	}
	if _, ok := m.signatures[communicationID]; !ok {
		m.signatures[communicationID] = make(map[string]communications.CommunicationSignature)
	}
	m.signatures[communicationID][userID] = communications.CommunicationSignature{
		ID:              fmt.Sprintf("sig-%d", len(m.signatures[communicationID])+1),
		CommunicationID: communicationID,
		UserID:          userID,
		SignedAt:        time.Now(),
		IPAddress:       ipAddress,
	}
	return nil
}

func (m *mockCommsLifecycleRepo) GetSignatures(ctx context.Context, communicationID string) ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var userIDs []string
	if sigs, ok := m.signatures[communicationID]; ok {
		for uid := range sigs {
			userIDs = append(userIDs, uid)
		}
	}
	return userIDs, nil
}

func (m *mockCommsLifecycleRepo) GetSignatureReport(ctx context.Context, communicationID string) (*communications.SignatureReportResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	msg, ok := m.messages[communicationID]
	if !ok {
		return nil, errors.New("message not found")
	}
	var sigList []communications.CommunicationSignature
	if sigs, ok := m.signatures[communicationID]; ok {
		for _, sig := range sigs {
			sigList = append(sigList, sig)
		}
	}
	return &communications.SignatureReportResponse{
		CommunicationID: communicationID,
		Subject:         msg.Subject,
		TotalRecipients: len(msg.ReceiverIDs),
		SignedCount:     len(sigList),
		PendingCount:    len(msg.ReceiverIDs) - len(sigList),
		Signatures:      sigList,
	}, nil
}

func (m *mockCommsLifecycleRepo) Get(ctx context.Context, id string) (*communications.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	msg, ok := m.messages[id]
	if !ok {
		return nil, errors.New("message not found")
	}
	return msg, nil
}

func (m *mockCommsLifecycleRepo) Update(ctx context.Context, id string, subject, body string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	msg, ok := m.messages[id]
	if !ok {
		return errors.New("message not found")
	}
	msg.Subject = subject
	msg.Body = body
	return nil
}

func (m *mockCommsLifecycleRepo) MarkAsRead(ctx context.Context, communicationID, userID, ipAddress string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.readBy[communicationID]; !ok {
		m.readBy[communicationID] = make(map[string]time.Time)
	}
	m.readBy[communicationID][userID] = time.Now()
	return nil
}

func (m *mockCommsLifecycleRepo) GetUnreadUsers(ctx context.Context, communicationID string) ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	msg, ok := m.messages[communicationID]
	if !ok {
		return nil, errors.New("message not found")
	}
	var unread []string
	readers := m.readBy[communicationID]
	for _, uid := range msg.ReceiverIDs {
		if _, read := readers[uid]; !read {
			unread = append(unread, uid)
		}
	}
	return unread, nil
}

func (m *mockCommsLifecycleRepo) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	count := 0
	for msgID, msg := range m.messages {
		isRecipient := false
		for _, r := range msg.ReceiverIDs {
			if r == userID {
				isRecipient = true
				break
			}
		}
		if isRecipient {
			if _, read := m.readBy[msgID][userID]; !read {
				count++
			}
		}
	}
	return count, nil
}

func (m *mockCommsLifecycleRepo) ListCircolari(ctx context.Context, schoolID, userID, year string) ([]*communications.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []*communications.Message
	for _, msg := range m.messages {
		if msg.Type == "circular" {
			list = append(list, msg)
		}
	}
	return list, nil
}

func (m *mockCommsLifecycleRepo) Ack(ctx context.Context, communicationID, userID string) error {
	return m.MarkAsRead(ctx, communicationID, userID, "127.0.0.1")
}

var _ communications.Repository = (*mockCommsLifecycleRepo)(nil)

// In-memory mock for notes
type mockNotesLifecycleRepo struct {
	mu           sync.Mutex
	notes        map[string]*notes.StudentNote
	classTeacher map[string]string // classID -> teacherID
}

func newMockNotesLifecycleRepo() *mockNotesLifecycleRepo {
	return &mockNotesLifecycleRepo{
		notes:        make(map[string]*notes.StudentNote),
		classTeacher: make(map[string]string),
	}
}

func (m *mockNotesLifecycleRepo) Create(ctx context.Context, note *notes.StudentNote) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	note.ID = fmt.Sprintf("note-%d", len(m.notes)+1)
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()
	m.notes[note.ID] = note
	return nil
}

func (m *mockNotesLifecycleRepo) Update(ctx context.Context, note *notes.StudentNote) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	note.UpdatedAt = time.Now()
	m.notes[note.ID] = note
	return nil
}

func (m *mockNotesLifecycleRepo) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.notes, id)
	return nil
}

func (m *mockNotesLifecycleRepo) DeleteWithReason(ctx context.Context, id string, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.notes, id)
	return nil
}

func (m *mockNotesLifecycleRepo) Get(ctx context.Context, id string) (*notes.StudentNote, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	n, ok := m.notes[id]
	if !ok {
		return nil, errors.New("note not found")
	}
	return n, nil
}

func (m *mockNotesLifecycleRepo) List(ctx context.Context, filter notes.NoteFilter) ([]notes.StudentNote, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []notes.StudentNote
	for _, n := range m.notes {
		if filter.StudentID != "" && n.StudentID != filter.StudentID {
			continue
		}
		if filter.ClassID != "" && n.ClassID != filter.ClassID {
			continue
		}
		list = append(list, *n)
	}
	return list, nil
}

func (m *mockNotesLifecycleRepo) ApproveNote(ctx context.Context, id, approverID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	n, ok := m.notes[id]
	if !ok {
		return errors.New("note not found")
	}
	n.IsApproved = true
	return nil
}

func (m *mockNotesLifecycleRepo) MarkAsViewedByParent(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	n, ok := m.notes[id]
	if !ok {
		return errors.New("note not found")
	}
	n.IsViewedByParent = true
	now := time.Now()
	n.ParentViewedAt = &now
	return nil
}

func (m *mockNotesLifecycleRepo) MarkManyAsViewedByParent(ctx context.Context, ids []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	for _, id := range ids {
		if n, ok := m.notes[id]; ok {
			n.IsViewedByParent = true
			n.ParentViewedAt = &now
		}
	}
	return nil
}

func (m *mockNotesLifecycleRepo) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return true, nil
}

var _ notes.Repository = (*mockNotesLifecycleRepo)(nil)

// ==========================================
// INTEGRATION TESTS
// ==========================================

func TestIntegration_CommunicationsAndNotesLifecycle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 1. Circular communication send, unread counter, read, and sign lifecycle
	t.Run("Circular_Send_Unread_Read_And_Sign_Lifecycle", func(t *testing.T) {
		commsRepo := newMockCommsLifecycleRepo()
		usersRepo := &mockUserRepoForComms{}
		commsSvc := communications.NewService(commsRepo, usersRepo)
		commsHandler := communications.NewHandler(commsSvc)

		// Admin/Principal context
		rAdmin := gin.New()
		rAdmin.Use(func(c *gin.Context) {
			c.Set("user_id", "admin-1")
			c.Set("role", "admin")
			c.Set("school_id", "school-1")
			c.Next()
		})
		commsHandler.RegisterRoutes(rAdmin.Group("/api/v1"))

		circReq := communications.CreateMessageRequest{
			Subject:           "Circolare 105: Viaggio di istruzione a Trieste",
			Body:              "Si comunica che sono aperte le adesioni al viaggio di istruzione...",
			Type:              "circular",
			Recipients:        []string{"parent-1", "teacher-1"},
			RequiresSignature: true,
		}
		cb, _ := json.Marshal(circReq)
		req1, _ := http.NewRequest(http.MethodPost, "/api/v1/communications", bytes.NewBuffer(cb))
		req1.Header.Set("Content-Type", "application/json")
		w1 := httptest.NewRecorder()
		rAdmin.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusCreated, w1.Code)

		var sentMsg communications.Message
		_ = json.Unmarshal(w1.Body.Bytes(), &sentMsg)
		assert.NotEmpty(t, sentMsg.ID)
		assert.True(t, sentMsg.RequiresSignature)

		// Parent context
		rParent := gin.New()
		rParent.Use(func(c *gin.Context) {
			c.Set("user_id", "parent-1")
			c.Set("role", "parent")
			c.Set("school_id", "school-1")
			c.Next()
		})
		commsHandler.RegisterRoutes(rParent.Group("/api/v1"))

		// Check unread count for parent
		req2, _ := http.NewRequest(http.MethodGet, "/api/v1/communications/unread-count", nil)
		w2 := httptest.NewRecorder()
		rParent.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)
		var unreadResp map[string]interface{}
		_ = json.Unmarshal(w2.Body.Bytes(), &unreadResp)
		assert.Equal(t, float64(1), unreadResp["count"])

		// Check unread users from management endpoint
		reqUnreadUsers, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/communications/%s/unread-users", sentMsg.ID), nil)
		wUnread := httptest.NewRecorder()
		rAdmin.ServeHTTP(wUnread, reqUnreadUsers)
		assert.Equal(t, http.StatusOK, wUnread.Code)
		var unreadUsers []string
		_ = json.Unmarshal(wUnread.Body.Bytes(), &unreadUsers)
		assert.Equal(t, 2, len(unreadUsers))

		// Parent marks circular as read
		req3, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/communications/%s/read", sentMsg.ID), nil)
		w3 := httptest.NewRecorder()
		rParent.ServeHTTP(w3, req3)
		assert.Equal(t, http.StatusOK, w3.Code)

		// Unread count now drops to 0 for parent
		w4 := httptest.NewRecorder()
		rParent.ServeHTTP(w4, req2)
		assert.Equal(t, http.StatusOK, w4.Code)
		var unreadRespAfter map[string]interface{}
		_ = json.Unmarshal(w4.Body.Bytes(), &unreadRespAfter)
		assert.Equal(t, float64(0), unreadRespAfter["count"])

		// Parent signs circular
		req5, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/communications/%s/sign", sentMsg.ID), nil)
		w5 := httptest.NewRecorder()
		rParent.ServeHTTP(w5, req5)
		assert.Equal(t, http.StatusOK, w5.Code)

		// Verify signature is stored in signature report
		req6, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/communications/%s/signature-report", sentMsg.ID), nil)
		w6 := httptest.NewRecorder()
		rAdmin.ServeHTTP(w6, req6)
		assert.Equal(t, http.StatusOK, w6.Code)
		var sigReport communications.SignatureReportResponse
		_ = json.Unmarshal(w6.Body.Bytes(), &sigReport)
		assert.Equal(t, 1, sigReport.SignedCount)
		assert.Equal(t, 1, sigReport.PendingCount)
	})

	// 2. Student Disciplinary Note: Teacher creates -> Pending -> Principal Approves -> Parent views
	t.Run("Disciplinary_Note_Creation_Approval_And_ParentView_Lifecycle", func(t *testing.T) {
		notesRepo := newMockNotesLifecycleRepo()
		usersRepo := &mockUserRepoForComms{}
		notesSvc := notes.NewService(notesRepo, usersRepo)
		notesHandler := notes.NewHandler(notesSvc)

		// 1. Teacher creates disciplinary note for student
		rTeacher := gin.New()
		rTeacher.Use(func(c *gin.Context) {
			c.Set("user_id", "teacher-1")
			c.Set("role", "teacher")
			c.Set("school_id", "school-1")
			c.Next()
		})
		notesHandler.RegisterRoutes(rTeacher.Group("/api/v1"))

		createNoteReq := notes.CreateNoteRequest{
			StudentID:  "student-1",
			ClassID:    "class-1",
			Type:       "disciplinary",
			Note:       "Comportamento irrispettoso e ripetute interruzioni della lezione",
			Date:       "2026-03-01",
			IsReserved: false,
		}
		nb, _ := json.Marshal(createNoteReq)
		req1, _ := http.NewRequest(http.MethodPost, "/api/v1/notes", bytes.NewBuffer(nb))
		req1.Header.Set("Content-Type", "application/json")
		w1 := httptest.NewRecorder()
		rTeacher.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusCreated, w1.Code)

		var createdNote notes.StudentNote
		_ = json.Unmarshal(w1.Body.Bytes(), &createdNote)
		assert.False(t, createdNote.IsApproved, "disciplinary note must start unapproved")

		// 2. Parent checks notes before approval -> empty list
		rParent := gin.New()
		rParent.Use(func(c *gin.Context) {
			c.Set("user_id", "parent-1")
			c.Set("role", "parent")
			c.Set("school_id", "school-1")
			c.Next()
		})
		notesHandler.RegisterRoutes(rParent.Group("/api/v1"))

		reqParent1, _ := http.NewRequest(http.MethodGet, "/api/v1/notes/student/student-1", nil)
		wParent1 := httptest.NewRecorder()
		rParent.ServeHTTP(wParent1, reqParent1)
		assert.Equal(t, http.StatusOK, wParent1.Code)
		var notesListBefore []notes.StudentNote
		_ = json.Unmarshal(wParent1.Body.Bytes(), &notesListBefore)
		assert.Empty(t, notesListBefore, "unapproved disciplinary note must not be visible to parents")

		// 3. Principal approves disciplinary note
		rPrincipal := gin.New()
		rPrincipal.Use(func(c *gin.Context) {
			c.Set("user_id", "principal-1")
			c.Set("role", "principal")
			c.Set("school_id", "school-1")
			c.Next()
		})
		notesHandler.RegisterRoutes(rPrincipal.Group("/api/v1"))

		reqApprove, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/notes/%s/approve", createdNote.ID), nil)
		wApprove := httptest.NewRecorder()
		rPrincipal.ServeHTTP(wApprove, reqApprove)
		assert.Equal(t, http.StatusOK, wApprove.Code)

		// 4. Parent views notes again -> note is returned and marked as viewed
		wParent2 := httptest.NewRecorder()
		rParent.ServeHTTP(wParent2, reqParent1)
		assert.Equal(t, http.StatusOK, wParent2.Code)
		var notesListAfter []notes.StudentNote
		_ = json.Unmarshal(wParent2.Body.Bytes(), &notesListAfter)
		assert.Len(t, notesListAfter, 1)
		assert.True(t, notesListAfter[0].IsApproved)
		assert.True(t, notesListAfter[0].IsViewedByParent)
		assert.NotNil(t, notesListAfter[0].ParentViewedAt)
	})
}
