package documents

import (
	"registro-backend/internal/models"
	"time"
)

type Document struct {
	models.Model
	Title       string    `json:"title"`
	Description string    `json:"description"`
	FilePath    string    `json:"file_path"`
	UploadedBy  uint      `json:"uploaded_by"`
	Type        string    `json:"type"`
	Date        time.Time `json:"date"`
}
