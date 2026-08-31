package orientamento

import (
	"context"
	"testing"
)

type mockRepo struct {
	capolavori []Capolavoro
}

func (m *mockRepo) CreateEvent(ctx context.Context, e *Event) error { return nil }
func (m *mockRepo) GetEvents(ctx context.Context, schoolID string) ([]Event, error) {
	return []Event{}, nil
}
func (m *mockRepo) RegisterStudent(ctx context.Context, p *Participation) error { return nil }
func (m *mockRepo) GetParticipations(ctx context.Context, studentID string) ([]Participation, error) {
	return []Participation{
		{
			ID:       "part-1",
			Attended: true,
			Event:    &Event{Hours: 15.0},
		},
	}, nil
}
func (m *mockRepo) MarkAttendance(ctx context.Context, eventID, studentID string, attended bool) error {
	return nil
}
func (m *mockRepo) SavePreference(ctx context.Context, p *StudentPreference) error { return nil }
func (m *mockRepo) GetPreference(ctx context.Context, studentID string) (*StudentPreference, error) {
	return &StudentPreference{StudentID: studentID}, nil
}
func (m *mockRepo) SaveCapolavoro(ctx context.Context, c *Capolavoro) error {
	c.ID = "cap-123"
	m.capolavori = append(m.capolavori, *c)
	return nil
}
func (m *mockRepo) GetCapolavori(ctx context.Context, studentID string) ([]Capolavoro, error) {
	return m.capolavori, nil
}
func (m *mockRepo) GetCurriculumStudente(ctx context.Context, studentID string) (*CurriculumStudenteSummary, error) {
	return &CurriculumStudenteSummary{
		StudentID:         studentID,
		StudentName:       "Mario Rossi",
		SchoolName:        "Liceo Scientifico Statale",
		PCTOHours:         120.0,
		OrientamentoHours: 30.0,
		Capolavori:        m.capolavori,
	}, nil
}

func TestEPortfolioAndCurriculum(t *testing.T) {
	repo := &mockRepo{}
	svc := NewService(repo)
	ctx := context.Background()

	// 1. Test SaveCapolavoro with validation
	errNoTitle := svc.SaveCapolavoro(ctx, "student-1", Capolavoro{})
	if errNoTitle == nil {
		t.Errorf("expected error when saving capolavoro with empty title")
	}

	cap := Capolavoro{
		Title:           "Sviluppo Registro Elettronico Open Source",
		Description:     "Progetto di sviluppo software per la scuola pubblica",
		SchoolYear:      "2025/2026",
		ReflectiveNotes: "Ho potenziato le mie competenze di problem solving e programmazione concorrente",
	}
	if err := svc.SaveCapolavoro(ctx, "student-1", cap); err != nil {
		t.Fatalf("failed to save capolavoro: %v", err)
	}

	// 2. Test GetCapolavori
	list, err := svc.GetCapolavori(ctx, "student-1")
	if err != nil || len(list) != 1 {
		t.Fatalf("expected 1 capolavoro, got %d (err: %v)", len(list), err)
	}
	if list[0].Title != cap.Title {
		t.Errorf("expected title '%s', got '%s'", cap.Title, list[0].Title)
	}

	// 3. Test GetCurriculumStudente
	curr, err := svc.GetCurriculumStudente(ctx, "student-1")
	if err != nil {
		t.Fatalf("failed to get curriculum: %v", err)
	}
	if curr.PCTOHours != 120.0 {
		t.Errorf("expected 120 PCTO hours, got %f", curr.PCTOHours)
	}
	if curr.OrientamentoHours < 30.0 {
		t.Errorf("expected at least 30 orientamento hours, got %f", curr.OrientamentoHours)
	}
	if len(curr.Capolavori) != 1 {
		t.Errorf("expected 1 capolavoro in curriculum, got %d", len(curr.Capolavori))
	}
}
