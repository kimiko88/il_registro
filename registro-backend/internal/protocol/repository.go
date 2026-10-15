package protocol

import (
	"context"
	"database/sql"
	"time"
)

type Repository interface {
	RegisterDocument(ctx context.Context, entry *ProtocolEntry, entityType, entityID string) error
	GetProtocolEntry(ctx context.Context, id string) (*ProtocolEntry, error)
	ListProtocolEntries(ctx context.Context, schoolID string, year int, flowDirection string) ([]ProtocolEntry, error)
	GetEntityProtocol(ctx context.Context, entityType, entityID string) (*ProtocolEntry, error)
}

type postgresRepo struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) RegisterDocument(ctx context.Context, entry *ProtocolEntry, entityType, entityID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	year := entry.ProtocolYear
	if year == 0 {
		year = time.Now().Year()
		entry.ProtocolYear = year
	}

	// 1. Compute next sequential number atomically with table lock per school/year
	var nextNum int
	lockQuery := `
		SELECT COALESCE(MAX(protocol_number), 0) + 1
		FROM agid_protocol_register
		WHERE school_id = $1 AND protocol_year = $2;
	`
	if err := tx.QueryRowContext(ctx, lockQuery, entry.SchoolID, year).Scan(&nextNum); err != nil {
		return err
	}
	entry.ProtocolNumber = nextNum
	entry.ProtocolDate = time.Now()

	// 2. Insert into protocol register
	insertQuery := `
		INSERT INTO agid_protocol_register (
			school_id, protocol_year, protocol_number, protocol_date, flow_direction,
			classification_title, classification_class, classification_fascicle,
			subject, sender, recipient, document_hash_sha256, document_file_url, protocolled_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id;
	`
	if err := tx.QueryRowContext(ctx, insertQuery,
		entry.SchoolID, entry.ProtocolYear, entry.ProtocolNumber, entry.ProtocolDate, entry.FlowDirection,
		entry.ClassificationTitle, entry.ClassificationClass, entry.ClassificationFascicle,
		entry.Subject, entry.Sender, entry.Recipient, entry.DocumentHashSHA256, entry.DocumentFileURL, entry.ProtocolledBy,
	).Scan(&entry.ID); err != nil {
		return err
	}

	// 3. Link entity if provided
	if entityType != "" && entityID != "" {
		linkQuery := `
			INSERT INTO entity_protocol_links (protocol_id, entity_type, entity_id)
			VALUES ($1, $2, $3::uuid);
		`
		if _, err := tx.ExecContext(ctx, linkQuery, entry.ID, entityType, entityID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *postgresRepo) GetProtocolEntry(ctx context.Context, id string) (*ProtocolEntry, error) {
	query := `
		SELECT id, school_id, protocol_year, protocol_number, protocol_date, flow_direction,
		       classification_title, classification_class, COALESCE(classification_fascicle, ''),
		       subject, sender, recipient, document_hash_sha256, document_file_url, protocolled_by
		FROM agid_protocol_register
		WHERE id = $1;
	`
	var e ProtocolEntry
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID, &e.SchoolID, &e.ProtocolYear, &e.ProtocolNumber, &e.ProtocolDate, &e.FlowDirection,
		&e.ClassificationTitle, &e.ClassificationClass, &e.ClassificationFascicle,
		&e.Subject, &e.Sender, &e.Recipient, &e.DocumentHashSHA256, &e.DocumentFileURL, &e.ProtocolledBy,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *postgresRepo) ListProtocolEntries(ctx context.Context, schoolID string, year int, flowDirection string) ([]ProtocolEntry, error) {
	query := `
		SELECT id, school_id, protocol_year, protocol_number, protocol_date, flow_direction,
		       classification_title, classification_class, COALESCE(classification_fascicle, ''),
		       subject, sender, recipient, document_hash_sha256, document_file_url, protocolled_by
		FROM agid_protocol_register
		WHERE school_id = $1
		  AND ($2 = 0 OR protocol_year = $2)
		  AND ($3 = '' OR flow_direction = $3)
		ORDER BY protocol_number DESC;
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, year, flowDirection)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ProtocolEntry
	for rows.Next() {
		var e ProtocolEntry
		if err := rows.Scan(
			&e.ID, &e.SchoolID, &e.ProtocolYear, &e.ProtocolNumber, &e.ProtocolDate, &e.FlowDirection,
			&e.ClassificationTitle, &e.ClassificationClass, &e.ClassificationFascicle,
			&e.Subject, &e.Sender, &e.Recipient, &e.DocumentHashSHA256, &e.DocumentFileURL, &e.ProtocolledBy,
		); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *postgresRepo) GetEntityProtocol(ctx context.Context, entityType, entityID string) (*ProtocolEntry, error) {
	query := `
		SELECT p.id, p.school_id, p.protocol_year, p.protocol_number, p.protocol_date, p.flow_direction,
		       p.classification_title, p.classification_class, COALESCE(p.classification_fascicle, ''),
		       p.subject, p.sender, p.recipient, p.document_hash_sha256, p.document_file_url, p.protocolled_by
		FROM agid_protocol_register p
		JOIN entity_protocol_links l ON l.protocol_id = p.id
		WHERE l.entity_type = $1 AND l.entity_id = $2::uuid;
	`
	var e ProtocolEntry
	err := r.db.QueryRowContext(ctx, query, entityType, entityID).Scan(
		&e.ID, &e.SchoolID, &e.ProtocolYear, &e.ProtocolNumber, &e.ProtocolDate, &e.FlowDirection,
		&e.ClassificationTitle, &e.ClassificationClass, &e.ClassificationFascicle,
		&e.Subject, &e.Sender, &e.Recipient, &e.DocumentHashSHA256, &e.DocumentFileURL, &e.ProtocolledBy,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}
