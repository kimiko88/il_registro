package testhelpers

import (
	"context"
	"registro-backend/internal/auth"
	"registro-backend/internal/classes"
	"registro-backend/internal/grades"
	"registro-backend/internal/schools"
	"registro-backend/internal/subjects"
	"registro-backend/internal/users"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockAuthRepository mocks auth.Repository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateUser(ctx context.Context, user *auth.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockAuthRepository) SchoolExists(ctx context.Context, schoolID string) (bool, error) {
	args := m.Called(ctx, schoolID)
	return args.Bool(0), args.Error(1)
}
func (m *MockAuthRepository) GetPasswordHistory(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockAuthRepository) AddPasswordHistory(ctx context.Context, userID, passwordHash string) error {
	args := m.Called(ctx, userID, passwordHash)
	return args.Error(0)
}
func (m *MockAuthRepository) GetUserByEmail(ctx context.Context, email string) (*auth.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.User), args.Error(1)
}
func (m *MockAuthRepository) GetUserByID(ctx context.Context, id string) (*auth.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.User), args.Error(1)
}
func (m *MockAuthRepository) UpdateUser(ctx context.Context, user *auth.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockAuthRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
func (m *MockAuthRepository) CreateRefreshToken(ctx context.Context, token *auth.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
func (m *MockAuthRepository) GetRefreshToken(ctx context.Context, token string) (*auth.RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.RefreshToken), args.Error(1)
}
func (m *MockAuthRepository) RevokeRefreshToken(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockAuthRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
func (m *MockAuthRepository) RotateRefreshTokenTx(ctx context.Context, oldID string, newRt *auth.RefreshToken) error {
	args := m.Called(ctx, oldID, newRt)
	return args.Error(0)
}
func (m *MockAuthRepository) RecordLoginAttempt(ctx context.Context, attempt *auth.LoginAttempt) error {
	args := m.Called(ctx, attempt)
	return args.Error(0)
}
func (m *MockAuthRepository) GetRecentLoginAttempts(ctx context.Context, email, ip string, since time.Time) (int, error) {
	args := m.Called(ctx, email, ip, since)
	return args.Int(0), args.Error(1)
}
func (m *MockAuthRepository) GetRecentLoginAttemptsByEmail(ctx context.Context, email string, since time.Time) (int, error) {
	args := m.Called(ctx, email, since)
	return args.Int(0), args.Error(1)
}
func (m *MockAuthRepository) GetRecentPasswordResets(ctx context.Context, userID string, since time.Time) (int, error) {
	args := m.Called(ctx, userID, since)
	return args.Int(0), args.Error(1)
}
func (m *MockAuthRepository) EnableMFA(ctx context.Context, userID, secret string) error {
	args := m.Called(ctx, userID, secret)
	return args.Error(0)
}
func (m *MockAuthRepository) ConfirmMFA(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
func (m *MockAuthRepository) ConfirmMFAAndSaveRecoveryCodesTx(ctx context.Context, userID string, codes []string) error {
	for _, call := range m.ExpectedCalls {
		if call.Method == "ConfirmMFAAndSaveRecoveryCodesTx" {
			args := m.Called(ctx, userID, codes)
			return args.Error(0)
		}
	}
	return nil
}
func (m *MockAuthRepository) DisableMFA(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
func (m *MockAuthRepository) GetMFASecret(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}
func (m *MockAuthRepository) SaveTempMFASecret(ctx context.Context, userID, secret string) error {
	args := m.Called(ctx, userID, secret)
	return args.Error(0)
}
func (m *MockAuthRepository) GetTempMFASecret(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}
func (m *MockAuthRepository) DeleteTempMFASecret(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
func (m *MockAuthRepository) CreateRecoveryCodes(ctx context.Context, userID string, hashedCodes []string) error {
	args := m.Called(ctx, userID, hashedCodes)
	return args.Error(0)
}
func (m *MockAuthRepository) GetRecoveryCodes(ctx context.Context, userID string) ([]*auth.MFARecoveryCode, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*auth.MFARecoveryCode), args.Error(1)
}
func (m *MockAuthRepository) UseRecoveryCode(ctx context.Context, userID, plainCode string) error {
	args := m.Called(ctx, userID, plainCode)
	return args.Error(0)
}
func (m *MockAuthRepository) CreatePasswordResetToken(ctx context.Context, token *auth.PasswordResetToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
func (m *MockAuthRepository) GetPasswordResetToken(ctx context.Context, token string) (*auth.PasswordResetToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.PasswordResetToken), args.Error(1)
}
func (m *MockAuthRepository) UsePasswordResetToken(ctx context.Context, tokenID string) error {
	args := m.Called(ctx, tokenID)
	return args.Error(0)
}
func (m *MockAuthRepository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	args := m.Called(ctx, userID, passwordHash)
	return args.Error(0)
}
func (m *MockAuthRepository) ResetPasswordTx(ctx context.Context, userID, passwordHash, tokenID string) error {
	args := m.Called(ctx, userID, passwordHash, tokenID)
	return args.Error(0)
}
func (m *MockAuthRepository) ChangePasswordTx(ctx context.Context, userID, passwordHash string) error {
	args := m.Called(ctx, userID, passwordHash)
	return args.Error(0)
}

// MockUsersRepository mocks users.Repository
type MockUsersRepository struct {
	mock.Mock
}

func (m *MockUsersRepository) Create(ctx context.Context, user *users.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockUsersRepository) GetByID(ctx context.Context, id string) (*users.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}
func (m *MockUsersRepository) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}
func (m *MockUsersRepository) Update(ctx context.Context, user *users.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockUsersRepository) IsActive(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}
func (m *MockUsersRepository) GetPasswordHistory(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockUsersRepository) AddPasswordHistory(ctx context.Context, userID, passwordHash string) error {
	args := m.Called(ctx, userID, passwordHash)
	return args.Error(0)
}
func (m *MockUsersRepository) SoftDelete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockUsersRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockUsersRepository) HardDelete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockUsersRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
func (m *MockUsersRepository) ClearTempMFASecret(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
func (m *MockUsersRepository) Restore(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockUsersRepository) List(ctx context.Context, filter users.UserFilter) ([]users.User, int, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]users.User), args.Int(1), args.Error(2)
}
func (m *MockUsersRepository) ListByIDs(ctx context.Context, ids []string) ([]users.User, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]users.User), args.Error(1)
}
func (m *MockUsersRepository) BulkCreate(ctx context.Context, usersList []users.User) (int, []string, error) {
	args := m.Called(ctx, usersList)
	return args.Int(0), args.Get(1).([]string), args.Error(2)
}
func (m *MockUsersRepository) BulkImport(ctx context.Context, usersList []users.User) (int, []string, error) {
	args := m.Called(ctx, usersList)
	return args.Int(0), args.Get(1).([]string), args.Error(2)
}
func (m *MockUsersRepository) BulkDelete(ctx context.Context, ids []string) (int, error) {
	args := m.Called(ctx, ids)
	return args.Int(0), args.Error(1)
}
func (m *MockUsersRepository) ChangePasswordTx(ctx context.Context, userID, newPasswordHash string) error {
	args := m.Called(ctx, userID, newPasswordHash)
	return args.Error(0)
}
func (m *MockUsersRepository) UpdatePassword(ctx context.Context, id string, passwordHash string) error {
	args := m.Called(ctx, id, passwordHash)
	return args.Error(0)
}
func (m *MockUsersRepository) LogAudit(ctx context.Context, audit *users.AuditLog) error {
	args := m.Called(ctx, audit)
	return args.Error(0)
}
func (m *MockUsersRepository) GetAuditLog(ctx context.Context, targetUserID string) ([]users.AuditLog, error) {
	args := m.Called(ctx, targetUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]users.AuditLog), args.Error(1)
}
func (m *MockUsersRepository) GetAuditLogs(ctx context.Context, userID string, limit, offset int) ([]users.AuditLog, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]users.AuditLog), args.Int(1), args.Error(2)
}
func (m *MockUsersRepository) AddGuardian(ctx context.Context, studentID, parentID, relationship string) error {
	args := m.Called(ctx, studentID, parentID, relationship)
	return args.Error(0)
}
func (m *MockUsersRepository) RemoveGuardian(ctx context.Context, studentID, parentID string) error {
	args := m.Called(ctx, studentID, parentID)
	return args.Error(0)
}
func (m *MockUsersRepository) GetGuardians(ctx context.Context, studentID string) ([]users.GuardianInfo, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]users.GuardianInfo), args.Error(1)
}
func (m *MockUsersRepository) GetFascicoloSummary(ctx context.Context, studentID string, isActive bool) (map[string]interface{}, error) {
	args := m.Called(ctx, studentID, isActive)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}
