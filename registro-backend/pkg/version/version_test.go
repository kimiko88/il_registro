package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionInfo(t *testing.T) {
	info := Get()
	assert.NotEmpty(t, info.Version)
	assert.NotEmpty(t, info.GoVersion)
	assert.NotEmpty(t, info.OS)
	assert.NotEmpty(t, info.Arch)

	str := String()
	assert.Contains(t, str, "v")
}
