package albopretorio

import (
	"encoding/xml"
	"time"
)

// Standard categories mandated by public education legal transparency standards
const (
	CategoryBandiGarePNRR       = "Bandi e Gare PNRR"
	CategoryDetermineDirigente  = "Determine Dirigenziali"
	CategoryDecreti             = "Decreti"
	CategoryDelibereConsiglio   = "Delibere Consiglio di Istituto"
	CategoryCircolariEsterne    = "Circolari Esterne"
	CategoryGraduatorieIstituto = "Graduatorie di Istituto"
)

type AlboItem struct {
	ID                      string     `json:"id" db:"id"`
	SchoolID                string     `json:"school_id" db:"school_id"`
	RepertoryYear           int        `json:"repertory_year" db:"repertory_year"`
	RepertoryNumber         int        `json:"repertory_number" db:"repertory_number"`
	RepertoryCode           string     `json:"repertory_code" db:"repertory_code"` // e.g. "2026/00001"
	Category                string     `json:"category" db:"category"`
	Subject                 string     `json:"subject" db:"subject"`
	PublishedAt             time.Time  `json:"published_at" db:"published_at"`
	ExpiresAt               time.Time  `json:"expires_at" db:"expires_at"` // PublishedAt + 15 solar days
	Status                  string     `json:"status" db:"status"`         // "in_pubblicazione", "defisso", "annullato"
	DocumentFileURL         string     `json:"document_file_url" db:"document_file_url"`
	DocumentSHA256          string     `json:"document_sha256" db:"document_sha256"`
	PublishedBy             string     `json:"published_by" db:"published_by"`
	RelataText              string     `json:"relata_text,omitempty" db:"relata_text"`
	RelataSignedBy          string     `json:"relata_signed_by,omitempty" db:"relata_signed_by"`
	RelataSignedAt          *time.Time `json:"relata_signed_at,omitempty" db:"relata_signed_at"`
	IsTransparencySection   bool       `json:"is_transparency_section" db:"is_transparency_section"`
	TransparencyMacroFamily string     `json:"transparency_macro_family,omitempty" db:"transparency_macro_family"`
	TransparencySubFamily   string     `json:"transparency_sub_family,omitempty" db:"transparency_sub_family"`
	CIGCode                 string     `json:"cig_code,omitempty" db:"cig_code"`
	AwardedAmount           float64    `json:"awarded_amount,omitempty" db:"awarded_amount"`
	CreatedAt               time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at" db:"updated_at"`
}

type PublishActRequest struct {
	Category                string  `json:"category" binding:"required"`
	Subject                 string  `json:"subject" binding:"required"`
	DocumentFileURL         string  `json:"document_file_url" binding:"required"`
	DocumentSHA256          string  `json:"document_sha256" binding:"required"`
	DaysDuration            int     `json:"days_duration"` // Defaults to 15 days as per Legge 69/2009
	IsTransparencySection   bool    `json:"is_transparency_section"`
	TransparencyMacroFamily string  `json:"transparency_macro_family"`
	TransparencySubFamily   string  `json:"transparency_sub_family"`
	CIGCode                 string  `json:"cig_code"`
	AwardedAmount           float64 `json:"awarded_amount"`
}

type CertificatoPubblicazione struct {
	RepertoryCode    string    `json:"repertory_code"`
	Subject          string    `json:"subject"`
	Category         string    `json:"category"`
	PublishedAt      time.Time `json:"published_at"`
	DefissoAt        time.Time `json:"defisso_at"`
	SolarDays        int       `json:"solar_days"`
	DocumentSHA256   string    `json:"document_sha256"`
	LegalAttestation string    `json:"legal_attestation"`
	DirigenteName    string    `json:"dirigente_name"`
	CertifiedAt      time.Time `json:"certified_at"`
}

// ANAC Legge 190/2012 XML Export Structs
type ANACDataset struct {
	XMLName  xml.Name     `xml:"legge190:dati"`
	XMLNS    string       `xml:"xmlns:legge190,attr"`
	Metadata ANACMetadata `xml:"metadata"`
	Lotti    []ANACLotto  `xml:"data>lotto"`
}

type ANACMetadata struct {
	Titolo            string `xml:"titolo"`
	CodiceFiscale     string `xml:"codiceFiscale"`
	AnnoRiferimento   int    `xml:"annoRiferimento"`
	DataPubblicazione string `xml:"dataPubblicazione"`
}

type ANACLotto struct {
	CIG                      string  `xml:"cig"`
	StrutturaProponente      string  `xml:"strutturaProponente>denominazione"`
	Oggetto                  string  `xml:"oggetto"`
	SceltaContraente         string  `xml:"sceltaContraente"`
	ImportoAggiudicazione    float64 `xml:"importoAggiudicazione"`
	TempiCompletamentoInizio string  `xml:"tempiCompletamento>dataInizio"`
	TempiCompletamentoFine   string  `xml:"tempiCompletamento>dataUltimazione"`
}
