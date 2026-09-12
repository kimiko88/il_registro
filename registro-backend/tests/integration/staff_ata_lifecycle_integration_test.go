package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"registro-backend/internal/personnel_desk"
	"registro-backend/internal/staff_attendance"
	"registro-backend/internal/visitors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// In-memory mock for staff attendance & leave repository
type mockStaffLifecycleRepo struct {
	mu           sync.Mutex
	attendances  map[string]*staff_attendance.StaffAttendance
	swipes       []*staff_attendance.BadgeSwipe
	badges       map[string]staff_attendance.UserBadge
	leaveCounter int
	leaves       map[string]*staff_attendance.LeaveRequest
}

func newMockStaffLifecycleRepo() *mockStaffLifecycleRepo {
	return &mockStaffLifecycleRepo{
		attendances: make(map[string]*staff_attendance.StaffAttendance),
		badges:      make(map[string]staff_attendance.UserBadge),
		leaves:      make(map[string]*staff_attendance.LeaveRequest),
	}
}

func (m *mockStaffLifecycleRepo) GetDailySummary(ctx context.Context, schoolID, date string) (*staff_attendance.DailyStaffSummary, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	summary := &staff_attendance.DailyStaffSummary{
		Date: date,
		TotalStaff: staff_attendance.RoleSummary{
			Total:   len(m.attendances),
			Present: len(m.attendances),
		},
	}
	return summary, nil
}

func (m *mockStaffLifecycleRepo) ListByDate(ctx context.Context, schoolID, date string) ([]staff_attendance.StaffAttendance, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []staff_attendance.StaffAttendance
	for _, a := range m.attendances {
		list = append(list, *a)
	}
	return list, nil
}

func (m *mockStaffLifecycleRepo) Upsert(ctx context.Context, schoolID, recordedBy string, req staff_attendance.UpsertStaffAttendanceRequest) (*staff_attendance.StaffAttendance, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	rec := recordedBy
	att := &staff_attendance.StaffAttendance{
		ID:         "att-" + req.UserID,
		SchoolID:   schoolID,
		UserID:     req.UserID,
		Date:       req.Date,
		Status:     req.Status,
		RecordedBy: &rec,
	}
	m.attendances[att.ID] = att
	return att, nil
}

func (m *mockStaffLifecycleRepo) BulkUpsert(ctx context.Context, schoolID, recordedBy string, req staff_attendance.BulkUpsertRequest) ([]staff_attendance.StaffAttendance, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []staff_attendance.StaffAttendance
	rec := recordedBy
	for _, item := range req.Attendances {
		att := &staff_attendance.StaffAttendance{
			ID:         "att-" + item.UserID,
			SchoolID:   schoolID,
			UserID:     item.UserID,
			Date:       req.Date,
			Status:     item.Status,
			RecordedBy: &rec,
		}
		m.attendances[att.ID] = att
		list = append(list, *att)
	}
	return list, nil
}

func (m *mockStaffLifecycleRepo) GetByUserAndDate(ctx context.Context, schoolID, userID, date string) (*staff_attendance.StaffAttendance, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, a := range m.attendances {
		if a.UserID == userID && a.Date == date {
			return a, nil
		}
	}
	return nil, nil
}

func (m *mockStaffLifecycleRepo) Delete(ctx context.Context, schoolID, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.attendances, id)
	return nil
}

func (m *mockStaffLifecycleRepo) RegisterBadgeSwipe(ctx context.Context, schoolID string, req staff_attendance.BadgeSwipeRequest) (*staff_attendance.BadgeSwipe, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	swipe := &staff_attendance.BadgeSwipe{
		ID:        fmt.Sprintf("swipe-%d", len(m.swipes)+1),
		SchoolID:  schoolID,
		BadgeCode: req.BadgeCode,
		SwipeType: req.SwipeType,
		SwipeTime: time.Now(),
		Processed: false,
	}
	m.swipes = append(m.swipes, swipe)
	return swipe, nil
}

func (m *mockStaffLifecycleRepo) ProcessPendingSwipes(ctx context.Context, schoolID string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	count := 0
	for _, s := range m.swipes {
		if !s.Processed {
			s.Processed = true
			count++
		}
	}
	return count, nil
}

