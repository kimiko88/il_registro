package grades

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"registro-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service   Service
	analytics AnalyticsService
}

func NewHandler(s Service, a AnalyticsService) *Handler {
	return &Handler{service: s, analytics: a}
}

// RegisterRoutes sets up the routes for the grades module
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	grades := router.Group("/grades")
	{
		grades.POST("", h.AddGrade)
		grades.PATCH("/:id", h.UpdateGrade)
		grades.DELETE("/:id", h.DeleteGrade)

		// Class Tests
		grades.POST("/tests", h.CreateTestWithGrades)
		grades.GET("/tests", h.GetClassTestsList)
		grades.GET("/tests/class/:classID", h.GetUpcomingClassTests)
		grades.DELETE("/tests/:id", h.DeleteClassTest)
		grades.PATCH("/tests/:id", h.UpdateClassTest)

		// Retrieval
		grades.GET("/export", h.Export)
		grades.POST("/bulk-import", h.BulkImport)
		// Weight Config — configurable per-category weights for weighted averages
		grades.GET("/weight-configs", h.ListWeightConfigs)
		grades.PUT("/weight-configs", h.UpsertWeightConfig)
		grades.DELETE("/weight-configs/:id", h.DeleteWeightConfig)
		// Legacy aliases kept for backwards compat
		grades.POST("/weight-config", h.UpsertWeightConfig)
		grades.GET("/weight-config/:subjectID", h.ListWeightConfigs)
		grades.POST("/weights", h.UpsertWeightConfig)
		grades.GET("/weights", h.ListWeightConfigs)

		grades.GET("/student/:studentID", h.GetStudentGrades)
		grades.GET("/student/:studentID/paged", h.GetStudentGradesPaged)
		grades.GET("/class/:classID", h.GetClassGrades)
		grades.GET("/subject/:subjectID", h.GetSubjectGrades)

		// Analytics
		grades.GET("/analytics/student/:studentID/average", h.GetStudentAverage)
		grades.GET("/analytics/class/:classID/average", h.GetClassAverage)
		grades.GET("/analytics/class/:classID/analysis", h.GetClassAnalysis)
		grades.GET("/analytics/subject/:subjectID/analysis", h.GetSubjectAnalysis)
		grades.GET("/analytics/student/:studentID/profile", h.GetStudentProfile)
		grades.GET("/analytics/statistics", h.GetSchoolStatistics)

		// Student Endpoints
		grades.GET("/my-grades", h.GetMyGrades)
		grades.GET("/my-grades/average", h.GetMyAverages)
		grades.GET("/my-grades/trend", h.GetMyTrend)
		grades.GET("/my-grades/semester/:semester", h.GetSemesterReport)
		grades.GET("/my-grades/semester/:semester/pdf", h.DownloadSemesterReportPDF)

		// Parent Endpoints
		grades.GET("/child-grades/:studentID", h.GetChildGrades)
		grades.GET("/child-grades/:studentID/average", h.GetChildGradesAverage)
		grades.GET("/child-grades/:studentID/semester/:semester", h.GetChildSemesterReport)
	}
}

// GetStudentGrades retrieves all grades for a specific student with optional filters.
// actorID and actorRole are forwarded to the service layer so it can enforce
// ownership rules (a student may only read their own grades; a teacher may
// only read grades for their classes; a parent must be a registered guardian).
func (h *Handler) GetStudentGrades(c *gin.Context) {
	studentID := c.Param("studentID")
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "studentID is required"})
		return
	}

	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	filter := h.parseFilter(c)

	if filter.Page > 0 {
		resp, err := h.service.GetStudentGradesPaged(c.Request.Context(), actorID, actorRole, studentID, filter)
		if err != nil {
			if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrNotGuardian) {
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	grades, err := h.service.GetStudentGradesWithFilter(c.Request.Context(), actorID, actorRole, studentID, filter)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrNotGuardian) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if grades == nil {
		grades = []GradeResponse{}
	}

	c.JSON(http.StatusOK, grades)
}

