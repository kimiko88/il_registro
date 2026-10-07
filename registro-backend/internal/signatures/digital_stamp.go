package signatures

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrTamperedStamp      = errors.New("il timbro digitale di sicurezza non corrisponde: documento contraffatto o manomesso (art. 23 CAD)")
	ErrInvalidStampFormat = errors.New("formato del contrassegno glifo non valido")
	ErrInvalidStampKey    = errors.New("chiave di sigillo elettronico non valida")
)

// GenerateDigitalStamp generates a high-density 2D barcode / QR cryptographic glifo as per Art. 23 CAD.
func GenerateDigitalStamp(baseURL, docID, docType, docSHA256, signerName, signerRole string, secretKey []byte) (*DigitalStamp, error) {
	if len(secretKey) == 0 {
		return nil, ErrInvalidStampKey
	}
	if docID == "" || docSHA256 == "" {
		return nil, errors.New("document_id e document_sha256 sono obbligatori")
	}

	now := time.Now().UTC()
	timestampStr := strconv.FormatInt(now.Unix(), 10)

	// Short hash (first 16 chars of hex sha256)
	shortHash := docSHA256
	if len(shortHash) > 16 {
		shortHash = shortHash[:16]
	}

	// Payload data to sign with HMAC-SHA256
	dataToSign := fmt.Sprintf("%s|%s|%s|%s|%s", docID, shortHash, signerRole, timestampStr, docSHA256)
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(dataToSign))
	sigHex := hex.EncodeToString(mac.Sum(nil))

	glifoToken := fmt.Sprintf("%s-%s-%s", shortHash, timestampStr, sigHex[:16])
	qrPayload := fmt.Sprintf("CAD-GLIFO|01|%s|%s|%s|%s|%s", docID, shortHash, signerRole, timestampStr, sigHex)
	verifyURL := fmt.Sprintf("%s/public/verifica-glifo/%s", strings.TrimRight(baseURL, "/"), glifoToken)

	return &DigitalStamp{
		DocumentID:      docID,
		DocumentType:    docType,
		DocumentSHA256:  docSHA256,
		GlifoToken:      glifoToken,
		QRCodePayload:   qrPayload,
		VerificationURL: verifyURL,
		SignerName:      signerName,
		SignerRole:      signerRole,
		SignedAt:        now,
		HMACSignature:   sigHex,
	}, nil
}

// VerifyDigitalStamp verifies the authenticity and integrity of a digital stamp token / payload.
func VerifyDigitalStamp(payload string, originalSHA256 string, secretKey []byte) (*DigitalStampVerification, error) {
	if len(secretKey) == 0 {
		return nil, ErrInvalidStampKey
	}

	parts := strings.Split(payload, "|")
	if len(parts) != 7 || parts[0] != "CAD-GLIFO" || parts[1] != "01" {
		return nil, ErrInvalidStampFormat
	}

	docID := parts[2]
	shortHash := parts[3]
	signerRole := parts[4]
	timestampStr := parts[5]
	expectedSigHex := parts[6]

	ts, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return nil, ErrInvalidStampFormat
	}

	// Recompute HMAC
	dataToSign := fmt.Sprintf("%s|%s|%s|%s|%s", docID, shortHash, signerRole, timestampStr, originalSHA256)
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(dataToSign))
	computedSigHex := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expectedSigHex), []byte(computedSigHex)) {
		return &DigitalStampVerification{
			IsValid:        false,
			DocumentID:     docID,
			LegalReference: "Art. 23 comma 2-bis D.Lgs. 82/2005 (CAD)",
			Message:        "ATTENZIONE: Sigillo elettronico non valido. Il documento cartaceo non è conforme all'originale informatico.",
		}, ErrTamperedStamp
	}

	return &DigitalStampVerification{
		IsValid:        true,
		DocumentID:     docID,
		DocumentSHA256: originalSHA256,
		SignerRole:     signerRole,
		SignedAt:       time.Unix(ts, 0).UTC(),
		LegalReference: "Art. 23 comma 2-bis D.Lgs. 82/2005 (CAD) - Piena conformità all'originale informatico",
		Message:        "Documento cartaceo pienamente conforme all'originale informatico conservato a norma.",
	}, nil
}
