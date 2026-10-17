package maturita

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateYearCredits(t *testing.T) {
	t.Run("3rd year credits by average", func(t *testing.T) {
		// M = 6.0: 7 (low), 8 (high)
		cLow, err := CalculateYearCredits(3, 6.0, false)
		require.NoError(t, err)
		assert.Equal(t, 7, cLow)

		cHigh, err := CalculateYearCredits(3, 6.0, true)
		require.NoError(t, err)
		assert.Equal(t, 8, cHigh)

		// 9 < M <= 10: 11 (low), 12 (high)
		cMax, err := CalculateYearCredits(3, 9.5, true)
		require.NoError(t, err)
		assert.Equal(t, 12, cMax)
	})

	t.Run("4th year credits by average", func(t *testing.T) {
		c, err := CalculateYearCredits(4, 8.2, true)
		require.NoError(t, err)
		assert.Equal(t, 12, c)
	})

	t.Run("5th year credits by average", func(t *testing.T) {
		c, err := CalculateYearCredits(5, 10.0, true)
		require.NoError(t, err)
		assert.Equal(t, 15, c)
	})

	t.Run("Invalid year or average", func(t *testing.T) {
		_, err := CalculateYearCredits(2, 8.0, true)
		assert.ErrorIs(t, err, ErrInvalidGradeYear)

		_, err = CalculateYearCredits(3, 5.5, true)
		assert.ErrorIs(t, err, ErrInvalidAverage)
	})
}

func TestCalculateFinalExamScores(t *testing.T) {
	credits := TrienniumCredits{
		Credit3rd:    12,
		Credit4th:    13,
		Credit5th:    15,
		TotalCredits: 40,
	}

	t.Run("Standard passing exam without bonus", func(t *testing.T) {
		scores, err := CalculateFinalExamScores(credits, 15.0, 16.0, 18.0, 0, false)
		require.NoError(t, err)
		// 40 + (15 + 16 + 18 = 49) = 89
		assert.Equal(t, 89, scores.FinalScore)
		assert.False(t, scores.Lode)
		assert.False(t, scores.BonusEligible) // ExamTotal 49 < 50
	})

	t.Run("Eligible candidate receives bonus points", func(t *testing.T) {
		scores, err := CalculateFinalExamScores(credits, 18.0, 18.0, 18.0, 4, false)
		require.NoError(t, err)
		// 40 + 54 + 4 = 98
		assert.True(t, scores.BonusEligible)
		assert.Equal(t, 4, scores.BonusPoints)
		assert.Equal(t, 98, scores.FinalScore)
	})

	t.Run("Bonus ineligible due to low exam score", func(t *testing.T) {
		_, err := CalculateFinalExamScores(credits, 15.0, 15.0, 15.0, 3, false)
		assert.ErrorIs(t, err, ErrBonusIneligible)
	})

	t.Run("Candidate awarded 100 e Lode", func(t *testing.T) {
		scores, err := CalculateFinalExamScores(credits, 20.0, 20.0, 20.0, 0, true)
		require.NoError(t, err)
		assert.Equal(t, 100, scores.FinalScore)
		assert.True(t, scores.Lode)
	})

	t.Run("Lode rejected when credits < 40 or scores < 60", func(t *testing.T) {
		subCredits := TrienniumCredits{TotalCredits: 38}
		_, err := CalculateFinalExamScores(subCredits, 20.0, 20.0, 20.0, 0, true)
		assert.ErrorIs(t, err, ErrLodeIneligible)
	})
}

func TestMaturitaLifecycleAndXML(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	// 1. Save Commission
	comm := &CommissioneMaturita{
		ClassID:       "class-5A",
		SchoolYear:    "2025/2026",
		PresidentName: "Prof. Alberto Angela",
		Commissioners: []Commissioner{
			{Name: "Prof.ssa Bianchi", Role: "commissario_interno", Subject: "Italiano"},
			{Name: "Prof. Neri", Role: "commissario_esterno", Subject: "Matematica"},
		},
	}
	err := svc.SaveCommission(ctx, comm)
	require.NoError(t, err)

	fetchedComm, err := svc.GetCommission(ctx, "class-5A", "2025/2026")
	require.NoError(t, err)
	assert.Equal(t, "Prof. Alberto Angela", fetchedComm.PresidentName)

	// 2. Save Student Record with Curriculum
	rec := &StudentMaturitaRecord{
		StudentID:   "std-1",
		ClassID:     "class-5A",
		StudentName: "Mario Rossi",
		SchoolYear:  "2025/2026",
		Credits: TrienniumCredits{
			Credit3rd:    11,
			Credit4th:    12,
			Credit5th:    14,
			TotalCredits: 37,
		},
		Scores: ExamScores{
			Written1Score: 18.0,
			Written2Score: 17.5,
			OralScore:     19.0,
			ExamTotal:     54.5,
			FinalScore:    92,
			Lode:          false,
		},
		CurriculumStudente: &CurriculumStudente{
			StudentCF:       "RSSMRA07A01H501U",
			StudentFullName: "Mario Rossi",
			PCTOHoursTotal:  210,
		},
		Status: "diplomato",
	}

	err = svc.SaveStudentRecord(ctx, rec)
	require.NoError(t, err)

	tabellone, err := svc.GetTabelloneClasse(ctx, "class-5A", "2025/2026")
	require.NoError(t, err)
	assert.Len(t, tabellone, 1)

	// 3. Export XML Curriculum dello Studente
	xmlOut, err := svc.ExportCurriculumXML(rec)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(xmlOut, "<?xml"))
	assert.Contains(t, xmlOut, "<codice_fiscale>RSSMRA07A01H501U</codice_fiscale>")
	assert.Contains(t, xmlOut, "<voto_finale>92</voto_finale>")
	assert.Contains(t, xmlOut, "<ore_totali>210</ore_totali>")
}
