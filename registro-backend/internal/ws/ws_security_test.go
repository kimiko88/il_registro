package ws

import (
	"encoding/json"
	"os"
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

func TestAllowedOriginsDynamicReload(t *testing.T) {
	os.Setenv("ALLOWED_ORIGINS", "http://domain1.com,http://domain2.com")
	origins := allowedOrigins()
	assert.True(t, origins["http://domain1.com"])
	assert.True(t, origins["http://domain2.com"])
	assert.False(t, origins["http://domain3.com"])

	// Update env and ensure dynamic reloading without server restart
	os.Setenv("ALLOWED_ORIGINS", "http://domain3.com")
	originsUpdated := allowedOrigins()
	assert.False(t, originsUpdated["http://domain1.com"])
	assert.True(t, originsUpdated["http://domain3.com"])
}
