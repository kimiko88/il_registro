package personnel_desk_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"registro-backend/internal/personnel_desk"
)

type mockDeskService struct {
	createResult *personnel_desk.DeskRequest
	createErr    error
	getResult    *personnel_desk.DeskRequest
	getErr       error
	listResult   []*personnel_desk.DeskRequest
	listErr      error
	submitResult *personnel_desk.DeskRequest
	submitErr    error
	aaResult     *personnel_desk.DeskRequest
	aaErr        error
	dsgaResult   *personnel_desk.DeskRequest
	dsgaErr      error
	dsResult     *personnel_desk.DeskRequest
	dsErr        error
	deleteErr    error
}

func (m *mockDeskService) CreateRequest(ctx context.Context, schoolID, applicantID string, input personnel_desk.CreateDeskRequestInput) (*personnel_desk.DeskRequest, error) {
	return m.createResult, m.createErr
}

func (m *mockDeskService) GetRequest(ctx context.Context, schoolID, id, currentUserID, currentUserRole string) (*personnel_desk.DeskRequest, error) {
	return m.getResult, m.getErr
}

func (m *mockDeskService) ListRequests(ctx context.Context, schoolID, currentUserID, currentUserRole, filterStatus string) ([]*personnel_desk.DeskRequest, error) {
	return m.listResult, m.listErr
}

func (m *mockDeskService) SubmitRequest(ctx context.Context, schoolID, id, currentUserID string) (*personnel_desk.DeskRequest, error) {
	return m.submitResult, m.submitErr
}

func (m *mockDeskService) AAReview(ctx context.Context, schoolID, id, reviewerID string, input personnel_desk.AAReviewInput) (*personnel_desk.DeskRequest, error) {
	return m.aaResult, m.aaErr
}

func (m *mockDeskService) DSGASign(ctx context.Context, schoolID, id, signerID string, input personnel_desk.DSGASignInput) (*personnel_desk.DeskRequest, error) {
	return m.dsgaResult, m.dsgaErr
}

func (m *mockDeskService) DSApprove(ctx context.Context, schoolID, id, approverID string, input personnel_desk.DSApproveInput) (*personnel_desk.DeskRequest, error) {
	return m.dsResult, m.dsErr
}

func (m *mockDeskService) DeleteRequest(ctx context.Context, schoolID, id, currentUserID string) error {
	return m.deleteErr
}

