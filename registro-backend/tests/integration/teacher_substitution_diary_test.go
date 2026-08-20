package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/lessons"
	"registro-backend/internal/substitutions"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_Teacher_Substitution_Diary(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSubRepo := &mockSubstitutionsRepoForDiary{}
	mockLessonsRepo := &mockLessonsRepoForSubDiary{}

	subSvc := substitutions.NewService(mockSubRepo)
	subHandler := substitutions.NewHandler(subSvc)

	lessonsSvc := lessons.NewService(mockLessonsRepo)
	lessonsHandler := lessons.NewHandler(lessonsSvc)

	r := gin.New()
	r.POST("/substitutions", func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
		subHandler.Create(c)
	})

	r.POST("/lessons", func(c *gin.Context) {
		c.Set("user_id", "sub-teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		lessonsHandler.CreateLesson(c)
	})

	// 1. Create substitution assignment
	subTeacherID := "sub-teacher-1"
	subReq := substitutions.CreateSubstitutionRequest{
		ClassID:             "class-1",
		AbsentTeacherID:     "teacher-1",
		SubstituteTeacherID: &subTeacherID,
		Date:                "2026-09-05",
		Hour:                2,
	}
	jsonBody1, _ := json.Marshal(subReq)

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/substitutions", bytes.NewBuffer(jsonBody1))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// 2. Substitute teacher logs lesson in diary
	lessReq := lessons.CreateLessonRequest{
		ClassID:        "class-1",
		SubjectID:      "subject-1",
		Date:           "2026-09-05",
		Hour:           2,
		Topic:          "Sostituzione - Ripasso Matematica",
		Type:           "Frontale",
		Duration:       60,
		IsSubstitution: true,
	}
	jsonBody2, _ := json.Marshal(lessReq)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/lessons", bytes.NewBuffer(jsonBody2))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusCreated, w2.Code)
}

type mockSubstitutionsRepoForDiary struct{}

func (m *mockSubstitutionsRepoForDiary) Create(ctx context.Context, sub *substitutions.Substitution) error {
	sub.ID = "sub-1"
	return nil
}
func (m *mockSubstitutionsRepoForDiary) GetByID(ctx context.Context, id string) (*substitutions.Substitution, error) {
	return &substitutions.Substitution{ID: id, SchoolID: "school-1"}, nil
}
func (m *mockSubstitutionsRepoForDiary) ListBySchool(ctx context.Context, schoolID, date string) ([]*substitutions.Substitution, error) {
	return []*substitutions.Substitution{}, nil
}
func (m *mockSubstitutionsRepoForDiary) ListByTeacher(ctx context.Context, teacherID string, date string) ([]*substitutions.Substitution, error) {
	return []*substitutions.Substitution{}, nil
}
func (m *mockSubstitutionsRepoForDiary) AssignSubstitute(ctx context.Context, id string, substituteTeacherID string, notes string) error {
	return nil
}
func (m *mockSubstitutionsRepoForDiary) ConfirmSubstitution(ctx context.Context, id string, substituteTeacherID string) error {
	return nil
}
func (m *mockSubstitutionsRepoForDiary) SignRegister(ctx context.Context, id string, sigHash string, notes string) error {
	return nil
}
func (m *mockSubstitutionsRepoForDiary) GetAvailableTeachers(ctx context.Context, schoolID string) ([]substitutions.TeacherCandidate, error) {
	return []substitutions.TeacherCandidate{}, nil
}
func (m *mockSubstitutionsRepoForDiary) GetTeacherProfileID(ctx context.Context, userID string) (string, error) {
	return "sub-teacher-prof-1", nil
}
func (m *mockSubstitutionsRepoForDiary) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	return true, nil
}
func (m *mockSubstitutionsRepoForDiary) IsTeacherAssignedToSubject(ctx context.Context, teacherID, subjectID string) (bool, error) {
	return true, nil
}
func (m *mockSubstitutionsRepoForDiary) GetWeeklySubstitutionCount(ctx context.Context, teacherID string) (int, error) {
	return 1, nil
}

type mockLessonsRepoForSubDiary struct{}

func (m *mockLessonsRepoForSubDiary) CreateLesson(l *lessons.Lesson) error {
	l.ID = "lesson-sub-1"
	return nil
}
func (m *mockLessonsRepoForSubDiary) GetLessonByID(id string) (*lessons.Lesson, error) {
	return &lessons.Lesson{ID: id}, nil
}
func (m *mockLessonsRepoForSubDiary) UpdateLesson(id string, req lessons.UpdateLessonRequest) (*lessons.Lesson, error) {
	return &lessons.Lesson{ID: id}, nil
}
func (m *mockLessonsRepoForSubDiary) DeleteLesson(id string) error { return nil }
func (m *mockLessonsRepoForSubDiary) GetLessonsByClass(classID string, date string) ([]lessons.Lesson, error) {
	return []lessons.Lesson{}, nil
}
func (m *mockLessonsRepoForSubDiary) GetLessonsByClassAndSubject(classID, subjectID string, date string) ([]lessons.Lesson, error) {
	return []lessons.Lesson{}, nil
}
func (m *mockLessonsRepoForSubDiary) GetLessonsByGroup(groupID string, date string) ([]lessons.Lesson, error) {
	return []lessons.Lesson{}, nil
}
func (m *mockLessonsRepoForSubDiary) GetLessonsByTeacher(teacherID string, fromDate, toDate string) ([]lessons.Lesson, error) {
	return []lessons.Lesson{}, nil
}
func (m *mockLessonsRepoForSubDiary) IsTeacherAssignedToClass(teacherID, classID string) (bool, error) {
	return false, nil // Substitute is not assigned to class directly
}
func (m *mockLessonsRepoForSubDiary) HasApprovedSubstitution(teacherID, classID, date string, hour int) (bool, error) {
	return true, nil // Approved substitution exists!
}
func (m *mockLessonsRepoForSubDiary) CreateHomework(h *lessons.Homework) error { return nil }
func (m *mockLessonsRepoForSubDiary) GetHomeworkByID(id string) (*lessons.Homework, error) {
	return &lessons.Homework{ID: id}, nil
}
func (m *mockLessonsRepoForSubDiary) UpdateHomework(id string, req lessons.UpdateHomeworkRequest) (*lessons.Homework, error) {
	return &lessons.Homework{ID: id}, nil
}
func (m *mockLessonsRepoForSubDiary) DeleteHomework(id string) error { return nil }
func (m *mockLessonsRepoForSubDiary) GetHomeworkByClass(classID string) ([]lessons.Homework, error) {
	return []lessons.Homework{}, nil
}
