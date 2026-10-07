package meals

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMealsAndWalletLifecycle(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	// 1. Register Special Diet for Student 1 (Celiachia con certificato medico)
	diet, err := svc.RegisterSpecialDiet(ctx, "school-1", SpecialDietRequest{
		StudentID:             "std-1",
		DietCategory:          "sanitaria",
		SpecificDiet:          "celiachia",
		MedicalCertificateURL: "https://scuola.edu.it/cert_celiachia.pdf",
		CertificateExpiryDate: "2027-06-30",
		AllergensList:         "Glutine",
		Notes:                 "Richiede kit pasto sigillato sterile",
	})
	require.NoError(t, err)
	assert.True(t, diet.IsApproved)
	assert.Equal(t, "celiachia", diet.SpecificDiet)

	// 2. Initial Wallet Balance Check & Top-Up via PagoPA
	wallet, err := svc.GetWallet(ctx, "std-1")
	require.NoError(t, err)
	assert.Equal(t, 0.0, wallet.Balance)

	wallet, err = svc.TopUpWallet(ctx, "std-1", 55.0, "00100000000000142")
	require.NoError(t, err)
	assert.Equal(t, 55.0, wallet.Balance)

	// 3. Roll call: 2 students present (std-1 with diet, std-2 standard)
	report, err := svc.RecordRollCallAndAggregate(ctx, "school-1", MealRollCallRequest{
		ClassID:    "class-3B",
		Date:       "2026-10-15",
		PresentIDs: []string{"std-1", "std-2"},
	})
	require.NoError(t, err)
	assert.Equal(t, 2, report.TotalMeals)
	assert.Equal(t, 1, report.TotalStandard)
	assert.Equal(t, 1, report.TotalDiets)
	assert.Equal(t, 1, report.DietsBreakdown["celiachia"])

	// 4. Verify Wallet Debited
	walletAfter, err := svc.GetWallet(ctx, "std-1")
	require.NoError(t, err)
	assert.Equal(t, 49.50, walletAfter.Balance) // 55.0 - 5.50
}
