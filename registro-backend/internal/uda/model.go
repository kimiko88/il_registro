package uda

import (
	"encoding/json"
	"time"
)

type UdaPlan struct {
	ID                 string          `json:"id" db:"id"`
	SchoolID           string          `json:"school_id" db:"school_id"`
	ClassID            string          `json:"class_id" db:"class_id"`
	ClassName          string          `json:"class_name,omitempty" db:"class_name"`
	SubjectID          string          `json:"subject_id" db:"subject_id"`
	SubjectName        string          `json:"subject_name,omitempty" db:"subject_name"`
	TeacherID          string          `json:"teacher_id" db:"teacher_id"`
	TeacherName        string          `json:"teacher_name,omitempty" db:"teacher_name"`
	Title              string          `json:"title" db:"title"`
	Description        string          `json:"description" db:"description"`
	Period             string          `json:"period" db:"period"` // primo_quadrimestre, secondo_quadrimestre, annuale
	StartDate          *time.Time      `json:"start_date,omitempty" db:"start_date"`
	EndDate            *time.Time      `json:"end_date,omitempty" db:"end_date"`
	CompetenciesRaw    json.RawMessage `json:"-" db:"competencies"`
	Competencies       []string        `json:"competencies"`
	Objectives         string          `json:"objectives" db:"objectives"`
	Methodologies      string          `json:"methodologies" db:"methodologies"`
	EvaluationCriteria string          `json:"evaluation_criteria" db:"evaluation_criteria"`
	Status             string          `json:"status" db:"status"` // draft, submitted, approved
	CreatedAt          time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at" db:"updated_at"`
}

type CreateUdaRequest struct {
	ClassID            string   `json:"class_id" binding:"required"`
	SubjectID          string   `json:"subject_id" binding:"required"`
	Title              string   `json:"title" binding:"required"`
	Description        string   `json:"description"`
	Period             string   `json:"period"`
	StartDate          string   `json:"start_date"`
	EndDate            string   `json:"end_date"`
	Competencies       []string `json:"competencies"`
	Objectives         string   `json:"objectives"`
	Methodologies      string   `json:"methodologies"`
	EvaluationCriteria string   `json:"evaluation_criteria"`
}

type UpdateUdaRequest struct {
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Period             string   `json:"period"`
	StartDate          string   `json:"start_date"`
	EndDate            string   `json:"end_date"`
	Competencies       []string `json:"competencies"`
	Objectives         string   `json:"objectives"`
	Methodologies      string   `json:"methodologies"`
	EvaluationCriteria string   `json:"evaluation_criteria"`
	Status             string   `json:"status"`
}
