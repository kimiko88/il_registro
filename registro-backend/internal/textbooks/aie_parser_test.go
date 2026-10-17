package textbooks

import (
	"strings"
	"testing"
)

func TestParseAIECatalogCSV(t *testing.T) {
	csvData := `CODICE_ISBN;TITOLO;AUTORI;EDITORE;DISCIPLINA;PREZZO;VOLUME;ANNO_EDIZIONE;ORDINE_SCUOLA;DIGITALE
9788808836243;Matematica.blu 2.0;Bergamini M., Barozzi G.;Zanichelli;Matematica;28,90;1;2023;secondaria_2;false
9788804724567;Divina Commedia - Inferno;Dante Alighieri;Mondadori Scuola;Italiano;18.50;1;2021;secondaria_2;false
9788839521118;Fisica Modelli;Amaldi U.;Zanichelli;Fisica;24,00;1;2022;secondaria_2;true
`
	books, err := ParseAIECatalog(strings.NewReader(csvData))
	if err != nil {
		t.Fatalf("unexpected error parsing AIE CSV: %v", err)
	}

	if len(books) != 3 {
		t.Fatalf("expected 3 books, got %d", len(books))
	}

	// First book check (Italian comma decimal)
	b1 := books[0]
	if b1.ISBN != "9788808836243" {
		t.Errorf("expected ISBN 9788808836243, got %s", b1.ISBN)
	}
	if b1.Price != 28.90 {
		t.Errorf("expected price 28.90, got %.2f", b1.Price)
	}
	if b1.Authors != "Bergamini M., Barozzi G." {
		t.Errorf("expected authors 'Bergamini M., Barozzi G.', got %s", b1.Authors)
	}
	if b1.SchoolOrder != "secondaria_2" {
		t.Errorf("expected school order 'secondaria_2', got %s", b1.SchoolOrder)
	}
	if b1.IsDigitalOnly {
		t.Errorf("expected IsDigitalOnly false, got true")
	}

	// Third book check (digital only)
	b3 := books[2]
	if !b3.IsDigitalOnly {
		t.Errorf("expected book 3 to be digital only")
	}
}

func TestParseAIECatalog_EmptyOrInvalid(t *testing.T) {
	_, err := ParseAIECatalog(strings.NewReader(""))
	if err == nil {
		t.Errorf("expected error for empty catalog content")
	}

	invalidCSV := `ISBN;TITOLO
123`
	_, err = ParseAIECatalog(strings.NewReader(invalidCSV))
	if err == nil {
		t.Errorf("expected error for invalid columns")
	}
}
