package may15

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Repository interface {
	GetByClassAndYear(ctx context.Context, classID, academicYear string) (*ClassMay15Document, error)
	Upsert(ctx context.Context, doc *ClassMay15Document) error
	Publish(ctx context.Context, classID, academicYear string) (*ClassMay15Document, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetByClassAndYear(ctx context.Context, classID, academicYear string) (*ClassMay15Document, error) {
	query := `
		SELECT d.id, d.school_id, d.class_id, COALESCE(c.name, '') as class_name,
		       d.academic_year, d.status,
		       COALESCE(d.class_presentation, ''), COALESCE(d.teaching_continuity, ''),
		       COALESCE(d.pcto_pathways, ''), COALESCE(d.exam_simulations, ''),
		       COALESCE(d.evaluation_rubrics, ''), COALESCE(d.clil_modules, ''),
		       d.approved_at, d.published_at, d.created_at, d.updated_at
		FROM class_may15_documents d
		JOIN classes c ON d.class_id = c.id
		WHERE d.class_id = $1::uuid AND d.academic_year = $2
	`
	row := r.db.QueryRowContext(ctx, query, classID, academicYear)
	var doc ClassMay15Document
	var approvedAt, publishedAt sql.NullTime

	err := row.Scan(
		&doc.ID, &doc.SchoolID, &doc.ClassID, &doc.ClassName,
		&doc.AcademicYear, &doc.Status,
		&doc.ClassPresentation, &doc.TeachingContinuity,
		&doc.PCTOPathways, &doc.ExamSimulations,
		&doc.EvaluationRubrics, &doc.CLILModules,
		&approvedAt, &publishedAt, &doc.CreatedAt, &doc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if approvedAt.Valid {
		doc.ApprovedAt = &approvedAt.Time
	}
	if publishedAt.Valid {
		doc.PublishedAt = &publishedAt.Time
	}

	return &doc, nil
}

func (r *repository) Upsert(ctx context.Context, doc *ClassMay15Document) error {
	query := `
		INSERT INTO class_may15_documents (
			school_id, class_id, academic_year, status,
			class_presentation, teaching_continuity, pcto_pathways,
			exam_simulations, evaluation_rubrics, clil_modules,
			approved_at, published_at, created_at, updated_at
		) VALUES (
			$1::uuid, $2::uuid, $3, $4,
			$5, $6, $7,
			$8, $9, $10,
			$11, $12, NOW(), NOW()
		)
		ON CONFLICT (class_id, academic_year) DO UPDATE SET
			status = EXCLUDED.status,
			class_presentation = EXCLUDED.class_presentation,
			teaching_continuity = EXCLUDED.teaching_continuity,
			pcto_pathways = EXCLUDED.pcto_pathways,
			exam_simulations = EXCLUDED.exam_simulations,
			evaluation_rubrics = EXCLUDED.evaluation_rubrics,
			clil_modules = EXCLUDED.clil_modules,
			approved_at = COALESCE(EXCLUDED.approved_at, class_may15_documents.approved_at),
			published_at = COALESCE(EXCLUDED.published_at, class_may15_documents.published_at),
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`

	var approvedAt, publishedAt *time.Time
	if doc.Status == StatusApprovatoCdC && doc.ApprovedAt == nil {
		now := time.Now()
		approvedAt = &now
		doc.ApprovedAt = approvedAt
	} else {
		approvedAt = doc.ApprovedAt
	}

	if doc.Status == StatusPubblicato && doc.PublishedAt == nil {
		now := time.Now()
		publishedAt = &now
		doc.PublishedAt = publishedAt
	} else {
		publishedAt = doc.PublishedAt
	}

	return r.db.QueryRowContext(ctx, query,
		doc.SchoolID, doc.ClassID, doc.AcademicYear, doc.Status,
		doc.ClassPresentation, doc.TeachingContinuity, doc.PCTOPathways,
		doc.ExamSimulations, doc.EvaluationRubrics, doc.CLILModules,
		approvedAt, publishedAt,
	).Scan(&doc.ID, &doc.CreatedAt, &doc.UpdatedAt)
}

func (r *repository) Publish(ctx context.Context, classID, academicYear string) (*ClassMay15Document, error) {
	now := time.Now()
	query := `
		UPDATE class_may15_documents
		SET status = 'pubblicato',
		    published_at = COALESCE(published_at, $3),
		    approved_at = COALESCE(approved_at, $3),
		    updated_at = NOW()
		WHERE class_id = $1::uuid AND academic_year = $2
		RETURNING id, school_id, class_id, academic_year, status,
		          class_presentation, teaching_continuity, pcto_pathways,
		          exam_simulations, evaluation_rubrics, clil_modules,
		          approved_at, published_at, created_at, updated_at
	`
	row := r.db.QueryRowContext(ctx, query, classID, academicYear, now)
	var doc ClassMay15Document
	var approvedAt, publishedAt sql.NullTime

	err := row.Scan(
		&doc.ID, &doc.SchoolID, &doc.ClassID, &doc.AcademicYear, &doc.Status,
		&doc.ClassPresentation, &doc.TeachingContinuity, &doc.PCTOPathways,
		&doc.ExamSimulations, &doc.EvaluationRubrics, &doc.CLILModules,
		&approvedAt, &publishedAt, &doc.CreatedAt, &doc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if approvedAt.Valid {
		doc.ApprovedAt = &approvedAt.Time
	}
	if publishedAt.Valid {
		doc.PublishedAt = &publishedAt.Time
	}
	return &doc, nil
}
