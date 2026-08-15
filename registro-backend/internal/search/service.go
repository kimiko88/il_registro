package search

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	if repo == nil {
		panic("search.NewService: repo must not be nil")
	}
	return &Service{repo: repo}
}

func (s *Service) Search(ctx context.Context, actorRole, schoolID, query, filterType string) (*SearchResponse, error) {
	results, err := s.repo.GlobalSearch(ctx, actorRole, schoolID, query, filterType)
	if err != nil {
		return nil, err
	}
	if results == nil {
		results = []SearchResultItem{}
	}

	return &SearchResponse{
		Query:   query,
		Total:   len(results),
		Results: results,
	}, nil
}
