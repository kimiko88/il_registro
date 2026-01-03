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
			name:       "secretary does not have user delete",
			role:       "secretary",
			permission: UserDelete,
			expected:   false,
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
			name:        "secretary missing one permission",
			role:        "secretary",
			permissions: []Permission{UserCreate, UserRead, UserDelete},
			expected:    false,
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
	t.Run("superadmin and admin have same permissions", func(t *testing.T) {
		superadminPerms := RoleDefinitions["superadmin"]
		adminPerms := RoleDefinitions["admin"]
		assert.Equal(t, len(superadminPerms), len(adminPerms))
	})

	// Verify student and parent have minimal permissions
	t.Run("student has minimal permissions", func(t *testing.T) {
		studentPerms := RoleDefinitions["student"]
		assert.Equal(t, 0, len(studentPerms))
	})

	t.Run("parent has minimal permissions", func(t *testing.T) {
		parentPerms := RoleDefinitions["parent"]
		assert.Equal(t, 0, len(parentPerms))
	})
}
