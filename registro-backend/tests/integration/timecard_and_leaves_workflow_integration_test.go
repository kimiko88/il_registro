package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/staff_attendance"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type memoryTimecardLeavesRepo struct {
	staff_attendance.Repository
	leaves    map[string]*staff_attendance.LeaveRequest
	timecards map[string]*staff_attendance.MonthlyTimecard
}

func newMemoryTimecardLeavesRepo() *memoryTimecardLeavesRepo {
	return &memoryTimecardLeavesRepo{
		leaves:    make(map[string]*staff_attendance.LeaveRequest),
		timecards: make(map[string]*staff_attendance.MonthlyTimecard),
	}
}

func (m *memoryTimecardLeavesRepo) CreateLeaveRequest(ctx context.Context, schoolID, userID string, req staff_attendance.CreateLeaveRequest) (*staff_attendance.LeaveRequest, error) {
	id := fmt.Sprintf("leave-%d", len(m.leaves)+1)
	now := time.Now()
	days := req.Days
	if days == 0 && req.Hours == 0 {
		days = 1
	}
	lr := &staff_attendance.LeaveRequest{
		ID:        id,
		SchoolID:  schoolID,
		UserID:    userID,
		Type:      req.Type,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Days:      days,
		Hours:     req.Hours,
		Notes:     req.Notes,
		Status:    staff_attendance.LeaveStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	m.leaves[id] = lr
	return lr, nil
}

func (m *memoryTimecardLeavesRepo) ListLeaveRequests(ctx context.Context, schoolID, userID, status string) ([]staff_attendance.LeaveRequest, error) {
	var out []staff_attendance.LeaveRequest
	for _, l := range m.leaves {
		if schoolID != "" && l.SchoolID != schoolID {
			continue
		}
		if userID != "" && l.UserID != userID {
			continue
		}
		if status != "" && string(l.Status) != status {
			continue
		}
		out = append(out, *l)
	}
	return out, nil
}

func (m *memoryTimecardLeavesRepo) GetLeaveRequest(ctx context.Context, schoolID, id string) (*staff_attendance.LeaveRequest, error) {
	lr, ok := m.leaves[id]
	if !ok || (schoolID != "" && lr.SchoolID != schoolID) {
		return nil, errors.New("richiesta non trovata")
	}
	return lr, nil
}

func (m *memoryTimecardLeavesRepo) ApproveLeaveRequest(ctx context.Context, schoolID, id, approvedBy, notes string) error {
	lr, ok := m.leaves[id]
	if !ok || (schoolID != "" && lr.SchoolID != schoolID) {
		return errors.New("richiesta non trovata")
	}
	now := time.Now()
	lr.Status = staff_attendance.LeaveStatusApproved
	lr.ApprovedBy = &approvedBy
	lr.ApprovedAt = &now
	lr.UpdatedAt = now
	return nil
}

func (m *memoryTimecardLeavesRepo) RejectLeaveRequest(ctx context.Context, schoolID, id, rejectedBy, reason string) error {
	lr, ok := m.leaves[id]
	if !ok || (schoolID != "" && lr.SchoolID != schoolID) {
		return errors.New("richiesta non trovata")
	}
	now := time.Now()
	lr.Status = staff_attendance.LeaveStatusRejected
	lr.RejectReason = reason
	lr.RejectedAt = &now
	lr.UpdatedAt = now
	return nil
}

func (m *memoryTimecardLeavesRepo) DeleteLeaveRequest(ctx context.Context, schoolID, id, userID string) error {
	lr, ok := m.leaves[id]
	if !ok || (schoolID != "" && lr.SchoolID != schoolID) {
		return errors.New("richiesta non trovata")
	}
	if lr.UserID != userID {
		return errors.New("non puoi eliminare una richiesta altrui")
	}
	if lr.Status != staff_attendance.LeaveStatusPending {
		return errors.New("non puoi eliminare una richiesta già processata")
	}
	delete(m.leaves, id)
	return nil
}

func (m *memoryTimecardLeavesRepo) GetMonthlyTimecard(ctx context.Context, schoolID, userID, month string) (*staff_attendance.MonthlyTimecard, error) {
	key := schoolID + ":" + userID + ":" + month
	if tc, ok := m.timecards[key]; ok {
		return tc, nil
	}
	return &staff_attendance.MonthlyTimecard{
		UserID:        userID,
		Month:         month,
		FirstName:     "Mario",
		LastName:      "Rossi",
		Role:          "collaboratore_scolastico",
		ContractHours: 156.0,
		WorkedHours:   160.0,
		OvertimeHours: 4.0,
		AbsenceDays:   0,
		LeaveDays:     2,
	}, nil
}

func (m *memoryTimecardLeavesRepo) GetAllMonthlyTimecards(ctx context.Context, schoolID, month string) ([]staff_attendance.MonthlyTimecard, error) {
	tc, _ := m.GetMonthlyTimecard(ctx, schoolID, "user-ata-1", month)
	return []staff_attendance.MonthlyTimecard{*tc}, nil
}

func setupTimecardLeavesIntegrationRouter(repo staff_attendance.Repository, currentUserID, currentRole, currentSchoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(func(c *gin.Context) {
		if currentUserID != "" {
			c.Set("user_id", currentUserID)
		}
		if currentRole != "" {
			c.Set("role", currentRole)
		}
		if currentSchoolID != "" {
			c.Set("school_id", currentSchoolID)
		}
		c.Next()
	})

	handler := staff_attendance.NewLeaveHandler(repo)
	v1 := r.Group("/api/v1")
	handler.RegisterLeaveRoutes(v1)

	return r
}

func TestTimecardAndLeavesWorkflow_Integration(t *testing.T) {
	repo := newMemoryTimecardLeavesRepo()
	schoolID := "school-liceo-virgilio"

	// 1. Unauthorized role (student) attempts to access timecard -> 403 Forbidden
	{
		rStudent := setupTimecardLeavesIntegrationRouter(repo, "student-1", "student", schoolID)
		req := httptest.NewRequest("GET", "/api/v1/staff-attendance/timecard", nil)
		w := httptest.NewRecorder()
		rStudent.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	}

	// 2. ATA Staff member (collaboratore_scolastico) consults their monthly timecard
	{
		rAta := setupTimecardLeavesIntegrationRouter(repo, "user-ata-1", "collaboratore_scolastico", schoolID)
		req := httptest.NewRequest("GET", "/api/v1/staff-attendance/timecard?month=2026-11", nil)
		w := httptest.NewRecorder()
		rAta.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var tc staff_attendance.MonthlyTimecard
		err := json.Unmarshal(w.Body.Bytes(), &tc)
		require.NoError(t, err)
		assert.Equal(t, "user-ata-1", tc.UserID)
		assert.Equal(t, 160.0, tc.WorkedHours)
		assert.Equal(t, 4.0, tc.OvertimeHours)
	}

	// 3. ATA Staff member submits a formal leave request
	var createdLeaveID string
	{
		rAta := setupTimecardLeavesIntegrationRouter(repo, "user-ata-1", "collaboratore_scolastico", schoolID)
		body, _ := json.Marshal(staff_attendance.CreateLeaveRequest{
			Type:      "ferie",
			StartDate: "2026-11-02",
			EndDate:   "2026-11-04",
			Days:      3,
			Notes:     "Ponte festività dei santi",
		})
		req := httptest.NewRequest("POST", "/api/v1/staff-attendance/leaves", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		rAta.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
		var lr staff_attendance.LeaveRequest
		err := json.Unmarshal(w.Body.Bytes(), &lr)
		require.NoError(t, err)
		createdLeaveID = lr.ID
		assert.NotEmpty(t, createdLeaveID)
		assert.Equal(t, staff_attendance.LeaveStatusPending, lr.Status)
		assert.Equal(t, staff_attendance.LeaveTypeFerie, lr.Type)
		assert.Equal(t, float64(3), lr.Days)
	}

	// 4. Unauthorized other staff member tries to inspect the leave request -> 403 Forbidden
	{
		rOtherAta := setupTimecardLeavesIntegrationRouter(repo, "user-ata-2", "assistente_tecnico", schoolID)
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/staff-attendance/leaves/%s", createdLeaveID), nil)
		w := httptest.NewRecorder()
		rOtherAta.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	}

	// 5. DSGA lists all pending leaves
	{
		rDsga := setupTimecardLeavesIntegrationRouter(repo, "dsga-1", "dsga", schoolID)
		req := httptest.NewRequest("GET", "/api/v1/staff-attendance/leaves?status=pending", nil)
		w := httptest.NewRecorder()
		rDsga.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var leaves []staff_attendance.LeaveRequest
		err := json.Unmarshal(w.Body.Bytes(), &leaves)
		require.NoError(t, err)
		require.Len(t, leaves, 1)
		assert.Equal(t, createdLeaveID, leaves[0].ID)
		assert.Equal(t, staff_attendance.LeaveStatusPending, leaves[0].Status)
	}

	// 6. DSGA approves the leave request
	{
		rDsga := setupTimecardLeavesIntegrationRouter(repo, "dsga-1", "dsga", schoolID)
		body, _ := json.Marshal(map[string]string{"notes": "Approvato nel piano ferie d'istituto"})
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/api/v1/staff-attendance/leaves/%s/approve", createdLeaveID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		rDsga.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var res map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &res)
		assert.Equal(t, "richiesta approvata", res["message"])

		// Verify state in repo
		lr, err := repo.GetLeaveRequest(context.Background(), schoolID, createdLeaveID)
		require.NoError(t, err)
		assert.Equal(t, staff_attendance.LeaveStatusApproved, lr.Status)
		assert.NotNil(t, lr.ApprovedBy)
		assert.Equal(t, "dsga-1", *lr.ApprovedBy)
	}

	// 7. ATA user attempts to delete an already approved leave -> 400 Bad Request
	{
		rAta := setupTimecardLeavesIntegrationRouter(repo, "user-ata-1", "collaboratore_scolastico", schoolID)
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/staff-attendance/leaves/%s", createdLeaveID), nil)
		w := httptest.NewRecorder()
		rAta.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	}

	// 8. DSGA exports the global monthly timecards in CSV format
	{
		rDsga := setupTimecardLeavesIntegrationRouter(repo, "dsga-1", "dsga", schoolID)
		req := httptest.NewRequest("GET", "/api/v1/staff-attendance/timecard/export?month=2026-11", nil)
		w := httptest.NewRecorder()
		rDsga.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "text/csv; charset=utf-8", w.Header().Get("Content-Type"))
		bodyStr := w.Body.String()
		assert.Contains(t, bodyStr, "Cognome,Nome,Ruolo,Mese,Ore Contratto,Ore Lavorate,Straordinari")
		assert.Contains(t, bodyStr, "Rossi,Mario,collaboratore_scolastico,2026-11")
	}
}
