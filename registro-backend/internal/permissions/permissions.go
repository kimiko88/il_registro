package permissions

import (
	"fmt"
	"strings"
)

// Permission represents a granular permission
type Permission string

// Resource represents a system resource
type Resource string

// Action represents an action on a resource
type Action string

const (
	// Resources
	ResourceUsers       Resource = "users"
	ResourceClasses     Resource = "classes"
	ResourceGrades      Resource = "grades"
	ResourceAssignments Resource = "assignments"
	ResourceAttendance  Resource = "attendance"
	ResourceAudit       Resource = "audit"
)

const (
	// Actions
	ActionCreate Action = "create"
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
	ActionExport Action = "export"
	ActionImport Action = "import"
)

// Permission format: "resource:action"
func New(resource Resource, action Action) Permission {
	return Permission(fmt.Sprintf("%s:%s", resource, action))
}

const (
	// User Permissions
	UserCreate Permission = "users:create"
	UserRead   Permission = "users:read"
	UserUpdate Permission = "users:update"
	UserDelete Permission = "users:delete"
	UserImport Permission = "users:import"
	UserExport Permission = "users:export"
	UserAudit  Permission = "users:audit"

	// Audit Permissions
	AuditRead Permission = "audit:read"
)

// RoleDefinitions maps roles to their default permissions
var RoleDefinitions = map[string][]Permission{
	"superadmin": {
		UserCreate, UserRead, UserUpdate, UserDelete, UserImport, UserExport, UserAudit,
		AuditRead,
	},
	"admin": {
		UserCreate, UserRead, UserUpdate, UserDelete, UserImport, UserExport, UserAudit,
		AuditRead,
	},
	"principal": {
		UserRead, UserExport, UserAudit,
		AuditRead,
	},
	"secretary": {
		UserRead, UserCreate, UserUpdate,
		UserExport,
	},
	"teacher": {
		UserRead,
	},
	"student": {
		// Minimal self-access is usually handled by logic, not generic permissions
	},
	"parent": {
		// Minimal self-access
	},
}

// Manager handles permission checks
type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

// HasPermission checks if a role has the required permission
func (m *Manager) HasPermission(role string, required Permission) bool {
	perms, ok := RoleDefinitions[strings.ToLower(role)]
	if !ok {
		return false
	}

	for _, p := range perms {
		if p == required {
			return true
		}
	}
	return false
}

// CheckPermissions checks if a role has ALL required permissions
func (m *Manager) CheckPermissions(role string, required ...Permission) bool {
	for _, req := range required {
		if !m.HasPermission(role, req) {
			return false
		}
	}
	return true
}