func (m *MockUsersRepository) GetChildren(ctx context.Context, parentID string) ([]users.StudentChild, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]users.StudentChild), args.Error(1)
}
func (m *MockUsersRepository) GetStudentsByClass(ctx context.Context, classID string) ([]users.User, error) {
	args := m.Called(ctx, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]users.User), args.Error(1)
}
func (m *MockUsersRepository) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	args := m.Called(ctx, parentID, studentID)
	return args.Bool(0), args.Error(1)
}
func (m *MockUsersRepository) GetStudentProfile(ctx context.Context, studentUserID string) (string, error) {
	args := m.Called(ctx, studentUserID)
	return args.String(0), args.Error(1)
}
func (m *MockUsersRepository) GetParentProfile(ctx context.Context, parentUserID string) (string, error) {
	args := m.Called(ctx, parentUserID)
	return args.String(0), args.Error(1)
}

// MockGradesRepository mocks grades.Repository
type MockGradesRepository struct {
	mock.Mock
}

func (m *MockGradesRepository) Create(ctx context.Context, grade *grades.Grade) error {
	for _, call := range m.ExpectedCalls {
		if call.Method == "Create" && len(call.Arguments) == 1 {
			return m.Called(grade).Error(0)
		}
	}
	args := m.Called(ctx, grade)
	return args.Error(0)
}
func (m *MockGradesRepository) BatchCreate(ctx context.Context, gradesList []*grades.Grade) error {
	for _, call := range m.ExpectedCalls {
		if call.Method == "BatchCreate" && len(call.Arguments) == 1 {
			return m.Called(gradesList).Error(0)
		}
	}
	args := m.Called(ctx, gradesList)
	return args.Error(0)
}
func (m *MockGradesRepository) GetByID(ctx context.Context, id string) (*grades.Grade, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "GetByID" && len(call.Arguments) == 1 {
			args := m.Called(id)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).(*grades.Grade), args.Error(1)
		}
	}
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*grades.Grade), args.Error(1)
}
func (m *MockGradesRepository) Update(ctx context.Context, grade *grades.Grade, history *grades.GradeHistory) error {
	for _, call := range m.ExpectedCalls {
		if call.Method == "Update" && len(call.Arguments) == 2 {
			return m.Called(grade, history).Error(0)
		}
	}
	args := m.Called(ctx, grade, history)
	return args.Error(0)
}
func (m *MockGradesRepository) Delete(ctx context.Context, id string, deletedBy string) error {
	for _, call := range m.ExpectedCalls {
		if call.Method == "Delete" && len(call.Arguments) == 2 {
			return m.Called(id, deletedBy).Error(0)
		}
	}
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}
func (m *MockGradesRepository) SoftDelete(ctx context.Context, id string, modifiedBy string) error {
	for _, call := range m.ExpectedCalls {
		if call.Method == "SoftDelete" && len(call.Arguments) == 2 {
			return m.Called(id, modifiedBy).Error(0)
		}
	}
	args := m.Called(ctx, id, modifiedBy)
	return args.Error(0)
}
func (m *MockGradesRepository) FindByID(ctx context.Context, id string) (*grades.Grade, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "FindByID" && len(call.Arguments) == 1 {
			args := m.Called(id)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).(*grades.Grade), args.Error(1)
		}
	}
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*grades.Grade), args.Error(1)
}
func (m *MockGradesRepository) FindByStudent(ctx context.Context, studentID string) ([]grades.Grade, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "FindByStudent" && len(call.Arguments) == 1 {
			args := m.Called(studentID)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).([]grades.Grade), args.Error(1)
		}
	}
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.Grade), args.Error(1)
}
func (m *MockGradesRepository) FindByClassAndSubject(ctx context.Context, classID string, subjectID string, semester int) ([]grades.Grade, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "FindByClassAndSubject" && len(call.Arguments) == 3 {
			args := m.Called(classID, subjectID, semester)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).([]grades.Grade), args.Error(1)
		}
	}
	args := m.Called(ctx, classID, subjectID, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.Grade), args.Error(1)
}
func (m *MockGradesRepository) FindByClass(ctx context.Context, classID string, semester int) ([]grades.Grade, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "FindByClass" && len(call.Arguments) == 2 {
			args := m.Called(classID, semester)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).([]grades.Grade), args.Error(1)
		}
	}
	args := m.Called(ctx, classID, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.Grade), args.Error(1)
}
func (m *MockGradesRepository) FindBySubject(ctx context.Context, subjectID string, semester int) ([]grades.Grade, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "FindBySubject" && len(call.Arguments) == 2 {
			args := m.Called(subjectID, semester)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).([]grades.Grade), args.Error(1)
		}
	}
	args := m.Called(ctx, subjectID, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.Grade), args.Error(1)
}
func (m *MockGradesRepository) FindWithFilter(ctx context.Context, filter grades.GradeFilter) ([]grades.Grade, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "FindWithFilter" && len(call.Arguments) == 1 {
			args := m.Called(filter)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).([]grades.Grade), args.Error(1)
		}
	}
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.Grade), args.Error(1)
}
func (m *MockGradesRepository) FindWithFilterPaginated(ctx context.Context, filter grades.GradeFilter) ([]grades.Grade, int, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "FindWithFilterPaginated" && len(call.Arguments) == 1 {
			args := m.Called(filter)
			if args.Get(0) == nil {
				return nil, args.Int(1), args.Error(2)
			}
			return args.Get(0).([]grades.Grade), args.Int(1), args.Error(2)
		}
	}
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]grades.Grade), args.Int(1), args.Error(2)
}
func (m *MockGradesRepository) FindByTeacher(ctx context.Context, teacherID string) ([]grades.Grade, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "FindByTeacher" && len(call.Arguments) == 1 {
			args := m.Called(teacherID)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).([]grades.Grade), args.Error(1)
		}
	}
	args := m.Called(ctx, teacherID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.Grade), args.Error(1)
}
func (m *MockGradesRepository) GetHistory(ctx context.Context, gradeID string) ([]grades.GradeHistory, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "GetHistory" && len(call.Arguments) == 1 {
			args := m.Called(gradeID)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).([]grades.GradeHistory), args.Error(1)
		}
	}
	args := m.Called(ctx, gradeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.GradeHistory), args.Error(1)
}
func (m *MockGradesRepository) FindEnrolledSubjects(ctx context.Context, studentID string, semester int) ([]string, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "FindEnrolledSubjects" && len(call.Arguments) == 2 {
			args := m.Called(studentID, semester)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).([]string), args.Error(1)
		}
	}
	args := m.Called(ctx, studentID, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockGradesRepository) CreateTest(ctx context.Context, test *grades.ClassTest) error {
	for _, call := range m.ExpectedCalls {
		if call.Method == "CreateTest" && len(call.Arguments) == 1 {
			return m.Called(test).Error(0)
		}
	}
	args := m.Called(ctx, test)
	return args.Error(0)
}
func (m *MockGradesRepository) FindTestByID(ctx context.Context, testID string) (*grades.ClassTest, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "FindTestByID" && len(call.Arguments) == 1 {
			args := m.Called(testID)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).(*grades.ClassTest), args.Error(1)
		}
	}
	args := m.Called(ctx, testID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*grades.ClassTest), args.Error(1)
}
func (m *MockGradesRepository) FindTestsByClassAndSubject(ctx context.Context, classID string, subjectID string) ([]grades.ClassTest, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "FindTestsByClassAndSubject" && len(call.Arguments) == 2 {
			args := m.Called(classID, subjectID)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).([]grades.ClassTest), args.Error(1)
		}
	}
	args := m.Called(ctx, classID, subjectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.ClassTest), args.Error(1)
}
func (m *MockGradesRepository) FindUpcomingTestsByClass(ctx context.Context, classID string) ([]grades.ClassTest, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "FindUpcomingTestsByClass" && len(call.Arguments) == 1 {
			args := m.Called(classID)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).([]grades.ClassTest), args.Error(1)
		}
	}
	args := m.Called(ctx, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.ClassTest), args.Error(1)
}
func (m *MockGradesRepository) DeleteTest(ctx context.Context, id string) error {
	for _, call := range m.ExpectedCalls {
		if call.Method == "DeleteTest" && len(call.Arguments) == 1 {
			return m.Called(id).Error(0)
		}
	}
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockGradesRepository) UpdateTest(ctx context.Context, test *grades.ClassTest) error {
	for _, call := range m.ExpectedCalls {
		if call.Method == "UpdateTest" && len(call.Arguments) == 1 {
			return m.Called(test).Error(0)
		}
	}
	args := m.Called(ctx, test)
	return args.Error(0)
}
func (m *MockGradesRepository) FindGradesByTestID(ctx context.Context, testID string) ([]grades.Grade, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "FindGradesByTestID" && len(call.Arguments) == 1 {
			args := m.Called(testID)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).([]grades.Grade), args.Error(1)
		}
	}
	args := m.Called(ctx, testID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.Grade), args.Error(1)
}
func (m *MockGradesRepository) DeleteWeightConfig(ctx context.Context, id string) error {
	for _, call := range m.ExpectedCalls {
		if call.Method == "DeleteWeightConfig" && len(call.Arguments) == 1 {
			return m.Called(id).Error(0)
		}
	}
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockGradesRepository) GetWeightConfigs(ctx context.Context, schoolID string, subjectID string, classID string) ([]grades.GradeWeightConfig, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "GetWeightConfigs" && len(call.Arguments) == 3 {
			args := m.Called(schoolID, subjectID, classID)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).([]grades.GradeWeightConfig), args.Error(1)
		}
	}
	args := m.Called(ctx, schoolID, subjectID, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.GradeWeightConfig), args.Error(1)
}
func (m *MockGradesRepository) SaveWeightConfig(ctx context.Context, cfg *grades.GradeWeightConfig) error {
	for _, call := range m.ExpectedCalls {
		if call.Method == "SaveWeightConfig" && len(call.Arguments) == 1 {
			return m.Called(cfg).Error(0)
		}
	}
	args := m.Called(ctx, cfg)
	return args.Error(0)
}
func (m *MockGradesRepository) UpsertWeightConfig(ctx context.Context, cfg *grades.GradeWeightConfig) (*grades.GradeWeightConfig, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "UpsertWeightConfig" && len(call.Arguments) == 1 {
			args := m.Called(cfg)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).(*grades.GradeWeightConfig), args.Error(1)
		}
	}
	args := m.Called(ctx, cfg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*grades.GradeWeightConfig), args.Error(1)
}


