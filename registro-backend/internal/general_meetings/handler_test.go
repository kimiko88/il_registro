package general_meetings

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupMeetingRouter(svc Service, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("school_id", "school-1")
		c.Set("user_id", "user-admin-1")
		c.Set("role", role)
		c.Next()
	})
	h := NewHandler(svc)
	rg := r.Group("/api/v1")
	h.RegisterRoutes(rg)
	return r
}

func TestMeetingHandler_RoleAccess(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	rStudent := setupMeetingRouter(svc, "student")

	// Student cannot create general meeting
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/general-meetings", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	rStudent.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden 403, got %d", w.Code)
	}

	// Admin can create general meeting
	rAdmin := setupMeetingRouter(svc, "admin")
	body, _ := json.Marshal(CreateGeneralMeetingRequest{
		Title:       "Collegio Docenti Straordinario",
		Location:    "Aula Magna",
		MeetingDate: time.Now().AddDate(0, 0, 7).Format("2006-01-02"),
		IsMandatory: true,
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/general-meetings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rAdmin.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 created, got %d: %s", w.Code, w.Body.String())
	}
}

func TestMeetingHandler_ListAndRegister(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	rTeacher := setupMeetingRouter(svc, "teacher")

	// Create meeting in repo
	repo.meetings["m-1"] = &GeneralMeeting{
		ID:          "m-1",
		SchoolID:    "school-1",
		Title:       "Assemblea d'Istituto",
		MeetingDate: time.Now().AddDate(0, 0, 3),
	}

	// List
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/general-meetings", nil)
	rTeacher.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 ok, got %d", w.Code)
	}

	// Register
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/general-meetings/m-1/register", nil)
	rTeacher.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 ok on registration, got %d", w.Code)
	}
}
