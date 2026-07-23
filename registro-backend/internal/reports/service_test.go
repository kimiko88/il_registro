package reports

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReportsInitialization(t *testing.T) {
	svc := NewService(nil)
	assert.NotNil(t, svc)
}
