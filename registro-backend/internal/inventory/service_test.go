package inventory

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateDepreciation(t *testing.T) {
	acq := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// After 1 year at 20% depreciation: 1000 - 200 = 800
	at1Year := acq.AddDate(1, 0, 0)
	val1 := CalculateDepreciation(1000.0, 20.0, acq, at1Year)
	assert.InDelta(t, 800.0, val1, 1.0)

	// After 5 years at 20%: 0.0 (fully depreciated)
	at5Years := acq.AddDate(5, 0, 0)
	val5 := CalculateDepreciation(1000.0, 20.0, acq, at5Years)
	assert.Equal(t, 0.0, val5)
}

func TestInventoryAndLoanLifecycle(t *testing.T) {
	ctx := context.Background()
	svc := NewService(NewMemoryRepository())

	// 1. Create Asset (Notebook PNRR)
	asset, err := svc.CreateAsset(ctx, "school-1", CreateAssetRequest{
		Category:         CategoryPNRR,
		Description:      "Notebook Lenovo ThinkPad L14 Gen 4",
		SerialNumber:     "PF4A999B",
		BuildingLocation: "Centrale",
		RoomLocation:     "Laboratorio STEAM",
		AcquisitionDate:  "2026-01-10",
		InitialValue:     850.00,
		DepreciationRate: 20.0,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, asset.InventoryNumber)
	assert.Contains(t, asset.BarcodeData, "INV")
	assert.Contains(t, asset.QRData, "INVENTARIO|")

	// 2. Get Barcode Label
	label, err := svc.GetLabel(ctx, asset.ID)
	require.NoError(t, err)
	assert.Equal(t, asset.InventoryNumber, label.InventoryNumber)
	assert.NotEmpty(t, label.QRCodeData)

	// 3. Create Loan Contract for Student
	loan, err := svc.CreateLoanContract(ctx, "school-1", CreateLoanRequest{
		AssetID:              asset.ID,
		StudentID:            "std-1",
		ParentUserID:         "parent-1",
		ExpectedReturnDate:   "2026-06-30",
		DeviceConditionNotes: "Dispositivo nuovo sigillato con alimentatore originale",
	})
	require.NoError(t, err)
	assert.Equal(t, "attivo", loan.Status)
	assert.NotEmpty(t, loan.ContractNumber)

	// 4. Return Device
	returned, err := svc.ReturnLoanDevice(ctx, loan.ID, "Restituito integro e funzionante")
	require.NoError(t, err)
	assert.Equal(t, "restituito", returned.Status)
	assert.NotNil(t, returned.ActualReturnDate)
	assert.Contains(t, returned.DeviceConditionNotes, "Restituzione: Restituito integro")
}
