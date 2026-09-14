package verbali

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ─── In-memory Repository mock ───────────────────────────────────────────────

type mockHandlerRepo struct {
	meetings  map[string]*CouncilMeeting
	verbali   map[string]*MeetingVerbale
	sigs      map[string][]VerbaleSignature
	templates map[string]*MeetingVerbaleTemplate
}

func newMockHandlerRepo() *mockHandlerRepo {
	return &mockHandlerRepo{
		meetings:  make(map[string]*CouncilMeeting),
		verbali:   make(map[string]*MeetingVerbale),
		sigs:      make(map[string][]VerbaleSignature),
		templates: make(map[string]*MeetingVerbaleTemplate),
	}
}

func (m *mockHandlerRepo) ResolveSchoolID(_ context.Context, _ string) string {
	return "school-test"
}
func (m *mockHandlerRepo) ClassBelongsToSchool(_ context.Context, _, _ string) (bool, error) {
	return true, nil
}
func (m *mockHandlerRepo) CreateMeeting(_ context.Context, cm *CouncilMeeting) error {
	cm.ID = "meeting-1"
	cm.CreatedAt = time.Now()
	m.meetings[cm.ID] = cm
	return nil
}
func (m *mockHandlerRepo) ListMeetings(_ context.Context, _, _ string) ([]*CouncilMeeting, error) {
	var list []*CouncilMeeting
	for _, cm := range m.meetings {
		list = append(list, cm)
	}
	return list, nil
}
func (m *mockHandlerRepo) GetMeetingByID(_ context.Context, id string) (*CouncilMeeting, error) {
	cm, ok := m.meetings[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return cm, nil
}
func (m *mockHandlerRepo) GetMeetingWithCoordinator(_ context.Context, meetingID string) (*CouncilMeeting, string, error) {
	cm, ok := m.meetings[meetingID]
	if !ok {
		return nil, "", errors.New("meeting not found")
	}
	return cm, "coord-1", nil
}
func (m *mockHandlerRepo) CreateVerbale(_ context.Context, v *MeetingVerbale) error {
	v.ID = "verb-1"
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()
	m.verbali[v.ID] = v
	return nil
}
func (m *mockHandlerRepo) GetVerbaleByID(_ context.Context, id, _ string) (*MeetingVerbale, error) {
	v, ok := m.verbali[id]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}
func (m *mockHandlerRepo) UpdateVerbale(_ context.Context, v *MeetingVerbale) error {
	v.UpdatedAt = time.Now()
	m.verbali[v.ID] = v
	return nil
}
func (m *mockHandlerRepo) DeleteVerbale(_ context.Context, id string) error {
	if _, ok := m.verbali[id]; !ok {
		return ErrNotFound
	}
	delete(m.verbali, id)
	return nil
}
func (m *mockHandlerRepo) ListVerbali(_ context.Context, meetingID, _ string, _ bool) ([]*MeetingVerbale, error) {
	var list []*MeetingVerbale
	for _, v := range m.verbali {
		if v.MeetingID == meetingID {
			list = append(list, v)
		}
	}
	return list, nil
}
func (m *mockHandlerRepo) ListAllVerbali(_ context.Context, _, _, _ string, _ bool) ([]*MeetingVerbale, error) {
	var list []*MeetingVerbale
	for _, v := range m.verbali {
		list = append(list, v)
	}
	return list, nil
}
func (m *mockHandlerRepo) SignVerbale(_ context.Context, verbaleID, userID, ipAddress string) error {
	v, ok := m.verbali[verbaleID]
	if !ok {
		return ErrNotFound
	}
	m.sigs[verbaleID] = append(m.sigs[verbaleID], VerbaleSignature{
		VerbaleID: verbaleID,
		UserID:    userID,
		IPAddress: ipAddress,
		SignedAt:  time.Now(),
	})
	v.IsSigned = true
	return nil
}
func (m *mockHandlerRepo) MarkVerbaleSigned(_ context.Context, verbaleID string) error {
	v, ok := m.verbali[verbaleID]
	if !ok {
		return ErrNotFound
	}
	now := time.Now()
	v.IsSigned = true
	v.Status = "signed"
	v.SignedAt = &now
	return nil
}
func (m *mockHandlerRepo) GetSignatures(_ context.Context, verbaleID string) ([]VerbaleSignature, error) {
	return m.sigs[verbaleID], nil
}
func (m *mockHandlerRepo) ListTemplates(_ context.Context, _, _ string) ([]*MeetingVerbaleTemplate, error) {
	var list []*MeetingVerbaleTemplate
	for _, t := range m.templates {
		list = append(list, t)
	}
	return list, nil
}
func (m *mockHandlerRepo) GetTemplateByID(_ context.Context, id string) (*MeetingVerbaleTemplate, error) {
	t, ok := m.templates[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}
func (m *mockHandlerRepo) CreateTemplate(_ context.Context, t *MeetingVerbaleTemplate) error {
	t.ID = "tpl-1"
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	m.templates[t.ID] = t
	return nil
}
func (m *mockHandlerRepo) UpdateTemplate(_ context.Context, t *MeetingVerbaleTemplate) error {
	m.templates[t.ID] = t
	return nil
}
func (m *mockHandlerRepo) DeleteTemplate(_ context.Context, id, _ string) error {
	if _, ok := m.templates[id]; !ok {
		return ErrNotFound
	}
	delete(m.templates, id)
	return nil
}

// ─── Router setup ────────────────────────────────────────────────────────────

func setupVerbaliRouter(h *Handler, userID, role, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Set("role", role)
		c.Set("school_id", schoolID)
		c.Next()
	})
	api := r.Group("/api/v1")
	h.RegisterRoutes(api)
	return r
}

