package reports

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"math"
	"sort"
	"strings"
)

type DropoutRiskItem struct {
	StudentID            string   `json:"student_id"`
	FirstName            string   `json:"first_name"`
	LastName             string   `json:"last_name"`
	ClassID              string   `json:"class_id"`
	ClassName            string   `json:"class_name"`
	AbsenceRate          float64  `json:"absence_rate"`
	TotalHours           int      `json:"total_hours"`
	AbsenceHours         int      `json:"absence_hours"`
	FailingSubjectsCount int      `json:"failing_subjects_count"`
	FailingSubjects      []string `json:"failing_subjects"`
	LatesCount           int      `json:"lates_count"`
	EarlyExitsCount      int      `json:"early_exits_count"`
	RiskScore            float64  `json:"risk_score"`
	RiskLevel            string   `json:"risk_level"` // Critico, Alto, Moderato, Basso
	RecommendedAction    string   `json:"recommended_action"`
}

// CalculateStudentRisk computes early warning risk metrics according to DPR 122/2009 art. 14
func CalculateStudentRisk(absenceRate float64, failingCount int, delays int) (score float64, level string, action string) {
	score = math.Round(((absenceRate*2.0)+float64(failingCount*12)+(float64(delays)*1.2))*10) / 10
	if score > 100 {
		score = 100
	}

	if absenceRate >= 25.0 || score >= 65 || (absenceRate >= 20.0 && failingCount >= 3) {
		level = "Critico"
		action = "Convocazione urgente famiglia e attivazione tutor dedicato per deroga/piano intensivo (DPR 122/2009)"
	} else if absenceRate >= 20.0 || failingCount >= 3 || score >= 48 {
		level = "Alto"
		action = "Sportello didattico di recupero, patto formativo e monitoraggio bisettimanale presenze"
	} else if absenceRate >= 15.0 || failingCount >= 1 || score >= 20 {
		level = "Moderato"
		action = "Colloquio con docente coordinatore e studio assistito nelle materie insufficienti"
	} else {
		level = "Basso"
		action = "Monitoraggio ordinario dell'andamento didattico e delle presenze"
	}
	return score, level, action
}

