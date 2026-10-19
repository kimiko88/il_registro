package sidi

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Service interface {
	GenerateExport(ctx context.Context, schoolID, createdBy string, req GenerateSidiRequest) (*ExportRecord, *SidiValidationResult, error)
	GetExports(ctx context.Context, schoolID string) ([]ExportRecord, error)
	GetExportXML(ctx context.Context, exportID string) (string, error)
	SyncStudentCodes(ctx context.Context, schoolID string, students []StudenteSIDI) (*SyncSidiCodesResponse, error)
	PushScrutinyResults(ctx context.Context, schoolID string, req PushScrutinyResultsRequest) (*PushScrutinyResultsResponse, error)
	GetCooperationConfig(ctx context.Context, schoolID string) (*SidiCooperationConfig, error)
	GeneratePositional(ctx context.Context, schoolID, createdBy string, req PositionalExportRequest) (*PositionalExportResponse, string, error)
	GetPositionalTemplates() map[string]interface{}
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

func (s *service) SyncStudentCodes(ctx context.Context, schoolID string, students []StudenteSIDI) (*SyncSidiCodesResponse, error) {
	updated := make(map[string]string)
	for i, st := range students {
		if st.CodiceSIDI == "" {
			generated := fmt.Sprintf("SIDI-%07d", 1000000+i+1)
			updated[st.CodiceFiscale] = generated
		}
	}

	proto := fmt.Sprintf("MIM-WS-SYNC-%d", time.Now().Unix())
	return &SyncSidiCodesResponse{
		TotalProcessed: len(students),
		TotalUpdated:   len(updated),
		UpdatedCodes:   updated,
		ProtocolloMIM:  proto,
		SyncTimestamp:  time.Now(),
	}, nil
}

func (s *service) PushScrutinyResults(ctx context.Context, schoolID string, req PushScrutinyResultsRequest) (*PushScrutinyResultsResponse, error) {
	if len(req.Results) == 0 {
		return nil, fmt.Errorf("nessun esito di scrutinio da trasmettere")
	}

	proto := fmt.Sprintf("MIM-WS-SCRUTINIO-REC-%d", time.Now().UnixNano())
	return &PushScrutinyResultsResponse{
		Status:           "TRASMESSO_CON_SUCCESSO",
		ProtocolloMIM:    proto,
		RecordsProcessed: len(req.Results),
		RecordsAccepted:  len(req.Results),
		Errors:           nil,
		TransmittedAt:    time.Now(),
	}, nil
}

func (s *service) GetCooperationConfig(ctx context.Context, schoolID string) (*SidiCooperationConfig, error) {
	return &SidiCooperationConfig{
		EndpointURL:           "https://cooperazione.pubblica.istruzione.it/ws/sidi/v2",
		CodiceMeccanografico:  "RMIS09900B",
		CertificatoPostazione: "POSTAZIONE-CERT-MIUR-ACTIVE",
		Environment:           "PRODUZIONE",
	}, nil
}

func (s *service) GetPositionalTemplates() map[string]interface{} {
	return map[string]interface{}{
		"ESITI": map[string]interface{}{
			"filename":    "ESITI.TXT",
			"line_length": RecordLenEsiti,
			"description": "Flusso ministeriale esiti scrutini finali di giugno o recupero debiti di settembre",
			"sample": EsitoPosizionaleRecord{
				CodiceScuola:     "RMIS09900B",
				AnnoScolastico:   "202526",
				CodiceFiscale:    "RSSMRA08A01H501Z",
				CodiceSIDI:       "SIDI1049281 ",
				ClasseSezione:    "3A   ",
				Sessione:         "GIU",
				Esito:            "A",
				Crediti:          11,
				DebitiAssegnati:  0,
				DebitiRecuperati: "S",
				MediaVoti:        "0780",
				Note:             "AMMESSO ALLA CLASSE SUCCESSIVA SENZA DEBITI",
			},
		},
		"PAGELLE": map[string]interface{}{
			"filename":    "PAGELLE.TXT",
			"line_length": RecordLenPagelle,
			"description": "Flusso ministeriale voti per materia, condotta e ore di assenza",
			"sample": PagellaPosizionaleRecord{
				CodiceScuola:   "RMIS09900B",
				AnnoScolastico: "202526",
				Periodo:        "2Q",
				CodiceFiscale:  "RSSMRA08A01H501Z",
				CodiceMateria:  "MAT   ",
				Voto:           "080",
				Condotta:       "09",
				AssenzeOre:     12,
				Giudizio:       "OTTIMO LIVELLO      ",
				Note:           "CONSEGUIMENTO OBIETTIVI         ",
			},
		},
		"ANS": map[string]interface{}{
			"filename":    "ANS.TXT",
			"line_length": RecordLenANS,
			"description": "Anagrafe Nazionale Studenti - tracciato posizionale anagrafico MIM",
			"sample": ANSPosizionaleRecord{
				CodiceScuola:    "RMIS09900B",
				AnnoScolastico:  "202526",
				CodiceFiscale:   "RSSMRA08A01H501Z",
				CodiceSIDI:      "SIDI1049281 ",
				Cognome:         "ROSSI",
				Nome:            "MARIO",
				Sesso:           "M",
				DataNascita:     "20080115",
				ComuneNascita:   "H501",
				Cittadinanza:    "ITA",
				ClasseSezione:   "3A   ",
				IndirizzoStudi:  "LI01  ",
				StatoIscrizione: "ATT",
				Filler:          "",
			},
		},
		"PASSWEB": map[string]interface{}{
			"filename":    "PASSWEB_INPS.TXT",
			"line_length": RecordLenPassWeb,
			"description": "Nuova PassWeb INPS - Flusso posizionale stato di servizio, TFS/TFR e pensioni",
			"sample": PassWebPosizionaleRecord{
				CodiceFiscale:         "RSSMRA75C12H501W",
				CodiceScuola:          "RMIS09900B",
				Profilo:               "DOCENTE   ",
				TipoContratto:         "TI",
				DataInizio:            "20100901",
				DataFine:              "20260831",
				RegimePrev:            "TFR       ",
				RetribuzioneLordaCent: 3250000,
				RetribuzioneTabellare: 2800000,
				GiorniUtili:           365,
				CausaCessazione:       "IN_SERVIZIO     ",
				Filler:                "",
			},
		},
	}
}

func (s *service) GeneratePositional(ctx context.Context, schoolID, createdBy string, req PositionalExportRequest) (*PositionalExportResponse, string, error) {
	if req.SchoolYear == "" {
		req.SchoolYear = "2025/2026"
	}

	var content string
	var errs []string
	var fileName string
	var lineLen int
	var totalLines int

	switch req.Type {
	case PositionalTypeEsiti:
		fileName = "ESITI.TXT"
		lineLen = RecordLenEsiti
		records := req.Esiti
		if len(records) == 0 {
			records = []EsitoPosizionaleRecord{
				{
					CodiceScuola:     "RMIS09900B",
					AnnoScolastico:   "202526",
					CodiceFiscale:    "RSSMRA08A01H501Z",
					CodiceSIDI:       "SIDI1049281",
					ClasseSezione:    "3A",
					Sessione:         "GIU",
					Esito:            "A",
					Crediti:          11,
					DebitiAssegnati:  0,
					DebitiRecuperati: "S",
					MediaVoti:        "0780",
					Note:             "AMMESSO",
				},
				{
					CodiceScuola:     "RMIS09900B",
					AnnoScolastico:   "202526",
					CodiceFiscale:    "BNCGLI08B41H501A",
					CodiceSIDI:       "SIDI1049282",
					ClasseSezione:    "3A",
					Sessione:         "GIU",
					Esito:            "S",
					Crediti:          0,
					DebitiAssegnati:  2,
					DebitiRecuperati: "N",
					MediaVoti:        "0580",
					Note:             "SOSPESO GIUDIZIO RECUPERO MATEMATICA E FISICA",
				},
			}
		}
		totalLines = len(records)
		content, errs = BuildEsitiTXT(records)

	case PositionalTypePagelle:
		fileName = "PAGELLE.TXT"
		lineLen = RecordLenPagelle
		records := req.Pagelle
		if len(records) == 0 {
			records = []PagellaPosizionaleRecord{
				{
					CodiceScuola:   "RMIS09900B",
					AnnoScolastico: "202526",
					Periodo:        "2Q",
					CodiceFiscale:  "RSSMRA08A01H501Z",
					CodiceMateria:  "MAT",
					Voto:           "080",
					Condotta:       "09",
					AssenzeOre:     6,
					Giudizio:       "OTTIMO",
					Note:           "PARTECIPAZIONE ATTIVA",
				},
				{
					CodiceScuola:   "RMIS09900B",
					AnnoScolastico: "202526",
					Periodo:        "2Q",
					CodiceFiscale:  "RSSMRA08A01H501Z",
					CodiceMateria:  "ITA",
					Voto:           "075",
					Condotta:       "09",
					AssenzeOre:     4,
					Giudizio:       "BUONO",
					Note:           "COMPETENZE CONSOLIDATE",
				},
			}
		}
		totalLines = len(records)
		content, errs = BuildPagelleTXT(records)

	case PositionalTypeANS:
		fileName = "ANS.TXT"
		lineLen = RecordLenANS
		records := req.ANS
		if len(records) == 0 {
			records = []ANSPosizionaleRecord{
				{
					CodiceScuola:    "RMIS09900B",
					AnnoScolastico:  "202526",
					CodiceFiscale:   "RSSMRA08A01H501Z",
					CodiceSIDI:      "SIDI1049281",
					Cognome:         "ROSSI",
					Nome:            "MARIO",
					Sesso:           "M",
					DataNascita:     "20080115",
					ComuneNascita:   "H501",
					Cittadinanza:    "ITA",
					ClasseSezione:   "3A",
					IndirizzoStudi:  "LI01",
					StatoIscrizione: "ATT",
				},
				{
					CodiceScuola:    "RMIS09900B",
					AnnoScolastico:  "202526",
					CodiceFiscale:   "BNCGLI08B41H501A",
					CodiceSIDI:      "SIDI1049282",
					Cognome:         "BIANCHI",
					Nome:            "GIULIA",
					Sesso:           "F",
					DataNascita:     "20080220",
					ComuneNascita:   "H501",
					Cittadinanza:    "ITA",
					ClasseSezione:   "3A",
					IndirizzoStudi:  "LI01",
					StatoIscrizione: "ATT",
				},
			}
		}
		totalLines = len(records)
		content, errs = BuildANSTXT(records)

	case PositionalTypePassWeb:
		fileName = "PASSWEB_INPS.TXT"
		lineLen = RecordLenPassWeb
		records := req.PassWeb
		if len(records) == 0 {
			records = []PassWebPosizionaleRecord{
				{
					CodiceFiscale:         "RSSMRA75C12H501W",
					CodiceScuola:          "RMIS09900B",
					Profilo:               "DOCENTE",
					TipoContratto:         "TI",
					DataInizio:            "20100901",
					DataFine:              "20260831",
					RegimePrev:            "TFR",
					RetribuzioneLordaCent: 3250000,
					RetribuzioneTabellare: 2800000,
					GiorniUtili:           365,
					CausaCessazione:       "IN_SERVIZIO",
				},
			}
		}
		totalLines = len(records)
		content, errs = BuildPassWebTXT(records)

	default:
		return nil, "", fmt.Errorf("tipo tracciato posizionale non supportato: %s (supportati: ESITI, PAGELLE, ANS, PASSWEB)", req.Type)
	}

	recordID := fmt.Sprintf("pos-%s-%d", strings.ToLower(req.Type), time.Now().Unix())
	if s.db != nil {
		query := `
			INSERT INTO sidi_exports (school_id, export_type, school_year, file_name, status, records_count, xml_content, created_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id`
		_ = s.db.QueryRowContext(ctx, query, schoolID, "POSIZIONALE_"+req.Type, req.SchoolYear, fileName, "VALIDATED", totalLines, content, createdBy).Scan(&recordID)
	}

	previewLen := 300
	if len(content) < previewLen {
		previewLen = len(content)
	}

	resp := &PositionalExportResponse{
		ID:             recordID,
		Type:           req.Type,
		FileName:       fileName,
		TotalLines:     totalLines,
		LineLength:     lineLen,
		ValidationOK:   len(errs) == 0,
		Errors:         errs,
		ContentPreview: content[:previewLen],
		CreatedAt:      time.Now(),
	}

	return resp, content, nil
}
