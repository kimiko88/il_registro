package albopretorio

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlboPretorioLifecycle(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	// 1. Publish Act (Determina dirigenziale)
	item1, err := svc.PublishAct(ctx, "school-1", "user-ds", PublishActRequest{
		Category:                CategoryDetermineDirigente,
		Subject:                 "Determina a contrarre fornitura arredi innovativi PNRR",
		DocumentFileURL:         "https://scuola.edu.it/docs/det_01.pdf",
		DocumentSHA256:          "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		DaysDuration:            15,
		IsTransparencySection:   true,
		TransparencyMacroFamily: "Bandi di gara e contratti",
		CIGCode:                 "Z1A2B3C4D5",
		AwardedAmount:           14500.00,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, item1.ID)
	assert.Contains(t, item1.RepertoryCode, "/00001")
	assert.Equal(t, "in_pubblicazione", item1.Status)

	// 2. Publish Second Act
	item2, err := svc.PublishAct(ctx, "school-1", "user-ds", PublishActRequest{
		Category:        CategoryDelibereConsiglio,
		Subject:         "Delibera approvazione PTOF triennale",
		DocumentFileURL: "https://scuola.edu.it/docs/delibera_ptof.pdf",
		DocumentSHA256:  "cf23df2207d99a74fbe169e3eba035e633b65d94",
		DaysDuration:    15,
	})
	require.NoError(t, err)
	assert.Contains(t, item2.RepertoryCode, "/00002")

	// 3. List public: both visible
	list, err := svc.ListPublicAlbo(ctx, "school-1", "", "", 0, false)
	require.NoError(t, err)
	assert.Len(t, list, 2)

	// 4. Filter by category
	delibere, err := svc.ListPublicAlbo(ctx, "school-1", CategoryDelibereConsiglio, "", 0, false)
	require.NoError(t, err)
	assert.Len(t, delibere, 1)
	assert.Equal(t, item2.ID, delibere[0].ID)

	// 5. Try defission before 15 days without force -> must fail
	_, err = svc.DefiggiAtto(ctx, item1.ID, "Prof. Giovanni Rossi", false)
	assert.ErrorIs(t, err, ErrActNotExpired)

	// 6. Force defission (or after 15 days) -> generates relata
	defisso, err := svc.DefiggiAtto(ctx, item1.ID, "Prof. Giovanni Rossi", true)
	require.NoError(t, err)
	assert.Equal(t, "defisso", defisso.Status)
	assert.Contains(t, defisso.RelataText, "RELATA DI PUBBLICAZIONE")
	assert.Equal(t, "Prof. Giovanni Rossi", defisso.RelataSignedBy)

	// 7. Verify Historical Archive
	archived, err := svc.ListPublicAlbo(ctx, "school-1", "", "", 0, true)
	require.NoError(t, err)
	assert.Len(t, archived, 1)
	assert.Equal(t, item1.ID, archived[0].ID)

	// 8. Generate Certificate (Certificato di Pubblicazione DS)
	cert, err := svc.GeneraCertificato(ctx, item1.ID, "Prof. Giovanni Rossi")
	require.NoError(t, err)
	assert.Equal(t, item1.RepertoryCode, cert.RepertoryCode)
	assert.Contains(t, cert.LegalAttestation, "CERTIFICATO DI AVVENUTA PUBBLICAZIONE")

	// 9. Generate ANAC XML Dataset
	xmlOutput, err := svc.GenerateANACXML(ctx, "school-1", item1.RepertoryYear)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(xmlOutput, "<?xml"))
	assert.Contains(t, xmlOutput, "<cig>Z1A2B3C4D5</cig>")
	assert.Contains(t, xmlOutput, "<importoAggiudicazione>14500</importoAggiudicazione>")
}
