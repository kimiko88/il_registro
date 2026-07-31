package reports

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHandler_ExportGradesExcel(t *testing.T) {
	handler := NewHandler(nil)

	tests := []struct {
		name           string
		userID         string
		role           string
		classID        string
		expectedStatus int
	}{
		{
			name:           "unauthorized",
			userID:         "",
			role:           "",
			classID:        "class-1",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "forbidden role",
			userID:         "user-1",
			role:           "student",
			classID:        "class-1",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "missing class_id",
			userID:         "user-1",
			role:           "admin",
			classID:        "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			if tt.userID != "" {
				c.Set("user_id", tt.userID)
			}
			if tt.role != "" {
				c.Set("role", tt.role)
			}

			url := "/reports/grades/excel"
			if tt.classID != "" {
				url += "?class_id=" + tt.classID
			}
			c.Request = httptest.NewRequest("GET", url, nil)

			handler.ExportGradesExcel(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
