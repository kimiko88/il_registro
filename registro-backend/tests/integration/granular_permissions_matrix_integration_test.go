package integration

import (
	"testing"

	"registro-backend/internal/permissions"

	"github.com/stretchr/testify/assert"
)

// ─── Tests ─────────────────────────────────────────────────────────────────

// PERM01 — SuperAdmin has full permissions across all system resources
func TestPermissions_SuperAdmin_FullAccess(t *testing.T) {
	mgr := permissions.NewManager()

	assert.True(t, mgr.HasPermission("superadmin", permissions.UserCreate))
	assert.True(t, mgr.HasPermission("superadmin", permissions.UserDelete))
	assert.True(t, mgr.HasPermission("superadmin", permissions.AuditRead))
	assert.True(t, mgr.HasPermission("superadmin", permissions.GradeCreate))
	assert.True(t, mgr.HasPermission("superadmin", permissions.AttendanceCreate))
	assert.True(t, mgr.HasPermission("superadmin", permissions.DocumentDelete))
	assert.True(t, mgr.HasPermission("superadmin", permissions.PCTODelete))
	assert.True(t, mgr.CheckPermissions("superadmin", permissions.UserRead, permissions.GradeRead, permissions.AuditRead))
}

// PERM02 — Teacher can manage grades, attendance, and create slots, but cannot manage users or audit
func TestPermissions_Teacher_AcademicScoping(t *testing.T) {
	mgr := permissions.NewManager()

	// Allowed
	assert.True(t, mgr.HasPermission("teacher", permissions.GradeCreate))
	assert.True(t, mgr.HasPermission("teacher", permissions.GradeUpdate))
	assert.True(t, mgr.HasPermission("teacher", permissions.AttendanceCreate))
	assert.True(t, mgr.HasPermission("teacher", permissions.SchedulingCreate))

	// Denied
	assert.False(t, mgr.HasPermission("teacher", permissions.UserCreate))
	assert.False(t, mgr.HasPermission("teacher", permissions.UserDelete))
	assert.False(t, mgr.HasPermission("teacher", permissions.AuditRead))
	assert.False(t, mgr.HasPermission("teacher", permissions.DocumentDelete))
}

// PERM03 — Parent can read grades and book slots, but cannot create slots or edit grades
func TestPermissions_Parent_GuardianScoping(t *testing.T) {
	mgr := permissions.NewManager()

	// Allowed
	assert.True(t, mgr.HasPermission("parent", permissions.GradeRead))
	assert.True(t, mgr.HasPermission("parent", permissions.AttendanceRead))
	assert.True(t, mgr.HasPermission("parent", permissions.SchedulingBook))

	// Denied
	assert.False(t, mgr.HasPermission("parent", permissions.SchedulingCreate))
	assert.False(t, mgr.HasPermission("parent", permissions.GradeCreate))
	assert.False(t, mgr.HasPermission("parent", permissions.AttendanceCreate))
	assert.False(t, mgr.HasPermission("parent", permissions.DocumentCreate))
}

// PERM04 — Secretary can manage documents, users, and PCTO, but cannot assign grades
func TestPermissions_Secretary_AdministrativeScoping(t *testing.T) {
	mgr := permissions.NewManager()

	// Allowed
	assert.True(t, mgr.HasPermission("secretary", permissions.UserCreate))
	assert.True(t, mgr.HasPermission("secretary", permissions.DocumentCreate))
	assert.True(t, mgr.HasPermission("secretary", permissions.DocumentDelete))
	assert.True(t, mgr.HasPermission("secretary", permissions.PCTOCreate))

	// Denied
	assert.False(t, mgr.HasPermission("secretary", permissions.GradeCreate))
	assert.False(t, mgr.HasPermission("secretary", permissions.GradeUpdate))
}

// PERM05 — Case-insensitivity and multi-permission checking
func TestPermissions_Manager_CaseInsensitiveAndMultiCheck(t *testing.T) {
	mgr := permissions.NewManager()

	// Case-insensitivity
	assert.True(t, mgr.HasPermission("TEACHER", permissions.GradeCreate))
	assert.True(t, mgr.HasPermission("SuperAdmin", permissions.AuditRead))
	assert.True(t, mgr.HasPermission("Parent", permissions.SchedulingBook))

	// Multi-permission checks
	assert.True(t, mgr.CheckPermissions("teacher", permissions.GradeRead, permissions.AttendanceRead))
	assert.False(t, mgr.CheckPermissions("teacher", permissions.GradeRead, permissions.AuditRead))
}
