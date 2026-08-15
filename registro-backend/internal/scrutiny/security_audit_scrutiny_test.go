package scrutiny

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScrutiny_SaveDeficiency_RBACRejection(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)

	ctx := context.Background()

	// Student role attempting to save deficiency MUST be rejected
	errStudent := svc.SaveDeficiency(ctx, "stu-1", "student", &SaveDeficiencyRequest{
		StudentID: "stu-1",
		ClassID:   "cls-1",
		SubjectID: "sub-1",
	})
	assert.Error(t, errStudent)
	assert.Contains(t, errStudent.Error(), "forbidden")

	// Parent role attempting to save deficiency MUST be rejected
	errParent := svc.SaveDeficiency(ctx, "parent-1", "parent", &SaveDeficiencyRequest{
		StudentID: "stu-1",
		ClassID:   "cls-1",
		SubjectID: "sub-1",
	})
	assert.Error(t, errParent)
	assert.Contains(t, errParent.Error(), "forbidden")
}

func TestScrutiny_GetClassDeficiencies_RBACRejection(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)

	ctx := context.Background()

	// Unprivileged student attempting to access class deficiencies MUST be rejected
	_, err := svc.GetClassDeficiencies(ctx, "stu-1", "student", "cls-1", 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}
