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

	"registro-backend/internal/recovery"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Mock Recovery Repository ──────────────────────────────────────────────

type mockRecoveryRepo struct {
	courses map[string]recovery.RecoveryCourse
	tests   map[string]recovery.RecoveryTest
}

func newMockRecoveryRepo() *mockRecoveryRepo {
	return &mockRecoveryRepo{
		courses: make(map[string]recovery.RecoveryCourse),
		tests:   make(map[string]recovery.RecoveryTest),
	}
}

func (m *mockRecoveryRepo) CreateCourse(_ context.Context, c *recovery.RecoveryCourse, sessions []recovery.RecoveryCourseSession, studentIDs []string) error {
	if c.ID == "" {
		c.ID = "rec-course-" + time.Now().Format("150405.000")
	}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	c.Status = "scheduled"
	c.Sessions = sessions
	for _, sid := range studentIDs {
		c.Students = append(c.Students, recovery.RecoveryCourseStudent{
			ID:        "rcs-" + sid,
			CourseID:  c.ID,
			StudentID: sid,
		})
	}
	m.courses[c.ID] = *c
	return nil
}

func (m *mockRecoveryRepo) GetCourseByID(_ context.Context, id string) (*recovery.RecoveryCourse, error) {
	c, ok := m.courses[id]
	if !ok {
		return nil, errors.New("course not found")
	}
	return &c, nil
}

func (m *mockRecoveryRepo) ListCourses(_ context.Context, schoolID, academicYear, _ string) ([]recovery.RecoveryCourse, error) {
	var list []recovery.RecoveryCourse
	for _, c := range m.courses {
		if schoolID != "" && c.SchoolID != schoolID {
			continue
		}
		if academicYear != "" && c.AcademicYear != academicYear {
			continue
		}
		list = append(list, c)
	}
	return list, nil
}

func (m *mockRecoveryRepo) UpdateCourseStatus(_ context.Context, id, status string) error {
	c, ok := m.courses[id]
	if !ok {
		return errors.New("course not found")
	}
	c.Status = status
	c.UpdatedAt = time.Now()
	m.courses[id] = c
	return nil
}

func (m *mockRecoveryRepo) UpdateStudentAttendance(_ context.Context, courseID, studentID string, hours float64, notes string) error {
	c, ok := m.courses[courseID]
	if !ok {
		return errors.New("course not found")
	}
	for i, s := range c.Students {
		if s.StudentID == studentID {
			c.Students[i].AttendanceHours = hours
			c.Students[i].Notes = notes
			m.courses[courseID] = c
			return nil
		}
	}
	return errors.New("student not found in course")
}

func (m *mockRecoveryRepo) RecordRecoveryTest(_ context.Context, t *recovery.RecoveryTest) error {
	if t.ID == "" {
		t.ID = "rec-test-" + time.Now().Format("150405.000")
	}
	m.tests[t.ID] = *t
	return nil
}

func (m *mockRecoveryRepo) ListRecoveryTests(_ context.Context, schoolID, classID, studentID string) ([]recovery.RecoveryTest, error) {
	var list []recovery.RecoveryTest
	for _, t := range m.tests {
		if schoolID != "" && t.SchoolID != schoolID {
			continue
		}
		if classID != "" && t.ClassID != classID {
			continue
		}
		if studentID != "" && t.StudentID != studentID {
			continue
		}
		list = append(list, t)
	}
	return list, nil
}

// ─── Setup Router Helper ───────────────────────────────────────────────────

func setupRecoveryRouter(repo *mockRecoveryRepo, role, userID, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := recovery.NewService(repo)
	h := recovery.NewHandler(svc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
			c.Set("role", role)
			c.Set("school_id", schoolID)
			c.Set("teacher_id", userID)
		}
		c.Next()
	})
	h.RegisterRoutes(api)
	return r
}

// ─── Tests ─────────────────────────────────────────────────────────────────

