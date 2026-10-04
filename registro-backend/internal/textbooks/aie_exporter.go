package textbooks

import (
	"bytes"
	"fmt"
	"strings"
)

// ExportAIEFormat outputs adoptions in standard AIE exchange format (.txt/.csv)
func ExportAIEFormat(records []AIEExportRecord) (string, error) {
	var buf bytes.Buffer

	// Official AIE header
	header := []string{
		"CODICE_SCUOLA",
		"ANNO_SCOLASTICO",
		"CLASSE",
		"CODICE_MATERIA",
		"ISBN",
		"TITOLO",
		"AUTORI",
		"EDITORE",
		"PREZZO",
		"VOLUME",
		"TIPO_ADOZIONE",
		"GIA_IN_POSSESSO",
	}
	buf.WriteString(strings.Join(header, ";") + "\r\n")

	for _, r := range records {
		alreadyOwnedStr := "NO"
		if r.AlreadyOwned {
			alreadyOwnedStr = "SI"
		}

		line := []string{
			r.SchoolCode,
			r.AcademicYear,
			r.ClassName,
			r.SubjectCode,
			r.ISBN,
			r.Title,
			r.Authors,
			r.Publisher,
			fmt.Sprintf("%.2f", r.Price),
			r.Volume,
			r.AdoptionType,
			alreadyOwnedStr,
		}
		buf.WriteString(strings.Join(line, ";") + "\r\n")
	}

	return buf.String(), nil
}
