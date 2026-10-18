package support

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Repository interface {
	CreateDiaryEntry(ctx context.Context, entry *SupportDiaryEntry) error
	ListDiaryEntries(ctx context.Context, schoolID, studentID, teacherID, classID string, isFamily bool) ([]SupportDiaryEntry, error)
	DeleteDiaryEntry(ctx context.Context, id, teacherID string) error

	CreatePeiGoal(ctx context.Context, goal *SupportPeiGoal) error
	UpdatePeiGoalProgress(ctx context.Context, id, progressStatus string) error
	ListPeiGoals(ctx context.Context, schoolID, studentID string) ([]SupportPeiGoal, error)
	DeletePeiGoal(ctx context.Context, id string) error

	CreateGloMeeting(ctx context.Context, meeting *SupportGloMeeting) error
	GetGloMeeting(ctx context.Context, id string) (*SupportGloMeeting, error)
	ListGloMeetings(ctx context.Context, schoolID, studentID string) ([]SupportGloMeeting, error)
	SignGloMeeting(ctx context.Context, meetingID string, sig *GloSignature) error
	SaveSupportHours(ctx context.Context, h *SupportHoursQuantification) error
	GetSupportHours(ctx context.Context, schoolID, studentID, schoolYear, peiStage string) (*SupportHoursQuantification, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateDiaryEntry(ctx context.Context, e *SupportDiaryEntry) error {
	query := `
		INSERT INTO support_diaries (
			id, school_id, teacher_id, student_id, class_id, entry_date,
			time_slot, co_teacher_id, activity_type, topic_and_activities,
			student_responses, educator_notes, is_shared_with_family,
			created_at, updated_at
		) VALUES (
			COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()), $2, $3, $4, $5, $6::date,
			$7, NULLIF($8, '')::uuid, $9, $10,
			$11, $12, $13,
			NOW(), NOW()
		) RETURNING id
	`
	coTeacher := ""
	if e.CoTeacherID != nil {
		coTeacher = *e.CoTeacherID
	}
	return r.db.QueryRowContext(ctx, query,
		e.ID, e.SchoolID, e.TeacherID, e.StudentID, e.ClassID, e.EntryDate,
		e.TimeSlot, coTeacher, e.ActivityType, e.TopicAndActivities,
		e.StudentResponses, e.EducatorNotes, e.IsSharedWithFamily,
	).Scan(&e.ID)
}

func (r *repository) ListDiaryEntries(ctx context.Context, schoolID, studentID, teacherID, classID string, isFamily bool) ([]SupportDiaryEntry, error) {
	familyFilter := ""
	if isFamily {
		familyFilter = " AND sd.is_shared_with_family = TRUE "
	}

	query := fmt.Sprintf(`
		SELECT sd.id, sd.school_id, sd.teacher_id,
		       TRIM(COALESCE(tu.first_name, '') || ' ' || COALESCE(tu.last_name, '')) as teacher_name,
		       sd.student_id,
		       TRIM(COALESCE(su.first_name, '') || ' ' || COALESCE(su.last_name, '')) as student_name,
		       sd.class_id, cls.name as class_name,
		       sd.entry_date::text, sd.time_slot,
		       sd.co_teacher_id::text,
		       TRIM(COALESCE(cu.first_name, '') || ' ' || COALESCE(cu.last_name, '')) as co_teacher_name,
		       sd.activity_type, sd.topic_and_activities, sd.student_responses,
		       sd.educator_notes, sd.is_shared_with_family, sd.created_at, sd.updated_at
		FROM support_diaries sd
		JOIN teachers t ON sd.teacher_id = t.id
		JOIN users tu ON t.user_id = tu.id
		JOIN students s ON sd.student_id = s.id
		JOIN users su ON s.user_id = su.id
		JOIN classes cls ON sd.class_id = cls.id
		LEFT JOIN teachers ct ON sd.co_teacher_id = ct.id
		LEFT JOIN users cu ON ct.user_id = cu.id
		WHERE sd.school_id = $1
		  %s
		  AND ($2 = '' OR sd.student_id = NULLIF($2, '')::uuid)
		  AND ($3 = '' OR sd.teacher_id = NULLIF($3, '')::uuid)
		  AND ($4 = '' OR sd.class_id = NULLIF($4, '')::uuid)
		ORDER BY sd.entry_date DESC, sd.created_at DESC
	`, familyFilter)

	rows, err := r.db.QueryContext(ctx, query, schoolID, studentID, teacherID, classID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var list []SupportDiaryEntry
	for rows.Next() {
		var e SupportDiaryEntry
		var coID, coName sql.NullString
		if err := rows.Scan(
			&e.ID, &e.SchoolID, &e.TeacherID, &e.TeacherName,
			&e.StudentID, &e.StudentName,
			&e.ClassID, &e.ClassName,
			&e.EntryDate, &e.TimeSlot,
			&coID, &coName,
			&e.ActivityType, &e.TopicAndActivities, &e.StudentResponses,
			&e.EducatorNotes, &e.IsSharedWithFamily, &e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if coID.Valid {
			e.CoTeacherID = &coID.String
			e.CoTeacherName = coName.String
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *repository) DeleteDiaryEntry(ctx context.Context, id, teacherID string) error {
	query := `DELETE FROM support_diaries WHERE id = $1 AND ($2 = '' OR teacher_id = NULLIF($2, '')::uuid)`
	_, err := r.db.ExecContext(ctx, query, id, teacherID)
	return err
}

func (r *repository) CreatePeiGoal(ctx context.Context, g *SupportPeiGoal) error {
	query := `
		INSERT INTO support_pei_goals (
			id, school_id, student_id, pei_type, axis, ministerial_dimension, pathway_type, title, description,
			expected_term, progress_status, glo_notes, created_at, updated_at
		) VALUES (
			COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()), $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, NOW(), NOW()
		) RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		g.ID, g.SchoolID, g.StudentID, g.PeiType, g.Axis, g.MinisterialDimension, g.PathwayType, g.Title, g.Description,
		g.ExpectedTerm, g.ProgressStatus, g.GloNotes,
	).Scan(&g.ID)
}

func (r *repository) UpdatePeiGoalProgress(ctx context.Context, id, progressStatus string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE support_pei_goals SET progress_status = $1, updated_at = NOW() WHERE id = $2`, progressStatus, id)
	return err
}

func (r *repository) ListPeiGoals(ctx context.Context, schoolID, studentID string) ([]SupportPeiGoal, error) {
	query := `
		SELECT pg.id, pg.school_id, pg.student_id,
		       TRIM(COALESCE(u.first_name, '') || ' ' || COALESCE(u.last_name, '')) as student_name,
		       pg.pei_type, pg.axis,
		       COALESCE(pg.ministerial_dimension, 'dimensione_autonomia') as ministerial_dimension,
		       COALESCE(pg.pathway_type, 'percorso_b_personalizzato') as pathway_type,
		       pg.title, pg.description, pg.expected_term, pg.progress_status,
		       COALESCE(pg.glo_notes, '') as glo_notes,
		       pg.created_at, pg.updated_at
		FROM support_pei_goals pg
		JOIN students s ON pg.student_id = s.id
		JOIN users u ON s.user_id = u.id
		WHERE pg.school_id = $1 AND ($2 = '' OR pg.student_id = NULLIF($2, '')::uuid)
		ORDER BY pg.axis ASC, pg.created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, studentID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var list []SupportPeiGoal
	for rows.Next() {
		var g SupportPeiGoal
		if err := rows.Scan(
			&g.ID, &g.SchoolID, &g.StudentID, &g.StudentName,
			&g.PeiType, &g.Axis, &g.MinisterialDimension, &g.PathwayType,
			&g.Title, &g.Description, &g.ExpectedTerm, &g.ProgressStatus,
			&g.GloNotes, &g.CreatedAt, &g.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, g)
	}
	return list, rows.Err()
}

func (r *repository) DeletePeiGoal(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM support_pei_goals WHERE id = $1`, id)
	return err
}

func (r *repository) CreateGloMeeting(ctx context.Context, m *SupportGloMeeting) error {
	if m.ID == "" {
		m.ID = fmt.Sprintf("glo-%d", time.Now().UnixNano())
	}
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	if r.db == nil {
		return nil
	}

	membersJSON, _ := json.Marshal(m.MembersPresent)
	query := `
		INSERT INTO support_glo_meetings (
			id, school_id, student_id, meeting_date, meeting_type,
			school_year, members_present, discussion_summary, decisions_taken, status,
			created_at, updated_at
		) VALUES (
			COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()), $2, $3, $4::date, $5,
			$6, $7::jsonb, $8, $9, $10,
			NOW(), NOW()
		) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		m.ID, m.SchoolID, m.StudentID, m.MeetingDate, m.MeetingType,
		m.SchoolYear, string(membersJSON), m.DiscussionSummary, m.DecisionsTaken, m.Status,
	).Scan(&m.ID)
}

func (r *repository) GetGloMeeting(ctx context.Context, id string) (*SupportGloMeeting, error) {
	if r.db == nil {
		return &SupportGloMeeting{
			ID:                id,
			MeetingType:       "INIZIALE",
			SchoolYear:        "2025/2026",
			Status:            "BOZZA",
			DiscussionSummary: "Sintesi GLO",
			CreatedAt:         time.Now(),
		}, nil
	}

	query := `
		SELECT id, school_id, student_id, meeting_date::text, meeting_type,
		       school_year, members_present, discussion_summary, decisions_taken, status,
		       created_at, updated_at
		FROM support_glo_meetings
		WHERE id = $1::uuid`

	var m SupportGloMeeting
	var membersRaw []byte
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.SchoolID, &m.StudentID, &m.MeetingDate, &m.MeetingType,
		&m.SchoolYear, &membersRaw, &m.DiscussionSummary, &m.DecisionsTaken, &m.Status,
		&m.CreatedAt, &m.UpdatedAt,
	); err != nil {
		return nil, err
	}
	_ = json.Unmarshal(membersRaw, &m.MembersPresent)

	// Fetch signatures
	sigQuery := `
		SELECT id, meeting_id, member_id, member_name, member_role, signed_at, signature_type, signature_hash, status
		FROM support_glo_signatures
		WHERE meeting_id = $1::uuid`
	rows, err := r.db.QueryContext(ctx, sigQuery, id)
	if err == nil {
		defer func() { _ = rows.Close() }()
		for rows.Next() {
			var sig GloSignature
			if err := rows.Scan(&sig.ID, &sig.MeetingID, &sig.MemberID, &sig.MemberName, &sig.MemberRole, &sig.SignedAt, &sig.SignatureType, &sig.SignatureHash, &sig.Status); err == nil {
				m.Signatures = append(m.Signatures, sig)
			}
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return &m, nil
}

func (r *repository) ListGloMeetings(ctx context.Context, schoolID, studentID string) ([]SupportGloMeeting, error) {
	if r.db == nil {
		return []SupportGloMeeting{
			{
				ID:                "glo-sample-1",
				SchoolID:          schoolID,
				StudentID:         studentID,
				MeetingType:       "INIZIALE",
				SchoolYear:        "2025/2026",
				Status:            "APPROVATO_CON_FIRME",
				DiscussionSummary: "Verbale di approvazione del Piano Educativo Individualizzato (PEI).",
				CreatedAt:         time.Now(),
			},
		}, nil
	}

	query := `
		SELECT m.id, m.school_id, m.student_id,
		       COALESCE(u.first_name || ' ' || u.last_name, '') as student_name,
		       m.meeting_date::text, m.meeting_type, m.school_year, m.members_present,
		       m.discussion_summary, m.decisions_taken, m.status, m.created_at, m.updated_at
		FROM support_glo_meetings m
		JOIN students s ON m.student_id = s.id
		JOIN users u ON s.user_id = u.id
		WHERE m.school_id = $1 AND ($2 = '' OR m.student_id = NULLIF($2, '')::uuid)
		ORDER BY m.meeting_date DESC`

	rows, err := r.db.QueryContext(ctx, query, schoolID, studentID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var list []SupportGloMeeting
	for rows.Next() {
		var m SupportGloMeeting
		var membersRaw []byte
		if err := rows.Scan(
			&m.ID, &m.SchoolID, &m.StudentID, &m.StudentName,
			&m.MeetingDate, &m.MeetingType, &m.SchoolYear, &membersRaw,
			&m.DiscussionSummary, &m.DecisionsTaken, &m.Status,
			&m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(membersRaw, &m.MembersPresent)
		list = append(list, m)
	}
	return list, rows.Err()
}

func (r *repository) SignGloMeeting(ctx context.Context, meetingID string, sig *GloSignature) error {
	if sig.ID == "" {
		sig.ID = fmt.Sprintf("sig-%d", time.Now().UnixNano())
	}
	now := time.Now()
	sig.SignedAt = &now
	sig.Status = "SIGNED"

	if r.db == nil {
		return nil
	}

	query := `
		INSERT INTO support_glo_signatures (
			id, meeting_id, member_id, member_name, member_role, signed_at, signature_type, signature_hash, status
		) VALUES (
			COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()), $2::uuid, $3, $4, $5, NOW(), $6, $7, $8
		) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		sig.ID, meetingID, sig.MemberID, sig.MemberName, sig.MemberRole,
		sig.SignatureType, sig.SignatureHash, sig.Status,
	).Scan(&sig.ID)
}

func (r *repository) SaveSupportHours(ctx context.Context, h *SupportHoursQuantification) error {
	if h.ID == "" {
		h.ID = fmt.Sprintf("sh-%d", time.Now().UnixNano())
	}
	h.CreatedAt = time.Now()
	h.UpdatedAt = time.Now()

	if r.db == nil {
		return nil
	}

	query := `
		INSERT INTO support_hours_quantifications (
			id, school_id, student_id, school_year, pei_stage,
			debt_relazione, debt_comunicazione, debt_autonomia, debt_cognitiva,
			weekly_support_hours_requested, weekly_support_hours_assigned,
			assistant_hours_requested, assistant_hours_assigned,
			transportation_needed, personal_care_needed, justification_summary,
			created_at, updated_at
		) VALUES (
			COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()), $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11,
			$12, $13,
			$14, $15, $16,
			NOW(), NOW()
		)
		ON CONFLICT (school_id, student_id, school_year, pei_stage)
		DO UPDATE SET
			debt_relazione = EXCLUDED.debt_relazione,
			debt_comunicazione = EXCLUDED.debt_comunicazione,
			debt_autonomia = EXCLUDED.debt_autonomia,
			debt_cognitiva = EXCLUDED.debt_cognitiva,
			weekly_support_hours_requested = EXCLUDED.weekly_support_hours_requested,
			weekly_support_hours_assigned = EXCLUDED.weekly_support_hours_assigned,
			assistant_hours_requested = EXCLUDED.assistant_hours_requested,
			assistant_hours_assigned = EXCLUDED.assistant_hours_assigned,
			transportation_needed = EXCLUDED.transportation_needed,
			personal_care_needed = EXCLUDED.personal_care_needed,
			justification_summary = EXCLUDED.justification_summary,
			updated_at = NOW()
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		h.ID, h.SchoolID, h.StudentID, h.SchoolYear, h.PeiStage,
		h.DebtRelazione, h.DebtComunicazione, h.DebtAutonomia, h.DebtCognitiva,
		h.WeeklySupportHoursRequested, h.WeeklySupportHoursAssigned,
		h.AssistantHoursRequested, h.AssistantHoursAssigned,
		h.TransportationNeeded, h.PersonalCareNeeded, h.JustificationSummary,
	).Scan(&h.ID)
}

func (r *repository) GetSupportHours(ctx context.Context, schoolID, studentID, schoolYear, peiStage string) (*SupportHoursQuantification, error) {
	if r.db == nil {
		return &SupportHoursQuantification{
			ID:                          "sh-sample-1",
			SchoolID:                    schoolID,
			StudentID:                   studentID,
			SchoolYear:                  schoolYear,
			PeiStage:                    peiStage,
			DebtRelazione:               3,
			DebtComunicazione:           2,
			DebtAutonomia:               3,
			DebtCognitiva:               2,
			WeeklySupportHoursRequested: 18.0,
			WeeklySupportHoursAssigned:  18.0,
			AssistantHoursRequested:     10.0,
			AssistantHoursAssigned:      10.0,
			TransportationNeeded:        false,
			PersonalCareNeeded:          true,
			JustificationSummary:        "Fabbisogno quantificato secondo Allegato C Linee Guida D.I. 182/2020.",
		}, nil
	}

	query := `
		SELECT id, school_id, student_id, school_year, pei_stage,
		       debt_relazione, debt_comunicazione, debt_autonomia, debt_cognitiva,
		       weekly_support_hours_requested, weekly_support_hours_assigned,
		       assistant_hours_requested, assistant_hours_assigned,
		       transportation_needed, personal_care_needed, justification_summary,
		       created_at, updated_at
		FROM support_hours_quantifications
		WHERE school_id = $1 AND student_id = $2::uuid AND school_year = $3 AND pei_stage = $4`

	var h SupportHoursQuantification
	if err := r.db.QueryRowContext(ctx, query, schoolID, studentID, schoolYear, peiStage).Scan(
		&h.ID, &h.SchoolID, &h.StudentID, &h.SchoolYear, &h.PeiStage,
		&h.DebtRelazione, &h.DebtComunicazione, &h.DebtAutonomia, &h.DebtCognitiva,
		&h.WeeklySupportHoursRequested, &h.WeeklySupportHoursAssigned,
		&h.AssistantHoursRequested, &h.AssistantHoursAssigned,
		&h.TransportationNeeded, &h.PersonalCareNeeded, &h.JustificationSummary,
		&h.CreatedAt, &h.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &h, nil
}
