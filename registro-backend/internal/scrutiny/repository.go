package scrutiny

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	SaveRecord(ctx context.Context, record *ScrutinyRecord) error
	GetRecord(ctx context.Context, studentID, classID string, semester int) (*ScrutinyRecord, error)
	ListRecordsByClass(ctx context.Context, classID string, semester int) ([]ScrutinyRecord, error)
	ValidateClassScrutiny(ctx context.Context, classID string, semester int, validatorID string) error
	UpdateClassScrutinyStatus(ctx context.Context, classID string, semester int, status string) error

	// Deficiency & Deferred Scrutiny Methods
	SaveDeficiency(ctx context.Context, def *StudentDeficiency) error
	GetDeficienciesByStudent(ctx context.Context, studentID string) ([]StudentDeficiency, error)
	GetDeficienciesByClass(ctx context.Context, classID string, semester int) ([]StudentDeficiency, error)
	SaveDeferredScrutiny(ctx context.Context, req *SaveDeferredScrutinyRequest) error
}

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) SaveRecord(ctx context.Context, rec *ScrutinyRecord) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if rec.ID == "" {
		rec.ID = uuid.New().String()
	}
	rec.UpdatedAt = time.Now()
	if rec.Status == "" {
		rec.Status = "in_progress"
	}

	query := `
		INSERT INTO scrutiny_records (id, student_id, class_id, semester, conduct_grade, final_decision, notes, coordinator_id, status, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (student_id, class_id, semester) 
		DO UPDATE SET conduct_grade = EXCLUDED.conduct_grade, final_decision = EXCLUDED.final_decision, notes = EXCLUDED.notes, status = EXCLUDED.status, updated_at = EXCLUDED.updated_at
		RETURNING id
	`
	err = tx.QueryRowContext(ctx, query, rec.ID, rec.StudentID, rec.ClassID, rec.Semester, rec.ConductGrade, rec.FinalDecision, rec.Notes, rec.CoordinatorID, rec.Status, rec.UpdatedAt).Scan(&rec.ID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM scrutiny_grades WHERE scrutiny_record_id = $1", rec.ID)
	if err != nil {
		return err
	}

	for _, g := range rec.Grades {
		g.ID = uuid.New().String()
		_, err = tx.ExecContext(ctx, "INSERT INTO scrutiny_grades (id, scrutiny_record_id, subject_id, final_grade, teacher_id) VALUES ($1, $2, $3, $4, $5)",
			g.ID, rec.ID, g.SubjectID, g.FinalGrade, g.TeacherID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *postgresRepository) GetRecord(ctx context.Context, studentID, classID string, semester int) (*ScrutinyRecord, error) {
	if studentID == "" || classID == "" {
		return nil, nil
	}
	query := `
		SELECT id, student_id, class_id, semester, conduct_grade, final_decision, notes, coordinator_id,
		       COALESCE(status, 'draft'), validated_by::text, validated_at, created_at, updated_at 
		FROM scrutiny_records WHERE student_id::text = $1 AND class_id::text = $2 AND semester = $3
	`
	rec := &ScrutinyRecord{}
	var valBy sql.NullString
	var valAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, studentID, classID, semester).Scan(
		&rec.ID, &rec.StudentID, &rec.ClassID, &rec.Semester, &rec.ConductGrade, &rec.FinalDecision,
		&rec.Notes, &rec.CoordinatorID, &rec.Status, &valBy, &valAt, &rec.CreatedAt, &rec.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if valBy.Valid {
		rec.ValidatedBy = valBy.String
	}
	if valAt.Valid {
		rec.ValidatedAt = &valAt.Time
	}

	gRows, err := r.db.QueryContext(ctx, "SELECT id, scrutiny_record_id, subject_id, final_grade, teacher_id FROM scrutiny_grades WHERE scrutiny_record_id::text = $1", rec.ID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = gRows.Close() }()

	for gRows.Next() {
		var g ScrutinyGrade
		if err := gRows.Scan(&g.ID, &g.ScrutinyRecordID, &g.SubjectID, &g.FinalGrade, &g.TeacherID); err != nil {
			return nil, err
		}
		rec.Grades = append(rec.Grades, g)
	}
	if err := gRows.Err(); err != nil {
		return nil, err
	}

	return rec, nil
}

func (r *postgresRepository) ListRecordsByClass(ctx context.Context, classID string, semester int) ([]ScrutinyRecord, error) {
	if classID == "" {
		return []ScrutinyRecord{}, nil
	}
	query := `
		SELECT id, student_id, class_id, semester, conduct_grade, final_decision, notes, coordinator_id,
		       COALESCE(status, 'draft'), validated_by::text, validated_at, created_at, updated_at 
		FROM scrutiny_records WHERE class_id::text = $1 AND semester = $2
	`
	rows, err := r.db.QueryContext(ctx, query, classID, semester)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []ScrutinyRecord
	for rows.Next() {
		var rec ScrutinyRecord
		var valBy sql.NullString
		var valAt sql.NullTime
		if err := rows.Scan(
			&rec.ID, &rec.StudentID, &rec.ClassID, &rec.Semester, &rec.ConductGrade, &rec.FinalDecision,
			&rec.Notes, &rec.CoordinatorID, &rec.Status, &valBy, &valAt, &rec.CreatedAt, &rec.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if valBy.Valid {
			rec.ValidatedBy = valBy.String
		}
		if valAt.Valid {
			rec.ValidatedAt = &valAt.Time
		}
		res = append(res, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *postgresRepository) ValidateClassScrutiny(ctx context.Context, classID string, semester int, validatorID string) error {
	query := `
		UPDATE scrutiny_records
		SET status = 'validated', validated_by = $1, validated_at = NOW(), updated_at = NOW()
		WHERE class_id = $2 AND semester = $3
	`
	_, err := r.db.ExecContext(ctx, query, validatorID, classID, semester)
	return err
}

func (r *postgresRepository) UpdateClassScrutinyStatus(ctx context.Context, classID string, semester int, status string) error {
	query := `
		UPDATE scrutiny_records
		SET status = $1, updated_at = NOW()
		WHERE class_id = $2 AND semester = $3
	`
	_, err := r.db.ExecContext(ctx, query, status, classID, semester)
	return err
}

func (r *postgresRepository) SaveDeficiency(ctx context.Context, def *StudentDeficiency) error {
	if def.ID == "" {
		def.ID = uuid.New().String()
	}
	def.UpdatedAt = time.Now()
	if def.PeriodType == "" {
		def.PeriodType = "semester_1"
		if def.Semester == 2 {
			def.PeriodType = "semester_2"
		}
	}
	if def.Status == "" {
		def.Status = "da_recuperare"
	}
	if def.RecoveryMode == "" {
		def.RecoveryMode = "studio_individuale"
	}

	query := `
		INSERT INTO student_deficiencies (
			id, school_id, student_id, class_id, subject_id, semester, period_type,
			topics, recovery_mode, status, recovery_grade, recovery_date, notes, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
		ON CONFLICT (id) DO UPDATE SET
			topics = EXCLUDED.topics,
			recovery_mode = EXCLUDED.recovery_mode,
			status = EXCLUDED.status,
			recovery_grade = EXCLUDED.recovery_grade,
			recovery_date = EXCLUDED.recovery_date,
			notes = EXCLUDED.notes,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.ExecContext(ctx, query,
		def.ID, def.SchoolID, def.StudentID, def.ClassID, def.SubjectID, def.Semester, def.PeriodType,
		def.Topics, def.RecoveryMode, def.Status, def.RecoveryGrade, def.RecoveryDate, def.Notes, def.UpdatedAt,
	)
	return err
}

func (r *postgresRepository) GetDeficienciesByStudent(ctx context.Context, studentID string) ([]StudentDeficiency, error) {
	if studentID == "" {
		return []StudentDeficiency{}, nil
	}
	query := `
		SELECT d.id, d.school_id, d.student_id, d.class_id, d.subject_id,
		       COALESCE(s.name, ''), d.semester, d.period_type, d.topics,
		       d.recovery_mode, d.status, d.recovery_grade, d.recovery_date,
		       COALESCE(d.notes, ''), d.created_at, d.updated_at
		FROM student_deficiencies d
		LEFT JOIN subjects s ON d.subject_id::text = s.id::text
		WHERE d.student_id::text = $1
		ORDER BY d.semester ASC, s.name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []StudentDeficiency
	for rows.Next() {
		var def StudentDeficiency
		var recGrade sql.NullFloat64
		var recDate sql.NullTime
		if err := rows.Scan(
			&def.ID, &def.SchoolID, &def.StudentID, &def.ClassID, &def.SubjectID,
			&def.SubjectName, &def.Semester, &def.PeriodType, &def.Topics,
			&def.RecoveryMode, &def.Status, &recGrade, &recDate,
			&def.Notes, &def.CreatedAt, &def.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if recGrade.Valid {
			def.RecoveryGrade = &recGrade.Float64
		}
		if recDate.Valid {
			def.RecoveryDate = &recDate.Time
		}
		res = append(res, def)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *postgresRepository) GetDeficienciesByClass(ctx context.Context, classID string, semester int) ([]StudentDeficiency, error) {
	if classID == "" {
		return []StudentDeficiency{}, nil
	}
	query := `
		SELECT d.id, d.school_id, d.student_id, COALESCE(u.first_name || ' ' || u.last_name, ''),
		       d.class_id, d.subject_id, COALESCE(s.name, ''), d.semester, d.period_type, d.topics,
		       d.recovery_mode, d.status, d.recovery_grade, d.recovery_date,
		       COALESCE(d.notes, ''), d.created_at, d.updated_at
		FROM student_deficiencies d
		LEFT JOIN users u ON d.student_id::text = u.id::text
		LEFT JOIN subjects s ON d.subject_id::text = s.id::text
		WHERE d.class_id::text = $1 AND ($2 = 0 OR d.semester = $2)
		ORDER BY u.last_name ASC, s.name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, classID, semester)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []StudentDeficiency
	for rows.Next() {
		var def StudentDeficiency
		var recGrade sql.NullFloat64
		var recDate sql.NullTime
		if err := rows.Scan(
			&def.ID, &def.SchoolID, &def.StudentID, &def.StudentName,
			&def.ClassID, &def.SubjectID, &def.SubjectName, &def.Semester, &def.PeriodType, &def.Topics,
			&def.RecoveryMode, &def.Status, &recGrade, &recDate,
			&def.Notes, &def.CreatedAt, &def.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if recGrade.Valid {
			def.RecoveryGrade = &recGrade.Float64
		}
		if recDate.Valid {
			def.RecoveryDate = &recDate.Time
		}
		res = append(res, def)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *postgresRepository) SaveDeferredScrutiny(ctx context.Context, req *SaveDeferredScrutinyRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// Update final decision in scrutiny_records for semester 2
	queryScrutiny := `
		UPDATE scrutiny_records
		SET final_decision = $1, notes = COALESCE(notes, '') || ' [Scrutinio Differito: ' || $1 || ']', updated_at = NOW()
		WHERE student_id::text = $2 AND class_id::text = $3 AND semester = 2
	`
	_, err = tx.ExecContext(ctx, queryScrutiny, req.FinalDecision, req.StudentID, req.ClassID)
	if err != nil {
		return err
	}

	// Update individual deficiency items
	queryDef := `
		UPDATE student_deficiencies
		SET status = $1, recovery_grade = $2, updated_at = NOW()
		WHERE id::text = $3
	`
	for _, item := range req.Deficiencies {
		_, err = tx.ExecContext(ctx, queryDef, item.Status, item.RecoveryGrade, item.DeficiencyID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
