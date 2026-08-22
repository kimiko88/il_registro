package classes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClassesSecurity_RemoveSubjectCrossSchoolForbidden(t *testing.T) {
	mockRepo := new(MockRepository)
	classObj := &Class{
		ID:       "class-1",
		SchoolID: "school-1",
		Name:     "1A",
	}
	mockRepo.On("Get", mock.Anything, "class-1").Return(classObj, nil)

	svc := NewService(mockRepo)

	// Attempting to remove subject from class-1 using school-2 context
	err := svc.RemoveSubject(context.Background(), "assignment-1", "school-2", "class-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden: class belongs to another school")
}

func TestClassesSecurity_RemoveSubjectSameSchoolAllowed(t *testing.T) {
	mockRepo := new(MockRepository)
	classObj := &Class{
		ID:       "class-1",
		SchoolID: "school-1",
		Name:     "1A",
	}
	mockRepo.On("Get", mock.Anything, "class-1").Return(classObj, nil)
	mockRepo.On("UnassignSubject", mock.Anything, "assignment-1").Return(nil)

	svc := NewService(mockRepo)

	err := svc.RemoveSubject(context.Background(), "assignment-1", "school-1", "class-1")
	assert.NoError(t, err)
}
