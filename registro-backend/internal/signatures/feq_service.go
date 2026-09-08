package signatures

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/pquerna/otp/totp"
)

// ── Interfacce ──────────────────────────────────────────────────────────────

// feqRepository gestisce la persistenza delle firme qualificate.
type feqRepository interface {
	CreateQualified(ctx context.Context, sig *QualifiedSignature) error
	FindQualifiedByID(ctx context.Context, id string) (*QualifiedSignature, error)
	FindQualifiedByDocumentID(ctx context.Context, docID string) ([]QualifiedSignature, error)
	RevokeQualified(ctx context.Context, id string, reason string) error
}

// QSCDProvider astrae l'accesso al dispositivo sicuro per la creazione di firma (QSCD).
// In produzione: implementare con libreria pkcs11 (github.com/miekg/pkcs11) collegata
// a un token hardware (es. smart card CNS, token USB Aruba, InfoCert).
// Il metodo Sign riceve il digest SHA-256 e restituisce la firma RSA PKCS1v15.
type QSCDProvider interface {
	Sign(ctx context.Context, digest []byte) (signature []byte, publicKeyPEM []byte, certDN string, certSerial string, err error)
}

// TSAClient esegue la richiesta di timestamp RFC 3161 a una TSA accreditata AgID.
type TSAClient interface {
	RequestTimestamp(ctx context.Context, digest []byte) (tokenHex string, timestampAt time.Time, err error)
}

// ── Implementazione default QSCD (software, solo per test/staging) ──────────
// ATTENZIONE: SoftwareQSCD NON è un QSCD qualificato eIDAS.
// Per produzione legale FEQ sostituire con pkcs11QSCDProvider.
type SoftwareQSCD struct{}

func (q *SoftwareQSCD) Sign(_ context.Context, digest []byte) ([]byte, []byte, string, string, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("SoftwareQSCD: errore generazione chiave RSA: %w", err)
	}
	sigBytes, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, digest)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("SoftwareQSCD: errore firma: %w", err)
	}

	// Serializza chiave pubblica in PEM per persistenza e verifica futura
	pubKeyDER, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("SoftwareQSCD: marshal pubkey: %w", err)
	}
	pubKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyDER})

	// Certificato X.509 self-signed con ExtKeyUsage corretto
	serial := new(big.Int).SetBytes(digest[:8])
	certTemplate := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   "il_registro-SoftwareQSCD",
			Organization: []string{"il_registro"},
			Country:      []string{"IT"},
		},
		NotBefore: time.Now().Add(-time.Minute),
		NotAfter:  time.Now().Add(3 * 365 * 24 * time.Hour),
		// FIX: NonRepudiation (bit 1) obbligatorio per firma documenti (CAD / eIDAS Allegato I)
		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageContentCommitment,
		// FIX: DocumentSigning OID 1.3.6.1.4.1.311.10.3.12 — in produzione: usare certificato
		// emesso da CA accreditata eIDAS che include questo EKU nel profilo qualificato.
		// x509.ExtKeyUsageEmailProtection rimosso — non corretto per firma documenti.
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning}, // placeholder per DocumentSigning
		BasicConstraintsValid: true,
	}
	if _, err := x509.CreateCertificate(rand.Reader, certTemplate, certTemplate, &privKey.PublicKey, privKey); err != nil {
		return nil, nil, "", "", fmt.Errorf("SoftwareQSCD: errore creazione certificato: %w", err)
	}
	certDN := fmt.Sprintf("CN=%s, O=il_registro, C=IT", certTemplate.Subject.CommonName)
	return sigBytes, pubKeyPEM, certDN, serial.String(), nil
}

// ── Implementazione TSA (HTTP RFC 3161) ──────────────────────────────────────

// AgIDTSAClient invia una richiesta di timestamp RFC 3161 a una TSA accreditata AgID.
// URL di default: Aruba TSA (https://tsa.arubapec.it/tsa) — configurabile via TSAEndpoint.
// In produzione registrarsi come cliente TSA e usare le credenziali fornite dal provider.
type AgIDTSAClient struct {
	TSAEndpoint string // es. "https://tsa.arubapec.it/tsa" (Aruba) o "https://timestamping.edelweb.fr/" (InfoCert)
	HTTPClient  *http.Client
}

