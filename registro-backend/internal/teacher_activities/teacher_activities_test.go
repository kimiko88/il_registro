package teacher_activities

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


type mockRepo struct {
	activities map[string]*TeacherFreeActivity
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		activities: make(map[string]*TeacherFreeActivity),
	}
}

func (m *mockRepo) Create(act *TeacherFreeActivity) error {
	act.ID = "act-1"
	act.CreatedAt = time.Now()
	act.UpdatedAt = time.Now()
	m.activities[act.ID] = act
	return nil
}

func (m *mockRepo) GetByID(id string) (*TeacherFreeActivity, error) {
	act, ok := m.activities[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return act, nil
}

func (m *mockRepo) Update(id string, req UpdateTeacherActivityRequest) (*TeacherFreeActivity, error) {
	act, ok := m.activities[id]
	if !ok {
		return nil, errors.New("not found")
	}
	if req.ActivityType != "" {
		act.ActivityType = req.ActivityType
	}
	if req.Description != "" {
		act.Description = req.Description
	}
	act.UpdatedAt = time.Now()
	return act, nil
}

func (m *mockRepo) Delete(id string) error {
	if _, ok := m.activities[id]; !ok {
		return errors.New("not found")
	}
	delete(m.activities, id)
	return nil
}

func (m *mockRepo) GetByTeacher(teacherID, fromDate, toDate string) ([]TeacherFreeActivity, error) {
	var result []TeacherFreeActivity
	for _, act := range m.activities {
		if act.TeacherID == teacherID {
			result = append(result, *act)
		}
	}
	return result, nil
}

func TestTeacherActivitiesService_CreateValidations(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	todayStr := time.Now().Format("2006-01-02")

	// 1. Missing description
	_, err := svc.Create("teacher-1", CreateTeacherActivityRequest{
		Date:         todayStr,
		StartHour:    1,
		Duration:     2,
		ActivityType: "riunione",
		Description:  "",
	})
	if err == nil {
		t.Error("expected error for empty description")
	}

	// 2. Invalid start_hour (< 1 or > 10)
	_, err = svc.Create("teacher-1", CreateTeacherActivityRequest{
		Date:         todayStr,
		StartHour:    11,
		Duration:     1,
		ActivityType: "riunione",
		Description:  "Riunione di dipartimento",
	})
	if err == nil {
		t.Error("expected error for start_hour > 10")
	}

	// 3. Invalid activity_type
	_, err = svc.Create("teacher-1", CreateTeacherActivityRequest{
		Date:         todayStr,
		StartHour:    2,
		Duration:     1,
		ActivityType: "tipo_inesistente",
		Description:  "Attività non valida",
	})
	if err == nil {
		t.Error("expected error for invalid activity_type")
	}

	// 4. Valid creation
	res, err := svc.Create("teacher-1", CreateTeacherActivityRequest{
		Date:         todayStr,
		StartHour:    3,
		Duration:     2,
		ActivityType: "formazione",
		Description:  "Corso sulla didattica digitale",
	})
	if err != nil {
		t.Fatalf("unexpected error creating activity: %v", err)
	}

	if res.ID != "act-1" || res.ActivityType != "formazione" {
		t.Errorf("unexpected response: %+v", res)
	}
}

func TestTeacherActivitiesService_UpdateAuthorization(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	todayStr := time.Now().Format("2006-01-02")

	created, err := svc.Create("teacher-1", CreateTeacherActivityRequest{
		Date:         todayStr,
		StartHour:    1,
		Duration:     2,
		ActivityType: "ptof",
		Description:  "Progetto PTOF",
	})
	if err != nil {
		t.Fatalf("failed to create setup activity: %v", err)
	}

	// Another teacher trying to update teacher-1's activity
	newDesc := "Attempted hijack"
	_, err = svc.Update("teacher-2", created.ID, UpdateTeacherActivityRequest{
		Description: newDesc,
	})
	if err == nil {
		t.Error("expected error when teacher-2 tries to edit teacher-1's activity")
	}

	// Right teacher updating activity
	updated, err := svc.Update("teacher-1", created.ID, UpdateTeacherActivityRequest{
		Description: newDesc,
	})
	if err != nil {
		t.Fatalf("unexpected error on valid update: %v", err)
	}
	if updated.Description != newDesc {
		t.Errorf("expected description '%s', got '%s'", newDesc, updated.Description)
	}
}

func TestTeacherActivitiesService_GetByTeacherDateRange(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	// fromDate > toDate
	_, err := svc.GetByTeacher("teacher-1", "2026-06-01", "2026-01-01")
	if err == nil {
		t.Error("expected error when fromDate is after toDate")
	}

	// date range > 1 year
	_, err = svc.GetByTeacher("teacher-1", "2024-01-01", "2026-01-01")
	if err == nil {
		t.Error("expected error when date range exceeds 1 year")
	}
}

func TestTeacherActivitiesService_GetByID_Authorization(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)
	todayStr := time.Now().Format("2006-01-02")

	created, err := svc.Create("teacher-1", CreateTeacherActivityRequest{
		Date:         todayStr,
		StartHour:    1,
		Duration:     2,
		ActivityType: "riunione",
		Description:  "Consiglio straordinario",
	})
	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	t.Run("owner can get by ID", func(t *testing.T) {
		res, err := svc.GetByID("teacher-1", "teacher", created.ID)
		if err != nil || res == nil {
			t.Fatalf("expected success, got: %v", err)
		}
	})

	t.Run("admin can get by ID", func(t *testing.T) {
		res, err := svc.GetByID("admin-1", "admin", created.ID)
		if err != nil || res == nil {
			t.Fatalf("expected success, got: %v", err)
		}
	})

	t.Run("other teacher cannot get by ID", func(t *testing.T) {
		_, err := svc.GetByID("teacher-2", "teacher", created.ID)
		if err == nil {
			t.Fatal("expected unauthorized error for other teacher")
		}
	})
}

func TestTeacherActivitiesHandler_List_RBAC(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockRepo()
	svc := NewService(repo)
	h := NewHandler(svc)

	t.Run("Unauthenticated is 401", func(t *testing.T) {
		r := gin.New()
		h.RegisterRoutes(r.Group("/api"))

		req := httptest.NewRequest("GET", "/api/teacher/free-activities", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Student is 403", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("user_id", "student-1")
			c.Set("role", "student")
			c.Next()
		})
		h.RegisterRoutes(r.Group("/api"))

		req := httptest.NewRequest("GET", "/api/teacher/free-activities", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Teacher is 200", func(t *testing.T) {
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("user_id", "teacher-1")
			c.Set("role", "teacher")
			c.Next()
		})
		h.RegisterRoutes(r.Group("/api"))

		req := httptest.NewRequest("GET", "/api/teacher/free-activities", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

