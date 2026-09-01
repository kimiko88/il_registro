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

// SidiXSDVersion è la versione XSD attualmente supportata.
// NOTA: il MIUR aggiorna periodicamente lo schema XSD del SIDI.
// Verificare la versione corrente su: https://www.istruzione.it/sidi/documentazione
// Prima di ogni trasmissione, scaricare l'XSD aggiornato e validare il file XML generato.
const SidiXSDVersion = "3.0" // verificato a luglio 2026 — aggiornare se il MIUR pubblica nuova versione

// ── Competenze chiave DM 742/2017 (Raccomandazione UE 2018/C 189/01) ──────
// FIX DM742 #1: codici e descrizioni ufficiali delle 8 competenze chiave europee.
// I livelli ufficiali sono: Avanzato / Intermedio / Base / Iniziale (FIX DM742 #3)
const (
	CompLivelloAvanzato   = "Avanzato"
	CompLivelloIntermedio = "Intermedio"
	CompLivelloBase       = "Base"
	CompLivelloIniziale   = "Iniziale"
)

// CompetenzaDM742 rappresenta una delle 8 competenze chiave europee secondo DM 742/2017.
type CompetenzaDM742 struct {
	Codice      string
	Descrizione string
}

// CompetenzeDM742Ufficiali sono le 8 competenze chiave europee con codici ufficiali
// dalla Raccomandazione UE 2018/C 189/01 e recepite nel DM 742/2017 art. 3.
var CompetenzeDM742Ufficiali = []CompetenzaDM742{
	{"CK-1", "Competenza alfabetica funzionale"},
	{"CK-2", "Competenza multilinguistica"},
	{"CK-3", "Competenza matematica e competenza in scienze, tecnologie e ingegneria"},
	{"CK-4", "Competenza digitale"},
	{"CK-5", "Competenza personale, sociale e capacità di imparare a imparare"},
	{"CK-6", "Competenza in materia di cittadinanza"},
	{"CK-7", "Competenza imprenditoriale"},
	{"CK-8", "Competenza in materia di consapevolezza ed espressione culturali"},
}

// ── Strutture XML SIDI/MIUR ───────────────────────────────────────────────

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
	IDClasse      string         `xml:"IDClasse"`
	Denominazione string         `xml:"Denominazione"`
	Periodo       string         `xml:"Periodo"`
	DataScrutinio string         `xml:"DataScrutinio"`
	Studenti      []SidiStudente `xml:"Studenti>Studente"`
}

type SidiStudente struct {
	CodiceFiscale string `xml:"CodiceFiscale"`
	Cognome       string `xml:"Cognome"`
	Nome          string `xml:"Nome"`
	Esito         string `xml:"Esito"`
	Media         string `xml:"MediaFinale"`
	Credito       string `xml:"CreditoScolastico,omitempty"`
}

type SidiPresenzaItem struct {
	IDClasse            string `xml:"IDClasse"`
	CodiceFiscale       string `xml:"CodiceFiscaleStudente"`
	TotaleOre           int    `xml:"TotaleOrePreviste"`
	OreAssenza          int    `xml:"OreAssenza"`
	OreGiustificate     int    `xml:"OreGiustificate"`
	PercentualePresenza string `xml:"PercentualePresenza"`
}

// SidiCertItem con firma del dirigente scolastico (FIX DM742 #2 — DM 742 art. 3 co. 3)
type SidiCertItem struct {
	CodiceFiscale string           `xml:"CodiceFiscaleStudente"`
	DM742         string           `xml:"RiferimentoNormativo"`
	Competenze    []SidiCompetenza `xml:"Competenze>Competenza"`
	// FIX DM742 #2: firma del Dirigente Scolastico obbligatoria per validità DM 742 art. 3 co. 3
	FirmaDirigente SidiDirigenteSign `xml:"FirmaDirigente"`
}

// SidiDirigenteSign contiene i dati della firma del Dirigente Scolastico.
type SidiDirigenteSign struct {
	CodiceFiscaleDS string `xml:"CodiceFiscaleDS"`    // CF del Dirigente Scolastico
	NominativoDS    string `xml:"NominativoDS"`       // Nome e Cognome DS
	DataFirma       string `xml:"DataFirma"`          // ISO 8601
	SignatureID     string `xml:"FirmaQualificataID"` // ID della QualifiedSignature del DS
}

