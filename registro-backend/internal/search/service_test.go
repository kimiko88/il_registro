package search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GlobalSearch(ctx context.Context, actorRole, schoolID, query, filterType string) ([]SearchResultItem, error) {
	args := m.Called(ctx, actorRole, schoolID, query, filterType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]SearchResultItem), args.Error(1)
}

func TestSearchService(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	expected := []SearchResultItem{
		{ID: "1", Type: "student", Title: "Rossi Mario"},
	}

	mockRepo.On("GlobalSearch", mock.Anything, "teacher", "school-1", "Mario", "").Return(expected, nil).Once()

	res, err := svc.Search(context.Background(), "teacher", "school-1", "Mario", "")
	assert.NoError(t, err)
	assert.Equal(t, 1, res.Total)
	assert.Equal(t, "Rossi Mario", res.Results[0].Title)
	mockRepo.AssertExpectations(t)
}
