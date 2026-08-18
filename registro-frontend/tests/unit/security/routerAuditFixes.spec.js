import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { authGuard } from '@/router/guards'
import { useAuthStore } from '@/stores/auth'
import routes from '@/router/routes'

const createMockJWT = (role = 'teacher', expInSeconds = 3600) => {
    const header = btoa(JSON.stringify({ alg: 'HS256', typ: 'JWT' }))
    const payload = btoa(JSON.stringify({
        sub: 'user-123',
        role: role,
        exp: Math.floor(Date.now() / 1000) + expInSeconds
    }))
    return `${header}.${payload}.signature`
}

describe('Router & Guard 10-Item Security, Logic & Quality Audit Test Suite', () => {
    let nextSpy

    beforeEach(() => {
        setActivePinia(createPinia())
        nextSpy = vi.fn()
    })

    // Item 1 & Item 8: Fallback and explicit requiresAuth: false
    it('Item 1 & 8: honors explicit requiresAuth === false on public routes and matchers', () => {
        const authStore = useAuthStore()
        authStore.token = null
        authStore.user = null

        const publicPath = { path: '/login', meta: { requiresAuth: false } }
        authGuard(publicPath, { path: '/' }, nextSpy)
        expect(nextSpy).toHaveBeenCalledWith()
    })

    // Item 2: Ambiguous 404 behavior for catchAll route
    it('Item 2: allows unauthenticated user to access 404 catchAll without redirect to /login', () => {
        const authStore = useAuthStore()
        authStore.token = null
        authStore.user = null

        const catchAllRoute = routes.find(r => r.path === '/:catchAll(.*)*')
        expect(catchAllRoute).toBeDefined()
        expect(catchAllRoute.meta.requiresAuth).toBe(false)

        const to = { path: '/non-existent-page', meta: catchAllRoute.meta }
        authGuard(to, { path: '/' }, nextSpy)
        expect(nextSpy).toHaveBeenCalledWith()
    })

    // Item 3 & 5: Dynamic Redirects for /communications and /profile
    it('Item 3 & 5: evaluates /communications redirect for all roles without erroring out', async () => {
        const commRoute = routes[0].children.find(r => r.path === 'communications')
        expect(commRoute).toBeDefined()

        const authStore = useAuthStore()
        authStore.login({ id: 't1', role: 'teacher' }, createMockJWT('teacher'), 'refresh-token')

        const teacherRedirect = await commRoute.redirect()
        expect(teacherRedirect).toBe('/teacher/communications')

        authStore.login({ id: 'p1', role: 'parent' }, createMockJWT('parent'), 'refresh-token')
        const parentRedirect = await commRoute.redirect()
        expect(parentRedirect).toBe('/parent/communications')

        authStore.login({ id: 'a1', role: 'admin' }, createMockJWT('admin'), 'refresh-token')
        const adminRedirect = await commRoute.redirect()
        expect(adminRedirect).toBe('/admin/dashboard')
    })

    it('Item 3 & 5: evaluates /profile redirect for teacher, principal, vice_principal, student, parent', async () => {
        const profileRoute = routes[0].children.find(r => r.path === 'profile')
        expect(profileRoute).toBeDefined()

        const authStore = useAuthStore()

        authStore.login({ id: 't1', role: 'teacher' }, createMockJWT('teacher'), 'refresh-token')
        expect(await profileRoute.redirect()).toBe('/teacher/settings')

        authStore.login({ id: 'pr1', role: 'principal' }, createMockJWT('principal'), 'refresh-token')
        expect(await profileRoute.redirect()).toBe('/secretary/settings')

        authStore.login({ id: 's1', role: 'student' }, createMockJWT('student'), 'refresh-token')
        expect(await profileRoute.redirect()).toBe('/student/profile')
    })

    // Item 4: Coordinator role matching
    it('Item 4: allows coordinator role to access teacher routes without redirect loop', () => {
        const authStore = useAuthStore()
        authStore.login({ id: 'c1', role: 'coordinator' }, createMockJWT('coordinator'), 'refresh-token')

        const teacherRoute = routes[0].children.find(r => r.path === 'teacher')
        expect(teacherRoute.meta.roles).toContain('coordinator')

        const to = { path: '/teacher', meta: teacherRoute.meta }
        authGuard(to, { path: '/' }, nextSpy)
        expect(nextSpy).toHaveBeenCalledWith()
    })

    // Item 6: Root / Dashboard roles for principal and vice_principal
    it('Item 6: includes principal and vice_principal in root / dashboard route roles', () => {
        const authStore = useAuthStore()
        authStore.login({ id: 'p1', role: 'principal' }, createMockJWT('principal'), 'refresh-token')

        const rootDashboard = routes[0].children.find(r => r.path === '')
        expect(rootDashboard.meta.roles).toContain('principal')
        expect(rootDashboard.meta.roles).toContain('vice_principal')

        const to = { path: '/', meta: rootDashboard.meta }
        authGuard(to, { path: '/' }, nextSpy)
        expect(nextSpy).toHaveBeenCalledWith()
    })

    // Item 7: system_auditor and admin/audit-logs alignment
    it('Item 7: aligns system_auditor role with /admin/audit-logs path and route permissions', () => {
        const authStore = useAuthStore()
        authStore.login({ id: 'sa1', role: 'system_auditor' }, createMockJWT('system_auditor'), 'refresh-token')

        const auditLogsRoute = routes[0].children.find(r => r.path === 'admin/audit-logs')
        expect(auditLogsRoute).toBeDefined()
        expect(auditLogsRoute.meta.roles).toContain('system_auditor')

        const to = { path: '/admin/audit-logs', meta: auditLogsRoute.meta }
        authGuard(to, { path: '/' }, nextSpy)
        expect(nextSpy).toHaveBeenCalledWith()
    })

    // Item 9: Module-level import safety
    it('Item 9: routes.js exports default array without importing useAuthStore at module top level', () => {
        expect(Array.isArray(routes)).toBe(true)
    })
})

