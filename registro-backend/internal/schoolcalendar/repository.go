package schoolcalendar

import (
	"context"
	"database/sql"
	"time"
)

type Repository interface {
	UpsertYear(s *SchoolYearSettings) error
	GetYear(schoolID string) (*SchoolYearSettings, error)

	AddNonTeachingDay(d *NonTeachingDay) error
	DeleteNonTeachingDay(schoolID, id string) error
	ListNonTeachingDays(schoolID string) ([]NonTeachingDay, error)

	// CountTeachingDays conta i giorni lavorativi (lun-ven) nell'intervallo
	// escludendo i giorni non didattici registrati.
	CountTeachingDays(schoolID string, from, to time.Time) (int, error)

	CreateAcademicPeriod(p *AcademicPeriod) error
	ListAcademicPeriods(schoolID string) ([]AcademicPeriod, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) UpsertYear(s *SchoolYearSettings) error {
	query := `
		INSERT INTO school_calendar_settings
			(school_id, year_label, start_date, end_date, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (school_id) DO UPDATE SET
			year_label = EXCLUDED.year_label,
			start_date = EXCLUDED.start_date,
			end_date   = EXCLUDED.end_date,
			updated_at = NOW()
		RETURNING id`
	return r.db.QueryRow(query, s.SchoolID, s.YearLabel, s.StartDate, s.EndDate, s.CreatedBy).Scan(&s.ID)
}

func (r *repository) GetYear(schoolID string) (*SchoolYearSettings, error) {
	var s SchoolYearSettings
	query := `SELECT id, school_id, year_label, start_date, end_date, created_by FROM school_calendar_settings WHERE school_id=$1`
	err := r.db.QueryRow(query, schoolID).Scan(&s.ID, &s.SchoolID, &s.YearLabel, &s.StartDate, &s.EndDate, &s.CreatedBy)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) AddNonTeachingDay(d *NonTeachingDay) error {
	query := `
		INSERT INTO school_non_teaching_days (school_id, date, label, created_by, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (school_id, date) DO UPDATE SET label = EXCLUDED.label
		RETURNING id`
	return r.db.QueryRow(query, d.SchoolID, d.Date, d.Label, d.CreatedBy).Scan(&d.ID)
}

func (r *repository) DeleteNonTeachingDay(schoolID, id string) error {
	_, err := r.db.ExecContext(context.Background(), `DELETE FROM school_non_teaching_days WHERE id=$1::uuid AND school_id=$2`, id, schoolID)
	return err
}

func (r *repository) ListNonTeachingDays(schoolID string) ([]NonTeachingDay, error) {
	rows, err := r.db.QueryContext(context.Background(), `SELECT id, school_id, date, label FROM school_non_teaching_days WHERE school_id=$1 ORDER BY date ASC`, schoolID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var res []NonTeachingDay
	for rows.Next() {
		var d NonTeachingDay
		if err := rows.Scan(&d.ID, &d.SchoolID, &d.Date, &d.Label); err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *repository) CountTeachingDays(schoolID string, from, to time.Time) (int, error) {
	// Genera la serie di date e sottrae i giorni non didattici e i weekend.
	query := `
		SELECT COUNT(*)
		FROM generate_series($1::date, $2::date, '1 day'::interval) AS d(date)
		WHERE
			EXTRACT(DOW FROM d.date) NOT IN (0, 6) -- Esclude sabato e domenica
			AND d.date NOT IN (
				SELECT date FROM school_non_teaching_days WHERE school_id = $3
			)`
	var count int
	err := r.db.QueryRowContext(context.Background(), query, from, to, schoolID).Scan(&count)
	return count, err
}

func (r *repository) CreateAcademicPeriod(p *AcademicPeriod) error {
	query := `
		INSERT INTO academic_periods (school_id, academic_year_id, name, code, start_date, end_date, is_current)
		VALUES ($1::uuid, $2::uuid, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	var yearUUID interface{} = nil
	if p.AcademicYearID != nil && *p.AcademicYearID != "" {
		yearUUID = *p.AcademicYearID
	}
	return r.db.QueryRowContext(context.Background(), query, p.SchoolID, yearUUID, p.Name, p.Code, p.StartDate, p.EndDate, p.IsCurrent).Scan(&p.ID, &p.CreatedAt)
}

func (r *repository) ListAcademicPeriods(schoolID string) ([]AcademicPeriod, error) {
	query := `
		SELECT id, school_id, academic_year_id, name, COALESCE(code, ''), start_date, end_date, is_current, created_at
		FROM academic_periods
		WHERE school_id = $1::uuid
		ORDER BY start_date ASC
	`
	rows, err := r.db.QueryContext(context.Background(), query, schoolID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var list []AcademicPeriod
	for rows.Next() {
		var p AcademicPeriod
		var yearID sql.NullString
		if err := rows.Scan(&p.ID, &p.SchoolID, &yearID, &p.Name, &p.Code, &p.StartDate, &p.EndDate, &p.IsCurrent, &p.CreatedAt); err != nil {
			return nil, err
		}
		if yearID.Valid {
			p.AcademicYearID = &yearID.String
		}
		list = append(list, p)
	}
	return list, rows.Err()
}
