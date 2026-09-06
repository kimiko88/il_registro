package classes

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

// MonthlyJournalData holds the complete information for the official monthly class journal PDF
type MonthlyJournalData struct {
	SchoolName        string
	ClassName         string
	CoordinatorName   string
	Month             int
	Year              int
	MonthName         string
	DaysInMonth       int
	Students          []MonthlyStudentAttendance
	Lessons           []MonthlyLesson
	DisciplinaryNotes []MonthlyDisciplinaryNote
	Stats             MonthlyJournalStats
}

// MonthlyStudentAttendance holds attendance statistics for a single student across the month
type MonthlyStudentAttendance struct {
	ID          string
	Name        string
	DailyStatus map[int]string // day number (1..31) -> "P", "A", "R", "U", "G"
	TotalP      int
	TotalA      int
	TotalR      int
	TotalU      int
}

// MonthlyLesson holds signed lesson details conducted in the class
type MonthlyLesson struct {
	Date        string
	Hour        int
	TeacherName string
	SubjectName string
	Topic       string
	Type        string
}

// MonthlyDisciplinaryNote holds note information recorded during the month
type MonthlyDisciplinaryNote struct {
	Date        string
	StudentName string
	TeacherName string
	Description string
	NoteType    string
}

// MonthlyJournalStats aggregates overall attendance metrics
type MonthlyJournalStats struct {
	TotalSchoolDays int
	TotalAbsences   int
	TotalLates      int
	TotalEarlyExits int
	AttendanceRate  float64
}

