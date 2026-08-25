package grades

import (
	"math"
	"sort"
	"strings"
)

type Calculator struct{}

func NewCalculator() *Calculator {
	return &Calculator{}
}

// isVotableGrade determines if a grade entry represents a valid evaluable score in the Italian 1-10 scale.
// Valid grades must be between 1.0 and 10.0.
// Unset/unrated grades (val < 1.0 or val == 0.0) or absences (< 0) are excluded.
func isVotableGrade(val float64) bool {
	return val >= 1.0 && val <= 10.0
}

// CalculateAverage computes the arithmetic mean
func (c *Calculator) CalculateAverage(grades []Grade) float64 {
	vals := c.extractValues(grades)
	if len(vals) == 0 {
		return 0
	}
	var total float64
	for _, v := range vals {
		total += v
	}
	return math.Round((total/float64(len(vals)))*100) / 100 // Round to 2 decimal places
}

// ConvertJudgmentToValue converts standard and composite Italian judgments to scale 1-10 with case-insensitive matching
func (c *Calculator) ConvertJudgmentToValue(judgment string) float64 {
	cleaned := strings.TrimSpace(strings.ToLower(judgment))
	if cleaned == "" {
		return 0.0
	}
	switch cleaned {
	case "eccellente", "ottimo", "avanzato", "ottimo+", "eccellente+":
		return 10.0
	case "ottimo-", "quasi ottimo", "distinto+":
		return 9.5
	case "distinto":
		return 9.0
	case "distinto-", "quasi distinto", "buono+", "più che buono":
		return 8.5
	case "buono", "intermedio":
		return 8.0
	case "buono-", "discreto+", "più che discreto":
		return 7.5
	case "discreto", "base":
		return 7.0
	case "discreto-", "sufficiente+", "più che sufficiente", "sopra la sufficienza":
		return 6.5
	case "sufficiente":
		return 6.0
	case "sufficiente-", "mediocre+":
		return 5.5
	case "mediocre", "quasi sufficiente", "iniziale":
		return 5.0
	case "mediocre-", "insufficiente+":
		return 4.5
	case "insufficiente", "non raggiunto":
		return 4.0
	case "insufficiente-", "gravemente insufficiente+":
		return 3.5
	case "gravemente insufficiente":
		return 3.0
	}

	// Fallback check for dynamic prefixes/suffixes
	if strings.HasSuffix(cleaned, "+") {
		baseVal := c.ConvertJudgmentToValue(strings.TrimSuffix(cleaned, "+"))
		if baseVal > 0 {
			return math.Min(10.0, baseVal+0.25)
		}
	} else if strings.HasSuffix(cleaned, "-") {
		baseVal := c.ConvertJudgmentToValue(strings.TrimSuffix(cleaned, "-"))
		if baseVal > 0 {
			return math.Max(1.0, baseVal-0.25)
		}
	} else if strings.HasPrefix(cleaned, "più che ") {
		baseVal := c.ConvertJudgmentToValue(strings.TrimPrefix(cleaned, "più che "))
		if baseVal > 0 {
			return math.Min(10.0, baseVal+0.5)
		}
	} else if strings.HasPrefix(cleaned, "quasi ") {
		baseVal := c.ConvertJudgmentToValue(strings.TrimPrefix(cleaned, "quasi "))
		if baseVal > 0 {
			return math.Max(1.0, baseVal-0.5)
		}
	}

	return 0.0
}

// CalculateWeightedAverage computes the weighted mean.
// If no explicit non-zero weights are configured across all grades, it falls back to the simple arithmetic average.
func (c *Calculator) CalculateWeightedAverage(grades []Grade) float64 {
	if len(grades) == 0 {
		return 0
	}

	var totalWeighted float64
	var totalWeights float64
	var hasNonZeroWeight bool

	for _, g := range grades {
		val := g.GradeValue
		if val == 0 && g.GradeType == GradeTypeJudgment {
			val = c.ConvertJudgmentToValue(g.Description)
		}
		if isVotableGrade(val) {
			if g.Weight > 0 {
				totalWeighted += val * g.Weight
				totalWeights += g.Weight
				hasNonZeroWeight = true
			}
		}
	}

	if !hasNonZeroWeight {
		return c.CalculateAverage(grades)
	}
	if totalWeights == 0 {
		return 0
	}

	return math.Round((totalWeighted/totalWeights)*100) / 100
}

// --- Advanced Statistics ---

func (c *Calculator) CalculateMedian(grades []Grade) float64 {
	vals := c.extractValues(grades)
	if len(vals) == 0 {
		return 0
	}

	sort.Float64s(vals)
	n := len(vals)
	if n%2 == 1 {
		return vals[n/2]
	}
	return (vals[n/2-1] + vals[n/2]) / 2
}

