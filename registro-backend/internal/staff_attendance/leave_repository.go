package staff_attendance

// leave_repository.go — implementazione dei metodi ferie/cartellino
// Si appoggia alla stessa struttura PostgresRepository già definita in repository.go

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CreateLeaveRequest registra una nuova richiesta di ferie/permesso
func (r *PostgresRepository) CreateLeaveRequest(ctx context.Context, schoolID, userID string, req CreateLeaveRequest) (*LeaveRequest, error) {
	id := uuid.New().String()
	now := time.Now()

	days := req.Days
	if days == 0 && req.Hours == 0 {
		days = 1 // default 1 giorno
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO staff_leave_requests
		    (id, school_id, user_id, type, start_date, end_date, days, hours, notes, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,'pending',$10,$10)`,
		id, schoolID, userID, req.Type, req.StartDate, req.EndDate, days, req.Hours, req.Notes, now)
	if err != nil {
		return nil, fmt.Errorf("errore creazione richiesta: %w", err)
	}

	return r.GetLeaveRequest(ctx, schoolID, id)
}

// ListLeaveRequests elenca le richieste (tutte o per utente specifico)
func (r *PostgresRepository) ListLeaveRequests(ctx context.Context, schoolID, userID, status string) ([]LeaveRequest, error) {
	query := `
		SELECT lr.id, lr.school_id, lr.user_id, lr.type, lr.start_date::text, lr.end_date::text,
		       COALESCE(lr.days, 0), COALESCE(lr.hours, 0), COALESCE(lr.notes, ''), lr.status,
		       lr.approved_by, lr.approved_at, lr.rejected_at, COALESCE(lr.reject_reason, ''),
		       lr.created_at, lr.updated_at,
		       COALESCE(u.first_name, '') AS user_first_name,
		       COALESCE(u.last_name, '')  AS user_last_name,
		       COALESCE(u.role, '')       AS user_role
		FROM staff_leave_requests lr
		JOIN users u ON u.id = lr.user_id
		WHERE 1=1`
	args := []interface{}{}
	idx := 1

	if schoolID != "" {
		query += fmt.Sprintf(" AND lr.school_id = $%d::uuid", idx)
		args = append(args, schoolID)
		idx++
	}
	if userID != "" {
		query += fmt.Sprintf(" AND lr.user_id = $%d", idx)
		args = append(args, userID)
		idx++
	}
	if status != "" {
		query += fmt.Sprintf(" AND lr.status = $%d", idx)
		args = append(args, status)
	}
	query += " ORDER BY lr.created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaves []LeaveRequest
	for rows.Next() {
		var l LeaveRequest
		err := rows.Scan(
			&l.ID, &l.SchoolID, &l.UserID, &l.Type, &l.StartDate, &l.EndDate,
			&l.Days, &l.Hours, &l.Notes, &l.Status,
			&l.ApprovedBy, &l.ApprovedAt, &l.RejectedAt, &l.RejectReason,
			&l.CreatedAt, &l.UpdatedAt,
			&l.UserFirstName, &l.UserLastName, &l.UserRole,
		)
		if err != nil {
			return nil, err
		}
		leaves = append(leaves, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return leaves, nil
}

// GetLeaveRequest recupera una singola richiesta per ID
func (r *PostgresRepository) GetLeaveRequest(ctx context.Context, schoolID, id string) (*LeaveRequest, error) {
	leaves, err := r.ListLeaveRequests(ctx, schoolID, "", "")
	if err != nil {
		return nil, err
	}
	for i := range leaves {
		if leaves[i].ID == id {
			return &leaves[i], nil
		}
	}
	return nil, fmt.Errorf("richiesta non trovata")
}

// ApproveLeaveRequest approva una richiesta
func (r *PostgresRepository) ApproveLeaveRequest(ctx context.Context, schoolID, id, approvedBy, notes string) error {
	now := time.Now()
	res, err := r.db.ExecContext(ctx, `
		UPDATE staff_leave_requests
		SET status='approved', approved_by=$1, approved_at=$2, notes=COALESCE(NULLIF($3,''),notes), updated_at=$2
		WHERE id=$4 AND school_id=$5 AND status='pending'`,
		approvedBy, now, notes, id, schoolID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("richiesta non trovata o già elaborata")
	}
	return nil
}

// RejectLeaveRequest rifiuta una richiesta
func (r *PostgresRepository) RejectLeaveRequest(ctx context.Context, schoolID, id, rejectedBy, reason string) error {
	now := time.Now()
	res, err := r.db.ExecContext(ctx, `
		UPDATE staff_leave_requests
		SET status='rejected', approved_by=$1, rejected_at=$2, reject_reason=$3, updated_at=$2
		WHERE id=$4 AND school_id=$5 AND status='pending'`,
		rejectedBy, now, reason, id, schoolID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("richiesta non trovata o già elaborata")
	}
	return nil
}

// DeleteLeaveRequest elimina una richiesta in stato pending (solo l'autore)
func (r *PostgresRepository) DeleteLeaveRequest(ctx context.Context, schoolID, id, userID string) error {
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM staff_leave_requests
		WHERE id=$1 AND school_id=$2 AND user_id=$3 AND status='pending'`,
		id, schoolID, userID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("richiesta non trovata, non sei l'autore o già elaborata")
	}
	return nil
}

