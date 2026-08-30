package communications

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repository interface {
	Create(ctx context.Context, msg *Message) error
	List(ctx context.Context, userID, schoolID string) ([]*Message, error)
	ListBacheca(ctx context.Context, schoolID, userID string) ([]*Message, error)
	Delete(ctx context.Context, id string) error
	Sign(ctx context.Context, communicationID string, userID string) error
	SignWithIP(ctx context.Context, communicationID string, userID string, ipAddress string) error
	GetSignatures(ctx context.Context, communicationID string) ([]string, error)
	GetSignatureReport(ctx context.Context, communicationID string) (*SignatureReportResponse, error)
	Get(ctx context.Context, id string) (*Message, error)

	Update(ctx context.Context, id string, subject, body string) error
	MarkAsRead(ctx context.Context, communicationID, userID, ipAddress string) error
	GetUnreadUsers(ctx context.Context, communicationID string) ([]string, error)
	GetUnreadCount(ctx context.Context, userID string) (int, error)
	ListCircolari(ctx context.Context, schoolID, userID, year string) ([]*Message, error)
	Ack(ctx context.Context, communicationID, userID string) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, msg *Message) error {
	msg.ID = uuid.New().String()
	msg.CreatedAt = time.Now()

	query := `
		INSERT INTO communications (id, school_id, sender_id, receiver_ids, subject, body, attachment_url, type, requires_signature, signature_deadline, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	var schoolUUID interface{} = nil
	if msg.SchoolID != nil && *msg.SchoolID != "" {
		schoolUUID = *msg.SchoolID
	}
	_, err := r.db.ExecContext(ctx, query,
		msg.ID, schoolUUID, msg.SenderID, pq.Array(msg.ReceiverIDs), msg.Subject, msg.Body, msg.AttachmentURL,
		msg.Type, msg.RequiresSignature, msg.SignatureDeadline, msg.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) List(ctx context.Context, userID, schoolID string) ([]*Message, error) {
	if userID == "" {
		return []*Message{}, nil
	}
	if _, err := uuid.Parse(userID); err != nil {
		return []*Message{}, nil
	}
	query := `
		SELECT c.id, c.school_id, c.sender_id, c.receiver_ids, c.subject, c.body, c.attachment_url, c.type,
		       COALESCE(c.requires_signature, false), c.signature_deadline, c.created_at,
		       EXISTS(SELECT 1 FROM communication_signatures cs WHERE cs.communication_id = c.id AND cs.user_id = NULLIF($1, '')::uuid) AS is_signed
		FROM communications c
		WHERE (c.sender_id = NULLIF($1, '')::uuid OR $1::text = ANY(c.receiver_ids))
		  AND (c.school_id = NULLIF($2, '')::uuid OR $2 = '')
		ORDER BY c.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []*Message
	for rows.Next() {
		m := &Message{}
		var receivers []string
		var schoolID, attachURL sql.NullString
		if err := rows.Scan(
			&m.ID, &schoolID, &m.SenderID, pq.Array(&receivers), &m.Subject, &m.Body, &attachURL, &m.Type,
			&m.RequiresSignature, &m.SignatureDeadline, &m.CreatedAt, &m.IsSigned,
		); err != nil {
			return nil, err
		}
		if schoolID.Valid {
			m.SchoolID = &schoolID.String
		}
		if attachURL.Valid {
			m.AttachmentURL = &attachURL.String
		}
		m.ReceiverIDs = receivers
		msgs = append(msgs, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if msgs == nil {
		msgs = []*Message{}
	}
	return msgs, nil
}

func (r *PostgresRepository) ListBacheca(ctx context.Context, schoolID, userID string) ([]*Message, error) {
	if userID != "" {
		if _, err := uuid.Parse(userID); err != nil {
			userID = ""
		}
	}
	if schoolID != "" {
		if _, err := uuid.Parse(schoolID); err != nil {
			schoolID = ""
		}
	}
	query := `
		SELECT c.id, c.school_id, c.sender_id, c.receiver_ids, c.subject, c.body, c.attachment_url, c.type,
		       COALESCE(c.requires_signature, false), c.signature_deadline, c.created_at,
		       EXISTS(SELECT 1 FROM communication_signatures cs WHERE cs.communication_id = c.id AND cs.user_id = NULLIF($2, '')::uuid) AS is_signed
		FROM communications c
		WHERE (c.type IN ('circular', 'notice', 'bacheca'))
		  AND ($1 = '' OR c.school_id IS NULL OR c.school_id = NULLIF($1, '')::uuid)
		  AND (array_length(c.receiver_ids, 1) IS NULL OR array_length(c.receiver_ids, 1) = 0 OR $2::text = ANY(c.receiver_ids) OR c.sender_id = NULLIF($2, '')::uuid)
		ORDER BY c.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []*Message
	for rows.Next() {
		m := &Message{}
		var receivers []string
		var schID, attachURL sql.NullString
		if err := rows.Scan(
			&m.ID, &schID, &m.SenderID, pq.Array(&receivers), &m.Subject, &m.Body, &attachURL, &m.Type,
			&m.RequiresSignature, &m.SignatureDeadline, &m.CreatedAt, &m.IsSigned,
		); err != nil {
			return nil, err
		}
		if schID.Valid {
			m.SchoolID = &schID.String
		}
		if attachURL.Valid {
			m.AttachmentURL = &attachURL.String
		}
		m.ReceiverIDs = receivers
		msgs = append(msgs, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if msgs == nil {
		msgs = []*Message{}
	}
	return msgs, nil
}


func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM communications WHERE id = $1::uuid", id)
	return err
}

func (r *PostgresRepository) Sign(ctx context.Context, communicationID string, userID string) error {
	return r.SignWithIP(ctx, communicationID, userID, "")
}

func (r *PostgresRepository) SignWithIP(ctx context.Context, communicationID string, userID string, ipAddress string) error {
	query := `
		INSERT INTO communication_signatures (communication_id, user_id, signed_at, ip_address)
		VALUES ($1::uuid, $2::uuid, NOW(), $3)
		ON CONFLICT (communication_id, user_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, communicationID, userID, ipAddress)
	return err
}

func (r *PostgresRepository) GetSignatures(ctx context.Context, communicationID string) ([]string, error) {
	query := `
		SELECT COALESCE(u.first_name || ' ' || u.last_name, '') AS name
		FROM communication_signatures cs
		JOIN users u ON cs.user_id = u.id
		WHERE cs.communication_id = $1::uuid
		ORDER BY cs.signed_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, communicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return names, nil
}

func (r *PostgresRepository) GetSignatureReport(ctx context.Context, communicationID string) (*SignatureReportResponse, error) {
	msg, err := r.Get(ctx, communicationID)
	if err != nil {
		return nil, err
	}

	report := &SignatureReportResponse{
		CommunicationID:   msg.ID,
		Subject:           msg.Subject,
		SignatureDeadline: msg.SignatureDeadline,
		TotalRecipients:   len(msg.ReceiverIDs),
		Signatures:        []CommunicationSignature{},
	}

	query := `
		SELECT cs.id, cs.communication_id, cs.user_id, cs.signed_at, COALESCE(cs.ip_address, ''),
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS user_name,
		       COALESCE(u.role, '') AS user_role
		FROM communication_signatures cs
		JOIN users u ON cs.user_id = u.id
		WHERE cs.communication_id = $1::uuid
		ORDER BY cs.signed_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, communicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sig CommunicationSignature
		if err := rows.Scan(&sig.ID, &sig.CommunicationID, &sig.UserID, &sig.SignedAt, &sig.IPAddress, &sig.UserName, &sig.UserRole); err != nil {
			return nil, err
		}
		// SEC-31: Redact personal IP address telemetry from API reports
		sig.IPAddress = ""
		report.Signatures = append(report.Signatures, sig)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	report.SignedCount = len(report.Signatures)
	if report.TotalRecipients > 0 {
		report.PendingCount = report.TotalRecipients - report.SignedCount
		if report.PendingCount < 0 {
			report.PendingCount = 0
		}
	}

	return report, nil
}


var ErrNotFound = errors.New("communication not found")

func (r *PostgresRepository) Get(ctx context.Context, id string) (*Message, error) {
	query := `
		SELECT id, school_id, sender_id, receiver_ids, subject, body, attachment_url, type, COALESCE(requires_signature, false), signature_deadline, created_at
		FROM communications
		WHERE id = $1::uuid
	`
	m := &Message{}
	var receivers []string
	var schID, attachURL sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &schID, &m.SenderID, pq.Array(&receivers), &m.Subject, &m.Body, &attachURL, &m.Type,
		&m.RequiresSignature, &m.SignatureDeadline, &m.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if schID.Valid {
		m.SchoolID = &schID.String
	}
	if attachURL.Valid {
		m.AttachmentURL = &attachURL.String
	}
	m.ReceiverIDs = receivers
	return m, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id string, subject, body string) error {
	query := `
		UPDATE communications 
		SET subject = COALESCE(NULLIF($2, ''), subject), body = COALESCE(NULLIF($3, ''), body) 
		WHERE id = $1::uuid 
		  AND NOT EXISTS (SELECT 1 FROM communication_signatures WHERE communication_id = $1::uuid)
	`
	res, err := r.db.ExecContext(ctx, query, id, subject, body)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("impossibile modificare un messaggio inesistente o che contiene già firme digitali")
	}
	return nil
}

func (r *PostgresRepository) MarkAsRead(ctx context.Context, communicationID, userID, ipAddress string) error {
	query := `
		INSERT INTO communication_read_receipts (communication_id, user_id, read_at, ip_address)
		VALUES ($1::uuid, $2::uuid, NOW(), $3)
		ON CONFLICT (communication_id, user_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, communicationID, userID, ipAddress)
	return err
}

func (r *PostgresRepository) GetUnreadUsers(ctx context.Context, communicationID string) ([]string, error) {
	query := `
		SELECT u.first_name || ' ' || u.last_name AS unread_name
		FROM communications c, UNNEST(c.receiver_ids) AS rid(uid)
		JOIN users u ON u.id = rid.uid
		LEFT JOIN communication_read_receipts crr ON crr.communication_id = c.id AND crr.user_id = u.id
		WHERE c.id = $1::uuid AND crr.read_at IS NULL
	`
	rows, err := r.db.QueryContext(ctx, query, communicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var unread []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			unread = append(unread, name)
		}
	}
	return unread, rows.Err()
}

func (r *PostgresRepository) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM communications c
		LEFT JOIN communication_read_receipts crr ON crr.communication_id = c.id AND crr.user_id = $1::uuid
		WHERE ($1 = ANY(c.receiver_ids) OR ARRAY_LENGTH(c.receiver_ids, 1) IS NULL)
		  AND crr.read_at IS NULL
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *PostgresRepository) ListCircolari(ctx context.Context, schoolID, userID, year string) ([]*Message, error) {
	query := `
		SELECT c.id, c.school_id, c.sender_id, c.receiver_ids, c.subject, c.body, c.attachment_url, c.type,
		       COALESCE(c.requires_signature, false), c.signature_deadline, c.created_at,
		       EXISTS(SELECT 1 FROM communication_signatures cs WHERE cs.communication_id = c.id AND cs.user_id = $2::uuid) AS is_signed
		FROM communications c
		WHERE c.type = 'circular'
		  AND ($1 = '' OR c.school_id IS NULL OR c.school_id = NULLIF($1, '')::uuid)
		  AND ($3 = '' OR EXTRACT(YEAR FROM c.created_at)::text = $3)
		ORDER BY c.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, userID, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []*Message
	for rows.Next() {
		m := &Message{}
		var receivers []string
		var schID, attachURL sql.NullString
		if err := rows.Scan(
			&m.ID, &schID, &m.SenderID, pq.Array(&receivers), &m.Subject, &m.Body, &attachURL, &m.Type,
			&m.RequiresSignature, &m.SignatureDeadline, &m.CreatedAt, &m.IsSigned,
		); err != nil {
			return nil, err
		}
		if schID.Valid {
			m.SchoolID = &schID.String
		}
		if attachURL.Valid {
			m.AttachmentURL = &attachURL.String
		}
		m.ReceiverIDs = receivers
		m.IsOfficialCircular = true
		msgs = append(msgs, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return msgs, nil
}


func (r *PostgresRepository) Ack(ctx context.Context, communicationID, userID string) error {
	query := `
		INSERT INTO communication_acks (communication_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (communication_id, user_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, communicationID, userID)
	return err
}
