package parents

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"registro-backend/internal/attendance"
	"registro-backend/internal/communications"
	"registro-backend/internal/grades"
	"registro-backend/internal/users"
)

type mockUserRepoForParentTest struct{ mock.Mock }

func (m *mockUserRepoForParentTest) GetChildren(ctx context.Context, parentID string) ([]users.StudentChild, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]users.StudentChild), args.Error(1)
}
func (m *mockUserRepoForParentTest) GetByID(ctx context.Context, id string) (*users.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}
func (m *mockUserRepoForParentTest) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	args := m.Called(ctx, parentID, studentID)
	return args.Bool(0), args.Error(1)
}
func (m *mockUserRepoForParentTest) ListByIDs(ctx context.Context, ids []string) ([]users.User, error) {
	return nil, nil
}
func (m *mockUserRepoForParentTest) Create(ctx context.Context, u *users.User) error { return nil }
func (m *mockUserRepoForParentTest) Update(ctx context.Context, u *users.User) error { return nil }
func (m *mockUserRepoForParentTest) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	return nil, nil
}
func (m *mockUserRepoForParentTest) GetPasswordHistory(ctx context.Context, id string) ([]string, error) {
	return nil, nil
}
func (m *mockUserRepoForParentTest) AddPasswordHistory(ctx context.Context, id, hash string) error {
	return nil
}
func (m *mockUserRepoForParentTest) RevokeAllUserTokens(ctx context.Context, id string) error {
	return nil
}
func (m *mockUserRepoForParentTest) ClearTempMFASecret(ctx context.Context, id string) error {
	return nil
}
func (m *mockUserRepoForParentTest) GetAllActive(ctx context.Context, schoolID string) ([]users.User, error) {
	return nil, nil
}
func (m *mockUserRepoForParentTest) SoftDelete(ctx context.Context, id string) error { return nil }
func (m *mockUserRepoForParentTest) Delete(ctx context.Context, id string) error     { return nil }
func (m *mockUserRepoForParentTest) Restore(ctx context.Context, id string) error    { return nil }
func (m *mockUserRepoForParentTest) HardDelete(ctx context.Context, id string) error { return nil }
func (m *mockUserRepoForParentTest) List(ctx context.Context, filter users.UserFilter) ([]users.User, int, error) {
	return nil, 0, nil
}
func (m *mockUserRepoForParentTest) LogAudit(ctx context.Context, log *users.AuditLog) error {
	return nil
}
func (m *mockUserRepoForParentTest) GetAuditLogs(ctx context.Context, userID string, limit, offset int) ([]users.AuditLog, int, error) {
	return nil, 0, nil
}
func (m *mockUserRepoForParentTest) AddGuardian(ctx context.Context, parentID, studentID, relType string) error {
	return nil
}
func (m *mockUserRepoForParentTest) RemoveGuardian(ctx context.Context, studentProfileID, parentProfileID string) error {
	return nil
}
func (m *mockUserRepoForParentTest) GetGuardians(ctx context.Context, studentProfileID string) ([]users.GuardianInfo, error) {
	return nil, nil
}
func (m *mockUserRepoForParentTest) GetStudentsByClass(ctx context.Context, classID string) ([]users.User, error) {
	return nil, nil
}
func (m *mockUserRepoForParentTest) GetStudentProfile(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *mockUserRepoForParentTest) GetParentProfile(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *mockUserRepoForParentTest) GetFascicoloSummary(ctx context.Context, studentID string, isActive bool) (map[string]interface{}, error) {
	return nil, nil
}
func (m *mockUserRepoForParentTest) IsActive(ctx context.Context, id string) (bool, error) {
	return true, nil
}
func (m *mockUserRepoForParentTest) BulkCreate(ctx context.Context, users []users.User) (int, []string, error) {
	return 0, nil, nil
}
func (m *mockUserRepoForParentTest) BulkDelete(ctx context.Context, ids []string) (int, error) {
	return 0, nil
}

type mockGradesRepoForParentTest struct{ mock.Mock }

