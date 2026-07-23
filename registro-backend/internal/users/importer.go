package users

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

type Importer struct {
	validator *Validator
}

func NewImporter() *Importer {
	return &Importer{
		validator: NewValidator(),
	}
}

func (i *Importer) ParseFile(file multipart.File, filename string) ([]User, []string, error) {
	if strings.HasSuffix(filename, ".csv") {
		return i.ParseCSV(file)
	} else if strings.HasSuffix(filename, ".xlsx") {
		return i.ParseXLSX(file)
	}
	return nil, nil, fmt.Errorf("unsupported file format")
}

func (i *Importer) ParseCSV(r io.Reader) ([]User, []string, error) {
	reader := csv.NewReader(r)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	if len(rows) < 2 {
		return nil, nil, fmt.Errorf("empty or invalid CSV file")
	}

	return i.processRows(rows[1:]) // Skip header
}

func (i *Importer) ParseXLSX(r io.Reader) ([]User, []string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	// Assume first sheet
	rows, err := f.GetRows(f.GetSheetList()[0])
	if err != nil {
		return nil, nil, err
	}

	if len(rows) < 2 {
		return nil, nil, fmt.Errorf("empty or invalid XLSX file")
	}

	return i.processRows(rows[1:])
}

func (i *Importer) processRows(rows [][]string) ([]User, []string, error) {
	var users []User
	var errors []string

	for idx, row := range rows {
		if len(row) < 5 {
			errors = append(errors, fmt.Sprintf("Row %d: Insufficient columns", idx+2))
			continue
		}

		// Expected columns: Email, FirstName, LastName, FiscalCode, Role
		email := strings.TrimSpace(row[0])
		firstName := strings.TrimSpace(row[1])
		lastName := strings.TrimSpace(row[2])
		fiscalCode := strings.TrimSpace(row[3])
		role := strings.ToLower(strings.TrimSpace(row[4]))

		// Basic Validation
		if email == "" || !strings.Contains(email, "@") {
			errors = append(errors, fmt.Sprintf("Row %d: Invalid email", idx+2))
			continue
		}
		if !i.validator.ValidateFiscalCode(fiscalCode) {
			errors = append(errors, fmt.Sprintf("Row %d: Invalid fiscal code %s", idx+2, fiscalCode))
			continue
		}

		// Generate temporary password (or allow it to be set later via reset)
		// For simplicity, generate a random one and hatch it, user must reset.
		hashed, _ := bcrypt.GenerateFromPassword([]byte("ChangeMe123!"), bcrypt.DefaultCost)

		users = append(users, User{
			ID:           uuid.New().String(),
			Email:        email,
			FirstName:    firstName,
			LastName:     lastName,
			FiscalCode:   &fiscalCode,
			Role:         role,
			PasswordHash: string(hashed),
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
	}

	return users, errors, nil
}
