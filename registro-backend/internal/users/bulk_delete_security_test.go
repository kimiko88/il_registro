package users

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBulkDeleteUsers_EnforcesTenantScopeForEveryNonSuperadmin(t *testing.T) {
	ctx := context.Background()
	schoolA, schoolB := "school-A", "school-B"
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.On("ListByIDs", mock.Anything, []string{"foreign-user"}).
		Return([]User{{ID: "foreign-user", Role: "student", SchoolID: &schoolB}}, nil).Once()

	count, err := service.BulkDeleteUsers(ctx, "principal", schoolA, []string{"foreign-user"})
	assert.NoError(t, err)
	assert.Zero(t, count)
	mockRepo.AssertNotCalled(t, "RevokeAllUserTokens", mock.Anything, mock.Anything)
	mockRepo.AssertNotCalled(t, "BulkDelete", mock.Anything, mock.Anything)
	mockRepo.AssertExpectations(t)
}

func TestBulkDeleteUsers_RevokesSessionsBeforeDeleting(t *testing.T) {
	ctx := context.Background()
	schoolA := "school-A"
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.On("ListByIDs", mock.Anything, []string{"local-user"}).
		Return([]User{{ID: "local-user", Role: "student", SchoolID: &schoolA}}, nil).Once()
	mockRepo.On("RevokeAllUserTokens", mock.Anything, "local-user").Return(nil).Once()
	mockRepo.On("BulkDelete", mock.Anything, []string{"local-user"}).Return(1, nil).Once()

	count, err := service.BulkDeleteUsers(ctx, "principal", schoolA, []string{"local-user"})
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
	mockRepo.AssertExpectations(t)
}

func TestBulkDeleteUsers_AbortsWhenSessionRevocationFails(t *testing.T) {
	ctx := context.Background()
	schoolA := "school-A"
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	revokeErr := errors.New("redis/database unavailable")

	mockRepo.On("ListByIDs", mock.Anything, []string{"local-user"}).
		Return([]User{{ID: "local-user", Role: "student", SchoolID: &schoolA}}, nil).Once()
	mockRepo.On("RevokeAllUserTokens", mock.Anything, "local-user").Return(revokeErr).Once()

	count, err := service.BulkDeleteUsers(ctx, "principal", schoolA, []string{"local-user"})
	assert.ErrorIs(t, err, revokeErr)
	assert.Zero(t, count)
	mockRepo.AssertNotCalled(t, "BulkDelete", mock.Anything, mock.Anything)
	mockRepo.AssertExpectations(t)
}

func TestBulkDeleteUsers_NonSuperadminWithoutSchoolScopeDeletesNothing(t *testing.T) {
	ctx := context.Background()
	schoolA := "school-A"
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.On("ListByIDs", mock.Anything, []string{"local-user"}).
		Return([]User{{ID: "local-user", Role: "student", SchoolID: &schoolA}}, nil).Once()

	count, err := service.BulkDeleteUsers(ctx, "vice_principal", "", []string{"local-user"})
	assert.NoError(t, err)
	assert.Zero(t, count)
	mockRepo.AssertNotCalled(t, "RevokeAllUserTokens", mock.Anything, mock.Anything)
	mockRepo.AssertNotCalled(t, "BulkDelete", mock.Anything, mock.Anything)
	mockRepo.AssertExpectations(t)
}
