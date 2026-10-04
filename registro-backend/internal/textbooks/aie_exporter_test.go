package textbooks

import (
	"strings"
	"testing"
)

func TestExportAIETXT(t *testing.T) {
	adoptions := []AIEExportRecord{
		{
			SchoolCode:   "RMPS010004",
			AcademicYear: "2026/2027",
			ClassName:    "1A",
			SubjectCode:  "MAT",
			ISBN:         "9788808836243",
			Title:        "Matematica.blu 2.0",
			Authors:      "Bergamini M.",
			Publisher:    "Zanichelli",
			Price:        28.90,
			Volume:       "1",
			AdoptionType: "nuova_adozione",
			AlreadyOwned: false,
		},
	}

	txtOutput, err := ExportAIEFormat(adoptions)
	if err != nil {
		t.Fatalf("unexpected export error: %v", err)
	}

	if !strings.Contains(txtOutput, "RMPS010004") {
		t.Errorf("export missing school code")
	}
	if !strings.Contains(txtOutput, "9788808836243") {
		t.Errorf("export missing ISBN")
	}
	if !strings.Contains(txtOutput, "Matematica.blu 2.0") {
		t.Errorf("export missing title")
	}
}
