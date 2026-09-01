package sidi

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Service interface {
	GenerateExport(ctx context.Context, schoolID, createdBy string, req GenerateSidiRequest) (*ExportRecord, *SidiValidationResult, error)
	GetExports(ctx context.Context, schoolID string) ([]ExportRecord, error)
	GetExportXML(ctx context.Context, exportID string) (string, error)
}

type service struct {
	db      *sql.DB
	builder *Builder
}

func NewService(db *sql.DB) Service {
	return &service{
		db:      db,
		builder: NewBuilder(),
	}
}

func (s *service) GenerateExport(ctx context.Context, schoolID, createdBy string, req GenerateSidiRequest) (*ExportRecord, *SidiValidationResult, error) {
	if req.SchoolYear == "" {
		req.SchoolYear = "2025/2026"
	}

	// Sample student data from DB or simulated
	students := []StudenteSIDI{
		{
			CodiceSIDI:    "SIDI-1049281",
			CodiceFiscale: "RSSMRA08A01H501Z",
			Cognome:       "Rossi",
			Nome:          "Mario",
			Classe:        "3ª A",
		},
		{
			CodiceSIDI:    "SIDI-1049282",
			CodiceFiscale: "BNCGLI08B41H501A",
			Cognome:       "Bianchi",
			Nome:          "Giulia",
			Classe:        "3ª A",
		},
	}

	var scrutini []ScrutinioSIDI
	if req.ExportType == "SCRUTINIO_GIUGNO" || req.ExportType == "SCRUTINIO_SETTEMBRE_DEBITI" {
		scrutini = []ScrutinioSIDI{
			{
				CodiceSIDI:       "SIDI-1049281",
				Classe:           "3ª A",
				EsitoFinale:      "AMMESSO",
				CreditiFormat:    11,
				DebitiRecuperati: true,
			},
		}
	}

	validation := s.builder.ValidateSidiData(students)

	xmlData, err := s.builder.BuildSidiXML("RMIS09900B", "Istituto Superiore Statale", req.SchoolYear, req.ExportType, students, scrutini)
	if err != nil {
		return nil, nil, fmt.Errorf("XML generation failed: %w", err)
	}

	fileName := fmt.Sprintf("FLUSSO_SIDI_%s_%d.xml", req.ExportType, time.Now().Unix())
	record := &ExportRecord{
		SchoolID:     schoolID,
		ExportType:   req.ExportType,
		SchoolYear:   req.SchoolYear,
		FileName:     fileName,
		Status:       "VALIDATED",
		RecordsCount: len(students),
		XMLContent:   string(xmlData),
		CreatedBy:    createdBy,
		CreatedAt:    time.Now(),
	}

	query := `
		INSERT INTO sidi_exports (school_id, export_type, school_year, file_name, status, records_count, xml_content, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`
	if s.db != nil {
		_ = s.db.QueryRowContext(ctx, query, record.SchoolID, record.ExportType, record.SchoolYear, record.FileName, record.Status, record.RecordsCount, record.XMLContent, record.CreatedBy).Scan(&record.ID, &record.CreatedAt)
	} else {
		record.ID = "simulated-sidi-id"
	}

	return record, &validation, nil
}

func (s *service) GetExports(ctx context.Context, schoolID string) ([]ExportRecord, error) {
	if s.db == nil {
		return []ExportRecord{
			{
				ID:           "exp-1",
				SchoolID:     schoolID,
				ExportType:   "ANS_ANAGRAFE",
				SchoolYear:   "2025/2026",
				FileName:     "FLUSSO_SIDI_ANS_ANAGRAFE.xml",
				Status:       "EXPORTED",
				RecordsCount: 2,
				CreatedAt:    time.Now(),
			},
		}, nil
	}

	query := `
		SELECT id, school_id, export_type, school_year, file_name, status, records_count, created_by, created_at
		FROM sidi_exports
		WHERE school_id = $1 OR $1 = ''
		ORDER BY created_at DESC`
	rows, err := s.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return []ExportRecord{}, nil
	}
	defer func() { _ = rows.Close() }()

	var list []ExportRecord
	for rows.Next() {
		var r ExportRecord
		if err := rows.Scan(&r.ID, &r.SchoolID, &r.ExportType, &r.SchoolYear, &r.FileName, &r.Status, &r.RecordsCount, &r.CreatedBy, &r.CreatedAt); err == nil {
			list = append(list, r)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (s *service) GetExportXML(ctx context.Context, exportID string) (string, error) {
	if s.db == nil {
		return "<FlussoSIDI></FlussoSIDI>", nil
	}
	var content string
	err := s.db.QueryRowContext(ctx, "SELECT xml_content FROM sidi_exports WHERE id = $1::uuid", exportID).Scan(&content)
	return content, err
}
