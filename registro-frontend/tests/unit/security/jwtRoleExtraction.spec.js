import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'

describe('Security: JWT Role Extraction & Storage Tampering Defense', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorage.clear()
        sessionStorage.clear()
    })

    const createDummyJWT = (role, expInSeconds = 3600) => {
        const header = btoa(JSON.stringify({ alg: 'HS256', typ: 'JWT' }))
        const payload = btoa(JSON.stringify({
            sub: 'user-123',
            role: role,
            exp: Math.floor(Date.now() / 1000) + expInSeconds
        }))
        return `${header}.${payload}.signature`
    }

    it('should derive userRole from JWT token rather than tampered localStorage user object', () => {
        const authStore = useAuthStore()
        const validStudentJWT = createDummyJWT('student')

        // Attacker tampers with localStorage 'user' object to claim admin role
        localStorage.setItem('user', JSON.stringify({ id: 'user-123', role: 'admin' }))
        localStorage.setItem('token', validStudentJWT)

        // Force store initialization
        authStore.token = validStudentJWT
        authStore.user = JSON.parse(localStorage.getItem('user'))

        // Expect userRole to evaluate to 'student' from signed JWT, ignoring tampered 'admin'
        expect(authStore.userRole).toBe('student')
    })

    it('should detect expired JWT tokens and mark isAuthenticated as false', () => {
        const authStore = useAuthStore()
        const expiredJWT = createDummyJWT('teacher', -3600) // Expired 1 hour ago

        authStore.token = expiredJWT
        expect(authStore.isAuthenticated).toBe(false)
    })
})
