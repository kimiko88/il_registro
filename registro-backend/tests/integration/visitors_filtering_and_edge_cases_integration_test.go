package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/visitors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_Visitors_Filtering_And_EdgeCases(t *testing.T) {
	repo := newMemoryVisitorsRepo()
	r := setupVisitorsRouter(repo)

	schoolID := "school-vis-filter-1"
	staffCS := "cs-user-1"
	staffAdmin := "admin-user-1"

	// 1. Register multiple visitors on different purposes
	visitor1 := visitors.RegisterVisitorRequest{
		Name:        "Ing. Roberto Neri",
		DocumentID:  "CI-554433",
		Purpose:     visitors.PurposeSupplier,
		HostName:    "DSGA",
		BadgeNumber: "B-201",
		Notes:       "Manutenzione ascensore",
	}
	b1, _ := json.Marshal(visitor1)
	req1, _ := http.NewRequest(http.MethodPost, "/api/v1/visitors", bytes.NewBuffer(b1))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("X-User-ID", staffCS)
	req1.Header.Set("X-User-Role", "collaboratore_scolastico")
	req1.Header.Set("X-School-ID", schoolID)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	require.Equal(t, http.StatusCreated, w1.Code)

	var v1 visitors.Visitor
	_ = json.Unmarshal(w1.Body.Bytes(), &v1)

	// 2. Reject visitor with empty name -> 400 Bad Request
	badVisitor := map[string]string{
		"name":    "",
		"purpose": "parent",
	}
	bBad, _ := json.Marshal(badVisitor)
	reqBad, _ := http.NewRequest(http.MethodPost, "/api/v1/visitors", bytes.NewBuffer(bBad))
	reqBad.Header.Set("Content-Type", "application/json")
	reqBad.Header.Set("X-User-ID", staffCS)
	reqBad.Header.Set("X-User-Role", "collaboratore_scolastico")
	reqBad.Header.Set("X-School-ID", schoolID)
	wBad := httptest.NewRecorder()
	r.ServeHTTP(wBad, reqBad)
	assert.Equal(t, http.StatusBadRequest, wBad.Code)

	// 3. Early Exit with missing student_id -> 400 Bad Request
	badExit := map[string]string{
		"delegatee_name": "Mario Rossi",
	}
	bBadExit, _ := json.Marshal(badExit)
	reqBadExit, _ := http.NewRequest(http.MethodPost, "/api/v1/visitors/early-exits", bytes.NewBuffer(bBadExit))
	reqBadExit.Header.Set("Content-Type", "application/json")
	reqBadExit.Header.Set("X-User-ID", staffCS)
	reqBadExit.Header.Set("X-User-Role", "collaboratore_scolastico")
	reqBadExit.Header.Set("X-School-ID", schoolID)
	wBadExit := httptest.NewRecorder()
	r.ServeHTTP(wBadExit, reqBadExit)
	assert.Equal(t, http.StatusBadRequest, wBadExit.Code)

	// 4. Create and filter maintenance reports by status
	maint1 := visitors.CreateMaintenanceReportRequest{
		Location:    "Palestra Plesso B",
		Category:    "strutturale",
		Description: "Maniglia porta antincendio bloccata",
		Priority:    visitors.PriorityUrgent,
	}
	bM1, _ := json.Marshal(maint1)
	reqM1, _ := http.NewRequest(http.MethodPost, "/api/v1/visitors/maintenance", bytes.NewBuffer(bM1))
	reqM1.Header.Set("Content-Type", "application/json")
	reqM1.Header.Set("X-User-ID", staffCS)
	reqM1.Header.Set("X-User-Role", "collaboratore_scolastico")
	reqM1.Header.Set("X-School-ID", schoolID)
	wM1 := httptest.NewRecorder()
	r.ServeHTTP(wM1, reqM1)
	require.Equal(t, http.StatusCreated, wM1.Code)

	var report1 visitors.MaintenanceReport
	_ = json.Unmarshal(wM1.Body.Bytes(), &report1)

	// 5. Update status to 'in_lavorazione' with assigned technician
	updateReq := visitors.UpdateMaintenanceStatusRequest{
		Status:     visitors.MaintenanceInProgress,
		AssignedTo: "Squadra Pronto Intervento Comune",
	}
	bUp, _ := json.Marshal(updateReq)
	reqUp, _ := http.NewRequest(http.MethodPatch, "/api/v1/visitors/maintenance/"+report1.ID+"/status", bytes.NewBuffer(bUp))
	reqUp.Header.Set("Content-Type", "application/json")
	reqUp.Header.Set("X-User-ID", staffAdmin)
	reqUp.Header.Set("X-User-Role", "admin")
	reqUp.Header.Set("X-School-ID", schoolID)
	wUp := httptest.NewRecorder()
	r.ServeHTTP(wUp, reqUp)
	assert.Equal(t, http.StatusOK, wUp.Code)

	// 6. Filter maintenance reports: ?status=in_lavorazione
	reqFilter, _ := http.NewRequest(http.MethodGet, "/api/v1/visitors/maintenance?status=in_lavorazione", nil)
	reqFilter.Header.Set("X-User-ID", staffCS)
	reqFilter.Header.Set("X-User-Role", "collaboratore_scolastico")
	reqFilter.Header.Set("X-School-ID", schoolID)
	wFilter := httptest.NewRecorder()
	r.ServeHTTP(wFilter, reqFilter)
	require.Equal(t, http.StatusOK, wFilter.Code)

	var filteredReports []*visitors.MaintenanceReport
	err := json.Unmarshal(wFilter.Body.Bytes(), &filteredReports)
	require.NoError(t, err)
	assert.Len(t, filteredReports, 1)
	assert.Equal(t, visitors.MaintenanceInProgress, filteredReports[0].Status)

	// 7. Verify student cannot access visitor maintenance -> 403 Forbidden
	reqStudent, _ := http.NewRequest(http.MethodGet, "/api/v1/visitors/maintenance", nil)
	reqStudent.Header.Set("X-User-ID", "student-user-1")
	reqStudent.Header.Set("X-User-Role", "student")
	reqStudent.Header.Set("X-School-ID", schoolID)
	wStudent := httptest.NewRecorder()
	r.ServeHTTP(wStudent, reqStudent)
	assert.Equal(t, http.StatusForbidden, wStudent.Code)

	// 8. Exit visitor v1
	reqExit, _ := http.NewRequest(http.MethodPatch, "/api/v1/visitors/"+v1.ID+"/exit", bytes.NewBufferString(`{"notes":"Terminato intervento"}`))
	reqExit.Header.Set("Content-Type", "application/json")
	reqExit.Header.Set("X-User-ID", staffCS)
	reqExit.Header.Set("X-User-Role", "collaboratore_scolastico")
	reqExit.Header.Set("X-School-ID", schoolID)
	wExit := httptest.NewRecorder()
	r.ServeHTTP(wExit, reqExit)
	assert.Equal(t, http.StatusOK, wExit.Code)

	// Listing visitors today confirms exit timestamp is now present
	todayStr := time.Now().Format("2006-01-02")
	reqToday, _ := http.NewRequest(http.MethodGet, "/api/v1/visitors?date="+todayStr, nil)
	reqToday.Header.Set("X-User-ID", staffCS)
	reqToday.Header.Set("X-User-Role", "collaboratore_scolastico")
	reqToday.Header.Set("X-School-ID", schoolID)
	wToday := httptest.NewRecorder()
	r.ServeHTTP(wToday, reqToday)
	require.Equal(t, http.StatusOK, wToday.Code)

	var todayList []*visitors.Visitor
	_ = json.Unmarshal(wToday.Body.Bytes(), &todayList)
	assert.NotEmpty(t, todayList)
	assert.NotNil(t, todayList[0].ExitTime)
}
