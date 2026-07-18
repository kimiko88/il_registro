package signatures

import "time"

type Signature struct {
	ID            string    `json:"id" db:"id"`
	DocumentID    string    `json:"document_id" db:"document_id"`
	SignerID      string    `json:"signer_id" db:"signer_id"`
	SignedAt      time.Time `json:"signed_at" db:"signed_at"`
	SignatureHash string    `json:"signature_hash" db:"signature_hash"`
	IPAddress     string    `json:"ip_address" db:"ip_address"`
	Metadata      string    `json:"metadata" db:"metadata"` // JSONB stored as string
}

type SignRequest struct {
	DocumentID string `json:"document_id"`
	Pin        string `json:"pin"` // Simulated MFA
}

type Service interface {
	SignDocument(signerID string, req SignRequest) (*Signature, error)
	GetSignatures(documentID string) ([]Signature, error)
}