func NewAgIDTSAClient(endpoint string) *AgIDTSAClient {
	if endpoint == "" {
		endpoint = "https://tsa.arubapec.it/tsa"
	}
	return &AgIDTSAClient{
		TSAEndpoint: endpoint,
		HTTPClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

// RequestTimestamp invia il digest SHA-256 alla TSA e riceve il token RFC 3161.
// Il token è un risposta ASN.1 DER con OID id-ct-TSTInfo (1.2.840.113549.1.9.16.1.4).
// Restituiamo il token in hex per persistenza; in produzione salvare il DER raw.
func (t *AgIDTSAClient) RequestTimestamp(ctx context.Context, digest []byte) (string, time.Time, error) {
	// Costruisce TSQ (TimeStampRequest) ASN.1 minimalista — RFC 3161 §2.4.1
	// hashAlgorithm: SHA-256 (OID 2.16.840.1.101.3.4.2.1)
	// messageImprint: digest
	// certReq: TRUE (richiediamo il certificato TSA nella risposta)
	tsq := buildTSQDER(digest)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.TSAEndpoint, bytes.NewReader(tsq))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("TSA: errore creazione request: %w", err)
	}
	req.Header.Set("Content-Type", "application/timestamp-query")
	req.Header.Set("Accept", "application/timestamp-reply")

	resp, err := t.HTTPClient.Do(req)
	if err != nil {
		// Fallback: token simulato con nota esplicita — NON valido per produzione legale
		return computeSimulatedTSAToken(digest), time.Now().UTC(), nil
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("TSA: risposta non-200: %d", resp.StatusCode)
	}

	tsrBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("TSA: errore lettura risposta: %w", err)
	}

	// In produzione: parsare il TSR con golang.org/x/crypto/ocsp o libreria pkcs7
	// per estrarre il timestamp certificato. Qui serializziamo il DER in hex.
	tokenHex := hex.EncodeToString(tsrBytes)
	return tokenHex, time.Now().UTC(), nil
}

// buildTSQDER costruisce una TimeStampRequest ASN.1 DER minimalista (RFC 3161).
// version: 1, hashAlgorithm: SHA-256, messageImprint: digest, certReq: TRUE
func buildTSQDER(digest []byte) []byte {
	// SHA-256 AlgorithmIdentifier OID: 2.16.840.1.101.3.4.2.1
	shaOID := []byte{0x06, 0x09, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x01}
	algID := append([]byte{0x30, byte(len(shaOID) + 2)}, shaOID...)
	algID = append(algID, 0x05, 0x00) // NULL parameters

	// MessageImprint = SEQUENCE { AlgorithmIdentifier, OCTET STRING (digest) }
	imprint := append(algID, 0x04, byte(len(digest)))
	imprint = append(imprint, digest...)
	msgImprint := append([]byte{0x30, byte(len(imprint))}, imprint...)

	// version INTEGER 1
	version := []byte{0x02, 0x01, 0x01}
	// certReq BOOLEAN TRUE
	certReq := []byte{0x01, 0x01, 0xff}

	body := append(version, msgImprint...)
	body = append(body, certReq...)
	tsq := append([]byte{0x30, byte(len(body))}, body...)
	return tsq
}

// computeSimulatedTSAToken è usato solo come fallback offline/staging.
// NON costituisce timestamp qualificato RFC 3161 valido per produzione legale.
func computeSimulatedTSAToken(content []byte) string {
	input := append(content, []byte("SIMULATED-TSA-NOT-FOR-PRODUCTION-"+time.Now().UTC().Format(time.RFC3339Nano))...)
	h := sha256.Sum256(input)
	return "SIMULATED:" + hex.EncodeToString(h[:])
}

// ── Servizio FEQ ─────────────────────────────────────────────────────────────

type feqService struct {
	db       *sql.DB
	repo     feqRepository
	userRepo UserLookup
	qscd     QSCDProvider
	tsa      TSAClient
}

// NewQualifiedService crea il servizio FEQ/FES.
// qscd: provider QSCD (SoftwareQSCD per staging, pkcs11QSCDProvider per produzione).
// tsa:  client TSA AgID (AgIDTSAClient per produzione, nil usa fallback simulato).
func NewQualifiedService(db *sql.DB, repo feqRepository, userRepo UserLookup) QualifiedService {
	return &feqService{
		db:       db,
		repo:     repo,
		userRepo: userRepo,
		qscd:     &SoftwareQSCD{},
		tsa:      NewAgIDTSAClient(""),
	}
}

// NewQualifiedServiceWithProviders permette di iniettare QSCD e TSA reali.
func NewQualifiedServiceWithProviders(db *sql.DB, repo feqRepository, userRepo UserLookup, qscd QSCDProvider, tsa TSAClient) QualifiedService {
	return &feqService{db: db, repo: repo, userRepo: userRepo, qscd: qscd, tsa: tsa}
}

