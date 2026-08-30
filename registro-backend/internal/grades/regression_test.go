package grades

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- Calculator Regression Tests ---

func TestCalculator_EdgeCases(t *testing.T) {
	c := NewCalculator()

	t.Run("Empty Slice", func(t *testing.T) {
		assert.Equal(t, 0.0, c.CalculateAverage(nil))
		assert.Equal(t, 0.0, c.CalculateWeightedAverage(nil))
		assert.Equal(t, 0.0, c.CalculateMedian(nil))
		assert.Equal(t, 0.0, c.CalculateStandardDeviation(nil))
	})

	t.Run("Single Element StdDev", func(t *testing.T) {
		g := []Grade{{GradeValue: 8.0}}
		// StdDev requires >= 2 elements usually, code return 0
		assert.Equal(t, 0.0, c.CalculateStandardDeviation(g))
	})

	t.Run("Weighted Average with Zero Weights", func(t *testing.T) {
		g := []Grade{
			{GradeValue: 10, Weight: 0},
			{GradeValue: 10, Weight: 0},
		}
		// Falls back to arithmetic average (10.0) when total weights is zero
		assert.Equal(t, 10.0, c.CalculateWeightedAverage(g))
	})

	t.Run("Weighted Average Mixed", func(t *testing.T) {
		g := []Grade{
			{GradeValue: 10, Weight: 1}, // 10
			{GradeValue: 6, Weight: 1},  // 6 -> 16 / 2 = 8
		}
		assert.Equal(t, 8.0, c.CalculateWeightedAverage(g))

		g2 := []Grade{
			{GradeValue: 10, Weight: 0.5}, // 5
			{GradeValue: 6, Weight: 0.5},  // 3 -> 8 / 1 = 8
		}
		assert.Equal(t, 8.0, c.CalculateWeightedAverage(g2))

		g3 := []Grade{
			{GradeValue: 10, Weight: 2}, // 20
			{GradeValue: 5, Weight: 1},  // 5 -> 25 / 3 = 8.33
		}
		assert.InDelta(t, 8.33, c.CalculateWeightedAverage(g3), 0.01)
	})
}

// --- Service Filtering Regression Tests ---

type RegressionMockRepo struct {
	data []Grade
}

// Implement necessary Repository interface methods stub
func (m *RegressionMockRepo) Create(ctx context.Context, g *Grade) error                  { return nil }
func (m *RegressionMockRepo) Update(ctx context.Context, g *Grade, h *GradeHistory) error { return nil }
func (m *RegressionMockRepo) Delete(ctx context.Context, id, teacherID string) error      { return nil }
func (m *RegressionMockRepo) SoftDelete(ctx context.Context, id, modifiedBy string) error { return nil }
func (m *RegressionMockRepo) FindByID(ctx context.Context, id string) (*Grade, error)     { return nil, nil }
func (m *RegressionMockRepo) FindByStudent(ctx context.Context, studentID string) ([]Grade, error) {
	return m.data, nil // Return all data for filtering check
}

// Stubs for others
func (m *RegressionMockRepo) FindByClass(ctx context.Context, classID string, semester int) ([]Grade, error) {
	return nil, nil
}
func (m *RegressionMockRepo) FindBySubject(ctx context.Context, subjectID string, semester int) ([]Grade, error) {
	return nil, nil
}
func (m *RegressionMockRepo) FindByClassAndSubject(ctx context.Context, classID, subjectID string, semester int) ([]Grade, error) {
	return nil, nil
}
func (m *RegressionMockRepo) FindWithFilter(ctx context.Context, f GradeFilter) ([]Grade, error) {
	var filtered []Grade
	for _, g := range m.data {
		if g.DeletedAt != nil {
			continue
		}
		if f.StudentID != "" && g.StudentID != f.StudentID {
			continue
		}
		if f.Semester > 0 && int(g.Semester) != f.Semester {
			continue
		}
		if f.SubjectID != "" && g.SubjectID != f.SubjectID {
			continue
		}
		if f.GradeType != "" && string(g.GradeType) != f.GradeType {
			continue
		}
		if f.IsPublished != nil && g.IsPublished != *f.IsPublished {
			continue
		}
		filtered = append(filtered, g)
	}
	return filtered, nil
}
func (m *RegressionMockRepo) FindWithFilterPaginated(ctx context.Context, f GradeFilter) ([]Grade, int, error) {
	res, err := m.FindWithFilter(ctx, f)
	return res, len(res), err
}
func (m *RegressionMockRepo) BatchCreate(ctx context.Context, grades []*Grade) error                 { return nil }
func (m *RegressionMockRepo) GetHistory(ctx context.Context, gradeID string) ([]GradeHistory, error) { return nil, nil }
func (m *RegressionMockRepo) FindByTeacher(ctx context.Context, teacherID string) ([]Grade, error)   { return nil, nil }
func (m *RegressionMockRepo) CreateTest(ctx context.Context, test *ClassTest) error                  { return nil }
func (m *RegressionMockRepo) FindTestsByClassAndSubject(ctx context.Context, classID string, subjectID string) ([]ClassTest, error) {
	return nil, nil
}
func (m *RegressionMockRepo) DeleteTest(ctx context.Context, id string) error       { return nil }
func (m *RegressionMockRepo) UpdateTest(ctx context.Context, test *ClassTest) error { return nil }
func (m *RegressionMockRepo) FindUpcomingTestsByClass(ctx context.Context, classID string) ([]ClassTest, error) {
	return nil, nil
}
func (m *RegressionMockRepo) FindGradesByTestID(ctx context.Context, testID string) ([]Grade, error) { return nil, nil }
func (m *RegressionMockRepo) FindTestByID(ctx context.Context, id string) (*ClassTest, error)        { return nil, nil }

