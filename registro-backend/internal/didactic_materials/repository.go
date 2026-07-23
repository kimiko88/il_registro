package didactic_materials

import (
	"database/sql"
)

type Repository interface {
	Create(m *DidacticMaterial) error
	GetByClass(classID string) ([]DidacticMaterial, error)
	Delete(id string, teacherID string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(m *DidacticMaterial) error {
	query := `
		WITH inserted AS (
			INSERT INTO didactic_materials (school_id, class_id, subject_id, teacher_id, title, description, attachment_url, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
			RETURNING id, teacher_id
		)
		SELECT inserted.id, COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name
		FROM inserted
		LEFT JOIN users u ON inserted.teacher_id = u.id
	`
	return r.db.QueryRow(query,
		m.SchoolID, m.ClassID, m.SubjectID, m.TeacherID,
		m.Title, m.Description, m.AttachmentURL,
	).Scan(&m.ID, &m.TeacherName)
}

func (r *repository) GetByClass(classID string) ([]DidacticMaterial, error) {
	query := `
		SELECT dm.id, dm.school_id, dm.class_id, dm.subject_id, dm.teacher_id, 
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name, 
		       dm.title, dm.description, dm.attachment_url, dm.created_at, dm.updated_at
		FROM didactic_materials dm
		LEFT JOIN users u ON dm.teacher_id = u.id
		WHERE dm.class_id = $1
		ORDER BY dm.created_at DESC
	`
	rows, err := r.db.Query(query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []DidacticMaterial
	for rows.Next() {
		var m DidacticMaterial
		err := rows.Scan(
			&m.ID, &m.SchoolID, &m.ClassID, &m.SubjectID, &m.TeacherID,
			&m.TeacherName, &m.Title, &m.Description, &m.AttachmentURL,
			&m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

func (r *repository) Delete(id string, teacherID string) error {
	query := `DELETE FROM didactic_materials WHERE id = $1 AND teacher_id = $2`
	_, err := r.db.Exec(query, id, teacherID)
	return err
}
