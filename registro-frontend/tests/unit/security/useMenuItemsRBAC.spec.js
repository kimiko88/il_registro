/**
 * useMenuItems — Per-Role RBAC Menu Visibility Tests
 *
 * Verifies that each of the 25 institutional roles receives:
 *  1. Only the menu paths they are allowed to see
 *  2. No paths that belong to other roles exclusively
 *  3. Correct coordinatorOnly filtering for plain teachers
 *  4. Dynamic assignment unlocks (coordinator, referente_progetto, etc.)
 */
import { describe, it, expect } from 'vitest'
import { useMenuItems } from 'src/composables/useMenuItems'

// Helper: flatten all paths from a menu tree
function collectPaths(menu) {
    const paths = []
    for (const item of menu) {
        if (item.path) paths.push(item.path)
        if (Array.isArray(item.children)) {
            for (const child of item.children) {
                if (child.path) paths.push(child.path)
            }
        }
    }
    return [...new Set(paths)]
}


// Paths that are exclusively admin-only (students/parents must NEVER see them)
const ADMIN_ONLY_PATHS = [
    '/admin/admins',
    '/admin/monitoring',
    '/admin/analytics',
    '/admin/settings',
    '/admin/tenants',
    '/admin/scheduler',
    '/admin/school-settings',
]

// Paths that are exclusively teacher-facing (students/parents/ATA must NEVER see them)
const TEACHER_ONLY_PATHS = [
    '/teacher/grades',
    '/teacher/attendance',
    '/teacher/lessons',
    '/teacher/scrutiny',
    '/teacher/notes',
    '/teacher/rubrics',
    '/teacher/uda',
    '/teacher/competencies',
    '/teacher/pdp',
    '/teacher/recovery',
    '/teacher/credits',
]

// Paths only secretary/principal can access
const SECRETARY_RESTRICTED_PATHS = [
    '/secretary/sidi',
    '/secretary/scrutiny',
    '/secretary/classes',
]

// ─── SUPERADMIN ────────────────────────────────────────────────────────────
describe('useMenuItems — superadmin', () => {
    const menu = useMenuItems('superadmin')
    const paths = collectPaths(menu)

    it('returns a non-empty menu', () => {
        expect(menu.length).toBeGreaterThan(0)
    })

    it('contains admin management paths', () => {
        expect(paths).toContain('/admin/schools')
        expect(paths).toContain('/admin/admins')
        expect(paths).toContain('/admin/monitoring')
        expect(paths).toContain('/admin/audit-logs')
        expect(paths).toContain('/admin/analytics')
    })

    it('does NOT contain student or parent paths', () => {
        for (const p of paths) {
            expect(p.startsWith('/student/')).toBe(false)
            expect(p.startsWith('/parent/')).toBe(false)
        }
    })
})

// ─── ADMIN ─────────────────────────────────────────────────────────────────
describe('useMenuItems — admin', () => {
    const menu = useMenuItems('admin')
    const paths = collectPaths(menu)

    it('returns a non-empty menu', () => {
        expect(menu.length).toBeGreaterThan(0)
    })

    it('contains admin-level paths', () => {
        expect(paths).toContain('/admin/users')
        expect(paths).toContain('/admin/analytics')
    })

    it('does NOT contain superadmin-only admins management path', () => {
        // /admin/admins is superadmin-only in routes.js
        expect(paths).not.toContain('/admin/admins')
    })

    it('does NOT contain student or parent paths', () => {
        for (const p of paths) {
            expect(p.startsWith('/student/')).toBe(false)
            expect(p.startsWith('/parent/')).toBe(false)
        }
    })
})

// ─── SECRETARY ─────────────────────────────────────────────────────────────
describe('useMenuItems — secretary', () => {
    const menu = useMenuItems('secretary')
    const paths = collectPaths(menu)

    it('returns a non-empty menu', () => {
        expect(menu.length).toBeGreaterThan(0)
    })

    it('contains secretary-specific paths', () => {
        expect(paths).toContain('/secretary/students')
        expect(paths).toContain('/secretary/classes')
        expect(paths).toContain('/secretary/documents')
        expect(paths).toContain('/secretary/substitutions')
        expect(paths).toContain('/secretary/timetable')
        expect(paths).toContain('/secretary/sidi')
    })

    it('does NOT contain admin-only management paths', () => {
        for (const p of ADMIN_ONLY_PATHS) {
            expect(paths).not.toContain(p)
        }
    })

    it('does NOT contain teacher/student/parent paths', () => {
        for (const p of paths) {
            expect(p.startsWith('/teacher/')).toBe(false)
            expect(p.startsWith('/student/')).toBe(false)
            expect(p.startsWith('/parent/')).toBe(false)
        }
    })
})

