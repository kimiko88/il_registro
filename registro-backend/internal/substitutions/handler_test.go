package substitutions

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupTestSubstitutionsRouter(h *Handler, authUserID, role, schoolID string, isStaff bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if authUserID != "" {
			c.Set("user_id", authUserID)
		}
		if role != "" {
			c.Set("role", role)
		}
		if schoolID != "" {
			c.Set("school_id", schoolID)
		}
		if isStaff {
			c.Set("is_staff", true)
		}
		c.Next()
	})
	api := r.Group("/api/v1")
	h.RegisterRoutes(api)
	return r
}

func TestHandler_Create(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "", "", "", false)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/substitutions", bytes.NewBufferString("{}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("forbidden_for_student", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "u1", "student", "s1", false)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/substitutions", bytes.NewBufferString("{}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})

	t.Run("bad_request_json", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "u1", "admin", "s1", false)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/substitutions", bytes.NewBufferString("{invalid"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("service_error_invalid_date", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "u1", "admin", "s1", false)

		body := CreateSubstitutionRequest{
			ClassID:         "c1",
			AbsentTeacherID: "t1",
			Date:            "invalid-date",
			Slot:            1,
		}
		b, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/substitutions", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("success_admin", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "u1", "admin", "s1", false)

		body := CreateSubstitutionRequest{
			ClassID:         "c1",
			AbsentTeacherID: "t1",
			Date:            "2026-03-01",
			Slot:            2,
			Subject:         "Storia",
		}
		b, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/substitutions", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("success_is_staff", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "u1", "other", "s1", true)

		body := CreateSubstitutionRequest{
			ClassID:         "c1",
			AbsentTeacherID: "t1",
			Date:            "2026-03-01",
			Hour:            3,
			Subject:         "Arte",
		}
		b, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/substitutions", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
		}
	})
}

func TestHandler_ListBySchool(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "", "", "", false)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/substitutions", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("forbidden_non_staff", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "u1", "parent", "s1", false)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/substitutions", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})

	t.Run("success_and_error", func(t *testing.T) {
		d, _ := time.Parse("2006-01-02", "2026-03-01")
		repo := &mockSubstitutionRepo{
			subs: []*Substitution{
				{ID: "sub1", SchoolID: "s1", Date: d, Status: StatusPending},
			},
		}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "u1", "teacher", "s1", false)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/substitutions?date=2026-03-01", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		// with error
		repo.getErr = errors.New("db error")
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req)
		if w2.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w2.Code)
		}
	})
}

func TestHandler_ListByTeacher(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "", "", "", false)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/substitutions/my", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("success_and_error", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "t1", "teacher", "s1", false)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/substitutions/my?date=2026-03-01", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		repo.getErr = errors.New("db failure")
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req)
		if w2.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w2.Code)
		}
	})
}

func TestHandler_ListMyToday(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "", "", "", false)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/substitutions/my-today", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("success_and_error", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "t1", "teacher", "s1", false)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/substitutions/my-today", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		repo.getErr = errors.New("db failure")
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req)
		if w2.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w2.Code)
		}
	})
}

func TestHandler_AssignSubstitute(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "", "", "", false)

		req, _ := http.NewRequest(http.MethodPut, "/api/v1/substitutions/sub1/assign", bytes.NewBufferString("{}"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("forbidden_non_staff", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "u1", "student", "s1", false)

		req, _ := http.NewRequest(http.MethodPut, "/api/v1/substitutions/sub1/assign", bytes.NewBufferString("{}"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})

	t.Run("bad_request_json", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "u1", "admin", "s1", false)

		req, _ := http.NewRequest(http.MethodPut, "/api/v1/substitutions/sub1/assign", bytes.NewBufferString("{invalid"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("success_and_error", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "u1", "admin", "s1", false)

		body := AssignSubstituteRequest{
			SubstituteTeacherID: "t-sub",
			Notes:               "Covering math class",
		}
		b, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/substitutions/sub1/assign", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		repo.saveErr = errors.New("cannot assign")
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest(http.MethodPut, "/api/v1/substitutions/sub1/assign", bytes.NewBuffer(b))
		req2.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w2, req2)
		if w2.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w2.Code)
		}
	})
}

