package payments

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateMod97(t *testing.T) {
	t.Run("Valid numeric strings", func(t *testing.T) {
		check, err := CalculateMod97("001000000000001")
		require.NoError(t, err)
		assert.Len(t, check, 2)
	})

	t.Run("Empty string error", func(t *testing.T) {
		_, err := CalculateMod97("")
		assert.Error(t, err)
	})

	t.Run("Non-numeric character error", func(t *testing.T) {
		_, err := CalculateMod97("00100A000000001")
		assert.Error(t, err)
	})
}

func TestGenerateAndValidateIUV(t *testing.T) {
	t.Run("Generate valid IUV", func(t *testing.T) {
		iuv, err := GenerateIUV(0, "12", 123456789)
		require.NoError(t, err)
		assert.Len(t, iuv, 17)
		assert.True(t, ValidateIUV(iuv))
	})

	t.Run("Invalid parameters", func(t *testing.T) {
		_, err := GenerateIUV(5, "12", 100)
		assert.Error(t, err)

		_, err = GenerateIUV(0, "123", 100)
		assert.Error(t, err)

		_, err = GenerateIUV(0, "12", -1)
		assert.Error(t, err)
	})

	t.Run("Validate tampered IUV", func(t *testing.T) {
		iuv, err := GenerateIUV(0, "01", 987654321)
		require.NoError(t, err)
		assert.True(t, ValidateIUV(iuv))

		// Corrupt one digit
		tampered := iuv[:10] + "9" + iuv[11:]
		if tampered != iuv {
			assert.False(t, ValidateIUV(tampered))
		}
	})
}

func TestPagoPAQRCodePayload(t *testing.T) {
	cfEnte := "97123456789"
	iuv := "00100000000000142"
	amountCents := int64(4500)

	payload := GeneratePagoPAQRCodePayload(cfEnte, iuv, amountCents)
	assert.Equal(t, "PAGOPA|002|00100000000000142|97123456789|4500", payload)

	parsedIUV, parsedCF, parsedAmount, err := ParsePagoPAQRCodePayload(payload)
	require.NoError(t, err)
	assert.Equal(t, iuv, parsedIUV)
	assert.Equal(t, cfEnte, parsedCF)
	assert.Equal(t, amountCents, parsedAmount)

	_, _, _, err = ParsePagoPAQRCodePayload("INVALID|FORMAT")
	assert.Error(t, err)
}

func TestGenerateBollettinoNotice(t *testing.T) {
	p := &SchoolPayment{
		ID:          "pay-123",
		StudentID:   "student-456",
		StudentName: "Mario Rossi",
		Title:       "Assicurazione Scolastica",
		Amount:      15.50,
		DueDate:     time.Now().Add(10 * 24 * time.Hour),
		Status:      "pending",
	}

	notice, err := GenerateBollettinoNotice(p, "80012345678", "Liceo Scientifico Leonardo")
	require.NoError(t, err)
	assert.NotEmpty(t, notice.IUV)
	assert.Contains(t, notice.QRCodePayload, "PAGOPA|002|")
	assert.Contains(t, notice.QRCodePayload, "80012345678")
	assert.Equal(t, 15.50, notice.Amount)
	assert.Equal(t, "Mario Rossi", notice.StudentName)
}
