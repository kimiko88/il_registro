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
	ResourceUsers           Resource = "users"
	ResourceClasses         Resource = "classes"
	ResourceGrades          Resource = "grades"
	ResourceAssignments     Resource = "assignments"
	ResourceAttendance      Resource = "attendance"
	ResourceStaffAttendance Resource = "staff_attendance" // Presenze personale (docenti + ATA)
	ResourceAudit           Resource = "audit"
	ResourceDocuments       Resource = "documents"
	ResourcePCTO            Resource = "pcto"
	ResourceCommunications  Resource = "communications"
	ResourceScrutiny        Resource = "scrutiny"
	ResourceVerbali         Resource = "verbali"
	ResourceInclusion       Resource = "inclusion"
	ResourceProjects        Resource = "projects"
	ResourceTechnical       Resource = "technical"
	ResourceAccounting      Resource = "accounting"
	ResourceProtocol        Resource = "protocol"
	ResourceMaintenance     Resource = "maintenance"
	ResourcePreservation    Resource = "preservation"
	ResourcePrivacy         Resource = "privacy"
)

const (
	// Actions
	ActionCreate   Action = "create"
	ActionRead     Action = "read"
	ActionUpdate   Action = "update"
	ActionDelete   Action = "delete"
	ActionExport   Action = "export"
	ActionImport   Action = "import"
	ActionPublish  Action = "publish"
	ActionSign     Action = "sign"
	ActionValidate Action = "validate"
	ActionCancel   Action = "cancel"
	ActionAudit    Action = "audit"
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

	// Staff Attendance Permissions (Presenze Personale)
	StaffAttendanceRead   Permission = "staff_attendance:read"
	StaffAttendanceCreate Permission = "staff_attendance:create"
	StaffAttendanceUpdate Permission = "staff_attendance:update"
	StaffAttendanceDelete Permission = "staff_attendance:delete"

	// Communications
	CommunicationsRead    Permission = "communications:read"
	CommunicationsCreate  Permission = "communications:create"
	CommunicationsUpdate  Permission = "communications:update"
	CommunicationsDelete  Permission = "communications:delete"
	CommunicationsPublish Permission = "communications:publish"

	// Scrutiny
	ScrutinyRead     Permission = "scrutiny:read"
	ScrutinyUpdate   Permission = "scrutiny:update"
	ScrutinyValidate Permission = "scrutiny:validate"

	// Verbali
	VerbaliRead     Permission = "verbali:read"
	VerbaliCreate   Permission = "verbali:create"
	VerbaliUpdate   Permission = "verbali:update"
	VerbaliSign     Permission = "verbali:sign"
	VerbaliValidate Permission = "verbali:validate"

	// Inclusion (PEI / PDP)
	InclusionRead   Permission = "inclusion:read"
	InclusionUpdate Permission = "inclusion:update"

	// Projects
	ProjectsRead   Permission = "projects:read"
	ProjectsCreate Permission = "projects:create"
	ProjectsUpdate Permission = "projects:update"

	// Technical & System
	TechnicalConfig  Permission = "technical:config"
	TechnicalSupport Permission = "technical:support"

	// Accounting
	AccountingRead   Permission = "accounting:read"
	AccountingUpdate Permission = "accounting:update"

	// Protocol
	ProtocolRead   Permission = "protocol:read"
	ProtocolCreate Permission = "protocol:create"
	ProtocolUpdate Permission = "protocol:update"
	ProtocolCancel Permission = "protocol:cancel"

	// Maintenance
	MaintenanceRead   Permission = "maintenance:read"
	MaintenanceCreate Permission = "maintenance:create"
	MaintenanceUpdate Permission = "maintenance:update"

	// Preservation
	PreservationRead   Permission = "preservation:read"
	PreservationUpdate Permission = "preservation:update"

	// Privacy / DPO
	PrivacyRead  Permission = "privacy:read"
	PrivacyAudit Permission = "privacy:audit"
)

