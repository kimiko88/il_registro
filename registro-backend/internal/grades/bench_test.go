package grades

import (
	"fmt"
	"testing"
)

func BenchmarkConvertJudgmentToValue(b *testing.B) {
	calc := NewCalculator()
	judgments := []string{
		"ottimo",
		"distinto+",
		"più che buono",
		"quasi sufficiente",
		"insufficiente-",
		"avanzato",
		"base",
	}

	idx := 0
	for b.Loop() {
		_ = calc.ConvertJudgmentToValue(judgments[idx%len(judgments)])
		idx++
	}
}

func BenchmarkCalculateWeightedAverage(b *testing.B) {
	calc := NewCalculator()
	gradesList := make([]Grade, 30)
	for i := 0; i < 30; i++ {
		gradesList[i] = Grade{
			ID:         fmt.Sprintf("g-%d", i),
			GradeValue: 4.0 + float64(i%6)*1.0,
			Weight:     0.5 + float64(i%3)*0.5,
		}
	}

	for b.Loop() {
		_ = calc.CalculateWeightedAverage(gradesList)
	}
}

func BenchmarkCalculateBellCurve(b *testing.B) {
	calc := NewCalculator()
	gradesList := make([]Grade, 50)
	for i := 0; i < 50; i++ {
		gradesList[i] = Grade{
			ID:         fmt.Sprintf("g-%d", i),
			GradeValue: 4.0 + float64(i%7)*0.85,
		}
	}

	for b.Loop() {
		_, _, _, _ = calc.CalculateBellCurve(gradesList)
	}
}

func BenchmarkDetectOutliers(b *testing.B) {
	calc := NewCalculator()
	gradesList := make([]Grade, 50)
	for i := 0; i < 50; i++ {
		gradesList[i] = Grade{
			ID:         fmt.Sprintf("g-%d", i),
			GradeValue: 6.0 + float64(i%4)*0.5,
		}
	}
	// Inject outlier
	gradesList[10].GradeValue = 2.0
	gradesList[40].GradeValue = 10.0

	for b.Loop() {
		_ = calc.DetectOutliers(gradesList)
	}
}
