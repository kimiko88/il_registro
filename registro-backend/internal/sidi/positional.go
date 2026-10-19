package sidi

import (
	"fmt"
	"strings"
	"time"
)

// Positional Export Types
const (
	PositionalTypeEsiti   = "ESITI"
	PositionalTypePagelle = "PAGELLE"
	PositionalTypeANS     = "ANS"
	PositionalTypePassWeb = "PASSWEB"
)

// Standard Record Lengths
const (
	RecordLenEsiti   = 120
	RecordLenPagelle = 100
	RecordLenANS     = 140
	RecordLenPassWeb = 160
)

type EsitoPosizionaleRecord struct {
	CodiceScuola     string `json:"codice_scuola"`     // 10 chars
	AnnoScolastico   string `json:"anno_scolastico"`   // 6 chars (es. 202526)
	CodiceFiscale    string `json:"codice_fiscale"`    // 16 chars
	CodiceSIDI       string `json:"codice_sidi"`       // 12 chars
	ClasseSezione    string `json:"classe_sezione"`    // 5 chars
	Sessione         string `json:"sessione"`          // 3 chars (GIU, SET)
	Esito            string `json:"esito"`             // 1 char (A=Ammesso, N=Non ammesso, S=Sospeso)
	Crediti          int    `json:"crediti"`           // 2 chars (00-12)
	DebitiAssegnati  int    `json:"debiti_assegnati"`  // 2 chars (00-99)
	DebitiRecuperati string `json:"debiti_recuperati"` // 1 char (S/N)
	MediaVoti        string `json:"media_voti"`        // 4 chars (es. 0750 per 7.50)
	Note             string `json:"note"`              // 58 chars
}

type PagellaPosizionaleRecord struct {
	CodiceScuola   string `json:"codice_scuola"`   // 10 chars
	AnnoScolastico string `json:"anno_scolastico"` // 6 chars
	Periodo        string `json:"periodo"`         // 2 chars (1Q, 2Q)
	CodiceFiscale  string `json:"codice_fiscale"`  // 16 chars
	CodiceMateria  string `json:"codice_materia"`  // 6 chars
	Voto           string `json:"voto"`            // 3 chars (es. 080)
	Condotta       string `json:"condotta"`        // 2 chars (es. 09)
	AssenzeOre     int    `json:"assenze_ore"`     // 3 chars
	Giudizio       string `json:"giudizio"`        // 20 chars
	Note           string `json:"note"`            // 32 chars
}

type ANSPosizionaleRecord struct {
	CodiceScuola    string `json:"codice_scuola"`    // 10 chars
	AnnoScolastico  string `json:"anno_scolastico"`  // 6 chars
	CodiceFiscale   string `json:"codice_fiscale"`   // 16 chars
	CodiceSIDI      string `json:"codice_sidi"`      // 12 chars
	Cognome         string `json:"cognome"`          // 30 chars
	Nome            string `json:"nome"`             // 30 chars
	Sesso           string `json:"sesso"`            // 1 char (M/F)
	DataNascita     string `json:"data_nascita"`     // 8 chars (YYYYMMDD)
	ComuneNascita   string `json:"comune_nascita"`   // 4 chars (Codice catastale)
	Cittadinanza    string `json:"cittadinanza"`     // 3 chars (es. ITA)
	ClasseSezione   string `json:"classe_sezione"`   // 5 chars
	IndirizzoStudi  string `json:"indirizzo_studi"`  // 6 chars (es. LI01)
	StatoIscrizione string `json:"stato_iscrizione"` // 3 chars (ATT, TRA, RIT)
	Filler          string `json:"filler"`           // 6 chars
}

type PassWebPosizionaleRecord struct {
	CodiceFiscale         string `json:"codice_fiscale"`          // 16 chars
	CodiceScuola          string `json:"codice_scuola"`           // 10 chars
	Profilo               string `json:"profilo"`                 // 10 chars (DOCENTE, COLL_SCOL, ASS_AMM, DSGA)
	TipoContratto         string `json:"tipo_contratto"`          // 2 chars (TI, TD)
	DataInizio            string `json:"data_inizio"`             // 8 chars (YYYYMMDD)
	DataFine              string `json:"data_fine"`               // 8 chars (YYYYMMDD)
	RegimePrev            string `json:"regime_prev"`             // 10 chars (TFS, TFR, SPERO)
	RetribuzioneLordaCent int64  `json:"retribuzione_lorda_cent"` // 10 chars (es. 0003200000 per 32.000,00)
	RetribuzioneTabellare int64  `json:"retribuzione_tabellare"`  // 10 chars
	GiorniUtili           int    `json:"giorni_utili"`            // 4 chars
	CausaCessazione       string `json:"causa_cessazione"`        // 16 chars (DIMISSIONI, COLLOC_RIPOSO, SCADENZA_TERMINE)
	Filler                string `json:"filler"`                  // 56 chars
}

type PositionalExportRequest struct {
	Type       string                     `json:"type" binding:"required"` // ESITI, PAGELLE, ANS, PASSWEB
	SchoolYear string                     `json:"school_year"`
	Esiti      []EsitoPosizionaleRecord   `json:"esiti,omitempty"`
	Pagelle    []PagellaPosizionaleRecord `json:"pagelle,omitempty"`
	ANS        []ANSPosizionaleRecord     `json:"ans,omitempty"`
	PassWeb    []PassWebPosizionaleRecord `json:"passweb,omitempty"`
}