// RoleDefinitions maps all 25 roles and functional duties to their permissions
var RoleDefinitions = map[string][]Permission{
	// 1. Amministratore applicativo (ruolo tecnico di sistema)
	"superadmin": {
		UserCreate, UserRead, UserUpdate, UserDelete, UserImport, UserExport, UserAudit,
		AuditRead,
		GradeRead, GradeCreate, GradeUpdate, GradeDelete,
		AttendanceRead, AttendanceCreate, AttendanceUpdate,
		SchedulingRead, SchedulingCreate, SchedulingBook,
		DocumentRead, DocumentCreate, DocumentUpdate, DocumentDelete,
		PCTORead, PCTOCreate, PCTOUpdate, PCTODelete,
		TechnicalConfig, TechnicalSupport,
		CommunicationsRead, CommunicationsCreate, CommunicationsUpdate, CommunicationsDelete,
	},
	"admin": {
		UserCreate, UserRead, UserUpdate, UserDelete, UserImport, UserExport, UserAudit,
		AuditRead,
		GradeRead, GradeCreate, GradeUpdate, GradeDelete,
		AttendanceRead, AttendanceUpdate,
		SchedulingRead, SchedulingCreate,
		DocumentRead, DocumentCreate, DocumentUpdate, DocumentDelete,
		PCTORead, PCTOCreate, PCTOUpdate, PCTODelete,
		TechnicalConfig, TechnicalSupport,
		CommunicationsRead, CommunicationsCreate, CommunicationsUpdate, CommunicationsDelete,
	},

	// 2. Dirigente scolastico (massima visibilità istituzionale e validazioni)
	"principal": {
		UserRead, UserExport, UserAudit,
		AuditRead,
		GradeRead,
		AttendanceRead,
		StaffAttendanceRead, StaffAttendanceUpdate,
		SchedulingRead,
		DocumentRead, DocumentCreate, DocumentUpdate, DocumentDelete,
		CommunicationsRead, CommunicationsCreate, CommunicationsUpdate, CommunicationsPublish,
		ScrutinyRead, ScrutinyValidate,
		VerbaliRead, VerbaliValidate,
		InclusionRead,
		PCTORead,
		AccountingRead,
		ProtocolRead,
	},

	// 3. Collaboratore del dirigente / Vicario
	"vice_principal": {
		UserRead, UserAudit,
		GradeRead,
		AttendanceRead, AttendanceCreate, AttendanceUpdate,
		StaffAttendanceRead, StaffAttendanceCreate, StaffAttendanceUpdate,
		DocumentRead,
		SchedulingRead,
		CommunicationsRead, CommunicationsCreate, CommunicationsUpdate,
		VerbaliRead, VerbaliUpdate,
	},
	"collaboratore_ds": {
		UserRead, UserAudit,
		GradeRead,
		AttendanceRead, AttendanceCreate, AttendanceUpdate,
		StaffAttendanceRead, StaffAttendanceCreate, StaffAttendanceUpdate,
		DocumentRead,
		SchedulingRead,
		CommunicationsRead, CommunicationsCreate, CommunicationsUpdate,
		VerbaliRead, VerbaliUpdate,
	},

	// 4. DSGA (Direttore dei Servizi Generali e Amministrativi)
	"dsga": {
		UserRead, UserCreate, UserUpdate, UserExport,
		AttendanceRead,
		StaffAttendanceRead, StaffAttendanceCreate, StaffAttendanceUpdate, StaffAttendanceDelete,
		DocumentRead, DocumentCreate, DocumentUpdate, DocumentDelete,
		AuditRead,
		SchedulingRead,
		AccountingRead, AccountingUpdate,
		ProtocolRead, ProtocolCreate, ProtocolUpdate,
		CommunicationsRead, CommunicationsCreate,
	},

	// Segreteria generale
	"secretary": {
		UserRead, UserCreate, UserUpdate, UserDelete,
		UserExport,
		AttendanceRead, AttendanceUpdate,
		DocumentRead, DocumentCreate, DocumentUpdate, DocumentDelete,
		PCTORead, PCTOCreate, PCTOUpdate, PCTODelete,
		CommunicationsRead, CommunicationsCreate,
	},

	// 5. Docente (ruolo base)
	"teacher": {
		UserRead,
		GradeCreate, GradeRead, GradeUpdate,
		AttendanceCreate, AttendanceRead, AttendanceUpdate,
		SchedulingRead, SchedulingCreate,
		DocumentRead, DocumentCreate,
		CommunicationsRead,
		VerbaliRead,
	},

	// 6. Coordinatore di classe (incarico aggiuntivo docente)
	"coordinatore_classe": {
		GradeRead,
		AttendanceRead,
		ScrutinyRead, ScrutinyUpdate,
		CommunicationsCreate,
		VerbaliRead, VerbaliUpdate,
		InclusionRead,
	},
	"coordinator": {
		GradeRead,
		AttendanceRead,
		ScrutinyRead, ScrutinyUpdate,
		CommunicationsCreate,
		VerbaliRead, VerbaliUpdate,
		InclusionRead,
	},

	// 7. Segretario del consiglio di classe (incarico aggiuntivo docente)
	"segretario_consiglio": {
		VerbaliRead, VerbaliCreate, VerbaliUpdate, VerbaliSign,
		DocumentRead,
	},

	// 8. Referente di progetto / PCTO (incarico aggiuntivo docente)
	"referente_progetto": {
		PCTORead, PCTOCreate, PCTOUpdate,
		ProjectsRead, ProjectsCreate, ProjectsUpdate,
		DocumentRead, DocumentCreate,
		CommunicationsCreate,
	},

	// 9. Referente inclusione / BES / DSA (incarico aggiuntivo docente)
	"referente_inclusione": {
		InclusionRead, InclusionUpdate,
		DocumentRead, DocumentCreate, DocumentUpdate,
		CommunicationsCreate,
	},

	// 10. Responsabile di dipartimento (incarico aggiuntivo docente)
	"responsabile_dipartimento": {
		DocumentRead, DocumentCreate, DocumentUpdate,
		VerbaliRead, VerbaliCreate, VerbaliUpdate,
		CommunicationsCreate,
	},

	// 11. Tutor / docente orientatore (incarico aggiuntivo docente)
	"tutor_orientatore": {
		PCTORead, PCTOUpdate,
		ProjectsRead,
		DocumentRead,
		SchedulingRead, SchedulingCreate,
		CommunicationsCreate,
	},

	// 12. Animatore digitale / referente tecnologico (incarico aggiuntivo docente)
	"animatore_digitale": {
		TechnicalConfig, TechnicalSupport,
		DocumentRead, DocumentCreate,
		CommunicationsCreate,
	},

	// 13. Assistente amministrativo – area alunni
	"assistente_alunni": {
		UserRead, UserCreate, UserUpdate,
		AttendanceRead,
		DocumentRead, DocumentCreate, DocumentUpdate,
		CommunicationsRead, CommunicationsCreate,
	},

	// 14. Assistente amministrativo – area personale
	"assistente_personale": {
		UserRead, UserCreate, UserUpdate,
		StaffAttendanceRead, StaffAttendanceCreate, StaffAttendanceUpdate,
		DocumentRead, DocumentCreate, DocumentUpdate,
		CommunicationsRead, CommunicationsCreate,
	},

	// 15. Assistente amministrativo – contabilità
	"assistente_contabilita": {
		AccountingRead, AccountingUpdate,
		DocumentRead, DocumentCreate, DocumentUpdate,
		ProjectsRead,
	},

	// 16. Assistente amministrativo – protocollo e gestione documentale
	"assistente_protocollo": {
		ProtocolRead, ProtocolCreate, ProtocolUpdate, ProtocolCancel,
		DocumentRead, DocumentCreate, DocumentUpdate,
	},

	// 17. Assistente amministrativo – sportello/famiglie
	"assistente_sportello": {
		UserRead,
		AttendanceRead,
		DocumentRead, DocumentCreate,
		CommunicationsRead, CommunicationsCreate,
	},

	// Profilo generale Assistente Amministrativo (comprensivo delle aree ordinarie)
	"assistente_amministrativo": {
		UserRead, UserCreate, UserUpdate,
		AttendanceRead,
		StaffAttendanceRead, StaffAttendanceCreate, StaffAttendanceUpdate,
		DocumentRead, DocumentCreate, DocumentUpdate,
		SchedulingRead,
		ProtocolRead, ProtocolCreate, ProtocolUpdate,
		CommunicationsRead, CommunicationsCreate,
	},

	// 18. Collaboratore scolastico (sorveglianza plesso, presenze sciopero, guasti)
	"collaboratore_scolastico": {
		UserRead,
		StaffAttendanceRead, StaffAttendanceCreate, StaffAttendanceUpdate,
		DocumentRead,
		MaintenanceRead, MaintenanceCreate,
		CommunicationsRead,
	},

	// 19. Assistente tecnico (laboratori, attrezzature, inventario)
	"assistente_tecnico": {
		MaintenanceRead, MaintenanceCreate, MaintenanceUpdate,
		DocumentRead, DocumentCreate,
		CommunicationsRead,
	},

	// 20. Personale ATA responsabile di un servizio (incarico aggiuntivo ATA)
	"responsabile_servizio": {
		MaintenanceRead, MaintenanceCreate, MaintenanceUpdate,
		DocumentRead, DocumentCreate, DocumentUpdate,
		CommunicationsCreate,
	},

	// 21. Responsabile della gestione documentale
	"responsabile_gestione_documentale": {
		ProtocolRead, ProtocolCreate, ProtocolUpdate,
		PreservationRead,
		DocumentRead, DocumentCreate, DocumentUpdate,
	},

	// 22. Responsabile della conservazione
	"responsabile_conservazione": {
		PreservationRead, PreservationUpdate,
		DocumentRead,
		AuditRead,
	},

	// 23. DPO / Responsabile della protezione dei dati
	"dpo": {
		PrivacyRead, PrivacyAudit,
		AuditRead,
		DocumentRead,
	},
	"system_auditor": {
		AuditRead,
		DocumentRead,
		PrivacyRead,
	},

	// 24. Studente
	"student": {
		GradeRead, AttendanceRead, SchedulingRead,
		DocumentRead,
		CommunicationsRead,
	},

	// 25. Genitore / tutore
	"parent": {
		GradeRead, AttendanceRead, SchedulingBook, SchedulingRead,
		DocumentRead,
		CommunicationsRead,
	},

	// Alias generico staff
	"staff": {
		UserRead,
		DocumentRead,
		CommunicationsRead,
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

// HasPermission checks if a single role has the required permission
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

// HasUserPermission checks if a user's base role OR any of their active assignments has the required permission
func (m *Manager) HasUserPermission(baseRole string, assignments []string, required Permission) bool {
	if m.HasPermission(baseRole, required) {
		return true
	}
	for _, assignment := range assignments {
		if m.HasPermission(assignment, required) {
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
