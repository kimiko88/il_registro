package grades

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ReligionJudgments elenca i giudizi validi per la materia IRC.
// Valori ammessi dalle disposizioni ministeriali italiane:
// Non classificabile, Insufficiente, Sufficiente, Buono, Distinto, Ottimo.
var ReligionJudgments = []string{
	"Non classificabile",
	"Insufficiente",
	"Sufficiente",
	"Buono",
	"Distinto",
	"Ottimo",
}

// ReligionJudgmentValues mappa i giudizi IRC al valore numerico interno.
// Usato solo per ordinamento e statistiche aggregate; nella visualizzazione
// si mostra sempre il giudizio testuale.
var ReligionJudgmentValues = map[string]float64{
	"Non classificabile": 0,
	"Insufficiente":      4,
	"Sufficiente":        6,
	"Buono":              7,
	"Distinto":           8,
	"Ottimo":             10,
}

type Validator struct {
	db *sql.DB
}

func NewValidator(db *sql.DB) *Validator {
	return &Validator{db: db}
}

// IsAvailable restituisce true se il validatore ha una connessione DB attiva.
// Usato dal service per evitare accesso diretto al campo db interno.
func (v *Validator) IsAvailable() bool {
	return v != nil && v.db != nil
}

// ValidateGradeValue checks if the grade value matches the grade type requirements.
// In the Italian school system, grades range from 1 to 10 (with -1 indicating absence).
func (v *Validator) ValidateGradeValue(value float64, gradeType string) error {
	switch GradeType(gradeType) {
	case GradeTypeNumeric, "":
		if (value < 1.0 && value != -1.0) || value > 10.0 {
			return errors.New("voto deve essere tra 1 e 10, o -1 per assenza")
		}
	case GradeTypeJudgment:
		if (value < 1.0 && value != 0.0) || value > 10.0 {
			return errors.New("valore giudizio fuori range (0-10)")
		}
	case GradeTypeCredit:
		if value < 1.0 || value > 25.0 {
			return errors.New("credito scolastico deve essere compreso tra 1 e 25")
		}
	case GradeTypeCompetence:
		if value < 1 || value > 4 {
			return errors.New("livello competenza non valido (1-4)")
		}
	default:
		if (value < 1.0 && value != -1.0) || value > 10.0 {
			return errors.New("voto deve essere tra 1 e 10, o -1 per assenza")
		}
	}
	return nil
}

// ValidateJudgmentString checks that a judgment description is one of the accepted values.
func (v *Validator) ValidateJudgmentString(judgment string) error {
	allowed := []string{
		"Insufficiente", "Mediocre", "Sufficiente", "Discreto", "Buono",
		"Distinto", "Ottimo", "Eccellente", "Gravemente Insufficiente", "Quasi Sufficiente",
		"Avanzato", "Intermedio", "Base", "Iniziale", "Non Raggiunto",
	}
	for _, a := range allowed {
		if strings.EqualFold(a, judgment) {
			return nil
		}
	}
	return errors.New("giudizio non valido")
}

// ValidateTeacherCanGrade verifies that the teacher is assigned to teach the given
// subject in the given class.
func (v *Validator) ValidateTeacherCanGrade(teacherUserID string, subjectID string, classID string) error {
	assigned, err := v.IsTeacherAssignedToSubject(teacherUserID, subjectID, classID)
	if err != nil {
		return fmt.Errorf("teacher validation error: %w", err)
	}
	if !assigned {
		return errors.New("insegnante non assegnato a questa materia in questa classe")
	}
	return nil
}

// ValidateStudentEnrolled verifies that the student has an active enrollment
// for the given semester (checks class_students, not just student existence).
func (v *Validator) ValidateStudentEnrolled(studentID string, semester int) error {
	var query string
	var args []interface{}
	if semester > 0 {
		query = `
			SELECT 1
			FROM class_students cs
			JOIN students s ON (cs.student_id = s.id OR cs.student_id = s.user_id)
			WHERE (s.id = $1::uuid OR s.user_id = $1::uuid)
			  AND cs.status = 'active'
			  AND (cs.semester = $2 OR cs.semester IS NULL OR cs.semester = 0)`
		args = []interface{}{studentID, semester}
	} else {
		query = `
			SELECT 1
			FROM class_students cs
			JOIN students s ON (cs.student_id = s.id OR cs.student_id = s.user_id)
			WHERE (s.id = $1::uuid OR s.user_id = $1::uuid)
			  AND cs.status = 'active'`
		args = []interface{}{studentID}
	}
	var exists int
	err := v.db.QueryRow(query, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("studente non iscritto ad alcuna classe attiva (semestre %d)", semester)
		}
		return fmt.Errorf("student enrollment validation error: %w", err)
	}
	return nil
}

