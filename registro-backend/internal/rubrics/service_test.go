package rubrics

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRubricsRepo struct {
	mock.Mock
}

func (m *MockRubricsRepo) CreateRubric(ctx context.Context, r *Rubric) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockRubricsRepo) GetRubricByID(ctx context.Context, id string) (*Rubric, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Rubric), args.Error(1)
}

func (m *MockRubricsRepo) ListRubrics(ctx context.Context, schoolID, teacherID string) ([]*Rubric, error) {
	args := m.Called(ctx, schoolID, teacherID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Rubric), args.Error(1)
}

func (m *MockRubricsRepo) UpdateRubric(ctx context.Context, r *Rubric) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockRubricsRepo) DeleteRubric(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRubricsRepo) CreateAssessment(ctx context.Context, a *RubricAssessment) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *MockRubricsRepo) ListAssessmentsByStudent(ctx context.Context, studentID string) ([]*RubricAssessment, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*RubricAssessment), args.Error(1)
}

func (m *MockRubricsRepo) ListAssessmentsByClass(ctx context.Context, classID string) ([]*RubricAssessment, error) {
	args := m.Called(ctx, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*RubricAssessment), args.Error(1)
}

func TestRubricsService_CreateAndAssess(t *testing.T) {
	mockRepo := new(MockRubricsRepo)
	svc := NewService(mockRepo)
	ctx := context.Background()

	t.Run("CreateRubric", func(t *testing.T) {
		mockRepo.On("CreateRubric", ctx, mock.Anything).Return(nil).Once()

		req := CreateRubricRequest{
			SubjectID:   "subj-1",
			Title:       "Rubrica Matematica",
			Description: "Valutazione problemi",
			Criteria: []CreateCriterionInput{
				{
					Name:     "Comprensione",
					MaxScore: 4,
					Levels: []CreateLevelInput{
						{Score: 4, Label: "Avanzato"},
						{Score: 2, Label: "Base"},
					},
				},
			},
		}

		rub, err := svc.CreateRubric(ctx, "school-1", "teacher-1", req)
		assert.NoError(t, err)
		assert.NotNil(t, rub)
		assert.Equal(t, "Rubrica Matematica", rub.Title)
	})

	t.Run("AssessStudent", func(t *testing.T) {
		mockRepo.On("CreateAssessment", ctx, mock.Anything).Return(nil).Once()

		req := CreateAssessmentRequest{
			StudentID: "std-1",
			ClassID:   "class-1",
			Notes:     "Ottimo lavoro",
			Scores: []CriterionScore{
				{CriterionID: "crit-1", LevelID: "lvl-1", Score: 4.0},
				{CriterionID: "crit-2", LevelID: "lvl-2", Score: 3.5},
			},
		}

		ass, err := svc.AssessStudent(ctx, "teacher-1", "rubric-1", req)
		assert.NoError(t, err)
		assert.NotNil(t, ass)
		assert.Equal(t, 7.5, ass.TotalScore)
	})
}
