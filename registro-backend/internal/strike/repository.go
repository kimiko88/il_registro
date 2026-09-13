package strike

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repository interface {
	CreateNotice(ctx context.Context, n *StrikeNotice) error
	GetNoticeByID(ctx context.Context, id, schoolID string) (*StrikeNotice, error)
	ListNotices(ctx context.Context, schoolID string) ([]*StrikeNotice, error)
	DeleteNotice(ctx context.Context, id, schoolID string) error
	UpsertDeclaration(ctx context.Context, d *StrikeDeclaration) error
	GetDeclaration(ctx context.Context, noticeID, userID string) (*StrikeDeclaration, error)
	GetNoticeSummary(ctx context.Context, noticeID, schoolID string) (*StrikeNoticeSummaryResponse, error)
	CreateBachecaCommunication(ctx context.Context, schoolID, senderID, title, content string, deadline time.Time) (string, error)
	ResolveSchoolID(ctx context.Context, userID string) string
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ResolveSchoolID(ctx context.Context, userID string) string {
	var schoolID string
	if userID != "" {
		_ = r.db.QueryRowContext(ctx, `SELECT COALESCE(school_id::text, '') FROM users WHERE id = $1`, userID).Scan(&schoolID)
	}
	if schoolID == "" {
		_ = r.db.QueryRowContext(ctx, `SELECT COALESCE(id::text, '') FROM schools LIMIT 1`).Scan(&schoolID)
	}
	return schoolID
}

func (r *PostgresRepository) CreateNotice(ctx context.Context, n *StrikeNotice) error {
	n.ID = uuid.New().String()
	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()
	if n.SchoolID == "" {
		n.SchoolID = r.ResolveSchoolID(ctx, n.CreatedBy)
	}

	query := `
		INSERT INTO strike_notices (
			id, school_id, title, proclaimed_by, strike_date, declaration_deadline,
			content, created_by, is_published, communication_id, created_at, updated_at
		) VALUES (
			$1, $2::uuid, $3, $4, $5::date, $6,
			$7, $8::uuid, $9, NULLIF($10, '')::uuid, $11, $12
		)
	`
	commID := ""
	if n.CommunicationID != nil {
		commID = *n.CommunicationID
	}

	_, err := r.db.ExecContext(ctx, query,
		n.ID, n.SchoolID, n.Title, n.ProclaimedBy, n.StrikeDate, n.DeclarationDeadline,
		n.Content, n.CreatedBy, n.IsPublished, commID, n.CreatedAt, n.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetNoticeByID(ctx context.Context, id, schoolID string) (*StrikeNotice, error) {
	query := `
		SELECT sn.id, sn.school_id, sn.title, sn.proclaimed_by, sn.strike_date::text,
		       sn.declaration_deadline, sn.content, sn.created_by,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS created_by_name,
		       sn.is_published, COALESCE(sn.communication_id::text, ''), sn.created_at, sn.updated_at
		FROM strike_notices sn
		LEFT JOIN users u ON sn.created_by = u.id
		WHERE sn.id = $1::uuid AND (sn.school_id = NULLIF($2, '')::uuid OR $2 = '')
	`
	n := &StrikeNotice{}
	var commID string
	err := r.db.QueryRowContext(ctx, query, id, schoolID).Scan(
		&n.ID, &n.SchoolID, &n.Title, &n.ProclaimedBy, &n.StrikeDate,
		&n.DeclarationDeadline, &n.Content, &n.CreatedBy, &n.CreatedByName,
		&n.IsPublished, &commID, &n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if commID != "" {
		n.CommunicationID = &commID
	}
	n.IsExpired = time.Now().After(n.DeclarationDeadline)
	return n, nil
}

func (r *PostgresRepository) ListNotices(ctx context.Context, schoolID string) ([]*StrikeNotice, error) {
	query := `
		SELECT sn.id, sn.school_id, sn.title, sn.proclaimed_by, sn.strike_date::text,
		       sn.declaration_deadline, sn.content, sn.created_by,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS created_by_name,
		       sn.is_published, COALESCE(sn.communication_id::text, ''), sn.created_at, sn.updated_at
		FROM strike_notices sn
		LEFT JOIN users u ON sn.created_by = u.id
		WHERE (sn.school_id = NULLIF($1, '')::uuid OR $1 = '')
		ORDER BY sn.strike_date DESC, sn.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var list []*StrikeNotice
	now := time.Now()
	for rows.Next() {
		n := &StrikeNotice{}
		var commID string
		if err := rows.Scan(
			&n.ID, &n.SchoolID, &n.Title, &n.ProclaimedBy, &n.StrikeDate,
			&n.DeclarationDeadline, &n.Content, &n.CreatedBy, &n.CreatedByName,
			&n.IsPublished, &commID, &n.CreatedAt, &n.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if commID != "" {
			n.CommunicationID = &commID
		}
		n.IsExpired = now.After(n.DeclarationDeadline)
		list = append(list, n)
	}
	if list == nil {
		list = []*StrikeNotice{}
	}
	return list, rows.Err()
}

func (r *PostgresRepository) DeleteNotice(ctx context.Context, id, schoolID string) error {
	query := `DELETE FROM strike_notices WHERE id = $1::uuid AND (school_id = NULLIF($2, '')::uuid OR $2 = '')`
	_, err := r.db.ExecContext(ctx, query, id, schoolID)
	return err
}

func (r *PostgresRepository) UpsertDeclaration(ctx context.Context, d *StrikeDeclaration) error {
	d.DeclaredAt = time.Now()
	query := `
		INSERT INTO strike_declarations (
			strike_notice_id, user_id, intention, declared_at, ip_address, notes
		) VALUES (
			$1::uuid, $2::uuid, $3, $4, $5, $6
		)
		ON CONFLICT (strike_notice_id, user_id)
		DO UPDATE SET
			intention = EXCLUDED.intention,
			declared_at = EXCLUDED.declared_at,
			ip_address = EXCLUDED.ip_address,
			notes = EXCLUDED.notes
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		d.StrikeNoticeID, d.UserID, string(d.Intention), d.DeclaredAt, d.IPAddress, d.Notes,
	).Scan(&d.ID)
}

func (r *PostgresRepository) GetDeclaration(ctx context.Context, noticeID, userID string) (*StrikeDeclaration, error) {
	query := `
		SELECT id, strike_notice_id, user_id, intention, declared_at, COALESCE(ip_address, ''), COALESCE(notes, '')
		FROM strike_declarations
		WHERE strike_notice_id = $1::uuid AND user_id = $2::uuid
	`
	d := &StrikeDeclaration{}
	err := r.db.QueryRowContext(ctx, query, noticeID, userID).Scan(
		&d.ID, &d.StrikeNoticeID, &d.UserID, &d.Intention, &d.DeclaredAt, &d.IPAddress, &d.Notes,
	)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func roundPct(count, total int) float64 {
	if total <= 0 {
		return 0.0
	}
	pct := (float64(count) / float64(total)) * 100.0
	return math.Round(pct*10) / 10
}

func (r *PostgresRepository) GetNoticeSummary(ctx context.Context, noticeID, schoolID string) (*StrikeNoticeSummaryResponse, error) {
	notice, err := r.GetNoticeByID(ctx, noticeID, schoolID)
	if err != nil {
		return nil, err
	}

	// Fetch all staff members (teachers, ATA roles) for the school
	// Roles included: teacher, coordinator, dsga, assistente_amministrativo, collaboratore_ds, collaboratore_scolastico, secretary
	queryStaff := `
		SELECT u.id, u.first_name, u.last_name, u.email, u.role,
		       COALESCE(sd.intention, 'unanswered') AS intention,
		       sd.declared_at,
		       COALESCE(sd.ip_address, '') AS ip_address,
		       COALESCE(sd.notes, '') AS notes
		FROM users u
		LEFT JOIN strike_declarations sd ON sd.user_id = u.id AND sd.strike_notice_id = $1::uuid
		WHERE (u.school_id = NULLIF($2, '')::uuid OR $2 = '')
		  AND u.role IN (
		      'teacher', 'coordinator', 'dsga', 'assistente_amministrativo',
		      'collaboratore_ds', 'collaboratore_scolastico', 'secretary'
		  )
		ORDER BY u.last_name ASC, u.first_name ASC
	`
	rows, err := r.db.QueryContext(ctx, queryStaff, noticeID, notice.SchoolID)
	if err != nil {
		return nil, fmt.Errorf("GetNoticeSummary queryStaff: %w", err)
	}
	defer func() { _ = rows.Close() }()

	roleNames := map[string]string{
		"teacher":                   "Docente",
		"coordinator":               "Coordinatore",
		"dsga":                      "DSGA",
		"assistente_amministrativo": "Assistente Amministrativo",
		"collaboratore_ds":          "Collaboratore D.S.",
		"collaboratore_scolastico":  "Collaboratore Scolastico",
		"secretary":                 "Segreteria",
	}

	roleStatsMap := make(map[string]*RoleDeclarationSummary)
	var declarations []NominativeDeclarationItem

	totalStaff := 0
	partCount := 0
	notPartCount := 0
	undecidedCount := 0
	unansweredCount := 0

	for rows.Next() {
		var item NominativeDeclarationItem
		var declaredAt sql.NullTime
		if err := rows.Scan(
			&item.UserID, &item.FirstName, &item.LastName, &item.Email, &item.Role,
			&item.Intention, &declaredAt, &item.IPAddress, &item.Notes,
		); err != nil {
			return nil, err
		}

		if declaredAt.Valid {
			item.DeclaredAt = &declaredAt.Time
		}

		item.RoleDisplay = roleNames[item.Role]
		if item.RoleDisplay == "" {
			item.RoleDisplay = item.Role
		}

		totalStaff++
		switch item.Intention {
		case string(IntentionParticipates):
			partCount++
		case string(IntentionNotParticipates):
			notPartCount++
		case string(IntentionUndecided):
			undecidedCount++
		default:
			item.Intention = string(IntentionUnanswered)
			unansweredCount++
		}

		// Role aggregation
		rs, exists := roleStatsMap[item.Role]
		if !exists {
			rs = &RoleDeclarationSummary{
				Role:        item.Role,
				RoleDisplay: item.RoleDisplay,
			}
			roleStatsMap[item.Role] = rs
		}
		rs.Total++
		switch item.Intention {
		case string(IntentionParticipates):
			rs.Participates++
		case string(IntentionNotParticipates):
			rs.NotParticipates++
		case string(IntentionUndecided):
			rs.Undecided++
		default:
			rs.Unanswered++
		}

		declarations = append(declarations, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetNoticeSummary rows error: %w", err)
	}

	roleOrder := []string{
		"teacher", "coordinator", "dsga", "assistente_amministrativo",
		"collaboratore_ds", "collaboratore_scolastico", "secretary",
	}
	var byRole []RoleDeclarationSummary
	for _, rKey := range roleOrder {
		if rs, ok := roleStatsMap[rKey]; ok {
			rs.Percent = roundPct(rs.Participates, rs.Total)
			byRole = append(byRole, *rs)
		}
	}

	resp := &StrikeNoticeSummaryResponse{
		Notice:                 notice,
		TotalStaff:             totalStaff,
		TotalAnswered:          partCount + notPartCount + undecidedCount,
		ParticipatesCount:      partCount,
		NotParticipatesCount:   notPartCount,
		UndecidedCount:         undecidedCount,
		UnansweredCount:        unansweredCount,
		ParticipatesPercent:    roundPct(partCount, totalStaff),
		NotParticipatesPercent: roundPct(notPartCount, totalStaff),
		UndecidedPercent:       roundPct(undecidedCount, totalStaff),
		UnansweredPercent:      roundPct(unansweredCount, totalStaff),
		ByRole:                 byRole,
		Declarations:           declarations,
	}

	if resp.Declarations == nil {
		resp.Declarations = []NominativeDeclarationItem{}
	}
	if resp.ByRole == nil {
		resp.ByRole = []RoleDeclarationSummary{}
	}

	return resp, nil
}

func (r *PostgresRepository) CreateBachecaCommunication(ctx context.Context, schoolID, senderID, title, content string, deadline time.Time) (string, error) {
	commID := uuid.New().String()
	now := time.Now()

	bodyWithDeadline := fmt.Sprintf("%s\n\n[AVVISO SCIOPERO]: Si prega di comunicare la propria intenzione preventiva (Aderisco / Non aderisco / Non ho ancora deciso) entro il %s.",
		content, deadline.Format("02/01/2006 15:04"),
	)

	// Receiver IDs: include staff group or all staff
	query := `
		INSERT INTO communications (
			id, school_id, sender_id, receiver_ids, subject, body, type, requires_signature, signature_deadline, created_at
		) VALUES (
			$1::uuid, $2::uuid, $3::uuid, $4, $5, $6, 'notice', false, $7, $8
		)
	`
	recipients := pq.StringArray{"staff", "teachers", "ata"}
	_, err := r.db.ExecContext(ctx, query,
		commID, schoolID, senderID, recipients, title, bodyWithDeadline, deadline, now,
	)
	if err != nil {
		return "", err
	}
	return commID, nil
}
