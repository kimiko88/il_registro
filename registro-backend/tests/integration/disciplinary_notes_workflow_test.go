package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/notes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_Disciplinary_Notes_Workflow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &mockNotesRepoForWorkflow{}
	svc := notes.NewService(mockRepo, &mockUserRepoForComms{})
	handler := notes.NewHandler(svc)

	r := gin.New()
	r.POST("/notes", func(c *gin.Context) {
		c.Set("user_id", "teacher-1")
		c.Set("role", "teacher")
		c.Set("school_id", "school-1")
		handler.Create(c)
	})

	r.POST("/notes/:id/approve", func(c *gin.Context) {
		c.Set("user_id", "admin-1")
		c.Set("role", "admin")
		c.Set("school_id", "school-1")
		handler.Approve(c)
	})

	// 1. Create disciplinary note
	body := notes.CreateNoteRequest{
		StudentID: "student-1",
		ClassID:   "class-1",
		Type:      "disciplinary",
		Note:      "Comportamento scorretto durante la lezione.",
		Date:      time.Now().Format("2006-01-02"),
	}
	jsonBody, _ := json.Marshal(body)

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/notes", bytes.NewBuffer(jsonBody))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// 2. Approve disciplinary note
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/notes/note-1/approve", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

type mockNotesRepoForWorkflow struct{}

func (m *mockNotesRepoForWorkflow) Create(ctx context.Context, note *notes.StudentNote) error {
	note.ID = "note-1"
	return nil
}
func (m *mockNotesRepoForWorkflow) Update(ctx context.Context, note *notes.StudentNote) error {
	return nil
}
func (m *mockNotesRepoForWorkflow) Delete(ctx context.Context, id string) error { return nil }
func (m *mockNotesRepoForWorkflow) DeleteWithReason(ctx context.Context, id string, reason string) error {
	return nil
}
func (m *mockNotesRepoForWorkflow) Get(ctx context.Context, id string) (*notes.StudentNote, error) {
	return &notes.StudentNote{ID: id, SchoolID: "school-1", TeacherID: "teacher-1", Type: "disciplinary", IsApproved: false}, nil
}
func (m *mockNotesRepoForWorkflow) List(ctx context.Context, filter notes.NoteFilter) ([]notes.StudentNote, error) {
	return []notes.StudentNote{}, nil
}
func (m *mockNotesRepoForWorkflow) ApproveNote(ctx context.Context, id, approverID string) error {
	return nil
}
func (m *mockNotesRepoForWorkflow) MarkAsViewedByParent(ctx context.Context, id string) error {
	return nil
}
func (m *mockNotesRepoForWorkflow) IsTeacherAssignedToClass(ctx context.Context, teacherID, classID string) (bool, error) {
	return true, nil
}
