package grades

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGradesFixes_CalcSemesterAverages_Weighted(t *testing.T) {
	grades := []GradeResponse{
		{
			GradeValue:    10.0,
			Weight:        2.0,
			Semester:      1,
			GradeCategory: string(GradeCategorySummative),
		},
		{
			GradeValue:    5.0,
			Weight:        1.0,
			Semester:      1,
			GradeCategory: string(GradeCategorySummative),
		},
		{
			GradeValue:    6.0,
			Weight:        1.0,
			Semester:      2,
			GradeCategory: string(GradeCategorySummative),
		},
	}

	// Semester 1: (10*2 + 5*1) / (2+1) = 25 / 3 = 8.3333...
	// Previously unweighted was (10 + 5) / 2 = 7.5
	avg1, avg2 := calcSemesterAverages(grades)
	assert.InDelta(t, 8.33, avg1, 0.01)
	assert.Equal(t, 6.0, avg2)
}

func TestGradesFixes_BatchCreateGrades_Deduplication(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo, nil, nil, nil)

	now := time.Now()
	evalType := EvaluationTypeWritten

	gradesInput := []*Grade{
		{
			StudentID:      "s1",
			SubjectID:      "sub1",
			Date:           now,
			GradeType:      GradeTypeNumeric,
			EvaluationType: &evalType,
			GradeValue:     8.0,
			Description:    "Test 1",
		},
		{
			StudentID:      "s1",
			SubjectID:      "sub1",
			Date:           now,
			GradeType:      GradeTypeNumeric,
			EvaluationType: &evalType,
			GradeValue:     8.0,
			Description:    "Test 1",
		},
	}

	mockRepo.On("BatchCreate", mock.Anything, mock.MatchedBy(func(deduped []*Grade) bool {
		return len(deduped) == 1
	})).Return(nil).Once()

	err := svc.BatchCreateGrades(context.Background(), "t1", "teacher", "sch1", gradesInput)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGradesFixes_GetMyTrend_UsesRepository(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo, nil, nil, nil)

	now := time.Now()
	mockGrades := []Grade{
		{
			ID:          "g1",
			StudentID:   "s1",
			SubjectID:   "sub1",
			GradeValue:  9.0,
			Semester:    1,
			IsPublished: true,
			Date:        now,
		},
	}

	mockRepo.On("FindByStudent", "s1").Return(mockGrades, nil).Once()
	mockRepo.On("GetStudentClassAndSchoolInfo", mock.Anything, "s1").Return("John Doe", "5A", "c100", "sch1", nil).Once()
	mockRepo.On("GetClassSubjectAverage", mock.Anything, "c100", "sub1", 1, "s1").Return(7.5, nil).Once()

	trend, err := svc.GetMyTrend(context.Background(), "s1", "student", "s1", "sub1")
	assert.NoError(t, err)
	assert.NotNil(t, trend)
	assert.NotEmpty(t, trend.Trends)
	assert.Equal(t, 7.5, trend.Trends[0].ClassAverage)
	mockRepo.AssertExpectations(t)
}
