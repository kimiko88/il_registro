package grades

import (
	"time"
)

// Helper Functions for Italian Grading System

// 1. ConvertJudgmentToNumeric converts judgment string to 0-10 scale
func ConvertJudgmentToNumeric(judgment string) float64 {
	switch judgment {
	case "Insufficiente":
		return 3.0 // Or 4.0 depending on school policy, using 3.0 as per prompt
	case "Mediocre":
		return 4.5
	case "Sufficiente":
		return 6.0
	case "Discreto":
		return 7.0
	case "Buono":
		return 8.0
	case "Distinto":
		return 9.0
	case "Ottimo":
		return 10.0
	default:
		return 0.0
	}
}

// 2. ConvertNumericToJudgment converts value 0-10 to judgment string
func ConvertNumericToJudgment(value float64) string {
	if value < 4.0 {
		return "Insufficiente"
	} else if value < 6.0 {
		return "Mediocre"
	} else if value < 7.0 {
		return "Sufficiente"
	} else if value < 8.0 {
		return "Discreto"
	} else if value < 9.0 {
		return "Buono"
	} else if value < 10.0 {
		return "Distinto"
	}
	return "Ottimo"
}

// 3. IsGradePublishable checks if a grade is valid for publication
func IsGradePublishable(grade Grade) bool {
	if !grade.IsPublished {
		return false
	}
	if grade.DeletedAt != nil {
		return false
	}
	return true
}

// 4. GetSemesterDateRange returns start and end dates for a semester dynamically
func GetSemesterDateRange(semester int) (time.Time, time.Time) {
	now := time.Now()
	yearStart := now.Year()
	if now.Month() < time.September {
		yearStart--
	}

	switch semester {
	case 1:
		// Sept 1 to Jan 31
		start := time.Date(yearStart, time.September, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(yearStart+1, time.January, 31, 23, 59, 59, 0, time.UTC)
		return start, end
	case 2:
		// Feb 1 to June 30
		start := time.Date(yearStart+1, time.February, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(yearStart+1, time.June, 30, 23, 59, 59, 0, time.UTC)
		return start, end
	default:
		return time.Time{}, time.Time{}
	}
}

// 5. IsInLockPeriod checks if the current date falls within a lock period for the semester
func IsInLockPeriod(semester int, date time.Time) bool {
	year := date.Year()

	var start, end time.Time

	switch semester {
	case 1:
		start = time.Date(year, time.January, 28, 0, 0, 0, 0, time.UTC)
		end = time.Date(year, time.February, 15, 23, 59, 59, 0, time.UTC)
	case 2:
		start = time.Date(year, time.June, 25, 0, 0, 0, 0, time.UTC)
		end = time.Date(year, time.July, 15, 23, 59, 59, 0, time.UTC)
	default:
		return false
	}

	return date.After(start) && date.Before(end)
}
