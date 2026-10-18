package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiChannelDispatcher(t *testing.T) {
	dispatcher := NewMultiChannelDispatcher(nil, nil, nil)
	ctx := context.Background()

	t.Run("Dispatch to App IO, WhatsApp, and SMS successfully", func(t *testing.T) {
		req := MultiChannelRequest{
			Channels: []string{ChannelAppIO, ChannelWhatsApp, ChannelSMS},
			Recipients: []RecipientInfo{
				{
					FiscalCode: "RSSMRA80A01H501U",
					Phone:      "+393331122334",
					Name:       "Mario Rossi",
				},
				{
					FiscalCode: "BNCGLI85B41H501K",
					Phone:      "+393339988776",
					Name:       "Giulia Bianchi",
				},
			},
			Subject:        "Allerta Meteo Rossa - Chiusura Plesso",
			Body:           "Si comunica che domani le lezioni saranno sospese causa allerta meteo.",
			Category:       "EMERGENCY_WEATHER",
			EmergencyAlert: true,
			SenderID:       "LICEO_ROMA",
		}

		resp, err := dispatcher.Dispatch(ctx, req)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.BatchID)
		assert.True(t, resp.EmergencyAlert)
		assert.Equal(t, 6, resp.TotalDispatches) // 2 recipients * 3 channels
		assert.Equal(t, 6, resp.SuccessCount)
		assert.Equal(t, 0, resp.FailureCount)
	})

	t.Run("Validation failure with invalid fiscal code or invalid phone", func(t *testing.T) {
		req := MultiChannelRequest{
			Channels: []string{ChannelAppIO, ChannelSMS},
			Recipients: []RecipientInfo{
				{
					FiscalCode: "SHORT_CF",
					Phone:      "123",
				},
			},
			Subject: "Test",
			Body:    "Test message",
		}

		resp, err := dispatcher.Dispatch(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, 2, resp.TotalDispatches)
		assert.Equal(t, 0, resp.SuccessCount)
		assert.Equal(t, 2, resp.FailureCount)
		assert.NotEmpty(t, resp.Results[0].Error)
		assert.NotEmpty(t, resp.Results[1].Error)
	})

	t.Run("Error when no recipients or no channels specified", func(t *testing.T) {
		_, err := dispatcher.Dispatch(ctx, MultiChannelRequest{
			Channels:   []string{ChannelSMS},
			Recipients: []RecipientInfo{},
		})
		assert.Error(t, err)

		_, err = dispatcher.Dispatch(ctx, MultiChannelRequest{
			Channels: []string{},
			Recipients: []RecipientInfo{
				{Phone: "+393331122334"},
			},
		})
		assert.Error(t, err)
	})
}

func TestNotificationsHandler_MultiChannelEndpoint(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)
	r := setupTestNotificationsRouter(svc, "admin-1", "principal")

	payload := MultiChannelRequest{
		Channels: []string{ChannelAppIO, ChannelSMS},
		Recipients: []RecipientInfo{
			{
				FiscalCode: "RSSMRA80A01H501U",
				Phone:      "+393331122334",
				Name:       "Mario Rossi",
			},
		},
		Subject:  "Convocazione Consiglio Straordinario",
		Body:     "Si comunica la convocazione del consiglio straordinario.",
		Category: "CIRCULAR",
	}

	bodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/notifications/multi-channel", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp MultiChannelResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 2, resp.TotalDispatches)
	assert.Equal(t, 2, resp.SuccessCount)

	// Unauthorized role test
	studentRouter := setupTestNotificationsRouter(svc, "student-1", "student")
	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/notifications/multi-channel", bytes.NewBuffer(bodyBytes))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	studentRouter.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusForbidden, w2.Code)
}
