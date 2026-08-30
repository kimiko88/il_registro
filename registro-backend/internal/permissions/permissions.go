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
	ResourceDocuments   Resource = "documents"
	ResourcePCTO        Resource = "pcto"
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

	// Grade Permissions
	GradeCreate Permission = "grades:create"
	GradeRead   Permission = "grades:read"
	GradeUpdate Permission = "grades:update"
	GradeDelete Permission = "grades:delete"

	// Attendance Permissions
	AttendanceCreate Permission = "attendance:create"
	AttendanceRead   Permission = "attendance:read"
	AttendanceUpdate Permission = "attendance:update"

	// Scheduling Permissions
	SchedulingRead   Permission = "scheduling:read"
	SchedulingCreate Permission = "scheduling:create"
	SchedulingBook   Permission = "scheduling:book"

	// Audit Permissions
	AuditRead Permission = "audit:read"

	// Document Permissions
	DocumentCreate Permission = "documents:create"
	DocumentRead   Permission = "documents:read"
	DocumentUpdate Permission = "documents:update"
	DocumentDelete Permission = "documents:delete"

	// PCTO Permissions
	PCTORead   Permission = "pcto:read"
	PCTOCreate Permission = "pcto:create"
	PCTOUpdate Permission = "pcto:update"
	PCTODelete Permission = "pcto:delete"
)

// RoleDefinitions maps roles to their default permissions
var RoleDefinitions = map[string][]Permission{
	"superadmin": {
		UserCreate, UserRead, UserUpdate, UserDelete, UserImport, UserExport, UserAudit,
		AuditRead,
		GradeRead, GradeCreate, GradeUpdate, GradeDelete,
		AttendanceRead, AttendanceCreate, AttendanceUpdate,
		SchedulingRead, SchedulingCreate, SchedulingBook,
		DocumentRead, DocumentCreate, DocumentUpdate, DocumentDelete,
		PCTORead, PCTOCreate, PCTOUpdate, PCTODelete,
	},
	"admin": {
		UserCreate, UserRead, UserUpdate, UserDelete, UserImport, UserExport, UserAudit,
		AuditRead,
		GradeRead, GradeCreate, GradeUpdate, GradeDelete,
		AttendanceRead, AttendanceUpdate,
		SchedulingRead, SchedulingCreate,
		DocumentRead, DocumentCreate, DocumentUpdate, DocumentDelete,
		PCTORead, PCTOCreate, PCTOUpdate, PCTODelete,
	},

	"principal": {
		UserRead, UserExport, UserAudit,
		AuditRead,
		GradeRead,
		AttendanceRead,
		SchedulingRead,
		DocumentRead,
	},
	"secretary": {
		UserRead, UserCreate, UserUpdate, UserDelete,
		UserExport,
		AttendanceRead, AttendanceUpdate,
		DocumentRead, DocumentCreate, DocumentUpdate, DocumentDelete,
		PCTORead, PCTOCreate, PCTOUpdate, PCTODelete,
	},
	"teacher": {
		UserRead,
		GradeCreate, GradeRead, GradeUpdate,
		AttendanceCreate, AttendanceRead, AttendanceUpdate,
		SchedulingRead, SchedulingCreate, // Teachers create slots
		DocumentRead,
	},
	"student": {
		// Minimal self-access is usually handled by logic, not generic permissions
		GradeRead, AttendanceRead, SchedulingRead,
	},
	"parent": {
		// Minimal self-access
		GradeRead, AttendanceRead, SchedulingBook, SchedulingRead,
	},
}

// Manager handles permission checks
type Manager struct {
	fastMap map[string]map[Permission]bool
}

func NewManager() *Manager {
	fm := make(map[string]map[Permission]bool)
	for role, perms := range RoleDefinitions {
		r := strings.ToLower(role)
		if fm[r] == nil {
			fm[r] = make(map[Permission]bool)
		}
		for _, p := range perms {
			fm[r][p] = true
		}
	}
	return &Manager{fastMap: fm}
}

// HasPermission checks if a role has the required permission
func (m *Manager) HasPermission(role string, required Permission) bool {
	r := strings.ToLower(role)
	if m != nil && m.fastMap != nil {
		if perms, ok := m.fastMap[r]; ok {
			return perms[required]
		}
		return false
	}
	perms, ok := RoleDefinitions[r]
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
