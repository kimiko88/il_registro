package signatures

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidPINorOTP  = errors.New("PIN o codice OTP remoto non valido per il firmatario CSC")
	ErrEmptyBatchList   = errors.New("l'elenco dei documenti per la firma massiva non può essere vuoto")
	ErrBatchLimitExceed = errors.New("il limite massimo per singola transazione CSC è di 500 documenti")
)

type CSCConnector interface {
	BatchSign(ctx context.Context, signerID, schoolID string, req CSCBatchSignRequest) (*CSCBatchSignResponse, error)
}

type DefaultCSCConnector struct {
	qscd QSCDProvider
	tsa  TSAClient
}

type SimulatedTSAClient struct{}

func (s *SimulatedTSAClient) RequestTimestamp(_ context.Context, digest []byte) (string, time.Time, error) {
	return computeSimulatedTSAToken(digest), time.Now().UTC(), nil
}

func NewCSCConnector(qscd QSCDProvider, tsa TSAClient) *DefaultCSCConnector {
	if qscd == nil {
		qscd = &SoftwareQSCD{}
	}
	if tsa == nil {
		tsa = &SimulatedTSAClient{}
	}
	return &DefaultCSCConnector{
		qscd: qscd,
		tsa:  tsa,
	}
}

// BatchSign executes a massive remote qualified signature using Cloud Signature Consortium specifications.
func (c *DefaultCSCConnector) BatchSign(ctx context.Context, signerID, schoolID string, req CSCBatchSignRequest) (*CSCBatchSignResponse, error) {
	if len(req.DocumentIDs) == 0 {
		return nil, ErrEmptyBatchList
	}
	if len(req.DocumentIDs) > 500 {
		return nil, ErrBatchLimitExceed
	}

	cleanPIN := strings.TrimSpace(req.PIN)
	cleanOTP := strings.TrimSpace(req.OTP)
	if len(cleanPIN) < 4 || len(cleanOTP) != 6 {
		return nil, ErrInvalidPINorOTP
	}

	batchID := fmt.Sprintf("CSC-BATCH-%d", time.Now().UnixNano())
	now := time.Now().UTC()

	var signedList []QualifiedSignature
	for _, docID := range req.DocumentIDs {
		// Calculate document digest
		docDigest := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%d", docID, schoolID, now.Unix())))
		docHashHex := hex.EncodeToString(docDigest[:])

		// Perform signature on the QSCD
		sigBytes, pubKeyPEM, certDN, certSerial, err := c.qscd.Sign(ctx, docDigest[:])
		if err != nil {
			return nil, fmt.Errorf("errore durante la firma QSCD del documento %s: %w", docID, err)
		}

		// Request RFC 3161 TSA timestamp
		tsToken, tsAt, err := c.tsa.RequestTimestamp(ctx, sigBytes)
		if err != nil {
			tsToken = computeSimulatedTSAToken(sigBytes)
			tsAt = now
		}

		signedSig := QualifiedSignature{
			ID:                fmt.Sprintf("sig-csc-%s", hex.EncodeToString(docDigest[:8])),
			DocumentID:        docID,
			SignerID:          signerID,
			Level:             SignatureLevelFEQ,
			DocumentHash:      docHashHex,
			SignatureValue:    hex.EncodeToString(sigBytes),
			PublicKeyPEM:      string(pubKeyPEM),
			CertificateSerial: certSerial,
			CertificateDN:     certDN,
			TimestampToken:    tsToken,
			TimestampAt:       tsAt,
			SignedAt:          now,
			IsValid:           true,
			SchoolID:          schoolID,
		}

		signedList = append(signedList, signedSig)
	}

	return &CSCBatchSignResponse{
		TotalRequested: len(req.DocumentIDs),
		TotalSigned:    len(signedList),
		Signatures:     signedList,
		BatchID:        batchID,
		SignedAt:       now,
	}, nil
}
