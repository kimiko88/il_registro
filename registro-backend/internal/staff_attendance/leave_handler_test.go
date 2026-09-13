package staff_attendance_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"registro-backend/internal/staff_attendance"
)

func setupLeaveRouter(h *staff_attendance.LeaveHandler, userID, schoolID, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		if schoolID != "" {
			c.Set("school_id", schoolID)
		}
		if role != "" {
			c.Set("role", role)
		}
		c.Next()
	})
	rg := r.Group("/api/v1")
	h.RegisterLeaveRoutes(rg)
	return r
}

func TestLeaveHandler_AuthCheck(t *testing.T) {
	mockRepo := &MockRepo{}
	h := staff_attendance.NewLeaveHandler(mockRepo)

	// Case 1: Unauthorized (no user_id) -> 401
	rNoUser := setupLeaveRouter(h, "", "school-1", "dsga")
	req, _ := http.NewRequest("GET", "/api/v1/staff-attendance/timecard", nil)
	w := httptest.NewRecorder()
	rNoUser.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}

	// Case 2: Forbidden (student role cannot read attendance) -> 403
	rStudent := setupLeaveRouter(h, "s-1", "school-1", "student")
	req, _ = http.NewRequest("GET", "/api/v1/staff-attendance/timecard", nil)
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestLeaveHandler_Timecard(t *testing.T) {
	mockRepo := &MockRepo{
		timecardResult: &staff_attendance.MonthlyTimecard{
			UserID:        "u-1",
			FirstName:     "Mario",
			LastName:      "Rossi",
			Role:          "collaboratore_scolastico",
			Month:         "2026-09",
			WorkedHours:   140.5,
			ContractHours: 144.0,
		},
		allTimecards: []staff_attendance.MonthlyTimecard{
			{
				UserID:        "u-1",
				FirstName:     "Mario",
				LastName:      "Rossi",
				Role:          "collaboratore_scolastico",
				WorkedHours:   140.5,
				ContractHours: 144.0,
			},
		},
	}
	h := staff_attendance.NewLeaveHandler(mockRepo)

	// 1. ATA user gets own timecard -> 200
	rAta := setupLeaveRouter(h, "u-1", "school-1", "collaboratore_scolastico")
	req, _ := http.NewRequest("GET", "/api/v1/staff-attendance/timecard", nil)
	w := httptest.NewRecorder()
	rAta.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 2. ATA user tries to view another user's timecard -> 403
	req, _ = http.NewRequest("GET", "/api/v1/staff-attendance/timecard?user_id=u-other", nil)
	w = httptest.NewRecorder()
	rAta.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}

	// 3. DSGA views another user's timecard -> 200
	rDsga := setupLeaveRouter(h, "dsga-1", "school-1", "dsga")
	req, _ = http.NewRequest("GET", "/api/v1/staff-attendance/timecard?user_id=u-1", nil)
	w = httptest.NewRecorder()
	rDsga.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 4. ExportTimecard by DSGA -> 200 (CSV)
	req, _ = http.NewRequest("GET", "/api/v1/staff-attendance/timecard/export?month=2026-09", nil)
	w = httptest.NewRecorder()
	rDsga.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Header().Get("Content-Type") == "" || w.Header().Get("Content-Type")[:8] != "text/csv" {
		t.Fatalf("expected text/csv prefix, got %s", w.Header().Get("Content-Type"))
	}
}

func TestLeaveHandler_LeavesCRUD(t *testing.T) {
	mockRepo := &MockRepo{
		createLeaveResult: &staff_attendance.LeaveRequest{
			ID:       "leave-1",
			SchoolID: "school-1",
			UserID:   "u-1",
			Type:     staff_attendance.LeaveTypeFerie,
			Status:   staff_attendance.LeaveStatusPending,
		},
		listLeavesResult: []staff_attendance.LeaveRequest{
			{ID: "leave-1", UserID: "u-1", Type: staff_attendance.LeaveTypeFerie, Status: staff_attendance.LeaveStatusPending},
		},
		getLeaveResult: &staff_attendance.LeaveRequest{
			ID:       "leave-1",
			SchoolID: "school-1",
			UserID:   "u-1",
			Type:     staff_attendance.LeaveTypeFerie,
			Status:   staff_attendance.LeaveStatusPending,
		},
	}
	h := staff_attendance.NewLeaveHandler(mockRepo)
	rDsga := setupLeaveRouter(h, "dsga-1", "school-1", "dsga")

	// 1. CreateLeave
	body, _ := json.Marshal(staff_attendance.CreateLeaveRequest{
		Type:      staff_attendance.LeaveTypeFerie,
		StartDate: "2026-09-15",
		EndDate:   "2026-09-18",
		Days:      4,
		Notes:     "Ferie estive residue",
	})
	req, _ := http.NewRequest("POST", "/api/v1/staff-attendance/leaves", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	rDsga.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}

	// 2. ListLeaves
	req, _ = http.NewRequest("GET", "/api/v1/staff-attendance/leaves", nil)
	w = httptest.NewRecorder()
	rDsga.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 3. GetLeave
	req, _ = http.NewRequest("GET", "/api/v1/staff-attendance/leaves/leave-1", nil)
	w = httptest.NewRecorder()
	rDsga.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 4. ApproveLeave
	appBody, _ := json.Marshal(map[string]string{"notes": "Approvate"})
	req, _ = http.NewRequest("PATCH", "/api/v1/staff-attendance/leaves/leave-1/approve", bytes.NewBuffer(appBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rDsga.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 5. RejectLeave
	rejBody, _ := json.Marshal(map[string]string{"reason": "Carenza personale"})
	req, _ = http.NewRequest("PATCH", "/api/v1/staff-attendance/leaves/leave-1/reject", bytes.NewBuffer(rejBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rDsga.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 6. DeleteLeave
	req, _ = http.NewRequest("DELETE", "/api/v1/staff-attendance/leaves/leave-1", nil)
	w = httptest.NewRecorder()
	rDsga.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
