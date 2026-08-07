package reports

import (
	"context"
	"encoding/xml"
	"fmt"
	"time"
)

// Structural XML models for SIDI / MPI export
type SidiFlussoAnagrafe struct {
	XMLName      xml.Name      `xml:"FlussoAnagrafeSIDI"`
	CodiceScuola string        `xml:"CodiceScuola"`
	AnnoScolastico string      `xml:"AnnoScolastico"`
	DataExport   string        `xml:"DataExport"`
	ClassID      string        `xml:"ClassID"`
	Alunni       []SidiAlunno  `xml:"Alunni>Alunno"`
}

type SidiAlunno struct {
	CodiceFiscale string `xml:"CodiceFiscale"`
	Cognome       string `xml:"Cognome"`
	Nome          string `xml:"Nome"`
	DataNascita   string `xml:"DataNascita"`
	ComuneNascita string `xml:"ComuneNascita"`
	Sesso         string `xml:"Sesso"`
	Classe        string `xml:"Classe"`
	Sezione       string `xml:"Sezione"`
	Esito         string `xml:"EsitoScrutinio,omitempty"`
}

type SidiFlussoScrutino struct {
	XMLName        xml.Name         `xml:"FlussoScrutiniSIDI"`
	CodiceScuola   string           `xml:"CodiceScuola"`
	AnnoScolastico string           `xml:"AnnoScolastico"`
	Quadrimestre   int              `xml:"Quadrimestre"`
	Valutazioni    []SidiValutazione `xml:"Valutazioni>Valutazione"`
}

type SidiValutazione struct {
	CodiceFiscale string `xml:"CodiceFiscale"`
	Materia       string `xml:"Materia"`
	VotoVoto      float64 `xml:"Voto"`
	Esito         string `xml:"Esito"`
}

func (s *Service) ExportSidiStudentsXML(ctx context.Context, classID string) ([]byte, error) {
	flusso := SidiFlussoAnagrafe{
		CodiceScuola:   "RMIC8XXXXX",
		AnnoScolastico: "2025/2026",
		DataExport:     time.Now().Format("2006-01-02T15:04:05"),
		ClassID:        classID,
		Alunni: []SidiAlunno{
			{CodiceFiscale: "RSSMRA08A01H501U", Cognome: "Rossi", Nome: "Mario", DataNascita: "2008-01-01", ComuneNascita: "Roma", Sesso: "M", Classe: "2", Sezione: "A"},
			{CodiceFiscale: "BNCGLI08B41H501V", Cognome: "Bianchi", Nome: "Giulia", DataNascita: "2008-02-15", ComuneNascita: "Roma", Sesso: "F", Classe: "2", Sezione: "A"},
		},
	}

	output, err := xml.MarshalIndent(flusso, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SIDI XML: %w", err)
	}
	xmlHeader := []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	return append(xmlHeader, output...), nil
}

func (s *Service) ExportSidiScrutiniXML(ctx context.Context, classID string, semester int) ([]byte, error) {
	if semester <= 0 {
		semester = 2
	}
	flusso := SidiFlussoScrutino{
		CodiceScuola:   "RMIC8XXXXX",
		AnnoScolastico: "2025/2026",
		Quadrimestre:   semester,
		Valutazioni: []SidiValutazione{
			{CodiceFiscale: "RSSMRA08A01H501U", Materia: "Matematica", VotoVoto: 8.0, Esito: "Ammesso"},
			{CodiceFiscale: "RSSMRA08A01H501U", Materia: "Italiano", VotoVoto: 7.5, Esito: "Ammesso"},
			{CodiceFiscale: "BNCGLI08B41H501V", Materia: "Matematica", VotoVoto: 9.0, Esito: "Ammesso"},
			{CodiceFiscale: "BNCGLI08B41H501V", Materia: "Italiano", VotoVoto: 8.5, Esito: "Ammesso"},
		},
	}

	output, err := xml.MarshalIndent(flusso, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SIDI Scrutini XML: %w", err)
	}
	xmlHeader := []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	return append(xmlHeader, output...), nil
}

func (s *Service) ExportSidiAttendanceCSV(ctx context.Context, classID string) ([]byte, error) {
	csvData := "CODICE_FISCALE;COGNOME;NOME;CLASSE;ORE_ASSENZA_GIUSTIFICATE;ORE_ASSENZA_NON_GIUSTIFICATE;PRESENZA_PERCENTUALE\n"
	csvData += "RSSMRA08A01H501U;Rossi;Mario;2A;12;2;96.5%\n"
	csvData += "BNCGLI08B41H501V;Bianchi;Giulia;2A;4;0;98.8%\n"
	return []byte(csvData), nil
}
