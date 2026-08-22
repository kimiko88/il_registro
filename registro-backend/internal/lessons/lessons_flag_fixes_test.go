package lessons

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLessons_MutualExclusivity_CoTeachingAndSubstitution(t *testing.T) {
	mockRepo := new(MockLessonsRepository)
	svc := NewService(mockRepo)

	req := CreateLessonRequest{
		ClassID:        "class-1",
		SubjectID:      "sub-1",
		Date:           "2026-09-01",
		Hour:           1,
		Duration:       1,
		IsCoTeaching:   true,
		IsSubstitution: true,
	}

	// Should reject payload with error because both flags are set
	_, err := svc.CreateLesson("t-1", req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non è possibile contrassegnare una lezione sia come Compresenza che come Sostituzione")
}

func TestLessons_UpdateLesson_SubstitutionToStandardTransition(t *testing.T) {
	mockRepo := new(MockLessonsRepository)
	svc := NewService(mockRepo)

	existingLesson := &Lesson{
		ID:             "les-1",
		ClassID:        "class-1",
		TeacherID:      "t-1",
		Date:           time.Now(),
		Hour:           1,
		Duration:       1,
		IsSubstitution: true,
	}

	classLessons := []Lesson{
		{
			ID:             "les-2",
			ClassID:        "class-1",
			TeacherID:      "t-2",
			Date:           time.Now(),
			Hour:           1,
			Duration:       1,
			IsSubstitution: false,
		},
	}

	mockRepo.On("GetLessonByID", "les-1").Return(existingLesson, nil).Once()
	mockRepo.On("GetLessonsByClass", "class-1", existingLesson.Date.Format("2006-01-02")).Return(classLessons, nil).Once()

	// Update lesson les-1 with IsSubstitution = &false (transitioning to standard lesson)
	isSubFalse := false
	req := UpdateLessonRequest{
		IsSubstitution: &isSubFalse,
	}

	// Should fail overlap check because les-1 is becoming a standard lesson and collides with les-2
	_, err := svc.UpdateLesson("t-1", "teacher", "les-1", req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "impossibile registrare più lezioni")
	mockRepo.AssertExpectations(t)
}
