package signatures

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"time"
)

// ── Strutture XML SIDI/MIUR ────────────────────────────────────────────────

// SidiTrasmissione è il documento XML radice per la trasmissione al SIDI.
type SidiTrasmissione struct {
	XMLName              xml.Name            `xml:"TrasmissioneSIDI"`
	Versione             string              `xml:"versione,attr"`
	Xmlns                string              `xml:"xmlns,attr"`
	CodiceMeccanografico string              `xml:"Intestazione>CodiceMeccanografico"`
	DenominazioneScuola  string              `xml:"Intestazione>DenominazioneScuola"`
	AnnoScolastico       string              `xml:"Intestazione>AnnoScolastico"`
	DataGenerazione      string              `xml:"Intestazione>DataGenerazione"`
	Tipologia            string              `xml:"Intestazione>Tipologia"`
	HashIntegrita        string              `xml:"Intestazione>HashIntegritaSHA256"`
	Scrutini             []SidiScrutinioItem `xml:"Dati>Scrutini>Scrutinio,omitempty"`
	Presenze             []SidiPresenzaItem  `xml:"Dati>Presenze>PresenzaClasse,omitempty"`
	Certificazioni       []SidiCertItem      `xml:"Dati>Certificazioni>Certificazione,omitempty"`
}

type SidiScrutinioItem struct {
	IDClasse      string        `xml:"IDClasse"`
	Denominazione string        `xml:"Denominazione"`
	Periodo       string        `xml:"Periodo"` // "Q1", "Q2", "FINALE"
	DataScrutinio string        `xml:"DataScrutinio"`
	Studenti       []SidiStudente `xml:"Studenti>Studente"`
}

type SidiStudente struct {
	CodiceFiscale string `xml:"CodiceFiscale"`
	Cognome       string `xml:"Cognome"`
	Nome          string `xml:"Nome"`
	Esito         string `xml:"Esito"`   // "AMMESSO", "NON_AMMESSO", "SOSPESO"
	Media         string `xml:"MediaFinale"`
	Credito       string `xml:"CreditoScolastico,omitempty"`
}

type SidiPresenzaItem struct {
	IDClasse        string `xml:"IDClasse"`
	CodiceFiscale   string `xml:"CodiceFiscaleStudente"`
	TotaleOre       int    `xml:"TotaleOrePreviste"`
	OreAssenza      int    `xml:"OreAssenza"`
	OreGiustificate int    `xml:"OreGiustificate"`
	PercentualePresenza string `xml:"PercentualePresenza"`
}

type SidiCertItem struct {
	CodiceFiscale string            `xml:"CodiceFiscaleStudente"`
	DM742         string            `xml:"RiferimentoNormativo"` // "DM 742/2017"
	Competenze    []SidiCompetenza  `xml:"Competenze>Competenza"`
}

type SidiCompetenza struct {
	Codice      string `xml:"Codice"`
	Descrizione string `xml:"Descrizione"`
	Livello     string `xml:"Livello"` // "A","B","C","D"
}

// ── Servizio SIDI ─────────────────────────────────────────────────────────

// SidiExportService genera il pacchetto di trasmissione SIDI/MIUR.
type SidiExportService struct{}

func NewSidiExportService() *SidiExportService {
	return &SidiExportService{}
}

