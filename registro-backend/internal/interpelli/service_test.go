package interpelli

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateScore(t *testing.T) {
	t.Run("Score without habilitation, grade 110 with lode", func(t *testing.T) {
		deg, srv, certs, total := CalculateScore(110, true, false, 6, "B2", 2)
		// 12 + 0.5 * 50 = 37; +4 (lode) = 41
		assert.Equal(t, 41.0, deg)
		// 6 * 2 = 12
		assert.Equal(t, 12.0, srv)
		// B2 (3) + 2*0.5 (1) = 4
		assert.Equal(t, 4.0, certs)
		assert.Equal(t, 57.0, total)
	})

	t.Run("Score with habilitation and C2", func(t *testing.T) {
		deg, _, certs, _ := CalculateScore(100, false, true, 0, "C2", 0)
		// 12 + 0.5 * 40 = 32; +24 (hab) = 56
		assert.Equal(t, 56.0, deg)
		assert.Equal(t, 6.0, certs)
	})

	t.Run("Cap on service and digital certs", func(t *testing.T) {
		_, srv, certs, _ := CalculateScore(60, false, false, 50, "", 10)
		assert.Equal(t, 36.0, srv)  // capped at 36
		assert.Equal(t, 2.0, certs) // capped at 2
	})
}

func TestInterpelliLifecycle(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	// 1. Create notice
	notice, err := svc.CreateNotice(ctx, "school-1", CreateNoticeRequest{
		Title:         "Ricerca Docente A026 Matematica",
		ConcorsoClass: "A026",
		PostType:      "comune",
		WeeklyHours:   18,
		StartDate:     "2026-10-15",
		EndDate:       "2027-06-30",
		Deadline:      time.Now().Add(48 * time.Hour).Format(time.RFC3339),
		Description:   "Cattedra intera fino al termine delle attività didattiche.",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, notice.ID)
	assert.Equal(t, "aperto", notice.Status)

	// 2. Submit Candidatura 1 (lower score)
	cand1, err := svc.SubmitCandidatura(ctx, notice.ID, SubmitCandidaturaRequest{
		CandidateName:    "Mario",
		CandidateSurname: "Rossi",
		FiscalCode:       "RSSMRA80A01H501U",
		Email:            "mario.rossi@example.com",
		Phone:            "3331234567",
		GraduationGrade:  80,
		HasHabilitation:  false,
		MonthsOfService:  2,
		DPR445Declared:   true,
	})
	require.NoError(t, err)

	// 3. Submit Candidatura 2 (higher score with habilitation)
	cand2, err := svc.SubmitCandidatura(ctx, notice.ID, SubmitCandidaturaRequest{
		CandidateName:    "Laura",
		CandidateSurname: "Bianchi",
		FiscalCode:       "BNCLRA85B41H501Z",
		Email:            "laura.bianchi@example.com",
		Phone:            "3339876543",
		GraduationGrade:  110,
		GraduationLode:   true,
		HasHabilitation:  true,
		MonthsOfService:  12,
		CertLanguage:     "C1",
		CertDigitalCount: 2,
		DPR445Declared:   true,
	})
	require.NoError(t, err)

	// 4. Test DPR 445 mandatory check
	_, err = svc.SubmitCandidatura(ctx, notice.ID, SubmitCandidaturaRequest{
		CandidateName:   "Pietro",
		GraduationGrade: 100,
		DPR445Declared:  false, // missing declaration
	})
	assert.ErrorIs(t, err, ErrDPR445Mandatory)

	// 5. Get Graduatoria: Laura Bianchi must be #1, Mario Rossi #2
	grad, err := svc.GetGraduatoria(ctx, notice.ID)
	require.NoError(t, err)
	require.Len(t, grad, 2)
	assert.Equal(t, cand2.ID, grad[0].ID)
	assert.Equal(t, cand1.ID, grad[1].ID)
	assert.Greater(t, grad[0].TotalScore, grad[1].TotalScore)

	// 6. Convocation with 24 hours deadline
	convocata, err := svc.ConvocaCandidato(ctx, cand2.ID, 24)
	require.NoError(t, err)
	assert.Equal(t, "convocata", convocata.Status)
	assert.NotNil(t, convocata.ConvocationDeadline)

	// 7. Accept convocation
	accettata, err := svc.RispondiConvocazione(ctx, cand2.ID, "accettata", "Accetto con presa di servizio il 15 ottobre")
	require.NoError(t, err)
	assert.Equal(t, "accettata", accettata.Status)
	assert.NotEmpty(t, accettata.ResponseNotes)
}
