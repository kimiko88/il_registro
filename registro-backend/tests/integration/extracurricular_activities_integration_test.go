package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/extracurricular"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Mock Extracurricular Repository ───────────────────────────────────────

type mockActivitiesExtraRepo struct {
	courses     map[string]extracurricular.Course
	enrollments map[string][]extracurricular.Enrollment
	attendance  map[string][]extracurricular.AttendanceRecord
}

func newMockActivitiesExtraRepo() *mockActivitiesExtraRepo {
	return &mockActivitiesExtraRepo{
		courses:     make(map[string]extracurricular.Course),
		enrollments: make(map[string][]extracurricular.Enrollment),
		attendance:  make(map[string][]extracurricular.AttendanceRecord),
	}
}

func (m *mockActivitiesExtraRepo) CreateCourse(_ context.Context, c *extracurricular.Course) error {
	if c.ID == "" {
		c.ID = "extra-" + time.Now().Format("150405.000")
	}
	c.CreatedAt = time.Now()
	m.courses[c.ID] = *c
	return nil
}

func (m *mockActivitiesExtraRepo) GetCourseByID(_ context.Context, id string) (*extracurricular.Course, error) {
	c, ok := m.courses[id]
	if !ok {
		return nil, errors.New("course not found")
	}
	return &c, nil
}

func (m *mockActivitiesExtraRepo) ListCourses(_ context.Context, schoolID, _ string) ([]*extracurricular.Course, error) {
	var list []*extracurricular.Course
	for _, c := range m.courses {
		if schoolID != "" && c.SchoolID != schoolID {
			continue
		}
		list = append(list, &c)
	}
	return list, nil
}

func (m *mockActivitiesExtraRepo) EnrollStudent(_ context.Context, courseID, studentID string) error {
	c, ok := m.courses[courseID]
	if !ok {
		return errors.New("course not found")
	}
	enrolled := m.enrollments[courseID]
	if c.MaxParticipants > 0 && len(enrolled) >= c.MaxParticipants {
		return errors.New("course is full")
	}
	m.enrollments[courseID] = append(enrolled, extracurricular.Enrollment{
		ID:         "enr-" + studentID,
		CourseID:   courseID,
		StudentID:  studentID,
		EnrolledAt: time.Now(),
	})
	return nil
}

func (m *mockActivitiesExtraRepo) ListEnrollments(_ context.Context, courseID string) ([]*extracurricular.Enrollment, error) {
	var list []*extracurricular.Enrollment
	for _, e := range m.enrollments[courseID] {
		list = append(list, &e)
	}
	return list, nil
}

func (m *mockActivitiesExtraRepo) MarkAttendance(_ context.Context, att *extracurricular.AttendanceRecord) error {
	if att.ID == "" {
		att.ID = "att-" + time.Now().Format("150405.000")
	}
	m.attendance[att.CourseID] = append(m.attendance[att.CourseID], *att)
	return nil
}

func (m *mockActivitiesExtraRepo) ListAttendance(_ context.Context, courseID string, _ time.Time) ([]*extracurricular.AttendanceRecord, error) {
	var list []*extracurricular.AttendanceRecord
	for _, a := range m.attendance[courseID] {
		list = append(list, &a)
	}
	return list, nil
}

// ─── Setup Router Helper ───────────────────────────────────────────────────

func setupExtraRouter(repo *mockActivitiesExtraRepo, role, userID, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := extracurricular.NewService(repo)
	h := extracurricular.NewHandler(svc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
			c.Set("role", role)
			c.Set("school_id", schoolID)
		}
		c.Next()
	})
	h.RegisterRoutes(api)
	return r
}

// ─── Tests ─────────────────────────────────────────────────────────────────

// EXT01 — Teacher creates an extracurricular course
func TestExtracurricular_CreateCourse_TeacherSuccess(t *testing.T) {
	repo := newMockActivitiesExtraRepo()
	r := setupExtraRouter(repo, "teacher", "teacher-1", "school-1")

	body, _ := json.Marshal(extracurricular.CreateCourseRequest{
		Title:           "Laboratorio di Robotica & Coding",
		Description:     "Corso pratico di programmazione Arduino e robotica educativa",
		StartDate:       "2026-10-01",
		EndDate:         "2026-12-15",
		MaxParticipants: 20,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/extracurricular/courses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var course extracurricular.Course
	err := json.Unmarshal(w.Body.Bytes(), &course)
	require.NoError(t, err)
	assert.Equal(t, "Laboratorio di Robotica & Coding", course.Title)
	assert.Equal(t, 20, course.MaxParticipants)
}

// EXT02 — Student is forbidden from creating extracurricular courses
func TestExtracurricular_CreateCourse_ForbiddenForStudent(t *testing.T) {
	repo := newMockActivitiesExtraRepo()
	r := setupExtraRouter(repo, "student", "student-1", "school-1")

	body, _ := json.Marshal(extracurricular.CreateCourseRequest{
		Title:     "Torneo Scacchi",
		StartDate: "2026-10-01",
		EndDate:   "2026-11-01",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/extracurricular/courses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// EXT03 — Student enrolls in an available course
func TestExtracurricular_EnrollStudent_Success(t *testing.T) {
	repo := newMockActivitiesExtraRepo()
	_ = repo.CreateCourse(context.Background(), &extracurricular.Course{
		ID:              "course-robotics",
		SchoolID:        "school-1",
		Title:           "Robotica",
		MaxParticipants: 15,
	})

	r := setupExtraRouter(repo, "student", "student-coder-1", "school-1")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/extracurricular/courses/course-robotics/enroll", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	enrollments, _ := repo.ListEnrollments(context.Background(), "course-robotics")
	assert.Len(t, enrollments, 1)
	assert.Equal(t, "student-coder-1", enrollments[0].StudentID)
}

// EXT04 — Teacher records attendance for extracurricular course session
func TestExtracurricular_MarkAttendance_Success(t *testing.T) {
	repo := newMockActivitiesExtraRepo()
	_ = repo.CreateCourse(context.Background(), &extracurricular.Course{
		ID:        "course-robotics",
		SchoolID:  "school-1",
		TeacherID: "teacher-1",
	})

	r := setupExtraRouter(repo, "teacher", "teacher-1", "school-1")

	body, _ := json.Marshal(extracurricular.MarkAttendanceRequest{
		CourseID:  "course-robotics",
		StudentID: "student-1",
		Date:      "2026-10-05",
		Status:    "present",
		Hours:     2.0,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/extracurricular/attendance", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	attList, _ := repo.ListAttendance(context.Background(), "course-robotics", time.Now())
	require.Len(t, attList, 1)
	assert.Equal(t, "present", attList[0].Status)
	assert.Equal(t, 2.0, attList[0].Hours)
}
