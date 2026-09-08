package students

import (
	"context"
	"database/sql"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StudentDashboardStatsResponse struct {
	AverageGrade   float64 `json:"average_grade"`
	AttendanceRate float64 `json:"attendance_rate"`
	HomeworkCount  int     `json:"homework_count"`
	DocumentsCount int     `json:"documents_count"`
	TotalGrades    int     `json:"total_grades"`
	PresenceRate   float64 `json:"presence_rate"`
	UpcomingTests  int     `json:"upcoming_tests"`
}

type DashboardService struct {
	db *sql.DB
}

func NewDashboardService(db *sql.DB) *DashboardService {
	return &DashboardService{db: db}
}

func (s *DashboardService) GetStats(ctx context.Context, userID string) (*StudentDashboardStatsResponse, error) {
	resp := &StudentDashboardStatsResponse{
		AverageGrade:   0,
		AttendanceRate: 100.0,
		HomeworkCount:  0,
		DocumentsCount: 0,
		TotalGrades:    0,
		PresenceRate:   100.0,
		UpcomingTests:  0,
	}

	if s.db == nil || userID == "" {
		return resp, nil
	}

	// 1. Identify student profile ID and class ID
	var studentProfileID, classID string
	queryProfile := `
		SELECT id::text, COALESCE(class_id::text, '') 
		FROM students 
		WHERE (user_id = $1::uuid OR id = $1::uuid) AND deleted_at IS NULL 
		LIMIT 1`
	_ = s.db.QueryRowContext(ctx, queryProfile, userID).Scan(&studentProfileID, &classID)
	if studentProfileID == "" {
		studentProfileID = userID
	}

	// 2. Calculate average grade and total grades count
	var avgGrade sql.NullFloat64
	var totalGrades int
	queryGrades := `
		SELECT AVG(grade_value), COUNT(*) 
		FROM grades 
		WHERE (student_id = $1::uuid OR student_id = $2::uuid)
		  AND is_published = TRUE 
		  AND deleted_at IS NULL`
	if err := s.db.QueryRowContext(ctx, queryGrades, userID, studentProfileID).Scan(&avgGrade, &totalGrades); err == nil {
		if avgGrade.Valid && totalGrades > 0 {
			resp.AverageGrade = math.Round(avgGrade.Float64*100) / 100
			resp.TotalGrades = totalGrades
		}
	}

	// 3. Calculate attendance / presence rate
	var presentCount, totalAttCount int
	queryAtt := `
		SELECT 
			COUNT(*) FILTER (WHERE status = 'present') AS present_count,
			COUNT(*) AS total_count
		FROM attendance 
		WHERE (student_id = $1::uuid OR student_id = $2::uuid)
		  AND deleted_at IS NULL`
	if err := s.db.QueryRowContext(ctx, queryAtt, userID, studentProfileID).Scan(&presentCount, &totalAttCount); err == nil {
		if totalAttCount > 0 {
			rate := (float64(presentCount) / float64(totalAttCount)) * 100.0
			resp.AttendanceRate = math.Round(rate*10) / 10
			resp.PresenceRate = resp.AttendanceRate
		}
	}

	// 4. Calculate homework and upcoming tests if class is known
	if classID != "" {
		var homeworkCount int
		queryHw := `
			SELECT COUNT(*) 
			FROM homeworks h
			JOIN class_lessons l ON h.lesson_id = l.id
			WHERE l.class_id = $1::uuid 
			  AND h.due_date >= CURRENT_DATE 
			  AND h.deleted_at IS NULL`
		if err := s.db.QueryRowContext(ctx, queryHw, classID).Scan(&homeworkCount); err == nil {
			resp.HomeworkCount = homeworkCount
		}

		var upcomingTests int
		queryTests := `
			SELECT COUNT(*) 
			FROM class_tests 
			WHERE class_id = $1::uuid 
			  AND test_date >= CURRENT_DATE 
			  AND deleted_at IS NULL`
		if err := s.db.QueryRowContext(ctx, queryTests, classID).Scan(&upcomingTests); err == nil {
			resp.UpcomingTests = upcomingTests
		}
	}

	// 5. Calculate documents count
	var docsCount int
	queryDocs := `
		SELECT COUNT(*) 
		FROM documents_enhanced 
		WHERE ((student_id = $1::uuid OR student_id = $2::uuid) 
		       OR (class_id IS NOT NULL AND class_id = NULLIF($3, '')::uuid))
		  AND status = 'published' 
		  AND deleted_at IS NULL`
	if err := s.db.QueryRowContext(ctx, queryDocs, userID, studentProfileID, classID).Scan(&docsCount); err == nil {
		resp.DocumentsCount = docsCount
	}

	return resp, nil
}

type DashboardHandler struct {
	svc *DashboardService
}

func NewDashboardHandler(svc *DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/students/dashboard/stats", h.GetStats)
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	stats, err := h.svc.GetStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
