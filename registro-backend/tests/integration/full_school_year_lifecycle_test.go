package integration

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/attendance"
	"registro-backend/internal/classes"
	"registro-backend/internal/grades"
	"registro-backend/internal/schools"
	"registro-backend/internal/scrutiny"
	"registro-backend/internal/users"
	"registro-backend/tests/testhelpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAttRepo struct {
	mock.Mock
}

func (m *mockAttRepo) Create(att *attendance.Attendance) error { return nil }
func (m *mockAttRepo) BatchCreate(atts []*attendance.Attendance) error {
	args := m.Called(atts)
	return args.Error(0)
}
func (m *mockAttRepo) Update(att *attendance.Attendance) error { return nil }
func (m *mockAttRepo) DeleteByClassDateHour(classID string, date time.Time, hour int) error { return nil }
func (m *mockAttRepo) FindByID(id string) (*attendance.Attendance, error) { return nil, nil }
func (m *mockAttRepo) FindByClassAndDate(classID string, date time.Time) ([]attendance.Attendance, error) {
	return nil, nil
}
func (m *mockAttRepo) FindByStudent(studentID string, startDate, endDate time.Time) ([]attendance.Attendance, error) {
	return nil, nil
}
func (m *mockAttRepo) GetStats(studentID string) (*attendance.SummaryResponse, error) {
	return nil, nil
}
func (m *mockAttRepo) GetStatsBatch(ctx context.Context, studentIDs []string) (map[string]*attendance.SummaryResponse, error) {
	return nil, nil
}
func (m *mockAttRepo) CountDistinctDays(studentID string) (int, error) { return 0, nil }
func (m *mockAttRepo) GetAnalytics(ctx context.Context, schoolID string) (*attendance.AnalyticsResponse, error) {
	return nil, nil
}
func (m *mockAttRepo) CreateJustification(j *attendance.Justification) error { return nil }
func (m *mockAttRepo) UpdateJustification(j *attendance.Justification) error { return nil }
func (m *mockAttRepo) ProcessJustificationTx(ctx context.Context, j *attendance.Justification, teacherID string, approve bool) error {
	return nil
}
func (m *mockAttRepo) FindJustificationByID(id string) (*attendance.Justification, error) {
	return nil, nil
}
func (m *mockAttRepo) FindPendingJustifications(classID string) ([]attendance.Justification, error) {
	return nil, nil
}
func (m *mockAttRepo) DeleteJustification(id string) error { return nil }
func (m *mockAttRepo) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	args := m.Called(ctx, teacherID, classID)
	return args.Bool(0), args.Error(1)
}
func (m *mockAttRepo) IsTeacherSubstitute(ctx context.Context, teacherID, classID string, date time.Time, hour int) (bool, error) {
	return false, nil
}
func (m *mockAttRepo) HasOverlappingJustification(ctx context.Context, studentID string, startDate, endDate time.Time) (bool, error) {
	return false, nil
}
func (m *mockAttRepo) GetMonthlyBreakdown(ctx context.Context, studentID, schoolYear string) ([]attendance.MonthlyBreakdownRow, error) {
	return nil, nil
}
func (m *mockAttRepo) FindUnjustifiedByStudent(studentID string) ([]attendance.Attendance, error) {
	return nil, nil
}
func (m *mockAttRepo) JustifyAbsenceByParent(attendanceID string, reason string, notes string) error {
	return nil
}
func (m *mockAttRepo) GetStudentAttendanceStats(studentID string) (*attendance.AttendanceStats, error) {
	return nil, nil
}

type mockScrutinyRepo struct {
	mock.Mock
}

func (m *mockScrutinyRepo) SaveRecord(ctx context.Context, record *scrutiny.ScrutinyRecord) error {
	return nil
}
func (m *mockScrutinyRepo) GetRecord(ctx context.Context, studentID, classID string, semester int) (*scrutiny.ScrutinyRecord, error) {
	return nil, nil
}
func (m *mockScrutinyRepo) ListRecordsByClass(ctx context.Context, classID string, semester int) ([]scrutiny.ScrutinyRecord, error) {
	return nil, nil
}
func (m *mockScrutinyRepo) ValidateClassScrutiny(ctx context.Context, classID string, semester int, validatorID string) error {
	return nil
}
func (m *mockScrutinyRepo) UpdateClassScrutinyStatus(ctx context.Context, classID string, semester int, status string) error {
	return nil
}

