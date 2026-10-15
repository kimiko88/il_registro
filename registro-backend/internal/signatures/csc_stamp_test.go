package signatures

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCSCBatchSigning(t *testing.T) {
	connector := NewCSCConnector(&SoftwareQSCD{}, nil)
	ctx := context.Background()

	t.Run("Valid batch signing of 5 documents", func(t *testing.T) {
		docs := []string{"doc-1", "doc-2", "doc-3", "doc-4", "doc-5"}
		resp, err := connector.BatchSign(ctx, "user-ds", "school-1", CSCBatchSignRequest{
			DocumentIDs: docs,
			PIN:         "1234",
			OTP:         "654321",
			SignerRole:  "Dirigente Scolastico",
		})
		require.NoError(t, err)
		assert.Equal(t, 5, resp.TotalRequested)
		assert.Equal(t, 5, resp.TotalSigned)
		assert.Len(t, resp.Signatures, 5)
		for _, s := range resp.Signatures {
			assert.NotEmpty(t, s.SignatureValue)
			assert.NotEmpty(t, s.PublicKeyPEM)
			assert.NotEmpty(t, s.TimestampToken)
			assert.True(t, s.IsValid)
		}
	})

	t.Run("Invalid PIN or OTP fails", func(t *testing.T) {
		_, err := connector.BatchSign(ctx, "user-ds", "school-1", CSCBatchSignRequest{
			DocumentIDs: []string{"doc-1"},
			PIN:         "12",  // too short
			OTP:         "123", // not 6 digits
		})
		assert.ErrorIs(t, err, ErrInvalidPINorOTP)
	})

	t.Run("Empty documents list fails", func(t *testing.T) {
		_, err := connector.BatchSign(ctx, "user-ds", "school-1", CSCBatchSignRequest{
			DocumentIDs: []string{},
			PIN:         "1234",
			OTP:         "123456",
		})
		assert.ErrorIs(t, err, ErrEmptyBatchList)
	})
}

func TestDigitalStampCAD(t *testing.T) {
	secretKey := []byte("secret-electronic-seal-key-1234")
	baseURL := "https://scuola.edu.it"
	docID := "pagella-std-42"
	docType := "pagella_fine_anno"
	docSHA256 := "a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e"

	stamp, err := GenerateDigitalStamp(baseURL, docID, docType, docSHA256, "Prof. Marco Bianchi", "Dirigente Scolastico", secretKey)
	require.NoError(t, err)
	assert.NotEmpty(t, stamp.QRCodePayload)
	assert.Contains(t, stamp.VerificationURL, "/public/verifica-glifo/")
	assert.NotEmpty(t, stamp.HMACSignature)

	t.Run("Verify authentic digital stamp", func(t *testing.T) {
		ver, err := VerifyDigitalStamp(stamp.QRCodePayload, docSHA256, secretKey)
		require.NoError(t, err)
		assert.True(t, ver.IsValid)
		assert.Equal(t, docID, ver.DocumentID)
		assert.Contains(t, ver.LegalReference, "Art. 23")
	})

	t.Run("Verify tampered document SHA256", func(t *testing.T) {
		tamperedSHA := "0000000000000000000000000000000000000000000000000000000000000000"
		_, err := VerifyDigitalStamp(stamp.QRCodePayload, tamperedSHA, secretKey)
		assert.ErrorIs(t, err, ErrTamperedStamp)
	})

	t.Run("Verify tampered payload string", func(t *testing.T) {
		tamperedPayload := stamp.QRCodePayload[:len(stamp.QRCodePayload)-4] + "ffff"
		_, err := VerifyDigitalStamp(tamperedPayload, docSHA256, secretKey)
		assert.ErrorIs(t, err, ErrTamperedStamp)
	})
}
