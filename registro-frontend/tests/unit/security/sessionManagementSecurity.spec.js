import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
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

describe('Session Security & State Sanitization on Logout', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorage.clear()
        sessionStorage.clear()

        // Populate mock local storage with sensitive user session tokens
        localStorage.setItem('token', createMockJWT('teacher'))
        localStorage.setItem('user', JSON.stringify({ id: 'u1', role: 'teacher', email: 'teacher@scuola.it' }))
        sessionStorage.setItem('current_class_id', 'class-99')
    })

    it('clears JWT token and user info from localStorage on logout', () => {
        const authStore = useAuthStore()
        authStore.token = createMockJWT('teacher')
        authStore.user = { id: 'u1', role: 'teacher' }

        authStore.logout()

        expect(authStore.token).toBeNull()
        expect(authStore.user).toBeNull()
        expect(localStorage.getItem('token')).toBeNull()
        expect(localStorage.getItem('user')).toBeNull()
    })

    it('resets isAuthenticated getter to false after logout', () => {
        const authStore = useAuthStore()
        authStore.token = createMockJWT('teacher')
        authStore.user = { id: 'u1', role: 'teacher' }

        expect(authStore.isAuthenticated).toBe(true)

        authStore.logout()

        expect(authStore.isAuthenticated).toBe(false)
    })
})

