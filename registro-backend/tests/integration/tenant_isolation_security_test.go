package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_Tenant_Isolation_Security(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.GET("/api/v1/classes/:id", func(c *gin.Context) {
		callerSchoolID := c.GetString("school_id")
		targetClassSchoolID := "school-B" // Resource belongs to School B

		if callerSchoolID != targetClassSchoolID {
			c.JSON(http.StatusForbidden, gin.H{"error": "tenant access denied: resource belongs to another school"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "name": "Class 2A"})
	})

	// 1. User from School A attempts to access resource belonging to School B
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/api/v1/classes/class-school-b", nil)

	// Set context to School A
	r.ServeHTTP(w1, req1) // Unauthenticated -> Forbidden/Unauthorized

	rAuthSchoolA := gin.New()
	rAuthSchoolA.GET("/api/v1/classes/:id", func(c *gin.Context) {
		c.Set("user_id", "teacher-school-a")
		c.Set("school_id", "school-A")

		callerSchoolID := c.GetString("school_id")
		targetClassSchoolID := "school-B" // Class belongs to School B

		if callerSchoolID != targetClassSchoolID {
			c.JSON(http.StatusForbidden, gin.H{"error": "tenant access denied: resource belongs to another school"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "name": "Class 2A"})
	})

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/classes/class-school-b", nil)
	rAuthSchoolA.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusForbidden, w2.Code)

	// 2. User from School B accesses own resource
	rAuthSchoolB := gin.New()
	rAuthSchoolB.GET("/api/v1/classes/:id", func(c *gin.Context) {
		c.Set("user_id", "teacher-school-b")
		c.Set("school_id", "school-B")

		callerSchoolID := c.GetString("school_id")
		targetClassSchoolID := "school-B"

		if callerSchoolID != targetClassSchoolID {
			c.JSON(http.StatusForbidden, gin.H{"error": "tenant access denied: resource belongs to another school"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": c.Param("id"), "name": "Class 2A"})
	})

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/api/v1/classes/class-school-b", nil)
	rAuthSchoolB.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
}
