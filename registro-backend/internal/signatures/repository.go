package signatures

import (
	"context"
	"database/sql"
)

type Repository interface {
	Create(ctx context.Context, sig *Signature) error
	FindByDocumentID(ctx context.Context, docID string) ([]Signature, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, sig *Signature) error {
	query := `
		INSERT INTO signatures (document_id, signer_id, signature_hash, ip_address, metadata)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, signed_at`

	return r.db.QueryRowContext(ctx, query,
		sig.DocumentID, sig.SignerID, sig.SignatureHash, sig.IPAddress, sig.Metadata,
	).Scan(&sig.ID, &sig.SignedAt)
}

func (r *repository) FindByDocumentID(ctx context.Context, docID string) ([]Signature, error) {
	query := `SELECT id, document_id, signer_id, signed_at, signature_hash, ip_address, metadata FROM signatures WHERE document_id = $1`
	rows, err := r.db.QueryContext(ctx, query, docID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sigs []Signature
	for rows.Next() {
		var s Signature
		if err := rows.Scan(&s.ID, &s.DocumentID, &s.SignerID, &s.SignedAt, &s.SignatureHash, &s.IPAddress, &s.Metadata); err != nil {
			return nil, err
		}
		sigs = append(sigs, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sigs, nil
}

