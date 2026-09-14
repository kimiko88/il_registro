package staff_attendance

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Repository definisce l'interfaccia del repository
type Repository interface {
	// Presenze staff
	GetDailySummary(ctx context.Context, schoolID, date string) (*DailyStaffSummary, error)
	ListByDate(ctx context.Context, schoolID, date string) ([]StaffAttendance, error)
	Upsert(ctx context.Context, schoolID, recordedBy string, req UpsertStaffAttendanceRequest) (*StaffAttendance, error)
	BulkUpsert(ctx context.Context, schoolID, recordedBy string, req BulkUpsertRequest) ([]StaffAttendance, error)
	GetByUserAndDate(ctx context.Context, schoolID, userID, date string) (*StaffAttendance, error)
	Delete(ctx context.Context, schoolID, id string) error
	SetStrikeMode(ctx context.Context, schoolID, actorID, date string, isStrikeDay bool) error

	// Badge
	RegisterBadgeSwipe(ctx context.Context, schoolID string, req BadgeSwipeRequest) (*BadgeSwipe, error)
	ProcessPendingSwipes(ctx context.Context, schoolID string) (int, error)
	AssignBadge(ctx context.Context, badge UserBadge) error
	RevokeBadge(ctx context.Context, schoolID, userID, badgeCode string) error
	ListBadges(ctx context.Context, schoolID string) ([]UserBadge, error)

	// Ferie & Permessi
	CreateLeaveRequest(ctx context.Context, schoolID, userID string, req CreateLeaveRequest) (*LeaveRequest, error)
	ListLeaveRequests(ctx context.Context, schoolID, userID, status string) ([]LeaveRequest, error)
	GetLeaveRequest(ctx context.Context, schoolID, id string) (*LeaveRequest, error)
	ApproveLeaveRequest(ctx context.Context, schoolID, id, approvedBy, notes string) error
	RejectLeaveRequest(ctx context.Context, schoolID, id, rejectedBy, reason string) error
	DeleteLeaveRequest(ctx context.Context, schoolID, id, userID string) error

	// Cartellino mensile
	GetMonthlyTimecard(ctx context.Context, schoolID, userID, month string) (*MonthlyTimecard, error)
	GetAllMonthlyTimecards(ctx context.Context, schoolID, month string) ([]MonthlyTimecard, error)
}

// PostgresRepository implementa Repository
type PostgresRepository struct {
	db *sql.DB
}

// NewRepository crea un nuovo repository PostgreSQL
func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