// GenerateSidiPackage produce uno ZIP contenente:
//   - TrasmissioneSIDI.xml  (dati scrutini/presenze/certificazioni)
//   - IntegritaHash.txt     (SHA-256 del file XML per verifica MIUR)
//   - ReadMe_SIDI.txt       (istruzioni di trasmissione)
func (s *SidiExportService) GenerateSidiPackage(
	_ context.Context,
	schoolID, schoolName, academicYear, tipologia string,
) ([]byte, *SidiExportRecord, error) {

	if schoolID == "" || academicYear == "" {
		return nil, nil, fmt.Errorf("schoolID e academicYear sono obbligatori")
	}

	tipologiaValida := map[string]bool{"SCRUTINI": true, "PRESENZE": true, "CERTIFICAZIONI": true}
	if !tipologiaValida[tipologia] {
		return nil, nil, fmt.Errorf("tipologia non valida: %s (attese: SCRUTINI, PRESENZE, CERTIFICAZIONI)", tipologia)
	}

	now := time.Now().UTC()

	// Dati di esempio; in produzione questi vengono caricati da DB tramite i rispettivi repository.
	var scrutini []SidiScrutinioItem
	var presenze []SidiPresenzaItem
	var certificazioni []SidiCertItem

	switch tipologia {
	case "SCRUTINI":
		scrutini = []SidiScrutinioItem{
			{
				IDClasse: "2A", Denominazione: "Classe 2A",
				Periodo: "Q2", DataScrutinio: now.Format("2006-01-02"),
				Studenti: []SidiStudente{
					{CodiceFiscale: "RSSMRA10A01F205Z", Cognome: "Rossi", Nome: "Mario", Esito: "AMMESSO", Media: "7.50", Credito: "10"},
					{CodiceFiscale: "VRDLGI10B02F205Z", Cognome: "Verdi", Nome: "Luigi", Esito: "NON_AMMESSO", Media: "4.80"},
				},
			},
		}
	case "PRESENZE":
		presenze = []SidiPresenzaItem{
			{IDClasse: "2A", CodiceFiscale: "RSSMRA10A01F205Z", TotaleOre: 990, OreAssenza: 45, OreGiustificate: 40, PercentualePresenza: "95.45"},
			{IDClasse: "2A", CodiceFiscale: "VRDLGI10B02F205Z", TotaleOre: 990, OreAssenza: 280, OreGiustificate: 100, PercentualePresenza: "71.72"},
		}
	case "CERTIFICAZIONI":
		certificazioni = []SidiCertItem{
			{
				CodiceFiscale: "RSSMRA10A01F205Z",
				DM742: "DM 742/2017",
				Competenze: []SidiCompetenza{
					{Codice: "COM-IT-1", Descrizione: "Comunicazione in madrelingua", Livello: "B"},
					{Codice: "COM-MAT-2", Descrizione: "Competenza matematica", Livello: "C"},
					{Codice: "DIG-3", Descrizione: "Competenza digitale", Livello: "A"},
				},
			},
		}
	}

	trasmissione := SidiTrasmissione{
		Versione:             "3.0",
		Xmlns:                "http://www.istruzione.it/sidi/trasmissione/2024",
		CodiceMeccanografico: schoolID,
		DenominazioneScuola:  schoolName,
		AnnoScolastico:       academicYear,
		DataGenerazione:      now.Format(time.RFC3339),
		Tipologia:            tipologia,
		Scrutini:             scrutini,
		Presenze:             presenze,
		Certificazioni:       certificazioni,
	}

	xmlData, err := xml.MarshalIndent(trasmissione, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("errore serializzazione XML SIDI: %w", err)
	}
	fullXML := append([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"), xmlData...)

	// Hash di integrità del payload XML
	xmlHash := sha256.Sum256(fullXML)
	xmlHashHex := hex.EncodeToString(xmlHash[:])
	trasmissione.HashIntegrita = xmlHashHex

	// Ri-serializza con hash incluso
	xmlData2, _ := xml.MarshalIndent(trasmissione, "", "  ")
	fullXML = append([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"), xmlData2...)

	// Costruisce ZIP
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	addFile := func(name string, content []byte) error {
		f, err := zw.Create(name)
		if err != nil {
			return err
		}
		_, err = f.Write(content)
		return err
	}

	if err := addFile("TrasmissioneSIDI.xml", fullXML); err != nil {
		return nil, nil, err
	}
	if err := addFile("IntegritaHash.txt", []byte(fmt.Sprintf(
		"File: TrasmissioneSIDI.xml\nAlgoritmo: SHA-256\nHash: %s\nData: %s\n",
		xmlHashHex, now.Format(time.RFC3339),
	))); err != nil {
		return nil, nil, err
	}
	readme := fmt.Sprintf(
		"PACCHETTO TRASMISSIONE SIDI/MIUR\n================================\n"+
			"Scuola:           %s\n"+
			"Cod. Meccanogr.:  %s\n"+
			"Anno Scolastico:  %s\n"+
			"Tipologia:        %s\n"+
			"Data Generazione: %s\n\n"+
			"ISTRUZIONI:\n"+
			"1. Accedere al portale SIDI (https://www.istruzione.it/sidi)\n"+
			"2. Navigare in Servizi > Trasmissioni > Carica Pacchetto\n"+
			"3. Selezionare il file TrasmissioneSIDI.xml\n"+
			"4. Verificare il codice di ricevuta e conservarlo\n",
		schoolName, schoolID, academicYear, tipologia, now.Format("02/01/2006 15:04"),
	)
	if err := addFile("ReadMe_SIDI.txt", []byte(readme)); err != nil {
		return nil, nil, err
	}

	_ = zw.Close()

	record := &SidiExportRecord{
		CodiceMeccanografico: schoolID,
		AnnoScolastico:       academicYear,
		TipologiaExport:      tipologia,
		DataGenerazione:      now.Format(time.RFC3339),
		HashIntegrità:        xmlHashHex,
		StatoTrasmissione:    "PRONTO",
	}

	return buf.Bytes(), record, nil
}