func (m *mockGradesRepoForParentTest) FindByStudent(studentID string) ([]grades.Grade, error) {
	args := m.Called(studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]grades.Grade), args.Error(1)
}
func (m *mockGradesRepoForParentTest) Create(g *grades.Grade) error             { return nil }
func (m *mockGradesRepoForParentTest) BatchCreate(grades []*grades.Grade) error { return nil }
func (m *mockGradesRepoForParentTest) Update(g *grades.Grade, h *grades.GradeHistory) error {
	return nil
}
func (m *mockGradesRepoForParentTest) Delete(id, deletedBy string) error         { return nil }
func (m *mockGradesRepoForParentTest) FindByID(id string) (*grades.Grade, error) { return nil, nil }
func (m *mockGradesRepoForParentTest) FindByClass(classID string, semester int) ([]grades.Grade, error) {
	return nil, nil
}
func (m *mockGradesRepoForParentTest) FindByClassAndSubject(classID, subjectID string, semester int) ([]grades.Grade, error) {
	return nil, nil
}
func (m *mockGradesRepoForParentTest) FindBySubject(subjectID string, semester int) ([]grades.Grade, error) {
	return nil, nil
}
func (m *mockGradesRepoForParentTest) FindWithFilter(filter grades.GradeFilter) ([]grades.Grade, error) {
	return nil, nil
}
func (m *mockGradesRepoForParentTest) FindWithFilterPaginated(filter grades.GradeFilter) ([]grades.Grade, int, error) {
	return nil, 0, nil
}
func (m *mockGradesRepoForParentTest) FindByTeacher(teacherID string) ([]grades.Grade, error) {
	return nil, nil
}
func (m *mockGradesRepoForParentTest) GetHistory(gradeID string) ([]grades.GradeHistory, error) {
	return nil, nil
}
func (m *mockGradesRepoForParentTest) FindEnrolledSubjects(studentID string, semester int) ([]string, error) {
	return nil, nil
}
func (m *mockGradesRepoForParentTest) CreateTest(test *grades.ClassTest) error { return nil }
func (m *mockGradesRepoForParentTest) FindTestsByClassAndSubject(classID, subjectID string) ([]grades.ClassTest, error) {
	return nil, nil
}
func (m *mockGradesRepoForParentTest) FindUpcomingTestsByClass(classID string) ([]grades.ClassTest, error) {
	return nil, nil
}
func (m *mockGradesRepoForParentTest) DeleteTest(id string) error              { return nil }
func (m *mockGradesRepoForParentTest) UpdateTest(test *grades.ClassTest) error { return nil }
func (m *mockGradesRepoForParentTest) FindGradesByTestID(testID string) ([]grades.Grade, error) {
	return nil, nil
}
func (m *mockGradesRepoForParentTest) FindTestByID(id string) (*grades.ClassTest, error) {
	return nil, nil
}
func (m *mockGradesRepoForParentTest) GetWeightConfigs(schoolID, subjectID, classID string) ([]grades.GradeWeightConfig, error) {
	return nil, nil
}
func (m *mockGradesRepoForParentTest) UpsertWeightConfig(cfg *grades.GradeWeightConfig) (*grades.GradeWeightConfig, error) {
	return nil, nil
}
func (m *mockGradesRepoForParentTest) DeleteWeightConfig(id string) error { return nil }
func (m *mockGradesRepoForParentTest) GetStudentClassAndSchoolInfo(ctx context.Context, studentID string) (studentName, className, classID, schoolID string, err error) {
	return "", "", "", "", nil
}
func (m *mockGradesRepoForParentTest) GetTeacherNamesByClass(ctx context.Context, classID string) (map[string]string, error) {
	return nil, nil
}
func (m *mockGradesRepoForParentTest) GetSubjectNamesMap(ctx context.Context, schoolID string) (map[string]string, error) {
	return nil, nil
}
func (m *mockGradesRepoForParentTest) GetScrutinyRecordSummary(ctx context.Context, studentID string, semester int) (behaviorGrade, scholasticCredit float64, found bool, err error) {
	return 0, 0, false, nil
}
func (m *mockGradesRepoForParentTest) GetStudentAbsenceCountForPeriod(ctx context.Context, studentID, startD, endD string) (int, error) {
	return 0, nil
}
func (m *mockGradesRepoForParentTest) GetClassSubjectAverage(ctx context.Context, classID, subjectID string, semester int, studentID string) (float64, error) {
	return 0, nil
}
func (m *mockGradesRepoForParentTest) CheckClassAccessPermission(ctx context.Context, actorID, actorRole, classID string) (bool, error) {
	return true, nil
}

