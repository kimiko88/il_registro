package teachers

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-pdf/fpdf"
)

// TeacherRegisterData holds the complete payload required for the official teacher register PDF
type TeacherRegisterData struct {
	TeacherName  string
	SchoolName   string
	ClassName    string
	SubjectName  string
	AcademicYear string
	Students     []TeacherRegisterStudent
	Lessons      []TeacherRegisterLesson
}

// TeacherRegisterStudent holds student evaluation and attendance summary
type TeacherRegisterStudent struct {
	ID          string
	Name        string
	Q1Written   string
	Q1Oral      string
	Q1Practical string
	Q1Avg       string
	Q2Written   string
	Q2Oral      string
	Q2Practical string
	Q2Avg       string
	FinalAvg    string
	Absences    int
}

// TeacherRegisterLesson holds signed lesson metadata
type TeacherRegisterLesson struct {
	Date     string
	Hour     int
	Topic    string
	Type     string
	SignedBy string
}

func sanitize(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if r < 128 {
			sb.WriteRune(r)
		} else {
			switch r {
			case 'à', 'á':
				sb.WriteString("a'")
			case 'è', 'é':
				sb.WriteString("e'")
			case 'ì', 'í':
				sb.WriteString("i'")
			case 'ò', 'ó':
				sb.WriteString("o'")
			case 'ù', 'ú':
				sb.WriteString("u'")
			case '°':
				sb.WriteString("o")
			case '€':
				sb.WriteString("EUR")
			default:
				sb.WriteRune('?')
			}
		}
	}
	return sb.String()
}

