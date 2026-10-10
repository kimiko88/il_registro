package middleschoolexam

import (
	"math"
	"testing"
)

func FuzzCalculateExamOutcome(f *testing.F) {
	// Seed realistic and boundary cases
	f.Add(6, 6.0, 6.0, 6.0, 6.0, 6.0, false)
	f.Add(10, 10.0, 10.0, 10.0, 10.0, 10.0, true)
	f.Add(6, 4.0, 5.0, 4.5, 5.0, 4.0, false)
	f.Add(8, 7.5, 8.25, 9.0, 6.5, 7.0, false)
	f.Add(10, 10.0, 10.0, 10.0, 10.0, 10.0, false)
	f.Add(0, -1.0, 15.0, 999.0, -50.0, 0.0, true)

	f.Fuzz(func(t *testing.T, admissionGrade int, ita, mathVal, eng, secondLang, interview float64, proposedHonors bool) {
		// Ignore NaN or infinite floats that would break calculation
		if math.IsNaN(ita) || math.IsNaN(mathVal) || math.IsNaN(eng) || math.IsNaN(secondLang) || math.IsNaN(interview) ||
			math.IsInf(ita, 0) || math.IsInf(mathVal, 0) || math.IsInf(eng, 0) || math.IsInf(secondLang, 0) || math.IsInf(interview, 0) {
			return
		}

		c := ExamCandidateGrades{
			AdmissionGrade:  admissionGrade,
			GradeItalian:    ita,
			GradeMath:       mathVal,
			GradeEnglish:    eng,
			GradeSecondLang: secondLang,
			GradeInterview:  interview,
			ProposedHonors:  proposedHonors,
		}

		res := CalculateExamOutcome(c)

		// Invariants:
		// 1. FinalGrade must be in [0, 10]
		if res.FinalGrade < 0 || res.FinalGrade > 10 {
			t.Fatalf("finalGrade out of bounds: %d", res.FinalGrade)
		}

		// 2. Honors is only permitted if FinalGrade == 10 and proposedHonors is true
		if res.HasHonors && (res.FinalGrade != 10 || !proposedHonors) {
			t.Fatalf("honors invariant violated: finalGrade=%d, proposed=%v, honors=%v", res.FinalGrade, proposedHonors, res.HasHonors)
		}

		// 3. Outcome must be "licenziato" if finalGrade >= 6, else "non_licenziato"
		if res.FinalGrade >= 6 && res.Outcome != "licenziato" {
			t.Fatalf("expected licenziato for finalGrade %d, got %s", res.FinalGrade, res.Outcome)
		}
		if res.FinalGrade < 6 && res.Outcome != "non_licenziato" {
			t.Fatalf("expected non_licenziato for finalGrade %d, got %s", res.FinalGrade, res.Outcome)
		}
	})
}