// GetDailySummary restituisce il riepilogo presenze per una data
func (r *PostgresRepository) GetDailySummary(ctx context.Context, schoolID, date string) (*DailyStaffSummary, error) {
	staffList, err := r.ListByDate(ctx, schoolID, date)
	if err != nil {
		return nil, fmt.Errorf("GetDailySummary: %w", err)
	}

	summary := &DailyStaffSummary{
		Date:      date,
		StaffList: staffList,
		ByRole:    []RoleSummary{},
	}

	// Aggregazione per ruolo
	roleMap := make(map[string]*RoleSummary)
	for _, s := range staffList {
		if summary.IsStrikeDay || s.IsStrikeDay {
			summary.IsStrikeDay = true
		}

		rs, ok := roleMap[s.Role]
		if !ok {
			roleMap[s.Role] = &RoleSummary{
				Role:        s.Role,
				RoleDisplay: RoleDisplayNames[s.Role],
			}
			rs = roleMap[s.Role]
		}
		rs.Total++

		switch s.Status {
		case StatusPresent:
			rs.Present++
		case StatusAbsent:
			rs.Absent++
		case StatusLate:
			rs.Late++
			rs.Present++ // in ritardo = comunque presente
		case StatusMission:
			rs.OnMission++
		case StatusSickLeave, StatusPermit:
			rs.OnLeave++
		case StatusOnStrike:
			rs.OnStrike++
		}
	}

	// Ordine predefinito dei ruoli nella dashboard
	roleOrder := []string{
		"teacher", "coordinator", "dsga", "assistente_amministrativo",
		"collaboratore_ds", "collaboratore_scolastico", "assistente_tecnico",
		"assistente_alunni", "assistente_personale", "assistente_contabilita",
		"assistente_protocollo", "assistente_sportello", "responsabile_servizio",
		"secretary", "principal", "vice_principal",
	}

	var totalRole RoleSummary
	totalRole.Role = "total"
	totalRole.RoleDisplay = "Tutto il Personale"

	for _, role := range roleOrder {
		if rs, ok := roleMap[role]; ok {
			summary.ByRole = append(summary.ByRole, *rs)
			totalRole.Total += rs.Total
			totalRole.Present += rs.Present
			totalRole.Absent += rs.Absent
			totalRole.Late += rs.Late
			totalRole.OnMission += rs.OnMission
			totalRole.OnLeave += rs.OnLeave
			totalRole.OnStrike += rs.OnStrike
		}
	}

	// Aggiungi eventuali ruoli non nell'ordine predefinito
	for role, rs := range roleMap {
		found := false
		for _, r := range roleOrder {
			if r == role {
				found = true
				break
			}
		}
		if !found {
			summary.ByRole = append(summary.ByRole, *rs)
		}
	}

	summary.TotalStaff = totalRole

	// Se non ancora rilevato, controlla se la giornata ha un avviso di sciopero attivo o se e segnata come sciopero
	if !summary.IsStrikeDay {
		var hasNotice bool
		_ = r.db.QueryRowContext(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM strike_notices
				WHERE (school_id = NULLIF($1, '')::uuid OR $1 = '')
				  AND strike_date = $2::date
				  AND is_published = true
			)
		`, schoolID, date).Scan(&hasNotice)
		if hasNotice {
			summary.IsStrikeDay = true
		}
	}
	if !summary.IsStrikeDay {
		var hasStrikeDay bool
		_ = r.db.QueryRowContext(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM staff_attendance
				WHERE (school_id = NULLIF($1, '')::uuid OR $1 = '')
				  AND date = $2::date
				  AND is_strike_day = true
			)
		`, schoolID, date).Scan(&hasStrikeDay)
		if hasStrikeDay {
			summary.IsStrikeDay = true
		}
	}

	return summary, nil
}

