package grades

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"registro-backend/internal/pdfworker"
	"registro-backend/pkg/logger"
	"registro-backend/pkg/upload"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var filenameParamRegex = regexp.MustCompile(`[^a-zA-Z0-9_-]`)

func sanitizeFilenameParam(input string) string {
	res := filenameParamRegex.ReplaceAllString(input, "")
	if res == "" {
		return "export"
	}
	return res
}

func respond500(c *gin.Context, msg string, err error) {
	logger.Log.Errorf("%s: %v", msg, err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}

type Handler struct {
	service         Service
	analytics       AnalyticsService
	validator       *Validator
	pdfWorkerClient *pdfworker.Client
}

func NewHandler(s Service, a AnalyticsService, v ...*Validator) *Handler {
	h := &Handler{service: s, analytics: a}
	if len(v) > 0 {
		h.validator = v[0]
	}
	return h
}

func (h *Handler) SetPdfWorkerClient(client *pdfworker.Client) {
	h.pdfWorkerClient = client
}

func (h *Handler) legacyWeightConfig(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Deprecation", "true")
		c.Header("Link", `</api/v1/grades/weight-configs>; rel="successor-version"`)
		handler(c)
	}
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
		grades.POST("/export/async-pdf", h.EnqueueAsyncRegisterPdf)
		grades.GET("/pdf-jobs/:job_id", h.GetPdfJobStatus)
		grades.POST("/bulk-import", h.BulkImport)
		// Weight Config — configurable per-category weights for weighted averages
		grades.GET("/weight-configs", h.ListWeightConfigs)
		grades.PUT("/weight-configs", h.UpsertWeightConfig)
		grades.DELETE("/weight-configs/:id", h.DeleteWeightConfig)
		// Legacy aliases kept for backwards compat with Deprecation header
		grades.POST("/weight-config", h.legacyWeightConfig(h.UpsertWeightConfig))
		grades.GET("/weight-config/:subjectID", h.legacyWeightConfig(h.ListWeightConfigs))
		grades.POST("/weights", h.legacyWeightConfig(h.UpsertWeightConfig))
		grades.GET("/weights", h.legacyWeightConfig(h.ListWeightConfigs))

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

	c.Header("Deprecation", "true")
	c.Header("Link", fmt.Sprintf(`</api/v1/grades/student/%s/paged>; rel="successor-version"`, studentID))

	filter := h.parseFilter(c)

	grades, err := h.service.GetStudentGradesWithFilter(c.Request.Context(), actorID, actorRole, studentID, filter)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrNotGuardian) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "GetStudentGradesWithFilter error", err)
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
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	filter := h.parseFilter(c)

	resp, err := h.service.GetStudentGradesPaged(c.Request.Context(), actorID, actorRole, studentID, filter)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrNotGuardian) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "GetStudentGradesPaged error", err)
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
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "GetClassGrades error", err)
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
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	resp, err := h.service.GetSubjectGrades(c.Request.Context(), actorID, actorRole, subjectID, filter)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "GetSubjectGrades error", err)
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
	defer func() { _ = file.Close() }()

	if err := upload.ValidateUpload(file, fileHeader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			respond500(c, "file seek error", err)
			return
		}
	}

	// In multipart/form-data, the form field takes explicit precedence over query parameter
	semester := 1 // default
	semStr := c.PostForm("semester")
	if semStr == "" {
		semStr = c.Query("semester")
	}
	if semStr != "" {
		if semStr != "1" && semStr != "2" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "semester must be 1 or 2"})
			return
		}
		if semStr == "2" {
			semester = 2
		}
	}

	schoolID := c.GetString("school_id")
	result, err := h.service.BulkImport(c.Request.Context(), teacherID, schoolID, file, semester)
	if err != nil {
		respond500(c, "BulkImport error", err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Export(c *gin.Context) {
	teacherID := c.GetString("user_id")
	role := c.GetString("role")
	schoolID := c.GetString("school_id")
	if teacherID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: export reserved to staff"})
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

	if format == "pdf" && h.pdfWorkerClient != nil && c.Query("sync") != "true" {
		h.EnqueueAsyncRegisterPdf(c)
		return
	}

	filter := h.parseFilter(c)

	data, contentType, err := h.service.Export(c.Request.Context(), teacherID, schoolID, filter, format)
	if err != nil {
		respond500(c, "Export error", err)
		return
	}

	if format == "json" || contentType == "" {
		switch format {
		case "json":
			contentType = "application/json"
		case "csv":
			contentType = "text/csv"
		case "pdf":
			contentType = "application/pdf"
		}
	}

	filename := fmt.Sprintf("grades_export_%s.%s", time.Now().Format("20060102_150405"), format)
	c.Header("Content-Disposition", upload.FormatContentDisposition(filename))
	c.Data(http.StatusOK, contentType, data)
}

// Helper to parse query params into GradeFilter
func (h *Handler) parseFilter(c *gin.Context) GradeFilter {
	var filter GradeFilter
	if err := c.BindQuery(&filter); err != nil {
		logger.Log.Warnf("parseFilter BindQuery error: %v", err)
	}
	if filter.Semester < 0 || filter.Semester > 2 {
		filter.Semester = 0
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
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
		if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrNotGuardian) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
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
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "UpdateGrade error", err)
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
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "DeleteGrade error", err)
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
	if actorRole == "teacher" {
		if h.validator == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: validator unconfigured"})
			return
		}
		assigned, err := h.validator.IsTeacherAssignedToStudent(c.Request.Context(), userID, studentID)
		if err != nil || !assigned {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: teacher is not assigned to this student"})
			return
		}
	}
	if actorRole == "parent" {
		if err := h.service.ValidateParentGuardian(c.Request.Context(), userID, studentID); err != nil {
			if errors.Is(err, ErrNotGuardian) || strings.Contains(err.Error(), "not a guardian") {
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: not a guardian of this student"})
			} else {
				respond500(c, "GetStudentAverage ValidateParentGuardian error", err)
			}
			return
		}
	}

	avg, err := h.analytics.GetStudentAverage(studentID, subjectID)
	if err != nil {
		respond500(c, "GetStudentAverage error", err)
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
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
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
		respond500(c, "GetClassAverage error", err)
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
	role := c.GetString("role")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: endpoint reserved for students"})
		return
	}

	filter := h.parseFilter(c)

	resp, err := h.service.GetMyGrades(c.Request.Context(), studentID, filter)
	if err != nil {
		respond500(c, "GetMyGrades error", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetMyAverages(c *gin.Context) {
	studentID := c.GetString("user_id")
	role := c.GetString("role")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: solo uno studente può accedere alle proprie medie"})
		return
	}

	resp, err := h.service.GetMyAverages(c.Request.Context(), studentID)
	if err != nil {
		respond500(c, "GetMyAverages error", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetMyTrend(c *gin.Context) {
	actorID := c.GetString("user_id")
	role := c.GetString("role")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	targetStudentID := actorID
	if (role == "admin" || role == "superadmin" || role == "secretary" || role == "principal" || role == "vice_principal" || role == "teacher" || role == "parent") && c.Query("student_id") != "" {
		targetStudentID = c.Query("student_id")
	}

	if role == "parent" && targetStudentID != actorID {
		if err := h.service.ValidateParentGuardian(c.Request.Context(), actorID, targetStudentID); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: not a guardian of this student"})
			return
		}
	}

	if role == "teacher" && targetStudentID != actorID {
		if h.validator == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: validator unconfigured"})
			return
		}
		assigned, err := h.validator.IsTeacherAssignedToStudent(c.Request.Context(), actorID, targetStudentID)
		if err != nil || !assigned {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: teacher is not assigned to this student"})
			return
		}
	}

	subjectID := c.Query("subject_id")

	resp, err := h.service.GetMyTrend(c.Request.Context(), actorID, role, targetStudentID, subjectID)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrNotGuardian) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "GetMyTrend error", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetSemesterReport(c *gin.Context) {
	actorID := c.GetString("user_id")
	role := c.GetString("role")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	targetStudentID := actorID
	if (role == "admin" || role == "superadmin" || role == "secretary" || role == "principal" || role == "vice_principal" || role == "teacher" || role == "parent") && c.Query("student_id") != "" {
		targetStudentID = c.Query("student_id")
	}

	if role == "parent" && targetStudentID != actorID {
		if err := h.service.ValidateParentGuardian(c.Request.Context(), actorID, targetStudentID); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: not a guardian of this student"})
			return
		}
	}

	if role == "teacher" && targetStudentID != actorID {
		if h.validator == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: validator unconfigured"})
			return
		}
		assigned, err := h.validator.IsTeacherAssignedToStudent(c.Request.Context(), actorID, targetStudentID)
		if err != nil || !assigned {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: teacher is not assigned to this student"})
			return
		}
	}

	semStr := c.Param("semester")
	sem, err := strconv.Atoi(semStr)
	if err != nil || (sem != 1 && sem != 2) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "semester must be 1 or 2"})
		return
	}

	resp, err := h.service.GetSemesterReport(c.Request.Context(), actorID, role, targetStudentID, sem)
	if err != nil {
		if errors.Is(err, ErrNotGuardian) || errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "GetSemesterReport error", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DownloadSemesterReportPDF(c *gin.Context) {
	if h.pdfWorkerClient != nil && c.Query("sync") != "true" {
		actorID := c.GetString("user_id")
		targetStudentID := actorID
		if c.Query("student_id") != "" {
			targetStudentID = c.Query("student_id")
		}
		semStr := c.Param("semester")
		jobID := uuid.New().String()
		payload := pdfworker.ReportCardPdfPayload{
			JobID:       jobID,
			StudentID:   targetStudentID,
			Period:      semStr,
			RequestedBy: actorID,
		}
		jobStatus, err := h.pdfWorkerClient.EnqueueReportCardPdf(c.Request.Context(), payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to enqueue async pdf job: %v", err)})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{
			"message":    "PDF generation job enqueued successfully",
			"job_id":     jobStatus.JobID,
			"status":     jobStatus.Status,
			"status_url": fmt.Sprintf("/api/v1/grades/pdf-jobs/%s", jobStatus.JobID),
		})
		return
	}

	actorID := c.GetString("user_id")
	role := c.GetString("role")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "student" && role != "parent" && role != "admin" && role != "superadmin" && role != "secretary" && role != "principal" && role != "vice_principal" && role != "teacher" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	targetStudentID := actorID
	if (role == "admin" || role == "superadmin" || role == "secretary" || role == "principal" || role == "vice_principal" || role == "teacher" || role == "parent") && c.Query("student_id") != "" {
		targetStudentID = c.Query("student_id")
	}

	if role == "parent" && targetStudentID != actorID {
		if err := h.service.ValidateParentGuardian(c.Request.Context(), actorID, targetStudentID); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: not a guardian of this student"})
			return
		}
	}

	if role == "teacher" && targetStudentID != actorID {
		if h.validator == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: validator unconfigured"})
			return
		}
		assigned, err := h.validator.IsTeacherAssignedToStudent(c.Request.Context(), actorID, targetStudentID)
		if err != nil || !assigned {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: teacher is not assigned to this student"})
			return
		}
	}

	semStr := c.Param("semester")
	sem, err := strconv.Atoi(semStr)
	if err != nil || (sem != 1 && sem != 2) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "semester must be 1 or 2"})
		return
	}

	pdfData, err := h.service.GenerateSemesterReportPDF(c.Request.Context(), actorID, role, targetStudentID, sem)
	if err != nil {
		respond500(c, "GenerateSemesterReportPDF error", err)
		return
	}

	filename := fmt.Sprintf("pagella_q%d_%s.pdf", sem, sanitizeFilenameParam(targetStudentID))
	c.Header("Content-Disposition", upload.FormatContentDisposition(filename))
	c.Data(http.StatusOK, "application/pdf", pdfData)
}

