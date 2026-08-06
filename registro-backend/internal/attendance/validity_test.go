package attendance

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAttendanceValidity_Within75PercentThreshold(t *testing.T) {
	totalSchoolDays := 200
	absentDays := 40 // 20% absence, 80% presence => VALID

	absencePercentage := (float64(absentDays) / float64(totalSchoolDays)) * 100
	isValidYear := absencePercentage <= 25.0

	assert.InDelta(t, 20.0, absencePercentage, 0.001)
	assert.True(t, isValidYear, "Student with 20% absence should have a valid school year")
}

func TestAttendanceValidity_Exceeds25PercentThreshold(t *testing.T) {
	totalSchoolDays := 200
	absentDays := 55 // 27.5% absence => INVALID

	absencePercentage := (float64(absentDays) / float64(totalSchoolDays)) * 100
	isValidYear := absencePercentage <= 25.0

	assert.InDelta(t, 27.5, absencePercentage, 0.001)
	assert.False(t, isValidYear, "Student with >25% absence should NOT have a valid school year")
}

func TestAttendanceValidity_Exact25PercentBoundary(t *testing.T) {
	totalSchoolDays := 200
	absentDays := 50 // Exactly 25% absence => VALID

	absencePercentage := (float64(absentDays) / float64(totalSchoolDays)) * 100
	isValidYear := absencePercentage <= 25.0

	assert.InDelta(t, 25.0, absencePercentage, 0.001)
	assert.True(t, isValidYear, "Student with exactly 25% absence is at the exact threshold and valid")
}
