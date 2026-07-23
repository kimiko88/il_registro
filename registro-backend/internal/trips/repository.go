package trips

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repository interface {
	CreateTrip(ctx context.Context, t *EducationalTrip) error
	GetTripByID(ctx context.Context, id string) (*EducationalTrip, error)
	ListTrips(ctx context.Context, schoolID, studentID string) ([]*EducationalTrip, error)

	SubmitConsent(ctx context.Context, c *TripConsent) error
	ListConsents(ctx context.Context, tripID string) ([]*TripConsent, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateTrip(ctx context.Context, t *EducationalTrip) error {
	t.ID = uuid.New().String()
	t.CreatedAt = time.Now()

	query := `
		INSERT INTO educational_trips (id, school_id, title, destination, departure_date, return_date, description, accompanying_teachers, class_ids, created_at)
		VALUES ($1, $2::uuid, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		t.ID, t.SchoolID, t.Title, t.Destination, t.DepartureDate, t.ReturnDate,
		t.Description, t.AccompanyingTeachers, pq.Array(t.ClassIDs), t.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) GetTripByID(ctx context.Context, id string) (*EducationalTrip, error) {
	query := `
		SELECT id, school_id, title, destination, departure_date, return_date,
		       COALESCE(description, ''), COALESCE(accompanying_teachers, ''), class_ids, created_at
		FROM educational_trips
		WHERE id = $1::uuid
	`
	t := &EducationalTrip{}
	var classes []string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.SchoolID, &t.Title, &t.Destination, &t.DepartureDate, &t.ReturnDate,
		&t.Description, &t.AccompanyingTeachers, pq.Array(&classes), &t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	t.ClassIDs = classes
	return t, nil
}

func (r *PostgresRepository) ListTrips(ctx context.Context, schoolID, studentID string) ([]*EducationalTrip, error) {
	stdUUID := studentID
	if stdUUID == "" {
		stdUUID = "00000000-0000-0000-0000-000000000000"
	}
	query := `
		SELECT t.id, t.school_id, t.title, t.destination, t.departure_date, t.return_date,
		       COALESCE(t.description, ''), COALESCE(t.accompanying_teachers, ''), t.class_ids, t.created_at,
		       COALESCE(tc.status, 'pending') AS consent_status
		FROM educational_trips t
		LEFT JOIN trip_consents tc ON tc.trip_id = t.id AND tc.student_id = $2::uuid
		WHERE t.school_id = $1::uuid
		ORDER BY t.departure_date ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, stdUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []*EducationalTrip
	for rows.Next() {
		t := &EducationalTrip{}
		var classes []string
		if err := rows.Scan(
			&t.ID, &t.SchoolID, &t.Title, &t.Destination, &t.DepartureDate, &t.ReturnDate,
			&t.Description, &t.AccompanyingTeachers, pq.Array(&classes), &t.CreatedAt, &t.ConsentStatus,
		); err != nil {
			return nil, err
		}
		t.ClassIDs = classes
		trips = append(trips, t)
	}
	return trips, rows.Err()
}

func (r *PostgresRepository) SubmitConsent(ctx context.Context, c *TripConsent) error {
	c.ID = uuid.New().String()
	c.SignedAt = time.Now()

	query := `
		INSERT INTO trip_consents (id, trip_id, student_id, parent_id, status, signed_at, ip_address)
		VALUES ($1, $2::uuid, $3::uuid, $4::uuid, $5, $6, $7)
		ON CONFLICT (trip_id, student_id) DO UPDATE
		SET parent_id = EXCLUDED.parent_id, status = EXCLUDED.status, signed_at = EXCLUDED.signed_at, ip_address = EXCLUDED.ip_address
	`
	var parentUUID interface{} = nil
	if c.ParentID != nil && *c.ParentID != "" {
		parentUUID = *c.ParentID
	}
	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.TripID, c.StudentID, parentUUID, c.Status, c.SignedAt, c.IPAddress,
	)
	return err
}

func (r *PostgresRepository) ListConsents(ctx context.Context, tripID string) ([]*TripConsent, error) {
	query := `
		SELECT tc.id, tc.trip_id, tc.student_id, tc.parent_id, tc.status, tc.signed_at, COALESCE(tc.ip_address, ''),
		       COALESCE(su.first_name || ' ' || su.last_name, '') AS student_name,
		       COALESCE(pu.first_name || ' ' || pu.last_name, '') AS parent_name
		FROM trip_consents tc
		JOIN users su ON tc.student_id = su.id
		LEFT JOIN users pu ON tc.parent_id = pu.id
		WHERE tc.trip_id = $1::uuid
		ORDER BY su.last_name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, tripID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*TripConsent
	for rows.Next() {
		c := &TripConsent{}
		var parentID sql.NullString
		if err := rows.Scan(&c.ID, &c.TripID, &c.StudentID, &parentID, &c.Status, &c.SignedAt, &c.IPAddress, &c.StudentName, &c.ParentName); err != nil {
			return nil, err
		}
		if parentID.Valid {
			c.ParentID = &parentID.String
		}
		list = append(list, c)
	}
	return list, rows.Err()
}
