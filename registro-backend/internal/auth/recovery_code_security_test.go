package auth

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUseRecoveryCode_ConsumesAtomicallyForUser(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewRepository(db)

	userID := "user-1"
	code := "ABCD1234"
	hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.MinCost)
	require.NoError(t, err)

	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, code, used, used_at, created_at FROM mfa_recovery_codes WHERE user_id = $1 AND used = false")).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "code", "used", "used_at", "created_at"}).
			AddRow("recovery-1", userID, string(hash), false, nil, time.Now()))
	mockDB.ExpectExec(regexp.QuoteMeta("UPDATE mfa_recovery_codes SET used = true, used_at = $1 WHERE id = $2 AND user_id = $3 AND used = false")).
		WithArgs(sqlmock.AnyArg(), "recovery-1", userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, repo.UseRecoveryCode(context.Background(), userID, code))
	assert.NoError(t, mockDB.ExpectationsWereMet())
}

func TestUseRecoveryCode_RejectsConcurrentReuse(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	repo := NewRepository(db)

	userID := "user-1"
	code := "ABCD1234"
	hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.MinCost)
	require.NoError(t, err)

	mockDB.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, code, used, used_at, created_at FROM mfa_recovery_codes WHERE user_id = $1 AND used = false")).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "code", "used", "used_at", "created_at"}).
			AddRow("recovery-1", userID, string(hash), false, nil, time.Now()))
	mockDB.ExpectExec(regexp.QuoteMeta("UPDATE mfa_recovery_codes SET used = true, used_at = $1 WHERE id = $2 AND user_id = $3 AND used = false")).
		WithArgs(sqlmock.AnyArg(), "recovery-1", userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	assert.ErrorIs(t, repo.UseRecoveryCode(context.Background(), userID, code), ErrInvalidMFAToken)
	assert.NoError(t, mockDB.ExpectationsWereMet())
}
