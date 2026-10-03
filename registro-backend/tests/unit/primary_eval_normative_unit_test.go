package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"registro-backend/internal/primaryeval"
)

type mockPrimaryEvalRepo struct {
	mock.Mock
}

func (m *mockPrimaryEvalRepo) CreateObjective(ctx context.Context, obj *primaryeval.LearningObjective) (*primaryeval.LearningObjective, error) {
	args := m.Called(ctx, obj)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*primaryeval.LearningObjective), args.Error(1)
}

func (m *mockPrimaryEvalRepo) ListObjectives(ctx context.Context, schoolID, subjectID string, classID *string, yearGrade int, academicYear string) ([]primaryeval.LearningObjective, error) {
	args := m.Called(ctx, schoolID, subjectID, classID, yearGrade, academicYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]primaryeval.LearningObjective), args.Error(1)
}

func (m *mockPrimaryEvalRepo) GetObjective(ctx context.Context, id string) (*primaryeval.LearningObjective, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*primaryeval.LearningObjective), args.Error(1)
}

func (m *mockPrimaryEvalRepo) DeleteObjective(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockPrimaryEvalRepo) SaveEvaluationsBatch(ctx context.Context, schoolID, teacherID string, req *primaryeval.SaveEvaluationsBatchRequest) error {
	args := m.Called(ctx, schoolID, teacherID, req)
	return args.Error(0)
}

func (m *mockPrimaryEvalRepo) ListEvaluationsByClassAndSubject(ctx context.Context, classID, subjectID string, semester int) ([]primaryeval.PrimaryEvaluation, error) {
	args := m.Called(ctx, classID, subjectID, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]primaryeval.PrimaryEvaluation), args.Error(1)
}

func (m *mockPrimaryEvalRepo) GetStudentEvaluations(ctx context.Context, studentID string, semester int) ([]primaryeval.PrimaryEvaluation, error) {
	args := m.Called(ctx, studentID, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]primaryeval.PrimaryEvaluation), args.Error(1)
}

func TestPrimaryEval_ValidLevels(t *testing.T) {
	assert.True(t, primaryeval.IsValidPrimaryLevel(primaryeval.LevelAvanzato))
	assert.True(t, primaryeval.IsValidPrimaryLevel(primaryeval.LevelIntermedio))
	assert.True(t, primaryeval.IsValidPrimaryLevel(primaryeval.LevelBase))
	assert.True(t, primaryeval.IsValidPrimaryLevel(primaryeval.LevelInViaPrimaAcquisiz))

	// Rejects numeric decimal or invalid grades
	assert.False(t, primaryeval.IsValidPrimaryLevel("8.5"))
	assert.False(t, primaryeval.IsValidPrimaryLevel("ottimo"))
	assert.False(t, primaryeval.IsValidPrimaryLevel("sufficiente"))
}

func TestPrimaryEval_CreateObjectiveValidation(t *testing.T) {
	repo := new(mockPrimaryEvalRepo)
	svc := primaryeval.NewService(repo)
	ctx := context.Background()

	// Missing title
	_, err := svc.CreateObjective(ctx, "school-1", &primaryeval.CreateObjectiveRequest{
		SubjectID: "sub-1",
		YearGrade: 3,
		Title:     "",
	})
	assert.ErrorIs(t, err, primaryeval.ErrMissingTitle)

	// Invalid year grade (e.g. 6 is middle school, not primary)
	_, err = svc.CreateObjective(ctx, "school-1", &primaryeval.CreateObjectiveRequest{
		SubjectID: "sub-1",
		YearGrade: 6,
		Title:     "Obiettivo",
	})
	assert.ErrorIs(t, err, primaryeval.ErrInvalidYearGrade)

	// Valid creation
	repo.On("CreateObjective", ctx, mock.AnythingOfType("*primaryeval.LearningObjective")).
		Return(&primaryeval.LearningObjective{ID: "obj-1", Title: "Geometria"}, nil)

	created, err := svc.CreateObjective(ctx, "school-1", &primaryeval.CreateObjectiveRequest{
		SubjectID:    "sub-1",
		YearGrade:    3,
		Title:        "Geometria",
		AcademicYear: "2025/2026",
	})
	assert.NoError(t, err)
	assert.Equal(t, "obj-1", created.ID)
}

func TestPrimaryEval_SaveEvaluationsBatch_RejectsInvalidLevel(t *testing.T) {
	repo := new(mockPrimaryEvalRepo)
	svc := primaryeval.NewService(repo)
	ctx := context.Background()

	// Rejected because level is "7.5" instead of ministerial level
	err := svc.SaveEvaluationsBatch(ctx, "school-1", "teacher-1", &primaryeval.SaveEvaluationsBatchRequest{
		ClassID:     "class-1",
		SubjectID:   "sub-1",
		ObjectiveID: "obj-1",
		Date:        "2026-03-15",
		Semester:    2,
		Evaluations: []primaryeval.StudentObjectiveLevelItem{
			{StudentID: "stud-1", Level: "7.5"},
		},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "livello non valido")

	// Accepted when using valid level
	repo.On("SaveEvaluationsBatch", ctx, "school-1", "teacher-1", mock.Anything).Return(nil)

	err = svc.SaveEvaluationsBatch(ctx, "school-1", "teacher-1", &primaryeval.SaveEvaluationsBatchRequest{
		ClassID:     "class-1",
		SubjectID:   "sub-1",
		ObjectiveID: "obj-1",
		Date:        "2026-03-15",
		Semester:    2,
		Evaluations: []primaryeval.StudentObjectiveLevelItem{
			{
				StudentID:            "stud-1",
				Level:                primaryeval.LevelAvanzato,
				DimensionAutonomy:    "autonomo",
				DimensionContinuity:  "continuo",
				DimensionFamiliarity: "nota",
				DimensionResources:   "risorse_proprie",
			},
		},
	})
	assert.NoError(t, err)
}

func TestPrimaryEval_GetPrimaryMatrix(t *testing.T) {
	repo := new(mockPrimaryEvalRepo)
	svc := primaryeval.NewService(repo)
	ctx := context.Background()

	classID := "class-1"
	repo.On("ListObjectives", ctx, "school-1", "sub-1", &classID, 0, "").
		Return([]primaryeval.LearningObjective{
			{ID: "obj-1", Title: "Frazioni"},
			{ID: "obj-2", Title: "Problemi"},
		}, nil)

	repo.On("ListEvaluationsByClassAndSubject", ctx, "class-1", "sub-1", 1).
		Return([]primaryeval.PrimaryEvaluation{
			{
				StudentID:   "stud-1",
				StudentName: "Mario Rossi",
				ObjectiveID: "obj-1",
				Level:       primaryeval.LevelAvanzato,
			},
		}, nil)

	matrix, err := svc.GetPrimaryMatrix(ctx, "school-1", "class-1", "sub-1", 1)
	assert.NoError(t, err)
	assert.NotNil(t, matrix)
	assert.Equal(t, 2, len(matrix.Objectives))
	assert.Equal(t, 1, len(matrix.Students))
	assert.Equal(t, "stud-1", matrix.Students[0].StudentID)
	assert.Equal(t, primaryeval.LevelAvanzato, matrix.Students[0].Evaluations["obj-1"].Level)
}
