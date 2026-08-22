package search

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct{ mock.Mock }

func (m *MockRepo) GlobalSearch(ctx context.Context, actorRole, schoolID, query, filterType string) ([]SearchResultItem, error) {
	args := m.Called(actorRole, schoolID, query, filterType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]SearchResultItem), args.Error(1)
}

func TestSearch_QueryTooLong_Rejected(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)
	longQuery := strings.Repeat("a", 201)
	_, err := svc.Search(context.Background(), "admin", "school-1", longQuery, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "troppo lunga")
}

func TestSearch_EmptyQuery_ReturnsEmpty(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)
	res, err := svc.Search(context.Background(), "admin", "school-1", "   ", "")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 0, res.Total)
	assert.Equal(t, "", res.Query)
}

func TestSearch_ValidQuery_PassedToRepo(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)
	repo.On("GlobalSearch", "admin", "school-1", "mario", "").Return([]SearchResultItem{{ID: "u-1", Title: "Mario Rossi"}}, nil)
	res, err := svc.Search(context.Background(), "admin", "school-1", "mario", "")
	assert.NoError(t, err)
	assert.Equal(t, 1, res.Total)
	repo.AssertExpectations(t)
}

func TestSearch_Exactly200Chars_Accepted(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)
	okQuery := strings.Repeat("b", 200)
	repo.On("GlobalSearch", "admin", "school-1", okQuery, "").Return([]SearchResultItem{}, nil)
	_, err := svc.Search(context.Background(), "admin", "school-1", okQuery, "")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
