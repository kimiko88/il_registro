package grades

import (
	"math"
	"sort"
)

type Calculator struct{}

func NewCalculator() *Calculator {
	return &Calculator{}
}

// CalculateAverage computes the arithmetic mean
func (c *Calculator) CalculateAverage(grades []Grade) float64 {
	if len(grades) == 0 {
		return 0
	}

	var total float64
	var count int

	for _, g := range grades {
		val := g.GradeValue
		if val == 0 && g.GradeType == GradeTypeJudgment {
			val = c.ConvertJudgmentToValue(g.Description)
		}
		if val > 0 {
			total += val
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return math.Round((total/float64(count))*100) / 100 // Round to 2 decimal places
}

// ConvertJudgmentToValue converts standard Italian judgments to scale 0-10
func (c *Calculator) ConvertJudgmentToValue(judgment string) float64 {
	switch judgment {
	case "Eccellente", "Ottimo":
		return 10.0
	case "Distinto":
		return 9.0
	case "Buono":
		return 8.0
	case "Discreto":
		return 7.0
	case "Sufficiente":
		return 6.0
	case "Mediocre", "Quasi Sufficiente":
		return 5.0
	case "Insufficiente":
		return 4.0
	case "Gravemente Insufficiente":
		return 3.0
	default:
		return 0.0
	}
}

// CalculateWeightedAverage computes the weighted mean
func (c *Calculator) CalculateWeightedAverage(grades []Grade) float64 {
	if len(grades) == 0 {
		return 0
	}

	var totalWeighted float64
	var totalWeights float64

	for _, g := range grades {
		val := g.GradeValue
		if val == 0 && g.GradeType == GradeTypeJudgment {
			val = c.ConvertJudgmentToValue(g.Description)
		}
		if val > 0 && g.Weight > 0 {
			totalWeighted += val * g.Weight
			totalWeights += g.Weight
		}
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

	var varianceSum float64
	for _, v := range vals {
		varianceSum += math.Pow(v-mean, 2)
	}
	return math.Sqrt(varianceSum / float64(len(vals)-1))
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
	stdDev = c.CalculateStandardDeviation(grades)

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
	// sample estimators, which was the previous inconsistency.
	if n > 3 {
		kurtosis = ((n*(n+1))/((n-1)*(n-2)*(n-3)))*kurtSum - ((3 * math.Pow(n-1, 2)) / ((n - 2) * (n - 3)))
	} else {
		// n == 3: insufficient degrees of freedom for the sample excess kurtosis.
		kurtosis = 0
	}

	return
}

func (c *Calculator) CalculatePercentile(grades []Grade, score float64) float64 {
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
	// Returns IDs of outlier grades (outside Mean +/- 2*StdDev)
	mean := c.CalculateAverage(grades)
	stdDev := c.CalculateStandardDeviation(grades)
	if stdDev == 0 {
		return nil
	}
	low := mean - 2*stdDev
	high := mean + 2*stdDev

	var outliers []string
	for _, g := range grades {
		val := g.GradeValue
		if val == 0 && g.GradeType == GradeTypeJudgment {
			val = c.ConvertJudgmentToValue(g.Description)
		}
		if val <= 0 {
			continue // Skip unrated / invalid non-positive grades
		}

		if val < low || val > high {
			outliers = append(outliers, g.ID)
		}
	}
	return outliers
}

// Helper to extract non-zero values
func (c *Calculator) extractValues(grades []Grade) []float64 {
	var vals []float64
	for _, g := range grades {
		val := g.GradeValue
		if val == 0 && g.GradeType == GradeTypeJudgment {
			val = c.ConvertJudgmentToValue(g.Description)
		}
		if val > 0 {
			vals = append(vals, val)
		}
	}
	return vals
}
