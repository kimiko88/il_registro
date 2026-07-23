package general_meetings

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Repository interface {
	Create(ctx context.Context, m *GeneralMeeting) error
	List(ctx context.Context, schoolID, userID, role string) ([]*GeneralMeeting, error)
	GetByID(ctx context.Context, id, userID string) (*GeneralMeeting, error)
	Delete(ctx context.Context, id string) error

	RegisterUser(ctx context.Context, meetingID, userID string) error
	UnregisterUser(ctx context.Context, meetingID, userID string) error
	ListRegistrations(ctx context.Context, meetingID string) ([]*GeneralMeetingRegistration, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, m *GeneralMeeting) error {
	query := `
		INSERT INTO general_meetings (
			school_id, title, description, location, meeting_date,
			registration_deadline, max_participants, is_mandatory, target_roles,
			created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		m.SchoolID, m.Title, m.Description, m.Location, m.MeetingDate,
		m.RegistrationDeadline, m.MaxParticipants, m.IsMandatory, pq.Array(m.TargetRoles),
		m.CreatedBy,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (r *repository) List(ctx context.Context, schoolID, userID, role string) ([]*GeneralMeeting, error) {
	query := `
		SELECT 
			m.id, m.school_id, m.title, m.description, m.location, m.meeting_date,
			m.registration_deadline, m.max_participants, m.is_mandatory, m.target_roles,
			m.created_by, m.created_at, m.updated_at,
			(SELECT COUNT(*) FROM general_meeting_registrations r WHERE r.meeting_id = m.id) AS reg_count,
			EXISTS (SELECT 1 FROM general_meeting_registrations r WHERE r.meeting_id = m.id AND r.user_id = $2::uuid) AS is_reg
		FROM general_meetings m
		WHERE m.school_id = $1::uuid
		  AND ($3 = '' OR $3 = ANY(m.target_roles) OR ARRAY_LENGTH(m.target_roles, 1) IS NULL OR ARRAY_LENGTH(m.target_roles, 1) = 0)
		ORDER BY m.meeting_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, userID, role)
	if err != nil {
		return nil, fmt.Errorf("List general_meetings: %w", err)
	}
	defer rows.Close()

	var res []*GeneralMeeting
	for rows.Next() {
		var m GeneralMeeting
		var targetRoles pq.StringArray
		if err := rows.Scan(
			&m.ID, &m.SchoolID, &m.Title, &m.Description, &m.Location, &m.MeetingDate,
			&m.RegistrationDeadline, &m.MaxParticipants, &m.IsMandatory, &targetRoles,
			&m.CreatedBy, &m.CreatedAt, &m.UpdatedAt,
			&m.RegistrationsCount, &m.IsRegistered,
		); err != nil {
			return nil, err
		}
		m.TargetRoles = []string(targetRoles)
		res = append(res, &m)
	}
	return res, rows.Err()
}

func (r *repository) GetByID(ctx context.Context, id, userID string) (*GeneralMeeting, error) {
	query := `
		SELECT 
			m.id, m.school_id, m.title, m.description, m.location, m.meeting_date,
			m.registration_deadline, m.max_participants, m.is_mandatory, m.target_roles,
			m.created_by, m.created_at, m.updated_at,
			(SELECT COUNT(*) FROM general_meeting_registrations r WHERE r.meeting_id = m.id) AS reg_count,
			EXISTS (SELECT 1 FROM general_meeting_registrations r WHERE r.meeting_id = m.id AND r.user_id = $2::uuid) AS is_reg
		FROM general_meetings m
		WHERE m.id = $1::uuid
	`
	var m GeneralMeeting
	var targetRoles pq.StringArray
	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(
		&m.ID, &m.SchoolID, &m.Title, &m.Description, &m.Location, &m.MeetingDate,
		&m.RegistrationDeadline, &m.MaxParticipants, &m.IsMandatory, &targetRoles,
		&m.CreatedBy, &m.CreatedAt, &m.UpdatedAt,
		&m.RegistrationsCount, &m.IsRegistered,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("meeting not found")
		}
		return nil, err
	}
	m.TargetRoles = []string(targetRoles)
	return &m, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM general_meetings WHERE id = $1::uuid`, id)
	return err
}

func (r *repository) RegisterUser(ctx context.Context, meetingID, userID string) error {
	query := `
		INSERT INTO general_meeting_registrations (meeting_id, user_id, registered_at)
		VALUES ($1::uuid, $2::uuid, NOW())
		ON CONFLICT (meeting_id, user_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, meetingID, userID)
	return err
}

func (r *repository) UnregisterUser(ctx context.Context, meetingID, userID string) error {
	query := `DELETE FROM general_meeting_registrations WHERE meeting_id = $1::uuid AND user_id = $2::uuid`
	_, err := r.db.ExecContext(ctx, query, meetingID, userID)
	return err
}

func (r *repository) ListRegistrations(ctx context.Context, meetingID string) ([]*GeneralMeetingRegistration, error) {
	query := `
		SELECT r.id, r.meeting_id, r.user_id, COALESCE(u.first_name || ' ' || u.last_name, u.email) AS user_name,
		       u.email AS user_email, u.role AS user_role, r.registered_at
		FROM general_meeting_registrations r
		JOIN users u ON u.id = r.user_id
		WHERE r.meeting_id = $1::uuid
		ORDER BY r.registered_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, meetingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []*GeneralMeetingRegistration
	for rows.Next() {
		var reg GeneralMeetingRegistration
		if err := rows.Scan(&reg.ID, &reg.MeetingID, &reg.UserID, &reg.UserName, &reg.UserEmail, &reg.UserRole, &reg.RegisteredAt); err != nil {
			return nil, err
		}
		res = append(res, &reg)
	}
	return res, rows.Err()
}
// Ensure time is used
var _ = time.Now
