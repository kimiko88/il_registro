package middleschoolexam

import (
	"testing"
)

func BenchmarkCalculateExamOutcome(b *testing.B) {
	c := ExamCandidateGrades{
		AdmissionGrade:  9,
		GradeItalian:    9.0,
		GradeMath:       8.5,
		GradeEnglish:    9.0,
		GradeSecondLang: 8.5,
		GradeInterview:  9.5,
		ProposedHonors:  false,
	}

	for b.Loop() {
		_ = CalculateExamOutcome(c)
	}
}

func BenchmarkCalculateExamOutcomeBatch(b *testing.B) {
	candidates := make([]ExamCandidateGrades, 100)
	for i := 0; i < 100; i++ {
		candidates[i] = ExamCandidateGrades{
			AdmissionGrade:  6 + (i % 5),
			GradeItalian:    6.0 + float64(i%5)*0.8,
			GradeMath:       6.0 + float64((i+1)%5)*0.8,
			GradeEnglish:    7.0,
			GradeSecondLang: 7.5,
			GradeInterview:  8.0,
			ProposedHonors:  i%10 == 0,
		}
	}

	for b.Loop() {
		for _, cand := range candidates {
			_ = CalculateExamOutcome(cand)
		}
	}
}
