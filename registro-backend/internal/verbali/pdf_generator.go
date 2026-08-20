package verbali

import (
	"bytes"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/go-pdf/fpdf"
)

// GenerateVerbale builds a proper PDF document for a verbale using fpdf.
// It includes meeting info, verbale content, a signatures table, and a footer.
func GenerateVerbale(verbale *MeetingVerbale, meeting *CouncilMeeting, sigs []VerbaleSignature) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.SetAutoPageBreak(true, 20)
	pdf.AddPage()

	// ---- Header ----
	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(0, 10, "VERBALE DI SEDUTA", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "I", 11)
	pdf.CellFormat(0, 7, "Consiglio di Classe — Istituto Scolastico", "", 1, "C", false, 0, "")
	pdf.Ln(3)

	// Horizontal rule
	pdf.SetDrawColor(100, 100, 100)
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(5)

	// ---- Meeting Info ----
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(40, 7, "Titolo:", "0", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 7, sanitize(verbale.Title), "0", 1, "L", false, 0, "")

	if meeting != nil {
		pdf.SetFont("Arial", "B", 11)
		pdf.CellFormat(40, 7, "Data seduta:", "0", 0, "L", false, 0, "")
		pdf.SetFont("Arial", "", 11)
		pdf.CellFormat(0, 7, meeting.Date.Format("02/01/2006"), "0", 1, "L", false, 0, "")

		if meeting.StartTime != "" {
			pdf.SetFont("Arial", "B", 11)
			pdf.CellFormat(40, 7, "Orario:", "0", 0, "L", false, 0, "")
			pdf.SetFont("Arial", "", 11)
			pdf.CellFormat(0, 7, fmt.Sprintf("%s – %s", meeting.StartTime, meeting.EndTime), "0", 1, "L", false, 0, "")
		}

		if meeting.ClassID != "" {
			pdf.SetFont("Arial", "B", 11)
			pdf.CellFormat(40, 7, "Classe:", "0", 0, "L", false, 0, "")
			pdf.SetFont("Arial", "", 11)
			pdf.CellFormat(0, 7, meeting.ClassID, "0", 1, "L", false, 0, "")
		}

		if meeting.Agenda != "" {
			pdf.Ln(3)
			pdf.SetFont("Arial", "B", 11)
			pdf.CellFormat(0, 7, "Ordine del Giorno:", "0", 1, "L", false, 0, "")
			pdf.SetFont("Arial", "", 10)
			pdf.MultiCell(0, 6, sanitize(meeting.Agenda), "1", "L", false)
		}
	}

	pdf.Ln(5)

	// ---- Verbale Content ----
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, "CONTENUTO DEL VERBALE", "0", 1, "L", false, 0, "")
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(3)

	pdf.SetFont("Arial", "", 10)
	// Split content into paragraphs preserving line breaks
	for _, para := range strings.Split(sanitize(verbale.Content), "\n") {
		if strings.TrimSpace(para) == "" {
			pdf.Ln(4)
		} else {
			pdf.MultiCell(0, 5.5, para, "0", "L", false)
		}
	}

	pdf.Ln(8)

	// ---- Signatures Table ----
	if len(sigs) > 0 {
		pdf.SetFont("Arial", "B", 12)
		pdf.CellFormat(0, 8, "FIRME DIGITALI", "0", 1, "L", false, 0, "")
		pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
		pdf.Ln(3)

		// Table header
		pdf.SetFillColor(230, 230, 230)
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(70, 7, "Firmatario", "1", 0, "L", true, 0, "")
		pdf.CellFormat(55, 7, "Data e Ora", "1", 0, "C", true, 0, "")
		pdf.CellFormat(45, 7, "Indirizzo IP", "1", 1, "C", true, 0, "")

		pdf.SetFont("Arial", "", 9)
		pdf.SetFillColor(255, 255, 255)
		for i, sig := range sigs {
			fillColor := i%2 == 0
			if fillColor {
				pdf.SetFillColor(248, 248, 248)
			} else {
				pdf.SetFillColor(255, 255, 255)
			}
			name := sig.UserName
			if name == "" {
				if len(sig.UserID) > 8 {
					name = sig.UserID[:8] + "…"
				} else {
					name = sig.UserID
				}
			}
			pdf.CellFormat(70, 6, sanitize(name), "1", 0, "L", fillColor, 0, "")
			pdf.CellFormat(55, 6, sig.SignedAt.Format("02/01/2006 15:04"), "1", 0, "C", fillColor, 0, "")
			pdf.CellFormat(45, 6, sig.IPAddress, "1", 1, "C", fillColor, 0, "")
		}
	} else {
		pdf.SetFont("Arial", "I", 10)
		pdf.CellFormat(0, 7, "Nessuna firma digitale registrata.", "0", 1, "L", false, 0, "")
	}

	pdf.Ln(10)

	// ---- Publication Status ----
	status := "In bozza"
	if verbale.IsPublished {
		status = "Pubblicato"
	}
	pdf.SetFont("Arial", "I", 9)
	pdf.CellFormat(0, 6, fmt.Sprintf("Stato documento: %s  |  Generato il: %s", status, time.Now().Format("02/01/2006 15:04")), "0", 1, "C", false, 0, "")

	// ---- Output ----
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("GenerateVerbale pdf output: %w", err)
	}
	return buf.Bytes(), nil
}

// sanitize replaces non-latin1 characters (not supported by Arial in fpdf) with '?'.
func sanitize(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r > 255 || !unicode.IsPrint(r) {
			// Replace smart quotes and em dashes with ASCII equivalents
			switch r {
			case '\u2018', '\u2019':
				b.WriteRune('\'')
			case '\u201c', '\u201d':
				b.WriteRune('"')
			case '\u2013', '\u2014':
				b.WriteRune('-')
			case '\u2026':
				b.WriteString("...")
			default:
				b.WriteRune('?')
			}
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
