package scrutiny

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-pdf/fpdf"
)

func GeneratePagellaPDF(matrix *ScrutinyMatrix, studentID string) ([]byte, error) {
	var targetStudent *StudentScrutinyRow
	for i := range matrix.Students {
		if matrix.Students[i].StudentID == studentID {
			targetStudent = &matrix.Students[i]
			break
		}
	}

	if targetStudent == nil {
		return nil, fmt.Errorf("student not found in scrutiny matrix")
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(0, 10, "DOCUMENTO DI VALUTAZIONE (PAGELLA)", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "I", 12)
	pdf.CellFormat(0, 8, fmt.Sprintf("Anno Scolastico - Semestre %d", matrix.Semester), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Student Info Box
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, fmt.Sprintf("Studente: %s", sanitize(targetStudent.StudentName)), "1", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, fmt.Sprintf("Assenze: %d | Ritardi: %d | Uscite Anticipate: %d",
		targetStudent.AttendanceStats.Absences,
		targetStudent.AttendanceStats.Lates,
		targetStudent.AttendanceStats.EarlyExits), "1", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Table Header
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(120, 8, "Materia", "1", 0, "L", false, 0, "")
	pdf.CellFormat(30, 8, "Media Voti", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 8, "Voto Finale", "1", 1, "C", false, 0, "")

	// Table Body
	pdf.SetFont("Arial", "", 10)
	var totalGrade float64
	var countSubjects int

	for _, sub := range matrix.Subjects {
		subData, exists := targetStudent.SubjectData[sub.ID]
		avgStr := "N/D"
		finalGradeStr := "N/D"

		if exists && subData.GradeCount > 0 {
			avgStr = fmt.Sprintf("%.2f", subData.Average)
			finalGradeStr = fmt.Sprintf("%.0f", subData.Proposed)
			totalGrade += subData.Proposed
			countSubjects++
		}

		pdf.CellFormat(120, 7, sanitize(sub.Name), "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 7, avgStr, "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 7, finalGradeStr, "1", 1, "C", false, 0, "")
	}

	// Conduct Grade
	if targetStudent.Record != nil && targetStudent.Record.ConductGrade > 0 {
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(120, 7, "Voto di Condotta", "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 7, "-", "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 7, fmt.Sprintf("%d", targetStudent.Record.ConductGrade), "1", 1, "C", false, 0, "")
	}

	pdf.Ln(5)

	// Final Decision
	pdf.SetFont("Arial", "B", 12)
	decision := "In Valutazione"
	if targetStudent.Record != nil && targetStudent.Record.FinalDecision != "" {
		decision = targetStudent.Record.FinalDecision
	}
	pdf.CellFormat(0, 8, fmt.Sprintf("Esito Finale / Decisione: %s", sanitize(decision)), "1", 1, "L", false, 0, "")

	if targetStudent.Record != nil && targetStudent.Record.Notes != "" {
		pdf.SetFont("Arial", "I", 10)
		pdf.MultiCell(0, 6, fmt.Sprintf("Note del Consiglio: %s", sanitize(targetStudent.Record.Notes)), "1", "L", false)
	}

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
	s = strings.ReplaceAll(s, "À", "A'")
	s = strings.ReplaceAll(s, "È", "E'")
	s = strings.ReplaceAll(s, "É", "E'")
	s = strings.ReplaceAll(s, "Ì", "I'")
	s = strings.ReplaceAll(s, "Ò", "O'")
	s = strings.ReplaceAll(s, "Ù", "U'")
	return s
}
