package signatures

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

type DocumentService interface {
	LockDocument(ctx context.Context, id string) error
}

type service struct {
	repo    Repository
	docsSvc DocumentService
}

func NewService(repo Repository, docsSvc DocumentService) Service {
	return &service{repo: repo, docsSvc: docsSvc}
}

func (s *service) SignDocument(signerID string, req SignRequest) (*Signature, error) {
	if req.DocumentID == "" {
		return nil, errors.New("document_id is required")
	}
	if req.Pin == "" {
		return nil, errors.New("pin is required")
	}

	// 1. Verify PIN (MFA Simulation)
	// In a real system, we'd check against a hashed secret for the user.
	// For MVP: assume "1234" is the correct PIN for everyone (or check something basic)
	if req.Pin != "1234" {
		return nil, errors.New("invalid pin")
	}

	// 2. Fetch Document Content (Mocked)
	// We need to hash the *content* to make it immutable.
	// In a real app, we'd fetch the document BLOB or text content.
	docContent := fmt.Sprintf("content-of-%s", req.DocumentID)

	// 3. Calculate Hash
	hash := sha256.Sum256([]byte(docContent))
	hashString := hex.EncodeToString(hash[:])

	// 4. Create Signature
	sig := &Signature{
		DocumentID:    req.DocumentID,
		SignerID:      signerID,
		SignatureHash: hashString,
		IPAddress:     "127.0.0.1", // In real handler, get from Context
		Metadata:      "{}",
	}

	if err := s.repo.Create(context.Background(), sig); err != nil {
		return nil, err
	}

	// 5. Lock Document
	if err := s.docsSvc.LockDocument(context.Background(), req.DocumentID); err != nil {
		// Log error but don't fail signature creation? Or fail?
		// Ideally transactional. For now, return error.
		return nil, fmt.Errorf("failed to lock document: %w", err)
	}

	return sig, nil
}

func (s *service) GetSignatures(documentID string) ([]Signature, error) {
	return s.repo.FindByDocumentID(context.Background(), documentID)
}
