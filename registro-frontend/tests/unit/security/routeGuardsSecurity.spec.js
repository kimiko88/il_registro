import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { authGuard } from '@/router/guards'
import { useAuthStore } from '@/stores/auth'

const createMockJWT = (role = 'teacher', expInSeconds = 3600) => {
    const header = btoa(JSON.stringify({ alg: 'HS256', typ: 'JWT' }))
    const payload = btoa(JSON.stringify({
        sub: 'user-123',
        role: role,
        exp: Math.floor(Date.now() / 1000) + expInSeconds
    }))
    return `${header}.${payload}.signature`
}

describe('Router Security Guards — Navigation & Role Access Control', () => {
    let nextSpy

    beforeEach(() => {
        setActivePinia(createPinia())
        localStorage.clear()
        sessionStorage.clear()
        nextSpy = vi.fn()
    })

    it('allows unauthenticated users to access public routes like /login', async () => {
        const authStore = useAuthStore()
        authStore.token = null
        authStore.user = null

        const to = { path: '/login' }
        const from = { path: '/' }

        await authGuard(to, from, nextSpy)

        expect(nextSpy).toHaveBeenCalledWith()
    })

    it('redirects authenticated users away from /login to their role dashboard', async () => {
        const authStore = useAuthStore()
        authStore.login({ id: 'u1', role: 'teacher' }, createMockJWT('teacher'), 'refresh-token')

        const to = { path: '/login' }
        const from = { path: '/' }

        await authGuard(to, from, nextSpy)

        expect(nextSpy).toHaveBeenCalledWith('/teacher')
    })

    it('redirects unauthenticated users attempting to access protected route to /login', async () => {
        const authStore = useAuthStore()
        authStore.token = null
        authStore.user = null

        const to = { path: '/admin/dashboard', meta: { roles: ['admin', 'superadmin'] } }
        const from = { path: '/' }

        await authGuard(to, from, nextSpy)

        expect(nextSpy).toHaveBeenCalledWith({ path: '/login', query: { reason: 'session_expired' } })
    })

    it('blocks user with insufficient role permissions and redirects to own dashboard', async () => {
        const authStore = useAuthStore()
        authStore.login({ id: 's1', role: 'student' }, createMockJWT('student'), 'refresh-token')

        const to = { path: '/admin/dashboard', meta: { roles: ['admin', 'superadmin'] } }
        const from = { path: '/student' }

        await authGuard(to, from, nextSpy)

        expect(nextSpy).toHaveBeenCalledWith('/student')
    })

    it('allows user with matching role permission to proceed to protected route', async () => {
        const authStore = useAuthStore()
        authStore.login({ id: 'a1', role: 'admin' }, createMockJWT('admin'), 'refresh-token')

        const to = { path: '/admin/dashboard', meta: { roles: ['admin', 'superadmin'] } }
        const from = { path: '/admin/dashboard' }

        await authGuard(to, from, nextSpy)

        expect(nextSpy).toHaveBeenCalledWith()
    })

    it('allows unauthenticated users to view 404 non-existent URLs directly', async () => {
        const authStore = useAuthStore()
        authStore.token = null
        authStore.user = null

        const to = { path: '/some-invalid-path', meta: { requiresAuth: false } }
        const from = { path: '/' }

        await authGuard(to, from, nextSpy)

        expect(nextSpy).toHaveBeenCalledWith()
    })

    it('restores session silently on page reload when user data exists in storage', async () => {
        const authStore = useAuthStore()
        authStore.token = null
        localStorage.setItem('user', JSON.stringify({ id: 't1', role: 'teacher' }))

        vi.spyOn(authStore, 'initAuth').mockImplementation(async () => {
            authStore.token = createMockJWT('teacher')
            authStore.user = { id: 't1', role: 'teacher' }
        })

        const to = { path: '/teacher', meta: { roles: ['teacher', 'coordinator'] } }
        const from = { path: '/' }

        await authGuard(to, from, nextSpy)

        expect(authStore.initAuth).toHaveBeenCalled()
        expect(nextSpy).toHaveBeenCalledWith()
    })

    describe('ATA Roles Route Permissions', () => {
        it('allows dsga to access authorized secretary routes', async () => {
            const authStore = useAuthStore()
            authStore.login({ id: 'd1', role: 'dsga' }, createMockJWT('dsga'), 'refresh-token')

            const to = { path: '/secretary/users', meta: { roles: ['secretary', 'principal', 'vice_principal', 'dsga', 'assistente_amministrativo'] } }
            await authGuard(to, { path: '/ata' }, nextSpy)
            expect(nextSpy).toHaveBeenCalledWith()
        })

        it('allows assistente_amministrativo to access authorized routes', async () => {
            const authStore = useAuthStore()
            authStore.login({ id: 'aa1', role: 'assistente_amministrativo' }, createMockJWT('assistente_amministrativo'), 'refresh-token')

            const to = { path: '/secretary/students', meta: { roles: ['secretary', 'principal', 'vice_principal', 'assistente_amministrativo'] } }
            await authGuard(to, { path: '/ata' }, nextSpy)
            expect(nextSpy).toHaveBeenCalledWith()
        })

        it('allows collaboratore_ds to access substitutions and timetable', async () => {
            const authStore = useAuthStore()
            authStore.login({ id: 'cds1', role: 'collaboratore_ds' }, createMockJWT('collaboratore_ds'), 'refresh-token')

            const to = { path: '/secretary/substitutions', meta: { roles: ['secretary', 'admin', 'superadmin', 'principal', 'vice_principal', 'dsga', 'collaboratore_ds'] } }
            await authGuard(to, { path: '/ata' }, nextSpy)
            expect(nextSpy).toHaveBeenCalledWith()
        })

        it('allows collaboratore_scolastico to access communications but blocks from sidi', async () => {
            const authStore = useAuthStore()
            authStore.login({ id: 'cs1', role: 'collaboratore_scolastico' }, createMockJWT('collaboratore_scolastico'), 'refresh-token')

            // Allowed to communications
            const toAllowed = { path: '/secretary/communications', meta: { roles: ['secretary', 'principal', 'vice_principal', 'dsga', 'assistente_amministrativo', 'collaboratore_ds', 'collaboratore_scolastico'] } }
            await authGuard(toAllowed, { path: '/ata' }, nextSpy)
            expect(nextSpy).toHaveBeenCalledWith()

            // Blocked from sidi -> redirected to /ata
            const nextSpyBlocked = vi.fn()
            const toBlocked = { path: '/secretary/sidi', meta: { roles: ['secretary', 'admin', 'superadmin', 'principal', 'vice_principal', 'dsga'] } }
            await authGuard(toBlocked, { path: '/ata' }, nextSpyBlocked)
            expect(nextSpyBlocked).toHaveBeenCalledWith('/ata')
        })
    })
})