// REC01 — Teacher/Admin creates a recovery course with session slots and enrolled students
func TestRecovery_CreateCourse_Success(t *testing.T) {
	repo := newMockRecoveryRepo()
	r := setupRecoveryRouter(repo, "teacher", "teacher-math-1", "school-1")

	body, _ := json.Marshal(recovery.CreateCourseRequest{
		SubjectID:    "subj-math",
		TeacherID:    "teacher-math-1",
		Title:        "Corso di Recupero Matematica - Trimestre 1",
		Description:  "Ripasso equazioni di secondo grado e trigonometria",
		AcademicYear: "2025/2026",
		Period:       "intermedio",
		TotalHours:   10,
		Room:         "Aula Magna",
		StudentIDs:   []string{"student-deb-1", "student-deb-2"},
		Sessions: []recovery.RecoveryCourseSession{
			{
				SessionDate: "2026-09-01",
				StartTime:   "15:00",
				EndTime:     "17:00",
				Room:        "Aula Magna",
				Topic:       "Equazioni di 2° grado",
			},
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/recovery/courses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var course recovery.RecoveryCourse
	err := json.Unmarshal(w.Body.Bytes(), &course)
	require.NoError(t, err)
	assert.Equal(t, "Corso di Recupero Matematica - Trimestre 1", course.Title)
	assert.Equal(t, "scheduled", course.Status)
	assert.Len(t, course.Students, 2)
	assert.Len(t, course.Sessions, 1)
}

// REC02 — Transition course status from scheduled to in_progress to completed
func TestRecovery_UpdateCourseStatus(t *testing.T) {
	repo := newMockRecoveryRepo()
	_ = repo.CreateCourse(context.Background(), &recovery.RecoveryCourse{
		ID:       "course-status-test",
		SchoolID: "school-1",
		Title:    "Fisica Recupero",
	}, nil, nil)

	r := setupRecoveryRouter(repo, "admin", "admin-1", "school-1")

	body, _ := json.Marshal(map[string]string{"status": "in_progress"})
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/recovery/courses/course-status-test/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	course, _ := repo.GetCourseByID(context.Background(), "course-status-test")
	assert.Equal(t, "in_progress", course.Status)
}

// REC03 — Update student attendance hours in recovery course
func TestRecovery_UpdateStudentAttendance(t *testing.T) {
	repo := newMockRecoveryRepo()
	_ = repo.CreateCourse(context.Background(), &recovery.RecoveryCourse{
		ID:       "course-att-test",
		SchoolID: "school-1",
		Title:    "Inglese Recupero",
	}, nil, []string{"student-1"})

	r := setupRecoveryRouter(repo, "teacher", "teacher-1", "school-1")

	body, _ := json.Marshal(map[string]interface{}{
		"attendance_hours": 8.0,
		"notes":            "Frequenza regolare, partecipazione attiva",
	})
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/recovery/courses/course-att-test/students/student-1/attendance", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	course, _ := repo.GetCourseByID(context.Background(), "course-att-test")
	assert.Equal(t, 8.0, course.Students[0].AttendanceHours)
}

// REC04 — Record recovery test exam result with passed / failed outcome
func TestRecovery_RecordRecoveryTest(t *testing.T) {
	repo := newMockRecoveryRepo()
	r := setupRecoveryRouter(repo, "teacher", "teacher-1", "school-1")

	// Passed test (grade 7.5)
	bodyPassed, _ := json.Marshal(recovery.RecordTestOutcomeRequest{
		StudentID: "student-1",
		SubjectID: "subj-math",
		ClassID:   "class-3A",
		TestDate:  "2026-09-10",
		TestType:  "written",
		Grade:     7.5,
		Notes:     "Debito saldato positivamente",
	})
	reqPassed := httptest.NewRequest(http.MethodPost, "/api/v1/recovery/tests", bytes.NewReader(bodyPassed))
	reqPassed.Header.Set("Content-Type", "application/json")
	wPassed := httptest.NewRecorder()
	r.ServeHTTP(wPassed, reqPassed)

	require.Equal(t, http.StatusCreated, wPassed.Code)
	var testResp recovery.RecoveryTest
	_ = json.Unmarshal(wPassed.Body.Bytes(), &testResp)
	assert.Equal(t, 7.5, testResp.Grade)
	assert.Equal(t, "recuperato", testResp.Outcome)
}

// REC05 — List recovery tests for a student
func TestRecovery_ListRecoveryTests_Student(t *testing.T) {
	repo := newMockRecoveryRepo()
	_ = repo.RecordRecoveryTest(context.Background(), &recovery.RecoveryTest{
		SchoolID:  "school-1",
		StudentID: "student-1",
		Grade:     8.0,
		Outcome:   "recuperato",
	})

	r := setupRecoveryRouter(repo, "student", "student-1", "school-1")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/recovery/tests?student_id=student-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var list []recovery.RecoveryTest
	err := json.Unmarshal(w.Body.Bytes(), &list)
	require.NoError(t, err)
	assert.Len(t, list, 1)
}
