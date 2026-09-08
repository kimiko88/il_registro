package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/notifications"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Mock Repository ───────────────────────────────────────────────────────

type mockNotifRepo struct {
	tokens        []notifications.PushToken
	notifications []notifications.DBNotification
}

func (m *mockNotifRepo) SaveToken(_ context.Context, t *notifications.PushToken) error {
	// upsert by device_token
	for i, existing := range m.tokens {
		if existing.UserID == t.UserID && existing.DeviceToken == t.DeviceToken {
			m.tokens[i] = *t
			return nil
		}
	}
	m.tokens = append(m.tokens, *t)
	return nil
}

func (m *mockNotifRepo) GetUserTokens(_ context.Context, userID string) ([]notifications.PushToken, error) {
	var out []notifications.PushToken
	for _, t := range m.tokens {
		if t.UserID == userID {
			out = append(out, t)
		}
	}
	return out, nil
}

func (m *mockNotifRepo) DeleteToken(_ context.Context, userID, deviceToken string) error {
	var remaining []notifications.PushToken
	for _, t := range m.tokens {
		if t.UserID != userID || t.DeviceToken != deviceToken {
			remaining = append(remaining, t)
		}
	}
	m.tokens = remaining
	return nil
}

func (m *mockNotifRepo) CreateDBNotification(_ context.Context, n *notifications.DBNotification) error {
	n.ID = "notif-" + n.UserID
	n.CreatedAt = time.Now()
	m.notifications = append(m.notifications, *n)
	return nil
}

func (m *mockNotifRepo) ListDBNotifications(_ context.Context, userID string, unreadOnly bool, limit, offset int) ([]notifications.DBNotification, error) {
	var out []notifications.DBNotification
	for _, n := range m.notifications {
		if n.UserID != userID {
			continue
		}
		if unreadOnly && n.ReadAt != nil {
			continue
		}
		out = append(out, n)
	}
	if offset >= len(out) {
		return []notifications.DBNotification{}, nil
	}
	out = out[offset:]
	if limit < len(out) {
		out = out[:limit]
	}
	return out, nil
}

func (m *mockNotifRepo) MarkAsRead(_ context.Context, userID, notificationID string) error {
	now := time.Now()
	for i, n := range m.notifications {
		if n.UserID == userID && n.ID == notificationID {
			m.notifications[i].ReadAt = &now
			return nil
		}
	}
	return nil
}

func (m *mockNotifRepo) MarkAllAsRead(_ context.Context, userID string) error {
	now := time.Now()
	for i, n := range m.notifications {
		if n.UserID == userID && n.ReadAt == nil {
			m.notifications[i].ReadAt = &now
		}
	}
	return nil
}

// ─── Helpers ───────────────────────────────────────────────────────────────

func setupNotifRouter(repo *mockNotifRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := notifications.NewService(repo)
	h := notifications.NewHandler(svc)
	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		c.Set("user_id", "user-abc")
		c.Set("role", "teacher")
		c.Next()
	})
	h.RegisterRoutes(api)
	return r
}

// ─── Tests ─────────────────────────────────────────────────────────────────

