package timetablegen

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Repository interface {
	// Preferences
	SavePreferencesBatch(ctx context.Context, schoolID, teacherID string, academicYearID *string, prefs []PreferenceEntry) error
	GetTeacherPreferences(ctx context.Context, schoolID, teacherID string, academicYearID *string) ([]TeacherPreference, error)
	LoadAllPreferences(ctx context.Context, schoolID string, academicYearID *string) ([]TeacherPreference, error)

	// Room Requirements
	ListRoomRequirements(ctx context.Context, schoolID string) ([]SubjectRoomRequirement, error)
	SaveRoomRequirement(ctx context.Context, schoolID string, req SaveRoomRequirementRequest) (*SubjectRoomRequirement, error)
	DeleteRoomRequirement(ctx context.Context, id string) error

	// Constraints
	ListConstraints(ctx context.Context, schoolID string) ([]TimetableConstraint, error)
	SaveConstraint(ctx context.Context, schoolID string, req SaveConstraintRequest) (*TimetableConstraint, error)
	DeleteConstraint(ctx context.Context, id string) error
	GetDesiderataWindow(ctx context.Context, schoolID string) (bool, error)
	SetDesiderataWindow(ctx context.Context, schoolID string, isOpen bool) error

	// Jobs
	CreateJob(ctx context.Context, job *TimetableJob) (*TimetableJob, error)
	UpdateJob(ctx context.Context, job *TimetableJob) error
	GetJob(ctx context.Context, id string) (*TimetableJob, error)

	// Data Loading for Generator
	LoadAssignments(ctx context.Context, schoolID string, academicYearID *string) ([]AssignmentData, error)
	LoadRooms(ctx context.Context, schoolID string) ([]RoomData, error)
	LoadRoomRequirements(ctx context.Context, schoolID string) (map[string]SubjectRoomRequirement, error)
	LoadAssociatedGroups(ctx context.Context, schoolID string) ([]AssociatedGroup, error)

	// Publish
	PublishGeneratedSchedule(ctx context.Context, schoolID string, slots []GeneratedSlot) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

// ----------------- Preferences -----------------

func (r *PostgresRepository) SavePreferencesBatch(ctx context.Context, schoolID, teacherID string, academicYearID *string, prefs []PreferenceEntry) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear existing preferences for this teacher & academic year
	deleteQuery := `
		DELETE FROM teacher_schedule_preferences
		WHERE school_id = $1 AND teacher_id = $2
		  AND ($3::uuid IS NULL OR academic_year_id = $3::uuid)
	`
	if _, err := tx.ExecContext(ctx, deleteQuery, schoolID, teacherID, academicYearID); err != nil {
		return fmt.Errorf("failed to clear existing preferences: %w", err)
	}

	insertQuery := `
		INSERT INTO teacher_schedule_preferences (
			id, school_id, teacher_id, academic_year_id, day_of_week, hour_index, preference_type, reason, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
		)
	`
	stmt, err := tx.PrepareContext(ctx, insertQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, p := range prefs {
		id := uuid.New().String()
		if _, err := stmt.ExecContext(ctx, id, schoolID, teacherID, academicYearID, p.DayOfWeek, p.HourIndex, p.PreferenceType, p.Reason); err != nil {
			return fmt.Errorf("failed to insert preference (day %d hour %d): %w", p.DayOfWeek, p.HourIndex, err)
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) GetTeacherPreferences(ctx context.Context, schoolID, teacherID string, academicYearID *string) ([]TeacherPreference, error) {
	query := `
		SELECT id, school_id, teacher_id, academic_year_id, day_of_week, hour_index, preference_type, reason, created_at, updated_at
		FROM teacher_schedule_preferences
		WHERE school_id = $1 AND (teacher_id = $2 OR teacher_id IN (SELECT user_id FROM teachers WHERE id = $2::uuid))
		  AND ($3::uuid IS NULL OR academic_year_id = $3::uuid)
		ORDER BY day_of_week ASC, hour_index ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, teacherID, academicYearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []TeacherPreference
	for rows.Next() {
		var p TeacherPreference
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.TeacherID, &p.AcademicYearID, &p.DayOfWeek, &p.HourIndex,
			&p.PreferenceType, &p.Reason, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *PostgresRepository) LoadAllPreferences(ctx context.Context, schoolID string, academicYearID *string) ([]TeacherPreference, error) {
	query := `
		SELECT id, school_id, teacher_id, academic_year_id, day_of_week, hour_index, preference_type, reason, created_at, updated_at
		FROM teacher_schedule_preferences
		WHERE school_id = $1
		  AND ($2::uuid IS NULL OR academic_year_id = $2::uuid)
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, academicYearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []TeacherPreference
	for rows.Next() {
		var p TeacherPreference
		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.TeacherID, &p.AcademicYearID, &p.DayOfWeek, &p.HourIndex,
			&p.PreferenceType, &p.Reason, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

// ----------------- Room Requirements -----------------

func (r *PostgresRepository) ListRoomRequirements(ctx context.Context, schoolID string) ([]SubjectRoomRequirement, error) {
	query := `
		SELECT srr.id, srr.school_id, srr.subject_id, COALESCE(s.name, ''), srr.required_room_type,
		       COALESCE(srr.lab_hours, 1), srr.is_mandatory, srr.created_at
		FROM subject_room_requirements srr
		JOIN subjects s ON srr.subject_id = s.id
		WHERE srr.school_id = $1
		ORDER BY s.name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []SubjectRoomRequirement
	for rows.Next() {
		var req SubjectRoomRequirement
		if err := rows.Scan(
			&req.ID, &req.SchoolID, &req.SubjectID, &req.SubjectName,
			&req.RequiredRoomType, &req.LabHours, &req.IsMandatory, &req.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, req)
	}
	return list, rows.Err()
}

func (r *PostgresRepository) SaveRoomRequirement(ctx context.Context, schoolID string, req SaveRoomRequirementRequest) (*SubjectRoomRequirement, error) {
	id := uuid.New().String()
	labHours := req.LabHours
	if labHours <= 0 {
		labHours = 1
	}
	query := `
		INSERT INTO subject_room_requirements (id, school_id, subject_id, required_room_type, lab_hours, is_mandatory, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (school_id, subject_id)
		DO UPDATE SET required_room_type = EXCLUDED.required_room_type, lab_hours = EXCLUDED.lab_hours, is_mandatory = EXCLUDED.is_mandatory
		RETURNING id, created_at
	`
	res := &SubjectRoomRequirement{
		ID:               id,
		SchoolID:         schoolID,
		SubjectID:        req.SubjectID,
		RequiredRoomType: req.RequiredRoomType,
		LabHours:         labHours,
		IsMandatory:      req.IsMandatory,
	}
	err := r.db.QueryRowContext(ctx, query, id, schoolID, req.SubjectID, req.RequiredRoomType, labHours, req.IsMandatory).
		Scan(&res.ID, &res.CreatedAt)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *PostgresRepository) DeleteRoomRequirement(ctx context.Context, id string) error {
	query := `DELETE FROM subject_room_requirements WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ----------------- Constraints -----------------

func (r *PostgresRepository) ListConstraints(ctx context.Context, schoolID string) ([]TimetableConstraint, error) {
	query := `
		SELECT id, school_id, constraint_type, target_type, target_id, parameters, is_hard, priority, is_active, created_at
		FROM timetable_constraints
		WHERE school_id = $1
		ORDER BY priority DESC, created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []TimetableConstraint
	for rows.Next() {
		var c TimetableConstraint
		var paramStr string
		if err := rows.Scan(
			&c.ID, &c.SchoolID, &c.ConstraintType, &c.TargetType, &c.TargetID,
			&paramStr, &c.IsHard, &c.Priority, &c.IsActive, &c.CreatedAt,
		); err != nil {
			return nil, err
		}
		c.Parameters = json.RawMessage(paramStr)
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r *PostgresRepository) SaveConstraint(ctx context.Context, schoolID string, req SaveConstraintRequest) (*TimetableConstraint, error) {
	id := uuid.New().String()
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}
	params := string(req.Parameters)
	if params == "" {
		params = "{}"
	}
	priority := req.Priority
	if priority <= 0 {
		priority = 5
	}

	query := `
		INSERT INTO timetable_constraints (
			id, school_id, constraint_type, target_type, target_id, parameters, is_hard, priority, is_active, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, NOW()
		)
		RETURNING created_at
	`
	res := &TimetableConstraint{
		ID:             id,
		SchoolID:       schoolID,
		ConstraintType: req.ConstraintType,
		TargetType:     req.TargetType,
		TargetID:       req.TargetID,
		Parameters:     json.RawMessage(params),
		IsHard:         req.IsHard,
		Priority:       priority,
		IsActive:       active,
	}
	err := r.db.QueryRowContext(ctx, query, id, schoolID, req.ConstraintType, req.TargetType, req.TargetID, params, req.IsHard, priority, active).
		Scan(&res.CreatedAt)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *PostgresRepository) DeleteConstraint(ctx context.Context, id string) error {
	query := `DELETE FROM timetable_constraints WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresRepository) GetDesiderataWindow(ctx context.Context, schoolID string) (bool, error) {
	query := `
		SELECT is_active
		FROM timetable_constraints
		WHERE school_id = $1 AND constraint_type = 'teacher_desiderata_window'
		LIMIT 1
	`
	var active bool
	err := r.db.QueryRowContext(ctx, query, schoolID).Scan(&active)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return active, nil
}

func (r *PostgresRepository) SetDesiderataWindow(ctx context.Context, schoolID string, isOpen bool) error {
	updateQuery := `
		UPDATE timetable_constraints
		SET is_active = $1
		WHERE school_id = $2 AND constraint_type = 'teacher_desiderata_window'
	`
	res, err := r.db.ExecContext(ctx, updateQuery, isOpen, schoolID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		id := uuid.New().String()
		insertQuery := `
			INSERT INTO timetable_constraints (
				id, school_id, constraint_type, parameters, is_hard, priority, is_active, created_at
			) VALUES (
				$1, $2, 'teacher_desiderata_window', '{}', false, 1, $3, NOW()
			)
		`
		_, err = r.db.ExecContext(ctx, insertQuery, id, schoolID, isOpen)
		if err != nil {
			return err
		}
	}
	return nil
}

// ----------------- Jobs -----------------

func (r *PostgresRepository) CreateJob(ctx context.Context, job *TimetableJob) (*TimetableJob, error) {
	if job.ID == "" {
		job.ID = uuid.New().String()
	}
	if job.Status == "" {
		job.Status = JobStatusPending
	}
	if job.Algorithm == "" {
		job.Algorithm = "greedy_local_search"
	}
	params := string(job.Parameters)
	if params == "" {
		params = "{}"
	}
	summary := string(job.ResultSummary)
	if summary == "" {
		summary = "{}"
	}

	query := `
		INSERT INTO timetable_generation_jobs (
			id, school_id, academic_year_id, triggered_by, status, algorithm, parameters, result_summary, started_at, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
		)
		RETURNING created_at
	`
	err := r.db.QueryRowContext(ctx, query,
		job.ID, job.SchoolID, job.AcademicYearID, job.TriggeredBy, job.Status, job.Algorithm, params, summary,
	).Scan(&job.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create timetable job: %w", err)
	}
	return job, nil
}

func (r *PostgresRepository) UpdateJob(ctx context.Context, job *TimetableJob) error {
	summary := string(job.ResultSummary)
	if summary == "" {
		summary = "{}"
	}

	query := `
		UPDATE timetable_generation_jobs
		SET status = $1, result_summary = $2, error_message = $3, completed_at = $4
		WHERE id = $5
	`
	_, err := r.db.ExecContext(ctx, query, job.Status, summary, job.ErrorMessage, job.CompletedAt, job.ID)
	return err
}

func (r *PostgresRepository) GetJob(ctx context.Context, id string) (*TimetableJob, error) {
	query := `
		SELECT id, school_id, academic_year_id, triggered_by, status, algorithm, parameters, result_summary, started_at, completed_at, error_message, created_at
		FROM timetable_generation_jobs
		WHERE id = $1
	`
	var j TimetableJob
	var paramStr, summaryStr string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&j.ID, &j.SchoolID, &j.AcademicYearID, &j.TriggeredBy, &j.Status, &j.Algorithm,
		&paramStr, &summaryStr, &j.StartedAt, &j.CompletedAt, &j.ErrorMessage, &j.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	j.Parameters = json.RawMessage(paramStr)
	j.ResultSummary = json.RawMessage(summaryStr)
	return &j, nil
}

// ----------------- Data Loading for Generator -----------------

func (r *PostgresRepository) LoadAssignments(ctx context.Context, schoolID string, academicYearID *string) ([]AssignmentData, error) {
	query := `
		SELECT cs.class_id, COALESCE(c.name || ' ' || COALESCE(c.section, ''), c.name, ''), c.building_id,
		       cs.subject_id, s.name,
		       COALESCE(t.id, 'unassigned-' || cs.id), COALESCE(t.user_id, 'unassigned-' || cs.id),
		       COALESCE(NULLIF(TRIM(COALESCE(u.last_name, '') || ' ' || COALESCE(u.first_name, '')), ''), 'Docente da Nominare (Cattedra non assegnata)'),
		       t.hiring_date,
		       CEIL(cs.hours_per_week)::int
		FROM class_subjects cs
		JOIN classes c ON cs.class_id = c.id
		JOIN subjects s ON cs.subject_id = s.id
		LEFT JOIN teachers t ON cs.teacher_id = t.id OR cs.teacher_id = t.user_id
		LEFT JOIN users u ON t.user_id = u.id
		WHERE c.school_id = $1
		  AND ($2::uuid IS NULL OR c.academic_year_id = $2::uuid)
		  AND cs.hours_per_week > 0
		ORDER BY c.name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID, academicYearID)
	if err != nil {
		return nil, fmt.Errorf("failed to load assignments: %w", err)
	}
	defer rows.Close()

	var list []AssignmentData
	for rows.Next() {
		var a AssignmentData
		var hiring sql.NullTime
		if err := rows.Scan(
			&a.ClassID, &a.ClassName, &a.BuildingID,
			&a.SubjectID, &a.SubjectName,
			&a.TeacherID, &a.TeacherUserID, &a.TeacherName,
			&hiring,
			&a.HoursPerWeek,
		); err != nil {
			return nil, err
		}
		if hiring.Valid {
			a.HiringDate = &hiring.Time
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate assignments: %w", err)
	}

	// 2. Load explicit unassigned hours ("spezzoni orari") from timetable_constraints
	cQuery := `
		SELECT tc.id, tc.target_id, COALESCE(c.name || ' ' || COALESCE(c.section, ''), c.name, ''), c.building_id, tc.parameters
		FROM timetable_constraints tc
		JOIN classes c ON tc.target_id = c.id
		WHERE tc.school_id = $1
		  AND tc.constraint_type = 'unassigned_hours'
		  AND tc.is_active = true
	`
	cRows, cErr := r.db.QueryContext(ctx, cQuery, schoolID)
	if cErr == nil {
		defer cRows.Close()
		for cRows.Next() {
			var cID, targetID, className string
			var bldID *string
			var rawParams []byte
			if scanErr := cRows.Scan(&cID, &targetID, &className, &bldID, &rawParams); scanErr == nil {
				var p struct {
					SubjectID       string `json:"subject_id"`
					SubjectName     string `json:"subject_name"`
					HoursPerWeek    int    `json:"hours_per_week"`
					PlaceholderName string `json:"placeholder_name"`
				}
				if json.Unmarshal(rawParams, &p) == nil && p.HoursPerWeek > 0 {
					placeholder := p.PlaceholderName
					if placeholder == "" {
						placeholder = fmt.Sprintf("Docente da Nominare (%s)", p.SubjectName)
					}
					list = append(list, AssignmentData{
						ClassID:       targetID,
						ClassName:     className,
						BuildingID:    bldID,
						SubjectID:     p.SubjectID,
						SubjectName:   p.SubjectName,
						TeacherID:     "spezzone-" + cID,
						TeacherUserID: "spezzone-" + cID,
						TeacherName:   placeholder,
						HoursPerWeek:  p.HoursPerWeek,
					})
				}
			}
		}
		if err := cRows.Err(); err != nil {
			return nil, fmt.Errorf("failed to iterate unassigned hours: %w", err)
		}
	}

	// 3. Link associated groups to assignments
	groups, _ := r.LoadAssociatedGroups(ctx, schoolID)
	for _, g := range groups {
		classSet := make(map[string]bool)
		for _, cid := range g.ClassIDs {
			classSet[cid] = true
		}
		for i := range list {
			if classSet[list[i].ClassID] && list[i].SubjectID == g.SubjectID && (list[i].TeacherID == g.TeacherID || list[i].TeacherUserID == g.TeacherID) {
				groupID := g.ID
				list[i].AssociatedGroupID = &groupID
				list[i].IsAssociatedGroup = true
			}
		}
	}

	return list, nil
}

func (r *PostgresRepository) LoadRooms(ctx context.Context, schoolID string) ([]RoomData, error) {
	query := `
		SELECT id, building_id, name, room_type, capacity, is_active
		FROM bookable_rooms
		WHERE school_id = $1 AND is_active = true
	`
	rows, err := r.db.QueryContext(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []RoomData
	for rows.Next() {
		var rd RoomData
		if err := rows.Scan(&rd.ID, &rd.BuildingID, &rd.Name, &rd.RoomType, &rd.Capacity, &rd.IsActive); err != nil {
			return nil, err
		}
		list = append(list, rd)
	}
	return list, rows.Err()
}

func (r *PostgresRepository) LoadRoomRequirements(ctx context.Context, schoolID string) (map[string]SubjectRoomRequirement, error) {
	reqs, err := r.ListRoomRequirements(ctx, schoolID)
	if err != nil {
		return nil, err
	}
	res := make(map[string]SubjectRoomRequirement)
	for _, req := range reqs {
		res[req.SubjectID] = req
	}
	return res, nil
}

func (r *PostgresRepository) LoadAssociatedGroups(ctx context.Context, schoolID string) ([]AssociatedGroup, error) {
	var list []AssociatedGroup

	// 1. Load from timetable_constraints with constraint_type = 'associated_group'
	queryConstraints := `
		SELECT tc.id, tc.parameters
		FROM timetable_constraints tc
		WHERE tc.school_id = $1
		  AND tc.constraint_type = 'associated_group'
		  AND tc.is_active = true
	`
	rows, err := r.db.QueryContext(ctx, queryConstraints, schoolID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id string
			var raw []byte
			if err := rows.Scan(&id, &raw); err == nil {
				var ag AssociatedGroup
				if err := json.Unmarshal(raw, &ag); err == nil {
					ag.ID = id
					ag.SchoolID = schoolID
					if ag.HoursPerWeek <= 0 {
						ag.HoursPerWeek = 2
					}
					list = append(list, ag)
				}
			}
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	// 2. Load from groups table where students belong to multiple classes
	queryGroups := `
		SELECT g.id, g.name, g.subject_id, COALESCE(s.name, ''),
		       COALESCE(g.teacher_id::text, ''),
		       COALESCE(u.last_name || ' ' || u.first_name, ''),
		       COALESCE(STRING_AGG(DISTINCT cs.class_id::text, ','), '') AS class_ids_str
		FROM groups g
		LEFT JOIN subjects s ON g.subject_id = s.id
		LEFT JOIN users u ON g.teacher_id = u.id
		JOIN group_students gs ON gs.group_id = g.id
		JOIN class_students cs ON (cs.student_id = gs.student_id OR cs.student_id IN (SELECT id FROM students WHERE user_id = gs.student_id))
		WHERE g.school_id = $1 AND g.teacher_id IS NOT NULL AND g.subject_id IS NOT NULL
		GROUP BY g.id, g.name, g.subject_id, s.name, g.teacher_id, u.last_name, u.first_name
		HAVING COUNT(DISTINCT cs.class_id) > 1
	`
	gRows, gErr := r.db.QueryContext(ctx, queryGroups, schoolID)
	if gErr == nil {
		defer gRows.Close()
		for gRows.Next() {
			var gID, gName, gSubID, gSubName, gTeacherID, gTeacherName, classIDsStr string
			if err := gRows.Scan(&gID, &gName, &gSubID, &gSubName, &gTeacherID, &gTeacherName, &classIDsStr); err == nil {
				var classIDs []string
				for _, cid := range strings.Split(classIDsStr, ",") {
					cid = strings.TrimSpace(cid)
					if cid != "" {
						classIDs = append(classIDs, cid)
					}
				}
				if len(classIDs) > 1 {
					exists := false
					for _, item := range list {
						if item.SubjectID == gSubID && item.TeacherID == gTeacherID {
							exists = true
							break
						}
					}
					if !exists {
						list = append(list, AssociatedGroup{
							ID:           gID,
							SchoolID:     schoolID,
							Name:         gName,
							SubjectID:    gSubID,
							SubjectName:  gSubName,
							TeacherID:    gTeacherID,
							TeacherName:  gTeacherName,
							ClassIDs:     classIDs,
							HoursPerWeek: 3,
						})
					}
				}
			}
		}
		if err := gRows.Err(); err != nil {
			return nil, err
		}
	}

	return list, nil
}

// ----------------- Publish Schedule -----------------

func (r *PostgresRepository) PublishGeneratedSchedule(ctx context.Context, schoolID string, slots []GeneratedSlot) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Collect unique classIDs to clear existing schedules
	classMap := make(map[string]bool)
	for _, s := range slots {
		classMap[s.ClassID] = true
	}

	for classID := range classMap {
		_, err := tx.ExecContext(ctx, "DELETE FROM class_schedules WHERE class_id = $1", classID)
		if err != nil {
			return fmt.Errorf("failed to clear existing schedule for class %s: %w", classID, err)
		}
	}

	// 2. Insert new slots into class_schedules
	insertStmt, err := tx.PrepareContext(ctx, `
		INSERT INTO class_schedules (
			id, class_id, day_of_week, hour_index, subject_id, teacher_id, room, room_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
		)
	`)
	if err != nil {
		return err
	}
	defer insertStmt.Close()

	for _, s := range slots {
		id := uuid.New().String()
		teacherID := s.TeacherUserID
		if teacherID == nil {
			teacherID = s.TeacherID
		}
		var insertTeacherID *string
		if teacherID != nil && !strings.HasPrefix(*teacherID, "unassigned-") && !strings.HasPrefix(*teacherID, "spezzone-") {
			insertTeacherID = teacherID
		}
		_, err := insertStmt.ExecContext(ctx,
			id, s.ClassID, s.DayOfWeek, s.HourIndex, s.SubjectID, insertTeacherID, s.RoomName, s.RoomID,
		)
		if err != nil {
			return fmt.Errorf("failed to insert schedule slot (class %s day %d hour %d): %w", s.ClassID, s.DayOfWeek, s.HourIndex, err)
		}
	}

	return tx.Commit()
}
