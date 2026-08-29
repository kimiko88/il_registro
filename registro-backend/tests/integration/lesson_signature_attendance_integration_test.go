package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/attendance"
	"registro-backend/internal/lessons"
	"registro-backend/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockLessonsRepoForSignature struct {
	lessons   map[string]*lessons.Lesson
	homeworks map[string]*lessons.Homework
}

func newMockLessonsRepoForSignature() *mockLessonsRepoForSignature {
	return &mockLessonsRepoForSignature{
		lessons:   make(map[string]*lessons.Lesson),
		homeworks: make(map[string]*lessons.Homework),
	}
}

func (m *mockLessonsRepoForSignature) CreateLesson(l *lessons.Lesson) error {
	if l.ID == "" {
		l.ID = "lesson-" + time.Now().Format("150405.000000")
	}
	l.CreatedAt = time.Now()
	l.UpdatedAt = time.Now()
	m.lessons[l.ID] = l
	return nil
}

func (m *mockLessonsRepoForSignature) GetLessonByID(id string) (*lessons.Lesson, error) {
	l, ok := m.lessons[id]
	if !ok {
		return nil, assert.AnError
	}
	return l, nil
}

func (m *mockLessonsRepoForSignature) UpdateLesson(id string, req lessons.UpdateLessonRequest) (*lessons.Lesson, error) {
	l, ok := m.lessons[id]
	if !ok {
		return nil, assert.AnError
	}
	if req.Topic != "" {
		l.Topic = req.Topic
	}
	return l, nil
}

func (m *mockLessonsRepoForSignature) DeleteLesson(id string) error {
	delete(m.lessons, id)
	return nil
}

func (m *mockLessonsRepoForSignature) GetLessonsByClass(classID string, date string) ([]lessons.Lesson, error) {
	var res []lessons.Lesson
	for _, l := range m.lessons {
		if l.ClassID == classID {
			res = append(res, *l)
		}
	}
	return res, nil
}

func (m *mockLessonsRepoForSignature) GetLessonsByClassAndSubject(classID, subjectID string, date string) ([]lessons.Lesson, error) {
	var res []lessons.Lesson
	for _, l := range m.lessons {
		if l.ClassID == classID && (subjectID == "" || l.SubjectID == subjectID) {
			res = append(res, *l)
		}
	}
	return res, nil
}

func (m *mockLessonsRepoForSignature) GetLessonsByGroup(groupID string, date string) ([]lessons.Lesson, error) {
	return []lessons.Lesson{}, nil
}

func (m *mockLessonsRepoForSignature) GetLessonsByTeacher(teacherID string, fromDate, toDate string) ([]lessons.Lesson, error) {
	var res []lessons.Lesson
	for _, l := range m.lessons {
		if l.TeacherID == teacherID {
			res = append(res, *l)
		}
	}
	return res, nil
}

func (m *mockLessonsRepoForSignature) CreateHomework(h *lessons.Homework) error {
	if h.ID == "" {
		h.ID = "hw-" + time.Now().Format("150405.000000")
	}
	m.homeworks[h.ID] = h
	return nil
}

func (m *mockLessonsRepoForSignature) GetHomeworkByID(id string) (*lessons.Homework, error) {
	h, ok := m.homeworks[id]
	if !ok {
		return nil, assert.AnError
	}
	return h, nil
}

func (m *mockLessonsRepoForSignature) UpdateHomework(id string, req lessons.UpdateHomeworkRequest) (*lessons.Homework, error) {
	h, ok := m.homeworks[id]
	if !ok {
		return nil, assert.AnError
	}
	if req.Description != "" {
		h.Description = req.Description
	}
	return h, nil
}

func (m *mockLessonsRepoForSignature) DeleteHomework(id string) error {
	delete(m.homeworks, id)
	return nil
}

func (m *mockLessonsRepoForSignature) GetHomeworkByClass(classID string, fromDate ...string) ([]lessons.Homework, error) {
	var res []lessons.Homework
	for _, h := range m.homeworks {
		if h.ClassID == classID {
			res = append(res, *h)
		}
	}
	return res, nil
}

func (m *mockLessonsRepoForSignature) IsTeacherAssignedToClass(teacherID, classID string) (bool, error) {
	return true, nil
}

func (m *mockLessonsRepoForSignature) HasApprovedSubstitution(teacherID, classID, date string, hour int) (bool, error) {
	return true, nil
}

type mockAttendanceRepoForLesson struct {
	records map[string]*attendance.Attendance
}

func (m *mockAttendanceRepoForLesson) Create(att *attendance.Attendance) error {
	if att.ID == "" {
		att.ID = "att-" + time.Now().Format("150405.000000")
	}
	m.records[att.ID] = att
	return nil
}

