package documents

import (
	"bytes"
	"fmt"
)

type Exporter struct{}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) ToPDF(doc *Document, content string) ([]byte, error) {
	// Mock PDF generation utilizing gofpdf normally
	// Returning a buffer for now
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("PDF Export\nTitle: %s\n\n%s", doc.Title, content))
	return b.Bytes(), nil
}

func (e *Exporter) ToDOCX(doc *Document, content string) ([]byte, error) {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("DOCX Export\nTitle: %s\n\n%s", doc.Title, content))
	return b.Bytes(), nil
}
