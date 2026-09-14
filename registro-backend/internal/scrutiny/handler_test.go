package scrutiny

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupScrutinyRouter creates a Gin test router wired to the handler,
// injecting the given identity claims via middleware.
func setupScrutinyRouter(h *Handler, userID, role, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Set("role", role)
		c.Set("school_id", schoolID)
		c.Next()
	})
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

// newScrutinyHandler returns a Handler whose Service has nil repos.
// Only handler-level auth guards (before any service call) are safe to invoke.
func newScrutinyHandler() *Handler {
	svc := NewService(nil, nil, nil, nil, nil)
	return NewHandler(svc, nil)
}

// ─── GetMatrix ────────────────────────────────────────────────────────────────
// Auth check lives in the service; only the 401 guard is at handler level.

func TestScrutinyHandler_GetMatrix_Unauthorized(t *testing.T) {
	h := newScrutinyHandler()
	// Empty user_id → 401 before reaching service
	r := setupScrutinyRouter(h, "", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/scrutiny/matrix/class-1?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ─── Save ────────────────────────────────────────────────────────────────────
// Auth check lives in the service; only handler-level guards are safe.

func TestScrutinyHandler_Save_Unauthorized(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "", "teacher", "school-1")
	body, _ := json.Marshal(SaveScrutinyRequest{StudentID: "s-1", ClassID: "c-1", Semester: 1})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/save", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestScrutinyHandler_Save_BadRequest_MissingFields(t *testing.T) {
	h := newScrutinyHandler()
	// Provide a user_id so we pass the 401 guard. Then binding fails → 400.
	r := setupScrutinyRouter(h, "teacher-1", "teacher", "school-1")
	body := []byte(`{}`) // student_id and class_id both missing (binding:"required" tags)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/save", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── Start / Validate / Close ─────────────────────────────────────────────────

func TestScrutinyHandler_Start_Unauthorized(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/class/c-1/start?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestScrutinyHandler_Validate_Unauthorized(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "", "principal", "school-1")
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/class/c-1/validate?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestScrutinyHandler_Close_Unauthorized(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "", "principal", "school-1")
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/class/c-1/close?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ─── GetOverview (handler-level RBAC check) ───────────────────────────────────

func TestScrutinyHandler_GetOverview_Unauthorized(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "", "principal", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/scrutiny/overview?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestScrutinyHandler_GetOverview_Forbidden_Teacher(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "t-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/scrutiny/overview?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestScrutinyHandler_GetOverview_Forbidden_Student(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "s-1", "student", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/scrutiny/overview?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ─── GetClassReport (handler-level RBAC check) ────────────────────────────────

func TestScrutinyHandler_GetClassReport_Unauthorized(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/scrutiny/class/c-1/report?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestScrutinyHandler_GetClassReport_Forbidden_Student(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "s-1", "student", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/scrutiny/class/c-1/report?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestScrutinyHandler_GetClassReport_Forbidden_Parent(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "p-1", "parent", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/scrutiny/class/c-1/report?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ─── FinalizeClass (handler-level RBAC check) ─────────────────────────────────

func TestScrutinyHandler_FinalizeClass_Unauthorized(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/class/c-1/finalize?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestScrutinyHandler_FinalizeClass_Forbidden_Teacher(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "t-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/class/c-1/finalize?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestScrutinyHandler_FinalizeClass_Forbidden_Student(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "s-1", "student", "school-1")
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/class/c-1/finalize?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ─── ExportAll (handler-level RBAC check) ─────────────────────────────────────

func TestScrutinyHandler_ExportAll_Unauthorized(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "", "principal", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/scrutiny/export?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestScrutinyHandler_ExportAll_Forbidden_Teacher(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "t-1", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/scrutiny/export?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ─── SaveDeficiency (handler-level RBAC check) ────────────────────────────────

func TestScrutinyHandler_SaveDeficiency_Unauthorized(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "", "admin", "school-1")
	body, _ := json.Marshal(SaveDeficiencyRequest{StudentID: "s-1", ClassID: "c-1", SubjectID: "sub-1", Semester: 1})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/deficiencies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestScrutinyHandler_SaveDeficiency_Forbidden_Student(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "s-1", "student", "school-1")
	body, _ := json.Marshal(SaveDeficiencyRequest{StudentID: "s-1", ClassID: "c-1", SubjectID: "sub-1", Semester: 1})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/deficiencies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestScrutinyHandler_SaveDeficiency_Forbidden_Parent(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "p-1", "parent", "school-1")
	body, _ := json.Marshal(SaveDeficiencyRequest{StudentID: "s-1", ClassID: "c-1", SubjectID: "sub-1", Semester: 1})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/deficiencies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ─── GetClassDeficiencies (handler-level RBAC check) ─────────────────────────

func TestScrutinyHandler_GetClassDeficiencies_Unauthorized(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/scrutiny/deficiencies/class/c-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestScrutinyHandler_GetClassDeficiencies_Forbidden_Student(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "s-1", "student", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/scrutiny/deficiencies/class/c-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ─── GetStudentDeficiencies ────────────────────────────────────────────────────

func TestScrutinyHandler_GetStudentDeficiencies_Unauthorized(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "", "teacher", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/scrutiny/deficiencies/student/s-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ─── SaveDeferredScrutiny ─────────────────────────────────────────────────────

func TestScrutinyHandler_SaveDeferredScrutiny_Unauthorized(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "", "admin", "school-1")
	body, _ := json.Marshal(SaveDeferredScrutinyRequest{StudentID: "s-1", ClassID: "c-1", FinalDecision: "promosso_con_debiti_saldati"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/deferred", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestScrutinyHandler_SaveDeferredScrutiny_BadRequest(t *testing.T) {
	// user_id present → passes 401 guard; binding fails → 400
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "admin-1", "admin", "school-1")
	body := []byte(`{}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/deferred", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── AsyncPDF endpoints ───────────────────────────────────────────────────────

func TestScrutinyHandler_EnqueueAsyncScrutinyPdf_Unauthorized(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/class/c-1/async-pdf?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestScrutinyHandler_EnqueueAsyncScrutinyPdf_Forbidden_Student(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "s-1", "student", "school-1")
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/class/c-1/async-pdf?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestScrutinyHandler_EnqueueAsyncScrutinyPdf_Forbidden_Parent(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "p-1", "parent", "school-1")
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/class/c-1/async-pdf?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestScrutinyHandler_EnqueueAsyncScrutinyPdf_NoPdfWorker(t *testing.T) {
	// Authorized user, nil pdfWorkerClient → 503
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "admin-1", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/class/c-1/async-pdf?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestScrutinyHandler_GetPdfJobStatus_Unauthorized(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/scrutiny/pdf-jobs/job-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestScrutinyHandler_GetPdfJobStatus_NoPdfWorker(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "admin-1", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/scrutiny/pdf-jobs/job-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestScrutinyHandler_EnqueueReportCardPdf_Unauthorized(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "", "admin", "school-1")
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/export/stud-1/async-pdf?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestScrutinyHandler_EnqueueReportCardPdf_NoPdfWorker(t *testing.T) {
	h := newScrutinyHandler()
	r := setupScrutinyRouter(h, "principal-1", "principal", "school-1")
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/scrutiny/export/stud-1/async-pdf?semester=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

// ─── Helper function unit tests ───────────────────────────────────────────────

func TestIsScrutinyManagementStaff(t *testing.T) {
	// All roles that should be allowed
	allowed := []string{
		"teacher", "coordinator", "coordinatore_classe",
		"admin", "superadmin", "principal", "vice_principal",
		"secretary", "assistente_alunni",
	}
	for _, role := range allowed {
		assert.True(t, isScrutinyManagementStaff(role), "role %q should be management staff", role)
	}

	// Roles that should NOT be allowed
	notAllowed := []string{"student", "parent", ""}
	for _, role := range notAllowed {
		assert.False(t, isScrutinyManagementStaff(role), "role %q should NOT be management staff", role)
	}
}

func TestIsScrutinyOverviewStaff(t *testing.T) {
	// Based on actual implementation:
	// admin, superadmin, principal, vice_principal, secretary, assistente_alunni
	allowed := []string{"admin", "superadmin", "principal", "vice_principal", "secretary", "assistente_alunni"}
	for _, role := range allowed {
		assert.True(t, isScrutinyOverviewStaff(role), "role %q should be overview staff", role)
	}

	notAllowed := []string{"teacher", "student", "parent", "coordinator", "dsga"}
	for _, role := range notAllowed {
		assert.False(t, isScrutinyOverviewStaff(role), "role %q should NOT be overview staff", role)
	}
}

func TestParseSemester(t *testing.T) {
	assert.Equal(t, 1, parseSemester("1"))
	assert.Equal(t, 2, parseSemester("2"))
	assert.Equal(t, 1, parseSemester(""))
	assert.Equal(t, 1, parseSemester("invalid"))
	assert.Equal(t, 2, parseSemester("second"))
	assert.Equal(t, 2, parseSemester("semester2"))
}
