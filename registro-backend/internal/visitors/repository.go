package visitors

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	// Visitors
	CreateVisitor(ctx context.Context, v *Visitor) error
	RecordVisitorExit(ctx context.Context, id, schoolID string, notes string) error
	ListTodayVisitors(ctx context.Context, schoolID, date string) ([]*Visitor, error)
	GetVisitor(ctx context.Context, id, schoolID string) (*Visitor, error)

	// Early exits
	CreateEarlyExit(ctx context.Context, e *EarlyExit) error
	RecordStudentReturn(ctx context.Context, id, schoolID, notes string) error
	ListTodayEarlyExits(ctx context.Context, schoolID, date string) ([]*EarlyExit, error)

	// Maintenance
	CreateMaintenanceReport(ctx context.Context, r *MaintenanceReport) error
	ListMaintenanceReports(ctx context.Context, schoolID string, statusFilter string) ([]*MaintenanceReport, error)
	UpdateMaintenanceStatus(ctx context.Context, id, schoolID string, req UpdateMaintenanceStatusRequest) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// --- Visitors ---

func (r *repository) CreateVisitor(ctx context.Context, v *Visitor) error {
	v.ID = uuid.New().String()
	v.EntryTime = time.Now()
	v.CreatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO visitors (id, school_id, name, document_id, purpose, host_name, badge_number, entry_time, notes, recorded_by, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		v.ID, v.SchoolID, v.Name, v.DocumentID, v.Purpose, v.HostName, v.BadgeNumber,
		v.EntryTime, v.Notes, v.RecordedBy, v.CreatedAt)
	return err
}

func (r *repository) RecordVisitorExit(ctx context.Context, id, schoolID, notes string) error {
	now := time.Now()
	res, err := r.db.ExecContext(ctx, `
		UPDATE visitors SET exit_time=$1, notes=CASE WHEN $2::text != '' THEN $2 ELSE notes END
		WHERE id=$3 AND school_id=$4 AND exit_time IS NULL`,
		now, notes, id, schoolID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("visitatore non trovato o già uscito")
	}
	return nil
}

