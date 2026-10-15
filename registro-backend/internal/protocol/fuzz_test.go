package protocol

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"
	"time"
)

func FuzzGenerateSegnaturaXML(f *testing.F) {
	f.Add(1, 2026, "in", 1, "01", "Circolare n.1", "Ministero", "Scuola", "a1b2c3d4e5f6")
	f.Add(999999, 2025, "out", 4, "02", "<XML>Special & Tag</XML>", "Scuola", "Comune", "")
	f.Add(0, 0, "", 0, "", "", "", "", "hash")

	f.Fuzz(func(t *testing.T, protoNum, protoYear int, direction string, title int, class, subject, sender, recipient, hash string) {
		entry := ProtocolEntry{
			ProtocolNumber:      protoNum,
			ProtocolYear:        protoYear,
			ProtocolDate:        time.Now(),
			FlowDirection:       direction,
			ClassificationTitle: title,
			ClassificationClass: class,
			Subject:             subject,
			Sender:              sender,
			Recipient:           recipient,
			DocumentHashSHA256:  hash,
		}

		xmlBytes, err := GenerateSegnaturaXML(entry)
		if err != nil {
			t.Fatalf("failed to generate segnatura XML: %v", err)
		}

		// Verify XML is well-formed by decoding it
		var decoded XMLSegnatura
		decoder := xml.NewDecoder(bytes.NewReader(xmlBytes))
		if err := decoder.Decode(&decoded); err != nil {
			t.Fatalf("generated invalid XML: %v\nXML:\n%s", err, string(xmlBytes))
		}

		if decoded.Intestazione.Oggetto != subject {
			t.Fatalf("subject mismatch in XML: expected %q, got %q", subject, decoded.Intestazione.Oggetto)
		}
	})
}

func FuzzFormatVisualStamp(f *testing.F) {
	f.Add(123, "I.C. Leonardo da Vinci")
	f.Add(0, "")
	f.Add(9999999, "Liceo Scientifico G. Galilei - Roma")

	f.Fuzz(func(t *testing.T, protoNum int, schoolName string) {
		entry := ProtocolEntry{
			ProtocolNumber: protoNum,
			ProtocolDate:   time.Now(),
		}
		stamp := FormatVisualStamp(entry, schoolName)
		if !strings.Contains(stamp, "REG. UFF. PROT.") {
			t.Fatalf("stamp missing standard protocol tag: %s", stamp)
		}
	})
}
