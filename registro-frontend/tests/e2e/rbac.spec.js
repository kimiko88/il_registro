import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createTestingPinia } from '@pinia/testing'
import { useAuthStore } from '@/stores/auth'
import { authGuard } from '@/router/guards'

function createFakeJwt(role, expiresInSec = 3600) {
    const header = btoa(JSON.stringify({ alg: 'HS256', typ: 'JWT' }))
    const payload = btoa(JSON.stringify({
        sub: 'user-1',
        role: role,
        user_role: role,
        exp: Math.floor(Date.now() / 1000) + expiresInSec
    }))
    return `${header}.${payload}.signature`
}

describe('RBAC Route Protection E2E', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        sessionStorage.clear()
        localStorage.clear()
        createTestingPinia({
            createSpy: vi.fn,
            initialState: {
                auth: {
                    user: null,
                    token: null,
                    isAuthenticated: true,
                    isInitializing: false
                }
            }
        })
    })

    it('R01 — student redirected from /admin/dashboard to /student', async () => {
        const authStore = useAuthStore()
        authStore.token = createFakeJwt('student')
        authStore.user = { role: 'student', email: 'student@school.it' }

        const next = vi.fn()
        const to = { path: '/admin/dashboard', meta: { requiresAuth: true, roles: ['superadmin', 'admin'] } }
        const from = { path: '/' }

        await authGuard(to, from, next)

        expect(next).toHaveBeenCalledWith('/student')
    })

    it('R02 — parent redirected from /teacher/grades to /parent', async () => {
        const authStore = useAuthStore()
        authStore.token = createFakeJwt('parent')
        authStore.user = { role: 'parent', email: 'parent@school.it' }

        const next = vi.fn()
        const to = { path: '/teacher/grades', meta: { requiresAuth: true, roles: ['teacher', 'coordinator'] } }
        const from = { path: '/' }

        await authGuard(to, from, next)

        expect(next).toHaveBeenCalledWith('/parent')
    })

    it('R03 — teacher redirected from /secretary/users to /teacher', async () => {
        const authStore = useAuthStore()
        authStore.token = createFakeJwt('teacher')
        authStore.user = { role: 'teacher', email: 'teacher@school.it' }

        const next = vi.fn()
        const to = { path: '/secretary/users', meta: { requiresAuth: true, roles: ['superadmin', 'admin', 'secretary'] } }
        const from = { path: '/' }

        await authGuard(to, from, next)

        expect(next).toHaveBeenCalledWith('/teacher')
    })

    it('R04 — admin can access /admin/dashboard', async () => {
        const authStore = useAuthStore()
        authStore.token = createFakeJwt('admin')
        authStore.user = { role: 'admin', email: 'admin@school.it' }

        const next = vi.fn()
        const to = { path: '/admin/dashboard', meta: { requiresAuth: true, roles: ['superadmin', 'admin'] } }
        const from = { path: '/' }

        await authGuard(to, from, next)

        expect(next).toHaveBeenCalledWith()
    })

    it('R05 — secretary can access /secretary', async () => {
        const authStore = useAuthStore()
        authStore.token = createFakeJwt('secretary')
        authStore.user = { role: 'secretary', email: 'secretary@school.it' }

        const next = vi.fn()
        const to = { path: '/secretary', meta: { requiresAuth: true, roles: ['secretary', 'admin', 'superadmin'] } }
        const from = { path: '/' }

        await authGuard(to, from, next)

        expect(next).toHaveBeenCalledWith()
    })

    it('R06 — principal redirects to /secretary dashboard', async () => {
        const authStore = useAuthStore()
        authStore.token = createFakeJwt('principal')
        authStore.user = { role: 'principal', email: 'preside@school.it' }

        const next = vi.fn()
        const to = { path: '/admin/schools', meta: { requiresAuth: true, roles: ['superadmin', 'admin'] } }
        const from = { path: '/' }

        await authGuard(to, from, next)

        expect(next).toHaveBeenCalledWith('/secretary')
    })
})
