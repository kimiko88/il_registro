package notes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupNotesRouter(h *Handler, userID, schoolID, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		if schoolID != "" {
			c.Set("school_id", schoolID)
		}
		if role != "" {
			c.Set("role", role)
		}
		c.Next()
	})
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestNotesHandler_Create(t *testing.T) {
	mockRepo := new(MockRepository)
	mockUsers := new(MockUserRepo)
	svc := NewService(mockRepo, mockUsers)
	h := NewHandler(svc)

	// 1. Unauthorized
	rNoUser := setupNotesRouter(h, "", "school-1", "teacher")
	req, _ := http.NewRequest("POST", "/api/v1/notes", bytes.NewBufferString("{}"))
	w := httptest.NewRecorder()
	rNoUser.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 2. Forbidden role (student cannot create notes)
	rStudent := setupNotesRouter(h, "s-1", "school-1", "student")
	req, _ = http.NewRequest("POST", "/api/v1/notes", bytes.NewBufferString("{}"))
	w = httptest.NewRecorder()
	rStudent.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 3. Invalid JSON
	rTeacher := setupNotesRouter(h, "t-1", "school-1", "teacher")
	req, _ = http.NewRequest("POST", "/api/v1/notes", bytes.NewBufferString("{invalid"))
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 4. Success -> 201
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	body, _ := json.Marshal(CreateNoteRequest{
		StudentID: "stud-1",
		ClassID:   "class-1",
		Type:      NoteTypeDisciplinary,
		Note:      "Comportamento scorretto durante la lezione",
		Date:      "2026-09-12",
	})
	req, _ = http.NewRequest("POST", "/api/v1/notes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestNotesHandler_ApproveUpdateDelete(t *testing.T) {
	mockRepo := new(MockRepository)
	mockUsers := new(MockUserRepo)
	svc := NewService(mockRepo, mockUsers)
	h := NewHandler(svc)
	rTeacher := setupNotesRouter(h, "t-1", "school-1", "teacher")

	existingNote := &StudentNote{
		ID:        "note-1",
		TeacherID: "t-1",
		StudentID: "stud-1",
		ClassID:   "class-1",
		Type:      NoteTypeDisciplinary,
		Note:      "Nota",
		CreatedAt: time.Now(),
	}
	mockRepo.On("Get", mock.Anything, "note-1").Return(existingNote, nil)

	// 1. Approve note (requires principal/admin)
	rPrincipal := setupNotesRouter(h, "p-1", "school-1", "principal")
	mockRepo.On("ApproveNote", mock.Anything, "note-1", "p-1").Return(nil)
	req, _ := http.NewRequest("POST", "/api/v1/notes/note-1/approve", nil)
	w := httptest.NewRecorder()
	rPrincipal.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. Update note
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
	updateBody, _ := json.Marshal(UpdateNoteRequest{
		Note: "Nota aggiornata",
	})
	req, _ = http.NewRequest("PATCH", "/api/v1/notes/note-1", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. Delete note invalid JSON -> 400
	req, _ = http.NewRequest("DELETE", "/api/v1/notes/note-1", bytes.NewBufferString("{invalid"))
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// 4. Delete note success -> 204
	mockRepo.On("Delete", mock.Anything, "note-1").Return(nil)
	delBody, _ := json.Marshal(DeleteNoteRequest{
		Reason: "Annullata per errore",
	})
	req, _ = http.NewRequest("DELETE", "/api/v1/notes/note-1", bytes.NewBuffer(delBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestNotesHandler_ListAndGetByStudent(t *testing.T) {
	mockRepo := new(MockRepository)
	mockUsers := new(MockUserRepo)
	svc := NewService(mockRepo, mockUsers)
	h := NewHandler(svc)
	rTeacher := setupNotesRouter(h, "t-1", "school-1", "teacher")

	mockRepo.On("List", mock.Anything, mock.Anything).Return([]StudentNote{
		{ID: "note-1", StudentID: "stud-1", Note: "Nota 1"},
	}, nil)

	// 1. List -> 200
	req, _ := http.NewRequest("GET", "/api/v1/notes?student_id=stud-1", nil)
	w := httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. GetByStudent -> 200
	req, _ = http.NewRequest("GET", "/api/v1/notes/student/stud-1", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
