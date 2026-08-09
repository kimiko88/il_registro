package grades

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Validator struct {
	db *sql.DB
}

func NewValidator(db *sql.DB) *Validator {
	return &Validator{db: db}
}

// ValidateGradeValue checks if the grade value matches the grade type requirements.
func (v *Validator) ValidateGradeValue(value float64, gradeType string) error {
	switch GradeType(gradeType) {
	case GradeTypeNumeric, "":
		if (value < 0 && value != -1) || value > 10 {
			return errors.New("Voto deve essere tra 0 e 10, o -1 per assenza")
		}
	case GradeTypeJudgment:
		if value < 0 || value > 10 {
			return errors.New("Valore giudizio fuori range")
		}
	case GradeTypeCredit:
		if value <= 0 {
			return errors.New("Credito deve essere positivo")
		}
	case GradeTypeCompetence:
		if value < 1 || value > 4 {
			return errors.New("Livello competenza non valido (1-4)")
		}
	default:
		if (value < 0 && value != -1) || value > 10 {
			return errors.New("Voto deve essere tra 0 e 10, o -1 per assenza")
		}
	}
	return nil
}

// ValidateJudgmentString checks that a judgment description is one of the accepted values.
func (v *Validator) ValidateJudgmentString(judgment string) error {
	allowed := []string{
		"Insufficiente", "Mediocre", "Sufficiente", "Discreto", "Buono",
		"Distinto", "Ottimo", "Eccellente", "Gravemente Insufficiente", "Quasi Sufficiente",
	}
	for _, a := range allowed {
		if strings.EqualFold(a, judgment) {
			return nil
		}
	}
	return errors.New("Giudizio non valido")
}

// ValidateTeacherCanGrade verifies that the teacher is assigned to teach the given
// subject in the given class.
func (v *Validator) ValidateTeacherCanGrade(teacherUserID string, subjectID string, classID string) error {
	assigned, err := v.IsTeacherAssignedToSubject(teacherUserID, subjectID, classID)
	if err != nil {
		return fmt.Errorf("teacher validation error: %w", err)
	}
	if !assigned {
		return errors.New("Insegnante non assegnato a questa materia in questa classe")
	}
	return nil
}

// ValidateStudentEnrolled verifies that the student has an active enrollment
// for the given semester (checks class_students, not just student existence).
func (v *Validator) ValidateStudentEnrolled(studentID string, semester int) error {
	query := `
		SELECT 1
		FROM class_students cs
		JOIN students s ON cs.student_id = s.id
		WHERE s.id = $1
		  AND cs.status = 'active'`
	var exists int
	err := v.db.QueryRow(query, studentID).Scan(&exists)
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
		return errors.New("Data non può essere nel futuro")
	}
	if date.Before(now.AddDate(-10, 0, 0)) {
		return errors.New("Data troppo vecchia")
	}

	start, end := GetSemesterDateRange(semester)
	if date.Before(start) || date.After(end) {
		return fmt.Errorf(
			"Data fuori dal range del quadrimestre (Sem %d: %s - %s)",
			semester, start.Format("2006-01-02"), end.Format("2006-01-02"),
		)
	}
	return nil
}

// ValidateDescription checks description length only.
func (v *Validator) ValidateDescription(desc string) error {
	if len(desc) > 500 {
		return errors.New("Descrizione max 500 caratteri")
	}
	return nil
}

// ValidateGradeNotLocked checks if grading is locked for the given semester.
func (v *Validator) ValidateGradeNotLocked(semester int) error {
	if IsInLockPeriod(semester, time.Now()) {
		return errors.New("Periodo in chiusura scrutini, modifiche non consentite")
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