type mockAttRepoForParentTest struct{ mock.Mock }

func (m *mockAttRepoForParentTest) GetStats(studentProfileID string) (*attendance.SummaryResponse, error) {
	return &attendance.SummaryResponse{TotalAbsences: 2, JustifiedCount: 2}, nil
}
func (m *mockAttRepoForParentTest) Create(att *attendance.Attendance) error         { return nil }
func (m *mockAttRepoForParentTest) BatchCreate(atts []*attendance.Attendance) error { return nil }
func (m *mockAttRepoForParentTest) Update(att *attendance.Attendance) error         { return nil }
func (m *mockAttRepoForParentTest) DeleteByClassDateHour(schoolID, classID string, date time.Time, hour int) error {
	return nil
}
func (m *mockAttRepoForParentTest) FindByID(id string) (*attendance.Attendance, error) {
	return nil, nil
}
func (m *mockAttRepoForParentTest) FindByClassAndDate(classID string, date time.Time) ([]attendance.Attendance, error) {
	return nil, nil
}
func (m *mockAttRepoForParentTest) FindByStudent(studentID string, startDate, endDate time.Time) ([]attendance.Attendance, error) {
	return nil, nil
}
func (m *mockAttRepoForParentTest) FindUnjustifiedByStudent(studentID string) ([]attendance.Attendance, error) {
	return nil, nil
}
func (m *mockAttRepoForParentTest) GetStatsBatch(ctx context.Context, studentIDs []string) (map[string]*attendance.SummaryResponse, error) {
	return nil, nil
}
func (m *mockAttRepoForParentTest) CountDistinctDays(studentID string) (int, error) { return 0, nil }
func (m *mockAttRepoForParentTest) GetAnalytics(ctx context.Context, schoolID string) (*attendance.AnalyticsResponse, error) {
	return nil, nil
}
func (m *mockAttRepoForParentTest) CreateJustification(j *attendance.Justification) error { return nil }
func (m *mockAttRepoForParentTest) UpdateJustification(j *attendance.Justification) error { return nil }
func (m *mockAttRepoForParentTest) ProcessJustificationTx(ctx context.Context, j *attendance.Justification, teacherID string, approve bool) error {
	return nil
}
func (m *mockAttRepoForParentTest) FindJustificationByID(id string) (*attendance.Justification, error) {
	return nil, nil
}
func (m *mockAttRepoForParentTest) FindPendingJustifications(classID, schoolID string) ([]attendance.Justification, error) {
	return nil, nil
}
func (m *mockAttRepoForParentTest) FindPendingJustificationsForTeacher(ctx context.Context, teacherID, schoolID string) ([]attendance.Justification, error) {
	return nil, nil
}
func (m *mockAttRepoForParentTest) DeleteJustification(id string) error        { return nil }
func (m *mockAttRepoForParentTest) DeletePendingJustification(id string) error { return nil }
func (m *mockAttRepoForParentTest) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	return true, nil
}
func (m *mockAttRepoForParentTest) IsStudentInClass(ctx context.Context, studentID, classID string) (bool, error) {
	return true, nil
}
func (m *mockAttRepoForParentTest) AreStudentsInClass(ctx context.Context, studentIDs []string, classID string) (map[string]bool, error) {
	return nil, nil
}
func (m *mockAttRepoForParentTest) IsClassInSchool(ctx context.Context, classID, schoolID string) (bool, error) {
	return true, nil
}
func (m *mockAttRepoForParentTest) IsTeacherSubstitute(ctx context.Context, teacherID, classID string, date time.Time, hour int) (bool, error) {
	return false, nil
}
func (m *mockAttRepoForParentTest) HasOverlappingJustification(ctx context.Context, studentID string, startDate, endDate time.Time) (bool, error) {
	return false, nil
}
func (m *mockAttRepoForParentTest) GetMonthlyBreakdown(ctx context.Context, studentID, schoolYear string) ([]attendance.MonthlyBreakdownRow, error) {
	return nil, nil
}
func (m *mockAttRepoForParentTest) JustifyAbsenceByParent(attendanceID, studentID, reason, notes string) error {
	return nil
}
func (m *mockAttRepoForParentTest) GetStudentAttendanceStats(studentID string) (*attendance.AttendanceStats, error) {
	return &attendance.AttendanceStats{DaysPresent: 50, DaysAbsent: 2}, nil
}