func (r *repository) ListTodayVisitors(ctx context.Context, schoolID, date string) ([]*Visitor, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, school_id, name, COALESCE(document_id, ''), purpose, COALESCE(host_name, ''),
		       badge_number, entry_time, exit_time, COALESCE(notes, ''), recorded_by, created_at
		FROM visitors
		WHERE school_id=$1 AND DATE(entry_time)=$2
		ORDER BY entry_time DESC`, schoolID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var visitors []*Visitor
	for rows.Next() {
		v := &Visitor{}
		var badgeNum sql.NullString
		var exitTime sql.NullTime
		if err := rows.Scan(
			&v.ID, &v.SchoolID, &v.Name, &v.DocumentID, &v.Purpose, &v.HostName,
			&badgeNum, &v.EntryTime, &exitTime, &v.Notes, &v.RecordedBy, &v.CreatedAt,
		); err != nil {
			return nil, err
		}
		if badgeNum.Valid {
			v.BadgeNumber = &badgeNum.String
		}
		if exitTime.Valid {
			v.ExitTime = &exitTime.Time
		}
		visitors = append(visitors, v)
	}
	return visitors, nil
}

func (r *repository) GetVisitor(ctx context.Context, id, schoolID string) (*Visitor, error) {
	v := &Visitor{}
	var badgeNum sql.NullString
	var exitTime sql.NullTime
	row := r.db.QueryRowContext(ctx, `
		SELECT id, school_id, name, COALESCE(document_id, ''), purpose, COALESCE(host_name, ''),
		       badge_number, entry_time, exit_time, COALESCE(notes, ''), recorded_by, created_at
		FROM visitors WHERE id=$1 AND school_id=$2`, id, schoolID)
	err := row.Scan(
		&v.ID, &v.SchoolID, &v.Name, &v.DocumentID, &v.Purpose, &v.HostName,
		&badgeNum, &v.EntryTime, &exitTime, &v.Notes, &v.RecordedBy, &v.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("visitatore non trovato")
	}
	if err != nil {
		return nil, err
	}
	if badgeNum.Valid {
		v.BadgeNumber = &badgeNum.String
	}
	if exitTime.Valid {
		v.ExitTime = &exitTime.Time
	}
	return v, nil
}

// --- Early exits ---

func (r *repository) CreateEarlyExit(ctx context.Context, e *EarlyExit) error {
	e.ID = uuid.New().String()
	e.ExitTime = time.Now()
	e.CreatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO student_early_exits (id, school_id, student_id, exit_time, delegatee_name, delegate_rel, reason_code, notes, recorded_by, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		e.ID, e.SchoolID, e.StudentID, e.ExitTime, e.DelegateeName, e.DelegateRel,
		e.ReasonCode, e.Notes, e.RecordedBy, e.CreatedAt)
	return err
}

func (r *repository) RecordStudentReturn(ctx context.Context, id, schoolID, notes string) error {
	now := time.Now()
	res, err := r.db.ExecContext(ctx, `
		UPDATE student_early_exits SET return_time=$1
		WHERE id=$2 AND school_id=$3 AND return_time IS NULL`,
		now, id, schoolID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("uscita non trovata o studente già rientrato")
	}
	return nil
}

func (r *repository) ListTodayEarlyExits(ctx context.Context, schoolID, date string) ([]*EarlyExit, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT see.id, see.school_id, see.student_id,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name,
		       COALESCE(c.name, '') AS class_name,
		       see.exit_time, see.return_time,
		       see.delegatee_name, COALESCE(see.delegate_rel, ''),
		       COALESCE(see.reason_code, ''), COALESCE(see.notes, ''),
		       see.recorded_by, see.created_at
		FROM student_early_exits see
		JOIN users u ON u.id = see.student_id
		LEFT JOIN classes c ON c.id = (
		    SELECT class_id FROM student_classes
		    WHERE student_id = see.student_id AND school_id = $1
		    LIMIT 1
		)
		WHERE see.school_id=$1 AND DATE(see.exit_time)=$2
		ORDER BY see.exit_time DESC`, schoolID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exits []*EarlyExit
	for rows.Next() {
		e := &EarlyExit{}
		var returnTime sql.NullTime
		if err := rows.Scan(
			&e.ID, &e.SchoolID, &e.StudentID,
			&e.StudentName, &e.ClassName,
			&e.ExitTime, &returnTime,
			&e.DelegateeName, &e.DelegateRel,
			&e.ReasonCode, &e.Notes,
			&e.RecordedBy, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		if returnTime.Valid {
			e.ReturnTime = &returnTime.Time
		}
		exits = append(exits, e)
	}
	return exits, nil
}

// --- Maintenance ---

func (r *repository) CreateMaintenanceReport(ctx context.Context, rep *MaintenanceReport) error {
	rep.ID = uuid.New().String()
	rep.Status = MaintenanceOpen
	rep.CreatedAt = time.Now()
	rep.UpdatedAt = time.Now()
	if rep.Priority == "" {
		rep.Priority = PriorityMedium
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO maintenance_reports (id, school_id, location, category, description, priority, status, reported_by, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		rep.ID, rep.SchoolID, rep.Location, rep.Category, rep.Description,
		rep.Priority, rep.Status, rep.ReportedBy, rep.CreatedAt, rep.UpdatedAt)
	return err
}

func (r *repository) ListMaintenanceReports(ctx context.Context, schoolID, statusFilter string) ([]*MaintenanceReport, error) {
	query := `SELECT id, school_id, location, category, description, priority, status,
	                 reported_by, assigned_to, closed_at, created_at, updated_at
	          FROM maintenance_reports WHERE school_id=$1`
	args := []interface{}{schoolID}
	if statusFilter != "" {
		query += ` AND status=$2`
		args = append(args, statusFilter)
	}
	query += ` ORDER BY priority DESC, created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*MaintenanceReport
	for rows.Next() {
		rep := &MaintenanceReport{}
		var assignedTo sql.NullString
		var closedAt sql.NullTime
		if err := rows.Scan(
			&rep.ID, &rep.SchoolID, &rep.Location, &rep.Category, &rep.Description,
			&rep.Priority, &rep.Status, &rep.ReportedBy, &assignedTo, &closedAt,
			&rep.CreatedAt, &rep.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if assignedTo.Valid {
			rep.AssignedTo = &assignedTo.String
		}
		if closedAt.Valid {
			rep.ClosedAt = &closedAt.Time
		}
		reports = append(reports, rep)
	}
	return reports, nil
}

func (r *repository) UpdateMaintenanceStatus(ctx context.Context, id, schoolID string, req UpdateMaintenanceStatusRequest) error {
	now := time.Now()
	var closedAt *time.Time
	if req.Status == MaintenanceClosed {
		closedAt = &now
	}
	res, err := r.db.ExecContext(ctx, `
		UPDATE maintenance_reports
		SET status=$1, assigned_to=COALESCE(NULLIF($2,''), assigned_to), closed_at=$3, updated_at=$4
		WHERE id=$5 AND school_id=$6`,
		req.Status, req.AssignedTo, closedAt, now, id, schoolID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("segnalazione non trovata")
	}
	return nil
}
