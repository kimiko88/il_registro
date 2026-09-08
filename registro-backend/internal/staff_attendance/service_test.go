package staff_attendance_test

import (
	"context"
	"testing"
	"time"

	"registro-backend/internal/staff_attendance"
)

// MockRepo implements staff_attendance.Repository for unit tests
type MockRepo struct {
	summary      *staff_attendance.DailyStaffSummary
	summaryErr   error
	list         []staff_attendance.StaffAttendance
	listErr      error
	upsertResult *staff_attendance.StaffAttendance
	upsertErr    error
	bulkResult   []staff_attendance.StaffAttendance
	bulkErr      error
	deleteErr    error
	swipeResult  *staff_attendance.BadgeSwipe
	swipeErr     error
	processCount int
	processErr   error
	assignErr    error
}

func (m *MockRepo) GetByUserAndDate(ctx context.Context, schoolID, userID, date string) (*staff_attendance.StaffAttendance, error) {
	return m.upsertResult, m.upsertErr
}

func (m *MockRepo) GetDailySummary(ctx context.Context, schoolID, date string) (*staff_attendance.DailyStaffSummary, error) {
	return m.summary, m.summaryErr
}

func (m *MockRepo) ListByDate(ctx context.Context, schoolID, date string) ([]staff_attendance.StaffAttendance, error) {
	return m.list, m.listErr
}

func (m *MockRepo) Upsert(ctx context.Context, schoolID, createdBy string, req staff_attendance.UpsertStaffAttendanceRequest) (*staff_attendance.StaffAttendance, error) {
	return m.upsertResult, m.upsertErr
}

func (m *MockRepo) BulkUpsert(ctx context.Context, schoolID, createdBy string, req staff_attendance.BulkUpsertRequest) ([]staff_attendance.StaffAttendance, error) {
	return m.bulkResult, m.bulkErr
}

func (m *MockRepo) Delete(ctx context.Context, schoolID, id string) error {
	return m.deleteErr
}

func (m *MockRepo) RegisterBadgeSwipe(ctx context.Context, schoolID string, req staff_attendance.BadgeSwipeRequest) (*staff_attendance.BadgeSwipe, error) {
	return m.swipeResult, m.swipeErr
}

func (m *MockRepo) ProcessPendingSwipes(ctx context.Context, schoolID string) (int, error) {
	return m.processCount, m.processErr
}

func (m *MockRepo) AssignBadge(ctx context.Context, badge staff_attendance.UserBadge) error {
	return m.assignErr
}

func (m *MockRepo) RevokeBadge(ctx context.Context, schoolID, userID, badgeCode string) error {
	return nil
}

func (m *MockRepo) ListBadges(ctx context.Context, schoolID string) ([]staff_attendance.UserBadge, error) {
	return nil, nil
}

func TestATARolesAndPermissions(t *testing.T) {
	roles := []string{"dsga", "assistente_amministrativo", "collaboratore_ds", "collaboratore_scolastico"}
	for _, r := range roles {
		if !staff_attendance.IsATARole(r) {
			t.Errorf("expected %s to be an ATA role", r)
		}
		if !staff_attendance.CanReadAttendance(r) {
			t.Errorf("expected %s to be able to read staff attendance", r)
		}
		if !staff_attendance.CanWriteAttendance(r) {
			t.Errorf("expected %s to be able to write staff attendance", r)
		}
	}

	// Student and parent must not have access
	if staff_attendance.CanReadAttendance("student") {
		t.Error("student should not be allowed to read staff attendance")
	}
	if staff_attendance.CanWriteAttendance("parent") {
		t.Error("parent should not be allowed to write staff attendance")
	}
	if staff_attendance.CanReadAttendance("teacher") {
		t.Error("regular teacher should not be allowed to read staff attendance dashboard")
	}
}

