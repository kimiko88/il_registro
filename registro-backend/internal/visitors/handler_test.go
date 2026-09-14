package visitors_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"registro-backend/internal/visitors"
)

type mockVisitorService struct {
	registerVisitorResult *visitors.Visitor
	registerVisitorErr    error
	recordExitErr         error
	todayVisitorsResult   []*visitors.Visitor
	todayVisitorsErr      error

	earlyExitResult       *visitors.EarlyExit
	earlyExitErr          error
	studentReturnErr      error
	todayEarlyExitsResult []*visitors.EarlyExit
	todayEarlyExitsErr    error

	maintenanceResult     *visitors.MaintenanceReport
	maintenanceErr        error
	listMaintenanceResult []*visitors.MaintenanceReport
	listMaintenanceErr    error
	updateStatusErr       error
}

func (m *mockVisitorService) RegisterVisitor(ctx context.Context, schoolID, recordedBy string, req visitors.RegisterVisitorRequest) (*visitors.Visitor, error) {
	return m.registerVisitorResult, m.registerVisitorErr
}

func (m *mockVisitorService) RecordVisitorExit(ctx context.Context, id, schoolID, notes string) error {
	return m.recordExitErr
}

func (m *mockVisitorService) ListTodayVisitors(ctx context.Context, schoolID, date string) ([]*visitors.Visitor, error) {
	return m.todayVisitorsResult, m.todayVisitorsErr
}

func (m *mockVisitorService) RecordEarlyExit(ctx context.Context, schoolID, recordedBy string, req visitors.RecordEarlyExitRequest) (*visitors.EarlyExit, error) {
	return m.earlyExitResult, m.earlyExitErr
}

func (m *mockVisitorService) RecordStudentReturn(ctx context.Context, id, schoolID, notes string) error {
	return m.studentReturnErr
}

func (m *mockVisitorService) ListTodayEarlyExits(ctx context.Context, schoolID, date string) ([]*visitors.EarlyExit, error) {
	return m.todayEarlyExitsResult, m.todayEarlyExitsErr
}

func (m *mockVisitorService) CreateMaintenanceReport(ctx context.Context, schoolID, reportedBy string, req visitors.CreateMaintenanceReportRequest) (*visitors.MaintenanceReport, error) {
	return m.maintenanceResult, m.maintenanceErr
}

func (m *mockVisitorService) ListMaintenanceReports(ctx context.Context, schoolID, statusFilter string) ([]*visitors.MaintenanceReport, error) {
	return m.listMaintenanceResult, m.listMaintenanceErr
}

func (m *mockVisitorService) UpdateMaintenanceStatus(ctx context.Context, id, schoolID string, req visitors.UpdateMaintenanceStatusRequest) error {
	return m.updateStatusErr
}

func setupVisitorRouter(h *visitors.Handler, userID, schoolID, role string) *gin.Engine {
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
	h.RegisterRoutes(rg)
	return r
}