// --- Parent Endpoints ---

func (h *Handler) GetChildGrades(c *gin.Context) {
	parentID := c.GetString("user_id")
	role := c.GetString("role")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "parent" && role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: endpoint reserved for parents"})
		return
	}

	studentID := c.Param("studentID")
	logger.Log.Debug("GetChildGrades requested by parent")
	filter := h.parseFilter(c)

	resp, err := h.service.GetChildGrades(c.Request.Context(), parentID, studentID, filter)
	if err != nil {
		if errors.Is(err, ErrNotGuardian) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "GetChildGrades error", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetChildGradesAverage(c *gin.Context) {
	parentID := c.GetString("user_id")
	role := c.GetString("role")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "parent" && role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: endpoint reserved for parents"})
		return
	}

	studentID := c.Param("studentID")
	resp, err := h.service.GetChildAverages(c.Request.Context(), parentID, studentID)
	if err != nil {
		if errors.Is(err, ErrNotGuardian) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "GetChildAverages error", err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetChildSemesterReport(c *gin.Context) {
	parentID := c.GetString("user_id")
	role := c.GetString("role")
	if parentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "parent" && role != "admin" && role != "superadmin" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: endpoint reserved for parents"})
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
		respond500(c, "GetChildSemesterReport error", err)
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
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" && role != "secretary" {
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
		respond500(c, "GetClassAnalysis error", err)
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
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" && role != "secretary" {
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
		respond500(c, "GetSubjectAnalysis error", err)
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

	if actorRole == "parent" {
		// Verify guardian relationship before exposing analytical student profile
		if err := h.service.ValidateParentGuardian(c.Request.Context(), userID, studentID); err != nil {
			if errors.Is(err, ErrNotGuardian) || strings.Contains(err.Error(), "guardian") {
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: not a guardian of this student"})
				return
			}
			respond500(c, "GetStudentProfile ValidateParentGuardian error", err)
			return
		}
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
		respond500(c, "GetStudentProfile error", err)
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
	if actorRole != "admin" && actorRole != "superadmin" && actorRole != "principal" && actorRole != "vice_principal" && actorRole != "secretary" && actorRole != "teacher" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: school statistics accessible only to staff and administrators"})
		return
	}

	schoolID := c.GetString("school_id")
	year := c.Query("year")
	resp, err := h.analytics.GetSchoolStatistics(year, schoolID)
	if err != nil {
		respond500(c, "GetSchoolStatistics error", err)
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

	test, err := h.service.CreateTestWithGrades(c.Request.Context(), teacherID, req)
	if err != nil {
		respond500(c, "CreateTestWithGrades error", err)
		return
	}

	testID := ""
	if test != nil {
		testID = test.ID
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "test and grades created successfully",
		"test_id": testID,
	})
}

func (h *Handler) GetClassTestsList(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
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

	resp, err := h.service.GetClassTests(c.Request.Context(), userID, role, classID, subjectID)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "GetClassTestsList error", err)
		return
	}
	if resp == nil {
		resp = []ClassTestResponse{}
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetUpcomingClassTests(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	classID := c.Param("classID")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "classID is required"})
		return
	}

	resp, err := h.service.GetUpcomingTestsByClass(c.Request.Context(), userID, role, classID)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		respond500(c, "GetUpcomingClassTests error", err)
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

	if err := h.service.DeleteClassTest(c.Request.Context(), teacherID, testID); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: not the author of this test"})
			return
		}
		respond500(c, "DeleteClassTest error", err)
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

	if err := h.service.UpdateClassTest(c.Request.Context(), teacherID, testID, req); err != nil {
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized: not the author of this test"})
			return
		}
		respond500(c, "UpdateClassTest error", err)
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

	configs, err := h.service.GetWeightConfigs(c.Request.Context(), schoolID, subjectID, classID)
	if err != nil {
		respond500(c, "ListWeightConfigs error", err)
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
	if role != "admin" && role != "superadmin" && role != "secretary" {
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
		if errors.Is(err, ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
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
	if role != "admin" && role != "superadmin" && role != "secretary" {
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

func (h *Handler) EnqueueAsyncRegisterPdf(c *gin.Context) {
	actorID := c.GetString("user_id")
	role := c.GetString("role")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != "teacher" && role != "admin" && role != "superadmin" && role != "principal" && role != "vice_principal" && role != "secretary" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: export reserved to staff"})
		return
	}

	classID := c.Query("class_id")
	subjectID := c.Query("subject_id")
	semester := c.DefaultQuery("semester", "1")

	if h.pdfWorkerClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "async pdf worker queue not configured"})
		return
	}

	jobID := uuid.New().String()
	payload := pdfworker.RegisterPdfPayload{
		JobID:       jobID,
		ClassID:     classID,
		SubjectID:   subjectID,
		Period:      semester,
		RequestedBy: actorID,
	}

	jobStatus, err := h.pdfWorkerClient.EnqueueRegisterPdf(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to enqueue async register pdf job: %v", err)})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":    "Register PDF generation job enqueued successfully",
		"job_id":     jobStatus.JobID,
		"status":     jobStatus.Status,
		"status_url": fmt.Sprintf("/api/v1/grades/pdf-jobs/%s", jobStatus.JobID),
	})
}

func (h *Handler) GetPdfJobStatus(c *gin.Context) {
	actorID := c.GetString("user_id")
	if actorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	jobID := c.Param("job_id")
	if h.pdfWorkerClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "async pdf worker queue not configured"})
		return
	}

	status, err := h.pdfWorkerClient.GetJobStatus(c.Request.Context(), jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found or expired"})
		return
	}

	c.JSON(http.StatusOK, status)
}