func TestService_GetDailySummary_Permissions(t *testing.T) {
	mock := &MockRepo{
		summary: &staff_attendance.DailyStaffSummary{
			Date: "2026-09-08",
		},
	}
	svc := staff_attendance.NewService(mock)
	ctx := context.Background()

	// Authorized role: dsga
	res, err := svc.GetDailySummary(ctx, "dsga", "school-1", "2026-09-08")
	if err != nil {
		t.Fatalf("unexpected error for dsga: %v", err)
	}
	if res == nil || res.Date != "2026-09-08" {
		t.Fatalf("unexpected response: %+v", res)
	}

	// Authorized role: collaboratore_scolastico
	res, err = svc.GetDailySummary(ctx, "collaboratore_scolastico", "school-1", "2026-09-08")
	if err != nil {
		t.Fatalf("unexpected error for collaboratore_scolastico: %v", err)
	}

	// Unauthorized role: student
	_, err = svc.GetDailySummary(ctx, "student", "school-1", "2026-09-08")
	if err == nil {
		t.Fatal("expected error for student role, got nil")
	}
}

func TestService_RecordAttendance_Validation(t *testing.T) {
	mock := &MockRepo{
		upsertResult: &staff_attendance.StaffAttendance{
			ID:     "att-1",
			UserID: "u-1",
			Date:   "2026-09-08",
			Status: staff_attendance.StatusPresent,
		},
	}
	svc := staff_attendance.NewService(mock)
	ctx := context.Background()

	// Unauthorized
	_, err := svc.RecordAttendance(ctx, "student", "actor-1", "school-1", staff_attendance.UpsertStaffAttendanceRequest{
		UserID: "u-1",
		Date:   "2026-09-08",
		Status: staff_attendance.StatusPresent,
	})
	if err == nil {
		t.Error("expected error for student, got nil")
	}

	// Missing user_id
	_, err = svc.RecordAttendance(ctx, "collaboratore_ds", "actor-1", "school-1", staff_attendance.UpsertStaffAttendanceRequest{
		UserID: "",
		Date:   "2026-09-08",
		Status: staff_attendance.StatusPresent,
	})
	if err == nil {
		t.Error("expected error for missing user_id")
	}

	// Invalid date
	_, err = svc.RecordAttendance(ctx, "collaboratore_ds", "actor-1", "school-1", staff_attendance.UpsertStaffAttendanceRequest{
		UserID: "u-1",
		Date:   "not-a-date",
		Status: staff_attendance.StatusPresent,
	})
	if err == nil {
		t.Error("expected error for invalid date")
	}

	// Invalid status
	_, err = svc.RecordAttendance(ctx, "collaboratore_ds", "actor-1", "school-1", staff_attendance.UpsertStaffAttendanceRequest{
		UserID: "u-1",
		Date:   "2026-09-08",
		Status: staff_attendance.AttendanceStatus("invalid_status"),
	})
	if err == nil {
		t.Error("expected error for invalid status")
	}

	// Success
	res, err := svc.RecordAttendance(ctx, "collaboratore_ds", "actor-1", "school-1", staff_attendance.UpsertStaffAttendanceRequest{
		UserID: "u-1",
		Date:   "2026-09-08",
		Status: staff_attendance.StatusPresent,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != "att-1" {
		t.Errorf("expected att-1, got %s", res.ID)
	}
}

func TestService_BadgeSwipe_Validation(t *testing.T) {
	mock := &MockRepo{
		swipeResult: &staff_attendance.BadgeSwipe{
			ID:        "swipe-1",
			BadgeCode: "BDG12345",
			SwipeTime: time.Now(),
		},
	}
	svc := staff_attendance.NewService(mock)
	ctx := context.Background()

	// Missing badge_code
	_, err := svc.RegisterBadgeSwipe(ctx, "school-1", staff_attendance.BadgeSwipeRequest{
		DeviceID: "DEV-1",
	})
	if err == nil {
		t.Error("expected error for missing badge_code")
	}

	// Missing device_id
	_, err = svc.RegisterBadgeSwipe(ctx, "school-1", staff_attendance.BadgeSwipeRequest{
		BadgeCode: "BDG12345",
	})
	if err == nil {
		t.Error("expected error for missing device_id")
	}

	// Success
	res, err := svc.RegisterBadgeSwipe(ctx, "school-1", staff_attendance.BadgeSwipeRequest{
		BadgeCode: "BDG12345",
		DeviceID:  "DEV-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != "swipe-1" {
		t.Errorf("expected swipe-1, got %s", res.ID)
	}
}
