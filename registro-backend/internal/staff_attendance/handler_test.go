package staff_attendance_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"registro-backend/internal/staff_attendance"
)

func setupTestRouter(handler *staff_attendance.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	rg := r.Group("/api/v1")
	handler.RegisterRoutes(rg)
	return r
}

func TestHandler_GetDailySummary(t *testing.T) {
	mockRepo := &MockRepo{
		summary: &staff_attendance.DailyStaffSummary{
			Date: "2026-09-12",
			TotalStaff: staff_attendance.RoleSummary{
				Total:   10,
				Present: 8,
			},
		},
	}
	service := staff_attendance.NewService(mockRepo)
	h := staff_attendance.NewHandler(service)
	router := setupTestRouter(h)

	// Case 1: Missing school_id -> 400
	req, _ := http.NewRequest("GET", "/api/v1/staff-attendance/summary?date=2026-09-12", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	// Case 2: Unauthorized role -> 403
	req, _ = http.NewRequest("GET", "/api/v1/staff-attendance/summary?date=2026-09-12&school_id=school-1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}

	// Case 3: Authorized (dsga) -> 200
	routerWithAuth := gin.New()
	routerWithAuth.Use(func(c *gin.Context) {
		c.Set("role", "dsga")
		c.Set("user_id", "dsga-1")
		c.Set("school_id", "school-1")
		c.Next()
	})
	h.RegisterRoutes(routerWithAuth.Group("/api/v1"))

	req, _ = http.NewRequest("GET", "/api/v1/staff-attendance/summary?date=2026-09-12", nil)
	w = httptest.NewRecorder()
	routerWithAuth.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_List(t *testing.T) {
	mockRepo := &MockRepo{
		list: []staff_attendance.StaffAttendance{
			{ID: "att-1", UserID: "u-1", Status: staff_attendance.StatusPresent},
		},
	}
	service := staff_attendance.NewService(mockRepo)
	h := staff_attendance.NewHandler(service)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", "admin")
		c.Set("user_id", "admin-1")
		c.Set("school_id", "school-1")
		c.Next()
	})
	h.RegisterRoutes(router.Group("/api/v1"))

	// Missing school_id check
	noSchoolRouter := gin.New()
	noSchoolRouter.Use(func(c *gin.Context) {
		c.Set("role", "admin")
		c.Next()
	})
	h.RegisterRoutes(noSchoolRouter.Group("/api/v1"))
	req, _ := http.NewRequest("GET", "/api/v1/staff-attendance?date=2026-09-12", nil)
	w := httptest.NewRecorder()
	noSchoolRouter.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing school_id, got %d", w.Code)
	}

	// Success
	req, _ = http.NewRequest("GET", "/api/v1/staff-attendance?date=2026-09-12", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestHandler_RecordAttendance(t *testing.T) {
	mockRepo := &MockRepo{
		upsertResult: &staff_attendance.StaffAttendance{
			ID:     "att-1",
			UserID: "u-1",
			Status: staff_attendance.StatusPresent,
		},
	}
	service := staff_attendance.NewService(mockRepo)
	h := staff_attendance.NewHandler(service)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", "dsga")
		c.Set("user_id", "dsga-1")
		c.Set("school_id", "school-1")
		c.Next()
	})
	h.RegisterRoutes(router.Group("/api/v1"))

	// Invalid JSON -> 400
	req, _ := http.NewRequest("POST", "/api/v1/staff-attendance", bytes.NewBufferString("{invalid"))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	// Valid JSON -> 200
	body, _ := json.Marshal(staff_attendance.UpsertStaffAttendanceRequest{
		UserID: "u-1",
		Date:   "2026-09-12",
		Status: staff_attendance.StatusPresent,
	})
	req, _ = http.NewRequest("POST", "/api/v1/staff-attendance", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_BulkRecordAttendance(t *testing.T) {
	mockRepo := &MockRepo{
		bulkResult: []staff_attendance.StaffAttendance{
			{ID: "att-1", UserID: "u-1", Status: staff_attendance.StatusPresent},
			{ID: "att-2", UserID: "u-2", Status: staff_attendance.StatusSickLeave},
		},
	}
	service := staff_attendance.NewService(mockRepo)
	h := staff_attendance.NewHandler(service)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", "dsga")
		c.Set("user_id", "dsga-1")
		c.Set("school_id", "school-1")
		c.Next()
	})
	h.RegisterRoutes(router.Group("/api/v1"))

	body, _ := json.Marshal(staff_attendance.BulkUpsertRequest{
		Date: "2026-09-12",
		Attendances: []staff_attendance.UpsertStaffAttendanceRequest{
			{UserID: "u-1", Date: "2026-09-12", Status: staff_attendance.StatusPresent},
			{UserID: "u-2", Date: "2026-09-12", Status: staff_attendance.StatusSickLeave},
		},
	})
	req, _ := http.NewRequest("POST", "/api/v1/staff-attendance/bulk", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_DeleteAttendance(t *testing.T) {
	mockRepo := &MockRepo{deleteErr: nil}
	service := staff_attendance.NewService(mockRepo)
	h := staff_attendance.NewHandler(service)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", "admin")
		c.Set("user_id", "admin-1")
		c.Set("school_id", "school-1")
		c.Next()
	})
	h.RegisterRoutes(router.Group("/api/v1"))

	req, _ := http.NewRequest("DELETE", "/api/v1/staff-attendance/att-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestHandler_BadgeEndpoints(t *testing.T) {
	mockRepo := &MockRepo{
		swipeResult: &staff_attendance.BadgeSwipe{
			ID:        "swipe-1",
			SchoolID:  "school-1",
			BadgeCode: "BDG123",
			SwipeType: "in",
			SwipeTime: time.Now(),
		},
		processCount: 5,
	}
	service := staff_attendance.NewService(mockRepo)
	h := staff_attendance.NewHandler(service)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("role", "dsga")
		c.Set("user_id", "dsga-1")
		c.Set("school_id", "school-1")
		c.Next()
	})
	h.RegisterRoutes(router.Group("/api/v1"))

	// 1. RegisterBadgeSwipe
	body, _ := json.Marshal(staff_attendance.BadgeSwipeRequest{
		BadgeCode: "BDG123",
		DeviceID:  "DEV001",
		SwipeType: "in",
	})
	req, _ := http.NewRequest("POST", "/api/v1/staff-attendance/badge-swipe", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	// 2. ProcessBadgeSwipes
	req, _ = http.NewRequest("POST", "/api/v1/staff-attendance/badge-swipe/process", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 3. AssignBadge
	assignBody, _ := json.Marshal(map[string]string{
		"user_id":    "u-1",
		"badge_code": "BDG123",
		"notes":      "Assegnato",
	})
	req, _ = http.NewRequest("POST", "/api/v1/staff-attendance/badges", bytes.NewBuffer(assignBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}
