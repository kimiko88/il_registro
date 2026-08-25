package agenda

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, item *AgendaItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}
func (m *MockRepository) GetByID(ctx context.Context, id string) (*AgendaItem, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AgendaItem), args.Error(1)
}
func (m *MockRepository) Update(ctx context.Context, item *AgendaItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}
func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockRepository) ListCalendar(ctx context.Context, schoolID string, filter CalendarFilter) ([]*AgendaItem, error) {
	args := m.Called(ctx, schoolID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*AgendaItem), args.Error(1)
}
func (m *MockRepository) SetCompletion(ctx context.Context, itemID, studentID string, completed bool) error {
	args := m.Called(ctx, itemID, studentID, completed)
	return args.Error(0)
}
func (m *MockRepository) IsStudentInClass(ctx context.Context, studentID, classID string) (bool, error) {
	args := m.Called(ctx, studentID, classID)
	return args.Bool(0), args.Error(1)
}

func TestCreateAgendaItem(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	req := CreateAgendaItemRequest{
		ClassID:     "class-1",
		Title:       "Math Homework",
		Description: "Pages 40-45",
		Type:        TypeHomework,
		Date:        "2025-11-15",
	}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(item *AgendaItem) bool {
		return item.Title == "Math Homework" && item.ClassID == "class-1"
	})).Return(nil).Once()

	item, err := svc.CreateAgendaItem(context.Background(), "teacher-1", "school-1", req)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "Math Homework", item.Title)
	mockRepo.AssertExpectations(t)
}

func TestTaskCompletion(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	item := &AgendaItem{ID: "item-1", TeacherID: "t-1"}
	mockRepo.On("GetByID", mock.Anything, "item-1").Return(item, nil).Twice()
	mockRepo.On("SetCompletion", mock.Anything, "item-1", "student-1", true).Return(nil).Once()
	mockRepo.On("SetCompletion", mock.Anything, "item-1", "student-1", false).Return(nil).Once()

	err := svc.SetTaskCompletion(context.Background(), "student-1", "student", "student-1", "item-1", true)
	assert.NoError(t, err)

	err = svc.SetTaskCompletion(context.Background(), "student-1", "student", "student-1", "item-1", false)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCalendarQuery(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	expectedItems := []*AgendaItem{
		{ID: "1", Title: "Verifica Fisica"},
	}

	mockRepo.On("ListCalendar", mock.Anything, "school-1", mock.Anything).Return(expectedItems, nil).Once()

	items, err := svc.GetCalendar(context.Background(), "school-1", "user-1", "student", CalendarFilter{
		From: time.Now(),
		To:   time.Now().AddDate(0, 1, 0),
	})

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "Verifica Fisica", items[0].Title)
}
