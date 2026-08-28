package accessibility_test

import (
	"context"
	"testing"

	"registro-backend/internal/accessibility"
)

type mockRepository struct {
	feedbacks []accessibility.AccessibilityFeedback
}

func (m *mockRepository) Create(ctx context.Context, f *accessibility.AccessibilityFeedback) error {
	f.ID = "test-uuid-1234"
	m.feedbacks = append(m.feedbacks, *f)
	return nil
}

func (m *mockRepository) List(ctx context.Context, schoolID, status string, limit, offset int) ([]accessibility.AccessibilityFeedback, int, error) {
	return m.feedbacks, len(m.feedbacks), nil
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*accessibility.AccessibilityFeedback, error) {
	for _, f := range m.feedbacks {
		if f.ID == id {
			return &f, nil
		}
	}
	return nil, nil
}

func (m *mockRepository) UpdateStatus(ctx context.Context, id, status, responseNotes string) error {
	return nil
}

func (m *mockRepository) GetUserPreferences(ctx context.Context, userID string) (string, error) {
	return `{"dsaFont":true}`, nil
}

func (m *mockRepository) UpsertUserPreferences(ctx context.Context, userID string, settingsJSON string) error {
	return nil
}

func TestSubmitFeedbackValidation(t *testing.T) {
	repo := &mockRepository{}
	svc := accessibility.NewService(repo)
	ctx := context.Background()

	// 1. Missing name
	req1 := &accessibility.CreateFeedbackRequest{
		Name:        "",
		Email:       "test@scuola.it",
		BarrierType: "contrast",
		Description: "Contrasto non sufficiente",
	}
	if _, err := svc.SubmitFeedback(ctx, req1, nil, nil, nil, nil); err == nil {
		t.Errorf("Expected error for missing name, got nil")
	}

	// 2. Invalid email
	req2 := &accessibility.CreateFeedbackRequest{
		Name:        "Mario Rossi",
		Email:       "email-non-valida",
		BarrierType: "contrast",
		Description: "Contrasto non sufficiente",
	}
	if _, err := svc.SubmitFeedback(ctx, req2, nil, nil, nil, nil); err == nil {
		t.Errorf("Expected error for invalid email, got nil")
	}

	// 3. Short description
	req3 := &accessibility.CreateFeedbackRequest{
		Name:        "Mario Rossi",
		Email:       "mario@scuola.it",
		BarrierType: "contrast",
		Description: "abc",
	}
	if _, err := svc.SubmitFeedback(ctx, req3, nil, nil, nil, nil); err == nil {
		t.Errorf("Expected error for short description, got nil")
	}

	// 4. Valid request
	req4 := &accessibility.CreateFeedbackRequest{
		Name:        "Mario Rossi",
		Email:       "mario@scuola.it",
		BarrierType: "contrast",
		Description: "Il testo delle circolari ha contrasto insufficiente in modalità scura.",
	}
	resp, err := svc.SubmitFeedback(ctx, req4, nil, nil, nil, nil)
	if err != nil {
		t.Fatalf("Unexpected error for valid submission: %v", err)
	}
	if resp.ProtocolNumber == "" {
		t.Errorf("Expected valid protocol number, got empty")
	}
}
