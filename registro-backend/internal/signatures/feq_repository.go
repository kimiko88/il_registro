package signatures

import (
	"context"
	"database/sql"
	"time"
)

type feqRepo struct {
	db *sql.DB
}

func NewFEQRepository(db *sql.DB) feqRepository {
	return &feqRepo{db: db}
}

func (r *feqRepo) CreateQualified(ctx context.Context, sig *QualifiedSignature) error {
	query := `
		INSERT INTO qualified_signatures (
			document_id, signer_id, school_id, level,
			document_hash, signature_value, public_key_pem,
			certificate_serial, certificate_dn,
			timestamp_token, timestamp_at,
			is_valid, ip_address, xades_envelope
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		RETURNING id, signed_at`
	return r.db.QueryRowContext(ctx, query,
		sig.DocumentID, sig.SignerID, sig.SchoolID, string(sig.Level),
		sig.DocumentHash, sig.SignatureValue, sig.PublicKeyPEM,
		sig.CertificateSerial, sig.CertificateDN,
		sig.TimestampToken, sig.TimestampAt,
		sig.IsValid, sig.IPAddress, sig.XAdESEnvelope,
	).Scan(&sig.ID, &sig.SignedAt)
}

func (r *feqRepo) FindQualifiedByID(ctx context.Context, id string) (*QualifiedSignature, error) {
	query := `
		SELECT id, document_id, signer_id, school_id, level,
		       document_hash, signature_value, public_key_pem,
		       certificate_serial, certificate_dn,
		       timestamp_token, timestamp_at,
		       signed_at, revoked_at, revocation_reason, is_valid, ip_address, xades_envelope
		FROM qualified_signatures WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var s QualifiedSignature
	var revokedAt sql.NullTime
	var revocationReason sql.NullString
	var xadesEnvelope sql.NullString
	if err := row.Scan(
		&s.ID, &s.DocumentID, &s.SignerID, &s.SchoolID, &s.Level,
		&s.DocumentHash, &s.SignatureValue, &s.PublicKeyPEM,
		&s.CertificateSerial, &s.CertificateDN,
		&s.TimestampToken, &s.TimestampAt,
		&s.SignedAt, &revokedAt, &revocationReason, &s.IsValid, &s.IPAddress, &xadesEnvelope,
	); err != nil {
		return nil, err
	}
	if revokedAt.Valid {
		t := revokedAt.Time
		s.RevokedAt = &t
	}
	if revocationReason.Valid {
		s.RevocationReason = revocationReason.String
	}
	if xadesEnvelope.Valid {
		s.XAdESEnvelope = xadesEnvelope.String
	}
	return &s, nil
}

func (r *feqRepo) FindQualifiedByDocumentID(ctx context.Context, docID string) ([]QualifiedSignature, error) {
	query := `
		SELECT id, document_id, signer_id, school_id, level,
		       document_hash, signature_value, public_key_pem,
		       certificate_serial, certificate_dn,
		       timestamp_token, timestamp_at,
		       signed_at, revoked_at, revocation_reason, is_valid, ip_address, xades_envelope
		FROM qualified_signatures WHERE document_id = $1 ORDER BY signed_at ASC`
	rows, err := r.db.QueryContext(ctx, query, docID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []QualifiedSignature
	for rows.Next() {
		var s QualifiedSignature
		var revokedAt sql.NullTime
		var revocationReason sql.NullString
		var xadesEnvelope sql.NullString
		if err := rows.Scan(
			&s.ID, &s.DocumentID, &s.SignerID, &s.SchoolID, &s.Level,
			&s.DocumentHash, &s.SignatureValue, &s.PublicKeyPEM,
			&s.CertificateSerial, &s.CertificateDN,
			&s.TimestampToken, &s.TimestampAt,
			&s.SignedAt, &revokedAt, &revocationReason, &s.IsValid, &s.IPAddress, &xadesEnvelope,
		); err != nil {
			return nil, err
		}
		if revokedAt.Valid {
			t := revokedAt.Time
			s.RevokedAt = &t
		}
		if revocationReason.Valid {
			s.RevocationReason = revocationReason.String
		}
		if xadesEnvelope.Valid {
			s.XAdESEnvelope = xadesEnvelope.String
		}
		result = append(result, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}


// RevokeQualified imposta revoked_at e revocation_reason (FIX #4 — CRL/OCSP).
// La RLS PostgreSQL consente UPDATE solo su is_valid/revoked_at/revocation_reason.
func (r *feqRepo) RevokeQualified(ctx context.Context, id string, reason string) error {
	query := `UPDATE qualified_signatures SET revoked_at = $1, revocation_reason = $2, is_valid = FALSE WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, time.Now().UTC(), reason, id)
	return err
}
