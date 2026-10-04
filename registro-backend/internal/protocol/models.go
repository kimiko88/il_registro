package protocol

import (
	"encoding/xml"
	"fmt"
	"time"
)

type ProtocolEntry struct {
	ID                     string    `json:"id"`
	SchoolID               string    `json:"school_id"`
	ProtocolYear           int       `json:"protocol_year"`
	ProtocolNumber         int       `json:"protocol_number"`
	ProtocolDate           time.Time `json:"protocol_date"`
	FlowDirection          string    `json:"flow_direction"`       // in, out, internal
	ClassificationTitle    int       `json:"classification_title"` // Titolo I..X
	ClassificationClass    string    `json:"classification_class"`
	ClassificationFascicle string    `json:"classification_fascicle,omitempty"`
	Subject                string    `json:"subject"`
	Sender                 string    `json:"sender"`
	Recipient              string    `json:"recipient"`
	DocumentHashSHA256     string    `json:"document_hash_sha256"`
	DocumentFileURL        string    `json:"document_file_url"`
	ProtocolledBy          string    `json:"protocolled_by"`
	VisualStamp            string    `json:"visual_stamp,omitempty"`
}

type EntityProtocolLink struct {
	ID         string    `json:"id"`
	ProtocolID string    `json:"protocol_id"`
	EntityType string    `json:"entity_type"` // circular, verbale, pagella, family_request, contract
	EntityID   string    `json:"entity_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type XMLSegnatura struct {
	XMLName      xml.Name        `xml:"Segnatura"`
	Intestazione XMLIntestazione `xml:"Intestazione"`
	Descrizione  XMLDescrizione  `xml:"Descrizione"`
}

type XMLIntestazione struct {
	Identificatore XMLIdentificatore `xml:"Identificatore"`
	Origine        XMLSoggetto       `xml:"Origine"`
	Destinazione   XMLSoggetto       `xml:"Destinazione"`
	Classifica     XMLClassifica     `xml:"Classifica"`
	Oggetto        string            `xml:"Oggetto"`
}

type XMLIdentificatore struct {
	Numero            string `xml:"Numero"`
	DataRegistrazione string `xml:"DataRegistrazione"`
	Flusso            string `xml:"Flusso"`
}

type XMLSoggetto struct {
	Denominazione string `xml:"Denominazione"`
}

type XMLClassifica struct {
	Titolo    int    `xml:"Titolo"`
	Classe    string `xml:"Classe"`
	Fascicolo string `xml:"Fascicolo,omitempty"`
}

type XMLDescrizione struct {
	Documento XMLDocumento `xml:"Documento"`
}

type XMLDocumento struct {
	Impronta XMLImpronta `xml:"Impronta"`
}

type XMLImpronta struct {
	Algoritmo string `xml:"algoritmo,attr"`
	Value     string `xml:",chardata"`
}

// GenerateSegnaturaXML outputs the official interoperability XML compliant with AgID guidelines
func GenerateSegnaturaXML(entry ProtocolEntry) ([]byte, error) {
	seg := XMLSegnatura{
		Intestazione: XMLIntestazione{
			Identificatore: XMLIdentificatore{
				Numero:            fmt.Sprintf("%07d", entry.ProtocolNumber),
				DataRegistrazione: entry.ProtocolDate.Format("2006-01-02T15:04:05"),
				Flusso:            entry.FlowDirection,
			},
			Origine: XMLSoggetto{
				Denominazione: entry.Sender,
			},
			Destinazione: XMLSoggetto{
				Denominazione: entry.Recipient,
			},
			Classifica: XMLClassifica{
				Titolo:    entry.ClassificationTitle,
				Classe:    entry.ClassificationClass,
				Fascicolo: entry.ClassificationFascicle,
			},
			Oggetto: entry.Subject,
		},
		Descrizione: XMLDescrizione{
			Documento: XMLDocumento{
				Impronta: XMLImpronta{
					Algoritmo: "SHA-256",
					Value:     entry.DocumentHashSHA256,
				},
			},
		},
	}

	return xml.MarshalIndent(seg, "", "  ")
}

// FormatVisualStamp creates a printable protocol stamp string
func FormatVisualStamp(entry ProtocolEntry, schoolName string) string {
	dateStr := entry.ProtocolDate.Format("02/01/2006")
	return fmt.Sprintf("%s - REG. UFF. PROT. N. %07d del %s", schoolName, entry.ProtocolNumber, dateStr)
}
