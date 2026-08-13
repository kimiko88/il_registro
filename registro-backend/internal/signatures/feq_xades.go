package signatures

import (
	"encoding/xml"
	"fmt"
	"time"
)

// ── XAdES-BES / XAdES-T envelope (eIDAS art. 37 + Decisione 2015/1506/UE) ──
//
// eIDAS art. 37 e la Decisione di esecuzione 2015/1506/UE richiedono che le
// firme elettroniche avanzate e qualificate usino formati AdES standardizzati:
//   - XAdES per documenti XML
//   - PAdES per documenti PDF
//   - CAdES per qualsiasi tipo di file (busta crittografica CMS)
//
// Questo file implementa XAdES-BES (livello base) e XAdES-T (con timestamp).
// Per piena conformità eIDAS usare la libreria esterna:
//   github.com/digitorus/pkcs7  +  go-xades (o DSS di e-CODEX)

// XAdESSignature è la struttura XML dell'envelope XAdES-BES conforme
// a ETSI EN 319 132-1 (XAdES Baseline Profile — livello B-B e B-T).
type XAdESSignature struct {
	XMLName        xml.Name            `xml:"Signature"`
	Xmlns          string              `xml:"xmlns,attr"`
	XmlnsXades     string              `xml:"xmlns:xades,attr"`
	ID             string              `xml:"Id,attr"`
	SignedInfo     XAdESSignedInfo     `xml:"SignedInfo"`
	SignatureValue XAdESSignatureValue `xml:"SignatureValue"`
	KeyInfo        XAdESKeyInfo        `xml:"KeyInfo"`
	Object         XAdESObject         `xml:"Object"`
}

type XAdESSignedInfo struct {
	CanonicalizationMethod XAdESMethod `xml:"CanonicalizationMethod"`
	SignatureMethod        XAdESMethod `xml:"SignatureMethod"`
	Reference              XAdESRef    `xml:"Reference"`
}

type XAdESMethod struct {
	Algorithm string `xml:"Algorithm,attr"`
}

type XAdESRef struct {
	URI          string      `xml:"URI,attr"`
	DigestMethod XAdESMethod `xml:"DigestMethod"`
	DigestValue  string      `xml:"DigestValue"`
}

type XAdESSignatureValue struct {
	ID    string `xml:"Id,attr"`
	Value string `xml:",chardata"`
}

type XAdESKeyInfo struct {
	X509Data XAdESX509Data `xml:"X509Data"`
}

type XAdESX509Data struct {
	X509SubjectName  string            `xml:"X509SubjectName"`
	X509IssuerSerial XAdESIssuerSerial `xml:"X509IssuerSerial"`
}

type XAdESIssuerSerial struct {
	X509IssuerName   string `xml:"X509IssuerName"`
	X509SerialNumber string `xml:"X509SerialNumber"`
}

type XAdESObject struct {
	QualifyingProperties XAdESQualifyingProps `xml:"xades:QualifyingProperties"`
}

type XAdESQualifyingProps struct {
	Target             string              `xml:"Target,attr"`
	SignedProperties   XAdESSignedProps    `xml:"xades:SignedProperties"`
	UnsignedProperties *XAdESUnsignedProps `xml:"xades:UnsignedProperties,omitempty"`
}

type XAdESSignedProps struct {
	ID                   string              `xml:"Id,attr"`
	SignedSignatureProps XAdESSignedSigProps `xml:"xades:SignedSignatureProperties"`
}

type XAdESSignedSigProps struct {
	SigningTime               string                `xml:"xades:SigningTime"`
	SigningCertificate        XAdESSigningCert      `xml:"xades:SigningCertificate"`
	SignaturePolicyIdentifier XAdESPolicyIdentifier `xml:"xades:SignaturePolicyIdentifier"`
}

type XAdESSigningCert struct {
	Cert XAdESCert `xml:"xades:Cert"`
}

type XAdESCert struct {
	CertDigest   XAdESDigest       `xml:"xades:CertDigest"`
	IssuerSerial XAdESIssuerSerial `xml:"xades:IssuerSerial"`
}

type XAdESDigest struct {
	DigestMethod XAdESMethod `xml:"ds:DigestMethod"`
	DigestValue  string      `xml:"ds:DigestValue"`
}

