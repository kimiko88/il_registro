package middleschoolexam

import (
	"testing"
)

func TestCalculateMiddleSchoolExamResult(t *testing.T) {
	tests := []struct {
		name              string
		admissionGrade    int
		gradeItalian      float64
		gradeMath         float64
		gradeEnglish      float64
		gradeSecondLang   float64
		gradeInterview    float64
		proposedHonors    bool
		expectedExamMean  float64
		expectedFinal     int
		expectedOutcome   string
		expectedHasHonors bool
	}{
		{
			name:              "Clean 10 with honors",
			admissionGrade:    10,
			gradeItalian:      10,
			gradeMath:         10,
			gradeEnglish:      10,
			gradeSecondLang:   10,
			gradeInterview:    10,
			proposedHonors:    true,
			expectedExamMean:  10.00,
			expectedFinal:     10,
			expectedOutcome:   "licenziato",
			expectedHasHonors: true,
		},
		{
			name:              "Rounding up: admission 7, exam mean 7.6 -> average 7.3 -> 7",
			admissionGrade:    7,
			gradeItalian:      8,
			gradeMath:         7,
			gradeEnglish:      8,
			gradeSecondLang:   7,
			gradeInterview:    8,
			proposedHonors:    false,
			expectedExamMean:  7.60,
			expectedFinal:     7, // (7 + 7.6)/2 = 7.3 -> 7
			expectedOutcome:   "licenziato",
			expectedHasHonors: false,
		},
		{
			name:              "Rounding up from >= 0.5: admission 8, exam mean 7.0 -> average 7.5 -> 8",
			admissionGrade:    8,
			gradeItalian:      7,
			gradeMath:         7,
			gradeEnglish:      7,
			gradeSecondLang:   7,
			gradeInterview:    7,
			proposedHonors:    false,
			expectedExamMean:  7.00,
			expectedFinal:     8, // (8 + 7.0)/2 = 7.5 -> rounds up to 8
			expectedOutcome:   "licenziato",
			expectedHasHonors: false,
		},
		{
			name:              "Honors rejected if final grade is not 10",
			admissionGrade:    9,
			gradeItalian:      9,
			gradeMath:         9,
			gradeEnglish:      9,
			gradeSecondLang:   9,
			gradeInterview:    9,
			proposedHonors:    true, // cannot have honors if not 10
			expectedExamMean:  9.00,
			expectedFinal:     9,
			expectedOutcome:   "licenziato",
			expectedHasHonors: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			candidate := ExamCandidateGrades{
				AdmissionGrade:  tc.admissionGrade,
				GradeItalian:    tc.gradeItalian,
				GradeMath:       tc.gradeMath,
				GradeEnglish:    tc.gradeEnglish,
				GradeSecondLang: tc.gradeSecondLang,
				GradeInterview:  tc.gradeInterview,
				ProposedHonors:  tc.proposedHonors,
			}

			res := CalculateExamOutcome(candidate)
			if res.ExamMean != tc.expectedExamMean {
				t.Errorf("expected exam mean %.2f, got %.2f", tc.expectedExamMean, res.ExamMean)
			}
			if res.FinalGrade != tc.expectedFinal {
				t.Errorf("expected final grade %d, got %d", tc.expectedFinal, res.FinalGrade)
			}
			if res.Outcome != tc.expectedOutcome {
				t.Errorf("expected outcome %s, got %s", tc.expectedOutcome, res.Outcome)
			}
			if res.HasHonors != tc.expectedHasHonors {
				t.Errorf("expected has honors %v, got %v", tc.expectedHasHonors, res.HasHonors)
			}
		})
	}
}
