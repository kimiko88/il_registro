package psychology

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPsychologyDeskLifecycle(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	studentID := "std-minor-1"
	parent1 := "parent-1"
	parent2 := "parent-2"
	psychologistID := "psycho-expert-1"
	principalID := "ds-user"

	// 1. Student attempts booking without parental consents -> blocked!
	_, err := svc.BookSession(ctx, "school-1", studentID, "2025/2026", true, BookSessionRequest{
		SlotTime:       time.Now().Add(48 * time.Hour).Format("2006-01-02 15:04"),
		PsychologistID: psychologistID,
	})
	assert.ErrorIs(t, err, ErrBothParentsConsentRequired)

	// 2. Only Parent 1 signs -> still blocked (pending parent 2)!
	consent, err := svc.SignParentConsent(ctx, "school-1", studentID, "2025/2026", parent1)
	require.NoError(t, err)
	assert.Equal(t, ConsentPending, consent.Status)

	_, err = svc.BookSession(ctx, "school-1", studentID, "2025/2026", true, BookSessionRequest{
		SlotTime:       time.Now().Add(48 * time.Hour).Format("2006-01-02 15:04"),
		PsychologistID: psychologistID,
	})
	assert.ErrorIs(t, err, ErrBothParentsConsentRequired)

	// 3. Parent 2 signs -> Consent Approved!
	consent, err = svc.SignParentConsent(ctx, "school-1", studentID, "2025/2026", parent2)
	require.NoError(t, err)
	assert.Equal(t, ConsentApproved, consent.Status)

	// 4. Minor student now books session successfully
	session, err := svc.BookSession(ctx, "school-1", studentID, "2025/2026", true, BookSessionRequest{
		SlotTime:       time.Now().Add(48 * time.Hour).Format("2006-01-02 15:04"),
		PsychologistID: psychologistID,
	})
	require.NoError(t, err)
	assert.Equal(t, "prenotato", session.Status)
	assert.Contains(t, session.AnonymousAlias, "ALUNNO-CIC-")

	// 5. Psychologist updates clinical notes
	_, err = svc.UpdateClinicalNotes(ctx, session.ID, psychologistID, "Colloquio orientativo svolto in clima di ascolto empatico")
	require.NoError(t, err)

	// 6. Test Professional Secrecy: Principal tries to read clinical notes -> strictly forbidden!
	_, err = svc.GetSessionSecure(ctx, session.ID, principalID, "principal")
	assert.ErrorIs(t, err, ErrProfessionalSecrecyViolation)

	// 7. Student reads their session -> notes are redacted
	stdView, err := svc.GetSessionSecure(ctx, session.ID, studentID, "student")
	require.NoError(t, err)
	assert.Empty(t, stdView.EncryptedClinicalNotes)

	// 8. Psychologist reads session -> notes are accessible
	psyView, err := svc.GetSessionSecure(ctx, session.ID, psychologistID, "psychologist")
	require.NoError(t, err)
	assert.NotEmpty(t, psyView.EncryptedClinicalNotes)
}
