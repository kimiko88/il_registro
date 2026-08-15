package grades

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetStudentGrades_EmptyActorRole_Denied(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := &service{
		repo: mockRepo,
	}

	res, err := svc.GetStudentGrades(context.Background(), "user-123", "", "student-456")
	assert.Nil(t, res)
	assert.ErrorIs(t, err, ErrUnauthorized)
}

func TestGetMyGrades_NonStudentRole_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &Handler{}

	r := gin.New()
	r.GET("/my-grades", func(c *gin.Context) {
		c.Set("user_id", "teacher-123")
		c.Set("role", "teacher")
		h.GetMyGrades(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/my-grades", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDownloadSemesterReportPDF_FilenameSanitization(t *testing.T) {
	sanitized := sanitizeFilenameParam("student-123\"\r\nContent-Type: text/html")
	assert.NotContains(t, sanitized, "\r")
	assert.NotContains(t, sanitized, "\n")
	assert.NotContains(t, sanitized, "\"")
	assert.Equal(t, "student-123Content-Typetexthtml", sanitized)
}