func setupDeskRouter(h *personnel_desk.Handler, userID, schoolID, role string) *gin.Engine {
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

func TestDeskHandler_Auth(t *testing.T) {
	mockSvc := &mockDeskService{}
	h := personnel_desk.NewHandler(mockSvc)

	rNoUser := setupDeskRouter(h, "", "school-1", "teacher")
	req, _ := http.NewRequest("GET", "/api/v1/personnel-desk/requests", nil)
	w := httptest.NewRecorder()
	rNoUser.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestDeskHandler_CreateAndGet(t *testing.T) {
	mockSvc := &mockDeskService{
		createResult: &personnel_desk.DeskRequest{
			ID:          "req-1",
			SchoolID:    "school-1",
			ApplicantID: "u-1",
			Category:    personnel_desk.CatPermessoBreve,
			Status:      personnel_desk.StatusDraft,
			CreatedAt:   time.Now(),
		},
		getResult: &personnel_desk.DeskRequest{
			ID:          "req-1",
			SchoolID:    "school-1",
			ApplicantID: "u-1",
			Category:    personnel_desk.CatPermessoBreve,
			Status:      personnel_desk.StatusDraft,
		},
		listResult: []*personnel_desk.DeskRequest{
			{ID: "req-1", ApplicantID: "u-1", Category: personnel_desk.CatPermessoBreve},
		},
	}
	h := personnel_desk.NewHandler(mockSvc)
	router := setupDeskRouter(h, "u-1", "school-1", "teacher")

	// 1. CreateRequest invalid -> 400
	req, _ := http.NewRequest("POST", "/api/v1/personnel-desk/requests", bytes.NewBufferString("{invalid"))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	// 2. CreateRequest success -> 201
	body, _ := json.Marshal(personnel_desk.CreateDeskRequestInput{
		Category:    personnel_desk.CatPermessoBreve,
		StartDate:   "2026-09-15",
		EndDate:     "2026-09-15",
		Hours:       2,
		Description: "Permesso breve motivi personali",
	})
	req, _ = http.NewRequest("POST", "/api/v1/personnel-desk/requests", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	// 3. ListRequests -> 200
	req, _ = http.NewRequest("GET", "/api/v1/personnel-desk/requests?status=draft", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 4. GetRequest -> 200
	req, _ = http.NewRequest("GET", "/api/v1/personnel-desk/requests/req-1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 5. GetRequest not found -> 404
	mockSvc.getErr = errors.New("richiesta non trovata")
	req, _ = http.NewRequest("GET", "/api/v1/personnel-desk/requests/req-non-existent", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeskHandler_WorkflowSteps(t *testing.T) {
	mockSvc := &mockDeskService{
		submitResult: &personnel_desk.DeskRequest{ID: "req-1", Status: personnel_desk.StatusSubmitted},
		aaResult:     &personnel_desk.DeskRequest{ID: "req-1", Status: personnel_desk.StatusDSGAReview},
		dsgaResult:   &personnel_desk.DeskRequest{ID: "req-1", Status: personnel_desk.StatusDSReview},
		dsResult:     &personnel_desk.DeskRequest{ID: "req-1", Status: personnel_desk.StatusApproved},
	}
	h := personnel_desk.NewHandler(mockSvc)

	// 1. SubmitRequest
	rTeacher := setupDeskRouter(h, "u-1", "school-1", "teacher")
	req, _ := http.NewRequest("PATCH", "/api/v1/personnel-desk/requests/req-1/submit", nil)
	w := httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 2. AAReview: teacher unauthorized -> 403
	aaBody, _ := json.Marshal(personnel_desk.AAReviewInput{Note: "Verificato", Approve: true})
	req, _ = http.NewRequest("PATCH", "/api/v1/personnel-desk/requests/req-1/aa-review", bytes.NewBuffer(aaBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}

	// AAReview: AA authorized -> 200
	rAA := setupDeskRouter(h, "aa-1", "school-1", "assistente_amministrativo")
	req, _ = http.NewRequest("PATCH", "/api/v1/personnel-desk/requests/req-1/aa-review", bytes.NewBuffer(aaBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAA.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 3. DSGASign: AA unauthorized -> 403
	dsgaBody, _ := json.Marshal(personnel_desk.DSGASignInput{Note: "Visto", Approve: true})
	req, _ = http.NewRequest("PATCH", "/api/v1/personnel-desk/requests/req-1/dsga-sign", bytes.NewBuffer(dsgaBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rAA.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}

	// DSGASign: DSGA authorized -> 200
	rDSGA := setupDeskRouter(h, "dsga-1", "school-1", "dsga")
	req, _ = http.NewRequest("PATCH", "/api/v1/personnel-desk/requests/req-1/dsga-sign", bytes.NewBuffer(dsgaBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rDSGA.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 4. DSApprove: DSGA unauthorized -> 403
	dsBody, _ := json.Marshal(personnel_desk.DSApproveInput{DecreeNum: "DEC/2026/01", Note: "Autorizzato", Approve: true})
	req, _ = http.NewRequest("PATCH", "/api/v1/personnel-desk/requests/req-1/ds-approve", bytes.NewBuffer(dsBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rDSGA.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}

	// DSApprove: Principal authorized -> 200
	rDS := setupDeskRouter(h, "ds-1", "school-1", "principal")
	req, _ = http.NewRequest("PATCH", "/api/v1/personnel-desk/requests/req-1/ds-approve", bytes.NewBuffer(dsBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	rDS.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 5. DeleteRequest -> 200
	req, _ = http.NewRequest("DELETE", "/api/v1/personnel-desk/requests/req-1", nil)
	w = httptest.NewRecorder()
	rTeacher.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