// N01 — Register a push token for an authenticated user
func TestNotifications_RegisterPushToken(t *testing.T) {
	repo := &mockNotifRepo{}
	r := setupNotifRouter(repo)

	body, _ := json.Marshal(map[string]string{
		"device_token": "fcm-token-xyz",
		"platform":     "android",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notifications/push-tokens", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	require.Len(t, repo.tokens, 1)
	assert.Equal(t, "fcm-token-xyz", repo.tokens[0].DeviceToken)
	assert.Equal(t, "android", repo.tokens[0].Platform)
}

// N02 — Register via /register-device endpoint validates platform
func TestNotifications_RegisterDevice_InvalidPlatform(t *testing.T) {
	repo := &mockNotifRepo{}
	r := setupNotifRouter(repo)

	body, _ := json.Marshal(map[string]string{
		"fcm_token": "tok",
		"platform":  "fax-machine",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notifications/register-device", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// N03 — Unregister a push token removes it from repository
func TestNotifications_UnregisterToken(t *testing.T) {
	repo := &mockNotifRepo{
		tokens: []notifications.PushToken{
			{UserID: "user-abc", DeviceToken: "tok-to-remove", Platform: "ios"},
		},
	}
	r := setupNotifRouter(repo)

	body, _ := json.Marshal(map[string]string{"device_token": "tok-to-remove"})
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/notifications/push-tokens", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, repo.tokens)
}

// N04 — List notifications returns only current user's entries
func TestNotifications_ListNotifications(t *testing.T) {
	now := time.Now()
	repo := &mockNotifRepo{
		notifications: []notifications.DBNotification{
			{ID: "n1", UserID: "user-abc", Title: "Voto inserito", Body: "8/10 in Matematica", Type: "grade", CreatedAt: now},
			{ID: "n2", UserID: "user-abc", Title: "Assenza giustificata", Body: "Assenza del 28/08", Type: "absence", CreatedAt: now},
			{ID: "n3", UserID: "user-other", Title: "Other user", Body: "...", Type: "info", CreatedAt: now},
		},
	}
	r := setupNotifRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/notifications", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var list []notifications.DBNotification
	err := json.Unmarshal(w.Body.Bytes(), &list)
	require.NoError(t, err)
	assert.Len(t, list, 2, "should only return user-abc's notifications")
}

// N05 — unread_only filter works correctly
func TestNotifications_ListNotifications_UnreadOnly(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)
	repo := &mockNotifRepo{
		notifications: []notifications.DBNotification{
			{ID: "n1", UserID: "user-abc", Title: "Unread", Type: "info", CreatedAt: now},
			{ID: "n2", UserID: "user-abc", Title: "Already read", Type: "grade", CreatedAt: now, ReadAt: &past},
		},
	}
	r := setupNotifRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/notifications?unread_only=true", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var list []notifications.DBNotification
	_ = json.Unmarshal(w.Body.Bytes(), &list)
	assert.Len(t, list, 1)
	assert.Equal(t, "Unread", list[0].Title)
}

// N06 — Mark a single notification as read
func TestNotifications_MarkAsRead(t *testing.T) {
	now := time.Now()
	repo := &mockNotifRepo{
		notifications: []notifications.DBNotification{
			{ID: "notif-user-abc", UserID: "user-abc", Title: "Test", Type: "info", CreatedAt: now},
		},
	}
	r := setupNotifRouter(repo)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/notifications/notif-user-abc/read", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, repo.notifications[0].ReadAt, "notification should have ReadAt set")
}

// N07 — Mark all notifications as read in one call
func TestNotifications_MarkAllAsRead(t *testing.T) {
	now := time.Now()
	repo := &mockNotifRepo{
		notifications: []notifications.DBNotification{
			{ID: "n1", UserID: "user-abc", Title: "A", Type: "info", CreatedAt: now},
			{ID: "n2", UserID: "user-abc", Title: "B", Type: "grade", CreatedAt: now},
		},
	}
	r := setupNotifRouter(repo)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/notifications/read-all", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	for _, n := range repo.notifications {
		assert.NotNil(t, n.ReadAt, "all notifications for user-abc should be read: %s", n.ID)
	}
}

// N08 — SendPush is restricted to admin/superadmin role
func TestNotifications_SendPush_ForbiddenForTeacher(t *testing.T) {
	repo := &mockNotifRepo{}
	r := setupNotifRouter(repo) // role = teacher

	body, _ := json.Marshal(map[string]string{
		"user_id": "target-user",
		"title":   "Avviso",
		"body":    "Riunione domani",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/notifications/send-push", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// N09 — PWA manifest endpoint returns correct structure
func TestNotifications_GetManifest(t *testing.T) {
	repo := &mockNotifRepo{}
	r := setupNotifRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/manifest.json", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var manifest notifications.PWAConfig
	err := json.Unmarshal(w.Body.Bytes(), &manifest)
	require.NoError(t, err)
	assert.Equal(t, "Registro Elettronico Scolastico", manifest.Name)
	assert.Equal(t, "standalone", manifest.Display)
}

// N10 — Unauthenticated user cannot access notifications
func TestNotifications_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	repo := &mockNotifRepo{}
	svc := notifications.NewService(repo)
	h := notifications.NewHandler(svc)
	api := r.Group("/api/v1")
	// No user_id middleware
	h.RegisterRoutes(api)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/notifications", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