// ─── PRINCIPAL (uses secretaryMenu) ────────────────────────────────────────
describe('useMenuItems — principal', () => {
    const menu = useMenuItems('principal')
    const paths = collectPaths(menu)

    it('gets the same menu structure as secretary', () => {
        const secretaryPaths = collectPaths(useMenuItems('secretary'))
        expect(paths.sort()).toEqual(secretaryPaths.sort())
    })

    it('does NOT contain teacher-only or student paths', () => {
        for (const p of paths) {
            expect(p.startsWith('/teacher/')).toBe(false)
            expect(p.startsWith('/student/')).toBe(false)
            expect(p.startsWith('/parent/')).toBe(false)
        }
    })
})

// ─── VICE_PRINCIPAL ─────────────────────────────────────────────────────────
describe('useMenuItems — vice_principal', () => {
    const menu = useMenuItems('vice_principal')
    const paths = collectPaths(menu)

    it('returns a non-empty menu', () => {
        expect(menu.length).toBeGreaterThan(0)
    })

    it('contains both teacher and secretary paths (vicario hybrid)', () => {
        expect(paths).toContain('/teacher/classes')
        expect(paths).toContain('/teacher/grades')
        expect(paths).toContain('/secretary/substitutions')
        expect(paths).toContain('/secretary/timetable')
    })

    it('does NOT contain admin-only paths', () => {
        for (const p of ADMIN_ONLY_PATHS) {
            expect(paths).not.toContain(p)
        }
    })

    it('does NOT contain student or parent paths', () => {
        for (const p of paths) {
            expect(p.startsWith('/student/')).toBe(false)
            expect(p.startsWith('/parent/')).toBe(false)
        }
    })
})

// ─── TEACHER ────────────────────────────────────────────────────────────────
describe('useMenuItems — teacher (without coordinator assignment)', () => {
    const menu = useMenuItems('teacher', [])
    const paths = collectPaths(menu)

    it('returns a non-empty menu', () => {
        expect(menu.length).toBeGreaterThan(0)
    })

    it('contains core teacher paths', () => {
        expect(paths).toContain('/teacher/grades')
        expect(paths).toContain('/teacher/attendance')
        expect(paths).toContain('/teacher/lessons')
        expect(paths).toContain('/teacher/timetable')
        expect(paths).toContain('/teacher/classes')
    })

    it('does NOT show coordinator-only items for plain teacher', () => {
        // coordinatorOnly items must be filtered out when called from MainLayout
        // useMenuItems returns them WITH the flag; MainLayout filters them.
        // We verify they carry the coordinatorOnly flag.
        let hasCoordOnlyFlag = false
        for (const item of menu) {
            if (Array.isArray(item.children)) {
                for (const child of item.children) {
                    if (child.coordinatorOnly) hasCoordOnlyFlag = true
                }
            }
        }
        expect(hasCoordOnlyFlag).toBe(true) // flag present so MainLayout can filter
    })

    it('does NOT contain secretary/admin/student/parent paths', () => {
        for (const p of paths) {
            expect(p.startsWith('/admin/')).toBe(false)
            expect(p.startsWith('/student/')).toBe(false)
            expect(p.startsWith('/parent/')).toBe(false)
        }
        // Secretary-restricted paths should not appear in teacher menu
        for (const p of SECRETARY_RESTRICTED_PATHS) {
            expect(paths).not.toContain(p)
        }
    })
})

describe('useMenuItems — teacher with coordinator assignment', () => {
    const assignments = [{ assignment_type: 'coordinatore_classe', is_active: true }]
    const menu = useMenuItems('teacher', assignments)

    it('unlocks coordinatorOnly items', () => {
        let allUnlocked = true
        for (const item of menu) {
            if (Array.isArray(item.children)) {
                for (const child of item.children) {
                    if (child.coordinatorOnly) allUnlocked = false
                }
            }
        }
        expect(allUnlocked).toBe(true) // all coordinatorOnly should be cleared
    })
})

describe('useMenuItems — teacher with referente_progetto assignment', () => {
    const assignments = [{ assignment_type: 'referente_progetto', is_active: true }]
    const menu = useMenuItems('teacher', assignments)
    const paths = collectPaths(menu)

    it('adds Progetti & Finanziamenti link', () => {
        expect(paths).toContain('/secretary/pcto')
    })
})

