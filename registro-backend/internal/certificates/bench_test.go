package certificates

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func BenchmarkIsValidCertificateType(b *testing.B) {
	types := []CertificateType{
		CertIscrizione,
		CertFrequenza,
		CertPromozione,
		CertBehavior,
		"invalid_type",
	}

	idx := 0
	for b.Loop() {
		t := types[idx%len(types)]
		idx++
		_ = IsValidCertificateType(t)
	}
}

func BenchmarkListCertificates(b *testing.B) {
	mockRepo := new(MockRepository)
	mockUsers := new(MockUsersRepository)

	svc := NewService(mockRepo, mockUsers)
	ctx := context.Background()

	certs := make([]Certificate, 50)
	for i := 0; i < 50; i++ {
		certs[i] = Certificate{
			ID:           "cert-bench",
			SchoolID:     "school-1",
			StudentID:    "student-1",
			Type:         CertIscrizione,
			IssuedAt:     time.Now(),
			AcademicYear: "2025/2026",
		}
	}

	mockRepo.On("List", mock.Anything, "school-1", "student-1", CertIscrizione, "2025/2026").
		Return(certs, nil)

	for b.Loop() {
		_, _ = svc.ListCertificates(ctx, "school-1", "student-1", CertIscrizione, "2025/2026")
	}
}
