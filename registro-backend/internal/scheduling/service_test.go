package scheduling

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServiceInitialization(t *testing.T) {
	repo := new(MockRepo)
	tRepo := new(MockTeacherRepo)

	svc := NewService(repo, tRepo, nil, nil, nil)
	assert.NotNil(t, svc)
}

func TestValidatorValidation(t *testing.T) {
	v := NewValidator()
	s := &ColloquioSlot{
		Date: time.Now().AddDate(0, 0, 1),
	}
	assert.NoError(t, v.ValidateSlot(s))
}