func (m *mockStaffLifecycleRepo) AssignBadge(ctx context.Context, badge staff_attendance.UserBadge) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.badges[badge.BadgeCode] = badge
	return nil
}

func (m *mockStaffLifecycleRepo) RevokeBadge(ctx context.Context, schoolID, userID, badgeCode string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.badges, badgeCode)
	return nil
}

func (m *mockStaffLifecycleRepo) ListBadges(ctx context.Context, schoolID string) ([]staff_attendance.UserBadge, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []staff_attendance.UserBadge
	for _, b := range m.badges {
		list = append(list, b)
	}
	return list, nil
}

// Leaves
func (m *mockStaffLifecycleRepo) CreateLeaveRequest(ctx context.Context, schoolID, userID string, req staff_attendance.CreateLeaveRequest) (*staff_attendance.LeaveRequest, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.leaveCounter++
	id := fmt.Sprintf("leave-%d", m.leaveCounter)
	l := &staff_attendance.LeaveRequest{
		ID:        id,
		SchoolID:  schoolID,
		UserID:    userID,
		Type:      req.Type,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Days:      req.Days,
		Notes:     req.Notes,
		Status:    staff_attendance.LeaveStatusPending,
		CreatedAt: time.Now(),
	}
	m.leaves[id] = l
	return l, nil
}

func (m *mockStaffLifecycleRepo) ListLeaveRequests(ctx context.Context, schoolID, userID, status string) ([]staff_attendance.LeaveRequest, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []staff_attendance.LeaveRequest
	for _, l := range m.leaves {
		if userID != "" && l.UserID != userID {
			continue
		}
		if status != "" && string(l.Status) != status {
			continue
		}
		list = append(list, *l)
	}
	return list, nil
}

func (m *mockStaffLifecycleRepo) GetLeaveRequest(ctx context.Context, schoolID, id string) (*staff_attendance.LeaveRequest, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if l, ok := m.leaves[id]; ok {
		return l, nil
	}
	return nil, errors.New("leave not found")
}

func (m *mockStaffLifecycleRepo) ApproveLeaveRequest(ctx context.Context, schoolID, id, approvedBy, notes string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.leaves[id]
	if !ok {
		return errors.New("leave not found")
	}
	l.Status = staff_attendance.LeaveStatusApproved
	l.ApprovedBy = &approvedBy
	now := time.Now()
	l.ApprovedAt = &now
	l.Notes = notes
	return nil
}

func (m *mockStaffLifecycleRepo) RejectLeaveRequest(ctx context.Context, schoolID, id, rejectedBy, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.leaves[id]
	if !ok {
		return errors.New("leave not found")
	}
	l.Status = staff_attendance.LeaveStatusRejected
	l.ApprovedBy = &rejectedBy
	now := time.Now()
	l.RejectedAt = &now
	l.RejectReason = reason
	return nil
}

func (m *mockStaffLifecycleRepo) DeleteLeaveRequest(ctx context.Context, schoolID, id, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.leaves, id)
	return nil
}

func (m *mockStaffLifecycleRepo) GetMonthlyTimecard(ctx context.Context, schoolID, userID, month string) (*staff_attendance.MonthlyTimecard, error) {
	return &staff_attendance.MonthlyTimecard{
		UserID:        userID,
		Month:         month,
		WorkedHours:   36,
		ContractHours: 36,
		LeaveDays:     0,
	}, nil
}

func (m *mockStaffLifecycleRepo) GetAllMonthlyTimecards(ctx context.Context, schoolID, month string) ([]staff_attendance.MonthlyTimecard, error) {
	return []staff_attendance.MonthlyTimecard{
		{UserID: "ata-1", Month: month, WorkedHours: 36},
	}, nil
}

var _ staff_attendance.Repository = (*mockStaffLifecycleRepo)(nil)

// In-memory mock for visitors
type mockVisitorLifecycleRepo struct {
	mu           sync.Mutex
	visitors     map[string]*visitors.Visitor
	earlyExits   map[string]*visitors.EarlyExit
	maintenances map[string]*visitors.MaintenanceReport
}

func newMockVisitorLifecycleRepo() *mockVisitorLifecycleRepo {
	return &mockVisitorLifecycleRepo{
		visitors:     make(map[string]*visitors.Visitor),
		earlyExits:   make(map[string]*visitors.EarlyExit),
		maintenances: make(map[string]*visitors.MaintenanceReport),
	}
}