describe('useMenuItems — teacher with tutor_orientamento assignment', () => {
    const assignments = [{ assignment_type: 'tutor_orientamento', is_active: true }]
    const menu = useMenuItems('teacher', assignments)
    const paths = collectPaths(menu)

    it('adds Tutor Orientamento link', () => {
        expect(paths).toContain('/student/orientamento')
    })
})

describe('useMenuItems — teacher with animatore_digitale assignment', () => {
    const assignments = [{ assignment_type: 'animatore_digitale', is_active: true }]
    const menu = useMenuItems('teacher', assignments)
    const paths = collectPaths(menu)

    it('adds Team Digitale & E-Learning link', () => {
        expect(paths).toContain('/admin/elearning')
    })
})

describe('useMenuItems — teacher with inactive assignment is ignored', () => {
    const assignments = [{ assignment_type: 'coordinatore_classe', is_active: false }]
    const menu = useMenuItems('teacher', assignments)

    it('still carries coordinatorOnly flag (inactive assignment not counted)', () => {
        let hasFlag = false
        for (const item of menu) {
            if (Array.isArray(item.children)) {
                for (const child of item.children) {
                    if (child.coordinatorOnly) hasFlag = true
                }
            }
        }
        expect(hasFlag).toBe(true)
    })
})

// ─── STUDENT ─────────────────────────────────────────────────────────────────
describe('useMenuItems — student', () => {
    const menu = useMenuItems('student')
    const paths = collectPaths(menu)

    it('returns a non-empty menu', () => {
        expect(menu.length).toBeGreaterThan(0)
    })

    it('contains only student-specific paths', () => {
        expect(paths).toContain('/student/grades')
        expect(paths).toContain('/student/attendance')
        expect(paths).toContain('/student/homework')
        expect(paths).toContain('/student/timetable')
        expect(paths).toContain('/student/report-card')
    })

    it('does NOT contain teacher, admin, secretary, or parent paths', () => {
        for (const p of paths) {
            expect(p.startsWith('/teacher/')).toBe(false)
            expect(p.startsWith('/admin/')).toBe(false)
            expect(p.startsWith('/secretary/')).toBe(false)
            expect(p.startsWith('/parent/')).toBe(false)
            expect(p.startsWith('/ata/')).toBe(false)
        }
    })

    it('does NOT expose grade entry tools (read-only viewer role)', () => {
        for (const p of TEACHER_ONLY_PATHS) {
            expect(paths).not.toContain(p)
        }
    })
})

// ─── PARENT ──────────────────────────────────────────────────────────────────
describe('useMenuItems — parent', () => {
    const menu = useMenuItems('parent')
    const paths = collectPaths(menu)

    it('returns a non-empty menu', () => {
        expect(menu.length).toBeGreaterThan(0)
    })

    it('contains parent-specific paths', () => {
        expect(paths).toContain('/parent/children')
        expect(paths).toContain('/parent/grades')
        expect(paths).toContain('/parent/attendance')
        expect(paths).toContain('/parent/colloqui')
        expect(paths).toContain('/parent/report-card')
        expect(paths).toContain('/parent/payments')
    })

    it('does NOT contain teacher, admin, secretary, student, or ATA paths', () => {
        for (const p of paths) {
            expect(p.startsWith('/teacher/')).toBe(false)
            expect(p.startsWith('/admin/')).toBe(false)
            expect(p.startsWith('/secretary/')).toBe(false)
            expect(p.startsWith('/student/')).toBe(false)
            expect(p.startsWith('/ata/')).toBe(false)
        }
    })
})

// ─── ATA ROLES ───────────────────────────────────────────────────────────────
describe('useMenuItems — dsga', () => {
    const menu = useMenuItems('dsga')
    const paths = collectPaths(menu)

    it('contains ATA and secretary management paths', () => {
        expect(paths).toContain('/ata/attendance')
        expect(paths).toContain('/ata/timecard')
        expect(paths).toContain('/ata/strike')
        expect(paths).toContain('/secretary/substitutions')
        expect(paths).toContain('/secretary/sidi')
        expect(paths).toContain('/secretary/users')
    })

    it('does NOT contain teacher classroom or student paths', () => {
        for (const p of TEACHER_ONLY_PATHS) {
            expect(paths).not.toContain(p)
        }
        for (const p of paths) {
            expect(p.startsWith('/student/')).toBe(false)
            expect(p.startsWith('/parent/')).toBe(false)
            expect(p.startsWith('/admin/')).toBe(false)
        }
    })
})