type PositionalExportResponse struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"`
	FileName       string    `json:"file_name"`
	TotalLines     int       `json:"total_lines"`
	LineLength     int       `json:"line_length"`
	ValidationOK   bool      `json:"validation_ok"`
	Errors         []string  `json:"errors,omitempty"`
	ContentPreview string    `json:"content_preview"`
	CreatedAt      time.Time `json:"created_at"`
}

// Helpers for fixed-width padding
func padRight(s string, length int) string {
	s = strings.TrimSpace(s)
	if len(s) >= length {
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

func padLeft(s string, length int, padChar byte) string {
	s = strings.TrimSpace(s)
	if len(s) >= length {
		return s[:length]
	}
	return strings.Repeat(string(padChar), length-len(s)) + s
}

// BuildEsitiTXT builds the 120-character fixed length ESITI.TXT file
func BuildEsitiTXT(records []EsitoPosizionaleRecord) (string, []string) {
	var sb strings.Builder
	var errs []string

	for i, r := range records {
		line := padRight(r.CodiceScuola, 10) +
			padRight(r.AnnoScolastico, 6) +
			padRight(r.CodiceFiscale, 16) +
			padRight(r.CodiceSIDI, 12) +
			padRight(r.ClasseSezione, 5) +
			padRight(r.Sessione, 3) +
			padRight(r.Esito, 1) +
			padLeft(fmt.Sprintf("%02d", r.Crediti), 2, '0') +
			padLeft(fmt.Sprintf("%02d", r.DebitiAssegnati), 2, '0') +
			padRight(r.DebitiRecuperati, 1) +
			padRight(r.MediaVoti, 4) +
			padRight(r.Note, 58)

		if len(line) != RecordLenEsiti {
			errs = append(errs, fmt.Sprintf("Riga %d (Esiti): lunghezza errata %d anziché %d", i+1, len(line), RecordLenEsiti))
		}
		sb.WriteString(line)
		sb.WriteString("\r\n")
	}
	return sb.String(), errs
}

// BuildPagelleTXT builds the 100-character fixed length PAGELLE.TXT file
func BuildPagelleTXT(records []PagellaPosizionaleRecord) (string, []string) {
	var sb strings.Builder
	var errs []string

	for i, r := range records {
		line := padRight(r.CodiceScuola, 10) +
			padRight(r.AnnoScolastico, 6) +
			padRight(r.Periodo, 2) +
			padRight(r.CodiceFiscale, 16) +
			padRight(r.CodiceMateria, 6) +
			padRight(r.Voto, 3) +
			padRight(r.Condotta, 2) +
			padLeft(fmt.Sprintf("%03d", r.AssenzeOre), 3, '0') +
			padRight(r.Giudizio, 20) +
			padRight(r.Note, 32)

		if len(line) != RecordLenPagelle {
			errs = append(errs, fmt.Sprintf("Riga %d (Pagelle): lunghezza errata %d anziché %d", i+1, len(line), RecordLenPagelle))
		}
		sb.WriteString(line)
		sb.WriteString("\r\n")
	}
	return sb.String(), errs
}

// BuildANSTXT builds the 140-character fixed length ANS.TXT file
func BuildANSTXT(records []ANSPosizionaleRecord) (string, []string) {
	var sb strings.Builder
	var errs []string

	for i, r := range records {
		line := padRight(r.CodiceScuola, 10) +
			padRight(r.AnnoScolastico, 6) +
			padRight(r.CodiceFiscale, 16) +
			padRight(r.CodiceSIDI, 12) +
			padRight(r.Cognome, 30) +
			padRight(r.Nome, 30) +
			padRight(r.Sesso, 1) +
			padRight(r.DataNascita, 8) +
			padRight(r.ComuneNascita, 4) +
			padRight(r.Cittadinanza, 3) +
			padRight(r.ClasseSezione, 5) +
			padRight(r.IndirizzoStudi, 6) +
			padRight(r.StatoIscrizione, 3) +
			padRight(r.Filler, 6)

		if len(line) != RecordLenANS {
			errs = append(errs, fmt.Sprintf("Riga %d (ANS): lunghezza errata %d anziché %d", i+1, len(line), RecordLenANS))
		}
		sb.WriteString(line)
		sb.WriteString("\r\n")
	}
	return sb.String(), errs
}

// BuildPassWebTXT builds the 160-character fixed length PASSWEB_INPS.TXT file
func BuildPassWebTXT(records []PassWebPosizionaleRecord) (string, []string) {
	var sb strings.Builder
	var errs []string

	for i, r := range records {
		line := padRight(r.CodiceFiscale, 16) +
			padRight(r.CodiceScuola, 10) +
			padRight(r.Profilo, 10) +
			padRight(r.TipoContratto, 2) +
			padRight(r.DataInizio, 8) +
			padRight(r.DataFine, 8) +
			padRight(r.RegimePrev, 10) +
			padLeft(fmt.Sprintf("%010d", r.RetribuzioneLordaCent), 10, '0') +
			padLeft(fmt.Sprintf("%010d", r.RetribuzioneTabellare), 10, '0') +
			padLeft(fmt.Sprintf("%04d", r.GiorniUtili), 4, '0') +
			padRight(r.CausaCessazione, 16) +
			padRight(r.Filler, 56)

		if len(line) != RecordLenPassWeb {
			errs = append(errs, fmt.Sprintf("Riga %d (PassWeb): lunghezza errata %d anziché %d", i+1, len(line), RecordLenPassWeb))
		}
		sb.WriteString(line)
		sb.WriteString("\r\n")
	}
	return sb.String(), errs
}
