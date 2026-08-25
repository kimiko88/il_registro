package subjects

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubjectsSecurity_UpdateCrossSchoolForbidden(t *testing.T) {
	mockRepo := new(MockRepository)
	existingSubject := &Subject{
		ID:       "subj-1",
		SchoolID: "school-1",
		Name:     "Math",
	}
	mockRepo.On("Get", mock.Anything, "subj-1").Return(existingSubject, nil)

	svc := NewService(mockRepo)

	// Admin of school-2 trying to update subject of school-1
	_, err := svc.UpdateSubject(context.Background(), "subj-1", CreateSubjectRequest{Name: "New Math"}, "school-2")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden: cannot update subject of another school")
}

func TestSubjectsSecurity_DeleteCrossSchoolForbidden(t *testing.T) {
	mockRepo := new(MockRepository)
	existingSubject := &Subject{
		ID:       "subj-1",
		SchoolID: "school-1",
		Name:     "History",
	}
	mockRepo.On("Get", mock.Anything, "subj-1").Return(existingSubject, nil)

	svc := NewService(mockRepo)

	// Admin of school-2 trying to delete subject of school-1
	err := svc.DeleteSubject(context.Background(), "subj-1", "school-2")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden: cannot delete subject of another school")
}