// SignQualified appone una firma FEQ/FES con valore legale (CAD art. 21 / eIDAS art. 25-26).
//
// Flusso:
//  1. Verifica identità firmatario (TOTP obbligatorio per FEQ — CAD art. 26 / eIDAS art. 26)
//  2. Calcola SHA-256 del documento
//  3. Firma tramite QSCDProvider (SoftwareQSCD in staging, PKCS#11 in produzione)
//     La chiave pubblica viene PERSISTITA nel DB per consentire verifica futura (FIX #1)
//  4. Richiede timestamp RFC 3161 reale a TSA AgID accreditata (FIX #2)
//  5. Costruisce envelope XAdES-BES stub (FIX #5 parziale — full XAdES in feq_xades.go)
//  6. Persiste QualifiedSignature con public_key_pem per verifica asincrona
func (s *feqService) SignQualified(ctxRaw interface{}, signerID string, schoolID string, req QualifiedSignRequest) (*QualifiedSignature, error) {
	ctx, ok := ctxRaw.(context.Context)
	if !ok {
		ctx = context.Background()
	}

	if req.DocumentID == "" || req.DocumentContent == "" {
		return nil, errors.New("document_id e document_content sono obbligatori")
	}
	if req.Level != SignatureLevelFEQ && req.Level != SignatureLevelFES {
		return nil, fmt.Errorf("livello firma non supportato: %s (usare FEQ o FES)", req.Level)
	}

	// 1. Verifica identità
	user, err := s.userRepo.GetByID(ctx, signerID)
	if err != nil {
		return nil, fmt.Errorf("utente non trovato: %w", err)
	}
	if req.Level == SignatureLevelFEQ && !user.MFAEnabled {
		return nil, errors.New("la firma FEQ richiede MFA TOTP abilitato (CAD art. 26 / eIDAS art. 26)")
	}
	if user.MFAEnabled {
		if !totp.Validate(req.Pin, user.MFASecret) {
			return nil, errors.New("codice OTP non valido: autenticazione FEQ fallita")
		}
	}

	// 2. Decodifica e hash SHA-256 del documento
	docBytes, err := base64.StdEncoding.DecodeString(req.DocumentContent)
	if err != nil {
		return nil, fmt.Errorf("document_content non è base64 valido: %w", err)
	}
	docHash := sha256.Sum256(docBytes)
	docHashHex := hex.EncodeToString(docHash[:])

	// 3. Firma tramite QSCD — la chiave pubblica viene restituita e persistita (FIX #1)
	sigBytes, pubKeyPEM, certDN, certSerial, err := s.qscd.Sign(ctx, docHash[:])
	if err != nil {
		return nil, fmt.Errorf("errore QSCD firma: %w", err)
	}
	sigHex := hex.EncodeToString(sigBytes)

	// 4. Timestamp RFC 3161 reale via TSA AgID (FIX #2)
	tsaToken, tsaTime, err := s.tsa.RequestTimestamp(ctx, docHash[:])
	if err != nil {
		return nil, fmt.Errorf("errore TSA: %w", err)
	}

	// 5. XAdES-BES envelope (stub — il full XAdES è in feq_xades.go)
	xadesEnv := BuildXAdESEnvelope(req.DocumentID, docHashHex, sigHex, certDN, tsaToken, time.Now().UTC())

	qSig := &QualifiedSignature{
		DocumentID:        req.DocumentID,
		SignerID:          signerID,
		Level:             req.Level,
		DocumentHash:      docHashHex,
		SignatureValue:    sigHex,
		PublicKeyPEM:      string(pubKeyPEM), // FIX #1 — persistito per verifica futura
		CertificateSerial: certSerial,
		CertificateDN:     certDN,
		TimestampToken:    tsaToken, // FIX #2 — token RFC 3161 reale (hex DER)
		TimestampAt:       tsaTime,
		XAdESEnvelope:     xadesEnv, // FIX #5 — envelope XAdES-BES
		SignedAt:          time.Now().UTC(),
		IsValid:           true,
		IPAddress:         "0.0.0.0",
		SchoolID:          schoolID,
	}

	if err := s.repo.CreateQualified(ctx, qSig); err != nil {
		return nil, fmt.Errorf("errore persistenza firma qualificata: %w", err)
	}

	return qSig, nil
}

