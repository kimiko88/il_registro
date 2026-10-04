package middleschoolexam

import (
	"math"
)

// CalculateExamOutcome applies D.Lgs. 62/2017 & D.M. 741/2017 rules for middle school state exam.
func CalculateExamOutcome(c ExamCandidateGrades) ExamOutcomeResult {
	// 1. Mean of 5 exam tests
	examSum := c.GradeItalian + c.GradeMath + c.GradeEnglish + c.GradeSecondLang + c.GradeInterview
	examMean := math.Round((examSum/5.0)*100) / 100

	// 2. Average of admission grade and exam mean
	avg := (float64(c.AdmissionGrade) + examMean) / 2.0

	// 3. Round to upper integer if decimal fraction >= 0.5
	floorVal := math.Floor(avg)
	fraction := avg - floorVal
	finalGrade := int(floorVal)
	if fraction >= 0.5 {
		finalGrade++
	}

	if finalGrade > 10 {
		finalGrade = 10
	}
	if finalGrade < 0 {
		finalGrade = 0
	}

	// 4. Outcome
	outcome := "non_licenziato"
	if finalGrade >= 6 {
		outcome = "licenziato"
	}

	// 5. Honors (Lode) only permitted with 10/10 final grade and unanimous subcommission vote
	hasHonors := false
	if finalGrade == 10 && c.ProposedHonors {
		hasHonors = true
	}

	return ExamOutcomeResult{
		ExamMean:   examMean,
		FinalGrade: finalGrade,
		HasHonors:  hasHonors,
		Outcome:    outcome,
	}
}
