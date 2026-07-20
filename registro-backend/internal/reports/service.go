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
	f.SetCellValue(sheet, "A1", "Studente")
	colChar := 'B'
	for _, sub := range matrix.Subjects {
		cell := fmt.Sprintf("%c1", colChar)
		f.SetCellValue(sheet, cell, sub.Name)
		colChar++
	}
	f.SetCellValue(sheet, fmt.Sprintf("%c1", colChar), "Media Generale")
	colChar++
	f.SetCellValue(sheet, fmt.Sprintf("%c1", colChar), "Esito")

	// Rows
	rowNum := 2
	for _, stu := range matrix.Students {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), stu.StudentName)
		colChar = 'B'
		var sum float64
		var count int

		for _, sub := range matrix.Subjects {
			cell := fmt.Sprintf("%c%d", colChar, rowNum)
			if data, ok := stu.SubjectData[sub.ID]; ok && data.GradeCount > 0 {
				f.SetCellValue(sheet, cell, data.Proposed)
				sum += data.Proposed
				count++
			} else {
				f.SetCellValue(sheet, cell, "N/D")
			}
			colChar++
		}

		genAvg := 0.0
		if count > 0 {
			genAvg = sum / float64(count)
		}
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", colChar, rowNum), fmt.Sprintf("%.2f", genAvg))
		colChar++

		decision := "In Valutazione"
		if stu.Record != nil && stu.Record.FinalDecision != "" {
			decision = stu.Record.FinalDecision
		}
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", colChar, rowNum), decision)

		rowNum++
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
