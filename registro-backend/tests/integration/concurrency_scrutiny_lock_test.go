package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_Concurrency_Scrutiny_Lock(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var mu sync.Mutex
	isLocked := false
	gradesCount := 0

	r := gin.New()
	r.POST("/scrutiny/lock", func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()
		isLocked = true
		c.JSON(http.StatusOK, gin.H{"status": "locked"})
	})
	r.POST("/grades", func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()
		if isLocked {
			c.JSON(http.StatusConflict, gin.H{"error": "scrutiny session locked: grade modification forbidden"})
			return
		}
		gradesCount++
		c.JSON(http.StatusCreated, gin.H{"status": "grade_added"})
	})

	// 1. Submit grade before scrutiny lock
	body1, _ := json.Marshal(map[string]interface{}{"value": 8.5})
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/grades", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// 2. Lock scrutiny session
	wLock := httptest.NewRecorder()
	reqLock, _ := http.NewRequest("POST", "/scrutiny/lock", nil)
	r.ServeHTTP(wLock, reqLock)
	assert.Equal(t, http.StatusOK, wLock.Code)

	// 3. Concurrent attempts by multiple teachers to add grades after lock
	var wg sync.WaitGroup
	concurrentAttempts := 10
	conflictCount := 0
	var conflictMu sync.Mutex

	for i := 0; i < concurrentAttempts; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/grades", bytes.NewBuffer(body1))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			if w.Code == http.StatusConflict {
				conflictMu.Lock()
				conflictCount++
				conflictMu.Unlock()
			}
		}()
	}

	wg.Wait()
	assert.Equal(t, concurrentAttempts, conflictCount)
	assert.Equal(t, 1, gradesCount) // No extra grades accepted after lock
}
