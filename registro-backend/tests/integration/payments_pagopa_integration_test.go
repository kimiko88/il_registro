package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"registro-backend/internal/payments"
	"registro-backend/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockPaymentsRepo struct {
	payments map[string]*payments.SchoolPayment
}

func newMockPaymentsRepo() *mockPaymentsRepo {
	return &mockPaymentsRepo{
		payments: make(map[string]*payments.SchoolPayment),
	}
}

func (m *mockPaymentsRepo) Create(ctx context.Context, p *payments.SchoolPayment) error {
	if p.ID == "" {
		p.ID = "pay-" + time.Now().Format("150405.000000")
	}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	p.Status = "pending"
	m.payments[p.ID] = p
	return nil
}

func (m *mockPaymentsRepo) GetByID(ctx context.Context, id string) (*payments.SchoolPayment, error) {
	p, ok := m.payments[id]
	if !ok {
		return nil, payments.ErrPaymentNotFound
	}
	return p, nil
}

func (m *mockPaymentsRepo) ListByParent(ctx context.Context, schoolID, parentID string) ([]*payments.SchoolPayment, error) {
	var res []*payments.SchoolPayment
	for _, p := range m.payments {
		res = append(res, p)
	}
	return res, nil
}

func (m *mockPaymentsRepo) ListByStudent(ctx context.Context, schoolID, studentID string) ([]*payments.SchoolPayment, error) {
	var res []*payments.SchoolPayment
	for _, p := range m.payments {
		if p.StudentID == studentID {
			res = append(res, p)
		}
	}
	return res, nil
}

func (m *mockPaymentsRepo) ListBySchool(ctx context.Context, schoolID string) ([]*payments.SchoolPayment, error) {
	var res []*payments.SchoolPayment
	for _, p := range m.payments {
		if p.SchoolID == schoolID {
			res = append(res, p)
		}
	}
	return res, nil
}

func (m *mockPaymentsRepo) Pay(ctx context.Context, id, payerUserID, method, transactionID, receiptNumber string) error {
	if p, ok := m.payments[id]; ok {
		now := time.Now()
		p.Status = "paid"
		p.PaidAt = &now
		p.PaymentMethod = method
		p.TransactionID = transactionID
		p.PayerUserID = &payerUserID
		p.ReceiptNumber = receiptNumber
		p.UpdatedAt = now
	}
	return nil
}

type mockUserRepoForPayments struct {
	users.Repository
}

func (m *mockUserRepoForPayments) IsGuardian(ctx context.Context, parentID, studentID string) (bool, error) {
	return true, nil
}

func setupPaymentsTestRouter(repo payments.Repository, uRepo users.Repository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	svc := payments.NewService(repo, uRepo)
	handler := payments.NewHandler(svc)

	api := r.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		role := c.GetHeader("X-Role")
		if role == "" {
			role = "parent"
		}
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = "parent-1"
		}
		c.Set("role", role)
		c.Set("user_id", userID)
		c.Set("school_id", "school-1")
		c.Next()
	})
	handler.RegisterRoutes(api)
	return r
}

func TestIntegration_Payments_PagoPA_Workflow(t *testing.T) {
	repo := newMockPaymentsRepo()
	uRepo := &mockUserRepoForPayments{}
	r := setupPaymentsTestRouter(repo, uRepo)

	// 1. Create a Payment Notice (Admin / Secretary creates notice for School Trip)
	createPayload := payments.CreatePaymentRequest{
		StudentID:   "student-1",
		Title:       "Gita Istituzionale - Visita Musei Vaticani Roma",
		Description: "Quota viaggio in pullman GT e biglietto d'ingresso guidato",
		Amount:      45.50,
		DueDate:     "2026-04-15",
	}
	body, _ := json.Marshal(createPayload)
	reqCreate, _ := http.NewRequest("POST", "/api/v1/payments", bytes.NewReader(body))
	reqCreate.Header.Set("Content-Type", "application/json")
	reqCreate.Header.Set("X-Role", "secretary")
	reqCreate.Header.Set("X-User-ID", "secr-1")
	wCreate := httptest.NewRecorder()
	r.ServeHTTP(wCreate, reqCreate)
	require.Equal(t, http.StatusCreated, wCreate.Code)

	var createdPay payments.SchoolPayment
	err := json.Unmarshal(wCreate.Body.Bytes(), &createdPay)
	require.NoError(t, err)
	assert.NotEmpty(t, createdPay.ID)
	assert.Equal(t, 45.50, createdPay.Amount)
	assert.Equal(t, "pending", createdPay.Status)

	// 2. Parent lists payments for child student-1
	reqList, _ := http.NewRequest("GET", "/api/v1/payments?student_id=student-1", nil)
	reqList.Header.Set("X-Role", "parent")
	reqList.Header.Set("X-User-ID", "parent-1")
	wList := httptest.NewRecorder()
	r.ServeHTTP(wList, reqList)
	require.Equal(t, http.StatusOK, wList.Code)

	var summary payments.PaymentSummaryResponse
	err = json.Unmarshal(wList.Body.Bytes(), &summary)
	require.NoError(t, err)
	assert.Equal(t, 1, summary.PendingCount)
	assert.Equal(t, 45.50, summary.PendingTotal)

	// 3. Parent executes PagoPA payment
	payReqBody, _ := json.Marshal(payments.PayRequest{PaymentMethod: "PagoPA"})
	reqPay, _ := http.NewRequest("POST", "/api/v1/payments/"+createdPay.ID+"/pay", bytes.NewReader(payReqBody))
	reqPay.Header.Set("Content-Type", "application/json")
	reqPay.Header.Set("X-Role", "parent")
	reqPay.Header.Set("X-User-ID", "parent-1")
	wPay := httptest.NewRecorder()
	r.ServeHTTP(wPay, reqPay)
	require.Equal(t, http.StatusOK, wPay.Code)

	var payResp struct {
		Message string                  `json:"message"`
		Payment *payments.SchoolPayment `json:"payment"`
	}
	err = json.Unmarshal(wPay.Body.Bytes(), &payResp)
	require.NoError(t, err)
	require.NotNil(t, payResp.Payment)
	assert.Equal(t, "paid", payResp.Payment.Status)
	assert.NotEmpty(t, payResp.Payment.ReceiptNumber)
	assert.NotNil(t, payResp.Payment.PaidAt)

	// 4. Verify updated list shows 0 pending and 1 paid
	wListAfter := httptest.NewRecorder()
	r.ServeHTTP(wListAfter, reqList)
	require.Equal(t, http.StatusOK, wListAfter.Code)

	var summaryAfter payments.PaymentSummaryResponse
	err = json.Unmarshal(wListAfter.Body.Bytes(), &summaryAfter)
	require.NoError(t, err)
	assert.Equal(t, 0, summaryAfter.PendingCount)
	assert.Equal(t, 1, summaryAfter.PaidCount)
}
