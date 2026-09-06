package reports

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"registro-backend/internal/scrutiny"

	"github.com/xuri/excelize/v2"
)

type Service struct {
	scrutinySvc *scrutiny.Service
	db          *sql.DB
}

func NewService(scrutinySvc *scrutiny.Service, db ...*sql.DB) *Service {
	s := &Service{
		scrutinySvc: scrutinySvc,
	}
	if len(db) > 0 {
		s.db = db[0]
	}
	return s
}

func (s *Service) SetDB(db *sql.DB) {
	s.db = db
}

// sanitizeExcelField prevents Excel formula injection by prepending a tab
// if the string begins with formula-triggering characters (=, +, -, @).
func sanitizeExcelField(v string) string {
	if len(v) > 0 {
		switch v[0] {
		case '=', '+', '-', '@':
			return "\t" + v
		}
	}
	return v
}

func (s *Service) ExportGradesExcel(ctx context.Context, actorID, actorRole, classID string, semester int) ([]byte, error) {
	matrix, err := s.scrutinySvc.GetMatrix(ctx, actorID, actorRole, classID, semester)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Matrice Voti"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)

	// Headers
	colIdx := 1
	cell, _ := excelize.CoordinatesToCellName(colIdx, 1)
	_ = f.SetCellValue(sheet, cell, "Studente")

	for _, sub := range matrix.Subjects {
		colIdx++
		cell, _ = excelize.CoordinatesToCellName(colIdx, 1)
		_ = f.SetCellValue(sheet, cell, sanitizeExcelField(sub.Name))
	}
	colIdx++
	cell, _ = excelize.CoordinatesToCellName(colIdx, 1)
	_ = f.SetCellValue(sheet, cell, "Media Generale")
	colIdx++
	cell, _ = excelize.CoordinatesToCellName(colIdx, 1)
	_ = f.SetCellValue(sheet, cell, "Esito")

	// Rows
	rowNum := 2
	for _, stu := range matrix.Students {
		cell, _ = excelize.CoordinatesToCellName(1, rowNum)
		_ = f.SetCellValue(sheet, cell, sanitizeExcelField(stu.StudentName))
		colIdx = 1
		var sum float64
		var count int

		for _, sub := range matrix.Subjects {
			colIdx++
			cell, _ = excelize.CoordinatesToCellName(colIdx, rowNum)
			val := 0.0
			found := false

			if stu.Record != nil {
				for _, g := range stu.Record.Grades {
					if g.SubjectID == sub.ID {
						val = g.FinalGrade
						found = true
						break
					}
				}
			}

			if !found {
				if data, ok := stu.SubjectData[sub.ID]; ok && data.GradeCount > 0 {
					val = data.Proposed
					found = true
				}
			}

			if found {
				_ = f.SetCellValue(sheet, cell, val)
				sum += val
				count++
			} else {
				_ = f.SetCellValue(sheet, cell, "N/D")
			}
		}

		genAvg := 0.0
		if count > 0 {
			genAvg = sum / float64(count)
		}
		colIdx++
		cell, _ = excelize.CoordinatesToCellName(colIdx, rowNum)
		_ = f.SetCellValue(sheet, cell, fmt.Sprintf("%.2f", genAvg))

		colIdx++
		cell, _ = excelize.CoordinatesToCellName(colIdx, rowNum)
		decision := "Nessun Record"
		if stu.Record != nil {
			if stu.Record.FinalDecision != "" {
				decision = stu.Record.FinalDecision
			} else {
				decision = "In Corso"
			}
		}
		_ = f.SetCellValue(sheet, cell, decision)

		rowNum++
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
