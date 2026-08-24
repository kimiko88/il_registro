package grades

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrades_CalcSemesterAverages_SentinelValueWhenNoGrades(t *testing.T) {
	grades := []GradeResponse{
		{GradeValue: 7.5, Weight: 1.0, Semester: 1, GradeCategory: string(GradeCategorySummative)},
	}

	avg1, avg2 := calcSemesterAverages(grades)
	assert.Equal(t, 7.5, avg1)
	assert.Equal(t, -1.0, avg2, "Semester 2 with no grades must return sentinel -1.0")
}

func TestGrades_ResolveTeacherProfileID_CachesLookups(t *testing.T) {
	svc := &service{}

	// Manual store in cache
	svc.teacherProfileCache.Store("user-teacher-1", "profile-teacher-1")

	profileID, err := svc.resolveTeacherProfileID(context.Background(), "user-teacher-1")
	assert.NoError(t, err)
	assert.Equal(t, "profile-teacher-1", profileID)
}

func TestGrades_ConvertJudgmentToValue_CompositeJudgments(t *testing.T) {
	calc := NewCalculator()

	// Exact judgments
	assert.Equal(t, 10.0, calc.ConvertJudgmentToValue("Ottimo"))
	assert.Equal(t, 6.0, calc.ConvertJudgmentToValue("Sufficiente"))

	// Composite judgments
	assert.Equal(t, 6.5, calc.ConvertJudgmentToValue("più che sufficiente"))
	assert.Equal(t, 8.5, calc.ConvertJudgmentToValue("buono+"))
	assert.Equal(t, 5.0, calc.ConvertJudgmentToValue("quasi sufficiente"))
	assert.Equal(t, 7.5, calc.ConvertJudgmentToValue("discreto+"))
}

func TestGrades_CalculateWeightedAverage_ZeroWeights(t *testing.T) {
	calc := NewCalculator()

	// No explicit weights -> simple average fallback
	gradesNoWeights := []Grade{
		{GradeValue: 6.0, Weight: 0},
		{GradeValue: 8.0, Weight: 0},
	}
	assert.Equal(t, 7.0, calc.CalculateWeightedAverage(gradesNoWeights))

	// Explicit weighted grades
	gradesWeighted := []Grade{
		{GradeValue: 6.0, Weight: 1.0},
		{GradeValue: 8.0, Weight: 2.0},
	}
	// (6*1 + 8*2) / 3 = 22 / 3 = 7.33
	assert.Equal(t, 7.33, calc.CalculateWeightedAverage(gradesWeighted))
}

func TestGrades_DetectOutliers_IQR(t *testing.T) {
	calc := NewCalculator()

	// Normal class distribution with 1 extreme low outlier
	grades := []Grade{
		{ID: "g1", GradeValue: 7.0},
		{ID: "g2", GradeValue: 7.5},
		{ID: "g3", GradeValue: 8.0},
		{ID: "g4", GradeValue: 7.0},
		{ID: "g5", GradeValue: 8.5},
		{ID: "g6", GradeValue: 7.5},
		{ID: "g7", GradeValue: 8.0},
		{ID: "g8", GradeValue: 2.0}, // clear outlier
	}

	outliers := calc.DetectOutliers(grades)
	assert.Contains(t, outliers, "g8")
	assert.NotContains(t, outliers, "g1")
}
