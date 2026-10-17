package enrollment

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repository interface {
	SaveApplications(ctx context.Context, apps []EnrollmentApplication, schoolID, academicYear string) (int, error)
	ListApplications(ctx context.Context, schoolID, academicYear, status string) ([]EnrollmentApplication, error)
	SaveDraft(ctx context.Context, draft *ClassFormationDraft) error
	GetDraft(ctx context.Context, id string) (*ClassFormationDraft, error)
	ListDrafts(ctx context.Context, schoolID string) ([]ClassFormationDraft, error)
	UpdateDraftAssignments(ctx context.Context, id string, assignments ClassFormationDraftResult) error
	FinalizeDraft(ctx context.Context, id string) error
}

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) SaveApplications(ctx context.Context, apps []EnrollmentApplication, schoolID, academicYear string) (int, error) {
	if len(apps) == 0 {
		return 0, nil
	}
	saved := 0
	query := `
		INSERT INTO enrollment_applications (
			school_id, academic_year, sidi_application_id, student_first_name, student_last_name,
			student_tax_code, birth_date, gender, origin_school, middle_school_grade,
			track_chosen, second_language, has_disability_l104, has_dsa, religion_choice,
			requested_classmates, incompatible_classmates, parent1_first_name, parent1_last_name,
			parent1_email, parent1_phone, status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7::date, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22
		)
		ON CONFLICT (sidi_application_id) DO UPDATE SET
			student_first_name = EXCLUDED.student_first_name,
			student_last_name = EXCLUDED.student_last_name,
			middle_school_grade = EXCLUDED.middle_school_grade,
			has_disability_l104 = EXCLUDED.has_disability_l104,
			has_dsa = EXCLUDED.has_dsa,
			second_language = EXCLUDED.second_language
	`
	for _, a := range apps {
		_, err := r.db.ExecContext(ctx, query,
			schoolID, academicYear, a.SidiApplicationID, a.StudentFirstName, a.StudentLastName,
			a.StudentTaxCode, a.BirthDate, a.Gender, a.OriginSchool, a.MiddleSchoolGrade,
			a.TrackChosen, a.SecondLanguage, a.HasDisabilityL104, a.HasDSA, a.ReligionChoice,
			pq.Array(a.RequestedClassmates), pq.Array(a.IncompatibleClassmates),
			a.Parent1FirstName, a.Parent1LastName, a.Parent1Email, a.Parent1Phone, a.Status,
		)
		if err == nil {
			saved++
		}
	}
	return saved, nil
}

func (r *postgresRepository) ListApplications(ctx context.Context, schoolID, academicYear, status string) ([]EnrollmentApplication, error) {
	query := `
		SELECT 
			id, school_id, academic_year, sidi_application_id, student_first_name, student_last_name,
			student_tax_code, to_char(birth_date, 'YYYY-MM-DD'), gender, COALESCE(origin_school, ''),
			COALESCE(middle_school_grade, 7), track_chosen, COALESCE(second_language, ''),
			has_disability_l104, has_dsa, COALESCE(religion_choice, 'irc'),
			COALESCE(requested_classmates, '{}'), COALESCE(incompatible_classmates, '{}'),
			parent1_first_name, parent1_last_name, parent1_email, COALESCE(parent1_phone, ''),
			status, COALESCE(assigned_class_id::text, ''), created_at
		FROM enrollment_applications
		WHERE school_id = $1 AND academic_year = $2
		ORDER BY student_last_name ASC, student_first_name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, academicYear)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []EnrollmentApplication
	for rows.Next() {
		var a EnrollmentApplication
		var assignedClass string
		var reqFriends, incompFriends []string
		if err := rows.Scan(
			&a.ID, &a.SchoolID, &a.AcademicYear, &a.SidiApplicationID, &a.StudentFirstName, &a.StudentLastName,
			&a.StudentTaxCode, &a.BirthDate, &a.Gender, &a.OriginSchool,
			&a.MiddleSchoolGrade, &a.TrackChosen, &a.SecondLanguage,
			&a.HasDisabilityL104, &a.HasDSA, &a.ReligionChoice,
			pq.Array(&reqFriends), pq.Array(&incompFriends),
			&a.Parent1FirstName, &a.Parent1LastName, &a.Parent1Email, &a.Parent1Phone,
			&a.Status, &assignedClass, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		a.RequestedClassmates = reqFriends
		a.IncompatibleClassmates = incompFriends
		a.AssignedClassID = assignedClass
		results = append(results, a)
	}
	return results, rows.Err()
}

func (r *postgresRepository) SaveDraft(ctx context.Context, draft *ClassFormationDraft) error {
	if draft.ID == "" {
		draft.ID = uuid.New().String()
	}
	paramsJSON, _ := json.Marshal(draft.Parameters)
	assignJSON, _ := json.Marshal(draft.Assignments)

	query := `
		INSERT INTO class_formation_drafts (id, school_id, academic_year, title, parameters, assignments, is_finalized, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query, draft.ID, draft.SchoolID, draft.AcademicYear, draft.Title, paramsJSON, assignJSON, draft.IsFinalized, time.Now())
	return err
}

