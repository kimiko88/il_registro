package permissions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		resource Resource
		action   Action
		expected Permission
	}{
		{
			name:     "create user permission",
			resource: ResourceUsers,
			action:   ActionCreate,
			expected: "users:create",
		},
		{
			name:     "read grades permission",
			resource: ResourceGrades,
			action:   ActionRead,
			expected: "grades:read",
		},
		{
			name:     "export attendance permission",
			resource: ResourceAttendance,
			action:   ActionExport,
			expected: "attendance:export",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := New(tt.resource, tt.action)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestManager_HasPermission(t *testing.T) {
	manager := NewManager()

	tests := []struct {
		name       string
		role       string
		permission Permission
		expected   bool
	}{
		// Superadmin tests
		{
			name:       "superadmin has user create",
			role:       "superadmin",
			permission: UserCreate,
			expected:   true,
		},
		{
			name:       "superadmin has user delete",
			role:       "superadmin",
			permission: UserDelete,
			expected:   true,
		},
		{
			name:       "superadmin has audit read",
			role:       "superadmin",
			permission: AuditRead,
			expected:   true,
		},
		// Admin tests
		{
			name:       "admin has user create",
			role:       "admin",
			permission: UserCreate,
			expected:   true,
		},
		{
			name:       "admin has user update",
			role:       "admin",
			permission: UserUpdate,
			expected:   true,
		},
		// Principal tests
		{
			name:       "principal has user read",
			role:       "principal",
			permission: UserRead,
			expected:   true,
		},
		{
			name:       "principal does not have user create",
			role:       "principal",
			permission: UserCreate,
			expected:   false,
		},
		{
			name:       "principal does not have user delete",
			role:       "principal",
			permission: UserDelete,
			expected:   false,
		},
		// Secretary tests
		{
			name:       "secretary has user read",
			role:       "secretary",
			permission: UserRead,
			expected:   true,
		},
		{
			name:       "secretary has user create",
			role:       "secretary",
			permission: UserCreate,
			expected:   true,
		},
		{
			name:       "secretary has user delete",
			role:       "secretary",
			permission: UserDelete,
			expected:   true,
		},
		{
			name:       "secretary does not have audit read",
			role:       "secretary",
			permission: AuditRead,
			expected:   false,
		},
		// Teacher tests
		{
			name:       "teacher has user read",
			role:       "teacher",
			permission: UserRead,
			expected:   true,
		},
		{
			name:       "teacher has grade create",
			role:       "teacher",
			permission: GradeCreate,
			expected:   true,
		},
		{
			name:       "teacher has attendance update",
			role:       "teacher",
			permission: AttendanceUpdate,
			expected:   true,
		},
		{
			name:       "teacher does not have user create",
			role:       "teacher",
			permission: UserCreate,
			expected:   false,
		},
		{
			name:       "teacher does not have user delete",
			role:       "teacher",
			permission: UserDelete,
			expected:   false,
		},
		// Student tests
		{
			name:       "student does not have user read",
			role:       "student",
			permission: UserRead,
			expected:   false,
		},
		{
			name:       "student has grade read",
			role:       "student",
			permission: GradeRead,
			expected:   true,
		},
		{
			name:       "student does not have user create",
			role:       "student",
			permission: UserCreate,
			expected:   false,
		},
		// Parent tests
		{
			name:       "parent does not have user read",
			role:       "parent",
			permission: UserRead,
			expected:   false,
		},
		{
			name:       "parent has scheduling book",
			role:       "parent",
			permission: SchedulingBook,
			expected:   true,
		},
		// Edge cases
		{
			name:       "invalid role returns false",
			role:       "invalid_role",
			permission: UserRead,
			expected:   false,
		},
		{
			name:       "empty role returns false",
			role:       "",
			permission: UserRead,
			expected:   false,
		},
		{
			name:       "case insensitive - ADMIN",
			role:       "ADMIN",
			permission: UserCreate,
			expected:   true,
		},
		{
			name:       "case insensitive - Admin",
			role:       "Admin",
			permission: UserCreate,
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.HasPermission(tt.role, tt.permission)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestManager_CheckPermissions(t *testing.T) {
	manager := NewManager()

	tests := []struct {
		name        string
		role        string
		permissions []Permission
		expected    bool
	}{
		{
			name:        "superadmin has all requested permissions",
			role:        "superadmin",
			permissions: []Permission{UserCreate, UserRead, UserUpdate, UserDelete},
			expected:    true,
		},
		{
			name:        "admin has all requested permissions",
			role:        "admin",
			permissions: []Permission{UserCreate, UserRead, UserUpdate},
			expected:    true,
		},
		{
			name:        "secretary has all requested permissions",
			role:        "secretary",
			permissions: []Permission{UserCreate, UserRead, UserDelete},
			expected:    true,
		},
		{
			name:        "teacher has only user read",
			role:        "teacher",
			permissions: []Permission{UserRead},
			expected:    true,
		},
		{
			name:        "teacher missing user create",
			role:        "teacher",
			permissions: []Permission{UserRead, UserCreate},
			expected:    false,
		},
		{
			name:        "student has no permissions",
			role:        "student",
			permissions: []Permission{UserRead},
			expected:    false,
		},
		{
			name:        "empty permissions list returns true",
			role:        "student",
			permissions: []Permission{},
			expected:    true,
		},
		{
			name:        "invalid role with permissions",
			role:        "invalid",
			permissions: []Permission{UserRead},
			expected:    false,
		},
		{
			name:        "principal has export but not delete",
			role:        "principal",
			permissions: []Permission{UserExport},
			expected:    true,
		},
		{
			name:        "principal missing delete permission",
			role:        "principal",
			permissions: []Permission{UserExport, UserDelete},
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.CheckPermissions(tt.role, tt.permissions...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewManager(t *testing.T) {
	manager := NewManager()
	assert.NotNil(t, manager)
	assert.IsType(t, &Manager{}, manager)
}

func TestRoleDefinitions(t *testing.T) {
	// Verify all expected roles exist
	expectedRoles := []string{"superadmin", "admin", "principal", "secretary", "teacher", "student", "parent"}

	for _, role := range expectedRoles {
		t.Run("role_"+role+"_exists", func(t *testing.T) {
			_, exists := RoleDefinitions[role]
			assert.True(t, exists, "Role %s should exist in RoleDefinitions", role)
		})
	}

	// Verify superadmin and admin have the same permissions
	// Verify superadmin has more or equal permissions than admin
	t.Run("superadmin has comprehensive permissions", func(t *testing.T) {
		superadminPerms := RoleDefinitions["superadmin"]
		adminPerms := RoleDefinitions["admin"]
		assert.GreaterOrEqual(t, len(superadminPerms), len(adminPerms))
	})

	// Verify student and parent have some read permissions
	t.Run("student has read permissions", func(t *testing.T) {
		studentPerms := RoleDefinitions["student"]
		assert.NotEmpty(t, studentPerms)
		assert.Contains(t, studentPerms, GradeRead)
	})

	t.Run("parent has read permissions", func(t *testing.T) {
		parentPerms := RoleDefinitions["parent"]
		assert.NotEmpty(t, parentPerms)
		assert.Contains(t, parentPerms, GradeRead)
	})

	// Verify all 25 official roles & duty profiles exist
	all25Roles := []string{
		"superadmin", "principal", "vice_principal", "dsga", "teacher",
		"coordinatore_classe", "segretario_consiglio", "referente_progetto", "referente_inclusione",
		"responsabile_dipartimento", "tutor_orientatore", "animatore_digitale",
		"assistente_alunni", "assistente_personale", "assistente_contabilita", "assistente_protocollo", "assistente_sportello",
		"collaboratore_scolastico", "assistente_tecnico", "responsabile_servizio",
		"responsabile_gestione_documentale", "responsabile_conservazione", "dpo",
		"student", "parent",
	}

	for _, r := range all25Roles {
		t.Run("verify_25_role_"+r, func(t *testing.T) {
			_, exists := RoleDefinitions[r]
			assert.True(t, exists, "Role or duty %s must exist in RoleDefinitions", r)
		})
	}
}

func TestManager_HasUserPermission(t *testing.T) {
	mgr := NewManager()

	t.Run("Teacher without assignments cannot sign verbali or edit scrutiny", func(t *testing.T) {
		assert.False(t, mgr.HasUserPermission("teacher", nil, VerbaliSign))
		assert.False(t, mgr.HasUserPermission("teacher", []string{}, ScrutinyUpdate))
	})

	t.Run("Teacher with coordinatore_classe assignment gains scrutiny update", func(t *testing.T) {
		assert.True(t, mgr.HasUserPermission("teacher", []string{"coordinatore_classe"}, ScrutinyUpdate))
		assert.True(t, mgr.HasUserPermission("teacher", []string{"coordinatore_classe"}, GradeRead))
	})

	t.Run("Teacher with segretario_consiglio assignment gains verbali sign", func(t *testing.T) {
		assert.True(t, mgr.HasUserPermission("teacher", []string{"segretario_consiglio"}, VerbaliSign))
	})

	t.Run("Teacher with multiple assignments gains combined permissions", func(t *testing.T) {
		assignments := []string{"coordinatore_classe", "referente_inclusione", "animatore_digitale"}
		assert.True(t, mgr.HasUserPermission("teacher", assignments, ScrutinyUpdate))
		assert.True(t, mgr.HasUserPermission("teacher", assignments, InclusionUpdate))
		assert.True(t, mgr.HasUserPermission("teacher", assignments, TechnicalConfig))
		assert.False(t, mgr.HasUserPermission("teacher", assignments, UserDelete))
	})

	t.Run("Principal has broad institutional permissions", func(t *testing.T) {
		assert.True(t, mgr.HasPermission("principal", ScrutinyValidate))
		assert.True(t, mgr.HasPermission("principal", VerbaliValidate))
		assert.True(t, mgr.HasPermission("principal", AuditRead))
	})
}
