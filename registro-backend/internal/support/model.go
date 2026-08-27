package support

import "time"

type SupportDiaryEntry struct {
	ID                 string    `json:"id"`
	SchoolID           string    `json:"school_id"`
	TeacherID          string    `json:"teacher_id"`
	TeacherName        string    `json:"teacher_name,omitempty"`
	StudentID          string    `json:"student_id"`
	StudentName        string    `json:"student_name,omitempty"`
	ClassID            string    `json:"class_id"`
	ClassName          string    `json:"class_name,omitempty"`
	EntryDate          string    `json:"entry_date"` // YYYY-MM-DD
	TimeSlot           string    `json:"time_slot"`  // e.g. "1ª Ora (08:00-09:00)"
	CoTeacherID        *string   `json:"co_teacher_id,omitempty"`
	CoTeacherName      string    `json:"co_teacher_name,omitempty"`
	ActivityType       string    `json:"activity_type"` // 'in_classe', 'laboratorio', 'aula_sostegno', 'individuale', 'piccolo_gruppo'
	TopicAndActivities string    `json:"topic_and_activities"`
	StudentResponses   string    `json:"student_responses"` // Livello partecipazione/autonomia
	EducatorNotes      string    `json:"educator_notes"`    // Note per educatore OEPA / ASACOM
	IsSharedWithFamily bool      `json:"is_shared_with_family"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type SupportPeiGoal struct {
	ID             string    `json:"id"`
	SchoolID       string    `json:"school_id"`
	StudentID      string    `json:"student_id"`
	StudentName    string    `json:"student_name,omitempty"`
	PeiType        string    `json:"pei_type"` // 'equipollente', 'differenziato'
	Axis           string    `json:"axis"`     // 'autonomia', 'cognitiva', 'comunicazionale', 'relazionale', 'linguistica', 'sensoriale'
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	ExpectedTerm   string    `json:"expected_term"`   // 'q1', 'q2', 'annuale'
	ProgressStatus string    `json:"progress_status"` // 'non_avviato', 'iniziale', 'intermedio', 'avanzato', 'raggiunto'
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateDiaryEntryRequest struct {
	StudentID          string  `json:"student_id" binding:"required"`
	ClassID            string  `json:"class_id" binding:"required"`
	EntryDate          string  `json:"entry_date" binding:"required"`
	TimeSlot           string  `json:"time_slot" binding:"required"`
	CoTeacherID        *string `json:"co_teacher_id"`
	ActivityType       string  `json:"activity_type"`
	TopicAndActivities string  `json:"topic_and_activities" binding:"required"`
	StudentResponses   string  `json:"student_responses"`
	EducatorNotes      string  `json:"educator_notes"`
	IsSharedWithFamily bool    `json:"is_shared_with_family"`
}

type CreatePeiGoalRequest struct {
	StudentID      string `json:"student_id" binding:"required"`
	PeiType        string `json:"pei_type" binding:"required"`
	Axis           string `json:"axis" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Description    string `json:"description"`
	ExpectedTerm   string `json:"expected_term"`
	ProgressStatus string `json:"progress_status"`
}
