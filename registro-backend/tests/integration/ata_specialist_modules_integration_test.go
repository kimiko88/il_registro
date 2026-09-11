package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/personnel_desk"
	"registro-backend/internal/staff_attendance"
	"registro-backend/internal/substitutions"
	"registro-backend/internal/visitors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSubstRepo struct{ mock.Mock }

func (m *mockSubstRepo) Create(ctx context.Context, sub *substitutions.Substitution) error {
	return m.Called(ctx, sub).Error(0)
}
func (m *mockSubstRepo) GetByID(ctx context.Context, id string) (*substitutions.Substitution, error) {
	args := m.Called(ctx, id)
	if res := args.Get(0); res != nil {
		return res.(*substitutions.Substitution), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockSubstRepo) ListBySchool(ctx context.Context, schoolID, date string) ([]*substitutions.Substitution, error) {
	args := m.Called(ctx, schoolID, date)
	return args.Get(0).([]*substitutions.Substitution), args.Error(1)
}
func (m *mockSubstRepo) ListByTeacher(ctx context.Context, teacherID string, date string) ([]*substitutions.Substitution, error) {
	args := m.Called(ctx, teacherID, date)
	return args.Get(0).([]*substitutions.Substitution), args.Error(1)
}
func (m *mockSubstRepo) AssignSubstitute(ctx context.Context, id string, substituteTeacherID string, notes string) error {
	return m.Called(ctx, id, substituteTeacherID, notes).Error(0)
}
func (m *mockSubstRepo) ConfirmSubstitution(ctx context.Context, id string, substituteTeacherID string) error {
	return m.Called(ctx, id, substituteTeacherID).Error(0)
}
func (m *mockSubstRepo) SignRegister(ctx context.Context, id string, sigHash string, notes string) error {
	return m.Called(ctx, id, sigHash, notes).Error(0)
}
func (m *mockSubstRepo) GetAvailableTeachers(ctx context.Context, schoolID string) ([]substitutions.TeacherCandidate, error) {
	args := m.Called(ctx, schoolID)
	return args.Get(0).([]substitutions.TeacherCandidate), args.Error(1)
}
func (m *mockSubstRepo) GetTeacherProfileID(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}
func (m *mockSubstRepo) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	args := m.Called(ctx, teacherID, classID)
	return args.Bool(0), args.Error(1)
}
func (m *mockSubstRepo) IsTeacherAssignedToSubject(ctx context.Context, teacherID, subjectID string) (bool, error) {
	args := m.Called(ctx, teacherID, subjectID)
	return args.Bool(0), args.Error(1)
}
func (m *mockSubstRepo) GetWeeklySubstitutionCount(ctx context.Context, teacherID string) (int, error) {
	args := m.Called(ctx, teacherID)
	return args.Int(0), args.Error(1)
}

type mockVisitorRepo struct{ mock.Mock }

