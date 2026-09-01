package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/users"
	"registro-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type DummyUserRepo struct {
	mock.Mock
}

func (m *DummyUserRepo) Create(ctx context.Context, u *users.User) error             { return nil }
func (m *DummyUserRepo) GetByID(ctx context.Context, id string) (*users.User, error) { return nil, nil }
func (m *DummyUserRepo) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	return nil, nil
}
func (m *DummyUserRepo) Update(ctx context.Context, u *users.User) error { return nil }
func (m *DummyUserRepo) GetPasswordHistory(ctx context.Context, userID string) ([]string, error) {
	return nil, nil
}
func (m *DummyUserRepo) AddPasswordHistory(ctx context.Context, userID, passwordHash string) error {
	return nil
}
func (m *DummyUserRepo) Delete(ctx context.Context, id string) error  { return nil }
func (m *DummyUserRepo) Restore(ctx context.Context, id string) error { return nil }
func (m *DummyUserRepo) List(ctx context.Context, filter users.UserFilter) ([]users.User, int, error) {
	return nil, 0, nil
}
func (m *DummyUserRepo) ListByIDs(ctx context.Context, ids []string) ([]users.User, error) {
	return nil, nil
}
func (m *DummyUserRepo) LogAudit(ctx context.Context, log *users.AuditLog) error { return nil }
func (m *DummyUserRepo) GetAuditLogs(ctx context.Context, userID string, limit, offset int) ([]users.AuditLog, int, error) {
	return nil, 0, nil
}
func (m *DummyUserRepo) BulkCreate(ctx context.Context, u []users.User) (int, []string, error) {
	return 0, nil, nil
}
func (m *DummyUserRepo) BulkDelete(ctx context.Context, ids []string) (int, error)    { return 0, nil }
func (m *DummyUserRepo) HardDelete(ctx context.Context, id string) error              { return nil }
func (m *DummyUserRepo) RevokeAllUserTokens(ctx context.Context, userID string) error { return nil }
func (m *DummyUserRepo) ClearTempMFASecret(ctx context.Context, userID string) error  { return nil }
func (m *DummyUserRepo) IsGuardian(ctx context.Context, parentUserID, studentUserID string) (bool, error) {
	return false, nil
}
func (m *DummyUserRepo) GetChildren(ctx context.Context, parentUserID string) ([]users.StudentChild, error) {
	return nil, nil
}
func (m *DummyUserRepo) GetStudentsByClass(ctx context.Context, classID string) ([]users.User, error) {
	return nil, nil
}
func (m *DummyUserRepo) GetStudentProfile(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *DummyUserRepo) GetParentProfile(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *DummyUserRepo) AddGuardian(ctx context.Context, studentProfileID, parentProfileID, relationship string) error {
	return nil
}
func (m *DummyUserRepo) RemoveGuardian(ctx context.Context, studentProfileID, parentProfileID string) error {
	return nil
}
func (m *DummyUserRepo) GetGuardians(ctx context.Context, studentProfileID string) ([]users.GuardianInfo, error) {
	return nil, nil
}
func (m *DummyUserRepo) GetFascicoloSummary(ctx context.Context, studentID string, isActive bool) (map[string]interface{}, error) {
	return nil, nil
}
func (m *DummyUserRepo) IsActive(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}
func (m *DummyUserRepo) ChangePasswordTx(ctx context.Context, userID, newPasswordHash string) error {
	return nil
}

func TestMiddleware_ActiveCache_CachesDisabledUserAndInvalidates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserRepo := new(DummyUserRepo)
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	tm := jwt.NewTokenManager(privateKey, &privateKey.PublicKey)

	mw := NewMiddleware(tm, mockUserRepo)

	userID := "user-disabled-cache-test"
	token, err := tm.GenerateAccessToken(userID, "test@example.com", RoleTeacher, "school-1")
	assert.NoError(t, err)

	// First request: DB returns false (user is disabled)
	mockUserRepo.On("IsActive", mock.Anything, userID).Return(false, nil).Once()

	r := gin.New()
	r.Use(mw.Authenticate())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Request 1: Should hit DB and return 403 Forbidden
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusForbidden, w1.Code)

	// Request 2: Should use cached false status without calling IsActive DB again
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusForbidden, w2.Code)

	// Invalidate cache for user
	mw.InvalidateUserActiveCache(userID)

	// Request 3: After invalidation, DB is called again. Configure DB to return true this time.
	mockUserRepo.On("IsActive", mock.Anything, userID).Return(true, nil).Once()

	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/test", nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)

	mockUserRepo.AssertExpectations(t)
}
