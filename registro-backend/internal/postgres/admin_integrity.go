package postgres

import (
	"context"
	"time"

	"registro-backend/internal/admin"
)

// CheckDataIntegrity performs comprehensive diagnostics on school data
func (r *AdminRepository) CheckDataIntegrity(ctx context.Context, schoolID *string) (*admin.DataIntegrityReport, error) {
	checks := make([]admin.DataIntegrityIssue, 0, 5)

	// 1. Orphaned students (no class assigned)
	orphanQuery := `
		SELECT s.id::text, COALESCE(u.first_name, ''), COALESCE(u.last_name, ''), COALESCE(u.email, '')
		FROM students s
		JOIN users u ON s.user_id = u.id
		WHERE u.role = 'student' AND u.deleted_at IS NULL AND s.deleted_at IS NULL AND s.class_id IS NULL`
	var orphanArgs []interface{}
	if schoolID != nil {
		orphanQuery += ` AND s.school_id = $1`
		orphanArgs = append(orphanArgs, *schoolID)
	}
	orphanQuery += ` ORDER BY u.last_name, u.first_name LIMIT 50`

	orphanItems := make([]map[string]interface{}, 0)
	if rows, err := r.db.QueryContext(ctx, orphanQuery, orphanArgs...); err == nil {
		defer rows.Close()
		for rows.Next() {
			var id, firstName, lastName, email string
			if err := rows.Scan(&id, &firstName, &lastName, &email); err == nil {
				orphanItems = append(orphanItems, map[string]interface{}{
					"id":    id,
					"name":  firstName + " " + lastName,
					"email": email,
				})
			}
		}
	}
	checks = append(checks, admin.DataIntegrityIssue{
		ID:          "orphaned_students",
		Category:    "students",
		Severity:    "high",
		Title:       "Studenti senza classe assegnata",
		Description: "Studenti iscritti attivi che non risultano associati ad alcuna classe scolastica.",
		Count:       len(orphanItems),
		Items:       orphanItems,
	})

	// 2. Uncoordinated classes (classes with no coordinator teacher)
	uncoordQuery := `
		SELECT c.id::text, COALESCE(al.name || ' ', '') || c.section
		FROM classes c
		LEFT JOIN academic_levels al ON c.level_id = al.id
		WHERE c.deleted_at IS NULL AND c.coordinator_id IS NULL`
	var uncoordArgs []interface{}
	if schoolID != nil {
		uncoordQuery += ` AND c.school_id = $1`
		uncoordArgs = append(uncoordArgs, *schoolID)
	}
	uncoordQuery += ` ORDER BY c.section LIMIT 50`

	uncoordItems := make([]map[string]interface{}, 0)
	if rows, err := r.db.QueryContext(ctx, uncoordQuery, uncoordArgs...); err == nil {
		defer rows.Close()
		for rows.Next() {
			var id, className string
			if err := rows.Scan(&id, &className); err == nil {
				uncoordItems = append(uncoordItems, map[string]interface{}{
					"id":    id,
					"class": className,
				})
			}
		}
	}
	checks = append(checks, admin.DataIntegrityIssue{
		ID:          "uncoordinated_classes",
		Category:    "classes",
		Severity:    "medium",
		Title:       "Classi senza docente coordinatore",
		Description: "Classi prive di coordinatore assegnato, necessario per la gestione scrutini e comunicazioni.",
		Count:       len(uncoordItems),
		Items:       uncoordItems,
	})

	// 3. Overlapping lessons (teacher assigned in two different classes at the same hour)
	overlapQuery := `
		SELECT l1.id::text, COALESCE(u.first_name || ' ' || u.last_name, 'Docente'),
		       TO_CHAR(l1.date, 'YYYY-MM-DD'), COALESCE(l1.hour, 1),
		       COALESCE(al1.name || ' ', '') || c1.section,
		       COALESCE(al2.name || ' ', '') || c2.section
		FROM class_lessons l1
		JOIN class_lessons l2 ON l1.teacher_id = l2.teacher_id 
		                     AND l1.date = l2.date 
		                     AND COALESCE(l1.hour, 1) = COALESCE(l2.hour, 1) 
		                     AND l1.class_id != l2.class_id
		                     AND l1.id < l2.id
		JOIN users u ON l1.teacher_id = u.id
		JOIN classes c1 ON l1.class_id = c1.id
		JOIN classes c2 ON l2.class_id = c2.id
		LEFT JOIN academic_levels al1 ON c1.level_id = al1.id
		LEFT JOIN academic_levels al2 ON c2.level_id = al2.id
		WHERE (l1.is_co_teaching IS NOT TRUE OR l2.is_co_teaching IS NOT TRUE)`
	var overlapArgs []interface{}
	if schoolID != nil {
		overlapQuery += ` AND c1.school_id = $1`
		overlapArgs = append(overlapArgs, *schoolID)
	}
	overlapQuery += ` LIMIT 50`

	overlapItems := make([]map[string]interface{}, 0)
	if rows, err := r.db.QueryContext(ctx, overlapQuery, overlapArgs...); err == nil {
		defer rows.Close()
		for rows.Next() {
			var id, teacherName, date string
			var hour int
			var class1, class2 string
			if err := rows.Scan(&id, &teacherName, &date, &hour, &class1, &class2); err == nil {
				overlapItems = append(overlapItems, map[string]interface{}{
					"id":      id,
					"teacher": teacherName,
					"date":    date,
					"hour":    hour,
					"class1":  class1,
					"class2":  class2,
				})
			}
		}
	}
	checks = append(checks, admin.DataIntegrityIssue{
		ID:          "overlapping_lessons",
		Category:    "lessons",
		Severity:    "high",
		Title:       "Lezioni sovrapposte per lo stesso docente",
		Description: "Docenti con due o più lezioni firmate contemporaneamente nella stessa data e ora in classi differenti.",
		Count:       len(overlapItems),
		Items:       overlapItems,
	})

	// 4. Weekend grades (grades recorded on Sunday)
	weekendQuery := `
		SELECT g.id::text, COALESCE(u.first_name || ' ' || u.last_name, 'Studente'),
		       COALESCE(sub.name, 'Materia'), TO_CHAR(g.date, 'YYYY-MM-DD'), g.grade_value
		FROM grades g
		JOIN students s ON g.student_id = s.id
		JOIN users u ON s.user_id = u.id
		LEFT JOIN subjects sub ON g.subject_id = sub.id
		WHERE g.deleted_at IS NULL AND EXTRACT(DOW FROM g.date) = 0`
	var weekendArgs []interface{}
	if schoolID != nil {
		weekendQuery += ` AND g.school_id = $1`
		weekendArgs = append(weekendArgs, *schoolID)
	}
	weekendQuery += ` LIMIT 50`

	weekendItems := make([]map[string]interface{}, 0)
	if rows, err := r.db.QueryContext(ctx, weekendQuery, weekendArgs...); err == nil {
		defer rows.Close()
		for rows.Next() {
			var id, studentName, subjectName, date string
			var gradeValue float64
			if err := rows.Scan(&id, &studentName, &subjectName, &date, &gradeValue); err == nil {
				weekendItems = append(weekendItems, map[string]interface{}{
					"id":      id,
					"student": studentName,
					"subject": subjectName,
					"date":    date,
					"grade":   gradeValue,
				})
			}
		}
	}
	checks = append(checks, admin.DataIntegrityIssue{
		ID:          "weekend_grades",
		Category:    "grades",
		Severity:    "medium",
		Title:       "Voti registrati di domenica",
		Description: "Valutazioni con data domenicale, possibile errore di digitazione durante l'inserimento.",
		Count:       len(weekendItems),
		Items:       weekendItems,
	})

	// 5. Unlinked guardians (students without parent link)
	unlinkedQuery := `
		SELECT s.id::text, COALESCE(u.first_name || ' ' || u.last_name, 'Studente'),
		       COALESCE(al.name || ' ', '') || COALESCE(c.section, 'Senza classe')
		FROM students s
		JOIN users u ON s.user_id = u.id
		LEFT JOIN classes c ON s.class_id = c.id
		LEFT JOIN academic_levels al ON c.level_id = al.id
		WHERE u.role = 'student' AND u.deleted_at IS NULL AND s.deleted_at IS NULL
		  AND NOT EXISTS (SELECT 1 FROM student_parents sp WHERE sp.student_id = s.id)`
	var unlinkedArgs []interface{}
	if schoolID != nil {
		unlinkedQuery += ` AND s.school_id = $1`
		unlinkedArgs = append(unlinkedArgs, *schoolID)
	}
	unlinkedQuery += ` LIMIT 50`

	unlinkedItems := make([]map[string]interface{}, 0)
	if rows, err := r.db.QueryContext(ctx, unlinkedQuery, unlinkedArgs...); err == nil {
		defer rows.Close()
		for rows.Next() {
			var id, studentName, className string
			if err := rows.Scan(&id, &studentName, &className); err == nil {
				unlinkedItems = append(unlinkedItems, map[string]interface{}{
					"id":      id,
					"student": studentName,
					"class":   className,
				})
			}
		}
	}
	checks = append(checks, admin.DataIntegrityIssue{
		ID:          "unlinked_guardians",
		Category:    "guardians",
		Severity:    "medium",
		Title:       "Studenti senza genitore/tutore collegato",
		Description: "Studenti per cui non risulta associato alcun genitore o tutore legale abilitato alle notifiche.",
		Count:       len(unlinkedItems),
		Items:       unlinkedItems,
	})

	// Calculate score
	score := 100
	totalIssues := 0
	for _, chk := range checks {
		totalIssues += chk.Count
		penalty := 0
		switch chk.Severity {
		case "high":
			penalty = chk.Count * 10
			if penalty > 35 {
				penalty = 35
			}
		case "medium":
			penalty = chk.Count * 5
			if penalty > 20 {
				penalty = 20
			}
		case "low":
			penalty = chk.Count * 2
			if penalty > 10 {
				penalty = 10
			}
		}
		score -= penalty
	}
	if score < 0 {
		score = 0
	}

	healthStatus := "healthy"
	if score < 70 {
		healthStatus = "critical"
	} else if score < 90 {
		healthStatus = "warning"
	}

	return &admin.DataIntegrityReport{
		Score:        score,
		HealthStatus: healthStatus,
		TotalIssues:  totalIssues,
		Checks:       checks,
		RunAt:        time.Now(),
	}, nil
}
