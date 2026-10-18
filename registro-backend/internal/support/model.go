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
	ID                   string    `json:"id"`
	SchoolID             string    `json:"school_id"`
	StudentID            string    `json:"student_id"`
	StudentName          string    `json:"student_name,omitempty"`
	PeiType              string    `json:"pei_type"`              // 'equipollente', 'differenziato'
	Axis                 string    `json:"axis"`                  // 'autonomia', 'cognitiva', 'comunicazionale', 'relazionale', 'linguistica', 'sensoriale'
	MinisterialDimension string    `json:"ministerial_dimension"` // 'dimensione_relazione', 'dimensione_comunicazione', 'dimensione_autonomia', 'dimensione_cognitiva' (D.I. 182/2020)
	PathwayType          string    `json:"pathway_type"`          // 'percorso_a_ordinario', 'percorso_b_personalizzato', 'percorso_c_differenziato'
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	ExpectedTerm         string    `json:"expected_term"`   // 'q1', 'q2', 'annuale'
	ProgressStatus       string    `json:"progress_status"` // 'non_avviato', 'iniziale', 'intermedio', 'avanzato', 'raggiunto'
	GloNotes             string    `json:"glo_notes"`       // Note e decisioni del GLO (Gruppo di Lavoro Operativo)
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
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
	StudentID            string `json:"student_id" binding:"required"`
	PeiType              string `json:"pei_type" binding:"required"`
	Axis                 string `json:"axis" binding:"required"`
	MinisterialDimension string `json:"ministerial_dimension"`
	PathwayType          string `json:"pathway_type"`
	Title                string `json:"title" binding:"required"`
	Description          string `json:"description"`
	ExpectedTerm         string `json:"expected_term"`
	ProgressStatus       string `json:"progress_status"`
	GloNotes             string `json:"glo_notes"`
}

// SupportGloMeeting represents a formal Gruppo di Lavoro Operativo meeting
type SupportGloMeeting struct {
	ID                string         `json:"id"`
	SchoolID          string         `json:"school_id"`
	StudentID         string         `json:"student_id"`
	StudentName       string         `json:"student_name,omitempty"`
	MeetingDate       string         `json:"meeting_date"` // YYYY-MM-DD
	MeetingType       string         `json:"meeting_type"` // 'PROVVISORIO', 'INIZIALE', 'INTERMEDIO', 'FINALE'
	SchoolYear        string         `json:"school_year"`  // e.g. '2025/2026'
	MembersPresent    []string       `json:"members_present"`
	DiscussionSummary string         `json:"discussion_summary"`
	DecisionsTaken    string         `json:"decisions_taken"`
	Status            string         `json:"status"` // 'BOZZA', 'IN_FIRMA', 'APPROVATO_CON_FIRME', 'ARCHIVIATO'
	Signatures        []GloSignature `json:"signatures,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

// GloSignature represents a digital or SPID-authenticated signature by a GLO member
type GloSignature struct {
	ID            string     `json:"id"`
	MeetingID     string     `json:"meeting_id"`
	MemberID      string     `json:"member_id"`
	MemberName    string     `json:"member_name"`
	MemberRole    string     `json:"member_role"` // 'DOCENTE_SOSTEGNO', 'DOCENTE_CURRICULARE', 'GENITORE', 'SPECIALISTA_ASL', 'EDUCATORE_OEPA'
	SignedAt      *time.Time `json:"signed_at,omitempty"`
	SignatureType string     `json:"signature_type"` // 'DIGITAL_AGID', 'SPID_OTP', 'PIN_ACCETTAZIONE'
	SignatureHash string     `json:"signature_hash,omitempty"`
	Status        string     `json:"status"` // 'PENDING', 'SIGNED', 'REJECTED'
}

// SupportHoursQuantification handles Sezioni 11 e 12 del PEI Ministeriale (D.I. 182/2020 e D.M. 153/2023)
type SupportHoursQuantification struct {
	ID                          string    `json:"id"`
	SchoolID                    string    `json:"school_id"`
	StudentID                   string    `json:"student_id"`
	SchoolYear                  string    `json:"school_year"`
	PeiStage                    string    `json:"pei_stage"`          // 'PROVVISORIO', 'DEFINITIVO', 'FINALE'
	DebtRelazione               int       `json:"debt_relazione"`     // 0=Assente, 1=Lieve, 2=Medio, 3=Elevato, 4=Molto Elevato
	DebtComunicazione           int       `json:"debt_comunicazione"` // 0-4
	DebtAutonomia               int       `json:"debt_autonomia"`     // 0-4
	DebtCognitiva               int       `json:"debt_cognitiva"`     // 0-4
	WeeklySupportHoursRequested float64   `json:"weekly_support_hours_requested"`
	WeeklySupportHoursAssigned  float64   `json:"weekly_support_hours_assigned"`
	AssistantHoursRequested     float64   `json:"assistant_hours_requested"` // Ore educatore OEPA / ASACOM
	AssistantHoursAssigned      float64   `json:"assistant_hours_assigned"`
	TransportationNeeded        bool      `json:"transportation_needed"`
	PersonalCareNeeded          bool      `json:"personal_care_needed"`
	JustificationSummary        string    `json:"justification_summary"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}

type CreateGloMeetingRequest struct {
	StudentID         string           `json:"student_id" binding:"required"`
	MeetingDate       string           `json:"meeting_date" binding:"required"`
	MeetingType       string           `json:"meeting_type" binding:"required"`
	SchoolYear        string           `json:"school_year"`
	MembersPresent    []string         `json:"members_present"`
	DiscussionSummary string           `json:"discussion_summary" binding:"required"`
	DecisionsTaken    string           `json:"decisions_taken"`
	Signers           []GloSignerDraft `json:"signers,omitempty"`
}

type GloSignerDraft struct {
	MemberID   string `json:"member_id"`
	MemberName string `json:"member_name"`
	MemberRole string `json:"member_role"`
}

type SignGloMeetingRequest struct {
	MeetingID     string `json:"meeting_id" binding:"required"`
	MemberRole    string `json:"member_role" binding:"required"`
	SignatureType string `json:"signature_type"` // 'DIGITAL_AGID', 'SPID_OTP', 'PIN_ACCETTAZIONE'
	PinOrOtp      string `json:"pin_or_otp"`
}

type SaveSupportHoursRequest struct {
	StudentID                   string  `json:"student_id" binding:"required"`
	SchoolYear                  string  `json:"school_year"`
	PeiStage                    string  `json:"pei_stage" binding:"required"`
	DebtRelazione               int     `json:"debt_relazione"`
	DebtComunicazione           int     `json:"debt_comunicazione"`
	DebtAutonomia               int     `json:"debt_autonomia"`
	DebtCognitiva               int     `json:"debt_cognitiva"`
	WeeklySupportHoursRequested float64 `json:"weekly_support_hours_requested"`
	WeeklySupportHoursAssigned  float64 `json:"weekly_support_hours_assigned"`
	AssistantHoursRequested     float64 `json:"assistant_hours_requested"`
	AssistantHoursAssigned      float64 `json:"assistant_hours_assigned"`
	TransportationNeeded        bool    `json:"transportation_needed"`
	PersonalCareNeeded          bool    `json:"personal_care_needed"`
	JustificationSummary        string  `json:"justification_summary"`
}
