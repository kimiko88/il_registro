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

	t.Run("Pay pending item", func(t *testing.T) {
		payment := &SchoolPayment{ID: "pay-1", Title: "Assicurazione", Amount: 10.0, Status: "pending"}
		paidPayment := &SchoolPayment{ID: "pay-1", Title: "Assicurazione", Amount: 10.0, Status: "paid", ReceiptNumber: "REC-1"}

		mockRepo.On("GetByID", mock.Anything, "pay-1").Return(payment, nil).Once()
		mockRepo.On("Pay", mock.Anything, "pay-1", "parent-1", "PagoPA", mock.Anything, mock.Anything).Return(nil).Once()
		mockRepo.On("GetByID", mock.Anything, "pay-1").Return(paidPayment, nil).Once()

		res, err := svc.Pay(context.Background(), "parent-1", "parent", "pay-1", "PagoPA")
		assert.NoError(t, err)
		assert.Equal(t, "paid", res.Status)
	})
}