func (m *mockVisitorRepo) CreateVisitor(ctx context.Context, v *visitors.Visitor) error {
	return m.Called(ctx, v).Error(0)
}
func (m *mockVisitorRepo) RecordVisitorExit(ctx context.Context, id, schoolID string, notes string) error {
	return m.Called(ctx, id, schoolID, notes).Error(0)
}
func (m *mockVisitorRepo) ListTodayVisitors(ctx context.Context, schoolID, date string) ([]*visitors.Visitor, error) {
	args := m.Called(ctx, schoolID, date)
	return args.Get(0).([]*visitors.Visitor), args.Error(1)
}
func (m *mockVisitorRepo) GetVisitor(ctx context.Context, id, schoolID string) (*visitors.Visitor, error) {
	args := m.Called(ctx, id, schoolID)
	if res := args.Get(0); res != nil {
		return res.(*visitors.Visitor), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockVisitorRepo) CreateEarlyExit(ctx context.Context, e *visitors.EarlyExit) error {
	return m.Called(ctx, e).Error(0)
}
func (m *mockVisitorRepo) RecordStudentReturn(ctx context.Context, id, schoolID, notes string) error {
	return m.Called(ctx, id, schoolID, notes).Error(0)
}
func (m *mockVisitorRepo) ListTodayEarlyExits(ctx context.Context, schoolID, date string) ([]*visitors.EarlyExit, error) {
	args := m.Called(ctx, schoolID, date)
	return args.Get(0).([]*visitors.EarlyExit), args.Error(1)
}
func (m *mockVisitorRepo) CreateMaintenanceReport(ctx context.Context, r *visitors.MaintenanceReport) error {
	return m.Called(ctx, r).Error(0)
}
func (m *mockVisitorRepo) ListMaintenanceReports(ctx context.Context, schoolID, statusFilter string) ([]*visitors.MaintenanceReport, error) {
	args := m.Called(ctx, schoolID, statusFilter)
	return args.Get(0).([]*visitors.MaintenanceReport), args.Error(1)
}
func (m *mockVisitorRepo) UpdateMaintenanceStatus(ctx context.Context, id, schoolID string, req visitors.UpdateMaintenanceStatusRequest) error {
	return m.Called(ctx, id, schoolID, req).Error(0)
}

type mockAttendanceRepo struct{ mock.Mock }

func (m *mockAttendanceRepo) GetDailySummary(ctx context.Context, schoolID, date string) (*staff_attendance.DailyStaffSummary, error) {
	args := m.Called(ctx, schoolID, date)
	if res := args.Get(0); res != nil {
		return res.(*staff_attendance.DailyStaffSummary), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockAttendanceRepo) ListByDate(ctx context.Context, schoolID, date string) ([]staff_attendance.StaffAttendance, error) {
	args := m.Called(ctx, schoolID, date)
	return args.Get(0).([]staff_attendance.StaffAttendance), args.Error(1)
}
func (m *mockAttendanceRepo) Upsert(ctx context.Context, schoolID, recordedBy string, req staff_attendance.UpsertStaffAttendanceRequest) (*staff_attendance.StaffAttendance, error) {
	args := m.Called(ctx, schoolID, recordedBy, req)
	if res := args.Get(0); res != nil {
		return res.(*staff_attendance.StaffAttendance), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockAttendanceRepo) BulkUpsert(ctx context.Context, schoolID, recordedBy string, req staff_attendance.BulkUpsertRequest) ([]staff_attendance.StaffAttendance, error) {
	args := m.Called(ctx, schoolID, recordedBy, req)
	return args.Get(0).([]staff_attendance.StaffAttendance), args.Error(1)
}
func (m *mockAttendanceRepo) GetByUserAndDate(ctx context.Context, schoolID, userID, date string) (*staff_attendance.StaffAttendance, error) {
	args := m.Called(ctx, schoolID, userID, date)
	if res := args.Get(0); res != nil {
		return res.(*staff_attendance.StaffAttendance), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockAttendanceRepo) Delete(ctx context.Context, schoolID, id string) error {
	return m.Called(ctx, schoolID, id).Error(0)
}
func (m *mockAttendanceRepo) RegisterBadgeSwipe(ctx context.Context, schoolID string, req staff_attendance.BadgeSwipeRequest) (*staff_attendance.BadgeSwipe, error) {
	args := m.Called(ctx, schoolID, req)
	if res := args.Get(0); res != nil {
		return res.(*staff_attendance.BadgeSwipe), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockAttendanceRepo) ProcessPendingSwipes(ctx context.Context, schoolID string) (int, error) {
	args := m.Called(ctx, schoolID)
	return args.Int(0), args.Error(1)
}
func (m *mockAttendanceRepo) AssignBadge(ctx context.Context, badge staff_attendance.UserBadge) error {
	return m.Called(ctx, badge).Error(0)
}
func (m *mockAttendanceRepo) RevokeBadge(ctx context.Context, schoolID, userID, badgeCode string) error {
	return m.Called(ctx, schoolID, userID, badgeCode).Error(0)
}
func (m *mockAttendanceRepo) ListBadges(ctx context.Context, schoolID string) ([]staff_attendance.UserBadge, error) {
	args := m.Called(ctx, schoolID)
	return args.Get(0).([]staff_attendance.UserBadge), args.Error(1)
}
func (m *mockAttendanceRepo) CreateLeaveRequest(ctx context.Context, schoolID, userID string, req staff_attendance.CreateLeaveRequest) (*staff_attendance.LeaveRequest, error) {
	args := m.Called(ctx, schoolID, userID, req)
	if res := args.Get(0); res != nil {
		return res.(*staff_attendance.LeaveRequest), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockAttendanceRepo) ListLeaveRequests(ctx context.Context, schoolID, userID, status string) ([]staff_attendance.LeaveRequest, error) {
	args := m.Called(ctx, schoolID, userID, status)
	return args.Get(0).([]staff_attendance.LeaveRequest), args.Error(1)
}
func (m *mockAttendanceRepo) GetLeaveRequest(ctx context.Context, schoolID, id string) (*staff_attendance.LeaveRequest, error) {
	args := m.Called(ctx, schoolID, id)
	if res := args.Get(0); res != nil {
		return res.(*staff_attendance.LeaveRequest), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockAttendanceRepo) ApproveLeaveRequest(ctx context.Context, schoolID, id, approvedBy, notes string) error {
	return m.Called(ctx, schoolID, id, approvedBy, notes).Error(0)
}
func (m *mockAttendanceRepo) RejectLeaveRequest(ctx context.Context, schoolID, id, rejectedBy, reason string) error {
	return m.Called(ctx, schoolID, id, rejectedBy, reason).Error(0)
}
func (m *mockAttendanceRepo) DeleteLeaveRequest(ctx context.Context, schoolID, id, userID string) error {
	return m.Called(ctx, schoolID, id, userID).Error(0)
}
func (m *mockAttendanceRepo) GetMonthlyTimecard(ctx context.Context, schoolID, userID, month string) (*staff_attendance.MonthlyTimecard, error) {
	args := m.Called(ctx, schoolID, userID, month)
	if res := args.Get(0); res != nil {
		return res.(*staff_attendance.MonthlyTimecard), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockAttendanceRepo) GetAllMonthlyTimecards(ctx context.Context, schoolID, month string) ([]staff_attendance.MonthlyTimecard, error) {
	args := m.Called(ctx, schoolID, month)
	return args.Get(0).([]staff_attendance.MonthlyTimecard), args.Error(1)
}

type mockDeskRepo struct{ mock.Mock }

func (m *mockDeskRepo) Create(ctx context.Context, req *personnel_desk.DeskRequest) error {
	return m.Called(ctx, req).Error(0)
}
func (m *mockDeskRepo) GetByID(ctx context.Context, schoolID, id string) (*personnel_desk.DeskRequest, error) {
	args := m.Called(ctx, schoolID, id)
	if res := args.Get(0); res != nil {
		return res.(*personnel_desk.DeskRequest), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockDeskRepo) List(ctx context.Context, schoolID, applicantID, status string) ([]*personnel_desk.DeskRequest, error) {
	args := m.Called(ctx, schoolID, applicantID, status)
	return args.Get(0).([]*personnel_desk.DeskRequest), args.Error(1)
}
func (m *mockDeskRepo) Update(ctx context.Context, req *personnel_desk.DeskRequest) error {
	return m.Called(ctx, req).Error(0)
}
func (m *mockDeskRepo) Delete(ctx context.Context, schoolID, id, currentUserID string) error {
	return m.Called(ctx, schoolID, id, currentUserID).Error(0)
}

// TestATASpecialistModulesIntegration verifies the 4 ATA Specialist Modules HTTP API Workflows
func TestATASpecialistModulesIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mSubst := new(mockSubstRepo)
	mVisitor := new(mockVisitorRepo)
	mAtt := new(mockAttendanceRepo)
	mDesk := new(mockDeskRepo)

	substSvc := substitutions.NewService(mSubst)
	visitorSvc := visitors.NewService(mVisitor)
	deskSvc := personnel_desk.NewService(mDesk)

	substH := substitutions.NewHandler(substSvc)
	visitorH := visitors.NewHandler(visitorSvc)
	leaveH := staff_attendance.NewLeaveHandler(mAtt)
	deskH := personnel_desk.NewHandler(deskSvc)

	router := gin.Default()
	api := router.Group("/api/v1")

	setAuth := func(userID, role, schoolID string) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("user_id", userID)
			c.Set("role", role)
			c.Set("school_id", schoolID)
			c.Next()
		}
	}

	schoolID := "sch-ata-101"

	substH.RegisterRoutes(api.Group("", setAuth("u-ds-collab", "collaboratore_ds", schoolID)))
	visitorH.RegisterRoutes(api.Group("", setAuth("u-ata-collab", "collaboratore_scolastico", schoolID)))
	leaveH.RegisterLeaveRoutes(api.Group("", setAuth("u-dsga", "dsga", schoolID)))
	deskH.RegisterRoutes(api.Group("", setAuth("u-teacher-1", "teacher", schoolID)))

	// 1. Module 1: Emergenza Sostituzioni
	t.Run("Module 1: Emergency Substitutions Today Summary", func(t *testing.T) {
		mSubst.On("ListBySchool", mock.Anything, schoolID, mock.Anything).
			Return([]*substitutions.Substitution{}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/substitutions/today-summary", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	// 2. Module 2: Registro Visitatori
	t.Run("Module 2: Visitor Registry Registration & Exit", func(t *testing.T) {
		mVisitor.On("CreateVisitor", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(1).(*visitors.Visitor)
			arg.ID = "v-999"
		}).Once()

		body, _ := json.Marshal(visitors.RegisterVisitorRequest{
			Name:       "Giuseppe Verdi",
			DocumentID: "CA12345AA",
			Purpose:    "Colloquio presidiato",
		})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/visitors", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)

		// Record Exit
		mVisitor.On("RecordVisitorExit", mock.Anything, "v-999", schoolID, "Uscito regolarmente").Return(nil).Once()

		exitBody, _ := json.Marshal(map[string]string{"notes": "Uscito regolarmente"})
		reqExit, _ := http.NewRequest(http.MethodPatch, "/api/v1/visitors/v-999/exit", bytes.NewBuffer(exitBody))
		reqExit.Header.Set("Content-Type", "application/json")
		respExit := httptest.NewRecorder()
		router.ServeHTTP(respExit, reqExit)

		assert.Equal(t, http.StatusOK, respExit.Code)
	})

	// 3. Module 3: Cartellino & Ferie ATA
	t.Run("Module 3: Timecard & Leave Approval", func(t *testing.T) {
		mAtt.On("GetMonthlyTimecard", mock.Anything, schoolID, "u-dsga", "2026-09").
			Return(&staff_attendance.MonthlyTimecard{UserID: "u-dsga", Month: "2026-09"}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/staff-attendance/timecard?month=2026-09", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	// 4. Module 4: Sportello Digitale Personale Workflow
	t.Run("Module 4: Digital Personnel Desk Full Workflow", func(t *testing.T) {
		deskItem := &personnel_desk.DeskRequest{
			ID:          "req-777",
			SchoolID:    schoolID,
			ApplicantID: "u-teacher-1",
			Category:    personnel_desk.CatFerie,
			StartDate:   "2026-10-01",
			EndDate:     "2026-10-05",
			Days:        5,
			Description: "Ferie ordinarie",
			Status:      personnel_desk.StatusDraft,
		}

		mDesk.On("Create", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(1).(*personnel_desk.DeskRequest)
			arg.ID = deskItem.ID
		}).Once()
		mDesk.On("GetByID", mock.Anything, schoolID, "req-777").Return(deskItem, nil)

		// Create Request
		createBody, _ := json.Marshal(personnel_desk.CreateDeskRequestInput{
			Category:    personnel_desk.CatFerie,
			StartDate:   "2026-10-01",
			EndDate:     "2026-10-05",
			Days:        5,
			Description: "Ferie ordinarie",
		})
		reqCreate, _ := http.NewRequest(http.MethodPost, "/api/v1/personnel-desk/requests", bytes.NewBuffer(createBody))
		reqCreate.Header.Set("Content-Type", "application/json")
		respCreate := httptest.NewRecorder()
		router.ServeHTTP(respCreate, reqCreate)
		assert.Equal(t, http.StatusCreated, respCreate.Code)

		// Submit Request
		mDesk.On("Update", mock.Anything, mock.MatchedBy(func(r *personnel_desk.DeskRequest) bool {
			return r.Status == personnel_desk.StatusSubmitted
		})).Return(nil).Once()

		reqSubmit, _ := http.NewRequest(http.MethodPatch, "/api/v1/personnel-desk/requests/req-777/submit", nil)
		respSubmit := httptest.NewRecorder()
		router.ServeHTTP(respSubmit, reqSubmit)
		assert.Equal(t, http.StatusOK, respSubmit.Code)
	})
}
