package extracurricular

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListEnrollments_CrossTenant_Blocked(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	course := &Course{
		ID:       "course-1",
		SchoolID: "school-1",
		Title:    "Corso Teatro",
	}
	mockRepo.On("GetCourseByID", mock.Anything, "course-1").Return(course, nil)

	// Teacher from school-2 attempts to view enrollments of school-1 course
	_, err := svc.ListEnrollments(context.Background(), "teacher", "school-2", "course-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestListEnrollments_SameSchool_Allowed(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	course := &Course{
		ID:       "course-1",
		SchoolID: "school-1",
		Title:    "Corso Teatro",
	}
	mockRepo.On("GetCourseByID", mock.Anything, "course-1").Return(course, nil)
	mockRepo.On("ListEnrollments", mock.Anything, "course-1").Return([]*Enrollment{
		{ID: "en-1", CourseID: "course-1", StudentID: "stu-1"},
	}, nil)

	// Teacher from school-1 views enrollments of school-1 course
	res, err := svc.ListEnrollments(context.Background(), "teacher", "school-1", "course-1")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}
