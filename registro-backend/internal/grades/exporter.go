package grades

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/go-pdf/fpdf"
)

type ExportOptions struct {
	Format     string // csv, json, pdf
	IncludeAll bool   // if true export even unpublished?
}

// 1. ExportToCSV
func ExportToCSV(grades []Grade, options ExportOptions) ([]byte, error) {
	// Sort
	sort.Slice(grades, func(i, j int) bool {
		if grades[i].StudentID != grades[j].StudentID {
			return grades[i].StudentID < grades[j].StudentID
		}
		return grades[i].Date.Before(grades[j].Date)
	})

	records := [][]string{
		{"StudentID", "SubjectID", "GradeValue", "GradeType", "Date", "Semester", "TeacherID", "Category", "Description"},
	}

	for _, g := range grades {
		records = append(records, []string{
			g.StudentID,
			g.SubjectID,
			fmt.Sprintf("%.2f", g.GradeValue),
			string(g.GradeType),
			g.Date.Format("2006-01-02"),
			strconv.Itoa(int(g.Semester)),
			g.TeacherID,
			string(g.GradeCategory),
			g.Description,
		})
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.WriteAll(records); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 2. ExportToJSON
type JSONExport struct {
	Metadata ExportMetadata `json:"metadata"`
	Grades   []Grade        `json:"grades"`
}

type ExportMetadata struct {
	ExportDate string `json:"exportDate"`
	Count      int    `json:"count"`
}

func ExportToJSON(grades []Grade) ([]byte, error) {
	export := JSONExport{
		Metadata: ExportMetadata{
			ExportDate: time.Now().Format(time.RFC3339),
			Count:      len(grades),
		},
		Grades: grades,
	}
	return json.MarshalIndent(export, "", "  ") // Indented for readability as requested
}

// 3. ExportToPDF
func ExportToPDF(grades []Grade, options ExportOptions) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(0, 0, 128) // Navy Blue
	pdf.Cell(40, 10, "Registro Elettronico - Report Voti")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(0, 10, fmt.Sprintf("Data Export: %s", time.Now().Format("2006-01-02 15:04")))
	pdf.Ln(10)

	// Table Header
	pdf.SetFillColor(240, 240, 240)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(30, 7, "Studente", "1", 0, "", true, 0, "")
	pdf.CellFormat(30, 7, "Materia", "1", 0, "", true, 0, "")
	pdf.CellFormat(25, 7, "Data", "1", 0, "", true, 0, "")
	pdf.CellFormat(15, 7, "Voto", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 7, "Categoria", "1", 0, "", true, 0, "")
	pdf.CellFormat(65, 7, "Descrizione", "1", 0, "", true, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 9)

	for _, g := range grades {
		// Color coding for grade
		if g.GradeValue < 6.0 {
			pdf.SetTextColor(200, 0, 0) // Red
		} else if g.GradeValue < 7.0 {
			pdf.SetTextColor(200, 165, 0) // Orange/Yellowish
		} else {
			pdf.SetTextColor(0, 128, 0) // Green
		}

		pdf.CellFormat(30, 6, g.StudentID, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 6, g.SubjectID, "1", 0, "", false, 0, "") // Truncate if needed
		pdf.CellFormat(25, 6, g.Date.Format("02/01/06"), "1", 0, "", false, 0, "")

		valStr := fmt.Sprintf("%.1f", g.GradeValue)
		if g.GradeType == GradeTypeJudgment {
			valStr = ConvertNumericToJudgment(g.GradeValue)
			if len(valStr) > 3 {
				valStr = valStr[:3] + "."
			} // Abbrev
		}
		pdf.CellFormat(15, 6, valStr, "1", 0, "C", false, 0, "")

		pdf.SetTextColor(0, 0, 0) // Reset for text
		pdf.CellFormat(25, 6, string(g.GradeCategory), "1", 0, "", false, 0, "")
		pdf.CellFormat(65, 6, g.Description, "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	return buf.Bytes(), err
}
