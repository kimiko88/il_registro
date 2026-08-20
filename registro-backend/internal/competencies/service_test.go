package competencies

import (
	"testing"
)

func TestCompetenceEvaluationStruct(t *testing.T) {
	eval := Evaluation{
		ID:             "comp-1",
		StudentID:      "stud-1",
		CompetenceCode: "COMP_L1_ITA",
		CompetenceName: "Comunicazione nella madrelingua",
		Level:          "A_Avanzato",
		Descriptor:     "Dimostra padronanza elevata",
	}

	if eval.Level != "A_Avanzato" {
		t.Fatalf("expected level A_Avanzato, got %s", eval.Level)
	}
	if eval.CompetenceCode != "COMP_L1_ITA" {
		t.Fatalf("expected code COMP_L1_ITA, got %s", eval.CompetenceCode)
	}
}

func TestSaveEvaluationRequest(t *testing.T) {
	req := SaveEvaluationRequest{
		StudentID:      "stu-123",
		ClassID:        "cls-456",
		Semester:       1,
		CompetenceCode: "COMP_STEM",
		CompetenceName: "Competenza matematica e scientifica",
		Level:          "B_Intermedio",
	}

	if req.StudentID != "stu-123" || req.Level != "B_Intermedio" {
		t.Fatalf("invalid request fields")
	}
}

func TestSaveEvaluation_Validation(t *testing.T) {
	svc := NewService(nil)

	// Test invalid level
	_, err := svc.SaveEvaluation(nil, "school-1", "evaluator-1", SaveEvaluationRequest{
		StudentID:      "stu-1",
		ClassID:        "cls-1",
		CompetenceCode: "COMP_1",
		Level:          "INVALID_LEVEL_XYZ",
	})
	if err == nil {
		t.Fatalf("expected error for invalid competency level, got nil")
	}

	// Test missing student_id
	_, err = svc.SaveEvaluation(nil, "school-1", "evaluator-1", SaveEvaluationRequest{
		StudentID:      "",
		ClassID:        "cls-1",
		CompetenceCode: "COMP_1",
		Level:          "A_Avanzato",
	})
	if err == nil {
		t.Fatalf("expected error for missing student_id, got nil")
	}
}
