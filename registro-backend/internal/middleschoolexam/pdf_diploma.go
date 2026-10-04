package middleschoolexam

import (
	"fmt"
	"strings"
	"time"
)

type DiplomaData struct {
	SchoolName      string    `json:"school_name"`
	SchoolCode      string    `json:"school_code"`
	AcademicYear    string    `json:"academic_year"`
	StudentFullName string    `json:"student_full_name"`
	BirthDate       string    `json:"birth_date"`
	BirthPlace      string    `json:"birth_place"`
	FinalGrade      int       `json:"final_grade"`
	HasHonors       bool      `json:"has_honors"`
	PresidentName   string    `json:"president_name"`
	IssueDate       time.Time `json:"issue_date"`
}

func GenerateDiplomaCertificateText(d DiplomaData) string {
	honorsStr := ""
	if d.HasHonors {
		honorsStr = " E LODE"
	}

	dateStr := d.IssueDate.Format("02/01/2006")

	var sb strings.Builder
	sb.WriteString("================================================================================\n")
	sb.WriteString("                  REPUBBLICA ITALIANA - MINISTERO DELL'ISTRUZIONE               \n")
	sb.WriteString(fmt.Sprintf("                        %s\n", strings.ToUpper(d.SchoolName)))
	sb.WriteString(fmt.Sprintf("                           Codice Meccanografico: %s\n", d.SchoolCode))
	sb.WriteString("================================================================================\n\n")
	sb.WriteString("               DIPLOMA CONCLUSIVO DEL PRIMO CICLO DI ISTRUZIONE                \n")
	sb.WriteString("                           (D.Lgs. 13 aprile 2017, n. 62)                      \n\n")
	sb.WriteString(fmt.Sprintf("Si certifica che l'alunno/a: %s\n", strings.ToUpper(d.StudentFullName)))
	sb.WriteString(fmt.Sprintf("nato/a a %s il %s\n\n", d.BirthPlace, d.BirthDate))
	sb.WriteString(fmt.Sprintf("ha superato nell'anno scolastico %s l'Esame di Stato conclusivo\n", d.AcademicYear))
	sb.WriteString("del primo ciclo di istruzione con la votazione finale di:\n\n")
	sb.WriteString(fmt.Sprintf("                              >>> %d/10%s <<<\n\n", d.FinalGrade, honorsStr))
	sb.WriteString("Ai sensi dell'art. 9 del D.Lgs. 62/2017 e del D.M. 741/2017.\n\n")
	sb.WriteString(fmt.Sprintf("Data di rilascio: %s\n\n", dateStr))
	sb.WriteString(fmt.Sprintf("Il Presidente della Commissione: %s\n", d.PresidentName))
	sb.WriteString("Il Dirigente Scolastico: _____________________________\n")
	sb.WriteString("================================================================================\n")

	return sb.String()
}
