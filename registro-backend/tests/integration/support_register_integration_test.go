package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/support"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSupportRepo struct {
	diaries map[string]*support.SupportDiaryEntry
	goals   map[string]*support.SupportPeiGoal
}

func newMockSupportRepo() *mockSupportRepo {
	return &mockSupportRepo{
		diaries: make(map[string]*support.SupportDiaryEntry),
		goals:   make(map[string]*support.SupportPeiGoal),
	}
}

func (m *mockSupportRepo) CreateDiaryEntry(ctx context.Context, entry *support.SupportDiaryEntry) error {
	if entry.ID == "" {
		entry.ID = "diary-" + time.Now().Format("150405.000000")
	}
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()
	m.diaries[entry.ID] = entry
	return nil
}

func (m *mockSupportRepo) ListDiaryEntries(ctx context.Context, schoolID, studentID, teacherID, classID string, isFamily bool) ([]support.SupportDiaryEntry, error) {
	var res []support.SupportDiaryEntry
	for _, d := range m.diaries {
		if schoolID != "" && d.SchoolID != schoolID {
			continue
		}
		if studentID != "" && d.StudentID != studentID {
			continue
		}
		if isFamily && !d.IsSharedWithFamily {
			continue
		}
		res = append(res, *d)
	}
	return res, nil
}

func (m *mockSupportRepo) DeleteDiaryEntry(ctx context.Context, id, teacherID string) error {
	delete(m.diaries, id)
	return nil
}

func (m *mockSupportRepo) CreatePeiGoal(ctx context.Context, goal *support.SupportPeiGoal) error {
	if goal.ID == "" {
		goal.ID = "goal-" + time.Now().Format("150405.000000")
	}
	goal.CreatedAt = time.Now()
	goal.UpdatedAt = time.Now()
	m.goals[goal.ID] = goal
	return nil
}

func (m *mockSupportRepo) UpdatePeiGoalProgress(ctx context.Context, id, progressStatus string) error {
	if g, ok := m.goals[id]; ok {
		g.ProgressStatus = progressStatus
		g.UpdatedAt = time.Now()
	}
	return nil
}

func (m *mockSupportRepo) ListPeiGoals(ctx context.Context, schoolID, studentID string) ([]support.SupportPeiGoal, error) {
	var res []support.SupportPeiGoal
	for _, g := range m.goals {
		if schoolID != "" && g.SchoolID != schoolID {
			continue
		}
		if studentID != "" && g.StudentID != studentID {
			continue
		}
		res = append(res, *g)
	}
	return res, nil
}

func (m *mockSupportRepo) DeletePeiGoal(ctx context.Context, id string) error {
	delete(m.goals, id)
	return nil
}

func setupSupportTestRouter(repo support.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := support.NewService(repo)
	handler := support.NewHandler(svc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		role := c.GetHeader("X-Role")
		if role == "" {
			role = "teacher"
		}
		c.Set("role", role)
		c.Set("user_id", "teacher-sostegno-1")
		c.Set("school_id", "school-1")
		c.Next()
	})
	handler.RegisterRoutes(api)
	return r
}

func TestIntegration_SupportRegister_FullWorkflow(t *testing.T) {
	repo := newMockSupportRepo()
	r := setupSupportTestRouter(repo)

	// 1. Create a support diary entry
	diaryPayload := support.CreateDiaryEntryRequest{
		StudentID:          "student-bes-1",
		ClassID:            "class-2A",
		EntryDate:          time.Now().Format("2006-01-02"),
		TimeSlot:           "2ª Ora (09:00 - 10:00)",
		ActivityType:       "in_classe",
		TopicAndActivities: "Attività di supporto personalizzato per Matematica (frazioni e percentuali)",
		StudentResponses:   "Buona autonomia con l'uso della mappa concettuale",
		EducatorNotes:      "Continuare con esercizi guidati",
		IsSharedWithFamily: true,
	}
	body, _ := json.Marshal(diaryPayload)
	req, _ := http.NewRequest("POST", "/api/v1/support/diaries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var createdDiary support.SupportDiaryEntry
	err := json.Unmarshal(w.Body.Bytes(), &createdDiary)
	require.NoError(t, err)
	assert.NotEmpty(t, createdDiary.ID)
	assert.Equal(t, "student-bes-1", createdDiary.StudentID)
	assert.True(t, createdDiary.IsSharedWithFamily)

	// 2. List diary entries as teacher
	reqList, _ := http.NewRequest("GET", "/api/v1/support/diaries?student_id=student-bes-1", nil)
	wList := httptest.NewRecorder()
	r.ServeHTTP(wList, reqList)
	require.Equal(t, http.StatusOK, wList.Code)

	var diaries []support.SupportDiaryEntry
	err = json.Unmarshal(wList.Body.Bytes(), &diaries)
	require.NoError(t, err)
	assert.Len(t, diaries, 1)

	// 3. Create a PEI Goal for the student
	goalPayload := support.CreatePeiGoalRequest{
		StudentID:    "student-bes-1",
		PeiType:      "differenziato",
		Axis:         "cognitiva",
		Title:        "Risoluzione autonoma di problemi con proporzioni",
		Description:  "Utilizzo di schemi visivi e calcolatrice per il calcolo percentuale",
		ExpectedTerm: "q2",
	}
	goalBody, _ := json.Marshal(goalPayload)
	reqGoal, _ := http.NewRequest("POST", "/api/v1/support/pei-goals", bytes.NewReader(goalBody))
	reqGoal.Header.Set("Content-Type", "application/json")
	wGoal := httptest.NewRecorder()
	r.ServeHTTP(wGoal, reqGoal)
	require.Equal(t, http.StatusCreated, wGoal.Code)

	var createdGoal support.SupportPeiGoal
	err = json.Unmarshal(wGoal.Body.Bytes(), &createdGoal)
	require.NoError(t, err)
	assert.NotEmpty(t, createdGoal.ID)
	assert.Equal(t, "non_avviato", createdGoal.ProgressStatus)

	// 4. Update PEI Goal progress to "raggiunto"
	patchReq, _ := http.NewRequest("PATCH", "/api/v1/support/pei-goals/"+createdGoal.ID+"/progress", bytes.NewReader([]byte(`{"progress_status": "raggiunto"}`)))
	patchReq.Header.Set("Content-Type", "application/json")
	wPatch := httptest.NewRecorder()
	r.ServeHTTP(wPatch, patchReq)
	require.Equal(t, http.StatusOK, wPatch.Code)

	// 5. Verify updated PEI goals list
	reqGoalsList, _ := http.NewRequest("GET", "/api/v1/support/pei-goals?student_id=student-bes-1", nil)
	wGoalsList := httptest.NewRecorder()
	r.ServeHTTP(wGoalsList, reqGoalsList)
	require.Equal(t, http.StatusOK, wGoalsList.Code)

	var goals []support.SupportPeiGoal
	err = json.Unmarshal(wGoalsList.Body.Bytes(), &goals)
	require.NoError(t, err)
	assert.Len(t, goals, 1)
	assert.Equal(t, "raggiunto", goals[0].ProgressStatus)
}
