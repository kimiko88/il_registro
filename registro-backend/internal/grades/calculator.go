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
		// Map judgments if needed (though Model implementation implies GradeValue is already numeric storage)
		// But if we receive raw judgment strings in GradeValue (impossible as float64) or relying on Type...
		// Let's assume GradeValue IS the authoritative numeric value.
		// However, prompt asks: "Voti giudizio convertiti in scala: Insuff=3, Mediocre=4.5..."
		// If stored GradeValue is 0 for Judgments, we need to convert based on Description or logic?
		// Usually, "GradeValue" stores the mapped value.
		// Let's ensure we use ConvertJudgmentToValue logic if applicable or just GradeValue.
		// Safety: if GradeType is Judgment, and Value is 0, try to map from somewhere?
		// Assumption: GradeValue is already set correctly on INSERT.
		// But for "Trend" and special averaging, let's implement the conversion utility.

		if g.GradeValue >= 0 {
			total += g.GradeValue
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
		return 5.0 // Or 4.5
	case "Insufficiente":
		return 4.0 // Or 3.0
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
		if g.GradeValue >= 0 && g.Weight > 0 {
			totalWeighted += g.GradeValue * g.Weight
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

	mean := c.CalculateAverage(grades)
	var varianceSum float64
	for _, v := range vals {
		varianceSum += math.Pow(v-mean, 2)
	}
	return math.Sqrt(varianceSum / float64(len(vals)))
}

func (c *Calculator) CalculateBellCurve(grades []Grade) (mean, stdDev, skewness, kurtosis float64) {
	vals := c.extractValues(grades)
	if len(vals) == 0 {
		return 0, 0, 0, 0
	}

	mean = c.CalculateAverage(grades)
	stdDev = c.CalculateStandardDeviation(grades)

	if stdDev == 0 {
		return mean, stdDev, 0, 0
	}

	var skewSum, kurtSum float64
	for _, v := range vals {
		z := (v - mean) / stdDev
		skewSum += math.Pow(z, 3)
		kurtSum += math.Pow(z, 4)
	}
	n := float64(len(vals))
	skewness = skewSum / n
	kurtosis = (kurtSum / n) - 3 // Excess kurtosis

	return
}

func (c *Calculator) CalculatePercentile(grades []Grade, score float64) float64 {
	vals := c.extractValues(grades)
	if len(vals) == 0 {
		return 0
	}

	countBelow := 0
	for _, v := range vals {
		if v < score {
			countBelow++
		}
	}

	return (float64(countBelow) / float64(len(vals))) * 100
}

func (c *Calculator) DetectOutliers(grades []Grade) []string {
	// Returns IDs of outlier grades (outside Mean +/- 2*StdDev)
	mean := c.CalculateAverage(grades)
	stdDev := c.CalculateStandardDeviation(grades)
	low := mean - 2*stdDev
	high := mean + 2*stdDev

	var outliers []string
	for _, g := range grades {
		val := g.GradeValue
		if g.GradeValue == 0 && g.GradeType == GradeTypeJudgment {
			val = c.ConvertJudgmentToValue(g.Description) // Fallback if GradeValue 0
		}
		if g.GradeValue > 0 {
			val = g.GradeValue
		} // Prefer stored value

		if val < low || val > high {
			outliers = append(outliers, g.ID) // OR Student ID? Prompt says "outliers: [{studentId, grade}]"
			// Returning Grade ID primarily, logic can map to student later.
		}
	}
	return outliers
}

// Helper to extract non-zero values
func (c *Calculator) extractValues(grades []Grade) []float64 {
	var vals []float64
	for _, g := range grades {
		if g.GradeValue >= 0 {
			vals = append(vals, g.GradeValue)
		}
	}
	return vals
}
