package grades

import (
	"math"
	"testing"
)

func FuzzConvertJudgmentToValue(f *testing.F) {
	seeds := []string{
		"ottimo",
		"distinto",
		"buono",
		"discreto",
		"sufficiente",
		"mediocre",
		"insufficiente",
		"gravemente insufficiente",
		"ottimo+",
		"buono-",
		"più che buono",
		"quasi distinto",
		"avanzato",
		"intermedio",
		"base",
		"",
		"sconosciuto",
		"\x00\xff",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	calc := NewCalculator()

	f.Fuzz(func(t *testing.T, judgment string) {
		val := calc.ConvertJudgmentToValue(judgment)
		if math.IsNaN(val) || math.IsInf(val, 0) {
			t.Fatalf("returned NaN or Inf for judgment %q", judgment)
		}
		if val < 0.0 || val > 10.0 {
			t.Fatalf("converted judgment out of bounds [0, 10]: %f for %q", val, judgment)
		}
	})
}

func FuzzCalculateBellCurve(f *testing.F) {
	f.Add(6.0, 7.0, 8.0, 9.0)
	f.Add(1.0, 1.0, 1.0, 1.0)
	f.Add(10.0, 10.0, 10.0, 10.0)
	f.Add(4.0, 5.5, 6.0, 8.5)

	calc := NewCalculator()

	f.Fuzz(func(t *testing.T, g1, g2, g3, g4 float64) {
		// Only test realistic grades in [1, 10]
		if g1 < 1.0 || g1 > 10.0 || g2 < 1.0 || g2 > 10.0 || g3 < 1.0 || g3 > 10.0 || g4 < 1.0 || g4 > 10.0 {
			return
		}

		gradesList := []Grade{
			{ID: "g1", GradeValue: g1},
			{ID: "g2", GradeValue: g2},
			{ID: "g3", GradeValue: g3},
			{ID: "g4", GradeValue: g4},
		}

		mean, stdDev, skewness, kurtosis := calc.CalculateBellCurve(gradesList)

		if math.IsNaN(mean) || math.IsInf(mean, 0) {
			t.Fatalf("mean is NaN or Inf: %f", mean)
		}
		if math.IsNaN(stdDev) || math.IsInf(stdDev, 0) {
			t.Fatalf("stdDev is NaN or Inf: %f", stdDev)
		}
		if math.IsNaN(skewness) || math.IsInf(skewness, 0) {
			t.Fatalf("skewness is NaN or Inf: %f", skewness)
		}
		if math.IsNaN(kurtosis) || math.IsInf(kurtosis, 0) {
			t.Fatalf("kurtosis is NaN or Inf: %f", kurtosis)
		}
	})
}

func FuzzCalculateWeightedAverage(f *testing.F) {
	f.Add(7.0, 1.0, 8.0, 2.0)
	f.Add(6.0, 0.0, 10.0, 1.0)
	f.Add(5.0, 0.5, 9.0, 1.5)

	calc := NewCalculator()

	f.Fuzz(func(t *testing.T, val1, w1, val2, w2 float64) {
		if val1 < 1.0 || val1 > 10.0 || val2 < 1.0 || val2 > 10.0 {
			return
		}
		if w1 < 0.0 || w1 > 10.0 || w2 < 0.0 || w2 > 10.0 {
			return
		}

		gradesList := []Grade{
			{ID: "g1", GradeValue: val1, Weight: w1},
			{ID: "g2", GradeValue: val2, Weight: w2},
		}

		avg := calc.CalculateWeightedAverage(gradesList)
		if math.IsNaN(avg) || math.IsInf(avg, 0) {
			t.Fatalf("weighted average is NaN or Inf: %f", avg)
		}
		if avg < 0.0 || avg > 10.0 {
			t.Fatalf("weighted average out of bounds: %f", avg)
		}
	})
}
