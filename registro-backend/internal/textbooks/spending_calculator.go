package textbooks

import (
	"math"
)

// CalculateSpendingReport evaluates adoptions against ministerial spending limits and tolerance.
func CalculateSpendingReport(adoptions []ClassAdoptionItem, limit SpendingLimit) SpendingReport {
	total := 0.0
	for _, item := range adoptions {
		// Only count mandatory books that are not already owned
		if item.AdoptionType != "consigliato" && !item.IsAlreadyOwned {
			total += item.Price
		}
	}

	total = math.Round(total*100) / 100
	maxAmount := limit.MaxAmount
	tolerancePct := limit.AllowedTolerancePct
	if tolerancePct <= 0 {
		tolerancePct = 10.0 // standard default 10%
	}

	toleranceThreshold := math.Round(maxAmount*(1.0+tolerancePct/100.0)*100) / 100
	diff := math.Round((total-maxAmount)*100) / 100

	var pct float64
	if maxAmount > 0 {
		pct = math.Round((total/maxAmount)*10000) / 100
	}

	var status SpendingStatus
	if total <= maxAmount {
		status = StatusWithinLimit
	} else if total <= toleranceThreshold {
		status = StatusWarningTolerance
	} else {
		status = StatusOverLimit
	}

	return SpendingReport{
		TotalSpending:       total,
		SpendingLimit:       maxAmount,
		ToleranceThreshold:  toleranceThreshold,
		AllowedTolerancePct: tolerancePct,
		Difference:          diff,
		Percentage:          pct,
		Status:              status,
		AdoptionsCount:      len(adoptions),
		Items:               adoptions,
	}
}
