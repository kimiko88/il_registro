package pdp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePlanRejectsCrossSchoolStudentOrClass(t *testing.T) {
	repo := new(MockRepository)
	service := NewService(repo)
	req := &CreatePdpRequest{StudentID: "student-school-b", ClassID: "class-school-b", AcademicYear: "2026/2027"}

	repo.On("ValidatePlanScope", mock.Anything, "school-a", req.StudentID, req.ClassID, req.CoordinatorID, req.ReferenteID).
		Return(false, nil).Once()

	plan, err := service.CreatePlan(context.Background(), "teacher-a", "teacher", "school-a", req)
	assert.ErrorIs(t, err, ErrUnauthorized)
	assert.Nil(t, plan)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
}

func TestUpdatePlanRejectsCrossSchoolCoordinator(t *testing.T) {
	repo := new(MockRepository)
	service := NewService(repo)
	plan := &PdpPlan{
		ID: "plan-a", StudentID: "student-a", ClassID: "class-a", SchoolID: "school-a",
	}
	foreignCoordinator := "coordinator-school-b"
	req := &UpdatePdpRequest{CoordinatorID: &foreignCoordinator}

	repo.On("GetByID", mock.Anything, "plan-a").Return(plan, nil).Once()
	repo.On("ValidatePlanScope", mock.Anything, "school-a", "student-a", "class-a", &foreignCoordinator, (*string)(nil)).
		Return(false, nil).Once()

	updated, err := service.UpdatePlan(context.Background(), "teacher", "school-a", "plan-a", req)
	assert.ErrorIs(t, err, ErrUnauthorized)
	assert.Nil(t, updated)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
}
