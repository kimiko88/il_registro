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
	GetMeetingWithCoordinator(ctx context.Context, meetingID string) (*CouncilMeeting, string, error)

	CreateVerbale(ctx context.Context, v *MeetingVerbale) error
	GetVerbaleByID(ctx context.Context, id string, userID string) (*MeetingVerbale, error)
	UpdateVerbale(ctx context.Context, v *MeetingVerbale) error
	DeleteVerbale(ctx context.Context, id string) error
	ListVerbali(ctx context.Context, meetingID string, userID string, onlySigned bool) ([]*MeetingVerbale, error)
	ListAllVerbali(ctx context.Context, schoolID, classID, userID string, onlySigned bool) ([]*MeetingVerbale, error)
	SignVerbale(ctx context.Context, verbaleID, userID, ipAddress string) error
	MarkVerbaleSigned(ctx context.Context, verbaleID string) error
	GetSignatures(ctx context.Context, verbaleID string) ([]VerbaleSignature, error)
	ClassBelongsToSchool(ctx context.Context, classID, schoolID string) (bool, error)

	CreateTemplate(ctx context.Context, t *MeetingVerbaleTemplate) error
	ListTemplates(ctx context.Context, schoolID, meetingType string) ([]*MeetingVerbaleTemplate, error)
	GetTemplateByID(ctx context.Context, id string) (*MeetingVerbaleTemplate, error)
	UpdateTemplate(ctx context.Context, t *MeetingVerbaleTemplate) error
	DeleteTemplate(ctx context.Context, id, schoolID string) error
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
	if m.MeetingType == "" {
		m.MeetingType = "consiglio_classe"
	}

	query := `
		INSERT INTO council_meetings (id, school_id, class_id, meeting_type, title, date, start_time, end_time, agenda, created_by, created_at)
		VALUES ($1, $2, NULLIF($3, '')::uuid, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		m.ID, m.SchoolID, m.ClassID, m.MeetingType, m.Title, m.Date, m.StartTime, m.EndTime, m.Agenda, m.CreatedBy, m.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) ListMeetings(ctx context.Context, schoolID, classID string) ([]*CouncilMeeting, error) {
	query := `
		SELECT id, school_id, COALESCE(class_id::text, ''), COALESCE(meeting_type, 'consiglio_classe'),
		       title, date, start_time, end_time, COALESCE(agenda, ''), created_by, created_at
		FROM council_meetings
		WHERE school_id = $1::uuid AND ($2 = '' OR class_id = NULLIF($2, '')::uuid)
		ORDER BY date DESC, start_time DESC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, classID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var meetings []*CouncilMeeting
	for rows.Next() {
		m := &CouncilMeeting{}
		if err := rows.Scan(&m.ID, &m.SchoolID, &m.ClassID, &m.MeetingType, &m.Title, &m.Date, &m.StartTime, &m.EndTime, &m.Agenda, &m.CreatedBy, &m.CreatedAt); err != nil {
			return nil, err
		}
		meetings = append(meetings, m)
	}
	return meetings, rows.Err()
}

