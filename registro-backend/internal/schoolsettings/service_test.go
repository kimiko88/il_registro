package schoolsettings

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetBySchoolID(ctx context.Context, schoolID string) (*SchoolSettings, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SchoolSettings), args.Error(1)
}

func (m *MockRepository) Upsert(ctx context.Context, settings *SchoolSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

func TestSchoolSettingsService_GetAndUpdate(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	schoolID := "school-123"
	initialSettings := &SchoolSettings{
		SchoolID:                         schoolID,
		RequirePrincipalApprovalForNotes: false,
		AllowParentsViewGrades:           true,
		UpdatedAt:                        time.Now(),
	}

	mockRepo.On("GetBySchoolID", mock.Anything, schoolID).Return(initialSettings, nil)
	mockRepo.On("Upsert", mock.Anything, mock.AnythingOfType("*schoolsettings.SchoolSettings")).Return(nil)

	// Get Settings
	st, err := service.GetSettings(context.Background(), schoolID)
	assert.NoError(t, err)
	assert.False(t, st.RequirePrincipalApprovalForNotes)

	// Update Settings
	trueVal := true
	updated, err := service.UpdateSettings(context.Background(), schoolID, UpdateSchoolSettingsRequest{
		RequirePrincipalApprovalForNotes: &trueVal,
	})
	assert.NoError(t, err)
	assert.True(t, updated.RequirePrincipalApprovalForNotes)
}
