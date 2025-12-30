package grades

import (
	"io"
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
		// Should handle zero division by returning 0
		assert.Equal(t, 0.0, c.CalculateWeightedAverage(g))
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
// Since we don't have a full mock repo in this specific test file easily without re-implementing it or using the package one
// We can test the logic if it was exposed.
// Since GetStudentGradesWithFilter logic is inside Service method using array manipulation,
// we can skip unit testing the PRIVATE internal loop unless we use an integration test with the mock.
// Assuming "comprehensive_test" style, we can use the MockRepository defined in "service_test.go" (if available in package).
// But `service_test.go` isn't in my viewed files. I'll make a local mock struct here for regression testing the filter logic.

type RegressionMockRepo struct {
	data []Grade
}

// Implement necessary Repository interface methods stub
func (m *RegressionMockRepo) Create(g *Grade) error                  { return nil }
func (m *RegressionMockRepo) Update(g *Grade, h *GradeHistory) error { return nil }
func (m *RegressionMockRepo) Delete(id, teacherID string) error      { return nil }
func (m *RegressionMockRepo) FindByID(id string) (*Grade, error)     { return nil, nil }
func (m *RegressionMockRepo) FindByStudent(studentID string) ([]Grade, error) {
	return m.data, nil // Return all data for filtering check
}

// Stubs for others
func (m *RegressionMockRepo) FindByClass(classID string, semester int) ([]Grade, error) {
	return nil, nil
}
func (m *RegressionMockRepo) FindBySubject(subjectID string, semester int) ([]Grade, error) {
	return nil, nil
}
func (m *RegressionMockRepo) FindByClassAndSubject(classID, subjectID string, semester int) ([]Grade, error) {
	return nil, nil
}
func (m *RegressionMockRepo) FindWithFilter(f GradeFilter) ([]Grade, error)     { return nil, nil }
func (m *RegressionMockRepo) BatchCreate(grades []*Grade) error                 { return nil }
func (m *RegressionMockRepo) GetHistory(gradeID string) ([]GradeHistory, error) { return nil, nil }
func (m *RegressionMockRepo) FindByTeacher(teacherID string) ([]Grade, error)   { return nil, nil }
func (m *RegressionMockRepo) BulkImport(t string, r io.Reader, s int) (*ImportResult, error) {
	return nil, nil
}

func TestService_FilterLogicRegex(t *testing.T) {
	// Setup specific data
	mockData := []Grade{
		{StudentID: "S1", GradeValue: 5, Semester: 1, SubjectID: "MATH", GradeType: GradeTypeNumeric, IsPublished: true},
		{StudentID: "S1", GradeValue: 8, Semester: 2, SubjectID: "HIST", GradeType: GradeTypeNumeric, IsPublished: true},
		{StudentID: "S1", GradeValue: 6, Semester: 1, SubjectID: "MATH", GradeType: GradeTypeJudgment, IsPublished: false},
	}

	repo := &RegressionMockRepo{data: mockData}
	svc := NewService(repo, nil, nil) // userRepo nil, db nil

	t.Run("Filter by Semester", func(t *testing.T) {
		res, err := svc.GetStudentGradesWithFilter("S1", GradeFilter{Semester: 1})
		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, 5.0, res[0].GradeValue)
		assert.Equal(t, 6.0, res[1].GradeValue)
	})

	t.Run("Filter by Subject", func(t *testing.T) {
		res, err := svc.GetStudentGradesWithFilter("S1", GradeFilter{SubjectID: "HIST"})
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, 8.0, res[0].GradeValue)
	})

	t.Run("Filter by Published", func(t *testing.T) {
		pub := true
		res, err := svc.GetStudentGradesWithFilter("S1", GradeFilter{IsPublished: &pub})
		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, 5.0, res[0].GradeValue)
	})
}