type mockCommsRepoForParentTest struct{ mock.Mock }

func (m *mockCommsRepoForParentTest) ListBacheca(ctx context.Context, schoolID, userID string) ([]*communications.Message, error) {
	args := m.Called(ctx, schoolID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*communications.Message), args.Error(1)
}
func (m *mockCommsRepoForParentTest) Create(ctx context.Context, msg *communications.Message) error {
	return nil
}
func (m *mockCommsRepoForParentTest) List(ctx context.Context, userID, schoolID string) ([]*communications.Message, error) {
	return nil, nil
}
func (m *mockCommsRepoForParentTest) Delete(ctx context.Context, id string) error { return nil }
func (m *mockCommsRepoForParentTest) Sign(ctx context.Context, commID, userID string) error {
	return nil
}
func (m *mockCommsRepoForParentTest) SignWithIP(ctx context.Context, commID, userID, ip string) error {
	return nil
}
func (m *mockCommsRepoForParentTest) GetSignatures(ctx context.Context, commID string) ([]string, error) {
	return nil, nil
}
func (m *mockCommsRepoForParentTest) GetSignatureReport(ctx context.Context, commID string) (*communications.SignatureReportResponse, error) {
	return nil, nil
}
func (m *mockCommsRepoForParentTest) Get(ctx context.Context, id string) (*communications.Message, error) {
	return nil, nil
}
func (m *mockCommsRepoForParentTest) Update(ctx context.Context, id, subject, body string) error {
	return nil
}
func (m *mockCommsRepoForParentTest) MarkAsRead(ctx context.Context, commID, userID, ip string) error {
	return nil
}
func (m *mockCommsRepoForParentTest) GetUnreadUsers(ctx context.Context, commID string) ([]string, error) {
	return nil, nil
}
func (m *mockCommsRepoForParentTest) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	return 0, nil
}
func (m *mockCommsRepoForParentTest) ListCircolari(ctx context.Context, schoolID, userID, year string) ([]*communications.Message, error) {
	return nil, nil
}
func (m *mockCommsRepoForParentTest) Ack(ctx context.Context, commID, userID string) error {
	return nil
}

func TestGetDashboard_UsesRealChildSchoolID(t *testing.T) {
	uRepo := new(mockUserRepoForParentTest)
	gRepo := new(mockGradesRepoForParentTest)
	aRepo := new(mockAttRepoForParentTest)
	cRepo := new(mockCommsRepoForParentTest)

	svc := NewService(nil, uRepo, gRepo, aRepo, cRepo)
	ctx := context.Background()

	realSchoolID := "school-real-uuid-42"
	studentUser := &users.User{
		ID:       "student-u1",
		SchoolID: &realSchoolID,
	}

	uRepo.On("GetChildren", ctx, "parent-1").Return([]users.StudentChild{
		{ID: "student-prof-1", UserID: "student-u1", FirstName: "Marco", LastName: "Rossi"},
	}, nil)

	uRepo.On("GetByID", ctx, "student-u1").Return(studentUser, nil)
	gRepo.On("FindByStudent", "student-u1").Return([]grades.Grade{}, nil)

	// Verify ListBacheca is called with the REAL schoolID, not "school-id"
	cRepo.On("ListBacheca", ctx, realSchoolID, "student-u1").Return([]*communications.Message{
		{ID: "msg-1", Subject: "Circolare Viaggio", RequiresSignature: true, IsSigned: false},
	}, nil)

	dash, err := svc.GetDashboard(ctx, "parent-1")
	assert.NoError(t, err)
	assert.NotNil(t, dash)
	assert.Len(t, dash.Children, 1)
	assert.Len(t, dash.Children[0].PendingCirculars, 1)

	cRepo.AssertCalled(t, "ListBacheca", ctx, realSchoolID, "student-u1")
}
