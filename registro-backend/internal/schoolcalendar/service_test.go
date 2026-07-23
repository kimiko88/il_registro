package schoolcalendar

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

// mockCalendarRepo implements Repository for unit tests.
type mockCalendarRepo struct {
	year         *SchoolYearSettings
	nonTeaching  []NonTeachingDay
	periods      []AcademicPeriod
	teachingDays int
	err          error
}

func (m *mockCalendarRepo) UpsertYear(s *SchoolYearSettings) error {
	if m.err != nil {
		return m.err
	}
	m.year = s
	return nil
}

func (m *mockCalendarRepo) GetYear(schoolID string) (*SchoolYearSettings, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.year == nil {
		return nil, errors.New("sql: no rows in result set")
	}
	return m.year, nil
}

func (m *mockCalendarRepo) AddNonTeachingDay(d *NonTeachingDay) error {
	if m.err != nil {
		return m.err
	}
	m.nonTeaching = append(m.nonTeaching, *d)
	return nil
}

func (m *mockCalendarRepo) DeleteNonTeachingDay(schoolID, id string) error {
	return m.err
}

func (m *mockCalendarRepo) ListNonTeachingDays(schoolID string) ([]NonTeachingDay, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.nonTeaching, nil
}

func (m *mockCalendarRepo) CountTeachingDays(schoolID string, from, to time.Time) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.teachingDays, nil
}

func (m *mockCalendarRepo) CreateAcademicPeriod(p *AcademicPeriod) error {
	if m.err != nil {
		return m.err
	}
	m.periods = append(m.periods, *p)
	return nil
}

func (m *mockCalendarRepo) ListAcademicPeriods(schoolID string) ([]AcademicPeriod, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.periods, nil
}

var _ Repository = (*mockCalendarRepo)(nil)

// --- SetSchoolYear tests ---

func TestSetSchoolYear_ForbiddenForNonSecretary(t *testing.T) {
	svc := NewService(&mockCalendarRepo{})
	_, err := svc.SetSchoolYear(context.Background(), "actor", "teacher", "school-1", SetYearRequest{
		StartDate: "2025-09-15",
		EndDate:   "2026-06-10",
		YearLabel: "2025/2026",
	})
	if err == nil || !strings.Contains(err.Error(), "forbidden") {
		t.Errorf("expected forbidden error for non-secretary, got: %v", err)
	}
}

func TestSetSchoolYear_EndBeforeStart(t *testing.T) {
	svc := NewService(&mockCalendarRepo{})
	_, err := svc.SetSchoolYear(context.Background(), "actor", "secretary", "school-1", SetYearRequest{
		StartDate: "2026-06-10",
		EndDate:   "2025-09-15",
		YearLabel: "2025/2026",
	})
	if err == nil || !strings.Contains(err.Error(), "end_date") {
		t.Errorf("expected date order error, got: %v", err)
	}
}

func TestSetSchoolYear_ValidRequest(t *testing.T) {
	repo := &mockCalendarRepo{teachingDays: 200}
	svc := NewService(repo)
	resp, err := svc.SetSchoolYear(context.Background(), "actor", "secretary", "school-1", SetYearRequest{
		StartDate: "2025-09-15",
		EndDate:   "2026-06-10",
		YearLabel: "2025/2026",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.YearLabel != "2025/2026" {
		t.Errorf("expected YearLabel '2025/2026', got '%s'", resp.YearLabel)
	}
	if resp.TeachingDays != 200 {
		t.Errorf("expected TeachingDays 200, got %d", resp.TeachingDays)
	}
}

func TestSetSchoolYear_InvalidDateFormat(t *testing.T) {
	svc := NewService(&mockCalendarRepo{})
	_, err := svc.SetSchoolYear(context.Background(), "actor", "secretary", "school-1", SetYearRequest{
		StartDate: "15/09/2025", // wrong format
		EndDate:   "2026-06-10",
		YearLabel: "2025/2026",
	})
	if err == nil {
		t.Error("expected error for invalid date format, got nil")
	}
}

// --- AddNonTeachingDay tests ---

func TestAddNonTeachingDay_ForbiddenForTeacher(t *testing.T) {
	svc := NewService(&mockCalendarRepo{})
	_, err := svc.AddNonTeachingDay(context.Background(), "actor", "teacher", "school-1", AddNonTeachingDayRequest{
		Date:  "2025-12-25",
		Label: "Natale",
	})
	if err == nil || !strings.Contains(err.Error(), "forbidden") {
		t.Errorf("expected forbidden error, got: %v", err)
	}
}

func TestAddNonTeachingDay_ValidRequest(t *testing.T) {
	repo := &mockCalendarRepo{}
	svc := NewService(repo)
	resp, err := svc.AddNonTeachingDay(context.Background(), "actor", "admin", "school-1", AddNonTeachingDayRequest{
		Date:  "2025-12-25",
		Label: "Natale",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Label != "Natale" {
		t.Errorf("expected Label 'Natale', got '%s'", resp.Label)
	}
	if resp.Date != "2025-12-25" {
		t.Errorf("expected Date '2025-12-25', got '%s'", resp.Date)
	}
}

// --- ListNonTeachingDays tests ---

func TestListNonTeachingDays_EmptyReturnsSlice(t *testing.T) {
	svc := NewService(&mockCalendarRepo{})
	days, err := svc.ListNonTeachingDays(context.Background(), "school-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if days == nil {
		t.Error("expected non-nil slice for empty days list")
	}
}

// --- GetCalendarEvents tests ---

func TestGetCalendarEvents_FilterByYear(t *testing.T) {
	repo := &mockCalendarRepo{
		nonTeaching: []NonTeachingDay{
			{ID: "h1", Date: time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC), Label: "Natale"},
			{ID: "h2", Date: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), Label: "Capodanno"},
		},
	}
	svc := NewService(repo)
	events, err := svc.GetCalendarEvents(context.Background(), "school-1", "2025", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 1 {
		t.Errorf("expected 1 event for year 2025, got %d", len(events))
	}
	if events[0].Title != "Natale" {
		t.Errorf("expected title 'Natale', got '%s'", events[0].Title)
	}
}

func TestGetCalendarEvents_EmptyReturnsSlice(t *testing.T) {
	svc := NewService(&mockCalendarRepo{})
	events, err := svc.GetCalendarEvents(context.Background(), "school-1", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if events == nil {
		t.Error("expected non-nil slice for no events")
	}
}

// --- isSecretary tests ---

func TestIsSecretary_Roles(t *testing.T) {
	secretaryRoles := []string{"secretary", "admin", "superadmin"}
	nonSecretaryRoles := []string{"teacher", "student", "parent", "principal"}

	for _, r := range secretaryRoles {
		if !isSecretary(r) {
			t.Errorf("expected role '%s' to be secretary", r)
		}
	}
	for _, r := range nonSecretaryRoles {
		if isSecretary(r) {
			t.Errorf("expected role '%s' to NOT be secretary", r)
		}
	}
}
