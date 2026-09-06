package subjects

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"registro-backend/internal/cache"
)

func TestSubjects_ListCacheAndInvalidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(MockRepository)
	svc := NewService(repo)
	memCache := cache.NewMemoryCache()
	defer func() { _ = memCache.Close() }()

	h := NewHandler(svc, memCache)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "u-123")
		c.Set("role", "admin")
		c.Set("school_id", "sch-1")
		c.Next()
	})
	r.GET("/subjects", h.List)
	r.POST("/subjects", h.Create)

	schoolID := "sch-1"
	subjList := []Subject{
		{ID: "sub-1", Name: "Matematica", Code: "MAT", SchoolID: schoolID},
	}

	// 1. First List: mock repo called once
	repo.On("List", mock.Anything, schoolID).Return(subjList, nil).Once()

	req1 := httptest.NewRequest(http.MethodGet, "/subjects", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, "MISS", w1.Header().Get("X-Cache"))
	assert.Contains(t, w1.Body.String(), "Matematica")

	// 2. Second List: cache hit, repo.List is NOT called again
	req2 := httptest.NewRequest(http.MethodGet, "/subjects", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, "HIT", w2.Header().Get("X-Cache"))
	assert.Contains(t, w2.Body.String(), "Matematica")

	// 3. Create a subject: invalidates cache
	repo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
	createPayload, _ := json.Marshal(CreateSubjectRequest{
		Name: "Fisica",
		Code: "FIS",
	})
	reqCreate := httptest.NewRequest(http.MethodPost, "/subjects", bytes.NewBuffer(createPayload))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	r.ServeHTTP(wCreate, reqCreate)
	assert.Equal(t, http.StatusCreated, wCreate.Code)

	// 4. Third List: cache was invalidated, so it is a MISS and repo.List is called again
	repo.On("List", mock.Anything, schoolID).Return(append(subjList, Subject{ID: "sub-2", Name: "Fisica", Code: "FIS", SchoolID: schoolID}), nil).Once()

	req3 := httptest.NewRequest(http.MethodGet, "/subjects", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusOK, w3.Code)
	assert.Equal(t, "MISS", w3.Header().Get("X-Cache"))
	assert.Contains(t, w3.Body.String(), "Fisica")

	repo.AssertExpectations(t)
}
