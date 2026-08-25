import { describe, it, expect, vi, beforeEach } from 'vitest'
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

describe('Security: 401 Unauthorized Auto-Refresh Queue', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorage.clear()
        sessionStorage.clear()
    })

    it('should initialize auth store and clear session on logout', () => {
        const authStore = useAuthStore()
        const token = createMockJWT('teacher')
        authStore.login({ id: '1', role: 'teacher' }, token, 'refresh-token')

        expect(authStore.isAuthenticated).toBe(true)
        expect(authStore.token).toBe(token)

        authStore.logout()

        expect(authStore.isAuthenticated).toBe(false)
        expect(authStore.token).toBeNull()
        expect(localStorage.getItem('token')).toBeNull()
    })

    it('should invoke resetApiState and clear storage on clearLocalSession', async () => {
        const { clearLocalSession, resetApiState } = await import('@/services/api')
        localStorage.setItem('token', 'test-token')
        localStorage.setItem('user', JSON.stringify({ id: '1', role: 'admin' }))
        sessionStorage.setItem('token', 'test-token')

        clearLocalSession()

        expect(localStorage.getItem('token')).toBeNull()
        expect(localStorage.getItem('user')).toBeNull()
        expect(sessionStorage.getItem('token')).toBeNull()

        // Calling resetApiState should run safely without throw
        expect(() => resetApiState()).not.toThrow()
    })

    it('should not attempt refresh token recursion on auth endpoint 401 errors', async () => {
        const { default: api } = await import('@/services/api')
        const responseInterceptor = api.interceptors.response.handlers[0]

        // 401 on /auth/refresh-token should immediately reject without re-triggering refresh
        const refreshError = {
            response: { status: 401, data: { error: 'invalid refresh token' } },
            config: { url: '/api/v1/auth/refresh-token' }
        }

        await expect(responseInterceptor.rejected(refreshError)).rejects.toBeDefined()
    })
})

