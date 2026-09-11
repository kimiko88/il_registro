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
		SELECT lr.*,
		       u.first_name AS user_first_name,
		       u.last_name  AS user_last_name,
		       u.role       AS user_role
		FROM staff_leave_requests lr
		JOIN users u ON u.id = lr.user_id
		WHERE lr.school_id = $1`
	args := []interface{}{schoolID}
	idx := 2

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

	// Recupera dati utente
	var firstName, lastName, role string
	err := r.db.QueryRowContext(ctx,
		`SELECT first_name, last_name, role FROM users WHERE id=$1 AND school_id=$2`,
		userID, schoolID).Scan(&firstName, &lastName, &role)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("utente non trovato")
	}
	if err != nil {
		return nil, err
	}

	// Calcolo ore da timbrature badge
	var workedMins sql.NullFloat64
	err = r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(
		    EXTRACT(EPOCH FROM (COALESCE(badge_exit_time, NOW()) - badge_entry_time)) / 60
		), 0) AS worked_minutes
		FROM staff_attendances
		WHERE school_id=$1 AND user_id=$2
		  AND TO_CHAR(date::date, 'YYYY-MM') = $3
		  AND badge_entry_time IS NOT NULL`,
		schoolID, userID, month).Scan(&workedMins)
	if err != nil {
		return nil, err
	}
	workedHours := workedMins.Float64 / 60.0

	// Conteggi assenze per tipo
	var absenceDays, leaveDays, sickDays int
	var permitHours float64

	rows, err := r.db.QueryContext(ctx, `
		SELECT type, COUNT(*) as cnt, COALESCE(SUM(hours),0) as hrs
		FROM staff_leave_requests
		WHERE school_id=$1 AND user_id=$2
		  AND TO_CHAR(start_date::date, 'YYYY-MM') = $3
		  AND status='approved'
		GROUP BY type`, schoolID, userID, month)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ltype string
			var cnt int
			var hrs float64
			if err := rows.Scan(&ltype, &cnt, &hrs); err == nil {
				switch LeaveType(ltype) {
				case LeaveTypeFerie:
					leaveDays += cnt
				case LeaveTypeMalattia:
					sickDays += cnt
				case LeaveTypePermesso, LeaveTypePermessoStudio, LeaveTypeRecupero:
					permitHours += hrs
					if hrs == 0 {
						permitHours += float64(cnt) * 8 // default 8h per giorno
					}
				default:
					absenceDays += cnt
				}
			}
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	// Ore contrattuali (CCNL Scuola: 36h/settimana = ~156h/mese)
	const ccnlHoursPerMonth = 156.0

	return &MonthlyTimecard{
		UserID:        userID,
		FirstName:     firstName,
		LastName:      lastName,
		Role:          role,
		Month:         month,
		ContractHours: ccnlHoursPerMonth,
		WorkedHours:   workedHours,
		OvertimeHours: max(0, workedHours-ccnlHoursPerMonth),
		AbsenceDays:   absenceDays,
		LeaveDays:     leaveDays,
		SickDays:      sickDays,
		PermitHours:   permitHours,
	}, nil
}

// GetAllMonthlyTimecards restituisce il riepilogo mensile per tutto il personale ATA
func (r *PostgresRepository) GetAllMonthlyTimecards(ctx context.Context, schoolID, month string) ([]MonthlyTimecard, error) {
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id FROM users
		WHERE school_id=$1 AND role IN ('dsga','assistente_amministrativo','collaboratore_ds','collaboratore_scolastico')
		ORDER BY last_name, first_name`, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timecards []MonthlyTimecard
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			continue
		}
		tc, err := r.GetMonthlyTimecard(ctx, schoolID, uid, month)
		if err == nil {
			timecards = append(timecards, *tc)
		}
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
