package scrutiny

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapGradeToReligionJudgment(t *testing.T) {
	cases := []struct {
		val      float64
		expected string
	}{
		{10.0, "Ottimo"},
		{9.75, "Ottimo"},
		{9.5, "Ottimo"},
		{9.0, "Distinto"},
		{8.0, "Distinto"},
		{7.5, "Distinto"},
		{7.0, "Buono"},
		{6.5, "Buono"},
		{6.0, "Sufficiente"},
		{5.5, "Sufficiente"},
		{5.0, "Insufficiente"},
		{4.0, "Insufficiente"},
		{1.0, "Insufficiente"},
		{0.0, "Non classificabile"},
		{-1.0, "Non classificabile"},
	}

	for _, tc := range cases {
		actual := mapGradeToReligionJudgment(tc.val)
		assert.Equal(t, tc.expected, actual, "val %f should map to %s", tc.val, tc.expected)
	}
}

func TestGeneratePagellaPDF_WithReligionSubjects(t *testing.T) {
	matrix := &ScrutinyMatrix{
		ClassID:    "cls-1",
		Semester:   1,
		PeriodType: "semester_1",
		Subjects: []SubjectInfo{
			{ID: "sub-ita", Name: "Italiano", IsReligion: false},
			{ID: "sub-rel", Name: "Religione Cattolica", IsReligion: true, IsJudgmentOnly: true},
		},
		Students: []StudentScrutinyRow{
			{
				StudentID:      "stu-avvalente",
				StudentName:    "Rossi Mario",
				ReligionChoice: "avvalente",
				SubjectData: map[string]SubjectAverages{
					"sub-ita": {Average: 8.0, GradeCount: 3, Proposed: 8.0},
					"sub-rel": {Average: 9.5, GradeCount: 2, Proposed: 10.0, ProposedJudgment: "Ottimo"},
				},
				Record: &ScrutinyRecord{
					ConductGrade:  9,
					FinalDecision: "Ammesso",
				},
			},
			{
				StudentID:      "stu-non-avvalente",
				StudentName:    "Bianchi Luca",
				ReligionChoice: "non_avvalente",
				SubjectData: map[string]SubjectAverages{
					"sub-ita": {Average: 6.5, GradeCount: 2, Proposed: 7.0},
					"sub-rel": {Average: 0, GradeCount: 0, Proposed: 0, ProposedJudgment: "Non avvalente"},
				},
				Record: &ScrutinyRecord{
					ConductGrade:  8,
					FinalDecision: "Ammesso",
				},
			},
			{
				StudentID:      "stu-alt",
				StudentName:    "Verdi Sara",
				ReligionChoice: "attivita_alternativa",
				SubjectData: map[string]SubjectAverages{
					"sub-ita": {Average: 7.0, GradeCount: 3, Proposed: 7.0},
					"sub-rel": {Average: 0, GradeCount: 0, Proposed: 0, ProposedJudgment: "Attività alternativa"},
				},
				Record: &ScrutinyRecord{
					ConductGrade:  10,
					FinalDecision: "Ammesso",
				},
			},
		},
	}

	t.Run("avvalente student PDF renders judgment without error", func(t *testing.T) {
		pdfBytes, err := GeneratePagellaPDF(matrix, "stu-avvalente")
		require.NoError(t, err)
		assert.NotEmpty(t, pdfBytes)
		// PDF magic header %PDF-
		assert.True(t, len(pdfBytes) > 100)
		assert.Equal(t, "%PDF-", string(pdfBytes[:5]))
	})

	t.Run("non-avvalente student PDF renders Non Avvalente without error", func(t *testing.T) {
		pdfBytes, err := GeneratePagellaPDF(matrix, "stu-non-avvalente")
		require.NoError(t, err)
		assert.NotEmpty(t, pdfBytes)
		assert.Equal(t, "%PDF-", string(pdfBytes[:5]))
	})

	t.Run("attivita alternativa student PDF renders Att. Alternativa without error", func(t *testing.T) {
		pdfBytes, err := GeneratePagellaPDF(matrix, "stu-alt")
		require.NoError(t, err)
		assert.NotEmpty(t, pdfBytes)
		assert.Equal(t, "%PDF-", string(pdfBytes[:5]))
	})

	t.Run("ZIP export with all student types succeeds", func(t *testing.T) {
		zipBytes, err := GenerateClassPagelleZIP(matrix)
		require.NoError(t, err)
		assert.NotEmpty(t, zipBytes)
		// ZIP magic header PK\x03\x04
		assert.Equal(t, "PK\x03\x04", string(zipBytes[:4]))
	})
}