// GetStudentGradesPaged returns a paginated list of grades for a student.
// Query params: page (1-based, default 1), page_size (default 50), plus the
// standard filter params (semester, subject_id, grade_type, published).
func (h *Handler) GetStudentGradesPaged(c *gin.Context) {
	studentID := c.Param("studentID")
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "studentID is required"})
		return
	}

	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")

	filter := h.parseFilter(c)

	resp, err := h.service.GetStudentGradesPaged(c.Request.Context(), actorID, actorRole, studentID, filter)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrNotGuardian) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetClassGrades(c *gin.Context) {
	classID := c.Param("classID")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "classID is required"})
		return
	}

	filter := h.parseFilter(c)
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	resp, err := h.service.GetClassGrades(c.Request.Context(), actorID, actorRole, classID, filter)
	if err != nil {
		logger.Log.Errorf("GetClassGrades error: %v", err)
		if errors.Is(err, ErrUnauthorized) || strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if resp == nil {
		resp = &ClassGradesResponse{ClassID: classID, Students: []StudentGradeSummary{}}
	} else if resp.Students == nil {
		resp.Students = []StudentGradeSummary{}
	}

	c.JSON(http.StatusOK, resp)
}

// GetSubjectGrades retrieves grades for a subject.
// actorID and actorRole are passed to the service so it can verify that
// the caller is a teacher assigned to that subject (or an admin/superadmin).
func (h *Handler) GetSubjectGrades(c *gin.Context) {
	subjectID := c.Param("subjectID")
	if subjectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subjectID is required"})
		return
	}

	filter := h.parseFilter(c)
	actorID := c.GetString("user_id")
	actorRole := c.GetString("role")

	resp, err := h.service.GetSubjectGrades(c.Request.Context(), actorID, actorRole, subjectID, filter)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) || strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if resp == nil {
		resp = &SubjectStatsResponse{Classes: []ClassStat{}}
	} else if resp.Classes == nil {
		resp.Classes = []ClassStat{}
	}

	c.JSON(http.StatusOK, resp)
}

// BulkImport imports grades from an uploaded file.
func (h *Handler) BulkImport(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file upload required"})
		return
	}

	if fileHeader.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds 10MB limit"})
		return
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".csv" && ext != ".xlsx" && ext != ".xls" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file format, only CSV and Excel (.xlsx, .xls) files are supported"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	semester := 1 // default
	semStr := c.PostForm("semester")
	if semStr != "" {
		if semStr != "1" && semStr != "2" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "semester must be 1 or 2"})
			return
		}
		if semStr == "2" {
			semester = 2
		}
	}

	result, err := h.service.BulkImport(teacherID, file, semester)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Export(c *gin.Context) {
	teacherID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	format := strings.ToLower(c.Query("format"))
	if format == "" {
		format = "json"
	}
	if format != "json" && format != "csv" && format != "pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported format: " + format})
		return
	}

	filter := h.parseFilter(c)

	data, contentType, err := h.service.Export(teacherID, schoolID, filter, format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if format == "json" {
		contentType = "application/json"
	}

	filename := fmt.Sprintf("grades_export_%s.%s", time.Now().Format("20060102_150405"), format)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Data(http.StatusOK, contentType, data)
}

// Helper to parse query params into GradeFilter
func (h *Handler) parseFilter(c *gin.Context) GradeFilter {
	var filter GradeFilter
	_ = c.BindQuery(&filter)
	return filter
}

// AddGrade creates a new grade entry
func (h *Handler) AddGrade(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.AddGrade(c.Request.Context(), teacherID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// UpdateGrade modifies an existing grade
func (h *Handler) UpdateGrade(c *gin.Context) {
	gradeID := c.Param("id")
	if gradeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gradeID is required"})
		return
	}

	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req UpdateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.UpdateGrade(c.Request.Context(), teacherID, gradeID, req)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) || strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// DeleteGrade removes a grade (soft delete)
func (h *Handler) DeleteGrade(c *gin.Context) {
	gradeID := c.Param("id")
	if gradeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gradeID is required"})
		return
	}

	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := h.service.DeleteGrade(c.Request.Context(), teacherID, gradeID); err != nil {
		if errors.Is(err, ErrUnauthorized) || strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "grade deleted successfully"})
}