// ListByDate restituisce la lista nominativa del personale per una data
func (r *PostgresRepository) ListByDate(ctx context.Context, schoolID, date string) ([]StaffAttendance, error) {
	query := `
		SELECT
			sa.id, sa.school_id, sa.user_id, sa.date::text, sa.status,
			sa.badge_entry_time, sa.badge_exit_time, sa.badge_device_id, sa.badge_imported_at,
			sa.is_strike_day, sa.strike_confirmed_by, sa.notes, sa.recorded_by,
			sa.created_at, sa.updated_at,
			u.first_name, u.last_name, u.email, u.role,
			COALESCE(ub.badge_code, '') AS badge_code
		FROM staff_attendance sa
		JOIN users u ON sa.user_id = u.id
		LEFT JOIN LATERAL (
			SELECT badge_code FROM user_badges
			WHERE user_id = u.id AND is_active = true
			ORDER BY assigned_at DESC LIMIT 1
		) ub ON true
		WHERE sa.school_id = $1 AND sa.date = $2::date
		ORDER BY
			CASE u.role
				WHEN 'principal' THEN 1
				WHEN 'vice_principal' THEN 2
				WHEN 'dsga' THEN 3
				WHEN 'collaboratore_ds' THEN 4
				WHEN 'assistente_amministrativo' THEN 5
				WHEN 'assistente_tecnico' THEN 6
				WHEN 'collaboratore_scolastico' THEN 7
				WHEN 'secretary' THEN 8
				WHEN 'teacher' THEN 9
				WHEN 'coordinator' THEN 10
				ELSE 11
			END,
			u.last_name, u.first_name
	`

	rows, err := r.db.QueryContext(ctx, query, schoolID, date)
	if err != nil {
		return nil, fmt.Errorf("ListByDate query: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var results []StaffAttendance
	for rows.Next() {
		var sa StaffAttendance
		var badgeEntry, badgeExit, badgeImported sql.NullTime
		var badgeDevice, strikeConfBy, recordedBy, notes, badgeCode sql.NullString

		err := rows.Scan(
			&sa.ID, &sa.SchoolID, &sa.UserID, &sa.Date, &sa.Status,
			&badgeEntry, &badgeExit, &badgeDevice, &badgeImported,
			&sa.IsStrikeDay, &strikeConfBy, &notes, &recordedBy,
			&sa.CreatedAt, &sa.UpdatedAt,
			&sa.FirstName, &sa.LastName, &sa.Email, &sa.Role,
			&badgeCode,
		)
		if err != nil {
			return nil, fmt.Errorf("ListByDate scan: %w", err)
		}

		if badgeEntry.Valid {
			sa.BadgeEntryTime = &badgeEntry.Time
		}
		if badgeExit.Valid {
			sa.BadgeExitTime = &badgeExit.Time
		}
		if badgeImported.Valid {
			sa.BadgeImportedAt = &badgeImported.Time
		}
		if badgeDevice.Valid {
			sa.BadgeDeviceID = &badgeDevice.String
		}
		if strikeConfBy.Valid {
			sa.StrikeConfirmedBy = &strikeConfBy.String
		}
		if notes.Valid {
			sa.Notes = notes.String
		}
		if recordedBy.Valid {
			sa.RecordedBy = &recordedBy.String
		}
		if badgeCode.Valid {
			sa.BadgeCode = badgeCode.String
		}

		sa.RoleDisplay = RoleDisplayNames[sa.Role]
		results = append(results, sa)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListByDate rows: %w", err)
	}

	if results == nil {
		results = []StaffAttendance{}
	}
	return results, nil
}

func parseAttendanceTime(dateStr, timeStr string) (*time.Time, error) {
	if timeStr == "" {
		return nil, nil
	}
	if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
		return &t, nil
	}
	parts := strings.Split(timeStr, ":")
	if len(parts) >= 2 {
		hour, errH := strconv.Atoi(strings.TrimSpace(parts[0]))
		min, errM := strconv.Atoi(strings.TrimSpace(parts[1]))
		sec := 0
		if len(parts) >= 3 {
			sec, _ = strconv.Atoi(strings.TrimSpace(parts[2]))
		}
		if errH == nil && errM == nil {
			d, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				t := time.Date(d.Year(), d.Month(), d.Day(), hour, min, sec, 0, time.Local)
				return &t, nil
			}
		}
	}
	return nil, fmt.Errorf("formato orario non valido: %s", timeStr)
}

// Upsert inserisce o aggiorna una presenza staff
func (r *PostgresRepository) Upsert(ctx context.Context, schoolID, recordedBy string, req UpsertStaffAttendanceRequest) (*StaffAttendance, error) {
	var entryTime, exitTime *time.Time
	if req.BadgeEntryTime != nil {
		entryTime = req.BadgeEntryTime
	} else if req.EntryTime != nil && *req.EntryTime != "" {
		entryTime, _ = parseAttendanceTime(req.Date, *req.EntryTime)
	}

	if req.BadgeExitTime != nil {
		exitTime = req.BadgeExitTime
	} else if req.ExitTime != nil && *req.ExitTime != "" {
		exitTime, _ = parseAttendanceTime(req.Date, *req.ExitTime)
	}

	query := `
		INSERT INTO staff_attendance (
			school_id, user_id, date, status, notes, is_strike_day, recorded_by,
			badge_entry_time, badge_exit_time
		)
		VALUES ($1, $2, $3::date, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (school_id, user_id, date)
		DO UPDATE SET
			status = EXCLUDED.status,
			notes = EXCLUDED.notes,
			is_strike_day = EXCLUDED.is_strike_day,
			recorded_by = EXCLUDED.recorded_by,
			badge_entry_time = CASE 
				WHEN EXCLUDED.badge_entry_time IS NOT NULL THEN EXCLUDED.badge_entry_time 
				ELSE staff_attendance.badge_entry_time 
			END,
			badge_exit_time = CASE 
				WHEN EXCLUDED.badge_exit_time IS NOT NULL THEN EXCLUDED.badge_exit_time 
				ELSE staff_attendance.badge_exit_time 
			END,
			updated_at = NOW()
		RETURNING id, school_id, user_id, date::text, status, is_strike_day, notes, recorded_by,
		          badge_entry_time, badge_exit_time, created_at, updated_at
	`

	var sa StaffAttendance
	var notes, recBy sql.NullString
	var bEntry, bExit sql.NullTime
	err := r.db.QueryRowContext(ctx, query,
		schoolID, req.UserID, req.Date, string(req.Status),
		req.Notes, req.IsStrikeDay, recordedBy,
		entryTime, exitTime,
	).Scan(
		&sa.ID, &sa.SchoolID, &sa.UserID, &sa.Date, &sa.Status,
		&sa.IsStrikeDay, &notes, &recBy,
		&bEntry, &bExit,
		&sa.CreatedAt, &sa.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("Upsert: %w", err)
	}

	if notes.Valid {
		sa.Notes = notes.String
	}
	if recBy.Valid {
		sa.RecordedBy = &recBy.String
	}
	if bEntry.Valid {
		sa.BadgeEntryTime = &bEntry.Time
	}
	if bExit.Valid {
		sa.BadgeExitTime = &bExit.Time
	}

	return &sa, nil
}

