package ws

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWSPingStrictParsing(t *testing.T) {
	// False positive payload containing "PING" string inside title
	payload := []byte(`{"type":"NOTE_ADDED","payload":{"title":"Avviso PING ritardato"}}`)

	var cm wsClientMsg
	err := json.Unmarshal(payload, &cm)
	assert.NoError(t, err)
	assert.NotEqual(t, "PING", cm.Type, "Payload with PING substring must not match cm.Type == PING")

	// Valid PING payload
	pingPayload := []byte(`{"type":"PING"}`)
	err = json.Unmarshal(pingPayload, &cm)
	assert.NoError(t, err)
	assert.Equal(t, "PING", cm.Type)
}

func TestAllowedOriginsOnceCaching(t *testing.T) {
	origins := allowedOrigins()
	assert.NotNil(t, origins)
}
