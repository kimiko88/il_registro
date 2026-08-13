import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'

describe('Session Security & State Sanitization on Logout', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorage.clear()
        sessionStorage.clear()

        // Populate mock local storage with sensitive user session tokens
        localStorage.setItem('token', 'jwt-secret-token-value')
        localStorage.setItem('user', JSON.stringify({ id: 'u1', role: 'teacher', email: 'teacher@scuola.it' }))
        sessionStorage.setItem('current_class_id', 'class-99')
    })

    it('clears JWT token and user info from localStorage on logout', () => {
        const authStore = useAuthStore()
        authStore.token = 'jwt-secret-token-value'
        authStore.user = { id: 'u1', role: 'teacher' }

        authStore.logout()

        expect(authStore.token).toBeNull()
        expect(authStore.user).toBeNull()
        expect(localStorage.getItem('token')).toBeNull()
        expect(localStorage.getItem('user')).toBeNull()
    })

    it('resets isAuthenticated getter to false after logout', () => {
        const authStore = useAuthStore()
        authStore.token = 'jwt-secret-token-value'
        authStore.user = { id: 'u1', role: 'teacher' }

        expect(authStore.isAuthenticated).toBe(true)

        authStore.logout()

        expect(authStore.isAuthenticated).toBe(false)
    })
})