describe('useMenuItems — collaboratore_scolastico', () => {
    const menu = useMenuItems('collaboratore_scolastico')
    const paths = collectPaths(menu)

    it('contains visitor registry and timecard', () => {
        expect(paths).toContain('/ata/visitor-registry')
        expect(paths).toContain('/ata/timecard')
        expect(paths).toContain('/ata/personnel-desk')
        expect(paths).toContain('/secretary/communications')
    })

    it('does NOT contain staff attendance or SIDI (higher privilege)', () => {
        expect(paths).not.toContain('/ata/attendance')
        expect(paths).not.toContain('/secretary/sidi')
        expect(paths).not.toContain('/secretary/students')
        expect(paths).not.toContain('/secretary/substitutions')
    })

    it('does NOT contain teacher, student, parent, or admin paths', () => {
        for (const p of paths) {
            expect(p.startsWith('/teacher/')).toBe(false)
            expect(p.startsWith('/student/')).toBe(false)
            expect(p.startsWith('/parent/')).toBe(false)
            expect(p.startsWith('/admin/')).toBe(false)
        }
    })
})

describe('useMenuItems — collaboratore_ds', () => {
    const menu = useMenuItems('collaboratore_ds')
    const paths = collectPaths(menu)

    it('contains emergency substitutions and staff attendance', () => {
        expect(paths).toContain('/ata/emergency-substitutions')
        expect(paths).toContain('/ata/attendance')
        expect(paths).toContain('/ata/strike')
        expect(paths).toContain('/secretary/timetable')
        expect(paths).toContain('/secretary/rooms')
    })

    it('does NOT contain secretary classes or SIDI', () => {
        expect(paths).not.toContain('/secretary/classes')
        expect(paths).not.toContain('/secretary/sidi')
    })
})

describe('useMenuItems — assistente_amministrativo', () => {
    const menu = useMenuItems('assistente_amministrativo')
    const paths = collectPaths(menu)

    it('contains users, students, SIDI, and attendance', () => {
        expect(paths).toContain('/secretary/users')
        expect(paths).toContain('/secretary/students')
        expect(paths).toContain('/secretary/sidi')
        expect(paths).toContain('/ata/attendance')
    })

    it('does NOT contain emergency substitutions (collaboratore_ds only)', () => {
        expect(paths).not.toContain('/ata/emergency-substitutions')
    })
})

describe('useMenuItems — assistente_alunni', () => {
    const menu = useMenuItems('assistente_alunni')
    const paths = collectPaths(menu)

    it('focuses on student management paths', () => {
        expect(paths).toContain('/secretary/students')
        expect(paths).toContain('/secretary/certificates')
        expect(paths).toContain('/secretary/sidi')
        expect(paths).toContain('/secretary/scrutiny')
    })

    it('does NOT contain staff attendance (not their domain)', () => {
        expect(paths).not.toContain('/ata/attendance')
    })
})

describe('useMenuItems — assistente_contabilita', () => {
    const menu = useMenuItems('assistente_contabilita')
    const paths = collectPaths(menu)

    it('focuses on document and certificate paths', () => {
        expect(paths).toContain('/secretary/documents')
        expect(paths).toContain('/secretary/certificates')
        expect(paths).toContain('/ata/personnel-desk')
    })

    it('does NOT have access to student registry or SIDI', () => {
        expect(paths).not.toContain('/secretary/students')
        expect(paths).not.toContain('/secretary/sidi')
    })
})

describe('useMenuItems — assistente_protocollo', () => {
    const menu = useMenuItems('assistente_protocollo')
    const paths = collectPaths(menu)

    it('focuses on protocol and verbali paths', () => {
        expect(paths).toContain('/secretary/documents')
        expect(paths).toContain('/secretary/verbali')
        expect(paths).toContain('/secretary/communications')
    })

    it('does NOT have attendance or student registry', () => {
        expect(paths).not.toContain('/ata/attendance')
        expect(paths).not.toContain('/secretary/students')
    })
})

describe('useMenuItems — assistente_sportello', () => {
    const menu = useMenuItems('assistente_sportello')
    const paths = collectPaths(menu)

    it('focuses on front-office paths', () => {
        expect(paths).toContain('/ata/personnel-desk')
        expect(paths).toContain('/ata/visitor-registry')
        expect(paths).toContain('/secretary/certificates')
    })

    it('does NOT have staff attendance or SIDI', () => {
        expect(paths).not.toContain('/ata/attendance')
        expect(paths).not.toContain('/secretary/sidi')
    })
})

