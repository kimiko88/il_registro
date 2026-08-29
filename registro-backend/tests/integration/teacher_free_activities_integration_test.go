package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/teacher_activities"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockTeacherActivitiesRepo struct {
	activities map[string]*teacher_activities.TeacherFreeActivity
}

func newMockTeacherActivitiesRepo() *mockTeacherActivitiesRepo {
	return &mockTeacherActivitiesRepo{
		activities: make(map[string]*teacher_activities.TeacherFreeActivity),
	}
}

func (m *mockTeacherActivitiesRepo) Create(act *teacher_activities.TeacherFreeActivity) error {
	if act.ID == "" {
		act.ID = "act-" + time.Now().Format("150405.000000")
	}
	act.TeacherName = "Prof. Mario Rossi"
	act.CreatedAt = time.Now()
	act.UpdatedAt = time.Now()
	m.activities[act.ID] = act
	return nil
}

func (m *mockTeacherActivitiesRepo) GetByID(id string) (*teacher_activities.TeacherFreeActivity, error) {
	act, ok := m.activities[id]
	if !ok {
		return nil, assert.AnError
	}
	return act, nil
}

func (m *mockTeacherActivitiesRepo) Update(id string, req teacher_activities.UpdateTeacherActivityRequest) (*teacher_activities.TeacherFreeActivity, error) {
	act, ok := m.activities[id]
	if !ok {
		return nil, assert.AnError
	}
	if req.Description != "" {
		act.Description = req.Description
	}
	if req.Notes != "" {
		act.Notes = req.Notes
	}
	act.UpdatedAt = time.Now()
	return act, nil
}

func (m *mockTeacherActivitiesRepo) Delete(id string) error {
	delete(m.activities, id)
	return nil
}

func (m *mockTeacherActivitiesRepo) GetByTeacher(teacherID, fromDate, toDate string) ([]teacher_activities.TeacherFreeActivity, error) {
	var res []teacher_activities.TeacherFreeActivity
	for _, act := range m.activities {
		if act.TeacherID == teacherID {
			res = append(res, *act)
		}
	}
	return res, nil
}

func setupTeacherActivitiesRouter(repo teacher_activities.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := teacher_activities.NewService(repo)
	handler := teacher_activities.NewHandler(svc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		c.Set("role", "teacher")
		c.Set("user_id", "teacher-1")
		c.Set("school_id", "school-1")
		c.Next()
	})
	handler.RegisterRoutes(api)
	return r
}

func TestIntegration_TeacherFreeActivities_Workflow(t *testing.T) {
	repo := newMockTeacherActivitiesRepo()
	r := setupTeacherActivitiesRouter(repo)

	today := time.Now().Format("2006-01-02")

	// 1. Teacher logs a non-teaching activity (Collegio Docenti / Riunione dipartimento)
	createReq := teacher_activities.CreateTeacherActivityRequest{
		Date:         today,
		StartHour:    3,
		Duration:     2,
		ActivityType: "riunione",
		Description:  "Riunione Dipartimento Scientifico per programmazione percorsi PCTO",
		Notes:        "Verbale n. 4 approvato all'unanimità",
	}
	body, _ := json.Marshal(createReq)
	reqCreate, _ := http.NewRequest("POST", "/api/v1/teacher/free-activities", bytes.NewReader(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	r.ServeHTTP(wCreate, reqCreate)
	require.Equal(t, http.StatusCreated, wCreate.Code)

	var created teacher_activities.TeacherFreeActivity
	err := json.Unmarshal(wCreate.Body.Bytes(), &created)
	require.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "riunione", created.ActivityType)
	assert.Equal(t, 2, created.Duration)

	// 2. Fetch logged free activities for the teacher
	reqList, _ := http.NewRequest("GET", "/api/v1/teacher/free-activities", nil)
	wList := httptest.NewRecorder()
	r.ServeHTTP(wList, reqList)
	require.Equal(t, http.StatusOK, wList.Code)

	var list []teacher_activities.TeacherFreeActivity
	err = json.Unmarshal(wList.Body.Bytes(), &list)
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, created.ID, list[0].ID)

	// 3. Update the activity notes
	updateReq := teacher_activities.UpdateTeacherActivityRequest{
		Notes: "Verbale n. 4 archiviato e protocollato",
	}
	updateBody, _ := json.Marshal(updateReq)
	reqUpdate, _ := http.NewRequest("PUT", "/api/v1/teacher/free-activities/"+created.ID, bytes.NewReader(updateBody))
	reqUpdate.Header.Set("Content-Type", "application/json")
	wUpdate := httptest.NewRecorder()
	r.ServeHTTP(wUpdate, reqUpdate)
	require.Equal(t, http.StatusOK, wUpdate.Code)

	// 4. Delete the logged activity
	reqDel, _ := http.NewRequest("DELETE", "/api/v1/teacher/free-activities/"+created.ID, nil)
	wDel := httptest.NewRecorder()
	r.ServeHTTP(wDel, reqDel)
	require.Equal(t, http.StatusNoContent, wDel.Code)
	assert.Empty(t, repo.activities)
}
