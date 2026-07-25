package grades

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type ImportRequest struct {
	StudentID     string
	SubjectID     string
	SchoolID      string
	TeacherID     string
	GradeValue    float64
	GradeType     GradeType
	Date          time.Time
	Description   string
	GradeCategory GradeCategory
	Weight        float64
	Semester      int
}

// ParseCSVGrades parses a CSV file into ImportRequest slice.
func ParseCSVGrades(r io.Reader, semester int) ([]ImportRequest, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1 // Allow variable fields
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if semester != 1 && semester != 2 {
		semester = 1 // safe default
	}

	var requests []ImportRequest
	// Expect Header: StudentID, SubjectID, Value, Date, Category...
	for i, row := range rows {
		if i == 0 {
			continue // Skip header
		}
		if len(row) < 3 {
			continue
		}

		val, _ := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)

		var date time.Time
		if len(row) > 3 {
			d, err := time.Parse("2006-01-02", strings.TrimSpace(row[3]))
			if err == nil {
				date = d
			}
		}

		cat := GradeCategorySummative
		if len(row) > 4 && strings.TrimSpace(row[4]) != "" {
			cat = GradeCategory(strings.TrimSpace(row[4]))
		}

		desc := ""
		if len(row) > 5 {
			desc = strings.TrimSpace(row[5])
		}

		requests = append(requests, ImportRequest{
			StudentID:     strings.TrimSpace(row[0]),
			SubjectID:     strings.TrimSpace(row[1]),
			GradeValue:    val,
			GradeType:     GradeTypeNumeric,
			Date:          date,
			GradeCategory: cat,
			Description:   desc,
			Weight:        1.0,
			Semester:      semester,
		})
	}
	return requests, nil
}

// ParseXLSXGrades parses an XLSX file into ImportRequest slice.
func ParseXLSXGrades(r io.Reader, semester int) ([]ImportRequest, error) {
	if semester != 1 && semester != 2 {
		semester = 1 // safe default
	}

	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("Grades")
	if err != nil {
		sheets := f.GetSheetList()
		if len(sheets) > 0 {
			rows, err = f.GetRows(sheets[0])
		}
	}
	if err != nil {
		return nil, err
	}

	var requests []ImportRequest
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 3 {
			continue
		}

		val, _ := strconv.ParseFloat(strings.TrimSpace(row[2]), 64)

		var date time.Time
		if len(row) > 3 {
			d, err := time.Parse("2006-01-02", strings.TrimSpace(row[3]))
			if err == nil {
				date = d
			}
		}

		cat := GradeCategorySummative
		if len(row) > 4 && strings.TrimSpace(row[4]) != "" {
			cat = GradeCategory(strings.TrimSpace(row[4]))
		}

		desc := ""
		if len(row) > 5 {
			desc = strings.TrimSpace(row[5])
		}

		requests = append(requests, ImportRequest{
			StudentID:     strings.TrimSpace(row[0]),
			SubjectID:     strings.TrimSpace(row[1]),
			GradeValue:    val,
			GradeType:     GradeTypeNumeric,
			Date:          date,
			Semester:      semester,
			GradeCategory: cat,
			Description:   desc,
			Weight:        1.0,
		})
	}
	return requests, nil
}

// ProcessBulkImport inserts the parsed ImportRequests via the repository.
func ProcessBulkImport(repo Repository, reqs []ImportRequest, teacherID string, teacherProfileID string, schoolID string) (ImportResult, error) {
	res := ImportResult{}

	var grades []*Grade
	for _, req := range reqs {
		sID := req.SchoolID
		if sID == "" {
			sID = schoolID
		}
		tID := req.TeacherID
		if tID == "" {
			tID = teacherProfileID
		}
		if tID == "" {
			tID = teacherID
		}

		grades = append(grades, &Grade{
			StudentID:     req.StudentID,
			SchoolID:      sID,
			SubjectID:     req.SubjectID,
			GradeValue:    req.GradeValue,
			GradeType:     req.GradeType,
			Date:          req.Date,
			TeacherID:     tID,
			CreatedBy:     teacherID,
			Semester:      Semester(req.Semester),
			GradeCategory: req.GradeCategory,
			Weight:        req.Weight,
			Description:   req.Description,
			IsPublished:   true,
		})
	}

	if len(grades) > 0 {
		err := repo.BatchCreate(grades)
		if err != nil {
			return res, fmt.Errorf("failed to import grades batch: %w", err)
		}
		res.Imported = len(grades)
	}

	return res, nil
}
