package protocol

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateAgIDSegnaturaXML(t *testing.T) {
	entry := ProtocolEntry{
		ProtocolYear:           2026,
		ProtocolNumber:         452,
		ProtocolDate:           time.Date(2026, 10, 4, 10, 0, 0, 0, time.UTC),
		FlowDirection:          "in",
		ClassificationTitle:    7, // Titolo VII - Alunni
		ClassificationClass:    "1",
		ClassificationFascicle: "Iscrizioni",
		Subject:                "Domanda di iscrizione alunno Rossi Mario",
		Sender:                 "Rossi Mario (Genitore)",
		Recipient:              "Liceo Statale Scientifico",
		DocumentHashSHA256:     "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	}

	xmlBytes, err := GenerateSegnaturaXML(entry)
	if err != nil {
		t.Fatalf("failed to generate Segnatura.xml: %v", err)
	}

	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, "<Segnatura") {
		t.Errorf("missing Segnatura root tag in XML")
	}
	if !strings.Contains(xmlStr, "<Numero>0000452</Numero>") {
		t.Errorf("missing zero-padded protocol number in XML: %s", xmlStr)
	}
	if !strings.Contains(xmlStr, "<Oggetto>Domanda di iscrizione alunno Rossi Mario</Oggetto>") {
		t.Errorf("missing Oggetto in XML")
	}
	if !strings.Contains(xmlStr, "<Impronta") {
		t.Errorf("missing document SHA-256 Impronta tag in XML")
	}
}

func TestFormatVisualProtocolStamp(t *testing.T) {
	entry := ProtocolEntry{
		ProtocolYear:   2026,
		ProtocolNumber: 89,
		ProtocolDate:   time.Date(2026, 10, 4, 0, 0, 0, 0, time.UTC),
	}

	stamp := FormatVisualStamp(entry, "I.C. LEONARDO DA VINCI")
	expected := "I.C. LEONARDO DA VINCI - REG. UFF. PROT. N. 0000089 del 04/10/2026"
	if stamp != expected {
		t.Errorf("expected stamp %q, got %q", expected, stamp)
	}
}