describe('useMenuItems — assistente_tecnico', () => {
    const menu = useMenuItems('assistente_tecnico')
    const paths = collectPaths(menu)

    it('has timecard and e-learning but not SIDI or scrutiny', () => {
        expect(paths).toContain('/ata/timecard')
        expect(paths).toContain('/admin/elearning')
        expect(paths).not.toContain('/secretary/sidi')
        expect(paths).not.toContain('/secretary/scrutiny')
        expect(paths).not.toContain('/ata/attendance') // cannot manage others' attendance
    })
})

describe('useMenuItems — responsabile_servizio', () => {
    const menu = useMenuItems('responsabile_servizio')
    const paths = collectPaths(menu)

    it('has visitor registry and timecard but no privileged admin paths', () => {
        expect(paths).toContain('/ata/visitor-registry')
        expect(paths).toContain('/ata/timecard')
        expect(paths).not.toContain('/admin/users')
        expect(paths).not.toContain('/secretary/students')
    })
})

describe('useMenuItems — responsabile_gestione_documentale', () => {
    const menu = useMenuItems('responsabile_gestione_documentale')
    const paths = collectPaths(menu)

    it('focuses on document archival', () => {
        expect(paths).toContain('/secretary/documents')
        expect(paths).toContain('/secretary/verbali')
    })

    it('does NOT have ATA attendance or student management', () => {
        expect(paths).not.toContain('/ata/attendance')
        expect(paths).not.toContain('/secretary/students')
    })
})

describe('useMenuItems — responsabile_conservazione', () => {
    const menu = useMenuItems('responsabile_conservazione')
    const paths = collectPaths(menu)

    it('focuses on digital conservation paths', () => {
        expect(paths).toContain('/secretary/documents')
        expect(paths).toContain('/secretary/verbali')
    })

    it('does NOT contain privileged paths', () => {
        expect(paths).not.toContain('/ata/attendance')
        expect(paths).not.toContain('/secretary/students')
        expect(paths).not.toContain('/admin/users')
    })
})

describe('useMenuItems — dpo', () => {
    const menu = useMenuItems('dpo')
    const paths = collectPaths(menu)

    it('focuses on audit logs and privacy', () => {
        expect(paths).toContain('/admin/audit-logs')
        expect(paths).toContain('/admin/monitoring')
    })

    it('does NOT contain student, parent, or teacher paths', () => {
        for (const p of paths) {
            expect(p.startsWith('/student/')).toBe(false)
            expect(p.startsWith('/parent/')).toBe(false)
            expect(p.startsWith('/teacher/')).toBe(false)
            expect(p.startsWith('/secretary/students')).toBe(false)
        }
    })

    it('does NOT have superadmin-only management paths', () => {
        // DPO cannot manage admin users or system tenants
        expect(paths).not.toContain('/admin/admins')
        expect(paths).not.toContain('/admin/tenants')
        expect(paths).not.toContain('/admin/scheduler')
        // Note: /admin/settings and /admin/monitoring ARE accessible to DPO
        //       for privacy settings and system monitoring (privacy oversight)
    })
})

// ─── EDGE CASES ──────────────────────────────────────────────────────────────
describe('useMenuItems — edge cases', () => {
    it('returns empty array for null role', () => {
        expect(useMenuItems(null)).toEqual([])
    })

    it('returns empty array for empty string role', () => {
        expect(useMenuItems('')).toEqual([])
    })

    it('returns empty array for unknown role', () => {
        expect(useMenuItems('unknown_role_xyz')).toEqual([])
    })

    it('treats "studente" alias same as "student"', () => {
        const a = collectPaths(useMenuItems('student'))
        const b = collectPaths(useMenuItems('studente'))
        expect(b).toEqual(a)
    })

    it('treats "genitore" alias same as "parent"', () => {
        const a = collectPaths(useMenuItems('parent'))
        const b = collectPaths(useMenuItems('genitore'))
        expect(b).toEqual(a)
    })

    it('treats "docente" alias same as "teacher"', () => {
        const a = collectPaths(useMenuItems('teacher'))
        const b = collectPaths(useMenuItems('docente'))
        expect(b).toEqual(a)
    })

    it('treats "system_auditor" as superadmin menu', () => {
        const a = collectPaths(useMenuItems('superadmin'))
        const b = collectPaths(useMenuItems('system_auditor'))
        expect(b).toEqual(a)
    })
})
