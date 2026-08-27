package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/credits"
	"registro-backend/internal/recovery"
	"registro-backend/internal/support"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCreditsRepo struct {
	mock.Mock
}

func (m *mockCreditsRepo) SaveCredit(ctx context.Context, c *credits.StudentSchoolCredit) error {
	return m.Called(ctx, c).Error(0)
}
func (m *mockCreditsRepo) GetCreditByStudentAndYear(ctx context.Context, studentID, academicYear string, gradeLevel int) (*credits.StudentSchoolCredit, error) {
	return nil, nil
}
func (m *mockCreditsRepo) ListCreditsByClass(ctx context.Context, classID, academicYear string) ([]credits.StudentSchoolCredit, error) {
	args := m.Called(ctx, classID, academicYear)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]credits.StudentSchoolCredit), args.Error(1)
}
func (m *mockCreditsRepo) GetStudentCreditSummary(ctx context.Context, studentID string) (*credits.StudentCreditSummary, error) {
	args := m.Called(ctx, studentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*credits.StudentCreditSummary), args.Error(1)
}

func TestIntegration_Credits_AssignAndSummary(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockCreditsRepo)
	svc := credits.NewService(repo)
	handler := credits.NewHandler(svc)

	r := gin.New()
	protected := r.Group("")
	protected.Use(func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("school_id", "school-1")
		c.Next()
	})
	handler.RegisterRoutes(protected)

	// 1. Calculate endpoint test
	wCalc := httptest.NewRecorder()
	reqCalc, _ := http.NewRequest("GET", "/credits/calculate?grade_level=5&average=9.5&conduct=10&pcto_hours=100&has_extracurricular=true", nil)
	r.ServeHTTP(wCalc, reqCalc)
	assert.Equal(t, http.StatusOK, wCalc.Code)

	var calcRes credits.CreditCalculationResult
	_ = json.Unmarshal(wCalc.Body.Bytes(), &calcRes)
	assert.Equal(t, 14, calcRes.BaseCreditRangeMin)
	assert.Equal(t, 15, calcRes.BaseCreditRangeMax)
	assert.Equal(t, 15, calcRes.SuggestedCredit)

	// 2. Assign endpoint test
	assignBody := credits.AssignCreditRequest{
		StudentID:          "student-1",
		ClassID:            "class-5A",
		AcademicYear:       "2024/2025",
		GradeLevel:         5,
		GradeAverage:       9.5,
		ConductGrade:       10,
		AssignedCredit:     15,
		PCTOHours:          100,
		HasExtracurricular: true,
		DeliberationNotes:  "Eccellente percorso di studi e condotta",
	}
	repo.On("SaveCredit", mock.Anything, mock.MatchedBy(func(c *credits.StudentSchoolCredit) bool {
		return c.StudentID == "student-1" && c.AssignedCredit == 15
	})).Return(nil).Once()

	jsonAssign, _ := json.Marshal(assignBody)
	wAssign := httptest.NewRecorder()
	reqAssign, _ := http.NewRequest("POST", "/credits/assign", bytes.NewBuffer(jsonAssign))
	reqAssign.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wAssign, reqAssign)
	assert.Equal(t, http.StatusOK, wAssign.Code)

	// 3. Summary endpoint test
	repo.On("GetStudentCreditSummary", mock.Anything, "student-1").Return(&credits.StudentCreditSummary{
		StudentID:            "student-1",
		TotalTrienniumCredit: 38,
		MaxPossibleCredit:    40,
	}, nil).Once()

	wSum := httptest.NewRecorder()
	reqSum, _ := http.NewRequest("GET", "/credits/student/student-1/summary", nil)
	r.ServeHTTP(wSum, reqSum)
	assert.Equal(t, http.StatusOK, wSum.Code)

	var sumRes credits.StudentCreditSummary
	_ = json.Unmarshal(wSum.Body.Bytes(), &sumRes)
	assert.Equal(t, 38, sumRes.TotalTrienniumCredit)
}

type mockRecoveryIntegrationRepo struct {
	mock.Mock
}