// GetMonthlyTimecard calcola il cartellino mensile per un utente
func (r *PostgresRepository) GetMonthlyTimecard(ctx context.Context, schoolID, userID, month string) (*MonthlyTimecard, error) {
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	t, err := time.Parse("2006-01", month)
	if err != nil {
		t = time.Now()
		month = t.Format("2006-01")
	}
	startDate := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	endDate := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")

	// Recupera dati utente e badge attivo in una singola query
	var firstName, lastName, role, badgeCode string
	var userSchoolID sql.NullString
	err = r.db.QueryRowContext(ctx, `
		SELECT 
			u.first_name, 
			u.last_name, 
			u.role, 
			u.school_id::text,
			COALESCE(ub.badge_code, '')
		FROM users u
		LEFT JOIN LATERAL (
			SELECT badge_code 
			FROM user_badges 
			WHERE user_id = u.id AND is_active = true 
			ORDER BY assigned_at DESC 
			LIMIT 1
		) ub ON true
		WHERE u.id = $1`, userID).Scan(&firstName, &lastName, &role, &userSchoolID, &badgeCode)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("utente non trovato")
	}
	if err != nil {
		return nil, err
	}

	if schoolID == "" && userSchoolID.Valid && userSchoolID.String != "" {
		schoolID = userSchoolID.String
	}

	// Dettaglio giornaliero per il cartellino con index range scan
	var dailyEntries []DailyTimecardEntry
	var totalWorkedMinutes float64

	dRows, err := r.db.QueryContext(ctx, `
		SELECT 
			date::text, 
			badge_entry_time, 
			badge_exit_time,
			COALESCE(EXTRACT(EPOCH FROM (COALESCE(badge_exit_time, NOW()) - badge_entry_time)) / 60, 0)::int as worked_minutes,
			status, 
			COALESCE(notes, '')
		FROM staff_attendance
		WHERE ($1 = '' OR school_id = NULLIF($1, '')::uuid) 
		  AND user_id = $2
		  AND date >= $3::date AND date < $4::date
		ORDER BY date ASC`, schoolID, userID, startDate, endDate)
	if err == nil {
		defer dRows.Close()
		for dRows.Next() {
			var de DailyTimecardEntry
			var entryTime, exitTime sql.NullTime
			if err := dRows.Scan(&de.Date, &entryTime, &exitTime, &de.WorkedMinutes, &de.Status, &de.Notes); err == nil {
				if entryTime.Valid {
					de.EntryTime = &entryTime.Time
				}
				if exitTime.Valid {
					de.ExitTime = &exitTime.Time
				}
				dailyEntries = append(dailyEntries, de)
				totalWorkedMinutes += float64(de.WorkedMinutes)
			}
		}
		if err := dRows.Err(); err != nil {
			return nil, err
		}
	}
	if dailyEntries == nil {
		dailyEntries = []DailyTimecardEntry{}
	}
	workedHours := totalWorkedMinutes / 60.0

	// Conteggi assenze per tipo con index range scan
	var absenceDays, leaveDays, sickDays int
	var permitHours float64

	lRows, err := r.db.QueryContext(ctx, `
		SELECT type, COUNT(*) as cnt, COALESCE(SUM(hours),0) as hrs
		FROM staff_leave_requests
		WHERE ($1 = '' OR school_id = NULLIF($1, '')::uuid) 
		  AND user_id = $2
		  AND start_date >= $3::date AND start_date < $4::date
		  AND status = 'approved'
		GROUP BY type`, schoolID, userID, startDate, endDate)
	if err == nil {
		defer lRows.Close()
		for lRows.Next() {
			var ltype string
			var cnt int
			var hrs float64
			if err := lRows.Scan(&ltype, &cnt, &hrs); err == nil {
				switch LeaveType(ltype) {
				case LeaveTypeFerie:
					leaveDays += cnt
				case LeaveTypeMalattia:
					sickDays += cnt
				case LeaveTypePermesso, LeaveTypePermessoStudio, LeaveTypeRecupero:
					permitHours += hrs
					if hrs == 0 {
						permitHours += float64(cnt) * 8.0 // default 8h per giorno
					}
				default:
					absenceDays += cnt
				}
			}
		}
		if err := lRows.Err(); err != nil {
			return nil, err
		}
	}

	const ccnlHoursPerMonth = 156.0

	return &MonthlyTimecard{
		UserID:        userID,
		FirstName:     firstName,
		LastName:      lastName,
		UserName:      fmt.Sprintf("%s %s", lastName, firstName),
		Role:          role,
		BadgeCode:     badgeCode,
		Month:         month,
		ContractHours: ccnlHoursPerMonth,
		WorkedHours:   workedHours,
		OvertimeHours: max(0, workedHours-ccnlHoursPerMonth),
		AbsenceDays:   absenceDays,
		LeaveDays:     leaveDays,
		SickDays:      sickDays,
		PermitHours:   permitHours,
		DailyEntries:  dailyEntries,
	}, nil
}

