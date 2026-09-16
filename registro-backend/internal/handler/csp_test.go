package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"registro-backend/internal/handler"
	"registro-backend/pkg/logger"
)

func TestHandleCSPReport(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init("test")

	r := gin.New()
	r.POST("/api/v1/public/csp-report", handler.HandleCSPReport)

	payload := `{
		"csp-report": {
			"document-uri": "https://school.example.com/grades",
			"referrer": "",
			"blocked-uri": "http://malicious-cdn.com/evil.js",
			"violated-directive": "script-src 'self'",
			"original-policy": "default-src 'self'"
		}
	}`

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/public/csp-report", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/csp-report")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