func (m *MockGradesRepository) GetStudentClassAndSchoolInfo(ctx context.Context, studentID string) (string, string, string, string, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "GetStudentClassAndSchoolInfo" {
			args := m.Called(ctx, studentID)
			return args.String(0), args.String(1), args.String(2), args.String(3), args.Error(4)
		}
	}
	return "", "", "", "", nil
}

func (m *MockGradesRepository) GetTeacherNamesByClass(ctx context.Context, classID string) (map[string]string, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "GetTeacherNamesByClass" {
			args := m.Called(ctx, classID)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).(map[string]string), args.Error(1)
		}
	}
	return nil, nil
}

func (m *MockGradesRepository) GetSubjectNamesMap(ctx context.Context, schoolID string) (map[string]string, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "GetSubjectNamesMap" {
			args := m.Called(ctx, schoolID)
			if args.Get(0) == nil {
				return nil, args.Error(1)
			}
			return args.Get(0).(map[string]string), args.Error(1)
		}
	}
	return nil, nil
}

func (m *MockGradesRepository) GetScrutinyRecordSummary(ctx context.Context, studentID string, semester int) (float64, float64, bool, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "GetScrutinyRecordSummary" {
			args := m.Called(ctx, studentID, semester)
			return args.Get(0).(float64), args.Get(1).(float64), args.Bool(2), args.Error(3)
		}
	}
	return 0, 0, false, nil
}

