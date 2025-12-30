package unit

import (
	"testing"
	"time"

	"registro-backend/internal/grades"
	"registro-backend/tests/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGradesService_GetMyGrades(t *testing.T) {
	mockRepo := new(testhelpers.MockGradesRepository)
	mockUserRepo := new(testhelpers.MockUsersRepository)

	// DB passed as nil
	service := grades.NewService(mockRepo, mockUserRepo, nil)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("FindByStudent", "student1").Return([]grades.Grade{
			{ID: "1", GradeValue: 8.0, Semester: 1, IsPublished: true, Date: time.Now()},
			{ID: "2", GradeValue: 7.0, Semester: 2, IsPublished: true, Date: time.Now()},
		}, nil)

		resp, err := service.GetMyGrades("student1", grades.GradeFilter{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		if len(resp.Semesters) > 0 {
			assert.Equal(t, "2025-09-01", resp.Semesters[0].StartDate)
		}
	})

	t.Run("ignores unpublished", func(t *testing.T) {
		// Fresh mock
		mockRepo2 := new(testhelpers.MockGradesRepository)
		service2 := grades.NewService(mockRepo2, mockUserRepo, nil)

		mockRepo2.On("FindByStudent", "student1").Return([]grades.Grade{
			{ID: "1", GradeValue: 9.0, IsPublished: false},
		}, nil)

		resp, err := service2.GetMyGrades("student1", grades.GradeFilter{})
		assert.NoError(t, err)
		assert.Empty(t, resp.Semesters)
	})
}

func TestGradesService_Export(t *testing.T) {
	mockRepo := new(testhelpers.MockGradesRepository)
	service := grades.NewService(mockRepo, new(testhelpers.MockUsersRepository), nil)

	t.Run("export csv", func(t *testing.T) {
		mockRepo.On("FindWithFilter", mock.Anything).Return([]grades.Grade{
			{ID: "1", StudentID: "s1", GradeValue: 10, Date: time.Now(), TeacherID: "t1"},
		}, nil)

		data, contentType, err := service.Export("t1", grades.GradeFilter{}, "csv")

		assert.NoError(t, err)
		assert.Equal(t, "text/csv", contentType)
		assert.Contains(t, string(data), "s1,")
		assert.Contains(t, string(data), "10.00")
	})
}
