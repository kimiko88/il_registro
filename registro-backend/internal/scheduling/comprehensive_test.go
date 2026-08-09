package scheduling

import (
	"context"
	"testing"
	"time"
	"registro-backend/internal/teachers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockRepo struct {
	mock.Mock
}

type MockTeacherRepo struct {
	mock.Mock
}

func (m *MockTeacherRepo) List(ctx context.Context, schoolID string) ([]teachers.Teacher, error) { return nil, nil }
func (m *MockTeacherRepo) Get(ctx context.Context, id string) (*teachers.Teacher, error) { return nil, nil }
func (m *MockTeacherRepo) GetByUserID(ctx context.Context, userID string) (*teachers.Teacher, error) {
	return &teachers.Teacher{ID: "teacher1", UserID: userID, SchoolID: "school1"}, nil
}
func (m *MockTeacherRepo) GetBySubject(ctx context.Context, subjectID string) ([]teachers.Teacher, error) { return nil, nil }
func (m *MockTeacherRepo) GetSubjects(ctx context.Context, teacherID string) ([]teachers.TeacherSubject, error) { return nil, nil }
func (m *MockTeacherRepo) AssignSubject(ctx context.Context, teacherID, subjectID string) error { return nil }
func (m *MockTeacherRepo) RemoveSubject(ctx context.Context, teacherID, subjectID string) error { return nil }
func (m *MockTeacherRepo) GetDashboardStats(ctx context.Context, teacherUserID string) (map[string]interface{}, error) {
	return nil, nil
}


func (m *MockRepo) CreateSlot(ctx context.Context, s *ColloquioSlot) error {
	s.ID = "new-id"
	return nil
}
func (m *MockRepo) CreateSlotsBatch(ctx context.Context, slots []ColloquioSlot) error {
	for i := range slots {
		slots[i].ID = "new-id"
	}
	return nil
}
func (m *MockRepo) DeleteSlot(ctx context.Context, id string) error { return nil }
func (m *MockRepo) GetSlots(ctx context.Context, tID string, f, t time.Time) ([]ColloquioSlot, error) {
	return []ColloquioSlot{}, nil
}
func (m *MockRepo) GetAvailableSlots(ctx context.Context, s, tID string, f, t time.Time) ([]ColloquioSlot, error) {
	return nil, nil
}
func (m *MockRepo) GetSlotByID(ctx context.Context, id string) (*ColloquioSlot, error) {
	args := m.Called(id)
	return args.Get(0).(*ColloquioSlot), args.Error(1)
}
func (m *MockRepo) UpdateSlot(ctx context.Context, s *ColloquioSlot) error { return nil }
func (m *MockRepo) CreateBooking(ctx context.Context, b *ColloquioBooking) error {
	b.ID = "booking-id"
	return nil
}
func (m *MockRepo) GetBooking(ctx context.Context, id string) (*ColloquioBooking, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ColloquioBooking), args.Error(1)
}
func (m *MockRepo) GetBookingsByParent(ctx context.Context, pID string) ([]ColloquioBooking, error) {
	return nil, nil
}
func (m *MockRepo) GetBookingsByTeacher(ctx context.Context, tID string) ([]ColloquioBooking, error) {
	return nil, nil
}
func (m *MockRepo) UpdateBooking(ctx context.Context, b *ColloquioBooking) error { return nil }
func (m *MockRepo) CountBookingsForParent(ctx context.Context, pID string, d time.Time, s, e time.Time) (int, error) {
	return 0, nil
}
func (m *MockRepo) GetSettings(ctx context.Context, sID string) (*ColloquioSettings, error) {
	return &ColloquioSettings{BookingBufferHours: 1, BookingWindowDays: 60}, nil
}
func (m *MockRepo) UpdateSettings(ctx context.Context, s *ColloquioSettings) error { return nil }
func (m *MockRepo) GetAnalytics(ctx context.Context, sID string) (*AnalyticsResponse, error) {
	return nil, nil
}
func (m *MockRepo) ResolveParentUserID(ctx context.Context, userID string) (string, error) {
	return userID, nil
}
func (m *MockRepo) IsGuardian(ctx context.Context, parentUserID, studentID string) (bool, error) {
	return true, nil
}

