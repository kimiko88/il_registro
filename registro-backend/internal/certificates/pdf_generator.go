package certificates

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-pdf/fpdf"
)

func GenerateCertificatePDF(cert *Certificate, schoolName string) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.SetAutoPageBreak(true, 25)
	pdf.AddPage()

	if schoolName == "" {
		schoolName = "ISTITUTO STATALE D'ISTRUZIONE SUPERIORE"
	}

	// ---- Header ----
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, sanitize(schoolName), "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "I", 10)
	pdf.CellFormat(0, 6, "Repubblica Italiana — Ministero dell'Istruzione e del Merito", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("Anno Scolastico %s", sanitize(cert.AcademicYear)), "", 1, "C", false, 0, "")
	pdf.Ln(4)

	pdf.SetDrawColor(40, 60, 100)
	pdf.SetLineWidth(0.8)
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(8)

	// ---- Protocol & Title ----
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, fmt.Sprintf("Prot. n. %s", cert.ProtocolNo), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("Data: %s", cert.IssuedAt.Format("02/01/2006")), "", 1, "L", false, 0, "")
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 16)
	titleText := ""
	switch cert.Type {
	case CertIscrizione:
		titleText = "CERTIFICATO DI ISCRIZIONE"
	case CertFrequenza:
		titleText = "CERTIFICATO DI FREQUENZA"
	case CertPromozione:
		titleText = "CERTIFICATO DI PROMOZIONE"
	case CertBehavior:
		titleText = "CERTIFICATO DI CONDOTTA"
	default:
		titleText = "CERTIFICATO SCOLASTICO"
	}
	pdf.CellFormat(0, 10, titleText, "", 1, "C", false, 0, "")
	pdf.Ln(8)

	// ---- Body Text ----
	pdf.SetFont("Arial", "", 12)
	studentName := sanitize(cert.StudentName)
	if studentName == "" {
		studentName = "[Studente]"
	}
	className := sanitize(cert.ClassName)
	if className == "" {
		className = "[Classe]"
	}

	var body string
	switch cert.Type {
	case CertIscrizione:
		body = fmt.Sprintf("Si certifica che lo/la studente/ssa %s e' regolarmente iscritto/a per l'Anno Scolastico %s alla classe %s dell'Istituto.", studentName, cert.AcademicYear, className)
	case CertFrequenza:
		body = fmt.Sprintf("Si certifica che lo/la studente/ssa %s ha frequentato le lezioni per l'Anno Scolastico %s nella classe %s con regolare assiduita'.", studentName, cert.AcademicYear, className)
	case CertPromozione:
		body = fmt.Sprintf("Si certifica che lo/la studente/ssa %s, iscritto/a alla classe %s per l'Anno Scolastico %s, e' stato/a promosso/a alla classe successiva.", studentName, className, cert.AcademicYear)
	case CertBehavior:
		body = fmt.Sprintf("Si certifica che lo/la studente/ssa %s, iscritto/a alla classe %s nell'Anno Scolastico %s, ha mantenuto un comportamento ineccepibile ed un corretto adempimento dei doveri scolastici.", studentName, className, cert.AcademicYear)
	}

	pdf.MultiCell(0, 7, body, "", "J", false)
	pdf.Ln(6)

	if cert.Notes != "" {
		pdf.SetFont("Arial", "I", 10)
		pdf.MultiCell(0, 6, fmt.Sprintf("Note: %s", sanitize(cert.Notes)), "", "L", false)
		pdf.Ln(6)
	}

	pdf.SetFont("Arial", "", 11)
	pdf.MultiCell(0, 7, "Si rilascia il presente certificato per tutti gli usi consentiti dalla legge.", "", "J", false)
	pdf.Ln(25)

	// ---- Signature ----
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(100, 6, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(70, 6, "IL DIRIGENTE SCOLASTICO", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "I", 9)
	pdf.CellFormat(100, 6, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(70, 6, "(Firma autografa omessa ai sensi dell'art. 3 del D.Lgs. n. 39/1993)", "", 1, "C", false, 0, "")

	// ---- Footer ----
	pdf.SetY(260)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, 5, fmt.Sprintf("Documento emesso digitalmente il %s — RegistroV2 Certifications", cert.IssuedAt.Format("02/01/2006 15:04")), "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func sanitize(s string) string {
	s = strings.ReplaceAll(s, "à", "a'")
	s = strings.ReplaceAll(s, "è", "e'")
	s = strings.ReplaceAll(s, "é", "e'")
	s = strings.ReplaceAll(s, "ì", "i'")
	s = strings.ReplaceAll(s, "ò", "o'")
	s = strings.ReplaceAll(s, "ù", "u'")
	return s
}
