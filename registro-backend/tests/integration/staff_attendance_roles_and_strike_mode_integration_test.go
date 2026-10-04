package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/staff_attendance"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupStaffAttendanceTestRouter(repo staff_attendance.Repository, role, userID, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Set("role", role)
		c.Set("school_id", schoolID)
		c.Next()
	})
	svc := staff_attendance.NewService(repo)
	handler := staff_attendance.NewHandler(svc)
	api := r.Group("/api/v1")
	handler.RegisterRoutes(api)
	return r
}

func TestIntegration_StaffAttendance_StrikeMode_And_RBAC(t *testing.T) {
	repo := newMockStaffLifecycleRepo()
	schoolID := "school-att-rbac-1"
	date := "2026-10-15"

	// 1. RBAC check on Daily Summary: DSGA and Secretary are permitted
	routerDSGA := setupStaffAttendanceTestRouter(repo, "dsga", "dsga-user-1", schoolID)
	reqSummary, _ := http.NewRequest(http.MethodGet, "/api/v1/staff-attendance/summary?date="+date, nil)
	wSummary := httptest.NewRecorder()
	routerDSGA.ServeHTTP(wSummary, reqSummary)
	assert.Equal(t, http.StatusOK, wSummary.Code)

	routerSec := setupStaffAttendanceTestRouter(repo, "secretary", "sec-user-1", schoolID)
	wSummarySec := httptest.NewRecorder()
	routerSec.ServeHTTP(wSummarySec, reqSummary)
	assert.Equal(t, http.StatusOK, wSummarySec.Code)

	// 2. RBAC check on Daily Summary: Collaboratore Scolastico, Teacher and Student are FORBIDDEN
	forbiddenRoles := []string{"collaboratore_scolastico", "teacher", "student", "parent"}
	for _, role := range forbiddenRoles {
		routerForb := setupStaffAttendanceTestRouter(repo, role, "unauth-user", schoolID)
		wForb := httptest.NewRecorder()
		routerForb.ServeHTTP(wForb, reqSummary)
		assert.Equal(t, http.StatusForbidden, wForb.Code, "Role %s should get 403 on staff-attendance summary", role)
	}

	// 3. Strike Mode toggle: DSGA sets strike mode
	strikeBody := map[string]interface{}{
		"date":          date,
		"is_strike_day": true,
	}
	bStrike, _ := json.Marshal(strikeBody)
	reqStrike, _ := http.NewRequest(http.MethodPost, "/api/v1/staff-attendance/strike-mode", bytes.NewBuffer(bStrike))
	reqStrike.Header.Set("Content-Type", "application/json")
	wStrike := httptest.NewRecorder()
	routerDSGA.ServeHTTP(wStrike, reqStrike)
	assert.Equal(t, http.StatusOK, wStrike.Code)

	// 4. Strike Mode toggle: Teacher attempting strike mode toggle gets 403 Forbidden
	routerTeacher := setupStaffAttendanceTestRouter(repo, "teacher", "prof-1", schoolID)
	wStrikeTeach := httptest.NewRecorder()
	reqStrikeTeach, _ := http.NewRequest(http.MethodPost, "/api/v1/staff-attendance/strike-mode", bytes.NewBuffer(bStrike))
	reqStrikeTeach.Header.Set("Content-Type", "application/json")
	routerTeacher.ServeHTTP(wStrikeTeach, reqStrikeTeach)
	assert.Equal(t, http.StatusForbidden, wStrikeTeach.Code)

	// 5. Bulk staff attendance recording by DSGA
	bulkReq := staff_attendance.BulkUpsertRequest{
		Date: date,
		Attendances: []staff_attendance.UpsertStaffAttendanceRequest{
			{
				UserID: "teacher-1",
				Date:   date,
				Status: staff_attendance.StatusPresent,
			},
			{
				UserID:      "teacher-2",
				Date:        date,
				Status:      staff_attendance.StatusOnStrike,
				IsStrikeDay: true,
				Notes:       "Adesione sciopero comparto scuola",
			},
			{
				UserID: "ata-1",
				Date:   date,
				Status: staff_attendance.StatusPresent,
				Notes:  "Contingente minimo garantito",
			},
		},
	}
	bBulk, _ := json.Marshal(bulkReq)
	reqBulk, _ := http.NewRequest(http.MethodPost, "/api/v1/staff-attendance/bulk", bytes.NewBuffer(bBulk))
	reqBulk.Header.Set("Content-Type", "application/json")
	wBulk := httptest.NewRecorder()
	routerDSGA.ServeHTTP(wBulk, reqBulk)
	require.Equal(t, http.StatusOK, wBulk.Code)

	var respObj struct {
		Data  []staff_attendance.StaffAttendance `json:"data"`
		Count int                                `json:"count"`
	}
	err := json.Unmarshal(wBulk.Body.Bytes(), &respObj)
	require.NoError(t, err)
	assert.Equal(t, 3, respObj.Count)
	assert.Len(t, respObj.Data, 3)

	// 6. Non-manager role (Teacher) attempting bulk attendance recording gets 403 Forbidden
	wBulkTeach := httptest.NewRecorder()
	reqBulkTeach, _ := http.NewRequest(http.MethodPost, "/api/v1/staff-attendance/bulk", bytes.NewBuffer(bBulk))
	reqBulkTeach.Header.Set("Content-Type", "application/json")
	routerTeacher.ServeHTTP(wBulkTeach, reqBulkTeach)
	assert.Equal(t, http.StatusForbidden, wBulkTeach.Code)
}
