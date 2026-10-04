package religion_alternative

import (
	"context"
	"database/sql"
)

type Repository interface {
	SaveOption(ctx context.Context, opt *StudentReligionOption) error
	GetOption(ctx context.Context, studentID, academicYear string) (*StudentReligionOption, error)
	ListOptionsBySchool(ctx context.Context, schoolID, academicYear string) ([]StudentReligionOption, error)
	SaveEvaluation(ctx context.Context, eval *AlternativeEvaluation) error
	ListEvaluations(ctx context.Context, schoolID, period, subjectKind string) ([]AlternativeEvaluation, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) SaveOption(ctx context.Context, opt *StudentReligionOption) error {
	query := `
		INSERT INTO student_religion_options (student_id, school_id, academic_year, option_type, notes)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (student_id, academic_year)
		DO UPDATE SET option_type = EXCLUDED.option_type, notes = EXCLUDED.notes, chosen_at = CURRENT_TIMESTAMP
		RETURNING id, chosen_at;
	`
	return r.db.QueryRowContext(ctx, query, opt.StudentID, opt.SchoolID, opt.AcademicYear, opt.OptionType, opt.Notes).
		Scan(&opt.ID, &opt.ChosenAt)
}

func (r *postgresRepo) GetOption(ctx context.Context, studentID, academicYear string) (*StudentReligionOption, error) {
	query := `
		SELECT id, student_id, school_id, academic_year, option_type, COALESCE(notes, ''), chosen_at
		FROM student_religion_options
		WHERE student_id = $1 AND academic_year = $2
	`
	var opt StudentReligionOption
	err := r.db.QueryRowContext(ctx, query, studentID, academicYear).
		Scan(&opt.ID, &opt.StudentID, &opt.SchoolID, &opt.AcademicYear, &opt.OptionType, &opt.Notes, &opt.ChosenAt)
	if err != nil {
		return nil, err
	}
	return &opt, nil
}

func (r *postgresRepo) ListOptionsBySchool(ctx context.Context, schoolID, academicYear string) ([]StudentReligionOption, error) {
	query := `
		SELECT id, student_id, school_id, academic_year, option_type, COALESCE(notes, ''), chosen_at
		FROM student_religion_options
		WHERE school_id = $1 AND academic_year = $2
		ORDER BY student_id ASC;
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, academicYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []StudentReligionOption
	for rows.Next() {
		var o StudentReligionOption
		if err := rows.Scan(&o.ID, &o.StudentID, &o.SchoolID, &o.AcademicYear, &o.OptionType, &o.Notes, &o.ChosenAt); err != nil {
			return nil, err
		}
		list = append(list, o)
	}
	return list, nil
}

func (r *postgresRepo) SaveEvaluation(ctx context.Context, eval *AlternativeEvaluation) error {
	query := `
		INSERT INTO alternative_evaluations (
			student_id, school_id, group_id, class_id, period, subject_kind, judgment_level, descriptive_notes, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (student_id, period, subject_kind)
		DO UPDATE SET
			judgment_level = EXCLUDED.judgment_level,
			descriptive_notes = EXCLUDED.descriptive_notes,
			group_id = EXCLUDED.group_id,
			class_id = EXCLUDED.class_id,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at;
	`
	return r.db.QueryRowContext(ctx, query,
		eval.StudentID, eval.SchoolID, eval.GroupID, eval.ClassID, eval.Period, eval.SubjectKind, eval.JudgmentLevel, eval.DescriptiveNotes, eval.CreatedBy,
	).Scan(&eval.ID, &eval.CreatedAt, &eval.UpdatedAt)
}

func (r *postgresRepo) ListEvaluations(ctx context.Context, schoolID, period, subjectKind string) ([]AlternativeEvaluation, error) {
	query := `
		SELECT id, student_id, school_id, group_id, class_id, period, subject_kind, judgment_level, COALESCE(descriptive_notes, ''), created_by, created_at, updated_at
		FROM alternative_evaluations
		WHERE school_id = $1 AND ($2 = '' OR period = $2) AND ($3 = '' OR subject_kind = $3)
		ORDER BY student_id ASC;
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, period, subjectKind)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []AlternativeEvaluation
	for rows.Next() {
		var e AlternativeEvaluation
		if err := rows.Scan(&e.ID, &e.StudentID, &e.SchoolID, &e.GroupID, &e.ClassID, &e.Period, &e.SubjectKind, &e.JudgmentLevel, &e.DescriptiveNotes, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}