func TestFullSchoolYearLifecycle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const schoolID = "school-piccola-01"
	const classID = "class-1a"
	const teacher1ID = "teacher-mario"
	const studentID = "student-luca"
	todayStr := time.Now().Format("2006-01-02")
	schoolIDPtr := schoolID

	setAuth := func(userID, role, sID string) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("user_id", userID)
			c.Set("role", role)
			c.Set("school_id", sID)
			c.Next()
		}
	}

	t.Run("Phase 1: Institution Setup & Multi-Role Enrollment", func(t *testing.T) {
		mockSchoolRepo := new(testhelpers.MockSchoolsRepository)
		mockSchoolRepo.On("Create", mock.Anything, mock.AnythingOfType("*schools.School")).Return(nil)
		schoolSvc := schools.NewService(mockSchoolRepo)
		schoolH := schools.NewHandler(schoolSvc)

		r1 := gin.New()
		g1 := r1.Group("/api/v1")
		g1.Use(setAuth("superadmin-1", "superadmin", ""))
		g1.POST("/schools", schoolH.Create)

		bodySchool := `{"name": "Liceo Piccola Scuola", "code": "LPS001", "city": "Torino", "email": "info@piccolascuola.it"}`
		req := httptest.NewRequest("POST", "/api/v1/schools", bytes.NewBufferString(bodySchool))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		r1.ServeHTTP(res, req)
		assert.Equal(t, http.StatusCreated, res.Code)

		mockUserRepo := new(testhelpers.MockUsersRepository)
		mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*users.User")).Return(nil)
		userSvc := users.NewService(mockUserRepo)
		userH := users.NewHandler(userSvc)

		r2 := gin.New()
		g2 := r2.Group("/api/v1")
		g2.Use(setAuth("superadmin-1", "superadmin", schoolID))
		g2.POST("/users", userH.Create)

		rolesToCreate := []struct {
			email     string
			firstName string
			lastName  string
			role      string
		}{
			{"admin@piccolascuola.it", "Alessandro", "Admin", "admin"},
			{"dirigente@piccolascuola.it", "Roberto", "Preside", "teacher"},
			{"mario.rossi@piccolascuola.it", "Mario", "Rossi", "teacher"},
			{"luca.bianchi@piccolascuola.it", "Luca", "Bianchi", "student"},
			{"giovanni.bianchi@piccolascuola.it", "Giovanni", "Bianchi", "parent"},
		}

		for _, r := range rolesToCreate {
			userPayload := fmt.Sprintf(`{"email":"%s","password":"Password123!","first_name":"%s","last_name":"%s","role":"%s","school_id":"%s"}`,
				r.email, r.firstName, r.lastName, r.role, schoolID)
			reqU := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(userPayload))
			reqU.Header.Set("Content-Type", "application/json")
			resU := httptest.NewRecorder()
			r2.ServeHTTP(resU, reqU)
			assert.Equal(t, http.StatusCreated, resU.Code)
		}

		mockClassRepo := new(testhelpers.MockClassesRepository)
		mockClassRepo.On("Create", mock.Anything, mock.AnythingOfType("*classes.Class")).Return(nil)
		classSvc := classes.NewService(mockClassRepo)
		classH := classes.NewHandler(classSvc)

		r3 := gin.New()
		g3 := r3.Group("/api/v1")
		g3.Use(setAuth("sec-1", "secretary", schoolID))
		g3.POST("/classes", classH.Create)

		bodyClass := fmt.Sprintf(`{"name":"1","section":"A","academic_year":"2025-2026","school_id":"%s"}`, schoolID)
		reqC := httptest.NewRequest("POST", "/api/v1/classes", bytes.NewBufferString(bodyClass))
		reqC.Header.Set("Content-Type", "application/json")
		resC := httptest.NewRecorder()
		r3.ServeHTTP(resC, reqC)
		assert.Equal(t, http.StatusCreated, resC.Code)
	})

	t.Run("Phase 2: Term 1 Daily Operations & Grade Entries (Sept - Jan)", func(t *testing.T) {
		mUser := new(testhelpers.MockUsersRepository)
		mUser.On("GetByID", mock.Anything, mock.Anything).Return(&users.User{
			ID:       teacher1ID,
			Role:     "teacher",
			SchoolID: &schoolIDPtr,
		}, nil).Maybe()

		mAtt := new(mockAttRepo)
		mAtt.On("IsTeacherAssignedToClass", mock.Anything, teacher1ID, classID).Return(true, nil)
		mAtt.On("BatchCreate", mock.Anything).Return(nil)

		attSvc := attendance.NewService(mAtt, mUser, nil, nil)
		attH := attendance.NewHandler(attSvc)

		rAtt := gin.New()
		gAtt := rAtt.Group("/api/v1")
		gAtt.Use(setAuth(teacher1ID, "teacher", schoolID))
		gAtt.POST("/attendance/mark-bulk", attH.MarkBulk)

		bodyMark := fmt.Sprintf(`{"class_id":"%s","date":"%s","hour":1,"statuses":[{"student_id":"%s","status":"Present"}]}`, classID, todayStr, studentID)
		reqAtt := httptest.NewRequest("POST", "/api/v1/attendance/mark-bulk", bytes.NewBufferString(bodyMark))
		reqAtt.Header.Set("Content-Type", "application/json")
		resAtt := httptest.NewRecorder()
		rAtt.ServeHTTP(resAtt, reqAtt)
		assert.Equal(t, http.StatusOK, resAtt.Code)

		mGrade := new(testhelpers.MockGradesRepository)
		createdGrade := &grades.Grade{
			ID:          "grade-q1-1",
			StudentID:   studentID,
			SubjectID:   "subj-math",
			TeacherID:   teacher1ID,
			SchoolID:    schoolID,
			GradeValue:  8.5,
			GradeType:   "numeric",
			Semester:    1,
			IsPublished: true,
		}
		mGrade.On("Create", mock.Anything).Return(nil)
		mGrade.On("FindByID", "grade-q1-1").Return(createdGrade, nil)

		gradeSvc := grades.NewService(mGrade, mUser, nil, nil)
		gradeH := grades.NewHandler(gradeSvc, nil)

		rGrade := gin.New()
		gGrade := rGrade.Group("/api/v1")
		gGrade.Use(setAuth(teacher1ID, "teacher", schoolID))
		gGrade.POST("/grades", gradeH.AddGrade)

		bodyGrade := fmt.Sprintf(`{"student_id":"%s","subject_id":"subj-math","grade_value":8.5,"grade_type":"numeric","semester":1,"date":"2025-10-15"}`, studentID)
		reqG := httptest.NewRequest("POST", "/api/v1/grades", bytes.NewBufferString(bodyGrade))
		reqG.Header.Set("Content-Type", "application/json")
		resG := httptest.NewRecorder()
		rGrade.ServeHTTP(resG, reqG)
		assert.Equal(t, http.StatusCreated, resG.Code)
	})

	t.Run("Phase 3: Intermediate Scrutiny (Q1 Board Meeting)", func(t *testing.T) {
		mUser := new(testhelpers.MockUsersRepository)
		mUser.On("GetByID", mock.Anything, mock.Anything).Return(&users.User{
			ID:       "dirigente-1",
			Role:     "principal",
			SchoolID: &schoolIDPtr,
		}, nil).Maybe()

		mGrade := new(testhelpers.MockGradesRepository)
		mAtt := new(mockAttRepo)
		mClass := new(testhelpers.MockClassesRepository)
		mClass.On("List", mock.Anything, mock.Anything, mock.Anything).Return([]classes.Class{}, nil).Maybe()

		mScrutiny := new(mockScrutinyRepo)
		mScrutiny.On("GetOverview", mock.Anything, schoolID).Return(map[string]interface{}{"status": "ok"}, nil)

		scrutinySvc := scrutiny.NewService(mScrutiny, mGrade, mClass, mUser, mAtt)
		scrutinyH := scrutiny.NewHandler(scrutinySvc)

		rScrutiny := gin.New()
		gScrutiny := rScrutiny.Group("/api/v1")
		gScrutiny.Use(setAuth("dirigente-1", "principal", schoolID))
		gScrutiny.GET("/scrutiny/overview", scrutinyH.GetOverview)

		reqO := httptest.NewRequest("GET", "/api/v1/scrutiny/overview", nil)
		resO := httptest.NewRecorder()
		rScrutiny.ServeHTTP(resO, reqO)
		assert.Equal(t, http.StatusOK, resO.Code)
	})

	t.Run("Phase 4: Term 2 & Final Scrutiny (Q2 Board Meeting)", func(t *testing.T) {
		mUser := new(testhelpers.MockUsersRepository)
		mUser.On("GetByID", mock.Anything, mock.Anything).Return(&users.User{
			ID:       "dirigente-1",
			Role:     "principal",
			SchoolID: &schoolIDPtr,
		}, nil).Maybe()

		mGrade := new(testhelpers.MockGradesRepository)
		mAtt := new(mockAttRepo)
		mClass := new(testhelpers.MockClassesRepository)
		mClass.On("List", mock.Anything, mock.Anything, mock.Anything).Return([]classes.Class{}, nil).Maybe()

		mScrutiny := new(mockScrutinyRepo)
		mScrutiny.On("ExportAll", mock.Anything, schoolID).Return([]byte("pdf-data"), nil)

		scrutinySvc := scrutiny.NewService(mScrutiny, mGrade, mClass, mUser, mAtt)
		scrutinyH := scrutiny.NewHandler(scrutinySvc)

		rFinal := gin.New()
		gFinal := rFinal.Group("/api/v1")
		gFinal.Use(setAuth("dirigente-1", "principal", schoolID))
		gFinal.GET("/scrutiny-final/export", scrutinyH.ExportAll)

		reqExp := httptest.NewRequest("GET", "/api/v1/scrutiny-final/export", nil)
		resExp := httptest.NewRecorder()
		rFinal.ServeHTTP(resExp, reqExp)
		assert.Equal(t, http.StatusOK, resExp.Code)
	})
}