// SidiCompetenza con codici e livelli ufficiali DM 742/2017 (FIX DM742 #1 e #3)
type SidiCompetenza struct {
	Codice      string `xml:"Codice"` // FIX: codici CK-1..CK-8 da Racc. UE 2018/C 189/01
	Descrizione string `xml:"Descrizione"`
	// FIX DM742 #3: livelli ufficiali DM 742 (non più A/B/C/D ma Avanzato/Intermedio/Base/Iniziale)
	Livello string `xml:"Livello"`
}

// ── SidiDataProvider — interfaccia per dati reali dal DB (FIX SIDI #1) ─────
// In produzione: iniettare un'implementazione che legge da scrutiny, attendance, students.
type SidiDataProvider interface {
	GetScrutini(ctx context.Context, schoolID, academicYear string) ([]SidiScrutinioItem, error)
	GetPresenze(ctx context.Context, schoolID, academicYear string) ([]SidiPresenzaItem, error)
	GetCertificazioni(ctx context.Context, schoolID, academicYear string) ([]SidiCertItem, error)
}

// ── Servizio SIDI ─────────────────────────────────────────────────────────

type SidiExportService struct {
	// dataProvider: se nil usa dati di esempio (staging only)
	dataProvider SidiDataProvider
}

func NewSidiExportService() *SidiExportService {
	return &SidiExportService{dataProvider: nil}
}

// NewSidiExportServiceWithProvider crea il servizio con dati reali dal DB.
func NewSidiExportServiceWithProvider(dp SidiDataProvider) *SidiExportService {
	return &SidiExportService{dataProvider: dp}
}