func TestVisitorHandler_AuthCheck(t *testing.T) {
	mockSvc := &mockVisitorService{}
	h := visitors.NewHandler(mockSvc)

	// 1. Missing user_id -> 401
	rNoUser := setupVisitorRouter(h, "", "school-1", "collaboratore_scolastico")
	req, _ := http.NewRequest("GET", "/api/v1/visitors", nil)
	w := httptest.NewRecorder()
	rNoUser.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}

	// 2. Unauthorized role (student) -> 403
	rStudent := setupVisitorRouter(h, "s-1", "school-1", "student")
	req, _ = http.NewRequest("GET", "/api/v1/visitors", nil)
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestVisitorHandler_VisitorsCRUD(t *testing.T) {
	badge := "B-101"
	mockSvc := &mockVisitorService{
		registerVisitorResult: &visitors.Visitor{
			ID:          "v-1",
			SchoolID:    "school-1",
			Name:        "Mario Rossi",
			Purpose:     visitors.PurposeSupplier,
			BadgeNumber: &badge,
			EntryTime:   time.Now(),
		},
		todayVisitorsResult: []*visitors.Visitor{
			{
				ID:          "v-1",
				Name:        "Mario Rossi",
				Purpose:     visitors.PurposeSupplier,
				BadgeNumber: &badge,
			},
		},
	}
	h := visitors.NewHandler(mockSvc)
	router := setupVisitorRouter(h, "staff-1", "school-1", "collaboratore_scolastico")

	// 1. RegisterVisitor: invalid JSON -> 400
	req, _ := http.NewRequest("POST", "/api/v1/visitors", bytes.NewBufferString("{invalid"))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	// 2. RegisterVisitor: success -> 201
	body, _ := json.Marshal(visitors.RegisterVisitorRequest{
		Name:        "Mario Rossi",
		Purpose:     visitors.PurposeSupplier,
		BadgeNumber: "B-101",
	})
	req, _ = http.NewRequest("POST", "/api/v1/visitors", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}

	// 3. ListTodayVisitors -> 200
	req, _ = http.NewRequest("GET", "/api/v1/visitors?date=2026-09-12", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 4. RecordVisitorExit -> 200
	exitBody, _ := json.Marshal(visitors.RecordExitRequest{Notes: "Uscito"})
	req, _ = http.NewRequest("PATCH", "/api/v1/visitors/v-1/exit", bytes.NewBuffer(exitBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 5. RecordVisitorExit error -> 400
	mockSvc.recordExitErr = errors.New("visitatore già uscito")
	req, _ = http.NewRequest("PATCH", "/api/v1/visitors/v-1/exit", bytes.NewBuffer(exitBody))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestVisitorHandler_EarlyExits(t *testing.T) {
	mockSvc := &mockVisitorService{
		earlyExitResult: &visitors.EarlyExit{
			ID:            "exit-1",
			SchoolID:      "school-1",
			StudentID:     "stud-1",
			DelegateeName: "Giuseppe Verdi",
			ExitTime:      time.Now(),
		},
		todayEarlyExitsResult: []*visitors.EarlyExit{
			{ID: "exit-1", StudentID: "stud-1", DelegateeName: "Giuseppe Verdi"},
		},
	}
	h := visitors.NewHandler(mockSvc)
	router := setupVisitorRouter(h, "staff-1", "school-1", "collaboratore_scolastico")

	// 1. RecordEarlyExit: invalid JSON -> 400
	req, _ := http.NewRequest("POST", "/api/v1/visitors/early-exits", bytes.NewBufferString("{invalid"))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	// 2. RecordEarlyExit: success -> 201
	body, _ := json.Marshal(visitors.RecordEarlyExitRequest{
		StudentID:     "stud-1",
		DelegateeName: "Giuseppe Verdi",
		DelegateRel:   "genitore",
		ReasonCode:    "medica",
	})
	req, _ = http.NewRequest("POST", "/api/v1/visitors/early-exits", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	// 3. ListTodayEarlyExits -> 200
	req, _ = http.NewRequest("GET", "/api/v1/visitors/early-exits", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 4. RecordStudentReturn -> 200
	retBody, _ := json.Marshal(visitors.RecordStudentReturnRequest{Notes: "Rientrato"})
	req, _ = http.NewRequest("PATCH", "/api/v1/visitors/early-exits/exit-1/return", bytes.NewBuffer(retBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestVisitorHandler_Maintenance(t *testing.T) {
	mockSvc := &mockVisitorService{
		maintenanceResult: &visitors.MaintenanceReport{
			ID:          "maint-1",
			SchoolID:    "school-1",
			Location:    "Aula 3B Piano 1",
			Category:    "strutturale",
			Description: "Finestra rotta aula 3B",
			Priority:    visitors.PriorityHigh,
			Status:      visitors.MaintenanceOpen,
		},
		listMaintenanceResult: []*visitors.MaintenanceReport{
			{ID: "maint-1", Location: "Aula 3B", Description: "Finestra rotta aula 3B", Status: visitors.MaintenanceOpen},
		},
	}
	h := visitors.NewHandler(mockSvc)
	router := setupVisitorRouter(h, "staff-1", "school-1", "assistente_amministrativo")

	// 1. CreateMaintenanceReport -> 201
	body, _ := json.Marshal(visitors.CreateMaintenanceReportRequest{
		Location:    "Aula 3B Piano 1",
		Category:    "strutturale",
		Description: "Vetro crepato",
		Priority:    visitors.PriorityHigh,
	})
	req, _ := http.NewRequest("POST", "/api/v1/visitors/maintenance", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	// 2. ListMaintenanceReports -> 200
	req, _ = http.NewRequest("GET", "/api/v1/visitors/maintenance?status=aperto", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 3. UpdateMaintenanceStatus -> 200
	statusBody, _ := json.Marshal(visitors.UpdateMaintenanceStatusRequest{
		Status:     visitors.MaintenanceClosed,
		AssignedTo: "tecnico-1",
	})
	req, _ = http.NewRequest("PATCH", "/api/v1/visitors/maintenance/maint-1/status", bytes.NewBuffer(statusBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
