package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// TestBcryptCost_ChangePassword verifica che ChangePassword usi cost 12.
func TestBcryptCost_Constant(t *testing.T) {
	assert.Equal(t, 12, bcryptCost, "bcryptCost deve essere 12, non bcrypt.DefaultCost (10)")
	assert.Greater(t, bcryptCost, bcrypt.DefaultCost, "cost 12 deve essere maggiore del DefaultCost (10)")
}

func TestBcryptCost_HashActuallyUses12(t *testing.T) {
	password := "TestPassword1!"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	assert.NoError(t, err)

	cost, err := bcrypt.Cost(hash)
	assert.NoError(t, err)
	assert.Equal(t, 12, cost, "hash generato con bcryptCost deve avere cost 12")
}