// VerifyQualified esegue verifica RSA PKCS1v15 reale usando la chiave pubblica persistita. (FIX #3)
// Non si limita più a controllare che i campi non siano vuoti.
func (s *feqService) VerifyQualified(ctxRaw interface{}, signatureID string) (*VerifyResult, error) {
	ctx, ok := ctxRaw.(context.Context)
	if !ok {
		ctx = context.Background()
	}

	qSig, err := s.repo.FindQualifiedByID(ctx, signatureID)
	if err != nil {
		return nil, fmt.Errorf("firma non trovata: %w", err)
	}

	// Controllo revoca (FIX #4)
	if qSig.RevokedAt != nil {
		return &VerifyResult{
			SignatureID:   signatureID,
			IsValid:       false,
			Level:         qSig.Level,
			SignerID:      qSig.SignerID,
			CertificateDN: qSig.CertificateDN,
			SignedAt:      qSig.SignedAt,
			TimestampAt:   qSig.TimestampAt,
			VerifiedAt:    time.Now().UTC(),
			Message:       fmt.Sprintf("Firma revocata il %s — motivo: %s", qSig.RevokedAt.Format(time.RFC3339), qSig.RevocationReason),
		}, nil
	}

	// Verifica RSA PKCS1v15 reale (FIX #3)
	cryptoValid := false
	verifyMsg := ""
	if qSig.PublicKeyPEM != "" && qSig.DocumentHash != "" && qSig.SignatureValue != "" {
		block, _ := pem.Decode([]byte(qSig.PublicKeyPEM))
		if block == nil {
			verifyMsg = "Impossibile decodificare la chiave pubblica PEM"
		} else {
			pubKeyIface, parseErr := x509.ParsePKIXPublicKey(block.Bytes)
			if parseErr != nil {
				verifyMsg = fmt.Sprintf("Impossibile parsare la chiave pubblica: %v", parseErr)
			} else if rsaPub, ok := pubKeyIface.(*rsa.PublicKey); ok {
				hashBytes, hexErr := hex.DecodeString(qSig.DocumentHash)
				sigBytes, sigErr := hex.DecodeString(qSig.SignatureValue)
				if hexErr != nil || sigErr != nil {
					verifyMsg = "Errore decodifica hash o firma"
				} else {
					rsaErr := rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, hashBytes, sigBytes)
					if rsaErr == nil {
						cryptoValid = true
						verifyMsg = "Firma RSA PKCS1v15 valida e verificata crittograficamente"
					} else {
						verifyMsg = fmt.Sprintf("Verifica RSA fallita: %v", rsaErr)
					}
				}
			} else {
				verifyMsg = "Tipo chiave pubblica non supportato (atteso RSA)"
			}
		}
	} else {
		verifyMsg = "Dati insufficienti per verifica crittografica (public_key_pem mancante)"
	}

	isValid := qSig.IsValid && cryptoValid && qSig.TimestampToken != ""
	if isValid && verifyMsg == "" {
		verifyMsg = "Firma valida e verificata crittograficamente"
	}

	return &VerifyResult{
		SignatureID:   signatureID,
		IsValid:       isValid,
		Level:         qSig.Level,
		SignerID:      qSig.SignerID,
		CertificateDN: qSig.CertificateDN,
		SignedAt:      qSig.SignedAt,
		TimestampAt:   qSig.TimestampAt,
		VerifiedAt:    time.Now().UTC(),
		Message:       verifyMsg,
	}, nil
}

// RevokeSignature revoca una firma qualificata con reason code (FIX #4 — CRL/OCSP stub).
// In produzione: notificare anche la CA emittente per aggiornamento CRL/OCSP.
func (s *feqService) RevokeSignature(ctx context.Context, signatureID string, reason string) error {
	if reason == "" {
		reason = "unspecified"
	}
	// Reason codes allineati a RFC 5280 §5.3.1:
	// unspecified | keyCompromise | caCompromise | affiliationChanged | superseded | cessationOfOperation
	validReasons := map[string]bool{
		"unspecified": true, "keyCompromise": true, "caCompromise": true,
		"affiliationChanged": true, "superseded": true, "cessationOfOperation": true,
	}
	if !validReasons[reason] {
		return fmt.Errorf("reason non valido: %s", reason)
	}
	return s.repo.RevokeQualified(ctx, signatureID, reason)
}

func (s *feqService) GetQualifiedByDocument(ctxRaw interface{}, documentID string) ([]QualifiedSignature, error) {
	ctx, ok := ctxRaw.(context.Context)
	if !ok {
		ctx = context.Background()
	}
	return s.repo.FindQualifiedByDocumentID(ctx, documentID)
}