func (c *Calculator) calculateStandardDeviationFromValues(vals []float64, mean float64) float64 {
	if len(vals) < 2 {
		return 0
	}
	var varianceSum float64
	for _, v := range vals {
		varianceSum += math.Pow(v-mean, 2)
	}
	return math.Sqrt(varianceSum / float64(len(vals)-1))
}

func (c *Calculator) CalculateStandardDeviation(grades []Grade) float64 {
	vals := c.extractValues(grades)
	if len(vals) < 2 {
		return 0
	}

	var sum float64
	for _, v := range vals {
		sum += v
	}
	mean := sum / float64(len(vals))
	return c.calculateStandardDeviationFromValues(vals, mean)
}

func (c *Calculator) CalculateBellCurve(grades []Grade) (mean, stdDev, skewness, kurtosis float64) {
	vals := c.extractValues(grades)
	n := float64(len(vals))
	if n < 2 {
		return 0, 0, 0, 0
	}

	var sum float64
	for _, v := range vals {
		sum += v
	}
	mean = sum / n
	stdDev = c.calculateStandardDeviationFromValues(vals, mean)

	if stdDev == 0 || n < 3 {
		return mean, stdDev, 0, 0
	}

	var skewSum, kurtSum float64
	for _, v := range vals {
		z := (v - mean) / stdDev
		skewSum += math.Pow(z, 3)
		kurtSum += math.Pow(z, 4)
	}

	// Adjusted sample skewness (Fisher-Pearson coefficient)
	skewness = (n / ((n - 1) * (n - 2))) * skewSum

	// Adjusted sample excess kurtosis (Joanes-Gill estimator).
	// Requires n > 3; for n == 3 the formula is undefined (division by zero
	// in the correction term). Return 0 rather than mixing population and
	// sample estimators.
	if n > 3 {
		kurtosis = ((n*(n+1))/((n-1)*(n-2)*(n-3)))*kurtSum - ((3 * math.Pow(n-1, 2)) / ((n - 2) * (n - 3)))
	} else {
		// n == 3: insufficient degrees of freedom for the sample excess kurtosis.
		kurtosis = 0
	}

	return
}

func (c *Calculator) CalculatePercentile(grades []Grade, score float64) float64 {
	if score < 1.0 || score > 10.0 {
		return 0
	}
	vals := c.extractValues(grades)
	if len(vals) == 0 {
		return 0
	}

	countBelow := 0
	countEqual := 0
	for _, v := range vals {
		if v < score {
			countBelow++
		} else if v == score {
			countEqual++
		}
	}

	return ((float64(countBelow) + 0.5*float64(countEqual)) / float64(len(vals))) * 100
}

func (c *Calculator) DetectOutliers(grades []Grade) []string {
	type indexedVal struct {
		id  string
		val float64
	}
	var vals []indexedVal
	for _, g := range grades {
		val := g.GradeValue
		if val == 0 && g.GradeType == GradeTypeJudgment {
			val = c.ConvertJudgmentToValue(g.Description)
		}
		if isVotableGrade(val) {
			vals = append(vals, indexedVal{id: g.ID, val: val})
		}
	}
	if len(vals) < 4 {
		return nil
	}

	// Calculate IQR (Interquartile Range)
	sortedVals := make([]float64, len(vals))
	for i, iv := range vals {
		sortedVals[i] = iv.val
	}
	sort.Float64s(sortedVals)

	n := len(sortedVals)
	q1 := sortedVals[n/4]
	q3 := sortedVals[(3*n)/4]
	iqr := q3 - q1

	var low, high float64
	if iqr > 0 {
		low = q1 - 1.5*iqr
		high = q3 + 1.5*iqr
	} else {
		var sum float64
		for _, v := range sortedVals {
			sum += v
		}
		mean := sum / float64(n)
		stdDev := c.calculateStandardDeviationFromValues(sortedVals, mean)
		if stdDev == 0 {
			return nil
		}
		low = mean - 2.5*stdDev
		high = mean + 2.5*stdDev
	}

	var outliers []string
	for _, iv := range vals {
		if iv.val < low || iv.val > high {
			outliers = append(outliers, iv.id)
		}
	}
	return outliers
}

// Helper to extract non-negative values
func (c *Calculator) extractValues(grades []Grade) []float64 {
	var vals []float64
	for _, g := range grades {
		val := g.GradeValue
		if val == 0 && g.GradeType == GradeTypeJudgment {
			val = c.ConvertJudgmentToValue(g.Description)
		}
		if isVotableGrade(val) {
			vals = append(vals, val)
		}
	}
	return vals
}
