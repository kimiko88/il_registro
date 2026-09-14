package strike

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ─── In-memory Repository mock ────────────────────────────────────────────────

type inMemStrikeRepo struct {
	notices      map[string]*StrikeNotice
	declarations map[string]*StrikeDeclaration
	schoolID     string
}

func newInMemStrikeRepo(schoolID string) *inMemStrikeRepo {
	return &inMemStrikeRepo{
		notices:      make(map[string]*StrikeNotice),
		declarations: make(map[string]*StrikeDeclaration),
		schoolID:     schoolID,
	}
}

func (r *inMemStrikeRepo) ResolveSchoolID(_ context.Context, _ string) string {
	return r.schoolID
}
func (r *inMemStrikeRepo) CreateNotice(_ context.Context, n *StrikeNotice) error {
	n.ID = "notice-1"
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()
	r.notices[n.ID] = n
	return nil
}
func (r *inMemStrikeRepo) GetNoticeByID(_ context.Context, id, _ string) (*StrikeNotice, error) {
	n, ok := r.notices[id]
	if !ok {
		return nil, ErrNotFound
	}
	return n, nil
}
func (r *inMemStrikeRepo) ListNotices(_ context.Context, _ string) ([]*StrikeNotice, error) {
	var list []*StrikeNotice
	for _, n := range r.notices {
		list = append(list, n)
	}
	return list, nil
}
func (r *inMemStrikeRepo) DeleteNotice(_ context.Context, id, _ string) error {
	if _, ok := r.notices[id]; !ok {
		return ErrNotFound
	}
	delete(r.notices, id)
	return nil
}
func (r *inMemStrikeRepo) UpsertDeclaration(_ context.Context, d *StrikeDeclaration) error {
	d.DeclaredAt = time.Now()
	r.declarations[d.StrikeNoticeID+":"+d.UserID] = d
	return nil
}
func (r *inMemStrikeRepo) GetDeclaration(_ context.Context, noticeID, userID string) (*StrikeDeclaration, error) {
	d, ok := r.declarations[noticeID+":"+userID]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}
func (r *inMemStrikeRepo) GetNoticeSummary(_ context.Context, noticeID, _ string) (*StrikeNoticeSummaryResponse, error) {
	n, ok := r.notices[noticeID]
	if !ok {
		return nil, ErrNotFound
	}
	return &StrikeNoticeSummaryResponse{Notice: n}, nil
}
func (r *inMemStrikeRepo) CreateBachecaCommunication(_ context.Context, _, _, _, _ string, _ time.Time) (string, error) {
	return "comm-xyz", nil
}

// ─── Router setup ─────────────────────────────────────────────────────────────

