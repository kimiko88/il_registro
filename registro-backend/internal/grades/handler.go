package grades

import (
	"net/http"

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

		// Retrieval
		grades.GET("/export", h.Export) // Export all/filtered
		grades.POST("/bulk-import", h.BulkImport)

		// Retrieval
		grades.GET("/student/:studentID", h.GetStudentGrades)
		grades.GET("/class/:classID", h.GetClassGrades)
		grades.GET("/subject/:subjectID", h.GetSubjectGrades)

		// Analytics
		grades.GET("/analytics/student/:studentID/average", h.GetStudentAverage)
		grades.GET("/analytics/class/:classID/average", h.GetClassAverage)

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

// GetStudentGrades retrieves all grades for a specific student
// GetStudentGrades retrieves all grades for a specific student with optional filters
func (h *Handler) GetStudentGrades(c *gin.Context) {
	studentID := c.Param("studentID")
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "studentID is required"})
		return
	}

	filter := h.parseFilter(c)

	grades, err := h.service.GetStudentGradesWithFilter(studentID, filter)
	if err != nil {
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
	actorID := c.GetString("userID")
	actorRole := c.GetString("role")

	resp, err := h.service.GetClassGrades(c.Request.Context(), actorID, actorRole, classID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetSubjectGrades(c *gin.Context) {
	subjectID := c.Param("subjectID")
	if subjectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subjectID is required"})
		return
	}

	filter := h.parseFilter(c)

	resp, err := h.service.GetSubjectGrades(subjectID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) BulkImport(c *gin.Context) {
	teacherID := c.GetString("userID")
	if teacherID == "" {
		teacherID = "dev-teacher-id"
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file upload required"})
		return
	}
	defer file.Close()

	semester := 1 // Default
	if s := c.PostForm("semester"); s != "" {
		// handle parsing or bind
		// simplicity:
		// ignore error for MVP or use 1
	}

	result, err := h.service.BulkImport(teacherID, file, semester)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Export(c *gin.Context) {
	teacherID := c.GetString("userID")
	if teacherID == "" {
		teacherID = "dev-teacher-id"
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

	// BindQuery requires struct tags `form:"name"` which we added to DTO
	if err := c.BindQuery(&filter); err != nil {
		// Log error or ignore
	}

	return filter
}

// AddGrade creates a new grade entry
func (h *Handler) AddGrade(c *gin.Context) {
	var req CreateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	teacherID := c.GetString("userID")
	if teacherID == "" {
		// For testing purpose if auth middleware missing, checking header or mocking
		// In prod this is fatal or handled by middleware
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		// return
		// ALLOW PASS for development if local
		teacherID = "dev-teacher-id"
	}

	if err := h.service.AddGrade(teacherID, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Changed to 400 as validations are common
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

	teacherID := c.GetString("userID")
	if teacherID == "" {
		teacherID = "dev-teacher-id"
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

	teacherID := c.GetString("userID")
	if teacherID == "" {
		teacherID = "dev-teacher-id"
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
	subjectID := c.Query("subject_id") // Optional filter

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
	subjectID := c.Query("subject_id") // Optional

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
	studentID := c.GetString("userID")
	if studentID == "" {
		studentID = "dev-student-id"
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
	studentID := c.GetString("userID")
	if studentID == "" {
		studentID = "dev-student-id"
	}

	resp, err := h.service.GetMyAverages(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetMyTrend(c *gin.Context) {
	studentID := c.GetString("userID")
	if studentID == "" {
		studentID = "dev-student-id"
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
	studentID := c.GetString("userID")
	if studentID == "" {
		studentID = "dev-student-id"
	}

	semStr := c.Param("semester")
	sem := 1
	if semStr == "2" {
		sem = 2
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
	parentID := c.GetString("userID")
	if parentID == "" {
		parentID = "dev-parent-id"
	}

	studentID := c.Param("studentID")
	filter := h.parseFilter(c)

	resp, err := h.service.GetChildGrades(parentID, studentID, filter)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetChildGradesAverage(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "GetChildAverages not fully wired"})
}

// --- Analytics Handlers ---

func (h *Handler) GetClassAnalysis(c *gin.Context) {
	classID := c.Param("classID")
	semStr := c.Query("semester")
	sem := 1
	if semStr == "2" {
		sem = 2
	}

	resp, err := h.analytics.GetClassAnalysis(classID, sem)
	// Note: Handler sees Service interface. `analyticsService` is separate struct in current mapping.
	// `Handler` struct has `analytics` field of type `AnalyticsService`.
	// I should call h.analytics.GetClassAnalysis

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetSubjectAnalysis(c *gin.Context) {
	subjectID := c.Param("subjectID")
	semStr := c.Query("semester")
	sem := 1
	if semStr == "2" {
		sem = 2
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
	semStr := c.Query("semester")
	sem := 0
	if semStr == "1" {
		sem = 1
	}
	if semStr == "2" {
		sem = 2
	}

	resp, err := h.analytics.GetStudentProfile(studentID, sem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetSchoolStatistics(c *gin.Context) {
	year := c.Query("year") // Optional
	resp, err := h.analytics.GetSchoolStatistics(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
