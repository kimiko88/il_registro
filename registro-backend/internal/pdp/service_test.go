package pdp

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

func (m *MockRepository) Create(ctx context.Context, plan *PdpPlan) (*PdpPlan, error) {
	args := m.Called(ctx, plan)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PdpPlan), args.Error(1)
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*PdpPlan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PdpPlan), args.Error(1)
}

func (m *MockRepository) GetByStudent(ctx context.Context, studentID, year string) ([]*PdpPlan, error) {
	args := m.Called(ctx, studentID, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*PdpPlan), args.Error(1)
}

func (m *MockRepository) GetByClass(ctx context.Context, classID, year string) ([]*PdpPlan, error) {
	args := m.Called(ctx, classID, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*PdpPlan), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, id string, req *UpdatePdpRequest) (*PdpPlan, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PdpPlan), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) SetSharedWithFamily(ctx context.Context, id string, share bool) error {
	args := m.Called(ctx, id, share)
	return args.Error(0)
}

func (m *MockRepository) ApproveByFamily(ctx context.Context, id, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func TestPDPService_CreatePlan(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	t.Run("successful creation by teacher", func(t *testing.T) {
		req := CreatePdpRequest{
			StudentID:    "student-1",
			ClassID:      "class-1",
			AcademicYear: "2025/2026",
			PlanType:     PlanTypePDP,
			Diagnosis:    "DSA Dislessia Evolutiva",
			Content: PdpContent{
				Compensative: []string{"calcolatrice", "tempo_aggiuntivo_30"},
			},
		}

		expectedPlan := &PdpPlan{
			ID:           "pdp-1",
			StudentID:    "student-1",
			ClassID:      "class-1",
			AcademicYear: "2025/2026",
			PlanType:     PlanTypePDP,
			Diagnosis:    "DSA Dislessia Evolutiva",
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*pdp.PdpPlan")).Return(expectedPlan, nil).Once()

		plan, err := svc.CreatePlan(ctx, "teacher-1", "teacher", &req)
		assert.NoError(t, err)
		assert.NotNil(t, plan)
		assert.Equal(t, "student-1", plan.StudentID)
		assert.Equal(t, PlanTypePDP, plan.PlanType)
		assert.Equal(t, "DSA Dislessia Evolutiva", plan.Diagnosis)
		mockRepo.AssertExpectations(t)
	})

	t.Run("forbidden for student role", func(t *testing.T) {
		req := CreatePdpRequest{
			StudentID: "student-1",
			ClassID:   "class-1",
		}
		plan, err := svc.CreatePlan(ctx, "student-1", "student", &req)
		assert.Error(t, err)
		assert.Nil(t, plan)
		assert.Equal(t, ErrUnauthorized, err)
	})
}

func TestPDPService_GetByID_DiagnosisRedaction(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	originalPlan := &PdpPlan{
		ID:               "pdp-100",
		StudentID:        "student-1",
		ClassID:          "class-1",
		SchoolID:         "school-1",
		PlanType:         PlanTypePDP,
		Diagnosis:        "Diagnosi Riservata Clinica",
		SharedWithFamily: true,
		CreatedAt:        time.Now(),
	}

	t.Run("teacher sees full diagnosis", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, "pdp-100").Return(originalPlan, nil).Once()

		plan, err := svc.GetByID(ctx, "teacher", "pdp-100")
		assert.NoError(t, err)
		assert.Equal(t, "Diagnosi Riservata Clinica", plan.Diagnosis)
	})

	t.Run("parent sees shared plan with redacted diagnosis", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, "pdp-100").Return(originalPlan, nil).Once()

		plan, err := svc.GetByID(ctx, "parent", "pdp-100")
		assert.NoError(t, err)
		assert.Equal(t, "", plan.Diagnosis)
	})

	t.Run("parent cannot see unshared plan", func(t *testing.T) {
		unsharedPlan := &PdpPlan{
			ID:               "pdp-200",
			StudentID:        "student-1",
			SharedWithFamily: false,
		}
		mockRepo.On("GetByID", ctx, "pdp-200").Return(unsharedPlan, nil).Once()

		plan, err := svc.GetByID(ctx, "parent", "pdp-200")
		assert.Error(t, err)
		assert.Nil(t, plan)
		assert.Equal(t, ErrNotSharedYet, err)
	})
}

func TestPDPService_ApproveByFamily(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	ctx := context.Background()

	t.Run("parent approves shared plan", func(t *testing.T) {
		plan := &PdpPlan{
			ID:               "pdp-100",
			StudentID:        "student-1",
			SharedWithFamily: true,
		}
		mockRepo.On("GetByID", ctx, "pdp-100").Return(plan, nil).Once()
		mockRepo.On("ApproveByFamily", ctx, "pdp-100", "parent-1").Return(nil).Once()

		err := svc.ApproveByFamily(ctx, "parent", "parent-1", "pdp-100")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("teacher cannot approve plan as family", func(t *testing.T) {
		err := svc.ApproveByFamily(ctx, "teacher", "teacher-1", "pdp-100")
		assert.Error(t, err)
		assert.Equal(t, ErrUnauthorized, err)
	})
}