func (s *Service) GetDropoutRisk(ctx context.Context, classID, riskFilter string) ([]DropoutRiskItem, error) {
	if s.db == nil {
		return []DropoutRiskItem{}, nil
	}

	// 1. Fetch Students via the students table (users.class_id does not exist;
	//    the class association is in students.class_id → classes.id)
	userQuery := `
		SELECT u.id, u.first_name, u.last_name,
		       COALESCE(st.class_id::text, '') AS class_id,
		       COALESCE(c.name || ' ' || COALESCE(c.section, ''), c.name, 'N/D') AS class_name,
		       st.id AS student_record_id
		FROM users u
		INNER JOIN students st ON st.user_id = u.id AND st.deleted_at IS NULL
		LEFT JOIN classes c ON c.id = st.class_id
		WHERE u.role = 'student' AND u.deleted_at IS NULL
	`
	var userArgs []interface{}
	if classID != "" {
		userQuery += " AND st.class_id = $1"
		userArgs = append(userArgs, classID)
	}
	userQuery += " ORDER BY u.last_name ASC, u.first_name ASC"

	rows, err := s.db.QueryContext(ctx, userQuery, userArgs...)
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query students error: %w", err)
	}
	defer rows.Close()

	type studentBasic struct {
		userID          string
		firstName       string
		lastName        string
		classID         string
		className       string
		studentRecordID string // students.id – used as FK in attendance and grades
	}
	var students []studentBasic

	for rows.Next() {
		var sb studentBasic
		if err := rows.Scan(&sb.userID, &sb.firstName, &sb.lastName, &sb.classID, &sb.className, &sb.studentRecordID); err != nil {
			return nil, err
		}
		students = append(students, sb)
	}
	if len(students) == 0 {
		return []DropoutRiskItem{}, nil
	}

	// 2. Fetch Attendance Aggregates per students.id
	type attAgg struct {
		total   int
		absence int
		late    int
		exit    int
	}
	attMap := make(map[string]*attAgg) // keyed by students.id

	attQuery := `
		SELECT student_id,
		       COUNT(*) AS total_rec,
		       COUNT(CASE WHEN status = 'Absent' THEN 1 END) AS abs_count,
		       COUNT(CASE WHEN status = 'Late' THEN 1 END) AS late_count,
		       COUNT(CASE WHEN status = 'LeftEarly' THEN 1 END) AS exit_count
		FROM attendance
		WHERE deleted_at IS NULL
	`
	var attArgs []interface{}
	if classID != "" {
		attQuery += " AND class_id = $1"
		attArgs = append(attArgs, classID)
	}
	attQuery += " GROUP BY student_id"

	attRows, err := s.db.QueryContext(ctx, attQuery, attArgs...)
	if err == nil {
		defer attRows.Close()
		for attRows.Next() {
			var sID string
			var a attAgg
			if err := attRows.Scan(&sID, &a.total, &a.absence, &a.late, &a.exit); err == nil {
				attMap[sID] = &a
			}
		}
	}

	// 3. Fetch Grades Average < 5.0 per students.id per subject
	//    grades.student_id references students.id; grade column is grade_value
	gradeMap := make(map[string][]string) // students.id -> list of failing subjects

	gradeQuery := `
		SELECT g.student_id, COALESCE(s.name, 'Materia') AS subject_name, AVG(g.grade_value) AS avg_val
		FROM grades g
		LEFT JOIN subjects s ON g.subject_id = s.id
		WHERE g.deleted_at IS NULL
		GROUP BY g.student_id, s.name
		HAVING AVG(g.grade_value) < 5.0
	`
	gradeRows, err := s.db.QueryContext(ctx, gradeQuery)
	if err == nil {
		defer gradeRows.Close()
		for gradeRows.Next() {
			var sID, sName string
			var avg float64
			if err := gradeRows.Scan(&sID, &sName, &avg); err == nil {
				gradeMap[sID] = append(gradeMap[sID], fmt.Sprintf("%s (%.1f)", sName, avg))
			}
		}
	}

	// 4. Build Result List
	var results []DropoutRiskItem
	for _, st := range students {
		// attMap and gradeMap are keyed by students.id (the student record primary key)
		agg := attMap[st.studentRecordID]
		totalHours := 0
		absHours := 0
		lates := 0
		earlyExits := 0
		absenceRate := 0.0

		if agg != nil && agg.total > 0 {
			totalHours = agg.total
			absHours = agg.absence
			lates = agg.late
			earlyExits = agg.exit
			absenceRate = math.Round((float64(absHours)/float64(totalHours))*1000) / 10
		}

		failingSubs := gradeMap[st.studentRecordID]
		if failingSubs == nil {
			failingSubs = []string{}
		}

		score, level, action := CalculateStudentRisk(absenceRate, len(failingSubs), lates+earlyExits)

		if riskFilter != "" && !strings.EqualFold(riskFilter, level) {
			continue
		}

		results = append(results, DropoutRiskItem{
			StudentID:            st.userID, // expose users.id to the frontend
			FirstName:            st.firstName,
			LastName:             st.lastName,
			ClassID:              st.classID,
			ClassName:            st.className,
			AbsenceRate:          absenceRate,
			TotalHours:           totalHours,
			AbsenceHours:         absHours,
			FailingSubjectsCount: len(failingSubs),
			FailingSubjects:      failingSubs,
			LatesCount:           lates,
			EarlyExitsCount:      earlyExits,
			RiskScore:            score,
			RiskLevel:            level,
			RecommendedAction:    action,
		})
	}

	// Sort by RiskScore desc, then LastName asc
	sort.Slice(results, func(i, j int) bool {
		if results[i].RiskScore != results[j].RiskScore {
			return results[i].RiskScore > results[j].RiskScore
		}
		return results[i].LastName < results[j].LastName
	})

	return results, nil
}

func (s *Service) ExportDropoutRiskCSV(ctx context.Context, classID, riskFilter string) ([]byte, error) {
	items, err := s.GetDropoutRisk(ctx, classID, riskFilter)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	// UTF-8 BOM for Excel compatibility
	buf.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(&buf)
	writer.Comma = ';'

	header := []string{
		"ID Studente",
		"Cognome",
		"Nome",
		"Classe",
		"Tasso Assenze (%)",
		"Ore Totali",
		"Ore Assenza",
		"Materie Grav. Insufficienti (<5.0)",
		"Ritardi",
		"Uscite Anticipate",
		"Punteggio Rischio (0-100)",
		"Livello Rischio",
		"Piano di Supporto Raccomandato (DPR 122/2009)",
	}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	for _, it := range items {
		failSubsStr := strings.Join(it.FailingSubjects, ", ")
		if failSubsStr == "" {
			failSubsStr = "Nessuna"
		}

		row := []string{
			it.StudentID,
			it.LastName,
			it.FirstName,
			it.ClassName,
			fmt.Sprintf("%.1f%%", it.AbsenceRate),
			fmt.Sprintf("%d", it.TotalHours),
			fmt.Sprintf("%d", it.AbsenceHours),
			failSubsStr,
			fmt.Sprintf("%d", it.LatesCount),
			fmt.Sprintf("%d", it.EarlyExitsCount),
			fmt.Sprintf("%.1f", it.RiskScore),
			it.RiskLevel,
			it.RecommendedAction,
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