// fix: FindEnrolledSubjects era mancante nel mock causando errore di compilazione
func (m *RegressionMockRepo) FindEnrolledSubjects(ctx context.Context, studentID string, semester int) ([]string, error) {
	return nil, nil
}

func (m *RegressionMockRepo) GetWeightConfigs(ctx context.Context, schoolID, subjectID, classID string) ([]GradeWeightConfig, error) {
	return nil, nil
}
func (m *RegressionMockRepo) UpsertWeightConfig(ctx context.Context, cfg *GradeWeightConfig) (*GradeWeightConfig, error) {
	return cfg, nil
}
func (m *RegressionMockRepo) DeleteWeightConfig(ctx context.Context, id string) error { return nil }

func (m *RegressionMockRepo) GetStudentClassAndSchoolInfo(ctx context.Context, studentID string) (string, string, string, string, error) {
	return "", "", "", "", nil
}
func (m *RegressionMockRepo) GetTeacherNamesByClass(ctx context.Context, classID string) (map[string]string, error) {
	return nil, nil
}
func (m *RegressionMockRepo) GetSubjectNamesMap(ctx context.Context, schoolID string) (map[string]string, error) {
	return nil, nil
}
func (m *RegressionMockRepo) GetScrutinyRecordSummary(ctx context.Context, studentID string, semester int) (float64, float64, bool, error) {
	return 0, 0, false, nil
}
func (m *RegressionMockRepo) GetStudentAbsenceCountForPeriod(ctx context.Context, studentID, startD, endD string) (int, error) {
	return 0, nil
}
func (m *RegressionMockRepo) GetClassSubjectAverage(ctx context.Context, classID, subjectID string, semester int, studentID string) (float64, error) {
	return 7.5, nil
}
func (m *RegressionMockRepo) CheckClassAccessPermission(ctx context.Context, actorID, actorRole, classID string) (bool, error) {
	return true, nil
}

func TestService_FilterLogicRegex(t *testing.T) {
	// Setup specific data
	mockData := []Grade{
		{StudentID: "S1", GradeValue: 5, Semester: 1, SubjectID: "MATH", GradeType: GradeTypeNumeric, IsPublished: true},
		{StudentID: "S1", GradeValue: 8, Semester: 2, SubjectID: "HIST", GradeType: GradeTypeNumeric, IsPublished: true},
		{StudentID: "S1", GradeValue: 6, Semester: 1, SubjectID: "MATH", GradeType: GradeTypeJudgment, IsPublished: false},
	}

	repo := &RegressionMockRepo{data: mockData}
	svc := NewService(repo, nil, nil, nil) // userRepo nil, db nil, broadcaster nil

	t.Run("Filter by Semester", func(t *testing.T) {
		res, err := svc.GetStudentGradesWithFilter(context.Background(), "admin-id", "admin", "S1", GradeFilter{Semester: 1})
		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, 5.0, res[0].GradeValue)
		assert.Equal(t, 6.0, res[1].GradeValue)
	})

	t.Run("Filter by Subject", func(t *testing.T) {
		res, err := svc.GetStudentGradesWithFilter(context.Background(), "admin-id", "admin", "S1", GradeFilter{SubjectID: "HIST"})
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, 8.0, res[0].GradeValue)
	})

	t.Run("Filter by Published", func(t *testing.T) {
		pub := true
		res, err := svc.GetStudentGradesWithFilter(context.Background(), "admin-id", "admin", "S1", GradeFilter{IsPublished: &pub})
		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, 5.0, res[0].GradeValue)
	})
}

func TestCalculator_AverageWithAbsence(t *testing.T) {
	c := NewCalculator()

	t.Run("Average excludes -1", func(t *testing.T) {
		grades := []Grade{
			{GradeValue: 8.0, Weight: 1.0},
			{GradeValue: -1.0, Weight: 1.0}, // Absence
			{GradeValue: 6.0, Weight: 1.0},
		}
		assert.Equal(t, 7.0, c.CalculateAverage(grades))
		assert.Equal(t, 7.0, c.CalculateWeightedAverage(grades))
	})

	t.Run("Zero grade is unrated and excluded from average", func(t *testing.T) {
		grades := []Grade{
			{GradeValue: 8.0, Weight: 1.0},
			{GradeValue: 0.0, Weight: 1.0}, // Unrated/unset grade (0)
			{GradeValue: 6.0, Weight: 1.0},
		}
		assert.Equal(t, 7.0, c.CalculateAverage(grades))
		assert.Equal(t, 7.0, c.CalculateWeightedAverage(grades))
	})
}
