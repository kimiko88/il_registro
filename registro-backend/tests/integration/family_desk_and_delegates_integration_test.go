package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/family_desk"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockFamilyDeskRepoForIntegration struct {
	requests  map[string]*family_desk.FamilyRequest
	delegates map[string]*family_desk.PermanentDelegate
}

func newMockFamilyDeskRepoForIntegration() *mockFamilyDeskRepoForIntegration {
	return &mockFamilyDeskRepoForIntegration{
		requests:  make(map[string]*family_desk.FamilyRequest),
		delegates: make(map[string]*family_desk.PermanentDelegate),
	}
}

func (m *mockFamilyDeskRepoForIntegration) CreateRequest(ctx context.Context, req *family_desk.FamilyRequest) error {
	req.ID = "req-" + time.Now().Format("150405.000000")
	req.Status = family_desk.StatusSubmitted
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	m.requests[req.ID] = req
	return nil
}

func (m *mockFamilyDeskRepoForIntegration) GetRequest(ctx context.Context, id string) (*family_desk.FamilyRequest, error) {
	req, ok := m.requests[id]
	if !ok {
		return nil, assert.AnError
	}
	return req, nil
}

func (m *mockFamilyDeskRepoForIntegration) ListRequests(ctx context.Context, schoolID, parentID, status string) ([]family_desk.FamilyRequest, error) {
	var list []family_desk.FamilyRequest
	for _, req := range m.requests {
		if schoolID != "" && req.SchoolID != schoolID {
			continue
		}
		if parentID != "" && req.ParentID != parentID {
			continue
		}
		if status != "" && req.Status != status {
			continue
		}
		list = append(list, *req)
	}
	return list, nil
}

func (m *mockFamilyDeskRepoForIntegration) UpdateRequestStatus(ctx context.Context, id, status, rejectionReason, protocolNumber, reviewerID string) error {
	req, ok := m.requests[id]
	if !ok {
		return assert.AnError
	}
	req.Status = status
	if rejectionReason != "" {
		req.RejectionReason = rejectionReason
	}
	if protocolNumber != "" {
		req.ProtocolNumber = protocolNumber
	}
	req.ReviewedBy = &reviewerID
	now := time.Now()
	req.ReviewedAt = &now
	req.UpdatedAt = now
	return nil
}

func (m *mockFamilyDeskRepoForIntegration) CreatePermanentDelegate(ctx context.Context, delegate *family_desk.PermanentDelegate) error {
	delegate.ID = "del-" + time.Now().Format("150405.000000")
	delegate.CreatedAt = time.Now()
	m.delegates[delegate.ID] = delegate
	return nil
}

func (m *mockFamilyDeskRepoForIntegration) ListPermanentDelegates(ctx context.Context, studentID, schoolID string) ([]family_desk.PermanentDelegate, error) {
	var list []family_desk.PermanentDelegate
	for _, d := range m.delegates {
		if studentID != "" && d.StudentID != studentID {
			continue
		}
		if schoolID != "" && d.SchoolID != schoolID {
			continue
		}
		list = append(list, *d)
	}
	return list, nil
}

func setupFamilyDeskTestServer(repo family_desk.Repository, role, userID, schoolID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("role", role)
		c.Set("user_id", userID)
		c.Set("school_id", schoolID)
		c.Next()
	})

	svc := family_desk.NewService(repo)
	handler := family_desk.NewHandler(svc)
	rg := r.Group("/api/v1")
	handler.RegisterRoutes(rg)

	return r
}

