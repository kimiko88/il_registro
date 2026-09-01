package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// WSH01 — Direct message delivery to targeted user
func TestWebSocket_Hub_DirectUserDelivery(t *testing.T) {
	hub := ws.NewHub("")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	clientChan := make(chan []byte, 10)
	client := &ws.Client{
		Hub:      hub,
		Send:     clientChan,
		UserID:   "teacher-u1",
		SchoolID: "school-A",
		Role:     "teacher",
	}

	// Register client
	hub.RegisterClient(client)
	time.Sleep(50 * time.Millisecond)

	// Broadcast to specific user
	hub.BroadcastToUser("teacher-u1", "school-A", "NEW_GRADE", map[string]string{
		"student": "Mario Rossi",
		"grade":   "8.5",
	})

	select {
	case msgBytes := <-clientChan:
		var msg ws.Message
		err := json.Unmarshal(msgBytes, &msg)
		require.NoError(t, err)
		assert.Equal(t, "NEW_GRADE", msg.Type)
		assert.Equal(t, "teacher-u1", msg.Recipient)
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for targeted WebSocket message")
	}

	// Unregister client
	hub.UnregisterClient(client)
}

// WSH02 — Role-targeted broadcast only reaches specified roles in school
func TestWebSocket_Hub_SchoolRolesDelivery(t *testing.T) {
	hub := ws.NewHub("")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	teacherChan := make(chan []byte, 10)
	teacherClient := &ws.Client{
		Hub:      hub,
		Send:     teacherChan,
		UserID:   "teacher-u2",
		SchoolID: "school-A",
		Role:     "teacher",
	}

	studentChan := make(chan []byte, 10)
	studentClient := &ws.Client{
		Hub:      hub,
		Send:     studentChan,
		UserID:   "student-u3",
		SchoolID: "school-A",
		Role:     "student",
	}

	hub.RegisterClient(teacherClient)
	hub.RegisterClient(studentClient)
	time.Sleep(50 * time.Millisecond)

	// Broadcast only to teachers
	hub.BroadcastToSchoolRoles("school-A", []string{"teacher"}, "FACULTY_MEETING", map[string]string{
		"time": "16:00",
	})

	// Teacher should receive message
	select {
	case msgBytes := <-teacherChan:
		var msg ws.Message
		err := json.Unmarshal(msgBytes, &msg)
		require.NoError(t, err)
		assert.Equal(t, "FACULTY_MEETING", msg.Type)
	case <-time.After(1 * time.Second):
		t.Fatal("teacher did not receive role-filtered message")
	}

	// Student should NOT receive message
	select {
	case <-studentChan:
		t.Fatal("student received message intended only for teachers")
	case <-time.After(100 * time.Millisecond):
		// Expected: no message delivered
	}

	hub.UnregisterClient(teacherClient)
	hub.UnregisterClient(studentClient)
}

// WSH03 — Multi-tenant cross-school message isolation
func TestWebSocket_Hub_CrossSchoolIsolation(t *testing.T) {
	hub := ws.NewHub("")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go hub.Run(ctx)

	schoolAChan := make(chan []byte, 10)
	clientA := &ws.Client{
		Hub:      hub,
		Send:     schoolAChan,
		UserID:   "user-a",
		SchoolID: "school-A",
		Role:     "teacher",
	}

	schoolBChan := make(chan []byte, 10)
	clientB := &ws.Client{
		Hub:      hub,
		Send:     schoolBChan,
		UserID:   "user-b",
		SchoolID: "school-B",
		Role:     "teacher",
	}

	hub.RegisterClient(clientA)
	hub.RegisterClient(clientB)
	time.Sleep(50 * time.Millisecond)

	// Broadcast to school A only
	hub.BroadcastToSchool("school-A", "EMERGENCY_DRILL", map[string]string{"type": "fire"})

	select {
	case msgBytes := <-schoolAChan:
		var msg ws.Message
		err := json.Unmarshal(msgBytes, &msg)
		require.NoError(t, err)
		assert.Equal(t, "EMERGENCY_DRILL", msg.Type)
	case <-time.After(1 * time.Second):
		t.Fatal("school A client did not receive school broadcast")
	}

	// School B client must NOT receive anything
	select {
	case <-schoolBChan:
		t.Fatal("cross-tenant leak: school B received broadcast meant for school A")
	case <-time.After(100 * time.Millisecond):
		// Expected
	}

	hub.UnregisterClient(clientA)
	hub.UnregisterClient(clientB)
}

// WSH04 — Unauthenticated WebSocket upgrade endpoint rejects connection with 401
func TestWebSocket_Handler_Unauthorized_NoUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	hub := ws.NewHub("")
	h := ws.NewHandler(hub)

	r.GET("/api/v1/ws", h.Listen)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ws", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
