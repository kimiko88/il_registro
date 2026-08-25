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
	assert.NotNil(t, origins1)

	// Memoized map remains cached
	origins2 := allowedOrigins()
	assert.Equal(t, origins1, origins2)
}