// ValidateGradeDate checks that the date is not in the future and falls within
// the expected semester date range.
func (v *Validator) ValidateGradeDate(date time.Time, semester int) error {
	now := time.Now()
	if date.After(now.Add(24 * time.Hour)) {
		return errors.New("data non può essere nel futuro")
	}
	if date.Before(now.AddDate(-10, 0, 0)) {
		return errors.New("data troppo vecchia")
	}

	start, end := GetSemesterDateRange(semester)
	if date.Before(start) || date.After(end) {
		return fmt.Errorf(
			"data fuori dal range del quadrimestre (Sem %d: %s - %s)",
			semester, start.Format("2006-01-02"), end.Format("2006-01-02"),
		)
	}
	return nil
}

// ValidateDescription checks description length only.
func (v *Validator) ValidateDescription(desc string) error {
	if len(desc) > 500 {
		return errors.New("descrizione max 500 caratteri")
	}
	return nil
}

// ValidateGradeNotLocked checks if grading is locked for the given semester.
func (v *Validator) ValidateGradeNotLocked(semester int) error {
	if IsInLockPeriod(semester, time.Now()) {
		return errors.New("periodo in chiusura scrutini, modifiche non consentite")
	}
	return nil
}

// --- Wrappers for Service ---

func (v *Validator) ValidateCreateRequest(req CreateGradeRequest, teacherID string) error {
	if err := v.ValidateGradeValue(req.GradeValue, req.GradeType); err != nil {
		return err
	}
	if err := v.ValidateDescription(req.Description); err != nil {
		return err
	}
	date, err := parseDate(req.Date)
	if err != nil {
		return fmt.Errorf("formato data non valido (atteso YYYY-MM-DD): %w", err)
	}
	if err := v.ValidateGradeDate(date, req.Semester); err != nil {
		return err
	}
	return nil
}

func (v *Validator) ValidateModification(grade Grade, req UpdateGradeRequest) error {
	if err := v.ValidateGradeNotLocked(int(grade.Semester)); err != nil {
		return err
	}
	if req.GradeValue != nil && req.GradeType != nil {
		if err := v.ValidateGradeValue(*req.GradeValue, *req.GradeType); err != nil {
			return err
		}
	} else if req.GradeValue != nil {
		if err := v.ValidateGradeValue(*req.GradeValue, string(grade.GradeType)); err != nil {
			return err
		}
	}
	if req.Description != nil {
		if err := v.ValidateDescription(*req.Description); err != nil {
			return err
		}
	}
	return nil
}

func parseDate(d string) (time.Time, error) {
	return time.Parse("2006-01-02", d)
}

