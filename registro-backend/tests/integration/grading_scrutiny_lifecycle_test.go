package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/classes"
	"registro-backend/internal/grades"
	"registro-backend/internal/scrutiny"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_Grading_Scrutiny_Lifecycle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockGradesRepo := &mockGradesRepoForScrutiny{}
	mockScrutinyRepo := &mockScrutinyRepoForGrading{}
	mockUserRepo := &mockUserRepoForComms{}
	mockAbsenceRepo := &mockAbsenceRepoForLifecycle{}

	gradesSvc := grades.NewService(mockGradesRepo, mockUserRepo, nil, nil)
	gradesHandler := grades.NewHandler(gradesSvc, nil)

	mockClassesRepo := &mockClassesRepoForScrutiny{}
	scrutinySvc := scrutiny.NewService(mockScrutinyRepo, mockGradesRepo, mockClassesRepo, mockUserRepo, mockAbsenceRepo)
	scrutinyHandler := scrutiny.NewHandler(scrutinySvc, nil)

	r := gin.New()
	r.GET("/grades/classes/:classID", func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
		gradesHandler.GetClassGrades(c)
	})

	r.GET("/scrutiny/overview", func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
		scrutinyHandler.GetOverview(c)
	})

	// 1. Get class grades
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/grades/classes/class-1", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// 2. Get scrutiny overview
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/scrutiny/overview?class_id=class-1&semester=1", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

type mockGradesRepoForScrutiny struct{}

