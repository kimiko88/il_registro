package scheduling

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSchedRepo struct {
	mock.Mock
}

func (m *MockSchedRepo) CreateSlot(ctx context.Context, s *ColloquioSlot) error {
	return m.Called(s).Error(0)
}
func (m *MockSchedRepo) GetSlots(ctx context.Context, tid string, f, t time.Time) ([]ColloquioSlot, error) {
	return nil, nil
}
func (m *MockSchedRepo) GetSlotByID(ctx context.Context, id string) (*ColloquioSlot, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ColloquioSlot), args.Error(1)
}
func (m *MockSchedRepo) UpdateSlot(ctx context.Context, s *ColloquioSlot) error { return nil }
func (m *MockSchedRepo) GetAvailableSlots(ctx context.Context, sid, tid string, f, t time.Time) ([]ColloquioSlot, error) {
	return nil, nil
}
func (m *MockSchedRepo) GetSettings(ctx context.Context, sid string) (*ColloquioSettings, error) {
	// Return valid defaults to pass validation
	return &ColloquioSettings{
		SchoolID:           sid,
		BookingWindowDays:  14,
		BookingBufferHours: 24,
	}, nil
}
func (m *MockSchedRepo) UpdateSettings(ctx context.Context, s *ColloquioSettings) error { return nil }
func (m *MockSchedRepo) CountBookingsForParent(ctx context.Context, pid string, d time.Time, s, e time.Time) (int, error) {
	return m.Called(pid).Int(0), nil
}
func (m *MockSchedRepo) CreateBooking(ctx context.Context, b *ColloquioBooking) error {
	return m.Called(b).Error(0)
}
func (m *MockSchedRepo) GetBookingsByParent(ctx context.Context, pid string) ([]ColloquioBooking, error) {
	return nil, nil
}
func (m *MockSchedRepo) GetBookingsByTeacher(ctx context.Context, tid string) ([]ColloquioBooking, error) {
	return nil, nil
}
func (m *MockSchedRepo) GetBooking(ctx context.Context, id string) (*ColloquioBooking, error) {
	return nil, nil
}
func (m *MockSchedRepo) UpdateBooking(ctx context.Context, b *ColloquioBooking) error { return nil }
func (m *MockSchedRepo) GetAnalytics(ctx context.Context, schoolID string) (*AnalyticsResponse, error) {
	return nil, nil
}

func TestRegression_DoubleBooking(t *testing.T) {
	repo := &MockSchedRepo{}
	svc := NewService(repo)
	ctx := context.Background()

	// 2 days in future -> Inside 14d window, Outside 24h buffer
	future := time.Now().AddDate(0, 0, 2)

	slot := &ColloquioSlot{
		ID:        "slot-1",
		SchoolID:  "school-1",
		Date:      future,
		StartTime: future,
		EndTime:   future.Add(15 * time.Minute),
	}

	repo.On("GetSlotByID", "slot-1").Return(slot, nil)
	repo.On("CountBookingsForParent", "parent-1").Return(1)

	sid := "student-1"
	req := BookSlotRequest{
		SlotID:    "slot-1",
		StudentID: &sid,
	}

	_, err := svc.BookSlot(ctx, "parent-1", req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "conflict")
}
