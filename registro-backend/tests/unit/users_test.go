package unit

import (
	"context"
	"errors"
	"testing"

	"registro-backend/internal/users"
	"registro-backend/tests/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUsersService_CreateUser(t *testing.T) {
	mockRepo := new(testhelpers.MockUsersRepository)
	service := users.NewService(mockRepo)

	t.Run("success", func(t *testing.T) {
		req := users.CreateUserRequest{
			Email:      "newuser@test.com",
			FirstName:  "Test",
			LastName:   "User",
			Role:       "student",
			Password:   "StrongPass1!",
			FiscalCode: "RSSMRA80A01H501U", // Valid-looking fiscal code
		}

		mockRepo.On("GetByEmail", mock.Anything, "newuser@test.com").Return(nil, errors.New("not found"))
		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *users.User) bool {
			return u.Email == "newuser@test.com" && u.FirstName == "Test"
		})).Return(nil)

		mockRepo.On("AddPasswordHistory", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		// Audit log mock - Must match signature in repository.go
		// LogAudit(ctx context.Context, log *AuditLog) error
		mockRepo.On("LogAudit", mock.Anything, mock.Anything).Return(nil)

		// CreateUser checks permManager.HasPermission(actorRole, permissions.UserCreate).
		// "admin" role usually has this.
		// If permManager is internal and checks hardcoded map, we assume "admin" works.
		user, err := service.CreateUser(context.Background(), "admin", req)

		if err != nil {
			t.Logf("Error: %v", err)
		}
		assert.NoError(t, err)
		assert.NotNil(t, user)
		if user != nil {
			assert.Equal(t, "newuser@test.com", user.Email)
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		req := users.CreateUserRequest{}
		_, err := service.CreateUser(context.Background(), "student", req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unauthorized")
	})
}