func (m *mockGradesRepoForScrutiny) Create(ctx context.Context, g *grades.Grade) error { return nil }
func (m *mockGradesRepoForScrutiny) BatchCreate(ctx context.Context, gradesList []*grades.Grade) error {
	return nil
}
func (m *mockGradesRepoForScrutiny) Update(ctx context.Context, g *grades.Grade, history *grades.GradeHistory) error {
	return nil
}
func (m *mockGradesRepoForScrutiny) Delete(ctx context.Context, id string, deletedBy string) error {
	return nil
}
func (m *mockGradesRepoForScrutiny) SoftDelete(ctx context.Context, id, modifiedBy string) error {
	return nil
}
func (m *mockGradesRepoForScrutiny) FindByID(ctx context.Context, id string) (*grades.Grade, error) {
	return &grades.Grade{ID: id, GradeValue: 8}, nil
}
func (m *mockGradesRepoForScrutiny) FindByStudent(ctx context.Context, studentID string) ([]grades.Grade, error) {
	return []grades.Grade{}, nil
}
func (m *mockGradesRepoForScrutiny) FindByClassAndSubject(ctx context.Context, classID, subjectID string, semester int) ([]grades.Grade, error) {
	return []grades.Grade{}, nil
}
func (m *mockGradesRepoForScrutiny) GetClassAverage(classID string) (float64, error) { return 7.5, nil }
func (m *mockGradesRepoForScrutiny) GetSubjectAverage(studentID, subjectID string) (float64, error) {
	return 8.0, nil
}
func (m *mockGradesRepoForScrutiny) IsTeacherAssignedToClassSubject(teacherID, classID, subjectID string) (bool, error) {
	return true, nil
}
func (m *mockGradesRepoForScrutiny) IsTeacherAssignedToClass(teacherID, classID string) (bool, error) {
	return true, nil
}
func (m *mockGradesRepoForScrutiny) IsClassCoordinator(teacherID, classID string) (bool, error) {
	return true, nil
}
func (m *mockGradesRepoForScrutiny) IsStudentInClass(studentID, classID string) (bool, error) {
	return true, nil
}
func (m *mockGradesRepoForScrutiny) IsClassInSchool(classID, schoolID string) (bool, error) {
	return true, nil
}
func (m *mockGradesRepoForScrutiny) UpsertWeightConfig(ctx context.Context, cfg *grades.GradeWeightConfig) (*grades.GradeWeightConfig, error) {
	return cfg, nil
}
func (m *mockGradesRepoForScrutiny) GetWeightConfig(schoolID, subjectID, classID string) (*grades.GradeWeightConfig, error) {
	return nil, nil
}
func (m *mockGradesRepoForScrutiny) DeleteWeightConfig(ctx context.Context, id string) error {
	return nil
}
func (m *mockGradesRepoForScrutiny) CreateTest(ctx context.Context, testData *grades.ClassTest) error {
	return nil
}
func (m *mockGradesRepoForScrutiny) DeleteTest(ctx context.Context, id string) error { return nil }
func (m *mockGradesRepoForScrutiny) FindByClass(ctx context.Context, classID string, semester int) ([]grades.Grade, error) {
	return []grades.Grade{}, nil
}
func (m *mockGradesRepoForScrutiny) FindBySubject(ctx context.Context, subjectID string, semester int) ([]grades.Grade, error) {
	return []grades.Grade{}, nil
}
func (m *mockGradesRepoForScrutiny) FindByTeacher(ctx context.Context, teacherID string) ([]grades.Grade, error) {
	return []grades.Grade{}, nil
}
func (m *mockGradesRepoForScrutiny) FindEnrolledSubjects(ctx context.Context, studentID string, semester int) ([]string, error) {
	return []string{}, nil
}
func (m *mockGradesRepoForScrutiny) FindGradesByTestID(ctx context.Context, testID string) ([]grades.Grade, error) {
	return []grades.Grade{}, nil
}
func (m *mockGradesRepoForScrutiny) FindTestByID(ctx context.Context, id string) (*grades.ClassTest, error) {
	return &grades.ClassTest{ID: id}, nil
}
func (m *mockGradesRepoForScrutiny) FindTestsByClassAndSubject(ctx context.Context, classID string, subjectID string) ([]grades.ClassTest, error) {
	return []grades.ClassTest{}, nil
}
func (m *mockGradesRepoForScrutiny) FindUpcomingTestsByClass(ctx context.Context, classID string) ([]grades.ClassTest, error) {
	return []grades.ClassTest{}, nil
}
func (m *mockGradesRepoForScrutiny) FindWithFilter(ctx context.Context, filter grades.GradeFilter) ([]grades.Grade, error) {
	return []grades.Grade{}, nil
}
func (m *mockGradesRepoForScrutiny) FindWithFilterPaginated(ctx context.Context, filter grades.GradeFilter) ([]grades.Grade, int, error) {
	return []grades.Grade{}, 0, nil
}
func (m *mockGradesRepoForScrutiny) GetHistory(ctx context.Context, gradeID string) ([]grades.GradeHistory, error) {
	return []grades.GradeHistory{}, nil
}
func (m *mockGradesRepoForScrutiny) GetWeightConfigs(ctx context.Context, schoolID, subjectID, classID string) ([]grades.GradeWeightConfig, error) {
	return []grades.GradeWeightConfig{}, nil
}
func (m *mockGradesRepoForScrutiny) UpdateTest(ctx context.Context, test *grades.ClassTest) error {
	return nil
}

func (m *mockGradesRepoForScrutiny) GetStudentClassAndSchoolInfo(ctx context.Context, studentID string) (string, string, string, string, error) {
	return "", "", "", "", nil
}
func (m *mockGradesRepoForScrutiny) GetTeacherNamesByClass(ctx context.Context, classID string) (map[string]string, error) {
	return nil, nil
}
func (m *mockGradesRepoForScrutiny) GetSubjectNamesMap(ctx context.Context, schoolID string) (map[string]string, error) {
	return nil, nil
}
func (m *mockGradesRepoForScrutiny) GetScrutinyRecordSummary(ctx context.Context, studentID string, semester int) (float64, float64, bool, error) {
	return 0, 0, false, nil
}
func (m *mockGradesRepoForScrutiny) GetStudentAbsenceCountForPeriod(ctx context.Context, studentID, startD, endD string) (int, error) {
	return 0, nil
}
func (m *mockGradesRepoForScrutiny) GetClassSubjectAverage(ctx context.Context, classID, subjectID string, semester int, studentID string) (float64, error) {
	return -1, nil
}
func (m *mockGradesRepoForScrutiny) CheckClassAccessPermission(ctx context.Context, actorID, actorRole, classID string) (bool, error) {
	return true, nil
}

