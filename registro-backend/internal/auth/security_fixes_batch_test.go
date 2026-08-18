package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIssueWSTicket_RequiresAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHandler(nil, nil, nil, nil)

	r := gin.New()
	r.POST("/api/v1/auth/ws-ticket", h.IssueWSTicket)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/ws-ticket", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	// Without authenticated user in context, should return 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_MalformedRequest_Returns400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewHandler(nil, nil, nil, nil)

	r := gin.New()
	r.POST("/api/v1/auth/login", h.Login)

	// Malformed JSON payload
	body := bytes.NewBufferString("{ invalid json ")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
