package pdp

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"registro-backend/pkg/logger"
)

// Repository defines the data access interface for PDP plans.
type Repository interface {
	Create(ctx context.Context, plan *PdpPlan) (*PdpPlan, error)
	GetByID(ctx context.Context, id string) (*PdpPlan, error)
	GetByStudent(ctx context.Context, studentID, academicYear string) ([]*PdpPlan, error)
	GetByClass(ctx context.Context, classID, academicYear string) ([]*PdpPlan, error)
	Update(ctx context.Context, id string, req *UpdatePdpRequest) (*PdpPlan, error)
	SetSharedWithFamily(ctx context.Context, id string, shared bool) error
	ApproveByFamily(ctx context.Context, id, approverID string) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *sql.DB
}

// NewRepository creates a new PDP repository.
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// scanPlan scans a row into a PdpPlan, deserializing the JSONB content field.
func scanPlan(row interface {
	Scan(...any) error
}) (*PdpPlan, error) {
	plan := &PdpPlan{}
	var contentRaw []byte
	var coordinatorID, referenteID, familyApprovedBy sql.NullString
	var createdBy sql.NullString

	err := row.Scan(
		&plan.ID,
		&plan.StudentID,
		&plan.ClassID,
		&plan.SchoolID,
		&plan.AcademicYear,
		&plan.PlanType,
		&plan.Diagnosis,
		&contentRaw,
		&coordinatorID,
		&referenteID,
		&plan.SharedWithFamily,
		&plan.FamilyApprovedAt,
		&familyApprovedBy,
		&createdBy,
		&plan.CreatedAt,
		&plan.UpdatedAt,
		// Joined
		&plan.StudentName,
	)
	if err != nil {
		return nil, err
	}

	if coordinatorID.Valid {
		plan.CoordinatorID = &coordinatorID.String
	}
	if referenteID.Valid {
		plan.ReferenteID = &referenteID.String
	}
	if familyApprovedBy.Valid {
		plan.FamilyApprovedBy = &familyApprovedBy.String
	}
	if createdBy.Valid {
		plan.CreatedBy = createdBy.String
	}

	if err := json.Unmarshal(contentRaw, &plan.Content); err != nil {
		logger.Log.Warnf("pdp: failed to unmarshal content for plan %s: %v", plan.ID, err)
		plan.Content = PdpContent{}
	}

	return plan, nil
}

const selectPlanSQL = `
	SELECT p.id, p.student_id, p.class_id, p.school_id, p.academic_year,
	       p.plan_type, p.diagnosis, p.content,
	       p.coordinator_id, p.referente_id,
	       p.shared_with_family, p.family_approved_at, p.family_approved_by,
	       p.created_by, p.created_at, p.updated_at,
	       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name
	FROM pdp_plans p
	LEFT JOIN users u ON u.id = p.student_id
`

func (r *repository) Create(ctx context.Context, plan *PdpPlan) (*PdpPlan, error) {
	contentJSON, err := json.Marshal(plan.Content)
	if err != nil {
		return nil, fmt.Errorf("pdp.Create: marshal content: %w", err)
	}

	query := `
		INSERT INTO pdp_plans
		  (student_id, class_id, school_id, academic_year, plan_type, diagnosis,
		   content, coordinator_id, referente_id, created_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id, created_at, updated_at
	`
	row := r.db.QueryRowContext(ctx, query,
		plan.StudentID, plan.ClassID, plan.SchoolID, plan.AcademicYear,
		plan.PlanType, plan.Diagnosis, contentJSON,
		plan.CoordinatorID, plan.ReferenteID, plan.CreatedBy,
	)
	if err := row.Scan(&plan.ID, &plan.CreatedAt, &plan.UpdatedAt); err != nil {
		return nil, fmt.Errorf("pdp.Create: scan returning: %w", err)
	}
	return plan, nil
}

func (r *repository) GetByID(ctx context.Context, id string) (*PdpPlan, error) {
	q := selectPlanSQL + ` WHERE p.id = $1`
	row := r.db.QueryRowContext(ctx, q, id)
	plan, err := scanPlan(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return plan, err
}

func (r *repository) GetByStudent(ctx context.Context, studentID, academicYear string) ([]*PdpPlan, error) {
	q := selectPlanSQL + ` WHERE p.student_id = $1 AND p.academic_year = $2 ORDER BY p.created_at DESC`
	rows, err := r.db.QueryContext(ctx, q, studentID, academicYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (r *repository) GetByClass(ctx context.Context, classID, academicYear string) ([]*PdpPlan, error) {
	q := selectPlanSQL + ` WHERE p.class_id = $1 AND p.academic_year = $2 ORDER BY student_name`
	rows, err := r.db.QueryContext(ctx, q, classID, academicYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func scanRows(rows *sql.Rows) ([]*PdpPlan, error) {
	var plans []*PdpPlan
	for rows.Next() {
		p, err := scanPlan(rows)
		if err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, rows.Err()
}

func (r *repository) Update(ctx context.Context, id string, req *UpdatePdpRequest) (*PdpPlan, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Build dynamic UPDATE
	args := []any{}
	set := ""
	i := 1

	if req.PlanType != nil {
		set += fmt.Sprintf("plan_type = $%d, ", i)
		args = append(args, *req.PlanType)
		i++
	}
	if req.Diagnosis != nil {
		set += fmt.Sprintf("diagnosis = $%d, ", i)
		args = append(args, *req.Diagnosis)
		i++
	}
	if req.Content != nil {
		b, _ := json.Marshal(req.Content)
		set += fmt.Sprintf("content = $%d, ", i)
		args = append(args, b)
		i++
	}
	if req.CoordinatorID != nil {
		set += fmt.Sprintf("coordinator_id = $%d, ", i)
		args = append(args, *req.CoordinatorID)
		i++
	}
	if req.ReferenteID != nil {
		set += fmt.Sprintf("referente_id = $%d, ", i)
		args = append(args, *req.ReferenteID)
		i++
	}

	if set == "" {
		return r.GetByID(ctx, id)
	}

	// Remove trailing comma
	set = set[:len(set)-2]
	args = append(args, id)
	query := fmt.Sprintf("UPDATE pdp_plans SET %s WHERE id = $%d", set, i)
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return nil, fmt.Errorf("pdp.Update: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *repository) SetSharedWithFamily(ctx context.Context, id string, shared bool) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE pdp_plans SET shared_with_family = $1 WHERE id = $2`,
		shared, id,
	)
	return err
}

func (r *repository) ApproveByFamily(ctx context.Context, id, approverID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE pdp_plans SET family_approved_at = NOW(), family_approved_by = $1 WHERE id = $2`,
		approverID, id,
	)
	return err
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM pdp_plans WHERE id = $1`, id)
	return err
}
