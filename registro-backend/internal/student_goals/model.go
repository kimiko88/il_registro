package student_goals

import "time"

type GoalStatus string

const (
	StatusPending    GoalStatus = "pending"
	StatusInProgress GoalStatus = "in_progress"
	StatusCompleted  GoalStatus = "completed"
)

type StudentGoal struct {
	ID          string     `json:"id" db:"id"`
	StudentID   string     `json:"student_id" db:"student_id"`
	TeacherID   string     `json:"teacher_id" db:"teacher_id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	BadgeName   string     `json:"badge_name,omitempty" db:"badge_name"`
	BadgeIcon   string     `json:"badge_icon,omitempty" db:"badge_icon"`
	Category    string     `json:"category" db:"category"` // 'academic', 'behavioral', 'social'
	Status      GoalStatus `json:"status" db:"status"`
	Points      int        `json:"points" db:"points"`
	DueDate     *time.Time `json:"due_date,omitempty" db:"due_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}

type CreateGoalRequest struct {
	StudentID   string `json:"student_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	BadgeName   string `json:"badge_name"`
	BadgeIcon   string `json:"badge_icon"`
	Category    string `json:"category"`
	Points      int    `json:"points"`
	DueDate     string `json:"due_date"`
}

type UpdateGoalStatusRequest struct {
	Status GoalStatus `json:"status" binding:"required"`
}
