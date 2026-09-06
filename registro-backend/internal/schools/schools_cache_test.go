package schools

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"registro-backend/internal/cache"
)

func TestSchools_ListPublicCache(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(MockRepository)
	svc := NewService(repo)
	memCache := cache.NewMemoryCache()
	defer func() { _ = memCache.Close() }()

	h := NewHandler(svc, memCache)

	r := gin.New()
	r.GET("/public/schools", h.ListPublic)

	// Setup mock expectation: List should only be called once by the DB/repo
	repo.On("List", mock.Anything, mock.Anything).Return([]*School{
		{ID: "s1", Name: "Liceo Scientifico", Code: "RMPS01000P", City: "Roma"},
	}, 1, nil).Once()

	// Request 1: cache miss
	req1 := httptest.NewRequest(http.MethodGet, "/public/schools", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, "MISS", w1.Header().Get("X-Cache"))
	assert.Contains(t, w1.Body.String(), "Liceo Scientifico")

	// Request 2: cache hit (repo.List is NOT called again because of .Once())
	req2 := httptest.NewRequest(http.MethodGet, "/public/schools", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, "HIT", w2.Header().Get("X-Cache"))
	assert.Contains(t, w2.Body.String(), "Liceo Scientifico")

	repo.AssertExpectations(t)
}
