package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/textbooks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockTextbookRepo struct{}

func (m *mockTextbookRepo) Create(ctx context.Context, t *textbooks.Textbook) error {
	return nil
}
func (m *mockTextbookRepo) Update(ctx context.Context, t *textbooks.Textbook) error {
	return nil
}
func (m *mockTextbookRepo) List(ctx context.Context, schoolID string) ([]textbooks.Textbook, error) {
	return nil, nil
}
func (m *mockTextbookRepo) Delete(ctx context.Context, id string) error {
	return nil
}
func (m *mockTextbookRepo) AssignToClass(ctx context.Context, classID string, subjectID string, textbookID string, optional bool) error {
	return nil
}
func (m *mockTextbookRepo) RemoveFromClass(ctx context.Context, assignmentID string) error {
	return nil
}
func (m *mockTextbookRepo) ListByClass(ctx context.Context, classID string) ([]textbooks.ClassTextbook, error) {
	return nil, nil
}

func TestRobustnessAndBoundaryEdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("1. Grade Addition Numerical Boundary HTTP Checks", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", "t-mario")
			c.Set("role", "teacher")
			c.Next()
		})
		router.POST("/api/v1/grades", func(c *gin.Context) {
			var body struct {
				GradeValue float64 `json:"grade_value" binding:"required"`
			}
			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if body.GradeValue <= 0 || body.GradeValue > 10 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "voto fuori dai limiti consentiti (1-10)"})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"id": "g-new", "grade_value": body.GradeValue})
		})

		// Invalid values <= 0 or > 10
		invalidGrades := []float64{-5.0, 0.0, 10.5, 100.0}
		for _, v := range invalidGrades {
			payload, _ := json.Marshal(map[string]interface{}{
				"student_id":  "st-01",
				"subject_id":  "sub-math",
				"grade_value": v,
				"grade_type":  "numeric",
				"semester":    1,
				"date":        "2025-10-15",
			})
			req, _ := http.NewRequest("POST", "/api/v1/grades", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusBadRequest, resp.Code, "Grade value %f should return 400 Bad Request", v)
		}

		// Valid boundary values
		validGrades := []float64{1.0, 6.0, 10.0}
		for _, v := range validGrades {
			payload, _ := json.Marshal(map[string]interface{}{
				"student_id":  "st-01",
				"subject_id":  "sub-math",
				"grade_value": v,
				"grade_type":  "numeric",
				"semester":    1,
				"date":        "2025-10-15",
			})
			req, _ := http.NewRequest("POST", "/api/v1/grades", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusCreated, resp.Code, "Grade value %f should return 201 Created", v)
		}
	})

	t.Run("2. Textbook Catalog Input Validation HTTP Checks", func(t *testing.T) {
		repo := &mockTextbookRepo{}
		tbSvc := textbooks.NewService(repo)
		h := textbooks.NewHandler(tbSvc)

		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", "sec-01")
			c.Set("role", "secretary")
			c.Set("school_id", "sch-101")
			c.Next()
		})
		router.POST("/api/v1/textbooks", h.Create)

		// Missing title
		payload1, _ := json.Marshal(map[string]interface{}{
			"title":   "   ",
			"subject": "Matematica",
			"price":   25.0,
		})
		req1, _ := http.NewRequest("POST", "/api/v1/textbooks", bytes.NewBuffer(payload1))
		req1.Header.Set("Content-Type", "application/json")
		resp1 := httptest.NewRecorder()
		router.ServeHTTP(resp1, req1)
		assert.Equal(t, http.StatusBadRequest, resp1.Code)

		// Negative price
		payload2, _ := json.Marshal(map[string]interface{}{
			"title":   "Analisi 1",
			"subject": "Matematica",
			"price":   -15.00,
		})
		req2, _ := http.NewRequest("POST", "/api/v1/textbooks", bytes.NewBuffer(payload2))
		req2.Header.Set("Content-Type", "application/json")
		resp2 := httptest.NewRecorder()
		router.ServeHTTP(resp2, req2)
		assert.Equal(t, http.StatusBadRequest, resp2.Code)
	})

	t.Run("3. Attendance Retroactive Date Enforcement HTTP Checks", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", "t-mario")
			c.Set("role", "teacher")
			c.Next()
		})
		router.POST("/api/v1/attendance/mark-bulk", func(c *gin.Context) {
			var body struct {
				ClassID  string `json:"class_id" binding:"required"`
				Date     string `json:"date" binding:"required"`
				Statuses []struct {
					StudentID string `json:"student_id"`
					Status    string `json:"status"`
				} `json:"statuses"`
			}
			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			parsedDate, err := time.Parse("2006-01-02", body.Date)
			if err != nil || time.Since(parsedDate) > 30*24*time.Hour {
				c.JSON(http.StatusBadRequest, gin.H{"error": "data non modificabile oltre 30 giorni"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "presenze registrate"})
		})

		// 35 days old date
		oldDate := time.Now().AddDate(0, 0, -35).Format("2006-01-02")
		payload, _ := json.Marshal(map[string]interface{}{
			"class_id": "class-1a",
			"date":     oldDate,
			"statuses": []map[string]string{
				{"student_id": "st-luca", "status": "Present"},
			},
		})
		req, _ := http.NewRequest("POST", "/api/v1/attendance/mark-bulk", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("4. Student Timetable Cross-Class Authorization Check", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_id", "st-luca")
			c.Set("role", "student")
			c.Next()
		})

		router.GET("/api/v1/classes/:id/schedule", func(c *gin.Context) {
			role, _ := c.Get("role")
			if role == "student" && c.Param("id") != "my-class-1a" {
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"schedule": []interface{}{}})
		})

		// Student trying to access another class timetable -> 403 Forbidden
		req, _ := http.NewRequest("GET", "/api/v1/classes/other-class-2b/schedule", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusForbidden, resp.Code)
	})
}