func setupStrikeRouter(h *Handler, userID, role, schoolID string) *gin.Engine {
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

func newStrikeHandlerAndRepo() (*Handler, *inMemStrikeRepo) {
	repo := newInMemStrikeRepo("school-1")
	svc := NewService(repo)
	return NewHandler(svc), repo
}

// ─── CreateNotice ─────────────────────────────────────────────────────────────

func TestStrikeHandler_CreateNotice_Unauthorized(t *testing.T) {
	h, _ := newStrikeHandlerAndRepo()
	r := setupStrikeRouter(h, "", "dsga", "school-1")
	body, _ := json.Marshal(CreateStrikeNoticeRequest{
		Title: "Sciopero", ProclaimedBy: "FLC", StrikeDate: "2026-11-01",
		DeclarationDeadline: "2026-10-30T12:00:00Z",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/strike-notices", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestStrikeHandler_CreateNotice_BadRequest_MissingFields(t *testing.T) {
	h, _ := newStrikeHandlerAndRepo()
	r := setupStrikeRouter(h, "dsga-1", "dsga", "school-1")
	body := []byte(`{"title":"Sciopero"}`) // missing required fields
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/strike-notices", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStrikeHandler_CreateNotice_Forbidden_TeacherRole(t *testing.T) {
	h, _ := newStrikeHandlerAndRepo()
	r := setupStrikeRouter(h, "teacher-1", "teacher", "school-1")
	body, _ := json.Marshal(CreateStrikeNoticeRequest{
		Title: "Sciopero", ProclaimedBy: "COBAS", StrikeDate: "2026-11-05",
		DeclarationDeadline: "2026-11-03T12:00:00Z",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/strike-notices", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestStrikeHandler_CreateNotice_Success_DSGA(t *testing.T) {
	h, _ := newStrikeHandlerAndRepo()
	r := setupStrikeRouter(h, "dsga-1", "dsga", "school-1")
	body, _ := json.Marshal(CreateStrikeNoticeRequest{
		Title:               "Sciopero Nazionale",
		ProclaimedBy:        "FLC CGIL",
		StrikeDate:          "2026-12-01",
		DeclarationDeadline: "2026-11-28T12:00:00Z",
		Content:             "Avviso ufficiale di sciopero",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/strike-notices", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp StrikeNotice
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "Sciopero Nazionale", resp.Title)
}

func TestStrikeHandler_CreateNotice_Success_Principal(t *testing.T) {
	h, _ := newStrikeHandlerAndRepo()
	r := setupStrikeRouter(h, "principal-1", "principal", "school-1")
	body, _ := json.Marshal(CreateStrikeNoticeRequest{
		Title:               "Sciopero Regionale",
		ProclaimedBy:        "SNALS",
		StrikeDate:          "2026-12-10",
		DeclarationDeadline: "2026-12-08T10:00:00Z",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/strike-notices", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

// ─── ListNotices ──────────────────────────────────────────────────────────────

func TestStrikeHandler_ListNotices_Unauthorized(t *testing.T) {
	h, _ := newStrikeHandlerAndRepo()
	r := setupStrikeRouter(h, "", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/strike-notices", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestStrikeHandler_ListNotices_Empty(t *testing.T) {
	h, _ := newStrikeHandlerAndRepo()
	r := setupStrikeRouter(h, "teacher-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/strike-notices", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestStrikeHandler_ListNotices_WithData(t *testing.T) {
	h, repo := newStrikeHandlerAndRepo()
	repo.notices["n-1"] = &StrikeNotice{
		ID: "n-1", SchoolID: "school-1", Title: "Sciopero 1",
		DeclarationDeadline: time.Now().Add(24 * time.Hour),
	}
	r := setupStrikeRouter(h, "teacher-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/strike-notices", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── GetNotice ────────────────────────────────────────────────────────────────

func TestStrikeHandler_GetNotice_Unauthorized(t *testing.T) {
	h, _ := newStrikeHandlerAndRepo()
	r := setupStrikeRouter(h, "", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/strike-notices/n-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestStrikeHandler_GetNotice_NotFound(t *testing.T) {
	h, _ := newStrikeHandlerAndRepo()
	r := setupStrikeRouter(h, "teacher-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/strike-notices/ghost", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestStrikeHandler_GetNotice_Success(t *testing.T) {
	h, repo := newStrikeHandlerAndRepo()
	repo.notices["notice-1"] = &StrikeNotice{
		ID: "notice-1", SchoolID: "school-1", Title: "Sciopero Test",
		DeclarationDeadline: time.Now().Add(48 * time.Hour),
	}
	r := setupStrikeRouter(h, "teacher-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/strike-notices/notice-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── DeleteNotice ─────────────────────────────────────────────────────────────

func TestStrikeHandler_DeleteNotice_Unauthorized(t *testing.T) {
	h, _ := newStrikeHandlerAndRepo()
	r := setupStrikeRouter(h, "", "dsga", "school-1")
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/strike-notices/n-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestStrikeHandler_DeleteNotice_Forbidden_TeacherRole(t *testing.T) {
	h, repo := newStrikeHandlerAndRepo()
	repo.notices["n-1"] = &StrikeNotice{ID: "n-1", SchoolID: "school-1"}
	r := setupStrikeRouter(h, "teacher-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/strike-notices/n-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestStrikeHandler_DeleteNotice_Success(t *testing.T) {
	h, repo := newStrikeHandlerAndRepo()
	repo.notices["n-del"] = &StrikeNotice{ID: "n-del", SchoolID: "school-1"}
	r := setupStrikeRouter(h, "dsga-1", "dsga", "school-1")
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/strike-notices/n-del", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotContains(t, repo.notices, "n-del")
}

// ─── SubmitDeclaration ────────────────────────────────────────────────────────

func TestStrikeHandler_SubmitDeclaration_Unauthorized(t *testing.T) {
	h, _ := newStrikeHandlerAndRepo()
	r := setupStrikeRouter(h, "", "teacher", "school-1")
	body, _ := json.Marshal(SubmitDeclarationRequest{Intention: IntentionParticipates})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/strike-notices/n-1/declare", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestStrikeHandler_SubmitDeclaration_NotFound(t *testing.T) {
	h, _ := newStrikeHandlerAndRepo()
	r := setupStrikeRouter(h, "teacher-1", "teacher", "school-1")
	body, _ := json.Marshal(SubmitDeclarationRequest{Intention: IntentionParticipates})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/strike-notices/ghost/declare", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestStrikeHandler_SubmitDeclaration_InvalidIntention(t *testing.T) {
	h, _ := newStrikeHandlerAndRepo()
	r := setupStrikeRouter(h, "teacher-1", "teacher", "school-1")
	body, _ := json.Marshal(SubmitDeclarationRequest{Intention: "maybe"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/strike-notices/n-1/declare", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStrikeHandler_SubmitDeclaration_DeadlinePassed(t *testing.T) {
	h, repo := newStrikeHandlerAndRepo()
	repo.notices["n-expired"] = &StrikeNotice{
		ID:                  "n-expired",
		SchoolID:            "school-1",
		DeclarationDeadline: time.Now().Add(-2 * time.Hour), // expired
	}
	r := setupStrikeRouter(h, "teacher-1", "teacher", "school-1")
	body, _ := json.Marshal(SubmitDeclarationRequest{Intention: IntentionNotParticipates})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/strike-notices/n-expired/declare", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestStrikeHandler_SubmitDeclaration_Success(t *testing.T) {
	h, repo := newStrikeHandlerAndRepo()
	repo.notices["n-active"] = &StrikeNotice{
		ID:                  "n-active",
		SchoolID:            "school-1",
		DeclarationDeadline: time.Now().Add(72 * time.Hour),
	}
	r := setupStrikeRouter(h, "teacher-1", "teacher", "school-1")
	body, _ := json.Marshal(SubmitDeclarationRequest{
		Intention: IntentionUndecided, Notes: "Non so ancora",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/strike-notices/n-active/declare", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── GetNoticeSummary ─────────────────────────────────────────────────────────

func TestStrikeHandler_GetNoticeSummary_Unauthorized(t *testing.T) {
	h, _ := newStrikeHandlerAndRepo()
	r := setupStrikeRouter(h, "", "dsga", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/strike-notices/n-1/summary", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestStrikeHandler_GetNoticeSummary_Forbidden_Teacher(t *testing.T) {
	h, _ := newStrikeHandlerAndRepo()
	r := setupStrikeRouter(h, "teacher-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/strike-notices/n-1/summary", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestStrikeHandler_GetNoticeSummary_Success_DSGA(t *testing.T) {
	h, repo := newStrikeHandlerAndRepo()
	repo.notices["n-sum"] = &StrikeNotice{
		ID: "n-sum", SchoolID: "school-1", Title: "Sciopero Riepilogo",
		DeclarationDeadline: time.Now().Add(24 * time.Hour),
	}
	r := setupStrikeRouter(h, "dsga-1", "dsga", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/strike-notices/n-sum/summary", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
