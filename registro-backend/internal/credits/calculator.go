package credits

import (
	"fmt"
	"math"
)

// CalculateCreditRange computes the official ministerial credit range for classes III, IV, V (D.Lgs. 62/2017 Tabella A).
func CalculateCreditRange(gradeLevel int, average float64, conductGrade int, pctoHours int, hasExtracurricular bool) CreditCalculationResult {
	if gradeLevel < 3 || gradeLevel > 5 {
		gradeLevel = 3
	}

	// Clamp average between 6.0 and 10.0
	avg := math.Round(average*100) / 100
	if avg < 6.0 {
		return CreditCalculationResult{
			GradeLevel:         gradeLevel,
			GradeAverage:       avg,
			ConductGrade:       conductGrade,
			BaseCreditRangeMin: 0,
			BaseCreditRangeMax: 0,
			SuggestedCredit:    0,
			MaxYearCredit:      getMaxYearCredit(gradeLevel),
			Motivation:         "Media inferiore a 6.00: nessun credito attribuibile",
		}
	}

	var minCredit, maxCredit int

	switch gradeLevel {
	case 3:
		// Classe Terza (Max 12)
		switch {
		case avg == 6.00:
			minCredit, maxCredit = 7, 8
		case avg > 6.00 && avg <= 7.00:
			minCredit, maxCredit = 8, 9
		case avg > 7.00 && avg <= 8.00:
			minCredit, maxCredit = 9, 10
		case avg > 8.00 && avg <= 9.00:
			minCredit, maxCredit = 10, 11
		case avg > 9.00:
			minCredit, maxCredit = 11, 12
		}
	case 4:
		// Classe Quarta (Max 13)
		switch {
		case avg == 6.00:
			minCredit, maxCredit = 8, 9
		case avg > 6.00 && avg <= 7.00:
			minCredit, maxCredit = 9, 10
		case avg > 7.00 && avg <= 8.00:
			minCredit, maxCredit = 10, 11
		case avg > 8.00 && avg <= 9.00:
			minCredit, maxCredit = 11, 12
		case avg > 9.00:
			minCredit, maxCredit = 12, 13
		}
	case 5:
		// Classe Quinta (Max 15)
		switch {
		case avg == 6.00:
			minCredit, maxCredit = 9, 10
		case avg > 6.00 && avg <= 7.00:
			minCredit, maxCredit = 10, 11
		case avg > 7.00 && avg <= 8.00:
			minCredit, maxCredit = 11, 12
		case avg > 8.00 && avg <= 9.00:
			minCredit, maxCredit = 12, 13
		case avg > 9.00:
			minCredit, maxCredit = 14, 15
		}
	}

	// Determine suggested score in range [minCredit, maxCredit]
	suggested := minCredit
	motivation := fmt.Sprintf("Fascia base D.Lgs 62/2017: %d - %d punti", minCredit, maxCredit)

	// Criteria for max of the band: high decimal portion (>= .50), high conduct (>= 9), or valid extracurricular/PCTO
	decimalPart := avg - math.Floor(avg)
	if decimalPart >= 0.50 || conductGrade >= 9 || (pctoHours >= 30 && hasExtracurricular) {
		suggested = maxCredit
		motivation = fmt.Sprintf("Attribuzione punteggio massimo della fascia (%d punti) per merito, condotta (%d) o crediti formativi", maxCredit, conductGrade)
	}

	return CreditCalculationResult{
		GradeLevel:         gradeLevel,
		GradeAverage:       avg,
		ConductGrade:       conductGrade,
		BaseCreditRangeMin: minCredit,
		BaseCreditRangeMax: maxCredit,
		SuggestedCredit:    suggested,
		MaxYearCredit:      getMaxYearCredit(gradeLevel),
		Motivation:         motivation,
	}
}

func getMaxYearCredit(gradeLevel int) int {
	switch gradeLevel {
	case 3:
		return 12
	case 4:
		return 13
	case 5:
		return 15
	default:
		return 12
	}
}
