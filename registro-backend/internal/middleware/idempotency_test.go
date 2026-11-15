package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestIdempotencyScopesResponsesAndRejectsPayloadReuse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var calls atomic.Int32
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", c.GetHeader("X-Test-User"))
		c.Set("school_id", c.GetHeader("X-Test-School"))
		c.Next()
	})
	router.Use(IdempotencyMiddleware())
	router.POST("/mutate", func(c *gin.Context) {
		userID := c.GetString("user_id")
		schoolID := c.GetString("school_id")
		calls.Add(1)
		var body map[string]string
		_ = c.ShouldBindJSON(&body)
		c.JSON(http.StatusCreated, gin.H{"user": userID, "school": schoolID, "value": body["value"]})
	})

	request := func(user, school, key, value string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/mutate", strings.NewReader(`{"value":"`+value+`"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Test-User", user)
		req.Header.Set("X-Test-School", school)
		req.Header.Set(idempotencyHeader, key)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec
	}

	first := request("user-1", "school-1", "operation-key-123", "one")
	require.Equal(t, http.StatusCreated, first.Code)
	replay := request("user-1", "school-1", "operation-key-123", "one")
	require.Equal(t, http.StatusCreated, replay.Code)
	require.Equal(t, "true", replay.Header().Get("X-Idempotency-Replayed"))
	require.Equal(t, first.Body.String(), replay.Body.String())
	require.EqualValues(t, 1, calls.Load())

	mismatch := request("user-1", "school-1", "operation-key-123", "two")
	require.Equal(t, http.StatusConflict, mismatch.Code)
	require.EqualValues(t, 1, calls.Load())

	otherUser := request("user-2", "school-1", "operation-key-123", "one")
	require.Equal(t, http.StatusCreated, otherUser.Code)
	require.NotContains(t, otherUser.Body.String(), "user-1")
	require.EqualValues(t, 2, calls.Load())
}
