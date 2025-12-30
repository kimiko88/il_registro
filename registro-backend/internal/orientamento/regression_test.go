package orientamento

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ValidationMockRepo struct{ mock.Mock }

func (m *ValidationMockRepo) CreateEvent(ctx context.Context, e *Event) error { return nil }
func (m *ValidationMockRepo) GetEvents(ctx context.Context, s string) ([]Event, error) {
	return nil, nil
}
func (m *ValidationMockRepo) RegisterStudent(ctx context.Context, p *Participation) error { return nil }
func (m *ValidationMockRepo) GetParticipations(ctx context.Context, s string) ([]Participation, error) {
	return nil, nil
}
func (m *ValidationMockRepo) MarkAttendance(ctx context.Context, e, s string, a bool) error {
	return nil
}

func TestRegression_CreateEvent_Validation(t *testing.T) {
	repo := new(ValidationMockRepo)
	svc := NewService(repo)

	t.Run("End Before Start", func(t *testing.T) {
		req := CreateEventRequest{
			Title: "Event", Category: "Uni",
			Date:    "2025-01-01T10:00:00Z",
			EndDate: "2024-12-31T10:00:00Z",
		}
		err := svc.CreateEvent(context.Background(), "t1", req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "end_date cannot be before")
	})
}
