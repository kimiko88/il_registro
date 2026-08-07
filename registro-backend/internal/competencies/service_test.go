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
