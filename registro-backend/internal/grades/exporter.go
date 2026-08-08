package grades

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"

	"github.com/go-pdf/fpdf"
)

type ReportCardPDFExporter interface {
	ExportReportCard(report *SemesterReportResponse, semester int) ([]byte, error)
}

type pdfExporter struct{}

func NewReportCardPDFExporter() ReportCardPDFExporter {
	return &pdfExporter{}
}

func (e *pdfExporter) ExportReportCard(report *SemesterReportResponse, semester int) ([]byte, error) {
	if report == nil {
		return nil, fmt.Errorf("report is nil")
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pdf.CellFormat(190, 10, "REGISTRO ELETTRONICO - SCHEDA DI VALUTAZIONE", "0", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	semText := fmt.Sprintf("%d° Quadrimestre", semester)
	pdf.CellFormat(190, 8, fmt.Sprintf("Valutazione Finale - %s - A.S. %s", semText, report.SchoolYear), "0", 1, "C", false, 0, "")
	pdf.Ln(4)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(95, 7, fmt.Sprintf("Studente: %s", report.StudentName), "1", 0, "L", false, 0, "")
	pdf.CellFormat(95, 7, fmt.Sprintf("Classe: %s", report.ClassName), "1", 1, "L", false, 0, "")
	pdf.CellFormat(95, 7, fmt.Sprintf("Media Generale: %.2f", report.OverallAverage), "1", 0, "L", false, 0, "")
	pdf.CellFormat(95, 7, fmt.Sprintf("Voto Comportamento: %.0f", report.BehaviorGrade), "1", 1, "L", false, 0, "")
	pdf.Ln(6)

	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(60, 8, "Materia", "1", 0, "L", true, 0, "")
	pdf.CellFormat(40, 8, "Docente", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 8, "Voti (N°)", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Media", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Voto Finale", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 10)
	for _, sub := range report.Subjects {
		subName := sub.Subject
		if len(subName) > 25 {
			subName = subName[:22] + "..."
		}
		teacherName := sub.Teacher
		if len(teacherName) > 18 {
			teacherName = teacherName[:15] + "..."
		}

		pdf.CellFormat(60, 7, subName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 7, teacherName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%d", sub.GradeCount), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%.2f", sub.SubjectAverage), "1", 0, "C", false, 0, "")

		finalGradeStr := fmt.Sprintf("%.0f", sub.FinalGrade)
		if sub.FinalGrade == 0 {
			finalGradeStr = fmt.Sprintf("%.1f", sub.SubjectAverage)
		}
		pdf.CellFormat(30, 7, finalGradeStr, "1", 1, "C", false, 0, "")
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type ExportOptions struct {
	Format string
}

func ExportToJSON(grades []Grade) ([]byte, error) {
	return json.Marshal(grades)
}

func ExportToCSV(grades []Grade, opts ExportOptions) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"ID", "StudentID", "SubjectID", "GradeValue", "GradeType", "Date", "Description"})
	for _, g := range grades {
		_ = w.Write([]string{g.ID, g.StudentID, g.SubjectID, fmt.Sprintf("%.2f", g.GradeValue), string(g.GradeType), g.Date.Format("2006-01-02"), g.Description})
	}
	w.Flush()
	return buf.Bytes(), nil
}

func ExportToPDF(grades []Grade, opts ExportOptions) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(190, 10, "EXPORT VOTI", "0", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	for _, g := range grades {
		pdf.CellFormat(190, 6, fmt.Sprintf("Data: %s | Studente: %s | Voto: %.2f | Tipo: %s", g.Date.Format("2006-01-02"), g.StudentID, g.GradeValue, string(g.GradeType)), "0", 1, "L", false, 0, "")
	}
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
