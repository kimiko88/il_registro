package users

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_ApplyDataRetention(t *testing.T) {
	ctx := context.Background()

	t.Run("Unauthorized role (teacher) returns ErrUnauthorized", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		schoolID := "school-123"
		res, err := service.ApplyDataRetention(ctx, "teacher", &schoolID, 5)
		assert.ErrorIs(t, err, ErrUnauthorized)
		assert.Nil(t, res)
	})

	t.Run("Admin successfully runs retention policy for own school", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		schoolID := "school-456"
		mockRepo.On("ApplyDataRetention", mock.Anything, &schoolID, mock.MatchedBy(func(cutoff time.Time) bool {
			// Expected ~5 years ago
			expectedYear := time.Now().AddDate(-5, 0, 0).Year()
			return cutoff.Year() == expectedYear
		})).Return(12, nil).Once()

		res, err := service.ApplyDataRetention(ctx, "admin", &schoolID, 5)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 12, res.ProcessedCount)
		assert.Equal(t, 5, res.RetentionYears)
		assert.Contains(t, res.Message, "12 student records pseudonymized")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Default retention years is 5 when 0 is provided", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		schoolID := "school-456"
		mockRepo.On("ApplyDataRetention", mock.Anything, &schoolID, mock.MatchedBy(func(cutoff time.Time) bool {
			expectedYear := time.Now().AddDate(-5, 0, 0).Year()
			return cutoff.Year() == expectedYear
		})).Return(3, nil).Once()

		res, err := service.ApplyDataRetention(ctx, "superadmin", &schoolID, 0)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 5, res.RetentionYears)
		assert.Equal(t, 3, res.ProcessedCount)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Superadmin runs retention policy globally with nil schoolID", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		mockRepo.On("ApplyDataRetention", mock.Anything, (*string)(nil), mock.Anything).Return(45, nil).Once()

		res, err := service.ApplyDataRetention(ctx, "superadmin", nil, 7)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 45, res.ProcessedCount)
		assert.Equal(t, 7, res.RetentionYears)
		mockRepo.AssertExpectations(t)
	})
}
