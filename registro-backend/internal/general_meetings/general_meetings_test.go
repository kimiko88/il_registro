package general_meetings

import (
	"context"
	"errors"
	"testing"
)

type mockRepo struct {
	meetings      map[string]*GeneralMeeting
	registrations map[string][]string
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		meetings:      make(map[string]*GeneralMeeting),
		registrations: make(map[string][]string),
	}
}

func (m *mockRepo) Create(ctx context.Context, gm *GeneralMeeting) error {
	gm.ID = "meeting-1"
	m.meetings[gm.ID] = gm
	return nil
}

func (m *mockRepo) List(ctx context.Context, schoolID, userID, role string) ([]*GeneralMeeting, error) {
	var list []*GeneralMeeting
	for _, gm := range m.meetings {
		if gm.SchoolID == schoolID {
			list = append(list, gm)
		}
	}
	return list, nil
}

func (m *mockRepo) GetByID(ctx context.Context, id, userID string) (*GeneralMeeting, error) {
	gm, ok := m.meetings[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return gm, nil
}

func (m *mockRepo) Delete(ctx context.Context, id string) error {
	delete(m.meetings, id)
	return nil
}

func (m *mockRepo) RegisterUser(ctx context.Context, meetingID, userID string) error {
	m.registrations[meetingID] = append(m.registrations[meetingID], userID)
	if gm, ok := m.meetings[meetingID]; ok {
		gm.RegistrationsCount++
	}
	return nil
}

func (m *mockRepo) UnregisterUser(ctx context.Context, meetingID, userID string) error {
	if regs, ok := m.registrations[meetingID]; ok {
		var updated []string
		for _, uid := range regs {
			if uid != userID {
				updated = append(updated, uid)
			}
		}
		m.registrations[meetingID] = updated
		if gm, ok := m.meetings[meetingID]; ok && gm.RegistrationsCount > 0 {
			gm.RegistrationsCount--
		}
	}
	return nil
}

func (m *mockRepo) ListRegistrations(ctx context.Context, meetingID string) ([]*GeneralMeetingRegistration, error) {
	var list []*GeneralMeetingRegistration
	for _, uid := range m.registrations[meetingID] {
		list = append(list, &GeneralMeetingRegistration{
			MeetingID: meetingID,
			UserID:    uid,
		})
	}
	return list, nil
}

func TestGeneralMeetings_CreateAndRegister(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	ctx := context.Background()

	maxParts := 2
	gm, err := svc.CreateMeeting(ctx, "user-creator", "school-1", CreateGeneralMeetingRequest{
		Title:           "Collegio Docenti Inizio Anno",
		Description:     "Pianificazione attività didattiche",
		Location:        "Aula Magna",
		MeetingDate:     "2026-09-01T09:00:00Z",
		MaxParticipants: &maxParts,
		IsMandatory:     true,
	})
	if err != nil {
		t.Fatalf("unexpected error creating meeting: %v", err)
	}

	if gm.Title != "Collegio Docenti Inizio Anno" || gm.CreatedBy != "user-creator" {
		t.Errorf("unexpected meeting data: %+v", gm)
	}

	// Register 2 users (up to max capacity)
	err = svc.RegisterUser(ctx, gm.ID, "teacher-1")
	if err != nil {
		t.Fatalf("unexpected error registering teacher-1: %v", err)
	}

	err = svc.RegisterUser(ctx, gm.ID, "teacher-2")
	if err != nil {
		t.Fatalf("unexpected error registering teacher-2: %v", err)
	}

	// 3rd user register attempt should fail (fully booked)
	err = svc.RegisterUser(ctx, gm.ID, "teacher-3")
	if err == nil {
		t.Error("expected error when registering for fully booked meeting, but succeeded")
	}
}

func TestGeneralMeetings_DeleteAuthorization(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	ctx := context.Background()

	gm, err := svc.CreateMeeting(ctx, "creator-id", "school-1", CreateGeneralMeetingRequest{
		Title:       "Assemblea d'Istituto",
		MeetingDate: "2026-10-15",
	})
	if err != nil {
		t.Fatalf("failed to create meeting: %v", err)
	}

	// Unauthorized user (not creator, not admin) trying to delete
	err = svc.DeleteMeeting(ctx, gm.ID, "other-user", "teacher", "school-1")
	if err == nil {
		t.Error("expected error when non-creator non-admin user deletes meeting")
	}

	// Admin of different school trying to delete
	err = svc.DeleteMeeting(ctx, gm.ID, "admin-user", "admin", "school-2")
	if err == nil {
		t.Error("expected error when admin of different school deletes meeting")
	}

	// Admin of same school deleting meeting
	err = svc.DeleteMeeting(ctx, gm.ID, "admin-user", "admin", "school-1")
	if err != nil {
		t.Errorf("unexpected error when admin deletes meeting: %v", err)
	}
}

func TestGeneralMeetings_ListRegistrationsMultiTenant(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	ctx := context.Background()

	gm, err := svc.CreateMeeting(ctx, "creator-id", "school-1", CreateGeneralMeetingRequest{
		Title:       "Collegio Docenti",
		MeetingDate: "2026-10-15",
	})
	if err != nil {
		t.Fatalf("failed to create meeting: %v", err)
	}

	// Different school admin listing registrations
	_, err = svc.ListRegistrations(ctx, gm.ID, "admin-2", "admin", "school-2")
	if err == nil {
		t.Error("expected error for admin of different school listing registrations")
	}

	// Same school admin listing registrations
	regs, err := svc.ListRegistrations(ctx, gm.ID, "admin-1", "admin", "school-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(regs) != 0 {
		t.Errorf("expected 0 registrations, got %d", len(regs))
	}
}