func TestIntegration_FamilyDesk_PermanentDelegates_Lifecycle(t *testing.T) {
	repo := newMockFamilyDeskRepoForIntegration()
	parentServer := setupFamilyDeskTestServer(repo, "parent", "parent-101", "school-001")
	secretaryServer := setupFamilyDeskTestServer(repo, "secretary", "secr-201", "school-001")
	studentServer := setupFamilyDeskTestServer(repo, "student", "stud-301", "school-001")

	var createdRequestID string

	// Step 1: Parent submits a delegate request
	t.Run("Parent submits valid delega_ritiro request", func(t *testing.T) {
		payload := map[string]interface{}{
			"student_id":   "student-555",
			"request_type": family_desk.TypeDelegaRitiro,
			"form_data": map[string]interface{}{
				"first_name":       "Giuseppe",
				"last_name":        "Verdi",
				"tax_code":         "VRDGPP60A01H501U",
				"relationship":     "nonno",
				"phone":            "+393331234567",
				"id_card_details":  "CI CA12345AA rilasciata da Roma",
				"id_card_file_url": "https://storage.school.it/idcards/verdi_ci.pdf",
			},
			"attachment_urls": []string{"https://storage.school.it/idcards/verdi_ci.pdf"},
		}
		body, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/family-desk/requests", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		parentServer.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Message string                    `json:"message"`
			Request family_desk.FamilyRequest `json:"request"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "Istanza presentata con successo", resp.Message)
		assert.NotEmpty(t, resp.Request.ID)
		assert.Equal(t, family_desk.StatusSubmitted, resp.Request.Status)
		createdRequestID = resp.Request.ID
	})

	// Step 2: Parent lists their requests and sees it submitted
	t.Run("Parent lists their requests and filters by submitted status", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/family-desk/requests?status=submitted", nil)
		parentServer.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp struct {
			Data []family_desk.FamilyRequest `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Len(t, resp.Data, 1)
		assert.Equal(t, createdRequestID, resp.Data[0].ID)
		assert.Equal(t, family_desk.TypeDelegaRitiro, resp.Data[0].RequestType)
	})

	// Step 3: Unauthorized user (student) attempts to review request -> 403 Forbidden
	t.Run("Unauthorized role (student) receives 403 on review", func(t *testing.T) {
		reviewBody := map[string]interface{}{
			"status":          "approved",
			"protocol_number": "PROT/2026/00142",
		}
		body, _ := json.Marshal(reviewBody)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/api/v1/family-desk/requests/"+createdRequestID+"/review", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		studentServer.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	// Step 4: Secretary reviews and approves request -> auto-creates permanent delegate
	t.Run("Secretary approves request with protocol number", func(t *testing.T) {
		reviewBody := map[string]interface{}{
			"status":          family_desk.StatusApproved,
			"protocol_number": "PROT/2026/00142",
		}
		body, _ := json.Marshal(reviewBody)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/api/v1/family-desk/requests/"+createdRequestID+"/review", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		secretaryServer.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Message string `json:"message"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Contains(t, resp.Message, "aggiornato con successo")
	})

	// Step 5: Listing permanent delegates returns the new delegate
	t.Run("Listing delegates for student verifies newly synced permanent delegate", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/family-desk/delegates?student_id=student-555", nil)
		secretaryServer.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Data []family_desk.PermanentDelegate `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Len(t, resp.Data, 1)

		del := resp.Data[0]
		assert.Equal(t, "Giuseppe", del.FirstName)
		assert.Equal(t, "Verdi", del.LastName)
		assert.Equal(t, "VRDGPP60A01H501U", del.TaxCode)
		assert.Equal(t, "nonno", del.Relationship)
		assert.Equal(t, "+393331234567", del.Phone)
		assert.True(t, del.IsValid)
		assert.NotNil(t, del.ApprovedRequestID)
		assert.Equal(t, createdRequestID, *del.ApprovedRequestID)
	})

	// Step 6: Parent submits autonomous exit for under-14 (legge 172/2017)
	t.Run("Parent submits valid autonomous exit under 14 with explicit consent", func(t *testing.T) {
		payload := map[string]interface{}{
			"student_id":   "student-555",
			"request_type": family_desk.TypeUscitaAutonomaUnder14,
			"form_data": map[string]interface{}{
				"consent_given":   true,
				"emergency_phone": "+393339998877",
			},
		}
		body, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/family-desk/requests", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		parentServer.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	// Step 7: Under 14 exit without consent is rejected with clear error
	t.Run("Under 14 exit without consent fails validation", func(t *testing.T) {
		payload := map[string]interface{}{
			"student_id":   "student-555",
			"request_type": family_desk.TypeUscitaAutonomaUnder14,
			"form_data": map[string]interface{}{
				"consent_given": false,
			},
		}
		body, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/family-desk/requests", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		parentServer.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Step 8: Request with invalid request type is rejected
	t.Run("Reject unknown request type", func(t *testing.T) {
		payload := map[string]interface{}{
			"student_id":   "student-555",
			"request_type": "tipo_inesistente_xyz",
			"form_data":    map[string]interface{}{"foo": "bar"},
		}
		body, _ := json.Marshal(payload)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/family-desk/requests", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		parentServer.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