func (m *mockAttendanceRepoForLesson) BatchCreate(atts []*attendance.Attendance) error {
	for _, a := range atts {
		if a.ID == "" {
			a.ID = "att-" + time.Now().Format("150405.000000")
		}
		m.records[a.ID] = a
	}
	return nil
}

func (m *mockAttendanceRepoForLesson) Update(att *attendance.Attendance) error { return nil }
func (m *mockAttendanceRepoForLesson) DeleteByClassDateHour(schoolID, classID string, date time.Time, hour int) error {
	return nil
}
func (m *mockAttendanceRepoForLesson) FindByID(id string) (*attendance.Attendance, error) {
	return &attendance.Attendance{ID: id}, nil
}
func (m *mockAttendanceRepoForLesson) FindByClassAndDate(classID string, date time.Time) ([]attendance.Attendance, error) {
	var res []attendance.Attendance
	for _, a := range m.records {
		if a.ClassID == classID {
			res = append(res, *a)
		}
	}
	return res, nil
}
func (m *mockAttendanceRepoForLesson) FindByStudent(studentID string, startDate, endDate time.Time) ([]attendance.Attendance, error) {
	return []attendance.Attendance{}, nil
}
func (m *mockAttendanceRepoForLesson) GetStats(studentID string) (*attendance.SummaryResponse, error) {
	return &attendance.SummaryResponse{}, nil
}
func (m *mockAttendanceRepoForLesson) GetStatsBatch(ctx context.Context, studentIDs []string) (map[string]*attendance.SummaryResponse, error) {
	return map[string]*attendance.SummaryResponse{}, nil
}
func (m *mockAttendanceRepoForLesson) CountDistinctDays(studentID string) (int, error) {
	return 0, nil
}
func (m *mockAttendanceRepoForLesson) GetAnalytics(ctx context.Context, schoolID string) (*attendance.AnalyticsResponse, error) {
	return &attendance.AnalyticsResponse{}, nil
}
func (m *mockAttendanceRepoForLesson) CreateJustification(j *attendance.Justification) error {
	return nil
}
func (m *mockAttendanceRepoForLesson) UpdateJustification(j *attendance.Justification) error {
	return nil
}
func (m *mockAttendanceRepoForLesson) ProcessJustificationTx(ctx context.Context, j *attendance.Justification, teacherID string, approve bool) error {
	return nil
}
func (m *mockAttendanceRepoForLesson) FindJustificationByID(id string) (*attendance.Justification, error) {
	return &attendance.Justification{ID: id}, nil
}
func (m *mockAttendanceRepoForLesson) FindPendingJustifications(classID, schoolID string) ([]attendance.Justification, error) {
	return []attendance.Justification{}, nil
}
func (m *mockAttendanceRepoForLesson) FindPendingJustificationsForTeacher(ctx context.Context, teacherID, schoolID string) ([]attendance.Justification, error) {
	return []attendance.Justification{}, nil
}
func (m *mockAttendanceRepoForLesson) DeleteJustification(id string) error        { return nil }
func (m *mockAttendanceRepoForLesson) DeletePendingJustification(id string) error { return nil }
func (m *mockAttendanceRepoForLesson) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	return true, nil
}
func (m *mockAttendanceRepoForLesson) IsTeacherSubstitute(ctx context.Context, teacherID, classID string, date time.Time, hour int) (bool, error) {
	return false, nil
}
func (m *mockAttendanceRepoForLesson) JustifyAbsenceByParent(attendanceID string, parentID string, reason string, documentURL string) error {
	return nil
}
func (m *mockAttendanceRepoForLesson) IsClassInSchool(ctx context.Context, classID, schoolID string) (bool, error) {
	return true, nil
}
func (m *mockAttendanceRepoForLesson) IsStudentInClass(ctx context.Context, studentID, classID string) (bool, error) {
	return true, nil
}
func (m *mockAttendanceRepoForLesson) AreStudentsInClass(ctx context.Context, studentIDs []string, classID string) (map[string]bool, error) {
	res := make(map[string]bool)
	for _, id := range studentIDs {
		res[id] = true
	}
	return res, nil
}
func (m *mockAttendanceRepoForLesson) HasOverlappingJustification(ctx context.Context, studentID string, startDate, endDate time.Time) (bool, error) {
	return false, nil
}
func (m *mockAttendanceRepoForLesson) FindUnjustifiedByStudent(studentID string) ([]attendance.Attendance, error) {
	return []attendance.Attendance{}, nil
}
func (m *mockAttendanceRepoForLesson) GetMonthlyBreakdown(ctx context.Context, studentID string, schoolID string) ([]attendance.MonthlyBreakdownRow, error) {
	return []attendance.MonthlyBreakdownRow{}, nil
}
func (m *mockAttendanceRepoForLesson) GetStudentAttendanceStats(studentID string) (*attendance.AttendanceStats, error) {
	return &attendance.AttendanceStats{}, nil
}

type mockUserRepoForLessonAttendance struct {
	users.Repository
}

