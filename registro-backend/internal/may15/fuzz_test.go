package may15

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
)

func FuzzSaveMay15Document(f *testing.F) {
	f.Add("2025/2026", "bozza", "Presentazione", "Continuità didattica")
	f.Add("2026/2027", "pubblicato", "", "")
	f.Add("", "approvato_cdc", "CLIL", "PCTO")

	f.Fuzz(func(t *testing.T, academicYear, status, pres, cont string) {
		repo := new(MockMay15Repo)
		svc := NewService(repo)

		req := SaveMay15Request{
			AcademicYear:       academicYear,
			Status:             DocumentStatus(status),
			ClassPresentation:  pres,
			TeachingContinuity: cont,
		}

		repo.On("Upsert", mock.Anything, mock.Anything).Return(nil).Maybe()

		doc, err := svc.SaveDocument(context.Background(), "school-1", "class-5b", req)
		if err != nil {
			t.Fatalf("unexpected error saving may15 document: %v", err)
		}

		if doc == nil {
			t.Fatalf("expected non-nil document when error is nil")
		}

		if academicYear == "" && doc.AcademicYear != "2025/2026" {
			t.Fatalf("expected default academic year '2025/2026', got '%s'", doc.AcademicYear)
		}
	})
}
