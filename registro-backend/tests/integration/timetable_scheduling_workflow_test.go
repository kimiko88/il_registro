package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"registro-backend/internal/timetables"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockTimetablesRepo struct {
	classSchedules map[string][]timetables.ClassSchedule
}

func (m *mockTimetablesRepo) GetByClass(ctx context.Context, classID string) ([]timetables.ClassSchedule, error) {
	return m.classSchedules[classID], nil
}

func (m *mockTimetablesRepo) GetStudentClassID(ctx context.Context, userID string) (string, error) {
	return "class-1", nil
}

func (m *mockTimetablesRepo) GetParentStudentClassID(ctx context.Context, userID string) (string, error) {
	return "class-1", nil
}

func (m *mockTimetablesRepo) GetTeacherSchedule(ctx context.Context, userID string) ([]timetables.ClassSchedule, error) {
	return []timetables.ClassSchedule{}, nil
}

func (m *mockTimetablesRepo) GetClassSchoolID(ctx context.Context, classID string) (string, error) {
	return "school-1", nil
}

func (m *mockTimetablesRepo) Update(ctx context.Context, classID string, entries []timetables.ScheduleEntry) error {
	schedules := make([]timetables.ClassSchedule, len(entries))
	for i, e := range entries {
		schedules[i] = timetables.ClassSchedule{
			ID:          "sched-" + string(rune(i)),
			ClassID:     classID,
			DayOfWeek:   e.DayOfWeek,
			HourIndex:   e.HourIndex,
			SubjectID:   e.SubjectID,
			TeacherID:   e.TeacherID,
			SubjectName: "Matematica",
			TeacherName: "Prof. Rossi",
			Room:        e.Room,
		}
	}
	m.classSchedules[classID] = schedules
	return nil
}

func (m *mockTimetablesRepo) UpdateTeacher(ctx context.Context, teacherID string, entries []timetables.TeacherScheduleEntry) error {
	return nil
}

func TestIntegration_Timetable_Scheduling_Workflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockTimetablesRepo{classSchedules: make(map[string][]timetables.ClassSchedule)}
	handler := timetables.NewHandler(repo)

	r := gin.New()
	r.GET("/timetables/class/:id", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.GetByClass(c)
	})
	r.PUT("/timetables/class/:id", func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
		handler.Update(c)
	})

	// 1. Update class timetable schedule
	teacherID := "teacher-1"
	reqPayload := timetables.UpdateScheduleRequest{
		Entries: []timetables.ScheduleEntry{
			{DayOfWeek: 1, HourIndex: 1, SubjectID: "subj-1", TeacherID: &teacherID, Room: "Aula 101"},
			{DayOfWeek: 1, HourIndex: 2, SubjectID: "subj-1", TeacherID: &teacherID, Room: "Aula 101"},
		},
	}
	body, _ := json.Marshal(reqPayload)
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("PUT", "/timetables/class/class-1", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// 2. Fetch updated class timetable
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/timetables/class/class-1", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}
