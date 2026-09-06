package classes

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

func TestClasses_GetClassSubjectsCacheAndInvalidation(t *testing.T) {
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
	r.GET("/classes/:id/subjects", h.GetClassSubjects)
	r.POST("/classes/:id/subjects", h.AssignSubject)

	classID := "class-42"
	subList := []ClassSubject{
		{ID: "cs-1", ClassID: classID, SubjectID: "sub-1", SubjectName: "Matematica"},
	}

	// 1. First call: cache miss, repo called once
	repo.On("Get", mock.Anything, classID).Return(&Class{ID: classID, SchoolID: "sch-1"}, nil).Maybe()
	repo.On("GetClassSubjects", mock.Anything, classID).Return(subList, nil).Once()

	req1 := httptest.NewRequest(http.MethodGet, "/classes/"+classID+"/subjects", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, "MISS", w1.Header().Get("X-Cache"))
	assert.Contains(t, w1.Body.String(), "Matematica")

	// 2. Second call: cache hit, repo.GetClassSubjects NOT called again
	req2 := httptest.NewRequest(http.MethodGet, "/classes/"+classID+"/subjects", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, "HIT", w2.Header().Get("X-Cache"))
	assert.Contains(t, w2.Body.String(), "Matematica")

	// 3. Assign Subject: invalidates cache
	prof := "prof-1"
	repo.On("AssignSubject", mock.Anything, classID, "sub-2", &prof, float64(0)).Return(nil).Once()
	assignPayload, _ := json.Marshal(AssignSubjectRequest{
		SubjectID: "sub-2",
		TeacherID: &prof,
	})
	reqAssign := httptest.NewRequest(http.MethodPost, "/classes/"+classID+"/subjects", bytes.NewBuffer(assignPayload))
	reqAssign.Header.Set("Content-Type", "application/json")
	wAssign := httptest.NewRecorder()
	r.ServeHTTP(wAssign, reqAssign)
	assert.Equal(t, http.StatusCreated, wAssign.Code)

	// 4. Third call: cache is a MISS after invalidation
	repo.On("GetClassSubjects", mock.Anything, classID).Return(append(subList, ClassSubject{ID: "cs-2", ClassID: classID, SubjectID: "sub-2", SubjectName: "Storia"}), nil).Once()

	req3 := httptest.NewRequest(http.MethodGet, "/classes/"+classID+"/subjects", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusOK, w3.Code)
	assert.Equal(t, "MISS", w3.Header().Get("X-Cache"))
	assert.Contains(t, w3.Body.String(), "Storia")

	repo.AssertExpectations(t)
}
