package protocol

import (
	"testing"
	"time"
)

func BenchmarkGenerateSegnaturaXML(b *testing.B) {
	entry := ProtocolEntry{
		ProtocolNumber:      1450,
		ProtocolYear:        2026,
		ProtocolDate:        time.Date(2026, 3, 15, 10, 30, 0, 0, time.UTC),
		FlowDirection:       "in",
		ClassificationTitle: 2,
		ClassificationClass: "04",
		Subject:             "Iscrizione alunno e trasmissione fascicolo personale",
		Sender:              "Istituto Comprensivo G. Pascoli",
		Recipient:           "Liceo Statale A. Manzoni",
		DocumentHashSHA256:  "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	}

	for b.Loop() {
		_, err := GenerateSegnaturaXML(entry)
		if err != nil {
			b.Fatalf("benchmark failed: %v", err)
		}
	}
}

func BenchmarkFormatVisualStamp(b *testing.B) {
	entry := ProtocolEntry{
		ProtocolNumber: 1450,
		ProtocolDate:   time.Date(2026, 3, 15, 10, 30, 0, 0, time.UTC),
	}
	school := "Liceo Statale A. Manzoni"

	for b.Loop() {
		_ = FormatVisualStamp(entry, school)
	}
}
