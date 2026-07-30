package reports

import (
	"bytes"
	"context"
	"fmt"
	"registro-backend/internal/scrutiny"

	"github.com/xuri/excelize/v2"
)

type Service struct {
	scrutinySvc *scrutiny.Service
}

func NewService(scrutinySvc *scrutiny.Service) *Service {
	return &Service{
		scrutinySvc: scrutinySvc,
	}
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
	_ = f.SetCellValue(sheet, "A1", "Studente")
	colChar := 'B'
	for _, sub := range matrix.Subjects {
		cell := fmt.Sprintf("%c1", colChar)
		_ = f.SetCellValue(sheet, cell, sub.Name)
		colChar++
	}
	_ = f.SetCellValue(sheet, fmt.Sprintf("%c1", colChar), "Media Generale")
	colChar++
	_ = f.SetCellValue(sheet, fmt.Sprintf("%c1", colChar), "Esito")

	// Rows
	rowNum := 2
	for _, stu := range matrix.Students {
		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), stu.StudentName)
		colChar = 'B'
		var sum float64
		var count int

		for _, sub := range matrix.Subjects {
			cell := fmt.Sprintf("%c%d", colChar, rowNum)
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
			colChar++
		}

		genAvg := 0.0
		if count > 0 {
			genAvg = sum / float64(count)
		}
		_ = f.SetCellValue(sheet, fmt.Sprintf("%c%d", colChar, rowNum), fmt.Sprintf("%.2f", genAvg))
		colChar++

		decision := "Nessun Record"
		if stu.Record != nil {
			if stu.Record.FinalDecision != "" {
				decision = stu.Record.FinalDecision
			} else {
				decision = "In Corso"
			}
		}
		_ = f.SetCellValue(sheet, fmt.Sprintf("%c%d", colChar, rowNum), decision)

		rowNum++
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
