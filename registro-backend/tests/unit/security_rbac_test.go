package unit

import (
	"context"
	"strings"
	"testing"

	"registro-backend/internal/users"
	"registro-backend/tests/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSecurity_RBAC_RolePermissions(t *testing.T) {
	mockRepo := new(testhelpers.MockUsersRepository)
	service := users.NewService(mockRepo)

	req := users.CreateUserRequest{
		Email:     "target@test.com",
		FirstName: "Target",
		LastName:  "User",
		Role:      "student",
		Password:  "Password123!",
	}

	t.Run("student cannot create user", func(t *testing.T) {
		_, err := service.CreateUser(context.Background(), "student", req)
		assert.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "unauthorized")
	})

	t.Run("parent cannot create user", func(t *testing.T) {
		_, err := service.CreateUser(context.Background(), "parent", req)
		assert.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "unauthorized")
	})

	t.Run("teacher cannot create user", func(t *testing.T) {
		_, err := service.CreateUser(context.Background(), "teacher", req)
		assert.Error(t, err)
		assert.Contains(t, strings.ToLower(err.Error()), "unauthorized")
	})
}

func TestSecurity_CSVFormulaInjectionProtection(t *testing.T) {
	mockRepo := new(testhelpers.MockUsersRepository)
	service := users.NewService(mockRepo)

	maliciousUser := &users.User{
		ID:        "u1",
		FirstName: "=CMD|'/C calc'!A1",
		LastName:  "+SUM(1+1)",
		Email:     "@malicious@example.com",
		Role:      "student",
	}

	mockRepo.On("List", mock.Anything, mock.Anything).Return([]users.User{*maliciousUser}, 1, nil)

	data, err := service.ExportUsers(context.Background(), "admin", "", users.UserFilter{}, "csv")
	assert.NoError(t, err)

	output := string(data)
	assert.Contains(t, output, "'=CMD")
	assert.Contains(t, output, "'+SUM")
	assert.Contains(t, output, "'@malicious")
}
