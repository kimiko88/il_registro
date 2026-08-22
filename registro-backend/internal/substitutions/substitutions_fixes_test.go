package substitutions

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct{ mock.Mock }

func (m *MockRepo) Create(ctx context.Context, s *Substitution) error { return m.Called(s).Error(0) }
func (m *MockRepo) ListBySchool(ctx context.Context, schoolID, date string) ([]*Substitution, error) {
	return nil, nil
}
func (m *MockRepo) ListByTeacher(ctx context.Context, teacherID, date string) ([]*Substitution, error) {
	return nil, nil
}
func (m *MockRepo) AssignSubstitute(ctx context.Context, id, substituteID, notes string) error {
	return nil
}
func (m *MockRepo) ConfirmSubstitution(ctx context.Context, id, teacherID string) error { return nil }
func (m *MockRepo) SignRegister(ctx context.Context, id, sigHash, notes string) error   { return nil }
func (m *MockRepo) GetByID(ctx context.Context, id string) (*Substitution, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Substitution), args.Error(1)
}
func (m *MockRepo) GetTeacherProfileID(ctx context.Context, userID string) (string, error) {
	return userID, nil
}
func (m *MockRepo) GetAvailableTeachers(ctx context.Context, schoolID string) ([]TeacherCandidate, error) {
	return m.Called(ctx, schoolID).Get(0).([]TeacherCandidate), m.Called(ctx, schoolID).Error(1)
}
func (m *MockRepo) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	args := m.Called(ctx, teacherID, classID)
	return args.Bool(0), args.Error(1)
}
func (m *MockRepo) IsTeacherAssignedToSubject(ctx context.Context, teacherID, subjectID string) (bool, error) {
	args := m.Called(ctx, teacherID, subjectID)
	return args.Bool(0), args.Error(1)
}
func (m *MockRepo) GetWeeklySubstitutionCount(ctx context.Context, teacherID string) (int, error) {
	args := m.Called(ctx, teacherID)
	return args.Int(0), args.Error(1)
}

func TestRecommendSubstitutes_SortedByScore(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)

	candidates := []TeacherCandidate{
		{TeacherID: "t1", TeacherName: "Rossi"},
		{TeacherID: "t2", TeacherName: "Bianchi"},
		{TeacherID: "t3", TeacherName: "Verdi"},
	}
	repo.On("GetAvailableTeachers", mock.Anything, "school-1").Return(candidates, nil)

	// t1: insegna la classe (+25) e la materia (+15) -> score 90
	repo.On("IsTeacherAssignedToClass", mock.Anything, "t1", "class-1").Return(true, nil)
	repo.On("IsTeacherAssignedToSubject", mock.Anything, "t1", "subj-1").Return(true, nil)
	repo.On("GetWeeklySubstitutionCount", mock.Anything, "t1").Return(0, nil) // +10

	// t2: nessun bonus -> score 50
	repo.On("IsTeacherAssignedToClass", mock.Anything, "t2", "class-1").Return(false, nil)
	repo.On("IsTeacherAssignedToSubject", mock.Anything, "t2", "subj-1").Return(false, nil)
	repo.On("GetWeeklySubstitutionCount", mock.Anything, "t2").Return(5, nil)

	// t3: insegna la materia (+15) -> score 65
	repo.On("IsTeacherAssignedToClass", mock.Anything, "t3", "class-1").Return(false, nil)
	repo.On("IsTeacherAssignedToSubject", mock.Anything, "t3", "subj-1").Return(true, nil)
	repo.On("GetWeeklySubstitutionCount", mock.Anything, "t3").Return(1, nil) // +5 (count <= 2)

	recs, err := svc.RecommendSubstitutes(context.Background(), "school-1", "class-1", "subj-1", "2026-09-01", 1)
	assert.NoError(t, err)
	assert.Len(t, recs, 3)

	// Verifica ordine decrescente per score
	for i := 1; i < len(recs); i++ {
		assert.GreaterOrEqual(t, recs[i-1].Score, recs[i].Score,
			"recs[%d].Score=%d >= recs[%d].Score=%d", i-1, recs[i-1].Score, i, recs[i].Score)
	}
	// Il primo deve essere t1 (score 90)
	assert.Equal(t, "t1", recs[0].TeacherID)
}
