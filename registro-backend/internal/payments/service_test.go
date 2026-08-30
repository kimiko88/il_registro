package payments

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, p *SchoolPayment) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*SchoolPayment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SchoolPayment), args.Error(1)
}

func (m *MockRepository) ListByStudent(ctx context.Context, schoolID, studentID string) ([]*SchoolPayment, error) {
	args := m.Called(ctx, schoolID, studentID)
	return args.Get(0).([]*SchoolPayment), args.Error(1)
}

func (m *MockRepository) ListByParent(ctx context.Context, schoolID, parentUserID string) ([]*SchoolPayment, error) {
	args := m.Called(ctx, schoolID, parentUserID)
	return args.Get(0).([]*SchoolPayment), args.Error(1)
}

func (m *MockRepository) ListBySchool(ctx context.Context, schoolID string) ([]*SchoolPayment, error) {
	args := m.Called(ctx, schoolID)
	return args.Get(0).([]*SchoolPayment), args.Error(1)
}

func (m *MockRepository) Pay(ctx context.Context, id, payerUserID, method, transactionID, receiptNumber string) error {
	args := m.Called(ctx, id, payerUserID, method, transactionID, receiptNumber)
	return args.Error(0)
}

func TestPaymentsService_List(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	t.Run("Parent list payments", func(t *testing.T) {
		mockRepo.On("ListByParent", mock.Anything, "school-1", "parent-1").Return([]*SchoolPayment{
			{ID: "pay-1", Title: "Assicurazione", Amount: 10.0, Status: "pending", DueDate: time.Now()},
			{ID: "pay-2", Title: "Gita", Amount: 50.0, Status: "paid", DueDate: time.Now()},
		}, nil).Once()

		res, err := svc.List(context.Background(), "parent-1", "parent", "school-1", "")
		assert.NoError(t, err)
		assert.Equal(t, 10.0, res.PendingTotal)
		assert.Equal(t, 1, res.PendingCount)
		assert.Equal(t, 1, res.PaidCount)
		assert.Len(t, res.Payments, 2)
	})

	t.Run("Unauthorized actor", func(t *testing.T) {
		_, err := svc.List(context.Background(), "", "parent", "school-1", "")
		assert.Error(t, err)
	})
}

func TestPaymentsService_Pay(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	t.Run("Pay pending item by student", func(t *testing.T) {
		payment := &SchoolPayment{ID: "pay-1", StudentID: "student-1", Title: "Assicurazione", Amount: 10.0, Status: "pending"}
		paidPayment := &SchoolPayment{ID: "pay-1", StudentID: "student-1", Title: "Assicurazione", Amount: 10.0, Status: "paid", ReceiptNumber: "REC-1"}

		mockRepo.On("GetByID", mock.Anything, "pay-1").Return(payment, nil).Once()
		mockRepo.On("Pay", mock.Anything, "pay-1", "student-1", "PagoPA", mock.Anything, mock.Anything).Return(nil).Once()
		mockRepo.On("GetByID", mock.Anything, "pay-1").Return(paidPayment, nil).Once()

		res, err := svc.Pay(context.Background(), "student-1", "student", "pay-1", "PagoPA")
		assert.NoError(t, err)
		assert.Equal(t, "paid", res.Status)
	})

	t.Run("Pay pending item by other student forbidden", func(t *testing.T) {
		payment := &SchoolPayment{ID: "pay-1", StudentID: "student-1", Title: "Assicurazione", Amount: 10.0, Status: "pending"}
		mockRepo.On("GetByID", mock.Anything, "pay-1").Return(payment, nil).Once()

		_, err := svc.Pay(context.Background(), "student-2", "student", "pay-1", "PagoPA")
		assert.ErrorIs(t, err, ErrUnauthorized)
	})
}

func TestPaymentsService_Create(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	t.Run("Create valid payment", func(t *testing.T) {
		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *SchoolPayment) bool {
			return p.Amount == 25.5 && p.SchoolID == "school-1"
		})).Return(nil).Once()

		req := CreatePaymentRequest{
			StudentID:   "student-1",
			Title:       "Assicurazione Scolastica",
			Description: "Quota annuale",
			Amount:      25.5,
			DueDate:     "2026-10-31",
		}
		p, err := svc.Create(context.Background(), "admin-1", "admin", "school-1", req)
		assert.NoError(t, err)
		assert.NotNil(t, p)
		assert.Equal(t, 25.5, p.Amount)
	})

	t.Run("Create payment with non-positive amount fails", func(t *testing.T) {
		req := CreatePaymentRequest{
			StudentID: "student-1",
			Title:     "Importo non valido",
			Amount:    -10.0,
			DueDate:   "2026-10-31",
		}
		_, err := svc.Create(context.Background(), "admin-1", "admin", "school-1", req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "maggiore di zero")
	})

	t.Run("Create payment by unauthorized role fails", func(t *testing.T) {
		req := CreatePaymentRequest{
			StudentID: "student-1",
			Title:     "Quota",
			Amount:    15.0,
		}
		_, err := svc.Create(context.Background(), "student-1", "student", "school-1", req)
		assert.ErrorIs(t, err, ErrUnauthorized)
	})
}

func TestPaymentsService_GetByID(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	t.Run("GetByID student success", func(t *testing.T) {
		payment := &SchoolPayment{ID: "p-1", StudentID: "s-1", SchoolID: "sch-1", Amount: 50.0}
		mockRepo.On("GetByID", mock.Anything, "p-1").Return(payment, nil).Once()

		res, err := svc.GetByID(context.Background(), "s-1", "student", "sch-1", "p-1")
		assert.NoError(t, err)
		assert.Equal(t, payment, res)
	})

	t.Run("GetByID other student unauthorized", func(t *testing.T) {
		payment := &SchoolPayment{ID: "p-1", StudentID: "s-1", SchoolID: "sch-1", Amount: 50.0}
		mockRepo.On("GetByID", mock.Anything, "p-1").Return(payment, nil).Once()

		_, err := svc.GetByID(context.Background(), "s-2", "student", "sch-1", "p-1")
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("GetByID unauthorized role rejected", func(t *testing.T) {
		payment := &SchoolPayment{ID: "p-1", StudentID: "s-1", SchoolID: "sch-1", Amount: 50.0}
		mockRepo.On("GetByID", mock.Anything, "p-1").Return(payment, nil).Once()

		_, err := svc.GetByID(context.Background(), "t-1", "teacher", "sch-1", "p-1")
		assert.ErrorIs(t, err, ErrUnauthorized)
	})

	t.Run("GetByID wrong school rejected", func(t *testing.T) {
		payment := &SchoolPayment{ID: "p-1", StudentID: "s-1", SchoolID: "sch-1", Amount: 50.0}
		mockRepo.On("GetByID", mock.Anything, "p-1").Return(payment, nil).Once()

		_, err := svc.GetByID(context.Background(), "adm-1", "admin", "sch-2", "p-1")
		assert.ErrorIs(t, err, ErrUnauthorized)
	})
}

func TestPaymentsService_Pay_AlreadyPaid(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := NewService(mockRepo)

	payment := &SchoolPayment{ID: "p-1", StudentID: "s-1", Status: "paid"}
	mockRepo.On("GetByID", mock.Anything, "p-1").Return(payment, nil).Once()

	_, err := svc.Pay(context.Background(), "s-1", "student", "p-1", "PagoPA")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "già stato pagato")
}