func (r *postgresRepository) GetDraft(ctx context.Context, id string) (*ClassFormationDraft, error) {
	query := `
		SELECT id, school_id, academic_year, title, parameters, assignments, is_finalized, created_at
		FROM class_formation_drafts
		WHERE id = $1
	`
	var row RawDraftRow
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&row.ID, &row.SchoolID, &row.AcademicYear, &row.Title,
		&row.Parameters, &row.Assignments, &row.IsFinalized, &row.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	var d ClassFormationDraft
	d.ID = row.ID
	d.SchoolID = row.SchoolID
	d.AcademicYear = row.AcademicYear
	d.Title = row.Title
	d.IsFinalized = row.IsFinalized
	d.CreatedAt = row.CreatedAt
	_ = json.Unmarshal(row.Parameters, &d.Parameters)
	_ = json.Unmarshal(row.Assignments, &d.Assignments)

	return &d, nil
}

func (r *postgresRepository) ListDrafts(ctx context.Context, schoolID string) ([]ClassFormationDraft, error) {
	query := `
		SELECT id, school_id, academic_year, title, parameters, assignments, is_finalized, created_at
		FROM class_formation_drafts
		WHERE school_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []ClassFormationDraft
	for rows.Next() {
		var row RawDraftRow
		if err := rows.Scan(
			&row.ID, &row.SchoolID, &row.AcademicYear, &row.Title,
			&row.Parameters, &row.Assignments, &row.IsFinalized, &row.CreatedAt,
		); err != nil {
			return nil, err
		}
		var d ClassFormationDraft
		d.ID = row.ID
		d.SchoolID = row.SchoolID
		d.AcademicYear = row.AcademicYear
		d.Title = row.Title
		d.IsFinalized = row.IsFinalized
		d.CreatedAt = row.CreatedAt
		_ = json.Unmarshal(row.Parameters, &d.Parameters)
		_ = json.Unmarshal(row.Assignments, &d.Assignments)
		results = append(results, d)
	}
	return results, rows.Err()
}

func (r *postgresRepository) UpdateDraftAssignments(ctx context.Context, id string, assignments ClassFormationDraftResult) error {
	assignJSON, err := json.Marshal(assignments)
	if err != nil {
		return err
	}
	query := `UPDATE class_formation_drafts SET assignments = $1 WHERE id = $2`
	_, err = r.db.ExecContext(ctx, query, assignJSON, id)
	return err
}

func (r *postgresRepository) FinalizeDraft(ctx context.Context, id string) error {
	draft, err := r.GetDraft(ctx, id)
	if err != nil {
		return err
	}
	if draft.IsFinalized {
		return fmt.Errorf("draft %s already finalized", id)
	}

	// Update draft finalized status
	_, err = r.db.ExecContext(ctx, "UPDATE class_formation_drafts SET is_finalized = true WHERE id = $1", id)
	return err
}
