package search

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGlobalSearchEndpoint_NonSuperAdmin_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &Handler{}

	r := gin.New()
	r.GET("/search/global", h.GlobalSearchEndpoint)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/search/global?q=test", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", "user-teacher")
	c.Set("role", "teacher")

	h.GlobalSearchEndpoint(c)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestValidFilterTypes_IncludesUsersAndLessons(t *testing.T) {
	assert.True(t, validFilterTypes["users"])
	assert.True(t, validFilterTypes["lessons"])
	assert.True(t, validSearchFilterTypes["users"])
	assert.True(t, validSearchFilterTypes["lessons"])
}