func TestService_CreateSlot(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo, new(MockTeacherRepo), nil, nil, nil)

	req := CreateSlotRequest{
		Dates:      []string{time.Now().AddDate(0, 0, 1).Format("2006-01-02")}, // Future
		StartTime: "10:00", EndTime: "11:00",
		Type: SlotIndividual,
	}

	err := svc.CreateSlots(context.Background(), "teacher1", req)
	assert.NoError(t, err)
}

func TestService_BookSlot(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo, new(MockTeacherRepo), nil, nil, nil)

	slotID := "slot-1"
	slotDate := time.Now().AddDate(0, 0, 2)
	slot := &ColloquioSlot{
		ID: slotID, Date: slotDate,
		StartTime: time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC),
		EndTime:   time.Date(0, 0, 0, 11, 0, 0, 0, time.UTC),
	}

	repo.On("GetSlotByID", slotID).Return(slot, nil)
	repo.On("GetBooking", "booking-id").Return(&ColloquioBooking{ID: "booking-id", ParentName: "Parent", StudentName: "Student"}, nil)

	res, err := svc.BookSlot(context.Background(), "parent1", BookSlotRequest{SlotID: slotID})
	assert.NoError(t, err)
	assert.Equal(t, "booking-id", res.ID)
}

func TestValidator_ValidateSlot(t *testing.T) {
	v := NewValidator()
	s := &ColloquioSlot{
		Date: time.Now().AddDate(0, 0, -1), // Past
	}
	assert.Error(t, v.ValidateSlot(s))
}

func TestService_CancelBooking_Ownership(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo, new(MockTeacherRepo), nil, nil, nil)
	ctx := context.Background()

	bookingID := "booking-1"
	parentID := "parent-1"
	otherParentID := "parent-2"

	booking := &ColloquioBooking{
		ID:       bookingID,
		ParentID: &parentID,
		SlotID:   "slot-1",
		Status:   StatusConfirmed,
	}

	repo.On("GetBooking", bookingID).Return(booking, nil)
	// Mock Try to cancel as other parent
	// The service will check parentID, fail, then try to load slot to check teacherID
	// We need to mock GetSlotByID for that fallback check
	slot := &ColloquioSlot{ID: "slot-1", TeacherID: "teacher-1"}
	repo.On("GetSlotByID", "slot-1").Return(slot, nil)

	err := svc.CancelBooking(ctx, otherParentID, bookingID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")

	// Try as correct parent
	repo.On("GetBooking", bookingID).Return(booking, nil)
	repo.On("UpdateBooking", mock.Anything).Return(nil)

	err = svc.CancelBooking(ctx, parentID, bookingID)
	assert.NoError(t, err)
}

func TestService_BookSlot_Past(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo, new(MockTeacherRepo), nil, nil, nil)
	ctx := context.Background()

	slotID := "slot-past"
	// Slot in the past
	slotDate := time.Now().AddDate(0, 0, -5)
	slot := &ColloquioSlot{
		ID: slotID, Date: slotDate,
		StartTime: slotDate, EndTime: slotDate.Add(time.Hour),
		SchoolID: "school-1",
	}

	repo.On("GetSlotByID", slotID).Return(slot, nil)
	// Return settings needed for validation
	repo.On("GetSettings", "school-1").Return(&ColloquioSettings{
		BookingBufferHours: 1, BookingWindowDays: 30,
	}, nil)

	req := BookSlotRequest{SlotID: slotID, StudentID: nil}
	_, err := svc.BookSlot(ctx, "parent-1", req)
	assert.Error(t, err)
	// The validator should catch "too late" or similar logic depending on implementation
	// ValidateBooking checks now.Add(buffer).After(slotTime)
	// Since slot is in past, now > slotTime, so it should error "too late to book this slot"
	assert.Contains(t, err.Error(), "too late")
}
