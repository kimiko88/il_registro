package colloqui

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateSlot(ctx context.Context, slot *ColloquioSlot) error {
	args := m.Called(ctx, slot)
	return args.Error(0)
}
func (m *MockRepository) GetSlotByID(ctx context.Context, id string) (*ColloquioSlot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ColloquioSlot), args.Error(1)
}
func (m *MockRepository) ListSlots(ctx context.Context, filter SlotFilter) ([]*ColloquioSlot, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ColloquioSlot), args.Error(1)
}
func (m *MockRepository) PatchSlot(ctx context.Context, slotID string, startTime, endTime string) error {
	args := m.Called(ctx, slotID, startTime, endTime)
	return args.Error(0)
}
func (m *MockRepository) CancelSlot(ctx context.Context, slotID string) error {
	args := m.Called(ctx, slotID)
	return args.Error(0)
}
func (m *MockRepository) CreateBooking(ctx context.Context, booking *ColloquioBooking) error {
	args := m.Called(ctx, booking)
	return args.Error(0)
}
func (m *MockRepository) GetBookingByID(ctx context.Context, id string) (*ColloquioBooking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ColloquioBooking), args.Error(1)
}
func (m *MockRepository) ListUserBookings(ctx context.Context, userID, schoolID string) ([]*ColloquioBooking, error) {
	args := m.Called(ctx, userID, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ColloquioBooking), args.Error(1)
}
func (m *MockRepository) ListSlotBookings(ctx context.Context, slotID string) ([]*ColloquioBooking, error) {
	args := m.Called(ctx, slotID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ColloquioBooking), args.Error(1)
}
func (m *MockRepository) UpdateBookingStatus(ctx context.Context, bookingID string, status BookingStatus, changedBy string, reason string) error {
	args := m.Called(ctx, bookingID, status, changedBy, reason)
	return args.Error(0)
}
func (m *MockRepository) GetTeacherProfileID(ctx context.Context, userID string) (string, error) {
	return userID, nil
}
func (m *MockRepository) GetParentProfileID(ctx context.Context, userID string) (string, error) {
	return userID, nil
}
func (m *MockRepository) GetStudentProfileID(ctx context.Context, userID string) (string, error) {
	return userID, nil
}
func (m *MockRepository) IsGuardian(ctx context.Context, parentUserID, studentUserID string) (bool, error) {
	args := m.Called(ctx, parentUserID, studentUserID)
	return args.Bool(0), args.Error(1)
}
func (m *MockRepository) ExistsOverlappingSlot(ctx context.Context, teacherID, date, startTime, endTime string) (bool, error) {
	args := m.Called(ctx, teacherID, date, startTime, endTime)
	return args.Bool(0), args.Error(1)
}

func TestCreateSlot(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	req := CreateSlotRequest{
		Date:        "2027-11-20",
		StartTime:   "15:00",
		EndTime:     "16:00",
		MaxBookings: 2,
		Type:        TypeIndividual,
		Location:    "Stanza 12 / Meet link",
	}

	mockRepo.On("ExistsOverlappingSlot", mock.Anything, "t-1", "2027-11-20", "15:00", "16:00").Return(false, nil).Once()
	mockRepo.On("CreateSlot", mock.Anything, mock.MatchedBy(func(s *ColloquioSlot) bool {
		return s.TeacherID == "t-1" && s.StartTime == "15:00"
	})).Return(nil).Once()

	slot, err := svc.CreateSlot(context.Background(), "t-1", "school-1", req)
	assert.NoError(t, err)
	assert.NotNil(t, slot)
	assert.Equal(t, "15:00", slot.StartTime)
	mockRepo.AssertExpectations(t)
}

func TestBookSlot(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	req := CreateBookingRequest{
		SlotID: "slot-1",
		Notes:  "Richiesta chiarimento voti",
	}

	parentID := "p-1"
	expectedBooking := &ColloquioBooking{
		ID:       "b-1",
		SlotID:   "slot-1",
		ParentID: &parentID,
		Status:   StatusConfirmed,
	}

	mockRepo.On("GetSlotByID", mock.Anything, "slot-1").Return(&ColloquioSlot{ID: "slot-1", MaxBookings: 10, BookingCount: 0}, nil).Once()
	mockRepo.On("CreateBooking", mock.Anything, mock.MatchedBy(func(b *ColloquioBooking) bool {
		return b.SlotID == "slot-1"
	})).Return(nil).Once()
	mockRepo.On("GetBookingByID", mock.Anything, mock.Anything).Return(expectedBooking, nil).Once()

	booking, err := svc.BookSlot(context.Background(), "p-1", req)
	assert.NoError(t, err)
	assert.NotNil(t, booking)
	assert.Equal(t, StatusConfirmed, booking.Status)
	mockRepo.AssertExpectations(t)
}

func TestCancelSlot(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	slot := &ColloquioSlot{ID: "slot-1", TeacherID: "t-1"}
	mockRepo.On("GetSlotByID", mock.Anything, "slot-1").Return(slot, nil).Once()
	mockRepo.On("CancelSlot", mock.Anything, "slot-1").Return(nil).Once()

	err := svc.CancelSlot(context.Background(), "t-1", "teacher", "slot-1")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetBookingByIDAuthorization(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	parentID := "parent-1"
	booking := &ColloquioBooking{
		ID:       "booking-1",
		SlotID:   "slot-1",
		ParentID: &parentID,
	}
	slot := &ColloquioSlot{
		ID:        "slot-1",
		TeacherID: "teacher-1",
	}

	mockRepo.On("GetBookingByID", mock.Anything, "booking-1").Return(booking, nil)
	mockRepo.On("GetSlotByID", mock.Anything, "slot-1").Return(slot, nil)

	// Admin access succeeds
	b, err := svc.GetBookingByID(context.Background(), "admin-1", "admin", "booking-1")
	assert.NoError(t, err)
	assert.NotNil(t, b)

	// Parent owner access succeeds
	b, err = svc.GetBookingByID(context.Background(), "parent-1", "parent", "booking-1")
	assert.NoError(t, err)
	assert.NotNil(t, b)

	// Other parent access fails
	_, err = svc.GetBookingByID(context.Background(), "parent-2", "parent", "booking-1")
	assert.Error(t, err)
	assert.Equal(t, ErrUnauthorized, err)
}
