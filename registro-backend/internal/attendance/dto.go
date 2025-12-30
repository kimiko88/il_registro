package attendance

import "time"

type CreateAttendanceRequest struct {
	StudentID uint      `json:"student_id" binding:"required"`
	ClassID   uint      `json:"class_id" binding:"required"`
	Date      time.Time `json:"date" binding:"required"`
	Status    string    `json:"status" binding:"required"`
}

type AttendanceResponse struct {
	ID        uint      `json:"id"`
	StudentID uint      `json:"student_id"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"`
	Justified bool      `json:"justified"`
}
