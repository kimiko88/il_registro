package security_passes

import (
	"database/sql"
	"time"
)

type Repository interface {
	CreateHallPass(p *StudentHallPass) error
	CloseHallPass(id string, returnedAt time.Time) error
	ListActiveHallPassesBySchool(schoolID string) ([]StudentHallPass, error)
	ListHallPassesByClass(classID string, limit int) ([]StudentHallPass, error)

	CreatePickupDelegate(d *StudentPickupDelegate) error
	ListDelegatesByStudent(studentID string) ([]StudentPickupDelegate, error)
	RevokeDelegate(id string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateHallPass(p *StudentHallPass) error {
	if p.DepartedAt.IsZero() {
		p.DepartedAt = time.Now()
	}
	query := `
		INSERT INTO student_hall_passes (
			school_id, class_id, student_id, teacher_id, pass_type, departed_at, status, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query, p.SchoolID, p.ClassID, p.StudentID, p.TeacherID, p.PassType, p.DepartedAt, p.Status, p.Notes,
	).Scan(&p.ID, &p.CreatedAt)
}

func (r *repository) CloseHallPass(id string, returnedAt time.Time) error {
	query := `
		UPDATE student_hall_passes
		SET returned_at = $2,
		    duration_minutes = GREATEST(1, ROUND(EXTRACT(EPOCH FROM ($2 - departed_at)) / 60)),
		    status = 'rientrato'
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id, returnedAt)
	return err
}

func (r *repository) ListActiveHallPassesBySchool(schoolID string) ([]StudentHallPass, error) {
	query := `
		SELECT p.id, p.school_id, p.class_id, p.student_id, p.teacher_id, p.pass_type,
		       p.departed_at, p.returned_at, p.duration_minutes, p.status, COALESCE(p.notes, ''), p.created_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name,
		       COALESCE(c.name, '') AS class_name
		FROM student_hall_passes p
		LEFT JOIN users u ON p.student_id = u.id
		LEFT JOIN classes c ON p.class_id = c.id
		WHERE p.school_id = $1 AND p.status = 'fuori_aula'
		ORDER BY p.departed_at ASC
	`
	rows, err := r.db.Query(query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passes []StudentHallPass
	now := time.Now()
	for rows.Next() {
		var p StudentHallPass
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.ClassID, &p.StudentID, &p.TeacherID, &p.PassType,
			&p.DepartedAt, &p.ReturnedAt, &p.DurationMinutes, &p.Status, &p.Notes, &p.CreatedAt,
			&p.StudentFullName, &p.ClassName,
		); err != nil {
			return nil, err
		}
		p.ElapsedMinutes = int(now.Sub(p.DepartedAt).Minutes())
		if p.ElapsedMinutes > 15 {
			p.IsAnomalyAlert = true
		}
		passes = append(passes, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return passes, nil
}

func (r *repository) ListHallPassesByClass(classID string, limit int) ([]StudentHallPass, error) {
	if limit <= 0 {
		limit = 30
	}
	query := `
		SELECT p.id, p.school_id, p.class_id, p.student_id, p.teacher_id, p.pass_type,
		       p.departed_at, p.returned_at, p.duration_minutes, p.status, COALESCE(p.notes, ''), p.created_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name,
		       COALESCE(c.name, '') AS class_name
		FROM student_hall_passes p
		LEFT JOIN users u ON p.student_id = u.id
		LEFT JOIN classes c ON p.class_id = c.id
		WHERE p.class_id = $1
		ORDER BY p.departed_at DESC
		LIMIT $2
	`
	rows, err := r.db.Query(query, classID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passes []StudentHallPass
	now := time.Now()
	for rows.Next() {
		var p StudentHallPass
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.ClassID, &p.StudentID, &p.TeacherID, &p.PassType,
			&p.DepartedAt, &p.ReturnedAt, &p.DurationMinutes, &p.Status, &p.Notes, &p.CreatedAt,
			&p.StudentFullName, &p.ClassName,
		); err != nil {
			return nil, err
		}
		if p.Status == "fuori_aula" {
			p.ElapsedMinutes = int(now.Sub(p.DepartedAt).Minutes())
			if p.ElapsedMinutes > 15 {
				p.IsAnomalyAlert = true
			}
		}
		passes = append(passes, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return passes, nil
}

func (r *repository) CreatePickupDelegate(d *StudentPickupDelegate) error {
	query := `
		INSERT INTO student_pickup_delegates (
			school_id, student_id, delegate_name, delegate_relation, fiscal_code,
			phone, id_card_number, id_card_file_url, is_permanent, valid_until, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query, d.SchoolID, d.StudentID, d.DelegateName, d.DelegateRelation, d.FiscalCode,
		d.Phone, d.IDCardNumber, d.IDCardFileURL, d.IsPermanent, d.ValidUntil, d.Status,
	).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
}

func (r *repository) ListDelegatesByStudent(studentID string) ([]StudentPickupDelegate, error) {
	query := `
		SELECT id, school_id, student_id, delegate_name, delegate_relation, fiscal_code,
		       phone, id_card_number, COALESCE(id_card_file_url, ''), is_permanent,
		       TO_CHAR(valid_until, 'YYYY-MM-DD'), status, created_at, updated_at
		FROM student_pickup_delegates
		WHERE student_id = $1
		ORDER BY status ASC, delegate_name ASC
	`
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []StudentPickupDelegate
	for rows.Next() {
		var d StudentPickupDelegate
		if err := rows.Scan(
			&d.ID, &d.SchoolID, &d.StudentID, &d.DelegateName, &d.DelegateRelation, &d.FiscalCode,
			&d.Phone, &d.IDCardNumber, &d.IDCardFileURL, &d.IsPermanent, &d.ValidUntil,
			&d.Status, &d.CreatedAt, &d.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) RevokeDelegate(id string) error {
	query := `UPDATE student_pickup_delegates SET status = 'revocato', updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