// BulkUpsert inserisce/aggiorna più presenze in un'unica transazione
func (r *PostgresRepository) BulkUpsert(ctx context.Context, schoolID, recordedBy string, req BulkUpsertRequest) ([]StaffAttendance, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("BulkUpsert begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	query := `
		INSERT INTO staff_attendance (school_id, user_id, date, status, notes, is_strike_day, recorded_by)
		VALUES ($1, $2, $3::date, $4, $5, $6, $7)
		ON CONFLICT (school_id, user_id, date)
		DO UPDATE SET
			status = EXCLUDED.status,
			notes = EXCLUDED.notes,
			is_strike_day = EXCLUDED.is_strike_day,
			recorded_by = EXCLUDED.recorded_by,
			updated_at = NOW()
		RETURNING id, school_id, user_id, date::text, status, is_strike_day, notes, recorded_by, created_at, updated_at
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("BulkUpsert prepare: %w", err)
	}
	defer func() { _ = stmt.Close() }()

	var results []StaffAttendance
	for _, att := range req.Attendances {
		var sa StaffAttendance
		var notes, recBy sql.NullString
		err := stmt.QueryRowContext(ctx,
			schoolID, att.UserID, req.Date, string(att.Status),
			att.Notes, req.IsStrikeDay, recordedBy,
		).Scan(
			&sa.ID, &sa.SchoolID, &sa.UserID, &sa.Date, &sa.Status,
			&sa.IsStrikeDay, &notes, &recBy, &sa.CreatedAt, &sa.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("BulkUpsert row %s: %w", att.UserID, err)
		}

		if notes.Valid {
			sa.Notes = notes.String
		}
		if recBy.Valid {
			sa.RecordedBy = &recBy.String
		}
		results = append(results, sa)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("BulkUpsert commit: %w", err)
	}

	return results, nil
}

// GetByUserAndDate restituisce la presenza di un utente per una data
func (r *PostgresRepository) GetByUserAndDate(ctx context.Context, schoolID, userID, date string) (*StaffAttendance, error) {
	query := `
		SELECT sa.id, sa.school_id, sa.user_id, sa.date::text, sa.status,
		       sa.badge_entry_time, sa.badge_exit_time, sa.is_strike_day, sa.notes, sa.recorded_by,
		       sa.created_at, sa.updated_at,
		       u.first_name, u.last_name, u.email, u.role
		FROM staff_attendance sa
		JOIN users u ON sa.user_id = u.id
		WHERE sa.school_id = $1 AND sa.user_id = $2 AND sa.date = $3::date
	`

	var sa StaffAttendance
	var badgeEntry, badgeExit sql.NullTime
	var notes, recBy sql.NullString

	err := r.db.QueryRowContext(ctx, query, schoolID, userID, date).Scan(
		&sa.ID, &sa.SchoolID, &sa.UserID, &sa.Date, &sa.Status,
		&badgeEntry, &badgeExit, &sa.IsStrikeDay, &notes, &recBy,
		&sa.CreatedAt, &sa.UpdatedAt,
		&sa.FirstName, &sa.LastName, &sa.Email, &sa.Role,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetByUserAndDate: %w", err)
	}

	if badgeEntry.Valid {
		sa.BadgeEntryTime = &badgeEntry.Time
	}
	if badgeExit.Valid {
		sa.BadgeExitTime = &badgeExit.Time
	}
	if notes.Valid {
		sa.Notes = notes.String
	}
	if recBy.Valid {
		sa.RecordedBy = &recBy.String
	}
	sa.RoleDisplay = RoleDisplayNames[sa.Role]
	return &sa, nil
}