func (r *PostgresRepository) GetMeetingByID(ctx context.Context, id string) (*CouncilMeeting, error) {
	query := `
		SELECT id, school_id, COALESCE(class_id::text, ''), COALESCE(meeting_type, 'consiglio_classe'),
		       title, date, start_time, end_time, COALESCE(agenda, ''), created_by, created_at
		FROM council_meetings
		WHERE id = $1::uuid
	`
	m := &CouncilMeeting{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.SchoolID, &m.ClassID, &m.MeetingType, &m.Title, &m.Date, &m.StartTime, &m.EndTime, &m.Agenda, &m.CreatedBy, &m.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *PostgresRepository) GetMeetingWithCoordinator(ctx context.Context, meetingID string) (*CouncilMeeting, string, error) {
	query := `
		SELECT cm.id, cm.school_id, COALESCE(cm.class_id::text, ''), COALESCE(cm.meeting_type, 'consiglio_classe'),
		       cm.title, cm.date, cm.start_time, cm.end_time, COALESCE(cm.agenda, ''), cm.created_by, cm.created_at,
		       COALESCE(c.coordinator_id::text, '')
		FROM council_meetings cm
		LEFT JOIN classes c ON cm.class_id = c.id
		WHERE cm.id = $1::uuid
	`
	m := &CouncilMeeting{}
	var coordinatorID string
	err := r.db.QueryRowContext(ctx, query, meetingID).Scan(
		&m.ID, &m.SchoolID, &m.ClassID, &m.MeetingType, &m.Title, &m.Date, &m.StartTime, &m.EndTime, &m.Agenda, &m.CreatedBy, &m.CreatedAt,
		&coordinatorID,
	)
	if err != nil {
		return nil, "", err
	}
	return m, coordinatorID, nil
}

func (r *PostgresRepository) CreateVerbale(ctx context.Context, v *MeetingVerbale) error {
	v.ID = uuid.New().String()
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()
	if v.Status == "" {
		v.Status = "draft"
	}

	query := `
		INSERT INTO meeting_verbali (id, meeting_id, title, content, secretary_id, president_id, is_published, is_signed, signed_at, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.ExecContext(ctx, query,
		v.ID, v.MeetingID, v.Title, v.Content, v.SecretaryID, v.PresidentID, v.IsPublished, v.IsSigned, v.SignedAt, v.Status, v.CreatedAt, v.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetVerbaleByID(ctx context.Context, id string, userID string) (*MeetingVerbale, error) {
	userUUID := userID
	if userUUID == "" {
		userUUID = "00000000-0000-0000-0000-000000000000"
	}
	query := `
		SELECT v.id, v.meeting_id, v.title, v.content, v.secretary_id, v.president_id, v.is_published,
		       COALESCE(v.is_signed, false), v.signed_at, COALESCE(v.status, 'draft'), v.created_at, v.updated_at,
		       EXISTS(SELECT 1 FROM verbale_signatures vs WHERE vs.verbale_id = v.id AND vs.user_id = $2::uuid) AS is_signed_by_me
		FROM meeting_verbali v
		WHERE v.id = $1::uuid
	`
	v := &MeetingVerbale{}
	var secID, presID sql.NullString
	err := r.db.QueryRowContext(ctx, query, id, userUUID).Scan(
		&v.ID, &v.MeetingID, &v.Title, &v.Content, &secID, &presID, &v.IsPublished,
		&v.IsSigned, &v.SignedAt, &v.Status, &v.CreatedAt, &v.UpdatedAt, &v.IsSignedByMe,
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

func (r *PostgresRepository) UpdateVerbale(ctx context.Context, v *MeetingVerbale) error {
	v.UpdatedAt = time.Now()
	query := `
		UPDATE meeting_verbali
		SET title = $1, content = $2, secretary_id = $3, president_id = $4, is_published = $5, updated_at = $6
		WHERE id = $7::uuid
	`
	_, err := r.db.ExecContext(ctx, query,
		v.Title, v.Content, v.SecretaryID, v.PresidentID, v.IsPublished, v.UpdatedAt, v.ID,
	)
	return err
}

func (r *PostgresRepository) DeleteVerbale(ctx context.Context, id string) error {
	query := `DELETE FROM meeting_verbali WHERE id = $1::uuid`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresRepository) MarkVerbaleSigned(ctx context.Context, verbaleID string) error {
	query := `
		UPDATE meeting_verbali
		SET is_signed = TRUE, status = 'signed', signed_at = NOW(), is_published = TRUE, updated_at = NOW()
		WHERE id = $1::uuid
	`
	_, err := r.db.ExecContext(ctx, query, verbaleID)
	return err
}

func (r *PostgresRepository) ListVerbali(ctx context.Context, meetingID string, userID string, onlySigned bool) ([]*MeetingVerbale, error) {
	userUUID := userID
	if userUUID == "" {
		userUUID = "00000000-0000-0000-0000-000000000000"
	}
	query := `
		SELECT v.id, v.meeting_id, v.title, v.content, v.secretary_id, v.president_id, v.is_published,
		       COALESCE(v.is_signed, false), v.signed_at, COALESCE(v.status, 'draft'), v.created_at, v.updated_at,
		       EXISTS(SELECT 1 FROM verbale_signatures vs WHERE vs.verbale_id = v.id AND vs.user_id = $2::uuid) AS is_signed_by_me
		FROM meeting_verbali v
		WHERE v.meeting_id = $1::uuid
		  AND (NOT $3 OR v.is_signed = TRUE)
		ORDER BY v.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, meetingID, userUUID, onlySigned)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var list []*MeetingVerbale
	for rows.Next() {
		v := &MeetingVerbale{}
		var secID, presID sql.NullString
		if err := rows.Scan(
			&v.ID, &v.MeetingID, &v.Title, &v.Content, &secID, &presID, &v.IsPublished,
			&v.IsSigned, &v.SignedAt, &v.Status, &v.CreatedAt, &v.UpdatedAt, &v.IsSignedByMe,
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

func (r *PostgresRepository) ListAllVerbali(ctx context.Context, schoolID, classID, userID string, onlySigned bool) ([]*MeetingVerbale, error) {
	userUUID := userID
	if userUUID == "" {
		userUUID = "00000000-0000-0000-0000-000000000000"
	}
	query := `
		SELECT v.id, v.meeting_id, v.title, v.content, v.secretary_id, v.president_id, v.is_published,
		       COALESCE(v.is_signed, false), v.signed_at, COALESCE(v.status, 'draft'), v.created_at, v.updated_at,
		       EXISTS(SELECT 1 FROM verbale_signatures vs WHERE vs.verbale_id = v.id AND vs.user_id = $2::uuid) AS is_signed_by_me
		FROM meeting_verbali v
		JOIN council_meetings cm ON v.meeting_id = cm.id
		WHERE cm.school_id = $1::uuid
		  AND ($3 = '' OR cm.class_id = NULLIF($3, '')::uuid)
		  AND (NOT $4 OR v.is_signed = TRUE)
		ORDER BY v.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, userUUID, classID, onlySigned)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var list []*MeetingVerbale
	for rows.Next() {
		v := &MeetingVerbale{}
		var secID, presID sql.NullString
		if err := rows.Scan(
			&v.ID, &v.MeetingID, &v.Title, &v.Content, &secID, &presID, &v.IsPublished,
			&v.IsSigned, &v.SignedAt, &v.Status, &v.CreatedAt, &v.UpdatedAt, &v.IsSignedByMe,
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
	defer func() { _ = rows.Close() }()

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

func (r *PostgresRepository) ClassBelongsToSchool(ctx context.Context, classID, schoolID string) (bool, error) {
	if r.db == nil || classID == "" || schoolID == "" {
		return false, nil
	}
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM classes WHERE id = $1::uuid AND school_id = $2::uuid AND deleted_at IS NULL)`
	err := r.db.QueryRowContext(ctx, query, classID, schoolID).Scan(&exists)
	return exists, err
}

// ----------------- Template CRUD -----------------

func (r *PostgresRepository) CreateTemplate(ctx context.Context, t *MeetingVerbaleTemplate) error {
	t.ID = uuid.New().String()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	query := `
		INSERT INTO meeting_verbale_templates (id, school_id, title, meeting_type, description, default_agenda, template_content, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		t.ID, t.SchoolID, t.Title, t.MeetingType, t.Description, t.DefaultAgenda, t.TemplateContent, t.CreatedBy, t.CreatedAt, t.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) ListTemplates(ctx context.Context, schoolID, meetingType string) ([]*MeetingVerbaleTemplate, error) {
	query := `
		SELECT id, school_id, title, meeting_type, COALESCE(description, ''), default_agenda, template_content, created_by, created_at, updated_at
		FROM meeting_verbale_templates
		WHERE school_id = $1::uuid AND ($2 = '' OR meeting_type = $2)
		ORDER BY title ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, meetingType)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var templates []*MeetingVerbaleTemplate
	for rows.Next() {
		t := &MeetingVerbaleTemplate{}
		if err := rows.Scan(
			&t.ID, &t.SchoolID, &t.Title, &t.MeetingType, &t.Description, &t.DefaultAgenda, &t.TemplateContent, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		templates = append(templates, t)
	}
	return templates, rows.Err()
}

func (r *PostgresRepository) GetTemplateByID(ctx context.Context, id string) (*MeetingVerbaleTemplate, error) {
	query := `
		SELECT id, school_id, title, meeting_type, COALESCE(description, ''), default_agenda, template_content, created_by, created_at, updated_at
		FROM meeting_verbale_templates
		WHERE id = $1::uuid
	`
	t := &MeetingVerbaleTemplate{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.SchoolID, &t.Title, &t.MeetingType, &t.Description, &t.DefaultAgenda, &t.TemplateContent, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *PostgresRepository) UpdateTemplate(ctx context.Context, t *MeetingVerbaleTemplate) error {
	t.UpdatedAt = time.Now()
	query := `
		UPDATE meeting_verbale_templates
		SET title = $1, meeting_type = $2, description = $3, default_agenda = $4, template_content = $5, updated_at = $6
		WHERE id = $7::uuid AND school_id = $8::uuid
	`
	_, err := r.db.ExecContext(ctx, query,
		t.Title, t.MeetingType, t.Description, t.DefaultAgenda, t.TemplateContent, t.UpdatedAt, t.ID, t.SchoolID,
	)
	return err
}

func (r *PostgresRepository) DeleteTemplate(ctx context.Context, id, schoolID string) error {
	query := `DELETE FROM meeting_verbale_templates WHERE id = $1::uuid AND school_id = $2::uuid`
	_, err := r.db.ExecContext(ctx, query, id, schoolID)
	return err
}
