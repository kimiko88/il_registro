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
	// GetDocumentContent restituisce il contenuto corrente del documento per il calcolo dell'hash crittografico.
	// Deve essere implementato da documents.Service.
	GetDocumentContent(ctx context.Context, id string) (string, error)
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
	return s.SignDocumentCtx(context.Background(), signerID, req)
}

// SignDocumentCtx è la versione context-aware di SignDocument.
// Viene chiamata internamente e dagli handler che dispongono di un context.
func (s *service) SignDocumentCtx(ctx context.Context, signerID string, req SignRequest) (*Signature, error) {
	if req.DocumentID == "" {
		return nil, errors.New("document_id is required")
	}
	if req.Pin == "" {
		return nil, errors.New("pin is required")
	}

	// 1. Verifica identità tramite MFA TOTP (obbligatorio per FEA/FES — CAD art. 21 / eIDAS art. 26)
	user, err := s.userRepo.GetByID(ctx, signerID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	if !user.MFAEnabled {
		// Fix: rimossa backdoor PIN "1234" — la firma richiede MFA abilitato.
		// Un utente senza MFA non può apporre una firma elettronica con valore legale.
		return nil, errors.New("firma non disponibile: abilitare l'autenticazione MFA (TOTP) prima di firmare documenti")
	}
	if !totp.Validate(req.Pin, user.MFASecret) {
		return nil, errors.New("codice OTP non valido")
	}

	// 2. Recupera il contenuto reale del documento per il calcolo crittografico dell'hash.
	// Fix: eliminato il mock "content-of-{id}" che rendeva la firma crittograficamente inutile.
	docContent, err := s.docsSvc.GetDocumentContent(ctx, req.DocumentID)
	if err != nil {
		return nil, fmt.Errorf("impossibile recuperare il contenuto del documento per la firma: %w", err)
	}
	if docContent == "" {
		return nil, errors.New("impossibile firmare un documento senza contenuto")
	}

	// 3. Calcola SHA-256 del contenuto reale
	hash := sha256.Sum256([]byte(docContent))
	hashString := hex.EncodeToString(hash[:])

	// 4. Crea il record firma
	sig := &Signature{
		DocumentID:    req.DocumentID,
		SignerID:      signerID,
		SignatureHash: hashString,
		IPAddress:     req.IPAddress, // impostato dall'handler dalla request reale
		Metadata:      "{}",
	}

	if err := s.repo.Create(ctx, sig); err != nil {
		return nil, err
	}

	// 5. Blocca il documento (non più modificabile dopo la firma)
	if err := s.docsSvc.LockDocument(ctx, req.DocumentID); err != nil {
		return nil, fmt.Errorf("failed to lock document: %w", err)
	}

	return sig, nil
}

func (s *service) GetSignatures(documentID string) ([]Signature, error) {
	return s.repo.FindByDocumentID(context.Background(), documentID)
}