// Analytics Handlers

func (h *Handler) GetStudentAverage(c *gin.Context) {
	userID := c.GetString("user_id")
	actorRole := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	studentID := c.Param("studentID")
	subjectID := c.Query("subject_id")

	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "studentID is required"})
		return
	}

	if actorRole == "student" && userID != studentID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	avg, err := h.analytics.GetStudentAverage(studentID, subjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"student_id": studentID,
		"subject_id": subjectID,
		"average":    avg,
	})
}

func (h *Handler) GetClassAverage(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	classID := c.Param("classID")
	subjectID := c.Query("subject_id")

	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "classID is required"})
		return
	}

	avg, err := h.analytics.GetClassAverage(classID, subjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"class_id":   classID,
		"subject_id": subjectID,
		"average":    avg,
	})
}

// --- Student Endpoints ---

func (h *Handler) GetMyGrades(c *gin.Context) {
	studentID := c.GetString("user_id")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	filter := h.parseFilter(c)

	resp, err := h.service.GetMyGrades(c.Request.Context(), studentID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetMyAverages(c *gin.Context) {
	studentID := c.GetString("user_id")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	resp, err := h.service.GetMyAverages(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetMyTrend(c *gin.Context) {
	studentID := c.GetString("user_id")
	role := c.GetString("role")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subjectID := c.Query("subject_id")

	resp, err := h.service.GetMyTrend(c.Request.Context(), studentID, role, studentID, subjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetSemesterReport(c *gin.Context) {
	studentID := c.GetString("user_id")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	semStr := c.Param("semester")
	sem, err := strconv.Atoi(semStr)
	if err != nil || (sem != 1 && sem != 2) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "semester must be 1 or 2"})
		return
	}

	resp, err := h.service.GetSemesterReport(c.Request.Context(), studentID, sem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DownloadSemesterReportPDF(c *gin.Context) {
	studentID := c.GetString("user_id")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	semStr := c.Param("semester")
	sem, err := strconv.Atoi(semStr)
	if err != nil || (sem != 1 && sem != 2) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "semester must be 1 or 2"})
		return
	}

	pdfData, err := h.service.GenerateSemesterReportPDF(studentID, sem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := fmt.Sprintf("pagella_q%d_%s.pdf", sem, studentID)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/pdf", pdfData)
}

// --- Parent Endpoints ---

func (h *Handler) GetChildGrades(c *gin.Context) {
	parentID := c.GetString("user_id")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	studentID := c.Param("studentID")
	logger.Log.Debug("GetChildGrades requested by parent")
	filter := h.parseFilter(c)

	resp, err := h.service.GetChildGrades(c.Request.Context(), parentID, studentID, filter)
	if err != nil {
		logger.Log.Errorf("GetChildGrades error: %v", err)
		if errors.Is(err, ErrNotGuardian) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetChildGradesAverage(c *gin.Context) {
	parentID := c.GetString("user_id")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	studentID := c.Param("studentID")
	resp, err := h.service.GetChildAverages(c.Request.Context(), parentID, studentID)
	if err != nil {
		if errors.Is(err, ErrNotGuardian) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetChildSemesterReport(c *gin.Context) {
	parentID := c.GetString("user_id")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	studentID := c.Param("studentID")
	semStr := c.Param("semester")
	sem, err := strconv.Atoi(semStr)
	if err != nil || (sem != 1 && sem != 2) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "semester must be 1 or 2"})
		return
	}

	resp, err := h.service.GetChildSemesterReport(c.Request.Context(), parentID, studentID, sem)
	if err != nil {
		if errors.Is(err, ErrNotGuardian) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// --- Analytics Handlers ---

func (h *Handler) GetClassAnalysis(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	classID := c.Param("classID")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "classID is required"})
		return
	}

	semStr := c.Query("semester")
	sem := 1
	if semStr != "" {
		parsed, err := strconv.Atoi(semStr)
		if err != nil || (parsed != 1 && parsed != 2) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "semester must be 1 or 2"})
			return
		}
		sem = parsed
	}

	resp, err := h.analytics.GetClassAnalysis(classID, sem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetSubjectAnalysis(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	subjectID := c.Param("subjectID")
	if subjectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subjectID is required"})
		return
	}

	semStr := c.Query("semester")
	sem := 1
	if semStr != "" {
		parsed, err := strconv.Atoi(semStr)
		if err != nil || (parsed != 1 && parsed != 2) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "semester must be 1 or 2"})
			return
		}
		sem = parsed
	}

	resp, err := h.analytics.GetSubjectAnalysis(subjectID, sem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetStudentProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	actorRole := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	studentID := c.Param("studentID")
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "studentID is required"})
		return
	}

	if actorRole == "student" && userID != studentID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	sem := 0
	semStr := c.Query("semester")
	if semStr != "" {
		parsed, err := strconv.Atoi(semStr)
		if err != nil || parsed < 0 || parsed > 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "semester must be 0, 1 or 2"})
			return
		}
		sem = parsed
	}

	resp, err := h.analytics.GetStudentProfile(studentID, sem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetSchoolStatistics(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	actorRole := c.GetString("role")
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "principal" && actorRole != "secretary" && actorRole != "teacher" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: school statistics accessible only to staff and administrators"})
		return
	}

	year := c.Query("year")
	resp, err := h.analytics.GetSchoolStatistics(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) CreateTestWithGrades(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateClassTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateTestWithGrades(teacherID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "test and grades created successfully"})
}

