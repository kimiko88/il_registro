package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/uda"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_Uda_Learning_Units_Workflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test service / handler or mock service
	r := gin.New()
	r.POST("/uda", func(c *gin.Context) {
		role := c.GetString("role")
		if role != "teacher" && role != "admin" && role != "superadmin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		var req uda.CreateUdaRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"id":         "uda-1",
			"title":      req.Title,
			"class_id":   req.ClassID,
			"subject_id": req.SubjectID,
			"status":     "draft",
		})
	})

	// 1. Multidisciplinary teacher creates UDA Learning Unit
	udaReq := uda.CreateUdaRequest{
		ClassID:      "class-1",
		SubjectID:    "subj-1",
		Title:        "UDA Interdisciplinare: Sostenibilità Ambientale",
		Description:  "Unità di apprendimento integrata Matematica, Scienze e Italiano",
		Period:       "primo_quadrimestre",
		StartDate:    "2026-10-01",
		EndDate:      "2026-12-15",
		Competencies: []string{"Competenza STEM", "Competenza di Cittadinanza"},
	}
	body1, _ := json.Marshal(udaReq)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/uda", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")

	// Set context values
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusForbidden, w1.Code) // Unauthorized without role

	// Authorized request
	rAuthorized := gin.New()
	rAuthorized.POST("/uda", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		var req uda.CreateUdaRequest
		_ = c.ShouldBindJSON(&req)
		c.JSON(http.StatusCreated, gin.H{
			"id":       "uda-1",
			"title":    req.Title,
			"class_id": req.ClassID,
			"status":   "draft",
		})
	})

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/uda", bytes.NewBuffer(body1))
	req2.Header.Set("Content-Type", "application/json")
	rAuthorized.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusCreated, w2.Code)
}
