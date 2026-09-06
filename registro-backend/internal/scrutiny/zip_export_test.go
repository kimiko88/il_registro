package scrutiny

import (
	"archive/zip"
	"bytes"
	"testing"
)

func TestGenerateClassPagelleZIP(t *testing.T) {
	// Test empty matrix
	_, err := GenerateClassPagelleZIP(&ScrutinyMatrix{})
	if err == nil {
		t.Fatal("expected error for empty matrix, got nil")
	}

	matrix := &ScrutinyMatrix{
		ClassID:  "class-1",
		Semester: 1,
		Subjects: []SubjectInfo{
			{ID: "sub-1", Name: "Matematica"},
			{ID: "sub-2", Name: "Italiano"},
		},
		Students: []StudentScrutinyRow{
			{
				StudentID:   "stud-1",
				StudentName: "Mario Rossi",
				SubjectData: map[string]SubjectAverages{
					"sub-1": {Proposed: 8, Average: 7.8, GradeCount: 5},
					"sub-2": {Proposed: 7, Average: 7.1, GradeCount: 4},
				},
				Record: &ScrutinyRecord{
					ConductGrade:  9,
					FinalDecision: "Ammesso",
					Notes:         "Ottimo rendimento",
				},
			},
			{
				StudentID:   "stud-2",
				StudentName: "Luigi Bianchi",
				SubjectData: map[string]SubjectAverages{
					"sub-1": {Proposed: 6, Average: 6.2, GradeCount: 3},
					"sub-2": {Proposed: 6, Average: 5.9, GradeCount: 4},
				},
				Record: &ScrutinyRecord{
					ConductGrade:  8,
					FinalDecision: "Ammesso",
				},
			},
		},
	}

	zipBytes, err := GenerateClassPagelleZIP(matrix)
	if err != nil {
		t.Fatalf("unexpected error generating class pagelle ZIP: %v", err)
	}

	if len(zipBytes) == 0 {
		t.Fatal("expected non-empty ZIP bytes")
	}

	// Read and verify the zip entries
	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		t.Fatalf("failed to read generated zip: %v", err)
	}

	if len(zipReader.File) != 2 {
		t.Fatalf("expected 2 files in zip, found %d", len(zipReader.File))
	}

	expectedFiles := map[string]bool{
		"Pagella_Mario_Rossi_Semestre_1.pdf":   false,
		"Pagella_Luigi_Bianchi_Semestre_1.pdf": false,
	}

	for _, f := range zipReader.File {
		if _, ok := expectedFiles[f.Name]; ok {
			expectedFiles[f.Name] = true
		} else {
			t.Errorf("unexpected file in ZIP: %s", f.Name)
		}
	}

	for name, found := range expectedFiles {
		if !found {
			t.Errorf("expected file %s was not found in ZIP", name)
		}
	}
}
