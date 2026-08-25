package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRefreshTokenCookieHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := &Handler{}

	r.POST("/auth/refresh-token", h.RefreshToken)

	// Case 1: Missing refresh_token in both body and cookie -> 400 Bad Request
	req, _ := http.NewRequest(http.MethodPost, "/auth/refresh-token", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "missing refresh_token")
}
