package grades

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

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

		grades.GET("/student/:studentID", h.GetStudentGrades)
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

		// Parent Endpoints
		grades.GET("/child-grades/:studentID", h.GetChildGrades)
		grades.GET("/child-grades/:studentID/average", h.GetChildGradesAverage)
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

	filter := h.parseFilter(c)

	grades, err := h.service.GetStudentGradesWithFilter(c.Request.Context(), actorID, actorRole, studentID, filter)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrNotGuardian) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, grades)
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

	resp, err := h.service.GetClassGrades(c.Request.Context(), actorID, actorRole, classID, filter)
	if err != nil {
		logger.Log.Errorf("GetClassGrades error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

	c.JSON(http.StatusOK, resp)
}

// BulkImport imports grades from an uploaded file.
// The semester parameter (1 or 2) is read from the multipart form; if absent
// it defaults to 1. Previously the form value was read but silently discarded
// with `_ = c.PostForm(...)`, causing all imports to land in semester 1.
func (h *Handler) BulkImport(c *gin.Context) {
	teacherID := c.GetString("user_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file upload required"})
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
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	format := c.Query("format")
	if format == "" {
		format = "json"
	}

	filter := h.parseFilter(c)

	data, contentType, err := h.service.Export(teacherID, filter, format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if format == "json" {
		contentType = "application/json"
	}

	c.Header("Content-Disposition", "attachment; filename=grades."+format)
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
	var req CreateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacherID := c.GetString("user_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.AddGrade(teacherID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "grade created successfully"})
}

// UpdateGrade modifies an existing grade
func (h *Handler) UpdateGrade(c *gin.Context) {
	gradeID := c.Param("id")
	if gradeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gradeID is required"})
		return
	}

	var req UpdateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacherID := c.GetString("user_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.UpdateGrade(teacherID, gradeID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "grade updated successfully"})
}

// DeleteGrade removes a grade (soft delete)
func (h *Handler) DeleteGrade(c *gin.Context) {
	gradeID := c.Param("id")
	if gradeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gradeID is required"})
		return
	}

	teacherID := c.GetString("user_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.DeleteGrade(teacherID, gradeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "grade deleted successfully"})
}

// Analytics Handlers

func (h *Handler) GetStudentAverage(c *gin.Context) {
	studentID := c.Param("studentID")
	subjectID := c.Query("subject_id")

	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "studentID is required"})
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

	resp, err := h.service.GetMyGrades(studentID, filter)
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

	resp, err := h.service.GetMyAverages(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetMyTrend(c *gin.Context) {
	studentID := c.GetString("user_id")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subjectID := c.Query("subject_id")

	resp, err := h.service.GetMyTrend(studentID, subjectID)
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

	resp, err := h.service.GetSemesterReport(studentID, sem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
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

	resp, err := h.service.GetChildGrades(parentID, studentID, filter)
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
	resp, err := h.service.GetChildAverages(parentID, studentID)
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
	studentID := c.Param("studentID")
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "studentID is required"})
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
	var req CreateClassTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacherID := c.GetString("user_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.CreateTestWithGrades(teacherID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "test and grades created successfully"})
}

func (h *Handler) GetClassTestsList(c *gin.Context) {
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

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetUpcomingClassTests(c *gin.Context) {
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
	testID := c.Param("id")
	if testID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "test ID is required"})
		return
	}

	teacherID := c.GetString("user_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
	testID := c.Param("id")
	if testID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "test ID is required"})
		return
	}

	teacherID := c.GetString("user_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