func newVerbaliHandler() (*Handler, *mockHandlerRepo) {
	repo := newMockHandlerRepo()
	svc := NewService(repo)
	return NewHandler(svc), repo
}

// ─── CreateMeeting ───────────────────────────────────────────────────────────

func TestVerbaliHandler_CreateMeeting_Unauthorized(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "", "teacher", "school-1")
	body, _ := json.Marshal(CreateMeetingRequest{Title: "Test", Date: "2026-10-10", StartTime: "09:00", EndTime: "11:00"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/verbali/meetings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestVerbaliHandler_CreateMeeting_Forbidden(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "u-1", "student", "school-1")
	body, _ := json.Marshal(CreateMeetingRequest{Title: "Test", Date: "2026-10-10", StartTime: "09:00", EndTime: "11:00"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/verbali/meetings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestVerbaliHandler_CreateMeeting_BadRequest(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "u-1", "teacher", "school-1")
	// missing required fields: date, start_time, end_time
	body := []byte(`{"title":"No Date"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/verbali/meetings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVerbaliHandler_CreateMeeting_Success(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "teacher-1", "teacher", "school-1")
	body, _ := json.Marshal(CreateMeetingRequest{
		Title: "Consiglio di Classe", Date: "2026-10-15", StartTime: "14:00", EndTime: "16:00",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/verbali/meetings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp CouncilMeeting
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "Consiglio di Classe", resp.Title)
}

// ─── ListMeetings ────────────────────────────────────────────────────────────

func TestVerbaliHandler_ListMeetings_Success(t *testing.T) {
	h, repo := newVerbaliHandler()
	repo.meetings["m-1"] = &CouncilMeeting{ID: "m-1", SchoolID: "school-1", Title: "Riunione 1"}
	r := setupVerbaliRouter(h, "teacher-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/meetings", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var meetings []*CouncilMeeting
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &meetings))
	assert.Len(t, meetings, 1)
}

func TestVerbaliHandler_ListMeetings_EmptyList(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "admin-1", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/meetings", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "[]")
}

func TestVerbaliHandler_ListMeetings_Unauthorized(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/meetings", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ─── CreateVerbale ────────────────────────────────────────────────────────────

func TestVerbaliHandler_CreateVerbale_Unauthorized(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "", "teacher", "school-1")
	body, _ := json.Marshal(CreateVerbaleRequest{MeetingID: "m-1", Title: "T", Content: "C"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/verbali", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestVerbaliHandler_CreateVerbale_Forbidden_Role(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "parent-1", "parent", "school-1")
	body, _ := json.Marshal(CreateVerbaleRequest{MeetingID: "m-1", Title: "T", Content: "C"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/verbali", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestVerbaliHandler_CreateVerbale_MeetingNotFound(t *testing.T) {
	h, _ := newVerbaliHandler()
	// meeting "missing-meeting" does not exist in the mock repo
	r := setupVerbaliRouter(h, "admin-1", "admin", "school-1")
	body, _ := json.Marshal(CreateVerbaleRequest{MeetingID: "missing-meeting", Title: "T", Content: "C"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/verbali", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestVerbaliHandler_CreateVerbale_Success_AsAdmin(t *testing.T) {
	h, repo := newVerbaliHandler()
	repo.meetings["m-100"] = &CouncilMeeting{
		ID: "m-100", SchoolID: "school-1", Title: "Meeting", CreatedBy: "admin-1", Date: time.Now(),
	}
	r := setupVerbaliRouter(h, "admin-1", "admin", "school-1")
	body, _ := json.Marshal(CreateVerbaleRequest{
		MeetingID: "m-100", Title: "Verbale del Consiglio", Content: "Ordine del giorno...",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/verbali", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestVerbaliHandler_CreateVerbale_OnlyCoordinatorOrSecretary(t *testing.T) {
	h, repo := newVerbaliHandler()
	// teacher-1 is not the coordinator (coord-1) and is not secretary
	repo.meetings["m-200"] = &CouncilMeeting{
		ID: "m-200", SchoolID: "school-1", Title: "Meeting", Date: time.Now(),
	}
	r := setupVerbaliRouter(h, "teacher-1", "teacher", "school-1")
	body, _ := json.Marshal(CreateVerbaleRequest{
		MeetingID: "m-200", Title: "Verbale", Content: "Contenuto",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/verbali", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ─── GetVerbale ───────────────────────────────────────────────────────────────

func TestVerbaliHandler_GetVerbale_Unauthorized(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/some-id", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestVerbaliHandler_GetVerbale_Forbidden(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "p-1", "parent", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/some-id", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestVerbaliHandler_GetVerbale_NotFound(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "teacher-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/nonexistent-id", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestVerbaliHandler_GetVerbale_Success(t *testing.T) {
	h, repo := newVerbaliHandler()
	repo.verbali["v-ok"] = &MeetingVerbale{
		ID: "v-ok", MeetingID: "m-1", Status: "draft", Title: "Verbale Test",
	}
	repo.meetings["m-1"] = &CouncilMeeting{ID: "m-1"}
	r := setupVerbaliRouter(h, "teacher-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/v-ok", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVerbaliHandler_GetVerbale_DraftHiddenFromPrincipal(t *testing.T) {
	h, repo := newVerbaliHandler()
	repo.verbali["v-draft"] = &MeetingVerbale{
		ID: "v-draft", MeetingID: "m-1", Status: "draft", IsSigned: false,
	}
	repo.meetings["m-1"] = &CouncilMeeting{ID: "m-1"}
	r := setupVerbaliRouter(h, "principal-1", "principal", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/v-draft", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── UpdateVerbale ────────────────────────────────────────────────────────────

func TestVerbaliHandler_UpdateVerbale_NotFound(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "admin-1", "admin", "school-1")
	body, _ := json.Marshal(UpdateVerbaleRequest{Title: "Updated", Content: "New content"})
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/verbali/ghost-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestVerbaliHandler_UpdateVerbale_Locked(t *testing.T) {
	h, repo := newVerbaliHandler()
	now := time.Now()
	repo.verbali["v-signed"] = &MeetingVerbale{
		ID: "v-signed", MeetingID: "m-1", IsSigned: true, Status: "signed", SignedAt: &now,
	}
	repo.meetings["m-1"] = &CouncilMeeting{ID: "m-1"}
	r := setupVerbaliRouter(h, "admin-1", "admin", "school-1")
	body, _ := json.Marshal(UpdateVerbaleRequest{Title: "T", Content: "C"})
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/verbali/v-signed", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "VERBALE_LOCKED")
}

// ─── DeleteVerbale ────────────────────────────────────────────────────────────

func TestVerbaliHandler_DeleteVerbale_Unauthorized(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/verbali/v-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestVerbaliHandler_DeleteVerbale_NotFound(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "admin-1", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/verbali/ghost-id", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ─── ListVerbali / ListAllVerbali ─────────────────────────────────────────────

func TestVerbaliHandler_ListVerbali_InvalidUUID(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "teacher-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/meeting/not-a-uuid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "[]")
}

func TestVerbaliHandler_ListAllVerbali_Forbidden(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "p-1", "parent", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestVerbaliHandler_ListAllVerbali_Success(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "teacher-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── Signatures ──────────────────────────────────────────────────────────────

func TestVerbaliHandler_GetSignatures_Unauthorized(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/v-1/signatures", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestVerbaliHandler_GetSignatures_EmptyList(t *testing.T) {
	h, repo := newVerbaliHandler()
	repo.verbali["v-sig-1"] = &MeetingVerbale{ID: "v-sig-1", Status: "draft", MeetingID: "m-1"}
	r := setupVerbaliRouter(h, "teacher-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/v-sig-1/signatures", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "[]")
}

// ─── Templates ────────────────────────────────────────────────────────────────

func TestVerbaliHandler_ListTemplates_Unauthorized(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/templates", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestVerbaliHandler_ListTemplates_Success(t *testing.T) {
	h, repo := newVerbaliHandler()
	repo.templates["tpl-1"] = &MeetingVerbaleTemplate{ID: "tpl-1", Title: "Modello Base", SchoolID: "school-1"}
	r := setupVerbaliRouter(h, "admin-1", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/templates", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var templates []*MeetingVerbaleTemplate
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &templates))
	assert.Len(t, templates, 1)
}

func TestVerbaliHandler_GetTemplate_NotFound(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "admin-1", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/templates/nonexistent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestVerbaliHandler_CreateTemplate_Success(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "dsga-1", "dsga", "school-1")
	body, _ := json.Marshal(CreateTemplateRequest{
		Title:           "Modello Consiglio Classe",
		MeetingType:     "consiglio_classe",
		DefaultAgenda:   "1. Appello 2. Comunicazioni",
		TemplateContent: "Contenuto verbale",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/verbali/templates", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestVerbaliHandler_CreateTemplate_ForbiddenForTeacher(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "teacher-1", "teacher", "school-1")
	body, _ := json.Marshal(CreateTemplateRequest{
		Title: "Modello", MeetingType: "consiglio_classe",
		DefaultAgenda: "1. Appello", TemplateContent: "...",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/verbali/templates", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestVerbaliHandler_DeleteTemplate_NotFound(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "admin-1", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/verbali/templates/ghost-tpl", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── PDF endpoints (no pdfWorkerClient) ──────────────────────────────────────

func TestVerbaliHandler_ExportPDF_Unauthorized(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/v-1/pdf?sync=true", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestVerbaliHandler_EnqueueAsyncPdf_NoPdfWorker(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "teacher-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/verbali/v-1/async-pdf", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestVerbaliHandler_GetPdfJobStatus_NoPdfWorker(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "teacher-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/pdf-jobs/some-job-id", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestVerbaliHandler_GetPdfJobStatus_Unauthorized(t *testing.T) {
	h, _ := newVerbaliHandler()
	r := setupVerbaliRouter(h, "", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/verbali/pdf-jobs/some-job", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