// GetAllMonthlyTimecards restituisce il riepilogo mensile per tutto il personale ATA in una singola query batch ad alte prestazioni
func (r *PostgresRepository) GetAllMonthlyTimecards(ctx context.Context, schoolID, month string) ([]MonthlyTimecard, error) {
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	t, err := time.Parse("2006-01", month)
	if err != nil {
		t = time.Now()
		month = t.Format("2006-01")
	}
	startDate := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	endDate := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")

	const ccnlHoursPerMonth = 156.0

	query := `
		WITH target_users AS (
			SELECT 
				u.id, 
				u.first_name, 
				u.last_name, 
				u.role,
				ub.badge_code
			FROM users u
			LEFT JOIN LATERAL (
				SELECT badge_code 
				FROM user_badges 
				WHERE user_id = u.id AND is_active = true 
				ORDER BY assigned_at DESC 
				LIMIT 1
			) ub ON true
			WHERE ($1 = '' OR u.school_id = NULLIF($1, '')::uuid)
			  AND u.role IN (
				  'dsga', 'assistente_amministrativo', 'collaboratore_ds', 'collaboratore_scolastico',
				  'assistente_tecnico', 'assistente_alunni', 'assistente_personale', 'assistente_contabilita',
				  'assistente_protocollo', 'assistente_sportello', 'responsabile_servizio', 'secretary'
			  )
		),
		worked_agg AS (
			SELECT 
				sa.user_id,
				COALESCE(SUM(
					EXTRACT(EPOCH FROM (COALESCE(sa.badge_exit_time, NOW()) - sa.badge_entry_time)) / 60
				), 0) AS total_worked_minutes
			FROM staff_attendance sa
			WHERE sa.user_id IN (SELECT id FROM target_users)
			  AND sa.date >= $2::date AND sa.date < $3::date
			  AND sa.badge_entry_time IS NOT NULL
			GROUP BY sa.user_id
		),
		leaves_agg AS (
			SELECT 
				slr.user_id,
				COUNT(*) FILTER (WHERE slr.type = 'ferie') AS leave_days,
				COUNT(*) FILTER (WHERE slr.type = 'malattia') AS sick_days,
				COALESCE(SUM(
					CASE 
						WHEN slr.type IN ('permesso', 'permesso_studio', 'recupero') THEN 
							CASE WHEN slr.hours > 0 THEN slr.hours ELSE (COALESCE(slr.days, 1) * 8.0) END
						ELSE 0
					END
				), 0) AS permit_hours,
				COUNT(*) FILTER (WHERE slr.type NOT IN ('ferie', 'malattia', 'permesso', 'permesso_studio', 'recupero')) AS absence_days
			FROM staff_leave_requests slr
			WHERE slr.user_id IN (SELECT id FROM target_users)
			  AND slr.start_date >= $2::date AND slr.start_date < $3::date
			  AND slr.status = 'approved'
			GROUP BY slr.user_id
		)
		SELECT 
			tu.id,
			tu.first_name,
			tu.last_name,
			tu.role,
			COALESCE(tu.badge_code, '') AS badge_code,
			COALESCE(wa.total_worked_minutes, 0) AS total_worked_minutes,
			COALESCE(la.leave_days, 0) AS leave_days,
			COALESCE(la.sick_days, 0) AS sick_days,
			COALESCE(la.permit_hours, 0) AS permit_hours,
			COALESCE(la.absence_days, 0) AS absence_days
		FROM target_users tu
		LEFT JOIN worked_agg wa ON wa.user_id = tu.id
		LEFT JOIN leaves_agg la ON la.user_id = tu.id
		ORDER BY tu.last_name, tu.first_name`

	rows, err := r.db.QueryContext(ctx, query, schoolID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	timecards := make([]MonthlyTimecard, 0)
	for rows.Next() {
		var (
			uid, fName, lName, uRole, bCode string
			totalMins, permitHrs            float64
			leaveD, sickD, absenceD         int
		)
		if err := rows.Scan(
			&uid, &fName, &lName, &uRole, &bCode,
			&totalMins, &leaveD, &sickD, &permitHrs, &absenceD,
		); err != nil {
			continue
		}

		workedHours := totalMins / 60.0
		overtimeHours := max(0, workedHours-ccnlHoursPerMonth)

		timecards = append(timecards, MonthlyTimecard{
			UserID:        uid,
			FirstName:     fName,
			LastName:      lName,
			UserName:      fmt.Sprintf("%s %s", lName, fName),
			Role:          uRole,
			BadgeCode:     bCode,
			Month:         month,
			ContractHours: ccnlHoursPerMonth,
			WorkedHours:   workedHours,
			OvertimeHours: overtimeHours,
			AbsenceDays:   absenceD,
			LeaveDays:     leaveD,
			SickDays:      sickD,
			PermitHours:   permitHrs,
			DailyEntries:  []DailyTimecardEntry{},
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return timecards, nil
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
