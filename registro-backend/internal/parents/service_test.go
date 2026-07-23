package parents

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParentsInitialization(t *testing.T) {
	svc := NewService(nil, nil, nil, nil, nil)
	assert.NotNil(t, svc)
}
