package signatures

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/pquerna/otp/totp"
	"registro-backend/internal/users"
)

type DocumentService interface {
	LockDocument(ctx context.Context, id string) error
}

type UserLookup interface {
	GetByID(ctx context.Context, id string) (*users.User, error)
}

type service struct {
	repo     Repository
	docsSvc  DocumentService
	userRepo UserLookup
}

func NewService(repo Repository, docsSvc DocumentService, userRepo UserLookup) Service {
	return &service{repo: repo, docsSvc: docsSvc, userRepo: userRepo}
}

func (s *service) SignDocument(signerID string, req SignRequest) (*Signature, error) {
	if req.DocumentID == "" {
		return nil, errors.New("document_id is required")
	}
	if req.Pin == "" {
		return nil, errors.New("pin is required")
	}

	// 1. Verify OTP/PIN (MFA FEA Verification)
	user, err := s.userRepo.GetByID(context.Background(), signerID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	if user.MFAEnabled {
		if !totp.Validate(req.Pin, user.MFASecret) {
			return nil, errors.New("codice OTP non valido")
		}
	} else {
		// Fallback to legacy PIN for backwards compatibility / dev mode
		if req.Pin != "1234" {
			return nil, errors.New("mfa non abilitato e PIN non valido")
		}
	}

	// 2. Fetch Document Content (Mocked)
	docContent := fmt.Sprintf("content-of-%s", req.DocumentID)

	// 3. Calculate Hash
	hash := sha256.Sum256([]byte(docContent))
	hashString := hex.EncodeToString(hash[:])

	// 4. Create Signature
	sig := &Signature{
		DocumentID:    req.DocumentID,
		SignerID:      signerID,
		SignatureHash: hashString,
		IPAddress:     "127.0.0.1",
		Metadata:      "{}",
	}

	if err := s.repo.Create(context.Background(), sig); err != nil {
		return nil, err
	}

	// 5. Lock Document
	if err := s.docsSvc.LockDocument(context.Background(), req.DocumentID); err != nil {
		return nil, fmt.Errorf("failed to lock document: %w", err)
	}

	return sig, nil
}

func (s *service) GetSignatures(documentID string) ([]Signature, error) {
	return s.repo.FindByDocumentID(context.Background(), documentID)
}
