package sidi

import (
	"encoding/xml"
	"time"
)

type ExportRecord struct {
	ID           string    `json:"id" db:"id"`
	SchoolID     string    `json:"school_id" db:"school_id"`
	ExportType   string    `json:"export_type" db:"export_type"`
	SchoolYear   string    `json:"school_year" db:"school_year"`
	FileName     string    `json:"file_name" db:"file_name"`
	Status       string    `json:"status" db:"status"`
	RecordsCount int       `json:"records_count" db:"records_count"`
	XMLContent   string    `json:"xml_content,omitempty" db:"xml_content"`
	CreatedBy    string    `json:"created_by" db:"created_by"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// Strutture XML conformi agli schemi XSD SIDI MIM
type FlussoSIDI struct {
	XMLName    xml.Name   `xml:"FlussoSIDI"`
	Versione   string     `xml:"versione,attr"`
	Testata    Testata    `xml:"Testata"`
	DatiScuola DatiScuola `xml:"DatiScuola"`
}

type Testata struct {
	CodiceScuola string    `xml:"CodiceScuola"`
	AnnoSco      string    `xml:"AnnoScolastico"`
	DataExport   time.Time `xml:"DataEsportazione"`
	TipoFlusso   string    `xml:"TipoFlusso"`
}

type DatiScuola struct {
	Denominazione string          `xml:"Denominazione"`
	Studenti      []StudenteSIDI  `xml:"Studenti>Studente"`
	Scrutini      []ScrutinioSIDI `xml:"Scrutini>Scrutinio,omitempty"`
}

type StudenteSIDI struct {
	CodiceSIDI    string `xml:"CodiceSIDI"`
	CodiceFiscale string `xml:"CodiceFiscale"`
	Cognome       string `xml:"Cognome"`
	Nome          string `xml:"Nome"`
	Classe        string `xml:"Classe"`
}

type ScrutinioSIDI struct {
	CodiceSIDI       string `xml:"CodiceSIDI"`
	Classe           string `xml:"Classe"`
	EsitoFinale      string `xml:"EsitoFinale"` // AMMESSO, NON_AMMESSO, SOSPESO_GIUDIZIO
	CreditiFormat    int    `xml:"CreditiFormativi,omitempty"`
	DebitiRecuperati bool   `xml:"DebitiRecuperati,omitempty"`
}

type GenerateSidiRequest struct {
	ExportType string `json:"export_type" binding:"required"` // ANS_ANAGRAFE, SCRUTINIO_GIUGNO, SCRUTINIO_SETTEMBRE_DEBITI
	SchoolYear string `json:"school_year"`
}

type SidiValidationResult struct {
	Valid          bool     `json:"valid"`
	TotalRecords   int      `json:"total_records"`
	MissingSidiIDs []string `json:"missing_sidi_ids"`
	Errors         []string `json:"errors"`
}

type SidiCooperationConfig struct {
	EndpointURL           string `json:"endpoint_url"`
	CodiceMeccanografico  string `json:"codice_meccanografico"`
	CertificatoPostazione string `json:"certificato_postazione"`
	Environment           string `json:"environment"` // "PRODUZIONE", "COLLAUDO"
}

type SyncSidiCodesResponse struct {
	TotalProcessed int               `json:"total_processed"`
	TotalUpdated   int               `json:"total_updated"`
	UpdatedCodes   map[string]string `json:"updated_codes"` // CF -> Codice SIDI
	ProtocolloMIM  string            `json:"protocollo_mim"`
	SyncTimestamp  time.Time         `json:"sync_timestamp"`
}

type PushScrutinyResultsRequest struct {
	SchoolYear string          `json:"school_year"`
	Sessione   string          `json:"sessione"` // "GIUGNO", "SETTEMBRE"
	ClassID    string          `json:"class_id"`
	Results    []ScrutinioSIDI `json:"results"`
}

type PushScrutinyResultsResponse struct {
	Status           string    `json:"status"` // "TRASMESSO_CON_SUCCESSO", "SCARTATO"
	ProtocolloMIM    string    `json:"protocollo_mim"`
	RecordsProcessed int       `json:"records_processed"`
	RecordsAccepted  int       `json:"records_accepted"`
	Errors           []string  `json:"errors"`
	TransmittedAt    time.Time `json:"transmitted_at"`
}
