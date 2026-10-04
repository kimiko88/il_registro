package textbooks

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ParseAIECatalog parses an AIE catalog in CSV or semicolon-delimited format.
func ParseAIECatalog(r io.Reader) ([]AIECatalogBook, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading catalog: %w", err)
	}
	if len(lines) < 2 {
		return nil, errors.New("catalog is empty or missing headers")
	}

	// Detect delimiter (; or , or \t)
	delimiter := ';'
	if strings.Contains(lines[0], ";") {
		delimiter = ';'
	} else if strings.Contains(lines[0], "\t") {
		delimiter = '\t'
	} else if strings.Contains(lines[0], ",") {
		delimiter = ','
	}

	reader := csv.NewReader(strings.NewReader(strings.Join(lines, "\n")))
	reader.Comma = delimiter
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parsing CSV: %w", err)
	}
	if len(records) < 2 {
		return nil, errors.New("insufficient records in catalog")
	}

	header := records[0]
	colMap := make(map[string]int)
	for i, h := range header {
		clean := strings.ToUpper(strings.TrimSpace(h))
		clean = strings.ReplaceAll(clean, " ", "_")
		colMap[clean] = i
	}

	getColIdx := func(keys ...string) int {
		for _, k := range keys {
			if idx, ok := colMap[k]; ok {
				return idx
			}
		}
		return -1
	}

	isbnIdx := getColIdx("CODICE_ISBN", "ISBN")
	titleIdx := getColIdx("TITOLO", "TITLE")
	authorIdx := getColIdx("AUTORI", "AUTORE", "AUTHORS", "AUTHOR")
	publisherIdx := getColIdx("EDITORE", "PUBLISHER")
	subjectIdx := getColIdx("DISCIPLINA", "MATERIA", "SUBJECT")
	priceIdx := getColIdx("PREZZO", "PRICE")
	volumeIdx := getColIdx("VOLUME", "VOL")
	yearIdx := getColIdx("ANNO_EDIZIONE", "ANNO", "EDITION_YEAR")
	orderIdx := getColIdx("ORDINE_SCUOLA", "ORDINE", "SCHOOL_ORDER")
	digitalIdx := getColIdx("DIGITALE", "IS_DIGITAL", "IS_DIGITAL_ONLY")

	if isbnIdx == -1 || titleIdx == -1 {
		return nil, errors.New("missing mandatory columns ISBN or TITOLO")
	}

	var books []AIECatalogBook
	for rowIdx, row := range records[1:] {
		if len(row) <= isbnIdx || len(row) <= titleIdx {
			continue
		}
		isbn := strings.TrimSpace(row[isbnIdx])
		if isbn == "" {
			continue
		}
		title := strings.TrimSpace(row[titleIdx])

		var author, publisher, subject, volume string
		if authorIdx != -1 && authorIdx < len(row) {
			author = strings.TrimSpace(row[authorIdx])
		}
		if publisherIdx != -1 && publisherIdx < len(row) {
			publisher = strings.TrimSpace(row[publisherIdx])
		}
		if subjectIdx != -1 && subjectIdx < len(row) {
			subject = strings.TrimSpace(row[subjectIdx])
		}
		if volumeIdx != -1 && volumeIdx < len(row) {
			volume = strings.TrimSpace(row[volumeIdx])
		}
		if volume == "" {
			volume = "1"
		}

		price := 0.0
		if priceIdx != -1 && priceIdx < len(row) {
			priceStr := strings.TrimSpace(row[priceIdx])
			priceStr = strings.ReplaceAll(priceStr, "€", "")
			priceStr = strings.ReplaceAll(priceStr, " ", "")
			priceStr = strings.ReplaceAll(priceStr, ",", ".")
			if p, err := strconv.ParseFloat(priceStr, 64); err == nil {
				price = p
			}
		}

		year := 0
		if yearIdx != -1 && yearIdx < len(row) {
			if y, err := strconv.Atoi(strings.TrimSpace(row[yearIdx])); err == nil {
				year = y
			}
		}

		schoolOrder := "secondaria_2"
		if orderIdx != -1 && orderIdx < len(row) {
			val := strings.TrimSpace(row[orderIdx])
			if val != "" {
				schoolOrder = val
			}
		}

		isDigital := false
		if digitalIdx != -1 && digitalIdx < len(row) {
			val := strings.ToLower(strings.TrimSpace(row[digitalIdx]))
			if val == "true" || val == "1" || val == "si" || val == "s" {
				isDigital = true
			}
		}

		books = append(books, AIECatalogBook{
			ISBN:          isbn,
			Title:         title,
			Authors:       author,
			Publisher:     publisher,
			Subject:       subject,
			Price:         price,
			Volume:        volume,
			EditionYear:   year,
			SchoolOrder:   schoolOrder,
			IsDigitalOnly: isDigital,
		})
		_ = rowIdx
	}

	return books, nil
}