// Delete rimuove una presenza staff
func (r *PostgresRepository) Delete(ctx context.Context, schoolID, id string) error {
	res, err := r.db.ExecContext(ctx,
		"DELETE FROM staff_attendance WHERE id = $1 AND school_id = $2",
		id, schoolID,
	)
	if err != nil {
		return fmt.Errorf("Delete: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("staff attendance record not found")
	}
	return nil
}

// RegisterBadgeSwipe registra una timbratura badge grezza
func (r *PostgresRepository) RegisterBadgeSwipe(ctx context.Context, schoolID string, req BadgeSwipeRequest) (*BadgeSwipe, error) {
	swipeTime := time.Now()
	if req.SwipeTime != "" {
		if t, err := time.Parse(time.RFC3339, req.SwipeTime); err == nil {
			swipeTime = t
		} else if t, err := parseAttendanceTime(req.Date, req.SwipeTime); err == nil && t != nil {
			swipeTime = *t
		}
	} else if req.Date != "" {
		now := time.Now()
		d, err := time.Parse("2006-01-02", req.Date)
		if err == nil {
			swipeTime = time.Date(d.Year(), d.Month(), d.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.Local)
		}
	}

	swipeType := req.SwipeType
	if swipeType == "" {
		swipeType = "in"
	}

	// Cerca l'utente associato al badge
	var userID sql.NullString
	if req.UserID != "" {
		userID = sql.NullString{String: req.UserID, Valid: true}
	} else {
		_ = r.db.QueryRowContext(ctx,
			"SELECT user_id FROM user_badges WHERE school_id = $1 AND badge_code = $2 AND is_active = true",
			schoolID, req.BadgeCode,
		).Scan(&userID)
		if !userID.Valid {
			_ = r.db.QueryRowContext(ctx,
				"SELECT user_id FROM user_badges WHERE badge_code = $1 AND is_active = true",
				req.BadgeCode,
			).Scan(&userID)
		}
	}

	// Se l'utente è noto e il badge non è registrato in user_badges, inseriscilo per le letture future
	if userID.Valid && req.BadgeCode != "" {
		_, _ = r.db.ExecContext(ctx, `
			INSERT INTO user_badges (school_id, user_id, badge_code, is_active, assigned_at)
			VALUES ($1, $2, $3, true, NOW())
			ON CONFLICT DO NOTHING
		`, schoolID, userID.String, req.BadgeCode)
	}

	var rawDataJSON *string
	if req.RawData != "" {
		// Valida che sia JSON valido
		if json.Valid([]byte(req.RawData)) {
			rawDataJSON = &req.RawData
		} else {
			escaped := `"` + strings.ReplaceAll(req.RawData, `"`, `\"`) + `"`
			rawDataJSON = &escaped
		}
	}

	query := `
		INSERT INTO badge_swipes (school_id, user_id, badge_code, device_id, swipe_time, swipe_type, raw_data)
		VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb)
		RETURNING id, school_id, user_id, badge_code, device_id, swipe_time, swipe_type, processed, created_at
	`

	var sw BadgeSwipe
	var uid sql.NullString
	err := r.db.QueryRowContext(ctx, query,
		schoolID, userID, req.BadgeCode, req.DeviceID, swipeTime, swipeType, rawDataJSON,
	).Scan(
		&sw.ID, &sw.SchoolID, &uid, &sw.BadgeCode,
		&sw.DeviceID, &sw.SwipeTime, &sw.SwipeType, &sw.Processed, &sw.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("RegisterBadgeSwipe: %w", err)
	}

	if uid.Valid {
		sw.UserID = &uid.String
	}
	return &sw, nil
}

// ProcessPendingSwipes elabora le timbrature non ancora processate e aggiorna staff_attendance
func (r *PostgresRepository) ProcessPendingSwipes(ctx context.Context, schoolID string) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("ProcessPendingSwipes begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Recupera le timbrature non processate con utente associato
	rows, err := tx.QueryContext(ctx, `
		SELECT id, user_id, swipe_time, swipe_type
		FROM badge_swipes
		WHERE school_id = $1 AND processed = false AND user_id IS NOT NULL
		ORDER BY swipe_time
	`, schoolID)
	if err != nil {
		return 0, fmt.Errorf("ProcessPendingSwipes query: %w", err)
	}

	type swipeRow struct {
		id        string
		userID    string
		swipeTime time.Time
		swipeType string
	}
	var swipes []swipeRow

	for rows.Next() {
		var sw swipeRow
		var uid sql.NullString
		if err := rows.Scan(&sw.id, &uid, &sw.swipeTime, &sw.swipeType); err != nil {
			_ = rows.Close()
			return 0, fmt.Errorf("ProcessPendingSwipes scan: %w", err)
		}
		if uid.Valid {
			sw.userID = uid.String
			swipes = append(swipes, sw)
		}
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return 0, fmt.Errorf("ProcessPendingSwipes rows: %w", err)
	}
	_ = rows.Close()

	processed := 0
	for _, sw := range swipes {
		dateStr := sw.swipeTime.Format("2006-01-02")

		if sw.swipeType == "in" || sw.swipeType == "break_in" {
			_, err = tx.ExecContext(ctx, `
				INSERT INTO staff_attendance (school_id, user_id, date, status, badge_entry_time, badge_device_id)
				VALUES ($1, $2, $3::date, 'present', $4, 'badge')
				ON CONFLICT (school_id, user_id, date)
				DO UPDATE SET
					badge_entry_time = CASE WHEN staff_attendance.badge_entry_time IS NULL THEN EXCLUDED.badge_entry_time
					                        ELSE LEAST(staff_attendance.badge_entry_time, EXCLUDED.badge_entry_time) END,
					status = CASE WHEN staff_attendance.status = 'absent' THEN 'present' ELSE staff_attendance.status END,
					updated_at = NOW()
			`, schoolID, sw.userID, dateStr, sw.swipeTime)
		} else {
			// out / break_out — aggiorna orario uscita
			_, err = tx.ExecContext(ctx, `
				INSERT INTO staff_attendance (school_id, user_id, date, status, badge_exit_time, badge_device_id)
				VALUES ($1, $2, $3::date, 'present', $4, 'badge')
				ON CONFLICT (school_id, user_id, date)
				DO UPDATE SET
					badge_exit_time = COALESCE(GREATEST(staff_attendance.badge_exit_time, EXCLUDED.badge_exit_time), EXCLUDED.badge_exit_time),
					updated_at = NOW()
			`, schoolID, sw.userID, dateStr, sw.swipeTime)
		}
		if err != nil {
			return processed, fmt.Errorf("ProcessPendingSwipes upsert for %s: %w", sw.userID, err)
		}

		// Marca come processata
		_, err = tx.ExecContext(ctx,
			"UPDATE badge_swipes SET processed = true, processed_at = NOW() WHERE id = $1",
			sw.id,
		)
		if err != nil {
			return processed, fmt.Errorf("ProcessPendingSwipes mark processed: %w", err)
		}
		processed++
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("ProcessPendingSwipes commit: %w", err)
	}
	return processed, nil
}

// AssignBadge associa un badge a un utente
func (r *PostgresRepository) AssignBadge(ctx context.Context, badge UserBadge) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO user_badges (school_id, user_id, badge_code, notes)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (school_id, badge_code)
		DO UPDATE SET user_id = EXCLUDED.user_id, is_active = true, notes = EXCLUDED.notes, assigned_at = NOW(), revoked_at = NULL
	`, badge.SchoolID, badge.UserID, badge.BadgeCode, badge.Notes)
	if err != nil {
		return fmt.Errorf("AssignBadge: %w", err)
	}
	return nil
}

// RevokeBadge revoca un badge
func (r *PostgresRepository) RevokeBadge(ctx context.Context, schoolID, userID, badgeCode string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE user_badges SET is_active = false, revoked_at = NOW()
		WHERE school_id = $1 AND user_id = $2 AND badge_code = $3
	`, schoolID, userID, badgeCode)
	if err != nil {
		return fmt.Errorf("RevokeBadge: %w", err)
	}
	return nil
}

