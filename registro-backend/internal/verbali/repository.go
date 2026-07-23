package verbali

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateMeeting(ctx context.Context, m *CouncilMeeting) error
	ListMeetings(ctx context.Context, schoolID, classID string) ([]*CouncilMeeting, error)
	GetMeetingByID(ctx context.Context, id string) (*CouncilMeeting, error)

	CreateVerbale(ctx context.Context, v *MeetingVerbale) error
	GetVerbaleByID(ctx context.Context, id string, userID string) (*MeetingVerbale, error)
	ListVerbali(ctx context.Context, meetingID string, userID string) ([]*MeetingVerbale, error)
	SignVerbale(ctx context.Context, verbaleID, userID, ipAddress string) error
	GetSignatures(ctx context.Context, verbaleID string) ([]VerbaleSignature, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateMeeting(ctx context.Context, m *CouncilMeeting) error {
	m.ID = uuid.New().String()
	m.CreatedAt = time.Now()

	query := `
		INSERT INTO council_meetings (id, school_id, class_id, title, date, start_time, end_time, agenda, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		m.ID, m.SchoolID, m.ClassID, m.Title, m.Date, m.StartTime, m.EndTime, m.Agenda, m.CreatedBy, m.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) ListMeetings(ctx context.Context, schoolID, classID string) ([]*CouncilMeeting, error) {
	query := `
		SELECT id, school_id, class_id, title, date, start_time, end_time, COALESCE(agenda, ''), created_by, created_at
		FROM council_meetings
		WHERE school_id = $1::uuid AND ($2 = '' OR class_id = $2::uuid)
		ORDER BY date DESC, start_time DESC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meetings []*CouncilMeeting
	for rows.Next() {
		m := &CouncilMeeting{}
		if err := rows.Scan(&m.ID, &m.SchoolID, &m.ClassID, &m.Title, &m.Date, &m.StartTime, &m.EndTime, &m.Agenda, &m.CreatedBy, &m.CreatedAt); err != nil {
			return nil, err
		}
		meetings = append(meetings, m)
	}
	return meetings, rows.Err()
}

func (r *PostgresRepository) GetMeetingByID(ctx context.Context, id string) (*CouncilMeeting, error) {
	query := `
		SELECT id, school_id, class_id, title, date, start_time, end_time, COALESCE(agenda, ''), created_by, created_at
		FROM council_meetings
		WHERE id = $1::uuid
	`
	m := &CouncilMeeting{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&m.ID, &m.SchoolID, &m.ClassID, &m.Title, &m.Date, &m.StartTime, &m.EndTime, &m.Agenda, &m.CreatedBy, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *PostgresRepository) CreateVerbale(ctx context.Context, v *MeetingVerbale) error {
	v.ID = uuid.New().String()
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()

	query := `
		INSERT INTO meeting_verbali (id, meeting_id, title, content, secretary_id, president_id, is_published, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		v.ID, v.MeetingID, v.Title, v.Content, v.SecretaryID, v.PresidentID, v.IsPublished, v.CreatedAt, v.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetVerbaleByID(ctx context.Context, id string, userID string) (*MeetingVerbale, error) {
	userUUID := userID
	if userUUID == "" {
		userUUID = "00000000-0000-0000-0000-000000000000"
	}
	query := `
		SELECT v.id, v.meeting_id, v.title, v.content, v.secretary_id, v.president_id, v.is_published, v.created_at, v.updated_at,
		       EXISTS(SELECT 1 FROM verbale_signatures vs WHERE vs.verbale_id = v.id AND vs.user_id = $2::uuid) AS is_signed_by_me
		FROM meeting_verbali v
		WHERE v.id = $1::uuid
	`
	v := &MeetingVerbale{}
	var secID, presID sql.NullString
	err := r.db.QueryRowContext(ctx, query, id, userUUID).Scan(
		&v.ID, &v.MeetingID, &v.Title, &v.Content, &secID, &presID, &v.IsPublished, &v.CreatedAt, &v.UpdatedAt, &v.IsSignedByMe,
	)
	if err != nil {
		return nil, err
	}
	if secID.Valid {
		v.SecretaryID = &secID.String
	}
	if presID.Valid {
		v.PresidentID = &presID.String
	}
	return v, nil
}

func (r *PostgresRepository) ListVerbali(ctx context.Context, meetingID string, userID string) ([]*MeetingVerbale, error) {
	userUUID := userID
	if userUUID == "" {
		userUUID = "00000000-0000-0000-0000-000000000000"
	}
	query := `
		SELECT v.id, v.meeting_id, v.title, v.content, v.secretary_id, v.president_id, v.is_published, v.created_at, v.updated_at,
		       EXISTS(SELECT 1 FROM verbale_signatures vs WHERE vs.verbale_id = v.id AND vs.user_id = $2::uuid) AS is_signed_by_me
		FROM meeting_verbali v
		WHERE v.meeting_id = $1::uuid
		ORDER BY v.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, meetingID, userUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*MeetingVerbale
	for rows.Next() {
		v := &MeetingVerbale{}
		var secID, presID sql.NullString
		if err := rows.Scan(
			&v.ID, &v.MeetingID, &v.Title, &v.Content, &secID, &presID, &v.IsPublished, &v.CreatedAt, &v.UpdatedAt, &v.IsSignedByMe,
		); err != nil {
			return nil, err
		}
		if secID.Valid {
			v.SecretaryID = &secID.String
		}
		if presID.Valid {
			v.PresidentID = &presID.String
		}
		list = append(list, v)
	}
	return list, rows.Err()
}

func (r *PostgresRepository) SignVerbale(ctx context.Context, verbaleID, userID, ipAddress string) error {
	query := `
		INSERT INTO verbale_signatures (verbale_id, user_id, signed_at, ip_address)
		VALUES ($1::uuid, $2::uuid, NOW(), $3)
		ON CONFLICT (verbale_id, user_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, verbaleID, userID, ipAddress)
	return err
}

func (r *PostgresRepository) GetSignatures(ctx context.Context, verbaleID string) ([]VerbaleSignature, error) {
	query := `
		SELECT vs.id, vs.verbale_id, vs.user_id, vs.signed_at, COALESCE(vs.ip_address, ''),
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS user_name
		FROM verbale_signatures vs
		JOIN users u ON vs.user_id = u.id
		WHERE vs.verbale_id = $1::uuid
		ORDER BY vs.signed_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, verbaleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sigs []VerbaleSignature
	for rows.Next() {
		var s VerbaleSignature
		if err := rows.Scan(&s.ID, &s.VerbaleID, &s.UserID, &s.SignedAt, &s.IPAddress, &s.UserName); err != nil {
			return nil, err
		}
		sigs = append(sigs, s)
	}
	return sigs, rows.Err()
}