func (m *MockGradesRepository) GetStudentAbsenceCountForPeriod(ctx context.Context, studentID, startD, endD string) (int, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "GetStudentAbsenceCountForPeriod" {
			args := m.Called(ctx, studentID, startD, endD)
			return args.Int(0), args.Error(1)
		}
	}
	return 0, nil
}

func (m *MockGradesRepository) GetClassSubjectAverage(ctx context.Context, classID, subjectID string, semester int, studentID string) (float64, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "GetClassSubjectAverage" {
			args := m.Called(ctx, classID, subjectID, semester, studentID)
			return args.Get(0).(float64), args.Error(1)
		}
	}
	return -1, nil
}

func (m *MockGradesRepository) CheckClassAccessPermission(ctx context.Context, actorID, actorRole, classID string) (bool, error) {
	for _, call := range m.ExpectedCalls {
		if call.Method == "CheckClassAccessPermission" {
			args := m.Called(ctx, actorID, actorRole, classID)
			return args.Bool(0), args.Error(1)
		}
	}
	return true, nil
}

type MockAnalyticsService struct {
	mock.Mock
}

func (m *MockAnalyticsService) GetStudentAverage(studentID string, subjectID string) (float64, error) {
	args := m.Called(studentID, subjectID)
	return args.Get(0).(float64), args.Error(1)
}
func (m *MockAnalyticsService) GetClassAverage(classID string, subjectID string) (float64, error) {
	args := m.Called(classID, subjectID)
	return args.Get(0).(float64), args.Error(1)
}
func (m *MockAnalyticsService) GetClassAnalysis(classID string, semester int) (*grades.AnalyticsClassResponse, error) {
	args := m.Called(classID, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*grades.AnalyticsClassResponse), args.Error(1)
}
func (m *MockAnalyticsService) GetSubjectAnalysis(subjectID string, semester int) (*grades.AnalyticsSubjectResponse, error) {
	args := m.Called(subjectID, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*grades.AnalyticsSubjectResponse), args.Error(1)
}
func (m *MockAnalyticsService) GetStudentProfile(studentID string, semester int) (*grades.AnalyticsStudentResponse, error) {
	args := m.Called(studentID, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*grades.AnalyticsStudentResponse), args.Error(1)
}
func (m *MockAnalyticsService) GetSchoolStatistics(year string, schoolID ...string) (*grades.SchoolStatisticsResponse, error) {
	var sid string
	if len(schoolID) > 0 {
		sid = schoolID[0]
	}
	args := m.Called(year, sid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*grades.SchoolStatisticsResponse), args.Error(1)
}

