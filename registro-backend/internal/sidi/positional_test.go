package sidi

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPositionalRecordLengths(t *testing.T) {
	t.Run("ESITI.TXT record length is exactly 120 chars", func(t *testing.T) {
		records := []EsitoPosizionaleRecord{
			{
				CodiceScuola:     "RMIS09900B",
				AnnoScolastico:   "202526",
				CodiceFiscale:    "RSSMRA08A01H501Z",
				CodiceSIDI:       "SIDI1049281",
				ClasseSezione:    "3A",
				Sessione:         "GIU",
				Esito:            "A",
				Crediti:          11,
				DebitiAssegnati:  0,
				DebitiRecuperati: "S",
				MediaVoti:        "0780",
				Note:             "AMMESSO",
			},
		}

		txt, errs := BuildEsitiTXT(records)
		assert.Empty(t, errs)
		lines := strings.Split(strings.TrimRight(txt, "\r\n"), "\r\n")
		require.Len(t, lines, 1)
		assert.Equal(t, RecordLenEsiti, len(lines[0]))
	})

	t.Run("PAGELLE.TXT record length is exactly 100 chars", func(t *testing.T) {
		records := []PagellaPosizionaleRecord{
			{
				CodiceScuola:   "RMIS09900B",
				AnnoScolastico: "202526",
				Periodo:        "2Q",
				CodiceFiscale:  "RSSMRA08A01H501Z",
				CodiceMateria:  "MAT",
				Voto:           "080",
				Condotta:       "09",
				AssenzeOre:     6,
				Giudizio:       "OTTIMO",
				Note:           "PARTECIPAZIONE ATTIVA",
			},
		}

		txt, errs := BuildPagelleTXT(records)
		assert.Empty(t, errs)
		lines := strings.Split(strings.TrimRight(txt, "\r\n"), "\r\n")
		require.Len(t, lines, 1)
		assert.Equal(t, RecordLenPagelle, len(lines[0]))
	})

	t.Run("ANS.TXT record length is exactly 140 chars", func(t *testing.T) {
		records := []ANSPosizionaleRecord{
			{
				CodiceScuola:    "RMIS09900B",
				AnnoScolastico:  "202526",
				CodiceFiscale:   "RSSMRA08A01H501Z",
				CodiceSIDI:      "SIDI1049281",
				Cognome:         "ROSSI",
				Nome:            "MARIO",
				Sesso:           "M",
				DataNascita:     "20080115",
				ComuneNascita:   "H501",
				Cittadinanza:    "ITA",
				ClasseSezione:   "3A",
				IndirizzoStudi:  "LI01",
				StatoIscrizione: "ATT",
			},
		}

		txt, errs := BuildANSTXT(records)
		assert.Empty(t, errs)
		lines := strings.Split(strings.TrimRight(txt, "\r\n"), "\r\n")
		require.Len(t, lines, 1)
		assert.Equal(t, RecordLenANS, len(lines[0]))
	})

	t.Run("PASSWEB_INPS.TXT record length is exactly 160 chars", func(t *testing.T) {
		records := []PassWebPosizionaleRecord{
			{
				CodiceFiscale:         "RSSMRA75C12H501W",
				CodiceScuola:          "RMIS09900B",
				Profilo:               "DOCENTE",
				TipoContratto:         "TI",
				DataInizio:            "20100901",
				DataFine:              "20260831",
				RegimePrev:            "TFR",
				RetribuzioneLordaCent: 3250000,
				RetribuzioneTabellare: 2800000,
				GiorniUtili:           365,
				CausaCessazione:       "IN_SERVIZIO",
			},
		}

		txt, errs := BuildPassWebTXT(records)
		assert.Empty(t, errs)
		lines := strings.Split(strings.TrimRight(txt, "\r\n"), "\r\n")
		require.Len(t, lines, 1)
		assert.Equal(t, RecordLenPassWeb, len(lines[0]))
	})
}

func TestServiceGeneratePositional(t *testing.T) {
	svc := NewService(nil)
	ctx := context.Background()

	types := []string{PositionalTypeEsiti, PositionalTypePagelle, PositionalTypeANS, PositionalTypePassWeb}
	expectedLengths := map[string]int{
		PositionalTypeEsiti:   RecordLenEsiti,
		PositionalTypePagelle: RecordLenPagelle,
		PositionalTypeANS:     RecordLenANS,
		PositionalTypePassWeb: RecordLenPassWeb,
	}

	for _, pt := range types {
		t.Run("Generate "+pt, func(t *testing.T) {
			resp, content, err := svc.GeneratePositional(ctx, "school-123", "user-1", PositionalExportRequest{
				Type:       pt,
				SchoolYear: "2025/2026",
			})
			require.NoError(t, err)
			assert.True(t, resp.ValidationOK)
			assert.Equal(t, pt, resp.Type)
			assert.Equal(t, expectedLengths[pt], resp.LineLength)
			assert.NotEmpty(t, content)
			assert.NotEmpty(t, resp.ContentPreview)
		})
	}

	t.Run("Generate Invalid Type fails", func(t *testing.T) {
		_, _, err := svc.GeneratePositional(ctx, "school-123", "user-1", PositionalExportRequest{
			Type: "UNKNOWN_TYPE",
		})
		assert.Error(t, err)
	})

	t.Run("Get Positional Templates returns all 4 types", func(t *testing.T) {
		tmpl := svc.GetPositionalTemplates()
		assert.Contains(t, tmpl, "ESITI")
		assert.Contains(t, tmpl, "PAGELLE")
		assert.Contains(t, tmpl, "ANS")
		assert.Contains(t, tmpl, "PASSWEB")
	})
}
