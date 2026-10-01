package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/timetablegen"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_TeacherQuickPreferences_FullLifecycle(t *testing.T) {
	repo := newMockIntegrationTimetableRepo()
	schoolID := "school-quick-prefs-1"
	ay := "2024/2025"

	// 1. RBAC Check: Unauthorized roles cannot access or mutate quick preferences
	unauthorizedRoles := []string{"teacher", "student", "parent", "collaboratore_scolastico", "assistente_tecnico"}
	for _, role := range unauthorizedRoles {
		router := setupTimetableIntegrationRouter(repo, role, "user-unauth", schoolID)

		// GET forbidden
		reqGet, _ := http.NewRequest(http.MethodGet, "/timetable/teachers-quick-preferences", nil)
		wGet := httptest.NewRecorder()
		router.ServeHTTP(wGet, reqGet)
		assert.Equal(t, http.StatusForbidden, wGet.Code, "Role %s should get 403 on GET", role)

		// POST forbidden
		reqPost, _ := http.NewRequest(http.MethodPost, "/timetable/teachers-quick-preferences", bytes.NewBufferString(`{"preferences":[]}`))
		reqPost.Header.Set("Content-Type", "application/json")
		wPost := httptest.NewRecorder()
		router.ServeHTTP(wPost, reqPost)
		assert.Equal(t, http.StatusForbidden, wPost.Code, "Role %s should get 403 on POST", role)

		// PUT forbidden
		reqPut, _ := http.NewRequest(http.MethodPut, "/timetable/teachers-quick-preferences/t-1", bytes.NewBufferString(`{"teacher_id":"t-1"}`))
		reqPut.Header.Set("Content-Type", "application/json")
		wPut := httptest.NewRecorder()
		router.ServeHTTP(wPut, reqPut)
		assert.Equal(t, http.StatusForbidden, wPut.Code, "Role %s should get 403 on PUT", role)
	}

	// 2. Initial state: Principal queries overview
	routerPrincipal := setupTimetableIntegrationRouter(repo, "principal", "ds-1", schoolID)
	reqInit, _ := http.NewRequest(http.MethodGet, "/timetable/teachers-quick-preferences?academic_year_id="+ay, nil)
	wInit := httptest.NewRecorder()
	routerPrincipal.ServeHTTP(wInit, reqInit)
	require.Equal(t, http.StatusOK, wInit.Code)

	var initResp timetablegen.TeacherQuickPreferencesOverviewResponse
	err := json.Unmarshal(wInit.Body.Bytes(), &initResp)
	require.NoError(t, err)
	assert.Empty(t, initResp.Teachers)

	// 3. Batch Save: Secretary saves preferences for 3 teachers
	routerSecretary := setupTimetableIntegrationRouter(repo, "secretary", "sec-1", schoolID)
	batchPayload := timetablegen.SaveTeacherQuickPreferencesRequest{
		AcademicYearID: &ay,
		Preferences: []timetablegen.TeacherQuickPreferenceItem{
			{
				TeacherID:      "t-rossi",
				TeacherName:    "Mario Rossi",
				SubjectName:    "Matematica",
				DayOff:         1, // Lunedì
				TimeSlotPref:   "early_hours",
				MaxHoursPerDay: 5,
			},
			{
				TeacherID:      "t-bianchi",
				TeacherName:    "Giulia Bianchi",
				SubjectName:    "Italiano",
				DayOff:         1, // Lunedì
				TimeSlotPref:   "late_hours",
				MaxHoursPerDay: 4,
			},
			{
				TeacherID:      "t-verdi",
				TeacherName:    "Luca Verdi",
				SubjectName:    "Scienze",
				DayOff:         3, // Mercoledì
				TimeSlotPref:   "none",
				MaxHoursPerDay: 6,
			},
		},
	}
	bodyBytes, _ := json.Marshal(batchPayload)
	reqSave, _ := http.NewRequest(http.MethodPost, "/timetable/teachers-quick-preferences", bytes.NewBuffer(bodyBytes))
	reqSave.Header.Set("Content-Type", "application/json")
	wSave := httptest.NewRecorder()
	routerSecretary.ServeHTTP(wSave, reqSave)
	assert.Equal(t, http.StatusOK, wSave.Code)

	// 4. Verification: Vice Principal retrieves updated overview and KPI distributions
	routerVP := setupTimetableIntegrationRouter(repo, "vice_principal", "vp-1", schoolID)
	reqOverview, _ := http.NewRequest(http.MethodGet, "/timetable/teachers-quick-preferences?academic_year_id="+ay, nil)
	wOverview := httptest.NewRecorder()
	routerVP.ServeHTTP(wOverview, reqOverview)
	require.Equal(t, http.StatusOK, wOverview.Code)

	var overviewResp timetablegen.TeacherQuickPreferencesOverviewResponse
	err = json.Unmarshal(wOverview.Body.Bytes(), &overviewResp)
	require.NoError(t, err)
	assert.Len(t, overviewResp.Teachers, 3)
	// Day off 1 (Lunedì) has 2 teachers
	assert.Equal(t, 2, overviewResp.DayOffCounts[1])
	// Day off 3 (Mercoledì) has 1 teacher
	assert.Equal(t, 1, overviewResp.DayOffCounts[3])

	// 5. Single Teacher Update: Collaboratore DS modifies t-bianchi's day off to Venerdì (5)
	routerCollab := setupTimetableIntegrationRouter(repo, "collaboratore_ds", "collab-1", schoolID)
	singleItem := timetablegen.TeacherQuickPreferenceItem{
		TeacherID:      "t-bianchi",
		DayOff:         5, // Spostato a Venerdì
		TimeSlotPref:   "early_hours",
		MaxHoursPerDay: 4,
	}
	singleBytes, _ := json.Marshal(singleItem)
	reqSingle, _ := http.NewRequest(http.MethodPut, "/timetable/teachers-quick-preferences/t-bianchi?academic_year_id="+ay, bytes.NewBuffer(singleBytes))
	reqSingle.Header.Set("Content-Type", "application/json")
	wSingle := httptest.NewRecorder()
	routerCollab.ServeHTTP(wSingle, reqSingle)
	assert.Equal(t, http.StatusOK, wSingle.Code)

	// 6. Verification: Overview now reflects rebalanced distribution
	reqRebalanced, _ := http.NewRequest(http.MethodGet, "/timetable/teachers-quick-preferences?academic_year_id="+ay, nil)
	wRebalanced := httptest.NewRecorder()
	routerPrincipal.ServeHTTP(wRebalanced, reqRebalanced)
	require.Equal(t, http.StatusOK, wRebalanced.Code)

	var rebalResp timetablegen.TeacherQuickPreferencesOverviewResponse
	err = json.Unmarshal(wRebalanced.Body.Bytes(), &rebalResp)
	require.NoError(t, err)
	assert.Len(t, rebalResp.Teachers, 3)
	assert.Equal(t, 1, rebalResp.DayOffCounts[1], "Lunedì should now have 1 teacher")
	assert.Equal(t, 1, rebalResp.DayOffCounts[3], "Mercoledì should have 1 teacher")
	assert.Equal(t, 1, rebalResp.DayOffCounts[5], "Venerdì should now have 1 teacher")

	// 7. Bad Request Validation
	reqBadJSON, _ := http.NewRequest(http.MethodPost, "/timetable/teachers-quick-preferences", bytes.NewBufferString(`{invalid-json`))
	reqBadJSON.Header.Set("Content-Type", "application/json")
	wBadJSON := httptest.NewRecorder()
	routerPrincipal.ServeHTTP(wBadJSON, reqBadJSON)
	assert.Equal(t, http.StatusBadRequest, wBadJSON.Code)
}
