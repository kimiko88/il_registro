package infanzia

import (
	"database/sql"
)

type Repository interface {
	CreateSection(s *InfanziaSection) error
	ListSections(schoolID, academicYear string) ([]InfanziaSection, error)
	GetSection(id string) (*InfanziaSection, error)

	SaveDailyAttendancesBatch(sectionID, date string, attendances []InfanziaDailyAttendance) error
	GetDailyAttendances(sectionID, date string) ([]InfanziaDailyAttendance, error)

	SaveObservation(obs *InfanziaExperienceObservation) error
	ListObservationsByStudent(studentID string) ([]InfanziaExperienceObservation, error)

	CreateDiarioEntry(d *InfanziaDiarioBordo) error
	ListDiarioEntries(sectionID string, limit int) ([]InfanziaDiarioBordo, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateSection(s *InfanziaSection) error {
	query := `
		INSERT INTO infanzia_sections (school_id, section_name, color_hex, room, academic_year)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query, s.SchoolID, s.SectionName, s.ColorHex, s.Room, s.AcademicYear,
	).Scan(&s.ID, &s.CreatedAt)
}

func (r *repository) ListSections(schoolID, academicYear string) ([]InfanziaSection, error) {
	query := `
		SELECT s.id, s.school_id, s.section_name, s.color_hex, s.room, s.academic_year, s.created_at,
		       COUNT(DISTINCT ida.student_id) AS student_count
		FROM infanzia_sections s
		LEFT JOIN infanzia_daily_attendances ida ON s.id = ida.section_id
		WHERE s.school_id = $1 AND ($2 = '' OR s.academic_year = $2)
		GROUP BY s.id
		ORDER BY s.section_name ASC
	`
	rows, err := r.db.Query(query, schoolID, academicYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []InfanziaSection
	for rows.Next() {
		var s InfanziaSection
		if err := rows.Scan(
			&s.ID, &s.SchoolID, &s.SectionName, &s.ColorHex, &s.Room, &s.AcademicYear,
			&s.CreatedAt, &s.StudentCount,
		); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) GetSection(id string) (*InfanziaSection, error) {
	query := `
		SELECT s.id, s.school_id, s.section_name, s.color_hex, s.room, s.academic_year, s.created_at,
		       COUNT(DISTINCT ida.student_id) AS student_count
		FROM infanzia_sections s
		LEFT JOIN infanzia_daily_attendances ida ON s.id = ida.section_id
		WHERE s.id = $1
		GROUP BY s.id
	`
	var s InfanziaSection
	err := r.db.QueryRow(query, id).Scan(
		&s.ID, &s.SchoolID, &s.SectionName, &s.ColorHex, &s.Room, &s.AcademicYear,
		&s.CreatedAt, &s.StudentCount,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) SaveDailyAttendancesBatch(sectionID, date string, attendances []InfanziaDailyAttendance) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO infanzia_daily_attendances (
			section_id, student_id, attendance_date, is_present, meal_booked, special_diet_note, entry_time, exit_time
		) VALUES ($1, $2, $3, $4, $5, $6, $7::time, $8::time)
		ON CONFLICT (section_id, student_id, attendance_date)
		DO UPDATE SET
			is_present = EXCLUDED.is_present,
			meal_booked = EXCLUDED.meal_booked,
			special_diet_note = EXCLUDED.special_diet_note,
			entry_time = EXCLUDED.entry_time,
			exit_time = EXCLUDED.exit_time
	`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, a := range attendances {
		if _, err := stmt.Exec(
			sectionID, a.StudentID, date, a.IsPresent, a.MealBooked, a.SpecialDietNote, a.EntryTime, a.ExitTime,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *repository) GetDailyAttendances(sectionID, date string) ([]InfanziaDailyAttendance, error) {
	query := `
		SELECT a.id, a.section_id, a.student_id, TO_CHAR(a.attendance_date, 'YYYY-MM-DD'),
		       a.is_present, a.meal_booked, COALESCE(a.special_diet_note, ''),
		       CASE WHEN a.entry_time IS NOT NULL THEN TO_CHAR(a.entry_time, 'HH24:MI') ELSE NULL END,
		       CASE WHEN a.exit_time IS NOT NULL THEN TO_CHAR(a.exit_time, 'HH24:MI') ELSE NULL END,
		       a.created_at,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS student_name
		FROM infanzia_daily_attendances a
		LEFT JOIN users u ON a.student_id = u.id
		WHERE a.section_id = $1 AND a.attendance_date = $2::date
		ORDER BY student_name ASC
	`
	rows, err := r.db.Query(query, sectionID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []InfanziaDailyAttendance
	for rows.Next() {
		var a InfanziaDailyAttendance
		if err := rows.Scan(
			&a.ID, &a.SectionID, &a.StudentID, &a.AttendanceDate, &a.IsPresent, &a.MealBooked,
			&a.SpecialDietNote, &a.EntryTime, &a.ExitTime, &a.CreatedAt, &a.StudentFullName,
		); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) SaveObservation(obs *InfanziaExperienceObservation) error {
	query := `
		INSERT INTO infanzia_experience_observations (
			student_id, campo_esperienza, observation_date, observation_notes, level_descriptor, recorded_by
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query, obs.StudentID, obs.CampoEsperienza, obs.ObservationDate,
		obs.ObservationNotes, obs.LevelDescriptor, obs.RecordedBy,
	).Scan(&obs.ID, &obs.CreatedAt)
}

func (r *repository) ListObservationsByStudent(studentID string) ([]InfanziaExperienceObservation, error) {
	query := `
		SELECT o.id, o.student_id, o.campo_esperienza, TO_CHAR(o.observation_date, 'YYYY-MM-DD'),
		       o.observation_notes, o.level_descriptor, o.recorded_by, o.created_at,
		       COALESCE(us.first_name || ' ' || us.last_name, '') AS student_name,
		       COALESCE(ut.first_name || ' ' || ut.last_name, '') AS teacher_name
		FROM infanzia_experience_observations o
		LEFT JOIN users us ON o.student_id = us.id
		LEFT JOIN users ut ON o.recorded_by = ut.id
		WHERE o.student_id = $1
		ORDER BY o.observation_date DESC, o.created_at DESC
	`
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []InfanziaExperienceObservation
	for rows.Next() {
		var obs InfanziaExperienceObservation
		if err := rows.Scan(
			&obs.ID, &obs.StudentID, &obs.CampoEsperienza, &obs.ObservationDate,
			&obs.ObservationNotes, &obs.LevelDescriptor, &obs.RecordedBy, &obs.CreatedAt,
			&obs.StudentFullName, &obs.TeacherFullName,
		); err != nil {
			return nil, err
		}
		list = append(list, obs)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) CreateDiarioEntry(d *InfanziaDiarioBordo) error {
	if d.MediaURLsJSON == "" {
		d.MediaURLsJSON = "[]"
	}
	query := `
		INSERT INTO infanzia_diario_bordo (
			section_id, story_date, title, narrative_story, media_urls, created_by
		) VALUES ($1, $2, $3, $4, $5::jsonb, $6)
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query, d.SectionID, d.StoryDate, d.Title, d.NarrativeStory, d.MediaURLsJSON, d.CreatedBy,
	).Scan(&d.ID, &d.CreatedAt)
}

func (r *repository) ListDiarioEntries(sectionID string, limit int) ([]InfanziaDiarioBordo, error) {
	if limit <= 0 {
		limit = 30
	}
	query := `
		SELECT d.id, d.section_id, TO_CHAR(d.story_date, 'YYYY-MM-DD'), d.title, d.narrative_story,
		       d.media_urls::text, d.created_by, d.created_at,
		       COALESCE(s.section_name, '') AS section_name,
		       COALESCE(u.first_name || ' ' || u.last_name, '') AS teacher_name
		FROM infanzia_diario_bordo d
		LEFT JOIN infanzia_sections s ON d.section_id = s.id
		LEFT JOIN users u ON d.created_by = u.id
		WHERE d.section_id = $1
		ORDER BY d.story_date DESC, d.created_at DESC
		LIMIT $2
	`
	rows, err := r.db.Query(query, sectionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []InfanziaDiarioBordo
	for rows.Next() {
		var d InfanziaDiarioBordo
		if err := rows.Scan(
			&d.ID, &d.SectionID, &d.StoryDate, &d.Title, &d.NarrativeStory,
			&d.MediaURLsJSON, &d.CreatedBy, &d.CreatedAt, &d.SectionName, &d.TeacherFullName,
		); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}
