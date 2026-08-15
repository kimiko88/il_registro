package ws

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWS_AllowedOriginsCaching(t *testing.T) {
	os.Setenv("ALLOWED_ORIGINS", "https://app.scuola.it")
	defer os.Unsetenv("ALLOWED_ORIGINS")

	origins1 := allowedOrigins()
	assert.True(t, origins1["https://app.scuola.it"])

	// Mutate env variable; cached Once map MUST preserve first value
	os.Setenv("ALLOWED_ORIGINS", "https://hacked.com")

	origins2 := allowedOrigins()
	assert.True(t, origins2["https://app.scuola.it"])
	assert.False(t, origins2["https://hacked.com"])
}
