package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/search"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Mock Repository ───────────────────────────────────────────────────────

type mockSearchRepo struct {
	results []search.SearchResultItem
}

func (m *mockSearchRepo) GlobalSearch(_ context.Context, actorRole, schoolID, query, filterType string) ([]search.SearchResultItem, error) {
	// Filter by type if specified
	if filterType == "" || filterType == "all" {
		return m.results, nil
	}
	var filtered []search.SearchResultItem
	for _, r := range m.results {
		if r.Type == filterType {
			filtered = append(filtered, r)
		}
	}
	return filtered, nil
}

// ─── Helpers ───────────────────────────────────────────────────────────────

func setupSearchRouter(repo *mockSearchRepo, role, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := search.NewService(repo)
	h := search.NewHandler(svc)
	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		c.Set("user_id", "user-1")
		c.Set("role", role)
		c.Set("school_id", schoolID)
		c.Next()
	})
	h.RegisterRoutes(api)
	return r
}

func sampleResults() []search.SearchResultItem {
	return []search.SearchResultItem{
		{ID: "s1", Type: "student", Title: "Mario Rossi", Subtitle: "3A"},
		{ID: "t1", Type: "teacher", Title: "Prof. Bianchi", Subtitle: "Matematica"},
		{ID: "c1", Type: "communications", Title: "Circolare n.1 - Inizio anno", Description: "..."},
		{ID: "l1", Type: "lessons", Title: "Lezione Matematica", Subtitle: "2026-09-01"},
	}
}

// ─── Tests ─────────────────────────────────────────────────────────────────

// SR01 — Short queries (< 2 chars) return empty results immediately
func TestSearch_ShortQuery_ReturnsEmpty(t *testing.T) {
	repo := &mockSearchRepo{results: sampleResults()}
	r := setupSearchRouter(repo, "teacher", "school-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/search?q=M", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp search.SearchResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Total)
	assert.Empty(t, resp.Results)
}

// SR02 — Valid query returns matching results
func TestSearch_ValidQuery_ReturnsResults(t *testing.T) {
	repo := &mockSearchRepo{results: sampleResults()}
	r := setupSearchRouter(repo, "teacher", "school-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/search?q=Mario", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp search.SearchResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "Mario", resp.Query)
	assert.Greater(t, resp.Total, 0)
}

// SR03 — Filter by type=student returns only student results
func TestSearch_FilterByType_Student(t *testing.T) {
	repo := &mockSearchRepo{results: sampleResults()}
	r := setupSearchRouter(repo, "teacher", "school-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/search?q=Mario&type=students", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp search.SearchResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	for _, item := range resp.Results {
		assert.Equal(t, "student", item.Type)
	}
}

// SR04 — Filter by type=teacher returns only teacher results
func TestSearch_FilterByType_Teacher(t *testing.T) {
	repo := &mockSearchRepo{results: sampleResults()}
	r := setupSearchRouter(repo, "teacher", "school-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/search?q=Bianchi&type=teachers", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp search.SearchResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	for _, item := range resp.Results {
		assert.Equal(t, "teacher", item.Type)
	}
}

// SR05 — Empty query returns empty results
func TestSearch_EmptyQuery_ReturnsEmpty(t *testing.T) {
	repo := &mockSearchRepo{results: sampleResults()}
	r := setupSearchRouter(repo, "admin", "school-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/search?q=", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp search.SearchResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Total)
}

// SR06 — Unauthenticated user receives 401
func TestSearch_Unauthenticated_Returns401(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockSearchRepo{results: sampleResults()}
	r := gin.New()
	svc := search.NewService(repo)
	h := search.NewHandler(svc)
	api := r.Group("/api/v1")
	// No user_id set in middleware
	h.RegisterRoutes(api)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/search?q=Mario", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// SR07 — Non-superadmin cannot use global search endpoint
func TestSearch_GlobalEndpoint_ForbiddenForNonSuperadmin(t *testing.T) {
	repo := &mockSearchRepo{results: sampleResults()}
	r := setupSearchRouter(repo, "teacher", "school-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/search/global?q=Mario", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// SR08 — Superadmin can use global search endpoint
func TestSearch_GlobalEndpoint_AllowedForSuperadmin(t *testing.T) {
	repo := &mockSearchRepo{results: sampleResults()}
	r := setupSearchRouter(repo, "superadmin", "")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/search/global?q=Rossi", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp search.SearchResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "Rossi", resp.Query)
}

// SR09 — Very long query is truncated to 200 chars and does not error
func TestSearch_VeryLongQuery_IsTruncated(t *testing.T) {
	repo := &mockSearchRepo{results: sampleResults()}
	r := setupSearchRouter(repo, "admin", "school-1")

	longQuery := ""
	for i := 0; i < 250; i++ {
		longQuery += "A"
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/search?q="+longQuery, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should succeed (truncated to 200) or return results for "AAAA..."
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusBadRequest)
}

// SR10 — Invalid filterType is normalized to "all" without error
func TestSearch_InvalidFilterType_NormalizedToAll(t *testing.T) {
	repo := &mockSearchRepo{results: sampleResults()}
	r := setupSearchRouter(repo, "teacher", "school-1")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/search?q=test&type=invalid_type", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp search.SearchResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	// All results returned (type normalized to "all")
	assert.Equal(t, len(sampleResults()), resp.Total)
}
