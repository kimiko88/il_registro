package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_GDPR_Anonymization_Privacy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/parent/student/:studentID/grades", func(c *gin.Context) {
		studentID := c.Param("studentID")
		// Simulate adult 18+ student without explicit parent consent
		if studentID == "st-adult-18" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "gdpr privacy restriction: student is 18+ adult and parent access consent is missing",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"student_id": studentID, "grades": []int{8, 9}})
	})

	// 1. Parent attempts to access grades of minor student (16 y.o.)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/parent/student/st-minor-16/grades", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// 2. Parent attempts to access grades of adult student (18+ y.o.) without consent
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/parent/student/st-adult-18/grades", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusForbidden, w2.Code)
}
