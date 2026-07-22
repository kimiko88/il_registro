package testhelpers

import (
	"context"
	"registro-backend/internal/auth"
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
func (m *MockAuthRepository) DisableMFA(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
func (m *MockAuthRepository) SaveTempMFASecret(ctx context.Context, userID, secret string) error {
	args := m.Called(ctx, userID, secret)
	return args.Error(0)
}
func (m *MockAuthRepository) ConfirmMFA(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
func (m *MockAuthRepository) GetMFASecret(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}
func (m *MockAuthRepository) CreateRecoveryCodes(ctx context.Context, userID string, codes []string) error {
	args := m.Called(ctx, userID, codes)
	return args.Error(0)
}
func (m *MockAuthRepository) GetRecoveryCodes(ctx context.Context, userID string) ([]*auth.MFARecoveryCode, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*auth.MFARecoveryCode), args.Error(1)
}
func (m *MockAuthRepository) UseRecoveryCode(ctx context.Context, userID, code string) error {
	args := m.Called(ctx, userID, code)
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
func (m *MockAuthRepository) UsePasswordResetToken(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockAuthRepository) UpdatePassword(ctx context.Context, userID, hash string) error {
	args := m.Called(ctx, userID, hash)
	return args.Error(0)
}

// MockGradesRepository mocks grades.Repository
type MockGradesRepository struct {
	mock.Mock
}

func (m *MockGradesRepository) Create(grade *grades.Grade) error {
	args := m.Called(grade)
	return args.Error(0)
}
func (m *MockGradesRepository) BatchCreate(gs []*grades.Grade) error {
	args := m.Called(gs)
	return args.Error(0)
}
func (m *MockGradesRepository) Update(grade *grades.Grade, history *grades.GradeHistory) error {
	args := m.Called(grade, history)
	return args.Error(0)
}
func (m *MockGradesRepository) Delete(id string, deletedBy string) error {
	args := m.Called(id, deletedBy)
	return args.Error(0)
}
func (m *MockGradesRepository) FindByID(id string) (*grades.Grade, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*grades.Grade), args.Error(1)
}
func (m *MockGradesRepository) FindByStudent(studentID string) ([]grades.Grade, error) {
	args := m.Called(studentID)
	// Return empty slice if nil to strictly match type
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.Grade), args.Error(1)
}
func (m *MockGradesRepository) FindByClassAndSubject(classID, subjectID string, semester int) ([]grades.Grade, error) {
	args := m.Called(classID, subjectID, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.Grade), args.Error(1)
}
func (m *MockGradesRepository) FindByClass(classID string, semester int) ([]grades.Grade, error) {
	args := m.Called(classID, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.Grade), args.Error(1)
}
func (m *MockGradesRepository) FindBySubject(subjectID string, semester int) ([]grades.Grade, error) {
	args := m.Called(subjectID, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.Grade), args.Error(1)
}
func (m *MockGradesRepository) FindWithFilter(filter grades.GradeFilter) ([]grades.Grade, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.Grade), args.Error(1)
}
func (m *MockGradesRepository) FindByTeacher(teacherID string) ([]grades.Grade, error) {
	args := m.Called(teacherID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.Grade), args.Error(1)
}
func (m *MockGradesRepository) GetHistory(gradeID string) ([]grades.GradeHistory, error) {
	args := m.Called(gradeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.GradeHistory), args.Error(1)
}
func (m *MockGradesRepository) CreateTest(test *grades.ClassTest) error {
	args := m.Called(test)
	return args.Error(0)
}
func (m *MockGradesRepository) FindTestsByClassAndSubject(classID string, subjectID string) ([]grades.ClassTest, error) {
	args := m.Called(classID, subjectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.ClassTest), args.Error(1)
}
func (m *MockGradesRepository) FindUpcomingTestsByClass(classID string) ([]grades.ClassTest, error) {
	args := m.Called(classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.ClassTest), args.Error(1)
}
func (m *MockGradesRepository) DeleteTest(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockGradesRepository) UpdateTest(test *grades.ClassTest) error {
	args := m.Called(test)
	return args.Error(0)
}
func (m *MockGradesRepository) FindGradesByTestID(testID string) ([]grades.Grade, error) {
	args := m.Called(testID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.Grade), args.Error(1)
}
func (m *MockGradesRepository) FindTestByID(id string) (*grades.ClassTest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*grades.ClassTest), args.Error(1)
}
func (m *MockGradesRepository) FindEnrolledSubjects(studentID string, semester int) ([]string, error) {
	args := m.Called(studentID, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockGradesRepository) GetWeightConfigs(schoolID, subjectID, classID string) ([]grades.GradeWeightConfig, error) {
	args := m.Called(schoolID, subjectID, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.GradeWeightConfig), args.Error(1)
}
func (m *MockGradesRepository) UpsertWeightConfig(cfg *grades.GradeWeightConfig) (*grades.GradeWeightConfig, error) {
	args := m.Called(cfg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*grades.GradeWeightConfig), args.Error(1)
}
func (m *MockGradesRepository) DeleteWeightConfig(id string) error {
	args := m.Called(id)
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
func (m *MockUsersRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
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
func (m *MockUsersRepository) LogAudit(ctx context.Context, log *users.AuditLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}
func (m *MockUsersRepository) GetAuditLogs(ctx context.Context, userID string, limit, offset int) ([]users.AuditLog, int, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]users.AuditLog), args.Int(1), args.Error(2)
}
func (m *MockUsersRepository) BulkCreate(ctx context.Context, usersList []users.User) (int, []string, error) {
	args := m.Called(ctx, usersList)
	// Return: count, errs, error
	var errs []string
	if args.Get(1) != nil {
		errs = args.Get(1).([]string)
	}
	return args.Int(0), errs, args.Error(2)
}
func (m *MockUsersRepository) BulkDelete(ctx context.Context, ids []string) (int, error) {
	args := m.Called(ctx, ids)
	return args.Int(0), args.Error(1)
}
func (m *MockUsersRepository) HardDelete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockUsersRepository) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	args := m.Called(ctx, parentID, studentID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUsersRepository) GetChildren(ctx context.Context, parentID string) ([]users.StudentChild, error) {
	args := m.Called(ctx, parentID)
	// Return empty slice if not mocked specifically to avoid panic on .Get(0) if not set up,
	// but standard testify Mock usage implies we set it up.
	// For compilation, signature match is enough.
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]users.StudentChild), args.Error(1)
}

func (m *MockUsersRepository) IsActive(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockUsersRepository) AddGuardian(ctx context.Context, studentID, parentID, relation string) error {
	args := m.Called(ctx, studentID, parentID, relation)
	return args.Error(0)
}
func (m *MockUsersRepository) GetStudentsByClass(ctx context.Context, classID string) ([]users.User, error) {
	args := m.Called(ctx, classID)
	return args.Get(0).([]users.User), args.Error(1)
}
func (m *MockUsersRepository) GetStudentProfile(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}
func (m *MockUsersRepository) GetParentProfile(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}
func (m *MockUsersRepository) RemoveGuardian(ctx context.Context, studentID, parentID string) error {
	args := m.Called(ctx, studentID, parentID)
	return args.Error(0)
}
func (m *MockUsersRepository) GetGuardians(ctx context.Context, studentID string) ([]users.GuardianInfo, error) {
	args := m.Called(ctx, studentID)
	return args.Get(0).([]users.GuardianInfo), args.Error(1)
}

// MockAnalyticsService mocks grades.AnalyticsService
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
func (m *MockAnalyticsService) GetSchoolStatistics(year string) (*grades.SchoolStatisticsResponse, error) {
	args := m.Called(year)
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
