package enrollment

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func ParseSIDIExport(r io.Reader) ([]EnrollmentApplication, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading SIDI stream: %w", err)
	}
	if len(lines) < 2 {
		return nil, errors.New("SIDI file is empty or missing headers")
	}

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
		return nil, fmt.Errorf("parsing SIDI CSV: %w", err)
	}
	if len(records) < 2 {
		return nil, errors.New("insufficient records in SIDI file")
	}

	header := records[0]
	colMap := make(map[string]int)
	for i, h := range header {
		clean := strings.ToUpper(strings.TrimSpace(h))
		clean = strings.ReplaceAll(clean, " ", "_")
		colMap[clean] = i
	}

	getIdx := func(keys ...string) int {
		for _, k := range keys {
			if idx, ok := colMap[k]; ok {
				return idx
			}
		}
		return -1
	}

	idIdx := getIdx("ID_DOMANDA", "ID", "SIDI_APPLICATION_ID")
	nameIdx := getIdx("NOME", "FIRST_NAME", "STUDENT_FIRST_NAME")
	surnameIdx := getIdx("COGNOME", "LAST_NAME", "STUDENT_LAST_NAME")
	taxCodeIdx := getIdx("CODICE_FISCALE", "CF", "TAX_CODE")
	birthIdx := getIdx("DATA_NASCITA", "BIRTH_DATE")
	genderIdx := getIdx("SESSO", "GENDER")
	gradeIdx := getIdx("VOTO_LICENZA", "VOTO_ESAME", "GRADE")
	trackIdx := getIdx("INDIRIZZO", "INDIRIZZO_STUDI", "TRACK")
	langIdx := getIdx("SECONDA_LINGUA", "LINGUA_2", "SECOND_LANGUAGE")
	l104Idx := getIdx("L104", "DISABILITA", "HAS_L104")
	dsaIdx := getIdx("DSA", "BES_DSA", "HAS_DSA")
	relIdx := getIdx("RELIGIONE", "SCELTA_RELIGIONE", "RELIGION")
	reqClassmatesIdx := getIdx("COMPAGNI_RICHIESTI", "DESIDERATA", "REQUESTED_CLASSMATES")
	p1NameIdx := getIdx("GENITORE1_NOME", "GENITORE_NOME", "PARENT1_FIRST_NAME")
	p1SurnameIdx := getIdx("GENITORE1_COGNOME", "GENITORE_COGNOME", "PARENT1_LAST_NAME")
	p1EmailIdx := getIdx("GENITORE1_EMAIL", "EMAIL_GENITORE", "PARENT1_EMAIL")
	p1PhoneIdx := getIdx("GENITORE1_TEL", "TELEFONO_GENITORE", "PARENT1_PHONE")

	if nameIdx == -1 || surnameIdx == -1 || taxCodeIdx == -1 {
		return nil, errors.New("missing required SIDI columns (NOME, COGNOME, CODICE_FISCALE)")
	}

	var results []EnrollmentApplication
	for _, row := range records[1:] {
		if len(row) <= taxCodeIdx || strings.TrimSpace(row[taxCodeIdx]) == "" {
			continue
		}

		sidiID := ""
		if idIdx != -1 && idIdx < len(row) {
			sidiID = strings.TrimSpace(row[idIdx])
		}
		if sidiID == "" {
			sidiID = fmt.Sprintf("SIDI-%d", time.Now().UnixNano())
		}

		firstName := strings.TrimSpace(row[nameIdx])
		lastName := strings.TrimSpace(row[surnameIdx])
		taxCode := strings.ToUpper(strings.TrimSpace(row[taxCodeIdx]))

		birthDate := "2012-01-01"
		if birthIdx != -1 && birthIdx < len(row) && strings.TrimSpace(row[birthIdx]) != "" {
			birthDate = strings.TrimSpace(row[birthIdx])
		}

		gender := "M"
		if genderIdx != -1 && genderIdx < len(row) {
			val := strings.ToUpper(strings.TrimSpace(row[genderIdx]))
			if val == "F" || strings.HasPrefix(val, "FEMM") {
				gender = "F"
			}
		}

		grade := 7
		if gradeIdx != -1 && gradeIdx < len(row) {
			if g, err := strconv.Atoi(strings.TrimSpace(row[gradeIdx])); err == nil && g >= 6 && g <= 10 {
				grade = g
			}
		}

		track := "Ordinario"
		if trackIdx != -1 && trackIdx < len(row) && strings.TrimSpace(row[trackIdx]) != "" {
			track = strings.TrimSpace(row[trackIdx])
		}

		lang := "Spagnolo"
		if langIdx != -1 && langIdx < len(row) && strings.TrimSpace(row[langIdx]) != "" {
			lang = strings.TrimSpace(row[langIdx])
		}

		l104 := false
		if l104Idx != -1 && l104Idx < len(row) {
			val := strings.ToLower(strings.TrimSpace(row[l104Idx]))
			if val == "true" || val == "1" || val == "si" || val == "s" {
				l104 = true
			}
		}

		dsa := false
		if dsaIdx != -1 && dsaIdx < len(row) {
			val := strings.ToLower(strings.TrimSpace(row[dsaIdx]))
			if val == "true" || val == "1" || val == "si" || val == "s" {
				dsa = true
			}
		}

		religion := "irc"
		if relIdx != -1 && relIdx < len(row) && strings.TrimSpace(row[relIdx]) != "" {
			religion = strings.TrimSpace(row[relIdx])
		}

		var reqClassmates []string
		if reqClassmatesIdx != -1 && reqClassmatesIdx < len(row) {
			raw := strings.TrimSpace(row[reqClassmatesIdx])
			if raw != "" {
				parts := strings.Split(raw, ",")
				for _, p := range parts {
					clean := strings.TrimSpace(p)
					if clean != "" {
						reqClassmates = append(reqClassmates, clean)
					}
				}
			}
		}

		p1First := "Genitore"
		if p1NameIdx != -1 && p1NameIdx < len(row) && strings.TrimSpace(row[p1NameIdx]) != "" {
			p1First = strings.TrimSpace(row[p1NameIdx])
		}

		p1Last := lastName
		if p1SurnameIdx != -1 && p1SurnameIdx < len(row) && strings.TrimSpace(row[p1SurnameIdx]) != "" {
			p1Last = strings.TrimSpace(row[p1SurnameIdx])
		}

		p1Email := "famiglia@example.com"
		if p1EmailIdx != -1 && p1EmailIdx < len(row) && strings.TrimSpace(row[p1EmailIdx]) != "" {
			p1Email = strings.TrimSpace(row[p1EmailIdx])
		}

		p1Phone := ""
		if p1PhoneIdx != -1 && p1PhoneIdx < len(row) {
			p1Phone = strings.TrimSpace(row[p1PhoneIdx])
		}

		results = append(results, EnrollmentApplication{
			SidiApplicationID:   sidiID,
			StudentFirstName:    firstName,
			StudentLastName:     lastName,
			StudentTaxCode:      taxCode,
			BirthDate:           birthDate,
			Gender:              gender,
			MiddleSchoolGrade:   grade,
			TrackChosen:         track,
			SecondLanguage:      lang,
			HasDisabilityL104:   l104,
			HasDSA:              dsa,
			ReligionChoice:      religion,
			RequestedClassmates: reqClassmates,
			Parent1FirstName:    p1First,
			Parent1LastName:     p1Last,
			Parent1Email:        p1Email,
			Parent1Phone:        p1Phone,
			Status:              "pending",
		})
	}

	return results, nil
}