type XAdESPolicyIdentifier struct {
	// ImpliedPolicy = nessun identificativo di policy (XAdES-BES baseline)
	ImpliedPolicy *struct{} `xml:"xades:ImpliedPolicy,omitempty"`
}

// XAdESUnsignedProps contiene il timestamp RFC 3161 (livello XAdES-T)
type XAdESUnsignedProps struct {
	UnsignedSigProps XAdESUnsignedSigProps `xml:"xades:UnsignedSignatureProperties"`
}

type XAdESUnsignedSigProps struct {
	SignatureTimeStamp XAdESTimestamp `xml:"xades:SignatureTimeStamp"`
}

type XAdESTimestamp struct {
	// EncapsulatedTimeStamp contiene il token RFC 3161 in Base64 (DER encoded)
	EncapsulatedTimeStamp string `xml:"xades:EncapsulatedTimeStamp"`
}

// BuildXAdESEnvelope costruisce un envelope XAdES-BES (senza timestamp) o XAdES-T (con timestamp).
// documentID:  URI del documento firmato
// docHashHex:  SHA-256 del documento in hex
// sigHex:      valore firma RSA PKCS1v15 in hex
// certDN:      Subject Distinguished Name del firmatario
// tsaToken:    token RFC 3161 in hex (se vuoto → XAdES-BES, altrimenti XAdES-T)
// signedAt:    orario di firma
func BuildXAdESEnvelope(documentID, docHashHex, sigHex, certDN, tsaToken string, signedAt time.Time) string {
	var unsignedProps *XAdESUnsignedProps
	if tsaToken != "" {
		unsignedProps = &XAdESUnsignedProps{
			UnsignedSigProps: XAdESUnsignedSigProps{
				SignatureTimeStamp: XAdESTimestamp{
					EncapsulatedTimeStamp: tsaToken,
				},
			},
		}
	}

	sig := XAdESSignature{
		Xmlns:      "http://www.w3.org/2000/09/xmldsig#",
		XmlnsXades: "http://uri.etsi.org/01903/v1.3.2#",
		ID:         fmt.Sprintf("Signature-%s", documentID),
		SignedInfo: XAdESSignedInfo{
			CanonicalizationMethod: XAdESMethod{Algorithm: "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"},
			SignatureMethod:        XAdESMethod{Algorithm: "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"},
			Reference: XAdESRef{
				URI:          "#" + documentID,
				DigestMethod: XAdESMethod{Algorithm: "http://www.w3.org/2001/04/xmlenc#sha256"},
				DigestValue:  docHashHex,
			},
		},
		SignatureValue: XAdESSignatureValue{
			ID:    fmt.Sprintf("SignatureValue-%s", documentID),
			Value: sigHex,
		},
		KeyInfo: XAdESKeyInfo{
			X509Data: XAdESX509Data{
				X509SubjectName: certDN,
				X509IssuerSerial: XAdESIssuerSerial{
					X509IssuerName:   certDN,
					X509SerialNumber: "0",
				},
			},
		},
		Object: XAdESObject{
			QualifyingProperties: XAdESQualifyingProps{
				Target: fmt.Sprintf("#Signature-%s", documentID),
				SignedProperties: XAdESSignedProps{
					ID: fmt.Sprintf("SignedProperties-%s", documentID),
					SignedSignatureProps: XAdESSignedSigProps{
						SigningTime: signedAt.Format(time.RFC3339),
						SigningCertificate: XAdESSigningCert{
							Cert: XAdESCert{
								CertDigest: XAdESDigest{
									DigestMethod: XAdESMethod{Algorithm: "http://www.w3.org/2001/04/xmlenc#sha256"},
									DigestValue:  docHashHex,
								},
								IssuerSerial: XAdESIssuerSerial{
									X509IssuerName:   certDN,
									X509SerialNumber: "0",
								},
							},
						},
						SignaturePolicyIdentifier: XAdESPolicyIdentifier{ImpliedPolicy: &struct{}{}},
					},
				},
				UnsignedProperties: unsignedProps,
			},
		},
	}

	xmlBytes, err := xml.MarshalIndent(sig, "", "  ")
	if err != nil {
		return fmt.Sprintf("<!-- XAdES build error: %v -->", err)
	}
	return string(xmlBytes)
}