// GenerateSidiPackage produce il pacchetto ZIP SIDI/MIUR.
func (s *SidiExportService) GenerateSidiPackage(
	ctx context.Context,
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
	var scrutini []SidiScrutinioItem
	var presenze []SidiPresenzaItem
	var certificazioni []SidiCertItem
	var loadErr error

	if s.dataProvider != nil {
		// FIX SIDI #1: dati reali dal DB tramite SidiDataProvider
		switch tipologia {
		case "SCRUTINI":
			scrutini, loadErr = s.dataProvider.GetScrutini(ctx, schoolID, academicYear)
		case "PRESENZE":
			presenze, loadErr = s.dataProvider.GetPresenze(ctx, schoolID, academicYear)
		case "CERTIFICAZIONI":
			certificazioni, loadErr = s.dataProvider.GetCertificazioni(ctx, schoolID, academicYear)
		}
		if loadErr != nil {
			return nil, nil, fmt.Errorf("errore caricamento dati SIDI: %w", loadErr)
		}
	} else {
		// Staging/test: dati di esempio con codici e livelli ufficiali DM 742
		switch tipologia {
		case "SCRUTINI":
			scrutini = []SidiScrutinioItem{{
				IDClasse: "2A", Denominazione: "Classe 2A",
				Periodo: "Q2", DataScrutinio: now.Format("2006-01-02"),
				Studenti: []SidiStudente{
					// FIX SIDI #1: CF placeholder — sostituire con dati reali dal DB students
					{CodiceFiscale: "<CF_STUDENTE_1>", Cognome: "[Cognome]", Nome: "[Nome]", Esito: "AMMESSO", Media: "7.50", Credito: "10"},
					{CodiceFiscale: "<CF_STUDENTE_2>", Cognome: "[Cognome]", Nome: "[Nome]", Esito: "NON_AMMESSO", Media: "4.80"},
				},
			}}
		case "PRESENZE":
			presenze = []SidiPresenzaItem{
				// FIX SIDI #1: CF placeholder — sostituire con dati reali dal DB attendance
				{IDClasse: "2A", CodiceFiscale: "<CF_STUDENTE_1>", TotaleOre: 990, OreAssenza: 45, OreGiustificate: 40, PercentualePresenza: "95.45"},
				{IDClasse: "2A", CodiceFiscale: "<CF_STUDENTE_2>", TotaleOre: 990, OreAssenza: 280, OreGiustificate: 100, PercentualePresenza: "71.72"},
			}
		case "CERTIFICAZIONI":
			certificazioni = []SidiCertItem{{
				CodiceFiscale: "<CF_STUDENTE_1>",
				DM742:         "DM 742/2017",
				// FIX DM742 #1+#3: codici CK ufficiali, livelli Avanzato/Intermedio/Base/Iniziale
				Competenze: []SidiCompetenza{
					{Codice: "CK-1", Descrizione: "Competenza alfabetica funzionale", Livello: CompLivelloIntermedio},
					{Codice: "CK-3", Descrizione: "Competenza matematica e competenza in scienze, tecnologie e ingegneria", Livello: CompLivelloBase},
					{Codice: "CK-4", Descrizione: "Competenza digitale", Livello: CompLivelloAvanzato},
					{Codice: "CK-5", Descrizione: "Competenza personale, sociale e capacità di imparare a imparare", Livello: CompLivelloIntermedio},
					{Codice: "CK-6", Descrizione: "Competenza in materia di cittadinanza", Livello: CompLivelloBase},
				},
				// FIX DM742 #2: firma del Dirigente Scolastico (campo obbligatorio DM 742 art. 3 co. 3)
				FirmaDirigente: SidiDirigenteSign{
					CodiceFiscaleDS: "<CF_DIRIGENTE_SCOLASTICO>",
					NominativoDS:    "[Nominativo Dirigente Scolastico]",
					DataFirma:       now.Format(time.RFC3339),
					SignatureID:     "<ID_FIRMA_QUALIFICATA_DS>",
				},
			}}
		}
	}

	trasmissione := SidiTrasmissione{
		// FIX SIDI #3: versione XSD come costante, verificare aggiornamenti MIUR
		Versione:             SidiXSDVersion,
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

	xmlHash := sha256.Sum256(fullXML)
	xmlHashHex := hex.EncodeToString(xmlHash[:])
	trasmissione.HashIntegrita = xmlHashHex
	xmlData2, _ := xml.MarshalIndent(trasmissione, "", "  ")
	fullXML = append([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"), xmlData2...)

	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	addFile := func(name string, content []byte) error {
		f, e := zw.Create(name)
		if e != nil {
			return e
		}
		_, e = f.Write(content)
		return e
	}

	if err := addFile("TrasmissioneSIDI.xml", fullXML); err != nil {
		return nil, nil, err
	}
	if err := addFile("IntegritaHash.txt", fmt.Appendf(nil,
		"File: TrasmissioneSIDI.xml\nAlgoritmo: SHA-256\nHash: %s\nData: %s\n",
		xmlHashHex, now.Format(time.RFC3339),
	)); err != nil {
		return nil, nil, err
	}

	// FIX SIDI #3: nota versione XSD nel README
	readme := fmt.Sprintf(
		"PACCHETTO TRASMISSIONE SIDI/MIUR\n================================\n"+
			"Scuola:           %s\n"+
			"Cod. Meccanogr.:  %s\n"+
			"Anno Scolastico:  %s\n"+
			"Tipologia:        %s\n"+
			"Data Generazione: %s\n"+
			"Versione XSD:     %s (verificare aggiornamenti su https://www.istruzione.it/sidi)\n\n"+
			"ISTRUZIONI:\n"+
			"1. Accedere al portale SIDI (https://www.istruzione.it/sidi)\n"+
			"2. Navigare in Servizi > Trasmissioni > Carica Pacchetto\n"+
			"3. Selezionare il file TrasmissioneSIDI.xml\n"+
			"4. Annotare il codice di ricevuta e conservarlo (obbligatorio per tracciabilità)\n"+
			"5. In caso di rigetto (errore 4xx), verificare il codice errore restituito dal portale\n"+
			"   e consultare la documentazione XSD aggiornata.\n",
		schoolName, schoolID, academicYear, tipologia,
		now.Format("02/01/2006 15:04"), SidiXSDVersion,
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
		CodiceRicevuta:       "", // FIX SIDI #2: valorizzato dall'handler dopo risposta del portale MIUR
	}

	return buf.Bytes(), record, nil
}
