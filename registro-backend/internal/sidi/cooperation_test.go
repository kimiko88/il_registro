package sidi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSIDICooperation(t *testing.T) {
	svc := NewService(nil)
	ctx := context.Background()

	t.Run("Sync Student SIDI Codes", func(t *testing.T) {
		students := []StudenteSIDI{
			{CodiceFiscale: "RSSMRA08A01H501Z", CodiceSIDI: ""},
			{CodiceFiscale: "BNCGLI08B41H501A", CodiceSIDI: "SIDI-EXISTING"},
			{CodiceFiscale: "VRDLGI08C12H501B", CodiceSIDI: ""},
		}

		resp, err := svc.SyncStudentCodes(ctx, "school-1", students)
		require.NoError(t, err)
		assert.Equal(t, 3, resp.TotalProcessed)
		assert.Equal(t, 2, resp.TotalUpdated)
		assert.NotEmpty(t, resp.UpdatedCodes["RSSMRA08A01H501Z"])
		assert.NotEmpty(t, resp.UpdatedCodes["VRDLGI08C12H501B"])
		assert.Empty(t, resp.UpdatedCodes["BNCGLI08B41H501A"])
	})

	t.Run("Push Scrutiny Results to MIM WebService", func(t *testing.T) {
		req := PushScrutinyResultsRequest{
			SchoolYear: "2025/2026",
			Sessione:   "GIUGNO",
			ClassID:    "class-3A",
			Results: []ScrutinioSIDI{
				{CodiceSIDI: "SIDI-1000001", Classe: "3A", EsitoFinale: "AMMESSO", CreditiFormat: 11},
				{CodiceSIDI: "SIDI-1000002", Classe: "3A", EsitoFinale: "SOSPESO_GIUDIZIO"},
			},
		}

		resp, err := svc.PushScrutinyResults(ctx, "school-1", req)
		require.NoError(t, err)
		assert.Equal(t, "TRASMESSO_CON_SUCCESSO", resp.Status)
		assert.Equal(t, 2, resp.RecordsAccepted)
		assert.NotEmpty(t, resp.ProtocolloMIM)
	})

	t.Run("Get Cooperation Config", func(t *testing.T) {
		cfg, err := svc.GetCooperationConfig(ctx, "school-1")
		require.NoError(t, err)
		assert.NotEmpty(t, cfg.EndpointURL)
		assert.NotEmpty(t, cfg.CertificatoPostazione)
	})
}
