package users

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/xuri/excelize/v2"
)

type Exporter struct{}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) ToCSV(data []map[string]interface{}) ([]byte, error) {
	if len(data) == 0 {
		return []byte{}, nil
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// Headers sorted deterministically
	var headers []string
	for k := range data[0] {
		headers = append(headers, k)
	}
	sort.Strings(headers)
	if err := w.Write(headers); err != nil {
		return nil, err
	}

	// Rows
	for _, row := range data {
		var record []string
		for _, h := range headers {
			val := fmt.Sprintf("%v", row[h])
			record = append(record, val)
		}
		if err := w.Write(record); err != nil {
			return nil, err
		}
	}

	w.Flush()
	return buf.Bytes(), nil
}

func (e *Exporter) ToXLSX(data []map[string]interface{}) ([]byte, error) {
	f := excelize.NewFile()
	sheet := "Sheet1"

	if len(data) == 0 {
		buf, _ := f.WriteToBuffer()
		return buf.Bytes(), nil
	}

	// Headers sorted deterministically
	var headers []string
	for k := range data[0] {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheet, cell, h)
	}

	// Rows
	for r, row := range data {
		for c, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(c+1, r+2)
			_ = f.SetCellValue(sheet, cell, row[h])
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e *Exporter) ToJSON(data []map[string]interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}
