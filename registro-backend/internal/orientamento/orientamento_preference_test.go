package orientamento

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPreferenceRepo struct {
	pref *StudentPreference
}

func (m *mockPreferenceRepo) CreateEvent(ctx context.Context, e *Event) error { return nil }
func (m *mockPreferenceRepo) GetEvents(ctx context.Context, schoolID string) ([]Event, error) {
	return nil, nil
}
func (m *mockPreferenceRepo) RegisterStudent(ctx context.Context, p *Participation) error {
	return nil
}
func (m *mockPreferenceRepo) GetParticipations(ctx context.Context, studentID string) ([]Participation, error) {
	return nil, nil
}
func (m *mockPreferenceRepo) MarkAttendance(ctx context.Context, eventID, studentID string, attended bool) error {
	return nil
}
func (m *mockPreferenceRepo) SavePreference(ctx context.Context, p *StudentPreference) error {
	m.pref = p
	return nil
}
func (m *mockPreferenceRepo) GetPreference(ctx context.Context, studentID string) (*StudentPreference, error) {
	if m.pref == nil {
		return &StudentPreference{StudentID: studentID}, nil
	}
	return m.pref, nil
}

func TestOrientamento_GetPreferenceEmpty(t *testing.T) {
	repo := &mockPreferenceRepo{}
	svc := NewService(repo)

	pref, err := svc.GetPreference(context.Background(), "student-1")
	assert.NoError(t, err)
	assert.NotNil(t, pref)
	assert.Equal(t, "student-1", pref.StudentID)
	assert.Empty(t, pref.PreferredTrack, "should not return hardcoded University")
	assert.Empty(t, pref.TargetField, "should not return hardcoded Ingegneria Informatica")
}
