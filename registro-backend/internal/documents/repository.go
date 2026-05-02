package documents

import (
	"database/sql"
)

type Repository interface {
	Create(doc *Document, initialContent string) error
	Update(doc *Document, newContent, changeLog string) error
	UpdateStatus(docID string, status DocStatus) error
	Delete(docID string) error
	ListAll(schoolID string, docType *DocType) ([]Document, error)

	FindByID(id string) (*Document, error)
	GetContent(docID string, version int) (string, error)
	GetVersions(docID string) ([]DocumentVersion, error)

	FindByClass(classID string) ([]Document, error)
	FindByStudent(studentID string) ([]Document, error)

	// Workflow Queues
	GetInbox(schoolID string) ([]Document, error)       // Submitted
	GetReviewQueue(schoolID string) ([]Document, error) // Review

	// Signatures
	AddSignature(sig *DocumentSignature) error
	GetSignatures(docID string) ([]DocumentSignature, error)

	// Templates
	GetTemplates(schoolID string) ([]DocumentTemplate, error)
	GetTemplate(id string) (*DocumentTemplate, error)
	CreateTemplate(tpl *DocumentTemplate) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(d *Document, content string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Insert Document Header
	query := `
		INSERT INTO documents_enhanced (
			school_id, title, type, student_id, class_id, status, 
			current_version, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, 1, $7, NOW(), NOW())
		RETURNING id`

	err = tx.QueryRow(query,
		d.SchoolID, d.Title, d.Type, d.StudentID, d.ClassID, d.Status, d.CreatedBy,
	).Scan(&d.ID)

	if err != nil {
		return err
	}

	// 2. Insert Version 1
	vQuery := `
		INSERT INTO document_versions (document_id, version_num, content, created_by, change_log)
		VALUES ($1, 1, $2, $3, 'Initial Creation')`

	_, err = tx.Exec(vQuery, d.ID, content, d.CreatedBy)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) Update(d *Document, newContent, changeLog string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Increment Version
	newVersion := d.CurrentVersion + 1

	// 1. Insert New Version
	vQuery := `
		INSERT INTO document_versions (document_id, version_num, content, created_by, change_log)
		VALUES ($1, $2, $3, $4, $5)`

	_, err = tx.Exec(vQuery, d.ID, newVersion, newContent, d.CreatedBy, changeLog) // CreatedBy here is the updater
	if err != nil {
		return err
	}

	// 2. Update Header
	uQuery := `
		UPDATE documents_enhanced SET 
			title=$1, current_version=$2, updated_at=NOW()
		WHERE id=$3`

	_, err = tx.Exec(uQuery, d.Title, newVersion, d.ID)
	if err != nil {
		return err
	}

	d.CurrentVersion = newVersion
	return tx.Commit()
}

func (r *repository) UpdateStatus(docID string, status DocStatus) error {
	_, err := r.db.Exec(`UPDATE documents_enhanced SET status=$1, updated_at=NOW() WHERE id=$2`, status, docID)
	return err
}

func (r *repository) Delete(docID string) error {
	_, err := r.db.Exec(`UPDATE documents_enhanced SET deleted_at=NOW() WHERE id=$1`, docID)
	return err
}

func (r *repository) ListAll(schoolID string, docType *DocType) ([]Document, error) {
	query := `SELECT id, school_id, title, type, student_id, class_id, status, current_version, is_signed, signed_by, signed_at, created_by, created_at, updated_at, deleted_at FROM documents_enhanced WHERE school_id = $1 AND deleted_at IS NULL`
	if docType != nil {
		query += ` AND type = '` + string(*docType) + `'`
	}
	query += ` ORDER BY updated_at DESC`
	return r.queryDocs(query, schoolID)
}

func (r *repository) FindByID(id string) (*Document, error) {
	query := `
		SELECT id, school_id, title, type, student_id, class_id, status, 
		       current_version, is_signed, signed_by, signed_at, created_by, created_at, updated_at
		FROM documents_enhanced WHERE id=$1 AND deleted_at IS NULL`

	var d Document
	err := r.db.QueryRow(query, id).Scan(
		&d.ID, &d.SchoolID, &d.Title, &d.Type, &d.StudentID, &d.ClassID, &d.Status,
		&d.CurrentVersion, &d.IsSigned, &d.SignedBy, &d.SignedAt,
		&d.CreatedBy, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *repository) GetContent(docID string, version int) (string, error) {
	var content string
	err := r.db.QueryRow(`SELECT content FROM document_versions WHERE document_id=$1 AND version_num=$2`, docID, version).Scan(&content)
	return content, err
}

func (r *repository) GetVersions(docID string) ([]DocumentVersion, error) {
	query := `SELECT id, version_num, created_by, created_at, change_log FROM document_versions WHERE document_id=$1 ORDER BY version_num DESC`
	rows, err := r.db.Query(query, docID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vers []DocumentVersion
	for rows.Next() {
		var v DocumentVersion
		if err := rows.Scan(&v.ID, &v.VersionNum, &v.CreatedBy, &v.CreatedAt, &v.ChangeLog); err != nil {
			return nil, err
		}
		vers = append(vers, v)
	}
	return vers, nil
}

func (r *repository) FindByClass(classID string) ([]Document, error) {
	return r.queryDocs(`SELECT * FROM documents_enhanced WHERE class_id=$1 AND deleted_at IS NULL`, classID)
}

func (r *repository) FindByStudent(studentID string) ([]Document, error) {
	return r.queryDocs(`SELECT * FROM documents_enhanced WHERE student_id=$1 AND deleted_at IS NULL`, studentID)
}

func (r *repository) GetInbox(schoolID string) ([]Document, error) {
	return r.queryDocs(`SELECT * FROM documents_enhanced WHERE school_id=$1 AND status='submitted' AND deleted_at IS NULL`, schoolID)
}

func (r *repository) GetReviewQueue(schoolID string) ([]Document, error) {
	return r.queryDocs(`SELECT * FROM documents_enhanced WHERE school_id=$1 AND status='review' AND deleted_at IS NULL`, schoolID)
}

func (r *repository) AddSignature(sig *DocumentSignature) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Insert Sig
	_, err = tx.Exec(`
		INSERT INTO document_signatures (document_id, version_id, signer_id, signature_data, certificate_data, signed_at)
		VALUES ($1, $2, $3, $4, $5, NOW())`,
		sig.DocumentID, sig.VersionID, sig.SignerID, sig.SignatureData, sig.Certificate)
	if err != nil {
		return err
	}

	// 2. Mark Doc as Signed
	_, err = tx.Exec(`
		UPDATE documents_enhanced SET is_signed=TRUE, signed_by=$1, signed_at=NOW(), status='signed'
		WHERE id=$2`, sig.SignerID, sig.DocumentID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) GetSignatures(docID string) ([]DocumentSignature, error) {
	rows, err := r.db.Query(`SELECT signer_id, signed_at FROM document_signatures WHERE document_id=$1`, docID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sigs []DocumentSignature
	for rows.Next() {
		var s DocumentSignature
		if err := rows.Scan(&s.SignerID, &s.SignedAt); err != nil {
			return nil, err
		}
		sigs = append(sigs, s)
	}
	return sigs, nil
}

// Templates
func (r *repository) GetTemplates(schoolID string) ([]DocumentTemplate, error) {
	rows, err := r.db.Query(`SELECT id, name, type, content FROM document_templates WHERE school_id=$1 AND is_active=TRUE`, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tpls []DocumentTemplate
	for rows.Next() {
		var t DocumentTemplate
		if err := rows.Scan(&t.ID, &t.Name, &t.Type, &t.Content); err != nil {
			return nil, err
		}
		tpls = append(tpls, t)
	}
	return tpls, nil
}

func (r *repository) GetTemplate(id string) (*DocumentTemplate, error) {
	var t DocumentTemplate
	err := r.db.QueryRow(`SELECT id, name, type, content FROM document_templates WHERE id=$1`, id).Scan(&t.ID, &t.Name, &t.Type, &t.Content)
	return &t, err
}

func (r *repository) CreateTemplate(t *DocumentTemplate) error {
	return r.db.QueryRow(`
		INSERT INTO document_templates (school_id, name, type, content) VALUES ($1, $2, $3, $4) RETURNING id`,
		t.SchoolID, t.Name, t.Type, t.Content).Scan(&t.ID)
}

// Helper
func (r *repository) queryDocs(query string, args ...interface{}) ([]Document, error) {
	// Note: SELECT * is dangerous if schema changes, explicit columns preferred in prod.
	// For brevity using explicit scan of known cols matching FindByID
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []Document
	for rows.Next() {
		var d Document
		if err := rows.Scan(
			&d.ID, &d.SchoolID, &d.Title, &d.Type, &d.StudentID, &d.ClassID, &d.Status,
			&d.CurrentVersion, &d.IsSigned, &d.SignedBy, &d.SignedAt,
			&d.CreatedBy, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt, // Scan DeletedAt since SELECT * includes it
		); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, nil
}