func (m *mockVisitorLifecycleRepo) CreateVisitor(ctx context.Context, v *visitors.Visitor) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	v.ID = fmt.Sprintf("vis-%d", len(m.visitors)+1)
	v.EntryTime = time.Now()
	m.visitors[v.ID] = v
	return nil
}

func (m *mockVisitorLifecycleRepo) RecordVisitorExit(ctx context.Context, id, schoolID string, notes string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.visitors[id]
	if !ok {
		return errors.New("visitor not found")
	}
	now := time.Now()
	v.ExitTime = &now
	v.Notes = notes
	return nil
}

func (m *mockVisitorLifecycleRepo) ListTodayVisitors(ctx context.Context, schoolID, date string) ([]*visitors.Visitor, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []*visitors.Visitor
	for _, v := range m.visitors {
		list = append(list, v)
	}
	return list, nil
}

func (m *mockVisitorLifecycleRepo) GetVisitor(ctx context.Context, id, schoolID string) (*visitors.Visitor, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.visitors[id]
	if !ok {
		return nil, errors.New("visitor not found")
	}
	return v, nil
}

func (m *mockVisitorLifecycleRepo) CreateEarlyExit(ctx context.Context, e *visitors.EarlyExit) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	e.ID = fmt.Sprintf("exit-%d", len(m.earlyExits)+1)
	e.ExitTime = time.Now()
	m.earlyExits[e.ID] = e
	return nil
}

func (m *mockVisitorLifecycleRepo) RecordStudentReturn(ctx context.Context, id, schoolID, notes string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.earlyExits[id]
	if !ok {
		return errors.New("early exit not found")
	}
	now := time.Now()
	e.ReturnTime = &now
	return nil
}

func (m *mockVisitorLifecycleRepo) ListTodayEarlyExits(ctx context.Context, schoolID, date string) ([]*visitors.EarlyExit, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []*visitors.EarlyExit
	for _, e := range m.earlyExits {
		list = append(list, e)
	}
	return list, nil
}

func (m *mockVisitorLifecycleRepo) CreateMaintenanceReport(ctx context.Context, r *visitors.MaintenanceReport) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	r.ID = fmt.Sprintf("maint-%d", len(m.maintenances)+1)
	m.maintenances[r.ID] = r
	return nil
}

func (m *mockVisitorLifecycleRepo) ListMaintenanceReports(ctx context.Context, schoolID, statusFilter string) ([]*visitors.MaintenanceReport, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []*visitors.MaintenanceReport
	for _, r := range m.maintenances {
		list = append(list, r)
	}
	return list, nil
}

func (m *mockVisitorLifecycleRepo) UpdateMaintenanceStatus(ctx context.Context, id, schoolID string, req visitors.UpdateMaintenanceStatusRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	r, ok := m.maintenances[id]
	if !ok {
		return errors.New("maintenance not found")
	}
	r.Status = req.Status
	return nil
}

var _ visitors.Repository = (*mockVisitorLifecycleRepo)(nil)

// In-memory mock for personnel desk
type mockPersonnelDeskLifecycleRepo struct {
	mu       sync.Mutex
	requests map[string]*personnel_desk.DeskRequest
}

func newMockPersonnelDeskLifecycleRepo() *mockPersonnelDeskLifecycleRepo {
	return &mockPersonnelDeskLifecycleRepo{
		requests: make(map[string]*personnel_desk.DeskRequest),
	}
}

func (m *mockPersonnelDeskLifecycleRepo) Create(ctx context.Context, req *personnel_desk.DeskRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	req.ID = fmt.Sprintf("req-%d", len(m.requests)+1)
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	m.requests[req.ID] = req
	return nil
}

func (m *mockPersonnelDeskLifecycleRepo) GetByID(ctx context.Context, schoolID, id string) (*personnel_desk.DeskRequest, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	r, ok := m.requests[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return r, nil
}

func (m *mockPersonnelDeskLifecycleRepo) List(ctx context.Context, schoolID, applicantID, status string) ([]*personnel_desk.DeskRequest, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []*personnel_desk.DeskRequest
	for _, r := range m.requests {
		if applicantID != "" && r.ApplicantID != applicantID {
			continue
		}
		if status != "" && string(r.Status) != status {
			continue
		}
		list = append(list, r)
	}
	return list, nil
}

func (m *mockPersonnelDeskLifecycleRepo) Update(ctx context.Context, req *personnel_desk.DeskRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	req.UpdatedAt = time.Now()
	m.requests[req.ID] = req
	return nil
}

func (m *mockPersonnelDeskLifecycleRepo) Delete(ctx context.Context, schoolID, id, applicantID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.requests, id)
	return nil
}

