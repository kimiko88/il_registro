package grades

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockServiceForDeprecationTest struct {
	mock.Mock
}

func (m *mockServiceForDeprecationTest) GetStudentGrades(ctx context.Context, actorID string, actorRole string, studentID string) ([]GradeResponse, error) {
	args := m.Called(ctx, actorID, actorRole, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]GradeResponse), args.Error(1)
}
func (m *mockServiceForDeprecationTest) GetStudentGradesWithFilter(ctx context.Context, actorID string, actorRole string, studentID string, filter GradeFilter) ([]GradeResponse, error) {
	args := m.Called(ctx, actorID, actorRole, studentID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]GradeResponse), args.Error(1)
}
func (m *mockServiceForDeprecationTest) GetStudentGradesPaged(ctx context.Context, actorID string, actorRole string, studentID string, filter GradeFilter) (*PaginatedGradesResponse, error) {
	args := m.Called(ctx, actorID, actorRole, studentID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*PaginatedGradesResponse), args.Error(1)
}
func (m *mockServiceForDeprecationTest) GetClassGrades(ctx context.Context, actorID string, actorRole string, classID string, filter GradeFilter) (*ClassGradesResponse, error) {
	args := m.Called(ctx, actorID, actorRole, classID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ClassGradesResponse), args.Error(1)
}
func (m *mockServiceForDeprecationTest) GetSubjectGrades(ctx context.Context, actorID string, actorRole string, subjectID string, filter GradeFilter) (*SubjectStatsResponse, error) {
	args := m.Called(ctx, actorID, actorRole, subjectID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SubjectStatsResponse), args.Error(1)
}
func (m *mockServiceForDeprecationTest) AddGrade(ctx context.Context, teacherID string, req CreateGradeRequest) (*GradeResponse, error) {
	args := m.Called(ctx, teacherID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GradeResponse), args.Error(1)
}
func (m *mockServiceForDeprecationTest) BatchCreateGrades(teacherID, actorRole, schoolID string, grades []*Grade) error {
	return m.Called(teacherID, actorRole, schoolID, grades).Error(0)
}
func (m *mockServiceForDeprecationTest) BulkImport(teacherID string, r io.Reader, semester int) (*ImportResult, error) {
	args := m.Called(teacherID, r, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ImportResult), args.Error(1)
}
func (m *mockServiceForDeprecationTest) Export(teacherID, schoolID string, filter GradeFilter, format string) ([]byte, string, error) {
	args := m.Called(teacherID, schoolID, filter, format)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).([]byte), args.String(1), args.Error(2)
}
func (m *mockServiceForDeprecationTest) UpdateGrade(ctx context.Context, teacherID string, gradeID string, req UpdateGradeRequest) (*GradeResponse, error) {
	args := m.Called(ctx, teacherID, gradeID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GradeResponse), args.Error(1)
}
func (m *mockServiceForDeprecationTest) DeleteGrade(ctx context.Context, teacherID string, gradeID string) error {
	return m.Called(ctx, teacherID, gradeID).Error(0)
}
func (m *mockServiceForDeprecationTest) GetMyGrades(ctx context.Context, studentID string, filter GradeFilter) (*MyGradesResponse, error) {
	args := m.Called(ctx, studentID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MyGradesResponse), args.Error(1)
}
func (m *mockServiceForDeprecationTest) GetMyAverages(ctx context.Context, studentID string) (*StudentAveragesResponse, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*StudentAveragesResponse), args.Error(1)
}
func (m *mockServiceForDeprecationTest) GetMyTrend(ctx context.Context, actorID string, actorRole string, studentID string, subjectID string) (*TrendResponse, error) {
	args := m.Called(ctx, actorID, actorRole, studentID, subjectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*TrendResponse), args.Error(1)
}
func (m *mockServiceForDeprecationTest) GetSemesterReport(ctx context.Context, actorID string, actorRole string, studentID string, semester int) (*SemesterReportResponse, error) {
	args := m.Called(ctx, actorID, actorRole, studentID, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SemesterReportResponse), args.Error(1)
}
func (m *mockServiceForDeprecationTest) GenerateSemesterReportPDF(studentID string, semester int) ([]byte, error) {
	args := m.Called(studentID, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}
func (m *mockServiceForDeprecationTest) ValidateParentGuardian(ctx context.Context, parentID, studentID string) error {
	return m.Called(ctx, parentID, studentID).Error(0)
}
func (m *mockServiceForDeprecationTest) GetChildGrades(ctx context.Context, parentID string, studentID string, filter GradeFilter) (*MyGradesResponse, error) {
	args := m.Called(ctx, parentID, studentID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MyGradesResponse), args.Error(1)
}
func (m *mockServiceForDeprecationTest) GetChildAverages(ctx context.Context, parentID string, studentID string) (*StudentAveragesResponse, error) {
	args := m.Called(ctx, parentID, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*StudentAveragesResponse), args.Error(1)
}
func (m *mockServiceForDeprecationTest) GetChildSemesterReport(ctx context.Context, parentID, studentID string, semester int) (*SemesterReportResponse, error) {
	args := m.Called(ctx, parentID, studentID, semester)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SemesterReportResponse), args.Error(1)
}
func (m *mockServiceForDeprecationTest) CreateTestWithGrades(teacherID string, req CreateClassTestRequest) (*ClassTest, error) {
	args := m.Called(teacherID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ClassTest), args.Error(1)
}
func (m *mockServiceForDeprecationTest) GetClassTests(ctx context.Context, actorID string, actorRole string, classID string, subjectID string) ([]ClassTestResponse, error) {
	args := m.Called(ctx, actorID, actorRole, classID, subjectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ClassTestResponse), args.Error(1)
}
func (m *mockServiceForDeprecationTest) GetUpcomingTestsByClass(ctx context.Context, actorID string, actorRole string, classID string) ([]ClassTestResponse, error) {
	args := m.Called(ctx, actorID, actorRole, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ClassTestResponse), args.Error(1)
}
func (m *mockServiceForDeprecationTest) DeleteClassTest(teacherID string, testID string) error {
	return m.Called(teacherID, testID).Error(0)
}
func (m *mockServiceForDeprecationTest) UpdateClassTest(teacherID string, testID string, req UpdateClassTestRequest) error {
	return m.Called(teacherID, testID, req).Error(0)
}
func (m *mockServiceForDeprecationTest) GetWeightConfigs(schoolID, subjectID, classID string) ([]GradeWeightConfig, error) {
	args := m.Called(schoolID, subjectID, classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]GradeWeightConfig), args.Error(1)
}
func (m *mockServiceForDeprecationTest) UpsertWeightConfig(actorID, actorRole, schoolID string, req UpsertWeightConfigRequest) (*GradeWeightConfig, error) {
	args := m.Called(actorID, actorRole, schoolID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GradeWeightConfig), args.Error(1)
}
func (m *mockServiceForDeprecationTest) DeleteWeightConfig(actorID, actorRole, schoolID, configID string) error {
	return m.Called(actorID, actorRole, schoolID, configID).Error(0)
}

func TestLegacyWeightConfig_DeprecationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockServiceForDeprecationTest)
	h := NewHandler(mockSvc, nil)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
	})
	h.RegisterRoutes(r.Group("/api/v1"))

	mockSvc.On("GetWeightConfigs", "school-1", "", "").
		Return([]GradeWeightConfig{}, nil).Once()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/grades/weights", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "true", w.Header().Get("Deprecation"))
	assert.Contains(t, w.Header().Get("Link"), "weight-configs")
}
