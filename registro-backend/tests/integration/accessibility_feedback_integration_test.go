package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/accessibility"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockAccessibilityRepo struct {
	feedbacks   map[string]*accessibility.AccessibilityFeedback
	preferences map[string]string
}

func newMockAccessibilityRepo() *mockAccessibilityRepo {
	return &mockAccessibilityRepo{
		feedbacks:   make(map[string]*accessibility.AccessibilityFeedback),
		preferences: make(map[string]string),
	}
}

func (m *mockAccessibilityRepo) Create(ctx context.Context, f *accessibility.AccessibilityFeedback) error {
	if f.ID == "" {
		f.ID = "fb-" + time.Now().Format("150405.000000")
	}
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
	m.feedbacks[f.ID] = f
	return nil
}

func (m *mockAccessibilityRepo) List(ctx context.Context, schoolID, status string, limit, offset int) ([]accessibility.AccessibilityFeedback, int, error) {
	var res []accessibility.AccessibilityFeedback
	for _, f := range m.feedbacks {
		if schoolID != "" && f.SchoolID != nil && *f.SchoolID != schoolID {
			continue
		}
		if status != "" && f.Status != status {
			continue
		}
		res = append(res, *f)
	}
	return res, len(res), nil
}

func (m *mockAccessibilityRepo) GetByID(ctx context.Context, id string) (*accessibility.AccessibilityFeedback, error) {
	f, ok := m.feedbacks[id]
	if !ok {
		return nil, assert.AnError
	}
	return f, nil
}

func (m *mockAccessibilityRepo) UpdateStatus(ctx context.Context, id, status, responseNotes string) error {
	if f, ok := m.feedbacks[id]; ok {
		f.Status = status
		f.ResponseNotes = &responseNotes
		f.UpdatedAt = time.Now()
	}
	return nil
}

func (m *mockAccessibilityRepo) GetUserPreferences(ctx context.Context, userID string) (string, error) {
	pref, ok := m.preferences[userID]
	if !ok {
		return `{"high_contrast":false,"dyslexia_font":false,"font_size":100}`, nil
	}
	return pref, nil
}

func (m *mockAccessibilityRepo) UpsertUserPreferences(ctx context.Context, userID string, settingsJSON string) error {
	m.preferences[userID] = settingsJSON
	return nil
}

func setupAccessibilityRouter(repo accessibility.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := accessibility.NewService(repo)
	handler := accessibility.NewHandler(svc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		role := c.GetHeader("X-Role")
		if role == "" {
			role = "admin"
		}
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = "user-parent-1"
		}
		c.Set("role", role)
		c.Set("user_id", userID)
		c.Set("school_id", "school-1")
		c.Next()
	})

	api.POST("/accessibility/feedback", handler.SubmitPublic)
	api.GET("/accessibility/feedback", handler.List)
	api.PATCH("/accessibility/feedback/:id/status", handler.UpdateStatus)
	api.GET("/accessibility/preferences", handler.GetMyPreferences)
	api.POST("/accessibility/preferences", handler.SaveMyPreferences)

	return r
}

func TestIntegration_Accessibility_Feedback_And_Preferences_Workflow(t *testing.T) {
	repo := newMockAccessibilityRepo()
	r := setupAccessibilityRouter(repo)

	// 1. Citizen/Parent submits an AgID digital barrier feedback
	feedbackReq := accessibility.CreateFeedbackRequest{
		Name:        "Giuseppe Verdi",
		Email:       "g.verdi@example.com",
		BarrierType: "screen_reader",
		Description: "Il contrasto dei pulsanti del registro presenze non supera il rapporto WCAG 2.1 AA",
	}
	body, _ := json.Marshal(feedbackReq)
	reqSubmit, _ := http.NewRequest("POST", "/api/v1/accessibility/feedback", bytes.NewReader(body))
	reqSubmit.Header.Set("Content-Type", "application/json")
	wSubmit := httptest.NewRecorder()
	r.ServeHTTP(wSubmit, reqSubmit)
	require.Equal(t, http.StatusCreated, wSubmit.Code)

	var submitResp accessibility.FeedbackResponse
	err := json.Unmarshal(wSubmit.Body.Bytes(), &submitResp)
	require.NoError(t, err)
	assert.NotEmpty(t, submitResp.ID)
	assert.NotEmpty(t, submitResp.ProtocolNumber)
	assert.Contains(t, submitResp.ProtocolNumber, "A11Y-")

	// 2. Admin queries submitted feedbacks
	reqList, _ := http.NewRequest("GET", "/api/v1/accessibility/feedback", nil)
	reqList.Header.Set("X-Role", "admin")
	wList := httptest.NewRecorder()
	r.ServeHTTP(wList, reqList)
	require.Equal(t, http.StatusOK, wList.Code)

	var listResp struct {
		Feedbacks []accessibility.AccessibilityFeedback `json:"feedbacks"`
		Total     int                                    `json:"total"`
	}
	err = json.Unmarshal(wList.Body.Bytes(), &listResp)
	require.NoError(t, err)
	assert.Equal(t, 1, listResp.Total)

	// 3. Admin updates feedback status to 'resolved'
	updateStatusPayload := map[string]string{
		"status":         "resolved",
		"response_notes": "Rapporto di contrasto elevato a 4.8:1 con il nuovo tema accessibile",
	}
	bodyStatus, _ := json.Marshal(updateStatusPayload)
	reqStatus, _ := http.NewRequest("PATCH", "/api/v1/accessibility/feedback/"+submitResp.ID+"/status", bytes.NewReader(bodyStatus))
	reqStatus.Header.Set("Content-Type", "application/json")
	reqStatus.Header.Set("X-Role", "admin")
	wStatus := httptest.NewRecorder()
	r.ServeHTTP(wStatus, reqStatus)
	require.Equal(t, http.StatusOK, wStatus.Code)

	// 4. Save and retrieve user custom accessibility preferences (Dyslexia font & High contrast)
	prefPayload := accessibility.SavePreferencesRequest{
		Settings: map[string]interface{}{
			"high_contrast":  true,
			"dyslexia_font":  true,
			"font_scaling":   1.2,
			"reduce_motion":  true,
		},
	}
	bodyPref, _ := json.Marshal(prefPayload)
	reqPref, _ := http.NewRequest("POST", "/api/v1/accessibility/preferences", bytes.NewReader(bodyPref))
	reqPref.Header.Set("Content-Type", "application/json")
	reqPref.Header.Set("X-User-ID", "user-parent-1")
	wPref := httptest.NewRecorder()
	r.ServeHTTP(wPref, reqPref)
	require.Equal(t, http.StatusOK, wPref.Code)

	// 5. Get saved preferences
	reqGetPref, _ := http.NewRequest("GET", "/api/v1/accessibility/preferences", nil)
	reqGetPref.Header.Set("X-User-ID", "user-parent-1")
	wGetPref := httptest.NewRecorder()
	r.ServeHTTP(wGetPref, reqGetPref)
	require.Equal(t, http.StatusOK, wGetPref.Code)

	var savedPrefs map[string]interface{}
	err = json.Unmarshal(wGetPref.Body.Bytes(), &savedPrefs)
	require.NoError(t, err)
	assert.Equal(t, true, savedPrefs["high_contrast"])
	assert.Equal(t, true, savedPrefs["dyslexia_font"])
}
