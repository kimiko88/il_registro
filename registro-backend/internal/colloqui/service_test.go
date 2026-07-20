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
func (m *MockRepository) ListUserBookings(ctx context.Context, userID string) ([]*ColloquioBooking, error) {
	args := m.Called(ctx, userID)
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

func TestCreateSlot(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	req := CreateSlotRequest{
		Date:        "2025-11-20",
		StartTime:   "15:00",
		EndTime:     "16:00",
		MaxBookings: 2,
		Type:        TypeIndividual,
		Location:    "Stanza 12 / Meet link",
	}

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
