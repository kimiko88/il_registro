package search

import (
	"context"
	"fmt"
	"strings"
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
	// Validazione input: previene ReDoS e query eccessive sul DB.
	query = strings.TrimSpace(query)
	if query == "" {
		return &SearchResponse{Query: "", Total: 0, Results: []SearchResultItem{}}, nil
	}
	const maxQueryLen = 200
	if len(query) > maxQueryLen {
		return nil, fmt.Errorf("query troppo lunga: massimo %d caratteri", maxQueryLen)
	}

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