func (h *Handler) GetClassTestsList(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	classID := c.Query("class_id")
	subjectID := c.Query("subject_id")
	if classID == "" || subjectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class_id and subject_id are required"})
		return
	}

	resp, err := h.service.GetClassTests(classID, subjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if resp == nil {
		resp = []ClassTestResponse{}
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetUpcomingClassTests(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	classID := c.Param("classID")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "classID is required"})
		return
	}

	resp, err := h.service.GetUpcomingTestsByClass(classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if resp == nil {
		resp = []ClassTestResponse{}
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteClassTest(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	testID := c.Param("id")
	if testID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "test ID is required"})
		return
	}

	if err := h.service.DeleteClassTest(teacherID, testID); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: not the author of this test"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "test and linked grades deleted successfully"})
}

func (h *Handler) UpdateClassTest(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	testID := c.Param("id")
	if testID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "test ID is required"})
		return
	}

	var req UpdateClassTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateClassTest(teacherID, testID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "test and linked grades updated successfully"})
}

// ListWeightConfigs returns grade weight configurations for a school/subject/class.
func (h *Handler) ListWeightConfigs(c *gin.Context) {
	if c.GetString("user_id") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	schoolID := c.GetString("school_id")
	subjectID := c.Query("subject_id")
	classID := c.Query("class_id")

	configs, err := h.service.GetWeightConfigs(schoolID, subjectID, classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, configs)
}

// UpsertWeightConfig creates or updates a weight config entry.
func (h *Handler) UpsertWeightConfig(c *gin.Context) {
	actorID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	role := c.GetString("role")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req UpsertWeightConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.UpsertWeightConfig(actorID, role, schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// DeleteWeightConfig removes a weight config entry by ID.
func (h *Handler) DeleteWeightConfig(c *gin.Context) {
	actorID := c.GetString("user_id")
	schoolID := c.GetString("school_id")
	role := c.GetString("role")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	configID := c.Param("id")
	if err := h.service.DeleteWeightConfig(actorID, role, schoolID, configID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// Deprecated: SetWeightConfig is a legacy alias kept for backwards compatibility.
func (h *Handler) SetWeightConfig(c *gin.Context) {
	h.UpsertWeightConfig(c)
}

// Deprecated: GetWeightConfig is a legacy alias kept for backwards compatibility.
func (h *Handler) GetWeightConfig(c *gin.Context) {
	h.ListWeightConfigs(c)
}