type mockScrutinyRepoForGrading struct{}

func (m *mockScrutinyRepoForGrading) SaveRecord(ctx context.Context, record *scrutiny.ScrutinyRecord) error {
	return nil
}
func (m *mockScrutinyRepoForGrading) GetRecord(ctx context.Context, studentID, classID string, semester int) (*scrutiny.ScrutinyRecord, error) {
	return &scrutiny.ScrutinyRecord{ID: "record-1", StudentID: studentID, ClassID: classID, Semester: semester}, nil
}
func (m *mockScrutinyRepoForGrading) ListRecordsByClass(ctx context.Context, classID string, semester int) ([]scrutiny.ScrutinyRecord, error) {
	return []scrutiny.ScrutinyRecord{}, nil
}
func (m *mockScrutinyRepoForGrading) ValidateClassScrutiny(ctx context.Context, classID string, semester int, validatorID string) error {
	return nil
}
func (m *mockScrutinyRepoForGrading) UpdateClassScrutinyStatus(ctx context.Context, classID string, semester int, status string) error {
	return nil
}
func (m *mockScrutinyRepoForGrading) SaveDeficiency(ctx context.Context, def *scrutiny.StudentDeficiency) error {
	return nil
}
func (m *mockScrutinyRepoForGrading) GetDeficienciesByStudent(ctx context.Context, studentID string) ([]scrutiny.StudentDeficiency, error) {
	return []scrutiny.StudentDeficiency{}, nil
}
func (m *mockScrutinyRepoForGrading) GetDeficienciesByClass(ctx context.Context, classID string, semester int) ([]scrutiny.StudentDeficiency, error) {
	return []scrutiny.StudentDeficiency{}, nil
}
func (m *mockScrutinyRepoForGrading) SaveDeferredScrutiny(ctx context.Context, req *scrutiny.SaveDeferredScrutinyRequest) error {
	return nil
}

type mockClassesRepoForScrutiny struct{}

func (m *mockClassesRepoForScrutiny) Create(ctx context.Context, class *classes.Class) error {
	return nil
}
func (m *mockClassesRepoForScrutiny) List(ctx context.Context, schoolID string, academicYear string) ([]classes.Class, error) {
	return []classes.Class{}, nil
}
func (m *mockClassesRepoForScrutiny) Get(ctx context.Context, id string) (*classes.Class, error) {
	return &classes.Class{ID: id}, nil
}
func (m *mockClassesRepoForScrutiny) Update(ctx context.Context, class *classes.Class) error {
	return nil
}
func (m *mockClassesRepoForScrutiny) Delete(ctx context.Context, id string) error { return nil }
func (m *mockClassesRepoForScrutiny) ListByTeacher(ctx context.Context, teacherUserID string, schoolYear string) ([]classes.Class, error) {
	return []classes.Class{}, nil
}
func (m *mockClassesRepoForScrutiny) AssignSubject(ctx context.Context, classID string, subjectID string, teacherID *string, hours float64) error {
	return nil
}
func (m *mockClassesRepoForScrutiny) UnassignSubject(ctx context.Context, assignmentID string) error {
	return nil
}
func (m *mockClassesRepoForScrutiny) GetClassSubjects(ctx context.Context, classID string) ([]classes.ClassSubject, error) {
	return []classes.ClassSubject{}, nil
}
func (m *mockClassesRepoForScrutiny) GetClassGuardians(ctx context.Context, classID string) ([]classes.GuardianInfo, error) {
	return []classes.GuardianInfo{}, nil
}
func (m *mockClassesRepoForScrutiny) GetLessonTopics(ctx context.Context, classID string) ([]classes.LessonTopic, error) {
	return []classes.LessonTopic{}, nil
}
func (m *mockClassesRepoForScrutiny) GetDisciplinaryNotes(ctx context.Context, classID string) ([]classes.DisciplinaryNoteReport, error) {
	return []classes.DisciplinaryNoteReport{}, nil
}
func (m *mockClassesRepoForScrutiny) BulkMigrateStudents(ctx context.Context, migrations []classes.StudentMigrationItem) error {
	return nil
}
