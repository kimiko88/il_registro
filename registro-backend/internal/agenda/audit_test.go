package agenda

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCalendar_ParentNoStudentID_Rejected(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	_, err := svc.GetCalendar(context.Background(), "school-1", "parent-1", "parent", CalendarFilter{
		From: time.Now(),
		To:   time.Now().AddDate(0, 1, 0),
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "student_id is required")
}

func TestNewService_NilUserRepo_ParentGetCalendar_Rejected(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo) // userRepo is nil

	_, err := svc.GetCalendar(context.Background(), "school-1", "parent-1", "parent", CalendarFilter{
		StudentID: "student-1",
		From:      time.Now(),
		To:        time.Now().AddDate(0, 1, 0),
	})
	assert.Error(t, err)
	assert.Equal(t, ErrUnauthorized, err)
}

func TestCreateAgendaItem_InvalidTimeFormat_Rejected(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	req := CreateAgendaItemRequest{
		ClassID:     "class-1",
		Title:       "Math Homework",
		Type:        TypeHomework,
		Date:        "2025-11-15",
		StartTime:   "invalid-time",
		EndTime:     "10:00",
		AllDay:      false,
	}

	_, err := svc.CreateAgendaItem(context.Background(), "teacher-1", "school-1", req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid start_time format")
}

func TestSetTaskCompletion_StudentNotInClass_Rejected(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	item := &AgendaItem{ID: "item-1", ClassID: "class-1", TeacherID: "t-1"}
	mockRepo.On("GetByID", mock.Anything, "item-1").Return(item, nil).Once()
	mockRepo.On("IsStudentInClass", mock.Anything, "student-1", "class-1").Return(false, nil).Once()

	err := svc.SetTaskCompletion(context.Background(), "student-1", "student", "student-1", "item-1", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non appartiene alla classe")
}

func TestCreateAgendaItem_MissingTitle_Rejected(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	req := CreateAgendaItemRequest{
		ClassID:   "class-1",
		Title:     "", // Empty title
		Type:      TypeHomework,
		Date:      "2025-11-15",
		StartTime: "09:00",
		EndTime:   "10:00",
	}

	_, err := svc.CreateAgendaItem(context.Background(), "teacher-1", "school-1", req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title is required")
}

func TestUpdateAgendaItem_NonOwnerTeacher_Rejected(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	item := &AgendaItem{ID: "item-1", SchoolID: "school-1", TeacherID: "teacher-creator"}
	mockRepo.On("GetByID", mock.Anything, "item-1").Return(item, nil).Once()

	newTitle := "Updated Title"
	req := UpdateAgendaItemRequest{Title: &newTitle}

	// Another teacher attempting update
	_, err := svc.UpdateAgendaItem(context.Background(), "teacher-other", "teacher", "school-1", "item-1", req)
	assert.Error(t, err)
	assert.Equal(t, ErrUnauthorized, err)
}
