package rubrics

import (
	"time"
)

type Rubric struct {
	ID          string      `json:"id"`
	SchoolID    string      `json:"school_id"`
	TeacherID   string      `json:"teacher_id"`
	SubjectID   string      `json:"subject_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Criteria    []Criterion `json:"criteria,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
}

type Criterion struct {
	ID          string  `json:"id"`
	RubricID    string  `json:"rubric_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	MaxScore    float64 `json:"max_score"`
	Levels      []Level `json:"levels,omitempty"`
}

type Level struct {
	ID          string  `json:"id"`
	CriterionID string  `json:"criterion_id"`
	Score       float64 `json:"score"`
	Label       string  `json:"label"`
	Description string  `json:"description"`
}

type RubricAssessment struct {
	ID         string           `json:"id"`
	RubricID   string           `json:"rubric_id"`
	StudentID  string           `json:"student_id"`
	ClassID    string           `json:"class_id"`
	TeacherID  string           `json:"teacher_id"`
	Date       time.Time        `json:"date"`
	Scores     []CriterionScore `json:"scores,omitempty"`
	TotalScore float64          `json:"total_score"`
	Notes      string           `json:"notes"`
	CreatedAt  time.Time        `json:"created_at"`

	// Joined metadata
	StudentName string `json:"student_name,omitempty"`
	RubricTitle string `json:"rubric_title,omitempty"`
}

type CriterionScore struct {
	CriterionID string  `json:"criterion_id"`
	LevelID     string  `json:"level_id"`
	Score       float64 `json:"score"`
}

type CreateRubricRequest struct {
	SubjectID   string                 `json:"subject_id" binding:"required"`
	Title       string                 `json:"title" binding:"required"`
	Description string                 `json:"description"`
	Criteria    []CreateCriterionInput `json:"criteria"`
}

type CreateCriterionInput struct {
	Name        string             `json:"name" binding:"required"`
	Description string             `json:"description"`
	MaxScore    float64            `json:"max_score"`
	Levels      []CreateLevelInput `json:"levels"`
}

type CreateLevelInput struct {
	Score       float64 `json:"score"`
	Label       string  `json:"label" binding:"required"`
	Description string  `json:"description"`
}

type CreateAssessmentRequest struct {
	StudentID string           `json:"student_id" binding:"required"`
	ClassID   string           `json:"class_id" binding:"required"`
	Date      string           `json:"date"` // YYYY-MM-DD
	Scores    []CriterionScore `json:"scores" binding:"required"`
	Notes     string           `json:"notes"`
}
