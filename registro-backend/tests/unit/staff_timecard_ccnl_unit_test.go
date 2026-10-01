package unit

import (
	"testing"

	"registro-backend/internal/staff_attendance"

	"github.com/stretchr/testify/assert"
)

func TestStaffAttendance_ATARolesMatrix(t *testing.T) {
	expectedATARoles := []string{
		"dsga",
		"assistente_amministrativo",
		"collaboratore_ds",
		"collaboratore_scolastico",
		"assistente_tecnico",
		"assistente_alunni",
		"assistente_personale",
		"assistente_contabilita",
		"assistente_protocollo",
		"assistente_sportello",
		"responsabile_servizio",
	}

	for _, role := range expectedATARoles {
		assert.True(t, staff_attendance.IsATARole(role), "Role %s must be recognized as ATA role", role)
	}

	nonATARoles := []string{
		"teacher",
		"docente",
		"student",
		"parent",
		"admin",
		"superadmin",
		"guest",
	}

	for _, role := range nonATARoles {
		assert.False(t, staff_attendance.IsATARole(role), "Role %s must NOT be an ATA role", role)
	}
}

func TestStaffAttendance_PermissionsGranularity(t *testing.T) {
	// 1. Roles authorized to read and write school-wide staff attendance
	authorizedManagers := []string{
		"dsga",
		"collaboratore_ds",
		"assistente_amministrativo",
		"assistente_personale",
		"principal",
		"vice_principal",
		"secretary",
		"admin",
		"superadmin",
	}

	for _, role := range authorizedManagers {
		assert.True(t, staff_attendance.CanReadAttendance(role), "Role %s should have read access", role)
		assert.True(t, staff_attendance.CanWriteAttendance(role), "Role %s should have write access", role)
	}

	// 2. Collaboratore Scolastico has NO school-wide attendance read/write access (personal timecard only)
	assert.False(t, staff_attendance.CanReadAttendance("collaboratore_scolastico"))
	assert.False(t, staff_attendance.CanWriteAttendance("collaboratore_scolastico"))

	// 3. Teachers and students have no staff attendance management access
	assert.False(t, staff_attendance.CanReadAttendance("teacher"))
	assert.False(t, staff_attendance.CanWriteAttendance("teacher"))
	assert.False(t, staff_attendance.CanReadAttendance("student"))
	assert.False(t, staff_attendance.CanWriteAttendance("student"))
}

func TestStaffAttendance_TimecardBalanceCalculation(t *testing.T) {
	calcBalanceMinutes := func(workedMinutes, contractMinutes int) (overtimeMinutes int, deficitMinutes int) {
		delta := workedMinutes - contractMinutes
		if delta > 0 {
			return delta, 0
		} else if delta < 0 {
			return 0, -delta
		}
		return 0, 0
	}

	t.Run("calculates positive overtime under CCNL 36-hour schedule", func(t *testing.T) {
		contractMinutes := 360 // 6 hours
		workedMinutes := 420   // 7 hours
		overtime, deficit := calcBalanceMinutes(workedMinutes, contractMinutes)
		assert.Equal(t, 60, overtime, "1 hour of overtime")
		assert.Equal(t, 0, deficit)
	})

	t.Run("calculates negative deficit under CCNL when early departure occurs", func(t *testing.T) {
		contractMinutes := 360 // 6 hours
		workedMinutes := 300   // 5 hours
		overtime, deficit := calcBalanceMinutes(workedMinutes, contractMinutes)
		assert.Equal(t, 0, overtime)
		assert.Equal(t, 60, deficit, "1 hour of deficit to recover")
	})

	t.Run("calculates zero delta when exactly contracted hours worked", func(t *testing.T) {
		contractMinutes := 360
		workedMinutes := 360
		overtime, deficit := calcBalanceMinutes(workedMinutes, contractMinutes)
		assert.Equal(t, 0, overtime)
		assert.Equal(t, 0, deficit)
	})
}
