package grades

import (
	"encoding/csv"
	"io"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

type ImportRequest struct {
	StudentID     string
	SubjectID     string
	SchoolID      string
	GradeValue    float64
	GradeType     GradeType
	Date          time.Time
	Description   string
	GradeCategory GradeCategory
	Weight        float64
	Semester      int
}

// ImportResult and ImportError are defined in dto.go

// ParseCSVGrades parses a CSV file into ImportRequest slice.
// The semester parameter is applied to every row so the caller controls
// which academic period the imported grades belong to.
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
			continue
		} // Skip header
		if len(row) < 4 {
			continue
		}

		val, _ := strconv.ParseFloat(row[2], 64)
		date, _ := time.Parse("2006-01-02", row[3])

		cat := GradeCategorySummative
		if len(row) > 4 {
			cat = GradeCategory(row[4])
		}

		desc := ""
		if len(row) > 5 {
			desc = row[5]
		}

		requests = append(requests, ImportRequest{
			StudentID:     row[0],
			SubjectID:     row[1],
			GradeValue:    val,
			GradeType:     GradeTypeNumeric, // Default
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
// The semester parameter is applied to every row so the caller controls
// which academic period the imported grades belong to.
// fix: semester era precedentemente hardcodato a 1, ignorando il valore reale.
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
		// Try first sheet if Grades not found
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
		if len(row) < 4 {
			continue
		}

		val, _ := strconv.ParseFloat(row[2], 64)
		date, _ := time.Parse("2006-01-02", row[3])

		cat := GradeCategorySummative
		if len(row) > 4 {
			cat = GradeCategory(row[4])
		}

		desc := ""
		if len(row) > 5 {
			desc = row[5]
		}

		requests = append(requests, ImportRequest{
			StudentID:     row[0],
			SubjectID:     row[1],
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
// teacherID is used as both TeacherID and CreatedBy on each Grade.
// fix: SchoolID e CreatedBy erano assenti causando potenziali constraint
// violation su colonne NOT NULL nel database.
func ProcessBulkImport(repo Repository, reqs []ImportRequest, teacherID string) (ImportResult, error) {
	res := ImportResult{}

	var grades []*Grade
	for _, req := range reqs {
		grades = append(grades, &Grade{
			StudentID:     req.StudentID,
			SchoolID:      req.SchoolID,
			SubjectID:     req.SubjectID,
			GradeValue:    req.GradeValue,
			GradeType:     req.GradeType,
			Date:          req.Date,
			TeacherID:     teacherID,
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
			return res, err
		}
		res.Imported = len(grades)
	}

	return res, nil
}
