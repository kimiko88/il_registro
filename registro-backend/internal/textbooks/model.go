package textbooks

import "time"

type Textbook struct {
	ID        string    `json:"id" db:"id"`
	SchoolID  string    `json:"school_id" db:"school_id"`
	Title     string    `json:"title" db:"title"`
	Author    string    `json:"author" db:"author"`
	Subject   string    `json:"subject" db:"subject"`
	ISBN      string    `json:"isbn" db:"isbn"`
	Publisher string    `json:"publisher" db:"publisher"`
	Price     float64   `json:"price" db:"price"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ClassTextbook struct {
	ID           string `json:"id" db:"id"`
	ClassID      string `json:"class_id" db:"class_id"`
	SubjectID    string `json:"subject_id" db:"subject_id"`
	TextbookID   string `json:"textbook_id" db:"textbook_id"`
	IsOptional   bool   `json:"is_optional" db:"is_optional"`
	
	// Joined fields
	Title       string  `json:"title,omitempty"`
	Author      string  `json:"author,omitempty"`
	SubjectName string  `json:"subject_name,omitempty"`
}

type CreateTextbookRequest struct {
	Title     string  `json:"title" binding:"required"`
	Author    string  `json:"author"`
	Subject   string  `json:"subject"`
	ISBN      string  `json:"isbn"`
	Publisher string  `json:"publisher"`
	Price     float64 `json:"price"`
}

type AssignTextbookRequest struct {
	SubjectID  string `json:"subject_id" binding:"required"`
	TextbookID string `json:"textbook_id" binding:"required"`
	IsOptional bool   `json:"is_optional"`
}