func sanitizeString(s string) string {
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

// GenerateMonthlyJournalPDF builds a vector PDF with ministerial styling for the monthly class register
func GenerateMonthlyJournalPDF(data *MonthlyJournalData) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	pdf := fpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(12, 10, 12)
	pdf.SetAutoPageBreak(true, 10)
	pdf.AddPage()

	// 1. Institutional Header
	pdf.SetFont("Arial", "B", 13)
	pdf.CellFormat(0, 6, "MINISTERO DELL'ISTRUZIONE E DEL MERITO", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 5, "GIORNALE DI CLASSE UFFICIALE DEL MESE - AGLI ATTI DELLA SCUOLA", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "I", 8.5)
	pdf.CellFormat(0, 4.5, fmt.Sprintf("Istituto: %s", sanitizeString(data.SchoolName)), "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// 2. Metadata Box
	pdf.SetFont("Arial", "B", 8.5)
	pdf.SetFillColor(245, 247, 250)
	pdf.CellFormat(70, 5.5, fmt.Sprintf(" Classe: %s", sanitizeString(data.ClassName)), "1", 0, "L", true, 0, "")
	pdf.CellFormat(95, 5.5, fmt.Sprintf(" Periodo: %s %d", sanitizeString(data.MonthName), data.Year), "1", 0, "L", true, 0, "")
	pdf.CellFormat(108, 5.5, fmt.Sprintf(" Docente Coordinatore: %s", sanitizeString(data.CoordinatorName)), "1", 1, "L", true, 0, "")
	pdf.Ln(3)

	// 3. Section 1: Daily Attendance Grid
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(0, 5, "1. PROSPETTO MENSILE DELLE PRESENZE, ASSENZE E GIUSTIFICAZIONI", "", 1, "L", false, 0, "")

	daysInMonth := data.DaysInMonth
	if daysInMonth <= 0 {
		daysInMonth = 31
	}

	// Dynamic day column width
	nameColWidth := 52.0
	dayColWidth := (273.0 - nameColWidth - 28.0) / float64(daysInMonth) // 28 mm for totals (P, A, R, U)
	totalColWidth := 7.0

	// Header Days Row
	pdf.SetFont("Arial", "B", 7)
	pdf.SetFillColor(230, 235, 245)
	pdf.CellFormat(nameColWidth, 5, "Alunno/a", "1", 0, "L", true, 0, "")

	for d := 1; d <= daysInMonth; d++ {
		// Detect sunday
		tDay := time.Date(data.Year, time.Month(data.Month), d, 0, 0, 0, 0, time.UTC)
		if tDay.Weekday() == time.Sunday {
			pdf.SetFillColor(215, 218, 224)
		} else {
			pdf.SetFillColor(230, 235, 245)
		}
		pdf.CellFormat(dayColWidth, 5, fmt.Sprintf("%d", d), "1", 0, "C", true, 0, "")
	}

	pdf.SetFillColor(230, 235, 245)
	pdf.CellFormat(totalColWidth, 5, "P", "1", 0, "C", true, 0, "")
	pdf.CellFormat(totalColWidth, 5, "A", "1", 0, "C", true, 0, "")
	pdf.CellFormat(totalColWidth, 5, "R", "1", 0, "C", true, 0, "")
	pdf.CellFormat(totalColWidth, 5, "U", "1", 1, "C", true, 0, "")

	// Student Rows
	pdf.SetFont("Arial", "", 6.5)
	for idx, st := range data.Students {
		fill := (idx % 2) == 1
		if fill {
			pdf.SetFillColor(250, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		name := sanitizeString(st.Name)
		if len(name) > 28 {
			name = name[:25] + "..."
		}
		pdf.CellFormat(nameColWidth, 4.5, fmt.Sprintf(" %d. %s", idx+1, name), "1", 0, "L", fill, 0, "")

		for d := 1; d <= daysInMonth; d++ {
			tDay := time.Date(data.Year, time.Month(data.Month), d, 0, 0, 0, 0, time.UTC)
			isSunday := tDay.Weekday() == time.Sunday

			status := st.DailyStatus[d]
			if status == "" {
				if isSunday {
					status = "D"
				} else {
					status = "."
				}
			}

			if isSunday {
				pdf.SetFillColor(235, 238, 242)
				pdf.CellFormat(dayColWidth, 4.5, "", "1", 0, "C", true, 0, "")
			} else {
				pdf.CellFormat(dayColWidth, 4.5, status, "1", 0, "C", fill, 0, "")
			}
		}

		pdf.SetFont("Arial", "B", 6.5)
		pdf.CellFormat(totalColWidth, 4.5, fmt.Sprintf("%d", st.TotalP), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(totalColWidth, 4.5, fmt.Sprintf("%d", st.TotalA), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(totalColWidth, 4.5, fmt.Sprintf("%d", st.TotalR), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(totalColWidth, 4.5, fmt.Sprintf("%d", st.TotalU), "1", 1, "C", fill, 0, "")
		pdf.SetFont("Arial", "", 6.5)
	}

	if len(data.Students) == 0 {
		pdf.CellFormat(273, 6, "Nessun dato di presenza registrato per il mese selezionato", "1", 1, "C", false, 0, "")
	}

	// 4. Monthly Stats Summary Bar
	pdf.Ln(3)
	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(240, 244, 250)
	statsText := fmt.Sprintf(" Riepilogo Mese: Giorni di lezione: %d | Totale Assenze: %d | Totale Ritardi: %d | Uscite Anticipate: %d | Tasso Presenza Medio: %.1f%%",
		data.Stats.TotalSchoolDays, data.Stats.TotalAbsences, data.Stats.TotalLates, data.Stats.TotalEarlyExits, data.Stats.AttendanceRate)
	pdf.CellFormat(273, 5.5, sanitizeString(statsText), "1", 1, "L", true, 0, "")

	// 5. Lessons Page
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(0, 5, "2. REGISTRO DELLE LEZIONI SVOLTE E ARGOMENTI FIRMATI NEL MESE", "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 7.5)
	pdf.SetFillColor(230, 235, 245)
	pdf.CellFormat(22, 5, "Data", "1", 0, "C", true, 0, "")
	pdf.CellFormat(12, 5, "Ora", "1", 0, "C", true, 0, "")
	pdf.CellFormat(38, 5, "Materia", "1", 0, "L", true, 0, "")
	pdf.CellFormat(42, 5, "Docente", "1", 0, "L", true, 0, "")
	pdf.CellFormat(127, 5, "Argomento Trattato", "1", 0, "L", true, 0, "")
	pdf.CellFormat(32, 5, "Tipologia", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 7.5)
	for i, l := range data.Lessons {
		fill := (i % 2) == 1
		if fill {
			pdf.SetFillColor(250, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		topic := sanitizeString(l.Topic)
		if len(topic) > 80 {
			topic = topic[:77] + "..."
		}

		pdf.CellFormat(22, 5, l.Date, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(12, 5, fmt.Sprintf("%dª", l.Hour), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(38, 5, sanitizeString(l.SubjectName), "1", 0, "L", fill, 0, "")
		pdf.CellFormat(42, 5, sanitizeString(l.TeacherName), "1", 0, "L", fill, 0, "")
		pdf.CellFormat(127, 5, topic, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(32, 5, sanitizeString(l.Type), "1", 1, "C", fill, 0, "")
	}

	if len(data.Lessons) == 0 {
		pdf.CellFormat(273, 6, "Nessuna lezione firmata nel mese selezionato", "1", 1, "C", false, 0, "")
	}

	// 6. Disciplinary Notes (if any)
	if len(data.DisciplinaryNotes) > 0 {
		pdf.Ln(4)
		pdf.SetFont("Arial", "B", 9)
		pdf.CellFormat(0, 5, "3. NOTE E PROVVEDIMENTI DISCIPLINARI DEL MESE", "", 1, "L", false, 0, "")

		pdf.SetFont("Arial", "B", 7.5)
		pdf.SetFillColor(230, 235, 245)
		pdf.CellFormat(22, 5, "Data", "1", 0, "C", true, 0, "")
		pdf.CellFormat(45, 5, "Alunno/a", "1", 0, "L", true, 0, "")
		pdf.CellFormat(40, 5, "Docente", "1", 0, "L", true, 0, "")
		pdf.CellFormat(30, 5, "Tipologia", "1", 0, "L", true, 0, "")
		pdf.CellFormat(136, 5, "Descrizione della Nota", "1", 1, "L", true, 0, "")

		pdf.SetFont("Arial", "", 7.5)
		for i, n := range data.DisciplinaryNotes {
			fill := (i % 2) == 1
			if fill {
				pdf.SetFillColor(250, 250, 252)
			} else {
				pdf.SetFillColor(255, 255, 255)
			}

			desc := sanitizeString(n.Description)
			if len(desc) > 85 {
				desc = desc[:82] + "..."
			}

			pdf.CellFormat(22, 5, n.Date, "1", 0, "C", fill, 0, "")
			pdf.CellFormat(45, 5, sanitizeString(n.StudentName), "1", 0, "L", fill, 0, "")
			pdf.CellFormat(40, 5, sanitizeString(n.TeacherName), "1", 0, "L", fill, 0, "")
			pdf.CellFormat(30, 5, sanitizeString(n.NoteType), "1", 0, "L", fill, 0, "")
			pdf.CellFormat(136, 5, desc, "1", 1, "L", fill, 0, "")
		}
	}

	// 7. Formal Closing & Signatures
	pdf.Ln(6)
	pdf.SetFont("Arial", "I", 7.5)
	pdf.MultiCell(0, 4, sanitizeString("Il presente giornale di classe mensile viene formalmente redatto, convalidato e sigillato agli atti della scuola conformemente al R.D. 965/1924 e alle disposizioni vigenti in materia di dematerializzazione documentale (CAD)."), "", "L", false)
	pdf.Ln(4)

	y := pdf.GetY()
	if y > 165 {
		pdf.AddPage()
	}

	pdf.SetFont("Arial", "B", 8.5)
	pdf.CellFormat(135, 5, "Il Docente Coordinatore di Classe", "0", 0, "C", false, 0, "")
	pdf.CellFormat(138, 5, "Visto: Il Dirigente Scolastico", "0", 1, "C", false, 0, "")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 8.5)
	pdf.CellFormat(135, 5, "________________________________________", "0", 0, "C", false, 0, "")
	pdf.CellFormat(138, 5, "________________________________________", "0", 1, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to output monthly journal pdf: %w", err)
	}

	return buf.Bytes(), nil
}
