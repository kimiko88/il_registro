package accessibility_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"registro-backend/internal/accessibility"
)

type fullMockRepo struct {
	feedbacks     []accessibility.AccessibilityFeedback
	createErr     error
	listErr       error
	statusErr     error
	getPrefErr    error
	upsertPrefErr error
	prefs         string
}

func (m *fullMockRepo) Create(ctx context.Context, f *accessibility.AccessibilityFeedback) error {
	if m.createErr != nil {
		return m.createErr
	}
	f.ID = "fb-123"
	m.feedbacks = append(m.feedbacks, *f)
	return nil
}

func (m *fullMockRepo) List(ctx context.Context, schoolID, status string, limit, offset int) ([]accessibility.AccessibilityFeedback, int, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	return m.feedbacks, len(m.feedbacks), nil
}

func (m *fullMockRepo) GetByID(ctx context.Context, id string) (*accessibility.AccessibilityFeedback, error) {
	for _, f := range m.feedbacks {
		if f.ID == id {
			return &f, nil
		}
	}
	return nil, nil
}

func (m *fullMockRepo) UpdateStatus(ctx context.Context, id, status, responseNotes string) error {
	return m.statusErr
}

func (m *fullMockRepo) GetUserPreferences(ctx context.Context, userID string) (string, error) {
	if m.getPrefErr != nil {
		return "", m.getPrefErr
	}
	if m.prefs != "" {
		return m.prefs, nil
	}
	return `{"highContrast":true}`, nil
}

func (m *fullMockRepo) UpsertUserPreferences(ctx context.Context, userID string, settingsJSON string) error {
	return m.upsertPrefErr
}

func setupA11yRouter(svc accessibility.Service, authUser string, authSchool string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if authUser != "" {
			c.Set("user_id", authUser)
		}
		if authSchool != "" {
			c.Set("school_id", authSchool)
		}
		c.Next()
	})
	h := accessibility.NewHandler(svc)
	rg := r.Group("/api/v1/accessibility")
	{
		rg.POST("/feedback", h.SubmitPublic)
		rg.GET("/feedbacks", h.List)
		rg.PUT("/feedbacks/:id/status", h.UpdateStatus)
		rg.GET("/preferences/me", h.GetMyPreferences)
		rg.PUT("/preferences/me", h.SaveMyPreferences)
	}
	return r
}

func TestAccessibilityService_Extended(t *testing.T) {
	repo := &fullMockRepo{}
	svc := accessibility.NewService(repo)
	ctx := context.Background()

	// Nil request
	_, err := svc.SubmitFeedback(ctx, nil, nil, nil, nil, nil)
	assert.Error(t, err)

	// Missing barrier type
	req := &accessibility.CreateFeedbackRequest{
		Name:        "Test User",
		Email:       "test@example.com",
		BarrierType: "",
		Description: "Problem description long enough",
	}
	_, err = svc.SubmitFeedback(ctx, req, nil, nil, nil, nil)
	assert.Error(t, err)

	// Valid submission
	req.BarrierType = "screen-reader"
	resp, err := svc.SubmitFeedback(ctx, req, nil, nil, nil, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.ProtocolNumber)

	// ListFeedbacks with negative & large limits
	fbs, total, err := svc.ListFeedbacks(ctx, "sch-1", "open", -1, -5)
	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, fbs, 1)

	// UpdateFeedbackStatus validation
	assert.Error(t, svc.UpdateFeedbackStatus(ctx, "", "open", "notes"))
	assert.Error(t, svc.UpdateFeedbackStatus(ctx, "fb-1", "invalid-status", "notes"))
	assert.NoError(t, svc.UpdateFeedbackStatus(ctx, "fb-1", "resolved", "Fixed contrast"))

	// GetUserPreferences
	_, err = svc.GetUserPreferences(ctx, "")
	assert.Error(t, err)
	prefs, err := svc.GetUserPreferences(ctx, "u1")
	assert.NoError(t, err)
	assert.Contains(t, prefs, "highContrast")

	// SaveUserPreferences
	assert.Error(t, svc.SaveUserPreferences(ctx, "", map[string]interface{}{}))
	assert.NoError(t, svc.SaveUserPreferences(ctx, "u1", nil))
	assert.NoError(t, svc.SaveUserPreferences(ctx, "u1", map[string]interface{}{"dsaFont": true}))
}

func TestAccessibilityHandler_Endpoints(t *testing.T) {
	repo := &fullMockRepo{}
	svc := accessibility.NewService(repo)
	r := setupA11yRouter(svc, "u-123", "sch-456")

	// 1. POST /feedback (Public)
	payload := `{"name":"Luigi Verdi","email":"luigi@test.com","barrier_type":"color-blindness","description":"Colori indistinguibili"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/accessibility/feedback", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	req, _ = http.NewRequest(http.MethodPost, "/api/v1/accessibility/feedback", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 2. GET /feedbacks (List)
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/accessibility/feedbacks?status=open", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. PUT /feedbacks/:id/status
	statusPayload := `{"status":"in_progress","response_notes":"Presa in carico"}`
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/accessibility/feedbacks/fb-123/status", bytes.NewBufferString(statusPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest(http.MethodPut, "/api/v1/accessibility/feedbacks/fb-123/status", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 4. GET /preferences/me (Authenticated)
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/accessibility/preferences/me", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. PUT /preferences/me (Save)
	prefPayload := `{"settings":{"fontSize":"large","screenReaderOptimized":true}}`
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/accessibility/preferences/me", bytes.NewBufferString(prefPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req, _ = http.NewRequest(http.MethodPut, "/api/v1/accessibility/preferences/me", bytes.NewBufferString(`{bad-json`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 6. Unauthenticated requests to preferences/me
	unauthRouter := setupA11yRouter(svc, "", "")
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/accessibility/preferences/me", nil)
	w = httptest.NewRecorder()
	unauthRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	req, _ = http.NewRequest(http.MethodPut, "/api/v1/accessibility/preferences/me", bytes.NewBufferString(prefPayload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	unauthRouter.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
