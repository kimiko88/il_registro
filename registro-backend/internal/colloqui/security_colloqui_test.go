package colloqui

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockColloquiRepo struct {
	mock.Mock
	Repository
}

func (m *mockColloquiRepo) GetSlotByID(ctx context.Context, id string) (*ColloquioSlot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ColloquioSlot), args.Error(1)
}

func (m *mockColloquiRepo) GetTeacherProfileID(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *mockColloquiRepo) CancelSlot(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestColloqui_CancelSlot_CrossTenantAdminRejection(t *testing.T) {
	repo := new(mockColloquiRepo)
	svc := NewService(repo)

	slotID := "slot-school-A"
	repo.On("GetSlotByID", mock.Anything, slotID).Return(&ColloquioSlot{
		ID:        slotID,
		SchoolID:  "school-A",
		TeacherID: "teacher-1",
	}, nil)

	// Admin with empty actorSchoolID MUST be rejected
	errEmptySchool := svc.CancelSlot(context.Background(), "admin-1", "admin", "", slotID)
	assert.ErrorIs(t, errEmptySchool, ErrUnauthorized)

	// Admin with mismatched schoolID MUST be rejected
	errDiffSchool := svc.CancelSlot(context.Background(), "admin-1", "admin", "school-B", slotID)
	assert.ErrorIs(t, errDiffSchool, ErrUnauthorized)

	repo.AssertExpectations(t)
}

func TestColloqui_ListSlots_RequiresFilter(t *testing.T) {
	repo := new(mockColloquiRepo)
	svc := NewService(repo)

	_, err := svc.ListSlots(context.Background(), "", "", time.Time{}, time.Time{}, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "obbligatorio")
}
