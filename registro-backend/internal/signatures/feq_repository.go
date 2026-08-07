package signatures

import (
	"context"
	"database/sql"
)

type feqRepo struct {
	db *sql.DB
}

// NewFEQRepository crea il repository PostgreSQL per le firme qualificate.
func NewFEQRepository(db *sql.DB) feqRepository {
	return &feqRepo{db: db}
}

func (r *feqRepo) CreateQualified(ctx context.Context, sig *QualifiedSignature) error {
	query := `
		INSERT INTO qualified_signatures (
			document_id, signer_id, school_id, level,
			document_hash, signature_value,
			certificate_serial, certificate_dn,
			timestamp_token, timestamp_at,
			is_valid, ip_address
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		RETURNING id, signed_at`
	return r.db.QueryRowContext(ctx, query,
		sig.DocumentID, sig.SignerID, sig.SchoolID, string(sig.Level),
		sig.DocumentHash, sig.SignatureValue,
		sig.CertificateSerial, sig.CertificateDN,
		sig.TimestampToken, sig.TimestampAt,
		sig.IsValid, sig.IPAddress,
	).Scan(&sig.ID, &sig.SignedAt)
}

func (r *feqRepo) FindQualifiedByID(ctx context.Context, id string) (*QualifiedSignature, error) {
	query := `
		SELECT id, document_id, signer_id, school_id, level,
		       document_hash, signature_value,
		       certificate_serial, certificate_dn,
		       timestamp_token, timestamp_at,
		       signed_at, revoked_at, is_valid, ip_address
		FROM qualified_signatures WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var s QualifiedSignature
	if err := row.Scan(
		&s.ID, &s.DocumentID, &s.SignerID, &s.SchoolID, &s.Level,
		&s.DocumentHash, &s.SignatureValue,
		&s.CertificateSerial, &s.CertificateDN,
		&s.TimestampToken, &s.TimestampAt,
		&s.SignedAt, &s.RevokedAt, &s.IsValid, &s.IPAddress,
	); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *feqRepo) FindQualifiedByDocumentID(ctx context.Context, docID string) ([]QualifiedSignature, error) {
	query := `
		SELECT id, document_id, signer_id, school_id, level,
		       document_hash, signature_value,
		       certificate_serial, certificate_dn,
		       timestamp_token, timestamp_at,
		       signed_at, revoked_at, is_valid, ip_address
		FROM qualified_signatures WHERE document_id = $1
		ORDER BY signed_at ASC`
	rows, err := r.db.QueryContext(ctx, query, docID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []QualifiedSignature
	for rows.Next() {
		var s QualifiedSignature
		if err := rows.Scan(
			&s.ID, &s.DocumentID, &s.SignerID, &s.SchoolID, &s.Level,
			&s.DocumentHash, &s.SignatureValue,
			&s.CertificateSerial, &s.CertificateDN,
			&s.TimestampToken, &s.TimestampAt,
			&s.SignedAt, &s.RevokedAt, &s.IsValid, &s.IPAddress,
		); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}
