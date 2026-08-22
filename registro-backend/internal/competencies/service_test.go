package competencies

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, "comp-1", eval.ID)
	assert.Equal(t, "stud-1", eval.StudentID)
	assert.Equal(t, "COMP_L1_ITA", eval.CompetenceCode)
	assert.Equal(t, "Comunicazione nella madrelingua", eval.CompetenceName)
	assert.Equal(t, "A_Avanzato", eval.Level)
	assert.Equal(t, "Dimostra padronanza elevata", eval.Descriptor)
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

	assert.Equal(t, "stu-123", req.StudentID)
	assert.Equal(t, "cls-456", req.ClassID)
	assert.Equal(t, 1, req.Semester)
	assert.Equal(t, "COMP_STEM", req.CompetenceCode)
	assert.Equal(t, "Competenza matematica e scientifica", req.CompetenceName)
	assert.Equal(t, "B_Intermedio", req.Level)
}

func TestSaveEvaluation_Validation(t *testing.T) {
	svc := NewService(nil)
	ctx := context.Background()

	// Test invalid level
	_, err := svc.SaveEvaluation(ctx, "school-1", "evaluator-1", SaveEvaluationRequest{
		StudentID:      "stu-1",
		ClassID:        "cls-1",
		CompetenceCode: "COMP_1",
		Level:          "INVALID_LEVEL_XYZ",
	})
	assert.Error(t, err)

	// Test missing student_id
	_, err = svc.SaveEvaluation(ctx, "school-1", "evaluator-1", SaveEvaluationRequest{
		StudentID:      "",
		ClassID:        "cls-1",
		CompetenceCode: "COMP_1",
		Level:          "A_Avanzato",
	})
	assert.Error(t, err)
}
