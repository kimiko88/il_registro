package textbooks

import (
	"context"
	"errors"
	"testing"
)

type mockRepo struct {
	textbooks   map[string]*Textbook
	assignments map[string]*ClassTextbook
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		textbooks:   make(map[string]*Textbook),
		assignments: make(map[string]*ClassTextbook),
	}
}

func (m *mockRepo) Create(ctx context.Context, t *Textbook) error {
	t.ID = "tb-1"
	m.textbooks[t.ID] = t
	return nil
}

func (m *mockRepo) Update(ctx context.Context, t *Textbook) error {
	if _, ok := m.textbooks[t.ID]; !ok {
		return errors.New("not found")
	}
	m.textbooks[t.ID] = t
	return nil
}

func (m *mockRepo) List(ctx context.Context, schoolID string) ([]Textbook, error) {
	var result []Textbook
	for _, t := range m.textbooks {
		if t.SchoolID == schoolID {
			result = append(result, *t)
		}
	}
	return result, nil
}

func (m *mockRepo) Delete(ctx context.Context, id string) error {
	delete(m.textbooks, id)
	return nil
}

func (m *mockRepo) AssignToClass(ctx context.Context, classID, subjectID, textbookID string, isOptional bool) error {
	assignID := "assign-1"
	m.assignments[assignID] = &ClassTextbook{
		ID:         assignID,
		ClassID:    classID,
		SubjectID:  subjectID,
		TextbookID: textbookID,
		IsOptional: isOptional,
	}
	return nil
}

func (m *mockRepo) RemoveFromClass(ctx context.Context, assignmentID string) error {
	delete(m.assignments, assignmentID)
	return nil
}

func (m *mockRepo) ListByClass(ctx context.Context, classID string) ([]ClassTextbook, error) {
	var result []ClassTextbook
	for _, a := range m.assignments {
		if a.ClassID == classID {
			result = append(result, *a)
		}
	}
	return result, nil
}

func TestTextbooksService_CreateValidations(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	ctx := context.Background()

	// 1. Empty title validation
	err := svc.CreateTextbook(ctx, "school-1", CreateTextbookRequest{
		Title: "   ",
		Price: 25.50,
	})
	if err == nil {
		t.Error("expected error for empty title")
	}

	// 2. Negative price validation
	err = svc.CreateTextbook(ctx, "school-1", CreateTextbookRequest{
		Title: "Matematica Blu",
		Price: -10.0,
	})
	if err == nil {
		t.Error("expected error for negative price")
	}

	// 3. Valid creation
	err = svc.CreateTextbook(ctx, "school-1", CreateTextbookRequest{
		Title:     "Matematica.blu 2.0",
		Author:    "Massimo Bergamini",
		Subject:   "Matematica",
		ISBN:      "9788808836243",
		Publisher: "Zanichelli",
		Price:     28.90,
	})
	if err != nil {
		t.Fatalf("unexpected error creating textbook: %v", err)
	}

	list, err := svc.ListTextbooks(ctx, "school-1")
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 textbook listed, got %d (err: %v)", len(list), err)
	}
	if list[0].Title != "Matematica.blu 2.0" || list[0].ISBN != "9788808836243" {
		t.Errorf("unexpected textbook data: %+v", list[0])
	}
}

func TestTextbooksService_AssignAndListByClass(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	ctx := context.Background()

	err := svc.AssignToClass(ctx, "class-101", AssignTextbookRequest{
		SubjectID:  "subj-math",
		TextbookID: "tb-1",
		IsOptional: false,
	})
	if err != nil {
		t.Fatalf("unexpected assignment error: %v", err)
	}

	assignedList, err := svc.ListByClass(ctx, "class-101")
	if err != nil || len(assignedList) != 1 {
		t.Fatalf("expected 1 assigned textbook, got %d", len(assignedList))
	}
	if assignedList[0].TextbookID != "tb-1" || assignedList[0].IsOptional {
		t.Errorf("unexpected assignment data: %+v", assignedList[0])
	}
}
