package uda

import (
	"context"
	"database/sql"
	"encoding/json"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, plan *UdaPlan) error {
	compJSON, _ := json.Marshal(plan.Competencies)
	query := `
		INSERT INTO uda_plans (school_id, class_id, subject_id, teacher_id, title, description, period, start_date, end_date, competencies, objectives, methodologies, evaluation_criteria, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW(), NOW())
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		plan.SchoolID, plan.ClassID, plan.SubjectID, plan.TeacherID,
		plan.Title, plan.Description, plan.Period, plan.StartDate, plan.EndDate,
		compJSON, plan.Objectives, plan.Methodologies, plan.EvaluationCriteria, plan.Status,
	).Scan(&plan.ID, &plan.CreatedAt, &plan.UpdatedAt)
}

func (r *Repository) ListByClass(ctx context.Context, classID string) ([]*UdaPlan, error) {
	if classID == "" {
		return r.ListAll(ctx)
	}
	query := `
		SELECT u.id, u.school_id, u.class_id, COALESCE(c.section, '') AS class_name,
		       u.subject_id, COALESCE(s.name, '') AS subject_name,
		       u.teacher_id, COALESCE(usr.first_name || ' ' || usr.last_name, '') AS teacher_name,
		       u.title, COALESCE(u.description, ''), COALESCE(u.period, 'annuale'),
		       u.start_date, u.end_date, COALESCE(u.competencies, '[]'::jsonb),
		       COALESCE(u.objectives, ''), COALESCE(u.methodologies, ''), COALESCE(u.evaluation_criteria, ''),
		       COALESCE(u.status, 'draft'), u.created_at, u.updated_at
		FROM uda_plans u
		LEFT JOIN classes c ON u.class_id = c.id
		LEFT JOIN subjects s ON u.subject_id = s.id
		LEFT JOIN users usr ON u.teacher_id = usr.id
		WHERE u.class_id = $1::uuid
		ORDER BY u.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*UdaPlan
	for rows.Next() {
		p := &UdaPlan{}
		var compBytes []byte
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.ClassID, &p.ClassName,
			&p.SubjectID, &p.SubjectName,
			&p.TeacherID, &p.TeacherName,
			&p.Title, &p.Description, &p.Period,
			&p.StartDate, &p.EndDate, &compBytes,
			&p.Objectives, &p.Methodologies, &p.EvaluationCriteria,
			&p.Status, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(compBytes, &p.Competencies)
		if p.Competencies == nil {
			p.Competencies = []string{}
		}
		plans = append(plans, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if plans == nil {
		plans = []*UdaPlan{}
	}
	return plans, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*UdaPlan, error) {
	query := `
		SELECT u.id, u.school_id, u.class_id, COALESCE(c.section, '') AS class_name,
		       u.subject_id, COALESCE(s.name, '') AS subject_name,
		       u.teacher_id, COALESCE(usr.first_name || ' ' || usr.last_name, '') AS teacher_name,
		       u.title, COALESCE(u.description, ''), COALESCE(u.period, 'annuale'),
		       u.start_date, u.end_date, COALESCE(u.competencies, '[]'::jsonb),
		       COALESCE(u.objectives, ''), COALESCE(u.methodologies, ''), COALESCE(u.evaluation_criteria, ''),
		       COALESCE(u.status, 'draft'), u.created_at, u.updated_at
		FROM uda_plans u
		LEFT JOIN classes c ON u.class_id = c.id
		LEFT JOIN subjects s ON u.subject_id = s.id
		LEFT JOIN users usr ON u.teacher_id = usr.id
		WHERE u.id = $1::uuid`

	p := &UdaPlan{}
	var compBytes []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.SchoolID, &p.ClassID, &p.ClassName,
		&p.SubjectID, &p.SubjectName,
		&p.TeacherID, &p.TeacherName,
		&p.Title, &p.Description, &p.Period,
		&p.StartDate, &p.EndDate, &compBytes,
		&p.Objectives, &p.Methodologies, &p.EvaluationCriteria,
		&p.Status, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(compBytes, &p.Competencies)
	if p.Competencies == nil {
		p.Competencies = []string{}
	}
	return p, nil
}

func (r *Repository) ListBySchool(ctx context.Context, schoolID string) ([]*UdaPlan, error) {
	if schoolID == "" {
		return r.ListAll(ctx)
	}
	query := `
		SELECT u.id, u.school_id, u.class_id, COALESCE(c.section, '') AS class_name,
		       u.subject_id, COALESCE(s.name, '') AS subject_name,
		       u.teacher_id, COALESCE(usr.first_name || ' ' || usr.last_name, '') AS teacher_name,
		       u.title, COALESCE(u.description, ''), COALESCE(u.period, 'annuale'),
		       u.start_date, u.end_date, COALESCE(u.competencies, '[]'::jsonb),
		       COALESCE(u.objectives, ''), COALESCE(u.methodologies, ''), COALESCE(u.evaluation_criteria, ''),
		       COALESCE(u.status, 'draft'), u.created_at, u.updated_at
		FROM uda_plans u
		LEFT JOIN classes c ON u.class_id = c.id
		LEFT JOIN subjects s ON u.subject_id = s.id
		LEFT JOIN users usr ON u.teacher_id = usr.id
		WHERE u.school_id = $1::uuid
		ORDER BY u.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*UdaPlan
	for rows.Next() {
		p := &UdaPlan{}
		var compBytes []byte
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.ClassID, &p.ClassName,
			&p.SubjectID, &p.SubjectName,
			&p.TeacherID, &p.TeacherName,
			&p.Title, &p.Description, &p.Period,
			&p.StartDate, &p.EndDate, &compBytes,
			&p.Objectives, &p.Methodologies, &p.EvaluationCriteria,
			&p.Status, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(compBytes, &p.Competencies)
		if p.Competencies == nil {
			p.Competencies = []string{}
		}
		plans = append(plans, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if plans == nil {
		plans = []*UdaPlan{}
	}
	return plans, nil
}

func (r *Repository) ListAll(ctx context.Context) ([]*UdaPlan, error) {
	query := `
		SELECT u.id, u.school_id, u.class_id, COALESCE(c.section, '') AS class_name,
		       u.subject_id, COALESCE(s.name, '') AS subject_name,
		       u.teacher_id, COALESCE(usr.first_name || ' ' || usr.last_name, '') AS teacher_name,
		       u.title, COALESCE(u.description, ''), COALESCE(u.period, 'annuale'),
		       u.start_date, u.end_date, COALESCE(u.competencies, '[]'::jsonb),
		       COALESCE(u.objectives, ''), COALESCE(u.methodologies, ''), COALESCE(u.evaluation_criteria, ''),
		       COALESCE(u.status, 'draft'), u.created_at, u.updated_at
		FROM uda_plans u
		LEFT JOIN classes c ON u.class_id = c.id
		LEFT JOIN subjects s ON u.subject_id = s.id
		LEFT JOIN users usr ON u.teacher_id = usr.id
		ORDER BY u.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*UdaPlan
	for rows.Next() {
		p := &UdaPlan{}
		var compBytes []byte
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.ClassID, &p.ClassName,
			&p.SubjectID, &p.SubjectName,
			&p.TeacherID, &p.TeacherName,
			&p.Title, &p.Description, &p.Period,
			&p.StartDate, &p.EndDate, &compBytes,
			&p.Objectives, &p.Methodologies, &p.EvaluationCriteria,
			&p.Status, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(compBytes, &p.Competencies)
		if p.Competencies == nil {
			p.Competencies = []string{}
		}
		plans = append(plans, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if plans == nil {
		plans = []*UdaPlan{}
	}
	return plans, nil
}


func (r *Repository) Update(ctx context.Context, id string, req UpdateUdaRequest) (*UdaPlan, error) {
	compJSON, _ := json.Marshal(req.Competencies)
	query := `
		UPDATE uda_plans
		SET title = COALESCE(NULLIF($2, ''), title),
		    description = COALESCE($3, description),
		    period = COALESCE(NULLIF($4, ''), period),
		    competencies = COALESCE($5, competencies),
		    objectives = COALESCE($6, objectives),
		    methodologies = COALESCE($7, methodologies),
		    evaluation_criteria = COALESCE($8, evaluation_criteria),
		    status = COALESCE(NULLIF($9, ''), status),
		    updated_at = NOW()
		WHERE id = $1::uuid
		RETURNING id, school_id, class_id, subject_id, teacher_id, title, status`

	p := &UdaPlan{}
	err := r.db.QueryRowContext(ctx, query, id, req.Title, req.Description, req.Period, compJSON, req.Objectives, req.Methodologies, req.EvaluationCriteria, req.Status).Scan(
		&p.ID, &p.SchoolID, &p.ClassID, &p.SubjectID, &p.TeacherID, &p.Title, &p.Status,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM uda_plans WHERE id = $1::uuid", id)
	return err
}
