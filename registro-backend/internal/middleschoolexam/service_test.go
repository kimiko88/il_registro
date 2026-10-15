package middleschoolexam

import (
	"context"
	"testing"
)

func TestService_AdmissionValidation(t *testing.T) {
	repo := newMockExamRepo()
	svc := NewService(repo)
	ctx := context.Background()

	// Grade < 6 must fail
	err := svc.SaveAdmission(ctx, "exam-1", "student-1", 5, "Insufficiente", false)
	if err == nil {
		t.Fatal("expected error for admission grade < 6")
	}

	// Grade > 10 must fail
	err = svc.SaveAdmission(ctx, "exam-1", "student-1", 11, "Superiore", true)
	if err == nil {
		t.Fatal("expected error for admission grade > 10")
	}

	// Grade 8 must succeed
	err = svc.SaveAdmission(ctx, "exam-1", "student-1", 8, "Buono", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_EvaluateCandidateWorkflow(t *testing.T) {
	repo := newMockExamRepo()
	svc := NewService(repo)
	ctx := context.Background()

	// Set candidate admission
	_ = svc.SaveAdmission(ctx, "exam-1", "student-1", 9, "Ottimo percorso", true)

	// Evaluate candidate
	grades := ExamCandidateGrades{
		GradeItalian:    9.0,
		GradeMath:       10.0,
		GradeEnglish:    9.0,
		GradeSecondLang: 10.0,
		GradeInterview:  10.0,
		ProposedHonors:  false,
	}

	cand, err := svc.EvaluateCandidate(ctx, "exam-1", "student-1", grades, "Colloquio interdisciplinare brillante")
	if err != nil {
		t.Fatalf("unexpected evaluate candidate error: %v", err)
	}

	if cand.ExamMean < 9.5 {
		t.Errorf("expected exam mean >= 9.5, got %f", cand.ExamMean)
	}
	if cand.FinalGrade != 9 {
		t.Errorf("expected final grade 9 (avg (9 + 9.6)/2 = 9.3 -> 9), got %d", cand.FinalGrade)
	}
	if cand.Outcome != "licenziato" {
		t.Errorf("expected outcome licenziato, got %s", cand.Outcome)
	}
}

func TestService_GenerateDiploma(t *testing.T) {
	repo := newMockExamRepo()
	svc := NewService(repo)
	ctx := context.Background()

	text, err := svc.GenerateDiploma(ctx, "cand-1")
	if err != nil {
		t.Fatalf("unexpected diploma error: %v", err)
	}
	if text == "" {
		t.Fatal("diploma text should not be empty")
	}
}

func TestService_UpdateStatus(t *testing.T) {
	repo := newMockExamRepo()
	svc := NewService(repo)
	ctx := context.Background()

	_, _ = svc.GetOrCreateExam(ctx, "school-1", "cls-1", "2026/2027", "Presidente")
	err := svc.UpdateStatus(ctx, "exam-1", "completed")
	if err != nil {
		t.Fatalf("unexpected error updating status: %v", err)
	}
}
