package infanzia

import "time"

type InfanziaSection struct {
	ID           string    `json:"id" db:"id"`
	SchoolID     string    `json:"school_id" db:"school_id"`
	SectionName  string    `json:"section_name" db:"section_name"` // Sezione Rossi, Blu, ecc.
	ColorHex     string    `json:"color_hex" db:"color_hex"`       // #EF4444
	Room         string    `json:"room" db:"room"`
	AcademicYear string    `json:"academic_year" db:"academic_year"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`

	StudentCount int `json:"student_count"`
}

type InfanziaDailyAttendance struct {
	ID              string    `json:"id" db:"id"`
	SectionID       string    `json:"section_id" db:"section_id"`
	StudentID       string    `json:"student_id" db:"student_id"`
	AttendanceDate  string    `json:"attendance_date" db:"attendance_date"` // YYYY-MM-DD
	IsPresent       bool      `json:"is_present" db:"is_present"`
	MealBooked      bool      `json:"meal_booked" db:"meal_booked"`
	SpecialDietNote string    `json:"special_diet_note" db:"special_diet_note"`
	EntryTime       *string   `json:"entry_time,omitempty" db:"entry_time"` // HH:MM
	ExitTime        *string   `json:"exit_time,omitempty" db:"exit_time"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`

	StudentFullName string `json:"student_full_name,omitempty"`
}

type InfanziaExperienceObservation struct {
	ID               string    `json:"id" db:"id"`
	StudentID        string    `json:"student_id" db:"student_id"`
	CampoEsperienza  string    `json:"campo_esperienza" db:"campo_esperienza"` // il_se_e_l_altro, il_corpo_e_il_movimento, immagini_suoni_colori, i_discorsi_e_le_parole, la_conoscenza_del_mondo
	ObservationDate  string    `json:"observation_date" db:"observation_date"` // YYYY-MM-DD
	ObservationNotes string    `json:"observation_notes" db:"observation_notes"`
	LevelDescriptor  string    `json:"level_descriptor" db:"level_descriptor"` // Iniziale, In evoluzione, Consolidato
	RecordedBy       string    `json:"recorded_by" db:"recorded_by"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`

	StudentFullName string `json:"student_full_name,omitempty"`
	TeacherFullName string `json:"teacher_full_name,omitempty"`
}

type InfanziaDiarioBordo struct {
	ID             string    `json:"id" db:"id"`
	SectionID      string    `json:"section_id" db:"section_id"`
	StoryDate      string    `json:"story_date" db:"story_date"` // YYYY-MM-DD
	Title          string    `json:"title" db:"title"`
	NarrativeStory string    `json:"narrative_story" db:"narrative_story"`
	MediaURLsJSON  string    `json:"media_urls" db:"media_urls"` // photo & video showcase for parents
	CreatedBy      string    `json:"created_by" db:"created_by"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`

	SectionName     string `json:"section_name,omitempty"`
	TeacherFullName string `json:"teacher_full_name,omitempty"`
}
