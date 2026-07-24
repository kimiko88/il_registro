import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'

describe('Security: 401 Unauthorized Auto-Refresh Queue', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorage.clear()
        sessionStorage.clear()
    })

    it('should initialize auth store and clear session on logout', () => {
        const authStore = useAuthStore()
        authStore.login({ id: '1', role: 'teacher' }, 'access-token', 'refresh-token')

        expect(authStore.isAuthenticated).toBe(true)
        expect(authStore.token).toBe('access-token')

        authStore.logout()

        expect(authStore.isAuthenticated).toBe(false)
        expect(authStore.token).toBeNull()
        expect(localStorage.getItem('token')).toBeNull()
    })
})
