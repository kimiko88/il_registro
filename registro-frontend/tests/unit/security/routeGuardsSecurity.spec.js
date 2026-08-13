import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { authGuard } from '@/router/guards'
import { useAuthStore } from '@/stores/auth'

describe('Router Security Guards — Navigation & Role Access Control', () => {
    let nextSpy

    beforeEach(() => {
        setActivePinia(createPinia())
        nextSpy = vi.fn()
    })

    it('allows unauthenticated users to access public routes like /login', () => {
        const authStore = useAuthStore()
        authStore.token = null
        authStore.user = null

        const to = { path: '/login' }
        const from = { path: '/' }

        authGuard(to, from, nextSpy)

        expect(nextSpy).toHaveBeenCalledWith()
    })

    it('redirects authenticated users away from /login to their role dashboard', () => {
        const authStore = useAuthStore()
        authStore.token = 'valid-token'
        authStore.user = { id: 'u1', role: 'teacher' }

        const to = { path: '/login' }
        const from = { path: '/' }

        authGuard(to, from, nextSpy)

        expect(nextSpy).toHaveBeenCalledWith('/teacher')
    })

    it('redirects unauthenticated users attempting to access protected route to /login', () => {
        const authStore = useAuthStore()
        authStore.token = null
        authStore.user = null

        const to = { path: '/admin/dashboard', meta: { roles: ['admin', 'superadmin'] } }
        const from = { path: '/' }

        authGuard(to, from, nextSpy)

        expect(nextSpy).toHaveBeenCalledWith('/login')
    })

    it('blocks user with insufficient role permissions and redirects to own dashboard', () => {
        const authStore = useAuthStore()
        authStore.token = 'valid-token'
        authStore.user = { id: 's1', role: 'student' }

        const to = { path: '/admin/dashboard', meta: { roles: ['admin', 'superadmin'] } }
        const from = { path: '/student' }

        authGuard(to, from, nextSpy)

        expect(nextSpy).toHaveBeenCalledWith('/student')
    })

    it('allows user with matching role permission to proceed to protected route', () => {
        const authStore = useAuthStore()
        authStore.token = 'valid-token'
        authStore.user = { id: 'a1', role: 'admin' }

        const to = { path: '/admin/dashboard', meta: { roles: ['admin', 'superadmin'] } }
        const from = { path: '/admin/dashboard' }

        authGuard(to, from, nextSpy)

        expect(nextSpy).toHaveBeenCalledWith()
    })
})