// ListBadges elenca i badge attivi per una scuola
func (r *PostgresRepository) ListBadges(ctx context.Context, schoolID string) ([]UserBadge, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, school_id, user_id, badge_code, is_active, assigned_at, revoked_at, notes
		FROM user_badges
		WHERE school_id = $1
		ORDER BY assigned_at DESC
	`, schoolID)
	if err != nil {
		return nil, fmt.Errorf("ListBadges: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var badges []UserBadge
	for rows.Next() {
		var b UserBadge
		var revoked sql.NullTime
		var notes sql.NullString
		if err := rows.Scan(&b.ID, &b.SchoolID, &b.UserID, &b.BadgeCode,
			&b.IsActive, &b.AssignedAt, &revoked, &notes); err != nil {
			return nil, fmt.Errorf("ListBadges scan: %w", err)
		}
		if revoked.Valid {
			b.RevokedAt = &revoked.Time
		}
		if notes.Valid {
			b.Notes = notes.String
		}
		badges = append(badges, b)
	}
	return badges, rows.Err()
}

// SetStrikeMode imposta o revoca la modalità sciopero per una determinata data
func (r *PostgresRepository) SetStrikeMode(ctx context.Context, schoolID, actorID, date string, isStrikeDay bool) error {
	var schoolUUID *string
	if s := strings.TrimSpace(schoolID); s != "" {
		schoolUUID = &s
	}
	var actorUUID *string
	if a := strings.TrimSpace(actorID); a != "" {
		actorUUID = &a
	}

	if isStrikeDay {
		// Aggiorna le presenze esistenti o inserisce una riga per il personale attivo indicando is_strike_day = true
		upsertQuery := `
			INSERT INTO staff_attendance (school_id, user_id, date, status, is_strike_day, recorded_by, created_at, updated_at)
			SELECT COALESCE($1::uuid, u.school_id), u.id, $2::date, 'absent', true, $3::uuid, NOW(), NOW()
			FROM users u
			WHERE ($1::uuid IS NULL OR u.school_id = $1::uuid)
			  AND (u.school_id IS NOT NULL OR $1::uuid IS NOT NULL)
			  AND u.is_active = true
			  AND u.role IN ('teacher', 'coordinator', 'dsga', 'collaboratore_ds', 'assistente_amministrativo', 'assistente_tecnico', 'collaboratore_scolastico', 'collaboratore_mensa', 'assistente_alunni', 'assistente_personale', 'assistente_contabilita', 'assistente_protocollo', 'assistente_sportello', 'responsabile_servizio', 'secretary')
			ON CONFLICT (school_id, user_id, date)
			DO UPDATE SET is_strike_day = true, updated_at = NOW()
		`
		_, err := r.db.ExecContext(ctx, upsertQuery, schoolUUID, date, actorUUID)
		if err != nil {
			return fmt.Errorf("SetStrikeMode upsert error: %w", err)
		}
	} else {
		// Disattiva la modalità sciopero per quella data
		updateQuery := `
			UPDATE staff_attendance
			SET is_strike_day = false, updated_at = NOW()
			WHERE ($1::uuid IS NULL OR school_id = $1::uuid) AND date = $2::date
		`
		_, err := r.db.ExecContext(ctx, updateQuery, schoolUUID, date)
		if err != nil {
			return fmt.Errorf("SetStrikeMode deactivate error: %w", err)
		}
	}
	return nil
}
