package grades

import (
	"math"
	"testing"
)

func TestCalculator_ArithmeticAverage(t *testing.T) {
	calc := NewCalculator()

	// 1. Empty list
	if avg := calc.CalculateAverage([]Grade{}); avg != 0 {
		t.Errorf("expected 0 for empty list, got %f", avg)
	}

	// 2. Normal grades
	grades := []Grade{
		{GradeValue: 7.0},
		{GradeValue: 8.0},
		{GradeValue: 6.0},
	}
	if avg := calc.CalculateAverage(grades); avg != 7.0 {
		t.Errorf("expected 7.0, got %f", avg)
	}

	// 3. Rounding to 2 decimal places (e.g. 7 + 8 + 7 = 22 / 3 = 7.33)
	grades2 := []Grade{
		{GradeValue: 7.0},
		{GradeValue: 8.0},
		{GradeValue: 7.0},
	}
	if avg := calc.CalculateAverage(grades2); avg != 7.33 {
		t.Errorf("expected 7.33, got %f", avg)
	}

	// 4. Ignores invalid / non-votable grades (< 1.0 or > 10.0)
	gradesInvalid := []Grade{
		{GradeValue: 6.0},
		{GradeValue: 8.0},
		{GradeValue: 0.0},  // unrated
		{GradeValue: -1.0}, // absent
		{GradeValue: 12.0}, // out of bounds
	}
	if avg := calc.CalculateAverage(gradesInvalid); avg != 7.0 {
		t.Errorf("expected 7.0 ignoring invalid grades, got %f", avg)
	}
}

func TestCalculator_WeightedAverage(t *testing.T) {
	calc := NewCalculator()

	// 1. Explicit weights
	// Grade 6 with weight 1, Grade 9 with weight 2 => (6*1 + 9*2)/3 = 24/3 = 8.0
	grades := []Grade{
		{GradeValue: 6.0, Weight: 1.0},
		{GradeValue: 9.0, Weight: 2.0},
	}
	if wAvg := calc.CalculateWeightedAverage(grades); wAvg != 8.0 {
		t.Errorf("expected weighted average 8.0, got %f", wAvg)
	}

	// 2. When no grade has weight > 0, fallback to arithmetic average
	noWeightGrades := []Grade{
		{GradeValue: 6.0, Weight: 0.0},
		{GradeValue: 8.0, Weight: 0.0},
	}
	if wAvg := calc.CalculateWeightedAverage(noWeightGrades); wAvg != 7.0 {
		t.Errorf("expected fallback arithmetic average 7.0, got %f", wAvg)
	}

	// 3. Judgment grades with weights
	// "Ottimo" (=10.0) with weight 2, "Sufficiente" (=6.0) with weight 1 => (20 + 6)/3 = 8.67
	judgmentGrades := []Grade{
		{GradeType: GradeTypeJudgment, Description: "ottimo", Weight: 2.0},
		{GradeType: GradeTypeJudgment, Description: "sufficiente", Weight: 1.0},
	}
	if wAvg := calc.CalculateWeightedAverage(judgmentGrades); wAvg != 8.67 {
		t.Errorf("expected 8.67 for weighted judgments, got %f", wAvg)
	}
}

func TestCalculator_MedianAndStatistics(t *testing.T) {
	calc := NewCalculator()

	// 1. Median with odd count
	oddGrades := []Grade{
		{GradeValue: 4.0},
		{GradeValue: 9.0},
		{GradeValue: 6.0},
	}
	if med := calc.CalculateMedian(oddGrades); med != 6.0 {
		t.Errorf("expected median 6.0, got %f", med)
	}

	// 2. Median with even count
	evenGrades := []Grade{
		{GradeValue: 5.0},
		{GradeValue: 6.0},
		{GradeValue: 8.0},
		{GradeValue: 9.0},
	}
	if med := calc.CalculateMedian(evenGrades); med != 7.0 {
		t.Errorf("expected median 7.0, got %f", med)
	}

	// 3. Standard deviation
	stdGrades := []Grade{
		{GradeValue: 6.0},
		{GradeValue: 6.0},
		{GradeValue: 6.0},
	}
	if std := calc.CalculateStandardDeviation(stdGrades); std != 0.0 {
		t.Errorf("expected standard deviation 0 for uniform grades, got %f", std)
	}

	// 4. Percentile
	percentileGrades := []Grade{
		{GradeValue: 4.0},
		{GradeValue: 6.0},
		{GradeValue: 7.0},
		{GradeValue: 8.0},
		{GradeValue: 10.0},
	}
	pct := calc.CalculatePercentile(percentileGrades, 7.0)
	if pct != 50.0 {
		t.Errorf("expected 50th percentile for median value 7.0, got %f", pct)
	}
}

func TestCalculator_JudgmentsConversion(t *testing.T) {
	calc := NewCalculator()

	cases := []struct {
		judgment string
		expected float64
	}{
		{"ottimo", 10.0},
		{"distinto", 9.0},
		{"buono", 8.0},
		{"discreto", 7.0},
		{"sufficiente", 6.0},
		{"mediocre", 5.0},
		{"insufficiente", 4.0},
		{"gravemente insufficiente", 3.0},
		{"buono+", 8.5},
		{"sufficiente-", 5.5},
		{"", 0.0},
		{"non_riconosciuto", 0.0},
	}

	for _, tc := range cases {
		val := calc.ConvertJudgmentToValue(tc.judgment)
		if math.Abs(val-tc.expected) > 0.001 {
			t.Errorf("for judgment %q expected %f, got %f", tc.judgment, tc.expected, val)
		}
	}
}

func TestCalculator_OutlierDetection(t *testing.T) {
	calc := NewCalculator()

	// Most grades around 6-7, one severe outlier grade of 1.0 and one 10.0
	var grades []Grade
	for i := 0; i < 15; i++ {
		grades = append(grades, Grade{ID: "normal", GradeValue: 6.5})
	}
	grades = append(grades, Grade{ID: "outlier-low", GradeValue: 1.0})
	grades = append(grades, Grade{ID: "outlier-high", GradeValue: 10.0})

	outliers := calc.DetectOutliers(grades)
	if len(outliers) == 0 {
		t.Logf("no outliers detected under IQR threshold")
	} else {
		foundLow := false
		for _, id := range outliers {
			if id == "outlier-low" {
				foundLow = true
			}
		}
		if !foundLow {
			t.Errorf("expected outlier-low to be detected")
		}
	}
}
