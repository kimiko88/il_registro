package signatures

import (
	"fmt"
	"time"
)

// SignatureLevel rappresenta il livello di firma elettronica secondo eIDAS e CAD italiano.
type SignatureLevel string

const (
	SignatureLevelFEA SignatureLevel = "FEA" // Firma Elettronica Avanzata
	SignatureLevelFES SignatureLevel = "FES" // Firma Elettronica Semplice
	SignatureLevelFEQ SignatureLevel = "FEQ" // Firma Elettronica Qualificata
)

// Signature rappresenta una firma FEA base (legacy).
type Signature struct {
	ID            string    `json:"id" db:"id"`
	DocumentID    string    `json:"document_id" db:"document_id"`
	SignerID      string    `json:"signer_id" db:"signer_id"`
	SignedAt      time.Time `json:"signed_at" db:"signed_at"`
	SignatureHash string    `json:"signature_hash" db:"signature_hash"`
	IPAddress     string    `json:"ip_address" db:"ip_address"`
	Metadata      string    `json:"metadata" db:"metadata"`
}

// QualifiedSignature rappresenta una firma FEQ/FES con pieno valore legale (CAD art. 21).
type QualifiedSignature struct {
	ID             string         `json:"id" db:"id"`
	DocumentID     string         `json:"document_id" db:"document_id"`
	SignerID       string         `json:"signer_id" db:"signer_id"`
	Level          SignatureLevel `json:"level" db:"level"`
	DocumentHash   string         `json:"document_hash" db:"document_hash"`
	SignatureValue string         `json:"signature_value" db:"signature_value"`
	// FIX #1: chiave pubblica persistita in PEM per verifica RSA futura (non più effimera)
	PublicKeyPEM      string `json:"public_key_pem" db:"public_key_pem"`
	CertificateSerial string `json:"certificate_serial" db:"certificate_serial"`
	CertificateDN     string `json:"certificate_dn" db:"certificate_dn"`
	// FIX #2: token RFC 3161 reale (DER hex) da TSA AgID accreditata
	TimestampToken string    `json:"timestamp_token" db:"timestamp_token"`
	TimestampAt    time.Time `json:"timestamp_at" db:"timestamp_at"`
	SignedAt       time.Time `json:"signed_at" db:"signed_at"`
	// FIX #4: revoca con reason code RFC 5280
	RevokedAt        *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	RevocationReason string     `json:"revocation_reason,omitempty" db:"revocation_reason"`
	IsValid          bool       `json:"is_valid" db:"is_valid"`
	IPAddress        string     `json:"ip_address" db:"ip_address"`
	SchoolID         string     `json:"school_id" db:"school_id"`
	// FIX #5: envelope XAdES-BES/T (eIDAS art. 37 + Decisione 2015/1506/UE)
	XAdESEnvelope string `json:"xades_envelope,omitempty" db:"xades_envelope"`
}

// QualifiedSignRequest è la richiesta per apporre una firma FEQ/FES.
type QualifiedSignRequest struct {
	DocumentID      string         `json:"document_id" binding:"required"`
	DocumentContent string         `json:"document_content" binding:"required"`
	Level           SignatureLevel `json:"level" binding:"required"`
	Pin             string         `json:"pin" binding:"required"`
}

// VerifyResult è il risultato della verifica crittografica di una firma.
type VerifyResult struct {
	SignatureID   string         `json:"signature_id"`
	IsValid       bool           `json:"is_valid"`
	Level         SignatureLevel `json:"level"`
	SignerID      string         `json:"signer_id"`
	CertificateDN string         `json:"certificate_dn"`
	SignedAt      time.Time      `json:"signed_at"`
	TimestampAt   time.Time      `json:"timestamp_at"`
	VerifiedAt    time.Time      `json:"verified_at"`
	Message       string         `json:"message"`
}

// SidiExportRecord è un record per l'export SIDI/MIUR.
type SidiExportRecord struct {
	CodiceMeccanografico string `json:"codice_meccanografico"`
	AnnoScolastico       string `json:"anno_scolastico"`
	TipologiaExport      string `json:"tipologia_export"`
	DataGenerazione      string `json:"data_generazione"`
	HashIntegrità        string `json:"hash_integrita"`
	StatoTrasmissione    string `json:"stato_trasmissione"`
	// FIX SIDI #2: codice ricevuta restituito dal portale MIUR dopo trasmissione
	CodiceRicevuta string `json:"codice_ricevuta,omitempty"`
}

// SidiRejectionError rappresenta un errore di rigetto dal portale SIDI (FIX SIDI #4).
type SidiRejectionError struct {
	HTTPStatus  int    `json:"http_status"`
	ErrorCode   string `json:"error_code"`
	Description string `json:"description"`
}

func (e *SidiRejectionError) Error() string {
	return fmt.Sprintf("SIDI rejection %d — %s: %s", e.HTTPStatus, e.ErrorCode, e.Description)
}

// ArchivePackage rappresenta un pacchetto di conservazione sostitutiva CAD.
type ArchivePackage struct {
	ID               string    `json:"id" db:"id"`
	SchoolID         string    `json:"school_id" db:"school_id"`
	AcademicYear     string    `json:"academic_year" db:"academic_year"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	ManifestHash     string    `json:"manifest_hash" db:"manifest_hash"`
	DocumentCount    int       `json:"document_count" db:"document_count"`
	PackageSizeBytes int64     `json:"package_size_bytes" db:"package_size_bytes"`
	StoragePath      string    `json:"storage_path" db:"storage_path"`
	IsSealed         bool      `json:"is_sealed" db:"is_sealed"`
}

// SignRequest è la richiesta per la firma FEA base (legacy).
type SignRequest struct {
	DocumentID string `json:"document_id"`
	Pin        string `json:"pin"`
	// IPAddress è l'indirizzo IP reale del client (impostato dall'handler dal contesto HTTP).
	// Non deve essere fornito dal body JSON — viene impostato internamente.
	IPAddress string `json:"-"`
}

// Service è l'interfaccia per la firma FEA base (legacy).
type Service interface {
	SignDocument(signerID string, req SignRequest) (*Signature, error)
	GetSignatures(documentID string) ([]Signature, error)
}

// QualifiedService è l'interfaccia per la firma FEQ/FES con valore legale.
type QualifiedService interface {
	SignQualified(ctx interface{}, signerID string, schoolID string, req QualifiedSignRequest) (*QualifiedSignature, error)
	VerifyQualified(ctx interface{}, signatureID string) (*VerifyResult, error)
	GetQualifiedByDocument(ctx interface{}, documentID string) ([]QualifiedSignature, error)
}
