package students

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetFascicolo_Unauthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewFascicoloHandler(nil)

	r := gin.New()
	handler.RegisterRoutes(r.Group("/"))

	req, _ := http.NewRequest("GET", "/students/student-123/fascicolo", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
