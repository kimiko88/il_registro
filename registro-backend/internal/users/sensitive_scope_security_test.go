package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSensitiveUserOperationsRejectCrossSchoolTargets(t *testing.T) {
	ctx := context.Background()
	schoolA, schoolB := "school-A", "school-B"
	foreignStudent := &User{ID: "student-B", Role: "student", SchoolID: &schoolB, IsActive: true}
	foreignParent := &User{ID: "parent-B", Role: "parent", SchoolID: &schoolB, IsActive: true}

	t.Run("disable MFA", func(t *testing.T) {
		repo := new(MockRepository)
		service := NewService(repo)
		repo.On("GetByID", mock.Anything, "student-B").Return(foreignStudent, nil).Once()

		err := service.DisableMFA(ctx, "admin", "student-B", schoolA)
		assert.ErrorIs(t, err, ErrUnauthorized)
		repo.AssertNotCalled(t, "RevokeAllUserTokens", mock.Anything, mock.Anything)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
	})

	t.Run("GDPR export", func(t *testing.T) {
		repo := new(MockRepository)
		service := NewService(repo)
		repo.On("GetByID", mock.Anything, "student-B").Return(foreignStudent, nil).Once()

		_, err := service.GDPRDataExport(ctx, "admin-A", "admin", "student-B", schoolA)
		assert.ErrorIs(t, err, ErrUnauthorized)
		repo.AssertNotCalled(t, "GetAuditLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
	})

	t.Run("GDPR deletion", func(t *testing.T) {
		repo := new(MockRepository)
		service := NewService(repo)
		repo.On("GetByID", mock.Anything, "student-B").Return(foreignStudent, nil).Once()

		err := service.GDPRDelete(ctx, "admin", "student-B", schoolA)
		assert.ErrorIs(t, err, ErrUnauthorized)
		repo.AssertNotCalled(t, "RevokeAllUserTokens", mock.Anything, mock.Anything)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
	})

	t.Run("audit log", func(t *testing.T) {
		repo := new(MockRepository)
		service := NewService(repo)
		repo.On("GetByID", mock.Anything, "student-B").Return(foreignStudent, nil).Once()

		_, _, err := service.GetAuditLogs(ctx, "admin-A", "admin", "student-B", 50, 0, schoolA)
		assert.ErrorIs(t, err, ErrUnauthorized)
		repo.AssertNotCalled(t, "GetAuditLogs", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
	})

	t.Run("guardian lookup", func(t *testing.T) {
		repo := new(MockRepository)
		service := NewService(repo)
		repo.On("GetByID", mock.Anything, "student-B").Return(foreignStudent, nil).Once()

		_, err := service.GetGuardians(ctx, "admin", "student-B", schoolA)
		assert.ErrorIs(t, err, ErrUnauthorized)
		repo.AssertNotCalled(t, "GetStudentProfile", mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
	})

	t.Run("guardian linking", func(t *testing.T) {
		repo := new(MockRepository)
		service := NewService(repo)
		repo.On("GetByID", mock.Anything, "student-B").Return(foreignStudent, nil).Once()
		repo.On("GetByID", mock.Anything, "parent-B").Return(foreignParent, nil).Once()

		err := service.AddGuardian(ctx, "admin", "student-B", "parent-B", "parent", schoolA)
		assert.ErrorIs(t, err, ErrUnauthorized)
		repo.AssertNotCalled(t, "AddGuardian", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
	})

	t.Run("assignment lookup", func(t *testing.T) {
		repo := new(MockRepository)
		service := NewService(repo)
		repo.On("GetByID", mock.Anything, "student-B").Return(foreignStudent, nil).Once()

		_, err := service.GetUserAssignments(ctx, "admin", "student-B", "admin-A", schoolA)
		assert.ErrorIs(t, err, ErrUnauthorized)
		repo.AssertNotCalled(t, "GetAssignments", mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
	})

	t.Run("student fascicolo", func(t *testing.T) {
		repo := new(MockRepository)
		service := NewService(repo)
		repo.On("GetByID", mock.Anything, "student-B").Return(foreignStudent, nil).Once()

		_, err := service.GetStudentFascicolo(ctx, "admin-A", "admin", "student-B", schoolA)
		assert.ErrorIs(t, err, ErrUnauthorized)
		repo.AssertNotCalled(t, "GetFascicoloSummary", mock.Anything, mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
	})

	t.Run("superadmin cannot create cross-school guardian relation", func(t *testing.T) {
		schoolC := "school-C"
		repo := new(MockRepository)
		service := NewService(repo)
		repo.On("GetByID", mock.Anything, "student-B").Return(foreignStudent, nil).Once()
		repo.On("GetByID", mock.Anything, "parent-C").Return(&User{ID: "parent-C", Role: "parent", SchoolID: &schoolC}, nil).Once()

		err := service.AddGuardian(ctx, "superadmin", "student-B", "parent-C", "parent")
		assert.ErrorIs(t, err, ErrUnauthorized)
		repo.AssertNotCalled(t, "AddGuardian", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		repo.AssertExpectations(t)
	})
}
