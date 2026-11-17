package family_desk

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/users"

	"github.com/gin-gonic/gin"
)

type mockFamilyRepo struct {
	requests  map[string]*FamilyRequest
	delegates map[string]*PermanentDelegate
}

func newMockFamilyRepo() *mockFamilyRepo {
	return &mockFamilyRepo{
		requests:  make(map[string]*FamilyRequest),
		delegates: make(map[string]*PermanentDelegate),
	}
}

func (m *mockFamilyRepo) CreateRequest(ctx context.Context, req *FamilyRequest) error {
	req.ID = "req-test-1"
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	m.requests[req.ID] = req
	return nil
}

func (m *mockFamilyRepo) GetRequest(ctx context.Context, id string) (*FamilyRequest, error) {
	if r, ok := m.requests[id]; ok {
		return r, nil
	}
	return nil, nil
}

func (m *mockFamilyRepo) ListRequests(ctx context.Context, schoolID, parentID, status string) ([]FamilyRequest, error) {
	var list []FamilyRequest
	for _, r := range m.requests {
		if r.SchoolID == schoolID {
			if (parentID == "" || r.ParentID == parentID) && (status == "" || r.Status == status) {
				list = append(list, *r)
			}
		}
	}
	return list, nil
}

func (m *mockFamilyRepo) UpdateRequestStatus(ctx context.Context, id, status, rejectionReason, protocolNumber, reviewerID string) error {
	if r, ok := m.requests[id]; ok {
		r.Status = status
		r.RejectionReason = rejectionReason
		r.ProtocolNumber = protocolNumber
		r.ReviewedBy = &reviewerID
		now := time.Now()
		r.ReviewedAt = &now
	}
	return nil
}

func (m *mockFamilyRepo) CreatePermanentDelegate(ctx context.Context, delegate *PermanentDelegate) error {
	delegate.ID = "del-test-1"
	delegate.CreatedAt = time.Now()
	m.delegates[delegate.ID] = delegate
	return nil
}

func (m *mockFamilyRepo) ListPermanentDelegates(ctx context.Context, studentID, schoolID string) ([]PermanentDelegate, error) {
	var list []PermanentDelegate
	for _, d := range m.delegates {
		if (studentID == "" || d.StudentID == studentID) && (schoolID == "" || d.SchoolID == schoolID) {
			list = append(list, *d)
		}
	}
	return list, nil
}

type mockFamilyDeskUsers struct {
	users.Repository
	byID map[string]*users.User
}

func newMockFamilyDeskUsers() *mockFamilyDeskUsers {
	schoolID := "school-123"
	return &mockFamilyDeskUsers{byID: map[string]*users.User{
		"parent-456": {ID: "parent-456", Role: "parent", SchoolID: &schoolID},
		"student-789": {ID: "student-789", Role: "student", SchoolID: &schoolID},
		"secretary-1": {ID: "secretary-1", Role: "secretary", SchoolID: &schoolID},
	}}
}

func (m *mockFamilyDeskUsers) GetByID(_ context.Context, id string) (*users.User, error) {
	user, ok := m.byID[id]
	if !ok {
		return nil, users.ErrUserNotFound
	}
	return user, nil
}

func (m *mockFamilyDeskUsers) IsGuardian(_ context.Context, parentID, studentID string) (bool, error) {
	return parentID == "parent-456" && studentID == "student-789", nil
}

func setupFamilyDeskRouter(repo Repository, userRepo users.Repository, role string, userID, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("role", role)
		c.Set("school_id", schoolID)
		c.Set("user_id", userID)
		c.Next()
	})

	svc := NewService(repo)
	handler := NewHandler(svc, userRepo)
	api := r.Group("/api/v1")
	handler.RegisterRoutes(api)

	return r
}

func TestFamilyDeskIntegrationWorkflow(t *testing.T) {
	repo := newMockFamilyRepo()
	userRepo := newMockFamilyDeskUsers()
	parentRouter := setupFamilyDeskRouter(repo, userRepo, "parent", "parent-456", "school-123")
	secretaryRouter := setupFamilyDeskRouter(repo, userRepo, "secretary", "secretary-1", "school-123")

	// Step 1: Parent submits 'delega_ritiro'
	delegaPayload := map[string]interface{}{
		"student_id":   "student-789",
		"request_type": TypeDelegaRitiro,
		"form_data": map[string]interface{}{
			"first_name":      "Giuseppe",
			"last_name":       "Verdi",
			"tax_code":        "VRDGPP60A01H501K",
			"relationship":    "nonno",
			"phone":           "+393339876543",
			"id_card_details": "CI AS98765AA",
		},
	}
	body, _ := json.Marshal(delegaPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/family-desk/requests", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	parentRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on submit request, got %d: %s", w.Code, w.Body.String())
	}

	// Step 2: Secretary approves request with protocol number
	reviewPayload := map[string]interface{}{
		"status":          StatusApproved,
		"protocol_number": "PROT-2026-00452",
	}
	body, _ = json.Marshal(reviewPayload)
	req = httptest.NewRequest(http.MethodPut, "/api/v1/family-desk/requests/req-test-1/review", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	secretaryRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on review request, got %d: %s", w.Code, w.Body.String())
	}

	// Step 3: Check permanent delegates for student (synced automatically)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/family-desk/delegates?student_id=student-789", nil)
	w = httptest.NewRecorder()
	parentRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on list delegates, got %d", w.Code)
	}

	var res struct {
		Data []PermanentDelegate `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	if len(res.Data) != 1 {
		t.Fatalf("expected 1 permanent delegate synced, got %d", len(res.Data))
	}
	if res.Data[0].FirstName != "Giuseppe" || res.Data[0].TaxCode != "VRDGPP60A01H501K" {
		t.Errorf("unexpected delegate data: %+v", res.Data[0])
	}
}
