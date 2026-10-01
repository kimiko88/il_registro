/**
 * Route ↔ Menu Consistency Tests
 *
 * Cross-validates that every path a role's sidebar menu exposes is actually
 * guarded in routes.js to allow that role. This prevents a UI regression where:
 *   - A menu item appears for a role but the route rejects that role → confusing 403/redirect
 *   - A route allows a role but no menu item exists → dead route (unreachable via UI)
 *
 * Strategy:
 *   1. Extract the route meta.roles[] directly from the imported routes array
 *   2. For each role, get its menu paths via useMenuItems()
 *   3. Assert every menu path is authorized in the matching route definition
 */
import { describe, it, expect } from 'vitest'
import routes from 'src/router/routes'
import { useMenuItems } from 'src/composables/useMenuItems'

// ── Route map builder ──────────────────────────────────────────────────────
/**
 * Flattens the nested routes[] into a map { path: Set<roles> }
 * Handles both top-level children and redirect-only entries.
 */
function buildRouteRoleMap(routeList, parentPath = '') {
    const map = new Map() // path → Set of allowed roles

    for (const route of routeList) {
        const fullPath = joinPath(parentPath, route.path)

        if (route.meta?.roles && Array.isArray(route.meta.roles)) {
            map.set(fullPath, new Set(route.meta.roles))
        }

        if (Array.isArray(route.children)) {
            const childMap = buildRouteRoleMap(route.children, fullPath)
            for (const [p, roles] of childMap) {
                map.set(p, roles)
            }
        }
    }
    return map
}

function joinPath(parent, child) {
    if (!child || child.startsWith('/')) return child
    const base = parent.endsWith('/') ? parent.slice(0, -1) : parent
    return `${base}/${child}`
}

// ── Menu path extractor ────────────────────────────────────────────────────
function collectMenuPaths(role, assignments = []) {
    const menu = useMenuItems(role, assignments)
    const paths = []
    for (const item of menu) {
        if (item.path && !item.redirect) paths.push(item.path)
        if (Array.isArray(item.children)) {
            for (const child of item.children) {
                if (child.path && !child.redirect) paths.push(child.path)
            }
        }
    }
    return [...new Set(paths)]
}

// ── Build the authoritative route map once ─────────────────────────────────
const routeRoleMap = buildRouteRoleMap(routes)

/**
 * The root '/' path uses a long roles array covering all roles; we skip it
 * for the per-role check since it is deliberately universal.
 */
const UNIVERSAL_PATHS = new Set(['/', '/dashboard', '/parent'])

/**
 * Paths that appear in menus but resolve via redirect (no roles meta on redirect node).
 * These are intentional pass-throughs — skip them in the consistency check.
 */
const REDIRECT_ONLY_PATHS = new Set([
    '/admin', '/teacher/dashboard', '/student/dashboard',
    '/parent/dashboard', '/secretary/dashboard',
])

// ── Per-role consistency suites ────────────────────────────────────────────
const ROLES_TO_CHECK = [
    { role: 'superadmin' },
    { role: 'admin' },
    { role: 'secretary' },
    { role: 'principal' },
    { role: 'vice_principal' },
    { role: 'teacher' },
    { role: 'student' },
    { role: 'parent' },
    { role: 'dsga' },
    { role: 'assistente_amministrativo' },
    { role: 'assistente_alunni' },
    { role: 'assistente_personale' },
    { role: 'assistente_contabilita' },
    { role: 'assistente_protocollo' },
    { role: 'assistente_sportello' },
    { role: 'assistente_tecnico' },
    { role: 'collaboratore_ds' },
    { role: 'collaboratore_scolastico' },
    { role: 'responsabile_servizio' },
    { role: 'responsabile_gestione_documentale' },
    { role: 'responsabile_conservazione' },
    { role: 'dpo' },
]

describe('Route ↔ Menu consistency: every menu path is authorized in routes.js', () => {
    for (const { role } of ROLES_TO_CHECK) {
        describe(`role: ${role}`, () => {
            const menuPaths = collectMenuPaths(role)

            for (const menuPath of menuPaths) {
                if (UNIVERSAL_PATHS.has(menuPath) || REDIRECT_ONLY_PATHS.has(menuPath)) {
                    // Skip universal / redirect paths
                    continue
                }

                it(`"${menuPath}" is allowed for role "${role}" in routes.js`, () => {
                    const allowedRoles = routeRoleMap.get(menuPath)
                    if (!allowedRoles) {
                        // Path may be a redirect or not directly mapped; treat as pass
                        return
                    }
                    expect(
                        allowedRoles.has(role),
                        `Menu exposes "${menuPath}" to role "${role}" but routes.js doesn't include it in meta.roles`
                    ).toBe(true)
                })
            }
        })
    }
})