// IsClassCoordinator checks if a teacher is the coordinator for the specified class.
func (v *Validator) IsClassCoordinator(teacherID string, classID string) (bool, error) {
	if teacherID == "" || classID == "" {
		return false, nil
	}
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM classes c
			LEFT JOIN teachers t ON (NULLIF(c.coordinator_id::text, '') = t.id::text OR NULLIF(c.coordinator_id::text, '') = t.user_id::text)
			WHERE c.id::text = $1 AND (c.coordinator_id::text = $2 OR t.id::text = $2 OR t.user_id::text = $2)
		)`
	err := v.db.QueryRow(query, classID, teacherID).Scan(&exists)
	return exists, err
}

// IsTeacherAssignedToSubject checks if a teacher is assigned to a class subject.
// Checks direct class_subjects assignment, class coordinator status, registered class lessons,
// active approved substitutions, or unassigned subjects in the same school.
func (v *Validator) IsTeacherAssignedToSubject(teacherID string, subjectID string, classID string) (bool, error) {
	if teacherID == "" || classID == "" || subjectID == "" {
		return false, nil
	}
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM class_subjects cs
			LEFT JOIN teachers t ON (NULLIF(cs.teacher_id::text, '') = t.id::text OR NULLIF(cs.teacher_id::text, '') = t.user_id::text)
			WHERE cs.class_id::text = $1 AND cs.subject_id::text = $2 
			  AND (NULLIF(cs.teacher_id::text, '') = $3 OR t.id::text = $3 OR t.user_id::text = $3)

			UNION ALL

			SELECT 1 FROM classes c
			LEFT JOIN teachers t ON (NULLIF(c.coordinator_id::text, '') = t.id::text OR NULLIF(c.coordinator_id::text, '') = t.user_id::text)
			WHERE c.id::text = $1 AND (c.coordinator_id::text = $3 OR t.id::text = $3 OR t.user_id::text = $3)

			UNION ALL

			SELECT 1 FROM class_lessons cl
			LEFT JOIN teachers t ON (NULLIF(cl.teacher_id::text, '') = t.id::text OR NULLIF(cl.teacher_id::text, '') = t.user_id::text)
			WHERE cl.class_id::text = $1 AND cl.subject_id::text = $2
			  AND (NULLIF(cl.teacher_id::text, '') = $3 OR t.id::text = $3 OR t.user_id::text = $3)

			UNION ALL

			SELECT 1 FROM substitutions s
			LEFT JOIN teachers t ON (NULLIF(s.substitute_teacher_id::text, '') = t.id::text OR NULLIF(s.substitute_teacher_id::text, '') = t.user_id::text)
			WHERE s.class_id::text = $1 AND s.subject_id::text = $2 
			  AND (NULLIF(s.substitute_teacher_id::text, '') = $3 OR t.id::text = $3 OR t.user_id::text = $3)
			  AND s.status IN ('assigned', 'confirmed')

			UNION ALL

			SELECT 1 FROM class_subjects cs
			JOIN classes c ON cs.class_id = c.id
			JOIN teachers t ON (t.id::text = $3 OR t.user_id::text = $3) AND t.school_id = c.school_id
			WHERE cs.class_id::text = $1 AND cs.subject_id::text = $2
			  AND (cs.teacher_id IS NULL OR cs.teacher_id::text = '')
		)`
	err := v.db.QueryRow(query, classID, subjectID, teacherID).Scan(&exists)
	return exists, err
}

// IsTeacherAssignedToStudent checks if a teacher is assigned to any of the student's class subjects or is class coordinator.
func (v *Validator) IsTeacherAssignedToStudent(ctx context.Context, teacherID string, studentID string) (bool, error) {
	if v == nil || v.db == nil || teacherID == "" || studentID == "" {
		return false, nil
	}
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM class_students cs
			JOIN class_subjects csub ON cs.class_id = csub.class_id
			LEFT JOIN teachers t ON (NULLIF(csub.teacher_id::text, '') = t.id::text OR NULLIF(csub.teacher_id::text, '') = t.user_id::text)
			WHERE cs.student_id::text = $1 AND (csub.teacher_id::text = $2 OR t.id::text = $2 OR t.user_id::text = $2)

			UNION ALL

			SELECT 1 FROM class_students cs
			JOIN classes c ON cs.class_id = c.id
			LEFT JOIN teachers t ON (NULLIF(c.coordinator_id::text, '') = t.id::text OR NULLIF(c.coordinator_id::text, '') = t.user_id::text)
			WHERE cs.student_id::text = $1 AND (c.coordinator_id::text = $2 OR t.id::text = $2 OR t.user_id::text = $2)
		)`
	err := v.db.QueryRowContext(ctx, query, studentID, teacherID).Scan(&exists)
	return exists, err
}

// IsTeacherAssignedToSubjectBySubjectID checks if a teacher teaches the given subject in any class.
// Unlike IsTeacherAssignedToSubject, no classID is required — used when only the subjectID is known (e.g. GetSubjectGrades).
func (v *Validator) IsTeacherAssignedToSubjectBySubjectID(ctx context.Context, teacherID, subjectID string) (bool, error) {
	if v == nil || v.db == nil || teacherID == "" || subjectID == "" {
		return false, nil
	}
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM class_subjects cs
			LEFT JOIN teachers t ON (NULLIF(cs.teacher_id::text, '') = t.id::text OR NULLIF(cs.teacher_id::text, '') = t.user_id::text)
			WHERE cs.subject_id::text = $1
			  AND (NULLIF(cs.teacher_id::text, '') = $2 OR t.id::text = $2 OR t.user_id::text = $2)
		)`
	err := v.db.QueryRowContext(ctx, query, subjectID, teacherID).Scan(&exists)
	return exists, err
}