func (m *mockRecoveryIntegrationRepo) CreateCourse(ctx context.Context, c *recovery.RecoveryCourse, sessions []recovery.RecoveryCourseSession, studentIDs []string) error {
	c.ID = "rec-course-1"
	return nil
}
func (m *mockRecoveryIntegrationRepo) GetCourseByID(ctx context.Context, id string) (*recovery.RecoveryCourse, error) {
	return &recovery.RecoveryCourse{ID: id, Title: "Corso PAI"}, nil
}
func (m *mockRecoveryIntegrationRepo) ListCourses(ctx context.Context, schoolID, academicYear, teacherID string) ([]recovery.RecoveryCourse, error) {
	return []recovery.RecoveryCourse{}, nil
}
func (m *mockRecoveryIntegrationRepo) UpdateCourseStatus(ctx context.Context, id, status string) error {
	return nil
}
func (m *mockRecoveryIntegrationRepo) UpdateStudentAttendance(ctx context.Context, courseID, studentID string, hours float64, notes string) error {
	return nil
}
func (m *mockRecoveryIntegrationRepo) RecordRecoveryTest(ctx context.Context, test *recovery.RecoveryTest) error {
	test.ID = "rec-test-1"
	return nil
}
func (m *mockRecoveryIntegrationRepo) ListRecoveryTests(ctx context.Context, schoolID, classID, studentID string) ([]recovery.RecoveryTest, error) {
	return []recovery.RecoveryTest{}, nil
}

func TestIntegration_Recovery_Workflow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockRecoveryIntegrationRepo)
	svc := recovery.NewService(repo)
	handler := recovery.NewHandler(svc)

	r := gin.New()
	protected := r.Group("")
	protected.Use(func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("teacher_id", "teacher-1")
		c.Set("school_id", "school-1")
		c.Next()
	})
	handler.RegisterRoutes(protected)

	// Record September test
	testReq := recovery.RecordTestOutcomeRequest{
		StudentID:     "student-1",
		SubjectID:     "sub-math",
		ClassID:       "class-2B",
		TestDate:      "2025-09-03",
		TestType:      "written",
		Grade:         7.0,
		VerbaleNumber: "VERB-REC-01",
	}
	jsonTest, _ := json.Marshal(testReq)
	wTest := httptest.NewRecorder()
	reqTest, _ := http.NewRequest("POST", "/recovery/tests", bytes.NewBuffer(jsonTest))
	reqTest.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wTest, reqTest)

	assert.Equal(t, http.StatusCreated, wTest.Code)
	var createdTest recovery.RecoveryTest
	_ = json.Unmarshal(wTest.Body.Bytes(), &createdTest)
	assert.Equal(t, "recuperato", createdTest.Outcome)
	assert.Equal(t, "Ammesso", createdTest.FinalDeliberation)
}

type mockSupportIntegrationRepo struct {
	mock.Mock
}

func (m *mockSupportIntegrationRepo) CreateDiaryEntry(ctx context.Context, entry *support.SupportDiaryEntry) error {
	entry.ID = "diary-1"
	return nil
}
func (m *mockSupportIntegrationRepo) ListDiaryEntries(ctx context.Context, schoolID, studentID, teacherID, classID string, isFamily bool) ([]support.SupportDiaryEntry, error) {
	return []support.SupportDiaryEntry{}, nil
}
func (m *mockSupportIntegrationRepo) DeleteDiaryEntry(ctx context.Context, id, teacherID string) error {
	return nil
}
func (m *mockSupportIntegrationRepo) CreatePeiGoal(ctx context.Context, goal *support.SupportPeiGoal) error {
	goal.ID = "goal-1"
	return nil
}
func (m *mockSupportIntegrationRepo) UpdatePeiGoalProgress(ctx context.Context, id, progressStatus string) error {
	return nil
}
func (m *mockSupportIntegrationRepo) ListPeiGoals(ctx context.Context, schoolID, studentID string) ([]support.SupportPeiGoal, error) {
	return []support.SupportPeiGoal{}, nil
}
func (m *mockSupportIntegrationRepo) DeletePeiGoal(ctx context.Context, id string) error {
	return nil
}

func TestIntegration_Support_Workflow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(mockSupportIntegrationRepo)
	svc := support.NewService(repo)
	handler := support.NewHandler(svc)

	r := gin.New()
	protected := r.Group("")
	protected.Use(func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("teacher_id", "teacher-1")
		c.Set("school_id", "school-1")
		c.Next()
	})
	handler.RegisterRoutes(protected)

	diaryReq := support.CreateDiaryEntryRequest{
		StudentID:          "student-1",
		ClassID:            "class-1A",
		EntryDate:          "2025-04-12",
		TimeSlot:           "2ª Ora (09:00-10:00)",
		ActivityType:       "laboratorio",
		TopicAndActivities: "Attività inclusiva di scienze con esperimenti pratici",
		StudentResponses:   "Grande entusiasmo e collaborazione con i compagni",
		IsSharedWithFamily: true,
	}
	jsonDiary, _ := json.Marshal(diaryReq)
	wDiary := httptest.NewRecorder()
	reqDiary, _ := http.NewRequest("POST", "/support/diaries", bytes.NewBuffer(jsonDiary))
	reqDiary.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wDiary, reqDiary)

	assert.Equal(t, http.StatusCreated, wDiary.Code)
}
