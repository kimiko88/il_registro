package competencies_invalsi

import (
	"database/sql"
)

type Repository interface {
	UpsertCertification(c *CompetencyCertificationDM14) error
	GetCertification(id string) (*CompetencyCertificationDM14, error)
	GetCertificationByStudent(studentID, schoolYear, schoolTier string) (*CompetencyCertificationDM14, error)
	ListCertificationsByClass(classID, schoolYear string) ([]CompetencyCertificationDM14, error)
	SignCertification(id, principalID string) error

	InsertInvalsiBatch(results []InvalsiResult) error
	ListInvalsiByClass(classID, schoolYear string) ([]InvalsiResult, error)
	GetSchoolInvalsiAnalytics(schoolID, schoolYear string) ([]InvalsiResult, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) resolveSchoolID(schoolID, studentID, classID string) string {
	if schoolID != "" {
		return schoolID
	}
	if classID != "" {
		var sID sql.NullString
		_ = r.db.QueryRow(`SELECT school_id::text FROM classes WHERE id::text = $1`, classID).Scan(&sID)
		if sID.Valid && sID.String != "" {
			return sID.String
		}
	}
	if studentID != "" {
		var sID sql.NullString
		_ = r.db.QueryRow(`SELECT school_id::text FROM users WHERE id::text = $1 AND school_id IS NOT NULL`, studentID).Scan(&sID)
		if sID.Valid && sID.String != "" {
			return sID.String
		}
	}
	var defaultSchoolID string
	_ = r.db.QueryRow(`SELECT id::text FROM schools ORDER BY created_at ASC LIMIT 1`).Scan(&defaultSchoolID)
	return defaultSchoolID
}

func (r *repository) UpsertCertification(c *CompetencyCertificationDM14) error {
	c.SchoolID = r.resolveSchoolID(c.SchoolID, c.StudentID, c.ClassID)
	query := `
		INSERT INTO competency_certifications_dm14 (
			school_id, student_id, class_id, school_tier, school_year,
			rubric_evaluations_json, final_considerations, council_deliberation_date, status
		) VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7, $8, $9)
		ON CONFLICT (student_id, school_year, school_tier)
		DO UPDATE SET
			rubric_evaluations_json = EXCLUDED.rubric_evaluations_json,
			final_considerations = EXCLUDED.final_considerations,
			council_deliberation_date = EXCLUDED.council_deliberation_date,
			status = EXCLUDED.status,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query, c.SchoolID, c.StudentID, c.ClassID, c.SchoolTier, c.SchoolYear,
		c.RubricEvaluationsJSON, c.FinalConsiderations, c.CouncilDeliberationDate, c.Status,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *repository) GetCertification(id string) (*CompetencyCertificationDM14, error) {
	query := `
		SELECT c.id, c.school_id, c.student_id, c.class_id, c.school_tier, c.school_year,
		       c.rubric_evaluations_json::text, COALESCE(c.final_considerations, ''),
		       TO_CHAR(c.council_deliberation_date, 'YYYY-MM-DD'), c.status,
		       c.signed_by_principal_id, c.signed_at, c.created_at, c.updated_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') as student_name
		FROM competency_certifications_dm14 c
		LEFT JOIN users u ON c.student_id = u.id
		WHERE c.id = $1
	`
	var cert CompetencyCertificationDM14
	err := r.db.QueryRow(query, id).Scan(
		&cert.ID, &cert.SchoolID, &cert.StudentID, &cert.ClassID, &cert.SchoolTier, &cert.SchoolYear,
		&cert.RubricEvaluationsJSON, &cert.FinalConsiderations, &cert.CouncilDeliberationDate, &cert.Status,
		&cert.SignedByPrincipalID, &cert.SignedAt, &cert.CreatedAt, &cert.UpdatedAt, &cert.StudentFullName,
	)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *repository) GetCertificationByStudent(studentID, schoolYear, schoolTier string) (*CompetencyCertificationDM14, error) {
	query := `
		SELECT c.id, c.school_id, c.student_id, c.class_id, c.school_tier, c.school_year,
		       c.rubric_evaluations_json::text, COALESCE(c.final_considerations, ''),
		       TO_CHAR(c.council_deliberation_date, 'YYYY-MM-DD'), c.status,
		       c.signed_by_principal_id, c.signed_at, c.created_at, c.updated_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') as student_name
		FROM competency_certifications_dm14 c
		LEFT JOIN users u ON c.student_id = u.id
		WHERE c.student_id = $1 AND c.school_year = $2 AND c.school_tier = $3
	`
	var cert CompetencyCertificationDM14
	err := r.db.QueryRow(query, studentID, schoolYear, schoolTier).Scan(
		&cert.ID, &cert.SchoolID, &cert.StudentID, &cert.ClassID, &cert.SchoolTier, &cert.SchoolYear,
		&cert.RubricEvaluationsJSON, &cert.FinalConsiderations, &cert.CouncilDeliberationDate, &cert.Status,
		&cert.SignedByPrincipalID, &cert.SignedAt, &cert.CreatedAt, &cert.UpdatedAt, &cert.StudentFullName,
	)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *repository) ListCertificationsByClass(classID, schoolYear string) ([]CompetencyCertificationDM14, error) {
	query := `
		SELECT c.id, c.school_id, c.student_id, c.class_id, c.school_tier, c.school_year,
		       c.rubric_evaluations_json::text, COALESCE(c.final_considerations, ''),
		       TO_CHAR(c.council_deliberation_date, 'YYYY-MM-DD'), c.status,
		       c.signed_by_principal_id, c.signed_at, c.created_at, c.updated_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') as student_name
		FROM competency_certifications_dm14 c
		LEFT JOIN users u ON c.student_id = u.id
		WHERE (c.class_id::text = $1 OR $1 = '') AND ($2 = '' OR c.school_year = $2)
		ORDER BY student_name ASC
	`
	rows, err := r.db.Query(query, classID, schoolYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []CompetencyCertificationDM14
	for rows.Next() {
		var cert CompetencyCertificationDM14
		if err := rows.Scan(
			&cert.ID, &cert.SchoolID, &cert.StudentID, &cert.ClassID, &cert.SchoolTier, &cert.SchoolYear,
			&cert.RubricEvaluationsJSON, &cert.FinalConsiderations, &cert.CouncilDeliberationDate, &cert.Status,
			&cert.SignedByPrincipalID, &cert.SignedAt, &cert.CreatedAt, &cert.UpdatedAt, &cert.StudentFullName,
		); err != nil {
			return nil, err
		}
		list = append(list, cert)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) SignCertification(id, principalID string) error {
	query := `
		UPDATE competency_certifications_dm14
		SET status = 'firmata_ds', signed_by_principal_id = $2, signed_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id, principalID)
	return err
}