// IsStudentInClass checks whether studentID belongs to classID.
func (v *Validator) IsStudentInClass(ctx context.Context, studentID string, classID string) (bool, error) {
	if v == nil || v.db == nil || studentID == "" || classID == "" {
		return true, nil
	}
	var exists bool
	query := `SELECT EXISTS(
		SELECT 1 FROM class_students WHERE student_id::text = $1 AND class_id::text = $2
	)`
	err := v.db.QueryRowContext(ctx, query, studentID, classID).Scan(&exists)
	return exists, err
}

// ─── Metodi specifici per la materia Religione Cattolica (IRC) ───────────────

// ValidateReligionJudgment verifica che il giudizio sia uno dei 6 valori ammessi per l'IRC.
// Restituisce il valore numerico interno corrispondente.
// Ammette sia il giudizio esatto ("Ottimo") sia con nota a seguire ("Ottimo - verifica").
func (v *Validator) ValidateReligionJudgment(judgment string) (float64, error) {
	trimmed := strings.TrimSpace(judgment)
	for j, val := range ReligionJudgmentValues {
		if strings.EqualFold(j, trimmed) {
			return val, nil
		}
	}
	for j, val := range ReligionJudgmentValues {
		if strings.HasPrefix(strings.ToLower(trimmed), strings.ToLower(j)) {
			rest := strings.TrimSpace(trimmed[len(j):])
			if rest == "" || strings.HasPrefix(rest, "-") || strings.HasPrefix(rest, ":") || strings.HasPrefix(rest, ",") {
				return val, nil
			}
		}
	}
	return 0, fmt.Errorf(
		"giudizio IRC non valido: '%s'. Valori ammessi: %s",
		judgment, strings.Join(ReligionJudgments, ", "),
	)
}

// IsReligionSubject verifica se la materia con l'ID indicato ha il flag is_religion = true.
func (v *Validator) IsReligionSubject(ctx context.Context, subjectID string) (bool, error) {
	if v == nil || v.db == nil || subjectID == "" {
		return false, nil
	}
	var isReligion bool
	query := `SELECT COALESCE(is_religion, FALSE) FROM subjects WHERE id = $1`
	err := v.db.QueryRowContext(ctx, query, subjectID).Scan(&isReligion)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("religion subject check error: %w", err)
	}
	return isReligion, nil
}

// IsStudentAvvalente verifica che lo studente si avvalga dell'IRC per la scuola indicata.
// Ritorna true se la scelta è 'avvalente', false altrimenti (non_avvalente, attivita_alternativa).
// Se non esiste ancora una scelta, si considera 'avvalente' per default (favor avvalimento).
func (v *Validator) IsStudentAvvalente(ctx context.Context, studentID, schoolID string) (bool, error) {
	if v == nil || v.db == nil || studentID == "" {
		return true, nil // default: avvalente
	}
	var choice string
	query := `
		SELECT choice::text
		FROM student_religion_choices
		WHERE (student_id = $1::uuid OR student_id IN (SELECT id FROM students WHERE user_id = $1::uuid OR id = $1::uuid))
		  AND school_id = $2::uuid
	`
	err := v.db.QueryRowContext(ctx, query, studentID, schoolID).Scan(&choice)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil // nessuna scelta registrata → default avvalente
		}
		return false, fmt.Errorf("religion choice check error: %w", err)
	}
	return choice == "avvalente", nil
}

// GetClassReligionChoices restituisce una mappa [studentID/userID] -> choice per tutti gli studenti di una classe.
func (v *Validator) GetClassReligionChoices(ctx context.Context, classID string) (map[string]string, error) {
	result := make(map[string]string)
	if v == nil || v.db == nil || classID == "" {
		return result, nil
	}
	query := `
		SELECT s.id::text, s.user_id::text, COALESCE(src.choice::text, 'avvalente')
		FROM students s
		LEFT JOIN student_religion_choices src ON (src.student_id = s.id AND src.school_id = s.school_id)
		WHERE s.class_id::text = $1
	`
	rows, err := v.db.QueryContext(ctx, query, classID)
	if err != nil {
		return result, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var sID, uID, choice string
		if err := rows.Scan(&sID, &uID, &choice); err != nil {
			return nil, err
		}
		result[sID] = choice
		result[uID] = choice
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
