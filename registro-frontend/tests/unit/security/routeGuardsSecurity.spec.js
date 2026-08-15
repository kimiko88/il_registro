import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { authGuard } from '@/router/guards'
import { useAuthStore } from '@/stores/auth'

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
        authStore.login({ id: 'u1', role: 'teacher' }, 'valid-token', 'refresh-token')

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
        authStore.login({ id: 's1', role: 'student' }, 'valid-token', 'refresh-token')

        const to = { path: '/admin/dashboard', meta: { roles: ['admin', 'superadmin'] } }
        const from = { path: '/student' }

        await authGuard(to, from, nextSpy)

        expect(nextSpy).toHaveBeenCalledWith('/student')
    })

    it('allows user with matching role permission to proceed to protected route', async () => {
        const authStore = useAuthStore()
        authStore.login({ id: 'a1', role: 'admin' }, 'valid-token', 'refresh-token')

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
            authStore.token = 'restored-access-token'
            authStore.user = { id: 't1', role: 'teacher' }
        })

        const to = { path: '/teacher', meta: { roles: ['teacher', 'coordinator'] } }
        const from = { path: '/' }

        await authGuard(to, from, nextSpy)

        expect(authStore.initAuth).toHaveBeenCalled()
        expect(nextSpy).toHaveBeenCalledWith()
    })
})