// MockSchoolsRepository mocks schools.Repository
type MockSchoolsRepository struct {
	mock.Mock
}

func (m *MockSchoolsRepository) Create(ctx context.Context, school *schools.School) error {
	args := m.Called(ctx, school)
	return args.Error(0)
}
func (m *MockSchoolsRepository) GetByID(ctx context.Context, id string) (*schools.School, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schools.School), args.Error(1)
}
func (m *MockSchoolsRepository) List(ctx context.Context, params *schools.ListParams) ([]*schools.School, int, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*schools.School), args.Int(1), args.Error(2)
}
func (m *MockSchoolsRepository) Update(ctx context.Context, id string, req *schools.UpdateSchoolRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}
func (m *MockSchoolsRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockSubjectsRepository mocks subjects.Repository
type MockSubjectsRepository struct {
	mock.Mock
}

func (m *MockSubjectsRepository) Create(ctx context.Context, subject *subjects.Subject) error {
	args := m.Called(ctx, subject)
	return args.Error(0)
}
func (m *MockSubjectsRepository) List(ctx context.Context, schoolID string) ([]subjects.Subject, error) {
	args := m.Called(ctx, schoolID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]subjects.Subject), args.Error(1)
}
func (m *MockSubjectsRepository) Get(ctx context.Context, id string) (*subjects.Subject, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*subjects.Subject), args.Error(1)
}
func (m *MockSubjectsRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockSubjectsRepository) Update(ctx context.Context, subject *subjects.Subject) error {
	args := m.Called(ctx, subject)
	return args.Error(0)
}