// GenerateTeacherRegisterPDF produces a ministerial-formatted PDF for the teacher's personal register
func GenerateTeacherRegisterPDF(data *TeacherRegisterData) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	pdf := fpdf.New("L", "mm", "A4", "") // Landscape for wide ministerial grade grids
	pdf.SetMargins(15, 12, 15)
	pdf.SetAutoPageBreak(true, 12)
	pdf.AddPage()

	// 1. Header & Institutional Heading
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 7, "MINISTERO DELL'ISTRUZIONE E DEL MERITO", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 6, "REGISTRO PERSONALE DEL DOCENTE - CHIUSURA ANNUALE AGLI ATTI", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "I", 9)
	pdf.CellFormat(0, 5, fmt.Sprintf("Istituto: %s | Anno Scolastico: %s", sanitize(data.SchoolName), sanitize(data.AcademicYear)), "", 1, "C", false, 0, "")
	pdf.Ln(3)

	// 2. Teacher & Subject Details Box
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(245, 247, 250)
	pdf.CellFormat(90, 6, fmt.Sprintf(" Docente: %s", sanitize(data.TeacherName)), "1", 0, "L", true, 0, "")
	pdf.CellFormat(90, 6, fmt.Sprintf(" Materia: %s", sanitize(data.SubjectName)), "1", 0, "L", true, 0, "")
	pdf.CellFormat(87, 6, fmt.Sprintf(" Classe: %s", sanitize(data.ClassName)), "1", 1, "L", true, 0, "")
	pdf.Ln(4)

	// 3. Section Title: Griglia Valutazioni Periodiche
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 6, "1. PROSPETTO GENERALE DEI VOTI E MEDIE DISCIPLINARI", "", 1, "L", false, 0, "")

	// Table Header
	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(230, 235, 245)
	pdf.CellFormat(55, 6, "Alunno/a", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 6, "1Q Scritti", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 6, "1Q Orali", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 6, "1Q Prat.", "1", 0, "C", true, 0, "")
	pdf.CellFormat(22, 6, "1Q Media", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 6, "2Q Scritti", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 6, "2Q Orali", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 6, "2Q Prat.", "1", 0, "C", true, 0, "")
	pdf.CellFormat(22, 6, "2Q Media", "1", 0, "C", true, 0, "")
	pdf.CellFormat(24, 6, "Media Finale", "1", 0, "C", true, 0, "")
	pdf.CellFormat(24, 6, "Assenze (ore)", "1", 1, "C", true, 0, "")

	// Table Rows
	pdf.SetFont("Arial", "", 8)
	for i, st := range data.Students {
		fill := (i % 2) == 1
		if fill {
			pdf.SetFillColor(250, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		name := sanitize(st.Name)
		if len(name) > 30 {
			name = name[:27] + "..."
		}
		pdf.CellFormat(55, 5.5, fmt.Sprintf(" %d. %s", i+1, name), "1", 0, "L", fill, 0, "")
		pdf.CellFormat(20, 5.5, st.Q1Written, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(20, 5.5, st.Q1Oral, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(20, 5.5, st.Q1Practical, "1", 0, "C", fill, 0, "")

		pdf.SetFont("Arial", "B", 8)
		pdf.CellFormat(22, 5.5, st.Q1Avg, "1", 0, "C", fill, 0, "")
		pdf.SetFont("Arial", "", 8)

		pdf.CellFormat(20, 5.5, st.Q2Written, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(20, 5.5, st.Q2Oral, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(20, 5.5, st.Q2Practical, "1", 0, "C", fill, 0, "")

		pdf.SetFont("Arial", "B", 8)
		pdf.CellFormat(22, 5.5, st.Q2Avg, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(24, 5.5, st.FinalAvg, "1", 0, "C", fill, 0, "")
		pdf.SetFont("Arial", "", 8)

		pdf.CellFormat(24, 5.5, fmt.Sprintf("%d", st.Absences), "1", 1, "C", fill, 0, "")
	}

	if len(data.Students) == 0 {
		pdf.CellFormat(267, 7, "Nessun dato registrato per la classe e materia selezionata", "1", 1, "C", false, 0, "")
	}

	pdf.Ln(4)

	// 4. Lessons Page / Section
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 6, "2. REGISTRO DELLE LEZIONI FIRMATE E DEGLI ARGOMENTI TRATTATI", "", 1, "L", false, 0, "")

	// Lessons Table Header
	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(230, 235, 245)
	pdf.CellFormat(25, 6, "Data", "1", 0, "C", true, 0, "")
	pdf.CellFormat(15, 6, "Ora", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 6, "Tipologia", "1", 0, "L", true, 0, "")
	pdf.CellFormat(137, 6, "Argomento e Attività Didattica Svolta", "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 6, "Firma Elettronica / Timestamp", "1", 1, "C", true, 0, "")

	// Lessons Rows
	pdf.SetFont("Arial", "", 8)
	for i, l := range data.Lessons {
		fill := (i % 2) == 1
		if fill {
			pdf.SetFillColor(250, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		topic := sanitize(l.Topic)
		if len(topic) > 85 {
			topic = topic[:82] + "..."
		}

		sig := sanitize(l.SignedBy)
		if sig == "" {
			sig = sanitize(data.TeacherName) + " (FEQ)"
		}

		pdf.CellFormat(25, 5.5, l.Date, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(15, 5.5, fmt.Sprintf("%dª", l.Hour), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(30, 5.5, sanitize(l.Type), "1", 0, "L", fill, 0, "")
		pdf.CellFormat(137, 5.5, topic, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(60, 5.5, sig, "1", 1, "C", fill, 0, "")
	}

	if len(data.Lessons) == 0 {
		pdf.CellFormat(267, 7, "Nessuna lezione firmata a registro per la materia e classe specificata", "1", 1, "C", false, 0, "")
	}

	pdf.Ln(6)

	// 5. Official Closing & Declaration Box
	pdf.SetFont("Arial", "I", 8)
	pdf.MultiCell(0, 4, sanitize("Il sottoscritto docente dichiara sotto la propria personale responsabilita', ai sensi del D.P.R. 445/2000, che i dati, i voti, le assenze e le lezioni sopra riportati corrispondono fedelmente alle annotazioni apposte durante l'anno scolastico nel registro elettronico."), "", "L", false)
	pdf.Ln(4)

	// Signature Blocks
	yBeforeSign := pdf.GetY()
	if yBeforeSign > 165 {
		pdf.AddPage()
	}

	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(130, 6, "Firma del Docente", "0", 0, "C", false, 0, "")
	pdf.CellFormat(137, 6, "Visto: Il Dirigente Scolastico", "0", 1, "C", false, 0, "")
	pdf.Ln(12)
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(130, 5, "________________________________________", "0", 0, "C", false, 0, "")
	pdf.CellFormat(137, 5, "________________________________________", "0", 1, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate teacher register pdf: %w", err)
	}

	return buf.Bytes(), nil
}