func (r *repository) InsertInvalsiBatch(results []InvalsiResult) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO invalsi_results (
			school_id, class_id, student_id, school_year, grade_level, subject,
			score_wle, level, school_mean, national_mean, report_class_grade
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range results {
		if _, err := stmt.Exec(
			item.SchoolID, item.ClassID, item.StudentID, item.SchoolYear, item.GradeLevel,
			item.Subject, item.ScoreWLE, item.Level, item.SchoolMean, item.NationalMean, item.ReportClassGrade,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *repository) ListInvalsiByClass(classID, schoolYear string) ([]InvalsiResult, error) {
	query := `
		SELECT id, school_id, class_id, student_id, school_year, grade_level, subject,
		       score_wle, level, school_mean, national_mean, report_class_grade, created_at
		FROM invalsi_results
		WHERE class_id = $1 AND ($2 = '' OR school_year = $2)
		ORDER BY subject ASC, score_wle DESC
	`
	rows, err := r.db.Query(query, classID, schoolYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []InvalsiResult
	for rows.Next() {
		var item InvalsiResult
		if err := rows.Scan(
			&item.ID, &item.SchoolID, &item.ClassID, &item.StudentID, &item.SchoolYear, &item.GradeLevel,
			&item.Subject, &item.ScoreWLE, &item.Level, &item.SchoolMean, &item.NationalMean,
			&item.ReportClassGrade, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		item.ScoreDeltaVsNational = item.ScoreWLE - item.NationalMean
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) GetSchoolInvalsiAnalytics(schoolID, schoolYear string) ([]InvalsiResult, error) {
	query := `
		SELECT id, school_id, class_id, student_id, school_year, grade_level, subject,
		       score_wle, level, school_mean, national_mean, report_class_grade, created_at
		FROM invalsi_results
		WHERE school_id = $1 AND ($2 = '' OR school_year = $2)
		ORDER BY grade_level ASC, subject ASC
	`
	rows, err := r.db.Query(query, schoolID, schoolYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []InvalsiResult
	for rows.Next() {
		var item InvalsiResult
		if err := rows.Scan(
			&item.ID, &item.SchoolID, &item.ClassID, &item.StudentID, &item.SchoolYear, &item.GradeLevel,
			&item.Subject, &item.ScoreWLE, &item.Level, &item.SchoolMean, &item.NationalMean,
			&item.ReportClassGrade, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		item.ScoreDeltaVsNational = item.ScoreWLE - item.NationalMean
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}
