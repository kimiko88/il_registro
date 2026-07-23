package auditlog

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// mockRepo implements Repository for unit tests.
type mockRepo struct {
	events []AuditEvent
	err    error
}

func (m *mockRepo) Insert(ctx context.Context, event *AuditEvent) error {
	if m.err != nil {
		return m.err
	}
	m.events = append(m.events, *event)
	return nil
}

func (m *mockRepo) InsertAsync(event AuditEvent) {
	m.events = append(m.events, event)
}

func (m *mockRepo) List(ctx context.Context, p FilterParams) ([]AuditEvent, int, error) {
	if m.err != nil {
		return nil, 0, m.err
	}

	var filtered []AuditEvent
	for _, e := range m.events {
		if p.ActorID != "" && e.ActorID != p.ActorID {
			continue
		}
		if p.Action != "" && e.Action != p.Action {
			continue
		}
		if p.EntityType != "" && e.EntityType != p.EntityType {
			continue
		}
		filtered = append(filtered, e)
	}

	total := len(filtered)
	start := (p.Page - 1) * p.Limit
	if start >= total {
		return []AuditEvent{}, total, nil
	}
	end := start + p.Limit
	if end > total {
		end = total
	}
	return filtered[start:end], total, nil
}

// Verify mockRepo satisfies Repository interface at compile time.
var _ Repository = (*mockRepo)(nil)

func TestService_Log_AddsEvent(t *testing.T) {
	repo := &mockRepo{}
	svc := NewService(repo)

	svc.Log(context.Background(), AuditEvent{
		ID:       "ev-1",
		ActorID:  "user-1",
		Action:   "login",
		SchoolID: "school-1",
	})

	if len(repo.events) != 1 {
		t.Errorf("expected 1 event, got %d", len(repo.events))
	}
	if repo.events[0].Action != "login" {
		t.Errorf("expected action 'login', got '%s'", repo.events[0].Action)
	}
}

func TestService_List_DefaultPagination(t *testing.T) {
	repo := &mockRepo{}
	for i := 0; i < 5; i++ {
		repo.events = append(repo.events, AuditEvent{ID: "ev", ActorID: "u1", Action: "read"})
	}
	svc := NewService(repo)

	// Unset page/limit — should default to page=1, limit=20
	result, err := svc.List(context.Background(), FilterParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Page != 1 {
		t.Errorf("expected Page=1, got %d", result.Page)
	}
	if result.Limit != 20 {
		t.Errorf("expected Limit=20, got %d", result.Limit)
	}
	if result.Total != 5 {
		t.Errorf("expected Total=5, got %d", result.Total)
	}
	if result.TotalPages != 1 {
		t.Errorf("expected TotalPages=1, got %d", result.TotalPages)
	}
}

func TestService_List_FilterByActorID(t *testing.T) {
	repo := &mockRepo{
		events: []AuditEvent{
			{ID: "e1", ActorID: "user-A", Action: "create"},
			{ID: "e2", ActorID: "user-B", Action: "delete"},
			{ID: "e3", ActorID: "user-A", Action: "update"},
		},
	}
	svc := NewService(repo)

	result, err := svc.List(context.Background(), FilterParams{ActorID: "user-A", Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 2 {
		t.Errorf("expected 2 events for user-A, got %d", result.Total)
	}
}

func TestService_List_FilterByAction(t *testing.T) {
	repo := &mockRepo{
		events: []AuditEvent{
			{ID: "e1", ActorID: "user-A", Action: "login"},
			{ID: "e2", ActorID: "user-B", Action: "logout"},
			{ID: "e3", ActorID: "user-C", Action: "login"},
		},
	}
	svc := NewService(repo)

	result, err := svc.List(context.Background(), FilterParams{Action: "login", Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 2 {
		t.Errorf("expected 2 login events, got %d", result.Total)
	}
}

func TestService_List_TotalPagesRounding(t *testing.T) {
	repo := &mockRepo{}
	for i := 0; i < 25; i++ {
		repo.events = append(repo.events, AuditEvent{ID: "ev", ActorID: "u"})
	}
	svc := NewService(repo)

	result, err := svc.List(context.Background(), FilterParams{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 25 items / 10 per page = ceil(2.5) = 3 pages
	if result.TotalPages != 3 {
		t.Errorf("expected 3 total pages, got %d", result.TotalPages)
	}
}

func TestService_List_EmptyResultIsNotNil(t *testing.T) {
	repo := &mockRepo{}
	svc := NewService(repo)

	result, err := svc.List(context.Background(), FilterParams{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Data == nil {
		t.Error("expected non-nil Data slice for empty results")
	}
}

func TestService_List_PropagatesRepositoryError(t *testing.T) {
	repo := &mockRepo{err: errors.New("db failure")}
	svc := NewService(repo)

	_, err := svc.List(context.Background(), FilterParams{Page: 1, Limit: 10})
	if err == nil {
		t.Error("expected error propagated from repository, got nil")
	}
}

func TestService_ExportCSV_ContainsHeader(t *testing.T) {
	repo := &mockRepo{
		events: []AuditEvent{
			{ID: "e1", ActorID: "user-A", ActorName: "Mario Rossi", ActorRole: "teacher", Action: "login"},
		},
	}
	svc := NewService(repo)

	data, err := svc.ExportCSV(context.Background(), FilterParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	csv := string(data)
	if !strings.Contains(csv, "ID") || !strings.Contains(csv, "Azione") {
		t.Errorf("CSV header missing expected columns, got: %s", csv[:min(len(csv), 100)])
	}
}

func TestService_ExportCSV_PropagatesError(t *testing.T) {
	repo := &mockRepo{err: errors.New("db failure")}
	svc := NewService(repo)

	_, err := svc.ExportCSV(context.Background(), FilterParams{})
	if err == nil {
		t.Error("expected error from ExportCSV when repo fails")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