// MockClassesRepository mocks classes.Repository
type MockClassesRepository struct {
	mock.Mock
}

func (m *MockClassesRepository) Create(ctx context.Context, class *classes.Class) error {
	args := m.Called(ctx, class)
	return args.Error(0)
}
func (m *MockClassesRepository) List(ctx context.Context, schoolID string, academicYear string) ([]classes.Class, error) {
	args := m.Called(ctx, schoolID, academicYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]classes.Class), args.Error(1)
}
func (m *MockClassesRepository) Get(ctx context.Context, id string) (*classes.Class, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*classes.Class), args.Error(1)
}
func (m *MockClassesRepository) Update(ctx context.Context, class *classes.Class) error {
	args := m.Called(ctx, class)
	return args.Error(0)
}
func (m *MockClassesRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockClassesRepository) ListByTeacher(ctx context.Context, teacherUserID string, schoolYear string) ([]classes.Class, error) {
	args := m.Called(ctx, teacherUserID, schoolYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]classes.Class), args.Error(1)
}
func (m *MockClassesRepository) AssignSubject(ctx context.Context, classID string, subjectID string, teacherID *string, hours float64) error {
	args := m.Called(ctx, classID, subjectID, teacherID, hours)
	return args.Error(0)
}
func (m *MockClassesRepository) UnassignSubject(ctx context.Context, assignmentID string) error {
	args := m.Called(ctx, assignmentID)
	return args.Error(0)
}
func (m *MockClassesRepository) GetClassSubjects(ctx context.Context, classID string) ([]classes.ClassSubject, error) {
	args := m.Called(ctx, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]classes.ClassSubject), args.Error(1)
}
func (m *MockClassesRepository) GetClassGuardians(ctx context.Context, classID string) ([]classes.GuardianInfo, error) {
	args := m.Called(ctx, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]classes.GuardianInfo), args.Error(1)
}
func (m *MockClassesRepository) GetLessonTopics(ctx context.Context, classID string) ([]classes.LessonTopic, error) {
	args := m.Called(ctx, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]classes.LessonTopic), args.Error(1)
}
func (m *MockClassesRepository) GetDisciplinaryNotes(ctx context.Context, classID string) ([]classes.DisciplinaryNoteReport, error) {
	args := m.Called(ctx, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]classes.DisciplinaryNoteReport), args.Error(1)
}
func (m *MockClassesRepository) BulkMigrateStudents(ctx context.Context, migrations []classes.StudentMigrationItem) error {
	args := m.Called(ctx, migrations)
	return args.Error(0)
}