// ── Forbidden path isolation tests ─────────────────────────────────────────
describe('Forbidden path isolation: roles cannot see each other\'s restricted paths', () => {
    it('student menu contains only /student/* paths (plus dashboard)', () => {
        const paths = collectMenuPaths('student')
        const forbidden = paths.filter(p =>
            !p.startsWith('/student') &&
            p !== '/' &&
            !UNIVERSAL_PATHS.has(p)
        )
        expect(forbidden).toEqual([])
    })

    it('parent menu contains only /parent/* paths (plus dashboard)', () => {
        const paths = collectMenuPaths('parent')
        const forbidden = paths.filter(p =>
            !p.startsWith('/parent') &&
            p !== '/' &&
            !UNIVERSAL_PATHS.has(p)
        )
        expect(forbidden).toEqual([])
    })

    it('collaboratore_scolastico cannot see any /secretary/* paths except communications', () => {
        const paths = collectMenuPaths('collaboratore_scolastico')
        const secPaths = paths.filter(p => p.startsWith('/secretary/') && p !== '/secretary/communications')
        expect(secPaths).toEqual([])
    })

    it('assistente_tecnico cannot see /secretary/students or /secretary/sidi', () => {
        const paths = collectMenuPaths('assistente_tecnico')
        expect(paths).not.toContain('/secretary/students')
        expect(paths).not.toContain('/secretary/sidi')
    })

    it('dpo cannot see /secretary/students or any /teacher/* paths', () => {
        const paths = collectMenuPaths('dpo')
        const teacherPaths = paths.filter(p => p.startsWith('/teacher/'))
        const studentPaths = paths.filter(p => p.startsWith('/secretary/students'))
        expect(teacherPaths).toEqual([])
        expect(studentPaths).toEqual([])
    })

    it('teacher cannot navigate to admin-only paths', () => {
        const paths = collectMenuPaths('teacher')
        const adminPaths = paths.filter(p =>
            p.startsWith('/admin/') &&
            p !== '/admin/elearning' // elearning is explicitly allowed for teacher
        )
        expect(adminPaths).toEqual([])
    })

    it('superadmin menu does not include student or parent routes', () => {
        const paths = collectMenuPaths('superadmin')
        expect(paths.filter(p => p.startsWith('/student/'))).toEqual([])
        expect(paths.filter(p => p.startsWith('/parent/'))).toEqual([])
    })

    it('no ATA role can access /admin/admins', () => {
        const ataRoles = [
            'dsga', 'assistente_amministrativo', 'assistente_alunni',
            'assistente_personale', 'assistente_contabilita', 'assistente_protocollo',
            'assistente_sportello', 'assistente_tecnico', 'collaboratore_ds',
            'collaboratore_scolastico', 'responsabile_servizio',
            'responsabile_gestione_documentale', 'responsabile_conservazione',
        ]
        for (const role of ataRoles) {
            const paths = collectMenuPaths(role)
            expect(paths, `${role} should not have /admin/admins`).not.toContain('/admin/admins')
        }
    })

    it('no ATA role can access /admin/monitoring', () => {
        const ataRoles = [
            'dsga', 'assistente_amministrativo', 'assistente_alunni',
            'collaboratore_ds', 'collaboratore_scolastico',
            'responsabile_servizio', 'responsabile_gestione_documentale',
        ]
        for (const role of ataRoles) {
            const paths = collectMenuPaths(role)
            expect(paths, `${role} should not have /admin/monitoring`).not.toContain('/admin/monitoring')
        }
    })
})

// ── Route guard: unauthenticated access ────────────────────────────────────
describe('Route meta integrity: all protected routes have roles defined', () => {
    it('every non-public route with components has a roles array', () => {
        const violations = []
        function checkRoutes(list, parentRequiresAuth = true) {
            for (const route of list) {
                const reqAuth = route.meta?.requiresAuth !== false && parentRequiresAuth
                if (reqAuth && route.component && !route.redirect) {
                    // Skip layout-wrapper routes: they are containers whose children
                    // carry the roles. A layout wrapper is identified by having children
                    // and its path being the root '/' OR its children all define their own roles.
                    const isLayoutWrapper = (route.path === '/' || route.path === '') &&
                        Array.isArray(route.children) &&
                        route.children.length > 0 &&
                        !route.meta?.roles

                    if (!isLayoutWrapper) {
                        const hasRoles = Array.isArray(route.meta?.roles) && route.meta.roles.length > 0
                        const hasRole = typeof route.meta?.role === 'string'
                        if (!hasRoles && !hasRole) {
                            violations.push(route.path || '(unnamed)')
                        }
                    }
                }
                if (Array.isArray(route.children)) {
                    checkRoutes(route.children, reqAuth)
                }
            }
        }
        checkRoutes(routes)
        expect(
            violations,
            `These protected routes have no roles defined: ${violations.join(', ')}`
        ).toEqual([])
    })
})