var _ personnel_desk.Repository = (*mockPersonnelDeskLifecycleRepo)(nil)

// ==========================================
// INTEGRATION TESTS
// ==========================================

func TestIntegration_StaffAndATALifecycle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 1. Staff Attendance and Badging Workflow
	t.Run("Staff_BadgeSwipe_And_Process_Lifecycle", func(t *testing.T) {
		repo := newMockStaffLifecycleRepo()
		svc := staff_attendance.NewService(repo)
		handler := staff_attendance.NewHandler(svc)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("user_id", "admin-1")
			c.Set("role", "admin")
			c.Set("school_id", "school-1")
			c.Next()
		})
		api := r.Group("/api/v1")
		handler.RegisterRoutes(api)

		// Assign a badge to ATA staff
		assignBody := staff_attendance.UserBadge{
			UserID:    "ata-user-1",
			SchoolID:  "school-1",
			BadgeCode: "BADGE-101",
		}
		b1, _ := json.Marshal(assignBody)
		req1, _ := http.NewRequest(http.MethodPost, "/api/v1/staff-attendance/badges", bytes.NewBuffer(b1))
		req1.Header.Set("Content-Type", "application/json")
		w1 := httptest.NewRecorder()
		r.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusCreated, w1.Code)

		// Record morning Check-in swipe
		inSwipe := staff_attendance.BadgeSwipeRequest{
			BadgeCode: "BADGE-101",
			DeviceID:  "DEV-01",
			SwipeType: "in",
		}
		b2, _ := json.Marshal(inSwipe)
		req2, _ := http.NewRequest(http.MethodPost, "/api/v1/staff-attendance/badge-swipe", bytes.NewBuffer(b2))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusCreated, w2.Code)

		// Record evening Check-out swipe
		outSwipe := staff_attendance.BadgeSwipeRequest{
			BadgeCode: "BADGE-101",
			DeviceID:  "DEV-01",
			SwipeType: "out",
		}
		b3, _ := json.Marshal(outSwipe)
		req3, _ := http.NewRequest(http.MethodPost, "/api/v1/staff-attendance/badge-swipe", bytes.NewBuffer(b3))
		req3.Header.Set("Content-Type", "application/json")
		w3 := httptest.NewRecorder()
		r.ServeHTTP(w3, req3)
		assert.Equal(t, http.StatusCreated, w3.Code)

		// Process pending swipes
		req4, _ := http.NewRequest(http.MethodPost, "/api/v1/staff-attendance/badge-swipe/process", nil)
		w4 := httptest.NewRecorder()
		r.ServeHTTP(w4, req4)
		assert.Equal(t, http.StatusOK, w4.Code)

		var procResp map[string]interface{}
		_ = json.Unmarshal(w4.Body.Bytes(), &procResp)
		assert.Equal(t, float64(2), procResp["processed"])
	})

	// 2. ATA Leave Request Approval Workflow
	t.Run("ATA_Leave_Request_Approve_And_Reject_Workflow", func(t *testing.T) {
		repo := newMockStaffLifecycleRepo()
		leaveHandler := staff_attendance.NewLeaveHandler(repo)

		router := gin.New()
		api := router.Group("/api/v1")

		// Create leave by ATA staff
		createReq := staff_attendance.CreateLeaveRequest{
			Type:      staff_attendance.LeaveTypeFerie,
			StartDate: "2026-04-01",
			EndDate:   "2026-04-05",
			Days:      5,
			Notes:     "Ferie pasquali programmate",
		}
		b, _ := json.Marshal(createReq)

		req1, _ := http.NewRequest(http.MethodPost, "/api/v1/staff-attendance/leaves", bytes.NewBuffer(b))
		req1.Header.Set("Content-Type", "application/json")

		// Serve with ATA user context
		subRouter1 := gin.New()
		subRouter1.Use(func(c *gin.Context) {
			c.Set("user_id", "ata-user-1")
			c.Set("role", "collaboratore_scolastico")
			c.Set("school_id", "school-1")
			c.Next()
		})
		leaveHandler.RegisterLeaveRoutes(subRouter1.Group("/api/v1"))
		w1 := httptest.NewRecorder()
		subRouter1.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusCreated, w1.Code)

		var createdLeave staff_attendance.LeaveRequest
		_ = json.Unmarshal(w1.Body.Bytes(), &createdLeave)
		assert.Equal(t, "leave-1", createdLeave.ID)
		assert.Equal(t, staff_attendance.LeaveStatusPending, createdLeave.Status)

		// DSGA reviews and approves leave
		approveBody := map[string]string{"notes": "Ferie approvate compatibilmente con turni"}
		ab, _ := json.Marshal(approveBody)
		req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/staff-attendance/leaves/%s/approve", createdLeave.ID), bytes.NewBuffer(ab))
		req2.Header.Set("Content-Type", "application/json")

		subRouter2 := gin.New()
		subRouter2.Use(func(c *gin.Context) {
			c.Set("user_id", "dsga-1")
			c.Set("role", "dsga")
			c.Set("school_id", "school-1")
			c.Next()
		})
		leaveHandler.RegisterLeaveRoutes(subRouter2.Group("/api/v1"))
		w2 := httptest.NewRecorder()
		subRouter2.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		// Verify state is approved
		getReq, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/staff-attendance/leaves/%s", createdLeave.ID), nil)
		w3 := httptest.NewRecorder()
		subRouter1.ServeHTTP(w3, getReq)
		assert.Equal(t, http.StatusOK, w3.Code)
		var approvedLeave staff_attendance.LeaveRequest
		_ = json.Unmarshal(w3.Body.Bytes(), &approvedLeave)
		assert.Equal(t, staff_attendance.LeaveStatusApproved, approvedLeave.Status)
		assert.Equal(t, "dsga-1", *approvedLeave.ApprovedBy)
		_ = router
		_ = api
	})

	// 3. Visitor Check-in, Badge Allocation and Check-out
	t.Run("Visitor_Lifecycle_CheckIn_And_CheckOut", func(t *testing.T) {
		repo := newMockVisitorLifecycleRepo()
		svc := visitors.NewService(repo)
		h := visitors.NewHandler(svc)

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("user_id", "ata-porter-1")
			c.Set("role", "collaboratore_scolastico")
			c.Set("is_staff", true)
			c.Set("school_id", "school-1")
			c.Next()
		})
		api := r.Group("/api/v1")
		h.RegisterRoutes(api)

		// 1. Check in visitor
		visReq := visitors.RegisterVisitorRequest{
			Name:        "Mario Rossi (Tecnico Stampanti)",
			DocumentID:  "CA123456",
			Purpose:     visitors.PurposeSupplier,
			HostName:    "Segreteria",
			BadgeNumber: "VIS-42",
		}
		vb, _ := json.Marshal(visReq)
		req1, _ := http.NewRequest(http.MethodPost, "/api/v1/visitors", bytes.NewBuffer(vb))
		req1.Header.Set("Content-Type", "application/json")
		w1 := httptest.NewRecorder()
		r.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusCreated, w1.Code)

		var visResp visitors.Visitor
		_ = json.Unmarshal(w1.Body.Bytes(), &visResp)
		assert.Equal(t, "Mario Rossi (Tecnico Stampanti)", visResp.Name)
		assert.Equal(t, "VIS-42", *visResp.BadgeNumber)
		assert.Nil(t, visResp.ExitTime)

		// 2. Check out visitor
		exitReq := map[string]string{"notes": "Riconsegnato badge integro"}
		eb, _ := json.Marshal(exitReq)
		req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/visitors/%s/exit", visResp.ID), bytes.NewBuffer(eb))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		// Verify exit is recorded
		assert.NotNil(t, repo.visitors[visResp.ID].ExitTime)
		assert.Equal(t, "Riconsegnato badge integro", repo.visitors[visResp.ID].Notes)
	})

	// 4. Personnel Desk Service Ticket Lifecycle
	t.Run("Personnel_Desk_Ticket_Full_Lifecycle", func(t *testing.T) {
		repo := newMockPersonnelDeskLifecycleRepo()
		svc := personnel_desk.NewService(repo)
		h := personnel_desk.NewHandler(svc)

		// Employee creates request as draft then submits
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("user_id", "emp-1")
			c.Set("role", "teacher")
			c.Set("school_id", "school-1")
			c.Next()
		})
		api := r.Group("/api/v1")
		h.RegisterRoutes(api)

		input := personnel_desk.CreateDeskRequestInput{
			Category:    personnel_desk.CatPermessoBreve,
			StartDate:   "2026-05-01",
			EndDate:     "2026-05-01",
			Hours:       2,
			Description: "Visita specialistica programmata",
			SubmitNow:   false,
		}
		ib, _ := json.Marshal(input)
		req1, _ := http.NewRequest(http.MethodPost, "/api/v1/personnel-desk/requests", bytes.NewBuffer(ib))
		req1.Header.Set("Content-Type", "application/json")
		w1 := httptest.NewRecorder()
		r.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusCreated, w1.Code)

		var ticket personnel_desk.DeskRequest
		_ = json.Unmarshal(w1.Body.Bytes(), &ticket)
		assert.Equal(t, personnel_desk.StatusDraft, ticket.Status)

		// Submit ticket
		req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/personnel-desk/requests/%s/submit", ticket.ID), nil)
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		// AA reviews ticket
		rAdmin := gin.New()
		rAdmin.Use(func(c *gin.Context) {
			c.Set("user_id", "aa-user-1")
			c.Set("role", "assistente_amministrativo")
			c.Set("school_id", "school-1")
			c.Next()
		})
		h.RegisterRoutes(rAdmin.Group("/api/v1"))

		aaBody := personnel_desk.AAReviewInput{
			Approve: true,
			Note:    "Verifica fascicolo personale completata con esito positivo",
		}
		aab, _ := json.Marshal(aaBody)
		req3, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/personnel-desk/requests/%s/aa-review", ticket.ID), bytes.NewBuffer(aab))
		req3.Header.Set("Content-Type", "application/json")
		w3 := httptest.NewRecorder()
		rAdmin.ServeHTTP(w3, req3)
		assert.Equal(t, http.StatusOK, w3.Code)

		// DSGA signs
		rDSGA := gin.New()
		rDSGA.Use(func(c *gin.Context) {
			c.Set("user_id", "dsga-1")
			c.Set("role", "dsga")
			c.Set("school_id", "school-1")
			c.Next()
		})
		h.RegisterRoutes(rDSGA.Group("/api/v1"))

		dsgaBody := personnel_desk.DSGASignInput{
			Approve: true,
			Note:    "Visto di regolarità contabile e amministrativa",
		}
		dsgab, _ := json.Marshal(dsgaBody)
		req4, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/personnel-desk/requests/%s/dsga-sign", ticket.ID), bytes.NewBuffer(dsgab))
		req4.Header.Set("Content-Type", "application/json")
		w4 := httptest.NewRecorder()
		rDSGA.ServeHTTP(w4, req4)
		assert.Equal(t, http.StatusOK, w4.Code)

		// Principal (DS) approves
		rDS := gin.New()
		rDS.Use(func(c *gin.Context) {
			c.Set("user_id", "principal-1")
			c.Set("role", "principal")
			c.Set("school_id", "school-1")
			c.Next()
		})
		h.RegisterRoutes(rDS.Group("/api/v1"))

		dsBody := personnel_desk.DSApproveInput{
			Approve:   true,
			DecreeNum: "DEC-2026-01",
			Note:      "Certificato concesso e firmato digitalmente",
		}
		dsb, _ := json.Marshal(dsBody)
		req5, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/personnel-desk/requests/%s/ds-approve", ticket.ID), bytes.NewBuffer(dsb))
		req5.Header.Set("Content-Type", "application/json")
		w5 := httptest.NewRecorder()
		rDS.ServeHTTP(w5, req5)
		assert.Equal(t, http.StatusOK, w5.Code)

		var finalTicket personnel_desk.DeskRequest
		_ = json.Unmarshal(w5.Body.Bytes(), &finalTicket)
		assert.Equal(t, personnel_desk.StatusApproved, finalTicket.Status)
	})
}