func (m *mockUserRepoForLessonAttendance) GetByID(ctx context.Context, id string) (*users.User, error) {
	schoolID := "school-1"
	return &users.User{
		ID:       id,
		SchoolID: &schoolID,
		Role:     "teacher",
	}, nil
}

func (m *mockUserRepoForLessonAttendance) IsTeacherOf(teacherID, classID, subjectID string) (bool, error) {
	return true, nil
}

func setupLessonSignatureRouter(lRepo lessons.Repository, aRepo attendance.Repository, uRepo users.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	lSvc := lessons.NewService(lRepo)
	lHandler := lessons.NewHandler(lSvc)

	aSvc := attendance.NewService(aRepo, uRepo, nil, nil)
	aHandler := attendance.NewHandler(aSvc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		c.Set("role", "teacher")
		c.Set("user_id", "teacher-math-1")
		c.Set("school_id", "school-1")
		c.Next()
	})
	lHandler.RegisterRoutes(api)
	aHandler.RegisterRoutes(api)
	return r
}

func TestIntegration_LessonSignature_And_Attendance_Workflow(t *testing.T) {
	lRepo := newMockLessonsRepoForSignature()
	aRepo := &mockAttendanceRepoForLesson{records: make(map[string]*attendance.Attendance)}
	uRepo := &mockUserRepoForLessonAttendance{}
	r := setupLessonSignatureRouter(lRepo, aRepo, uRepo)

	today := time.Now().Format("2006-01-02")

	// 1. Teacher Signs Current Hour Lesson with 1-Click Signature & Topic
	lessReq := lessons.CreateLessonRequest{
		ClassID:      "class-2A",
		SubjectID:    "sub-math",
		Date:         today,
		Hour:         2,
		Duration:     1,
		Type:         "Frontale",
		Topic:        "Scomposizione di polinomi: Raccoglimento a fattore comune e totale",
		IsCoTeaching: false,
	}
	bodyLess, _ := json.Marshal(lessReq)
	reqLess, _ := http.NewRequest("POST", "/api/v1/lessons", bytes.NewReader(bodyLess))
	reqLess.Header.Set("Content-Type", "application/json")
	wLess := httptest.NewRecorder()
	r.ServeHTTP(wLess, reqLess)
	require.Equal(t, http.StatusCreated, wLess.Code)

	var createdLesson lessons.LessonResponse
	err := json.Unmarshal(wLess.Body.Bytes(), &createdLesson)
	require.NoError(t, err)
	assert.NotEmpty(t, createdLesson.ID)
	assert.Equal(t, 2, createdLesson.Hour)
	assert.Equal(t, "Scomposizione di polinomi: Raccoglimento a fattore comune e totale", createdLesson.Topic)

	// 2. Assign Homework linked for next lesson
	hwReq := lessons.CreateHomeworkRequest{
		ClassID:     "class-2A",
		SubjectID:   "sub-math",
		Description: "Esercizi pag. 142 n. 12, 14, 18 e studio teoria",
		DueDate:     time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
	}
	bodyHw, _ := json.Marshal(hwReq)
	reqHw, _ := http.NewRequest("POST", "/api/v1/homeworks", bytes.NewReader(bodyHw))
	reqHw.Header.Set("Content-Type", "application/json")
	wHw := httptest.NewRecorder()
	r.ServeHTTP(wHw, reqHw)
	require.Equal(t, http.StatusCreated, wHw.Code)

	// 3. Mark Hourly Attendance / Absences / Delays for the Class
	bulkAttPayload := attendance.BulkAttendanceRequest{
		ClassID: "class-2A",
		Date:    today,
		Hour:    2,
		Statuses: []attendance.StudentStatusRequest{
			{StudentID: "student-1", Status: attendance.StatusPresent},
			{StudentID: "student-2", Status: attendance.StatusAbsent, Notes: "Assente giustificazione richiesta"},
			{StudentID: "student-3", Status: attendance.StatusLate, EntryTime: "09:15", Notes: "Entrata in 2ª ora ore 09:15"},
		},
	}
	bodyAtt, _ := json.Marshal(bulkAttPayload)
	reqAtt, _ := http.NewRequest("POST", "/api/v1/attendance/mark-bulk", bytes.NewReader(bodyAtt))
	reqAtt.Header.Set("Content-Type", "application/json")
	wAtt := httptest.NewRecorder()
	r.ServeHTTP(wAtt, reqAtt)
	require.Equal(t, http.StatusOK, wAtt.Code)

	// 4. Verify recorded attendance for class and date
	reqGetAtt, _ := http.NewRequest("GET", "/api/v1/attendance/class/class-2A?date="+today, nil)
	wGetAtt := httptest.NewRecorder()
	r.ServeHTTP(wGetAtt, reqGetAtt)
	require.Equal(t, http.StatusOK, wGetAtt.Code)
}
