package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/timetablegen"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestIntegration_TimetableAlternatives_PublishAndRBAC verifies end-to-end publishing of
// the 3 alternatives, custom slot overrides, and RBAC / multi-tenant boundary checks.
func TestIntegration_TimetableAlternatives_PublishAndRBAC(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockIntegrationTimetableRepo()
	cfg := timetablegen.DefaultConfig()
	cfg.MaxIterations = 200
	cfg.TimeLimitSeconds = 3
	svc := timetablegen.NewService(repo, timetablegen.NewGenerator(cfg))
	h := timetablegen.NewHandler(svc)

	schoolA := "school-uuid-tenant-a"
	schoolB := "school-uuid-tenant-b"
	userVP := "user-vp-1"
	userStudent := "user-student-1"

	// Setup pre-completed Job with 3 alternatives in school A
	jobID := uuid.New().String()
	t1 := "teacher-1"
	c1 := "class-1"
	s1 := "sub-math"

	alt1Slots := []timetablegen.GeneratedSlot{
		{DayOfWeek: 1, HourIndex: 1, ClassID: c1, SubjectID: s1, TeacherID: &t1},
		{DayOfWeek: 1, HourIndex: 2, ClassID: c1, SubjectID: s1, TeacherID: &t1},
	}
	alt2Slots := []timetablegen.GeneratedSlot{
		{DayOfWeek: 1, HourIndex: 2, ClassID: c1, SubjectID: s1, TeacherID: &t1},
		{DayOfWeek: 1, HourIndex: 3, ClassID: c1, SubjectID: s1, TeacherID: &t1},
	}
	alt3Slots := []timetablegen.GeneratedSlot{
		{DayOfWeek: 2, HourIndex: 1, ClassID: c1, SubjectID: s1, TeacherID: &t1},
		{DayOfWeek: 2, HourIndex: 2, ClassID: c1, SubjectID: s1, TeacherID: &t1},
	}

	resultSummary := timetablegen.TimetableGenerationResult{
		TotalSlots:    2,
		AssignedSlots: 2,
		CoveragePct:   100.0,
		Slots:         alt1Slots,
		Alternatives: []timetablegen.TimetableAlternative{
			{
				ID:            1,
				Label:         "Proposta 1: Bilanciata",
				Strategy:      "balanced",
				AssignedSlots: 2,
				Slots:         alt1Slots,
				Score:         850.0,
			},
			{
				ID:            2,
				Label:         "Proposta 2: Didattica",
				Strategy:      "didactic_first",
				AssignedSlots: 2,
				Slots:         alt2Slots,
				Score:         910.0,
			},
			{
				ID:            3,
				Label:         "Proposta 3: Compatta",
				Strategy:      "compact_teacher",
				AssignedSlots: 2,
				Slots:         alt3Slots,
				Score:         890.0,
			},
		},
	}
	rawSummary, _ := json.Marshal(resultSummary)
	now := time.Now()

	repo.jobs[jobID] = &timetablegen.TimetableJob{
		ID:            jobID,
		SchoolID:      schoolA,
		Status:        timetablegen.JobStatusCompleted,
		ResultSummary: rawSummary,
		CompletedAt:   &now,
	}

	// Routers with different RBAC contexts
	setupRouter := func(role, schoolID, userID string) *gin.Engine {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("role", role)
			c.Set("school_id", schoolID)
			c.Set("user_id", userID)
			c.Next()
		})
		h.RegisterRoutes(r.Group(""))
		return r
	}

	routerVP := setupRouter("vice_principal", schoolA, userVP)
	routerAdmin := setupRouter("admin", schoolA, "admin-1")
	routerStudent := setupRouter("student", schoolA, userStudent)
	routerForeignVP := setupRouter("vice_principal", schoolB, "foreign-vp")

	// 1. Student attempts to publish -> 403 Forbidden
	alt1ID := 1
	body, _ := json.Marshal(timetablegen.PublishScheduleRequest{AlternativeID: &alt1ID})
	req := httptest.NewRequest(http.MethodPost, "/timetable/generate/"+jobID+"/publish", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	routerStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 2. Foreign school VP attempts to publish school A's job -> 403 Forbidden
	req = httptest.NewRequest(http.MethodPost, "/timetable/generate/"+jobID+"/publish", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	routerForeignVP.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 3. Publishing an invalid alternative ID (99) -> 400 Bad Request
	badAltID := 99
	badBody, _ := json.Marshal(timetablegen.PublishScheduleRequest{AlternativeID: &badAltID})
	req = httptest.NewRequest(http.MethodPost, "/timetable/generate/"+jobID+"/publish", bytes.NewBuffer(badBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	routerVP.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 4. VP publishes Alternative 1 -> Success
	req = httptest.NewRequest(http.MethodPost, "/timetable/generate/"+jobID+"/publish", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	routerVP.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, repo.published, 2)
	assert.Equal(t, 1, repo.published[0].DayOfWeek)
	assert.Equal(t, 1, repo.published[0].HourIndex)

	// 5. Admin publishes Alternative 3 -> Success and updates published slots
	alt3ID := 3
	body3, _ := json.Marshal(timetablegen.PublishScheduleRequest{AlternativeID: &alt3ID})
	req = httptest.NewRequest(http.MethodPost, "/timetable/generate/"+jobID+"/publish", bytes.NewBuffer(body3))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	routerAdmin.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, repo.published, 2)
	assert.Equal(t, 2, repo.published[0].DayOfWeek)

	// 6. VP publishes custom edited slots -> Success and published contains custom slots
	customSlots := []timetablegen.GeneratedSlot{
		{DayOfWeek: 5, HourIndex: 5, ClassID: c1, SubjectID: s1, TeacherID: &t1},
	}
	customBody, _ := json.Marshal(timetablegen.PublishScheduleRequest{Slots: customSlots})
	req = httptest.NewRequest(http.MethodPost, "/timetable/generate/"+jobID+"/publish", bytes.NewBuffer(customBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	routerVP.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Len(t, repo.published, 1)
	assert.Equal(t, 5, repo.published[0].DayOfWeek)
	assert.Equal(t, 5, repo.published[0].HourIndex)
}