func TestHandler_Confirm(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "", "", "", false)

		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/substitutions/sub1/confirm", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("forbidden_role", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "u1", "student", "s1", false)

		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/substitutions/sub1/confirm", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})

	t.Run("success_and_service_error", func(t *testing.T) {
		subTeacher := "t-sub"
		repo := &mockSubstitutionRepo{
			subs: []*Substitution{
				{
					ID:                  "sub1",
					SubstituteTeacherID: &subTeacher,
					Status:              StatusAssigned,
				},
			},
		}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "t-sub", "teacher", "s1", false)

		req, _ := http.NewRequest(http.MethodPatch, "/api/v1/substitutions/sub1/confirm", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		// Wrong teacher error
		r2 := setupTestSubstitutionsRouter(h, "t-other", "teacher", "s1", false)
		w2 := httptest.NewRecorder()
		r2.ServeHTTP(w2, req)
		if w2.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w2.Code)
		}
	})
}

func TestHandler_SignRegister(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "", "", "", false)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/substitutions/sub1/sign-register", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("forbidden_role", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "u1", "parent", "s1", false)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/substitutions/sub1/sign-register", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})

	t.Run("success_and_service_error", func(t *testing.T) {
		subTeacher := "t-sub"
		repo := &mockSubstitutionRepo{
			subs: []*Substitution{
				{
					ID:                  "sub1",
					SubstituteTeacherID: &subTeacher,
					Status:              StatusConfirmed,
				},
			},
		}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "t-sub", "teacher", "s1", false)

		body := map[string]string{"notes": "Lezione svolta regolarmente"}
		b, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/substitutions/sub1/sign-register", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		// Wrong teacher
		r2 := setupTestSubstitutionsRouter(h, "t-other", "teacher", "s1", false)
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest(http.MethodPost, "/api/v1/substitutions/sub1/sign-register", bytes.NewBuffer(b))
		r2.ServeHTTP(w2, req2)
		if w2.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w2.Code)
		}
	})
}

func TestHandler_RecommendSubstitutes(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "", "", "", false)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/substitutions/recommend-substitutes", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("forbidden_non_staff", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "u1", "student", "s1", false)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/substitutions/recommend-substitutes", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})

	t.Run("success_and_error", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "u1", "admin", "s1", false)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/substitutions/recommend-substitutes?class_id=c1&subject_id=sub1&date=2026-03-01&hour=2", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		repo.getErr = errors.New("failed getting available teachers")
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req)
		if w2.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w2.Code)
		}
	})
}

func TestHandler_TodaySummary(t *testing.T) {
	t.Run("unauthorized", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "", "", "", false)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/substitutions/today-summary", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", w.Code)
		}
	})

	t.Run("forbidden", func(t *testing.T) {
		repo := &mockSubstitutionRepo{}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "u1", "student", "s1", false)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/substitutions/today-summary", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", w.Code)
		}
	})

	t.Run("success_summary_counts", func(t *testing.T) {
		repo := &mockSubstitutionRepo{
			subs: []*Substitution{
				{ID: "s1", SchoolID: "sch1", Status: StatusPending},
				{ID: "s2", SchoolID: "sch1", Status: StatusAssigned},
				{ID: "s3", SchoolID: "sch1", Status: StatusConfirmed},
				{ID: "s4", SchoolID: "sch1", Status: StatusCancelled},
			},
		}
		h := NewHandler(NewService(repo))
		r := setupTestSubstitutionsRouter(h, "u1", "collaboratore_ds", "sch1", false)

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/substitutions/today-summary?date=2026-03-01", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed unmarshal: %v", err)
		}

		if resp["total"].(float64) != 4 {
			t.Errorf("expected total 4, got %v", resp["total"])
		}
		if resp["pending"].(float64) != 1 {
			t.Errorf("expected pending 1, got %v", resp["pending"])
		}
		if resp["assigned"].(float64) != 1 {
			t.Errorf("expected assigned 1, got %v", resp["assigned"])
		}
		if resp["confirmed"].(float64) != 1 {
			t.Errorf("expected confirmed 1, got %v", resp["confirmed"])
		}
		if resp["cancelled"].(float64) != 1 {
			t.Errorf("expected cancelled 1, got %v", resp["cancelled"])
		}

		// with error
		repo.getErr = errors.New("db error")
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req)
		if w2.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w2.Code)
		}
	})
}
